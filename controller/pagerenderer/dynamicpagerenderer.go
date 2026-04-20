package pagerenderer

import (
	"fmt"
	"html"
	"regexp"
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
		if v == nil {
			continue
		}
		cssClass := strings.ToLower(meta.Name + "-" + prop.Name)

		if prop.Name == meta.SecondaryIdentifier {
			b.WriteString(fmt.Sprintf("<h1 class=\"%s\">%s</h1>", cssClass, v))
			continue
		}

		// TODO: Check if value is a string with references
		if strVal, ok := v.(string); ok && strings.Contains(strVal, "#ZHERO#") {
			rendered := renderReferences(strVal)
			b.WriteString(fmt.Sprintf("<p class=\"%s\">%s</p>", cssClass, rendered))
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
		delete(listableProperties, listable.SecondaryIdentifier)
		delete(listableProperties, listable.Identifier)

		var image string
		details := strings.Builder{}
		for k, v := range listableProperties {
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

// TODO: static page renderer and preview page

func renderReferences(text string) string {
	re := regexp.MustCompile(`#ZHERO#([^#]+)#\{([^}]*)\}#`)
	return re.ReplaceAllStringFunc(text, func(match string) string {
		submatches := re.FindStringSubmatch(match)
		if len(submatches) < 3 {
			return match
		}
		refPath := submatches[1]   // "Thing/123"
		propsJSON := submatches[2] // "{'linkText':'...','altText':'...'}"

		// Parse properties (simple single-quoted JS object)
		altText := extractProp(propsJSON, "altText")
		linkText := extractProp(propsJSON, "linkText")
		if linkText == "" {
			linkText = refPath // fallback
		}

		// HTML escape for safety
		linkText = html.EscapeString(linkText)
		altText = html.EscapeString(altText)

		return fmt.Sprintf(`<a href="/%s" title="%s">%s</a>`, refPath, altText, linkText)
	})
}

func extractProp(props, key string) string {
	pattern := fmt.Sprintf(`%s:\s*'([^']*)'`, key)
	re := regexp.MustCompile(pattern)
	if m := re.FindStringSubmatch(props); len(m) > 1 {
		return m[1]
	}
	return ""
}
