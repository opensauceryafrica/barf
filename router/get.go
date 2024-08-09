package router

import (
	"net/http"

	"github.com/opensaucerer/barf/typing"
)

// Get registers a route with the GET HTTP method
func Get(path string, handler func(http.ResponseWriter, *http.Request), m ...typing.Middleware) {
	route := &Route{
		Path:    path,
		Method:  get,
		Handler: handler,
		stack:   m,
	}
	route.Register()
}

// fget registers a route with the GET HTTP method and sets the RetroFrame flag
func fget(path string, handler func(http.ResponseWriter, *http.Request), entry string, m ...typing.Middleware) {
	route := &Route{
		Path:            path,
		Method:          get,
		Handler:         handler,
		RetroFrame:      true,
		RetroFrameEntry: entry,
		stack:           m,
	}
	route.Register()
}
