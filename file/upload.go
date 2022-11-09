package file

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
)

func UploadHandler(c *fiber.Ctx) error {

	data, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Please upload a file (multipart/form-data)",
			"error":   err.Error(),
		})
	}

	// Get Buffer from file
	buffer, err := data.Open()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "can't chage file to buffer",
			"error":   err.Error(),
		})
	}
	defer buffer.Close()

	objectName := data.Filename
	contentType := data.Header["Content-Type"][0]
	fileSize := data.Size

	// Upload the zip file with FPutObject
	info, err := MinioClient.PutObject(context.Background(), BucketName, objectName, buffer, fileSize, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "minio upload error",
			"error":   err.Error(),
		})
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)

	return c.JSON(fiber.Map{
		"message":      "upload success",
		"success":      true,
		"filename":     objectName,
		"size":         info.Size,
		"expires":      info.Expiration.Format(time.RFC3339),
		"isEncrypted":  false,
		"lastModified": info.LastModified.Format(time.RFC3339),
		"delete_url":   fmt.Sprintf("%s/del/%s", os.Getenv("BACKEND_BASEURL"), info.Key),
		"download_url": fmt.Sprintf("%s/dl/%s", os.Getenv("BACKEND_BASEURL"), info.Key),
	})

}
