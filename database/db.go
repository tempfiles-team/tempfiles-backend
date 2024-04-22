package database

import (
	"fmt"
	"os"

	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"

	"time"

	"xorm.io/xorm"
)

type FileListResponse struct {
	FileName    string `json:"fileName"`
	FileSize    int64  `json:"fileSize"`
	DownloadUrl string `json:"downloadUrl"`
}

type FileTracking struct {
	Id         int64  `json:"-"`
	FolderHash string `json:"-"`
	IsDeleted  bool   `json:"-"`

	IsHidden      bool      `json:"isHidden"`
	FolderId      string    `json:"folderId"`
	FileCount     int       `json:"fileCount"`
	DownloadCount int64     `json:"downloadCount"`
	DownloadLimit int64     `json:"downloadLimit"`
	UploadDate    time.Time `json:"uploadDate"`
	ExpireTime    int64     `json:"expireTime"`

	Files []FileListResponse `json:"files"`
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

	Engine.SetMaxIdleConns(20)
	Engine.SetMaxOpenConns(20)
	Engine.SetConnMaxLifetime(time.Minute * 5)

	if err := Engine.Ping(); err != nil {
		return err
	}
	if err := Engine.Sync(new(FileTracking)); err != nil {
		return err
	}

	return nil
}
