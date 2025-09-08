package main

import (
	"fmt"
	"os"
	"time"

	"github.com/domahidizoltan/zhero/config"
	"github.com/domahidizoltan/zhero/controller"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func zerologMiddleware(logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)

		logger.Info().
			Str("client_ip", c.ClientIP()).
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Int("status", c.Writer.Status()).
			Int("size", c.Writer.Size()).
			Dur("duration", duration).
			Msg("incoming request")
	}
}

func main() {
	gin.SetMode(gin.ReleaseMode)

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	log := zerolog.New(output).With().Timestamp().Logger()
	// log := zerolog.New(os.Stdout).With().Timestamp().Logger() //for JSON logging

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	router := gin.New()
	router.Use(zerologMiddleware(log))
	router.Use(gin.Recovery())

	controller.SetRoutes(router)

	serverAddr := fmt.Sprintf(":%d", cfg.Admin.Server.Port)
	log.Info().Int("port", cfg.Admin.Server.Port).Msg("Server starting on port")
	if err := router.Run(serverAddr); err != nil {
		log.Fatal().Err(err).Msg("failed to start server")
	}
}
