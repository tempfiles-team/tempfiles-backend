package file

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
)

func upload(objectName, contentType string, fileBuffer io.Reader, fileSize int64) (fiber.Map, error) {
	ctx := context.Background()

	// Upload the zip file with FPutObject
	info, err := MinioClient.PutObject(ctx, BucketName, objectName, fileBuffer, fileSize, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return nil, err
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)

	return fiber.Map{
		"message":      "upload success",
		"success":      true,
		"filename":     objectName,
		"size":         info.Size,
		"filetype":     contentType,
		"expires":      info.Expiration,
		"delete_url":   fmt.Sprintf("%s/delete/%s", os.Getenv("BACKEND_BASEURL"), info.Key),
		"download_url": fmt.Sprintf("%s/dl/%s", os.Getenv("BACKEND_BASEURL"), info.Key),
	}, nil
}

func UploadHandler(c *fiber.Ctx) error {
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

	result, err := upload(objectName, contentType, fileBuffer, fileSize)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status": "minio upload error",
			"error":  err.Error(),
		})
	}
	return c.JSON(result)
}
