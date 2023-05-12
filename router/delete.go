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

// fdelete registers a route with the DELETE HTTP method and sets the RetroFrame flag
func fdelete(path string, handler func(http.ResponseWriter, *http.Request), entry string) {
	route := &Route{
		Path:            path,
		Method:          delete,
		Handler:         handler,
		RetroFrame:      true,
		RetroFrameEntry: entry,
	}
	route.Register()
}
