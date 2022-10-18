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

func Upload(objectName, contentType string, fileBuffer io.Reader, fileSize int64) (fiber.Map, error) {
	ctx := context.Background()

	// Upload the zip file with FPutObject
	info, err := MinioClient.PutObject(ctx, BucketName, objectName, fileBuffer, fileSize, minio.PutObjectOptions{ContentType: contentType})
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
