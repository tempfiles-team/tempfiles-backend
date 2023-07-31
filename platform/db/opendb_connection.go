package db

import (
	"fmt"
	"log"
	"os"

	"github.com/tempfiles-Team/tempfiles-backend/app/models"
	"gorm.io/gorm"
)

var Engine *gorm.DB

// OpenDBConnection func for opening database connection.
func OpenDBConnection() error {

	var db *gorm.DB
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
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("error, not sent ping to database, %w", err)
	}
	err = sqlDB.Ping()
	if err != nil {
		sqlDB.Close()
		return fmt.Errorf("error, not sent ping to database, %w", err)
	}

	err = db.AutoMigrate(&models.FileTracking{}, &models.TextTracking{})
	if err != nil {
		sqlDB.Close()
		return fmt.Errorf("error, auto migration failed, %w", err)
	}

	Engine = db

	log.Println("Database connection successful. ⚡")

	return nil
}
