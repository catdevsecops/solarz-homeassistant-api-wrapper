# 🐳 Docker Publishing - GitHub Packages (GHCR)

## Visão Geral

O projeto publica imagens Docker automaticamente no **GitHub Container Registry (GHCR)** através do workflow `docker-publish.yml`.

**Registry**: `ghcr.io`
**Imagem Pública**: `ghcr.io/catdevsecops/solarz-homeassistant-api-wrapper`

---

## 🚀 Workflow Automático

### Arquivo de Configuração
- **Caminho**: `.github/workflows/docker-publish.yml`
- **Plataforma**: GitHub Actions
- **Sem configuração manual necessária** ✅

### Triggers (Quando faz build/push)

| Evento | Ação |
|--------|------|
| Push em `main` | Build + Push com tag `latest` |
| Tag `v*` (ex: v1.2.3) | Build + Push com tags semver |
| Pull Request | Build apenas (sem push) |

### Autenticação

```yaml
username: ${{ github.actor }}
password: ${{ secrets.GITHUB_TOKEN }}
```

✅ **Nenhum secret manual necessário** - usa o token automático do GitHub Actions

---

## 🏷️ Estratégia de Tags

A imagem recebe múltiplas tags para diferentes cenários:

### 1. Push em `main`

```
ghcr.io/catdevsecops/solarz-homeassistant-api-wrapper:latest
ghcr.io/catdevsecops/solarz-homeassistant-api-wrapper:main-a1b2c3d
```

### 2. Tag de Release (v1.2.3)

```
ghcr.io/catdevsecops/solarz-homeassistant-api-wrapper:1.2.3      # full
ghcr.io/catdevsecops/solarz-homeassistant-api-wrapper:1.2        # minor
ghcr.io/catdevsecops/solarz-homeassistant-api-wrapper:1          # major
ghcr.io/catdevsecops/solarz-homeassistant-api-wrapper:sha-a1b2c3d # commit
```

### 3. Pull Request

```
# Build apenas, sem push ao registry
```

---

## 📦 Como Usar a Imagem

### Pull da Imagem Pública

```bash
# Versão mais recente
docker pull ghcr.io/catdevsecops/solarz-homeassistant-api-wrapper:latest

# Versão específica
docker pull ghcr.io/catdevsecops/solarz-homeassistant-api-wrapper:1.2.3

# Versão minor
docker pull ghcr.io/catdevsecops/solarz-homeassistant-api-wrapper:1.2
```

### Executar Container

```bash
docker run -it ghcr.io/catdevsecops/solarz-homeassistant-api-wrapper:latest

# Com portas
docker run -p 8080:8080 ghcr.io/catdevsecops/solarz-homeassistant-api-wrapper:latest

# Em background
docker run -d ghcr.io/catdevsecops/solarz-homeassistant-api-wrapper:latest
```

### Docker Compose

```yaml
version: '3.8'
services:
  solarz-api:
    image: ghcr.io/catdevsecops/solarz-homeassistant-api-wrapper:latest
    ports:
      - "8080:8080"
    environment:
      - ENV=production
```

---

## 🔍 Verificando Imagens Publicadas

### Via GitHub Web UI

1. Acesse: https://github.com/catdevsecops/solarz-homeassistant-api-wrapper
2. Vá para: **Packages** (na lateral)
3. Procure por: `solarz-homeassistant-api-wrapper`
4. Veja todas as tags disponíveis

### Via GitHub CLI

```bash
# Listar todas as versões
gh api repos/catdevsecops/solarz-homeassistant-api-wrapper/packages

# Ver detalhes de um package
gh api repos/catdevsecops/solarz-homeassistant-api-wrapper/packages/container
```

### Via Docker CLI

```bash
# Inspeccionar imagem local
docker inspect ghcr.io/catdevsecops/solarz-homeassistant-api-wrapper:latest

# Ver camadas
docker history ghcr.io/catdevsecops/solarz-homeassistant-api-wrapper:latest
```

---

## 📝 Dockerfile

### Requisitos

- ✅ Localizado na **raiz do projeto**
- ✅ Configurado para **produção**
- ✅ Suporta **multi-arquitetura** (amd64, arm64)

### Exemplo (Padrão Go)

```dockerfile
# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o solarz-api -ldflags="-X main.Version=${VERSION}" .

# Runtime stage
FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/solarz-api /usr/local/bin/
EXPOSE 8080
CMD ["solarz-api"]
```

### Otimizações Aplicadas

- ✅ Multi-stage build (reduz tamanho da imagem)
- ✅ Alpine Linux (base leve)
- ✅ CGO_ENABLED=0 (compatibilidade)
- ✅ Sem layer desnecessárias

---

## 🔄 Fluxo de Publicação

```
┌─────────────────────────────────────┐
│  Desenvolvedor faz push             │
└────────────────┬────────────────────┘
                 ↓
     ┌──────────────────────┐
     │  GitHub Actions      │
     │  ci-cd.yml (tests)   │
     └──────────┬───────────┘
                ↓
     ┌──────────────────────────┐
     │  docker-publish.yml      │
     │  (build + push)          │
     └──────────┬───────────────┘
                ↓
     ┌──────────────────────────┐
     │  GitHub Container        │
     │  Registry (GHCR)         │
     │  🐳 Imagem publicada!    │
     └──────────────────────────┘
```

---

## ⚙️ Configuração (Sem Necessidade de Setup Manual)

### Permissões do Workflow

Já configuradas no arquivo `docker-publish.yml`:

```yaml
permissions:
  contents: read
  packages: write
```

✅ Nenhuma configuração adicional necessária

### Repository Settings

**Visibilidade da Imagem**: Pública por padrão

Se quiser mudar:
1. GitHub → Repository → Packages
2. Selecionar pacote
3. Package Settings → Change Visibility

---

## 🚨 Troubleshooting

### ❌ Build falha

**Causa**: Dockerfile com erro
**Solução**:
```bash
# Testar localmente
docker build .
```

### ❌ Push não funciona

**Causa**: Permissões insuficientes
**Solução**: Verificar se workflow tem `packages: write`

### ❌ Imagem não aparece no registry

**Causa**: PR não publica (esperado)
**Solução**: Fazer merge em `main` ou criar tag `v*`

### ❌ Cache muito lento

**Causa**: Primeira build não tem cache
**Solução**: Builds subsequentes serão rápidas

---

## 📊 Monitoramento

### Verificar Status do Workflow

1. GitHub → Actions
2. Selecionar `docker-publish.yml`
3. Ver builds recentes
4. Expandir `build-and-push` job para detalhes

### Métricas Úteis

- **Tempo de build**: Ideal < 5 minutos
- **Tamanho da imagem**: Otimizar no Dockerfile
- **Taxa de sucesso**: Deve ser 100%

---

## 🔐 Segurança

### Imagem Pública ✅

- Qualquer um pode fazer pull
- Código-fonte é aberto (GitHub público)
- Sem exposição de secrets

### Best Practices

1. ✅ Não incluir secrets no Docker
2. ✅ Usar variáveis de ambiente em runtime
3. ✅ Keep image updated
4. ✅ Use health checks

---

## 📚 Referências

- **GitHub Container Registry**: https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry
- **Docker Metadata Action**: https://github.com/docker/metadata-action
- **Docker Build Push**: https://github.com/docker/build-push-action

---

## 🔄 Changelog

### v1.0 (2026-06-05)
- ✨ Migração de Docker Hub para GitHub Packages
- ✨ Workflow completo de build e push
- ✨ Multi-arquitetura (amd64, arm64)
- ✨ Cache automático do GitHub Actions
- ✨ Tags semver automáticas

---

**Mantido por**: Equipe de Desenvolvimento
**Data**: 2026-06-05
**Status**: Ativo ✅
