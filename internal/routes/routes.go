package routes

import (
	"EffectiveMobile/internal/config"
	"EffectiveMobile/internal/handlers"

	_ "EffectiveMobile/api/swagger"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoutes(conf *config.Config) *gin.Engine {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Filter - "songs/?group=Eminem&song=&releaseDate=&page="
	r.GET("/songs/", func(c *gin.Context) {
		handlers.RenderSongsListHandler(c)
	})

	r.POST("/songs/", func(c *gin.Context) {
		handlers.AddSongHandler(c, conf)
	})

	r.PUT("/songs/", func(c *gin.Context) {
		handlers.UpdateSongHandler(c)
	})

	r.DELETE("/songs/", func(c *gin.Context) {
		handlers.DeleteSongHandler(c)
	})

	r.GET("/groups/:groupName/songs/:songName", func(c *gin.Context) {
		handlers.RenderSongTextHandler(c)
	})

	return r
}
