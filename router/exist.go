package router

// Exists returns true if a route with the given path and method exists
func (r *Route) Exists() bool {
	r.Func()
	return r.Handler != nil
}
