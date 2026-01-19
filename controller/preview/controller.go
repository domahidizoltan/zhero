// Package preview contains the controllers for the page previews
package preview

import (
	"net/http"

	"github.com/domahidizoltan/zhero/controller"
	page_domain "github.com/domahidizoltan/zhero/domain/page"
	schema_domain "github.com/domahidizoltan/zhero/domain/schema"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	schemaSvc      schema_domain.Service
	pageSvc        page_domain.Service
	dynamicPageRdr controller.UserFacingPageRenderer
}

func NewController(
	schemaSvc schema_domain.Service,
	pageSvc page_domain.Service,
	dynamicPageRdr controller.UserFacingPageRenderer,
) Controller {
	return Controller{
		schemaSvc:      schemaSvc,
		pageSvc:        pageSvc,
		dynamicPageRdr: dynamicPageRdr,
	}
}

func (ctrl *Controller) LoadPage(c *gin.Context) {
	class := c.Param("class")
	identifier := c.Param("identifier")

	page, err := ctrl.pageSvc.GetPageBySchemaNameAndIdentifier(c, class, identifier)
	if err != nil {
		controller.InternalServerError(c, "failed to load page", err)
		return
	}

	dataFn := func(schema_domain.SchemaMeta) map[string]any { return page.Data }
	ctrl.render(c, class, dataFn)
}

func (ctrl *Controller) InFlightPage(c *gin.Context) {
	class := c.Param("class")

	dataFn := func(meta schema_domain.SchemaMeta) map[string]any {
		data := map[string]any{}
		for _, prop := range meta.Properties {
			data[prop.Name] = c.PostForm("field-" + prop.Name)
		}
		return data
	}
	ctrl.render(c, class, dataFn)
}

func (ctrl *Controller) render(c *gin.Context, class string, dataFn func(schema_domain.SchemaMeta) map[string]any) {
	meta, err := ctrl.schemaSvc.GetSchemaMetaByName(c, class)
	if err != nil {
		controller.InternalServerError(c, "failed to get schema data", err)
		return
	}

	content, err := ctrl.dynamicPageRdr.Render(*meta, dataFn(*meta))
	if err != nil {
		controller.InternalServerError(c, "failed to generate page", err)
	}
	c.Data(http.StatusOK, "text/html", []byte(content))
}
