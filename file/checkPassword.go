package file

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tempfiles-Team/tempfiles-backend/database"
	"github.com/tempfiles-Team/tempfiles-backend/jwt"
	"github.com/tempfiles-Team/tempfiles-backend/response"
	"golang.org/x/crypto/bcrypt"
)

func CheckPasswordHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	pw := c.Query("pw", "")

	if id == "" || pw == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewFailMessageResponse("Please provide a file id and password"))
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

	if err := bcrypt.CompareHashAndPassword([]byte(FileTracking.Password), []byte(pw)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(response.NewFailMessageResponse("password incorrect"))
	}

	token, _, err := jwt.CreateJWTToken(FileTracking)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewFailMessageResponse("jwt create error"))
	}

	return c.JSON(response.NewSuccessDataResponse(fiber.Map{
		"token": token,
	}))
}
