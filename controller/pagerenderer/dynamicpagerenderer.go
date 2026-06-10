package pagerenderer

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/domahidizoltan/zhero/controller"
	"github.com/domahidizoltan/zhero/controller/file"
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
	fieldMeta := map[string]string{}
	if fm, found := data["fieldMeta"].(map[string]string); found {
		fieldMeta = fm
	}

	for _, prop := range meta.Properties {
		if prop.Name == "thumbnail" {
			continue
		}

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

		strVal, ok := v.(string)
		if !ok {
			continue
		}

		for _, s := range controller.RefPattern.FindAllString(strVal, -1) {
			if link, found := refs[s]; found {
				strVal = strings.ReplaceAll(strVal, s, link)
			}
		}

		altText := ""
		if t := fieldMeta[prop.Name+":altText"]; t != "" {
			altText = t
		}

		switch prop.Component {
		case "File":
			fmt.Println("file")

			filePath := fmt.Sprintf("%s/%s/%s/%s", file.UploadsPath, meta.Name, data["identifier"], strVal)
			b.WriteString(fmt.Sprintf("<p class=\"%s\"><a href=\"%s\" alt=\"%s\" target=\"_blank\">%s</a></p>", cssClass, filePath, altText, altText))
		case "Image":
			filePath := fmt.Sprintf("%s/%s/%s/%s", file.UploadsPath, meta.Name, data["identifier"], strVal)
			b.WriteString(fmt.Sprintf("<p class=\"%s\"><img src=\"%s\" alt=\"%s\" /></p>", cssClass, filePath, altText))
		default:
			b.WriteString(fmt.Sprintf("<p class=\"%s\">%s</p>", cssClass, strVal))
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

		refs := map[string]string{}
		if _, found := listableProperties["references"]; found {
			refs = listableProperties["references"].(map[string]string)
			delete(listableProperties, "references")
		}

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
				if v == "" {
					continue
				}
				src, alt := "", ""
				if !strings.HasPrefix(src, "/") {
					src = fmt.Sprintf("%s/%s/%s/%s", file.UploadsPath, listable.Name, d[listable.Identifier], v)
				}
				image = fmt.Sprintf("<img src=\"%s\" alt=\"%s\" />", src, alt)
				continue
			}

			if strVal, ok := v.(string); ok {
				linkText := filepath.Base(strVal)
				if altVal, found := listableProperties[k+"_alt"]; found {
					if altStr, ok := altVal.(string); ok && altStr != "" {
						linkText = altStr
					}
					delete(listableProperties, k+"_alt")
				}
				details.WriteString(fmt.Sprintf("<br/><span><a href=\"/%s\" target=\"_blank\">%s</a></span>", strVal, linkText))
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
