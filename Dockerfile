# Multi-stage build Dockerfile - Multi-architecture support
# Supports: linux/amd64, linux/arm64, and others

# Stage 1: Builder
FROM --platform=$BUILDPLATFORM golang:1.23-alpine AS builder

ARG BUILDPLATFORM
ARG TARGETPLATFORM
ARG TARGETARCH
ARG TARGETOS

WORKDIR /build

# Instalar dependências
RUN apk add --no-cache git ca-certificates tzdata

# Copiar módulos
COPY go.mod ./

# Baixar dependências
RUN go mod download

# Copiar código
COPY . .

# Build com arquitetura alvo
# Usa variáveis do buildkit para compilar para a arquitetura correta
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build \
    -ldflags="-w -s" \
    -o solarz-api \
    .

# Stage 2: Runtime
FROM alpine:latest

LABEL maintainer="Clayton Silva"
LABEL description="Solarz API - REST API para gerenciamento de dados de geração solar"

RUN apk add --no-cache ca-certificates tzdata

# Criar usuário não-root
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /home/appuser

# Copiar certificados do builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copiar binário do builder
COPY --from=builder --chown=appuser:appuser /build/solarz-api /home/appuser/

# Trocar para usuário não-root
USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

ENTRYPOINT [ "/home/appuser/solarz-api" ]
