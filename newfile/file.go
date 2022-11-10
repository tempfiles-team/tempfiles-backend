package newfile

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/minpeter/tempfiles-backend/database"
)

func FileHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	fileName := c.Params("filename")

	FileTracking := database.FileTracking{
		FileName: fileName,
		FileId:   id,
	}

	// var user = User{ID: 27}
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
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":      "file found",
		"filename":     FileTracking.FileName,
		"size":         FileTracking.FileSize,
		"isEncrypted":  FileTracking.IsEncrypted,
		"uploadDate":   FileTracking.UploadDate.Format("2000-00-00"),
		"delete_url":   fmt.Sprintf("%s/del/%s/%s", os.Getenv("BACKEND_BASEURL"), FileTracking.FileId, FileTracking.FileName),
		"download_url": fmt.Sprintf("%s/dl/%s/%s", os.Getenv("BACKEND_BASEURL"), FileTracking.FileId, FileTracking.FileName),
	})

}
