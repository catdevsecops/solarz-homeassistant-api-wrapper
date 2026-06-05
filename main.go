// Package main is the entry point for the Solarz API server.
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/catdevsecops/solarz-api/internal/handler"
)

// ServerConfig contains the server configuration.
type ServerConfig struct {
	Addr    string
	Handler http.Handler
}

// startServer starts the server (extracted from main for testability).
func startServer(config ServerConfig) error {
	server := &http.Server{
		Addr:         config.Addr,
		Handler:      config.Handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Server starting on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
		return err
	}

	return nil
}

// setupRouter creates and configures the router (extracted from main for testability).
func setupRouter() *http.ServeMux {
	router := http.NewServeMux()

	// Registrar rotas
	router.HandleFunc("GET /api/v1/data", handler.GetData)

	// Health check
	router.HandleFunc("GET /health", healthHandler)

	return router
}

func main() {
	router := setupRouter()

	config := ServerConfig{
		Addr:    ":8080",
		Handler: router,
	}

	if err := startServer(config); err != nil {
		log.Printf("Failed to start server: %v", err)
	}
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "ok"}); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}
