package file

import (
	"github.com/gofiber/fiber/v2"
	"github.com/minpeter/tempfiles-backend/database"
)

func ListHandler(c *fiber.Ctx) error {

	var files []database.FileTracking
	err := database.Engine.Find(&files)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "file list error",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":      "File list successfully",
		"list":         files,
		"numberOfList": len(files),
	})
}
