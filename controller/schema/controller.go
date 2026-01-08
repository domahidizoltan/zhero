// Package schema defines the handlers for managing schema
package schema

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/aymerick/raymond"
	"github.com/domahidizoltan/zhero/controller"
	"github.com/domahidizoltan/zhero/controller/template"
	"github.com/domahidizoltan/zhero/domain/schema"
	"github.com/domahidizoltan/zhero/pkg/collection"
	tpl "github.com/domahidizoltan/zhero/template"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

type Controller struct {
	schemaSvc schema.Service
}

func NewController(schemaSvc schema.Service) Controller {
	return Controller{
		schemaSvc: schemaSvc,
	}
}

func (sc *Controller) Search(c *gin.Context) {
	output, err := tpl.SchemaorgSearch.Exec(nil)
	if err != nil {
		controller.TemplateRenderError(c, err)
		return
	}

	c.Data(http.StatusOK, gin.MIMEHTML, []byte(output))
}

func (sc *Controller) GetClassHierarchy(c *gin.Context) {
	c.JSON(http.StatusOK, sc.schemaSvc.GetClassHierarchy())
}

func (sc *Controller) Edit(c *gin.Context) {
	clsName := c.Param("class")
	output, hasError := sc.edit(c, clsName, false)
	if hasError {
		return
	}
	c.Data(http.StatusOK, gin.MIMEHTML, []byte(output))
}

func (sc *Controller) Save(c *gin.Context) {
	if err := c.Request.ParseForm(); err != nil {
		controller.BadRequest(c, "invalid form", err)
		return
	}

	clsName := c.Param("class")
	output, hasError := sc.edit(c, clsName, true)
	if hasError {
		c.Data(http.StatusBadRequest, gin.MIMEHTML, []byte(output))
		return
	}

	c.Redirect(http.StatusSeeOther, "/")
}

func (sc *Controller) edit(c *gin.Context, clsName string, hasFormSubmitted bool) (string, bool) {
	if clsName == "" {
		controller.BadRequest(c, "class is missing", nil)
		return "", true
	}

	orgSchema := sc.schemaSvc.GetSchemaClassByName(clsName)
	savedSchema, err := sc.schemaSvc.GetSchemaMetaByName(c, clsName)
	if err != nil {
		controller.InternalServerError(c, "failed to get existing schema metadata", err)
		return "", true
	}

	errorMsg, successMsg := "", ""
	if hasFormSubmitted {
		schemaToSave, validationErrs, err := sc.schemaFromForm(c, clsName)
		if schemaToSave == nil {
			if len(validationErrs) > 0 {
				errorMsg = "Validation errors:\\n" + strings.Join(validationErrs, "\\n")
			} else if err != nil {
				errorMsg = err.Error()
			}
		} else if err := sc.schemaSvc.SaveSchemaMeta(c.Request.Context(), *schemaToSave); err != nil {
			log.Error().Err(err).Msg("failed to save schema")
			errorMsg = err.Error()
		} else {
			successMsg = fmt.Sprintf("Schema %s saved successfully", clsName)
		}
	}

	dto := schemaDtoFrom(*orgSchema, savedSchema)
	ctx := map[string]any{
		"class":       dto,
		"breadcrumbs": sc.classBreadcrumbs(clsName),
		"components":  []string{"TODO"},
	}
	body, err := tpl.SchemaorgEdit.Exec(ctx)
	if err != nil {
		controller.TemplateRenderError(c, err)
		return "", true
	}

	operation := "Create"
	output, err := template.AdminIndex(c, template.Content{
		Title:    operation + " schema: " + clsName,
		Body:     raymond.SafeString(body),
		ErrorMsg: errorMsg,
		FlashMsg: successMsg,
	})
	if err != nil {
		controller.TemplateRenderError(c, err)
		return "", true
	}

	return output, len(errorMsg) > 0
}

var (
	setPropMandatory  = func(p schema.Property, v bool) schema.Property { p.Mandatory = v; return p }
	setPropSearchable = func(p schema.Property, v bool) schema.Property { p.Searchable = v; return p }
)

func (sc *Controller) schemaFromForm(c *gin.Context, clsName string) (*schema.SchemaMeta, []string, error) {
	// TODO refactor to bind to dto instead
	var schemaToSave schema.SchemaMeta
	if err := c.Bind(&schemaToSave); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			errs := slices.Collect(collection.MapValues(validationErrors, func(e validator.FieldError) string {
				return fmt.Sprintf("- %s is %s", e.Field(), e.Tag())
			}))
			return nil, errs, nil
		}
		return nil, nil, err
	}

	schemaToSave.Name = clsName
	props := map[string]schema.Property{}
	for i, name := range c.PostFormArray("property-name") {
		props[name] = schema.Property{
			Name:      name,
			Type:      c.PostFormArray("property-type")[i],
			Component: c.PostFormArray("property-component")[i],
		}
	}
	alterMap(c, "property-mandatory", props, func(p schema.Property) schema.Property { return setPropMandatory(p, true) })
	alterMap(c, "property-searchable", props, func(p schema.Property) schema.Property { return setPropSearchable(p, true) })

	props[schemaToSave.Identifier] = setPropMandatory(props[schemaToSave.Identifier], false)
	props[schemaToSave.Identifier] = setPropSearchable(props[schemaToSave.Identifier], false)
	props[schemaToSave.SecondaryIdentifier] = setPropMandatory(props[schemaToSave.SecondaryIdentifier], true)
	props[schemaToSave.SecondaryIdentifier] = setPropSearchable(props[schemaToSave.SecondaryIdentifier], true)

	propertyOrder := strings.Split(c.PostForm("property-order"), ",")
	for i, p := range slices.Collect(collection.Unique(propertyOrder)) {
		if prop, found := props[p]; found {
			prop.Order = uint(i)
			schemaToSave.Properties = append(schemaToSave.Properties, prop)
		}
	}
	return &schemaToSave, nil, nil
}

func alterMap[T any](c *gin.Context, key string, itemsByName map[string]T, setter func(p T) T) {
	for _, name := range c.PostFormArray(key) {
		if p, found := itemsByName[name]; found {
			itemsByName[name] = setter(p)
		}
	}
}

func (sc *Controller) classBreadcrumbs(clsName string) []string {
	var breadcrumbs []string
	for _, ch := range sc.schemaSvc.GetClassHierarchy() {
		if ch[len(ch)-1] == clsName {
			breadcrumbs = append(breadcrumbs, ch...)
			break
		}
	}
	return breadcrumbs
}
