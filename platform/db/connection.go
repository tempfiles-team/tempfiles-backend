package db

import (
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/tempfiles-Team/tempfiles-backend/app/models"
	_ "modernc.org/sqlite"

	"xorm.io/xorm"
)

func SQLiteConnection() (*xorm.Engine, error) {
	db, err := xorm.NewEngine("sqlite", "tmp/data.db")
	if err != nil {
		return nil, fmt.Errorf("error, not connected to database, %w", err)
	}

	if err := db.Ping(); err != nil {
		defer db.Close()
		return nil, fmt.Errorf("error, not sent ping to database, %w", err)
	}

	return db, nil
}

func PostgreSQLConnection() (*xorm.Engine, error) {
	// Define database connection for PostgreSQL.
	connectionInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	db, err := xorm.NewEngine("postgres", connectionInfo)
	if err != nil {
		return nil, fmt.Errorf("error, not connected to database, %w", err)
	}

	// Try to ping database.
	if err := db.Ping(); err != nil {
		defer db.Close()
		return nil, fmt.Errorf("error, not sent ping to database, %w", err)
	}
	if err := db.Sync(new(models.FileTracking)); err != nil {
		defer db.Close()
		return nil, fmt.Errorf("error, not sent ping to database, %w", err)
	}
	if err := db.Sync(new(models.TextTracking)); err != nil {
		defer db.Close()
		return nil, fmt.Errorf("error, not sent ping to database, %w", err)
	}

	return db, nil
}
