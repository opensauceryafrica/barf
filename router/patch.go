package router

import "net/http"

// Patch registers a route with the PATCH HTTP method
func Patch(path string, handler func(http.ResponseWriter, *http.Request)) {
	route := &Route{
		Path:    path,
		Method:  patch,
		Handler: handler,
	}
	route.Register()
}

// fpatch registers a route with the PATCH HTTP method and sets the RetroFrame flag
func fpatch(path string, handler func(http.ResponseWriter, *http.Request), entry string) {
	route := &Route{
		Path:            path,
		Method:          patch,
		Handler:         handler,
		RetroFrame:      true,
		RetroFrameEntry: entry,
	}
	route.Register()
}
