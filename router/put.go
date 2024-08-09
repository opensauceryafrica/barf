package router

import (
	"net/http"

	"github.com/opensaucerer/barf/typing"
)

// Put registers a route with the PUT HTTP method
func Put(path string, handler func(http.ResponseWriter, *http.Request), m ...typing.Middleware) {
	route := &Route{
		Path:    path,
		Method:  put,
		Handler: handler,
		stack:   m,
	}
	route.Register()
}

// fput registers a route with the PUT HTTP method and sets the RetroFrame flag
func fput(path string, handler func(http.ResponseWriter, *http.Request), entry string, m ...typing.Middleware) {
	route := &Route{
		Path:            path,
		Method:          put,
		Handler:         handler,
		RetroFrame:      true,
		RetroFrameEntry: entry,
		stack:           m,
	}
	route.Register()
}
