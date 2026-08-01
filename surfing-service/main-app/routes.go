package main

import (
	"surfing-service/api"

	"github.com/gin-gonic/gin"
)

func initializeRoutes(router *gin.Engine) {
	router.GET("/isalive", api.IsAlive)
	router.GET("/days", api.ListDays)
	router.POST("/days", api.CreateDay)
	router.DELETE("/days/:id", api.DeleteDay)
	router.POST("/days/:id/media", api.UploadMedia)
	router.POST("/days/:id/media/link", api.AddMediaLink)
	router.PATCH("/days/:id/media/:mediaId", api.UpdateMedia)
	router.DELETE("/days/:id/media/:mediaId", api.DeleteMedia)
	router.POST("/days/:id/media/:mediaId/play", api.RecordMediaPlay)
	router.POST("/days/:id/media/:mediaId/tags", api.ProposeMediaTag)
	router.POST("/days/:id/media/:mediaId/tags/:tagId/:action", api.ModerateMediaTag)
	router.POST("/days/:id/publish", api.PublishDay)
	router.POST("/days/:id/theme/generate", api.GenerateTheme)
	router.POST("/shares", api.CreateShare)
	router.GET("/shares/:token", api.ShareMeta)
	router.GET("/s/:token", api.ResolveShare)
	router.GET("/serve", api.ServeMedia)
}
