package logger

import (
	"net/http"
	"strconv"
	"time"

	"github.com/opensaucerer/barf/server"
	"github.com/opensaucerer/barf/typing"
)

// Morgan logs the request to the console
func Morgan(rr typing.RequestResponse) {

	// format: utc timestamp: user-agent - http/version: method - path - status code - status text
	msg := time.Now().UTC().Format(time.RFC3339) + ": " + rr.Request.UserAgent() + " - " + rr.Request.Proto + ": " + rr.Request.Method + " - " + rr.Request.URL.Path + " - " + strconv.Itoa(rr.Code) + " - " + http.StatusText(rr.Code)
	// log request
	Code(msg, rr.Code)

	// call next middleware? (the logger is the last middleware in the stack)
	// next.ServeHTTP(w, r)

}

// Winston listens on a channel for log messages and calls morgan for each message
func Winston() {
	for msg := range server.RequestResponseChan {
		Morgan(msg)
	}
}
