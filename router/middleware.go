package router

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/opensaucerer/barf/server"
	"github.com/opensaucerer/barf/typing"
)

// Router routes requests to the correct handler
func Router(respond func(w http.ResponseWriter, status bool, statusCode int, message string, data map[string]interface{})) func(next http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if !server.Loaded(w) {
				w = server.Load(w)
			}

			// clone the handler
			sh := h

			// get route function
			route := Route{
				Path:    Path(r.URL),
				Method:  strings.ToLower(r.Method),
				Handler: nil,
				Params:  map[string]string{},
			}

			// check if route exists
			if !route.Exists() {
				respond(w, false, http.StatusNotFound, fmt.Sprintf("Path /%s for method %s not found", regexp.MustCompile("^/+|/+$").ReplaceAllString(route.Path, ""), strings.ToUpper(route.Method)), nil)
			} else {
				// load params into context if any
				ctx := context.WithValue(r.Context(), typing.ParamsCtxKey{}, route.Params)

				// call middleware stack if route is registered on a SubRoute
				if route.RetroFrame {
					s := route.Reframe()
					if s != nil {
						for i := range s.stack {
							sh = s.stack[len(s.stack)-1-i](sh)
						}
					}
				}

				// call route specific middleware(s)
				for i := range route.stack {
					sh = route.stack[len(route.stack)-1-i](sh)
				}

				// call route handler
				route.Handler(w, r.WithContext(ctx))
			}
			// call the next middleware
			sh.ServeHTTP(w, r)
		})
	}
}
