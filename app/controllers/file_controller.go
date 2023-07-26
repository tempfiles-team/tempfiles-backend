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
	"github.com/tempfiles-Team/tempfiles-backend/database"
	"github.com/tempfiles-Team/tempfiles-backend/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

// CheckPasswordFile godoc
// @Summary Check password of file item.
// @Description Check password of file item.
// @Tags file
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

	FileTracking := database.FileTracking{
		FileId: id,
	}

	has, err := database.Engine.Get(&FileTracking)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailMessageResponse("db query error"))
	}

	if !has {
		return c.Status(fiber.StatusNotFound).JSON(utils.NewFailMessageResponse("file not found"))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(FileTracking.Password), []byte(pw)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.NewFailMessageResponse("password incorrect"))
	}

	token, _, err := utils.CreateJWTToken(FileTracking)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailMessageResponse("jwt create error"))
	}

	return c.JSON(utils.NewSuccessDataResponse(fiber.Map{
		"token": token,
	}))
}

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

	FileTracking := database.FileTracking{
		FileId: id,
	}

	has, err := database.Engine.Get(&FileTracking)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailMessageResponse("db query error"))
	}

	if !has {
		return c.Status(fiber.StatusNotFound).JSON(utils.NewFailMessageResponse("file not found"))
	}

	if err := os.RemoveAll("tmp/" + FileTracking.FileId); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailMessageResponse("file delete error"))
	}

	//db에서 삭제
	if _, err := database.Engine.Delete(&FileTracking); err != nil {
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

	FileTracking := database.FileTracking{
		FileId: id,
	}

	has, err := database.Engine.Get(&FileTracking)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailMessageResponse("db query error"))
	}

	if !has {
		return c.Status(fiber.StatusNotFound).JSON(utils.NewFailMessageResponse("file not found"))
	}

	// db DownloadCount +1
	FileTracking.DownloadCount++
	if _, err := database.Engine.ID(FileTracking.Id).Update(&FileTracking); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailMessageResponse("db update error"))
	}

	// Download Limit check
	if FileTracking.DownloadLimit != 0 && FileTracking.DownloadCount >= FileTracking.DownloadLimit {
		// Download Limit exceeded -> check IsDelete
		FileTracking.IsDeleted = true

		log.Printf("check IsDeleted file: %s/%s \n", FileTracking.FileId, FileTracking.FileName)
		if _, err := database.Engine.ID(FileTracking.Id).Cols("Is_deleted").Update(&FileTracking); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailMessageResponse("db update error"))
		}
	}

	c.Response().Header.Set("Content-Disposition", "attachment; filename="+strings.ReplaceAll(url.PathEscape(FileTracking.FileName), "+", "%20"))
	return c.SendFile("tmp/" + FileTracking.FileId + "/" + FileTracking.FileName)
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

	FileTracking := database.FileTracking{
		FileId: id,
	}

	has, err := database.Engine.Get(&FileTracking)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailMessageResponse("db query error"))
	}

	if !has {
		return c.Status(fiber.StatusNotFound).JSON(utils.NewFailMessageResponse("file not found"))
	}

	backendUrl := c.BaseURL()

	return c.Status(fiber.StatusOK).JSON(utils.NewSuccessDataResponse(fiber.Map{
		"filename":      FileTracking.FileName,
		"size":          FileTracking.FileSize,
		"isEncrypted":   FileTracking.IsEncrypted,
		"uploadDate":    FileTracking.UploadDate.Format(time.RFC3339),
		"delete_url":    fmt.Sprintf("%s/del/%s", backendUrl, FileTracking.FileId),
		"download_url":  fmt.Sprintf("%s/dl/%s", backendUrl, FileTracking.FileId),
		"provide_token": c.Query("token") != "",
		"downloadLimit": FileTracking.DownloadLimit,
		"downloadCount": FileTracking.DownloadCount,
		"expireTime":    FileTracking.ExpireTime.Format(time.RFC3339),
	}))
}

// ListFile godoc
// @Summary List file items.
// @Description List file items.
// @Tags file
// @Accept */*
// @Produce json
// @Success 200 {object} utils.Response{data=database.FileTracking}
// @Router /files [get]
func ListFile(c *fiber.Ctx) error {

	var files []database.FileTracking
	// IsDeleted가 false인 파일만 가져옴
	if err := database.Engine.Where("is_deleted = ?", false).Find(&files); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailMessageResponse("file list error"))
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
// @Router /file/new [post]
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
			return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailMessageResponse("bcrypt hash error"))
		}
		FileTracking.Password = string(hash)
		token, _, err = utils.CreateJWTToken(*FileTracking)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailMessageResponse("jwt create error"))
		}
	}

	if utils.CheckFileFolder(FileTracking.FileId) != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailMessageResponse("file folder create error"))
	}

	if err := c.SaveFile(data, fmt.Sprintf("tmp/%s/%s", FileTracking.FileId, FileTracking.FileName)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailMessageResponse("file save error"))
	}

	_, err = database.Engine.Insert(FileTracking)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailMessageResponse("database insert error"))
	}

	log.Printf("Successfully uploaded %s of size %d, download limit %d\n", FileTracking.FileName, FileTracking.FileSize, FileTracking.DownloadLimit)

	return c.Status(fiber.StatusOK).JSON(utils.NewSuccessDataResponse(fiber.Map{
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
