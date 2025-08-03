package server

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

const templateDir = "templates/"

var templatePaths = getTemplatePaths(
	"blocks/footer", "blocks/header", "blocks/sidebar",
	"base", "add_origin", "origin", "origins",
)

// Parsing files at app initialization to avoid reading files on each request
var templates = template.Must(template.ParseFiles(templatePaths...))

// DRYing template paths
func getTemplatePaths(templateNames ...string) (paths []string) {
	paths = make([]string, len(templateNames))
	for idx, tmpltName := range templateNames {
		paths[idx] = templateDir + tmpltName + ".html"
	}
	return paths
}

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/", s.mainHandler) //s.HelloWorldHandler)

	// origin-related routes
	s.registerOriginRoutes(r)

	r.Get("/health", s.healthHandler)

	return r
}

type Page struct {
	Title   string
	Content any
	// TODO add URL parameter
}

func (s *Server) mainHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "base.html", "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, _ := json.Marshal(s.db.Health())
	_, _ = w.Write(jsonResp)
}
