package adminpage

import (
	"time"

	page_domain "github.com/domahidizoltan/zhero/domain/page"
	"github.com/domahidizoltan/zhero/domain/schema"
	"github.com/gin-gonic/gin"
)

type (
	pageDto struct {
		Route               string
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
		Order        uint
		IsMandatory  bool
		IsSearchable bool
		Type         string
		Component    string
		Value        any
	}
)

func PageDtoFrom(meta *schema.SchemaMeta) pageDto {
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
		})
	}

	return dto
}

func (dto *pageDto) EnhanceFromForm(c *gin.Context) {
	for i, f := range dto.Fields {
		dto.Fields[i].Value = c.PostForm("field-" + f.Name)
	}
	dto.IsEnabled = c.PostForm("is-enabled") == "on"
	dto.Route = c.PostForm("route")
}

func (dto *pageDto) enhanceFromModel(p *page_domain.Page) {
	if p == nil {
		return
	}

	dto.IsEnabled = p.IsEnabled
	dto.Route = p.Route
	for i, f := range dto.Fields {
		if val, ok := p.Data[f.Name]; ok {
			dto.Fields[i].Value = val
		}
	}
}

func (dto *pageDto) ToModel() page_domain.Page {
	searchVals := [page_domain.MaxSearchVals]any{}
	data := make(map[string]any, len(dto.Fields))
	scIdx := 0
	for _, f := range dto.Fields {
		val := f.Value
		if f.IsSearchable && f.Name != dto.SecondaryIdentifier && scIdx < 5 {
			searchVals[scIdx] = f.Value
			scIdx++
		}
		data[f.Name] = val
	}

	return page_domain.Page{
		Route:               dto.Route,
		SchemaName:          dto.SchemaName,
		Identifier:          data[dto.Identifier].(string),
		SecondaryIdentifier: data[dto.SecondaryIdentifier].(string),
		Data:                data,
		IsEnabled:           dto.IsEnabled,
		SearchVals:          searchVals,
	}
}
