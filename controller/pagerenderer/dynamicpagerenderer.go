package pagerenderer

import (
	"fmt"
	"strings"

	"github.com/domahidizoltan/zhero/controller"
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
	refs := data["references"].(map[string]string)
	for _, prop := range meta.Properties {
		if prop.Name == meta.Identifier || strings.HasPrefix(prop.Name, "@") {
			continue
		}

		v := data[prop.Name]
		if v == nil {
			continue
		}
		cssClass := strings.ToLower(meta.Name + "-" + prop.Name)

		if prop.Name == meta.SecondaryIdentifier {
			b.WriteString(fmt.Sprintf("<h1 class=\"%s\">%s</h1>", cssClass, v))
			continue
		}

		if strVal, ok := v.(string); ok {
			for _, s := range controller.RefPattern.FindAllString(strVal, -1) {
				if link, found := refs[s]; found {
					strVal = strings.ReplaceAll(strVal, s, link)
				}
			}
			b.WriteString(fmt.Sprintf("<p class=\"%s\">%s</p>", cssClass, strVal))
		} else {
			b.WriteString(fmt.Sprintf("<p class=\"%s\">%s</p>", cssClass, v))
		}
	}
	return b.String(), nil
}

func (DynamicPageRenderer) List(listable schema.SchemaMeta, data []map[string]any, paging paging.Meta) (string, error) {
	b := strings.Builder{}

	cssClass := strings.ToLower("list-item " + listable.Name)
	b.WriteString("<div class=\"list\">")
	for _, d := range data {
		link := fmt.Sprintf("/%s/%s", listable.Name, d[listable.Identifier])
		secID := d[listable.SecondaryIdentifier]
		listableProperties := d["listableProperties"].(map[string]any)
		refs := listableProperties["references"].(map[string]string)
		delete(listableProperties, listable.SecondaryIdentifier)
		delete(listableProperties, listable.Identifier)
		delete(listableProperties, "references")

		var image string
		details := strings.Builder{}
		for k, v := range listableProperties {
			if strVal, ok := v.(string); ok {
				for _, s := range controller.RefPattern.FindAllString(strVal, -1) {
					if link, found := refs[s]; found {
						strVal = strings.ReplaceAll(strVal, s, link)
					}
				}
				v = strVal
			}
			key := strings.ToLower(k)
			if strings.Contains(key, "thumbnail") || strings.Contains(key, "image") {
				delete(listableProperties, key)
				image = fmt.Sprintf("<img src=\"%s\" />", v)
				continue
			}

			details.WriteString(fmt.Sprintf("<br/><span>%s</span>", v))
		}

		b.WriteString(fmt.Sprintf(`<div class="%[1]s">
    <div class="img skeleton" onclick="window.location.href='%[2]s'">%[4]s</div>
    <div class="content">
			<b><a href="%[2]s">%[3]s</a></b>
			%[5]v
		</div>
</div>`, cssClass, link, secID, image, details.String()))
	}
	b.WriteString("</div>")

	baseURL := fmt.Sprintf("/%s?", listable.Name)
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
