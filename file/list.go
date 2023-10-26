package file

import (
	"github.com/gin-gonic/gin"
	"github.com/tempfiles-Team/tempfiles-backend/database"
)

func ListHandler(c *gin.Context) {

	var files []database.FileTracking

	if err := database.Engine.Where("is_deleted = ? AND is_hidden = ?", false, false).Find(&files); err != nil {
		c.JSON(500, gin.H{
			"message": "db query error",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "File list successfully",
		"list":    files,
	})
}
