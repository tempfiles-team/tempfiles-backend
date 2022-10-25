package file

import (
	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
	"golang.org/x/net/context"
)

func delete(objectName string) (fiber.Map, error) {
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

func DeleteHandler(c *fiber.Ctx) error {
	fileName := c.Params("filename")
	result, err := delete(fileName)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "minio delete error",
			"error":   err.Error(),
		})
	}
	return c.JSON(result)
}
