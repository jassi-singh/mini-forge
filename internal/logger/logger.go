package logger

import (
	"fmt"
	"log"
	"os"
)

// LogLevel represents the logging level
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	ERROR
)

// Logger is a custom logger with support for INFO, DEBUG, and ERROR levels
type Logger struct {
	debugEnabled bool
	infoLogger   *log.Logger
	errorLogger  *log.Logger
	debugLogger  *log.Logger
}

var defaultLogger *Logger

// InitLogger initializes the default logger
func InitLogger(debugEnabled bool) {
	defaultLogger = &Logger{
		debugEnabled: debugEnabled,
		infoLogger:   log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger:  log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		debugLogger:  log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

// Info logs an informational message (always logged)
func Info(format string, v ...interface{}) {
	if defaultLogger == nil {
		InitLogger(false)
	}
	defaultLogger.infoLogger.Output(2, fmt.Sprintf(format, v...))
}

// Error logs an error message (always logged)
func Error(format string, v ...interface{}) {
	if defaultLogger == nil {
		InitLogger(false)
	}
	defaultLogger.errorLogger.Output(2, fmt.Sprintf(format, v...))
}

// Debug logs a debug message (only if debug is enabled)
func Debug(format string, v ...interface{}) {
	if defaultLogger == nil {
		InitLogger(false)
	}
	if defaultLogger.debugEnabled {
		defaultLogger.debugLogger.Output(2, fmt.Sprintf(format, v...))
	}
}

// Fatal logs an error message and exits the program
func Fatal(format string, v ...interface{}) {
	if defaultLogger == nil {
		InitLogger(false)
	}
	defaultLogger.errorLogger.Output(2, fmt.Sprintf(format, v...))
	os.Exit(1)
}
