// Package template is to define common templates
package template

import (
	"github.com/aymerick/raymond"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Content struct {
	Title    string
	Style    string
	Script   string
	Body     raymond.SafeString
	ErrorMsg string
}

const (
	TemplatesPath = "templates/"
	tplIndex      = TemplatesPath + "index.hbs"
)

var indexTpl *raymond.Template

func init() {
	var err error
	indexTpl, err = raymond.ParseFile(tplIndex)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse template")
	}
}

func Index(c *gin.Context, content Content) (string, error) {
	output, err := indexTpl.Exec(content)
	if err != nil {
		log.Error().Err(err).Msg("error rendering template")
		return "", err
	}
	return output, nil
}
