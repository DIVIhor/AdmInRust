package server

import (
	"adminrust/internal/database"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func (s *Server) registerPluginDocRoutes(r chi.Router) {
	r.Route("/doc", func(r chi.Router) {
		// retrieving
		r.Get("/", s.getPluginDoc)
		// deleting
		r.Delete("/", s.deletePluginDoc)
		// adding
		r.Get("/add", s.addPluginDocForm)
		r.Post("/add", s.addPluginDoc)
		// editing
		r.Get("/edit", s.updatePluginDocForm)
		r.Post("/edit", s.updatePluginDoc)
	})
}

// Get plugin documentation
func (s *Server) getPluginDoc(w http.ResponseWriter, r *http.Request) {
	pluginSlug := r.PathValue("pluginSlug")
	doc, err := s.db.Queries().GetPluginDoc(r.Context(), pluginSlug)
	if err != nil {
		log.Println(err)
	}

	// parse HTML for clean template rendering
	cleanDoc := struct {
		PluginSlug string
		Doc        template.HTML
	}{
		pluginSlug,
		template.HTML(doc.Doc),
	}

	renderPage(w, "plugin_doc", "", cleanDoc, nil)
}

// Render page with form for adding plugin documentation
func (s *Server) addPluginDocForm(w http.ResponseWriter, r *http.Request) {
	renderPage(w, "add_plugin_doc", "Add Plugin Doc", nil, nil)
}

// Add plugin commands
func (s *Server) addPluginDoc(w http.ResponseWriter, r *http.Request) {
	// doc expected to be in HTML format
	receivedDoc := r.FormValue("doc")
	if receivedDoc == "" {
		log.Println("Empty command list")
		return
	}

	// add HTML validation
	validate = validator.New(validator.WithRequiredStructEnabled())
	err := validate.Var(receivedDoc, "html")
	if err != nil {
		log.Println("Invalid Doc input")
		return
	}

	// prepare data
	pluginSlug := r.PathValue("pluginSlug")
	doc := database.AddPluginDocParams{
		Doc:  receivedDoc,
		Slug: pluginSlug,
	}

	// save doc to DB
	_, err = s.db.Queries().AddPluginDoc(r.Context(), doc)
	if err != nil {
		log.Println(err)
		internalServerErr(w)
	}

	// redirect to a detailed plugin page
	http.Redirect(w, r, fmt.Sprintf("/plugins/%s", pluginSlug), http.StatusFound)
}

// Render form for updating plugin documentation
func (s *Server) updatePluginDocForm(w http.ResponseWriter, r *http.Request) {
	pluginSlug := r.PathValue("pluginSlug")
	// retrieve a related doc or return Not Found error
	pluginDoc, err := s.db.Queries().GetPluginDoc(r.Context(), pluginSlug)
	if err != nil {
		log.Println(err)
		notFound(w, r)
		return
	}

	renderPage(w, "add_plugin_doc", "Update Plugin Doc", pluginDoc, nil)
}

// Update plugin documentation
func (s *Server) updatePluginDoc(w http.ResponseWriter, r *http.Request) {
	// check if the retrieved form contains hidden PUT method
	if r.FormValue("_method") != "PUT" {
		log.Println("post with no PUT input")
		notAllowed(w, r)
		return
	}

	// get and validate doc
	receivedDoc := r.FormValue("doc")
	validate = validator.New(validator.WithRequiredStructEnabled())
	err := validate.Var(receivedDoc, "html")
	if err != nil {
		log.Println("Invalid Doc input")
		return
	}

	pluginSlug := r.PathValue("pluginSlug")
	// convert and save doc updates
	_, err = s.db.Queries().UpdatePluginDoc(r.Context(), database.UpdatePluginDocParams{
		Doc:  receivedDoc,
		Slug: pluginSlug,
	})
	if err != nil {
		log.Println(err)
		internalServerErr(w)
	}

	// redirect to a detailed plugin page
	http.Redirect(w, r, fmt.Sprintf("/plugins/%s", pluginSlug), http.StatusFound)
}

// Delete plugin documentation
func (s *Server) deletePluginDoc(w http.ResponseWriter, r *http.Request) {
	pluginSlug := r.PathValue("pluginSlug")

	_, err := s.db.Queries().DeletePluginDoc(r.Context(), pluginSlug)
	if err != nil {
		log.Println(err)
		internalServerErr(w)
		return
	}

	w.Header().Set("HX-Redirect", fmt.Sprintf("/plugins/%s", pluginSlug))
	w.WriteHeader(http.StatusNoContent)
}
