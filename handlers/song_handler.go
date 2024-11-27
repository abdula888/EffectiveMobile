package handlers

import (
	"EffectiveMobile/api"
	"EffectiveMobile/config"
	"EffectiveMobile/models"
	"EffectiveMobile/repository"
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Функция для парсинга даты в формат Go
func parseReleaseDate(dateStr string) (string, error) {
	// Пример: "02 Jan 2006, 15:04"
	parsedDate, err := time.Parse("02 Jan 2006, 15:04", dateStr)
	if err != nil {
		return "", err
	}
	return parsedDate.Format("2006-01-02"), nil // Возвращаем в формате YYYY-MM-DD
}

type Response struct {
	Message string `json:"message"`
}

// Структура для получения ссылки из ответа AudD API
type Media struct {
	Provider string `json:"provider"`
	Type     string `json:"type"`
	URL      string `json:"url"`
}

func RenderSongsList(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	// Получаем параметры фильтра
	group := r.URL.Query().Get("group")
	song := r.URL.Query().Get("song")
	releaseDate := r.URL.Query().Get("releaseDate")

	// Получаем параметр страницы
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}

	pageNumber, err := strconv.Atoi(page)
	if err != nil || pageNumber < 1 {
		pageNumber = 1
	}

	// Указываем количество песен на одной странице
	songsPerPage := 20
	offset := (pageNumber - 1) * songsPerPage

	db, err := config.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	defer db.Close()

	// Получаем песни с учётом лимита и смещения
	songs, err := repository.GetSongsWithPagination(db, songsPerPage, offset, group, song, releaseDate)
	if err != nil {
		http.Error(w, "Error loading songs", http.StatusInternalServerError)
		return
	}

	// Передаём список песен в шаблон
	data := struct {
		Songs             []models.Song
		CurrentPage       int
		HasPrevPage       bool
		HasNextPage       bool
		FilterGroup       string
		FilterSong        string
		FilterReleaseDate string
	}{
		Songs:             songs,
		CurrentPage:       pageNumber,
		HasPrevPage:       pageNumber > 1,
		HasNextPage:       len(songs) == songsPerPage,
		FilterGroup:       group,
		FilterSong:        song,
		FilterReleaseDate: releaseDate,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

// RenderSongText отображает полный текст песни
func RenderSongText(w http.ResponseWriter, r *http.Request) {
	// Извлекаем часть пути после "/songs/"
	path := r.URL.Path[len("/songs/"):]

	// Разбиваем строку на группу и песню
	parts := strings.SplitN(path, "+", 2)
	if len(parts) != 2 {
		http.Error(w, "Invalid song URL", http.StatusBadRequest)
		return
	}

	groupName := parts[0]
	songName := parts[1]
	// Подключаемся к базе данных
	db, err := config.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	defer db.Close()

	// Получаем песню через репозиторий
	song, err := repository.GetSongByName(db, groupName, songName)
	if err != nil {
		http.Error(w, "Song not found", http.StatusNotFound)
		return
	}

	// Разделяем текст песни на куплеты
	verses := strings.Split(song.Text, "\n\n")

	// Получаем параметр текущего куплета из URL
	verseParam := r.URL.Query().Get("verse")
	if verseParam == "" {
		verseParam = "1" // Если куплет не указан, по умолчанию отображаем первый
	}

	verseParam = strings.Trim(verseParam, "/")
	verseNumber, err := strconv.Atoi(verseParam)
	if err != nil || verseNumber < 1 || verseNumber > len(verses) {
		http.Error(w, "Invalid verse number", http.StatusBadRequest)
		return
	}

	// Текущий куплет
	currentVerse := verses[verseNumber-1]

	// Проверка наличия предыдущего и следующего куплета
	hasPrev := verseNumber > 1
	hasNext := verseNumber < len(verses)

	// Подготовка данных для шаблона
	data := struct {
		GroupName   string
		Song        string
		Text        string
		ReleaseDate string
		Link        string
		Verse       string
		VerseNumber int
		HasPrev     bool
		HasNext     bool
		PrevVerse   int
		NextVerse   int
	}{
		GroupName:   song.GroupName,
		Song:        song.Song,
		Text:        song.Text,
		ReleaseDate: song.ReleaseDate,
		Link:        song.Link,
		Verse:       currentVerse,
		VerseNumber: verseNumber,
		HasPrev:     hasPrev,
		HasNext:     hasNext,
		PrevVerse:   verseNumber - 1,
		NextVerse:   verseNumber + 1,
	}

	// Парсим шаблон для отображения куплета
	tmpl, err := template.ParseFiles("templates/song_text.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	// Отображаем шаблон
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

func UpdateSong(w http.ResponseWriter, r *http.Request) {
	var song models.Song
	err := json.NewDecoder(r.Body).Decode(&song)
	if err != nil {
		log.Println("Error decoding JSON: ", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	db, err := config.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	defer db.Close()

	// Получаем песни с учётом лимита и смещения
	err = repository.UpdateSongByName(db, song)
	if err != nil {
		http.Error(w, "Error loading songs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Song updated successfully!"})

}

// Основная функция добавления песни
func AddSong(w http.ResponseWriter, r *http.Request) {
	var song models.Song
	err := json.NewDecoder(r.Body).Decode(&song)
	if err != nil {
		log.Println("Error decoding JSON: ", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if song.GroupName == "" || song.Song == "" {
		http.Error(w, "Group and Song fields cannot be empty", http.StatusBadRequest)
		return
	}

	// Сначала получаем данные через AudD API
	audDData, err := api.GetAudDData(song.GroupName, song.Song)
	if err != nil {
		log.Println("Error fetching song data from AudD: ", err)
		http.Error(w, "Error fetching song data", http.StatusInternalServerError)
		return
	}

	lastFmData, err := api.GetLastFmData(song.GroupName, song.Song)
	if err != nil {
		log.Println("Error fetching song data from Last.fm: ", err)
		http.Error(w, "Error fetching song data", http.StatusInternalServerError)
		return
	}

	// Парсим дату релиза из Last.fm
	if lastFmData.Track.Wiki.Published != "" {
		song.ReleaseDate, err = parseReleaseDate(lastFmData.Track.Wiki.Published)
		if err != nil {
			log.Println("Error parsing release date: ", err)
			http.Error(w, "Error parsing release date", http.StatusInternalServerError)
			return
		}
	} else {
		log.Println("No release date found in Last.fm API response")
		song.ReleaseDate = ""
	}

	// Получаем текст песни и медиа ссылки из AudD API
	song.Text = audDData.Result[0].Lyrics

	// Парсим поле media как JSON
	var media []Media
	err = json.Unmarshal([]byte(audDData.Result[0].Media), &media)
	if err != nil {
		log.Println("Error parsing media field: ", err)
		http.Error(w, "Error parsing media field", http.StatusInternalServerError)
		return
	}

	// Ищем ссылку на YouTube
	for _, m := range media {
		if m.Provider == "youtube" {
			song.Link = m.URL
			break
		}
	}

	if song.Link == "" {
		log.Println("YouTube link not found")
	}

	db, err := config.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	defer db.Close()

	var groupID int

	// Проверить, существует ли группа
	queryCheckGroup := `SELECT group_id FROM groups WHERE group_name = $1`
	err = db.QueryRow(queryCheckGroup, song.GroupName).Scan(&groupID)
	if err == sql.ErrNoRows {
		// Группа не найдена, добавляем её
		queryInsertGroup := `INSERT INTO groups (group_name) VALUES ($1) RETURNING group_id`
		err = db.QueryRow(queryInsertGroup, song.GroupName).Scan(&groupID)
		if err != nil {
			log.Println("Error inserting song into database: ", err)
			http.Error(w, "Error saving song", http.StatusInternalServerError)
			return
		}
	} else if err != nil {
		log.Println("Error inserting song into database: ", err)
		http.Error(w, "Error saving song", http.StatusInternalServerError)
		return
	}

	// Добавить песню
	queryInsertSong := `INSERT INTO songs (group_id, song_name, text, releaseDate, link) VALUES ($1, $2, $3, $4, $5)`
	_, err = db.Exec(queryInsertSong, groupID, song.Song, song.Text, song.ReleaseDate, song.Link)
	if err != nil {
		log.Println("Error inserting song into database: ", err)
		http.Error(w, "Error saving song", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Response{Message: "Song added successfully"})
}

func DeleteSong(w http.ResponseWriter, r *http.Request) {
	var song models.Song
	err := json.NewDecoder(r.Body).Decode(&song)
	if err != nil {
		log.Println("Error decoding JSON: ", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Проверка на пустые поля
	if song.GroupName == "" || song.Song == "" {
		http.Error(w, "Group and Song fields cannot be empty", http.StatusBadRequest)
		return
	}

	db, err := config.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	defer db.Close()

	// Проверка, существует ли песня
	var exists bool
	query := `
    SELECT EXISTS (
        SELECT 1 
        FROM songs 
        WHERE group_id = (SELECT group_id FROM groups WHERE group_name = $1) 
          AND song_name = $2
    )
`
	err = db.QueryRow(query, song.GroupName, song.Song).Scan(&exists)
	if err != nil {
		log.Println("Error checking if song exists: ", err)
		http.Error(w, "Error checking song", http.StatusInternalServerError)
		return
	}

	if !exists {
		http.Error(w, "Song not found", http.StatusNotFound)
		return
	}

	// Удаление песни
	query = `
	DELETE 
	FROM songs 
	WHERE group_id = (SELECT group_id FROM groups WHERE group_name = $1 LIMIT 1)
	  AND song=$2
	  `
	_, err = db.Exec(query, song.GroupName, song.Song)
	if err != nil {
		log.Println("Error executing delete query: ", err)
		http.Error(w, "Error deleting song", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Song deleted successfully"})
}
