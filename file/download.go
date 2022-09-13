package file

import (
	"context"
	"log"

	"github.com/minio/minio-go/v7"
)

// MinioConnection func for opening minio connection.
func Download(objectName string) (string, error) {
	ctx := context.Background()

	filePath := "tmp/" + objectName

	err := MinioClient.FGetObject(ctx, BucketName, objectName, filePath, minio.GetObjectOptions{})
	if err != nil {
		return "", err
	}

	log.Printf("Successfully downloaded %s\n", objectName)

	return filePath, nil
}
