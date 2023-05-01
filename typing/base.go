package typing

import "net/http"

type Health struct {
	Version     string `json:"version"`
	Status      bool   `json:"status"`
	Description string `json:"description"`
}

type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type M map[string]string

type Middleware func(http.Handler) http.Handler
