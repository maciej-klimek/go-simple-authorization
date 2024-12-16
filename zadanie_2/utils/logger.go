package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func init() {
	Logger = logrus.New()

	Logger.Out = os.Stdout

	Logger.Level = logrus.DebugLevel
	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: false,
		ForceColors:   true,
		DisableColors: false,
	})

}
