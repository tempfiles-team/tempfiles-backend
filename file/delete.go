package file

import (
	"log"
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
	"golang.org/x/net/context"
)

func DeleteHandler(c *fiber.Ctx) error {

	fileName, err := url.QueryUnescape(c.Params("filename"))

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "minio delete error",
			"error":   err.Error(),
		})
	}

	if err := MinioClient.RemoveObject(context.Background(), BucketName, fileName, minio.RemoveObjectOptions{}); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "minio delete error",
			"error":   err.Error(),
		})
	}
	log.Printf("Successfully deleted %s\n", fileName)

	return c.JSON(fiber.Map{
		"message": "delete success",
		"success": true,
	})
}
