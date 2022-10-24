package file

import (
	"context"
	"log"
	"net/url"

	"github.com/minio/minio-go/v7"
)

// MinioConnection func for opening minio connection.
func Download(objectName string) (string, string, error) {
	ctx := context.Background()

	decodedObjectName, err := url.QueryUnescape(objectName)
	if err != nil {
		log.Fatal(err)
		return "", "", err
	}

	filePath := "tmp/" + decodedObjectName
	log.Print("Downloading file: ", filePath)

	if err := MinioClient.FGetObject(ctx, BucketName, decodedObjectName, filePath, minio.GetObjectOptions{}); err != nil {
		log.Println(err)
		return "", "", err
	}

	log.Printf("Successfully downloaded %s\n", decodedObjectName)

	return filePath, decodedObjectName, nil
}
