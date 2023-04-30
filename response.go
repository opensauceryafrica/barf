package barf

import (
	"github.com/opensaucerer/barf/server"
	"github.com/opensaucerer/barf/typing"
)

// Response prepares a barf response with the given writer
var Response = server.Response

// Res is a simple struct for a status based response
type Res = typing.Response
