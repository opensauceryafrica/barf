package router

import "net/http"

// Post registers a route with the POST HTTP method
func Post(path string, handler func(http.ResponseWriter, *http.Request)) {
	route := &Route{
		Path:    path,
		Method:  post,
		Handler: handler,
	}
	route.Register()
}
