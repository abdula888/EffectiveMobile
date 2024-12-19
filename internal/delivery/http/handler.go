package http

import (
	"EffectiveMobile/internal/domain/entity"
	"EffectiveMobile/pkg/log"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type usecase interface {
	GetSongs(filter entity.SongsFilter) ([]entity.SongsList, error)
	GetSongText(groupName, songName, verse string) (entity.SongText, error)
	AddSong(groupName, songName string) error
	UpdateSong(song entity.Song) error
	DeleteSong(groupName, songName string) error
}

type Handler struct {
	usecase usecase
}

func NewRouter(usecase usecase) *gin.Engine {
	h := Handler{
		usecase,
	}

	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/groups/:groupName/songs/:songName", h.GetSongText)

	// Filter - "songs?group=Eminem&song=&releaseDate=&page="
	r.GET("/songs", h.GetSongs)

	r.POST("/songs", h.AddSong)

	r.PUT("/songs", h.UpdateSong)

	r.DELETE("/songs", h.DeleteSong)

	return r
}

func getSongsFilter(c *gin.Context) entity.SongsFilter {
	// Получаем параметры фильтра
	group, _ := c.GetQuery("group")
	song, _ := c.GetQuery("song")
	releaseDate, _ := c.GetQuery("releaseDate")
	page, _ := c.GetQuery("page")

	return entity.SongsFilter{Group: group, Song: song, ReleaseDate: releaseDate, Page: page}
}

// @Summary GetSongs
// @Tags songs
// @Description display list of songs
// @Produce      json
// @Success      200  {object}  []model.Song
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router /songs/ [get]
func (h Handler) GetSongs(c *gin.Context) {
	w := c.Writer

	filter := getSongsFilter(c)

	songsList, err := h.usecase.GetSongs(filter)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var songsListJSON []SongsListJSON
	for _, song := range songsList {
		songJSON := SongsListJSON{GroupName: song.GroupName, SongName: song.SongName, ReleaseDate: song.ReleaseDate, Link: song.Link}
		songsListJSON = append(songsListJSON, songJSON)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(songsListJSON)
}

// @Summary GetSongText
// @Tags song
// @Description display song's info
// @Produce      json
// @Success      200  {object}  model.Song
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router /groups/:groupName/songs/:songName [get]
func (h Handler) GetSongText(c *gin.Context) {
	w := c.Writer
	groupName := c.Param("groupName")
	songName := c.Param("songName")
	verse, _ := c.GetQuery("verse")

	song, err := h.usecase.GetSongText(groupName, songName, verse)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	songTextJSON := SongTextJSON(song)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(songTextJSON)
}

// @Summary AddSong
// @Tags add_song
// @Description add song
// @Accept      json
// @Success      201  {object}  model.Song
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object} error
// @Router /songs [post]
func (h Handler) AddSong(c *gin.Context) {
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

	err = h.usecase.AddSong(songJSON.GroupName, songJSON.SongName)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Song added successfully"})
}

// @Summary UpdateSong
// @Tags songs
// @Description update song
// @Accept      json
// @Success      200  {object}  model.Song
// @Failure      400  {object}  error
// @Failure      500  {object}  error
// @Router /songs [put]
func (h Handler) UpdateSong(c *gin.Context) {
	r, w := c.Request, c.Writer
	var songJSON SongJSON
	err := json.NewDecoder(r.Body).Decode(&songJSON)
	if err != nil {
		log.Logger.Warn("Error decoding JSON: ", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	log.Logger.Debugf("Decoded song: %+v", songJSON)

	song := entity.Song(songJSON)
	err = h.usecase.UpdateSong(song)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Song updated successfully!"})
}

// @Summary DeleteSong
// @Tags delete_song
// @Description delete song
// @Accept      json
// @Success      200  {object}  model.Song
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router /songs [delete]
func (h Handler) DeleteSong(c *gin.Context) {
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

	err = h.usecase.DeleteSong(songJSON.GroupName, songJSON.SongName)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Song deleted successfully"})
}
