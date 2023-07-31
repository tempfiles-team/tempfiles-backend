package queries

import (
	"errors"
	"log"
	"time"

	"github.com/tempfiles-Team/tempfiles-backend/app/models"
	"github.com/tempfiles-Team/tempfiles-backend/platform/db"
	"gorm.io/gorm"
)

type TextState struct {
	Model models.TextTracking
}

func (s *TextState) GetText(textId string) (bool, error) {
	result := db.Engine.Where("text_id = ?", textId).First(&s.Model)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return true, result.Error
}

func (s *TextState) DelText() error {
	if s.Model.TextId == "" {
		return errors.New("textId is empty, GetText first")
	}
	result := db.Engine.Delete(&s.Model)
	return result.Error
}

func (s *TextState) IncreaseDLCount() error {
	if s.Model.TextId == "" {
		return errors.New("textId is empty, GetText first")
	}
	s.Model.DownloadCount++
	result := db.Engine.Save(&s.Model)
	return result.Error
}

func (s *TextState) IsExpiredText() (bool, error) {
	if s.Model.TextId == "" {
		return false, errors.New("textId is empty, GetText first")
	}

	if s.Model.DownloadLimit != 0 && s.Model.DownloadCount >= s.Model.DownloadLimit {
		s.Model.IsDeleted = true

		log.Printf("check IsDeleted text: %s \n", s.Model.TextId)

		result := db.Engine.Model(&s.Model).Where("id = ?", s.Model.Id).Update("is_deleted", s.Model.IsDeleted)
		return true, result.Error
	}
	return false, nil
}

func (s *TextState) InsertFile() error {
	if s.Model.TextId == "" {
		return errors.New("textId is empty, GetText first")
	}

	result := db.Engine.Create(&s.Model)
	return result.Error
}

func GetTexts() ([]models.TextTracking, error) {
	var texts []models.TextTracking
	err := db.Engine.Where("is_deleted = ?", false).Find(&texts).Error
	return texts, err
}

func IsExpiredTexts() error {
	var texts []models.TextTracking
	if err := db.Engine.Where("expire_time < ? and is_deleted = ?", time.Now(), false).Find(&texts).Error; err != nil {
		log.Println("cron db query error", err.Error())
	}

	for _, text := range texts {
		log.Printf("check IsDeleted text: %s \n", text.TextId)
		text.IsDeleted = true

		result := db.Engine.Model(&text).Where("id = ?", text.Id).Update("is_deleted", text.IsDeleted)

		if result.Error != nil {
			log.Printf("cron db update error, text: %s, error: %s\n", text.TextId, result.Error.Error())
		}
	}

	return nil
}

func DelExpireTexts() error {
	var texts []models.FileTracking
	// IsDeleted가 false인 파일만 가져옴
	if err := db.Engine.Where("is_deleted = ?", true).Find(&texts).Error; err != nil {
		log.Println("text list error: ", err.Error())
	}
	for _, text := range texts {
		if err := db.Engine.Delete(&text).Error; err != nil {
			log.Println("delete text error: ", err.Error())
		}
	}
	return nil
}
