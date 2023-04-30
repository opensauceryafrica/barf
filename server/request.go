package server

import (
	"net/http"

	"github.com/opensaucerer/barf/server/body"
	"github.com/opensaucerer/barf/server/param"
	"github.com/opensaucerer/barf/server/query"
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
