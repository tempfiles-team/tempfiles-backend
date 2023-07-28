package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tempfiles-Team/tempfiles-backend/app/queries"
	"github.com/tempfiles-Team/tempfiles-backend/pkg/utils"
)

// HealthCheck godoc
// @Summary Show the status of server.
// @Description get the status of server.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} utils.Response
// @Router / [get]
func HealthCheck(c *fiber.Ctx) error {
	res := utils.NewSuccessMessageResponse("Server is up and running")

	if err := c.JSON(res); err != nil {
		return err
	}

	return nil
}

// ListAll godoc
// @Summary List all items.
// @Description List all items.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} utils.Response
// @Router /list [get]
func ListAll(c *fiber.Ctx) error {
	f, err := queries.GetFiles()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailMessageResponse("db query error"))
	}
	t, err := queries.GetTexts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewFailMessageResponse("db query error"))
	}

	return c.JSON(utils.NewSuccessDataResponse(fiber.Map{
		"files": f,
		"texts": t,
	}))
}
