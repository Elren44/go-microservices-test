# ---- Build stage ----
FROM golang:1.24-alpine AS builder

# Установим зависимости
RUN apk add --no-cache git

# Рабочая директория для сборки
WORKDIR /app

# Только go.mod и go.sum сначала (для кэширования go mod download)
COPY go.mod go.sum ./
RUN go mod download

# Копируем всё нужное: cmd, internal, pkg, config и т.д.
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY pkg/ ./pkg/
COPY config/ ./config/
# COPY .env .env  # если используется

# Сборка конкретного сервиса
RUN go build -o auth-service ./cmd/app

# ---- Final stage ----
FROM alpine:latest

# Добавим CA сертификаты (если сервис делает внешние HTTPS-запросы)
RUN apk add --no-cache ca-certificates

WORKDIR /root/

# Копируем бинарник
COPY --from=builder /app/auth-service .


# Команда запуска
CMD ["./auth-service"]
