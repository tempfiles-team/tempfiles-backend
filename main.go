package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/robfig/cron"
	"github.com/tempfiles-Team/tempfiles-backend/database"
	"github.com/tempfiles-Team/tempfiles-backend/pkg/configs"
	"github.com/tempfiles-Team/tempfiles-backend/pkg/middleware"
	"github.com/tempfiles-Team/tempfiles-backend/pkg/routes"
	"github.com/tempfiles-Team/tempfiles-backend/pkg/utils"

	_ "github.com/joho/godotenv/autoload"
)

func main() {

	config := configs.FiberConfig()

	app := fiber.New(config)

	middleware.FiberMiddleware(app)

	terminator := cron.New()
	terminator.AddFunc("* */1 * * *", func() {
		var files []database.FileTracking
		//현재 시간보다 expire_time이 작고 is_deleted가 false인 파일을 가져옴
		if err := database.Engine.Where("expire_time < ? and is_deleted = ?", time.Now(), false).Find(&files); err != nil {
			log.Println("cron db query error", err.Error())
		}
		for _, file := range files {
			log.Printf("check IsDeleted file: %s/%s \n", file.FileId, file.FileName)
			//is_deleted를 true로 바꿔줌
			file.IsDeleted = true
			if _, err := database.Engine.ID(file.Id).Cols("Is_deleted").Update(&file); err != nil {
				log.Printf("cron db update error, file: %s/%s, error: %s\n", file.FileId, file.FileName, err.Error())
			}
		}
	})

	terminator.AddFunc("* */5 * * *", func() {
		var files []database.FileTracking
		// IsDeleted가 false인 파일만 가져옴
		if err := database.Engine.Where("is_deleted = ?", true).Find(&files); err != nil {
			log.Println("file list error: ", err.Error())
		}
		for _, file := range files {
			log.Printf("delete file: %s/%s\n", file.FileId, file.FileName)
			if err := os.RemoveAll("./tmp/" + file.FileId); err != nil {
				log.Println("delete file error: ", err.Error())
			}
			if _, err := database.Engine.Delete(&file); err != nil {
				log.Println("delete file error: ", err.Error())
			}
		}
	})

	terminator.Start()

	var err error

	if utils.CheckTmpFolder() != nil {
		log.Fatalf("tmp folder error: %v", err)
	}

	if database.CreateDBEngine() != nil {
		log.Fatalf("failed to create db engine: %v", err)
	}

	root := app.Group("/")

	routes.PublicRoutes(root)
	routes.PrivateRouter(root)

	utils.StartServer(app)
	terminator.Stop()
}
