package migrations

import (
	"EffectiveMobile/internal/models"
	"EffectiveMobile/pkg/log"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	err := db.AutoMigrate(&models.Group{}, &models.Song{})
	if err != nil {
		return err
	}

	log.Logger.Info("Migrations applied successfully!")
	return nil
}
