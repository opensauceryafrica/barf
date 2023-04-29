package logger

import (
	"net/http"
	"strconv"
	"time"
	// "github.com/opensaucerer/barf/typing"
)

// Request logs the request to the console
func Request(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: properly handle status code
		// create an addressable copy of the response writer
		// aw := reflect.New(reflect.TypeOf(w)).Interface().(http.ResponseWriter)
		// statusCode := r.Context().Value(typing.StatusCodeCtxKey{}).(int)
		statusCode := 200
		// format: iso timestamp: user-agent - http/version: method - path - status code - status text
		msg := time.Now().Format(time.RFC3339) + ": " + r.UserAgent() + " - " + r.Proto + ": " + r.Method + " - " + r.URL.Path + " - " + strconv.Itoa(statusCode) + " - " + http.StatusText(statusCode)
		// log request
		Code(msg, statusCode)
		// call next middleware
		next.ServeHTTP(w, r)
	})
}
