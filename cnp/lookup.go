package cnp

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/tempfiles-Team/tempfiles-backend/database"
	"github.com/tempfiles-Team/tempfiles-backend/response"
)

func ListHandler(c *fiber.Ctx) error {

	var texts []database.TextTracking
	// IsDeleted가 false인 파일만 가져옴
	if err := database.Engine.Where("is_deleted = ?", false).Find(&texts); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewFailMessageResponse("text list error"))
	}

	return c.Status(fiber.StatusOK).JSON(response.NewSuccessDataResponse(texts))
}

func DownloadHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewFailMessageResponse("Please provide a text id"))
	}

	TextTracking := database.TextTracking{
		TextId: id,
	}

	has, err := database.Engine.Get(&TextTracking)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewFailMessageResponse("db query error"))
	}

	if !has {
		return c.Status(fiber.StatusNotFound).JSON(response.NewFailMessageResponse("text not found"))
	}

	// db DownloadCount +1
	TextTracking.DownloadCount++
	if _, err := database.Engine.ID(TextTracking.Id).Update(&TextTracking); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewFailMessageResponse("db update error"))
	}

	// Download Limit check
	if TextTracking.DownloadLimit != 0 && TextTracking.DownloadCount >= TextTracking.DownloadLimit {
		// Download Limit exceeded -> check IsDelete
		TextTracking.IsDeleted = true

		log.Printf("check IsDeleted file: %s\n", TextTracking.TextId)
		if _, err := database.Engine.ID(TextTracking.Id).Cols("Is_deleted").Update(&TextTracking); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.NewFailMessageResponse("db update error"))
		}
	}

	return c.Status(fiber.StatusOK).JSON(response.NewSuccessDataResponse(fiber.Map{
		"textId":        TextTracking.TextId,
		"textData":      TextTracking.TextData,
		"uploadDate":    TextTracking.UploadDate.Format(time.RFC3339),
		"downloadLimit": TextTracking.DownloadLimit,
		"downloadCount": TextTracking.DownloadCount,
		"expireTime":    TextTracking.ExpireTime.Format(time.RFC3339),
	}))
}
