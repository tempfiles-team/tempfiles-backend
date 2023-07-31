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

// DeleteFile godoc
// @Summary Delete file item.
// @Description Delete file item.
// @Tags file
// @Accept */*
// @Produce json
// @Param id path string true "file id"
// @Success 200 {object} utils.Response
// @Router /file/{id} [delete]
func DeleteFile(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewFailMessageResponse("Please provide a file id"))
	}

	FileS := new(queries.FileState)
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

// DownloadFile godoc
// @Summary Download file item.
// @Description Download file item.
// @Tags file
// @Accept */*
// @Produce json
// @Param id path string true "file id"
// @Param token query string false "token"
// @Success 200 {object} utils.Response
// @Router /dl/{id} [get]
func DownloadFile(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewFailMessageResponse("Please provide a file id"))
	}

	FileS := new(queries.FileState)
	has, err := FileS.GetFile(id)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailMessageResponse("db query error"))
	}

	if !has {
		return c.Status(fiber.StatusNotFound).JSON(utils.NewFailMessageResponse("file not found"))
	}

	isExp, err := FileS.IsExpiredFile()
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.NewFailDataResponse(nil))
	}

	if isExp {
		return c.Status(fiber.StatusNotFound).JSON(utils.NewFailMessageResponse("file is expired"))
	}

	if err := FileS.IncreaseDLCount(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailMessageResponse("db update error"))
	}

	c.Response().Header.Set("Content-Disposition", "attachment; filename="+strings.ReplaceAll(url.PathEscape(FileS.Model.FileName), "+", "%20"))
	return c.SendFile("tmp/" + FileS.Model.FileId + "/" + FileS.Model.FileName)
}

// GetFile godoc
// @Summary Get file item.
// @Description Get file item.
// @Tags file
// @Accept */*
// @Produce json
// @Param id path string true "file id"
// @Param token query string false "token"
// @Success 200 {object} utils.Response
// @Router /file/{id} [get]
func GetFile(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewFailMessageResponse("Please provide a file id"))
	}

	FileS := new(queries.FileState)
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

// ListFile godoc
// @Summary List file items.
// @Description List file items.
// @Tags file
// @Accept */*
// @Produce json
// @Success 200 {object} utils.Response{data=models.FileTracking}
// @Router /files [get]
func ListFile(c *fiber.Ctx) error {

	files, err := queries.GetFiles()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailDataResponse(nil))
	}

	return c.Status(fiber.StatusOK).JSON(utils.NewSuccessDataResponse(files))
}

// UploadFile godoc
// @Summary Upload file item.
// @Description Upload file item.
// @Tags file
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "file"
// @Param pw query string false "password"
// @Param X-Download-Limit header string false "download limit"
// @Param X-Time-Limit header string false "time limit"
// @Success 200 {object} utils.Response
// @Router /upload [post]
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

	FileS := new(queries.FileState)

	FileS.Model = models.FileTracking{
		FileName: data.Filename,
		FileSize: data.Size,
		FileId:   utils.RandString(),
	}

	FileS.Model.UploadDate = time.Now()
	FileS.Model.IsEncrypted = password != ""
	FileS.Model.DownloadLimit = int64(downloadLimit)
	FileS.Model.ExpireTime = expireTimeDate

	var token string = ""
	if FileS.Model.IsEncrypted {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailMessageResponse("bcrypt hash error"))
		}
		FileS.Model.Password = string(hash)
		token, _, err = utils.CreateJWTToken(FileS.Model.FileId)
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

	if err := FileS.InsertFile(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailMessageResponse("database insert error"))
	}
	log.Printf("Successfully uploaded %s of size %d, download limit %d\n", FileS.Model.FileName, FileS.Model.FileSize, FileS.Model.DownloadLimit)

	return c.Status(fiber.StatusOK).JSON(utils.NewSuccessDataResponse(fiber.Map{
		"id":            FileS.Model.FileId,
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
