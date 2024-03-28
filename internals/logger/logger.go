package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func Init(output *os.File) *logrus.Logger {
	var level logrus.Level
	if os.Getenv("ENV") == "development" {
		level = logrus.DebugLevel
	} else {
		level = logrus.InfoLevel
	}

	Logger = &logrus.Logger{
		Out:   io.MultiWriter(output, os.Stdout),
		Level: level,
		Formatter: &logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		},
	}
	Logger.SetReportCaller(true)

	return Logger
}

func Fatal(err error, message string) {
	Logger.WithError(err).Fatal(message)
}

func Error(err error, message string) {
	Logger.WithError(err).Error(message)
}

func Warn(err error, message string) {
	Logger.WithError(err).Warn(message)
}

func Info(message string) {
	Logger.Info(message)
}

func Debug(message string) {
	Logger.Debug(message)
}

