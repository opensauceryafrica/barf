package server

import (
	"errors"
	"net/http"
)

// Recover is a middleware that recovers from panics and sends a 500 response.
func Recover(response func(w http.ResponseWriter, status bool, statusCode int, message string, data map[string]interface{})) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// overload the writer - do not remove this unless you know what you are doing
			w = Load(w)

			defer func() {
				rr := recover()
				if rr != nil {
					var err error
					switch t := rr.(type) {
					case string:
						err = errors.New(t)
					case error:
						err = t
					default:
						err = errors.New("unknown error")
					}

					// @TODO: prevent write into response until Recover has reported error or nil, if none.
					if Loaded(w) && Written(w) {
						return
					}

					response(w, false, http.StatusInternalServerError, "Internal Server Error: "+err.Error(), nil)
					return
				}
			}()
			h.ServeHTTP(w, r)
		})
	}
}
