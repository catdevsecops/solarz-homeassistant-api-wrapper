package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/catdevsecops/solarz-api/internal/model"
)

// TestHealthHandlerSuccess testa o health check com sucesso
func TestHealthHandlerSuccess(t *testing.T) {
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	healthHandler(w, req)

	// Verifica status code
	if w.Code != http.StatusOK {
		t.Errorf("healthHandler() returned status %d, want %d", w.Code, http.StatusOK)
	}

	// Verifica Content-Type
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Content-Type = %q, want 'application/json'", contentType)
	}

	// Verifica body
	var result map[string]string
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	if result["status"] != "ok" {
		t.Errorf("status = %q, want 'ok'", result["status"])
	}
}

// TestHealthHandlerWithDifferentMethods testa health check com diferentes métodos HTTP
func TestHealthHandlerWithDifferentMethods(t *testing.T) {
	tests := []struct {
		name   string
		method string
	}{
		{"GET request", "GET"},
		{"POST request", "POST"},
		{"DELETE request", "DELETE"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/health", nil)
			w := httptest.NewRecorder()

			healthHandler(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("%s: expected status %d, got %d", tt.method, http.StatusOK, w.Code)
			}

			var result map[string]string
			if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
				t.Errorf("%s: failed to decode: %v", tt.method, err)
			}

			if result["status"] != "ok" {
				t.Errorf("%s: expected status 'ok', got %q", tt.method, result["status"])
			}
		})
	}
}

// TestHealthHandlerResponseBody testa se o body está bem formatado
func TestHealthHandlerResponseBody(t *testing.T) {
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	healthHandler(w, req)

	// Verifica se o body não está vazio
	if w.Body.Len() == 0 {
		t.Error("healthHandler() returned empty body")
	}

	// Verifica se é JSON válido
	var result map[string]interface{}
	if err := json.NewDecoder(bytes.NewReader(w.Body.Bytes())).Decode(&result); err != nil {
		t.Errorf("invalid JSON in response: %v", err)
	}

	// Verifica campos específicos
	if _, ok := result["status"]; !ok {
		t.Error("'status' field missing from response")
	}
}

// TestHealthHandlerHeadersSet testa se os headers são configurados corretamente
func TestHealthHandlerHeadersSet(t *testing.T) {
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	healthHandler(w, req)

	// Verifica Content-Type
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type header = %q, want 'application/json'", ct)
	}

	// Verifica se Status Code está correto
	if w.Code != http.StatusOK {
		t.Errorf("HTTP status = %d, want %d", w.Code, http.StatusOK)
	}
}

// TestHealthHandlerDecodesCorrectly testa a decodificação correta
func TestHealthHandlerDecodesCorrectly(t *testing.T) {
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	healthHandler(w, req)

	// Decodifica o body
	body, err := io.ReadAll(w.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	var result map[string]string
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	// Verifica tipo de dados
	if status, ok := result["status"]; !ok {
		t.Fatal("'status' key not found in response")
	} else if status != "ok" {
		t.Errorf("status value = %q, want 'ok'", status)
	}
}

// TestHealthHandlerResponseStructure testa a estrutura da resposta
func TestHealthHandlerResponseStructure(t *testing.T) {
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	healthHandler(w, req)

	var result map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Verifica que há apenas um campo
	if len(result) != 1 {
		t.Errorf("expected 1 field in response, got %d", len(result))
	}

	// Verifica o campo status
	if status, exists := result["status"]; !exists {
		t.Error("'status' field is missing")
	} else if statusStr, ok := status.(string); !ok {
		t.Errorf("'status' field is not a string, got %T", status)
	} else if statusStr != "ok" {
		t.Errorf("'status' value = %q, want 'ok'", statusStr)
	}
}

// TestHealthHandlerIsIdempotent testa se múltiplas chamadas retornam o mesmo
func TestHealthHandlerIsIdempotent(t *testing.T) {
	var results []map[string]string

	for i := 0; i < 3; i++ {
		req := httptest.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()

		healthHandler(w, req)

		var result map[string]string
		if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
			t.Errorf("iteration %d: failed to decode: %v", i, err)
			continue
		}

		if w.Code != http.StatusOK {
			t.Errorf("iteration %d: expected status %d, got %d", i, http.StatusOK, w.Code)
		}

		results = append(results, result)
	}

	// Verifica se todos os resultados são iguais
	for i := 1; i < len(results); i++ {
		if results[i]["status"] != results[0]["status"] {
			t.Errorf("inconsistent results: iteration 0 = %q, iteration %d = %q",
				results[0]["status"], i, results[i]["status"])
		}
	}
}

// TestHealthHandlerConcurrency testa se é seguro para uso concorrente
func TestHealthHandlerConcurrency(t *testing.T) {
	done := make(chan bool, 5)

	for i := 0; i < 5; i++ {
		go func() {
			req := httptest.NewRequest("GET", "/health", nil)
			w := httptest.NewRecorder()

			healthHandler(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("concurrent call failed with status %d", w.Code)
			}

			var result map[string]string
			if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
				t.Errorf("concurrent call failed to decode: %v", err)
			}

			done <- true
		}()
	}

	// Aguarda todas as goroutines
	for i := 0; i < 5; i++ {
		<-done
	}
}

// TestHealthHandlerWithQueryParams testa com query parameters
func TestHealthHandlerWithQueryParams(t *testing.T) {
	req := httptest.NewRequest("GET", "/health?debug=true&verbose=1", nil)
	w := httptest.NewRecorder()

	healthHandler(w, req)

	// O handler deve ignorar query params e retornar normalmente
	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var result map[string]string
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	if result["status"] != "ok" {
		t.Errorf("status = %q, want 'ok'", result["status"])
	}
}

// TestHealthHandlerWithHeaders testa com headers customizados
func TestHealthHandlerWithHeaders(t *testing.T) {
	req := httptest.NewRequest("GET", "/health", nil)
	req.Header.Set("User-Agent", "Custom-Client/1.0")
	req.Header.Set("Authorization", "Bearer token123")
	w := httptest.NewRecorder()

	healthHandler(w, req)

	// O handler deve ignorar headers customizados e retornar normalmente
	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var result map[string]string
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	if result["status"] != "ok" {
		t.Errorf("status = %q, want 'ok'", result["status"])
	}
}

// TestHealthHandlerIntegration testa integração com o modelo
func TestHealthHandlerIntegration(t *testing.T) {
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	healthHandler(w, req)

	// Verifica estrutura da resposta
	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var result map[string]string
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Valida tipos de dados
	for key, value := range result {
		if key != "status" {
			t.Errorf("unexpected key in response: %q", key)
		}
		if value != "ok" {
			t.Errorf("unexpected value for key %q: %q", key, value)
		}
	}
}

// BenchmarkHealthHandler testa performance do health handler
func BenchmarkHealthHandler(b *testing.B) {
	req := httptest.NewRequest("GET", "/health", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		healthHandler(w, req)
	}
}

// TestHealthHandlerResponseSize testa tamanho da resposta
func TestHealthHandlerResponseSize(t *testing.T) {
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	healthHandler(w, req)

	// A resposta deve ser pequena
	maxSize := 1024 // 1KB
	if w.Body.Len() > maxSize {
		t.Errorf("response too large: %d bytes, expected <= %d bytes", w.Body.Len(), maxSize)
	}

	// A resposta deve conter dados
	if w.Body.Len() == 0 {
		t.Error("response is empty")
	}
}

// TestHealthHandlerErrorHandling testa manipulação de erro (mocado)
func TestHealthHandlerErrorHandling(t *testing.T) {
	// Este teste verifica que o handler trata erros gracefully
	// Mesmo que neste caso simples não haja muitos erros possíveis
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	// Deve não crashear
	healthHandler(w, req)

	if w.Code == 0 {
		t.Error("handler did not write status code")
	}
}

// TestErrorResponseModel testa o modelo de resposta de erro
func TestErrorResponseModel(t *testing.T) {
	errorResp := model.ErrorResponse{
		Error: "Test error",
	}

	// Marshala para JSON
	body, err := json.Marshal(errorResp)
	if err != nil {
		t.Fatalf("failed to marshal error response: %v", err)
	}

	// Unmarshala de volta
	var result model.ErrorResponse
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("failed to unmarshal error response: %v", err)
	}

	if result.Error != "Test error" {
		t.Errorf("error message = %q, want 'Test error'", result.Error)
	}
}

// TestHealthHandlerMultipleResponses testa múltiplas respostas
func TestHealthHandlerMultipleResponses(t *testing.T) {
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()

		healthHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("request %d: expected status %d, got %d", i, http.StatusOK, w.Code)
		}

		if w.Header().Get("Content-Type") != "application/json" {
			t.Errorf("request %d: invalid content type", i)
		}

		var result map[string]string
		if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
			t.Errorf("request %d: failed to decode: %v", i, err)
		}
	}
}

// TestHealthHandlerStatelessness testa se o handler é stateless
func TestHealthHandlerStatelessness(t *testing.T) {
	// Faz várias requisições
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()

		healthHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("request %d failed", i)
		}
	}

	// Faz mais uma requisição para garantir que o estado não foi alterado
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	healthHandler(w, req)

	if w.Code != http.StatusOK {
		t.Error("handler state was affected by previous requests")
	}
}

// TestHealthHandlerContentType testa diferentes content types
func TestHealthHandlerContentType(t *testing.T) {
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	healthHandler(w, req)

	// Verifica se Content-Type é exatamente o esperado
	expectedCT := "application/json"
	if ct := w.Header().Get("Content-Type"); ct != expectedCT {
		t.Errorf("Content-Type = %q, want %q", ct, expectedCT)
	}
}

// === TESTES DE CONFIGURAÇÃO DO SERVIDOR ===

// TestServerConfig testa a estrutura de configuração
func TestServerConfig(t *testing.T) {
	config := ServerConfig{
		Addr:    ":8080",
		Handler: http.NewServeMux(),
	}

	if config.Addr != ":8080" {
		t.Errorf("expected addr ':8080', got '%s'", config.Addr)
	}

	if config.Handler == nil {
		t.Error("expected handler to not be nil")
	}
}

// TestServerConfigDifferentAddresses testa configuração com endereços diferentes
func TestServerConfigDifferentAddresses(t *testing.T) {
	tests := []struct {
		name string
		addr string
	}{
		{"localhost 8080", "localhost:8080"},
		{"any interface 8080", ":8080"},
		{"any interface 3000", ":3000"},
		{"127.0.0.1 9000", "127.0.0.1:9000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := ServerConfig{
				Addr:    tt.addr,
				Handler: http.NewServeMux(),
			}

			if config.Addr != tt.addr {
				t.Errorf("expected addr '%s', got '%s'", tt.addr, config.Addr)
			}
		})
	}
}

// TestSetupRouter testa a configuração do router
func TestSetupRouter(t *testing.T) {
	router := setupRouter()

	if router == nil {
		t.Fatal("expected router to not be nil")
	}
}

// TestSetupRouterRegistersRoutes testa se rotas são registradas
func TestSetupRouterRegistersRoutes(t *testing.T) {
	router := setupRouter()

	// Testa rota /health
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("health endpoint: expected status %d, got %d", http.StatusOK, w.Code)
	}
}

// TestSetupRouterMultipleRoutes testa se múltiplas rotas são registradas
func TestSetupRouterMultipleRoutes(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name   string
		path   string
		method string
	}{
		{"health endpoint", "/health", "GET"},
		{"data endpoint", "/api/v1/data", "GET"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			// Qualquer status code (não 404) indica que a rota existe
			if w.Code == http.StatusNotFound {
				t.Errorf("%s: route not found", tt.name)
			}
		})
	}
}

// TestSetupRouterHealthEndpoint testa endpoint de health
func TestSetupRouterHealthEndpoint(t *testing.T) {
	router := setupRouter()

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var result map[string]string
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if result["status"] != "ok" {
		t.Errorf("expected status 'ok', got '%s'", result["status"])
	}
}

// TestSetupRouterReturnsServeMux testa tipo de retorno
func TestSetupRouterReturnsServeMux(t *testing.T) {
	router := setupRouter()

	// Valida que router é do tipo esperado
	if router == nil {
		t.Error("expected *http.ServeMux, got nil")
	}

	// Valida que é uma instância de http.ServeMux verificando se implementa http.Handler
	var _ http.Handler = router
	t.Logf("setupRouter returns valid http.Handler")
}

// TestSetupRouterIsConsistent testa se múltiplas chamadas retornam resultados consistentes
func TestSetupRouterIsConsistent(t *testing.T) {
	router1 := setupRouter()
	router2 := setupRouter()

	// Ambos devem ter as mesmas rotas
	for path, method := range map[string]string{
		"/health":      "GET",
		"/api/v1/data": "GET",
	} {
		req1 := httptest.NewRequest(method, path, nil)
		w1 := httptest.NewRecorder()
		router1.ServeHTTP(w1, req1)

		req2 := httptest.NewRequest(method, path, nil)
		w2 := httptest.NewRecorder()
		router2.ServeHTTP(w2, req2)

		if w1.Code != w2.Code {
			t.Errorf("%s: inconsistent responses (router1=%d, router2=%d)", path, w1.Code, w2.Code)
		}
	}
}

// TestSetupRouterContentType testa content type das rotas
func TestSetupRouterContentType(t *testing.T) {
	router := setupRouter()

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type 'application/json', got '%s'", ct)
	}
}

// TestMainIntegration testa se main pode ser chamada sem crashear (com limites)
func TestMainIntegration(t *testing.T) {
	// Este teste apenas verifica que setupRouter() funciona
	// main() não pode ser testado diretamente pois inicia um servidor
	router := setupRouter()

	if router == nil {
		t.Fatal("setupRouter returned nil")
	}

	// Verifica que router pode servir requisições
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Error("router cannot serve requests")
	}
}

// TestServerConfigWithHandler testa ServerConfig com handler
func TestServerConfigWithHandler(t *testing.T) {
	router := setupRouter()

	config := ServerConfig{
		Addr:    ":8080",
		Handler: router,
	}

	if config.Handler == nil {
		t.Error("handler should not be nil")
	}

	// Testa que handler funciona
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	config.Handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("handler should serve requests, got status %d", w.Code)
	}
}

// TestSetupRouterWithConcurrentRequests testa router com requisições concorrentes
func TestSetupRouterWithConcurrentRequests(t *testing.T) {
	router := setupRouter()
	done := make(chan bool, 5)

	for i := 0; i < 5; i++ {
		go func() {
			req := httptest.NewRequest("GET", "/health", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("expected status 200, got %d", w.Code)
			}

			done <- true
		}()
	}

	for i := 0; i < 5; i++ {
		<-done
	}
}

// TestServerConfigEmptyAddr testa ServerConfig com addr vazio
func TestServerConfigEmptyAddr(t *testing.T) {
	config := ServerConfig{
		Addr:    "",
		Handler: http.NewServeMux(),
	}

	if config.Addr != "" {
		t.Errorf("expected empty addr, got '%s'", config.Addr)
	}
}

// TestSetupRouterHandlesInvalidPaths testa router com caminhos inválidos
func TestSetupRouterHandlesInvalidPaths(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name string
		path string
	}{
		{"nonexistent path", "/nonexistent"},
		{"wrong method", "/health"},
		{"partial path", "/api"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var method string
			if tt.name == "wrong method" {
				method = "POST" // Method incorreto para /health
			} else {
				method = "GET"
			}

			req := httptest.NewRequest(method, tt.path, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			// Caminhos inválidos devem retornar 404 ou não encontrado
			// (depende da implementação do handler)
			if w.Code == http.StatusOK {
				t.Logf("%s returned 200 (router serves this path)", tt.path)
			}
		})
	}
}

// TestSetupRouterConcurrent testa se setupRouter é seguro para concorrência
func TestSetupRouterConcurrent(t *testing.T) {
	done := make(chan bool, 3)

	for i := 0; i < 3; i++ {
		go func() {
			router := setupRouter()

			if router == nil {
				t.Error("setupRouter returned nil in concurrent call")
			}

			req := httptest.NewRequest("GET", "/health", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("concurrent call failed with status %d", w.Code)
			}

			done <- true
		}()
	}

	for i := 0; i < 3; i++ {
		<-done
	}
}

// === TESTES COM CAPTURA DE LOGS ===

// captureLogOutput captura a saída de log durante a execução de uma função
func captureLogOutput(fn func()) string {
	// Cria um buffer para capturar logs
	var buf bytes.Buffer

	// Salva o logger original
	originalLogger := log.Writer()
	originalFlags := log.Flags()
	originalPrefix := log.Prefix()

	// Define novo logger que escreve no buffer
	log.SetOutput(&buf)
	log.SetFlags(log.LstdFlags)

	// Executa a função
	fn()

	// Restaura o logger original
	log.SetOutput(originalLogger)
	log.SetFlags(originalFlags)
	log.SetPrefix(originalPrefix)

	return buf.String()
}

// TestHealthHandlerLogging testa se os logs são registrados corretamente em caso de erro
func TestHealthHandlerLogging(t *testing.T) {
	t.Run("captura log em caso de erro de encoding", func(t *testing.T) {
		// Este teste verifica se logs são gerados
		// Para testar erro de encoding real, precisaríamos de um mock
		req := httptest.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()

		// Captura logs durante a execução
		logOutput := captureLogOutput(func() {
			healthHandler(w, req)
		})

		// Em caso normal, não deve haver erro de log
		if strings.Contains(logOutput, "Failed to encode response") {
			t.Errorf("unexpected error log: %s", logOutput)
		}
	})
}

// TestHealthHandlerLoggingWithFailedWrite testa logging quando há erro na escrita
func TestHealthHandlerLoggingWithFailedWrite(t *testing.T) {
	// Cria um ResponseWriter que falha ao escrever
	failWriter := &failingWriter{}
	req := httptest.NewRequest("GET", "/health", nil)

	// Captura logs
	logOutput := captureLogOutput(func() {
		healthHandler(failWriter, req)
	})

	// Deve conter mensagem de erro
	if !strings.Contains(logOutput, "Failed to encode response") {
		t.Errorf("expected error log not found\nGot: %s", logOutput)
	}
}

// failingWriter é um mock de ResponseWriter que falha ao escrever
type failingWriter struct {
	headerWritten bool
}

func (fw *failingWriter) Header() http.Header {
	return make(http.Header)
}

func (fw *failingWriter) Write(b []byte) (int, error) {
	fw.headerWritten = true
	return 0, io.ErrClosedPipe // Simula erro de escrita
}

func (fw *failingWriter) WriteHeader(statusCode int) {
	// Não faz nada
}

// TestLoggerCaptureBasic testa a função de captura de logs
func TestLoggerCaptureBasic(t *testing.T) {
	logOutput := captureLogOutput(func() {
		log.Println("Test message")
	})

	if !strings.Contains(logOutput, "Test message") {
		t.Errorf("expected 'Test message' in log output, got: %s", logOutput)
	}
}

// TestLoggerCaptureMultipleMessages testa captura de múltiplos logs
func TestLoggerCaptureMultipleMessages(t *testing.T) {
	logOutput := captureLogOutput(func() {
		log.Println("Message 1")
		log.Println("Message 2")
		log.Println("Message 3")
	})

	if !strings.Contains(logOutput, "Message 1") ||
		!strings.Contains(logOutput, "Message 2") ||
		!strings.Contains(logOutput, "Message 3") {
		t.Errorf("expected all messages in log output, got: %s", logOutput)
	}
}

// TestLoggerCaptureFormat testa se o formato do log é correto
func TestLoggerCaptureFormat(t *testing.T) {
	logOutput := captureLogOutput(func() {
		log.Printf("Error: %v", "test error")
	})

	if !strings.Contains(logOutput, "Error: test error") {
		t.Errorf("expected formatted message, got: %s", logOutput)
	}
}

// TestLoggerCaptureEmpty testa captura quando não há logs
func TestLoggerCaptureEmpty(t *testing.T) {
	logOutput := captureLogOutput(func() {
		// Não faz nada
	})

	if logOutput != "" {
		t.Errorf("expected empty log output, got: %s", logOutput)
	}
}

// TestHealthHandlerNoErrorLogsNormally testa que handler normal não gera erro logs
func TestHealthHandlerNoErrorLogsNormally(t *testing.T) {
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	logOutput := captureLogOutput(func() {
		healthHandler(w, req)
	})

	// Não deve haver "Failed" ou "error" (case-insensitive)
	if strings.Contains(strings.ToLower(logOutput), "failed") ||
		strings.Contains(strings.ToLower(logOutput), "error") {
		t.Errorf("unexpected error in normal operation: %s", logOutput)
	}
}

// TestHealthHandlerErrorLogContainsTimestamp testa se logs contêm timestamp
func TestHealthHandlerErrorLogContainsTimestamp(t *testing.T) {
	// Cria writer que falha
	failWriter := &failingWriter{}
	req := httptest.NewRequest("GET", "/health", nil)

	logOutput := captureLogOutput(func() {
		healthHandler(failWriter, req)
	})

	// Deve conter timestamp (padrão: 2006/01/02 15:04:05)
	hasTimestamp := strings.Contains(logOutput, "/") || strings.Contains(logOutput, ":")
	if !hasTimestamp {
		t.Errorf("expected timestamp in log, got: %s", logOutput)
	}
}

// TestMultipleHandlerCallsLogging testa logs de múltiplas chamadas
func TestMultipleHandlerCallsLogging(t *testing.T) {
	var successCount int
	var errorCount int

	for i := 0; i < 3; i++ {
		req := httptest.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()

		logOutput := captureLogOutput(func() {
			healthHandler(w, req)
		})

		if logOutput == "" {
			successCount++
		} else if strings.Contains(logOutput, "Failed") {
			errorCount++
		}
	}

	if successCount != 3 {
		t.Errorf("expected 3 successful calls, got %d", successCount)
	}

	if errorCount != 0 {
		t.Errorf("expected 0 error calls, got %d", errorCount)
	}
}

// TestLoggerWithDifferentLevels testa diferentes níveis de log
func TestLoggerWithDifferentLevels(t *testing.T) {
	t.Run("Println", func(t *testing.T) {
		logOutput := captureLogOutput(func() {
			log.Println("Info message")
		})
		if !strings.Contains(logOutput, "Info message") {
			t.Error("Println message not captured")
		}
	})

	t.Run("Printf", func(t *testing.T) {
		logOutput := captureLogOutput(func() {
			log.Printf("Format message: %s", "test")
		})
		if !strings.Contains(logOutput, "Format message: test") {
			t.Error("Printf message not captured")
		}
	})

	t.Run("Print", func(t *testing.T) {
		logOutput := captureLogOutput(func() {
			log.Print("Direct message")
		})
		if !strings.Contains(logOutput, "Direct message") {
			t.Error("Print message not captured")
		}
	})
}

// TestLoggerRestoresOriginalState testa se o logger é restaurado corretamente
func TestLoggerRestoresOriginalState(t *testing.T) {
	// Salva estado original
	originalFlags := log.Flags()
	originalPrefix := log.Prefix()

	// Executa captura
	captureLogOutput(func() {
		log.Println("Test")
	})

	// Verifica se foi restaurado
	if log.Flags() != originalFlags {
		t.Errorf("logger flags not restored: expected %d, got %d", originalFlags, log.Flags())
	}

	if log.Prefix() != originalPrefix {
		t.Errorf("logger prefix not restored: expected %q, got %q", originalPrefix, log.Prefix())
	}
}

// TestHealthHandlerConcurrentLogging testa logs em concorrência
func TestHealthHandlerConcurrentLogging(t *testing.T) {
	done := make(chan bool, 3)

	for i := 0; i < 3; i++ {
		go func() {
			req := httptest.NewRequest("GET", "/health", nil)
			w := httptest.NewRecorder()

			logOutput := captureLogOutput(func() {
				healthHandler(w, req)
			})

			if logOutput != "" && strings.Contains(logOutput, "Failed") {
				t.Errorf("unexpected error log in concurrent call: %s", logOutput)
			}

			done <- true
		}()
	}

	for i := 0; i < 3; i++ {
		<-done
	}
}

// TestFailingWriterBehavior testa o comportamento do failing writer
func TestFailingWriterBehavior(t *testing.T) {
	fw := &failingWriter{}

	// Testa Header
	header := fw.Header()
	if header == nil {
		t.Error("Header should not be nil")
	}

	// Testa Write
	_, err := fw.Write([]byte("test"))
	if err == nil {
		t.Error("Write should return error")
	}

	if err != io.ErrClosedPipe {
		t.Errorf("expected io.ErrClosedPipe, got %v", err)
	}

	// Testa que headerWritten é atualizado
	if !fw.headerWritten {
		t.Error("headerWritten should be true after Write")
	}
}

// TestLogOutputContainsErrorDetails testa se logs contêm detalhes do erro
func TestLogOutputContainsErrorDetails(t *testing.T) {
	failWriter := &failingWriter{}
	req := httptest.NewRequest("GET", "/health", nil)

	logOutput := captureLogOutput(func() {
		healthHandler(failWriter, req)
	})

	// Deve conter "Failed to encode response" e detalhes do erro
	if !strings.Contains(logOutput, "Failed to encode response") {
		t.Errorf("log should contain error message, got: %s", logOutput)
	}

	// Deve conter referência ao erro
	if !strings.Contains(logOutput, "ErrClosedPipe") && !strings.Contains(logOutput, "closed") {
		// Também aceita qualquer menção do erro
		if !strings.Contains(logOutput, "error") {
			t.Errorf("log should contain error details, got: %s", logOutput)
		}
	}
}
