package main

import (
	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/tempfiles-Team/tempfiles-backend/docs"
	"github.com/tempfiles-Team/tempfiles-backend/pkg/configs"
	"github.com/tempfiles-Team/tempfiles-backend/pkg/middleware"
	"github.com/tempfiles-Team/tempfiles-backend/pkg/routes"
	"github.com/tempfiles-Team/tempfiles-backend/pkg/utils"
)

// @title Tempfiles API
// @version 2.0
// @description This is a Tempfiles server for file and text sharing.
// @termsOfService https://tmpf.me/terms/

// @contact.name API Support
// @contact.url https://tmpf.me/credit
// @contact.email dev@tmpf.me

// @BasePath /
func main() {
	config := configs.FiberConfig()

	app := fiber.New(config)

	middleware.FiberMiddleware(app)

	utils.ReadyComponent()

	router := app.Group("/")

	routes.SwaggerRoute(router)
	routes.PublicRoutes(router)
	routes.PrivateRouter(router)
	routes.NotFoundRoute(router)

	utils.StartServer(app)
}
