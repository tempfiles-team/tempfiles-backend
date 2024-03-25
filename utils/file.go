package utils

import (
	"io"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/tempfiles-Team/tempfiles-backend/database"
)

func CheckIsFileExist(folderId, fileName string) bool {
	if _, err := os.Stat("tmp/" + folderId + "/" + fileName); os.IsNotExist(err) {
		return false
	}
	return true
}

func GetFiles(folderId string) ([]database.FileListResponse, error) {
	// return filenames, file sizes
	var files []database.FileListResponse

	err := filepath.Walk("tmp/"+folderId, func(path string, info os.FileInfo, err error) error {
		if path != "tmp/"+folderId {
			files = append(files, database.FileListResponse{
				FileName:    filepath.Base(path),
				FileSize:    info.Size(),
				DownloadUrl: "/dl/" + folderId + "/" + strings.ReplaceAll(url.PathEscape(filepath.Base(path)), "+", "%20"),
			})
		}
		return nil
	})
	return files, err

}

func SaveFile(folderId, fileName string, file multipart.File) error {
	// tmpf/debug/filename write
	os.MkdirAll("tmp/"+folderId, os.ModePerm)
	tmpf, err := os.Create("tmp/" + folderId + "/" + fileName)

	if err != nil {
		return err
	}

	defer tmpf.Close()

	_, err = io.Copy(tmpf, file)
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
