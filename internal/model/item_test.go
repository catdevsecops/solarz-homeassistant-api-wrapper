package model_test

import (
	"encoding/json"
	"testing"

	"github.com/catdevsecops/solarz-api/internal/model"
)

const (
	// Test IDs and names.
	testID          = "test-id"
	testName        = "Test Name"
	testValue       = "Test Value"
	testDescription = "Test description"

	// Data values.
	dataDate1  = "2024-01-01"
	dataDate2  = "2024-01-02"
	dataDate3  = "2024-01-03"
	usinaDenom = "Usina A"

	// Climate.
	climaDesc    = "Sunny"
	climaCreated = "2024-01-01T10:00:00Z"
	extIDValue   = "ext-1"

	// Generation.
	genLabel1   = "Generation 1"
	genLabel2   = "Generation 2"
	genLabelStr = "Generation"

	// Error messages.
	errMarshalItem             = "Failed to marshal *model.Item"
	errUnmarshalItem           = "Failed to unmarshal *model.Item"
	errMarshalErrorResp        = "Failed to marshal *model.ErrorResponse"
	errUnmarshalErrorResp      = "Failed to unmarshal *model.ErrorResponse"
	errMarshalSolarzResp       = "Failed to marshal *model.SolarzResponse"
	errUnmarshalSolarzResp     = "Failed to unmarshal *model.SolarzResponse"
	errMarshalDadoGeracao      = "Failed to marshal *model.DadoGeracao"
	errUnmarshalDadoGeracao    = "Failed to unmarshal *model.DadoGeracao"
	errMarshalClimaInfo        = "Failed to marshal *model.InformacaoClima"
	errUnmarshalClimaInfo      = "Failed to unmarshal *model.InformacaoClima"
	errMarshalGeracaoDetalhe   = "Failed to marshal *model.GeracaoDetalhe"
	errUnmarshalGeracaoDetalhe = "Failed to unmarshal *model.GeracaoDetalhe"
	errMarshalLabelValue       = "Failed to marshal *model.LabelValue"
	errUnmarshalLabelValue     = "Failed to unmarshal *model.LabelValue"
	errMarshalComplexResp      = "Failed to marshal complex *model.SolarzResponse"
	errUnmarshalComplexResp    = "Failed to unmarshal complex *model.SolarzResponse"

	// Test case names.
	testCaseValidItem          = "Valid item"
	testCaseItemEmptyID        = "Item with empty ID"
	testCaseNilItem            = "Nil item"
	testCaseEmptyItem          = "Empty item"
	testCaseItemWithID         = "Item with ID"
	testCaseItemWithName       = "Item with Name"
	testCaseItemWithValue      = "Item with Value"
	testCaseCompleteItem       = "Complete item"
	testCaseWithError          = "With error message"
	testCaseWithEmptyError     = "With empty error"
	testCaseNilErrorResponse   = "Nil error response"
	testCaseOneDato            = "With one dato"
	testCaseMultipleDados      = "With multiple dados"
	testCaseNoDados            = "With no dados"
	testCaseNilResponse        = "Nil response"
	testCaseNormalCalc         = "Normal calculation"
	testCaseAbovePrognosis     = "Above prognosis"
	testCaseBelowPrognosis     = "Below prognosis"
	testCaseZeroPrognosis      = "Zero prognosis"
	testCaseNilDado            = "Nil dado"
	testCaseManualEntry        = "Manual entry"
	testCaseAutomaticEntry     = "Automatic entry"
	testCaseWithDescription    = "With description"
	testCaseWithEmptyDesc      = "With empty description"
	testCaseWithoutDescription = "Without description"
	testCaseNilGeracao         = "Nil geracao"
	testCaseValidLabel         = "Valid label"
	testCaseEmptyLabel         = "With empty label"
	testCaseNilLabel           = "Nil label"

	// Expected values in assertions.
	expectedTestID         = "Expected ID 'test-id', got '%s'"
	expectedTestName       = "Expected Name 'Test Name', got '%s'"
	expectedTestValue      = "Expected Value 'Test Value', got '%s'"
	expectedEmptyID        = "Expected empty ID, got '%s'"
	expectedEmptyName      = "Expected empty Name, got '%s'"
	expectedEmptyValue     = "Expected empty Value, got '%s'"
	expectedErrorMsg       = "Expected error message 'Test error message', got '%s'"
	expectedDataDate       = "Expected Data '%s', got '%s'"
	expectedQuantity       = "Expected Quantidade %f, got %f"
	expectedPrognostico    = "Expected Prognostico %f, got %f"
	expectedQuantity25_0   = "Expected Quantidade 25.0, got %f"
	expectedVsGot          = "Expected %v, got %v"
	errorTestMessage       = "Test error message"
	expectedWentWrong      = "Something went wrong"
	expectedDenominacao    = "Expected Denominacao '%s', got '%s'"
	expectedDescription    = "Expected Descricao '%s', got '%s'"
	expectedClimaDesc      = "Expected Descricao 'Sunny', got '%s'"
	expectedClimaCreatedAt = "Expected CreatedAt '2024-01-01T10:00:00Z', got '%s'"
	expectedIDExterno      = "Expected IdExterno 'ext-1', got '%s'"
	expectedLabel          = "Expected Label '%s', got '%s'"
	expectedLabelValue     = "Expected Label 'Generation', got '%s'"
	expectedValue100_5     = "Expected Value 100.5, got %f"
	expectedValue50_0      = "Expected gen1 Value 50.0, got %f"
	expectedValue100_0     = "Expected jan value 100.0, got %f"
	expectedDados1         = "Expected 1 *model.DadoGeracao, got %d"
	expectedDados          = "Expected 1 Geracao, got %d"
	expectedGeracoes1      = "Expected 1 *model.GeracaoDetalhe, got %d"
	expectedLabeledGen2    = "Expected 2 labeled generations, got %d"
	expectedPrognosticos2  = "Expected 2 prognósticos, got %d"
	expectedGenKey         = "Expected 'gen1' key in LabeledGenerations"
	expectedJanKey         = "Expected 'jan' key in Prognosticos"
)

// TestItem testa a estrutura *model.Item.
func TestItem(t *testing.T) {
	item := model.Item{
		ID:    testID,
		Name:  testName,
		Value: testValue,
	}

	if item.ID != testID {
		t.Errorf(expectedTestID, item.ID)
	}

	if item.Name != testName {
		t.Errorf(expectedTestName, item.Name)
	}

	if item.Value != testValue {
		t.Errorf(expectedTestValue, item.Value)
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
		t.Fatalf(errMarshalItem+": %v", err)
	}

	// Desserializar de volta
	var deserialized model.Item
	err = json.Unmarshal(jsonData, &deserialized)
	if err != nil {
		t.Fatalf(errUnmarshalItem+": %v", err)
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
		t.Errorf(expectedEmptyID, item.ID)
	}

	if item.Name != "" {
		t.Errorf(expectedEmptyName, item.Name)
	}

	if item.Value != "" {
		t.Errorf(expectedEmptyValue, item.Value)
	}
}

// TestErrorResponse testa a estrutura *model.ErrorResponse.
func TestErrorResponse(t *testing.T) {
	errResp := model.ErrorResponse{
		Error: errorTestMessage,
	}

	if errResp.Error != errorTestMessage {
		t.Errorf(expectedErrorMsg, errResp.Error)
	}
}

// TestErrorResponseJSON testa a serialização/desserialização JSON do *model.ErrorResponse.
func TestErrorResponseJSON(t *testing.T) {
	original := model.ErrorResponse{
		Error: errorTestMessage,
	}

	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Fatalf(errMarshalErrorResp+": %v", err)
	}

	var deserialized model.ErrorResponse
	err = json.Unmarshal(jsonData, &deserialized)
	if err != nil {
		t.Fatalf(errUnmarshalErrorResp+": %v", err)
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
		t.Errorf(expectedDados1, len(resp.Dados))
	}

	if resp.Dados[0].Data != dataDate1 {
		t.Errorf(expectedDataDate, dataDate1, resp.Dados[0].Data)
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
		t.Fatalf(errMarshalSolarzResp+": %v", err)
	}

	var deserialized model.SolarzResponse
	err = json.Unmarshal(jsonData, &deserialized)
	if err != nil {
		t.Fatalf(errUnmarshalSolarzResp+": %v", err)
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
		t.Errorf(expectedDataDate, dataDate1, dado.Data)
	}

	if dado.Quantidade != 50.0 {
		t.Errorf(expectedQuantity, 50.0, dado.Quantidade)
	}

	if dado.Prognostico != 48.0 {
		t.Errorf(expectedPrognostico, 48.0, dado.Prognostico)
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
	descricao := testDescription
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
		t.Errorf(expectedGeracoes1, len(dado.Geracoes))
	}

	if dado.Geracoes[0].Quantidade != 25.0 {
		t.Errorf(expectedQuantity25_0, dado.Geracoes[0].Quantidade)
	}

	if *dado.Geracoes[0].Descricao != testDescription {
		t.Errorf(expectedDescription, testDescription, *dado.Geracoes[0].Descricao)
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
		t.Fatalf(errMarshalDadoGeracao+": %v", err)
	}

	var deserialized model.DadoGeracao
	err = json.Unmarshal(jsonData, &deserialized)
	if err != nil {
		t.Fatalf(errUnmarshalDadoGeracao+": %v", err)
	}

	if deserialized.Data != original.Data {
		t.Errorf(expectedDataDate, original.Data, deserialized.Data)
	}

	if deserialized.Denominacao != original.Denominacao {
		t.Errorf(expectedDenominacao, original.Denominacao, deserialized.Denominacao)
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
		t.Errorf(expectedClimaDesc, clima.Descricao)
	}

	if clima.CreatedAt != climaCreated {
		t.Errorf(expectedClimaCreatedAt, clima.CreatedAt)
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
		t.Fatalf(errMarshalClimaInfo+": %v", err)
	}

	var deserialized model.InformacaoClima
	err = json.Unmarshal(jsonData, &deserialized)
	if err != nil {
		t.Fatalf(errUnmarshalClimaInfo+": %v", err)
	}

	if deserialized.ID != original.ID {
		t.Errorf("Expected Id %d, got %d", original.ID, deserialized.ID)
	}

	if deserialized.Descricao != original.Descricao {
		t.Errorf(expectedDescription, original.Descricao, deserialized.Descricao)
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
		t.Errorf(expectedIDExterno, geracao.IDExterno)
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
		t.Errorf(expectedDescription, testDescription, *geracao.Descricao)
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
		t.Fatalf(errMarshalGeracaoDetalhe+": %v", err)
	}

	var deserialized model.GeracaoDetalhe
	err = json.Unmarshal(jsonData, &deserialized)
	if err != nil {
		t.Fatalf(errUnmarshalGeracaoDetalhe+": %v", err)
	}

	if deserialized.Quantidade != original.Quantidade {
		t.Errorf(expectedQuantity, original.Quantidade, deserialized.Quantidade)
	}

	if deserialized.IDExterno != original.IDExterno {
		t.Errorf(expectedIDExterno, deserialized.IDExterno)
	}
}

// TestLabelValue testa a estrutura *model.LabelValue.
func TestLabelValue(t *testing.T) {
	label := model.LabelValue{
		Label: genLabelStr,
		Value: 100.5,
	}

	if label.Label != genLabelStr {
		t.Errorf(expectedLabelValue, label.Label)
	}

	if label.Value != 100.5 {
		t.Errorf(expectedValue100_5, label.Value)
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
		t.Fatalf(errMarshalLabelValue+": %v", err)
	}

	var deserialized model.LabelValue
	err = json.Unmarshal(jsonData, &deserialized)
	if err != nil {
		t.Fatalf(errUnmarshalLabelValue+": %v", err)
	}

	if deserialized.Label != original.Label {
		t.Errorf(expectedLabel, original.Label, deserialized.Label)
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
		t.Errorf(expectedLabeledGen2, len(resp.LabeledGenerations))
	}

	gen1, ok := resp.LabeledGenerations["gen1"]
	if !ok {
		t.Fatal(expectedGenKey)
	}

	if gen1.Value != 50.0 {
		t.Errorf(expectedValue50_0, gen1.Value)
	}
}

// TestSolarzResponseWithPrognosticos testa SolarzResponse com progósticos.
func TestSolarzResponseWithPrognosticos(t *testing.T) {
	resp := model.SolarzResponse{
		Prognosticos: map[string]float64{
			"jan": 100.0,
			"fev": 95.5,
		},
	}

	if len(resp.Prognosticos) != 2 {
		t.Errorf(expectedPrognosticos2, len(resp.Prognosticos))
	}

	janValue, ok := resp.Prognosticos["jan"]
	if !ok {
		t.Fatal(expectedJanKey)
	}

	if janValue != 100.0 {
		t.Errorf(expectedValue100_0, janValue)
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
		t.Fatalf(errMarshalComplexResp+": %v", err)
	}

	var deserialized model.SolarzResponse
	err = json.Unmarshal(jsonData, &deserialized)
	if err != nil {
		t.Fatalf(errUnmarshalComplexResp+": %v", err)
	}

	if len(deserialized.Dados) != 1 {
		t.Errorf(expectedDados1, len(deserialized.Dados))
	}

	if deserialized.Dados[0].Denominacao != usinaDenom {
		t.Errorf(expectedDenominacao, usinaDenom, deserialized.Dados[0].Denominacao)
	}

	if len(deserialized.Dados[0].Geracoes) != 1 {
		t.Errorf(expectedDados, len(deserialized.Dados[0].Geracoes))
	}
}

// TestItemIsValid testa o método IsValid de *model.Item.
func TestItemIsValid(t *testing.T) {
	tests := []struct {
		name     string
		item     *model.Item
		expected bool
	}{
		{testCaseValidItem, &model.Item{ID: "123", Name: "Test TestItemIsValid"}, true},
		{testCaseItemEmptyID, &model.Item{ID: "", Name: "Test TestItemIsValid 2"}, false},
		{testCaseNilItem, nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.item.IsValid()
			if result != tt.expected {
				t.Errorf(expectedVsGot, tt.expected, result)
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
		{testCaseEmptyItem, &model.Item{}, true},
		{testCaseItemWithID, &model.Item{ID: "123"}, false},
		{testCaseItemWithName, &model.Item{Name: "Test TestItemIsEmpty"}, false},
		{testCaseItemWithValue, &model.Item{Value: "Value"}, false},
		{testCaseNilItem, nil, true},
		{testCaseCompleteItem, &model.Item{ID: "123", Name: "Test TestItemIsEmpty 2", Value: "Value"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.item.IsEmpty()
			if result != tt.expected {
				t.Errorf(expectedVsGot, tt.expected, result)
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
		{testCaseWithError, &model.ErrorResponse{Error: expectedWentWrong}, true},
		{testCaseWithEmptyError, &model.ErrorResponse{Error: ""}, false},
		{testCaseNilErrorResponse, nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.errResp.HasError()
			if result != tt.expected {
				t.Errorf(expectedVsGot, tt.expected, result)
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
		{testCaseOneDato, &model.SolarzResponse{Dados: []model.DadoGeracao{{Data: dataDate1}}}, 1},
		{testCaseMultipleDados, &model.SolarzResponse{Dados: []model.DadoGeracao{
			{Data: dataDate1}, {Data: dataDate2}, {Data: dataDate3},
		}}, 3},
		{testCaseNoDados, &model.SolarzResponse{Dados: []model.DadoGeracao{}}, 0},
		{testCaseNilResponse, nil, 0},
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
		{testCaseNormalCalc, &model.DadoGeracao{Quantidade: 100, Prognostico: 100}, 1.0},
		{testCaseAbovePrognosis, &model.DadoGeracao{Quantidade: 120, Prognostico: 100}, 1.2},
		{testCaseBelowPrognosis, &model.DadoGeracao{Quantidade: 80, Prognostico: 100}, 0.8},
		{testCaseZeroPrognosis, &model.DadoGeracao{Quantidade: 100, Prognostico: 0}, 0},
		{testCaseNilDado, nil, 0},
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
		{testCaseManualEntry, &model.DadoGeracao{Manual: true}, true},
		{testCaseAutomaticEntry, &model.DadoGeracao{Manual: false}, false},
		{testCaseNilDado, nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.dado.IsManualEntry()
			if result != tt.expected {
				t.Errorf(expectedVsGot, tt.expected, result)
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
		{testCaseWithDescription, &model.GeracaoDetalhe{Descricao: &descricaoNaoVazia}, true},
		{testCaseWithEmptyDesc, &model.GeracaoDetalhe{Descricao: &descricaoVazia}, false},
		{testCaseWithoutDescription, &model.GeracaoDetalhe{}, false},
		{testCaseNilGeracao, nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.geracao.HasDescription()
			if result != tt.expected {
				t.Errorf(expectedVsGot, tt.expected, result)
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
		{testCaseValidLabel, &model.LabelValue{Label: genLabelStr, Value: 100.5}, true},
		{testCaseEmptyLabel, &model.LabelValue{Label: "", Value: 100.5}, false},
		{testCaseNilLabel, nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.label.IsValid()
			if result != tt.expected {
				t.Errorf(expectedVsGot, tt.expected, result)
			}
		})
	}
}
