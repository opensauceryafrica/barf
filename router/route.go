package router

import (
	"net/http"

	"github.com/opensaucerer/barf/typing"
)

type Route struct {
	Path    string
	Method  string
	Handler func(http.ResponseWriter, *http.Request)
	Query   map[string]string
	Params  map[string]string
}

type Router struct {
	Entry  string
	Routes []*Route
	Stack  []typing.Middleware
}

func (r Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
}

type Hippocampus interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
