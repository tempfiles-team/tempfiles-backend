package cnp

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tempfiles-Team/tempfiles-backend/database"
)

func DeleteHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Please provide a text id",
			"error":   nil,
			"delete":  false,
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

	if _, err := database.Engine.Delete(&TextTracking); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "db delete error",
			"error":   err.Error(),
			"delete":  false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "text deleted successfully",
		"error":   nil,
		"delete":  true,
	})

}
