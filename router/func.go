package router

import (
	"regexp"
	"strings"
)

// Func retrieves the handler function for the given path and method
func (r *Route) Func() {
	// remove preceding and trailing slashes
	r.Path = regexp.MustCompile("^/+|/+$").ReplaceAllString(r.Path, "")
	if r.Path == "" {
		r.Path = "/"
	}
	// if path found in top level of table
	if table[r.Path] != nil && table[r.Path][r.Method] != nil {
		r.Handler = table[r.Path][r.Method]
		return
	}

	// handle path with parameters
	paths := strings.Split(r.Path, "/")
TLoop:
	for path, methods := range table {
		variables := strings.Split(path, "/")
		if len(paths) == len(variables) && methods[r.Method] != nil {
			match := true
		VLoop:
			for i, variable := range variables {
				if variable != paths[i] && variable[0] != ':' {
					match = false
					break VLoop
				}
			}
			if match {
				r.Handler = methods[r.Method]
				break TLoop
			}
		}
	}
}
