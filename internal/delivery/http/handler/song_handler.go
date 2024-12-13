package handler

import (
	"EffectiveMobile/internal/config"
	"EffectiveMobile/internal/infrastructure/postgres/model"
	"EffectiveMobile/internal/infrastructure/postgres/repository"
	"EffectiveMobile/pkg/api/audd"
	"EffectiveMobile/pkg/api/lastfm"
	"EffectiveMobile/pkg/db/conn"
	"EffectiveMobile/pkg/log"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// @Summary RenderSongsListHandler
// @Tags songs
// @Description display list of songs
// @Produce      json
// @Success      200  {object}  []model.Song
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router /songs/ [get]
func RenderSongsListHandler(c *gin.Context) {
	r, w := c.Request, c.Writer

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

	db, err := conn.InitDB("postgres://test_user:password@localhost:5432/test_db?sslmode=disable")
	if err != nil {
		log.Logger.Error("Failed to connect to the database:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	log.Logger.Debug("Successfully connected to the database")

	// Получаем песни с учётом лимита и смещения
	songs, err := repository.GetSongsList(db, songsPerPage, offset, group, song, releaseDate)
	if err != nil {
		log.Logger.Error("Error fetching songs from database:", err)
		http.Error(w, "Error loading songs", http.StatusInternalServerError)
		return
	}
	log.Logger.Infof("Fetched %d songs for page %d", len(songs), pageNumber)

	songsListJSON := make([]DataJSON, len(songs))
	for i, elem := range songs {
		songsListJSON[i] = DataJSON{GroupName: elem.GroupName, SongName: elem.SongName, ReleaseDate: elem.ReleaseDate, Link: elem.Link}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(songsListJSON)
}

// @Summary RenderSongTextHandler
// @Tags song
// @Description display song's info
// @Produce      json
// @Success      200  {object}  model.Song
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router /groups/:groupName/songs/:songName [get]
func RenderSongTextHandler(c *gin.Context) {
	r, w := c.Request, c.Writer
	groupName := c.Param("groupName")
	songName := c.Param("songName")
	log.Logger.Debugf("Group: %s, Song: %s", groupName, songName)

	// Подключаемся к базе данных
	db, err := conn.InitDB("postgres://test_user:password@localhost:5432/test_db?sslmode=disable")
	if err != nil {
		log.Logger.Error("Failed to connect to the database:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	log.Logger.Debug("Successfully connected to the database")

	// Получаем песню через репозиторий
	song, err := repository.GetSongText(db, groupName, songName)
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

	// Подготовка данных для шаблона
	data := struct {
		GroupName   string
		SongName    string
		VerseNumber int
		Verse       string
		ReleaseDate time.Time
		Link        string
		FullText    string
	}{
		GroupName:   song.GroupName,
		SongName:    song.SongName,
		VerseNumber: verseNumber,
		Verse:       currentVerse,
		ReleaseDate: song.ReleaseDate,
		Link:        song.Link,
		FullText:    song.Text,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// @Summary UpdateSongHandler
// @Tags songs
// @Description update song
// @Accept      json
// @Success      200  {object}  model.Song
// @Failure      400  {object}  error
// @Failure      500  {object}  error
// @Router /songs/ [put]
func UpdateSongHandler(c *gin.Context) {
	r, w := c.Request, c.Writer
	var songJSON SongJSON
	err := json.NewDecoder(r.Body).Decode(&songJSON)
	if err != nil {
		log.Logger.Warn("Error decoding JSON: ", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	log.Logger.Debugf("Decoded song: %+v", songJSON)

	song := model.Song{GroupName: songJSON.GroupName, SongName: songJSON.SongName,
		Text: songJSON.Text, ReleaseDate: songJSON.ReleaseDate, Link: songJSON.Link}
	db, err := conn.InitDB("postgres://test_user:password@localhost:5432/test_db?sslmode=disable")
	if err != nil {
		log.Logger.Error("Failed to connect to the database:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	log.Logger.Debug("Successfully connected to the database")

	err = repository.UpdateSong(db, song)
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

// @Summary AddSongHandler
// @Tags add_song
// @Description add song
// @Accept      json
// @Success      201  {object}  model.Song
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object} error
// @Router /songs/add_song/ [post]
func AddSongHandler(c *gin.Context, conf *config.Config) {
	r, w := c.Request, c.Writer
	var songJSON SongJSON
	err := json.NewDecoder(r.Body).Decode(&songJSON)
	if err != nil {
		log.Logger.Warn("Error decoding JSON: ", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	log.Logger.Debugf("Decoded song: %+v", songJSON)

	if songJSON.GroupName == "" || songJSON.SongName == "" {
		log.Logger.Warn("GroupName or Song fields are empty")
		http.Error(w, "Group and Song fields cannot be empty", http.StatusBadRequest)
		return
	}
	song := model.Song{GroupName: songJSON.GroupName, SongName: songJSON.SongName}
	audDData, err := audd.GetAudDData(song.GroupName, song.SongName, conf.AuddAPI.AuddAPIKey, conf.AuddAPI.AuddAPIURL)
	if err != nil {
		log.Logger.Error("Error fetching song data from AudD: ", err)
		http.Error(w, "Error fetching song data", http.StatusInternalServerError)
		return
	}
	log.Logger.Debugf("AudD data retrieved")

	lastFmData, err := lastfm.GetLastFmData(song.GroupName, song.SongName, conf.LastFMAPI.LastFMAPIKey, conf.LastFMAPI.LastFMAPIURL)
	if err != nil {
		log.Logger.Error("Error fetching song data from Last.fm: ", err)
		http.Error(w, "Error fetching song data", http.StatusInternalServerError)
		return
	}
	log.Logger.Debugf("Last.fm data retrieved")

	if lastFmData.Track.Wiki.Published != "" {
		song.ReleaseDate, err = time.Parse("02 Jan 2006, 15:04", lastFmData.Track.Wiki.Published)
		if err != nil {
			log.Logger.Error("Error parsing release date: ", err)
			http.Error(w, "Error parsing release date", http.StatusInternalServerError)
			return
		}
		log.Logger.Debugf("Parsed release date: %s", song.ReleaseDate)
	} else {
		log.Logger.Info("No release date found in Last.fm API response")
		song.ReleaseDate = time.Time{}
	}

	song.Text = audDData.Result[0].Lyrics
	log.Logger.Debugf("Fetched song lyrics")

	var media []audd.Media
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

	db, err := conn.InitDB("postgres://test_user:password@localhost:5432/test_db?sslmode=disable")
	if err != nil {
		log.Logger.Error("Failed to connect to the database:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	log.Logger.Debug("Successfully connected to the database")

	err = repository.AddSong(db, song)
	if err != nil {
		log.Logger.Warnf("Error adding song: %s", err)
		http.Error(w, "Error adding song", http.StatusNotFound)
		return
	}

	log.Logger.Infof("Song inserted into database: %s - %s", song.GroupName, song.SongName)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Song added successfully"})
}

// @Summary DeleteSongHandler
// @Tags delete_song
// @Description delete song
// @Accept      json
// @Success      200  {object}  model.Song
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router /songs/delete_song/ [delete]
func DeleteSongHandler(c *gin.Context) {
	r, w := c.Request, c.Writer
	var songJSON SongJSON
	err := json.NewDecoder(r.Body).Decode(&songJSON)
	if err != nil {
		log.Logger.Warn("Error decoding JSON: ", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	log.Logger.Debugf("Decoded song: %+v", songJSON)

	if songJSON.GroupName == "" || songJSON.SongName == "" {
		log.Logger.Warn("GroupName or Song fields are empty")
		http.Error(w, "Group and Song fields cannot be empty", http.StatusBadRequest)
		return
	}
	song := model.Song{GroupName: songJSON.GroupName, SongName: songJSON.SongName}

	db, err := conn.InitDB("postgres://test_user:password@localhost:5432/test_db?sslmode=disable")
	if err != nil {
		log.Logger.Error("Failed to connect to the database:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	log.Logger.Debug("Successfully connected to the database")

	err = repository.DeleteSong(db, song.SongName, song.GroupName)
	if err != nil {
		log.Logger.Warnf("Error delete song: %s", err)
		http.Error(w, "Error delete song", http.StatusNotFound)
		return
	}

	log.Logger.Infof("Song deleted successfully: Group=%s, Song=%s", song.GroupName, song.SongName)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Song deleted successfully"})
}
