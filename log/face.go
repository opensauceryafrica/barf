package logger

type logger struct{}

// Logger returns a new barf logger instance
func Logger() *logger {
	return &logger{}
}

// Info logs a message with the green color
func (l *logger) Info(msg string) {
	Info(msg)
}

// Error logs a message with the red color
func (l *logger) Error(msg string) {
	Error(msg)
}

// Debug logs a message with the blue color
func (l *logger) Debug(msg string) {
	Debug(msg)
}

// Code is a function that logs based on the status code
func (l *logger) Code(msg string, code int) {
	Code(msg, code)
}

// Warn logs a message with the yellow color
func (l *logger) Warn(msg string) {
	Warn(msg)
}
