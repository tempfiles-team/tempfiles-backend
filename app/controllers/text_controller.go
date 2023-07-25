package controllers

import (
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/tempfiles-Team/tempfiles-backend/database"
	"github.com/tempfiles-Team/tempfiles-backend/pkg/utils"
)

func DeleteText(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewFailMessageResponse("Please provide a text id"))
	}

	TextTracking := database.TextTracking{
		TextId: id,
	}

	has, err := database.Engine.Get(&TextTracking)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailMessageResponse("db query error"))
	}

	if !has {
		return c.Status(fiber.StatusNotFound).JSON(utils.NewFailMessageResponse("text not found"))
	}

	if _, err := database.Engine.Delete(&TextTracking); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailMessageResponse("db delete error"))
	}

	return c.Status(fiber.StatusOK).JSON(utils.NewSuccessMessageResponse("Text deleted successfully"))

}

func ListText(c *fiber.Ctx) error {

	var texts []database.TextTracking
	// IsDeleted가 false인 파일만 가져옴
	if err := database.Engine.Where("is_deleted = ?", false).Find(&texts); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailMessageResponse("text list error"))
	}

	return c.Status(fiber.StatusOK).JSON(utils.NewSuccessDataResponse(texts))
}

func DownloadText(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewFailMessageResponse("Please provide a text id"))
	}

	TextTracking := database.TextTracking{
		TextId: id,
	}

	has, err := database.Engine.Get(&TextTracking)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailMessageResponse("db query error"))
	}

	if !has {
		return c.Status(fiber.StatusNotFound).JSON(utils.NewFailMessageResponse("text not found"))
	}

	// db DownloadCount +1
	TextTracking.DownloadCount++
	if _, err := database.Engine.ID(TextTracking.Id).Update(&TextTracking); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailMessageResponse("db update error"))
	}

	// Download Limit check
	if TextTracking.DownloadLimit != 0 && TextTracking.DownloadCount >= TextTracking.DownloadLimit {
		// Download Limit exceeded -> check IsDelete
		TextTracking.IsDeleted = true

		log.Printf("check IsDeleted file: %s\n", TextTracking.TextId)
		if _, err := database.Engine.ID(TextTracking.Id).Cols("Is_deleted").Update(&TextTracking); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailMessageResponse("db update error"))
		}
	}

	return c.Status(fiber.StatusOK).JSON(utils.NewSuccessDataResponse(fiber.Map{
		"textId":        TextTracking.TextId,
		"textData":      TextTracking.TextData,
		"uploadDate":    TextTracking.UploadDate.Format(time.RFC3339),
		"downloadLimit": TextTracking.DownloadLimit,
		"downloadCount": TextTracking.DownloadCount,
		"expireTime":    TextTracking.ExpireTime.Format(time.RFC3339),
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
	TextTracking := &database.TextTracking{
		TextId:        database.RandString(),
		TextData:      pasteText,
		UploadDate:    time.Now(),
		DownloadLimit: int64(downloadLimit),
		ExpireTime:    expireTimeDate,
	}

	_, err = database.Engine.Insert(TextTracking)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewSuccessDataResponse(fiber.Map{
			"message": "database insert error",
			"error":   err.Error(),
		}))
	}

	log.Printf("Successfully uploaded %s, download limit %d\n", TextTracking.TextId, TextTracking.DownloadLimit)

	return c.Status(fiber.StatusOK).JSON(utils.NewSuccessDataResponse(fiber.Map{
		"textId":        TextTracking.TextId,
		"uploadDate":    TextTracking.UploadDate.Format(time.RFC3339),
		"downloadLimit": TextTracking.DownloadLimit,
		"downloadCount": TextTracking.DownloadCount,
		"expireTime":    TextTracking.ExpireTime.Format(time.RFC3339),
	}))
}
