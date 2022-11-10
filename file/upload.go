package file

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
	"github.com/minpeter/tempfiles-backend/database"
	"github.com/minpeter/tempfiles-backend/jwt"
	"golang.org/x/crypto/bcrypt"
)

func UploadHandler(c *fiber.Ctx) error {

	pw := c.Query("pw", "")

	hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "bcrypt hash error",
			"error":   err.Error(),
		})
	}

	data, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Please upload a file (multipart/form-data)",
			"error":   err.Error(),
		})
	}

	FileTracking := &database.FileTracking{
		FileName: data.Filename,
		// FileType:    data.Header["Content-Type"][0],
		FileSize:    data.Size,
		IsEncrypted: pw != "",
		Password:    string(hash),
	}

	_, err = database.Engine.Insert(FileTracking)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "database insert error",
			"error":   err.Error(),
		})
	}

	token, exp, err := jwt.CreateJWTToken(*FileTracking)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "jwt token creation error",
			"error":   err.Error(),
		})
	}

	// Get Buffer from file
	buffer, err := data.Open()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "can't chage file to buffer",
			"error":   err.Error(),
		})
	}
	defer buffer.Close()

	objectName := data.Filename
	contentType := data.Header["Content-Type"][0]
	fileSize := data.Size

	// Upload the zip file with FPutObject
	info, err := MinioClient.PutObject(context.Background(), BucketName, objectName, buffer, fileSize, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "minio upload error",
			"error":   err.Error(),
		})
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)

	if FileTracking.IsEncrypted {
		return c.JSON(fiber.Map{
			"message":      "upload success",
			"success":      true,
			"filename":     objectName,
			"size":         info.Size,
			"expires":      info.Expiration.Format(time.RFC3339),
			"filetype":     contentType,
			"isEncrypted":  FileTracking.IsEncrypted,
			"lastModified": info.LastModified.Format(time.RFC3339),
			"token":        token,
			"tokenExpires": exp,
			"delete_url":   fmt.Sprintf("%s/del/%s?token=%s", os.Getenv("BACKEND_BASEURL"), info.Key, token),
			"download_url": fmt.Sprintf("%s/dl/%s?token=%s", os.Getenv("BACKEND_BASEURL"), info.Key, token),
		})
	} else {
		return c.JSON(fiber.Map{
			"message":      "upload success",
			"success":      true,
			"filename":     objectName,
			"size":         info.Size,
			"expires":      info.Expiration.Format(time.RFC3339),
			"filetype":     contentType,
			"isEncrypted":  FileTracking.IsEncrypted,
			"lastModified": info.LastModified.Format(time.RFC3339),
			"delete_url":   fmt.Sprintf("%s/del/%s", os.Getenv("BACKEND_BASEURL"), info.Key),
			"download_url": fmt.Sprintf("%s/dl/%s", os.Getenv("BACKEND_BASEURL"), info.Key),
		})
	}
}
