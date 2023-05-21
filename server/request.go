package server

import (
	"net/http"

	"github.com/opensaucerer/barf/router/body"
	"github.com/opensaucerer/barf/router/form"
	"github.com/opensaucerer/barf/router/param"
	"github.com/opensaucerer/barf/router/query"
)

type request struct {
	request *http.Request
}

// Request prepares a barf request with the given http request
func Request(r *http.Request) *request {
	return &request{
		request: r,
	}
}

// Form prepares the barf request with the request form or multipart/form-data for further formatting.
//
// It takes an optional maxMemory parameter as the maximum memory to be used by the parser in bytes and defaults to 32MB (32 << 20) for file parts and 10MB (10 << 20) for non-file parts.
func (r *request) Form(maxMemory ...int64) form.M {
	if len(maxMemory) == 0 {
		maxMemory = []int64{32 << 20}
	}
	return form.Form(r.request, maxMemory[0])
}

// Body prepares the barf request with the request body for further formatting
func (r *request) Body() body.B {
	return body.Body(r.request)
}

// Params prepares the barf request with the request params for further formatting
func (r *request) Params() param.P {
	return param.Params(r.request)
}

// Query prepares the barf request with the request query for further formatting
func (r *request) Query() query.Q {
	return query.Query(r.request)
}
