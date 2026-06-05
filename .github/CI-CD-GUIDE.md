# 📋 Guia de CI/CD - GitHub Actions

Este documento descreve o pipeline CI/CD completo do projeto Solarz API usando **GitHub Actions**.

## 📚 Índice

1. [Visão Geral](#visão-geral)
2. [Workflows Disponíveis](#workflows-disponíveis)
3. [Docker Publishing](#docker-publishing)
4. [Fluxo de Desenvolvimento](#fluxo-de-desenvolvimento)
5. [Troubleshooting](#troubleshooting)

---

## Visão Geral

O projeto utiliza **GitHub Actions** para automação completa de CI/CD:

```
┌──────────────┐
│ GitHub Repo  │
└──────┬───────┘
       │
       ├─→ [Test & Lint] → Test, Security, Format Check
       │
       ├─→ [Docker Build] → Build & Push to GHCR
       │
       └─→ [Release] → Build binários multi-plataforma
```

### Características

- ✅ **Zero-config**: Sem tokens ou secrets adicionais
- ✅ **Automático**: Dispara em push, PR, tags
- ✅ **Multi-arquitetura**: Linux, macOS, Windows
- ✅ **Docker**: Publicado no GitHub Container Registry
- ✅ **Seguro**: Sem exposição de credenciais

---

## Workflows Disponíveis

### 1️⃣ CI/CD Pipeline (`ci-cd.yml`)

**Quando executa**: Push em `main`/`develop` e Pull Requests

**Jobs**:

#### Test & Coverage
```
go test -v -race -timeout=5m ./...
go test -coverprofile=coverage.out ./...
```
- ✅ Testes unitários
- ✅ Testes com race detector
- ✅ Coverage >= 50% (obrigatório)
- ✅ Upload para Codecov

#### Code Quality
```
golangci-lint run --timeout=5m
```
- ✅ Go vet
- ✅ golangci-lint (latest)
- ✅ Formato com gofmt

#### Security Scan
```
gosec ./...
```
- ✅ Vulnerabilidades de segurança
- ✅ Upload SARIF para GitHub Security

#### Build Binary
```
go build -o build/solarz-api .
```
- ✅ Compilation check
- ✅ Linux amd64 (CGO_ENABLED=0)
- ✅ Upload artifact (7 dias)

#### Benchmarks (main only)
```
go test -bench=. -benchmem ./...
```
- ✅ Performance tests
- ✅ Apenas em push para `main`

---

### 2️⃣ Docker Publishing (`docker-publish.yml`)

**Quando executa**: 
- Push em `main` → `latest` + `sha`
- Tag `v*` → versão semver
- Pull Request → build only (sem push)

**Actions**:
- Docker Buildx (multi-arch)
- GitHub Container Registry
- Automatic tagging

**Detalhes**: Ver [DOCKER_PUBLISHING.md](.github/DOCKER_PUBLISHING.md)

---

### 3️⃣ Build & Release (`build-release.yml`)

**Quando executa**: Tag `v*` (ex: v1.2.3)

**Plataformas**:
```
Linux:   amd64, arm64
macOS:   amd64, arm64
Windows: amd64
```

**Saída**:
- Binários compilados
- GitHub Release automático
- Release Notes geradas

---

## Docker Publishing

### 📍 Registry

- **URL**: `ghcr.io`
- **Imagem**: `ghcr.io/catdevsecops/solarz-homeassistant-api-wrapper`
- **Visibilidade**: Pública ✅

### 🏷️ Tags

```
# Main branch
ghcr.io/catdevsecops/solarz-homeassistant-api-wrapper:latest

# Version tags (v1.2.3)
ghcr.io/catdevsecops/solarz-homeassistant-api-wrapper:1.2.3
ghcr.io/catdevsecops/solarz-homeassistant-api-wrapper:1.2
ghcr.io/catdevsecops/solarz-homeassistant-api-wrapper:1

# Commit sha
ghcr.io/catdevsecops/solarz-homeassistant-api-wrapper:sha-a1b2c3d
```

### 🐳 Como Usar

```bash
# Pull
docker pull ghcr.io/catdevsecops/solarz-homeassistant-api-wrapper:latest

# Run
docker run -p 8080:8080 \
  ghcr.io/catdevsecops/solarz-homeassistant-api-wrapper:latest

# Compose
image: ghcr.io/catdevsecops/solarz-homeassistant-api-wrapper:latest
```

---

## Fluxo de Desenvolvimento

### 1. Feature Branch

```bash
git checkout -b feature/my-feature
# ... desenvolvimento ...
git push origin feature/my-feature
```

**Workflow executado**: `ci-cd.yml` (testes + linting)

### 2. Pull Request

```
GH Web → Create Pull Request
```

**Workflows executados**:
- `ci-cd.yml` (testes obrigatórios ✅)
- `docker-publish.yml` (build docker, sem push)

**Requisitos para merge**:
- ✅ Testes passando
- ✅ Coverage >= 50%
- ✅ Sem lint errors

### 3. Merge para Main

```bash
gh pr merge --squash
# ou via GitHub Web
```

**Workflows executados**:
- `ci-cd.yml` (final check)
- `docker-publish.yml` (push imagem latest + sha)

### 4. Release

```bash
git tag v1.2.3
git push origin v1.2.3
```

**Workflows executados**:
- `docker-publish.yml` (push com tags semver)
- `build-release.yml` (binários + GitHub Release)

---

## Monitorando Workflows

### GitHub Web UI

1. Repository → **Actions**
2. Selecione workflow
3. Veja runs recentes
4. Expanda job para logs

### GitHub CLI

```bash
# Listar runs recentes
gh run list --workflow=ci-cd.yml

# Ver detalhes de um run
gh run view <run-id> --log

# Cancelar run
gh run cancel <run-id>
```

### Badges (README.md)

```markdown
[![CI/CD Pipeline](https://github.com/catdevsecops/solarz-homeassistant-api-wrapper/actions/workflows/ci-cd.yml/badge.svg)](https://github.com/catdevsecops/solarz-homeassistant-api-wrapper/actions)

[![Docker Publish](https://github.com/catdevsecops/solarz-homeassistant-api-wrapper/actions/workflows/docker-publish.yml/badge.svg)](https://github.com/catdevsecops/solarz-homeassistant-api-wrapper/packages)
```

---

## Troubleshooting

### 🔴 Tests Failing

**Logs**:
- Actions → ci-cd.yml → test job → output

**Causas comuns**:
```go
// ❌ Coverage abaixo de 50%
// ✅ Solução: Adicionar testes

// ❌ Race condition detectada
// ✅ Solução: Sincronizar acesso a variáveis

// ❌ Lint error
// ✅ Solução: go fmt -w . && golangci-lint run
```

### 🔴 Docker Build Failing

**Causas**:
```dockerfile
# ❌ Dockerfile em local errado
# ✅ Deve estar na raiz

# ❌ Dependência de build faltando
# ✅ Adicionar RUN apk add

# ❌ Port não exposto
# ✅ Adicionar EXPOSE 8080
```

### 🔴 Push não ocorrendo

**Causa**: PR não publica (esperado)
**Solução**: Fazer merge em `main` ou criar tag

**Verificar**:
```bash
# Se é PR
if [[ $GITHUB_EVENT_NAME == "pull_request" ]]; then
  echo "PRs não fazem push"
fi
```

### 🔴 Imagem não encontrada no registry

**Possibilidades**:
1. Build ainda rodando → aguardar
2. Job falhou → ver logs
3. Não é main/tag → fazer commit correto

**Testar**:
```bash
docker pull ghcr.io/catdevsecops/solarz-homeassistant-api-wrapper:latest
```

---

## Segurança

### Secrets & Tokens

- ✅ Workflow usa `GITHUB_TOKEN` (automático)
- ✅ Nenhum secret manual necessário
- ✅ Tokens expiram automaticamente

### Permissões

```yaml
permissions:
  contents: read
  packages: write  # Necessário para publicar Docker
```

### Vulnerabilidades

- ✅ `gosec` verifica segurança
- ✅ Dependências auditadas com `go mod tidy`
- ✅ Container baseado em Alpine (mínimo)

---

## Best Practices

### Commits

```bash
# ✅ Commit específico
git commit -m "feat: add health check endpoint"

# ✅ Commit com corpo
git commit -m "fix: race condition in logger
  
  - Adicionar mutex
  - Sincronizar acesso
"

# ❌ Commit vago
git commit -m "fix bug"
```

### Pull Requests

```bash
# ✅ PR descritivo
gh pr create --title "Add Docker healthcheck" \
  --body "Adiciona healthcheck para melhor monitoramento"

# ✅ Esperar CI/CD passar
# ✅ Solicitar review
```

### Tags

```bash
# ✅ Semver correto
git tag v1.2.3

# ✅ Com mensagem
git tag -a v1.2.3 -m "Release v1.2.3 - Add feature X"

# ❌ Tag sem versão
git tag release
```

---

## Referências

- **GitHub Actions**: https://docs.github.com/en/actions
- **Container Registry**: https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry
- **Docker Buildx**: https://github.com/docker/buildx
- **Go Testing**: https://golang.org/pkg/testing/
- **golangci-lint**: https://golangci-lint.run/

---

## 🔄 Changelog

### v1.0 (2026-06-05)
- ✨ Pipeline completo com GitHub Actions
- ✨ Docker Publishing no GHCR
- ✨ Multi-plataforma (Go binaries)
- ✨ Testes + Security + Linting
- ✨ Zero-config CI/CD

---

**Status**: ✅ Ativo
**Mantido por**: Equipe de Desenvolvimento
**Data**: 2026-06-05
