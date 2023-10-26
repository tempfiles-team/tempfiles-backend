package file

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tempfiles-Team/tempfiles-backend/database"
)

func DeleteHandler(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(400, gin.H{
			"message": "Please provide a file id",
			"error":   nil,
			"delete":  false,
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

	if err := os.RemoveAll("tmp/" + FileTracking.FolderId); err != nil {
		c.JSON(500, gin.H{
			"message": "file delete error",
			"error":   err.Error(),
			"delete":  false,
		})
	}

	if _, err := database.Engine.Delete(&FileTracking); err != nil {
		c.JSON(500, gin.H{
			"message": "db delete error",
			"error":   err.Error(),
			"delete":  false,
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "File deleted successfully",
		"error":   nil,
		"delete":  true,
	})
}
