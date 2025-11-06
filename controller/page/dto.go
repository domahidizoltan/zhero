package page

import (
	"maps"
	"time"

	page_domain "github.com/domahidizoltan/zhero/domain/page"
	"github.com/domahidizoltan/zhero/domain/schema"
	"github.com/domahidizoltan/zhero/pkg/collection"
	"github.com/gin-gonic/gin"
)

type (
	pageDto struct {
		SchemaName          string
		Fields              []fieldDto
		Identifier          string
		SecondaryIdentifier string
		CreatedBy           string
		CreatedAt           time.Time
		UpdatedBy           string
		UpdatedAt           time.Time
		IsEnabled           bool
	}

	fieldDto struct {
		Name         string
		Order        int
		IsMandatory  bool
		IsSearchable bool
		Type         string
		Component    string
		Value        any
	}
)

func pageDtoFrom(meta *schema.SchemaMeta) pageDto {
	if meta == nil {
		return pageDto{}
	}

	dto := pageDto{
		SchemaName:          meta.Name,
		Identifier:          meta.Identifier,
		SecondaryIdentifier: meta.SecondaryIdentifier,
		IsEnabled:           false,
	}

	dto.Fields = make([]fieldDto, 0, len(meta.Properties))
	for _, p := range meta.Properties {
		dto.Fields = append(dto.Fields, fieldDto{
			Name:         p.Name,
			Order:        p.Order,
			IsMandatory:  p.Mandatory,
			IsSearchable: p.Searchable,
			Type:         p.Type,
			Component:    p.Component,
			// Value:        ,
		})
	}

	return dto
}

func (dto *pageDto) enhanceFromForm(c *gin.Context) {
	for i, f := range dto.Fields {
		dto.Fields[i].Value = c.PostForm("field-" + f.Name)
	}
}

func (dto *pageDto) enhanceFromModel(p *page_domain.Page) {
	if p == nil {
		return
	}

	modelFieldsByName := maps.Collect(collection.MapBy(p.Fields, func(f page_domain.Field) (string, page_domain.Field) {
		return f.Name, f
	}))

	for i, f := range dto.Fields {
		if f, ok := modelFieldsByName[f.Name]; ok {
			dto.Fields[i].Value = f.Value
		}
	}
}

func (dto *pageDto) toModel() page_domain.Page {
	fields := make([]page_domain.Field, 0, len(dto.Fields))
	for _, f := range dto.Fields {
		field := page_domain.Field{
			Name:  f.Name,
			Order: f.Order,
			Value: f.Value,
		}
		if f.IsSearchable {
			field.SearchColumn = f.Name
		}
		fields = append(fields, field)
	}

	return page_domain.Page{
		SchemaName:          dto.SchemaName,
		Fields:              fields,
		Identifier:          dto.Identifier,
		SecondaryIdentifier: dto.SecondaryIdentifier,
		IsEnabled:           dto.IsEnabled,
	}
}
