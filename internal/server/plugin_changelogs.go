package server

import (
	"log"
	"net/http"
)

// Get plugin version info
func (s *Server) getPluginChangelog(w http.ResponseWriter, r *http.Request) {
	pluginSlug := r.PathValue("pluginSlug")
	changelog, err := s.db.Queries().GetPluginChangelog(r.Context(), pluginSlug)
	if err != nil {
		log.Println(err)
		internalServerErr(w)
		return
	}

	renderPage(w, "plugin_changelogs", "", changelog, nil)
}
