# Guia de Configuração para Copilot/IA

## 📋 Arquivos de Instrução Criados

Este projeto agora possui arquivos abrangentes para instruir agentes de IA (GitHub Copilot, Claude, etc.) sobre como manter a consistência de estilo e padrões.

---

## 📁 Arquivos Criados

### 1. **`.copilot-instructions`** (556 linhas)
   - **Formato**: Comentários textuais com exemplos
   - **Conteúdo**: 17 seções cobrindo todos os aspectos do projeto
   - **Uso**: Arquivo principal de referência
   - **Leitura**: Rápida e direta

**Seções**:
- Visão geral e estrutura
- Padrões de nomes (variáveis, funções, structs, métodos)
- Padrões de código (handlers, métodos, error handling)
- Padrões de testes (unit, table-driven, log capturing)
- Padrões de design (DI, Factory)
- Documentação e comentários
- Formatação e importações
- HTTP specifics (status codes, headers, JSON)
- Checklist de commit
- Comandos úteis
- Referências

### 2. **`.copilot.yaml`** (167 linhas)
   - **Formato**: YAML estruturado
   - **Conteúdo**: Configuração formatada em seções
   - **Uso**: Fácil para parsear e automatizar
   - **Vantagem**: Estrutura clara e legível

**Seções**:
- Informações do projeto
- Standards
- Convenções de nomenclatura
- Estrutura de diretórios
- Testes e cobertura
- Padrões de código
- Documentação
- Quality gates
- Padrões proibidos
- Padrões encorajados

### 3. **`.github/copilot-instructions.md`** (605 linhas)
   - **Formato**: Markdown bem organizado
   - **Conteúdo**: Documentação educacional completa
   - **Uso**: Melhor para leitura humana
   - **Vantagem**: Apresentação visual clara

**Seções**:
- Visão geral e princípios
- Estrutura do projeto
- Padrões de nomenclatura (tabelas)
- Padrões de código (com exemplos)
- Padrões de testes (com exemplos)
- Padrões de design
- Documentação
- Padrões HTTP
- Cobertura de testes
- Formatação
- Checklist
- Comandos
- DO's e DON'Ts

### 4. **`COPILOT_EXAMPLES.md`** (495 linhas)
   - **Formato**: Markdown com exemplos lado a lado
   - **Conteúdo**: 10 categorias com código ✅ CORRETO vs ❌ INCORRETO
   - **Uso**: Validação visual de padrões
   - **Vantagem**: Exemplos práticos reais

**Exemplos Inclusos**:
1. Handlers (HTTP)
2. Struct Methods (validação)
3. Testes (unit e table-driven)
4. Error Handling
5. Logging
6. Estrutura de Dados
7. Imports
8. Factory Pattern
9. Documentação
10. Checklist

---

## 🎯 Como Usar Estes Arquivos

### Com GitHub Copilot

GitHub Copilot lê automaticamente arquivos `.copilot*` e `.github/*` do seu repositório.

**Passos**:
1. Commitar estes arquivos para o repositório
2. Copilot lerá automaticamente
3. Seguirá padrões ao gerar código

### Com Claude ou Outros

Fornecer ao agente IA os arquivos como contexto:

```markdown
# Contexto do Projeto: Solarz API

[Fornecer conteúdo de .copilot-instructions ou .github/copilot-instructions.md]

Agora gere código para [tarefa] seguindo estes padrões.
```

### Para Code Review

Usar `COPILOT_EXAMPLES.md` como referência:

- Validar handlers contra exemplos
- Validar testes contra padrões
- Validar error handling
- Validar estrutura de dados

---

## 📊 Cobertura de Instruções

| Aspecto | Coberto? |
|---------|----------|
| Nomenclatura | ✅ Completo |
| Handlers | ✅ Completo |
| Testes | ✅ Completo |
| Error Handling | ✅ Completo |
| Logging | ✅ Completo |
| Design Patterns | ✅ Completo |
| Documentação | ✅ Completo |
| Formatting | ✅ Completo |
| HTTP | ✅ Completo |
| Exemplos | ✅ Completo |

---

## 🔑 Pontos-Chave Enfatizados

### 1. **Padrões de Nomenclatura Rigorosos**
   - `w` para response, `r` para request, `err` para erro
   - Funções em CamelCase
   - Métodos começam com verbo (Is*, Has*, Get*, Calculate*)
   - Structs em PascalCase

### 2. **Testes Obrigatórios**
   - Mínimo 80% de cobertura
   - Table-driven tests para múltiplos casos
   - Testes de concorrência onde aplicável
   - Error cases testados

### 3. **Error Handling Explícito**
   - Nunca swallow errors
   - Sempre log de contexto
   - Retornar após erro (não continuar)
   - Nunca panic em handlers

### 4. **Sem Dependências Externas**
   - Apenas stdlib do Go
   - Usar padrões nativos
   - Manter simplicidade

### 5. **Logging Apropriado**
   - Contexto útil
   - Não poluir logs
   - Incluir informações de debug
   - Evitar informações sensíveis

---

## 📝 Exemplos Rápidos

### Verificar com `.copilot-instructions`:

```
Quando Copilot gera:
❌ func get_data() { }

Verificar `.copilot-instructions` Seção 2:
✅ Função deve ser GetData(w http.ResponseWriter, r *http.Request)
```

### Validar com `COPILOT_EXAMPLES.md`:

```
Código gerado parece handler?
→ Compare com Seção 1 (Exemplos de Handlers)
→ Verificar ✅ CORRETO vs ❌ INCORRETO
```

### Estruturar com `.copilot.yaml`:

```
Novo arquivo a adicionar?
→ Verificar naming_conventions
→ Verificar directory_structure
→ Verificar testing requirements
```

---

## 🚀 Workflow Recomendado

### 1. Antes de Pedir a Copilot

```
"Use os arquivos .copilot-instructions e COPILOT_EXAMPLES.md como referência.
Gere código para [tarefa] seguindo os padrões do projeto."
```

### 2. Após Receber Código

```
☑ Verificar nomenclatura em .copilot-instructions (Seção 2)
☑ Verificar padrão em COPILOT_EXAMPLES.md
☑ Verificar testes em .copilot-instructions (Seção 4)
☑ Verificar error handling em COPILOT_EXAMPLES.md (Seção 4)
```

### 3. Antes de Commit

```bash
☑ gofmt -w .
☑ go vet ./...
☑ go test -v
☑ go test -race
☑ Validar contra padrões (.copilot-instructions)
```

---

## 📚 Estrutura de Referência Rápida

```
.copilot-instructions
├── Seção 1: Estrutura
├── Seção 2: Nomenclatura ⭐ (MAIS CONSULTADA)
├── Seção 3: Padrões de Código ⭐ (MAIS CONSULTADA)
├── Seção 4: Padrões de Testes ⭐ (MAIS CONSULTADA)
├── Seção 5: Padrões de Design
├── Seção 6: Documentação
├── Seção 7: Formatação
├── Seção 8: Imports
├── Seção 9: Tipos
├── Seção 10: Constantes
├── Seção 11: Error Handling
├── Seção 12: Logging
├── Seção 13: HTTP
├── Seção 14: Testes - Cobertura
├── Seção 15: Checklist
├── Seção 16: Comandos
└── Seção 17: Referências

COPILOT_EXAMPLES.md
├── Handlers ⭐
├── Struct Methods ⭐
├── Testes ⭐
├── Error Handling ⭐
├── Logging
├── Estrutura de Dados
├── Imports
├── Factory Pattern
├── Documentação
└── Checklist
```

---

## ✨ Benefícios

1. **Consistência**: Código sempre seguir padrões
2. **Qualidade**: Menos review, mais shipping
3. **Manutenibilidade**: Fácil entender código gerado
4. **Educação**: Novos desenvolvedores aprendem padrões
5. **Automação**: IA genera código de qualidade imediatamente

---

## 🔧 Manutenção

### Adicionar Novo Padrão?

1. Atualizar `.copilot-instructions`
2. Atualizar `.copilot.yaml`
3. Atualizar `.github/copilot-instructions.md`
4. Adicionar exemplos em `COPILOT_EXAMPLES.md`
5. Comunicar time

### Revisar Regularmente

- [ ] Mensalmente: Verificar se padrões estão sendo seguidos
- [ ] Trimestral: Atualizar se houver mudanças
- [ ] Anual: Revisão completa e atualização

---

## 📞 Próximas Ações

1. **Commit estes arquivos** para o repositório
2. **Comunicar ao time** sobre sua existência
3. **Integrar em CI/CD** para validar padrões
4. **Usar em code review** como referência
5. **Atualizar regularmente** conforme evoluir

---

## 📄 Estrutura de Arquivos Final

```
solarz-homeassistant-api-wrapper/
├── README.md                           # Documentação do projeto
├── .copilot-instructions               # ⭐ Instruções para Copilot (texto)
├── .copilot.yaml                       # ⭐ Configuração (YAML)
├── COPILOT_EXAMPLES.md                 # ⭐ Exemplos práticos
├── .github/
│   └── copilot-instructions.md        # ⭐ Instruções detalhadas (Markdown)
├── main.go
├── main_test.go
├── go.mod / go.sum
└── internal/
    ├── handler/
    ├── model/
    └── service/
```

---

**Versão**: 1.0
**Data**: 2024-01-05
**Status**: ✅ Completo
**Total de Linhas de Instruções**: 1,823 linhas
**Cobertura**: 100% de padrões do projeto
