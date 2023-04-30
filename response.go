package barf

import (
	"github.com/opensaucerer/barf/server"
	"github.com/opensaucerer/barf/typing"
)

// Status prepares a response with the given writer and status code
var Status = server.Status

// Res is a simple struct for a status based response
type Res = typing.Response
