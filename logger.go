package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func ConfigureLogger(level string) {
	log.SetReportCaller(true)
	log.SetFormatter(&logrus.TextFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			filename := filepath.Base(frame.File)
			padding := 20
			fileLine := filename + ":" + strconv.Itoa(frame.Line)

			if len(fileLine) < padding {
				fileLine = fileLine + strings.Repeat(" ", padding-len(fileLine))
			}

			return ">", " < " + fileLine
		},
	})
	switch level {
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	default:
		log.SetLevel(logrus.InfoLevel)
	}
	log.SetOutput(os.Stdout)
}
