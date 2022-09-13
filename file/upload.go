package file

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
)

// MinioConnection func for opening minio connection.
func Upload(objectName, filePath, contentType string) (fiber.Map, error) {
	ctx := context.Background()

	// Upload the zip file with FPutObject
	info, err := MinioClient.FPutObject(ctx, BucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return nil, err
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)

	return fiber.Map{
		"message":    "upload success",
		"success":    true,
		"filename":   objectName,
		"size":       info.Size,
		"filetype":   contentType,
		"url":        info.Location,
		"expires":    info.Expiration,
		"delete_url": fmt.Sprintf("%s/delete/%s", os.Getenv("BACKEND_BASEURL"), info.Key),
	}, nil
}
