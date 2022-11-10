package newfile

import (
	"github.com/gofiber/fiber/v2"
	"github.com/minpeter/tempfiles-backend/database"
	"github.com/minpeter/tempfiles-backend/jwt"
	"golang.org/x/crypto/bcrypt"
)

func CheckPasswordHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	fileName := c.Params("filename")

	pw := c.Query("pw", "")

	if fileName == "" || id == "" || pw == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Please provide a file id, filename and password",
			"error":   nil,
			"unlock":  false,
		})
	}

	FileTracking := database.FileTracking{
		FileName: fileName,
		FileId:   id,
	}

	has, err := database.Engine.Get(&FileTracking)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "db query error",
			"error":   err.Error(),
			"unlock":  false,
		})
	}

	if !has {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "file not found",
			"error":   nil,
			"unlock":  false,
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(FileTracking.Password), []byte(pw)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "password incorrect",
			"error":   err.Error(),
			"unlock":  false,
		})
	}

	token, _, err := jwt.CreateJWTToken(FileTracking)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "jwt token creation error",
			"error":   err.Error(),
			"unlock":  false,
		})
	}

	// jwt 토큰 생성
	return c.JSON(fiber.Map{
		"message": "password correct",
		"token":   token,
		"unlock":  true,
	})
}
