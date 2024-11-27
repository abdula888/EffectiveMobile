package config

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq" // Подключение к Postgres
)

// Функция для инициализации подключения к базе данных
func InitDB() (*sql.DB, error) {
	databaseURL := os.Getenv("DATABASE_URL") // URL базы данных из .env
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	// Проверка подключения
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
