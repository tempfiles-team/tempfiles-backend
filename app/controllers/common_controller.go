package controllers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/tempfiles-Team/tempfiles-backend/app/queries"
	"github.com/tempfiles-Team/tempfiles-backend/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

// HealthCheck godoc
// @Summary Show the status of server.
// @Description get the status of server.
// @Tags common
// @Accept */*
// @Produce json
// @Success 200 {object} utils.Response
// @Router / [get]
func HealthCheck(c *fiber.Ctx) error {
	res := utils.NewSuccessMessageResponse("Server is up and running")

	if err := c.JSON(res); err != nil {
		return err
	}

	return nil
}

// ListAll godoc
// @Summary List all items.
// @Description List all items.
// @Tags common
// @Accept */*
// @Produce json
// @Success 200 {object} utils.Response
// @Router /list [get]
func ListAll(c *fiber.Ctx) error {
	f, err := queries.GetFiles()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailMessageResponse("db query error"))
	}
	t, err := queries.GetTexts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailMessageResponse("db query error"))
	}

	return c.JSON(utils.NewSuccessDataResponse(fiber.Map{
		"files": f,
		"texts": t,
	}))
}

// CheckPasswordFile godoc
// @Summary Check password of file item.
// @Description Check password of file item.
// @Tags common
// @Accept */*
// @Produce json
// @Param id path string true "file id"
// @Param pw query string true "password"
// @Success 200 {object} utils.Response
// @Router /checkpw/{id} [get]
func CheckPasswordFile(c *fiber.Ctx) error {
	id := c.Params("id")

	pw := c.Query("pw", "")

	if id == "" || pw == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewFailMessageResponse("Please provide a file id and password"))
	}

	FileS := new(queries.FileState)
	has, err := FileS.GetFile(id)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailMessageResponse("db query error"))
	}

	// 존재하고 비번이 맞으면 토큰 발급
	if err := bcrypt.CompareHashAndPassword([]byte(FileS.Model.Password), []byte(pw)); has && err != nil {
		token, _, err := utils.CreateJWTToken(id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailMessageResponse("jwt create error"))
		}

		return c.JSON(utils.NewSuccessDataResponse(fiber.Map{
			"token": token,
		}))
	}

	TextS := new(queries.TextState)
	has, err = TextS.GetText(id)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailMessageResponse("db query error"))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(TextS.Model.Password), []byte(pw)); has && err != nil {

		token, _, err := utils.CreateJWTToken(id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailMessageResponse("jwt create error"))
		}

		return c.JSON(utils.NewSuccessDataResponse(fiber.Map{
			"token": token,
		}))
	}

	return c.Status(fiber.StatusUnauthorized).JSON(utils.NewFailMessageResponse("password is incorrect"))

}

func GetDetail(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewFailMessageResponse("Please provide a file id"))
	}

	FileS := new(queries.FileState)
	has, err := FileS.GetFile(id)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailMessageResponse("db query error"))
	}
	backendUrl := c.BaseURL()

	if has {
		return c.JSON(utils.NewSuccessDataResponse(fiber.Map{
			"delete_url":    fmt.Sprintf("%s/del/%s", backendUrl, FileS.Model.FileId),
			"download_url":  fmt.Sprintf("%s/dl/%s", backendUrl, FileS.Model.FileId),
			"type":          "file",
			"filename":      FileS.Model.FileName,
			"size":          FileS.Model.FileSize,
			"uploadDate":    FileS.Model.UploadDate.Format(time.RFC3339),
			"isEncrypted":   FileS.Model.IsEncrypted,
			"provide_token": c.Query("token") != "",
			"downloadLimit": FileS.Model.DownloadLimit,
			"downloadCount": FileS.Model.DownloadCount,
			"expireTime":    FileS.Model.ExpireTime.Format(time.RFC3339),
		}))
	}

	TextS := new(queries.TextState)
	has, err = TextS.GetText(id)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailMessageResponse("db query error"))
	}

	if has {
		return c.JSON(utils.NewSuccessDataResponse(fiber.Map{
			"delete_url":    fmt.Sprintf("%s/text/%s", backendUrl, TextS.Model.TextId),
			"download_url":  fmt.Sprintf("%s/text/%s", backendUrl, TextS.Model.TextId),
			"type":          "text",
			"size":          TextS.Model.TextCount,
			"uploadDate":    TextS.Model.UploadDate.Format(time.RFC3339),
			"isEncrypted":   TextS.Model.IsEncrypted,
			"provide_token": c.Query("token") != "",
			"downloadLimit": TextS.Model.DownloadLimit,
			"downloadCount": TextS.Model.DownloadCount,
			"expireTime":    TextS.Model.ExpireTime.Format(time.RFC3339),
		}))
	}

	return c.Status(fiber.StatusNotFound).JSON(utils.NewFailMessageResponse("not found"))
}
