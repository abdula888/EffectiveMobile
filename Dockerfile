# Используем образ Golang
FROM golang:1.23-alpine

# Устанавливаем рабочую директорию
WORKDIR /app

COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

CMD ["go", "run", "./cmd/effective-mobile/"]
