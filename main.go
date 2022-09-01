package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload"
	"github.com/minpeter/tempfiles-backend/file"
)

func main() {

	app := fiber.New()

	var err error
	file.MinioClient, err = file.Connection()
	if err != nil {
		log.Fatalf("minio connection error: %v", err)
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "api is working normally :)",
		})
	})

	app.Post("/upload", upload)

	log.Fatal(app.Listen(":3000"))
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
