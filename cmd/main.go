package main

import (
	"EffectiveMobile/config"
	"EffectiveMobile/migrations"
	"EffectiveMobile/routes"
	"html/template"
	"log"
	"net/http"
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
		log.Fatal("Error loading .env file")
	}

	// Инициализация базы данных
	db, err := config.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	defer db.Close()

	// Применение миграций
	if err := migrations.RunMigrations(db); err != nil {
		log.Fatal("Error applying migration: ", err)
	}

	// Парсинг шаблона HTML для добавления песни
	tmplAddSong, err := template.ParseFiles(filepath.Join("templates", "add_song.html"))
	if err != nil {
		log.Fatal("Error parsing add_song template: ", err)
	}

	// Парсинг шаблона HTML для отображения списка песен
	tmplSongs := template.Must(template.New("songs.html").Funcs(templateFuncs).ParseFiles("templates/songs.html"))
	if err != nil {
		log.Fatal("Error parsing songs template: ", err)
	}

	// Парсинг шаблона HTML для удаления песни
	tmplDeleteSong, err := template.ParseFiles(filepath.Join("templates", "delete_song.html"))
	if err != nil {
		log.Fatal("Error parsing add_song template: ", err)
	}

	// Регистрация маршрутов
	routes.RegisterRoutes(tmplAddSong, tmplSongs, tmplDeleteSong)

	// Запуск сервера
	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
