# ==============================================================================
# DOCKERFILE MULTI-STAGE: VAMOS-LANG
# Imagem de produção ultra-leve (< 20MB) compilando código VAMOS para binário nativo.
# ==============================================================================

# Estágio 1: Compilação com o transpilador VAMOS-LANG
FROM golang:1.22-alpine AS construtor

WORKDIR /app

# Instala certificados e ferramentas básicas de build
RUN apk add --no-cache git ca-certificates tzdata

# Copia dependências do módulo Go
COPY go.mod go.sum ./
RUN go mod download

# Copia todo o código-fonte do projeto
COPY . .

# Compila o binário da CLI do VAMOS-LANG
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /app/bin/vamos .

# Compila a aplicação VAMOS-LANG (exemplo API REST) para binário de produção
RUN /app/bin/vamos build -o /app/bin/servidor_app exemplos/api_rest.vamos

# Estágio 2: Imagem Final de Execução Minimalista
FROM alpine:3.20

WORKDIR /app

# Copia certificados SSL e zona de horário
COPY --from=construtor /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=construtor /usr/share/zoneinfo /usr/share/zoneinfo

# Copia apenas o binário executável final compilado
COPY --from=construtor /app/bin/servidor_app /app/servidor_app

# Define usuário não-root por segurança
RUN addgroup -S vamosgroup && adduser -S vamosuser -G vamosgroup
USER vamosuser

EXPOSE 8080

ENTRYPOINT ["/app/servidor_app"]
