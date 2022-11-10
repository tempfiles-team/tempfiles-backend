package newfile

import (
	"net/url"

	"github.com/gofiber/fiber/v2"
)

func DownloadHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	fileName := c.Params("filename")

	return c.Download("tmp/"+id+"/"+fileName, url.PathEscape(fileName))
}
