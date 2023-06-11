package router

import (
	"net/http"

	"github.com/opensaucerer/barf/typing"
)

// Post registers a route with the POST HTTP method
func Post(path string, handler func(http.ResponseWriter, *http.Request), m ...typing.Middleware) {
	route := &Route{
		Path:    path,
		Method:  post,
		Handler: handler,
		stack:   m,
	}
	route.Register()
}

// fpost registers a route with the POST HTTP method and sets the RetroFrame flag
func fpost(path string, handler func(http.ResponseWriter, *http.Request), entry string, m ...typing.Middleware) {
	route := &Route{
		Path:            path,
		Method:          post,
		Handler:         handler,
		RetroFrame:      true,
		RetroFrameEntry: entry,
		stack:           m,
	}
	route.Register()
}
