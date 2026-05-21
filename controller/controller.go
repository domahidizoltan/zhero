// Package controller collects common controller level functions
package controller

import (
	"net/http"
	"regexp"

	"github.com/domahidizoltan/zhero/domain/schema"
	"github.com/domahidizoltan/zhero/pkg/paging"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

var (
	RefPattern      = regexp.MustCompile(`#ZHERO#([^#]+)#\{([^}]*)\}#`)
	ShortRefPattern = regexp.MustCompile(`#ref(\d+)`)
)

type UserFacingPageRenderer interface {
	Render(schemaMeta schema.SchemaMeta, data map[string]any) (string, error)
}

type UserFacingPageListRenderer interface {
	UserFacingPageRenderer
	List(schemaMeta schema.SchemaMeta, data []map[string]any, paging paging.Meta) (string, error)
}

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
