package file

import (
	"log"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/tempfiles-Team/tempfiles-backend/database"
)

func DownloadHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	name, err := url.PathUnescape(c.Params("name"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid file name",
			"error":   err.Error(),
		})
	}

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message":  "Please provide a file id",
			"error":    nil,
			"download": false,
		})
	}

	FileTracking := database.FileTracking{
		FolderId: id,
	}

	has, err := database.Engine.Get(&FileTracking)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "db query error",
			"error":   err.Error(),
		})
	}

	if !has {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "folder not found",
			"error":   nil,
		})
	}

	if !CheckIsFileExist(FileTracking.FolderId, name) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "file not found!!",
			"error":   nil,
		})
	}

	// db DownloadCount +1
	FileTracking.DownloadCount++
	if _, err := database.Engine.ID(FileTracking.Id).Update(&FileTracking); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "db update error",
			"error":   err.Error(),
		})
	}

	// Download Limit check
	if FileTracking.DownloadLimit != 0 && FileTracking.DownloadCount >= FileTracking.DownloadLimit {
		// Download Limit exceeded -> check IsDelete
		FileTracking.IsDeleted = true

		log.Printf("check IsDeleted file: %s \n", FileTracking.FolderId)
		if _, err := database.Engine.ID(FileTracking.Id).Cols("Is_deleted").Update(&FileTracking); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "db update error",
				"error":   err.Error(),
			})
		}
	}

	c.Response().Header.Set("Content-Disposition", "attachment; filename="+strings.ReplaceAll(url.PathEscape(name), "+", "%20"))
	return c.SendFile("tmp/" + FileTracking.FolderId + "/" + name)
}
