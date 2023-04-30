package router

import "net/http"

// Put registers a route with the PUT HTTP method
func Put(path string, handler func(http.ResponseWriter, *http.Request)) {
	route := &Route{
		Path:    path,
		Method:  put,
		Handler: handler,
	}
	route.Register()
}
