// Package handler contains HTTP request handlers.
package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/catdevsecops/solarz-api/internal/model"
	"github.com/catdevsecops/solarz-api/internal/service"
)

// GetData returns the list of items.
func GetData(responseWriter http.ResponseWriter, _ *http.Request) {
	responseWriter.Header().Set("Content-Type", "application/json")

	items, err := service.GetData()
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		errorResponse := model.ErrorResponse{Error: "Failed to fetch items"}
		if encErr := json.NewEncoder(responseWriter).Encode(errorResponse); encErr != nil {
			// Log error if needed
			log.Printf("ERROR: %v", err)
		}
		return
	}

	responseWriter.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(responseWriter).Encode(items); err != nil {
		// Log error if needed
		log.Printf("ERROR: %v", err)
	}
}
