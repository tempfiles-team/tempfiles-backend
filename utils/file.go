package utils

import (
	"io"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/tempfiles-Team/tempfiles-backend/file"
)

func CheckIsFileExist(folderId, fileName string) bool {
	if _, err := os.Stat("tmp/" + folderId + "/" + fileName); os.IsNotExist(err) {
		return false
	}
	return true
}

func GetFiles(folderId, baseUrl string) ([]file.FileListResponse, error) {
	// return filenames, file sizes
	var files []file.FileListResponse

	err := filepath.Walk("tmp/"+folderId, func(path string, info os.FileInfo, err error) error {
		if path != "tmp/"+folderId {
			files = append(files, file.FileListResponse{
				FileName:    filepath.Base(path),
				FileSize:    info.Size(),
				DownloadUrl: baseUrl + "/dl/" + folderId + "/" + strings.ReplaceAll(url.PathEscape(filepath.Base(path)), "+", "%20"),
			})
		}
		return nil
	})
	return files, err

}

func SaveFile(folderId, fileName string, file *multipart.File) error {

	// Create a new file
	newFile, err := os.Create("tmp/" + folderId + "/" + fileName)
	if err != nil {
		return err

	}

	defer newFile.Close()

	// Copy the file to the new file
	_, err = io.Copy(newFile, *file)
	if err != nil {
		return err

	}

	return nil

}

func CheckTmpFolder() error {
	if _, err := os.Stat("tmp"); os.IsNotExist(err) {
		err := os.Mkdir("tmp", 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func CheckFileFolder(folderId string) error {
	if _, err := os.Stat("tmp/" + folderId); os.IsNotExist(err) {
		err := os.MkdirAll("tmp/"+folderId, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}
