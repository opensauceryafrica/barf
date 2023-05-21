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
		entry:  r.entry + "/" + regexp.MustCompile("^/+|/+$").ReplaceAllString(path, ""),
		routes: []*Route{},
		stack:  []typing.Middleware{},
	}
	s.key = fmt.Sprintf("%p", s)
	stable[s.key] = s
	return s
}

// Get registers a route with the GET HTTP method against the path of the RetroFrame router instance.
func (r *SubRoute) Get(path string, handler func(http.ResponseWriter, *http.Request)) {
	fget(fmt.Sprintf("/%s/%s", r.entry, regexp.MustCompile("^/+|/+$").ReplaceAllString(path, "")), handler, r.key)
}

// Post registers a route with the POST HTTP method against the path of the RetroFrame router instance.
func (r *SubRoute) Post(path string, handler func(http.ResponseWriter, *http.Request)) {
	fpost(fmt.Sprintf("/%s/%s", r.entry, regexp.MustCompile("^/+|/+$").ReplaceAllString(path, "")), handler, r.key)
}

// Put registers a route with the PUT HTTP method against the path of the RetroFrame router instance.
func (r *SubRoute) Put(path string, handler func(http.ResponseWriter, *http.Request)) {
	fput(fmt.Sprintf("/%s/%s", r.entry, regexp.MustCompile("^/+|/+$").ReplaceAllString(path, "")), handler, r.key)
}

// Patch registers a route with the PATCH HTTP method against the path of the RetroFrame router instance.
func (r *SubRoute) Patch(path string, handler func(http.ResponseWriter, *http.Request)) {
	fpatch(fmt.Sprintf("/%s/%s", r.entry, regexp.MustCompile("^/+|/+$").ReplaceAllString(path, "")), handler, r.key)
}

// Delete registers a route with the DELETE HTTP method against the path of the RetroFrame router instance.
func (r *SubRoute) Delete(path string, handler func(http.ResponseWriter, *http.Request)) {
	fdelete(fmt.Sprintf("/%s/%s", r.entry, regexp.MustCompile("^/+|/+$").ReplaceAllString(path, "")), handler, r.key)
}

// Any registers a route with the all HTTP method against the path of the RetroFrame router instance.
func (r *SubRoute) Any(path string, handler func(http.ResponseWriter, *http.Request)) {
	fany(fmt.Sprintf("/%s/%s", r.entry, regexp.MustCompile("^/+|/+$").ReplaceAllString(path, "")), handler, r.key)
}
