package adminpage

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/domahidizoltan/zhero/controller"
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
		ListableData             map[string]any
		References               []string
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
		IsListable   bool
		Type         string
		Component    string
		InputType    bool
		Value        any
	}

	pageMeta struct {
		Title         string
		Description   string
		OGTitle       string
		OGDescription string
		Rating        string
		Robots        []string
		FieldMeta     map[string]string
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
		component := p.Component
		if component == "" || component == "TODO" {
			component = determineComponent(p.Type, p.Name)
		}
		dto.Fields = append(dto.Fields, fieldDto{
			Name:         p.Name,
			Order:        p.Order,
			IsMandatory:  p.Mandatory,
			IsSearchable: p.Searchable,
			IsListable:   p.Listable,
			Type:         p.Type,
			Component:    component,
			InputType:    slices.Contains([]string{"Color", "Email", "Tel", "Number", "Date", "DateTime", "Time"}, component),
		})
	}

	return dto
}

func (dto *pageDto) EnhanceFromForm(c *gin.Context) {
	dto.IsEnabled = c.PostForm("is-enabled") == "on"
	dto.Route = c.PostForm("route")

	dto.Meta = pageMeta{
		Title:         c.PostForm("meta-title"),
		Description:   c.PostForm("meta-description"),
		OGTitle:       c.PostForm("meta-og-title"),
		OGDescription: c.PostForm("meta-og-description"),
		FieldMeta:     map[string]string{},
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

	for i, f := range dto.Fields {
		fieldValue := c.PostForm("field-" + f.Name)
		dto.Fields[i].Value = fieldValue
		if altText, found := c.GetPostForm("alt-text-" + f.Name); found && fieldValue != "" {
			dto.Meta.FieldMeta[f.Name+":altText"] = altText
		}
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
	dto.ListableData = make(map[string]any, len(p.ListableData))

	refs := make(map[string]string, len(p.References))
	for idx, ref := range p.References {
		refs[fmt.Sprintf("#ref%d", idx)] = ref
	}

	for k, v := range p.ListableData {
		dto.ListableData[k] = replaceShortRefs(v, refs)
	}

	for i, f := range dto.Fields {
		if val, ok := p.Data[f.Name]; ok {
			dto.Fields[i].Value = replaceShortRefs(val, refs)
		}
	}
}

func replaceShortRefs(field any, refs map[string]string) any {
	vString, ok := field.(string)
	if ok {
		for k, v := range refs {
			vString = strings.ReplaceAll(vString, k, v)
		}
		field = vString
	}
	return field
}

func (dto *pageDto) ToModel() page_domain.Page {
	searchVals := [page_domain.MaxSearchVals]any{}
	data := make(map[string]any, len(dto.Fields))
	scIdx := 0
	listableData := make(map[string]any)

	for _, f := range dto.Fields {
		val := f.Value
		data[f.Name] = val

		if f.IsListable {
			listableData[f.Name] = val
		}

		if f.IsSearchable && f.Name != dto.SecondaryIdentifier && scIdx < 5 {
			searchVals[scIdx] = f.Value
			scIdx++
		}
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
		ListableData:        listableData,
		References:          dto.References,
	}
}

func (dto *pageDto) extractReferences() {
	fieldRefs := map[string]int{}
	refCounter := 0
	for i, f := range dto.Fields {
		if f.Value == nil {
			continue
		}
		strVal, ok := f.Value.(string)
		if !ok {
			continue
		}

		if strings.Contains(strVal, "#") {
			matches := controller.RefPattern.FindAllStringSubmatch(strVal, -1)
			for _, match := range matches {
				ref := match[0]
				replaceRef := ""
				if r, found := fieldRefs[ref]; found {
					replaceRef = fmt.Sprintf("#ref%d", r)
				} else {
					replaceRef = fmt.Sprintf("#ref%d", refCounter)
					fieldRefs[ref] = refCounter
					refCounter++
				}
				strVal = strings.ReplaceAll(strVal, ref, replaceRef)
			}
		}
		dto.Fields[i].Value = any(strVal)
	}

	dto.References = make([]string, len(fieldRefs))
	for link, idx := range fieldRefs {
		dto.References[idx] = link
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
		"listableData":             dto.ListableData,
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
		FieldMeta:     pm.FieldMeta,
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
		FieldMeta:     dm.FieldMeta,
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
