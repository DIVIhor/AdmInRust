package server

import (
	"adminrust/internal/database"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/go-chi/chi/v5"
)

// Routes to get one/many, add, and delete origins
func (s *Server) registerOriginRoutes(r *chi.Mux) {
	r.Route("/origins", func(r chi.Router) {
		r.Get("/", s.getOrigins)

		r.Get("/add", s.addOriginForm)
		r.Post("/add", s.addOrigin)

		r.Route("/edit/{originSlug:[a-z0-9-]+}", func(r chi.Router) {
			r.Get("/", s.updateOriginForm)
			r.Post("/", s.updateOrigin)
		})

		r.Route("/{originSlug:[a-z0-9-]+}", func(r chi.Router) {
			r.Get("/", s.getOrigin)
			r.Delete("/", s.deleteOrigin)
		})
	})
}

// Get from DB and render a list of origins
func (s *Server) getOrigins(w http.ResponseWriter, r *http.Request) {
	origins, err := s.db.Queries().GetOrigins(r.Context())
	if err != nil {
		log.Println(err)
		internalServerErr(w)
		return
	}

	// populate and render origins page
	renderPage(w, "origins", "Origins", origins, nil)
}

// Render a detailed page for a specific origin by its ID
func (s *Server) getOrigin(w http.ResponseWriter, r *http.Request) {
	originSlug := r.PathValue("originSlug")
	origin, err := s.db.Queries().GetOrigin(r.Context(), originSlug)
	if err != nil {
		log.Println(err)
		http.Error(w, "page not found", http.StatusNotFound)
		return
	}

	// populate and render detailed origin page
	renderPage(w, "origin", origin.Name, origin, nil)
}

// Render the page with origin addition form
func (s *Server) addOriginForm(w http.ResponseWriter, r *http.Request) {
	// populate and render origin addition form
	renderPage(w, "add_origin", "Add Origin", nil, nil)
}

// Post a new origin.
//
// Redirects to a detailed page for newly created origin.
func (s *Server) addOrigin(w http.ResponseWriter, r *http.Request) {
	// since plugin origins usually are uMod and Codefling,
	// origin name shouldn't be less than 3 symbols long
	name := r.FormValue("name")
	nameTemplate := regexp.MustCompile(`^[\w ]{3,}$`)
	validName := nameTemplate.FindString(name)
	if validName == "" {
		log.Println("name error:", name)
		http.Error(w, "Wrong name format", http.StatusBadRequest)
		return
	}

	// URL must match the regex with alphanumerical format
	url := r.FormValue("url")
	urlTemplate := regexp.MustCompile("^https://[a-zA-Z0-9-]+.[a-z]{1,5}$")
	validUrl := urlTemplate.FindString(url)
	if validUrl == "" {
		log.Println("invalid URL:", url)
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	// since path to plugin list is a URL path, it must match the regex
	// (!) perhaps the full URL should also be a valid path with further
	// processing, but not for now
	pathToPluginList := r.FormValue("pathToPluginList")
	pathTemplate := regexp.MustCompile("^/([a-zA-Z0-9/%?=&_-]+)$")
	validPath := pathTemplate.FindString(pathToPluginList)
	if validPath == "" {
		log.Println("invalid path to plugins:", pathToPluginList)
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	slug := slugify(validName)
	originParams := database.AddOriginParams{
		Name:             validName,
		Slug:             slug,
		Url:              url,
		PathToPluginList: pathToPluginList,
	}
	hasAPI := r.FormValue("hasApi")
	if hasAPI == "yes" {
		originParams.HasApi = 1
	}

	origin, err := s.db.Queries().AddOrigin(r.Context(), originParams)
	if err != nil {
		log.Println(err)
		internalServerErr(w)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/origins/%s", origin.Slug), http.StatusFound)
}

// Render an origin updating form
func (s *Server) updateOriginForm(w http.ResponseWriter, r *http.Request) {
	originSlug := r.PathValue("originSlug")
	// get origin data from DB
	origin, err := s.db.Queries().GetOrigin(r.Context(), originSlug)
	if err != nil {
		log.Println(err)
		http.Error(w, "page not found", http.StatusNotFound)
		return
	}

	// populate and render origin updating form
	renderPage(w, "add_origin", "Update Origin", origin, nil)
}

// Update origin details
func (s *Server) updateOrigin(w http.ResponseWriter, r *http.Request) {
	// check if the retrieved form contains hidden PUT method
	if r.FormValue("_method") != "PUT" {
		log.Println("post with no PUT input")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	originSlug := r.PathValue("originSlug")

	// URL must match the regex with alphanumerical format
	url := r.FormValue("url")
	urlTemplate := regexp.MustCompile("^https://[a-zA-Z0-9-]+.[a-z]{1,5}$")
	validUrl := urlTemplate.FindString(url)
	if validUrl == "" {
		log.Println("invalid URL:", url)
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	// since path to plugin list is a URL path, it must match the regex
	// (!) perhaps the full URL should also be a valid path with further
	// processing, but not for now
	pathToPluginList := r.FormValue("pathToPluginList")
	pathTemplate := regexp.MustCompile("^/([a-zA-Z0-9/%?=&_-]+)$")
	validPath := pathTemplate.FindString(pathToPluginList)
	if validPath == "" {
		log.Println("invalid path to plugins:", pathToPluginList)
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	// prepare data for updating the origin in DB
	updOriginParams := database.UpdateOriginParams{
		Url:              url,
		PathToPluginList: pathToPluginList,
		Slug:             originSlug,
	}
	hasAPI := r.FormValue("hasApi")
	if hasAPI == "yes" {
		updOriginParams.HasApi = 1
	}

	// update the origin in DB
	origin, err := s.db.Queries().UpdateOrigin(r.Context(), updOriginParams)
	if err != nil {
		log.Println(err)
		internalServerErr(w)
		return
	}

	// redirect to an origin detailed page
	http.Redirect(w, r, fmt.Sprintf("/origins/%s", origin.Slug), http.StatusFound)
}

// Delete origin by its ID and redirect to the origin list page
func (s *Server) deleteOrigin(w http.ResponseWriter, r *http.Request) {
	originSlug := r.PathValue("originSlug")
	_, err := s.db.Queries().DeleteOrigin(r.Context(), originSlug)
	if err != nil {
		log.Println(err)
		internalServerErr(w)
		return
	}

	w.Header().Set("redirectPath", "/origins")
	w.WriteHeader(http.StatusNoContent)
}
