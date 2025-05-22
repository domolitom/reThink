package utils

import (
	"log"
	"os"
)

// Logger represents a simple logger
type Logger struct {
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	DebugLogger *log.Logger
}

// NewLogger creates a new logger instance
func NewLogger() *Logger {
	return &Logger{
		InfoLogger:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		ErrorLogger: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		DebugLogger: log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

// Info logs an info message
func (l *Logger) Info(v ...interface{}) {
	l.InfoLogger.Println(v...)
}

// Error logs an error message
func (l *Logger) Error(v ...interface{}) {
	l.ErrorLogger.Println(v...)
}

// Debug logs a debug message if debug mode is enabled
func (l *Logger) Debug(v ...interface{}) {
	// Only log if DEBUG=true is set
	if os.Getenv("DEBUG") == "true" {
		l.DebugLogger.Println(v...)
	}
}
