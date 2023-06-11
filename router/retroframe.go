package router

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/opensaucerer/barf/typing"
)

// RetroFrame returns a new subrouter instance registered against the given entry path.
func RetroFrame(path string) *SubRoute {
	s := &SubRoute{
		entry:  regexp.MustCompile("^/+|/+$").ReplaceAllString(path, ""),
		routes: []*Route{},
		stack:  []typing.Middleware{},
	}
	s.key = fmt.Sprintf("%p", s)
	stable[s.key] = s
	return s
}

// RetroFrame returns a new subrouter instance registered against the given entry path and the RetroFrame instance it is called on.
func (r *SubRoute) RetroFrame(path string) *SubRoute {
	s := &SubRoute{
		entry:  regexp.MustCompile("^/+|/+$").ReplaceAllString(r.entry+"/"+regexp.MustCompile("^/+|/+$").ReplaceAllString(path, ""), ""),
		routes: []*Route{},
		stack:  append([]typing.Middleware{}, r.stack...),
	}
	s.key = fmt.Sprintf("%p", s)
	stable[s.key] = s
	return s
}

// Get registers a route with the GET HTTP method against the path of the RetroFrame router instance.
func (r *SubRoute) Get(path string, handler func(http.ResponseWriter, *http.Request)) *SubRoute {
	fget(fmt.Sprintf("/%s/%s", r.entry, regexp.MustCompile("^/+|/+$").ReplaceAllString(path, "")), handler, r.key)
	return r
}

// Post registers a route with the POST HTTP method against the path of the RetroFrame router instance.
func (r *SubRoute) Post(path string, handler func(http.ResponseWriter, *http.Request)) *SubRoute {
	fpost(fmt.Sprintf("/%s/%s", r.entry, regexp.MustCompile("^/+|/+$").ReplaceAllString(path, "")), handler, r.key)
	return r
}

// Put registers a route with the PUT HTTP method against the path of the RetroFrame router instance.
func (r *SubRoute) Put(path string, handler func(http.ResponseWriter, *http.Request)) *SubRoute {
	fput(fmt.Sprintf("/%s/%s", r.entry, regexp.MustCompile("^/+|/+$").ReplaceAllString(path, "")), handler, r.key)
	return r
}

// Patch registers a route with the PATCH HTTP method against the path of the RetroFrame router instance.
func (r *SubRoute) Patch(path string, handler func(http.ResponseWriter, *http.Request)) *SubRoute {
	fpatch(fmt.Sprintf("/%s/%s", r.entry, regexp.MustCompile("^/+|/+$").ReplaceAllString(path, "")), handler, r.key)
	return r
}

// Delete registers a route with the DELETE HTTP method against the path of the RetroFrame router instance.
func (r *SubRoute) Delete(path string, handler func(http.ResponseWriter, *http.Request)) *SubRoute {
	fdelete(fmt.Sprintf("/%s/%s", r.entry, regexp.MustCompile("^/+|/+$").ReplaceAllString(path, "")), handler, r.key)
	return r
}

// Any registers a route with the all HTTP method against the path of the RetroFrame router instance.
func (r *SubRoute) Any(path string, handler func(http.ResponseWriter, *http.Request)) *SubRoute {
	fany(fmt.Sprintf("/%s/%s", r.entry, regexp.MustCompile("^/+|/+$").ReplaceAllString(path, "")), handler, r.key)
	return r
}

// Use registers a middleware for a specific route and subroute
func (r *SubRoute) Use(handler ...typing.Middleware) *SubRoute {
	r.stack = append(r.stack, handler...)
	return r
}
