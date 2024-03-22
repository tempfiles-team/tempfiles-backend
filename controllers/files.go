package controller

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-fuego/fuego"
	"github.com/tempfiles-Team/tempfiles-backend/database"
	"github.com/tempfiles-Team/tempfiles-backend/utils"
)

type FilesRessources struct {
	// TODO add ressources
	FilesService RealFilesService
}

type Files struct {
	Data    interface{} `json:"data"`
	Error   string      `json:"error"`
	Message string      `json:"message"`
}

type FilesCreate struct {
	Message string `json:"message"`
}

func (rs FilesRessources) Routes(s *fuego.Server) {
	filesGroup := fuego.Group(s, "/files")

	fuego.Get(filesGroup, "/", rs.getAllFiles)
	fuego.Post(filesGroup, "/", rs.postFiles)

	fuego.Get(filesGroup, "/{id}/{name}", rs.getFiles)
	fuego.Get(filesGroup, "/{id}", rs.getFiles)
	fuego.Delete(filesGroup, "/{id}", rs.deleteFiles)
}

func (rs FilesRessources) getAllFiles(c fuego.ContextNoBody) (Files, error) {
	return rs.FilesService.GetAllFiles()
}

func (rs FilesRessources) postFiles(c *fuego.ContextWithBody[any]) (Files, error) {

	return rs.FilesService.CreateFiles(c)
}

func (rs FilesRessources) getFiles(c fuego.ContextNoBody) (Files, error) {
	return rs.FilesService.GetFiles(c.PathParam("id"))
}

func (rs FilesRessources) deleteFiles(c *fuego.ContextNoBody) (any, error) {
	return rs.FilesService.DeleteFiles(c.PathParam("id"))
}

type FilesService interface {
	GetFiles(id string) (Files, error)
	CreateFiles(FilesCreate) (Files, error)
	GetAllFiles() ([]Files, error)
	DeleteFiles(id string) (any, error)
}

type RealFilesService struct {
	FilesService
}

func (s RealFilesService) GetFiles(id string) (Files, error) {
	// TODO implement

	return Files{}, nil
}

func (s RealFilesService) CreateFiles(c *fuego.ContextWithBody[any]) (Files, error) {

	err := c.Request().ParseMultipartForm(10 << 20) // limit file size to 10MB
	if err != nil {
		return Files{
			Message: fmt.Sprintf("Error parsing file: %v", err),
		}, nil
	}

	isHidden, err := strconv.ParseBool(c.Header("X-Hidden"))
	if err != nil {
		isHidden = false
	}

	downloadLimit, err := strconv.Atoi(c.Header("X-Download-Limit"))
	if err != nil {
		downloadLimit = 100
	}
	expireTime, err := strconv.Atoi(c.Header("X-Time-Limit"))

	var expireTimeDate time.Time

	if err != nil || expireTime <= 0 {
		expireTimeDate = time.Now().Add(time.Duration(60*3) * time.Minute)
	} else {
		expireTimeDate = time.Now().Add(time.Duration(expireTime) * time.Minute)
	}

	// Multipart File And Header
	MFAH, err := utils.FormFiles(c.Request(), "file")
	if err != nil {
		return Files{
			Message: fmt.Sprintf("Please send the file using the “file” field in multipart/form-data.: %v", err),
		}, nil
	}

	FolderHash, err := utils.GenIdFormMulitpart(MFAH)
	if err != nil {
		return Files{
			Message: fmt.Sprintf("folder id generation error: %v", err),
		}, nil
	}

	isExist, err := database.Engine.Exist(&database.FileTracking{FolderHash: FolderHash})

	if err != nil {
		return Files{
			Error:   err.Error(),
			Message: fmt.Sprintf("database exist error: %v", err),
		}, nil
	}

	if isExist {
		FileTracking := database.FileTracking{
			FolderHash: FolderHash,
		}
		_, err := database.Engine.Get(&FileTracking)
		if err != nil {

			return Files{
				Error:   err.Error(),
				Message: fmt.Sprintf("database get error: %v", err),
			}, nil
		}

		return Files{
			Data:    FileTracking,
			Message: fmt.Sprintf("File %s already exists", FileTracking.FolderHash),
		}, nil
	}

	FileTracking := &database.FileTracking{
		FileCount:     len(MFAH),
		FolderId:      FolderHash[:5],
		IsHidden:      isHidden,
		FolderHash:    FolderHash,
		UploadDate:    time.Now(),
		DownloadLimit: int64(downloadLimit),
		ExpireTime:    expireTimeDate,
	}

	if utils.CheckFileFolder(FileTracking.FolderId) != nil {

		return Files{
			Error:   err.Error(),
			Message: fmt.Sprintf("file folder duplication error: %v", err),
		}, nil

	}

	for _, file := range MFAH {

		if err := utils.SaveFile(FileTracking.FolderId, file.Header.Filename, &file.File); err != nil {

			return Files{
				Error:   err.Error(),
				Message: fmt.Sprintf("file save error: %v", err),
			}, nil
		}
	}

	_, err = database.Engine.Insert(FileTracking)
	if err != nil {

		return Files{
			Error:   err.Error(),
			Message: fmt.Sprintf("database insert error: %v", err),
		}, nil
	}

	log.Printf("🥰  Successfully uploaded %s, %d files\n", FileTracking.FolderId, FileTracking.FileCount)

	return Files{
		Data:    FileTracking,
		Message: fmt.Sprintf("File %s uploaded successfully", FileTracking.FolderHash),
	}, nil

}

func (s RealFilesService) GetAllFiles() (Files, error) {
	// TODO implement

	var files []database.FileTracking

	if err := database.Engine.Where("is_deleted = ? AND is_hidden = ?", false, false).Find(&files); err != nil {
		return Files{
			Message: "db query error",
			Error:   err.Error(),
		}, nil
	}

	return Files{
		Data:    files,
		Message: "File list successfully",
	}, nil
}

func (s RealFilesService) DeleteFiles(id string) (any, error) {
	// TODO implement
	return nil, nil
}
