package pagerenderer

import (
	"fmt"
	"strings"

	page_ctrl "github.com/domahidizoltan/zhero/controller/page"
	"github.com/domahidizoltan/zhero/domain/page"
	"github.com/domahidizoltan/zhero/domain/schema"
)

type DynamicPageRenderer struct{}

func NewDynamicPageRenderer() DynamicPageRenderer {
	return DynamicPageRenderer{}
}

func (DynamicPageRenderer) Render(meta schema.SchemaMeta, data map[string]any) (string, error) {
	b := strings.Builder{}
	for _, prop := range meta.Properties {
		if prop.Name == meta.Identifier || strings.HasPrefix(prop.Name, "@") {
			continue
		}

		v := data[prop.Name]
		cssClass := strings.ToLower(meta.Name + "-" + prop.Name)
		if prop.Name == meta.SecondaryIdentifier {
			b.WriteString(fmt.Sprintf("<h1 class=\"%s\">%s</h1>", cssClass, v))
			continue
		}

		b.WriteString(fmt.Sprintf("<p class=\"%s\">%s</p>", cssClass, v))
	}
	return b.String(), nil
}

func (DynamicPageRenderer) List(meta schema.SchemaMeta, data []map[string]any, paging page.PagingMeta) (string, error) {
	b := strings.Builder{}

	cssClass := strings.ToLower("list-item " + meta.Name)
	b.WriteString("<div class=\"list\">")
	for _, d := range data {
		link := fmt.Sprintf("/%s/%s", meta.Name, d[meta.Identifier])
		secID := d[meta.SecondaryIdentifier]
		b.WriteString(fmt.Sprintf(`<div class="%[1]s">
    <div class="img skeleton" onclick="window.location.href='%[2]s'"></div>
    <div class="content">
			<a href="%[2]s">%[3]s</a>
		</div>
</div>`, cssClass, link, secID))
	}
	b.WriteString("</div>")

	if paging.TotalPages <= 1 {
		return b.String(), nil
	}

	baseURL := fmt.Sprintf("/%s?", meta.Name)
	dto := page_ctrl.PagingDtoFrom(paging, baseURL)
	if dto != nil {
		b.WriteString(`<div class="pagination mt-6 flex justify-between items-center"><div class="join">`)

		if dto.First != "" {
			b.WriteString(fmt.Sprintf(`<a href="%s&page=%s" class="join-item btn btn-sm">«</a>`, baseURL, dto.First))
		}

		for _, prev := range dto.Prev {
			b.WriteString(fmt.Sprintf(`<a href="%[1]s&page=%[2]d" class="join-item btn btn-sm">%[2]d</a>`, baseURL, prev))
		}

		b.WriteString(fmt.Sprintf(`<a class="join-item btn btn-sm btn-active">%d</a>`, dto.Current))

		for _, next := range dto.Next {
			b.WriteString(fmt.Sprintf(`<a href="%[1]s&page=%[2]d" class="join-item btn btn-sm">%[2]d</a>`, baseURL, next))
		}

		if dto.Last != "" {
			b.WriteString(fmt.Sprintf(`<a hx-get="%s&page=%s" class="join-item btn btn-sm">»</a>`, baseURL, dto.Last))
		}

		b.WriteString("</div></div>")
	}

	return b.String(), nil
}

// TODO static page renderer and preview page
