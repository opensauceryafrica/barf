package router

import (
	"net/http"

	"github.com/opensaucerer/barf/typing"
)

type Route struct {
	Path            string
	Method          string
	Handler         func(http.ResponseWriter, *http.Request)
	Query           map[string]string
	Params          map[string]string
	RetroFrame      bool   // true for routes registered on a SubRoute
	RetroFrameEntry string // entry path of the SubRoute
	stack           []typing.Middleware
}

type SubRoute struct {
	entry  string
	routes []*Route
	stack  []typing.Middleware
	key    string
}

func (r SubRoute) ServeHTTP(w http.ResponseWriter, req *http.Request) {
}

type Frame interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

var (
	Barf *(struct {
		Router Frame
		Stack  []typing.Middleware
	}) = &struct {
		Router Frame
		Stack  []typing.Middleware
	}{}
)
