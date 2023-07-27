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

	// HealthCheck godoc
	// @Summary Show the status of server.
	// @Description get the status of server.
	// @Tags root
	// @Accept */*
	// @Produce json
	// @Success 200 {object} utils.Response
	// @Router / [get]
	r.Get("/", controllers.HealthCheck)
	// GetInfo godoc
	// @Summary Show the information of server.
	// @Description get the information of server.
	// @Tags root
	// @Accept */*
	// @Produce json
	// @Param api query string false "api name"
	// @Success 200 {object} utils.Response
	// @Router /info [get]
	r.Get("/info", controllers.GetInfo)
	// ListFile godoc
	// @Summary List file items.
	// @Description List file items.
	// @Tags file
	// @Accept */*
	// @Produce json
	// @Success 200 {object} utils.Response{data=models.FileTracking}
	// @Router /files [get]
	r.Get("/files", controllers.ListFile)
	// UploadFile godoc
	// @Summary Upload file item.
	// @Description Upload file item.
	// @Tags file
	// @Accept multipart/form-data
	// @Produce json
	// @Param file formData file true "file"
	// @Param pw query string false "password"
	// @Param X-Download-Limit header string false "download limit"
	// @Param X-Time-Limit header string false "time limit"
	// @Success 200 {object} utils.Response
	// @Router /upload [post]
	r.Post("/upload", controllers.UploadFile)

	// UploadText godoc
	// @Summary Upload text item.
	// @Description Upload text item.
	// @Tags text
	// @Accept */*
	// @Produce json
	// @Param X-Download-Limit header string false "download limit"
	// @Param X-Time-Limit header string false "time limit"
	// @Success 200 {object} utils.Response
	// @Router /text/new [post]
	r.Post("/text/new", controllers.UploadText)
	// ListText godoc
	// @Summary List text items.
	// @Description List text items.
	// @Tags text
	// @Accept */*
	// @Produce json
	// @Success 200 {object} utils.Response{data=models.TextTracking}
	// @Router /texts [get]
	r.Get("/texts", controllers.ListText)
	// DownloadText godoc
	// @Summary Download text item.
	// @Description Download text item.
	// @Tags text
	// @Accept */*
	// @Produce json
	// @Param id path string true "text id"
	// @Param token query string false "token"
	// @Success 200 {object} utils.Response
	// @Router /text/{id} [get]
	r.Get("/text/:id", controllers.DownloadText)
	// DeleteText godoc
	// @Summary Delete text item.
	// @Description Delete text item.
	// @Tags text
	// @Accept */*
	// @Produce json
	// @Param id path string true "text id"
	// @Success 200 {object} utils.Response
	// @Router /text/{id} [delete]
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
	// GetFile godoc
	// @Summary Get file item.
	// @Description Get file item.
	// @Tags file
	// @Accept */*
	// @Produce json
	// @Param id path string true "file id"
	// @Param token query string false "token"
	// @Success 200 {object} utils.Response
	// @Router /file/{id} [get]
	r.Get("/file/:id", controllers.GetFile)
	// CheckPasswordFile godoc
	// @Summary Check password of file item.
	// @Description Check password of file item.
	// @Tags file
	// @Accept */*
	// @Produce json
	// @Param id path string true "file id"
	// @Param pw query string true "password"
	// @Success 200 {object} utils.Response
	// @Router /checkpw/{id} [get]
	r.Get("/checkpw/:id", controllers.CheckPasswordFile)
}
