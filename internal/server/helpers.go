package server

import (
	"html/template"
	"io/fs"
	"log"
	"path/filepath"
	"strings"
)

const templateDir = "templates/"

var templates *template.Template

// Walk through template directory to parse and cache HTML templates
func loadTemplates() {
	// resolve the template directory path
	absTemplateDir, err := filepath.Abs(templateDir)
	if err != nil {
		log.Fatal(err)
	}

	var templatePaths []string

	// walk through template folder to collect all template paths
	err = filepath.Walk(absTemplateDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Fatalf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if strings.HasSuffix(info.Name(), ".html") {
			templatePaths = append(templatePaths, path)
		}
		return nil
	})
	// log and exit if can't properly walk through the directory
	if err != nil {
		log.Fatal(err)
	}

	// parse the templates
	templates, err = template.ParseFiles(templatePaths...)
	if err != nil {
		log.Fatal(err)
	}
}
