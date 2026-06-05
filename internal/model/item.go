// Package model contains the data structures.
package model

// Item represents an item in the system.
type Item struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

// ErrorResponse represents an error response.
type ErrorResponse struct {
	Error string `json:"error"`
}

// SolarzResponse represents the response structure of the Solarz API.
type SolarzResponse struct {
	Dados              []DadoGeracao         `json:"dados"`
	TotalGerado        float64               `json:"totalGerado"`
	TotalPrognostico   float64               `json:"totalPrognostico"`
	Desempenho         float64               `json:"desempenho"`
	LabeledGenerations map[string]LabelValue `json:"labeledGenerations"`
	Prognosticos       map[string]float64    `json:"prognosticos"`
	MaisPortais        bool                  `json:"morePortais"`
}

// DadoGeracao represents a day of generation.
type DadoGeracao struct {
	Data            string           `json:"data"`
	Quantidade      float64          `json:"quantidade"`
	Prognostico     float64          `json:"prognostico"`
	InformacaoClima InformacaoClima  `json:"informacaoClima"`
	Manual          bool             `json:"manual"`
	UsinaID         int              `json:"usinaId"`
	Denominacao     string           `json:"denominacao"`
	Geracoes        []GeracaoDetalhe `json:"geracoes"`
	PlantShutdown   bool             `json:"plantShutdown"`
}

// InformacaoClima represents climate information.
type InformacaoClima struct {
	ID        int    `json:"id"`
	Descricao string `json:"descricao"`
	CreatedAt string `json:"createdAt"`
}

// GeracaoDetalhe represents generation details.
type GeracaoDetalhe struct {
	Quantidade float64 `json:"quantidade"`
	IDExterno  string  `json:"idExterno"`
	Descricao  *string `json:"descricao"`
}

// LabelValue represents a label with its value.
type LabelValue struct {
	Label string  `json:"label"`
	Value float64 `json:"value"`
}

// IsValid checks if the Item has a valid ID.
func (i *Item) IsValid() bool {
	return i != nil && i.ID != ""
}

// IsEmpty checks if the Item is empty.
func (i *Item) IsEmpty() bool {
	return i == nil || (i.ID == "" && i.Name == "" && i.Value == "")
}

// HasError checks if the error response contains a message.
func (e *ErrorResponse) HasError() bool {
	return e != nil && e.Error != ""
}

// GetTotalDados returns the number of data in DadoGeracao.
func (sr *SolarzResponse) GetTotalDados() int {
	if sr == nil {
		return 0
	}
	return len(sr.Dados)
}

// CalculateDesempenho calculates performance based on quantity and forecast.
func (dg *DadoGeracao) CalculateDesempenho() float64 {
	if dg == nil || dg.Prognostico == 0 {
		return 0
	}
	return dg.Quantidade / dg.Prognostico
}

// IsManualEntry checks if it is a manual entry.
func (dg *DadoGeracao) IsManualEntry() bool {
	return dg != nil && dg.Manual
}

// HasDescription checks if it has a description.
func (gd *GeracaoDetalhe) HasDescription() bool {
	return gd != nil && gd.Descricao != nil && *gd.Descricao != ""
}

// IsValid checks if LabelValue is valid.
func (lv *LabelValue) IsValid() bool {
	return lv != nil && lv.Label != ""
}
