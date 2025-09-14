// Package schemaorg defines the Schema.org handlers
package schemaorg

import (
	"net/http"
	"strings"

	"github.com/aymerick/raymond"
	"github.com/domahidizoltan/zhero/controller/template"
	"github.com/domahidizoltan/zhero/service/schemaorg"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

const (
	tpl       = template.TemplatesPath + "schemaorg/"
	tplSearch = tpl + "search.hbs"
	tplEdit   = tpl + "edit.hbs"
)

var (
	searchTpl *raymond.Template
	editTpl   *raymond.Template
)

func init() {
	var err error
	searchTpl, err = raymond.ParseFile(tplSearch)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse template")
	}
	editTpl, err = raymond.ParseFile(tplEdit)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse template")
	}
}

type SchemaorgCtrl struct {
	schemaorgSvc   schemaorg.Service
	classHierarchy [][]string
}

func New(schemaorgSvc schemaorg.Service) SchemaorgCtrl {
	return SchemaorgCtrl{
		schemaorgSvc: schemaorgSvc,
	}
}

func (sc *SchemaorgCtrl) Search(c *gin.Context) {
	body, err := searchTpl.Exec(nil)
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

func (sc *SchemaorgCtrl) Edit(c *gin.Context) {
	cls := c.Param("class")
	if cls == "" {
		c.String(http.StatusBadRequest, "missing class")
		return
	}

	body, err := editTpl.Exec(nil)
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

func (sc *SchemaorgCtrl) GetClassHierarchy(c *gin.Context) {
	marker := ">"

	if len(sc.classHierarchy) == 0 {
		lines := sc.schemaorgSvc.GetSubClassesHierarchyOf(schemaorg.RootClass, marker, 0)
		parents := []string{lines[0]}
		sc.classHierarchy = append(sc.classHierarchy, []string{lines[0]})

		for _, l := range lines[1:] {
			level := strings.Count(l, marker)
			switch {
			case level == len(parents):
				parents = append(parents, l[level:])
			case level == len(parents)-1:
				parents[len(parents)-1] = l[level:]
			case level < len(parents)-1:
				parents = parents[:level]
				parents = append(parents, l[level:])
			}

			tmp := make([]string, len(parents))
			copy(tmp, parents)
			sc.classHierarchy = append(sc.classHierarchy, tmp)
		}
	}

	c.JSON(http.StatusOK, sc.classHierarchy)
}
