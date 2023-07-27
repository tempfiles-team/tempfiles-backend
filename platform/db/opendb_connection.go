package db

import (
	"os"

	"xorm.io/xorm"
)

var Engine *xorm.Engine

// OpenDBConnection func for opening database connection.
func OpenDBConnection() error {

	var db *xorm.Engine
	var err error

	if os.Getenv("DB_TYPE") == "sqlite" {
		db, err = SQLiteConnection()
	} else {
		db, err = PostgreSQLConnection()
	}

	if err != nil {
		return err
	}

	Engine = db

	return nil
}
