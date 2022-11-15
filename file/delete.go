package file

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/minpeter/tempfiles-backend/database"
)

func DeleteHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Please provide a file id",
			"error":   nil,
			"delete":  false,
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

	if err := os.Remove("tmp/" + FileTracking.FileId + "/" + FileTracking.FileName); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "file delete error",
			"error":   err.Error(),
			"delete":  false,
		})
	}

	//db에서 삭제
	if _, err := database.Engine.Delete(&FileTracking); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "db delete error",
			"error":   err.Error(),
			"delete":  false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "File deleted successfully",
		"error":   nil,
		"delete":  true,
	})

}
