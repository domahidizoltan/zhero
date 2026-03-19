// Package url provides URL-related helper functions.
package url

import (
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
