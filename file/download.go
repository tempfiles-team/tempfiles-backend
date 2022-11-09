package file

import (
	"context"
	"log"
	"net/url"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
)

func DownloadHandler(c *fiber.Ctx) error {

	fileName, err := url.QueryUnescape((c.Params("filename")))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "minio download error",
			"error":   err.Error(),
		})
	}

	filePath := "tmp/" + fileName

	if err := MinioClient.FGetObject(context.Background(), BucketName, fileName, filePath, minio.GetObjectOptions{}); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "minio download error",
			"error":   err.Error(),
		})
	}

	// filePath와 fileName 로깅
	log.Printf("downloaded file %s saved to %s", fileName, filePath)

	defer os.Remove(filePath)

	return c.Download(filePath, fileName)
}
