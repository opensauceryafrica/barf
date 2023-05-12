package barf

import (
	"github.com/opensaucerer/barf/router"
)

/*
Hippocampus prepares the given barf router or base barf handler for hijacking. To take over the base barf handler, omit the router argument.

Note: the base barf handler is the one that is created by the barf.Stark() function and can only be hijacked before the barf.Beck() function is called.
*/
var Hippocampus = router.Hippocampus
