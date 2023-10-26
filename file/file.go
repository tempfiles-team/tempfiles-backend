package file

import (
	"github.com/gin-gonic/gin"
	"github.com/tempfiles-Team/tempfiles-backend/database"
)

func FileHandler(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(400, gin.H{
			"message": "Please provide a file id",
			"error":   nil,
		})
		return
	}

	FileTracking := database.FileTracking{
		FolderId: id,
	}

	has, err := database.Engine.Get(&FileTracking)

	if err != nil {
		c.JSON(500, gin.H{
			"message": "db query error",
			"error":   err.Error(),
		})
		return
	}

	if !has {
		c.JSON(404, gin.H{
			"message": "file not found",
			"error":   nil,
		})
		return
	}

	baseUrl := c.Request.Host

	if files, err := GetFiles(FileTracking.FolderId, baseUrl); err != nil {
		c.JSON(500, gin.H{
			"message": "folder not found",
			"error":   nil,
		})
		return
	} else {
		c.JSON(200, new(FileResponse).NewFileResponse(FileTracking, files, "file found", baseUrl))
	}
}
