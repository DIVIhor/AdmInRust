package server

import (
	"log"
	"net/http"
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

// Prepare, populate, and render page with many entries
// or return Internal Server Error
func renderPage(w http.ResponseWriter, tmpltName, pageTitle string, pageContent, pageMeta any) {
	// prepare data for template population
	page := Page{
		Title: pageTitle,
	}
	if pageContent != nil {
		page.Content = pageContent
	}
	if pageMeta != nil {
		page.Meta = pageMeta
	}

	// populate and render template or return HTTP 500
	err := templates[tmpltName].Execute(w, page)
	if err != nil {
		log.Println(err)
		internalServerErr(w)
	}
}
