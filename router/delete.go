package router

import "net/http"

// Delete registers a route with the DELETE HTTP method
func Delete(path string, handler func(http.ResponseWriter, *http.Request)) {
	route := &Route{
		Path:    path,
		Method:  delete,
		Handler: handler,
	}
	route.Register()
}
