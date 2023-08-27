package logger

import (
	"fmt"
	"os"
)

type logger struct{}

// Logger returns a new barf logger instance
func Logger() *logger {
	return &logger{}
}

// Info logs a message with the green color
func (l *logger) Info(msg string) {
	Info(msg)
}

// Infof logs a formatted message with the green color
func (l *logger) Infof(msg string, v ...interface{}) {
	Info(fmt.Sprintf(msg, v...))
}

// Error logs a message with the red color
func (l *logger) Error(msg string) {
	Error(msg)
}

// Errorf logs a formatted message with the red color
func (l *logger) Errorf(msg string, v ...interface{}) {
	Error(fmt.Sprintf(msg, v...))
}

// Fatal logs a message with the red color and exits the program
func (l *logger) Fatal(msg string) {
	Error(msg)
	os.Exit(1)
}

// Fatalf logs a formatted message with the red color and exits the program
func (l *logger) Fatalf(msg string, v ...interface{}) {
	Error(fmt.Sprintf(msg, v...))
	os.Exit(1)
}

// Debug logs a message with the blue color
func (l *logger) Debug(msg string) {
	Debug(msg)
}

// Debugf logs a formatted message with the blue color
func (l *logger) Debugf(msg string, v ...interface{}) {
	Debug(fmt.Sprintf(msg, v...))
}

// Code is a function that logs based on the status code
func (l *logger) Code(msg string, code int) {
	Code(msg, code)
}

// Codef is a function that logs based on the status code
func (l *logger) Codef(msg string, code int, v ...interface{}) {
	Code(fmt.Sprintf(msg, v...), code)
}

// Warn logs a message with the yellow color
func (l *logger) Warn(msg string) {
	Warn(msg)
}

// Warnf logs a formatted message with the yellow color
func (l *logger) Warnf(msg string, v ...interface{}) {
	Warn(fmt.Sprintf(msg, v...))
}
