package file

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tempfiles-Team/tempfiles-backend/database"
	"github.com/tempfiles-Team/tempfiles-backend/response"
)

func ListHandler(c *fiber.Ctx) error {

	var files []database.FileTracking
	// IsDeleted가 false인 파일만 가져옴
	if err := database.Engine.Where("is_deleted = ?", false).Find(&files); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewFailMessageResponse("file list error"))
	}

	return c.Status(fiber.StatusOK).JSON(response.NewSuccessDataResponse(files))
}
