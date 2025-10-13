package server

import (
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

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

// Compile a regexp pattern and return a validator function that checks
// if the input matches the required pattern
func validateByPattern(pattern string) (validator func(string) bool) {
	re := regexp.MustCompile(pattern)
	return func(input string) bool {
		return re.MatchString(input)
	}
}

// Form field validators
var (
	validateName           = validateByPattern(`^[\w -]{3,50}$`)
	validateOriginURL      = validateByPattern(`^https?://[a-zA-Z0-9-]+\.[a-z]{2,5}/?$`)
	validatePluginURL      = validateByPattern(`^(https?://[a-zA-Z0-9-]+\.[a-z]{2,5}(/[a-zA-Z0-9%?=&_-]+)+)$`)
	validatePluginsURLPath = validateByPattern(`^(https?://[a-zA-Z0-9-]+\.[a-z]{2,5}(/[a-zA-Z0-9%?=&_-]+)+|(/[a-zA-Z0-9%?=&_-]+)+)$`)
)
