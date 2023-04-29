package router

import "net/http"

// All registers a route with the all HTTP method
func All(path string, handler func(http.ResponseWriter, *http.Request)) {
	for _, method := range methods {
		route := &Route{
			path,
			method,
			handler,
		}
		route.Register()
	}
}
