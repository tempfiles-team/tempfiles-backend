package database

import (
	// _ "github.com/lib/pq"

	_ "modernc.org/sqlite"
	"xorm.io/xorm"
)

type FileRow struct {
	Id       int64
	FileName string
	FileType string
	FileSize int64
	Password string `json:"-"`
	Encrypto bool
}

var Engine *xorm.Engine

func CreateDBEngine() (*xorm.Engine, error) {
	// connectionInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", "localhost", 5432, "postgres", "qwer1234", "localAuth")
	// engine, err := xorm.NewEngine("postgres", connectionInfo)
	engine, err := xorm.NewEngine("sqlite", "./data.db")

	if err != nil {
		return nil, err
	}
	if err := engine.Ping(); err != nil {
		return nil, err
	}
	if err := engine.Sync(new(FileRow)); err != nil {
		return nil, err
	}

	return engine, nil
}
