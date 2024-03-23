package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

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

	port, _ := strconv.Atoi(os.Getenv("BACKEND_PORT"))
	s := fuego.NewServer(
		fuego.WithPort(port),
		fuego.WithCorsMiddleware(cors.New(cors.Options{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodDelete},
			AllowedHeaders: []string{"Origin", "Content-Type", "Accept", "X-Download-Limit", "X-Time-Limit", "X-Hidden"},
		}).Handler),
	)

	fuego.Get(s, "/", func(c fuego.ContextNoBody) (string, error) {
		return "TEMPFILES API WORKING 🚀\nIf you want to use the API, go to '/swagger'", nil
	})

	v2 := fuego.Group(s, "/v2").Tags("files v2")
	controller.FilesRessources{}.RoutesV2(v2)

	v1 := fuego.Group(s, "/v1").Tags("files v1")
	v1wv := fuego.Group(s, "/").Tags("files v1 (without version)")
	controller.FilesRessources{}.RoutesV1(v1wv)
	controller.FilesRessources{}.RoutesV1(v1)

	s.Run()
}
