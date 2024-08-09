package router

import (
	"net/http"

	"github.com/opensaucerer/barf/typing"
)

// Delete registers a route with the DELETE HTTP method
func Delete(path string, handler func(http.ResponseWriter, *http.Request), m ...typing.Middleware) {
	route := &Route{
		Path:    path,
		Method:  delete,
		Handler: handler,
		stack:   m,
	}
	route.Register()
}

// fdelete registers a route with the DELETE HTTP method and sets the RetroFrame flag
func fdelete(path string, handler func(http.ResponseWriter, *http.Request), entry string, m ...typing.Middleware) {
	route := &Route{
		Path:            path,
		Method:          delete,
		Handler:         handler,
		RetroFrame:      true,
		RetroFrameEntry: entry,
		stack:           m,
	}
	route.Register()
}
