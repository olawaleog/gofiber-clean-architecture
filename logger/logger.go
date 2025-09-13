package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

// LogLevel defines the log level
type LogLevel string

const (
	// Debug level
	Debug LogLevel = "debug"
	// Info level
	Info LogLevel = "info"
	// Warn level
	Warn LogLevel = "warn"
	// Error level
	Error LogLevel = "error"
	// Fatal level
	Fatal LogLevel = "fatal"
)

// LogContext stores contextual information for logging
type LogContext struct {
	File     string
	Line     int
	Function string
}

// NewLogger initializes and returns a configured logrus logger
func NewLogger() *logrus.Logger {
	logger := logrus.New()

	// Default to info level
	logger.SetLevel(logrus.InfoLevel)

	// Use JSON formatter for structured logging
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "message",
		},
	})

	// Create logs directory if it doesn't exist
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		err := os.Mkdir("logs", 0770)
		if err != nil {
			panic(fmt.Sprintf("Failed to create logs directory: %v", err))
		}
	}

	// Create log file with current date and hour
	date := time.Now()
	logFilePath := filepath.Join("logs", "log_"+date.Format("01-02-2006_15")+".log")
	logFile, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err == nil {
		// Output to both stdout and file
		multiWriter := io.MultiWriter(os.Stdout, logFile)
		logger.SetOutput(multiWriter)
	} else {
		fmt.Printf("Failed to open log file: %v. Logging to stdout only.\n", err)
	}

	Logger = logger
	return logger
}

// GetLoggerWithContext returns a logrus entry with file, line and function context
func GetLoggerWithContext() *logrus.Entry {
	// Skip 2 frames to get the caller of this function
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return Logger.WithFields(logrus.Fields{
			"file": "unknown",
			"line": 0,
		})
	}

	// Get function name
	fn := runtime.FuncForPC(pc)
	var fnName string
	if fn == nil {
		fnName = "unknown"
	} else {
		fnName = fn.Name()
	}

	// Use short file path
	shortFile := filepath.Base(file)

	return Logger.WithFields(logrus.Fields{
		"file":     shortFile,
		"line":     line,
		"function": fnName,
	})
}

// LogWithLevel logs a message with the specified level and error context
func LogWithLevel(level LogLevel, msg string, err error, context ...interface{}) {
	entry := GetLoggerWithContext()

	// Add error information if provided
	if err != nil {
		entry = entry.WithField("error", err.Error())
	}

	// Add additional context if provided
	if len(context) > 0 {
		if len(context)%2 == 0 {
			for i := 0; i < len(context); i += 2 {
				if key, ok := context[i].(string); ok {
					entry = entry.WithField(key, context[i+1])
				}
			}
		} else {
			entry = entry.WithField("context", context)
		}
	}

	switch level {
	case Debug:
		entry.Debug(msg)
	case Info:
		entry.Info(msg)
	case Warn:
		entry.Warn(msg)
	case Error:
		entry.Error(msg)
	case Fatal:
		entry.Fatal(msg)
	default:
		entry.Info(msg)
	}
}
