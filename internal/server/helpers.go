package server

import (
	"regexp"
	"strings"
)

// Convert string to appropriate slug
//
// • convert string to lower case
//
// • trim spaces around string
//
// • replace spaces with hyphens
//
// • cut special symbols
func slugify(name string) (slug string) {
	// NOTE accent translation (e.g. 'à' → 'a') may be useful here,
	// but since names should come from already filtered resources,
	// for now it's unnecessary

	slug = strings.ToLower(name)
	// replace symbols that not match lowercase alphanumerical format with hyphens
	slug = regexp.MustCompile("[^a-z0-9]+").ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")

	return slug
}
