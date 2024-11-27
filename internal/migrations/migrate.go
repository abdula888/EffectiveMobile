package migrations

import (
	"database/sql"

	"EffectiveMobile/pkg/log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func RunMigrations(db *sql.DB) error {
	// Настраиваем драйвер для работы с PostgreSQL
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	// Создаём мигратор
	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/migrations", // Путь к папке с миграциями
		"postgres",                   // Имя базы данных
		driver,
	)
	if err != nil {
		return err
	}

	// Применяем миграции
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	log.Logger.Info("Migrations applied successfully!")
	return nil
}
