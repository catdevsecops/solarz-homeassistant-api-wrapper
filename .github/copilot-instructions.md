# Copilot Instructions - Guia de Estilo do Projeto Solarz API

## 📋 Visão Geral

Este documento orienta agentes de IA (GitHub Copilot, Claude, etc.) sobre como manter a consistência de estilo e padrões no projeto **Solarz API**.

**Projeto**: Solarz API (API REST em Go)
**Versão**: 1.0
**Status**: Production-Ready
**Cobertura de Testes**: ~81% média

---

## 🎯 Princípios Fundamentais

1. **Go Idiomático**: Seguir convenções padrão de Go
2. **Sem Dependências Externas**: Apenas biblioteca padrão do Go
3. **Camadas Bem Definidas**: Handler → Service → Model
4. **Testes Abrangentes**: Mínimo 80% de cobertura
5. **Documentação Clara**: Comentários explicam "por quê"
6. **Padrões de Design**: DI, Factory, etc.
7. **Error Handling Explícito**: Nunca swallow errors
8. **Logging Apropriado**: Informações úteis sem poluição

---

## 📁 Estrutura do Projeto

```
solarz-homeassistant-api-wrapper/
├── main.go                     # Entry point (refatorado)
├── main_test.go                # 62 testes
├── go.mod / go.sum            # Dependências
├── README.md                   # Documentação
├── .copilot-instructions       # Este arquivo
└── internal/
    ├── handler/
    │   ├── item.go             # HTTP handlers
    │   └── item_test.go        # Testes de handler
    ├── model/
    │   ├── item.go             # Estruturas de dados
    │   └── item_test.go        # 38 testes (100% cobertura)
    └── service/
        ├── item.go             # Lógica de negócio
        └── item_test.go        # Testes de serviço
```

### Convenções de Diretório

- **`internal/`**: Código privado do módulo (sem exportação)
- **`handler/`**: HTTP request handlers (camada HTTP)
- **`service/`**: Business logic (camada de lógica)
- **`model/`**: Data structures (camada de dados)

---

## 🎯 Padrões de Nomeclatura de Variáveis

### Regra Fundamental: Legibilidade

**IMPORTANTE**: Variáveis devem ser fáceis de ler e entender. Nunca use nomes curtos ou abreviados (exceto para parâmetros padrão do Go).

### Parâmetros Padrão do Go (Aceitos)

Os únicos nomes curtos permitidos são parâmetros padrão amplamente reconhecidos:

| Parâmetro | Significado | Por quê | Aceito |
|-------------|-------------|---------|--------|
| `w` | http.ResponseWriter | Convenção de Go | ✅ |
| `r` | *http.Request | Convenção de Go | ✅ |
| `err` | error | Convenção de Go | ✅ |
| `i`, `j`, `k` | Loop counters | Apenas em loops simples | ✅ |
| `ctx` | context.Context | Convenção de Go | ✅ |
| `t` | *testing.T | Convenção de Go | ✅ |

### Variáveis Locais: Sempre Descritivas

```go
// ✅ CORRETO - Nomes descritivos
responseRecorder := httptest.NewRecorder()
logOutput := captureLogOutput(func() { ... })
originalLogger := log.Writer()
serverConfiguration := ServerConfig{...}
success Count := 0
userAuthenticationToken := "Bearer token123"

// ❌ INCORRETO - Nomes curtos
rec := httptest.NewRecorder()
out := captureLogOutput(func() { ... })
orig := log.Writer()
cfg := ServerConfig{...}
success := 0
token := "Bearer token123"
```

### Regras Específicas

1. **Variáveis de Teste**: Sempre descritivas
   ```go
   // ✅ CORRETO
   testCaseName := "GET request"
   expectedStatusCode := http.StatusOK
   responseRecorder := httptest.NewRecorder()
   
   // ❌ INCORRETO
   name := "GET request"
   status := http.StatusOK
   rec := httptest.NewRecorder()
   ```

2. **Variáveis em Loops**: Use índice só em loops de range simples
   ```go
   // ✅ CORRETO
   for iteration := range 5 { ... }  // ou for i := range 5
   for index, testCase := range testCases { ... }
   
   // ❌ INCORRETO
   for i := 0; i < len(items); i++ { ...  // use range
   ```

3. **Variáveis de Estado/Config**: Completamente descritivas
   ```go
   // ✅ CORRETO
   originalHttpHandler := log.Writer()
   capturedLogBuffer := &bytes.Buffer{}
   requestHeader := req.Header
   
   // ❌ INCORRETO
   h := log.Writer()
   buf := &bytes.Buffer{}
   hdr := req.Header
   ```

---

| Uso | Correto | Incorreto |
|-----|---------|----------|
| HTTP Response | `w` | `response_writer`, `rw` |
| HTTP Request | `r` | `request`, `req` |
| Erros | `err` | `error_value`, `e` |
| Configuração | `config` | `cfg`, `conf` |
| Router | `router` | `mux`, `r` |
| Context | `ctx` | `context`, `c` |

### Funções e Métodos

```go
// ✅ CORRETO
func GetData(w http.ResponseWriter, r *http.Request)
func setupRouter() *http.ServeMux
func startServer(config ServerConfig) error

// ❌ INCORRETO
func get_data(w http.ResponseWriter, r *http.Request)
func Setup_Router() *http.ServeMux
func StartServer(config ServerConfig) error
```

### Structs

```go
// ✅ CORRETO
type ServerConfig struct { ... }
type SolarzResponse struct { ... }
type DadoGeracao struct { ... }

// ❌ INCORRETO
type serverConfig struct { ... }
type solarzResponse struct { ... }
type dadoGeracao struct { ... }
```

### Métodos em Structs

```go
// ✅ CORRETO - Nomes começam com ação
func (i *Item) IsValid() bool { ... }
func (sr *SolarzResponse) GetTotalDados() int { ... }
func (dg *DadoGeracao) CalculateDesempenho() float64 { ... }

// ❌ INCORRETO
func (i *Item) isValid() bool { ... }
func (sr *SolarzResponse) getTotalDados() int { ... }
func (dg *DadoGeracao) calculateDesempenho() float64 { ... }
```

---

## 📝 Padrões de Código

### Handler Functions

```go
// ✅ PADRÃO A SEGUIR
func healthHandler(w http.ResponseWriter, _ *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(map[string]string{"status": "ok"}); err != nil {
        log.Printf("Failed to encode response: %v", err)
    }
}
```

**Características Obrigatórias**:
1. Parâmetros: `(w http.ResponseWriter, r *http.Request)`
2. Ignorar param não usado: `_`
3. Setar Content-Type antes de WriteHeader
4. Fazer WriteHeader explicitamente
5. Usar `json.NewEncoder` (não Marshal)
6. Fazer log de erros com `log.Printf`
7. Nunca panic ou retornar erro (está no response)

### Struct Methods - Validação

```go
// ✅ PADRÃO
func (i *Item) IsValid() bool {
    return i != nil && i.ID != ""
}

func (sr *SolarzResponse) GetTotalDados() int {
    if sr == nil {
        return 0
    }
    return len(sr.Dados)
}
```

**Características**:
- Sempre verificar nil primeiro
- Nomes: `Is*`, `Has*`, `Get*`, `Calculate*`
- Retornar tipos simples (bool, int, string, float64)
- Sem side effects
- Receptores são pointers (`*Type`)

### Error Handling

```go
// ✅ PADRÃO
if err := someFunction(); err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    if encErr := json.NewEncoder(w).Encode(model.ErrorResponse{
        Error: "Failed to fetch items",
    }); encErr != nil {
        log.Printf("ERROR: %v", err)
    }
    return
}
```

**Regras**:
1. Sempre verificar `err != nil`
2. Retornar após erro (não continuar)
3. Usar `model.ErrorResponse` para respostas
4. Usar status codes apropriados
5. Fazer log de contexto importante

### Struct Definitions

```go
// ✅ PADRÃO
type Item struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Value string `json:"value"`
}

type ServerConfig struct {
    Addr    string
    Handler http.Handler
}
```

**Regras**:
- Tags JSON obrigatórias em structs públicas
- Campos maiúsculos (exportados)
- Comentário documentando struct
- Organizar por relacionamento

---

## 🚀 Otimização de Memória: Pré-alocação de Arrays

### Regra Fundamental: Sempre Prealocalize Arrays Quando o Tamanho É Conhecido

A pré-alocação de arrays melhora a performance eliminando realocations de memória.

### Padrão: Arrays com Tamanho Desconhecido

```go
// ✅ CORRETO - Fazer append com tamanho conhecido previamente
var resultItems []map[string]string
for i := range numberOfIterations {
    // Processar item
    resultItems = append(resultItems, item)
}

// ✅ MELHOR - Prealocalize se souber o tamanho
resultItems := make([]map[string]string, 0, numberOfIterations)
for i := range numberOfIterations {
    // Processar item
    resultItems = append(resultItems, item)
}

// ❌ INCORRETO - Não prealocizar
var resultItems []map[string]string
for i := range 100000 {
    resultItems = append(resultItems, processItem(i))
}
```

### Variações de Uso

#### 1. Quantidade Exata Conhecida

```go
// ✅ CORRETO
appendedTestResults := make([]TestResult, numberOfTests)
for testIndex := range numberOfTests {
    appendedTestResults[testIndex] = runTest(testIndex)
}
```

#### 2. Quantidade Aproximada

```go
// ✅ CORRETO - Use capacity hint
collectedHttpResponses := make([]http.Response, 0, estimatedResponseCount)
for responseItem := range responseChannel {
    collectedHttpResponses = append(collectedHttpResponses, responseItem)
}
```

#### 3. Array de Structs (Teste)

```go
// ✅ CORRETO
testCaseScenarios := []struct {
    testName   string
    httpMethod string
}{
    {"GET request", "GET"},
    {"POST request", "POST"},
    {"DELETE request", "DELETE"},
}

// ✅ CORRETO - Com prealocisão se dinâmico
processedResults := make([]ProcessingResult, 0, len(inputItems))
for _, inputItem := range inputItems {
    processedResults = append(processedResults, processItem(inputItem))
}
```

#### 4. Acumulativo em Loops

```go
// ✅ CORRETO - Quando tamanho total conhecido
collectedResponses := make([]ResponseData, 0, totalExpectedResponses)
for iterationIndex := range numberOfIterations {
    request := httptest.NewRequest(http.MethodGet, "/health", nil)
    responseRecorder := httptest.NewRecorder()
    
    healthHandler(responseRecorder, request)
    
    collectedResponses = append(collectedResponses, parseResponse(responseRecorder))
}

// ✅ MELHOR - Se o tamanho final é exato
finalResponses := make([]ResponseData, numberOfIterations)
for iterationIndex := range numberOfIterations {
    // Processa
    finalResponses[iterationIndex] = processItem(iterationIndex)
}
```

### Quando NÃO Prealocizar

```go
// ✅ OK - Poucos itens, tamanho desconhecido
smallList := []string{}
for item := range smallCollection {
    smallList = append(smallList, item.String())
}

// ✅ OK - Uma única operação
quickResult := append(existingList, newItem)
```

### Impacto de Performance

```go
// Sem pré-alocação: Reallocations múltiplas
// 1 item: 1 alocão
// 10 items: 4-5 reallocations
// 1000 items: 10-12 reallocations
// 100000 items: 17-18 reallocations

// Com pré-alocação: Apenas 1 alocação
// 1 item: 1 alocão
// 10 items: 1 alocão
// 1000 items: 1 alocão
// 100000 items: 1 alocão
```

---

### Teste Básico (Unit Test)

```go
func TestHealthHandlerSuccess(t *testing.T) {
    req := httptest.NewRequest("GET", "/health", nil)
    w := httptest.NewRecorder()

    healthHandler(w, req)

    if w.Code != http.StatusOK {
        t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
    }
}
```

**Padrão de Nomenclatura**: `Test<Component><Scenario>`

**Obrigatórios**:
- `httptest.NewRequest` e `httptest.NewRecorder`
- Comparar com constantes esperadas
- Usar `t.Errorf` com mensagem clara
- Um aspecto por teste (Single Responsibility)

### Table-Driven Tests

```go
func TestSetupRouterWithDifferentMethods(t *testing.T) {
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
            
            router.ServeHTTP(w, req)
            
            if w.Code != http.StatusOK {
                t.Errorf("expected %d, got %d", http.StatusOK, w.Code)
            }
        })
    }
}
```

**Vantagens**:
- Casos múltiplos em um teste
- Cada caso é isolado
- Fácil adicionar novos casos
- Nomes descritivos

### Logging Capture Pattern

```go
func captureLogOutput(fn func()) string {
    var buf bytes.Buffer
    
    originalLogger := log.Writer()
    originalFlags := log.Flags()
    originalPrefix := log.Prefix()
    
    log.SetOutput(&buf)
    log.SetFlags(log.LstdFlags)
    
    fn()
    
    log.SetOutput(originalLogger)
    log.SetFlags(originalFlags)
    log.SetPrefix(originalPrefix)
    
    return buf.String()
}
```

**Padrão Crítico**:
- Salvar estado original
- Restaurar sempre (mesmo após erro)
- Usar `bytes.Buffer`

---

## 🎨 Padrões de Design

### Dependency Injection

```go
type ServerConfig struct {
    Addr    string
    Handler http.Handler
}

func startServer(config ServerConfig) error {
    server := &http.Server{
        Addr:    config.Addr,
        Handler: config.Handler,
    }
    return server.ListenAndServe()
}
```

**Benefícios**:
- Testável
- Reutilizável
- Sem globals
- Inversão de controle

### Factory Pattern

```go
func setupRouter() *http.ServeMux {
    router := http.NewServeMux()
    router.HandleFunc("GET /api/v1/data", handler.GetData)
    router.HandleFunc("GET /health", healthHandler)
    return router
}
```

**Benefícios**:
- Encapsula criação
- Facilita manutenção
- Facilita testes

---

## 💬 Documentação e Comentários

### Comentários de Arquivo

```go
// getdata retorna dados de geração solar da API Solarz
func GetData(w http.ResponseWriter, _ *http.Request) {
    // ...
}
```

**Regras**:
- Começar com o nome da função
- Uma linha se possível
- Explicar "o quê", não "como"
- Sempre em funções públicas

### Comentários Inline (Raro)

```go
// Restaura o logger original para não afetar outros testes
log.SetOutput(originalLogger)
```

**Regras**:
- Explicar "por quê", não "o quê"
- Apenas para lógica complexa
- Evitar ao máximo
- Código deve ser auto-explicativo

---

## 📋 Padrões HTTP

### Status Codes

```go
// ✅ CORRETO - Usar constantes
http.StatusOK                  // 200
http.StatusCreated             // 201
http.StatusBadRequest          // 400
http.StatusNotFound            // 404
http.StatusInternalServerError // 500
```

### Headers

```go
// ✅ PADRÃO
w.Header().Set("Content-Type", "application/json")
```

**Regras**:
- Setar antes de WriteHeader
- Sempre setar Content-Type
- WriteHeader é explícito

### JSON Encoding

```go
// ✅ CORRETO
if err := json.NewEncoder(w).Encode(data); err != nil {
    log.Printf("Failed to encode: %v", err)
}

// ❌ INCORRETO (para responses)
w.Write([]byte(jsonString))
json.NewEncoder(os.Stdout).Encode(data) // stdout, não response
```

---

## 📊 Cobertura de Testes Esperada

| Componente | Mínimo |
|-----------|--------|
| Models | 100% |
| Handlers | 80% |
| Service | 95% |
| Main | 50% |
| **Média** | **~81%** |

**Tipos Obrigatórios de Testes**:
- ✅ Unit tests
- ✅ Integration tests
- ✅ Concurrency tests
- ✅ Error cases
- ✅ Edge cases

---

## 🔍 Formatação

### Indentação

```go
// ✅ CORRETO
- Use tabs (não espaços)
- 1 nível = 1 tab
- Máximo ~120 caracteres por linha

// Exemplo
func complexFunction(
    param1 string,
    param2 int,
) error {
    // ...
}
```

### Espaçamento

```go
// ✅ PADRÃO
- Linha em branco entre funções
- Linha em branco entre grupos lógicos
- Sem linhas em branco excessivas
- Usar gofmt automaticamente
```

### Imports

```go
// ✅ PADRÃO
import (
    "bytes"
    "context"
    "encoding/json"
    "log"
    "net/http"
    
    "github.com/catdevsecops/solarz-api/internal/handler"
)
```

**Ordem**:
1. Stdlib (alfabético)
2. Linha em branco
3. Dependências externas (alfabético)

---

## 🛠️ Checklist Antes de Commit

- [ ] Código segue padrões de nomenclatura
- [ ] Funções têm comentários descritivos
- [ ] Cobertura de testes >= padrão mínimo
- [ ] Todos os casos de erro são testados
- [ ] Sem variáveis globais
- [ ] Error handling está correto (nunca swallow)
- [ ] Logging é apropriado (não poluído)
- [ ] Formatação OK (gofmt)
- [ ] `go test -v` passa
- [ ] `go test -race` sem problemas

---

## 🚀 Comandos Úteis

```bash
# Formatar código
gofmt -w .

# Verificar problemas
go vet ./...

# Rodar testes
go test -v

# Cobertura
go test -cover

# Race conditions
go test -race

# Benchmarks
go test -bench=. -benchmem
```

---

## 📚 Resumo Rápido

### DO's ✅

```go
// ✅ Usar constantes HTTP
w.WriteHeader(http.StatusOK)

// ✅ Usar json.NewEncoder
json.NewEncoder(w).Encode(data)

// ✅ Verificar erros
if err != nil {
    return err
}

// ✅ Documentar público
// GetData retorna dados de geração solar
func GetData(...) { ... }

// ✅ Testes abrangentes
func TestComponentScenario(t *testing.T) { ... }

// ✅ Error handling explícito
if err := doSomething(); err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    // ...
    return
}
```

### DON'Ts ❌

```go
// ❌ Não usar números mágicos
w.WriteHeader(200)

// ❌ Não usar json.Marshal para response
data, _ := json.Marshal(obj)
w.Write(data)

// ❌ Não ignorar erros
someFunction() // falta verificação

// ❌ Não deixar comentários óbvios
// i = i + 1  // incrementa i
i++

// ❌ Não usar variáveis globais
var globalConfig ServerConfig

// ❌ Não swallow errors
if err != nil {
    // não faz nada
}
```

---

## 📞 Contato e Referências

- **Effective Go**: https://golang.org/doc/effective_go
- **Code Review Comments**: https://golang.org/wiki/CodeReviewComments
- **Testing**: https://golang.org/pkg/testing/
- **HTTP**: https://golang.org/pkg/net/http/

---

## 🎯 Última Verificação

Antes de usar Copilot/IA para gerar código, certifique-se de que:

1. ✅ Seguirá padrões de nomenclatura
2. ✅ Usará padrões de código corretos
3. ✅ Incluirá testes apropriados
4. ✅ Terá documentação clara
5. ✅ Não terá variáveis globais
6. ✅ Terá error handling explícito
7. ✅ Será formatado com gofmt

---

**Versão**: 1.1
**Data**: 2024-01-05
**Mantido por**: Equipe de Desenvolvimento
**Status**: Ativo

## 🔄 Últimas Mudanças (v1.1)

### Novas Seções Adicionadas

1. **Padrões de Nomenclatura de Variáveis (Expandido)**
   - Guia completo sobre nomes descritivos
   - Exceções permitidas (parâmetros padrão Go)
   - Exemplos de correto vs incorreto
   - Regras específicas por contexto

2. **Otimização de Memória: Pré-alocação de Arrays**
   - Quando e como prealocizar
   - Variações de uso
   - Impacto de performance
   - Casos de uso

### Motivação

Essas mudanças surgiram de melhorias reais implementadas no código do projeto:

- **Variabilidade de Nomes**: A revisão do `main_test.go` revelou variáveis muito descritivas como `responseRecorder`, `logOutput`, `originalLogger` que melhoram muito a legibilidade
- **Pré-alocação**: Padrão observado em `TestHealthHandlerIsIdempotent` com `var results []map[string]string`

### Conformidade com Projeto

Os padrões agora refletem:
- Estado atual do código base
- Boas práticas observadas
- Performanceáncia de memória
- Legibilidade (Copilot Instructions)


