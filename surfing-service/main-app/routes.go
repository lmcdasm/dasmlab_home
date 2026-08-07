package main

import (
	"surfing-service/api"

	"github.com/gin-gonic/gin"
)

func initializeRoutes(router *gin.Engine) {
	router.GET("/isalive", api.IsAlive)

	router.GET("/auth/config", api.AuthConfig)
	router.GET("/auth/login", api.AuthLogin)
	router.GET("/auth/callback", api.AuthCallback)
	router.GET("/auth/logout", api.AuthLogout)
	router.GET("/auth/me", api.AuthMe)

	router.POST("/activity", api.RequireAuth(), api.PostActivity)
	router.GET("/activity", api.RequireActivityViewer(), api.ListActivity)

	router.GET("/days", api.ListDays)
	router.POST("/days", api.RequireOwner(), api.CreateDay)
	router.PATCH("/days/:id", api.RequireOwner(), api.PatchDay)
	router.DELETE("/days/:id", api.RequireOwner(), api.DeleteDay)
	router.POST("/days/:id/media", api.RequireOwner(), api.UploadMedia)
	router.POST("/days/:id/media/presign", api.RequireOwner(), api.PresignMediaUpload)
	router.POST("/days/:id/media/:mediaId/complete", api.RequireOwner(), api.CompleteMediaUpload)
	router.POST("/days/:id/media/link", api.RequireOwner(), api.AddMediaLink)
	router.PATCH("/days/:id/media/:mediaId", api.RequireOwner(), api.UpdateMedia)
	router.DELETE("/days/:id/media/:mediaId", api.RequireOwner(), api.DeleteMedia)
	router.POST("/days/:id/media/:mediaId/unhide", api.RequireOwner(), api.UnhideMedia)
	router.GET("/days/:id/media/:mediaId/download", api.DownloadMedia)
	router.POST("/days/:id/media/:mediaId/play", api.RecordMediaPlay)
	router.POST("/days/:id/media/:mediaId/tags", api.ProposeMediaTag)
	router.POST("/days/:id/media/:mediaId/tags/:tagId/:action", api.RequireOwner(), api.ModerateMediaTag)
	router.POST("/days/:id/media/:mediaId/transcode", api.RequireOwner(), api.TranscodeMedia)
	router.POST("/days/:id/publish", api.RequireOwner(), api.PublishDay)
	router.POST("/days/:id/curate/publish", api.RequireOwner(), api.CuratePublish)
	router.POST("/days/:id/ai/curate", api.RequireOwner(), api.AICurate)
	router.POST("/days/:id/theme/generate", api.RequireOwner(), api.GenerateTheme)
	router.POST("/shares", api.CreateShare)
	router.GET("/shares/:token", api.ShareMeta)
	router.GET("/s/:token", api.ResolveShare)
	router.GET("/serve", api.ServeMedia)
}
