package file

import (
	"github.com/gin-gonic/gin"
	"github.com/tempfiles-Team/tempfiles-backend/database"
)

func ListHandler(c *gin.Context) {

	var files []database.FileTracking
	// IsDeleted가 false인 파일만 가져옴
	if err := database.Engine.Where("is_deleted = ?", false).Find(&files); err != nil {
		c.JSON(400, gin.H{
			"message": "file list error",
			"error":   err.Error(),
		})
	}

	c.JSON(200, gin.H{
		"message":      "File list successfully",
		"list":         files,
		"numberOfList": len(files),
	})
}
