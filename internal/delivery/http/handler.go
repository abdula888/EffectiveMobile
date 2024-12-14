package http

import (
	"EffectiveMobile/internal/domain/entity"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type usecase interface {
	GetSongs(filter entity.SongsFilter) ([]entity.SongsList, error)
	GetSongText()
	AddSong()
	UpdateSong()
	DeleteSong()
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

	// Filter - "songs/?group=Eminem&song=&releaseDate=&page="
	r.GET("/songs/", h.GetSongs)

	//r.POST("/songs/", h.AddSong)

	r.PUT("/songs/", h.UpdateSong)

	r.DELETE("/songs/", h.DeleteSong)

	return r
}

func getSongsFilter(url url.URL) entity.SongsFilter {
	// Получаем параметры фильтра
	group := url.Query().Get("group")
	song := url.Query().Get("song")
	releaseDate := url.Query().Get("releaseDate")
	page := url.Query().Get("page")

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
	r, w := c.Request, c.Writer

	filter := getSongsFilter(*r.URL)

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
