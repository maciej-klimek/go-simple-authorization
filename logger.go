package main

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func init() {
	Log = logrus.New()

	Log.Out = os.Stdout

	Log.Level = logrus.DebugLevel
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: false,
		ForceColors:   true,
		DisableColors: false,
	})

	// JSON formatted logs instead
	// log.SetFormatter(&logrus.JSONFormatter{
	// 	PrettyPrint: true,
	// })
}

func ExampleLog() {
	Log.WithFields(logrus.Fields{
		"user_id":   1234,
		"operation": "fetch data",
	}).Info("Operation successful")
	Log.Debug("This is a debug message")
	Log.Warn("This is a warning message")
	Log.Error("This is an error message")
}
