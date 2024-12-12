package routes

import (
	"EffectiveMobile/internal/config"
	"EffectiveMobile/internal/handlers"
	"html/template"
	"net/http"

	_ "EffectiveMobile/api/swagger"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoutes(tmplAddSong, tmplSongs, tmplDeleteSong *template.Template, conf *config.Config) *gin.Engine {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/songs/add_song/", func(c *gin.Context) {
		err := tmplAddSong.Execute(c.Writer, nil) // Отображаем HTML-страницу
		if err != nil {
			http.Error(c.Writer, "Error rendering template", http.StatusInternalServerError)
			return
		}
	})

	r.POST("/songs/add_song/", func(c *gin.Context) {
		handlers.AddSongHandler(c, conf)
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
