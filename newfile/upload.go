package newfile

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/minpeter/tempfiles-backend/database"
	"golang.org/x/crypto/bcrypt"
)

func UploadHandler(c *fiber.Ctx) error {
	data, err := c.FormFile("file")
	password := c.Query("pw", "")

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Please upload a file (multipart/form-data)",
			"error":   err.Error(),
		})
	}

	FileTracking := &database.FileTracking{
		FileName:    data.Filename,
		FileSize:    data.Size,
		UploadDate:  time.Now(),
		FileId:      database.RandString(),
		IsEncrypted: password != "",
	}

	if FileTracking.IsEncrypted {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "bcrypt hash error",
				"error":   err.Error(),
			})
		}
		FileTracking.Password = string(hash)
	}

	if CheckFileFolder(FileTracking.FileId) != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "file folder creation error",
			"error":   err.Error(),
		})
	}

	if err := c.SaveFile(data, fmt.Sprintf("tmp/%s/%s", FileTracking.FileId, FileTracking.FileName)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "file save error",
			"error":   err.Error(),
		})
	}

	_, err = database.Engine.Insert(FileTracking)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "database insert error",
			"error":   err.Error(),
		})
	}

	log.Printf("Successfully uploaded %s of size %d\n", FileTracking.FileName, FileTracking.FileSize)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":      "File uploaded successfully",
		"filename":     FileTracking.FileName,
		"size":         FileTracking.FileSize,
		"isEncrypted":  FileTracking.IsEncrypted,
		"uploadDate":   FileTracking.UploadDate.Format(time.RFC3339),
		"delete_url":   fmt.Sprintf("%s/del/%s/%s", os.Getenv("BACKEND_BASEURL"), FileTracking.FileId, FileTracking.FileName),
		"download_url": fmt.Sprintf("%s/dl/%s/%s", os.Getenv("BACKEND_BASEURL"), FileTracking.FileId, FileTracking.FileName),
	})
}
