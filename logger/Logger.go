package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

// Instances a new Global Logger
func NewLogger(stage string) *logrus.Logger {
	logger = logrus.New()
	logger.SetOutput(os.Stdout)

	switch stage {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	default:
		logger.SetLevel(logrus.ErrorLevel)
	}

	return logger
}

// Gets the Global Logger
func GetLogger() *logrus.Logger {
	return logger
}
