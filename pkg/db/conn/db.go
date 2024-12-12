package conn

import (
	_ "github.com/lib/pq" // Подключение к Postgres
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Функция для инициализации подключения к базе данных
func InitDB(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
