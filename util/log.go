package util

import "github.com/sirupsen/logrus"

// NewLogger creates a new instance of logrus with project specific
// configuration
func NewLogger() *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	log.SetFormatter(&logrus.JSONFormatter{})

	return log
}
