package file

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/tempfiles-Team/tempfiles-backend/database"
	"github.com/tempfiles-Team/tempfiles-backend/jwt"
	"golang.org/x/crypto/bcrypt"
)

func UploadHandler(c *fiber.Ctx) error {

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Please upload a file (multipart/form-data)",
			"error":   err.Error(),
		})
	}

	password := c.Query("pw", "")
	downloadLimit, err := strconv.Atoi(string(c.Request().Header.Peek("X-Download-Limit")))
	if err != nil {
		downloadLimit = 100
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
		FileCount:     len(form.File["file"]),
		FolderId:      database.RandString(),
		UploadDate:    time.Now(),
		IsEncrypted:   password != "",
		DownloadLimit: int64(downloadLimit),
		ExpireTime:    expireTimeDate,
	}

	var token string = ""
	if FileTracking.IsEncrypted {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "bcrypt hash error",
				"error":   err.Error(),
			})
		}
		FileTracking.Password = string(hash)
		token, _, err = jwt.CreateJWTToken(*FileTracking)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "jwt token creation error",
				"error":   err.Error(),
			})
		}
	}

	if CheckFileFolder(FileTracking.FolderId) != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "file folder creation error",
			"error":   err.Error(),
		})
	}

	for _, file := range form.File["file"] {
		if err := c.SaveFile(file, fmt.Sprintf("tmp/%s/%s", FileTracking.FolderId, file.Filename)); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "file save error",
				"error":   err.Error(),
			})
		}
	}

	_, err = database.Engine.Insert(FileTracking)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "database insert error",
			"error":   err.Error(),
		})
	}

	log.Printf("Successfully uploaded %s, download limit %d\n", FileTracking.FolderId, FileTracking.DownloadLimit)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":       "File uploaded successfully",
		"folderId":      FileTracking.FolderId,
		"isEncrypted":   FileTracking.IsEncrypted,
		"uploadDate":    FileTracking.UploadDate.Format(time.RFC3339),
		"token":         token,
		"downloadLimit": FileTracking.DownloadLimit,
		"downloadCount": FileTracking.DownloadCount,
		"expireTime":    FileTracking.ExpireTime.Format(time.RFC3339),
	})
}
