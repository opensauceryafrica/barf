package server

import (
	"net/http"

	"github.com/opensaucerer/barf/router"
	"github.com/opensaucerer/barf/typing"
)

var (
	HTTP *http.Server

	Mux *http.ServeMux

	Augment *typing.Augment

	Beckoned *bool

	Barf *(struct {
		Router router.Hippocampus
		Stack  []typing.Middleware
	}) = &struct {
		Router router.Hippocampus
		Stack  []typing.Middleware
	}{}
)
