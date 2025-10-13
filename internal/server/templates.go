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
	var baseTemplatePaths []string
	// base template must be the first in the list
	baseTemplatePaths = append(baseTemplatePaths, makeAbsTemplPath(absTemplateDir, baseTemplateName))
	for _, path := range blockNames {
		baseTemplatePaths = makeAbsTemplPaths(filepath.Join(absTemplateDir, templateBlocksDir), path, baseTemplatePaths)
	}

	// a list of names for specific templates
	templateNames := []string{
		"add_origin", "origin", "origins",
		"add_plugin", "plugin", "plugins",
		"add_plugin_cmds",
		"add_plugin_doc",
		"http_error",
	}
	// populate the base template with content templates and cache each one
	for _, name := range templateNames {
		// get an absolute path for content template
		tempTemplates := makeAbsTemplPaths(absTemplateDir, name, baseTemplatePaths)
		// parse templates in order base → blocks → content template
		tmplt := template.Must(template.ParseFiles(tempTemplates...))
		// cache template to global templates
		templates[name] = tmplt
	}

	// process templates for inner-page tabs
	tabTemplateNames := []string{
		"plugin_changelogs", "plugin_commands",
		"plugin_doc",
	}
	for _, tabTempl := range tabTemplateNames {
		absPath := makeAbsTemplPath(absTemplateDir, tabTempl)
		templates[tabTempl] = template.Must(template.ParseFiles(absPath))
	}
}

// Make an absolute path to HTML template
func makeAbsTemplPath(absTemplateDir, fName string) (absPath string) {
	return filepath.Join(absTemplateDir, fName+".html")
}

// Append an absolute path to a list and return it
func makeAbsTemplPaths(absTemplateDir, fName string, absPathList []string) (extendedList []string) {
	return append(absPathList, makeAbsTemplPath(absTemplateDir, fName))
}
