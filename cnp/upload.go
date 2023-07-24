package cnp

import (
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/tempfiles-Team/tempfiles-backend/database"
	"github.com/tempfiles-Team/tempfiles-backend/response"
)

func UploadHandler(c *fiber.Ctx) error {

	pasteText := string(c.Body())

	if pasteText == "" {
		return c.Status(fiber.StatusOK).JSON(response.NewFailMessageResponse("Please provide a text"))
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
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewSuccessDataResponse(fiber.Map{
			"message": "database insert error",
			"error":   err.Error(),
		}))
	}

	log.Printf("Successfully uploaded %s, download limit %d\n", TextTracking.TextId, TextTracking.DownloadLimit)

	return c.Status(fiber.StatusOK).JSON(response.NewSuccessDataResponse(fiber.Map{
		"textId":        TextTracking.TextId,
		"uploadDate":    TextTracking.UploadDate.Format(time.RFC3339),
		"downloadLimit": TextTracking.DownloadLimit,
		"downloadCount": TextTracking.DownloadCount,
		"expireTime":    TextTracking.ExpireTime.Format(time.RFC3339),
	}))
}
