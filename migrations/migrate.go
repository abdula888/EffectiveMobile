package migrations

import (
	"database/sql"
	"io/ioutil"
	"log"
)

func ApplyMigration(db *sql.DB, migrationPath string) error {
	var migrationExists bool
	err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'songs')").Scan(&migrationExists)
	if err != nil {
		return err
	}

	// Проверяем, была ли уже применена миграция
	if !migrationExists {
		log.Println("Migration not found, applying migration...")

		// Чтение файла миграции
		migrationSQL, err := ioutil.ReadFile(migrationPath)
		if err != nil {
			return err
		}

		// Выполнение миграции
		_, err = db.Exec(string(migrationSQL))
		if err != nil {
			return err
		}

		log.Println("Migration applied successfully!")
	} else {
		log.Println("Migration already applied, skipping...")
	}

	return nil
}
