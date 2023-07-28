package queries

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/tempfiles-Team/tempfiles-backend/app/models"
	"github.com/tempfiles-Team/tempfiles-backend/platform/db"
)

type FileState struct {
	Model models.FileTracking
}

func (s *FileState) GetFile(fileId string) (bool, error) {
	s.Model = models.FileTracking{FileId: fileId}
	has, err := db.Engine.Get(&s.Model)
	return has, err
}

func (s *FileState) DelFile() error {
	if s.Model.FileId == "" {
		return errors.New("fileId is empty, GetFile first")
	}
	_, err := db.Engine.Delete(&s.Model)
	return err
}

func (s *FileState) IncreaseDLCount() error {
	if s.Model.FileId == "" {
		return errors.New("fileId is empty, GetFile first")
	}
	s.Model.DownloadCount++
	_, err := db.Engine.ID(s.Model.Id).Update(&s.Model)
	return err
}

func (s *FileState) IsExpiredFile() (bool, error) {
	if s.Model.FileId == "" {
		return false, errors.New("fileId is empty, GetFile first")
	}

	if s.Model.DownloadLimit != 0 && s.Model.DownloadCount >= s.Model.DownloadLimit {
		s.Model.IsDeleted = true

		log.Printf("check IsDeleted file: %s/%s \n", s.Model.FileId, s.Model.FileName)

		_, err := db.Engine.ID(s.Model.Id).Cols("is_deleted").Update(&s.Model)
		return true, err
	}
	return false, nil
}

func IsExpiredFiles() error {
	var files []models.FileTracking
	if err := db.Engine.Where("expire_time < ? and is_deleted = ?", time.Now(), false).Find(&files); err != nil {
		log.Println("cron db query error", err.Error())
	}

	for _, file := range files {
		log.Printf("check IsDeleted file: %s/%s \n", file.FileId, file.FileName)
		file.IsDeleted = true

		_, err := db.Engine.ID(file.Id).Cols("is_deleted").Update(&file)

		if err != nil {
			log.Printf("cron db update error, file: %s/%s, error: %s\n", file.FileId, file.FileName, err.Error())
		}
	}

	return nil
}

func DelExpireFiles() error {
	var files []models.FileTracking
	// IsDeleted가 false인 파일만 가져옴
	if err := db.Engine.Where("is_deleted = ?", true).Find(&files); err != nil {
		log.Println("file list error: ", err.Error())
	}
	for _, file := range files {
		log.Printf("delete file: %s/%s\n", file.FileId, file.FileName)
		if err := os.RemoveAll("./tmp/" + file.FileId); err != nil {
			log.Println("delete file error: ", err.Error())
		}
		if _, err := db.Engine.Delete(&file); err != nil {
			log.Println("delete file error: ", err.Error())
		}
	}
	return nil
}

func GetFiles() ([]models.FileTracking, error) {
	var files []models.FileTracking
	err := db.Engine.Where("is_deleted = ?", false).Find(&files)
	return files, err
}

func (s *FileState) InsertFile() error {
	_, err := db.Engine.Insert(&s.Model)
	return err
}
