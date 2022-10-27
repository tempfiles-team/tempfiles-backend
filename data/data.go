package data

import (
	// _ "github.com/lib/pq"
	_ "modernc.org/sqlite"
	"xorm.io/xorm"
)

type User struct {
	Id       int64
	Name     string
	Email    string
	Password string `json:"-"`
}

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
	if err := engine.Sync(new(User)); err != nil {
		return nil, err
	}

	return engine, nil
}
