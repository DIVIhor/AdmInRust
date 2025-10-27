package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Page struct {
	Title   string
	Content any
	Meta    any
	// TODO add URL parameter
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

	// r.Get("/", s.mainHandler) //s.HelloWorldHandler)

	// origin-related routes
	s.registerOriginRoutes(r)

	// plugin-related routes
	s.registerPluginRoutes(r)

	r.Get("/health", s.healthHandler)

	r.NotFound(notFound)
	r.MethodNotAllowed(notAllowed)

	return r
}

// func (s *Server) mainHandler(w http.ResponseWriter, r *http.Request) {
// 	err := templates.ExecuteTemplate(w, "base.html", "")
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// }

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

// Write HTTP error status code in header then render error template
func errorHandler(w http.ResponseWriter, httpErrCode int, httpErr string) {
	page := Page{
		Title:   fmt.Sprintf("(%d) %s", httpErrCode, httpErr),
		Content: struct{ Error string }{httpErr},
	}

	w.WriteHeader(httpErrCode)
	err := templates["http_error"].Execute(w, page)
	if err != nil {
		http.Error(w, httpErr, httpErrCode)
	}
}

// HTTP 400 handler
func badRequest(w http.ResponseWriter) {
	errorHandler(w, http.StatusBadRequest, "Bad request")
}

// HTTP 404 handler
func notFound(w http.ResponseWriter, r *http.Request) {
	errorHandler(w, http.StatusNotFound, "Page not found")
}

// HTTP 405 handler
func notAllowed(w http.ResponseWriter, r *http.Request) {
	errorHandler(w, http.StatusMethodNotAllowed, "Not allowed")
}

// HTTP 500 handler
func internalServerErr(w http.ResponseWriter) {
	errorHandler(w, http.StatusInternalServerError, "Internal server error")
}
