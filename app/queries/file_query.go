package queries

import (
	"errors"
	"log"
	"os"
	"time"

	"gorm.io/gorm"

	"github.com/tempfiles-Team/tempfiles-backend/app/models"
	db "github.com/tempfiles-Team/tempfiles-backend/database"
)

type FileState struct {
	Model models.FileTracking
}

func (s *FileState) GetFile(fileId string) (bool, error) {
	result := db.Engine.Where("file_id = ?", fileId).First(&s.Model)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return true, result.Error
}

func (s *FileState) DelFile() error {
	if s.Model.FileId == "" {
		return errors.New("fileId is empty, GetFile first")
	}
	result := db.Engine.Delete(&s.Model)
	return result.Error
}

func (s *FileState) IncreaseDLCount() error {
	if s.Model.FileId == "" {
		return errors.New("fileId is empty, GetFile first")
	}
	s.Model.DownloadCount++
	result := db.Engine.Save(&s.Model)
	return result.Error
}

func GetFiles() ([]models.FileTracking, error) {
	var files []models.FileTracking
	err := db.Engine.Where("is_deleted = ?", false).Find(&files).Error
	return files, err
}

func (s *FileState) InsertFile() error {
	result := db.Engine.Create(&s.Model)
	return result.Error
}

func (s *FileState) IsExpiredFile() (bool, error) {
	if s.Model.FileId == "" {
		return false, errors.New("fileId is empty, GetFile first")
	}

	if s.Model.DownloadLimit != 0 && s.Model.DownloadCount >= s.Model.DownloadLimit {
		s.Model.IsDeleted = true

		log.Printf("check IsDeleted file: %s/%s \n", s.Model.FileId, s.Model.FileName)

		result := db.Engine.Model(&s.Model).Where("id = ?", s.Model.Id).Update("is_deleted", s.Model.IsDeleted)
		return true, result.Error
	}
	return false, nil
}

func IsExpiredFiles() error {
	var files []models.FileTracking
	if err := db.Engine.Where("expire_time < ? and is_deleted = ?", time.Now(), false).Find(&files).Error; err != nil {
		log.Println("cron db query error", err.Error())
	}

	for _, file := range files {
		log.Printf("check IsDeleted file: %s/%s \n", file.FileId, file.FileName)
		file.IsDeleted = true

		result := db.Engine.Model(&file).Where("id = ?", file.Id).Update("is_deleted", file.IsDeleted)

		if result.Error != nil {
			log.Printf("cron db update error, file: %s/%s, error: %s\n", file.FileId, file.FileName, result.Error.Error())
		}
	}

	return nil
}

func DelExpireFiles() error {
	var files []models.FileTracking
	// IsDeleted가 false인 파일만 가져옴
	if err := db.Engine.Where("is_deleted = ?", true).Find(&files).Error; err != nil {
		log.Println("file list error: ", err.Error())
	}
	for _, file := range files {
		log.Printf("delete file: %s/%s\n", file.FileId, file.FileName)
		if err := os.RemoveAll("./tmp/" + file.FileId); err != nil {
			log.Println("delete file error: ", err.Error())
		}
		if err := db.Engine.Delete(&file).Error; err != nil {
			log.Println("delete file error: ", err.Error())
		}
	}
	return nil
}
