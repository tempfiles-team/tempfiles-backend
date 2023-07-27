package controllers

import (
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/tempfiles-Team/tempfiles-backend/app/models"
	"github.com/tempfiles-Team/tempfiles-backend/app/queries"
	"github.com/tempfiles-Team/tempfiles-backend/pkg/utils"
)

func DeleteText(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewFailMessageResponse("Please provide a text id"))
	}

	TextS := queries.TextState{}
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

func ListText(c *fiber.Ctx) error {

	TextS := queries.TextState{}
	texts, err := TextS.GetTexts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailMessageResponse("text list error"))
	}

	return c.Status(fiber.StatusOK).JSON(utils.NewSuccessDataResponse(texts))
}

func DownloadText(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewFailMessageResponse("Please provide a text id"))
	}

	TextS := queries.TextState{}
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
		"textId":        TextS.Model.TextId,
		"textData":      TextS.Model.TextData,
		"uploadDate":    TextS.Model.UploadDate.Format(time.RFC3339),
		"downloadLimit": TextS.Model.DownloadLimit,
		"downloadCount": TextS.Model.DownloadCount,
		"expireTime":    TextS.Model.ExpireTime.Format(time.RFC3339),
	}))
}

func UploadText(c *fiber.Ctx) error {

	pasteText := string(c.Body())

	if pasteText == "" {
		return c.Status(fiber.StatusOK).JSON(utils.NewFailMessageResponse("Please provide a text"))
	}

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

	TextS := queries.TextState{}
	TextS.Model = models.TextTracking{
		TextId:        utils.RandString(),
		TextData:      pasteText,
		UploadDate:    time.Now(),
		DownloadLimit: int64(downloadLimit),
		ExpireTime:    expireTimeDate,
	}

	if err := TextS.InsertFile(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewSuccessMessageResponse("database insert error"))
	}
	log.Printf("Successfully uploaded %s, download limit %d\n", TextS.Model.TextId, TextS.Model.DownloadLimit)

	return c.Status(fiber.StatusOK).JSON(utils.NewSuccessDataResponse(fiber.Map{
		"textId":        TextS.Model.TextId,
		"uploadDate":    TextS.Model.UploadDate.Format(time.RFC3339),
		"downloadLimit": TextS.Model.DownloadLimit,
		"downloadCount": TextS.Model.DownloadCount,
		"expireTime":    TextS.Model.ExpireTime.Format(time.RFC3339),
	}))
}
