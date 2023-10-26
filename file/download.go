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
	}

	if id == "" {
		c.JSON(400, gin.H{
			"message":  "Please provide a file id",
			"error":    nil,
			"download": false,
		})
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
	}

	if !has {
		c.JSON(404, gin.H{
			"message": "folder not found",
			"error":   nil,
		})
	}

	if !CheckIsFileExist(FileTracking.FolderId, name) {
		c.JSON(404, gin.H{
			"message": "file not found!!",
			"error":   nil,
		})
	}

	// db DownloadCount +1
	FileTracking.DownloadCount++
	if _, err := database.Engine.ID(FileTracking.Id).Update(&FileTracking); err != nil {
		c.JSON(500, gin.H{
			"message": "db update error",
			"error":   err.Error(),
		})
	}

	// Download Limit check
	if FileTracking.DownloadLimit != 0 && FileTracking.DownloadCount >= FileTracking.DownloadLimit {
		// Download Limit exceeded -> check IsDelete
		FileTracking.IsDeleted = true

		log.Printf("check IsDeleted file: %s \n", FileTracking.FolderId)
		if _, err := database.Engine.ID(FileTracking.Id).Cols("Is_deleted").Update(&FileTracking); err != nil {
			c.JSON(500, gin.H{
				"message": "db update error",
				"error":   err.Error(),
			})
		}
	}

	c.Header("Content-Disposition", "attachment; filename="+strings.ReplaceAll(url.PathEscape(name), "+", "%20"))
	c.File("tmp/" + FileTracking.FolderId + "/" + name)
}
