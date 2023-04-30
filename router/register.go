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
	if table[r.Path] == nil {
		table[r.Path] = make(map[string]func(http.ResponseWriter, *http.Request))
	}
	table[r.Path][r.Method] = r.Handler
}
