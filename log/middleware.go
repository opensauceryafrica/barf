package logger

import (
	"net/http"
	"strconv"
	"time"

	"github.com/opensaucerer/barf/server"
)

// Morgan logs the request to the console
func Morgan(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// get status code
		code := server.Status(w)
		// format: utc timestamp: user-agent - http/version: method - path - status code - status text
		msg := time.Now().UTC().Format(time.RFC3339) + ": " + r.UserAgent() + " - " + r.Proto + ": " + r.Method + " - " + r.URL.Path + " - " + strconv.Itoa(code) + " - " + http.StatusText(code)
		// log request
		Code(msg, code)

		// call next middleware? (the logger is the last middleware in the stack)
		// next.ServeHTTP(w, r)
	})
}
