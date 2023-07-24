package file

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/tempfiles-Team/tempfiles-backend/database"
	"github.com/tempfiles-Team/tempfiles-backend/jwt"
	"github.com/tempfiles-Team/tempfiles-backend/response"
	"golang.org/x/crypto/bcrypt"
)

func UploadHandler(c *fiber.Ctx) error {
	data, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewFailMessageResponse("Please upload a file (multipart/form-data)"))
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

	FileTracking := &database.FileTracking{
		FileName:      data.Filename,
		FileSize:      data.Size,
		UploadDate:    time.Now(),
		FileId:        database.RandString(),
		IsEncrypted:   password != "",
		DownloadLimit: int64(downloadLimit),
		ExpireTime:    expireTimeDate,
	}

	var token string = ""
	if FileTracking.IsEncrypted {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.NewFailMessageResponse("bcrypt hash error"))
		}
		FileTracking.Password = string(hash)
		token, _, err = jwt.CreateJWTToken(*FileTracking)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.NewFailMessageResponse("jwt create error"))
		}
	}

	if CheckFileFolder(FileTracking.FileId) != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewFailMessageResponse("file folder create error"))
	}

	if err := c.SaveFile(data, fmt.Sprintf("tmp/%s/%s", FileTracking.FileId, FileTracking.FileName)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewFailMessageResponse("file save error"))
	}

	_, err = database.Engine.Insert(FileTracking)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewFailMessageResponse("database insert error"))
	}

	log.Printf("Successfully uploaded %s of size %d, download limit %d\n", FileTracking.FileName, FileTracking.FileSize, FileTracking.DownloadLimit)

	return c.Status(fiber.StatusOK).JSON(response.NewSuccessDataResponse(fiber.Map{
		"fileId":        FileTracking.FileId,
		"filename":      FileTracking.FileName,
		"size":          FileTracking.FileSize,
		"isEncrypted":   FileTracking.IsEncrypted,
		"uploadDate":    FileTracking.UploadDate.Format(time.RFC3339),
		"token":         token,
		"downloadLimit": FileTracking.DownloadLimit,
		"downloadCount": FileTracking.DownloadCount,
		"expireTime":    FileTracking.ExpireTime.Format(time.RFC3339),
	}))
}
