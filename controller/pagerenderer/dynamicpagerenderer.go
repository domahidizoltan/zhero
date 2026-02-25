package pagerenderer

import (
	"fmt"
	"strings"

	"github.com/domahidizoltan/zhero/domain/schema"
	"github.com/domahidizoltan/zhero/pkg/paging"
	"github.com/domahidizoltan/zhero/template"
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

func (DynamicPageRenderer) List(meta schema.SchemaMeta, data []map[string]any, paging paging.Meta) (string, error) {
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

	baseURL := fmt.Sprintf("/%s?", meta.Name)
	dto := paging.ToDto(baseURL, "")
	if dto != nil {
		if pagination, err := template.PaginationPartial.Exec(map[string]any{"paging": dto}); err != nil {
			return "", err
		} else {
			b.WriteString(pagination)
		}
	}

	return b.String(), nil
}

// TODO static page renderer and preview page
