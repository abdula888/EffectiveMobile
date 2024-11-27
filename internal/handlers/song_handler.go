package handlers

import (
	"EffectiveMobile/internal/config"
	"EffectiveMobile/internal/models"
	"EffectiveMobile/internal/repository"
	"EffectiveMobile/pkg/api"
	"EffectiveMobile/pkg/log"
	"database/sql"
	"encoding/json"
	"html/template"
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

	// Логируем параметры фильтра
	log.Logger.Debugf("Filter parameters: group=%s, song=%s, releaseDate=%s, page=%s", group, song, releaseDate, page)

	// Указываем количество песен на одной странице
	songsPerPage := 20
	offset := (pageNumber - 1) * songsPerPage

	db, err := config.InitDB()
	if err != nil {
		log.Logger.Error("Failed to connect to the database:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	log.Logger.Debug("Successfully connected to the database")
	defer db.Close()

	// Получаем песни с учётом лимита и смещения
	songs, err := repository.GetSongsWithPagination(db, songsPerPage, offset, group, song, releaseDate)
	if err != nil {
		log.Logger.Error("Error fetching songs from database:", err)
		http.Error(w, "Error loading songs", http.StatusInternalServerError)
		return
	}
	log.Logger.Infof("Fetched %d songs for page %d", len(songs), pageNumber)

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
		log.Logger.Error("Error rendering template:", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
	log.Logger.Info("Template rendered successfully")
}

// RenderSongText отображает полный текст песни
func RenderSongText(w http.ResponseWriter, r *http.Request) {
	// Извлекаем часть пути после "/songs/"
	path := r.URL.Path[len("/songs/"):]
	log.Logger.Debugf("Request path: %s", path)

	// Разбиваем строку на группу и песню
	parts := strings.SplitN(path, "+", 2)
	if len(parts) != 2 {
		log.Logger.Warnf("Invalid song URL: %s", path)
		http.Error(w, "Invalid song URL", http.StatusBadRequest)
		return
	}

	groupName := parts[0]
	songName := parts[1]
	log.Logger.Debugf("Group: %s, Song: %s", groupName, songName)

	// Подключаемся к базе данных
	db, err := config.InitDB()
	if err != nil {
		log.Logger.Error("Failed to connect to the database:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	log.Logger.Debug("Successfully connected to the database")
	defer db.Close()

	// Получаем песню через репозиторий
	song, err := repository.GetSongByName(db, groupName, songName)
	if err != nil {
		log.Logger.Warnf("Song not found: Group=%s, Song=%s", groupName, songName)
		http.Error(w, "Song not found", http.StatusNotFound)
		return
	}
	log.Logger.Infof("Song retrieved: %s - %s", groupName, songName)

	// Разделяем текст песни на куплеты
	verses := strings.Split(song.Text, "\n\n")
	log.Logger.Debugf("Total verses: %d", len(verses))

	// Получаем параметр текущего куплета из URL
	verseParam := r.URL.Query().Get("verse")
	if verseParam == "" {
		verseParam = "1"
	}
	log.Logger.Debugf("Verse parameter: %s", verseParam)

	verseParam = strings.Trim(verseParam, "/")
	verseNumber, err := strconv.Atoi(verseParam)
	if err != nil || verseNumber < 1 || verseNumber > len(verses) {
		log.Logger.Warnf("Invalid verse number: %s", verseParam)
		http.Error(w, "Invalid verse number", http.StatusBadRequest)
		return
	}
	log.Logger.Infof("Displaying verse %d for song %s - %s", verseNumber, groupName, songName)

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
	tmpl, err := template.ParseFiles("internal/templates/song_text.html")
	if err != nil {
		log.Logger.Error("Error loading template:", err)
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}
	log.Logger.Debug("Template loaded successfully")

	// Отображаем шаблон
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Logger.Error("Error rendering template:", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
	log.Logger.Infof("Template rendered successfully for song %s - %s, verse %d", groupName, songName, verseNumber)
}

func UpdateSong(w http.ResponseWriter, r *http.Request) {
	var song models.Song
	err := json.NewDecoder(r.Body).Decode(&song)
	if err != nil {
		log.Logger.Warn("Error decoding JSON: ", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	log.Logger.Debugf("Decoded song: %+v", song)

	db, err := config.InitDB()
	if err != nil {
		log.Logger.Error("Failed to connect to the database:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	log.Logger.Debug("Successfully connected to the database")
	defer db.Close()

	err = repository.UpdateSongByName(db, song)
	if err != nil {
		log.Logger.Error("Error updating song in database:", err)
		http.Error(w, "Error updating song", http.StatusInternalServerError)
		return
	}
	log.Logger.Infof("Song updated successfully: %+v", song)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Song updated successfully!"})
}

// Основная функция добавления песни
func AddSong(w http.ResponseWriter, r *http.Request) {
	var song models.Song
	err := json.NewDecoder(r.Body).Decode(&song)
	if err != nil {
		log.Logger.Warn("Error decoding JSON: ", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	log.Logger.Debugf("Decoded song: %+v", song)

	if song.GroupName == "" || song.Song == "" {
		log.Logger.Warn("GroupName or Song fields are empty")
		http.Error(w, "Group and Song fields cannot be empty", http.StatusBadRequest)
		return
	}

	audDData, err := api.GetAudDData(song.GroupName, song.Song)
	if err != nil {
		log.Logger.Error("Error fetching song data from AudD: ", err)
		http.Error(w, "Error fetching song data", http.StatusInternalServerError)
		return
	}
	log.Logger.Debugf("AudD data retrieved: %+v", audDData)

	lastFmData, err := api.GetLastFmData(song.GroupName, song.Song)
	if err != nil {
		log.Logger.Error("Error fetching song data from Last.fm: ", err)
		http.Error(w, "Error fetching song data", http.StatusInternalServerError)
		return
	}
	log.Logger.Debugf("Last.fm data retrieved: %+v", lastFmData)

	if lastFmData.Track.Wiki.Published != "" {
		song.ReleaseDate, err = parseReleaseDate(lastFmData.Track.Wiki.Published)
		if err != nil {
			log.Logger.Error("Error parsing release date: ", err)
			http.Error(w, "Error parsing release date", http.StatusInternalServerError)
			return
		}
		log.Logger.Debugf("Parsed release date: %s", song.ReleaseDate)
	} else {
		log.Logger.Info("No release date found in Last.fm API response")
		song.ReleaseDate = ""
	}

	song.Text = audDData.Result[0].Lyrics
	log.Logger.Debugf("Fetched song lyrics: %s", song.Text)

	var media []Media
	err = json.Unmarshal([]byte(audDData.Result[0].Media), &media)
	if err != nil {
		log.Logger.Error("Error parsing media field: ", err)
		http.Error(w, "Error parsing media field", http.StatusInternalServerError)
		return
	}
	log.Logger.Debugf("Parsed media field: %+v", media)

	for _, m := range media {
		if m.Provider == "youtube" {
			song.Link = m.URL
			break
		}
	}
	if song.Link == "" {
		log.Logger.Warn("YouTube link not found")
	} else {
		log.Logger.Debugf("YouTube link found: %s", song.Link)
	}

	db, err := config.InitDB()
	if err != nil {
		log.Logger.Error("Failed to connect to the database:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	log.Logger.Debug("Successfully connected to the database")
	defer db.Close()

	var groupID int
	queryCheckGroup := `SELECT group_id FROM groups WHERE group_name = $1`
	err = db.QueryRow(queryCheckGroup, song.GroupName).Scan(&groupID)
	if err == sql.ErrNoRows {
		queryInsertGroup := `INSERT INTO groups (group_name) VALUES ($1) RETURNING group_id`
		err = db.QueryRow(queryInsertGroup, song.GroupName).Scan(&groupID)
		if err != nil {
			log.Logger.Error("Error inserting group into database: ", err)
			http.Error(w, "Error saving song", http.StatusInternalServerError)
			return
		}
		log.Logger.Infof("Group inserted into database: %s (ID: %d)", song.GroupName, groupID)
	} else if err != nil {
		log.Logger.Error("Error checking group in database: ", err)
		http.Error(w, "Error saving song", http.StatusInternalServerError)
		return
	} else {
		log.Logger.Debugf("Group exists in database: %s (ID: %d)", song.GroupName, groupID)
	}

	queryInsertSong := `INSERT INTO songs (group_id, song_name, text, releaseDate, link) VALUES ($1, $2, $3, $4, $5)`
	_, err = db.Exec(queryInsertSong, groupID, song.Song, song.Text, song.ReleaseDate, song.Link)
	if err != nil {
		log.Logger.Error("Error inserting song into database: ", err)
		http.Error(w, "Error saving song", http.StatusInternalServerError)
		return
	}
	log.Logger.Infof("Song inserted into database: %s - %s", song.GroupName, song.Song)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Response{Message: "Song added successfully"})
}

func DeleteSong(w http.ResponseWriter, r *http.Request) {
	var song models.Song
	err := json.NewDecoder(r.Body).Decode(&song)
	if err != nil {
		log.Logger.Warn("Error decoding JSON: ", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	log.Logger.Debugf("Decoded song for deletion: %+v", song)

	if song.GroupName == "" || song.Song == "" {
		log.Logger.Warn("GroupName or Song fields are empty")
		http.Error(w, "Group and Song fields cannot be empty", http.StatusBadRequest)
		return
	}

	db, err := config.InitDB()
	if err != nil {
		log.Logger.Error("Failed to connect to the database:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	log.Logger.Debug("Successfully connected to the database")
	defer db.Close()

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
		log.Logger.Error("Error checking if song exists: ", err)
		http.Error(w, "Error checking song", http.StatusInternalServerError)
		return
	}
	log.Logger.Debugf("Song existence check: Group=%s, Song=%s, Exists=%t", song.GroupName, song.Song, exists)

	if !exists {
		log.Logger.Warnf("Song not found for deletion: Group=%s, Song=%s", song.GroupName, song.Song)
		http.Error(w, "Song not found", http.StatusNotFound)
		return
	}

	query = `
	DELETE 
	FROM songs 
	WHERE group_id = (SELECT group_id FROM groups WHERE group_name = $1 LIMIT 1)
	  AND song_name=$2
	`
	_, err = db.Exec(query, song.GroupName, song.Song)
	if err != nil {
		log.Logger.Error("Error executing delete query: ", err)
		http.Error(w, "Error deleting song", http.StatusInternalServerError)
		return
	}
	log.Logger.Infof("Song deleted successfully: Group=%s, Song=%s", song.GroupName, song.Song)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Song deleted successfully"})
}
