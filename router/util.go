package router

import (
	"net/url"
	"strings"
)

// Path returns the path of the URL
func Path(u *url.URL) string {
	return u.Path
}

// Params returns a map of the path parameters
func Params(path, route string) map[string]string {
	params := map[string]string{}
	pathParts := strings.Split(path, "/")
	routeParts := strings.Split(route, "/")
	if len(pathParts) == len(routeParts) {
		for i, part := range routeParts {
			if strings.HasPrefix(part, ":") {
				params[part[1:]] = pathParts[i]
			}
		}
	}
	return params
}
