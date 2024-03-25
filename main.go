package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-fuego/fuego"
	"github.com/rs/cors"

	_ "github.com/joho/godotenv/autoload"
	controller "github.com/tempfiles-Team/tempfiles-backend/controllers"
	"github.com/tempfiles-Team/tempfiles-backend/crontab"
	"github.com/tempfiles-Team/tempfiles-backend/database"
	"github.com/tempfiles-Team/tempfiles-backend/utils"
)

func main() {

	crontab.Crontab()

	if err := utils.CheckTmpFolder(); err != nil {
		log.Fatalf("tmp folder error: %v", err)
	}

	if err := database.CreateDBEngine(); err != nil {
		log.Fatalf("failed to create db engine: %v", err)
	}

	// ======================== SERVER ========================

	if os.Getenv("BACKEND_PORT") == "" {
		os.Setenv("BACKEND_PORT", "5000")
	}
	port := os.Getenv("BACKEND_PORT")

	s := fuego.NewServer(
		fuego.WithAddr("0.0.0.0:"+port),
		fuego.WithCorsMiddleware(cors.New(cors.Options{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodDelete},
			AllowedHeaders: []string{"Origin", "Content-Type", "Accept", "X-Download-Limit", "X-Time-Limit", "X-Hidden"},
		}).Handler),
	)

	fuego.Get(s, "/", func(c fuego.ContextNoBody) (string, error) {
		return "TEMPFILES API WORKING 🚀\nIf you want to use the API, go to '/swagger'", nil
	})

	v1wv := fuego.Group(s, "/").Tags("files")
	controller.FilesRessources{}.RoutesV1(v1wv)

	s.Run()
}
