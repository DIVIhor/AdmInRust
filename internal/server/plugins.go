package server

import (
	"adminrust/internal/database"
	"encoding/json"
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

		r.Route("/edit/{pluginSlug:[a-z0-9-]+}", func(r chi.Router) {
			r.Get("/", s.updatePluginForm)
			r.Post("/", s.updatePlugin)
		})

		r.Route("/{pluginSlug:[a-z0-9-]+}", func(r chi.Router) {
			r.Get("/", s.getPlugin)
			r.Delete("/", s.deletePlugin)
		})
	})
}

// Render a list of plugins
func (s *Server) getPlugins(w http.ResponseWriter, r *http.Request) {
	plugins, err := s.db.Queries().GetPlugins(r.Context())
	if err != nil {
		log.Println(err)
		internalServerErr(w)
		return
	}

	// render plugins page
	renderPage(w, "plugins", "Plugins", plugins, nil)
}

// Render a detailed page for a specific plugin by its ID
func (s *Server) getPlugin(w http.ResponseWriter, r *http.Request) {
	pluginSlug := r.PathValue("pluginSlug")
	plugin, err := s.db.Queries().GetPlugin(r.Context(), pluginSlug)
	if err != nil {
		log.Println(err)
		http.Error(w, "page not found", http.StatusNotFound)
		return
	}

	// populate and render detailed origin page
	renderPage(w, "plugin", plugin.Name, plugin, nil)
}

// Render the page with plugin addition form
func (s *Server) addPluginForm(w http.ResponseWriter, r *http.Request) {
	// get available origins to use as meta data in form
	origins, err := s.db.Queries().GetOrigins(r.Context())
	if err != nil {
		log.Println(err)
		internalServerErr(w)
		return
	}
	meta := struct{ Origins []database.PluginOrigin }{}
	if origins != nil {
		meta = struct{ Origins []database.PluginOrigin }{origins}
	}

	// populate and render plugin addition form
	renderPage(w, "add_plugin", "Add Plugin", nil, meta)
}

// Post a new plugin.
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
	urlTemplate := regexp.MustCompile("^https://[a-zA-Z0-9-]+.[a-z]{1,5}/([a-zA-Z0-9/%?=&_-]+)$")
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
		internalServerErr(w)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/plugins/%s", plugin.Slug), http.StatusFound)
}

// Render a plugin updating form
func (s *Server) updatePluginForm(w http.ResponseWriter, r *http.Request) {
	pluginSlug := r.PathValue("pluginSlug")

	// get plugin with whole list of origins from DB
	pluginWithOrigins, err := s.db.Queries().GetPluginWithOriginsJson(r.Context(), pluginSlug)
	if err != nil {
		log.Println(err)
		http.Error(w, "page not found", http.StatusNotFound)
		return
	}

	// extract origins from JSON string
	type origin struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}
	origins := []origin{}
	rawStr := pluginWithOrigins.OriginOptions.(string)
	rawJSON := []byte(rawStr)
	err = json.Unmarshal([]byte(rawJSON), &origins)
	if err != nil {
		log.Println(err)
		internalServerErr(w)
		return
	}
	// use available origins as meta data in form
	meta := struct{ Origins []origin }{}
	if origins != nil {
		meta = struct{ Origins []origin }{origins}
	}

	// convert plugin with origins to a proper Plugin struct
	plugin := database.Plugin{
		ID:                pluginWithOrigins.ID,
		Name:              pluginWithOrigins.Name,
		Slug:              pluginWithOrigins.Slug,
		Description:       pluginWithOrigins.Description,
		Url:               pluginWithOrigins.Url,
		OriginID:          pluginWithOrigins.OriginID,
		IsUpdatedOnServer: pluginWithOrigins.IsUpdatedOnServer,
		CreatedAt:         pluginWithOrigins.CreatedAt,
		UpdatedAt:         pluginWithOrigins.UpdatedAt,
	}

	// populate and render plugin updating form
	renderPage(w, "add_plugin", "Update Plugin", plugin, meta)
}

// Update plugin details
func (s *Server) updatePlugin(w http.ResponseWriter, r *http.Request) {
	// check if the retrieved form contains hidden PUT method
	if r.FormValue("_method") != "PUT" {
		log.Println("post with no PUT input")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	pluginSlug := r.PathValue("pluginSlug")

	// description can be empty since Codefling doesn't provide any
	// for future improvements this should have some limit
	descr := r.FormValue("description")

	// URL must match the regex with alphanumerical format
	url := r.FormValue("url")
	urlTemplate := regexp.MustCompile("^https://[a-zA-Z0-9-]+.[a-z]{1,5}/([a-zA-Z0-9/%?=&_-]+)$")
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

	// prepare data for updating the plugin in DB
	updPluginParams := database.UpdatePluginParams{
		Description: descr,
		Url:         url,
		OriginID:    int64(originId),
		Slug:        pluginSlug,
	}
	isUpdatedOnServer := r.FormValue("isUpdatedOnServer")
	if isUpdatedOnServer == "yes" {
		updPluginParams.IsUpdatedOnServer = 1
	}

	// update the plugin in DB
	plugin, err := s.db.Queries().UpdatePlugin(r.Context(), updPluginParams)
	if err != nil {
		log.Println(err)
		internalServerErr(w)
		return
	}

	// redirect to a plugin detailed page
	http.Redirect(w, r, fmt.Sprintf("/plugins/%s", plugin.Slug), http.StatusFound)
}

// Delete plugin by its ID and redirect to the plugin list page
func (s *Server) deletePlugin(w http.ResponseWriter, r *http.Request) {
	pluginSlug := r.PathValue("pluginSlug")
	_, err := s.db.Queries().DeletePlugin(r.Context(), pluginSlug)
	if err != nil {
		log.Println(err)
		internalServerErr(w)
		return
	}

	w.Header().Set("redirectPath", "/plugins")
	w.WriteHeader(http.StatusNoContent)
}
