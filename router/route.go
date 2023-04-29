package router

import "net/http"

type Route struct {
	path    string
	method  string
	handler func(http.ResponseWriter, *http.Request)
}
