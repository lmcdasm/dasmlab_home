package main

import (
	"os"

	"github.com/Depado/ginprom"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"surfing-service/api"
	"surfing-service/logutil"
)

const version = "0.1.0"

var (
	componentName = "surfing-main"
	log           = logutil.InitLogger(componentName)
)

func main() {
	log.Infof("DASMLAB Surfing Service starting v%s", version)
	gin.SetMode(gin.ReleaseMode)

	if err := api.Initialize(); err != nil {
		log.Fatalf("Storage initialization failed: %v", err)
	}

	mainRouter := gin.Default()
	mainRouter.Use(cors.Default())

	metricsRouter := gin.Default()
	p := ginprom.New(
		ginprom.Engine(metricsRouter),
		ginprom.Subsystem("gin"),
		ginprom.Path("/metrics"),
	)
	mainRouter.Use(p.Instrument())

	initializeRoutes(mainRouter)

	go func() {
		log.Infof("Starting metrics server on :9222")
		if err := metricsRouter.Run(":9222"); err != nil {
			log.Fatalf("Metrics server error: %v", err)
		}
	}()

	port := os.Getenv("SURFING_API_PORT")
	if port == "" {
		port = "10023"
	}
	log.Infof("Starting API server on :%s", port)
	if err := mainRouter.Run(":" + port); err != nil {
		log.Fatalf("Main server error: %v", err)
	}
}
