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

// fput registers a route with the PUT HTTP method and sets the RetroFrame flag
func fput(path string, handler func(http.ResponseWriter, *http.Request), entry string) {
	route := &Route{
		Path:            path,
		Method:          put,
		Handler:         handler,
		RetroFrame:      true,
		RetroFrameEntry: entry,
	}
	route.Register()
}
