package controllers

import (
	"github.com/gofiber/fiber/v2"
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
