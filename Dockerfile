# Используем образ Golang
FROM golang

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы
COPY ./go.mod .
COPY ./go.sum .

# Загружаем зависимости
RUN go mod download

COPY . .

# Собираем приложение
RUN go build -o main ./cmd/effective-mobile/

# Указываем команду для запуска приложения
CMD ["./main"]

