package routes

import (
	"EffectiveMobile/handlers"
	"html/template"
	"net/http"
)

func RegisterRoutes(tmplAddSong, tmplSongs, tmplDeleteSong *template.Template) {
	// Обработчик для отображения страницы добавления песни
	http.HandleFunc("/songs/add_song/", func(w http.ResponseWriter, r *http.Request) {
		err := tmplAddSong.Execute(w, nil) // Отображаем HTML-страницу
		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			return
		}
	})

	// Обработчик для отображения страницы удаления песни
	http.HandleFunc("/songs/delete_song/", func(w http.ResponseWriter, r *http.Request) {
		err := tmplDeleteSong.Execute(w, nil) // Отображаем HTML-страницу
		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			return
		}
	})

	// Обработчик для отображения списка песен в HTML-формате
	http.HandleFunc("/songs/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			path := r.URL.Path
			if path == "/songs/" {
				// Получаем песни и передаем их в шаблон
				handlers.RenderSongsList(w, r, tmplSongs)
			} else {
				// Здесь мы обрабатываем путь вида /songs/{group_name}+{song_name}/
				handlers.RenderSongText(w, r)
			}
		case "POST":
			handlers.AddSong(w, r)
		case "DELETE":
			handlers.DeleteSong(w, r)
		case "PUT":
			handlers.UpdateSong(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
