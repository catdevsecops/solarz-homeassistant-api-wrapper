package model_test

import (
	"encoding/json"
	"testing"

	"github.com/catdevsecops/solarz-api/internal/model"
)

const testDescription = "Test description"

const (
	// Test IDs and names.
	testID    = "test-id"
	testName  = "Test Name"
	testValue = "Test Value"

	// Data values.
	dataDate1  = "2024-01-01"
	usinaDenom = "Usina A"

	// Climate.
	climaDesc    = "Sunny"
	climaCreated = "2024-01-01T10:00:00Z"
	extIDValue   = "ext-1"

	// Generation.
	genLabel1   = "Generation 1"
	genLabel2   = "Generation 2"
	genLabelStr = "Generation"
)

// TestItem testa a estrutura *model.Item.
func TestItem(t *testing.T) {
	item := model.Item{
		ID:    testID,
		Name:  testName,
		Value: testValue,
	}

	if item.ID != testID {
		t.Errorf("Expected ID 'test-id', got '%s'", item.ID)
	}

	if item.Name != testName {
		t.Errorf("Expected Name 'Test Name', got '%s'", item.Name)
	}

	if item.Value != testValue {
		t.Errorf("Expected Value 'Test Value', got '%s'", item.Value)
	}
}

// TestItemJSON testa a serialização/desserialização JSON do *model.Item.
func TestItemJSON(t *testing.T) {
	original := model.Item{
		ID:    testID,
		Name:  testName,
		Value: testValue,
	}

	// Serializar para JSON
	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal *model.Item: %v", err)
	}

	// Desserializar de volta
	var deserialized model.Item
	err = json.Unmarshal(jsonData, &deserialized)
	if err != nil {
		t.Fatalf("Failed to unmarshal *model.Item: %v", err)
	}

	if deserialized.ID != original.ID {
		t.Errorf("Expected ID '%s', got '%s'", original.ID, deserialized.ID)
	}

	if deserialized.Name != original.Name {
		t.Errorf("Expected Name '%s', got '%s'", original.Name, deserialized.Name)
	}

	if deserialized.Value != original.Value {
		t.Errorf("Expected Value '%s', got '%s'", original.Value, deserialized.Value)
	}
}

// TestItemEmptyFields testa Item com campos vazios.
func TestItemEmptyFields(t *testing.T) {
	item := model.Item{}

	if item.ID != "" {
		t.Errorf("Expected empty ID, got '%s'", item.ID)
	}

	if item.Name != "" {
		t.Errorf("Expected empty Name, got '%s'", item.Name)
	}

	if item.Value != "" {
		t.Errorf("Expected empty Value, got '%s'", item.Value)
	}
}

// TestErrorResponse testa a estrutura *model.ErrorResponse.
func TestErrorResponse(t *testing.T) {
	errResp := model.ErrorResponse{
		Error: "Test error message",
	}

	if errResp.Error != "Test error message" {
		t.Errorf("Expected error message 'Test error message', got '%s'", errResp.Error)
	}
}

// TestErrorResponseJSON testa a serialização/desserialização JSON do *model.ErrorResponse.
func TestErrorResponseJSON(t *testing.T) {
	original := model.ErrorResponse{
		Error: "Test error message",
	}

	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal *model.ErrorResponse: %v", err)
	}

	var deserialized model.ErrorResponse
	err = json.Unmarshal(jsonData, &deserialized)
	if err != nil {
		t.Fatalf("Failed to unmarshal *model.ErrorResponse: %v", err)
	}

	if deserialized.Error != original.Error {
		t.Errorf("Expected error '%s', got '%s'", original.Error, deserialized.Error)
	}
}

// TestSolarzResponse testa a estrutura *model.SolarzResponse.
func TestSolarzResponse(t *testing.T) {
	resp := model.SolarzResponse{
		TotalGerado:      100.5,
		TotalPrognostico: 95.3,
		Desempenho:       0.95,
		MaisPortais:      true,
	}

	if resp.TotalGerado != 100.5 {
		t.Errorf("Expected TotalGerado 100.5, got %f", resp.TotalGerado)
	}

	if resp.TotalPrognostico != 95.3 {
		t.Errorf("Expected TotalPrognostico 95.3, got %f", resp.TotalPrognostico)
	}

	if resp.Desempenho != 0.95 {
		t.Errorf("Expected Desempenho 0.95, got %f", resp.Desempenho)
	}

	if !resp.MaisPortais {
		t.Error("Expected MaisPortais to be true")
	}
}

// TestSolarzResponseWithDados testa SolarzResponse com dados.
func TestSolarzResponseWithDados(t *testing.T) {
	resp := model.SolarzResponse{
		Dados: []model.DadoGeracao{
			{
				Data:        dataDate1,
				Quantidade:  50.0,
				Prognostico: 48.0,
				UsinaID:     1,
			},
		},
		TotalGerado: 50.0,
	}

	if len(resp.Dados) != 1 {
		t.Errorf("Expected 1 *model.DadoGeracao, got %d", len(resp.Dados))
	}

	if resp.Dados[0].Data != dataDate1 {
		t.Errorf("Expected Data '2024-01-01', got '%s'", resp.Dados[0].Data)
	}
}

// TestSolarzResponseJSON testa a serialização/desserialização JSON do *model.SolarzResponse.
func TestSolarzResponseJSON(t *testing.T) {
	original := model.SolarzResponse{
		TotalGerado:      100.5,
		TotalPrognostico: 95.3,
		Desempenho:       0.95,
		MaisPortais:      true,
	}

	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal *model.SolarzResponse: %v", err)
	}

	var deserialized model.SolarzResponse
	err = json.Unmarshal(jsonData, &deserialized)
	if err != nil {
		t.Fatalf("Failed to unmarshal *model.SolarzResponse: %v", err)
	}

	if deserialized.TotalGerado != original.TotalGerado {
		t.Errorf("Expected TotalGerado %f, got %f", original.TotalGerado, deserialized.TotalGerado)
	}

	if deserialized.MaisPortais != original.MaisPortais {
		t.Errorf("Expected MaisPortais %v, got %v", original.MaisPortais, deserialized.MaisPortais)
	}
}

// TestDadoGeracao testa a estrutura *model.DadoGeracao.
func TestDadoGeracao(t *testing.T) {
	dado := model.DadoGeracao{
		Data:          dataDate1,
		Quantidade:    50.0,
		Prognostico:   48.0,
		Manual:        false,
		UsinaID:       1,
		Denominacao:   usinaDenom,
		PlantShutdown: false,
	}

	if dado.Data != dataDate1 {
		t.Errorf("Expected Data '2024-01-01', got '%s'", dado.Data)
	}

	if dado.Quantidade != 50.0 {
		t.Errorf("Expected Quantidade 50.0, got %f", dado.Quantidade)
	}

	if dado.Prognostico != 48.0 {
		t.Errorf("Expected Prognostico 48.0, got %f", dado.Prognostico)
	}

	if dado.Manual {
		t.Error("Expected Manual to be false")
	}

	if dado.PlantShutdown {
		t.Error("Expected PlantShutdown to be false")
	}
}

// TestDadoGeracaoWithGeracoes testa DadoGeracao com detalhes de geração.
func TestDadoGeracaoWithGeracoes(t *testing.T) {
	descricao := "Test generation"
	dado := model.DadoGeracao{
		Data:       dataDate1,
		Quantidade: 50.0,
		Geracoes: []model.GeracaoDetalhe{
			{
				Quantidade: 25.0,
				IDExterno:  extIDValue,
				Descricao:  &descricao,
			},
		},
	}

	if len(dado.Geracoes) != 1 {
		t.Errorf("Expected 1 *model.GeracaoDetalhe, got %d", len(dado.Geracoes))
	}

	if dado.Geracoes[0].Quantidade != 25.0 {
		t.Errorf("Expected Quantidade 25.0, got %f", dado.Geracoes[0].Quantidade)
	}

	if *dado.Geracoes[0].Descricao != "Test generation" {
		t.Errorf("Expected Descricao 'Test generation', got '%s'", *dado.Geracoes[0].Descricao)
	}
}

// TestDadoGeracaoJSON testa a serialização/desserialização JSON do *model.DadoGeracao.
func TestDadoGeracaoJSON(t *testing.T) {
	original := model.DadoGeracao{
		Data:        dataDate1,
		Quantidade:  50.0,
		Prognostico: 48.0,
		UsinaID:     1,
		Denominacao: usinaDenom,
	}

	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal *model.DadoGeracao: %v", err)
	}

	var deserialized model.DadoGeracao
	err = json.Unmarshal(jsonData, &deserialized)
	if err != nil {
		t.Fatalf("Failed to unmarshal *model.DadoGeracao: %v", err)
	}

	if deserialized.Data != original.Data {
		t.Errorf("Expected Data '%s', got '%s'", original.Data, deserialized.Data)
	}

	if deserialized.Denominacao != original.Denominacao {
		t.Errorf("Expected Denominacao '%s', got '%s'", original.Denominacao, deserialized.Denominacao)
	}
}

// TestInformacaoClima testa a estrutura *model.InformacaoClima.
func TestInformacaoClima(t *testing.T) {
	clima := model.InformacaoClima{
		ID:        1,
		Descricao: climaDesc,
		CreatedAt: climaCreated,
	}

	if clima.ID != 1 {
		t.Errorf("Expected Id 1, got %d", clima.ID)
	}

	if clima.Descricao != climaDesc {
		t.Errorf("Expected Descricao 'Sunny', got '%s'", clima.Descricao)
	}

	if clima.CreatedAt != climaCreated {
		t.Errorf("Expected CreatedAt '2024-01-01T10:00:00Z', got '%s'", clima.CreatedAt)
	}
}

// TestInformacaoClimaJSON testa a serialização/desserialização JSON do *model.InformacaoClima.
func TestInformacaoClimaJSON(t *testing.T) {
	original := model.InformacaoClima{
		ID:        1,
		Descricao: climaDesc,
		CreatedAt: climaCreated,
	}

	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal *model.InformacaoClima: %v", err)
	}

	var deserialized model.InformacaoClima
	err = json.Unmarshal(jsonData, &deserialized)
	if err != nil {
		t.Fatalf("Failed to unmarshal *model.InformacaoClima: %v", err)
	}

	if deserialized.ID != original.ID {
		t.Errorf("Expected Id %d, got %d", original.ID, deserialized.ID)
	}

	if deserialized.Descricao != original.Descricao {
		t.Errorf("Expected Descricao '%s', got '%s'", original.Descricao, deserialized.Descricao)
	}
}

// TestGeracaoDetalhe testa a estrutura *model.GeracaoDetalhe.
func TestGeracaoDetalhe(t *testing.T) {
	geracao := model.GeracaoDetalhe{
		Quantidade: 25.0,
		IDExterno:  extIDValue,
	}

	if geracao.Quantidade != 25.0 {
		t.Errorf("Expected Quantidade 25.0, got %f", geracao.Quantidade)
	}

	if geracao.IDExterno != extIDValue {
		t.Errorf("Expected IdExterno 'ext-1', got '%s'", geracao.IDExterno)
	}

	if geracao.Descricao != nil {
		t.Error("Expected Descricao to be nil")
	}
}

// TestGeracaoDetalheWithDescricao testa GeracaoDetalhe com descrição.
func TestGeracaoDetalheWithDescricao(t *testing.T) {
	descricao := testDescription
	geracao := model.GeracaoDetalhe{
		Quantidade: 25.0,
		IDExterno:  extIDValue,
		Descricao:  &descricao,
	}

	if geracao.Descricao == nil {
		t.Error("Expected Descricao to not be nil")
	}

	if *geracao.Descricao != testDescription {
		t.Errorf("Expected Descricao 'Test description', got '%s'", *geracao.Descricao)
	}
}

// TestGeracaoDetalheJSON testa a serialização/desserialização JSON do *model.GeracaoDetalhe.
func TestGeracaoDetalheJSON(t *testing.T) {
	descricao := testDescription
	original := model.GeracaoDetalhe{
		Quantidade: 25.0,
		IDExterno:  extIDValue,
		Descricao:  &descricao,
	}

	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal *model.GeracaoDetalhe: %v", err)
	}

	var deserialized model.GeracaoDetalhe
	err = json.Unmarshal(jsonData, &deserialized)
	if err != nil {
		t.Fatalf("Failed to unmarshal *model.GeracaoDetalhe: %v", err)
	}

	if deserialized.Quantidade != original.Quantidade {
		t.Errorf("Expected Quantidade %f, got %f", original.Quantidade, deserialized.Quantidade)
	}

	if deserialized.IDExterno != original.IDExterno {
		t.Errorf("Expected IdExterno '%s', got '%s'", original.IDExterno, deserialized.IDExterno)
	}
}

// TestLabelValue testa a estrutura *model.LabelValue.
func TestLabelValue(t *testing.T) {
	label := model.LabelValue{
		Label: genLabelStr,
		Value: 100.5,
	}

	if label.Label != genLabelStr {
		t.Errorf("Expected Label 'Generation', got '%s'", label.Label)
	}

	if label.Value != 100.5 {
		t.Errorf("Expected Value 100.5, got %f", label.Value)
	}
}

// TestLabelValueJSON testa a serialização/desserialização JSON do *model.LabelValue.
func TestLabelValueJSON(t *testing.T) {
	original := model.LabelValue{
		Label: genLabelStr,
		Value: 100.5,
	}

	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal *model.LabelValue: %v", err)
	}

	var deserialized model.LabelValue
	err = json.Unmarshal(jsonData, &deserialized)
	if err != nil {
		t.Fatalf("Failed to unmarshal *model.LabelValue: %v", err)
	}

	if deserialized.Label != original.Label {
		t.Errorf("Expected Label '%s', got '%s'", original.Label, deserialized.Label)
	}

	if deserialized.Value != original.Value {
		t.Errorf("Expected Value %f, got %f", original.Value, deserialized.Value)
	}
}

// TestSolarzResponseWithLabeledGenerations testa SolarzResponse com geração rotulada.
func TestSolarzResponseWithLabeledGenerations(t *testing.T) {
	resp := model.SolarzResponse{
		LabeledGenerations: map[string]model.LabelValue{
			"gen1": {Label: genLabel1, Value: 50.0},
			"gen2": {Label: genLabel2, Value: 45.5},
		},
	}

	if len(resp.LabeledGenerations) != 2 {
		t.Errorf("Expected 2 labeled generations, got %d", len(resp.LabeledGenerations))
	}

	gen1, ok := resp.LabeledGenerations["gen1"]
	if !ok {
		t.Fatal("Expected 'gen1' key in LabeledGenerations")
	}

	if gen1.Value != 50.0 {
		t.Errorf("Expected gen1 Value 50.0, got %f", gen1.Value)
	}
}

// TestSolarzResponseWithPrognosticos testa SolarzResponse com prognósticos.
func TestSolarzResponseWithPrognosticos(t *testing.T) {
	resp := model.SolarzResponse{
		Prognosticos: map[string]float64{
			"jan": 100.0,
			"fev": 95.5,
		},
	}

	if len(resp.Prognosticos) != 2 {
		t.Errorf("Expected 2 prognósticos, got %d", len(resp.Prognosticos))
	}

	janValue, ok := resp.Prognosticos["jan"]
	if !ok {
		t.Fatal("Expected 'jan' key in Prognosticos")
	}

	if janValue != 100.0 {
		t.Errorf("Expected jan value 100.0, got %f", janValue)
	}
}

// TestComplexSolarzResponseJSON testa a serialização/desserialização JSON de uma *model.SolarzResponse complexa.
func TestComplexSolarzResponseJSON(t *testing.T) {
	descricao := testDescription
	original := model.SolarzResponse{
		Dados: []model.DadoGeracao{
			{
				Data:        dataDate1,
				Quantidade:  50.0,
				Prognostico: 48.0,
				InformacaoClima: model.InformacaoClima{
					ID:        1,
					Descricao: climaDesc,
					CreatedAt: climaCreated,
				},
				Manual:      false,
				UsinaID:     1,
				Denominacao: usinaDenom,
				Geracoes: []model.GeracaoDetalhe{
					{
						Quantidade: 25.0,
						IDExterno:  extIDValue,
						Descricao:  &descricao,
					},
				},
				PlantShutdown: false,
			},
		},
		TotalGerado:      50.0,
		TotalPrognostico: 48.0,
		Desempenho:       0.96,
		LabeledGenerations: map[string]model.LabelValue{
			"gen1": {Label: genLabel1, Value: 50.0},
		},
		Prognosticos: map[string]float64{
			"jan": 100.0,
		},
		MaisPortais: false,
	}

	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal complex *model.SolarzResponse: %v", err)
	}

	var deserialized model.SolarzResponse
	err = json.Unmarshal(jsonData, &deserialized)
	if err != nil {
		t.Fatalf("Failed to unmarshal complex *model.SolarzResponse: %v", err)
	}

	if len(deserialized.Dados) != 1 {
		t.Errorf("Expected 1 *model.DadoGeracao, got %d", len(deserialized.Dados))
	}

	if deserialized.Dados[0].Denominacao != usinaDenom {
		t.Errorf("Expected Denominacao 'Usina A', got '%s'", deserialized.Dados[0].Denominacao)
	}

	if len(deserialized.Dados[0].Geracoes) != 1 {
		t.Errorf("Expected 1 Geracao, got %d", len(deserialized.Dados[0].Geracoes))
	}
}

// TestItemIsValid testa o método IsValid de *model.Item.
func TestItemIsValid(t *testing.T) {
	tests := []struct {
		name     string
		item     *model.Item
		expected bool
	}{
		{"Valid item", &model.Item{ID: "123", Name: "Test"}, true},
		{"Item with empty ID", &model.Item{ID: "", Name: "Test"}, false},
		{"Nil item", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.item.IsValid()
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestItemIsEmpty testa o método IsEmpty de *model.Item.
func TestItemIsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		item     *model.Item
		expected bool
	}{
		{"Empty item", &model.Item{}, true},
		{"Item with ID", &model.Item{ID: "123"}, false},
		{"Item with Name", &model.Item{Name: "Test"}, false},
		{"Item with Value", &model.Item{Value: "Value"}, false},
		{"Nil item", nil, true},
		{"Complete item", &model.Item{ID: "123", Name: "Test", Value: "Value"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.item.IsEmpty()
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestErrorResponseHasError testa o método HasError de *model.ErrorResponse.
func TestErrorResponseHasError(t *testing.T) {
	tests := []struct {
		name     string
		errResp  *model.ErrorResponse
		expected bool
	}{
		{"With error message", &model.ErrorResponse{Error: "Something went wrong"}, true},
		{"With empty error", &model.ErrorResponse{Error: ""}, false},
		{"Nil error response", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.errResp.HasError()
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestSolarzResponseGetTotalDados testa o método GetTotalDados.
func TestSolarzResponseGetTotalDados(t *testing.T) {
	tests := []struct {
		name     string
		resp     *model.SolarzResponse
		expected int
	}{
		{"With one dato", &model.SolarzResponse{Dados: []model.DadoGeracao{{Data: dataDate1}}}, 1},
		{"With multiple dados", &model.SolarzResponse{Dados: []model.DadoGeracao{
			{Data: dataDate1}, {Data: "2024-01-02"}, {Data: "2024-01-03"},
		}}, 3},
		{"With no dados", &model.SolarzResponse{Dados: []model.DadoGeracao{}}, 0},
		{"Nil response", nil, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.resp.GetTotalDados()
			if result != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result)
			}
		})
	}
}

// TestDadoGeracaoCalculateDesempenho testa o método CalculateDesempenho.
func TestDadoGeracaoCalculateDesempenho(t *testing.T) {
	tests := []struct {
		name     string
		dado     *model.DadoGeracao
		expected float64
	}{
		{"Normal calculation", &model.DadoGeracao{Quantidade: 100, Prognostico: 100}, 1.0},
		{"Above prognosis", &model.DadoGeracao{Quantidade: 120, Prognostico: 100}, 1.2},
		{"Below prognosis", &model.DadoGeracao{Quantidade: 80, Prognostico: 100}, 0.8},
		{"Zero prognosis", &model.DadoGeracao{Quantidade: 100, Prognostico: 0}, 0},
		{"Nil dado", nil, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.dado.CalculateDesempenho()
			if result != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, result)
			}
		})
	}
}

// TestDadoGeracaoIsManualEntry testa o método IsManualEntry.
func TestDadoGeracaoIsManualEntry(t *testing.T) {
	tests := []struct {
		name     string
		dado     *model.DadoGeracao
		expected bool
	}{
		{"Manual entry", &model.DadoGeracao{Manual: true}, true},
		{"Automatic entry", &model.DadoGeracao{Manual: false}, false},
		{"Nil dado", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.dado.IsManualEntry()
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestGeracaoDetalheHasDescription testa o método HasDescription.
func TestGeracaoDetalheHasDescription(t *testing.T) {
	descricaoNaoVazia := testDescription
	descricaoVazia := ""
	tests := []struct {
		name     string
		geracao  *model.GeracaoDetalhe
		expected bool
	}{
		{"With description", &model.GeracaoDetalhe{Descricao: &descricaoNaoVazia}, true},
		{"With empty description", &model.GeracaoDetalhe{Descricao: &descricaoVazia}, false},
		{"Without description", &model.GeracaoDetalhe{}, false},
		{"Nil geracao", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.geracao.HasDescription()
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestLabelValueIsValid testa o método IsValid de *model.LabelValue.
func TestLabelValueIsValid(t *testing.T) {
	tests := []struct {
		name     string
		label    *model.LabelValue
		expected bool
	}{
		{"Valid label", &model.LabelValue{Label: genLabelStr, Value: 100.5}, true},
		{"With empty label", &model.LabelValue{Label: "", Value: 100.5}, false},
		{"Nil label", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.label.IsValid()
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}
