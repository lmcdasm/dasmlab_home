package logutil

import (
	"os"

	"github.com/sirupsen/logrus"
)

func InitLogger(component string) *logrus.Entry {
	logger := logrus.New()
	logger.Out = os.Stdout
	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&logrus.JSONFormatter{})
	return logger.WithField("component", component)
}
