package db

import (
	"github.com/tempfiles-Team/tempfiles-backend/app/models"
)

func DropAllTables() error {

	if err := Engine.Migrator().DropTable(new(models.FileTracking), new(models.TextTracking)); err != nil {
		return err
	}
	return nil

}

func MigrateAllTables() error {

	if err := Engine.AutoMigrate(new(models.FileTracking), new(models.TextTracking)); err != nil {
		return err
	}

	return nil
}
