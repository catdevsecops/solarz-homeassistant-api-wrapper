# COPILOT_EXAMPLES.md
# Exemplos de Código Seguindo o Estilo do Projeto

Este arquivo contém exemplos práticos de como o Copilot deve gerar código para este projeto.

---

## 1. Exemplos de Handlers

### ✅ CORRETO - Handler com Error Handling

```go
// GetData retorna dados de geração solar da API Solarz
func GetData(w http.ResponseWriter, _ *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    items, err := service.GetData()
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        if encErr := json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Failed to fetch items"}); encErr != nil {
            log.Printf("ERROR: %v", err)
        }
        return
    }

    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(items); err != nil {
        log.Printf("ERROR: %v", err)
    }
}
```

### ❌ INCORRETO - Versão com problemas

```go
// GetData retorna dados
func get_data(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    
    items, err := service.GetData()
    if err != nil {
        // Não faz nada com o erro - BÃO!
        w.Write([]byte("{\"error\": \"failed\"}"))
        return
    }

    data, _ := json.Marshal(items) // Ignora erro - BÃO!
    w.Write(data)
}
```

---

## 2. Exemplos de Struct Methods

### ✅ CORRETO - Método de validação

```go
// IsValid verifica se o Item tem ID válido
func (i *Item) IsValid() bool {
    return i != nil && i.ID != ""
}

// HasError verifica se a resposta contém erro
func (e *ErrorResponse) HasError() bool {
    return e != nil && e.Error != ""
}

// GetTotalDados retorna a quantidade de dados
func (sr *SolarzResponse) GetTotalDados() int {
    if sr == nil {
        return 0
    }
    return len(sr.Dados)
}

// CalculateDesempenho calcula o desempenho baseado em quantidade e prognóstico
func (dg *DadoGeracao) CalculateDesempenho() float64 {
    if dg == nil || dg.Prognostico == 0 {
        return 0
    }
    return dg.Quantidade / dg.Prognostico
}
```

### ❌ INCORRETO - Versão problemática

```go
// Função em package, não método - BÃO!
func isValid(i *Item) bool {
    return i.ID != "" // Não verifica nil - BÃO!
}

// Método retorna erro em vez de tipo simples - BÃO!
func (i *Item) IsValid() (bool, error) {
    return i != nil && i.ID != "", nil
}

// Método com side effects - BÃO!
func (dg *DadoGeracao) CalculateDesempenho() float64 {
    dg.lastCalculated = time.Now() // Side effect!
    return dg.Quantidade / dg.Prognostico
}
```

---

## 3. Exemplos de Testes

### ✅ CORRETO - Unit Test

```go
func TestHealthHandlerSuccess(t *testing.T) {
    req := httptest.NewRequest("GET", "/health", nil)
    w := httptest.NewRecorder()

    healthHandler(w, req)

    if w.Code != http.StatusOK {
        t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
    }

    if ct := w.Header().Get("Content-Type"); ct != "application/json" {
        t.Errorf("expected Content-Type 'application/json', got '%s'", ct)
    }

    var result map[string]string
    if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
        t.Fatalf("failed to decode response: %v", err)
    }

    if result["status"] != "ok" {
        t.Errorf("expected status 'ok', got '%s'", result["status"])
    }
}
```

### ✅ CORRETO - Table-Driven Test

```go
func TestSetupRouterWithDifferentMethods(t *testing.T) {
    tests := []struct {
        name   string
        method string
        path   string
    }{
        {"health GET", "GET", "/health"},
        {"health POST", "POST", "/health"},
        {"data GET", "GET", "/api/v1/data"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            router := setupRouter()
            req := httptest.NewRequest(tt.method, tt.path, nil)
            w := httptest.NewRecorder()

            router.ServeHTTP(w, req)

            if w.Code == http.StatusInternalServerError {
                t.Errorf("%s returned 500", tt.name)
            }
        })
    }
}
```

### ❌ INCORRETO - Teste problemático

```go
// Teste sem assertções claras - BÃO!
func TestHealth(t *testing.T) {
    req := httptest.NewRequest("GET", "/health", nil)
    w := httptest.NewRecorder()
    healthHandler(w, req)
    // Sem verificações!
}

// Teste que faz muitas coisas - BÃO!
func TestEverything(t *testing.T) {
    // Testa handler
    // Testa service
    // Testa model
    // Testa concorrência
    // Não é claro o que está testando!
}
```

---

## 4. Exemplos de Error Handling

### ✅ CORRETO - Error Handling Explícito

```go
// Em um handler
func GetData(w http.ResponseWriter, _ *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    data, err := service.GetData()
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        if encErr := json.NewEncoder(w).Encode(model.ErrorResponse{
            Error: "Failed to fetch data",
        }); encErr != nil {
            log.Printf("Failed to encode error response: %v", encErr)
        }
        return
    }

    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(data); err != nil {
        log.Printf("Failed to encode response: %v", err)
    }
}

// Em um serviço
func GetData() (interface{}, error) {
    if err := validateInput(); err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }

    data, err := fetchData()
    if err != nil {
        return nil, fmt.Errorf("fetch failed: %w", err)
    }

    return data, nil
}
```

### ❌ INCORRETO - Error Handling Ruim

```go
// Ignora erro - BÃO!
func GetData(w http.ResponseWriter, _ *http.Request) {
    data, _ := service.GetData() // Ignora erro!
    w.Write(data)
}

// Retorna nil sem erro - BÃO!
func GetData() (interface{}, error) {
    data, err := fetchData()
    if err != nil {
        return nil, nil // Perdeu o erro!
    }
    return data, nil
}

// Panic em handler - BÃO!
func GetData(w http.ResponseWriter, _ *http.Request) {
    data := service.GetData() // Pode fazer panic!
    w.Write(data)
}
```

---

## 5. Exemplos de Logging

### ✅ CORRETO - Logging Apropriado

```go
// No main
log.Printf("Server starting on %s", server.Addr)

// Em handlers
if err != nil {
    log.Printf("Failed to encode response: %v", err)
}

// Com contexto útil
if err := service.ProcessData(data); err != nil {
    log.Printf("Failed to process data for userID=%s: %v", userID, err)
}
```

### ❌ INCORRETO - Logging Ruim

```go
// Muito verbose - BÃO!
log.Printf("Starting") // Não diz onde
log.Printf("Processing") // Não diz o quê

// Informações sensíveis - BÃO!
log.Printf("Processing user password: %s", password)

// Logging em library - BÃO!
func internalFunction() {
    log.Printf("Called internalFunction") // Poluição de logs
}
```

---

## 6. Exemplos de Estrutura de Dados

### ✅ CORRETO - Struct Bem Definida

```go
// Struct com documentação e tags JSON
type Item struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Value string `json:"value"`
}

type ServerConfig struct {
    Addr         string
    Handler      http.Handler
    ReadTimeout  time.Duration
    WriteTimeout time.Duration
}

type DadoGeracao struct {
    Data          string           `json:"data"`
    Quantidade    float64          `json:"quantidade"`
    Prognostico   float64          `json:"prognostico"`
    Manual        bool             `json:"manual"`
    UsinaId       int              `json:"usinaId"`
    Denominacao   string           `json:"denominacao"`
    Geracoes      []GeracaoDetalhe `json:"geracoes"`
    PlantShutdown bool             `json:"plantShutdown"`
}
```

### ❌ INCORRETO - Struct Problemática

```go
// Sem tags JSON - BÃO!
type Item struct {
    ID    string
    Name  string
    Value string
}

// Nomes inconsistentes - BÃO!
type item struct {
    id    string
    name  string
    value string
}

// Campos privados para API pública - BÃO!
type APIResponse struct {
    data   interface{} // Deve ser Data
    error  error       // Deve ser Error
}
```

---

## 7. Exemplos de Imports

### ✅ CORRETO - Imports Organizados

```go
package main

import (
    "encoding/json"
    "log"
    "net/http"
    "time"

    "github.com/catdevsecops/solarz-api/internal/handler"
    "github.com/catdevsecops/solarz-api/internal/model"
    "github.com/catdevsecops/solarz-api/internal/service"
)
```

### ❌ INCORRETO - Imports Bagunçados

```go
// Imports fora de ordem - BÃO!
import (
    "github.com/catdevsecops/solarz-api/internal/handler"
    "encoding/json"
    "github.com/catdevsecops/solarz-api/internal/model"
    "log"
)

// Imports não utilizados - BÃO!
import (
    "crypto/sha256"
    "database/sql"
    "log"
)
```

---

## 8. Exemplos de Factory Pattern

### ✅ CORRETO - Factory para Configuração

```go
// Factory que retorna object configurado
func setupRouter() *http.ServeMux {
    router := http.NewServeMux()
    router.HandleFunc("GET /api/v1/data", handler.GetData)
    router.HandleFunc("GET /health", healthHandler)
    return router
}

// Uso em main
func main() {
    router := setupRouter()
    
    config := ServerConfig{
        Addr:    ":8080",
        Handler: router,
    }
    
    startServer(config)
}
```

### ❌ INCORRETO - Sem Factory

```go
// Lógica de setup em main - BÃO!
func main() {
    router := http.NewServeMux()
    router.HandleFunc("GET /api/v1/data", handler.GetData)
    router.HandleFunc("GET /health", healthHandler)
    
    server := &http.Server{
        Addr:    ":8080",
        Handler: router,
    }
    
    server.ListenAndServe()
}
```

---

## 9. Exemplos de Documentação

### ✅ CORRETO - Documentação Clara

```go
// GetData retorna dados de geração solar da API Solarz
func GetData(w http.ResponseWriter, _ *http.Request) {
    // Implementação
}

// IsValid verifica se o Item tem ID válido
func (i *Item) IsValid() bool {
    return i != nil && i.ID != ""
}

// setupRouter cria e configura o router HTTP com todas as rotas
func setupRouter() *http.ServeMux {
    // Implementação
}
```

### ❌ INCORRETO - Documentação Ruim

```go
// getData (não começa com nome - BÃO!)
func GetData() { ... }

// função que faz coisas (vago demais - BÃO!)
func Process() { ... }

// TODO: implementar (não é documentação - BÃO!)
func DoSomething() { ... }
```

---

## 10. Checklist para Code Review

Quando o Copilot gera código, verificar:

- [ ] Segue padrões de nomenclatura
- [ ] Tem comentários nas funções públicas
- [ ] Error handling está explícito (nunca swallow)
- [ ] Testes estão inclusos
- [ ] Sem variáveis globais
- [ ] JSON tags em structs públicas
- [ ] Usa json.NewEncoder (não Marshal)
- [ ] Logging apropriado (não poluído)
- [ ] Formatação OK (tabs, não espaços)
- [ ] Sem dependências externas

---

**Uso**: Use esses exemplos para validar se o Copilot está gerando código no estilo do projeto.

**Data**: 2024-01-05
**Versão**: 1.0
