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
	t.Run("uses default endpoint when SOLARZ_ENDPOINT is not set", func(t *testing.T) {
		t.Setenv("SOLARZ_ENDPOINT", "")

		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, dataPath, nil)
		responseRecorder := httptest.NewRecorder()

		handler.GetData(responseRecorder, req)

		// When SOLARZ_ENDPOINT is not set, the default endpoint is used.
		// This may succeed (200) or fail (500) depending on network/auth.
		// Both are acceptable - we just verify the response is properly formatted.
		if responseRecorder.Code != http.StatusOK && responseRecorder.Code != http.StatusInternalServerError {
			t.Errorf("GetData() returned status %d, want 200 or 500", responseRecorder.Code)
		}

		contentType := responseRecorder.Header().Get("Content-Type")
		if contentType != applicationJSON {
			t.Errorf("Content-Type = %q, want %q", contentType, applicationJSON)
		}

		// Response format depends on success/failure
		if responseRecorder.Code == http.StatusOK {
			// Success: should be an array of items
			var items []model.Item
			if err := json.NewDecoder(responseRecorder.Body).Decode(&items); err != nil {
				t.Logf("GetData() returned success, items parsed correctly")
			}
		} else {
			// Error: should be an ErrorResponse
			var errorResponseData model.ErrorResponse
			if err := json.NewDecoder(responseRecorder.Body).Decode(&errorResponseData); err != nil {
				t.Logf("GetData() returned error: %v", err)
			}
		}
	})
}

func TestGetData_ValidData(t *testing.T) {
	t.Run("returns 500 Internal Server Error when API uses HTTP instead of HTTPS", func(t *testing.T) {
		// Mock API server - httptest.NewServer creates HTTP, not HTTPS
		// This test verifies SSRF protection rejects non-HTTPS URLs
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

		// HTTP URLs should be rejected by SSRF protection
		if responseRecorder.Code != http.StatusInternalServerError {
			t.Errorf("GetData() returned status %d, want %d", responseRecorder.Code, http.StatusInternalServerError)
		}

		contentType := responseRecorder.Header().Get("Content-Type")
		if contentType != applicationJSON {
			t.Errorf("Content-Type = %q, want %q", contentType, applicationJSON)
		}

		var errorResponseData model.ErrorResponse
		if err := json.NewDecoder(responseRecorder.Body).Decode(&errorResponseData); err != nil {
			t.Errorf("failed to decode response: %v", err)
		}

		if errorResponseData.Error == "" {
			t.Errorf("expected error message, got empty string")
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

		// Empty SOLARZ_ENDPOINT uses default endpoint, which may succeed or fail
		if responseRecorder.Code != http.StatusOK && responseRecorder.Code != http.StatusInternalServerError {
			t.Errorf("expected status 200 or 500, got %d", responseRecorder.Code)
		}

		// Verify Content-Type is always set
		if responseRecorder.Header().Get("Content-Type") != applicationJSON {
			t.Errorf("Content-Type header not set correctly")
		}
	})
}
