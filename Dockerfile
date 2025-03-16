FROM golang:1.23-alpine AS builder

WORKDIR /app

# Устанавливаем часовой пояс
ENV TZ=Europe/Moscow

RUN apk add --no-cache tzdata
RUN ln -sf /usr/share/zoneinfo/Europe/Moscow /etc/localtime

# Копируем файлы go.mod и go.sum и устанавливаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем бинарник
RUN go build -o bot ./cmd/main.go

# Финальный образ
FROM alpine:latest

WORKDIR /root/

# Устанавливаем часовой пояс
ENV TZ=Europe/Moscow

RUN apk add --no-cache tzdata && ln -sf /usr/share/zoneinfo/Europe/Moscow /etc/localtime

# Копируем собранное приложение
COPY --from=builder /app/bot .
COPY --from=builder /app/internal/app/loader/enriched_cities.json internal/app/loader/


# Запуск бота
CMD ["./bot"]