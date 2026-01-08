package pagerenderer

import (
	"github.com/aymerick/raymond"
	"github.com/domahidizoltan/zhero/template"
)

type DynamicPageRenderer struct{}

func NewDynamicPageRenderer() DynamicPageRenderer {
	return DynamicPageRenderer{}
}

// TODO maybe move to pkg/handlebars
func (DynamicPageRenderer) Render(content string) (string, error) {
	return template.Index.Exec(map[string]any{"body": raymond.SafeString(content)})
}

// TODO static page renderer and preview page
