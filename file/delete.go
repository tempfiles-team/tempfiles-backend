package file

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/tempfiles-Team/tempfiles-backend/database"
	"github.com/tempfiles-Team/tempfiles-backend/response"
)

func DeleteHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewFailMessageResponse("Please provide a file id"))
	}

	FileTracking := database.FileTracking{
		FileId: id,
	}

	has, err := database.Engine.Get(&FileTracking)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewFailMessageResponse("db query error"))
	}

	if !has {
		return c.Status(fiber.StatusNotFound).JSON(response.NewFailMessageResponse("file not found"))
	}

	if err := os.RemoveAll("tmp/" + FileTracking.FileId); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewFailMessageResponse("file delete error"))
	}

	//db에서 삭제
	if _, err := database.Engine.Delete(&FileTracking); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewFailMessageResponse("db delete error"))
	}

	return c.Status(fiber.StatusOK).JSON(response.NewSuccessMessageResponse("File deleted successfully"))

}
