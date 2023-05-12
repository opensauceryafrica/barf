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

// fpost registers a route with the POST HTTP method and sets the RetroFrame flag
func fpost(path string, handler func(http.ResponseWriter, *http.Request), entry string) {
	route := &Route{
		Path:            path,
		Method:          post,
		Handler:         handler,
		RetroFrame:      true,
		RetroFrameEntry: entry,
	}
	route.Register()
}
