// internal/server/server.go
package server

import (
	"net/http"

	handler "github.com/elangreza/e-commerce/api/internal/handlers"
)

type Server struct {
	httpServer *http.Server
}

func New(address string, handler *handler.ProductHandler) *Server {

	// Setup routes
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/products", handler.ListProducts())

	httpServer := &http.Server{
		Addr:    address,
		Handler: mux,
	}

	return &Server{httpServer: httpServer}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}
