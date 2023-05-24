package router

import (
	"net/http"

	"github.com/opensaucerer/barf/typing"
)

// Patch registers a route with the PATCH HTTP method
func Patch(path string, handler func(http.ResponseWriter, *http.Request), m ...typing.Middleware) {
	route := &Route{
		Path:    path,
		Method:  patch,
		Handler: handler,
		stack:   m,
	}
	route.Register()
}

// fpatch registers a route with the PATCH HTTP method and sets the RetroFrame flag
func fpatch(path string, handler func(http.ResponseWriter, *http.Request), entry string, m ...typing.Middleware) {
	route := &Route{
		Path:            path,
		Method:          patch,
		Handler:         handler,
		RetroFrame:      true,
		RetroFrameEntry: entry,
		stack:           m,
	}
	route.Register()
}
