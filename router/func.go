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
	if rtable[r.Path] != nil && rtable[r.Path][r.Method] != nil {
		*r = *rtable[r.Path][r.Method]
		return
	}

	// handle path with parameters
	paths := strings.Split(r.Path, "/")
TLoop:
	for path, methods := range rtable {
		variables := strings.Split(path, "/")
		if len(paths) == len(variables) && methods[r.Method] != nil {
			match := true
		VLoop:
			for i, variable := range variables {
				if variable != paths[i] && (len(variable) == 0 || variable[0] != ':') {
					match = false
					break VLoop
				}
			}
			if match {
				params := Params(r.Path, path)
				*r = *methods[r.Method]
				r.Params = params
				break TLoop
			}
		}
	}
}

// Reframe retrieves the SubRoute for the given route from the stable
func (r *Route) Reframe() *SubRoute {
	return stable[r.RetroFrameEntry]
}
