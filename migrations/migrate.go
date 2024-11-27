package migrations

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4" // Migrate library
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // File source for migrations
	_ "github.com/lib/pq"                                // PostgreSQL driver
)

func RunMigrations(db *sql.DB) error {
	// Настраиваем драйвер для работы с PostgreSQL
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	// Создаём мигратор
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations", // Путь к папке с миграциями
		"postgres",          // Имя базы данных
		driver,
	)
	if err != nil {
		return err
	}

	// Применяем миграции
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	log.Println("Migrations applied successfully!")
	return nil
}
