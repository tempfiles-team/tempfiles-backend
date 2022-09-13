package file

import (
	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
	"golang.org/x/net/context"
)

func Delete(objectName string) (fiber.Map, error) {
	ctx := context.Background()

	err := MinioClient.RemoveObject(ctx, BucketName, objectName, minio.RemoveObjectOptions{})

	if err != nil {
		return nil, err
	}
	return fiber.Map{
		"message": "delete success",
		"success": true,
	}, nil
}
