package queries

import (
	"errors"
	"log"

	"github.com/tempfiles-Team/tempfiles-backend/app/models"
	"github.com/tempfiles-Team/tempfiles-backend/platform/db"
)

type TextState struct {
	Model models.TextTracking
}

func (s *TextState) GetText(fileId string) (bool, error) {
	s.Model = models.TextTracking{TextId: fileId}
	has, err := db.Engine.Get(&s.Model)
	return has, err
}

func (s *TextState) DelText() error {
	if s.Model.TextId == "" {
		return errors.New("textId is empty, GetText first")
	}
	_, err := db.Engine.Delete(&s.Model)
	return err
}

func (s *TextState) IncreaseDLCount() error {
	if s.Model.TextId == "" {
		return errors.New("textId is empty, GetText first")
	}
	s.Model.DownloadCount++
	_, err := db.Engine.ID(s.Model.Id).Update(&s.Model)
	return err
}

func (s *TextState) IsExpiredText() (bool, error) {
	if s.Model.TextId == "" {
		return false, errors.New("textId is empty, GetText first")
	}

	if s.Model.DownloadLimit != 0 && s.Model.DownloadCount >= s.Model.DownloadLimit {
		s.Model.IsDeleted = true

		log.Printf("check IsDeleted text: %s \n", s.Model.TextId)

		_, err := db.Engine.ID(s.Model.Id).Cols("is_deleted").Update(&s.Model)
		return true, err
	}
	return false, nil
}

func (s *TextState) InsertFile() error {
	if s.Model.TextId == "" {
		return errors.New("textId is empty, GetText first")
	}

	_, err := db.Engine.Insert(&s.Model)
	return err
}

func GetTexts() ([]models.TextTracking, error) {
	var texts []models.TextTracking
	err := db.Engine.Where("is_deleted = ?", false).Find(&texts)
	return texts, err
}
