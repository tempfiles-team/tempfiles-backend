package file

import (
	"github.com/gofiber/fiber/v2"
	"github.com/minpeter/tempfiles-backend/database"
)

func DownloadHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message":  "Please provide a file id",
			"error":    nil,
			"download": false,
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

	c.Attachment()
	return c.SendFile("tmp/" + FileTracking.FileId + "/" + FileTracking.FileName)
}
