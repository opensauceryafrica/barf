package barf

import "github.com/opensaucerer/barf/router"

// Get registers a route with the GET HTTP method
var Get = router.Get

// Post registers a route with the POST HTTP method
var Post = router.Post

// Put registers a route with the PUT HTTP method
var Put = router.Put

// Patch registers a route with the PATCH HTTP method
var Patch = router.Patch

// Delete registers a route with the DELETE HTTP method
var Delete = router.Delete

// // Head registers a route with the HEAD HTTP method
// var Head = router.Head

// // Options registers a route with the OPTIONS HTTP method
// var Options = router.Options

// Any registers a route with all HTTP methods
var Any = router.Any

// RetroFrame returns a new RetroFrame instance registered against the given entry path.
var RetroFrame = router.RetroFrame

// SubRoute is the type of a barf RetroFrame instance.
type SubRoute = router.SubRoute
