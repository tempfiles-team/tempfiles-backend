package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tempfiles-Team/tempfiles-backend/app/controllers"
	"github.com/tempfiles-Team/tempfiles-backend/pkg/middleware"
)

// PublicRoutes func for describe group of public routes.
func PublicRoutes(r fiber.Router) {
	// Create routes group.

	r.Get("/", controllers.HealthCheck)
	r.Get("/list", controllers.ListAll)

	r.Get("/files", controllers.ListFile)
	r.Post("/upload", controllers.UploadFile)

	r.Get("/texts", controllers.ListText)
	r.Post("/text/new", controllers.UploadText)

	r.Get("/file/:id", middleware.IdVaildator(), controllers.GetFile)

	r.Get("/detail/:id", middleware.IdVaildator(), controllers.GetDetail)

	r.Get("/checkpw/:id", middleware.IdVaildator(), controllers.CheckPasswordFile)
}
