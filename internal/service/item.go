package service

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/catdevsecops/solarz-api/internal/model"
)

// GetData returns all items from Solarz API
func GetData() ([]model.Item, error) {
	solarzURL := os.Getenv("SOLARZ_ENDPOINT")
	if solarzURL == "" {
		return []model.Item{}, nil
	}

	resp, err := http.Get(solarzURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

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

// formatFloat converts float64 to formatted string
func formatFloat(val float64) string {
	return strconv.FormatFloat(val, 'f', 2, 64)
}

