package file

import (
	"log"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/tempfiles-Team/tempfiles-backend/database"
	"github.com/tempfiles-Team/tempfiles-backend/response"
)

func DownloadHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewFailMessageResponse("Please provide a file id"))
	}

	FileTracking := database.FileTracking{
		FileId: id,
	}

	has, err := database.Engine.Get(&FileTracking)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewFailMessageResponse("db query error"))
	}

	if !has {
		return c.Status(fiber.StatusNotFound).JSON(response.NewFailMessageResponse("file not found"))
	}

	// db DownloadCount +1
	FileTracking.DownloadCount++
	if _, err := database.Engine.ID(FileTracking.Id).Update(&FileTracking); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewFailMessageResponse("db update error"))
	}

	// Download Limit check
	if FileTracking.DownloadLimit != 0 && FileTracking.DownloadCount >= FileTracking.DownloadLimit {
		// Download Limit exceeded -> check IsDelete
		FileTracking.IsDeleted = true

		log.Printf("check IsDeleted file: %s/%s \n", FileTracking.FileId, FileTracking.FileName)
		if _, err := database.Engine.ID(FileTracking.Id).Cols("Is_deleted").Update(&FileTracking); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.NewFailMessageResponse("db update error"))
		}
	}

	c.Response().Header.Set("Content-Disposition", "attachment; filename="+strings.ReplaceAll(url.PathEscape(FileTracking.FileName), "+", "%20"))
	return c.SendFile("tmp/" + FileTracking.FileId + "/" + FileTracking.FileName)
}
