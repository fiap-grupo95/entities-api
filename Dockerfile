FROM golang:1.24.4-alpine AS builder

WORKDIR /app

# Instala swag CLI para gerar docs
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY . .

RUN go mod download

# Gera documentação Swagger (ajuste o caminho do main.go se necessário)
RUN swag init -g cmd/api/main.go -o docs

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app/cmd/api

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} go build -o /mecanica-xpto-api

# Build do binário de migration
WORKDIR /app/cmd/migrate

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} go build -o /mecanica-xpto-migrate

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /mecanica-xpto-api .

COPY --from=builder /mecanica-xpto-migrate .

COPY --from=builder /app/docs ./docs

EXPOSE 8080

CMD ["./mecanica-xpto-api"]