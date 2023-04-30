package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/opensaucerer/barf/router"
	"github.com/opensaucerer/barf/server"
)

// Router routes requests to the correct handler
func Router(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get route function
		route := router.Route{
			Path:    server.Path(r.URL),
			Method:  strings.ToLower(r.Method),
			Handler: nil,
		}
		// check if route exists
		if !route.Exists() {
			server.JSON(w, false, http.StatusNotFound, fmt.Sprintf("Path /%s for method %s not found", route.Path, strings.ToUpper(route.Method)), nil)
		} else {
			// call route handler
			route.Handler(w, r)
		}
		// call the next middleware
		next.ServeHTTP(w, r)
	})
}
