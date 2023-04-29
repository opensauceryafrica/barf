package router

import (
	"net/http"
	"regexp"
)

// Register registers a route in the router table
func (r *Route) Register() {
	// remove preceding and trailing slashes
	r.path = regexp.MustCompile("^/+|/+$").ReplaceAllString(r.path, "")
	if r.path == "" {
		r.path = "/"
	}
	if table[r.path] == nil {
		table[r.path] = make(map[string]func(http.ResponseWriter, *http.Request))
	}
	table[r.path][r.method] = r.handler
}
