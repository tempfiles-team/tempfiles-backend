package file

import (
	"github.com/tempfiles-Team/tempfiles-backend/database"
)

type FileListResponse struct {
	FileName    string `json:"fileName"`
	FileSize    int64  `json:"fileSize"`
	DownloadUrl string `json:"downloadUrl"`
}

type FileResponse struct {
	database.FileTracking
	DeleteUrl string             `json:"deleteUrl"`
	Files     []FileListResponse `json:"files"`
}

func (f *FileResponse) NewFileResponse(fileTracking database.FileTracking, files []FileListResponse, message, baseUrl string) *FileResponse {

	f.DeleteUrl = baseUrl + "/del/" + fileTracking.FolderId
	f.FileTracking = fileTracking
	f.Files = files
	return f

}
