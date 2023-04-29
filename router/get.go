package router

import "net/http"

// Get registers a route with the GET HTTP method
func Get(path string, handler func(http.ResponseWriter, *http.Request)) {
	route := &Route{
		path,
		get,
		handler,
	}
	route.Register()
}
