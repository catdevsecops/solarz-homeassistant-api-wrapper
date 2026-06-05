# GitHub Actions - Pipeline CI/CD

## 📋 Visão Geral

Este projeto possui dois pipelines completos de CI/CD executados automaticamente via GitHub Actions.

---

## 🚀 Pipelines Implementados

### 1. **CI/CD Pipeline** (`.github/workflows/ci-cd.yml`)

Executado em **push** e **pull_request** para `main` e `develop`.

#### Jobs:

##### **Test & Coverage**
- ✅ Checkout do código
- ✅ Setup Go 1.21
- ✅ Download e verificação de dependências
- ✅ Verificação de formatação (gofmt)
- ✅ Lint com `go vet`
- ✅ Rodar testes com race detector
- ✅ Gerar relatório de cobertura
- ✅ Upload para Codecov
- ✅ Verificar threshold mínimo (50%)

##### **Security Scan**
- ✅ Executar `gosec` para verificar vulnerabilidades
- ✅ Upload de SARIF para GitHub Security
- ✅ Não falha o pipeline (apenas relatório)

##### **Build Binary**
- ✅ Build binário Linux x86_64
- ✅ Verificar integridade
- ✅ Upload como artifact (7 dias)

##### **Code Quality**
- ✅ `golangci-lint` com configuração customizada
- ✅ 30+ verificações de qualidade
- ✅ Timeout de 5 minutos

##### **Benchmark Tests**
- ✅ Rodar benchmarks apenas em `main`
- ✅ Comparar com benchmarks anteriores
- ✅ Auto-push de resultados

##### **Notify on Failure**
- ✅ Notificação se algum job falhar
- ✅ Resumo de status

---

### 2. **Build & Release Pipeline** (`.github/workflows/build-release.yml`)

Executado em **push** para `main` e em **tags** `v*`.

#### Jobs:

##### **Build Multi-Platform**
- ✅ Build para Linux (amd64, arm64)
- ✅ Build para macOS (amd64, arm64)
- ✅ Build para Windows (amd64)
- ✅ Incluir versão no binário
- ✅ Upload como artifacts

##### **Create Release**
- ✅ Apenas para tags `v*`
- ✅ Preparar assets de release
- ✅ Gerar release notes automáticas
- ✅ Upload binários para GitHub Releases

##### **Build Docker Image**
- ✅ Build imagem Docker multi-stage
- ✅ Login no DockerHub (se configurado)
- ✅ Push para registry
- ✅ Cache de build layers

---

## 📊 Configurações

### Go
- **Versão**: 1.21
- **Cache**: Habilitado (mais rápido)

### Lint (golangci-lint)
- **Timeout**: 5 minutos
- **Linters**: 40+ habilitados
- **Complexidade Ciclomática**: Máximo 15
- **Comprimento de Função**: Máximo 100 linhas
- **Comprimento de Linha**: Máximo 120 caracteres

### Testes
- **Race Detector**: Habilitado
- **Timeout**: 5 minutos
- **Cobertura Mínima**: 50%
- **Codecov**: Upload automático

### Security
- **gosec**: Ativo em todos os pushes
- **SARIF Upload**: Para GitHub Security

---

## 🔧 Arquivos de Configuração

### `.golangci.yml`
Configuração detalhada para golangci-lint com:
- 40+ linters habilitados
- Exclusões para testes
- Limites de complexidade
- Variáveis ignoradas (w, r, i, etc.)

### `Dockerfile`
Multi-stage Dockerfile com:
- Stage 1: Builder (Go alpine)
- Stage 2: Runtime (Alpine mínimo)
- Usuário não-root
- Health check
- Otimizações de tamanho

---

## 📌 Triggers

### CI/CD Pipeline
```yaml
push:
  branches: [ main, develop ]
pull_request:
  branches: [ main, develop ]
```

### Build & Release Pipeline
```yaml
push:
  branches: [ main ]
  tags: [ 'v*' ]
```

### Benchmark (apenas main)
```yaml
if: github.event_name == 'push' && github.ref == 'refs/heads/main'
```

### Release (apenas tags)
```yaml
if: startsWith(github.ref, 'refs/tags/v')
```

---

## 🔐 Secrets Necessários

### Para Docker Push (Opcional)

Adicionar em **Settings → Secrets and variables → Actions**:

```
DOCKER_USERNAME  = seu_usuario_dockerhub
DOCKER_PASSWORD  = seu_token_dockerhub
```

Se não configurado, Docker build é skipped.

### GitHub Token

Automático (GITHUB_TOKEN), usado para:
- Release notes
- Benchmark push
- Code scanning

---

## 📈 Outputs e Artifacts

### Coverage
- Codecov dashboard (se integrado)
- Relatório em logs do GitHub

### Artifacts
- **solarz-api-binary** (7 dias): Build Linux
- **solarz-api-linux-amd64** (7 dias): Release build
- **solarz-api-linux-arm64** (7 dias): Release build
- **solarz-api-darwin-amd64** (7 dias): Release build
- **solarz-api-darwin-arm64** (7 dias): Release build
- **solarz-api-windows-amd64.exe** (7 dias): Release build

### Releases
- GitHub Releases com binários (em tags)
- Release notes automáticas

---

## 🚀 Como Usar

### Dispara automaticamente

1. **Push para main/develop**: CI/CD roda
2. **Pull Request**: CI/CD roda (blocking)
3. **Push tag v***: Build & Release roda

### Monitorar

1. Ir para **Actions** no GitHub
2. Ver status do pipeline
3. Clicar em job para detalhes
4. Ver logs de cada step

### Fazer Release

```bash
# Criar tag
git tag -a v1.0.0 -m "Release version 1.0.0"

# Push tag
git push origin v1.0.0
```

Isso disparará o pipeline de release automaticamente.

---

## ✅ Checklist de Teste

Cada push verifica:

- [ ] Formatação (gofmt)
- [ ] Lint (go vet + golangci-lint)
- [ ] Testes passam
- [ ] Race conditions (race detector)
- [ ] Cobertura >= 50%
- [ ] Segurança (gosec)
- [ ] Build sucede
- [ ] Benchmarks (main apenas)

**Status**: ✅ Todos devem passar para merge

---

## 🔍 Linters Habilitados

### Categorias

**Segurança**:
- gosec, g302, g304, g305, g306, g307, g308

**Estilo**:
- gofmt, goimports, stylecheck, revive, godot

**Erros**:
- errcheck, errchkjson, errorlint, nilerr, nilnil

**Performance**:
- ineffassign, prealloc, unconvert, wastedassign

**Padrões**:
- gocritic, govet, staticcheck, unused, unusedresult

**Tipos**:
- asasalint, typecheck, exhaustive, ireturn, musttag

**Complexidade**:
- cyclomatic (max 15), gocognit (max 20), funlen (max 100 linhas)

**Outras**:
- 30+ linters adicionais

---

## 📊 Exemplo de Saída

```
✅ Formatting: PASSED
✅ Linting: PASSED (0 issues)
✅ Testing: PASSED (62/62 tests)
✅ Coverage: 81.5%
✅ Race Detection: PASSED (no races)
✅ Security: PASSED (0 vulnerabilities)
✅ Build: PASSED (solarz-api binary)
✅ Benchmarks: Updated
```

---

## 🎯 Próximas Melhorias

### Recomendado
- [ ] Adicionar SONARQUBE para análise
- [ ] Adicionar renovate para dependências
- [ ] Adicionar sentry para tracking de erros
- [ ] Deploy automático em staging
- [ ] Notificações Slack/Discord

### Nice to Have
- [ ] SBOM (Software Bill of Materials)
- [ ] Attestation de assinatura
- [ ] Load tests
- [ ] E2E tests
- [ ] Database migration tests

---

## 📚 Referências

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Go Testing](https://golang.org/doc/effective_go#testing)
- [golangci-lint](https://golangci-lint.run/)
- [Codecov](https://codecov.io/)
- [Docker Multi-stage](https://docs.docker.com/build/building/multi-stage/)

---

**Status**: ✅ Pronto para Produção
**Data**: 2024-01-05
**Versão**: 1.0
