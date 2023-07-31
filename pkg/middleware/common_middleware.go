package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/tempfiles-Team/tempfiles-backend/app/queries"
	"github.com/tempfiles-Team/tempfiles-backend/pkg/utils"
)

// FiberMiddleware provide Fiber's built-in middlewares.
// See: https://docs.gofiber.io/api/middleware
func FiberMiddleware(a *fiber.App) {
	a.Use(

		cors.New(cors.Config{
			AllowOrigins: "*",
			AllowHeaders: "Origin, Content-Type, Accept, X-Download-Limit, X-Time-Limit",
			AllowMethods: "GET, POST, DELETE",
		}),

		recover.New(),

		logger.New(),
	)
}

func IdVaildator() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		if len(strings.Split(c.OriginalURL(), "/")) != 3 {
			return c.Status(fiber.StatusBadRequest).JSON(utils.NewFailMessageResponse("sorry, endpoint is not found"))
		}

		id := strings.Split(c.OriginalURL(), "/")[2]
		if strings.Contains(id, "?") {
			id = strings.Split(id, "?")[0]
		}

		FileS := new(queries.FileState)
		has, err := FileS.GetFile(id)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailDataResponse(nil))
		}
		if has && !FileS.Model.IsDeleted {
			return c.Next()
		}

		TextS := new(queries.TextState)
		has, err = TextS.GetText(id)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailDataResponse(nil))
		}

		if has && !TextS.Model.IsDeleted {
			return c.Next()
		}

		return c.Status(fiber.StatusNotFound).JSON(utils.NewFailMessageResponse("The item has been deleted or does not exist."))

	}
}
