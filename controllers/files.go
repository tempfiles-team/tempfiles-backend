package controller

import (
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-fuego/fuego"
	"github.com/tempfiles-Team/tempfiles-backend/database"
	"github.com/tempfiles-Team/tempfiles-backend/utils"
)

type FilesRessources struct {
	// TODO add ressources
	FilesService RealFilesService
}

type File struct {
	database.FileTracking

	Error   string `json:"error"`
	Message string `json:"message"`
}

type Files struct {
	List []database.FileTracking `json:"list"`

	Error   string `json:"error"`
	Message string `json:"message"`
}

type FilesCreate struct {
	Message string `json:"message"`
}

func (rs FilesRessources) RoutesV1(s *fuego.Server) {
	fuego.Get(s, "/list", rs.getAllFiles)

	fuego.Post(s, "/upload", rs.postFiles).
		Description("Upload files").
		Header("X-Download-Limit", "Download limit").
		Header("X-Time-Limit", "Time limit").
		Header("X-Hidden", "Hidden")

	fuego.GetStd(s, "/dl/{id}/{name}", rs.downloadFile)
	fuego.GetStd(s, "/view/{id}/{name}", rs.ViewFile)
	fuego.Get(s, "/file/{id}", rs.getFiles)
	fuego.Delete(s, "/del/{id}", rs.deleteFiles)
}

func (rs FilesRessources) getAllFiles(c fuego.ContextNoBody) (Files, error) {
	return rs.FilesService.GetAllFiles()
}

func (rs FilesRessources) postFiles(c *fuego.ContextWithBody[any]) (File, error) {

	return rs.FilesService.CreateFiles(c)
}

func (rs FilesRessources) getFiles(c fuego.ContextNoBody) (File, error) {
	return rs.FilesService.GetFiles(c.PathParam("id"))
}

func (rs FilesRessources) deleteFiles(c *fuego.ContextNoBody) (File, error) {
	return rs.FilesService.DeleteFiles(c.PathParam("id"))
}

func (rs FilesRessources) downloadFile(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	name := r.PathValue("name")

	path, err := rs.FilesService.DownloadFile(id, name)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+strings.ReplaceAll(url.PathEscape(name), "+", "%20"))
	http.ServeFile(w, r, path)
}

func (rs FilesRessources) ViewFile(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	name := r.PathValue("name")

	path, err := rs.FilesService.DownloadFile(id, name)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	mimeType := mime.TypeByExtension(filepath.Ext(name))

	// 만약 텍스트 파일인 경우에 (코드 및 기타 텍스트 파일) 파일을 열어서 plain text로 보여줍니다.
	if strings.Contains(mimeType, "text") {
		file, err := os.Open(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer file.Close()
		w.Header().Set("Content-Type", "text/plain")
		// file를 읽어서 response에 쓰기
		if _, err := io.Copy(w, file); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	} else {
		http.ServeFile(w, r, path)
	}

}

type FilesService interface {
	GetFiles(id string) (File, error)
	CreateFiles(FilesCreate) (File, error)
	GetAllFiles() (Files, error)
	DeleteFiles(id string) (any, error)
	DownloadFile(id string, name string) (File, error)
}

type RealFilesService struct {
	FilesService
}

func (s RealFilesService) GetFiles(id string) (File, error) {

	FileTracking := database.FileTracking{
		FolderId: id,
	}

	has, err := database.Engine.Get(&FileTracking)

	if err != nil {
		return File{
			Message: "db query error",
		}, err
	}

	if !has {
		return File{
			Message: "file not found",
		}, nil
	}

	if files, err := utils.GetFiles(FileTracking.FolderId); err != nil {
		return File{
			Message: "folder not found",
		}, nil
	} else {
		log.Println("✨  File found: ", FileTracking.FolderId)

		FileTracking.Files = files

		repons := File{
			Message: "File found",
		}

		repons.FileTracking = FileTracking

		return repons, nil
	}
}

func (s RealFilesService) DownloadFile(id string, name string) (path string, error error) {

	FileTracking := database.FileTracking{
		FolderId: id,
	}

	has, err := database.Engine.Get(&FileTracking)

	if err != nil {
		return "", err
	}

	if !has {
		return "", fmt.Errorf("folder not found")
	}

	if !utils.CheckIsFileExist(FileTracking.FolderId, name) {
		return "", fmt.Errorf("file not found")
	}

	// db DownloadCount +1
	FileTracking.DownloadCount++
	if _, err := database.Engine.ID(FileTracking.Id).Update(&FileTracking); err != nil {

		return "", err
	}

	if FileTracking.DownloadLimit != 0 && FileTracking.DownloadCount >= FileTracking.DownloadLimit {

		FileTracking.IsDeleted = true

		log.Printf("🗑️  Set this folder for deletion: %s \n", FileTracking.FolderId)
		if _, err := database.Engine.ID(FileTracking.Id).Cols("Is_deleted").Update(&FileTracking); err != nil {

			return "", err
		}
	}

	log.Printf("📥️  Successfully downloaded %s, %s\n", FileTracking.FolderId, name)

	return "tmp/" + FileTracking.FolderId + "/" + name, nil
}

func (s RealFilesService) CreateFiles(c *fuego.ContextWithBody[any]) (File, error) {

	err := c.Request().ParseMultipartForm(10 << 20) // limit file size to 10MB
	if err != nil {
		return File{
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
	MFAHASH, err := utils.FormFiles(c.Request(), "file")
	if err != nil {
		return File{
			Message: fmt.Sprintf("Please send the file using the “file” field in multipart/form-data.: %v", err),
		}, nil
	}

	FolderHash, err := utils.GenIdFormMulitpart(MFAHASH)
	if err != nil {
		return File{
			Message: fmt.Sprintf("folder id generation error: %v", err),
		}, nil
	}

	isExist, err := database.Engine.Exist(&database.FileTracking{FolderHash: FolderHash})

	if err != nil {
		return File{
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

			return File{
				Error:   err.Error(),
				Message: fmt.Sprintf("database get error: %v", err),
			}, nil
		}

		FileTracking.ExpireTime = expireTimeDate.Unix()
		if _, err := database.Engine.ID(FileTracking.Id).Cols("expire_time").Update(&FileTracking); err != nil {
			resp := File{
				Error:   err.Error(),
				Message: fmt.Sprintf("database update error: %v", err),
			}

			resp.FileTracking = FileTracking

			return resp, nil
		}

		FileTracking.DownloadCount = 0
		if _, err := database.Engine.ID(FileTracking.Id).Cols("download_count").Update(&FileTracking); err != nil {
			resp := File{
				Error:   err.Error(),
				Message: fmt.Sprintf("database update error: %v", err),
			}

			resp.FileTracking = FileTracking

			return resp, nil
		}

		resp := File{
			Message: fmt.Sprintf("File %s already exists, reset time limit and download count.", FileTracking.FolderHash[:5]),
		}

		resp.FileTracking = FileTracking

		return resp, nil
	}

	MFAH, err := utils.FormFiles(c.Request(), "file")
	if err != nil {
		return File{
			Message: fmt.Sprintf("Please send the file using the “file” field in multipart/form-data.: %v", err),
		}, nil
	}

	filesList := make([]database.FileListResponse, 0)
	for _, file := range MFAH {
		filesList = append(filesList, database.FileListResponse{
			FileName: file.Header.Filename,
			FileSize: file.Header.Size,
		})
	}

	FileTracking := &database.FileTracking{
		FileCount:     len(MFAH),
		FolderId:      FolderHash[:5],
		IsHidden:      isHidden,
		FolderHash:    FolderHash,
		UploadDate:    time.Now(),
		DownloadLimit: int64(downloadLimit),
		ExpireTime:    expireTimeDate.Unix(),
		Files:         filesList,
	}

	if utils.CheckFileFolder(FileTracking.FolderId) != nil {

		return File{
			Error:   err.Error(),
			Message: fmt.Sprintf("file folder duplication error: %v", err),
		}, nil

	}

	for _, file := range MFAH {

		if err := utils.SaveFile(FileTracking.FolderId, file.Header.Filename, file.File); err != nil {

			return File{
				Error:   err.Error(),
				Message: fmt.Sprintf("file save error: %v", err),
			}, nil
		}
	}

	_, err = database.Engine.Insert(FileTracking)
	if err != nil {

		return File{
			Error:   err.Error(),
			Message: fmt.Sprintf("database insert error: %v", err),
		}, nil
	}

	log.Printf("🥰  Successfully uploaded %s, %d files\n", FileTracking.FolderId, FileTracking.FileCount)

	resp := File{
		Message: fmt.Sprintf("File %s uploaded successfully", FileTracking.FolderHash[:5]),
	}

	resp.FileTracking = *FileTracking

	return resp, nil
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
		List:    files,
		Message: "File list successfully",
	}, nil
}

func (s RealFilesService) DeleteFiles(id string) (File, error) {
	FileTracking := database.FileTracking{
		FolderId: id,
	}

	has, err := database.Engine.Get(&FileTracking)

	if err != nil {
		return File{
			Message: "db query error",
			Error:   err.Error(),
		}, nil
	}

	if !has {
		return File{
			Message: "file not found",
			Error:   "true",
		}, nil
	}

	if err := os.RemoveAll("tmp/" + FileTracking.FolderId); err != nil {

		return File{
			Message: "file delete error",
			Error:   err.Error(),
		}, nil
	}

	if _, err := database.Engine.Delete(&FileTracking); err != nil {

		return File{
			Message: "db delete error",
			Error:   err.Error(),
		}, nil
	}

	log.Printf("🗑️  Delete this folder: %s\n", FileTracking.FolderId)

	return File{
		Message: "File deleted successfully",
	}, nil
}
