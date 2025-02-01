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

	if stage == "" {
		stage = os.Getenv("STAGE")
	}

	switch stage {
	case "lab", "dev", "development", "local":
		logger.SetLevel(logrus.DebugLevel)
	case "prod", "production":
		logger.SetLevel(logrus.ErrorLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}

	return logger
}

// Gets the Global Logger
func GetLogger() *logrus.Logger {
	return logger
}
