package newfile

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/minpeter/tempfiles-backend/database"
)

func DeleteHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	fileName := c.Params("filename")

	if err := os.Remove("tmp/" + id + "/" + fileName); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "file delete error",
			"error":   err.Error(),
		})
	}

	//db에서 삭제
	if _, err := database.Engine.Delete(&database.FileTracking{FileId: id, FileName: fileName}); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "db delete error",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "File deleted successfully",
	})
}
