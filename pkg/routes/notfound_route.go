package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tempfiles-Team/tempfiles-backend/pkg/utils"
)

// NotFoundRoute func for describe 404 Error route.
func NotFoundRoute(r fiber.Router) {
	// Register new special route.
	r.Use(
		// Anonymous function.
		func(c *fiber.Ctx) error {
			// Return HTTP 404 status and JSON response.
			return c.Status(fiber.StatusNotFound).JSON(utils.NewFailMessageResponse("sorry, endpoint is not found"))
		},
	)
}
