package main

import (
	"EffectiveMobile/internal/config"
	"EffectiveMobile/internal/migrations"
	"EffectiveMobile/internal/routes"
	"EffectiveMobile/pkg/log"
	"html/template"
	"path/filepath"

	"github.com/joho/godotenv" // для загрузки конфигурации из .env
)

// Функции для шаблонов
var templateFuncs = template.FuncMap{
	"add": func(a, b int) int {
		return a + b
	},
	"minus": func(a, b int) int {
		return b - a
	},
}

func main() {
	// Загрузка переменных окружения из файла .env
	if err := godotenv.Load(".env"); err != nil {
		log.Logger.Warn("Error loading .env file, using defaults")
	}

	// Инициализация базы данных
	db, err := config.InitDB()
	if err != nil {
		log.Logger.Fatal("Failed to connect to the database:", err)
	}
	log.Logger.Info("Successfully connected to the database")
	defer db.Close()

	// Применение миграций
	if err := migrations.RunMigrations(db); err != nil {
		log.Logger.Fatal("Error applying migration: ", err)
	}

	// Парсинг шаблона HTML для добавления песни
	tmplAddSong, err := template.ParseFiles(filepath.Join("internal/templates", "add_song.html"))
	if err != nil {
		log.Logger.Fatal("Error parsing add_song template: ", err)
	}
	log.Logger.Debug("add_song template parsed successfully")

	// Парсинг шаблона HTML для отображения списка песен
	tmplSongs := template.Must(template.New("songs.html").Funcs(templateFuncs).ParseFiles("internal/templates/songs.html"))
	log.Logger.Debug("songs template parsed successfully")

	// Парсинг шаблона HTML для удаления песни
	tmplDeleteSong, err := template.ParseFiles(filepath.Join("internal/templates", "delete_song.html"))
	if err != nil {
		log.Logger.Fatal("Error parsing delete_song template: ", err)
	}
	log.Logger.Debug("delete_song template parsed successfully")

	// Регистрация маршрутов
	r := routes.RegisterRoutes(tmplAddSong, tmplSongs, tmplDeleteSong)
	log.Logger.Info("Routes registered successfully")

	// Запуск сервера
	log.Logger.Info("Server started on port 8080")
	log.Logger.Fatal(r.Run(":8080"))
}
