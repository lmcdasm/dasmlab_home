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
	router.DELETE("/days/:id/media/:mediaId", api.DeleteMedia)
	router.GET("/serve", api.ServeMedia)
}
