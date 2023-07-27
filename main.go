package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload"
	"github.com/robfig/cron"
	"github.com/tempfiles-Team/tempfiles-backend/app/queries"
	_ "github.com/tempfiles-Team/tempfiles-backend/docs"
	"github.com/tempfiles-Team/tempfiles-backend/pkg/configs"
	"github.com/tempfiles-Team/tempfiles-backend/pkg/middleware"
	"github.com/tempfiles-Team/tempfiles-backend/pkg/routes"
	"github.com/tempfiles-Team/tempfiles-backend/pkg/utils"
	"github.com/tempfiles-Team/tempfiles-backend/platform/db"
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

	terminator := cron.New()
	terminator.AddFunc("0 * * * *", func() {
		if err := queries.IsExpiredFiles(); err != nil {
			log.Println("cron db query error", err.Error())
		}
	})

	terminator.AddFunc("0 0 * * *", func() {
		if err := queries.DelExpireFiles(); err != nil {
			log.Println("cron db query error", err.Error())
		}
	})
	terminator.Start()

	var err error

	if utils.CheckTmpFolder() != nil {
		log.Fatalf("tmp folder error: %v", err)
	}

	if db.OpenDBConnection() != nil {
		log.Fatalf("db connection error: %v", err)
	}

	root := app.Group("/")

	routes.SwaggerRoute(app)
	routes.PublicRoutes(root)
	routes.PrivateRouter(root)
	routes.NotFoundRoute(app)

	utils.StartServer(app)
	terminator.Stop()
}
