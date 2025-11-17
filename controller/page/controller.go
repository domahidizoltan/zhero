// Package page contains the controllers for the pages
package page

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/aymerick/raymond"
	"github.com/domahidizoltan/zhero/controller"
	"github.com/domahidizoltan/zhero/controller/template"
	"github.com/domahidizoltan/zhero/domain/page"
	"github.com/domahidizoltan/zhero/domain/schema"
	tpl "github.com/domahidizoltan/zhero/template"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Controller struct {
	schemaSvc schema.Service
	pageSvc   page.Service
}

func NewController(schemaSvc schema.Service, pageSvc page.Service) Controller {
	return Controller{
		schemaSvc: schemaSvc,
		pageSvc:   pageSvc,
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
	pageNo, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || pageNo < 1 {
		pageNo = 1
	}

	search := c.Query("search")
	opts := page.ListOptions{
		SecondaryIdentifierLike: search,
		Page:                    uint(pageNo),
	}

	sortQuery := c.DefaultQuery("sort", "identifier:asc")
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

	pages, paging, err := pc.pageSvc.List(c, clsName, opts)
	if err != nil {
		controller.InternalServerError(c, "failed to list pages", err)
		return
	}

	sort := fmt.Sprintf("%s:%s", opts.SortBy, opts.SortDir)
	urlQuery := fmt.Sprintf("search=%s&sort=%s", search, sort)
	ctx := map[string]any{
		"class":    clsName,
		"paging":   pagingDtoFrom(paging, fmt.Sprintf("/page/list/%s?%s", clsName, urlQuery)),
		"urlQuery": fmt.Sprintf("%s&page=%d", urlQuery, pageNo),
		"pages":    pages,
		"search":   search,
		"sort":     sortQuery,
		"listOpts": opts,
	}

	output, err := tpl.PageList.Exec(ctx)
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

	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/page/list/%s?search=%s&sort=%s&page=%s", class, search, sort, page))
}

func (pc *Controller) Save(c *gin.Context) {
	output, hasError := pc.edit(c, true)
	if hasError {
		c.Data(http.StatusBadRequest, gin.MIMEHTML, []byte(output))
		return
	}
	class := c.Param("class")
	c.Redirect(http.StatusSeeOther, "/page/list?schema="+class)
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

	dto = pageDtoFrom(meta)

	errorMsg, successMsg := "", ""
	if hasFormSubmitted {
		dto.enhanceFromForm(c)
		page := dto.toModel()

		var err error
		if len(identifier) == 0 {
			identifier, err = pc.pageSvc.Create(c, page)
		} else {
			err = pc.pageSvc.Update(c, identifier, page)
		}

		if err != nil {
			log.Error().Err(err).Msg("failed to save page")
			errorMsg = err.Error()
		} else {
			successMsg = fmt.Sprintf("\"%s\" page saved successfully with ID %s", class, identifier)
		}

	}

	pageModel, err := pc.pageSvc.GetPageBySchemaNameAndIdentifier(c.Request.Context(), class, identifier)
	if err != nil {
		controller.InternalServerError(c, "failed to load page date", err)
	}
	dto.enhanceFromModel(pageModel)

	ctx := map[string]any{
		"class":      class,
		"identifier": identifier,
		"page":       dto,
	}
	body, err := tpl.PageEdit.Exec(ctx)
	if err != nil {
		controller.TemplateRenderError(c, err)
		return "", true
	}

	output, err := template.Index(c, template.Content{
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
