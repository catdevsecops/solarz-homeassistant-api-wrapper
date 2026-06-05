package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/catdevsecops/solarz-api/internal/handler"
)

// ServerConfig contém a configuração do servidor
type ServerConfig struct {
	Addr    string
	Handler http.Handler
}

// startServer inicia o servidor (extraído de main para testabilidade)
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

// setupRouter cria e configura o router (extraído de main para testabilidade)
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

	startServer(config)
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "ok"}); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}
