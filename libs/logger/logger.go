package logger

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

type MyFormatter struct {
	logrus.JSONFormatter
}

var levelList = []string{
	"PANIC",
	"FATAL",
	"ERROR",
	"WARN",
	"INFO",
	"DEBUG",
	"TRACE",
}

func (mf *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	level := levelList[int(entry.Level)]
	switch level {
	case "ERROR":
		level = color.RedString("%s ", " ERROR")
	case "WARN":
		level = color.YellowString("%s ", " WARN")
	case "INFO":
		level = color.GreenString("%s ", " INFO")
	case "DEBUG":
		level = color.HiBlueString("%s ", " DEBUG")
	default:
		level = color.HiMagentaString("%s ", level)
	}
	b.WriteString(fmt.Sprintf("%s - %s - [line:%d] - %s msg ▶ %s\n",
		entry.Time.Format("2006-01-02 15:04:05,678"), entry.Caller.File,
		entry.Caller.Line, level, entry.Message))
	return b.Bytes(), nil
}

func NewLogger() *logrus.Logger {
	logger := logrus.New()

	// Set log level from environment variable
	logLevel := strings.ToLower(os.Getenv("LOG_LEVEL"))
	switch logLevel {
	case "trace":
		logger.SetLevel(logrus.TraceLevel)
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warn", "warning":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	case "fatal":
		logger.SetLevel(logrus.FatalLevel)
	case "panic":
		logger.SetLevel(logrus.PanicLevel)
	default:
		// Default to info level if not specified or invalid
		logger.SetLevel(logrus.InfoLevel)
	}

	logger.SetReportCaller(true)

	// Use JSON formatter for production/staging (Cloud Run friendly)
	// Use colored formatter for development
	logger.SetFormatter(&MyFormatter{
		JSONFormatter: logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "severity",
				logrus.FieldKeyMsg:   "message",
				logrus.FieldKeyFunc:  "function",
				logrus.FieldKeyFile:  "file",
			},
		},
	})

	return logger
}
