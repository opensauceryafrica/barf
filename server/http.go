package server

import (
	"net/http"

	"github.com/opensaucerer/barf/typing"
)

var (
	HTTP *http.Server

	Mux *http.ServeMux

	Augment *typing.Augment

	Beckoned *bool
)
