package file

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
	"github.com/minpeter/tempfiles-backend/jwt"
)

type ResultStruct struct {
	Name         string
	LastModified string
	Size         int64
	Expires      string
	IsEncrypted  bool
}

func ListHandler(c *fiber.Ctx) error {

	objectCh := MinioClient.ListObjects(context.Background(), BucketName, minio.ListObjectsOptions{})

	var result []ResultStruct

	for object := range objectCh {
		if object.Err != nil {
			return c.Status(500).JSON(fiber.Map{
				"message": "minio list error",
				"error":   object.Err.Error(),
			})
		}

		result = append(
			result,
			ResultStruct{
				Name:         object.Key,
				Size:         object.Size,
				IsEncrypted:  jwt.IsEncrypted(object.Key),
				Expires:      object.Expires.Format(time.RFC3339),
				LastModified: object.LastModified.Format(time.RFC3339),
			})
	}

	return c.JSON(fiber.Map{
		"message":      "list success",
		"success":      true,
		"list":         result,
		"numberOfList": len(result),
	})
}
