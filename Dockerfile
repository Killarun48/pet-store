# Используем официальный образ Go как базовый
FROM golang:1.22-alpine as builder

ENV CGO_ENABLED=1

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем исходники приложения в рабочую директорию
COPY . .

# Удаляем лишние файлы
RUN rm -rf /app/go.mod
RUN rm -rf /app/go.sum

# Скачиваем все зависимости
RUN go mod init app && go mod tidy

RUN apk add --no-cache git gcc musl-dev make
ENV GO111MODULE=on

# Собираем приложение
RUN go build -o main

# Начинаем новую стадию сборки на основе минимального образа
FROM alpine:latest

# Добавляем исполняемый файл из первой стадии в корневую директорию контейнера
COPY --from=builder /app/main /main

COPY .env .
COPY /internal/infrastructure/db/migrations/sqlite3 /internal/infrastructure/db/migrations/sqlite3

# Открываем порт 8080
#EXPOSE 8080

# Запускаем приложение
CMD ["/main"]
