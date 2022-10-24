package main

import (
	"log"
	"math"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/joho/godotenv/autoload"
	"github.com/minpeter/tempfiles-backend/file"
)

func main() {

	VER := "1.1.3"
	app := fiber.New(fiber.Config{
		AppName:   "tempfiles-backend",
		BodyLimit: int(math.Pow(1024, 3)), // 1 == 1byte
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	var err error
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
	app.Delete("/delete/:filename", delete)
	app.Get("/dl/:filename", download)

	log.Fatal(app.Listen(":5000"))
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

	objectName := strings.Replace(data.Filename, " ", "-", -1) // replace spaces with -
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
	filePath, err := file.Download(fileName)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "minio download error",
		})
	}

	defer os.Remove(filePath)

	return c.Download(filePath, fileName)
}
