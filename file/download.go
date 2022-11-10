package file

import (
	"net/url"

	"github.com/gofiber/fiber/v2"
)

func DownloadHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	fileName := c.Params("filename")

	if fileName == "" || id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message":  "Please provide a file id and filename",
			"error":    nil,
			"download": false,
		})
	}

	return c.Download("tmp/"+id+"/"+fileName, url.PathEscape(fileName))
}
