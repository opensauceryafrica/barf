package router

import "net/http"

// Post registers a route with the POST HTTP method
func Post(path string, handler func(http.ResponseWriter, *http.Request)) {
	route := &Route{
		path,
		post,
		handler,
	}
	route.Register()
}
