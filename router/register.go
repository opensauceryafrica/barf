package router

import (
	"net/http"
	"regexp"
)

// Register registers a route in the router table
func (r *Route) Register() {
	// remove preceding and trailing slashes
	r.Path = regexp.MustCompile("^/+|/+$").ReplaceAllString(r.Path, "")
	if r.Path == "" {
		r.Path = "/"
	}
	if rtable[r.Path] == nil {
		rtable[r.Path] = make(map[string]*Route)
	}
	rtable[r.Path][r.Method] = r
	// the handler is added as just another middleware
	r.stack = append(r.stack, func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			r.Handler(w, req)
			h.ServeHTTP(w, req)
		})
	})
}
