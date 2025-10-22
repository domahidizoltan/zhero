// Package page contains the controllers for the pages
package page

import (
	"net/http"

	"github.com/aymerick/raymond"
	"github.com/domahidizoltan/zhero/controller"
	"github.com/domahidizoltan/zhero/controller/template"
	"github.com/domahidizoltan/zhero/domain/schema"
	tpl "github.com/domahidizoltan/zhero/template"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	schemaSvc schema.Service
}

func NewController(schemaSvc schema.Service) Controller {
	return Controller{
		schemaSvc: schemaSvc,
	}
}

func (pc *Controller) Main(c *gin.Context) {
	schemas, err := pc.schemaSvc.GetSchemaMetaNames(c)
	if err != nil {
		controller.InternalServerError(c, "failed to get schemas", err)
		return
	}

	ctx := map[string]any{
		"schemas": schemas,
	}
	body, err := tpl.PageMain.Exec(ctx)
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

func (pc *Controller) List(c *gin.Context) {
	clsName := c.Param("class")
	ctx := map[string]any{
		"class": clsName,
	}
	output, err := tpl.PageList.Exec(ctx)
	if err != nil {
		controller.TemplateRenderError(c, err)
		return
	}

	c.Data(http.StatusOK, gin.MIMEHTML, []byte(output))
}
