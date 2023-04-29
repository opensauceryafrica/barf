package router

import "net/http"

// Put registers a route with the PUT HTTP method
func Put(path string, handler func(http.ResponseWriter, *http.Request)) {
	route := &Route{
		path,
		put,
		handler,
	}
	route.Register()
}
