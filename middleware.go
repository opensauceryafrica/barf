package barf

import (
	"github.com/opensaucerer/barf/router"
	"github.com/opensaucerer/barf/typing"
)

/*
Hippocampus prepares the given barf router or base barf handler for hijacking. To take over the base barf handler, omit the router argument.

Note: the base barf handler is the one that is created by the barf.Stark() function and can only be hijacked before the barf.Beck() function is called.
*/
var Hippocampus = router.Hippocampus

// Middleware is the type of a function that can be used as a barf middleware.
type Middleware = typing.Middleware
