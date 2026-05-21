package dynamicpage

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"

	"github.com/domahidizoltan/zhero/controller"
	"github.com/domahidizoltan/zhero/controller/template"
	"github.com/domahidizoltan/zhero/domain/page"
	"github.com/domahidizoltan/zhero/domain/schema"
	"github.com/domahidizoltan/zhero/pkg/collection"
	"github.com/domahidizoltan/zhero/pkg/paging"
	"github.com/domahidizoltan/zhero/pkg/url"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
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
		ctrl.transformShortRef(c.Request.Context(), p.ListableData, p.References, true)
		return map[string]any{
			id:                   p.Identifier,
			secID:                p.SecondaryIdentifier,
			"listableProperties": p.ListableData,
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

	dataFn := func(schema.SchemaMeta) map[string]any {
		ctrl.transformShortRef(c.Request.Context(), page.Data, page.References, false)
		return page.Data
	}

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

	data := dataFn(*schemaMeta)
	if err := ctrl.transformReferences(c.Request.Context(), data, false); err != nil {
		controller.InternalServerError(c, "failed to transform references", err)
		return
	}

	body, err := ctrl.dynamicPageRdr.Render(*schemaMeta, data)
	if err != nil {
		controller.InternalServerError(c, "failed to generate page", err)
		return
	}

	template.WithLayout(c, pageMeta, body)
}

func (ctrl *Controller) transformReferences(ctx context.Context, data map[string]any, disableAllLinks bool) error {
	dataRefs := map[string]string{}
	refs := map[string]string{}
	for k, v := range data {
		if k == "references" {
			continue
		}

		for _, ref := range controller.RefPattern.FindAllStringSubmatch(fmt.Sprintf("%s", v), -1) {
			dataRefs[ref[0]] = ref[1]
			refs[ref[1]] = ref[2]
		}
	}

	refSlice := slices.Collect(maps.Keys(refs))
	results, err := ctrl.pageSvc.ListEnabledRoutesByRef(ctx, refSlice)
	if err != nil {
		return err
	}

	for k, v := range dataRefs {
		text := v
		var metaObj map[string]string
		if err := json.Unmarshal([]byte(fmt.Sprintf("{%s}", refs[v])), &metaObj); err != nil {
			return err
		}
		if lt, found := metaObj["linkText"]; found && lt != "" {
			text = lt
		}

		pg, found := results[v]
		if disableAllLinks || !found || !pg.IsEnabled {
			dataRefs[k] = text
		} else {
			route := v
			if pg.Route != "" {
				route = pg.Route
			}
			dataRefs[k] = fmt.Sprintf("<a href=\"%s\" alt=\"%s\">%s</a>", route, metaObj["altText"], metaObj["linkText"])
		}
	}

	data["references"] = dataRefs

	return nil
}

func (ctrl *Controller) transformShortRef(ctx context.Context, data map[string]any, refs []string, disableAllLinks bool) {
	for k, v := range data {
		if vStr, ok := v.(string); ok {
			for _, match := range controller.ShortRefPattern.FindAllStringSubmatch(vStr, -1) {
				idx, err := strconv.Atoi(match[1])
				if err != nil {
					log.Err(err).Strs("match", match).Msg("failed to parse shortref index")
					continue
				}

				vStr = strings.ReplaceAll(vStr, match[0], refs[idx])
			}
			v = vStr
		}
		data[k] = v
	}
	ctrl.transformReferences(ctx, data, disableAllLinks)
}
