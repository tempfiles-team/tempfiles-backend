package file

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
)

type ResultStruct struct {
	Name         string
	LastModified string
	Size         int64
	Expires      string
}

func list() (fiber.Map, error) {

	ctx := context.Background()
	objectCh := MinioClient.ListObjects(ctx, BucketName, minio.ListObjectsOptions{})

	var result []ResultStruct
	for object := range objectCh {
		if object.Err != nil {
			return nil, object.Err
		}
		result = append(result, ResultStruct{Name: object.Key, LastModified: object.LastModified.Format(time.RFC3339),
			Size: object.Size, Expires: object.Expires.Format(time.RFC3339)})
	}

	return fiber.Map{
		"message":      "list success",
		"success":      true,
		"list":         result,
		"numberOfList": len(result),
	}, nil
}

func ListHandler(c *fiber.Ctx) error {
	result, err := list()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "minio list error",
			"error":   err.Error(),
		})
	}
	return c.JSON(result)
}
