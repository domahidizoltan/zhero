// Package controller collects common controller level functions
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func TemplateRenderError(c *gin.Context, err error) {
	InternalServerError(c, "failed to render template", err)
}

func BadRequest(c *gin.Context, msg string, err error) {
	log.Error().
		Err(err).
		Str("status", "BadRequest").
		Msg(msg)
	c.String(http.StatusBadRequest, msg)
}

func InternalServerError(c *gin.Context, msg string, err error) {
	log.Error().
		Err(err).
		Str("status", "InternalServerError").
		Msg(msg)
	c.String(http.StatusInternalServerError, msg)
}
