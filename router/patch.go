package router

import "net/http"

// Patch registers a route with the PATCH HTTP method
func Patch(path string, handler func(http.ResponseWriter, *http.Request)) {
	route := &Route{
		path,
		patch,
		handler,
	}
	route.Register()
}
