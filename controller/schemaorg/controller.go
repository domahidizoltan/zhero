// Package schemaorg defines the Schema.org handlers
package schemaorg

import (
	"net/http"
	"strings"

	"github.com/aymerick/raymond"
	"github.com/domahidizoltan/zhero/controller/template"
	"github.com/domahidizoltan/zhero/domain/schemametadata"
	"github.com/domahidizoltan/zhero/domain/schemaorg"
	"github.com/domahidizoltan/zhero/pkg/handlebars"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

const (
	tpl                    = template.TemplatesPath + "schemaorg/"
	tplSearch              = tpl + "search.hbs"
	tplEdit                = tpl + "edit.hbs"
	tplEditPropertyPartial = tpl + "edit-property.partial.hbs"
)

var (
	searchTpl        *raymond.Template
	editTpl          *raymond.Template
	editPropertyPart *raymond.Template
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

	editPropertyPart, err = raymond.ParseFile(tplEditPropertyPartial)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse template")
	}
	editTpl.RegisterPartialTemplate("editProperty", editPropertyPart)

	handlebars.InitHelpers()
}

type Controller struct {
	schemaorgSvc      schemaorg.Service
	schemametadataSvc schemametadata.Service
	classHierarchy    [][]string
}

func NewController(schemaorgSvc schemaorg.Service, schemametadataSvc schemametadata.Service) Controller {
	return Controller{
		schemaorgSvc:      schemaorgSvc,
		schemametadataSvc: schemametadataSvc,
	}
}

func (sc *Controller) Search(c *gin.Context) {
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

func (sc *Controller) Edit(c *gin.Context) {
	clsName := c.Param("class")
	if clsName == "" {
		c.String(http.StatusBadRequest, "missing class")
		return
	}

	cls := sc.schemaorgSvc.GetSchemaClassByName(clsName)
	sc.initClassHierarchy()
	var breadcrumbs []string
	for _, ch := range sc.classHierarchy {
		if ch[len(ch)-1] == clsName {
			breadcrumbs = append(breadcrumbs, ch...)
			break
		}
	}

	ctx := map[string]any{
		"class":       cls,
		"breadcrumbs": breadcrumbs,
		"components":  []string{"TODO"},
	}
	body, err := editTpl.Exec(ctx)
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

func (sc *Controller) initClassHierarchy() {
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
}

func (sc *Controller) GetClassHierarchy(c *gin.Context) {
	sc.initClassHierarchy()
	c.JSON(http.StatusOK, sc.classHierarchy)
}

// Schema represents a schema definition to be saved
type Schema struct {
	Name                string
	Identifier          string
	SecondaryIdentifier string
	Properties          []Property
}

// Property represents a schema property to be saved
type Property struct {
	Name       string
	Mandatory  bool
	Searchable bool
	Type       string
	Component  string
	Order      int
}

func (sc *Controller) Save(c *gin.Context) {
	if err := c.Request.ParseForm(); err != nil {
		c.String(http.StatusBadRequest, "invalid form")
		return
	}

	clsName := c.Param("class")
	if clsName == "" {
		c.String(http.StatusBadRequest, "missing class name")
		return
	}
	schemaToSave := Schema{
		Name:                clsName,
		Identifier:          c.PostForm("identifier"),
		SecondaryIdentifier: c.PostForm("secondary-identifier"),
	}

	orderedProperties := strings.Split(c.PostForm("property-order"), ",")
	for i, propName := range orderedProperties {
		if propName == "" {
			continue
		}

		if c.PostForm(propName+"-hide") == "true" {
			continue
		}

		prop := Property{
			Name:       propName,
			Mandatory:  c.PostForm(propName+"-mandatory") == "true",
			Searchable: c.PostForm(propName+"-searchable") == "true",
			Type:       c.PostForm(propName + "-type"),
			Component:  c.PostForm(propName + "-component"),
			Order:      i,
		}
		schemaToSave.Properties = append(schemaToSave.Properties, prop)
	}

	repoSchema := schemametadata.Schema{
		Name:                schemaToSave.Name,
		Identifier:          schemaToSave.Identifier,
		SecondaryIdentifier: schemaToSave.SecondaryIdentifier,
	}
	for _, p := range schemaToSave.Properties {
		repoSchema.Properties = append(repoSchema.Properties, schemametadata.Property{
			Name:       p.Name,
			Mandatory:  p.Mandatory,
			Searchable: p.Searchable,
			Type:       p.Type,
			Component:  p.Component,
			Order:      p.Order,
		})
	}

	saveErr := sc.schemametadataSvc.Save(c.Request.Context(), repoSchema)

	if saveErr == nil {
		c.Redirect(http.StatusSeeOther, "/")
	} else {
		log.Err(saveErr).Msg("failed to save schema")
		cls := sc.schemaorgSvc.GetSchemaClassByName(clsName)
		sc.initClassHierarchy()
		var breadcrumbs []string
		for _, ch := range sc.classHierarchy {
			if ch[len(ch)-1] == clsName {
				breadcrumbs = append(breadcrumbs, ch...)
				break
			}
		}
		ctx := map[string]any{
			"class":       cls,
			"breadcrumbs": breadcrumbs,
			"components":  []string{"TODO"},
		}
		body, err := editTpl.Exec(ctx)
		if err != nil {
			c.String(http.StatusInternalServerError, "error rendering template")
		}

		output, err := template.Index(c, template.Content{
			Title:    "Welcome to Zhero",
			Body:     raymond.SafeString(body),
			ErrorMsg: "Failed to save schema: " + saveErr.Error(),
		})
		if err != nil {
			c.String(http.StatusInternalServerError, "error rendering template")
			return
		}

		c.Data(http.StatusOK, gin.MIMEHTML, []byte(output))
	}
}
