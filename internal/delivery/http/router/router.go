package router

import (
	"EffectiveMobile/internal/config"
	"EffectiveMobile/internal/delivery/http/handler"

	_ "EffectiveMobile/api/swagger"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(conf *config.Config) *gin.Engine {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/groups/:groupName/songs/:songName", func(c *gin.Context) {
		handler.RenderSongTextHandler(c)
	})

	// Filter - "songs/?group=Eminem&song=&releaseDate=&page="
	r.GET("/songs/", func(c *gin.Context) {
		handler.RenderSongsListHandler(c)
	})

	r.POST("/songs/", func(c *gin.Context) {
		handler.AddSongHandler(c, conf)
	})

	r.PUT("/songs/", func(c *gin.Context) {
		handler.UpdateSongHandler(c)
	})

	r.DELETE("/songs/", func(c *gin.Context) {
		handler.DeleteSongHandler(c)
	})

	return r
}
