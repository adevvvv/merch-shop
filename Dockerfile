# Stage 1: сборка приложения
FROM golang:1.22.2-alpine AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файлы go.mod и go.sum и загружаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем бинарный файл (без поддержки CGO для минимального образа)
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd


# Stage 2: минимальный образ для запуска приложения
FROM alpine:latest

WORKDIR /app

# Копируем собранный бинарник из предыдущего этапа
COPY --from=builder /app/main .

# Пробрасываем порт 8080
EXPOSE 8080

# Запускаем приложение
CMD ["./main"]
