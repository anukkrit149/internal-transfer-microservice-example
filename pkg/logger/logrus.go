package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
)

// LogrusLogger is a logrus implementation of the Logger interface
type LogrusLogger struct {
	entry *logrus.Entry
}

// New creates a new logger with the given configuration
func New(config Config) (Logger, error) {
	log := logrus.New()

	// Set log level
	level, err := logrus.ParseLevel(string(config.Level))
	if err != nil {
		return nil, fmt.Errorf("invalid log level: %s", config.Level)
	}
	log.SetLevel(level)

	// Set log format
	switch config.Format {
	case JSONFormat:
		log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				filename := filepath.Base(f.File)
				return "", fmt.Sprintf("%s:%d", filename, f.Line)
			},
		})
	case TextFormat:
		log.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: time.RFC3339Nano,
			FullTimestamp:   true,
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				filename := filepath.Base(f.File)
				return "", fmt.Sprintf("%s:%d", filename, f.Line)
			},
		})
	default:
		return nil, fmt.Errorf("invalid log format: %s", config.Format)
	}

	// Set log output
	switch config.Output {
	case "stdout":
		log.SetOutput(os.Stdout)
	case "stderr":
		log.SetOutput(os.Stderr)
	default:
		// Assume it's a file path
		file, err := os.OpenFile(config.Output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %v", err)
		}
		log.SetOutput(file)
	}

	// Set caller reporting
	log.SetReportCaller(config.ReportCaller)

	return &LogrusLogger{
		entry: logrus.NewEntry(log),
	}, nil
}

// NewWithWriter creates a new logger with a custom writer
func NewWithWriter(writer io.Writer, level LogLevel, format LogFormat) Logger {
	log := logrus.New()

	// Set log level
	parsedLevel, _ := logrus.ParseLevel(string(level))
	log.SetLevel(parsedLevel)

	// Set log format
	if format == JSONFormat {
		log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
		})
	} else {
		log.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: time.RFC3339Nano,
			FullTimestamp:   true,
		})
	}

	log.SetOutput(writer)
	log.SetReportCaller(true)

	return &LogrusLogger{
		entry: logrus.NewEntry(log),
	}
}

// WithField adds a single field to the log entry
func (l *LogrusLogger) WithField(key string, value interface{}) Logger {
	return &LogrusLogger{
		entry: l.entry.WithField(key, value),
	}
}

// WithFields adds multiple fields to the log entry
func (l *LogrusLogger) WithFields(fields Fields) Logger {
	return &LogrusLogger{
		entry: l.entry.WithFields(logrus.Fields(fields)),
	}
}

// WithError adds an error field to the log entry
func (l *LogrusLogger) WithError(err error) Logger {
	return &LogrusLogger{
		entry: l.entry.WithError(err),
	}
}

// WithContext adds context to the log entry
func (l *LogrusLogger) WithContext(ctx context.Context) Logger {
	return &LogrusLogger{
		entry: l.entry.WithContext(ctx),
	}
}

// Debug logs a message at level Debug
func (l *LogrusLogger) Debug(args ...interface{}) {
	l.entry.Debug(args...)
}

// Debugf logs a formatted message at level Debug
func (l *LogrusLogger) Debugf(format string, args ...interface{}) {
	l.entry.Debugf(format, args...)
}

// Info logs a message at level Info
func (l *LogrusLogger) Info(args ...interface{}) {
	l.entry.Info(args...)
}

// Infof logs a formatted message at level Info
func (l *LogrusLogger) Infof(format string, args ...interface{}) {
	l.entry.Infof(format, args...)
}

// Warn logs a message at level Warn
func (l *LogrusLogger) Warn(args ...interface{}) {
	l.entry.Warn(args...)
}

// Warnf logs a formatted message at level Warn
func (l *LogrusLogger) Warnf(format string, args ...interface{}) {
	l.entry.Warnf(format, args...)
}

// Error logs a message at level Error
func (l *LogrusLogger) Error(args ...interface{}) {
	l.entry.Error(args...)
}

// Errorf logs a formatted message at level Error
func (l *LogrusLogger) Errorf(format string, args ...interface{}) {
	l.entry.Errorf(format, args...)
}

// Fatal logs a message at level Fatal then the process will exit with status set to 1
func (l *LogrusLogger) Fatal(args ...interface{}) {
	l.entry.Fatal(args...)
}

// Fatalf logs a formatted message at level Fatal then the process will exit with status set to 1
func (l *LogrusLogger) Fatalf(format string, args ...interface{}) {
	l.entry.Fatalf(format, args...)
}

// Panic logs a message at level Panic then panics
func (l *LogrusLogger) Panic(args ...interface{}) {
	l.entry.Panic(args...)
}

// Panicf logs a formatted message at level Panic then panics
func (l *LogrusLogger) Panicf(format string, args ...interface{}) {
	l.entry.Panicf(format, args...)
}
