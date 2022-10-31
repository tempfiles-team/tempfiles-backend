package file

import (
	"github.com/gofiber/fiber/v2"
	"github.com/minpeter/tempfiles-backend/database"
	"github.com/minpeter/tempfiles-backend/jwt"
	"golang.org/x/crypto/bcrypt"
)

func CheckPasswordHandler(c *fiber.Ctx) error {
	fileName := c.Params("filename")
	pw := c.Query("pw", "")

	if fileName == "" || pw == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "파일 이름 또는 비밀번호가 비어있습니다.",
			"error":   nil,
			"unlock":  false,
		})
	}

	fileRow := new(database.FileRow)
	has, err := database.Engine.Where("file_name = ?", fileName).Desc("id").Get(fileRow)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "데이터베이스 에러",
			"error":   err.Error(),
			"unlock":  false,
		})
	}
	if !has {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "파일이 존재하지 않습니다.",
			"error":   nil,
			"unlock":  false,
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(fileRow.Password), []byte(pw)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "비밀번호가 일치하지 않습니다.",
			"error":   err.Error(),
			"unlock":  false,
		})
	}

	token, exp, err := jwt.CreateJWTToken(*fileRow)
	if err != err {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "JWT 토큰 생성 에러",
			"error":   err.Error(),
			"unlock":  false,
		})
	}

	// jwt 토큰 생성
	return c.JSON(fiber.Map{
		"message": "파일 비밀번호가 일치합니다.",
		"token":   token, "tokenExpires": exp,
		"unlock": true,
	})
}
