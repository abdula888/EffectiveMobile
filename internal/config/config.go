package config

import (
	_ "github.com/lib/pq" // Подключение к Postgres
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Функция для инициализации подключения к базе данных
func InitDB() (*gorm.DB, error) {
	databaseURL := "host=localhost user=test_user dbname=test_db password=password sslmode=disable" // для go run
	//databaseURL := os.Getenv("DATABASE_URL") // URL базы данных из .env для Docker
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
