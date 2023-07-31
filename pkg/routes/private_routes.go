package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tempfiles-Team/tempfiles-backend/app/controllers"
	"github.com/tempfiles-Team/tempfiles-backend/pkg/middleware"
)

func PrivateRouter(r fiber.Router) {

	r.Use(middleware.IdVaildator(), middleware.JWTProtected())

	r.Get("/dl/:id", controllers.DownloadFile)
	r.Delete("/file/:id", controllers.DeleteFile)

	r.Get("/text/:id", controllers.DownloadText)
	r.Delete("/text/:id", controllers.DeleteText)

}
