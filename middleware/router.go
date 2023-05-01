package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/opensaucerer/barf/router"

	"github.com/opensaucerer/barf/typing"
)

// Router routes requests to the correct handler
func Router(respond func(w http.ResponseWriter, status bool, statusCode int, message string, data map[string]interface{})) func(next http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// get route function
			route := router.Route{
				Path:    router.Path(r.URL),
				Method:  strings.ToLower(r.Method),
				Handler: nil,
				Params:  map[string]string{},
			}
			// check if route exists
			if !route.Exists() {
				respond(w, false, http.StatusNotFound, fmt.Sprintf("Path /%s for method %s not found", route.Path, strings.ToUpper(route.Method)), nil)
			} else {
				// load params into context if any
				ctx := context.WithValue(r.Context(), typing.ParamsCtxKey{}, route.Params)

				// call route handler
				route.Handler(w, r.WithContext(ctx))
			}
			// call the next middleware
			h.ServeHTTP(w, r)
		})
	}
}
