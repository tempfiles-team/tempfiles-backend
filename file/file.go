package file

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/tempfiles-Team/tempfiles-backend/database"
	"github.com/tempfiles-Team/tempfiles-backend/response"
)

func FileHandler(c *fiber.Ctx) error {
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

	backendUrl := c.BaseURL()

	return c.Status(fiber.StatusOK).JSON(response.NewSuccessDataResponse(fiber.Map{
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
