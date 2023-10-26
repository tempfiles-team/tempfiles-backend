package main

import (
	"fmt"
	"log"
	"math"
	"os"

	"github.com/gin-contrib/cors"
	limits "github.com/gin-contrib/size"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/tempfiles-Team/tempfiles-backend/database"
	"github.com/tempfiles-Team/tempfiles-backend/file"
)

func main() {
	app := gin.Default()
	app.Use(limits.RequestSizeLimiter(int64(math.Pow(1024, 3)))) // 1 == 1byte, = 1GB

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "X-Download-Limit", "X-Time-Limit", "X-Hidden"}
	config.AllowMethods = []string{"GET", "POST", "DELETE"}
	app.Use(cors.New(config))

	// terminator := cron.New()
	// terminator.AddFunc("* */1 * * *", func() {
	// 	var files []database.FileTracking
	// 	//현재 시간보다 expire_time이 작고 is_deleted가 false인 파일을 가져옴
	// 	if err := database.Engine.Where("expire_time < ? and is_deleted = ?", time.Now(), false).Find(&files); err != nil {
	// 		log.Println("cron db query error", err.Error())
	// 	}
	// 	for _, file := range files {
	// 		log.Printf("check IsDeleted file: %s \n", file.FolderId)
	// 		//is_deleted를 true로 바꿔줌
	// 		file.IsDeleted = true
	// 		if _, err := database.Engine.ID(file.Id).Cols("Is_deleted").Update(&file); err != nil {
	// 			log.Printf("cron db update error, file: %s, error: %s\n", file.FolderId, err.Error())
	// 		}
	// 	}
	// })

	// // terminator.AddFunc("@daily", func() {
	// terminator.AddFunc("* */5 * * *", func() {
	// 	var files []database.FileTracking
	// 	// IsDeleted가 false인 파일만 가져옴
	// 	if err := database.Engine.Where("is_deleted = ?", true).Find(&files); err != nil {
	// 		log.Println("file list error: ", err.Error())
	// 	}
	// 	for _, file := range files {
	// 		log.Printf("delete file: %s\n", file.FolderId)
	// 		if err := os.RemoveAll("./tmp/" + file.FolderId); err != nil {
	// 			log.Println("delete file error: ", err.Error())
	// 		}
	// 		if _, err := database.Engine.Delete(&file); err != nil {
	// 			log.Println("delete file error: ", err.Error())
	// 		}
	// 	}
	// })

	// terminator.Start()

	var err error

	if file.CheckTmpFolder() != nil {
		log.Fatalf("tmp folder error: %v", err)
	}

	if database.CreateDBEngine() != nil {
		log.Fatalf("failed to create db engine: %v", err)
	}

	app.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "api is working normally :)",
		})
	})

	app.GET("/info", func(c *gin.Context) {
		apiName := c.Query("api")
		backendUrl := c.Request.Host
		switch apiName {
		case "upload":
			c.JSON(200, gin.H{
				"apiName": "/upload",
				"method":  "POST",
				"desc":    "특정 파일을 서버에 업로드합니다.",
				"command": "curl -X POST -F 'file=@[filepath or filename]' " + backendUrl + "/upload",
			})
		case "list":
			c.JSON(200, gin.H{
				"apiName": "/list",
				"method":  "GET",
				"desc":    "서버에 존재하는 파일 리스트를 반환합니다.",
				"command": "curl " + backendUrl + "/list",
			})
		case "file":
			c.JSON(200, gin.H{
				"apiName": "/file/[file_id]",
				"method":  "GET",
				"desc":    "서버에 존재하는 특정 파일에 대한 세부 정보를 반환합니다.",
				"command": "curl " + backendUrl + "/file/[file_id]",
			})
		case "del":
			c.JSON(200, gin.H{
				"apiName": "/del/[file_id]",
				"method":  "DELETE",
				"desc":    "서버에 존재하는 특정 파일을 삭제합니다.",
				"command": "curl -X DELETE " + backendUrl + "/del/[file_id]",
			})
		case "dl":
			c.JSON(200, gin.H{
				"apiName": "/dl/[file_id]",
				"method":  "GET",
				"desc":    "서버에 존재하는 특정 파일을 다운로드 합니다.",
				"command": "curl -O " + backendUrl + "/dl/[file_id]",
			})
		default:
			c.JSON(200, []gin.H{
				{
					"apiUrl":     backendUrl + "/upload",
					"apiHandler": "upload",
				},
				{

					"apiUrl":     backendUrl + "/list",
					"apiHandler": "list",
				},
				{
					"apiUrl":     backendUrl + "/file/[file_id]",
					"apiHandler": "file",
				},
				{
					"apiUrl":     backendUrl + "/del/[file_id]",
					"apiHandler": "del",
				},
				{
					"apiUrl":     backendUrl + "/dl/[file_id]",
					"apiHandler": "dl",
				},
			})
		}
	})

	app.GET("/list", file.ListHandler)
	app.POST("/upload", file.UploadHandler)

	// TODO: check id is valid (check is not delete)

	app.GET("/file/:id", file.FileHandler)

	app.GET("/dl/:id/:name", file.DownloadHandler)
	app.DELETE("/del/:id", file.DeleteHandler)

	if os.Getenv("BACKEND_PORT") == "" {
		os.Setenv("BACKEND_PORT", "5000")
	}

	log.Fatal(app.Run(fmt.Sprintf(":%s", os.Getenv("BACKEND_PORT"))))

	// terminator.Stop()
}
