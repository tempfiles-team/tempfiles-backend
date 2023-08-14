package controllers

import (
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/tempfiles-Team/tempfiles-backend/app/models"
	"github.com/tempfiles-Team/tempfiles-backend/app/queries"
	"github.com/tempfiles-Team/tempfiles-backend/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

// DeleteText godoc
// @Summary Delete text item.
// @Description Delete text item.
// @Tags text
// @Accept */*
// @Produce json
// @Param id path string true "text id"
// @Success 200 {object} utils.Response
// @Router /text/{id} [delete]
func DeleteText(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewFailMessageResponse("Please provide a text id"))
	}

	TextS := new(queries.TextState)
	has, err := TextS.GetText(id)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailMessageResponse("db query error"))
	}

	if !has {
		return c.Status(fiber.StatusNotFound).JSON(utils.NewFailMessageResponse("text not found"))
	}

	if err := TextS.DelText(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailMessageResponse("db delete error"))
	}

	return c.Status(fiber.StatusOK).JSON(utils.NewSuccessMessageResponse("Text deleted successfully"))

}

// ListText godoc
// @Summary List text items.
// @Description List text items.
// @Tags text
// @Accept */*
// @Produce json
// @Success 200 {object} utils.Response{data=models.TextTracking}
// @Router /texts [get]
func ListText(c *fiber.Ctx) error {

	texts, err := queries.GetTexts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailMessageResponse("text list error"))
	}

	return c.Status(fiber.StatusOK).JSON(utils.NewSuccessDataResponse(texts))
}

// DownloadText godoc
// @Summary Download text item.
// @Description Download text item.
// @Tags text
// @Accept */*
// @Produce json
// @Param id path string true "text id"
// @Param token query string false "token"
// @Success 200 {object} utils.Response
// @Router /text/{id} [get]
func DownloadText(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewFailMessageResponse("Please provide a text id"))
	}

	TextS := new(queries.TextState)
	has, err := TextS.GetText(id)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailMessageResponse("db query error"))
	}

	if !has {
		return c.Status(fiber.StatusNotFound).JSON(utils.NewFailMessageResponse("text not found"))
	}
	if err := TextS.IncreaseDLCount(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailMessageResponse("db update error"))
	}

	isExp, err := TextS.IsExpiredText()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailMessageResponse("db update error"))
	}

	if isExp {
		return c.Status(fiber.StatusNotFound).JSON(utils.NewFailMessageResponse("text is expired"))
	}

	return c.Status(fiber.StatusOK).JSON(utils.NewSuccessDataResponse(fiber.Map{
		"id":            TextS.Model.TextId,
		"textData":      TextS.Model.TextData,
		"uploadDate":    TextS.Model.UploadDate.Format(time.RFC3339),
		"downloadLimit": TextS.Model.DownloadLimit,
		"downloadCount": TextS.Model.DownloadCount,
		"expireTime":    TextS.Model.ExpireTime.Format(time.RFC3339),
	}))
}

// UploadText godoc
// @Summary Upload text item.
// @Description Upload text item.
// @Tags text
// @Accept */*
// @Produce json
// @Accept text/plain
// @Param X-Download-Limit header string false "download limit"
// @Param X-Time-Limit header string false "time limit"
// @Success 200 {object} utils.Response
// @Router /text/new [post]
func UploadText(c *fiber.Ctx) error {

	pasteText := string(c.Body())

	if pasteText == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewFailMessageResponse("Please provide a text"))
	}

	password := c.Query("pw", "")

	downloadLimit, err := strconv.Atoi(string(c.Request().Header.Peek("X-Download-Limit")))
	if err != nil {
		downloadLimit = 0
	}
	expireTime, err := strconv.Atoi(string(c.Request().Header.Peek("X-Time-Limit")))
	var expireTimeDate time.Time
	if err != nil || expireTime < 0 || expireTime == 0 {
		// 기본 3시간 후 만료
		expireTimeDate = time.Now().Add(time.Duration(60*3) * time.Minute)
	} else {
		expireTimeDate = time.Now().Add(time.Duration(expireTime) * time.Minute)
	}

	TextS := new(queries.TextState)
	TextS.Model = models.TextTracking{
		TextId:   utils.RandString(),
		TextData: pasteText,
	}

	TextS.Model.UploadDate = time.Now()
	TextS.Model.IsEncrypted = password != ""
	TextS.Model.DownloadLimit = int64(downloadLimit)
	TextS.Model.ExpireTime = expireTimeDate

	var token string = ""
	if TextS.Model.IsEncrypted {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailMessageResponse("bcrypt hash error"))
		}
		TextS.Model.Password = string(hash)
		token, _, err = utils.CreateJWTToken(TextS.Model.TextId)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailMessageResponse("jwt create error"))
		}
	}

	if err := TextS.InsertFile(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewSuccessMessageResponse("database insert error"))
	}
	log.Printf("Successfully uploaded %s, download limit %d\n", TextS.Model.TextId, TextS.Model.DownloadLimit)

	return c.Status(fiber.StatusOK).JSON(utils.NewSuccessDataResponse(fiber.Map{
		"id":            TextS.Model.TextId,
		"textData":      TextS.Model.TextData,
		"isEncrypted":   TextS.Model.IsEncrypted,
		"uploadDate":    TextS.Model.UploadDate.Format(time.RFC3339),
		"token":         token,
		"downloadLimit": TextS.Model.DownloadLimit,
		"downloadCount": TextS.Model.DownloadCount,
		"expireTime":    TextS.Model.ExpireTime.Format(time.RFC3339),
	}))
}
