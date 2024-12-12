package log

import (
	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func SetUpLogger(logLevel string) {
	// Создаем новый экземпляр логгера
	Logger = logrus.New()

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
