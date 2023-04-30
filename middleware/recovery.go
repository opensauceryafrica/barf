package middleware

import (
	"errors"
	"net/http"

	"github.com/opensaucerer/barf/server"
)

// Recover is a middleware that recovers from panics and sends a 500 response.
func Recover(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if server.Augment.Recovery {
			defer func() {
				r := recover()
				if r != nil {
					var err error
					switch t := r.(type) {
					case string:
						err = errors.New(t)
					case error:
						err = t
					default:
						err = errors.New("unknown error")
					}
					server.JSON(w, false, http.StatusInternalServerError, "Internal Server Error: "+err.Error(), nil)
				}
			}()
		}
		h.ServeHTTP(w, r)
	})
}
