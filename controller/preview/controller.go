// Package preview contains the controllers for the page previews
package preview

import (
	"github.com/domahidizoltan/zhero/controller/adminpage"
	"github.com/domahidizoltan/zhero/controller/dynamicpage"
	schema_domain "github.com/domahidizoltan/zhero/domain/schema"
	"github.com/domahidizoltan/zhero/pkg/url"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	dynamicPageCtrl dynamicpage.Controller
}

func NewController(
	dynamicPageCtrl dynamicpage.Controller,
) Controller {
	return Controller{
		dynamicPageCtrl,
	}
}

func (ctrl *Controller) LoadPage(c *gin.Context) {
	ctrl.dynamicPageCtrl.LoadPage(c, false)
}

func (ctrl *Controller) InFlightPage(c *gin.Context) {
	class := c.Param("class")

	dto := adminpage.PageDtoFrom(nil)
	dto.EnhanceFromForm(c)

	dataFn := func(meta schema_domain.SchemaMeta) map[string]any {
		data := dto.ToMap()
		for _, prop := range meta.Properties {
			data[prop.Name] = c.PostForm("field-" + prop.Name)
		}
		return data
	}

	pageMeta := dto.Meta.ToMap()
	pageMeta["canonicalURL"] = url.Canonical(c.Request)
	ctrl.dynamicPageCtrl.Render(c, class, pageMeta, dataFn)
}
