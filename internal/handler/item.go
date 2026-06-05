package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/catdevsecops/solarz-api/internal/model"
	"github.com/catdevsecops/solarz-api/internal/service"
)

// GetItems retorna a lista de itens
func GetData(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	items, err := service.GetData()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if encErr := json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Failed to fetch items"}); encErr != nil {
			// Log error if needed
			log.Printf("ERROR: %v", err)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(items); err != nil {
		// Log error if needed
		log.Printf("ERROR: %v", err)
	}
}
