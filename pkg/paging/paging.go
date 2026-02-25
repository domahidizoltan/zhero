// Package paging is for pagination related structs and functions
package paging

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type (
	Meta struct {
		TotalItems  uint
		PageSize    uint
		TotalPages  uint
		CurrentPage uint
	}

	PageOpts struct {
		search   string
		SortBy   string
		SortDir  SortDir
		Page     uint
		PageSize uint
	}

	Dto struct {
		BaseURL    string
		HtmxTarget string
		First      string
		Prev       []uint
		Current    uint
		Next       []uint
		Last       string
	}

	SortDir string
)

var (
	SortDirAsc  SortDir = "asc"
	SortDirDesc SortDir = "desc"
)

var jump uint

func SetJump(j uint) {
	jump = j
}

func (p Meta) ToDto(baseURL, htmxTarget string) *Dto {
	if p.TotalPages < 1 {
		return nil
	}

	pg := &Dto{
		BaseURL: baseURL,
		Current: p.CurrentPage,
	}

	if strings.HasPrefix(htmxTarget, "#") {
		pg.HtmxTarget = htmxTarget
	}

	if p.CurrentPage > jump+1 {
		pg.First = strconv.Itoa(1)
	}
	for i := range jump {
		if this := int(p.CurrentPage) - int(jump-i); this > 0 {
			pg.Prev = append(pg.Prev, uint(this))
		}
	}

	if limit := p.TotalPages - jump; limit >= 0 && p.CurrentPage < limit {
		pg.Last = strconv.Itoa(int(p.TotalPages))
	}
	for i := range jump {
		if this := p.CurrentPage + uint(i) + 1; this <= p.TotalPages {
			pg.Next = append(pg.Next, this)
		}
	}

	return pg
}

func RequestToPageOpts(c *gin.Context, defaultSearchParam string) PageOpts {
	pageNo, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || pageNo < 1 {
		pageNo = 1
	}

	search := c.Query("search")
	opts := PageOpts{
		search: search,
		Page:   uint(pageNo),
	}

	defaultSearch := ""
	if len(defaultSearchParam) > 0 {
		defaultSearch = defaultSearchParam + ":asc"
	}
	sortQuery := c.DefaultQuery("sort", defaultSearch)
	unescape, err := url.QueryUnescape(sortQuery)
	if err != nil {
		unescape = sortQuery
	}
	sortQuery = unescape

	if sortBy, sortOrder, found := strings.Cut(sortQuery, ":"); found {
		opts.SortBy = sortBy
		opts.SortDir = SortDir(sortOrder)
	}

	return opts
}

func (po PageOpts) SearchParam() string {
	return po.search
}

func (po PageOpts) SortQuery() string {
	if len(po.SortBy) > 0 {
		return fmt.Sprintf("%s:%s", po.SortBy, po.SortDir)
	}
	return ""
}

func (po PageOpts) GetURL(baseURL string) (string, string) {
	sort := fmt.Sprintf("%s:%s", po.SortBy, po.SortDir)
	urlQuery := fmt.Sprintf("search=%s&sort=%s", po.search, sort)

	urlSeparator := "?"
	if strings.Contains(baseURL, "?") {
		urlSeparator = "&"
	}
	pagingBaseURL := baseURL + urlSeparator + urlQuery
	return pagingBaseURL, urlQuery
}

func (po PageOpts) ToMeta(total int, defaultPageSize uint) Meta {
	meta := Meta{
		PageSize:    po.PageSize,
		CurrentPage: po.Page,
		TotalItems:  uint(total),
	}

	if po.PageSize == 0 && defaultPageSize > 0 {
		meta.PageSize = defaultPageSize
	}

	meta.TotalPages = uint(meta.TotalItems / meta.PageSize)

	return meta
}
