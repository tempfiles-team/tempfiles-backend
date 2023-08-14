package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Engine *gorm.DB

func NewConnection() error {

	var err error

	if os.Getenv("DB_TYPE") == "sqlite" {
		Engine, err = SQLiteConnection()
	} else if os.Getenv("DB_TYPE") == "postgres" {
		Engine, err = PostgreSQLConnection()
	}
	if err != nil {
		return err
	}

	sqlDB, err := Engine.DB()
	if err != nil {
		return fmt.Errorf("error, not connected to database, %w", err)
	}
	if err := sqlDB.Ping(); err != nil {
		sqlDB.Close()
		return fmt.Errorf("error, not connected to database, %w", err)
	}

	if err := MigrateAllTables(); err != nil {
		return fmt.Errorf("error, migration failed, %w", err)
	}

	log.Println("Database connection successful. ⚡")

	return nil
}

func SQLiteConnection() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("tmp/data.db"), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error, not connected to database, %w", err)
	}

	return db, nil
}

func PostgreSQLConnection() (*gorm.DB, error) {
	connectionInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	db, err := gorm.Open(postgres.Open(connectionInfo), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error, not connected to database, %w", err)
	}

	return db, nil
}
