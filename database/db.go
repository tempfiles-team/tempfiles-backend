package database

import (
	"fmt"
	"os"

	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"

	"time"

	"xorm.io/xorm"
)

type FileTracking struct {
	Id            int64     `json:"-"`
	FolderHash    string    `json:"-"`
	FolderId      string    `json:"folderId"`
	FolderSize    int64     `json:"folderSize"`
	FileCount     int       `json:"fileCount"`
	UploadDate    time.Time `json:"uploadDate"`
	IsEncrypted   bool      `json:"isEncrypted"`
	Password      string    `json:"-"`
	DownloadCount int64     `json:"downloadCount"`
	DownloadLimit int64     `json:"downloadLimit"`
	ExpireTime    time.Time `json:"expireTime"`
	IsDeleted     bool      `json:"-"`
}

var Engine *xorm.Engine

func CreateDBEngine() error {
	var err error
	if os.Getenv("DB_TYPE") == "sqlite" {
		Engine, err = xorm.NewEngine("sqlite", "tmp/data.db")
	} else {
		connectionInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
		Engine, err = xorm.NewEngine("postgres", connectionInfo)
	}

	if err != nil {
		return err
	}
	if err := Engine.Ping(); err != nil {
		return err
	}
	if err := Engine.Sync(new(FileTracking)); err != nil {
		return err
	}

	return nil
}
