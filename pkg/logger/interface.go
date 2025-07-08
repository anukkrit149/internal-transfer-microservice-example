package logger

import (
	"context"
)

// Fields type to pass multiple fields to the logger
type Fields map[string]interface{}

// Logger interface defines the logging methods
type Logger interface {
	// WithField adds a single field to the log entry
	WithField(key string, value interface{}) Logger
	// WithFields adds multiple fields to the log entry
	WithFields(fields Fields) Logger
	// WithError adds an error field to the log entry
	WithError(err error) Logger
	// WithContext adds context to the log entry
	WithContext(ctx context.Context) Logger
	// Debug logs a message at level Debug
	Debug(args ...interface{})
	// Debugf logs a formatted message at level Debug
	Debugf(format string, args ...interface{})
	// Info logs a message at level Info
	Info(args ...interface{})
	// Infof logs a formatted message at level Info
	Infof(format string, args ...interface{})
	// Warn logs a message at level Warn
	Warn(args ...interface{})
	// Warnf logs a formatted message at level Warn
	Warnf(format string, args ...interface{})
	// Error logs a message at level Error
	Error(args ...interface{})
	// Errorf logs a formatted message at level Error
	Errorf(format string, args ...interface{})
	// Fatal logs a message at level Fatal then the process will exit with status set to 1
	Fatal(args ...interface{})
	// Fatalf logs a formatted message at level Fatal then the process will exit with status set to 1
	Fatalf(format string, args ...interface{})
	// Panic logs a message at level Panic then panics
	Panic(args ...interface{})
	// Panicf logs a formatted message at level Panic then panics
	Panicf(format string, args ...interface{})
}

// LogLevel represents the level of logging
type LogLevel string

// Log levels
const (
	DebugLevel LogLevel = "debug"
	InfoLevel  LogLevel = "info"
	WarnLevel  LogLevel = "warn"
	ErrorLevel LogLevel = "error"
	FatalLevel LogLevel = "fatal"
	PanicLevel LogLevel = "panic"
)

// LogFormat represents the format of logging
type LogFormat string

// Log formats
const (
	JSONFormat LogFormat = "json"
	TextFormat LogFormat = "text"
)

// Config holds the configuration for the logger
type Config struct {
	// Level is the log level (debug, info, warn, error, fatal, panic)
	Level LogLevel
	// Format is the log format (json, text)
	Format LogFormat
	// Output is the log output (stdout, stderr, file)
	Output string
	// ReportCaller adds the file and line number to the log
	ReportCaller bool
}

// DefaultConfig returns a default configuration for the logger
func DefaultConfig() Config {
	return Config{
		Level:        InfoLevel,
		Format:       JSONFormat,
		Output:       "stdout",
		ReportCaller: true,
	}
}
