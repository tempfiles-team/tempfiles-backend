package cnp

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tempfiles-Team/tempfiles-backend/database"
	"github.com/tempfiles-Team/tempfiles-backend/response"
)

func DeleteHandler(c *fiber.Ctx) error {
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

	if _, err := database.Engine.Delete(&TextTracking); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewFailMessageResponse("db delete error"))
	}

	return c.Status(fiber.StatusOK).JSON(response.NewSuccessMessageResponse("Text deleted successfully"))

}
