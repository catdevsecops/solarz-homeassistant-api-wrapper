package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/catdevsecops/solarz-api/internal/handler"
	"github.com/catdevsecops/solarz-api/internal/model"
)

const applicationJSON = "application/json"

const (
	// Paths.
	dataPath = "/api/v1/data"

	// Response fields.
	statusField    = "status"
	errorField     = "Error"
	errorMessage   = "Failed to fetch items"
	invalidURLAddr = "http://invalid-url-that-does-not-exist:9999"

	// Data values for testing.
	dataDate1       = "2026-06-01"
	dataDate2       = "2026-06-04"
	dataDenominacao = "(3633) Clayton - Mogi"
	expectedValue   = "25.50"
	invalidJSONData = "not valid json"
)

func TestGetData_EmptyEndpoint(t *testing.T) {
	t.Run("returns 200 OK with empty data when SOLARZ_ENDPOINT is not set", func(t *testing.T) {
		t.Setenv("SOLARZ_ENDPOINT", "")

		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, dataPath, nil)
		responseRecorder := httptest.NewRecorder()

		handler.GetData(responseRecorder, req)

		if responseRecorder.Code != http.StatusOK {
			t.Errorf("GetData() returned status %d, want %d", responseRecorder.Code, http.StatusOK)
		}

		contentType := responseRecorder.Header().Get("Content-Type")
		if contentType != applicationJSON {
			t.Errorf("Content-Type = %q, want %q", contentType, applicationJSON)
		}

		var result []model.Item
		if err := json.NewDecoder(responseRecorder.Body).Decode(&result); err != nil {
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
		apiServer := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, _ *http.Request) {
			response := model.SolarzResponse{
				Dados: []model.DadoGeracao{
					{
						Data:        dataDate1,
						Quantidade:  12.6,
						Prognostico: 19.79,
						Manual:      false,
						UsinaID:     487759,
						Denominacao: dataDenominacao,
					},
					{
						Data:        dataDate2,
						Quantidade:  25.5,
						Prognostico: 19.79,
						Manual:      false,
						UsinaID:     487759,
						Denominacao: dataDenominacao,
					},
				},
				TotalGerado:      75.4,
				TotalPrognostico: 138.53,
				Desempenho:       54.42,
			}
			if err := json.NewEncoder(responseWriter).Encode(response); err != nil {
				t.Logf("failed to encode response: %v", err)
			}
		}))
		defer apiServer.Close()

		t.Setenv("SOLARZ_ENDPOINT", apiServer.URL)

		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, dataPath, nil)
		responseRecorder := httptest.NewRecorder()

		handler.GetData(responseRecorder, req)

		if responseRecorder.Code != http.StatusOK {
			t.Errorf("GetData() returned status %d, want %d", responseRecorder.Code, http.StatusOK)
		}

		contentType := responseRecorder.Header().Get("Content-Type")
		if contentType != applicationJSON {
			t.Errorf("Content-Type = %q, want %q", contentType, applicationJSON)
		}

		var result []model.Item
		if err := json.NewDecoder(responseRecorder.Body).Decode(&result); err != nil {
			t.Errorf("failed to decode response: %v", err)
		}

		if len(result) != 1 {
			t.Errorf("expected 1 item, got %d", len(result))
		}

		if result[0].ID != dataDate2 {
			t.Errorf("result[0].ID = %q, want '2026-06-04'", result[0].ID)
		}

		if result[0].Value != expectedValue {
			t.Errorf("result[0].Value = %q, want '25.50'", result[0].Value)
		}
	})
}

func TestGetData_APIError(t *testing.T) {
	t.Run("returns 500 Internal Server Error when API request fails", func(t *testing.T) {
		t.Setenv("SOLARZ_ENDPOINT", invalidURLAddr)

		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, dataPath, nil)
		responseRecorder := httptest.NewRecorder()

		handler.GetData(responseRecorder, req)

		if responseRecorder.Code != http.StatusInternalServerError {
			t.Errorf("GetData() returned status %d, want %d", responseRecorder.Code, http.StatusInternalServerError)
		}

		var result model.ErrorResponse
		if err := json.NewDecoder(responseRecorder.Body).Decode(&result); err != nil {
			t.Errorf("failed to decode response: %v", err)
		}

		if result.Error == "" {
			t.Errorf("expected error message, got empty string")
		}
	})
}

func TestGetData_InvalidJSON(t *testing.T) {
	t.Run("returns 500 Internal Server Error when API returns invalid JSON", func(t *testing.T) {
		apiServer := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, _ *http.Request) {
			responseWriter.Header().Set("Content-Type", applicationJSON)
			if _, err := responseWriter.Write([]byte(invalidJSONData)); err != nil {
				t.Logf("failed to write response: %v", err)
			}
		}))
		defer apiServer.Close()

		t.Setenv("SOLARZ_ENDPOINT", apiServer.URL)

		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, dataPath, nil)
		responseRecorder := httptest.NewRecorder()

		handler.GetData(responseRecorder, req)

		if responseRecorder.Code != http.StatusInternalServerError {
			t.Errorf("GetData() returned status %d, want %d", responseRecorder.Code, http.StatusInternalServerError)
		}

		var result model.ErrorResponse
		if err := json.NewDecoder(responseRecorder.Body).Decode(&result); err != nil {
			t.Errorf("failed to decode response: %v", err)
		}

		if result.Error != errorMessage {
			t.Errorf("result.Error = %q, want 'Failed to fetch items'", result.Error)
		}
	})
}

func TestGetData_HeaderValidation(t *testing.T) {
	t.Run("validates Content-Type header is set before WriteHeader", func(t *testing.T) {
		t.Setenv("SOLARZ_ENDPOINT", "")

		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, dataPath, nil)
		responseRecorder := httptest.NewRecorder()

		handler.GetData(responseRecorder, req)

		if responseRecorder.Header().Get("Content-Type") != applicationJSON {
			t.Errorf("Content-Type header not set correctly")
		}
	})
}

func TestGetData_DifferentMethods(t *testing.T) {
	t.Run("handles GET request correctly", func(t *testing.T) {
		t.Setenv("SOLARZ_ENDPOINT", "")

		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, dataPath, nil)
		responseRecorder := httptest.NewRecorder()

		handler.GetData(responseRecorder, req)

		if responseRecorder.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", responseRecorder.Code)
		}
	})
}
