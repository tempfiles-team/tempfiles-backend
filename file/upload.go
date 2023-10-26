package file

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tempfiles-Team/tempfiles-backend/database"
)

func UploadHandler(c *gin.Context) {

	form, err := c.MultipartForm()

	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	if form == nil || len(form.File["file"]) == 0 {
		c.JSON(400, gin.H{
			"error": "Please send the file using the “file” field in multipart/form-data.",
		})
		return
	}

	isHidden, err := strconv.ParseBool(c.GetHeader("X-Hidden"))
	if err != nil {
		isHidden = false
	}

	downloadLimit, err := strconv.Atoi(c.GetHeader("X-Download-Limit"))
	if err != nil {
		downloadLimit = 100
	}
	expireTime, err := strconv.Atoi(c.GetHeader("X-Time-Limit"))

	var expireTimeDate time.Time

	if err != nil || expireTime <= 0 {
		expireTimeDate = time.Now().Add(time.Duration(60*3) * time.Minute)
	} else {
		expireTimeDate = time.Now().Add(time.Duration(expireTime) * time.Minute)
	}

	FolderHash, err := GenerateFolderId(form.File["file"])

	if err != nil {
		c.JSON(500, gin.H{
			"message": "folder id generation error",
			"error":   err.Error(),
		})
		return
	}

	isExist, err := database.Engine.Exist(&database.FileTracking{FolderHash: FolderHash})

	if err != nil {
		c.JSON(500, gin.H{
			"message": "database exist error",
			"error":   err.Error(),
		})
		return
	}

	if isExist {
		FileTracking := database.FileTracking{
			FolderHash: FolderHash,
		}
		_, err := database.Engine.Get(&FileTracking)
		if err != nil {
			c.JSON(500, gin.H{
				"message": "database get error",
				"error":   err.Error(),
			})
			return
		}
		c.JSON(500, gin.H{
			"message":       "Duplicate entries with matching file names and contents have already been uploaded.",
			"folderId":      FileTracking.FolderId,
			"uploadDate":    FileTracking.UploadDate.Format(time.RFC3339),
			"downloadLimit": FileTracking.DownloadLimit,
			"downloadCount": FileTracking.DownloadCount,
			"expireTime":    FileTracking.ExpireTime.Format(time.RFC3339),
			"error":         nil,
		})
		return
	}

	FileTracking := &database.FileTracking{
		FileCount:     len(form.File["file"]),
		FolderId:      FolderHash[:5],
		IsHidden:      isHidden,
		FolderHash:    FolderHash,
		UploadDate:    time.Now(),
		DownloadLimit: int64(downloadLimit),
		ExpireTime:    expireTimeDate,
	}

	if CheckFileFolder(FileTracking.FolderId) != nil {
		c.JSON(500, gin.H{
			"message": "file folder creation error",
			"error":   err.Error(),
		})
		return
	}

	for _, file := range form.File["file"] {
		if err := c.SaveUploadedFile(file, fmt.Sprintf("tmp/%s/%s", FileTracking.FolderId, file.Filename)); err != nil {
			c.JSON(500, gin.H{
				"message": "file save error",
				"error":   err.Error(),
			})
			return
		}
	}

	_, err = database.Engine.Insert(FileTracking)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "database insert error",
			"error":   err.Error(),
		})
		return
	}

	log.Printf("Successfully uploaded %s, download limit %d\n", FileTracking.FolderId, FileTracking.DownloadLimit)

	c.JSON(200, gin.H{
		"message":       "File uploaded successfully",
		"folderId":      FileTracking.FolderId,
		"uploadDate":    FileTracking.UploadDate.Format(time.RFC3339),
		"downloadLimit": FileTracking.DownloadLimit,
		"downloadCount": FileTracking.DownloadCount,
		"expireTime":    FileTracking.ExpireTime.Format(time.RFC3339),
	})
}
