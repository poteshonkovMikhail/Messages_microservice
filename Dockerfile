# Используем официальный образ Golang как базовый
FROM golang:1.22 AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы go.mod и go.sum для установки зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Используем легкий образ для запуска приложения
FROM alpine:latest  

WORKDIR /root/

# Копируем собранное приложение из этапа сборки
COPY --from=builder /app/main .

# Устанавливаем необходимые зависимости
RUN apk --no-cache add ca-certificates && update-ca-certificates

# Задаем команду для запуска приложения
CMD ["./main"]