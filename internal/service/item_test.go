package service_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/catdevsecops/solarz-api/internal/model"
	"github.com/catdevsecops/solarz-api/internal/service"
)

const (
	// Format float test values and expectations.
	posNumber   = "positive number"
	wholNumber  = "whole number"
	decNumber   = "decimal number"
	zeroValue   = "zero"
	negNumber   = "negative number"
	largeNumber = "large number"

	posExpected   = "12.60"
	wholExpected  = "23.00"
	decExpected   = "75.40"
	zeroExpected  = "0.00"
	negExpected   = "-15.55"
	largeExpected = "1234567.89"

	// Data for API tests.
	dataDate1       = "2026-06-01"
	dataDate2       = "2026-06-04"
	dataDate3       = "2026-06-02"
	userQuantity    = "16.20"
	dataDenominacao = "(3633) Clayton - Mogi"
	invalidJSONData = "invalid json"
	invalidURLAddr  = "http://invalid-url-that-does-not-exist:9999"
	expectString    = "2026-06-04 - (3633) Clayton - Mogi"
)

func TestFormatFloat(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected string
	}{
		{
			name:     posNumber,
			input:    12.6,
			expected: posExpected,
		},
		{
			name:     wholNumber,
			input:    23.0,
			expected: wholExpected,
		},
		{
			name:     decNumber,
			input:    75.4,
			expected: decExpected,
		},
		{
			name:     zeroValue,
			input:    0.0,
			expected: zeroExpected,
		},
		{
			name:     negNumber,
			input:    -15.55,
			expected: negExpected,
		},
		{
			name:     largeNumber,
			input:    1234567.89,
			expected: largeExpected,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			result := service.FormatFloat(testCase.input)
			if result != testCase.expected {
				t.Errorf("formatFloat(%v) = %q, want %q", testCase.input, result, testCase.expected)
			}
		})
	}
}

func TestGetData_EmptyEndpoint(t *testing.T) {
	t.Run("returns empty slice when SOLARZ_ENDPOINT is not set", func(t *testing.T) {
		t.Setenv("SOLARZ_ENDPOINT", "")
		result, err := service.GetData()
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
						Quantidade:  16.2,
						Prognostico: 19.79,
						Manual:      false,
						UsinaID:     487759,
						Denominacao: dataDenominacao,
					},
					{
						Data:        dataDate3,
						Quantidade:  23.3,
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
		result, err := service.GetData()
		if err != nil {
			t.Errorf("GetData() error = %v, want nil", err)
		}

		if len(result) != 1 {
			t.Errorf("GetData() returned %d items, want 1", len(result))
		}

		if result[0].ID != dataDate2 {
			t.Errorf("result[0].ID = %q, want '2026-06-04'", result[0].ID)
		}

		if result[0].Value != userQuantity {
			t.Errorf("result[0].Value = %q, want '16.20'", result[0].Value)
		}

		expectedName := expectString
		if result[0].Name != expectedName {
			t.Errorf("result[0].Name = %q, want %q", result[0].Name, expectedName)
		}
	})
}

func TestGetData_EmptyData(t *testing.T) {
	t.Run("returns empty slice when API returns no data", func(t *testing.T) {
		apiServer := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, _ *http.Request) {
			response := model.SolarzResponse{
				Dados: []model.DadoGeracao{},
			}
			if err := json.NewEncoder(responseWriter).Encode(response); err != nil {
				t.Logf("failed to encode response: %v", err)
			}
		}))
		defer apiServer.Close()

		t.Setenv("SOLARZ_ENDPOINT", apiServer.URL)
		result, err := service.GetData()
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
		apiServer := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, _ *http.Request) {
			responseWriter.Header().Set("Content-Type", "application/json")
			if _, err := responseWriter.Write([]byte(invalidJSONData)); err != nil {
				t.Logf("failed to write response: %v", err)
			}
		}))
		defer apiServer.Close()

		t.Setenv("SOLARZ_ENDPOINT", apiServer.URL)
		result, err := service.GetData()

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
		t.Setenv("SOLARZ_ENDPOINT", invalidURLAddr)
		result, err := service.GetData()

		if err == nil {
			t.Errorf("GetData() error = nil, want error")
		}

		if len(result) != 0 {
			t.Errorf("GetData() returned %d items, want 0", len(result))
		}
	})
}
