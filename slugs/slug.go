package slugs

import (
	"regexp"
	"strings"
)

func GenerateSlug(title string) string {
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")

	reg := regexp.MustCompile(`[^\p{L}\p{N}\-]`)
	slug = reg.ReplaceAllString(slug, "")

	return slug
}
