package file

import (
	"context"
	"log"
	"net/url"

	"github.com/minio/minio-go/v7"
)

func Download(objectName string) (*minio.Object, minio.ObjectInfo, string, error) {
	decodedObjectName, err := url.QueryUnescape(objectName)
	if err != nil {
		log.Printf("Error decoding object name: %s", err)
		return nil, minio.ObjectInfo{}, "", err
	}

	object, err := MinioClient.GetObject(context.Background(), BucketName, decodedObjectName, minio.GetObjectOptions{})
	if err != nil {
		log.Println(err)
		return nil, minio.ObjectInfo{}, "", err
	}
	stat, nil := object.Stat()
	log.Printf("Successfully downloaded %s\n", decodedObjectName)

	return object, stat, decodedObjectName, nil
}
