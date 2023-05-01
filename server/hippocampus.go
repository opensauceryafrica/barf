package server

import (
	"github.com/opensaucerer/barf/middleware"
	"github.com/opensaucerer/barf/router"
	"github.com/opensaucerer/barf/typing"
)

type hippocampus struct {
	router router.Hippocampus
	stack  []typing.Middleware
}

/*
Hippocampus prepares the given barf router or base barf handler for hijacking. To take over the base barf handler, omit the router argument.

Note: the base barf handler is the one that is created by the barf.Stark() function and can only be hijacked before the barf.Beck() function is called.
*/
func Hippocampus(r ...router.Hippocampus) *hippocampus {
	h := &hippocampus{}
	if len(r) > 0 {
		h.router = r[0].(*router.Router)
		h.stack = r[0].(*router.Router).Stack
	} else {
		if Barf.Router != nil {
			h.stack = Barf.Stack
		}
	}
	return h
}

// Hijack takes over the given barf router or base barf handler by injecting the given middleware.
func (h *hippocampus) Hijack(m ...typing.Middleware) {
	if len(m) > 0 {
		if h.stack != nil {
			h.stack = append(h.stack, m...)
		}
	}
	// hijack base barf handler
	if h.router == nil {

		Barf.Stack = h.stack

		// hijacking the base barf handler is only possible before the barf.Beck() function is called
		if Beckoned == nil || !*Beckoned {
			// create a copy of Barf.Router
			r := Barf.Router
			for i := range h.stack {
				r = h.stack[len(h.stack)-1-i](r)
			}
			// add cors middleware such that it is called first before any user-defined middleware
			r = middleware.CORS(middleware.Prepare(*Augment.CORS))(r)
			// add recovery middleware
			if Augment.Recovery != nil && *Augment.Recovery {
				r = middleware.Recover(JSON)(r)
			}
			HTTP.Handler = r
		}
	} else {
		h.router.(*router.Router).Stack = h.stack
		// TODO: implement hijacking of router
	}
}
