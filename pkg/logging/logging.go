// Package logging is for log related functions
package logging

import (
	"os"
	"strings"
	"time"

	"github.com/domahidizoltan/zhero/config"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func ConfigureLogging(cfg *config.Config) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	logger := zerolog.New(output)
	lvl := zerolog.InfoLevel

	if cfg != nil {
		switch strings.ToLower(cfg.Log.Format) {
		case "json":
			logger = zerolog.New(os.Stdout)
		}

		switch strings.ToLower(cfg.Log.Level) {
		case "debug":
			lvl = zerolog.DebugLevel
		case "error":
			lvl = zerolog.ErrorLevel
		case "fatal":
			lvl = zerolog.FatalLevel
		case "info":
			lvl = zerolog.InfoLevel
		case "panic":
			lvl = zerolog.PanicLevel
		case "warn":
			lvl = zerolog.WarnLevel
		}
	}

	log.Logger = logger.With().Timestamp().Logger().Level(lvl)
}

func ZerologMiddleware(logger zerolog.Logger) gin.HandlerFunc {
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
