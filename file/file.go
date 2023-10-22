package file

import (
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
		FolderId: id,
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

	if files, err := GetFiles(FileTracking.FolderId, c.BaseURL()); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "folder not found",
			"error":   nil,
		})
	} else {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message":       "file found",
			"isEncrypted":   FileTracking.IsEncrypted,
			"uploadDate":    FileTracking.UploadDate.Format(time.RFC3339),
			"files":         files,
			"folderId":      FileTracking.FolderId,
			"provideToken":  c.Query("token") != "",
			"downloadLimit": FileTracking.DownloadLimit,
			"downloadCount": FileTracking.DownloadCount,
			"expireTime":    FileTracking.ExpireTime.Format(time.RFC3339),
		})
	}
}
