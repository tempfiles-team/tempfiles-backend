package database

import (
	"fmt"
	"os"

	_ "github.com/lib/pq"

	"time"

	"xorm.io/xorm"
)

type FileTracking struct {
	Id          int64     `json:"-"`
	FileId      string    `json:"fileId"`
	FileName    string    `json:"filename"`
	FileSize    int64     `json:"size"`
	UploadDate  time.Time `json:"uploadDate"`
	IsEncrypted bool      `json:"isEncrypted"`
	Password    string    `json:"-"`
}

var Engine *xorm.Engine

func CreateDBEngine() (*xorm.Engine, error) {
	connectionInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	engine, err := xorm.NewEngine("postgres", connectionInfo)
	randInit()

	if err != nil {
		return nil, err
	}
	if err := engine.Ping(); err != nil {
		return nil, err
	}
	if err := engine.Sync(new(FileTracking)); err != nil {
		return nil, err
	}

	return engine, nil
}
