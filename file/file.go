package file

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/tempfiles-Team/tempfiles-backend/database"
)

func FileHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Please provide a file id",
			"error":   nil,
		})
	}

	FileTracking := database.FileTracking{
		FileId: id,
	}

	has, err := database.Engine.Get(&FileTracking)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "db query error",
			"error":   err.Error(),
		})
	}

	if !has {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "file not found",
			"error":   nil,
		})
	}

	backendUrl := c.BaseURL()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":       "file found",
		"filename":      FileTracking.FileName,
		"size":          FileTracking.FileSize,
		"isEncrypted":   FileTracking.IsEncrypted,
		"uploadDate":    FileTracking.UploadDate.Format(time.RFC3339),
		"delete_url":    fmt.Sprintf("%s/del/%s", backendUrl, FileTracking.FileId),
		"download_url":  fmt.Sprintf("%s/dl/%s", backendUrl, FileTracking.FileId),
		"provide_token": c.Query("token") != "",
		"downloadLimit": FileTracking.DownloadLimit,
		"downloadCount": FileTracking.DownloadCount,
		"expireTime":    FileTracking.ExpireTime.Format(time.RFC3339),
	})
}
