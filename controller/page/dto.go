package page

import (
	"strconv"
	"time"

	page_domain "github.com/domahidizoltan/zhero/domain/page"
	"github.com/domahidizoltan/zhero/domain/schema"
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
		Order        uint
		IsMandatory  bool
		IsSearchable bool
		Type         string
		Component    string
		Value        any
	}

	pagingDto struct {
		BaseURL string
		First   string
		Prev    []uint
		Current uint
		Next    []uint
		Last    string
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
}

func (dto *pageDto) enhanceFromModel(p *page_domain.Page) {
	if p == nil {
		return
	}

	dto.IsEnabled = p.IsEnabled
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
		SchemaName:          dto.SchemaName,
		Identifier:          data[dto.Identifier].(string),
		SecondaryIdentifier: data[dto.SecondaryIdentifier].(string),
		Data:                data,
		IsEnabled:           dto.IsEnabled,
		SearchVals:          searchVals,
	}
}

const jump = 3

func pagingDtoFrom(p page_domain.PagingMeta, baseURL string) *pagingDto {
	if p.TotalPages < 1 {
		return nil
	}

	pg := &pagingDto{
		BaseURL: baseURL,
		Current: p.CurrentPage,
	}

	if p.CurrentPage > jump+1 {
		pg.First = strconv.Itoa(1)
	}
	for i := range jump {
		if this := int(p.CurrentPage) - (jump - i); this > 0 {
			pg.Prev = append(pg.Prev, uint(this))
		}
	}

	if limit := int(p.TotalPages) - jump; limit >= 0 && p.CurrentPage < uint(limit) {
		pg.Last = strconv.Itoa(int(p.TotalPages))
	}
	for i := range jump {
		if this := p.CurrentPage + uint(i) + 1; this <= p.TotalPages {
			pg.Next = append(pg.Next, this)
		}
	}

	return pg
}
