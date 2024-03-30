package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func Init(output *os.File) *logrus.Logger {
	var target io.Writer
	var level logrus.Level
	if os.Getenv("ENV") == "development" {
		level = logrus.DebugLevel
		target = io.MultiWriter(output, os.Stdout)
	} else {
		level = logrus.InfoLevel
		target = output
	}

	Logger = &logrus.Logger{
		Out:   target,
		Level: level,
		Formatter: &logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		},
	}
	Logger.SetReportCaller(true)

	return Logger
}

// Logs fattal error and exits the program fields are optional
func Fatal(err error, fields ...logrus.Fields) {
	if len(fields) > 0 {
		Logger.WithFields(fields[0]).Fatal(err)
	} else {
		Logger.Fatal(err)
	}
}

// Logs error fields are optional
func Error(err error, fields ...logrus.Fields) {
	if len(fields) > 0 {
		Logger.WithFields(fields[0]).Error(err)
	} else {
		Logger.Error(err)
	}
}

// Logs warning fields are optional
func Warn(err error, fields ...logrus.Fields) {
	if len(fields) > 0 {
		Logger.WithFields(fields[0]).Warn(err)
	} else {
		Logger.Warn(err)
	}
}

// Logs info fields are optional
func Info(message string, fields ...logrus.Fields) {
	if len(fields) > 0 {
		Logger.WithFields(fields[0]).Info(message)
	} else {
		Logger.Info(message)
	}
}

func Debug(message string, fields ...logrus.Fields) {
	if len(fields) > 0 {
		Logger.WithFields(fields[0]).Debug(message)
	} else {
		Logger.Debug(message)
	}
}
