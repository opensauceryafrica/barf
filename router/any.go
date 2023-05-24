package router

import (
	"net/http"

	"github.com/opensaucerer/barf/typing"
)

// Any registers a route with all HTTP method
func Any(path string, handler func(http.ResponseWriter, *http.Request), m ...typing.Middleware) {
	for _, method := range methods {
		route := &Route{
			Path:    path,
			Method:  method,
			Handler: handler,
			stack:   m,
		}
		route.Register()
	}
}

// fany registers a route with all HTTP method and sets the RetroFrame flag
func fany(path string, handler func(http.ResponseWriter, *http.Request), entry string, m ...typing.Middleware) {
	for _, method := range methods {
		route := &Route{
			Path:    path,
			Method:  method,
			Handler: handler,
			stack:   m,
		}
		route.Register()
	}
}
