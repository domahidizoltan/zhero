// Package page contains the controllers for the pages
package page

import (
	"fmt"
	"net/http"

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

func (pc *Controller) Save(c *gin.Context) {
	output, hasError := pc.edit(c, true)
	if hasError {
		c.Data(http.StatusBadRequest, gin.MIMEHTML, []byte(output))
		return
	}
	// class := c.Param("class")
	c.Redirect(http.StatusSeeOther, "/page/list")
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
		// schemaToSave, validationErrs, err := pc.fromForm(c)
		// if schemaToSave == nil {
		// 	if len(validationErrs) > 0 {
		// 		errorMsg = "Validation errors:\\n" + strings.Join(validationErrs, "\\n")
		// 	} else if err != nil {
		// 		errorMsg = err.Error()
		// 	}
		// } else if err := sc.schemaSvc.SaveSchemaMeta(c.Request.Context(), *schemaToSave); err != nil {
		// 	log.Error().Err(err).Msg("failed to save schema")
		// 	errorMsg = err.Error()
		// } else {
		// 	successMsg = fmt.Sprintf("Schema %s saved successfully", clsName)
		// }
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
