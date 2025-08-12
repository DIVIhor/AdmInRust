package server

import (
	"adminrust/internal/database"
	"fmt"
	"log"
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

		r.Route("/{pluginSlug:[a-z0-9-]+}", func(r chi.Router) {
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
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	page := Page{
		Title:   "Plugins",
		Content: plugins,
	}

	err = templates["plugins"].ExecuteTemplate(w, "base.html", page)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

// Render a detailed page for a specific plugin by its ID
func (s *Server) getPlugin(w http.ResponseWriter, r *http.Request) {
	pluginSlug := r.PathValue("pluginSlug")
	plugin, err := s.db.Queries().GetPlugin(r.Context(), pluginSlug)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "page not found", http.StatusNotFound)
		return
	}

	page := Page{
		Title:   plugin.Name,
		Content: plugin,
	}

	err = templates["plugin"].ExecuteTemplate(w, "base.html", page)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

// Render the page with adding plugin form
func (s *Server) addPluginForm(w http.ResponseWriter, r *http.Request) {
	page := Page{
		Title: "Add Plugin",
	}

	origins, err := s.db.Queries().GetOrigins(r.Context())
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if origins != nil {
		meta := struct {
			Origins []database.PluginOrigin
		}{origins}

		page.Meta = meta
	}

	err = templates["add_plugin"].ExecuteTemplate(w, "base.html", page)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

// Post the new plugin.
//
// Redirects to a detailed page for newly created plugin.
func (s *Server) addPlugin(w http.ResponseWriter, r *http.Request) {
	// since plugin plugins usually are uMod and Codefling,
	// plugin name shouldn't be less than 3 symbols long
	name := r.FormValue("name")
	nameTemplate := regexp.MustCompile(`^[\w ]{3,}$`)
	validName := nameTemplate.FindString(name)
	if validName == "" {
		log.Println("name error", name)
		http.Error(w, "Wrong name format", http.StatusBadRequest)
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
		log.Println("invalid URL:", url)
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	// origin cannot be empty and should be an integer
	originIdStr := r.FormValue("origin")
	if originIdStr == "" {
		log.Println("empty originID")
		http.Error(w, "origin ID cannot be empty", http.StatusBadRequest)
		return
	}
	originId, err := strconv.Atoi(originIdStr)
	if err != nil {
		log.Println("invalid originID:", originIdStr)
		http.Error(w, "origin ID should be a number", http.StatusBadRequest)
		return
	}

	slug := slugify(name)
	pluginParams := database.AddPluginParams{
		Name:        name,
		Slug:        slug,
		Description: descr,
		Url:         url,
		OriginID:    int64(originId),
	}
	isUpdatedOnServer := r.FormValue("isUpdatedOnServer")
	if isUpdatedOnServer == "yes" {
		pluginParams.IsUpdatedOnServer = 1
	}

	plugin, err := s.db.Queries().AddPlugin(r.Context(), pluginParams)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/plugins/%s", plugin.Slug), http.StatusFound)
}

// Update plugin details
// func (s *Server) updatePlugin(w http.ResponseWriter, r *http.Request) {
// 	just a placeholder for now
// }

// Delete plugin by its ID and redirect to the plugin list page
func (s *Server) deletePlugin(w http.ResponseWriter, r *http.Request) {
	pluginSlug := r.PathValue("pluginSlug")
	_, err := s.db.Queries().DeletePlugin(r.Context(), pluginSlug)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("redirectPath", "/plugins")
	w.WriteHeader(http.StatusNoContent)
}
