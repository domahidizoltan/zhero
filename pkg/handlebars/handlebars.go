// Package handlebars is for Handlebars helpers and variables
package handlebars

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"sync"

	"github.com/aymerick/raymond"
	"github.com/domahidizoltan/zhero/pkg/logging"
	"github.com/rs/zerolog/log"
	"github.com/russross/blackfriday"
)

var (
	replacer              = strings.NewReplacer("\\n", "<br/>", "<a ", "<a target='_blank' ")
	schemaLinkPlaceholder = regexp.MustCompile(`\[\[([^]]*)\]\]`)
	schemaLinkReplacer    = strings.NewReplacer("(/docs/", "(https://schema.org/docs/")

	helpers = map[string]any{
		"concat":         concat,
		"beautify":       beautify,
		"use":            use,
		"compareAndUse":  compareAndUse,
		"htmxSortButton": htmxSortButton,
	}
)

func InitHelpers() {
	raymond.RegisterHelpers(helpers)
}

func concat(c1, c2 string) string {
	return c1 + c2
}

func beautify(text string) string {
	text = schemaLinkReplacer.Replace(text)
	text = string(blackfriday.MarkdownBasic([]byte(text)))
	text = schemaLinkPlaceholder.ReplaceAllString(text, "<a href='https://schema.org/$1'>$1</a>")
	text = replacer.Replace(text)
	return text
}

func use(use string, enabled, condition bool) string {
	if enabled && condition {
		return use
	}
	return ""
}

func compareAndUse(use string, enabled bool, optionVal, comparisonVal any) string {
	if enabled && reflect.DeepEqual(optionVal, comparisonVal) {
		return use
	}
	return ""
}

func htmxSortButton(getURL, targetID, class, label, sortField, actualQuery string) string {
	actualField, actualDir, _ := strings.Cut(actualQuery, ":")
	dir := ""

	if sortField == actualField {
		switch actualDir {
		case "asc":
			dir = "-up"
			sortField += ":desc"
		case "desc":
			dir = "-down"
			sortField += ":asc"
		}
	}
	if !strings.Contains(sortField, ":") {
		sortField += ":asc"
	}

	sort := fmt.Sprintf("<i class=\"fa-solid fa-sort%s\"></i>", dir)
	output := fmt.Sprintf("<a hx-get=\"%s&sort=%s\" hx-target=\"#%s\" hx-swap=\"innerHTML\" class=\"%s\">%s%s</a>",
		getURL, sortField, targetID, class, label, sort)
	return output
}

var (
	once         sync.Once
	apOnce       sync.Once
	absolutePath string
)

func SetAbsolutePath(path string) {
	apOnce.Do(func() {
		absolutePath = path
	})
}

func MustParse(filePath string) *raymond.Template {
	once.Do(func() {
		logging.ConfigureLogging(nil)
	})

	tpl, err := raymond.ParseFile(absolutePath + filePath)
	if err != nil {
		log.Error().
			Err(err).
			Str("file", filePath).
			Msg("failed to parse template")
	}
	return tpl
}
