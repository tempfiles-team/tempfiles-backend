package routes

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/tempfiles-Team/tempfiles-backend/app/controllers"
	"github.com/tempfiles-Team/tempfiles-backend/pkg/utils"
)

func PrivateRouter(r fiber.Router) {

	r.Use(jwtware.New(jwtware.Config{
		TokenLookup: "query:token",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(utils.NewFailMessageResponse("invalid token, file is password protected / Unauthorized"))
		},

		Filter: func(c *fiber.Ctx) bool {
			if len(strings.Split(c.OriginalURL(), "/")) != 3 {
				return false
			}

			id := strings.Split(c.OriginalURL(), "/")[2]
			if strings.Contains(id, "?") {
				id = strings.Split(id, "?")[0]
			}

			fmt.Println(id)

			utils.FileId = id

			return utils.IsEncrypted(id)
		},
		KeyFunc: utils.IsMatched(),
	}))

	r.Get("/dl/:id", controllers.DownloadFile)

	r.Delete("/file/:id", controllers.DeleteFile)

}
