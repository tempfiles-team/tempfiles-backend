package controllers

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/tempfiles-Team/tempfiles-backend/app/models"
	"github.com/tempfiles-Team/tempfiles-backend/app/queries"
	"github.com/tempfiles-Team/tempfiles-backend/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

func CheckPasswordFile(c *fiber.Ctx) error {
	id := c.Params("id")

	pw := c.Query("pw", "")

	if id == "" || pw == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewFailMessageResponse("Please provide a file id and password"))
	}

	FileS := queries.FileState{}
	has, err := FileS.GetFile(id)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailMessageResponse("db query error"))
	}

	if !has {
		return c.Status(fiber.StatusNotFound).JSON(utils.NewFailMessageResponse("file not found"))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(FileS.Model.Password), []byte(pw)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.NewFailMessageResponse("password incorrect"))
	}

	token, _, err := utils.CreateJWTToken(FileS.Model)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailMessageResponse("jwt create error"))
	}

	return c.JSON(utils.NewSuccessDataResponse(fiber.Map{
		"token": token,
	}))
}

func DeleteFile(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewFailMessageResponse("Please provide a file id"))
	}

	FileS := queries.FileState{}
	has, err := FileS.GetFile(id)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailMessageResponse("db query error"))
	}

	if !has {
		return c.Status(fiber.StatusNotFound).JSON(utils.NewFailMessageResponse("file not found"))
	}

	if err := os.RemoveAll("tmp/" + FileS.Model.FileId); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailMessageResponse("file delete error"))
	}

	if err := FileS.DelFile(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailMessageResponse("db delete error"))
	}

	return c.Status(fiber.StatusOK).JSON(utils.NewSuccessMessageResponse("File deleted successfully"))
}

func DownloadFile(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewFailMessageResponse("Please provide a file id"))
	}

	FileS := queries.FileState{}
	has, err := FileS.GetFile(id)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailMessageResponse("db query error"))
	}

	if !has {
		return c.Status(fiber.StatusNotFound).JSON(utils.NewFailMessageResponse("file not found"))
	}

	if err := FileS.IncreaseDLCount(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailMessageResponse("db update error"))
	}

	isExp, err := FileS.IsExpiredFile()
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.NewFailDataResponse(nil))
	}

	if isExp {
		return c.Status(fiber.StatusNotFound).JSON(utils.NewFailMessageResponse("file is expired"))
	}

	c.Response().Header.Set("Content-Disposition", "attachment; filename="+strings.ReplaceAll(url.PathEscape(FileS.Model.FileName), "+", "%20"))
	return c.SendFile("tmp/" + FileS.Model.FileId + "/" + FileS.Model.FileName)
}

func GetFile(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewFailMessageResponse("Please provide a file id"))
	}

	FileS := queries.FileState{}
	has, err := FileS.GetFile(id)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailMessageResponse("db query error"))
	}

	if !has {
		return c.Status(fiber.StatusNotFound).JSON(utils.NewFailMessageResponse("file not found"))
	}

	backendUrl := c.BaseURL()

	return c.Status(fiber.StatusOK).JSON(utils.NewSuccessDataResponse(fiber.Map{
		"delete_url":    fmt.Sprintf("%s/del/%s", backendUrl, FileS.Model.FileId),
		"download_url":  fmt.Sprintf("%s/dl/%s", backendUrl, FileS.Model.FileId),
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

func ListFile(c *fiber.Ctx) error {

	FileS := queries.FileState{}
	files, err := FileS.GetFiles()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailDataResponse(nil))
	}

	return c.Status(fiber.StatusOK).JSON(utils.NewSuccessDataResponse(files))
}

func UploadFile(c *fiber.Ctx) error {
	data, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewFailMessageResponse("Please upload a file (multipart/form-data)"))
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

	FileS := queries.FileState{}

	FileS.Model = models.FileTracking{
		FileName:      data.Filename,
		FileSize:      data.Size,
		UploadDate:    time.Now(),
		FileId:        utils.RandString(),
		IsEncrypted:   password != "",
		DownloadLimit: int64(downloadLimit),
		ExpireTime:    expireTimeDate,
	}

	var token string = ""
	if FileS.Model.IsEncrypted {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailMessageResponse("bcrypt hash error"))
		}
		FileS.Model.Password = string(hash)
		token, _, err = utils.CreateJWTToken(FileS.Model)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailMessageResponse("jwt create error"))
		}
	}

	if utils.CheckFileFolder(FileS.Model.FileId) != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailMessageResponse("file folder create error"))
	}

	if err := c.SaveFile(data, fmt.Sprintf("tmp/%s/%s", FileS.Model.FileId, FileS.Model.FileName)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailMessageResponse("file save error"))
	}

	// _, err = queries.Engine.Insert(FileS.Model)
	// if err != nil {
	// 	return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailMessageResponse("database insert error"))
	// }

	if err := FileS.InsertFile(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailMessageResponse("database insert error"))
	}
	log.Printf("Successfully uploaded %s of size %d, download limit %d\n", FileS.Model.FileName, FileS.Model.FileSize, FileS.Model.DownloadLimit)

	return c.Status(fiber.StatusOK).JSON(utils.NewSuccessDataResponse(fiber.Map{
		"fileId":        FileS.Model.FileId,
		"filename":      FileS.Model.FileName,
		"size":          FileS.Model.FileSize,
		"isEncrypted":   FileS.Model.IsEncrypted,
		"uploadDate":    FileS.Model.UploadDate.Format(time.RFC3339),
		"token":         token,
		"downloadLimit": FileS.Model.DownloadLimit,
		"downloadCount": FileS.Model.DownloadCount,
		"expireTime":    FileS.Model.ExpireTime.Format(time.RFC3339),
	}))
}
