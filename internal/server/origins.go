package server

import (
	"adminrust/internal/database"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
)

// Render a list of origins
func (s *Server) getOrigins(w http.ResponseWriter, r *http.Request) {
	origins, err := s.db.Queries().GetOrigins(r.Context())
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	page := Page{
		Title:   "Origins",
		Content: origins,
	}

	err = templates.ExecuteTemplate(w, "origins.html", page)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

// Render a detailed page for a specific origin by its ID
func (s *Server) getOrigin(w http.ResponseWriter, r *http.Request) {
	originIdStr := r.PathValue("originId")
	originId, err := strconv.Atoi(originIdStr)
	if err != nil {
		http.Error(w, "origin ID should be a number", http.StatusBadRequest)
		return
	}
	origin, err := s.db.Queries().GetOrigin(r.Context(), int64(originId))
	if err != nil {
		http.Error(w, "page not found", http.StatusNotFound)
		return
	}

	page := Page{
		Title:   origin.Name,
		Content: origin,
	}

	err = templates.ExecuteTemplate(w, "origin.html", page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Render the page with adding origin form
func (s *Server) addOriginForm(w http.ResponseWriter, r *http.Request) {
	page := Page{
		Title: "Add Origin",
	}
	err := templates.ExecuteTemplate(w, "base.html", page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Post the new origin.
//
// Redirects to a detailed page for newly created origin.
func (s *Server) addOrigin(w http.ResponseWriter, r *http.Request) {
	// since plugin origins usually are uMod and Codefling,
	// origin name shouldn't be less than 3 symbols long
	name := r.FormValue("name")
	if len(name) < 3 {
		http.Error(w, "Wrong name format", http.StatusBadRequest)
		fmt.Println("name error", name)
		return
	}

	// URL must match the regex with alphanumerical format
	url := r.FormValue("url")
	urlTemplate := regexp.MustCompile("^https://[a-zA-Z0-9]+.[a-z]{1,5}$")
	validUrl := urlTemplate.FindString(url)
	if validUrl == "" {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		fmt.Println("invalid URL", url)
		return
	}

	// since path to plugin list is a URL path, it must match the regex
	// (!) perhaps the full URL should also be a valid path with further
	// processing, but not for now
	pathToPluginList := r.FormValue("pathToPluginList")
	pathTemplate := regexp.MustCompile("^/([a-zA-Z0-9/%?=&_-]+)$")
	validPath := pathTemplate.FindString(pathToPluginList)
	if validPath == "" {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		fmt.Println("invalid path to plugins", pathToPluginList)
		return
	}

	hasAPI := r.FormValue("hasApi")
	originParams := database.AddOriginParams{
		Name:             r.FormValue("name"),
		Url:              r.FormValue("url"),
		PathToPluginList: r.FormValue("pathToPluginList"),
	}
	if hasAPI == "yes" {
		originParams.HasApi = 1
	}

	origin, err := s.db.Queries().AddOrigin(r.Context(), originParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}

	fmt.Println(origin) // delete this in future

	http.Redirect(w, r, fmt.Sprintf("/origins/%d", origin.ID), http.StatusFound)
}

// Update origin details
// func (s *Server) updateOrigin(w http.ResponseWriter, r *http.Request) {
// 	just a placeholder for now
// }

// Delete origin by its ID and redirect to the origin list page
func (s *Server) deleteOrigin(w http.ResponseWriter, r *http.Request) {
	originIdStr := r.PathValue("originId")
	originId, err := strconv.Atoi(originIdStr)
	if err != nil {
		http.Error(w, "origin ID should be a number", http.StatusBadRequest)
		return
	}

	_, err = s.db.Queries().DeleteOrigin(r.Context(), int64(originId))
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("redirectPath", "/origins")
	w.WriteHeader(http.StatusNoContent)
}
