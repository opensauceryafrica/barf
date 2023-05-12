package router

import (
	"github.com/opensaucerer/barf/cors"
	"github.com/opensaucerer/barf/server"

	"github.com/opensaucerer/barf/typing"
)

type hippocampus struct {
	router Frame
	stack  []typing.Middleware
}

/*
Hippocampus prepares the given barf router or base barf handler for hijacking. To take over the base barf handler, omit the router argument.

Note: the base barf handler is the one that is created by the barf.Stark() function and can only be hijacked before the barf.Beck() function is called.
*/
func Hippocampus(r ...Frame) *hippocampus {
	h := &hippocampus{}
	if len(r) > 0 {
		if _, ok := r[0].(*SubRoute); ok {
			h.router = r[0].(*SubRoute)
			h.stack = r[0].(*SubRoute).stack
		}
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
		if server.Beckoned == nil || !*server.Beckoned {
			// create a copy of Barf.Router
			r := Barf.Router
			for i := range h.stack {
				r = h.stack[len(h.stack)-1-i](r)
			}
			// add cors middleware such that it is called first before any user-defined middleware
			r = cors.CORS(cors.Prepare(*server.Augment.CORS))(r)
			// add recovery middleware
			if server.Augment.Recovery != nil && *server.Augment.Recovery {
				r = server.Recover(server.JSON)(r)
			}
			server.HTTP.Handler = r
		}
	} else {
		h.router.(*SubRoute).stack = h.stack
	}
}
