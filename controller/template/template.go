// Package template is to define common templates
package template

import (
	"github.com/aymerick/raymond"
	"github.com/domahidizoltan/zhero/pkg/session"
	"github.com/domahidizoltan/zhero/template"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Content struct {
	Style    string
	Script   string
	Title    string
	Body     raymond.SafeString
	ErrorMsg string
	FlashMsg string
}

func Index(c *gin.Context, content Content) (string, error) {
	handleFlash(c, &content)
	output, err := template.Index.Exec(content)
	if err != nil {
		log.Error().Err(err).Msg("error rendering template")
		return "", err
	}
	return output, nil
}

func handleFlash(c *gin.Context, content *Content) {
	if len(content.FlashMsg) != 0 {
		if err := session.SetFlash(c, content.FlashMsg); err != nil {
			log.Error().Err(err).Msg("failed to save flash message")
		}
		content.FlashMsg = ""
	} else {
		var err error
		content.FlashMsg, err = session.GetFlash(c)
		if err != nil {
			log.Error().Err(err).Msg("failed to update flash session")
		}
	}
}
