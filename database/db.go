package database

import (
	// _ "github.com/lib/pq"

	"time"

	_ "modernc.org/sqlite"
	"xorm.io/xorm"
)

type FileTracking struct {
	Id          int64     `json:"-"`
	FileId      string    `json:"fileId"`
	FileName    string    `json:"fileName"`
	FileSize    int64     `json:"size"`
	UploadDate  time.Time `json:"uploadDate"`
	IsEncrypted bool      `json:"isEncrypted"`
	Password    string    `json:"-"`
}

var Engine *xorm.Engine

func CreateDBEngine() (*xorm.Engine, error) {
	// connectionInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", "localhost", 5432, "postgres", "qwer1234", "localAuth")
	// engine, err := xorm.NewEngine("postgres", connectionInfo)
	engine, err := xorm.NewEngine("sqlite", "./data.db")
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
