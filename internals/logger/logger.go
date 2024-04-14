package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

type ILogger interface {
	Fatal(err error, fields ...LogFields)
	Error(err error, fields ...LogFields)
	Warn(err error, fields ...LogFields)
	Info(message string, fields ...LogFields)
	Debug(message string, fields ...LogFields)
}

type Logger struct {
	logger *logrus.Logger
}

type LogFields logrus.Fields

var logger *Logger

// If logger used in a package that should be tested, you should use Logger class methods instead of global functions.
func NewLogger(output *os.File) *Logger {
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

	if logger == nil {
		logger = &Logger{
			logger: &logrus.Logger{
				Out:       target,
				Level:     level,
				Formatter: format,
			},
		}
	}

	return logger
}

// Logs fattal error and exits the program fields are optional
func (l *Logger) Fatal(err error, fields ...LogFields) {
	if len(fields) > 0 {
		l.logger.WithFields(logrus.Fields(fields[0])).Fatal(err)
	} else {
		l.logger.Fatal(err)
	}
}

// Logs error fields are optional
func (l *Logger) Error(err error, fields ...LogFields) {
	if len(fields) > 0 {
		l.logger.WithFields(logrus.Fields(fields[0])).Error(err)
	} else {
		l.logger.Error(err)
	}
}

// Logs warning fields are optional
func (l *Logger) Warn(err error, fields ...LogFields) {
	if len(fields) > 0 {
		l.logger.WithFields(logrus.Fields(fields[0])).Warn(err)
	} else {
		l.logger.Warn(err)
	}
}

// Logs info fields are optional
func (l *Logger) Info(message string, fields ...LogFields) {
	if len(fields) > 0 {
		l.logger.WithFields(logrus.Fields(fields[0])).Info(message)
	} else {
		l.logger.Info(message)
	}
}

func (l *Logger) Debug(message string, fields ...LogFields) {
	if len(fields) > 0 {
		l.logger.WithFields(logrus.Fields(fields[0])).Debug(message)
	} else {
		l.logger.Debug(message)
	}
}

// Logs fatal error and exits the program fields are optional
func Fatal(err error, fields ...LogFields) {
	logger.Fatal(err, fields...)
}

// Logs error fields are optional
func Error(err error, fields ...LogFields) {
	logger.Error(err, fields...)
}

// Logs warning fields are optional
func Warn(err error, fields ...LogFields) {
	logger.Warn(err, fields...)
}

// Logs info fields are optional
func Info(message string, fields ...LogFields) {
	logger.Info(message, fields...)
}

// Logs debug fields are optional
func Debug(message string, fields ...LogFields) {
	logger.Debug(message, fields...)
}
