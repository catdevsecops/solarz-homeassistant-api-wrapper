package service

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/catdevsecops/solarz-api/internal/model"
)

func TestFormatFloat(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected string
	}{
		{
			name:     "positive number",
			input:    12.6,
			expected: "12.60",
		},
		{
			name:     "whole number",
			input:    23.0,
			expected: "23.00",
		},
		{
			name:     "decimal number",
			input:    75.4,
			expected: "75.40",
		},
		{
			name:     "zero",
			input:    0.0,
			expected: "0.00",
		},
		{
			name:     "negative number",
			input:    -15.55,
			expected: "-15.55",
		},
		{
			name:     "large number",
			input:    1234567.89,
			expected: "1234567.89",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatFloat(tt.input)
			if result != tt.expected {
				t.Errorf("formatFloat(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGetData_EmptyEndpoint(t *testing.T) {
	t.Run("returns empty slice when SOLARZ_ENDPOINT is not set", func(t *testing.T) {
		t.Setenv("SOLARZ_ENDPOINT", "")
		result, err := GetData()

		if err != nil {
			t.Errorf("GetData() error = %v, want nil", err)
		}

		if len(result) != 0 {
			t.Errorf("GetData() returned %d items, want 0", len(result))
		}
	})
}

func TestGetData_ValidAPI(t *testing.T) {
	t.Run("returns latest item when API returns valid data", func(t *testing.T) {
		// Mock server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
						Quantidade:  16.2,
						Prognostico: 19.79,
						Manual:      false,
						UsinaId:     487759,
						Denominacao: "(3633) Clayton - Mogi",
					},
					{
						Data:        "2026-06-02",
						Quantidade:  23.3,
						Prognostico: 19.79,
						Manual:      false,
						UsinaId:     487759,
						Denominacao: "(3633) Clayton - Mogi",
					},
				},
				TotalGerado:    75.4,
				TotalPrognostico: 138.53,
				Desempenho:     54.42,
			}
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		t.Setenv("SOLARZ_ENDPOINT", server.URL)
		result, err := GetData()

		if err != nil {
			t.Errorf("GetData() error = %v, want nil", err)
		}

		if len(result) != 1 {
			t.Errorf("GetData() returned %d items, want 1", len(result))
		}

		if result[0].ID != "2026-06-04" {
			t.Errorf("result[0].ID = %q, want '2026-06-04'", result[0].ID)
		}

		if result[0].Value != "16.20" {
			t.Errorf("result[0].Value = %q, want '16.20'", result[0].Value)
		}

		expectedName := "2026-06-04 - (3633) Clayton - Mogi"
		if result[0].Name != expectedName {
			t.Errorf("result[0].Name = %q, want %q", result[0].Name, expectedName)
		}
	})
}

func TestGetData_EmptyData(t *testing.T) {
	t.Run("returns empty slice when API returns no data", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response := model.SolarzResponse{
				Dados: []model.DadoGeracao{},
			}
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		t.Setenv("SOLARZ_ENDPOINT", server.URL)
		result, err := GetData()

		if err != nil {
			t.Errorf("GetData() error = %v, want nil", err)
		}

		if len(result) != 0 {
			t.Errorf("GetData() returned %d items, want 0", len(result))
		}
	})
}

func TestGetData_InvalidJSON(t *testing.T) {
	t.Run("returns error when API returns invalid JSON", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte("invalid json"))
		}))
		defer server.Close()

		t.Setenv("SOLARZ_ENDPOINT", server.URL)
		result, err := GetData()

		if err == nil {
			t.Errorf("GetData() error = nil, want error")
		}

		if len(result) != 0 {
			t.Errorf("GetData() returned %d items, want 0", len(result))
		}
	})
}

func TestGetData_NetworkError(t *testing.T) {
	t.Run("returns error when cannot reach API", func(t *testing.T) {
		t.Setenv("SOLARZ_ENDPOINT", "http://invalid-url-that-does-not-exist:9999")
		result, err := GetData()

		if err == nil {
			t.Errorf("GetData() error = nil, want error")
		}

		if len(result) != 0 {
			t.Errorf("GetData() returned %d items, want 0", len(result))
		}
	})
}