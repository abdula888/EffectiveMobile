package log

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func init() {
	// Загружаем .env файл (если он существует)
	_ = godotenv.Load()

	// Создаем новый экземпляр логгера
	Logger = logrus.New()

	// Читаем уровень логирования из переменных окружения
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	// Устанавливаем уровень логирования
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		Logger.Fatalf("Invalid log level: %s", logLevel)
	}
	Logger.SetLevel(level)

	// Настраиваем формат вывода (например, JSON или текст)
	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}
