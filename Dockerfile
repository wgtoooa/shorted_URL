# Builder stage
FROM golang:alpine AS builder

# Создаём рабочую директорию
WORKDIR /app

# Копируем только файлы зависимостей (для кэширования)
COPY go.mod go.sum ./
RUN go mod download

# Копируем ВСЕ файлы проекта (включая .go)
COPY . .

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main .

# Final stage
FROM alpine:3.19
WORKDIR /app

# Копируем бинарник и .env
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

# Устанавливаем tzdata
RUN apk add --no-cache tzdata

# Настраиваем пользователя
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s \
  CMD wget --spider http://localhost:8080/health || exit 1

CMD ["./main"]