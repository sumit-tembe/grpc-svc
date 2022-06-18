package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logger ...
var Logger = newLogger("debug")

func newLogger(logLevelStr string) *logrus.Logger {
	logger := logrus.New()
	// Output to stdout instead of the default stderr
	logger.SetOutput(os.Stdout)
	logLevel, err := logrus.ParseLevel(logLevelStr)
	if err != nil {
		// Set Info level as a default
		logLevel = logrus.InfoLevel
	}
	logger.SetLevel(logLevel)
	logger.Formatter = &logrus.TextFormatter{FullTimestamp: true}
	return logger
}
