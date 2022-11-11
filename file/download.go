package file

import (
	"net/url"

	"github.com/gofiber/fiber/v2"
)

func DownloadHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	fileName, _ := url.PathUnescape(c.Params("filename"))

	if fileName == "" || id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message":  "Please provide a file id and filename",
			"error":    nil,
			"download": false,
		})
	}

	c.Attachment()
	return c.SendFile("tmp/" + id + "/" + fileName)
}
