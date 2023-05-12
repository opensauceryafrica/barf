package router

import (
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
}
