package router

import "net/http"

type Route struct {
	Path    string
	Method  string
	Handler func(http.ResponseWriter, *http.Request)
	Query   map[string]string
	Params  map[string]string
}

type Router struct {
	Routes []*Route
	Entry  string
}
