package router

import "net/http"

// Patch registers a route with the PATCH HTTP method
func Patch(path string, handler func(http.ResponseWriter, *http.Request)) {
	route := &Route{
		Path:    path,
		Method:  patch,
		Handler: handler,
	}
	route.Register()
}
