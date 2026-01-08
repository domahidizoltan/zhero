// Package preview contains the controllers for the page previews
package preview

import (
	"net/http"

	"github.com/domahidizoltan/zhero/controller"
	page "github.com/domahidizoltan/zhero/controller/page"
	schema_domain "github.com/domahidizoltan/zhero/domain/schema"
	"github.com/domahidizoltan/zhero/pkg/jsonld"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	schemaSvc      schema_domain.Service
	dynamicPageRdr controller.UserFacingPageRenderer
}

func NewController(
	schemaSvc schema_domain.Service,
	dynamicPageRdr controller.UserFacingPageRenderer,
) Controller {
	return Controller{
		schemaSvc:      schemaSvc,
		dynamicPageRdr: dynamicPageRdr,
	}
}

func (ctrl *Controller) Page(c *gin.Context) {
	class := c.Query("class")
	// id := c.PostForm("identifier")

	meta, err := ctrl.schemaSvc.GetSchemaMetaByName(c, class)
	if err != nil {
		controller.InternalServerError(c, "failed to get schema data", err)
		return
	}

	dto := page.PageDtoFrom(meta)
	dto.EnhanceFromForm(c)
	json, err := jsonld.FromPage(dto.ToModel())
	if err != nil {
		controller.InternalServerError(c, "failed to generate JSON-LD", err)
		return
	}
	// c.Data(http.StatusOK, "application/ld+json", json)

	content, err := ctrl.dynamicPageRdr.Render(string(json))
	if err != nil {
		controller.InternalServerError(c, "failed to generate page", err)
	}
	c.Data(http.StatusOK, "text/html", []byte(content))
}
