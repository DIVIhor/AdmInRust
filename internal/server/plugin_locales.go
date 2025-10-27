package server

import (
	"adminrust/internal/config"
	"adminrust/internal/database"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

// var availableLangs config.LangConfig
var availableLangs map[string]string

func (s *Server) registerPluginLocaleRoutes(r chi.Router) {
	r.Route("/loc", func(r chi.Router) {
		// retrieve all
		r.Get("/", s.getPluginLocales)
		// delete one
		r.Delete("/{lang-code:[a-zA-Z-]{2,5}}", s.deletePluginLocale)
		// add one
		r.Get("/add", s.addPluginLocaleForm)
		r.Post("/add", s.addPluginLocale)
		// edit one
		r.Get("/edit/{lang-code:[a-zA-Z-]{2,5}}", s.updatePluginLocaleForm)
		r.Post("/edit/{lang-code:[a-zA-Z-]{2,5}}", s.updatePluginLocale)
	})
}

// Retrieve all available plugin locales
func (s *Server) getPluginLocales(w http.ResponseWriter, r *http.Request) {
	pluginSlug := chi.URLParam(r, "pluginSlug")

	// get locales or 500 error
	locales, err := s.db.Queries().GetPluginLocales(r.Context(), pluginSlug)
	if err != nil {
		log.Println(err)
		internalServerErr(w)
		return
	}

	// prepare metadata
	metaData := struct {
		AddURL     string
		CurrentURL string
	}{
		AddURL:     fmt.Sprintf("%s/add", r.URL.Path),
		CurrentURL: r.URL.Path,
	}

	renderPage(w, "plugin_locales", "Plugin Locales", locales, metaData)
}

// Render form for adding plugin locale
func (s *Server) addPluginLocaleForm(w http.ResponseWriter, r *http.Request) {
	// check if languages have been loaded from config JSON
	if availableLangs == nil {
		langCfg, err := config.ReadLangs()
		if err != nil {
			log.Println(err)
			internalServerErr(w)
			return
		}
		availableLangs = map[string]string{}
		for _, lang := range langCfg {
			availableLangs[lang.Code] = lang.Name
		}
	}
	// prepare metadata
	meta := struct{ Locales map[string]string }{availableLangs}

	renderPage(w, "add_plugin_locale", "Add Plugin Locale", nil, meta)
}

// Add locale for plugin. Expects JSON input
func (s *Server) addPluginLocale(w http.ResponseWriter, r *http.Request) {
	// receive andvalidate JSON input
	recievedLocale := r.FormValue("content")
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Var(recievedLocale, "json")
	if err != nil {
		log.Printf("Error: invalid JSON: %s\n", err)
		badRequest(w)
		return
	}

	// get and prepare locale parameters for querying
	langCode := r.FormValue("lang-code")
	langName, exists := availableLangs[langCode]
	if !exists {
		log.Printf("No languages available for language code: %s\n", langCode)
		badRequest(w)
		return
	}
	pluginSlug := chi.URLParam(r, "pluginSlug")
	params := database.AddPluginLocaleParams{
		LangCode:    langCode,
		LangName:    langName,
		ContentJson: recievedLocale,
		Slug:        pluginSlug,
	}

	// write locales or 500 error
	_, err = s.db.Queries().AddPluginLocale(r.Context(), params)
	if err != nil {
		log.Printf("Error adding plugin locale: %s\n", err)
		internalServerErr(w)
		return
	}

	// redirect to a detailed plugin page
	http.Redirect(w, r, fmt.Sprintf("/plugins/%s", pluginSlug), http.StatusFound)
}

// Render plugin updating form
func (s *Server) updatePluginLocaleForm(w http.ResponseWriter, r *http.Request) {
	// get locale data
	pluginSlug := chi.URLParam(r, "pluginSlug")
	langCode := chi.URLParam(r, "lang-code")

	// check if languages have been loaded from config JSON
	if availableLangs == nil {
		langCfg, err := config.ReadLangs()
		if err != nil {
			log.Println(err)
			internalServerErr(w)
			return
		}
		availableLangs = map[string]string{}
		for _, lang := range langCfg {
			availableLangs[lang.Code] = lang.Name
		}
	}
	// validate lang code
	_, exists := availableLangs[langCode]
	if !exists {
		log.Printf("No languages available for language code: %s\n", langCode)
		badRequest(w)
		return
	}

	// prepare received parameters for querying
	params := database.GetPluginLocaleParams{
		Slug:     pluginSlug,
		LangCode: langCode,
	}

	// get plugin locale or 500 error
	locale, err := s.db.Queries().GetPluginLocale(r.Context(), params)
	if err != nil {
		log.Printf("Error getting plugin locale: %s\n", err)
		internalServerErr(w)
		return
	}

	// render plugin locale addition form
	renderPage(w, "add_plugin_locale", "Update Plugin Locale", locale, nil)
}

// Update plugin locale
func (s *Server) updatePluginLocale(w http.ResponseWriter, r *http.Request) {
	// check if the retrieved form contains hidden PUT method
	if r.FormValue("_method") != "PUT" {
		log.Println("post with no PUT input")
		notAllowed(w, r)
		return
	}

	// receive andvalidate JSON input
	recievedLocale := r.FormValue("content")
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Var(recievedLocale, "json")
	if err != nil {
		log.Printf("Error: invalid JSON: %s\n", err)
		badRequest(w)
		return
	}

	// get locale parameters
	pluginSlug := chi.URLParam(r, "pluginSlug")
	langCode := r.FormValue("lang-code")

	// check if languages have been loaded from config JSON
	if availableLangs == nil {
		langCfg, err := config.ReadLangs()
		if err != nil {
			log.Println(err)
			internalServerErr(w)
			return
		}
		availableLangs = map[string]string{}
		for _, lang := range langCfg {
			availableLangs[lang.Code] = lang.Name
		}
	}

	// validate lang code
	_, exists := availableLangs[langCode]
	if !exists {
		log.Printf("No languages available for language code: %s\n", langCode)
		badRequest(w)
		return
	}

	// prepare locale for querying
	params := database.UpdatePluginLocaleParams{
		ContentJson: recievedLocale,
		LangCode:    langCode,
		Slug:        pluginSlug,
	}

	// write locales or 500 error
	_, err = s.db.Queries().UpdatePluginLocale(r.Context(), params)
	if err != nil {
		log.Printf("Error adding plugin locale: %s\n", err)
		internalServerErr(w)
		return
	}

	// redirect to a detailed plugin page
	http.Redirect(w, r, fmt.Sprintf("/plugins/%s", pluginSlug), http.StatusFound)
}

// Delete plugin locale
func (s *Server) deletePluginLocale(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Called")
	// get and prepare locale parameters for querying
	pluginSlug := chi.URLParam(r, "pluginSlug")
	langCode := chi.URLParam(r, "lang-code")
	params := database.DeletePluginLocaleParams{
		LangCode: langCode,
		Slug:     pluginSlug,
	}
	fmt.Println(params)

	// send query
	_, err := s.db.Queries().DeletePluginLocale(r.Context(), params)
	if err != nil {
		log.Println(err)
		internalServerErr(w)
		return
	}

	// redirect to plugin page on success with HTMX
	w.Header().Set("HX-Redirect", fmt.Sprintf("/plugins/%s", pluginSlug))
	w.WriteHeader(http.StatusNoContent)
}
