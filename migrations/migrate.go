package migrations

import (
	"EffectiveMobile/internal/infrastructure/postgres/model"
	"EffectiveMobile/pkg/log"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	err := db.AutoMigrate(&model.Group{}, &model.Song{})
	if err != nil {
		return err
	}

	log.Logger.Info("Migrations applied successfully!")
	return nil
}
