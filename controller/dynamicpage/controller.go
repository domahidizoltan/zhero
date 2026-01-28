package dynamicpage

import (
	"net/url"
	"slices"
	"strconv"
	"strings"

	"github.com/domahidizoltan/zhero/controller"
	"github.com/domahidizoltan/zhero/controller/template"
	"github.com/domahidizoltan/zhero/domain/page"
	"github.com/domahidizoltan/zhero/domain/schema"
	"github.com/domahidizoltan/zhero/pkg/collection"
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

	pageNo, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || pageNo < 1 {
		pageNo = 1
	}

	opts := page.ListOptions{
		Page: uint(pageNo),
	}

	sortQuery := c.DefaultQuery("sort", "identifier:desc")
	unescape, err := url.QueryUnescape(sortQuery)
	if err != nil {
		controller.BadRequest(c, "failed to parse search query", err)
		return
	}
	sortQuery = unescape

	if sortBy, sortOrder, found := strings.Cut(sortQuery, ":"); found {
		opts.SortBy = sortBy
		opts.SortDir = page.SortDir(sortOrder)
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

	template.WithLayout(c, content)
}

func (ctrl *Controller) Page(c *gin.Context) {
	ctrl.LoadPage(c, true)
}

func (ctrl *Controller) LoadPage(c *gin.Context, onlyEnabled bool) {
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
	ctrl.Render(c, class, dataFn)
}

func (ctrl *Controller) Render(c *gin.Context, class string, dataFn func(schema.SchemaMeta) map[string]any) {
	meta, err := ctrl.schemaSvc.GetSchemaMetaByName(c, class)
	if err != nil {
		controller.InternalServerError(c, "failed to get schema data", err)
		return
	}

	content, err := ctrl.dynamicPageRdr.Render(*meta, dataFn(*meta))
	if err != nil {
		controller.InternalServerError(c, "failed to generate page", err)
		return
	}
	template.WithLayout(c, content)
}
