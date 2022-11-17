package file

import (
	"net/url"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/tempfiles-Team/tempfiles-backend/database"
)

func DownloadHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message":  "Please provide a file id",
			"error":    nil,
			"download": false,
		})
	}

	FileTracking := database.FileTracking{
		FileId: id,
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
			"message": "file not found",
			"error":   nil,
		})
	}

	// Download Limit check
	if FileTracking.DownloadLimit > 0 {
		if FileTracking.DownloadLimit <= FileTracking.DownloadCount {
			// Download Limit exceeded -> file delete
			if err := os.RemoveAll("tmp/" + FileTracking.FileId); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "file delete error",
					"error":   err.Error(),
					"delete":  false,
				})
			}

			//db에서 삭제
			if _, err := database.Engine.Delete(&FileTracking); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "db delete error",
					"error":   err.Error(),
					"delete":  false,
				})
			}

			return c.Status(fiber.StatusGone).JSON(fiber.Map{
				"message": "download limit exceeded",
				"error":   nil,
				"delete":  true,
			})
		}
	}

	// db DownloadCount +1
	FileTracking.DownloadCount++
	if _, err := database.Engine.ID(FileTracking.Id).Update(&FileTracking); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "db update error",
			"error":   err.Error(),
		})
	}

	c.Response().Header.Set("Content-Disposition", "attachment; filename="+strings.ReplaceAll(url.PathEscape(FileTracking.FileName), "+", "%20"))
	return c.SendFile("tmp/" + FileTracking.FileId + "/" + FileTracking.FileName)
}
