package server

import (
	"github.com/opensaucerer/barf/router"
	"github.com/opensaucerer/barf/typing"
)

type hippocampus struct {
	Router router.Hippocampus
	Stack  []typing.Middleware
}

/*
Hippocampus prepares the given barf router or base barf handler for hijacking. To take over the base barf handler, omit the router argument.

Note: the base barf handler is the one that is created by the barf.Stark() function and can only be hijacked before the barf.Stark() function is called.
*/
func Hippocampus(r ...router.Hippocampus) *hippocampus {
	h := &hippocampus{}
	if len(r) > 0 {
		h.Router = r[0].(router.Router)
		h.Stack = r[0].(router.Router).Stack
	} else {
		h.Router = Barf.Router
		h.Stack = Barf.Stack
	}
	return h
}

// Hijack takes over the given barf router or base barf handler by injecting the given middleware.
func (h *hippocampus) Hijack(middleware ...typing.Middleware) {
	if len(middleware) > 0 {
		h.Stack = append(h.Stack, middleware...)
	}
	// hijack base barf handler
	if h.Router == nil {
		Barf.Stack = h.Stack
	}

	// TODO: implement hijacking of router
}
