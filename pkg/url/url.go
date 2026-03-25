// Package url provides URL-related helper functions.
package url

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

var (
	validSlugCharactersRegex = regexp.MustCompile(`[^a-z0-9]+`)
	repeatingHyphensRegex    = regexp.MustCompile(`-+`)
	repeatingSlashesRegex    = regexp.MustCompile(`\/+`)
)

func Slugify(s string) string {
	s = strings.ToLower(s)
	s = validSlugCharactersRegex.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	s = repeatingHyphensRegex.ReplaceAllString(s, "-")
	s = repeatingSlashesRegex.ReplaceAllString(s, "/")
	return s
}

func Canonical(req *http.Request) string {
	if req == nil {
		return "http://localhost"
	}

	scheme := "http"
	if req.TLS != nil || req.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}

	return fmt.Sprintf("%s://%s%s", scheme, req.Host, req.URL.Path)
}
