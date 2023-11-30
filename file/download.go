package file

import (
	"log"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tempfiles-Team/tempfiles-backend/database"
)

func DownloadHandler(c *gin.Context) {
	id := c.Param("id")
	name, err := url.PathUnescape(c.Param("name"))
	if err != nil {
		c.JSON(400, gin.H{
			"message": "invalid file name",
			"error":   err.Error(),
		})
		return
	}

	if id == "" {
		c.JSON(400, gin.H{
			"message":  "Please provide a file id",
			"error":    nil,
			"download": false,
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
			"message": "folder not found",
			"error":   nil,
		})
		return
	}

	if !CheckIsFileExist(FileTracking.FolderId, name) {
		c.JSON(404, gin.H{
			"message": "file not found!!",
			"error":   nil,
		})
		return
	}

	// db DownloadCount +1
	FileTracking.DownloadCount++
	if _, err := database.Engine.ID(FileTracking.Id).Update(&FileTracking); err != nil {
		c.JSON(500, gin.H{
			"message": "db update error",
			"error":   err.Error(),
		})
		return
	}

	if FileTracking.DownloadLimit != 0 && FileTracking.DownloadCount >= FileTracking.DownloadLimit {

		FileTracking.IsDeleted = true

		log.Printf("🗑️  Set this folder for deletion: %s \n", FileTracking.FolderId)
		if _, err := database.Engine.ID(FileTracking.Id).Cols("Is_deleted").Update(&FileTracking); err != nil {
			c.JSON(500, gin.H{
				"message": "db update error",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(404, gin.H{
			"message": "file not found",
		})
		return
	}

	log.Printf("📥️  Successfully downloaded %s, %s\n", FileTracking.FolderId, name)

	c.Header("Content-Disposition", "attachment; filename="+strings.ReplaceAll(url.PathEscape(name), "+", "%20"))
	c.File("tmp/" + FileTracking.FolderId + "/" + name)
}
