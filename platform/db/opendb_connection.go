package db

import (
	"fmt"
	"os"

	"github.com/tempfiles-Team/tempfiles-backend/app/models"
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
	// Try to ping database.
	if err := db.Ping(); err != nil {
		defer db.Close()
		return fmt.Errorf("error, not sent ping to database, %w", err)
	}
	if err := db.Sync(new(models.FileTracking)); err != nil {
		defer db.Close()
		return fmt.Errorf("error, sync to database, %w", err)
	}
	if err := db.Sync(new(models.TextTracking)); err != nil {
		defer db.Close()
		return fmt.Errorf("error, sync to database, %w", err)
	}

	Engine = db

	return nil
}
