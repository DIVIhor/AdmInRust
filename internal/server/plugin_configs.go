package server

import (
	"adminrust/internal/database"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

// config routing
func (s *Server) registerPluginCfgRoutes(r chi.Router) {
	r.Route("/config", func(r chi.Router) {
		// retrieving
		r.Get("/", s.getPluginCfg)
		// deleting
		r.Delete("/", s.deletePluginCfg)
		// adding
		r.Get("/add", s.addPluginCfgForm)
		r.Post("/add", s.addPluginCfg)
		// editing
		r.Get("/edit", s.updatePluginCfgForm)
		r.Post("/edit", s.updatePluginCfg)
	})
}

// Show plugin configuration
func (s *Server) getPluginCfg(w http.ResponseWriter, r *http.Request) {
	pluginSlug := r.PathValue("pluginSlug")
	// get plugin configuration data
	configData, err := s.db.Queries().GetPluginConfig(r.Context(), pluginSlug)
	if err != nil {
		log.Println(err)
	}

	metaData := struct {
		AddURL     string
		CurrentURL string
		EditURL    string
	}{
		AddURL:     fmt.Sprintf("%s/add", r.URL.Path),
		CurrentURL: r.URL.Path,
		EditURL:    fmt.Sprintf("%s/edit", r.URL.Path),
	}

	// render populated page
	renderPage(w, "plugin_cfg", "Plugin Configuration", configData, metaData)
}

// Render form for adding plugin configuration
func (s *Server) addPluginCfgForm(w http.ResponseWriter, r *http.Request) {
	renderPage(w, "add_plugin_cfg", "Add Plugin Configuration", nil, nil)
}

// Add configuratoin for plugin. Expects JSON input
func (s *Server) addPluginCfg(w http.ResponseWriter, r *http.Request) {
	// validate JSON input
	receivedCfg := r.FormValue("config")
	validate = validator.New(validator.WithRequiredStructEnabled())
	err := validate.Var(receivedCfg, "json")
	if err != nil {
		log.Println("Invalid configuration input")
		return
	}

	pluginSlug := r.PathValue("pluginSlug")
	// prepare data and save configuration or Internal Server Error
	config := database.AddPluginConfigParams{
		ConfigJson: receivedCfg,
		Slug:       pluginSlug,
	}
	_, err = s.db.Queries().AddPluginConfig(r.Context(), config)
	if err != nil {
		log.Println(err)
		internalServerErr(w)
		return
	}

	// redirect to plugin page
	http.Redirect(w, r, fmt.Sprintf("/plugins/%s", pluginSlug), http.StatusFound)
}

// Render form for updating plugin configuration
func (s *Server) updatePluginCfgForm(w http.ResponseWriter, r *http.Request) {
	pluginSlug := r.PathValue("pluginSlug")

	// get plugin configuration
	config, err := s.db.Queries().GetPluginConfig(r.Context(), pluginSlug)
	if err != nil {
		log.Println(err)
		internalServerErr(w)
		return
	}

	// show pre-populated form
	renderPage(w, "add_plugin_cfg", "Update Plugin Configuration", config, nil)
}

// Update plugin configuration
func (s *Server) updatePluginCfg(w http.ResponseWriter, r *http.Request) {
	// check if the retrieved form contains hidden PUT method
	if r.FormValue("_method") != "PUT" {
		log.Println("post with no PUT input")
		notAllowed(w, r)
		return
	}

	// validate JSON input
	receivedCfg := r.FormValue("config")
	validate = validator.New(validator.WithRequiredStructEnabled())
	err := validate.Var(receivedCfg, "json")
	if err != nil {
		log.Println("Invalid configuration input")
		return
	}

	pluginSlug := r.PathValue("pluginSlug")
	// prepare data and save configuration or Internal Server Error
	config := database.UpdatePluginConfigParams{
		ConfigJson: receivedCfg,
		Slug:       pluginSlug,
	}
	_, err = s.db.Queries().UpdatePluginConfig(r.Context(), config)
	if err != nil {
		log.Println(err)
		internalServerErr(w)
		return
	}

	// redirect to plugin page
	http.Redirect(w, r, fmt.Sprintf("/plugins/%s", pluginSlug), http.StatusFound)
}

// Delete plugin configuration
func (s *Server) deletePluginCfg(w http.ResponseWriter, r *http.Request) {
	pluginSlug := r.PathValue("pluginSlug")

	_, err := s.db.Queries().DeletePluginConfig(r.Context(), pluginSlug)
	if err != nil {
		log.Println(err)
		internalServerErr(w)
		return
	}

	// redirect to plugin page on success with HTMX
	w.Header().Set("HX-Redirect", fmt.Sprintf("/plugins/%s", pluginSlug))
	w.WriteHeader(http.StatusNoContent)
}
