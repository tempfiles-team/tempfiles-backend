package file

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/tempfiles-Team/tempfiles-backend/database"
	"github.com/tempfiles-Team/tempfiles-backend/jwt"
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

	var token string = ""
	if FileTracking.IsEncrypted {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "bcrypt hash error",
				"error":   err.Error(),
			})
		}
		FileTracking.Password = string(hash)
		token, _, err = jwt.CreateJWTToken(*FileTracking)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "jwt token creation error",
				"error":   err.Error(),
			})
		}
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
		"message":     "File uploaded successfully",
		"fileId":      FileTracking.FileId,
		"filename":    FileTracking.FileName,
		"size":        FileTracking.FileSize,
		"isEncrypted": FileTracking.IsEncrypted,
		"uploadDate":  FileTracking.UploadDate.Format(time.RFC3339),
		"token":       token,
	})
}
