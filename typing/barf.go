package typing

import "net/http"

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
	// Host is the host for the server to listen on
	Host string
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
	// Logging is for defining whether or not to enable request logging
	// default is true
	Logging *bool
	// Recovery is for defining whether or not to enable panic recovery
	// default is true
	Recovery *bool
	// CORS is the configuration for Cross-Origin Resource Sharing
	CORS *CORS
	// UseHTTPS specifies that all barf connections listen on the TCP network address for inbound HTTPS requests.
	UseHTTPS bool
	// SSLCertFile is the path to the certificate file.
	//
	// If the certificate is signed by a certificate authority, the SSLCertfFile should be the concatenation of the server's certificate, any intermediates, and the CA's certificate.
	SSLCertFile string
	// SSLKeyFile is the path to the private key of the certificate in use.
	SSLKeyFile string
}

// CORS holds configuration for Cross-Origin Resource Sharing
type CORS struct {
	// AllowedOrigins is a list of origins a cross-domain request can be executed from
	AllowedOrigins []string
	// AllowMethods is a list of methods the client is allowed to use with
	// cross-domain requests
	AllowedMethods []string
	// AllowHeaders is list of non simple headers the client is allowed to use with
	// cross-domain requests
	AllowedHeaders []string
	// ExposeHeaders indicates which headers are safe to expose to the API of a CORS
	// API specification
	ExposedHeaders []string
	// AllowCredentials indicates whether or not the response to the request can be exposed
	// when the credentials flag is true. When used as part of a response to a preflight
	// request, this indicates whether or not the actual request can be made using credentials.
	AllowCredentials bool
	// MaxAge indicates how long (in seconds) the results of a preflight request
	// can be cached
	MaxAge int
	// OptionsPassthrough instructs preflight to let other potential next handlers to process the OPTIONS method
	OptionsPassthrough bool
	// OptionsSuccessStatus sets the statusCode for OPTIONS requests, defaults to 204
	OptionsSuccessStatus int
	// AllowedOriginFunc is a callback for handling user defined origin checks.
	AllowedOriginFunc func(origin string) bool
	// AllowedOriginWithRequestFunc is a callback for handling user defined origin checks with access to the http request object.
	AllowedOriginWithRequestFunc func(origin string, r *http.Request) bool
}
