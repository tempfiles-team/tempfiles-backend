package crontab

import (
	"log"
	"os"
	"time"

	"github.com/robfig/cron"
	"github.com/tempfiles-Team/tempfiles-backend/database"
)

func Crontab() {
	terminator := cron.New()

	terminator.AddFunc("@every 30s", func() {
		log.Println("⏲️  Check for expired files", time.Now().Format("2006-01-02 15:04:05"))
		var files []database.FileTracking
		if err := database.Engine.Where("expire_time < ? and is_deleted = ?", time.Now().Unix(), false).Find(&files); err != nil {
			log.Println("cron db query error", err.Error())
		}
		for _, file := range files {
			log.Printf("🗑️  Set this folder for deletion: %s \n", file.FolderId)
			file.IsDeleted = true
			if _, err := database.Engine.ID(file.Id).Cols("Is_deleted").Update(&file); err != nil {
				log.Printf("cron db update error, file: %s, error: %s\n", file.FolderId, err.Error())
			}
		}
	})

	terminator.AddFunc("@every 1m", func() {
		log.Println("⏲️  Check which files need to be deleted", time.Now().Format("2006-01-02 15:04:05"))
		var files []database.FileTracking
		if err := database.Engine.Where("is_deleted = ?", true).Find(&files); err != nil {
			log.Println("file list error: ", err.Error())
		}
		for _, file := range files {
			log.Printf("🗑️  Delete this folder: %s\n", file.FolderId)
			if err := os.RemoveAll("./tmp/" + file.FolderId); err != nil {
				log.Println("delete file error: ", err.Error())
			}
			if _, err := database.Engine.Delete(&file); err != nil {
				log.Println("delete file error: ", err.Error())
			}
		}
	})

	terminator.Start()
}
