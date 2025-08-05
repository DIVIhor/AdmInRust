package server

import (
	"adminrust/internal/database"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// Routes to get one/many, add, and delete plugins
func (s *Server) registerPluginRoutes(r *chi.Mux) {
	r.Route("/plugins", func(r chi.Router) {
		r.Get("/", s.getPlugins)

		r.Get("/add", s.addPluginForm)
		r.Post("/add", s.addPlugin)

		r.Route("/{pluginId:[0-9]+}", func(r chi.Router) {
			r.Get("/", s.getPlugin)
			// r.Put("/", s.updatePlugin)
			r.Delete("/", s.deletePlugin)
		})
	})
}

// Render a list of plugins
func (s *Server) getPlugins(w http.ResponseWriter, r *http.Request) {
	plugins, err := s.db.Queries().GetPlugins(r.Context())
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	page := Page{
		Title:   "Plugins",
		Content: plugins,
	}

	err = templates.ExecuteTemplate(w, "plugins.html", page)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

// Render a detailed page for a specific plugin by its ID
func (s *Server) getPlugin(w http.ResponseWriter, r *http.Request) {
	pluginIdStr := r.PathValue("pluginId")
	pluginId, err := strconv.Atoi(pluginIdStr)
	if err != nil {
		http.Error(w, "plugin ID should be a number", http.StatusBadRequest)
		return
	}
	plugin, err := s.db.Queries().GetPlugin(r.Context(), int64(pluginId))
	if err != nil {
		http.Error(w, "page not found", http.StatusNotFound)
		return
	}

	page := Page{
		Title:   plugin.Name,
		Content: plugin,
	}

	err = templates.ExecuteTemplate(w, "plugin.html", page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Render the page with adding plugin form
func (s *Server) addPluginForm(w http.ResponseWriter, r *http.Request) {
	page := Page{
		Title: "Add Plugin",
	}

	origins, err := s.db.Queries().GetOrigins(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if origins != nil {
		meta := struct {
			Origins []database.PluginOrigin
		}{origins}

		page.Meta = meta
	}

	err = templates.ExecuteTemplate(w, "add_plugin.html", page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Post the new plugin.
//
// Redirects to a detailed page for newly created plugin.
func (s *Server) addPlugin(w http.ResponseWriter, r *http.Request) {
	// since plugin plugins usually are uMod and Codefling,
	// plugin name shouldn't be less than 3 symbols long
	name := r.FormValue("name")
	if len(name) < 3 {
		http.Error(w, "Wrong name format", http.StatusBadRequest)
		fmt.Println("name error", name)
		return
	}

	// description can be empty since Codefling doesn't provide any
	// for future improvements this should have some limit
	descr := r.FormValue("description")

	// URL must match the regex with alphanumerical format
	url := r.FormValue("url")
	urlTemplate := regexp.MustCompile("^https://[a-zA-Z0-9]+.[a-z]{1,5}/([a-zA-Z0-9/%?=&_-]+)$")
	validUrl := urlTemplate.FindString(url)
	if validUrl == "" {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		fmt.Println("invalid URL", url)
		return
	}

	// origin cannot be empty and should be an integer
	originIdStr := r.FormValue("origin")
	if originIdStr == "" {
		http.Error(w, "origin ID cannot be empty", http.StatusBadRequest)
		return
	}
	originId, err := strconv.Atoi(originIdStr)
	if err != nil {
		http.Error(w, "origin ID should be a number", http.StatusBadRequest)
		return
	}

	isUpdatedOnServer := r.FormValue("isUpdatedOnServer")
	pluginParams := database.AddPluginParams{
		Name:        name,
		Description: descr,
		Url:         url,
		OriginID:    int64(originId),
	}
	if isUpdatedOnServer == "yes" {
		pluginParams.IsUpdatedOnServer = 1
	}

	plugin, err := s.db.Queries().AddPlugin(r.Context(), pluginParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/plugins/%d", plugin.ID), http.StatusFound)
}

// Update plugin details
// func (s *Server) updatePlugin(w http.ResponseWriter, r *http.Request) {
// 	just a placeholder for now
// }

// Delete plugin by its ID and redirect to the plugin list page
func (s *Server) deletePlugin(w http.ResponseWriter, r *http.Request) {
	pluginIdStr := r.PathValue("pluginId")
	pluginId, err := strconv.Atoi(pluginIdStr)
	if err != nil {
		http.Error(w, "plugin ID should be a number", http.StatusBadRequest)
		return
	}

	_, err = s.db.Queries().DeletePlugin(r.Context(), int64(pluginId))
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("redirectPath", "/plugins")
	w.WriteHeader(http.StatusNoContent)
}
