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

// GetInfo godoc
// @Summary Show the information of server.
// @Description get the information of server.
// @Tags root
// @Accept */*
// @Produce json
// @Param api query string false "api name"
// @Success 200 {object} utils.Response
// @Router /info [get]
func GetInfo(c *fiber.Ctx) error {
	apiName := c.Query("api", "")
	backendUrl := c.BaseURL()
	switch apiName {
	case "upload":
		return c.JSON(utils.NewSuccessDataResponse(fiber.Map{
			"apiName": "/upload",
			"method":  "POST",
			"desc":    "특정 파일을 서버에 업로드합니다.",
			"command": "curl -X POST -F 'file=@[filepath or filename]' " + backendUrl + "/upload",
		}))
	case "list":
		return c.JSON(utils.NewSuccessDataResponse(fiber.Map{
			"apiName": "/list",
			"method":  "GET",
			"desc":    "서버에 존재하는 모든 파일에 대한 세부 정보를 반환합니다.",
			"command": "curl " + backendUrl + "/list",
		}))
	case "file":
		return c.JSON(utils.NewSuccessDataResponse(fiber.Map{
			"apiName": "/file/[file_id]",
			"method":  "GET",
			"desc":    "서버에 존재하는 특정 파일에 대한 세부 정보를 반환합니다.",
			"command": "curl " + backendUrl + "/file/[file_id]",
		}))
	case "del":
		return c.JSON(utils.NewSuccessDataResponse(fiber.Map{
			"apiName": "/del/[file_id]",
			"method":  "DELETE",
			"desc":    "서버에 존재하는 특정 파일을 삭제합니다.",
			"command": "curl -X DELETE " + backendUrl + "/del/[file_id]",
		}))
	case "dl":
		return c.JSON(utils.NewSuccessDataResponse(fiber.Map{
			"apiName": "/dl/[file_id]",
			"method":  "GET",
			"desc":    "서버에 존재하는 특정 파일을 다운로드 합니다.",
			"command": "curl -O " + backendUrl + "/dl/[file_id]",
		}))
	case "":
		return c.JSON(utils.NewSuccessDataResponse([]fiber.Map{
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
		}))
	default:
		return c.JSON(utils.NewFailMessageResponse("invalid api name"))

	}
}
