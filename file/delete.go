package file

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/minpeter/tempfiles-backend/database"
)

func DeleteHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	fileName := c.Params("filename")

	if fileName == "" || id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Please provide a file id and filename",
			"error":   nil,
			"delete":  false,
		})
	}

	if err := os.Remove("tmp/" + id + "/" + fileName); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "file delete error",
			"error":   err.Error(),
			"delete":  false,
		})
	}

	//db에서 삭제
	if _, err := database.Engine.Delete(&database.FileTracking{FileId: id, FileName: fileName}); err != nil {
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
