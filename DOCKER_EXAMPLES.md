# 🎯 Exemplos Práticos: Docker Publishing

Este arquivo contém exemplos reais de como usar o novo pipeline de Docker Publishing.

---

## 📝 Exemplo 1: Fluxo Básico de Desenvolvimento

### Passo 1: Criar feature branch

```bash
git checkout -b feature/new-endpoint
```

### Passo 2: Fazer mudanças no código

```go
// internal/handler/items.go
func GetAllItems(w http.ResponseWriter, r *http.Request) {
    // implementação
}
```

### Passo 3: Fazer commit e push

```bash
git add internal/handler/items.go
git commit -m "feat: add new endpoint for items"
git push origin feature/new-endpoint
```

**Automação ativada**:
- ✅ `ci-cd.yml` → Testes + linting
- ✅ `docker-publish.yml` → Build (sem push, é PR)

### Passo 4: Criar Pull Request

```bash
gh pr create --title "Add new items endpoint" \
  --body "Implementa novo endpoint para listar items"
```

**Verificação no GitHub**:
```
Checks running...
  ✓ ci-cd.yml (tests passing)
  ✓ docker-publish.yml (docker build successful)
```

### Passo 5: Merge

```bash
gh pr merge --squash
```

**Automação ativada**:
- ✅ `ci-cd.yml` → Final test check
- ✅ `docker-publish.yml` → **PUSH para GHCR**

**Imagem publicada com tags**:
```
ghcr.io/catdevsecops/solarz-homeassistant-api-wrapper:latest
ghcr.io/catdevsecops/solarz-homeassistant-api-wrapper:main-a1b2c3d4
```

---

## 📝 Exemplo 2: Release com Versionamento

### Passo 1: Atualizar versão no código

```go
// main.go
const Version = "1.2.3"
```

### Passo 2: Criar tag de release

```bash
git tag -a v1.2.3 -m "Release v1.2.3 - Add new endpoint"
git push origin v1.2.3
```

**Automação ativada**:
- ✅ `build-release.yml` → Compila binários multi-plataforma
- ✅ `docker-publish.yml` → **Build + PUSH com tags semver**

**Imagens publicadas**:
```
ghcr.io/catdevsecops/solarz-homeassistant-api-wrapper:1.2.3  ← full version
ghcr.io/catdevsecops/solarz-homeassistant-api-wrapper:1.2    ← minor
ghcr.io/catdevsecops/solarz-homeassistant-api-wrapper:1      ← major
ghcr.io/catdevsecops/solarz-homeassistant-api-wrapper:sha-a1b2c3d4
```

**Release no GitHub**:
- Binários: `solarz-api-linux-amd64`, `solarz-api-darwin-arm64`, etc
- Release Notes automáticas

---

## 🐳 Exemplo 3: Usando Imagem Publicada

### Puxar a imagem

```bash
# Versão mais recente
docker pull ghcr.io/catdevsecops/solarz-homeassistant-api-wrapper:latest

# Versão específica
docker pull ghcr.io/catdevsecops/solarz-homeassistant-api-wrapper:1.2.3

# Versão minor
docker pull ghcr.io/catdevsecops/solarz-homeassistant-api-wrapper:1.2
```

### Executar a imagem

```bash
# Básico
docker run \
  ghcr.io/catdevsecops/solarz-homeassistant-api-wrapper:latest

# Com porta exposta
docker run -p 8080:8080 \
  ghcr.io/catdevsecops/solarz-homeassistant-api-wrapper:latest

# Com variáveis de ambiente
docker run -p 8080:8080 \
  -e ENV=production \
  -e LOG_LEVEL=info \
  ghcr.io/catdevsecops/solarz-homeassistant-api-wrapper:latest

# Em background
docker run -d -p 8080:8080 \
  --name solarz-api \
  ghcr.io/catdevsecops/solarz-homeassistant-api-wrapper:latest

# Ver logs
docker logs -f solarz-api

# Parar container
docker stop solarz-api
```

---

## 🐳 Exemplo 4: Docker Compose

### docker-compose.yml

```yaml
version: '3.8'

services:
  solarz-api:
    image: ghcr.io/catdevsecops/solarz-homeassistant-api-wrapper:latest
    container_name: solarz-api
    ports:
      - "8080:8080"
    environment:
      - ENV=production
      - LOG_LEVEL=info
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 5s
    networks:
      - solarz-network

  # Exemplo: Nginx reverse proxy
  nginx:
    image: nginx:latest
    ports:
      - "80:80"
    depends_on:
      - solarz-api
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
    networks:
      - solarz-network

networks:
  solarz-network:
    driver: bridge
```

### Executar

```bash
# Iniciar serviços
docker-compose up -d

# Ver status
docker-compose ps

# Ver logs da API
docker-compose logs -f solarz-api

# Parar serviços
docker-compose down
```

---

## 🐳 Exemplo 5: Kubernetes Deployment

### deploy.yaml

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: solarz-api
  namespace: default
spec:
  replicas: 3
  selector:
    matchLabels:
      app: solarz-api
  template:
    metadata:
      labels:
        app: solarz-api
    spec:
      containers:
      - name: solarz-api
        image: ghcr.io/catdevsecops/solarz-homeassistant-api-wrapper:latest
        imagePullPolicy: Always
        ports:
        - containerPort: 8080
        env:
        - name: ENV
          value: "production"
        - name: LOG_LEVEL
          value: "info"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 30
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 10
        resources:
          requests:
            memory: "64Mi"
            cpu: "100m"
          limits:
            memory: "256Mi"
            cpu: "500m"

---
apiVersion: v1
kind: Service
metadata:
  name: solarz-api
spec:
  type: LoadBalancer
  selector:
    app: solarz-api
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
```

### Deploy

```bash
# Criar deployment
kubectl apply -f deploy.yaml

# Ver pods
kubectl get pods

# Ver logs
kubectl logs -f deployment/solarz-api

# Escalar
kubectl scale deployment solarz-api --replicas=5

# Atualizar para nova versão
kubectl set image deployment/solarz-api \
  solarz-api=ghcr.io/catdevsecops/solarz-homeassistant-api-wrapper:1.2.3
```

---

## 🚀 Exemplo 6: CI/CD Pipeline Completo

### Cenário: Desenvolvimento + Deploy

#### Semana 1: Feature Development

```bash
# 2ª: Feature branch
git checkout -b feature/solar-data-api

# 3ª: Push code
git push origin feature/solar-data-api
# → ci-cd.yml: testes passam ✅
# → docker-publish.yml: build ok (sem push)

# 4ª: PR
gh pr create --title "Add solar data API"
# → Checks passam ✅

# 5ª: Merge
gh pr merge --squash
# → ci-cd.yml: final test ✅
# → docker-publish.yml: push imagem latest ✅

# ✅ Imagem live: ghcr.io/.../solarz-homeassistant-api-wrapper:latest
```

#### Semana 2: Release para Produção

```bash
# Preparar release
git tag v2.0.0
git push origin v2.0.0

# → build-release.yml: compila binários
# → docker-publish.yml: publica tags semver
#   - :2.0.0
#   - :2.0
#   - :2

# ✅ Tudo live em produção!

# Produção usa:
docker pull ghcr.io/catdevsecops/.../solarz-homeassistant-api-wrapper:2.0.0
```

---

## 📊 Exemplo 7: Monitorando o Build

### Via GitHub CLI

```bash
# Listar últimos runs
gh run list --workflow=docker-publish.yml

# Ver detalhes de um run
gh run view 12345 --log

# Cancelar run que travou
gh run cancel 12345

# Re-executar run
gh run rerun 12345
```

### Via GitHub Web

1. Repository → **Actions**
2. Selecione **docker-publish.yml**
3. Veja os runs
4. Expanda um run para ver detalhes
5. Veja logs de cada step

### Saída esperada do build bem-sucedido

```
✅ Build & Push Docker Image

  ✓ Checkout repository
  ✓ Set up Docker Buildx
  ✓ Log in to GitHub Container Registry
  ✓ Extract metadata
    - tags: [latest, main-a1b2c3d]
    - labels: [...]
  ✓ Build and push Docker image
    - Built image in 45s
    - Pushed to ghcr.io/catdevsecops/solarz-homeassistant-api-wrapper:latest
    - Pushed to ghcr.io/catdevsecops/solarz-homeassistant-api-wrapper:main-a1b2c3d
  ✓ Image digest
    - sha256:abc123def456...
```

---

## 🔍 Exemplo 8: Inspecionando Imagem Publicada

### Ver informações da imagem

```bash
# Pull e inspeccionar localmente
docker pull ghcr.io/catdevsecops/solarz-homeassistant-api-wrapper:latest
docker inspect ghcr.io/catdevsecops/solarz-homeassistant-api-wrapper:latest

# Ver histórico de layers
docker history ghcr.io/catdevsecops/solarz-homeassistant-api-wrapper:latest

# Executar interativo
docker run -it ghcr.io/catdevsecops/solarz-homeassistant-api-wrapper:latest sh

# Verificar health
docker run --rm \
  ghcr.io/catdevsecops/solarz-homeassistant-api-wrapper:latest \
  ./solarz-api --version
```

---

## 🚨 Exemplo 9: Troubleshooting

### Build está falhando?

```bash
# Ver logs detalhados
gh run view <run-id> --log

# Testar Dockerfile localmente
docker build .

# Validar sintaxe YAML
cat .github/workflows/docker-publish.yml | yq -P
```

### Imagem não está sendo publicada?

```bash
# Verificar: É push em main ou tag?
git branch -a          # Está em main?
git tag -l             # Criar tag v* para release

# Verificar permissões
gh api repos/{owner}/{repo}/actions/permissions

# Verificar registry
gh api repos/{owner}/{repo}/packages
```

### Imagem muito grande?

```bash
# Ver tamanho de cada layer
docker history ghcr.io/catdevsecops/.../solarz-homeassistant-api-wrapper:latest

# Otimizar Dockerfile (remover arquivos desnecessários)
# Usar alpine ao invés de full ubuntu
# Usar multi-stage builds
```

---

## 💡 Dicas e Tricks

### Usar version tag dinamicamente

```dockerfile
ARG VERSION=latest

FROM golang:1.23-alpine AS builder
ARG VERSION
WORKDIR /build
COPY . .
RUN go build -ldflags="-X main.Version=${VERSION}" -o app .

FROM alpine:latest
COPY --from=builder /build/app /app
CMD ["/app"]
```

### Build local para testar

```bash
# Simular build do workflow
docker buildx build --platform linux/amd64,linux/arm64 .

# Build com cache
docker build --cache-from=ghcr.io/.../solarz-homeassistant-api-wrapper:main .
```

### Atualizar todas imagens locais

```bash
# Pull todas as versões
docker pull ghcr.io/catdevsecops/.../solarz-homeassistant-api-wrapper:latest
docker pull ghcr.io/catdevsecops/.../solarz-homeassistant-api-wrapper:1
docker pull ghcr.io/catdevsecops/.../solarz-homeassistant-api-wrapper:1.2
docker pull ghcr.io/catdevsecops/.../solarz-homeassistant-api-wrapper:1.2.3

# Listar
docker images | grep ghcr.io
```

---

## 📚 Referências

- GitHub Actions: https://docs.github.com/en/actions
- Docker CLI: https://docs.docker.com/engine/reference/commandline/
- GHCR: https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry

---

**Mais exemplos?** Ver:
- `.github/CI-CD-GUIDE.md` - Guia completo
- `.github/DOCKER_PUBLISHING.md` - Docker specifics
- `.github/copilot-instructions.md` - Padrões de código
