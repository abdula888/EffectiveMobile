# Используем образ Golang
FROM golang

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы
COPY . .

# Загружаем зависимости
RUN go mod download

# Собираем приложение
RUN go build -o main ./cmd

# Указываем команду для запуска приложения
CMD ["./main"]

