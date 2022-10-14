package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/joho/godotenv/autoload"
	"github.com/minpeter/tempfiles-backend/file"
)

func main() {

	VER := 1.1
	app := fiber.New()

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
	tempfile := os.Getenv("BACKEND_TEMPPATH")
	data, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "The file field of multipart is required.",
		})
	}

	tempfilePath := fmt.Sprintf("./%s/%s", tempfile, data.Filename)
	if err := c.SaveFile(data, tempfilePath); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Error while saving file.",
		})
	}

	result, err := file.Upload(data.Filename, tempfilePath, data.Header["Content-Type"][0])
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "minio upload error",
		})
	}
	if err := os.Remove(tempfilePath); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Error while remove file.",
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

	return c.SendFile(filePath)
}
