package server

import (
	"html/template"
	"log"
	"path/filepath"
)

// Template folder names
const templateDir = "templates/"
const templateBlocksDir = "blocks/"

// A map of name-to-template pairs for easier template calls in handlers
var templates = make(map[string]*template.Template)

// Walk through template directory to parse HTML templates,
// extend the base template, and cache them in a global map
func loadTemplates() {
	// resolve the template directory path
	absTemplateDir, err := filepath.Abs(templateDir)
	if err != nil {
		log.Fatal(err)
	}

	// template block names
	blockNames := []string{"header", "sidebar", "footer"}
	baseTemplateName := "base"

	// create a list of base template and its blocks for further filling with specific templates
	var templatePaths []string
	// base template must be the first in the list
	templatePaths = append(templatePaths, makeAbsTemplPath(absTemplateDir, baseTemplateName))
	for _, path := range blockNames {
		templatePaths = makeAbsTemplPaths(filepath.Join(absTemplateDir, templateBlocksDir), path, templatePaths)
	}

	// a list of names for specific templates
	templateNames := []string{
		"add_origin", "origin", "origins",
		"add_plugin", "plugin", "plugins",
	}
	// populate the base template with content templates and cache each one
	for _, name := range templateNames {
		// get an absolute path for content template
		tempTemplates := makeAbsTemplPaths(absTemplateDir, name, templatePaths)
		// parse templates in order base → blocks → content template
		template, err := template.ParseFiles(tempTemplates...)
		if err != nil {
			log.Fatal(err)
		}
		// cache template to global templates
		templates[name] = template
	}
}

// DRYing functions

// Make an absolute path to HTML template
func makeAbsTemplPath(absTemplateDir, fName string) (absPath string) {
	return filepath.Join(absTemplateDir, fName+".html")
}

// Append an absolute path to a list and return it
func makeAbsTemplPaths(absTemplateDir, fName string, absPathList []string) (extendedList []string) {
	return append(absPathList, makeAbsTemplPath(absTemplateDir, fName))
}
