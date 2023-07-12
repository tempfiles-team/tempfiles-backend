package cnp

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/tempfiles-Team/tempfiles-backend/database"
)

func ListHandler(c *fiber.Ctx) error {

	var texts []database.TextTracking
	// IsDeleted가 false인 파일만 가져옴
	if err := database.Engine.Where("is_deleted = ?", false).Find(&texts); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "file list error",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":      "Text list successfully",
		"list":         texts,
		"numberOfList": len(texts),
	})
}

func DownloadHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message":  "Please provide a text id",
			"error":    nil,
			"download": false,
		})
	}

	TextTracking := database.TextTracking{
		TextId: id,
	}

	has, err := database.Engine.Get(&TextTracking)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "db query error",
			"error":   err.Error(),
		})
	}

	if !has {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "text not found",
			"error":   nil,
		})
	}

	// db DownloadCount +1
	TextTracking.DownloadCount++
	if _, err := database.Engine.ID(TextTracking.Id).Update(&TextTracking); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "db update error",
			"error":   err.Error(),
		})
	}

	// Download Limit check
	if TextTracking.DownloadLimit != 0 && TextTracking.DownloadCount >= TextTracking.DownloadLimit {
		// Download Limit exceeded -> check IsDelete
		TextTracking.IsDeleted = true

		log.Printf("check IsDeleted file: %s\n", TextTracking.TextId)
		if _, err := database.Engine.ID(TextTracking.Id).Cols("Is_deleted").Update(&TextTracking); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "db update error",
				"error":   err.Error(),
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":       "Paste text uploaded successfully",
		"textId":        TextTracking.TextId,
		"textData":      TextTracking.TextData,
		"uploadDate":    TextTracking.UploadDate.Format(time.RFC3339),
		"downloadLimit": TextTracking.DownloadLimit,
		"downloadCount": TextTracking.DownloadCount,
		"expireTime":    TextTracking.ExpireTime.Format(time.RFC3339),
	})
}
