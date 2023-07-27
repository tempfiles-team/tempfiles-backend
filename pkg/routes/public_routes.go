package routes

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/tempfiles-Team/tempfiles-backend/app/controllers"
	"github.com/tempfiles-Team/tempfiles-backend/app/queries"
	"github.com/tempfiles-Team/tempfiles-backend/pkg/utils"
)

// PublicRoutes func for describe group of public routes.
func PublicRoutes(r fiber.Router) {
	// Create routes group.

	r.Get("/", controllers.HealthCheck)

	r.Get("/files", controllers.ListFile)
	r.Post("/upload", controllers.UploadFile)

	r.Post("/text/new", controllers.UploadText)
	r.Get("/texts", controllers.ListText)
	r.Get("/text/:id", controllers.DownloadText)
	r.Delete("/text/:id", controllers.DeleteText)

	r.Use(func(c *fiber.Ctx) error {
		if len(strings.Split(c.OriginalURL(), "/")) != 3 {
			return c.Status(fiber.StatusBadRequest).JSON(utils.NewFailMessageResponse("sorry, endpoint is not found"))
		}

		id := strings.Split(c.OriginalURL(), "/")[2]
		if strings.Contains(id, "?") {
			id = strings.Split(id, "?")[0]
		}

		log.Printf("id: %v", id)

		FileS := queries.FileState{}
		has, err := FileS.GetFile(id)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailDataResponse(nil))
		}
		if !has {
			return c.Status(fiber.StatusNotFound).JSON(utils.NewFailMessageResponse("file not exist"))
		}
		if FileS.Model.IsDeleted {
			return c.Status(fiber.StatusNotFound).JSON(utils.NewFailMessageResponse("file is deleted"))
		}
		return c.Next()
	})

	r.Get("/file/:id", controllers.GetFile)
	r.Get("/checkpw/:id", controllers.CheckPasswordFile)
}
