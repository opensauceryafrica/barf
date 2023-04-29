package router

import (
	"regexp"
	"strings"
)

// Func retrieves the handler function for the given path and method
func (r *Route) Func() {
	// remove preceding and trailing slashes
	r.path = regexp.MustCompile("^/+|/+$").ReplaceAllString(r.path, "")
	if r.path == "" {
		r.path = "/"
	}
	// if path found in top level of table
	if table[r.path] != nil && table[r.path][r.method] != nil {
		r.handler = table[r.path][r.method]
		return
	}

	// handle path with parameters
	paths := strings.Split(r.path, "/")
TLoop:
	for path, methods := range table {
		variables := strings.Split(path, "/")
		if len(paths) == len(variables) && methods[r.method] != nil {
			match := true
		VLoop:
			for i, variable := range variables {
				if variable != paths[i] && variable[0] != ':' {
					match = false
					break VLoop
				}
			}
			if match {
				r.handler = methods[r.method]
				break TLoop
			}
		}
	}
}
