package logger

import (
	"context"
	"sync"
)

// Global logger instance
var globalLogger Logger
var once sync.Once

// Initialize sets up the global logger with the given configuration
func Initialize(config Config) error {
	logger, err := New(config)
	if err != nil {
		return err
	}
	globalLogger = logger
	return nil
}

// GetLogger returns the global logger instance
func GetLogger() Logger {
	once.Do(func() {
		if globalLogger == nil {
			// Initialize with default config if not initialized
			config := DefaultConfig()
			_ = Initialize(config)
		}
	})
	return globalLogger
}

// Debug logs a message at level Debug using the global logger
func Debug(args ...interface{}) {
	GetLogger().Debug(args...)
}

// Debugf logs a formatted message at level Debug using the global logger
func Debugf(format string, args ...interface{}) {
	GetLogger().Debugf(format, args...)
}

// Info logs a message at level Info using the global logger
func Info(args ...interface{}) {
	GetLogger().Info(args...)
}

// Infof logs a formatted message at level Info using the global logger
func Infof(format string, args ...interface{}) {
	GetLogger().Infof(format, args...)
}

// Warn logs a message at level Warn using the global logger
func Warn(args ...interface{}) {
	GetLogger().Warn(args...)
}

// Warnf logs a formatted message at level Warn using the global logger
func Warnf(format string, args ...interface{}) {
	GetLogger().Warnf(format, args...)
}

// Error logs a message at level Error using the global logger
func Error(args ...interface{}) {
	GetLogger().Error(args...)
}

// Errorf logs a formatted message at level Error using the global logger
func Errorf(format string, args ...interface{}) {
	GetLogger().Errorf(format, args...)
}

// Fatal logs a message at level Fatal using the global logger then the process will exit with status set to 1
func Fatal(args ...interface{}) {
	GetLogger().Fatal(args...)
}

// Fatalf logs a formatted message at level Fatal using the global logger then the process will exit with status set to 1
func Fatalf(format string, args ...interface{}) {
	GetLogger().Fatalf(format, args...)
}

// Panic logs a message at level Panic using the global logger then panics
func Panic(args ...interface{}) {
	GetLogger().Panic(args...)
}

// Panicf logs a formatted message at level Panic using the global logger then panics
func Panicf(format string, args ...interface{}) {
	GetLogger().Panicf(format, args...)
}

// WithField adds a single field to the global logger
func WithField(key string, value interface{}) Logger {
	return GetLogger().WithField(key, value)
}

// WithFields adds multiple fields to the global logger
func WithFields(fields Fields) Logger {
	return GetLogger().WithFields(fields)
}

// WithError adds an error field to the global logger
func WithError(err error) Logger {
	return GetLogger().WithError(err)
}

// WithContext adds context to the global logger
func WithContext(ctx context.Context) Logger {
	return GetLogger().WithContext(ctx)
}
