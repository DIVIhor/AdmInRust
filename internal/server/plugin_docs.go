package server

import (
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *Server) registerPluginDocRoutes(r chi.Router) {
	r.Route("/doc", func(r chi.Router) {
		r.Get("/", s.getPluginDoc)
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
		PluginID int64
		Doc      template.HTML
	}{
		doc.PluginID,
		template.HTML(doc.Doc),
	}

	renderPage(w, "plugin_doc", "", cleanDoc, nil)
}
