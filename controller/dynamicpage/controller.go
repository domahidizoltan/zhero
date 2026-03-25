package dynamicpage

import (
	"slices"

	"github.com/domahidizoltan/zhero/controller"
	"github.com/domahidizoltan/zhero/controller/template"
	"github.com/domahidizoltan/zhero/domain/page"
	"github.com/domahidizoltan/zhero/domain/schema"
	"github.com/domahidizoltan/zhero/pkg/collection"
	"github.com/domahidizoltan/zhero/pkg/paging"
	"github.com/domahidizoltan/zhero/pkg/url"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	dynamicPageRdr controller.UserFacingPageListRenderer
	schemaSvc      schema.Service
	pageSvc        page.Service
}

func NewController(pageRenderer controller.UserFacingPageListRenderer, schemaSvc schema.Service, pageSvc page.Service) Controller {
	return Controller{
		dynamicPageRdr: pageRenderer,
		schemaSvc:      schemaSvc,
		pageSvc:        pageSvc,
	}
}

func (ctrl *Controller) List(c *gin.Context) {
	clsName := c.Param("class")

	pageOpts := paging.RequestToPageOpts(c, "identifier")
	opts := page.ListOptions{
		PageOpts: pageOpts,
	}

	pages, paging, err := ctrl.pageSvc.List(c, clsName, opts, true)
	if err != nil {
		controller.InternalServerError(c, "failed to list pages", err)
		return
	}

	meta, err := ctrl.schemaSvc.GetSchemaMetaByName(c, clsName)
	if err != nil {
		controller.InternalServerError(c, "failed to get schema data", err)
		return
	}

	data := slices.Collect(collection.MapValues(pages, func(p page.Page) map[string]any {
		id, secID := meta.Identifier, meta.SecondaryIdentifier
		return map[string]any{
			id:    p.Identifier,
			secID: p.SecondaryIdentifier,
		}
	}))

	if meta == nil { // why is this called second time with nil?
		return
	}

	content, err := ctrl.dynamicPageRdr.List(*meta, data, paging)
	if err != nil {
		controller.TemplateRenderError(c, err)
		return
	}

	listMeta := map[string]any{} // TODO list page meta
	listMeta["canonicalURL"] = url.Canonical(c.Request)
	template.WithLayout(c, listMeta, content)
}

func (ctrl *Controller) Page(c *gin.Context) {
	ctrl.LoadPage(c, true)
}

func (ctrl *Controller) LoadPage(c *gin.Context, onlyEnabled bool) {
	if c.Param("skipLoadPage") != "" {
		return
	}
	class := c.Param("class")
	identifier := c.Param("identifier")

	page, err := ctrl.pageSvc.GetPageBySchemaNameAndIdentifier(c, class, identifier, onlyEnabled)
	if err != nil {
		controller.InternalServerError(c, "failed to load page", err)
		return
	}

	if page == nil {
		template.PageNotFoundLayout(c)
		return
	}

	dataFn := func(schema.SchemaMeta) map[string]any { return page.Data }

	if page.Meta.Title == "" {
		page.Meta.Title = page.SecondaryIdentifier
	}
	if page.Meta.OGTitle == "" {
		page.Meta.OGTitle = page.SecondaryIdentifier
	}
	pageMeta := page.Meta.ToMap()
	pageMeta["canonicalURL"] = url.Canonical(c.Request)

	ctrl.Render(c, class, pageMeta, dataFn)
}

func (ctrl *Controller) Render(c *gin.Context, class string, pageMeta map[string]any, dataFn func(schema.SchemaMeta) map[string]any) {
	schemaMeta, err := ctrl.schemaSvc.GetSchemaMetaByName(c, class)
	if err != nil {
		controller.InternalServerError(c, "failed to get schema data", err)
		return
	}

	body, err := ctrl.dynamicPageRdr.Render(*schemaMeta, dataFn(*schemaMeta))
	if err != nil {
		controller.InternalServerError(c, "failed to generate page", err)
		return
	}

	template.WithLayout(c, pageMeta, body)
}
