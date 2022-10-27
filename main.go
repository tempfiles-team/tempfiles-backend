package main

import (
	"fmt"
	"log"
	"math"

	"github.com/minpeter/tempfiles-backend/data"
	"github.com/minpeter/tempfiles-backend/file"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/joho/godotenv/autoload"
)

type SignupRequest struct {
	Name     string
	Email    string
	Password string
}

type LoginRequest struct {
	Email    string
	Password string
}

func main() {

	VER := "1.1.6"
	app := fiber.New(fiber.Config{
		AppName:   "tempfiles-backend",
		BodyLimit: int(math.Pow(1024, 3)), // 1 == 1byte
	})

	app.Use(cache.New(cache.Config{StoreResponseHeaders: true}), cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, DELETE",
	}))

	engine, err := data.CreateDBEngine()
	if err != nil {
		log.Fatal(err)
	}

	file.MinioClient, err = file.Connection()
	if err != nil {
		log.Fatalf("minio connection error: %v", err)
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message":    "api is working normally :)",
			"apiVersion": VER,
		})
	})

	app.Post("/upload", upload)
	app.Get("/list", list)
	app.Delete("/del/:filename", delete)
	app.Get("/dl/:filename", download)

	log.Fatal(app.Listen(fmt.Sprintf(":%s", os.Getenv("BACKEND_PORT"))))
}

func upload(c *fiber.Ctx) error {
	data, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Get Buffer from file
	buffer, err := data.Open()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "can't chage file to buffer",
			"error":  err.Error(),
		})
	}
	defer buffer.Close()

	objectName := data.Filename
	fileBuffer := buffer
	contentType := data.Header["Content-Type"][0]
	fileSize := data.Size

	result, err := file.Upload(objectName, contentType, fileBuffer, fileSize)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status": "minio upload error",
			"error":  err.Error(),
		})
	}
	return c.JSON(result)
}

func list(c *fiber.Ctx) error {
	result, err := file.List()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "minio list error",
		})
	}
	return c.JSON(result)
}

func delete(c *fiber.Ctx) error {
	fileName := c.Params("filename")
	result, err := file.Delete(fileName)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "minio delete error",
		})
	}
	return c.JSON(result)
}

func download(c *fiber.Ctx) error {
	fileName := c.Params("filename")
	object, objectInfo, fileName, err := file.Download(fileName)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "minio download error",
			"err":     err.Error(),
		})
	}

	c.Response().Header.Set("Content-Type", objectInfo.ContentType)
	c.Response().Header.Set("Content-Disposition", "attachment; filename="+fileName)
	c.Response().Header.Set("Accept-Ranges", "bytes")
	defer object.Close()

	return c.SendStream(object)
}
