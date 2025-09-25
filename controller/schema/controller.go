// Package schema defines the handlers for managing schema
package schema

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/aymerick/raymond"
	"github.com/domahidizoltan/zhero/controller"
	"github.com/domahidizoltan/zhero/controller/template"
	"github.com/domahidizoltan/zhero/domain/schema"
	"github.com/domahidizoltan/zhero/pkg/handlebars"
	"github.com/domahidizoltan/zhero/templates"
	"github.com/gin-gonic/gin"
)

var editTpl = templates.SchemaorgEdit

func init() {
	editTpl.RegisterPartialTemplate("editProperty", templates.SchemaorgEditPropertyPartial)
	handlebars.InitHelpers()
}

type Controller struct {
	schemaSvc schema.Service
}

func NewController(schemaSvc schema.Service) Controller {
	return Controller{
		schemaSvc: schemaSvc,
	}
}

func (sc *Controller) Search(c *gin.Context) {
	body, err := templates.SchemaorgSearch.Exec(nil)
	if err != nil {
		controller.TemplateRenderError(c, err)
		return
	}

	output, err := template.Index(c, template.Content{
		Title: "Welcome to Zhero",
		Body:  raymond.SafeString(body),
	})
	if err != nil {
		controller.TemplateRenderError(c, err)
		return
	}

	c.Data(http.StatusOK, gin.MIMEHTML, []byte(output))
}

func (sc *Controller) GetClassHierarchy(c *gin.Context) {
	c.JSON(http.StatusOK, sc.schemaSvc.GetClassHierarchy())
}

func (sc *Controller) Edit(c *gin.Context) {
	clsName := c.Param("class")
	output, hasError := sc.edit(c, clsName, false)
	if hasError {
		return
	}
	c.Data(http.StatusOK, gin.MIMEHTML, []byte(output))
}

func (sc *Controller) Save(c *gin.Context) {
	if err := c.Request.ParseForm(); err != nil {
		controller.BadRequest(c, "invalid form", err)
		return
	}

	clsName := c.Param("class")
	output, hasError := sc.edit(c, clsName, true)
	if hasError {
		c.Data(http.StatusBadRequest, gin.MIMEHTML, []byte(output))
		return
	}

	c.Redirect(http.StatusSeeOther, "/")
}

func (sc *Controller) edit(c *gin.Context, clsName string, hasFormSubmitted bool) (string, bool) {
	if clsName == "" {
		controller.BadRequest(c, "class is missing", nil)
		return "", true
	}

	ctx := map[string]any{
		"class":       sc.schemaSvc.GetSchemaClassByName(clsName),
		"breadcrumbs": sc.classBreadcrumbs(clsName),
		"components":  []string{"TODO"},
	}
	body, err := editTpl.Exec(ctx)
	if err != nil {
		controller.TemplateRenderError(c, err)
		return "", true
	}

	errorMsg, successMsg := "", ""
	if hasFormSubmitted {
		schemaToSave := sc.schemaFromForm(c, clsName)
		if err := sc.schemaSvc.Save(c.Request.Context(), schemaToSave); err != nil {
			errorMsg = err.Error()
		} else {
			successMsg = fmt.Sprintf("Schema %s saved successfully", clsName)
		}
	}

	// TODO display form errors
	// TODO load schema for save or edit
	operation := "Create"
	output, err := template.Index(c, template.Content{
		Title:    operation + " schema: " + clsName,
		Body:     raymond.SafeString(body),
		ErrorMsg: errorMsg,
		FlashMsg: successMsg,
	})
	if err != nil {
		controller.TemplateRenderError(c, err)
		return "", true
	}

	return output, false
}

func (sc *Controller) schemaFromForm(c *gin.Context, clsName string) schema.Schema {
	schemaToSave := schema.Schema{
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

		prop := schema.Property{
			Name:       propName,
			Mandatory:  c.PostForm(propName+"-mandatory") == "true",
			Searchable: c.PostForm(propName+"-searchable") == "true",
			Type:       c.PostForm(propName + "-type"),
			Component:  c.PostForm(propName + "-component"),
			Order:      i,
		}
		schemaToSave.Properties = append(schemaToSave.Properties, prop)
	}

	return schemaToSave
}

func (sc *Controller) classBreadcrumbs(clsName string) []string {
	var breadcrumbs []string
	for _, ch := range sc.schemaSvc.GetClassHierarchy() {
		if ch[len(ch)-1] == clsName {
			breadcrumbs = append(breadcrumbs, ch...)
			break
		}
	}
	return breadcrumbs
}
