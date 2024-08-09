package logger

// Code is a function that logs based on the status code
func Code(msg string, code int) {
	switch {
	case code >= 500:
		Error(msg)
	case code >= 400:
		Warn(msg)
	case code >= 300:
		Debug(msg)
	default:
		Info(msg)
	}
}
