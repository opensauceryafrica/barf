package router

import "net/http"

// Any registers a route with the all HTTP method
func Any(path string, handler func(http.ResponseWriter, *http.Request)) {
	for _, method := range methods {
		route := &Route{
			Path:    path,
			Method:  method,
			Handler: handler,
		}
		route.Register()
	}
}
