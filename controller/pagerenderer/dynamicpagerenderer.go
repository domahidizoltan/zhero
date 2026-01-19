package pagerenderer

import (
	"fmt"
	"strings"

	"github.com/aymerick/raymond"
	"github.com/domahidizoltan/zhero/domain/schema"
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
		if prop.Name == meta.SecondaryIdentifier {
			b.WriteString(fmt.Sprintf("<h1>%s</h1>", v))
			continue
		}

		b.WriteString(fmt.Sprintf("<p>%s</p>", v))
	}
	return template.Index.Exec(map[string]any{"body": raymond.SafeString(b.String())})
}

// TODO static page renderer and preview page
