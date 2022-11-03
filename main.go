package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strings"

	"github.com/minpeter/tempfiles-backend/database"
	"github.com/minpeter/tempfiles-backend/file"
	"github.com/minpeter/tempfiles-backend/jwt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/joho/godotenv/autoload"

	jwtware "github.com/gofiber/jwt/v3"
)

type LoginRequest struct {
	Email    string
	Password string
}

func main() {

	app := fiber.New(fiber.Config{
		AppName:   "tempfiles-backend",
		BodyLimit: int(math.Pow(1024, 3)), // 1 == 1byte
	})

	app.Use(
		cache.New(cache.Config{
			StoreResponseHeaders: true,
			Next: func(c *fiber.Ctx) bool {
				return c.Route().Path != "/dl/:filename"
			},
		}),
		cors.New(cors.Config{
			AllowOrigins: "*",
			AllowHeaders: "Origin, Content-Type, Accept",
			AllowMethods: "GET, POST, DELETE",
		}))

	var err error

	file.MinioClient, err = file.Connection()
	if err != nil {
		log.Fatalf("minio connection error: %v", err)
	}

	database.Engine, err = database.CreateDBEngine()
	if err != nil {
		log.Fatalf("failed to create db engine: %v", err)
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "api is working normally :)",
		})
	})

	app.Get("/info", func(c *fiber.Ctx) error {
		apiName := c.Query("api", "")
		switch apiName {
		case "upload":
			return c.JSON(fiber.Map{
				"apiName": "/upload",
				"method":  "POST",
				"desc":    "특정 파일을 서버에 업로드합니다.",
				"command": "curl -X POST -F 'file=@[filepath or filename]' https://tfb.minpeter.cf/upload",
			})
		case "list":
			return c.JSON(fiber.Map{
				"apiName": "/list",
				"method":  "GET",
				"desc":    "서버에 존재하는 파일 리스트를 반환합니다.",
				"command": "curl https://tfb.minpeter.cf/list",
			})
		case "del":
			return c.JSON(fiber.Map{
				"apiName": "/del/[filename]",
				"method":  "DELETE",
				"desc":    "서버에 존재하는 특정 파일을 삭제합니다.",
				"command": "curl -X DELETE https://tfb.minpeter.cf/del/[filename]",
			})
		case "dl":
			return c.JSON(fiber.Map{
				"apiName": "/dl/[filename]",
				"method":  "GET",
				"desc":    "서버에 존재하는 특정 파일을 다운로드 합니다.",
				"command": "curl -O https://tfb.minpeter.cf/dl/[filename]",
			})
		case "":
			backendUrl := os.Getenv("BACKEND_BASEURL")
			return c.JSON([]fiber.Map{
				{
					"apiUrl":     backendUrl + "/upload",
					"apiHandler": "upload",
				},
				{
					"apiUrl":     backendUrl + "/list",
					"apiHandler": "list",
				},
				{
					"apiUrl":     backendUrl + "/del/[filename]",
					"apiHandler": "del",
				},
				{
					"apiUrl":     backendUrl + "/dl/[filename]",
					"apiHandler": "dl",
				},
			})
		default:
			return c.JSON(fiber.Map{
				"message": "invalid api name",
			})

		}
	})

	app.Get("/list", file.ListHandler)

	app.Post("/upload", file.UploadHandler)
	app.Get("/checkpw/:filename", file.CheckPasswordHandler)

	app.Use(jwtware.New(jwtware.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized",
				"error":   err.Error(),
			})
		},
		SigningKey:  []byte(os.Getenv("JWT_SECRET")),
		TokenLookup: "query:token",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.JSON(fiber.Map{
				"message": "Missing or malformed JWT",
			})
		},
		Filter: func(c *fiber.Ctx) bool {
			fileName := strings.Split(strings.Split(c.OriginalURL(), "/")[2], "?")[0]

			log.Printf("c : %s\n c.OriginalURL() : %s\n", c.url, fileName)
			return jwt.IsEncrypted(fileName)
		},
	}))

	app.Get("/dl/:filename", file.OldDownloadHandler)
	app.Delete("/del/:filename", file.DeleteHandler)

	log.Fatal(app.Listen(fmt.Sprintf(":%s", os.Getenv("BACKEND_PORT"))))

}
