package file

import (
	"time"

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
		c.JSON(200, gin.H{
			"message":       "file found",
			"uploadDate":    FileTracking.UploadDate.Format(time.RFC3339),
			"files":         files,
			"folderId":      FileTracking.FolderId,
			"downloadLimit": FileTracking.DownloadLimit,
			"downloadCount": FileTracking.DownloadCount,
			"expireTime":    FileTracking.ExpireTime.Format(time.RFC3339),
			"deleteUrl":     baseUrl + "/del/" + FileTracking.FolderId,
		})
	}
}
