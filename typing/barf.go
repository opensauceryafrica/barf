package typing

// Augment holds refrence to all of barf's config
type Augment struct {
	// MaxHeaderBytes is the maximum number of bytes the server will
	// read parsing the request header's keys and values, including the
	// request line. It does not limit the size of the request body.
	// default is 1 << 20 (1 MB)
	MaxHeaderBytes int
	// ReadTimeout is the maximum duration in seconds for reading the entire
	// request, including the body.
	// default is 10 seconds
	ReadTimeout int
	// WriteTimeout is the maximum duration in seconds before timing out
	// writes of the response. It is reset whenever a new
	// request's header is read. Like ReadTimeout, it does not
	// let Handlers make decisions on a per-request basis.
	// default is 10 seconds
	WriteTimeout int
	// ShutdownTimeout is the time in seconds to wait for the server to shutdown gracefully
	// default is 5 seconds
	ShutdownTimeout int
	// Port is the port for the server to listen on
	Port string
	// ReadHeaderTimeout is the amount of time allowed to read
	// request headers. The connection's read deadline is reset
	// after reading the headers and the Handler can decide what
	// is considered too slow for the body. If ReadHeaderTimeout
	// is zero, the value of ReadTimeout is used. If both are
	// zero, there is no timeout.
	// default is 10 seconds
	ReadHeaderTimeout int
}
