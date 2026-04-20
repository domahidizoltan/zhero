// Package adminpage contains the controllers for the pages
package adminpage

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/aymerick/raymond"
	"github.com/domahidizoltan/zhero/controller"
	"github.com/domahidizoltan/zhero/controller/template"
	"github.com/domahidizoltan/zhero/domain/page"
	"github.com/domahidizoltan/zhero/domain/route"
	"github.com/domahidizoltan/zhero/domain/schema"
	"github.com/domahidizoltan/zhero/pkg/paging"
	tpl "github.com/domahidizoltan/zhero/template"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Controller struct {
	schemaSvc schema.Service
	pageSvc   page.Service
	routeSvc  route.Service
}

func NewController(schemaSvc schema.Service, pageSvc page.Service, routeSvc route.Service) Controller {
	return Controller{
		schemaSvc: schemaSvc,
		pageSvc:   pageSvc,
		routeSvc:  routeSvc,
	}
}

func (pc *Controller) Main(c *gin.Context) {
	schemas, err := pc.schemaSvc.GetSchemaMetaNames(c)
	if err != nil {
		controller.InternalServerError(c, "failed to get schemas", err)
		return
	}

	var selectedSchema string
	if len(schemas) > 0 {
		selectedSchema = schemas[0]
	}
	if s, ok := c.GetQuery("schema"); ok {
		selectedSchema = s
	}

	ctx := map[string]any{
		"schemas":        schemas,
		"selectedSchema": selectedSchema,
	}
	body, err := tpl.AdminPageMain.Exec(ctx)
	if err != nil {
		controller.TemplateRenderError(c, err)
		return
	}

	output, err := template.AdminIndex(c, template.Content{
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

	pageOpts := paging.RequestToPageOpts(c, "identifier")
	opts := page.ListOptions{
		SecondaryIdentifierLike: pageOpts.SearchParam(),
		PageOpts:                pageOpts,
	}

	pages, paging, err := pc.pageSvc.List(c, clsName, opts, false)
	if err != nil {
		controller.InternalServerError(c, "failed to list pages", err)
		return
	}

	pagingBaseURL, urlQuery := pageOpts.GetURL("/admin/page/list/" + clsName)

	ctx := map[string]any{
		"class":    clsName,
		"paging":   paging.ToDto(pagingBaseURL, "#page-list-content"),
		"urlQuery": urlQuery,

		"pages":    pages,
		"search":   pageOpts.SearchParam(),
		"sort":     pageOpts.SortQuery(),
		"listOpts": opts,
	}

	output, err := tpl.AdminPageList.Exec(ctx)
	if err != nil {
		controller.TemplateRenderError(c, err)
		return
	}

	c.Data(http.StatusOK, gin.MIMEHTML, []byte(output))
}

func (pc *Controller) Create(c *gin.Context) {
	output, hasError := pc.edit(c, false)
	if hasError {
		c.Data(http.StatusBadRequest, gin.MIMEHTML, []byte(output))
		return
	}
	c.Data(http.StatusOK, gin.MIMEHTML, []byte(output))
}

func (pc *Controller) Edit(c *gin.Context) {
	output, hasError := pc.edit(c, false)
	if hasError {
		c.Data(http.StatusBadRequest, gin.MIMEHTML, []byte(output))
		return
	}
	c.Data(http.StatusOK, gin.MIMEHTML, []byte(output))
}

func (pc *Controller) EditAction(c *gin.Context) {
	identifier := c.PostForm("item-identifier")
	action := c.PostForm("item-action")
	class := c.Param("class")

	sendPopupError := func(status int, msg string, err error) {
		log.Error().Err(err).Str("status", http.StatusText(status)).Msg(msg)
		jsonPayload, _ := json.Marshal(map[string]string{"showError": msg})
		c.Header("HX-Trigger", string(jsonPayload))
		c.Status(status)
	}

	if identifier == "" || action == "" {
		sendPopupError(http.StatusBadRequest, "identifier and action are mandatory", nil)
		return
	}

	var err error
	switch action {
	case "enable":
		err = pc.pageSvc.Enable(c, class, identifier, true)
	case "disable":
		err = pc.pageSvc.Enable(c, class, identifier, false)
	case "delete":
		err = pc.pageSvc.Delete(c, class, identifier)
	default:
		sendPopupError(http.StatusBadRequest, fmt.Sprintf("invalid action: %s", action), nil)
		return
	}

	if err != nil {
		sendPopupError(http.StatusInternalServerError, fmt.Sprintf("failed to perform action '%s'", action), err)
		return
	}

	search := c.Query("search")
	sort := c.Query("sort")
	page := c.Query("page")

	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/admin/page/list/%s?search=%s&sort=%s&page=%s", class, search, sort, page))
}

func (pc *Controller) Save(c *gin.Context) {
	output, hasError := pc.edit(c, true)
	if hasError {
		c.Data(http.StatusBadRequest, gin.MIMEHTML, []byte(output))
		return
	}
	class := c.Param("class")
	c.Redirect(http.StatusSeeOther, "/admin/page/list?schema="+class)
}

func (pc *Controller) edit(c *gin.Context, hasFormSubmitted bool) (string, bool) {
	class := c.Param("class")
	identifier := c.Param("identifier")
	if id, ok := c.GetPostForm("identifier"); ok {
		identifier = id
	}

	var dto pageDto
	meta, err := pc.schemaSvc.GetSchemaMetaByName(c, class)
	if err != nil {
		controller.InternalServerError(c, "failed to get schema data", err)
		return "", true
	}

	dto = PageDtoFrom(meta)
	errorMsg, successMsg := "", ""
	if hasFormSubmitted {
		dto.EnhanceFromForm(c)
		dto.extractReferences()
		page := dto.ToModel()

		var err error
		if len(identifier) == 0 {
			identifier, err = pc.pageSvc.Create(c, page, dto.Identifier)
		} else {
			err = pc.pageSvc.Update(c, identifier, page, dto.Identifier)
		}

		if err != nil {
			log.Error().Err(err).Msg("failed to save page")
			errorMsg = err.Error()
		} else {
			successMsg = fmt.Sprintf("\"%s\" page saved successfully with ID %s", class, identifier)
		}

	}

	pageModel, err := pc.pageSvc.GetPageBySchemaNameAndIdentifier(c.Request.Context(), class, identifier, false)
	if err != nil {
		controller.InternalServerError(c, "failed to load page data", err)
		return "", true
	}
	if pageModel != nil {
		dto.enhanceFromModel(pageModel)
		pageKey := pageModel.SchemaName + "/" + pageModel.Identifier
		if latestRoute, err := pc.routeSvc.GetLatestVersion(c.Request.Context(), pageKey); err != nil {
			log.Error().
				Err(err).
				Str("pageKey", pageKey).
				Msg("failed to get latest route")
		} else if latestRoute != nil {
			dto.Route = latestRoute.Route
		}
	}

	// TODO: Build listable properties list for template (slice of {Name, Value})
	listableProperties := make([]map[string]any, 0)
	for _, field := range dto.Fields {
		if field.IsListable {
			if val, ok := dto.ListableData[field.Name]; ok {
				listableProperties = append(listableProperties, map[string]any{
					"name":  field.Name,
					"value": val,
				})
			}
		}
	}

	ctx := map[string]any{
		"class":              class,
		"identifier":         identifier,
		"page":               dto,
		"listableData":       dto.ListableData,
		"listableProperties": listableProperties,
	}

	body, err := tpl.AdminPageEdit.Exec(ctx)
	if err != nil {
		controller.TemplateRenderError(c, err)
		return "", true
	}

	output, err := template.AdminIndex(c, template.Content{
		Title:    "Welcome to Zhero",
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

func (pc *Controller) GetValidSlug(c *gin.Context) {
	customRoute := c.PostForm("route")

	slug, err := pc.routeSvc.GetValidSlug(c.Request.Context(), customRoute)
	if err != nil {
		c.Data(http.StatusBadRequest, "text/plain", []byte(err.Error()))
		return
	}

	c.Data(http.StatusOK, "text/plain", []byte(slug))
}

// TODO: Helper to get listable property names for a schema
func getListablePropertyNames(meta *schema.SchemaMeta) []string {
	names := []string{}
	for _, p := range meta.Properties {
		if p.Listable {
			names = append(names, p.Name)
		}
	}
	return names
}

// TODO: determineComponent automatically selects the appropriate HTML component based on Type and Name
func determineComponent(propType, propName string) string {
	nameLower := strings.ToLower(propName)
	switch {
	case strings.Contains(nameLower, "color"):
		return "Color"
	case strings.Contains(nameLower, "email"):
		return "Email"
	case strings.Contains(nameLower, "file"):
		return "File"
	case strings.Contains(nameLower, "phone") || strings.Contains(nameLower, "tel"):
		return "Tel"
	}

	switch propType {
	case "Boolean":
		return "Checkbox"
	case "Date":
		return "Date"
	case "DateTime":
		return "DateTime"
	case "Number", "Integer", "Float":
		return "Number"
	case "Quantity":
		return "TextInput"
	case "Text":
		return "TextInput"
	case "URL":
		return "URL"
	case "Time":
		return "Time"
	default:
		// For any other type (Object, or schema.org types like "Person", "Organization")
		return "ReferenceSearch"
	}
}

// TODO: SearchReferences handles HTMX/JSON search for references
func (pc *Controller) SearchReferences(c *gin.Context) {
	schema := c.Query("schema")
	query := c.Query("q")
	field := c.Query("field")

	refs, err := pc.pageSvc.SearchReferences(c, schema, query)
	if err != nil {
		controller.InternalServerError(c, "search failed", err)
		return
	}

	// Render partial HTML for HTMX
	output, err := tpl.AdminReferenceSearchResults.Exec(map[string]any{
		"references": refs,
		"field":      field,
	})
	if err != nil {
		controller.TemplateRenderError(c, err)
		return
	}
	c.Data(http.StatusOK, gin.MIMEHTML, []byte(output))
}

// TODO: ReferenceModal returns the reference selection modal
func (pc *Controller) ReferenceModal(c *gin.Context) {
	schema := c.Query("schema")
	field := c.Query("field")
	output, err := tpl.AdminReferenceModal.Exec(map[string]any{
		"schema": schema,
		"field":  field,
	})
	if err != nil {
		controller.TemplateRenderError(c, err)
		return
	}
	c.Data(http.StatusOK, gin.MIMEHTML, []byte(output))
}

// TODO: ReferenceSelect inserts reference into field and returns updated input
func (pc *Controller) ReferenceSelect(c *gin.Context) {
	field := c.Query("field")
	identifier := c.Query("identifier")
	secondary := c.Query("secondary")
	linkText := c.Query("link-text")
	altText := c.Query("alt-text")

	// Use secondary identifier as default linkText if not provided
	if linkText == "" {
		linkText = secondary
	}
	if altText == "" {
		altText = secondary
	}

	refPath := fmt.Sprintf("%s/%s", c.Query("schema"), identifier)
	props := fmt.Sprintf(`{'linkText':'%s','altText':'%s'}`,
		url.QueryEscape(linkText), url.QueryEscape(altText))
	reference := fmt.Sprintf("#ZHERO#%s#%s#", refPath, props)

	c.String(http.StatusOK,
		fmt.Sprintf(`<input type="text" id="field-%s" name="field-%s" class="input input-bordered w-full" value="%s" />`,
			field, field, reference))
}
