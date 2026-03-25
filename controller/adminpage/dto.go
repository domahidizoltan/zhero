package adminpage

import (
	"time"

	page_domain "github.com/domahidizoltan/zhero/domain/page"
	"github.com/domahidizoltan/zhero/domain/schema"
	"github.com/gin-gonic/gin"
)

type (
	pageDto struct {
		Route                    string
		SchemaName               string
		Fields                   []fieldDto
		Identifier               string
		SecondaryIdentifier      string
		SecondaryIdentifierValue any
		CreatedBy                string
		CreatedAt                time.Time
		UpdatedBy                string
		UpdatedAt                time.Time
		IsEnabled                bool
		Meta                     pageMeta
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

	pageMeta struct {
		Title         string
		Description   string
		OGTitle       string
		OGDescription string
		Rating        string
		Robots        []string
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

	dto.Meta = pageMeta{
		Title:         c.PostForm("meta-title"),
		Description:   c.PostForm("meta-description"),
		OGTitle:       c.PostForm("meta-og-title"),
		OGDescription: c.PostForm("meta-og-description"),
	}

	if c.PostForm("meta-robots-noindex") == "on" {
		dto.Meta.Robots = append(dto.Meta.Robots, "noindex")
	}
	if c.PostForm("meta-robots-nofollow") == "on" {
		dto.Meta.Robots = append(dto.Meta.Robots, "nofollow")
	}

	if c.PostForm("meta-rating-adult") == "on" {
		dto.Meta.Rating = "adult"
	}
}

func (dto *pageDto) enhanceFromModel(p *page_domain.Page) {
	if p == nil {
		return
	}

	dto.IsEnabled = p.IsEnabled
	dto.Route = p.Route
	dto.Meta.FromModel(p.Meta)
	dto.SecondaryIdentifierValue = p.SecondaryIdentifier
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
		Meta:                dto.Meta.ToModel(),
	}
}

func (dto *pageDto) ToMap() map[string]any {
	if dto == nil {
		return nil
	}

	fields := make(map[string]any, len(dto.Fields))
	for _, f := range dto.Fields {
		fields[f.Name] = f.Value
	}
	return map[string]any{
		"route":                    dto.Route,
		"schemaName":               dto.SchemaName,
		"fields":                   fields,
		"identifier":               dto.Identifier,
		"secondaryIdentifier":      dto.SecondaryIdentifier,
		"secondaryIdentifierValue": dto.SecondaryIdentifierValue,
		"isEnabled":                dto.IsEnabled,
		"meta":                     dto.Meta.ToMap(),
	}
}

func (dm *pageMeta) FromModel(pm page_domain.PageMeta) {
	*dm = pageMeta{
		Title:         pm.Title,
		Description:   pm.Description,
		OGTitle:       pm.OGTitle,
		OGDescription: pm.OGDescription,
		Rating:        pm.Rating,
		Robots:        pm.Robots,
	}
}

func (dm *pageMeta) ToModel() page_domain.PageMeta {
	if dm == nil {
		return page_domain.PageMeta{}
	}

	return page_domain.PageMeta{
		Title:         dm.Title,
		Description:   dm.Description,
		OGTitle:       dm.OGTitle,
		OGDescription: dm.OGDescription,
		Rating:        dm.Rating,
		Robots:        dm.Robots,
	}
}

func (dm *pageMeta) ToMap() map[string]any {
	if dm == nil {
		return nil
	}

	return map[string]any{
		"title":         dm.Title,
		"description":   dm.Description,
		"rating":        dm.Rating,
		"robots":        dm.Robots,
		"ogTitle":       dm.OGTitle,
		"ogDescription": dm.OGDescription,
	}
}
