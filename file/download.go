package file

import (
	"context"
	"log"
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
)

func download(objectName string) (*minio.Object, minio.ObjectInfo, string, error) {
	decodedObjectName, err := url.QueryUnescape(objectName)
	if err != nil {
		log.Printf("Error decoding object name: %s", err)
		return nil, minio.ObjectInfo{}, "", err
	}

	object, err := MinioClient.GetObject(context.Background(), BucketName, decodedObjectName, minio.GetObjectOptions{})
	if err != nil {
		log.Printf("Error getting object: %s", err)
		return nil, minio.ObjectInfo{}, "", err
	}
	stat, nil := object.Stat()
	log.Printf("Successfully downloaded %s\n", decodedObjectName)

	return object, stat, decodedObjectName, nil
}

func DownloadHandler(c *fiber.Ctx) error {
	fileName := c.Params("filename")
	object, objectInfo, fileName, err := download(fileName)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "minio download error",
			"error":   err.Error(),
		})
	}

	c.Response().Header.Set("Content-Type", objectInfo.ContentType)
	c.Response().Header.Set("Content-Disposition", "attachment; filename="+fileName)
	c.Response().Header.Set("Accept-Ranges", "bytes")
	defer object.Close()

	return c.SendStream(object)
}
