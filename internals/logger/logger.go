package logger

import (
	"fmt"
	"io"
	"os"

	"github.com/rs/zerolog"
)

type ILogger interface {
	Fatal(err error, fields ...map[string]any)
	Error(err error, fields ...map[string]any)
	Warn(err error, fields ...map[string]any)
	Info(message string, fields ...map[string]any)
	Debug(message string, fields ...map[string]any)
}

type Logger struct {
	logger zerolog.Logger
}

var logger *Logger

// If logger used in a package that should be tested, you should use Logger class methods instead of global functions.
func NewLogger(output *os.File) *Logger {
	var target io.Writer
	env := os.Getenv("ENV")

	if env == "development" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		consoleWriter := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "2006-01-02 15:04:05",
		}
		target = zerolog.MultiLevelWriter(consoleWriter, output)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		target = output
	}

	if logger == nil {
		logger = &Logger{
			logger: zerolog.New(target).With().Timestamp().Logger(),
		}
	}

	return logger
}

// Logs fattal error and exits the program fields are optional
// Field with key "message" will be ignored
func (l *Logger) Fatal(err error, fields ...map[string]any) {
	if len(fields) > 0 {
		l.logger.Fatal().Err(err).Fields(fields[0]).Msg("")
	} else {
		l.logger.Fatal().Err(err).Msg("")
	}
}

// Logs error fields are optional
// Field with key "message" will be ignored
func (l *Logger) Error(err error, fields ...map[string]any) {
	if len(fields) > 0 {
		l.logger.Error().Err(err).Fields(fields[0]).Msg("")
	} else {
		l.logger.Error().Err(err).Msg("")
	}
}

// Logs warning fields are optional
// Field with key "message" will be ignored
func (l *Logger) Warn(err error, fields ...map[string]any) {
	if len(fields) > 0 {
		l.logger.Warn().Err(err).Fields(fields[0]).Msg("")
	} else {
		l.logger.Warn().Err(err).Msg("")
	}
}

// Logs info fields are optional
// Field with key "message" will be ignored
func (l *Logger) Info(message string, fields ...map[string]any) {
	if len(fields) > 0 {
		l.logger.Info().Fields(fields[0]).Msg(message)
	} else {
		l.logger.Info().Msg(message)
	}
}

// Logs debug fields are optional
// Field with key "message" will be ignored
func (l *Logger) Debug(message string, fields ...map[string]any) {
	if len(fields) > 0 {
		fmt.Println(fields[0])
		l.logger.Debug().Fields(fields[0]).Msg(message)
	} else {
		l.logger.Debug().Msg(message)
	}
}

// Logs fatal error and exits the program fields are optional
func Fatal(err error, fields ...map[string]any) {
	logger.Fatal(err, fields...)
}

// Logs error fields are optional
func Error(err error, fields ...map[string]any) {
	logger.Error(err, fields...)
}

// Logs warning fields are optional
func Warn(err error, fields ...map[string]any) {
	logger.Warn(err, fields...)
}

// Logs info fields are optional
func Info(message string, fields ...map[string]any) {
	logger.Info(message, fields...)
}

// Logs debug fields are optional
func Debug(message string, fields ...map[string]any) {
	logger.Debug(message, fields...)
}

func ToMap[T map[string]string | map[string]any](fields T) map[string]any {
	if fields == nil {
		return map[string]any{}
	}

	af := any(fields)

	if asf, ok := af.(map[string]interface{}); ok {
		return map[string]any(asf)
	}

	afm := af.(map[string]string)

	afmRes := make(map[string]any, len(afm))

	for dick, ass := range afm {
		afmRes[dick] = ass
	}

	return afmRes
}
