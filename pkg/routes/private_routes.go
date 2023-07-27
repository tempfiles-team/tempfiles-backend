package routes

import (
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
			//id or filename이 없으면 jwt 검사 안함
			if len(strings.Split(c.OriginalURL(), "/")) != 3 {
				// 핸들러가 알아서 에러를 반환함
				return false
			}

			id := strings.Split(c.OriginalURL(), "/")[2]
			if strings.Contains(id, "?") {
				id = strings.Split(id, "?")[0]
			}

			utils.FileId = id

			return utils.IsEncrypted(id)
		},
		KeyFunc: utils.IsMatched(),
	}))

	// DownloadFile godoc
	// @Summary Download file item.
	// @Description Download file item.
	// @Tags file
	// @Accept */*
	// @Produce json
	// @Param id path string true "file id"
	// @Param token query string false "token"
	// @Success 200 {object} utils.Response
	// @Router /dl/{id} [get]
	r.Get("/dl/:id", controllers.DownloadFile)
	// DeleteFile godoc
	// @Summary Delete file item.
	// @Description Delete file item.
	// @Tags file
	// @Accept */*
	// @Produce json
	// @Param id path string true "file id"
	// @Success 200 {object} utils.Response
	// @Router /file/{id} [delete]
	r.Delete("/file/:id", controllers.DeleteFile)

}
