package file

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
)

type ResultStruct struct {
	Filename     string `json:"filename"`
	Size         int64  `json:"size"`
	Expires      string `json:"expires"`
	Filetype     string `json:"filetype"`
	IsEncrypted  bool   `json:"isEncrypted"`
	LastModified string `json:"lastModified"`
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
				Filename:     object.Key,
				Size:         object.Size,
				Expires:      object.Expires.Format(time.RFC3339),
				IsEncrypted:  false,
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
