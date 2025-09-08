// Package schemaorg defines the Schema.org handlers
package schemaorg

import (
	"net/http"

	"github.com/aymerick/raymond"
	"github.com/domahidizoltan/zhero/controller/template"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

const (
	tpl       = template.TemplatesPath + "schemaorg/"
	tplCreate = tpl + "create.hbs"
)

var createTpl *raymond.Template

func init() {
	var err error
	createTpl, err = raymond.ParseFile(tplCreate)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse template")
	}
}

func Create(c *gin.Context) {
	body, err := createTpl.Exec(nil)
	if err != nil {
		c.String(http.StatusInternalServerError, "error rendering template")
	}

	output, err := template.Index(c, template.Content{
		Title: "Welcome to Zhero",
		Body:  raymond.SafeString(body),
	})
	if err != nil {
		c.String(http.StatusInternalServerError, "error rendering template")
		return
	}

	c.Data(http.StatusOK, gin.MIMEHTML, []byte(output))
}

func Search(c *gin.Context) {
	c.Data(http.StatusOK, gin.MIMEHTML, []byte("<p>testing1</p><br/><p>testing2"))
}
