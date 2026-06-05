package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/catdevsecops/solarz-api/internal/model"
)

func TestGetData_EmptyEndpoint(t *testing.T) {
	t.Run("returns 200 OK with empty data when SOLARZ_ENDPOINT is not set", func(t *testing.T) {
		t.Setenv("SOLARZ_ENDPOINT", "")

		req := httptest.NewRequest("GET", "/api/v1/data", nil)
		w := httptest.NewRecorder()

		GetData(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("GetData() returned status %d, want %d", w.Code, http.StatusOK)
		}

		contentType := w.Header().Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Content-Type = %q, want 'application/json'", contentType)
		}

		var result []model.Item
		if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
			t.Errorf("failed to decode response: %v", err)
		}

		if len(result) != 0 {
			t.Errorf("expected empty array, got %d items", len(result))
		}
	})
}

func TestGetData_ValidData(t *testing.T) {
	t.Run("returns 200 OK with item data from valid API", func(t *testing.T) {
		// Mock API server
		apiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response := model.SolarzResponse{
				Dados: []model.DadoGeracao{
					{
						Data:        "2026-06-01",
						Quantidade:  12.6,
						Prognostico: 19.79,
						Manual:      false,
						UsinaId:     487759,
						Denominacao: "(3633) Clayton - Mogi",
					},
					{
						Data:        "2026-06-04",
						Quantidade:  25.5,
						Prognostico: 19.79,
						Manual:      false,
						UsinaId:     487759,
						Denominacao: "(3633) Clayton - Mogi",
					},
				},
				TotalGerado:      75.4,
				TotalPrognostico: 138.53,
				Desempenho:       54.42,
			}
			if err := json.NewEncoder(w).Encode(response); err != nil {
				t.Logf("failed to encode response: %v", err)
			}
		}))
		defer apiServer.Close()

		t.Setenv("SOLARZ_ENDPOINT", apiServer.URL)

		req := httptest.NewRequest("GET", "/api/v1/data", nil)
		w := httptest.NewRecorder()

		GetData(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("GetData() returned status %d, want %d", w.Code, http.StatusOK)
		}

		contentType := w.Header().Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Content-Type = %q, want 'application/json'", contentType)
		}

		var result []model.Item
		if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
			t.Errorf("failed to decode response: %v", err)
		}

		if len(result) != 1 {
			t.Errorf("expected 1 item, got %d", len(result))
		}

		if result[0].ID != "2026-06-04" {
			t.Errorf("result[0].ID = %q, want '2026-06-04'", result[0].ID)
		}

		if result[0].Value != "25.50" {
			t.Errorf("result[0].Value = %q, want '25.50'", result[0].Value)
		}
	})
}

func TestGetData_APIError(t *testing.T) {
	t.Run("returns 500 Internal Server Error when API request fails", func(t *testing.T) {
		t.Setenv("SOLARZ_ENDPOINT", "http://invalid-url-that-does-not-exist:9999")

		req := httptest.NewRequest("GET", "/api/v1/data", nil)
		w := httptest.NewRecorder()

		GetData(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("GetData() returned status %d, want %d", w.Code, http.StatusInternalServerError)
		}

		var result model.ErrorResponse
		if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
			t.Errorf("failed to decode response: %v", err)
		}

		if result.Error == "" {
			t.Errorf("expected error message, got empty string")
		}
	})
}

func TestGetData_InvalidJSON(t *testing.T) {
	t.Run("returns 500 Internal Server Error when API returns invalid JSON", func(t *testing.T) {
		apiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte("not valid json"))
		}))
		defer apiServer.Close()

		t.Setenv("SOLARZ_ENDPOINT", apiServer.URL)

		req := httptest.NewRequest("GET", "/api/v1/data", nil)
		w := httptest.NewRecorder()

		GetData(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("GetData() returned status %d, want %d", w.Code, http.StatusInternalServerError)
		}

		var result model.ErrorResponse
		if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
			t.Errorf("failed to decode response: %v", err)
		}

		if result.Error != "Failed to fetch items" {
			t.Errorf("result.Error = %q, want 'Failed to fetch items'", result.Error)
		}
	})
}

func TestGetData_HeaderValidation(t *testing.T) {
	t.Run("validates Content-Type header is set before WriteHeader", func(t *testing.T) {
		t.Setenv("SOLARZ_ENDPOINT", "")

		req := httptest.NewRequest("GET", "/api/v1/data", nil)
		w := httptest.NewRecorder()

		GetData(w, req)

		if w.Header().Get("Content-Type") != "application/json" {
			t.Errorf("Content-Type header not set correctly")
		}
	})
}

func TestGetData_DifferentMethods(t *testing.T) {
	t.Run("handles GET request correctly", func(t *testing.T) {
		t.Setenv("SOLARZ_ENDPOINT", "")

		req := httptest.NewRequest("GET", "/api/v1/data", nil)
		w := httptest.NewRecorder()

		GetData(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})
}

