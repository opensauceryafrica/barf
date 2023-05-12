package router

import "net/http"

// Any registers a route with all HTTP method
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

// fany registers a route with all HTTP method and sets the RetroFrame flag
func fany(path string, handler func(http.ResponseWriter, *http.Request), entry string) {
	for _, method := range methods {
		route := &Route{
			Path:    path,
			Method:  method,
			Handler: handler,
		}
		route.Register()
	}
}
