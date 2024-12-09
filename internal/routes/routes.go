package routes

import (
	"EffectiveMobile/internal/handlers"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(tmplAddSong, tmplSongs, tmplDeleteSong *template.Template) *gin.Engine {
	r := gin.Default()

	r.GET("/songs/add_song/", func(c *gin.Context) {
		err := tmplAddSong.Execute(c.Writer, nil) // Отображаем HTML-страницу
		if err != nil {
			http.Error(c.Writer, "Error rendering template", http.StatusInternalServerError)
			return
		}
	})

	r.POST("/songs/add_song/", func(c *gin.Context) {
		handlers.AddSongHandler(c)
	})

	r.GET("/songs/delete_song/", func(c *gin.Context) {
		err := tmplDeleteSong.Execute(c.Writer, nil) // Отображаем HTML-страницу
		if err != nil {
			http.Error(c.Writer, "Error rendering template", http.StatusInternalServerError)
			return
		}
	})

	r.DELETE("/songs/delete_song/", func(c *gin.Context) {
		handlers.DeleteSongHandler(c)
	})

	r.GET("/groups/:groupName/songs/:songName", func(c *gin.Context) {
		handlers.RenderSongTextHandler(c)
	})

	r.GET("/songs/", func(c *gin.Context) {
		handlers.RenderSongsListHandler(c, tmplSongs)
	})

	r.PUT("/songs/", func(c *gin.Context) {
		handlers.UpdateSongHandler(c)
	})
	return r
}
