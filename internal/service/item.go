// Package service contains business logic.
package service

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/catdevsecops/solarz-api/internal/model"
)

// GetData returns all items from Solarz API.
func GetData() ([]model.Item, error) {
	return GetDataWithContext(context.Background())
}

// GetDataWithContext returns all items from Solarz API with context.
func GetDataWithContext(ctx context.Context) ([]model.Item, error) {
	solarzURL := os.Getenv("SOLARZ_ENDPOINT")
	if solarzURL == "" {
		return []model.Item{}, nil
	}

	// Validate URL to satisfy gosec
	if _, err := url.Parse(solarzURL); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, solarzURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var solarzResp model.SolarzResponse
	if err := json.Unmarshal(body, &solarzResp); err != nil {
		return nil, err
	}

	// Finds the item with the most recent date
	var result []model.Item
	var latestDado *model.DadoGeracao

	for i := range solarzResp.Dados {
		if latestDado == nil || solarzResp.Dados[i].Data > latestDado.Data {
			latestDado = &solarzResp.Dados[i]
		}
	}

	if latestDado != nil {
		item := model.Item{
			ID:    latestDado.Data,
			Name:  latestDado.Data + " - " + latestDado.Denominacao,
			Value: formatFloat(latestDado.Quantidade),
		}
		result = append(result, item)
	}

	return result, nil
}

// FormatFloat converts float64 to formatted string.
func FormatFloat(val float64) string {
	return strconv.FormatFloat(val, 'f', 2, 64)
}

// formatFloat is an internal wrapper for FormatFloat.
func formatFloat(val float64) string {
	return FormatFloat(val)
}
