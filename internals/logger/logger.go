package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

type LogFields logrus.Fields

var Logger *logrus.Logger

func Init(output *os.File) *logrus.Logger {
	var target io.Writer
	var level logrus.Level
	var format logrus.Formatter
	env := os.Getenv("ENV")

	if env == "development" {
		level = logrus.DebugLevel
		target = io.MultiWriter(output, os.Stdout)
		format = &logrus.TextFormatter{
			ForceColors:     true,
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		}
	} else {
		level = logrus.InfoLevel
		target = output
		format = &logrus.JSONFormatter{}
	}

	Logger = &logrus.Logger{
		Out:       target,
		Level:     level,
		Formatter: format,
	}

	return Logger
}

// Logs fattal error and exits the program fields are optional
func Fatal(err error, fields ...LogFields) {
	if len(fields) > 0 {
		Logger.WithFields(logrus.Fields(fields[0])).Fatal(err)
	} else {
		Logger.Fatal(err)
	}
}

// Logs error fields are optional
func Error(err error, fields ...LogFields) {
	if len(fields) > 0 {
		Logger.WithFields(logrus.Fields(fields[0])).Error(err)
	} else {
		Logger.Error(err)
	}
}

// Logs warning fields are optional
func Warn(err error, fields ...LogFields) {
	if len(fields) > 0 {
		Logger.WithFields(logrus.Fields(fields[0])).Warn(err)
	} else {
		Logger.Warn(err)
	}
}

// Logs info fields are optional
func Info(message string, fields ...LogFields) {
	if len(fields) > 0 {
		Logger.WithFields(logrus.Fields(fields[0])).Info(message)
	} else {
		Logger.Info(message)
	}
}

func Debug(message string, fields ...LogFields) {
	if len(fields) > 0 {
		Logger.WithFields(logrus.Fields(fields[0])).Debug(message)
	} else {
		Logger.Debug(message)
	}
}
