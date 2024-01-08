package typing

import (
	"fmt"
	"net/http"
)

// Augment holds reference to all of barf's config
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
	// AllowHotReload if true allows application to listen for file changes and restart the server.
	AllowHotReload *bool
	// HotReload contains configuration for server hot-re
	HotReload *HotReload
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

// RequestResponse holds the request and response objects
type RequestResponse struct {
	Request  *http.Request
	Response http.ResponseWriter
	Code     int
}

type HotReload struct {
	Root        string   // working dir, defaults to .
	IncludeExt  []string // file extensions to watch, defaults to all
	ExcludeExt  []string // file extensions to exclude
	ExcludeDir  []string // directories to exclude
	IncludeDir  []string // directories to include, defaults to all
	IncludeFile []string // files to include only
	Delay       uint     // delay reload if changes occur too frequently
	StopOnError bool     // if application should exit on error or ignore changes
	BuildCmd    string   // user specified build command
	Bin         string   // binary file path
	TmpDir      string   // temporary files directory
}

// GetRoot return the given root, defaults to .
func (h *HotReload) GetRoot() string {
	if h.Root == "" {
		return "."
	}
	return h.Root
}

// IsIncludeFile returns true if a given filename is in IncludeFile config
func (h *HotReload) IsIncludeFile(filename string) bool {
	if len(h.IncludeFile) < 1 {
		return false
	}
	for _, file := range h.IncludeFile {
		if file == filename {
			return true
		}
	}
	return false
}

// IsExcludeDir returns true if a given directory is in the ExcludeDir config
func (h *HotReload) IsExcludeDir(dir string) bool {
	if len(h.ExcludeDir) < 1 {
		return false
	}
	for _, d := range h.ExcludeDir {
		if d == dir {
			return true
		}
	}
	return false
}

// IsIncludeDir returns true if a given directory is in the IncludeDir config
// empty config returns true
func (h *HotReload) IsIncludeDir(dir string) bool {
	if len(h.IncludeDir) < 1 {
		return true
	}
	for _, d := range h.IncludeDir {
		if d == dir {
			return true
		}
	}
	return false
}

func (h *HotReload) GetBuildCmd() string {
	if h.BuildCmd == "" {
		return fmt.Sprintf("go build -o %s/main %s", h.GetTmpDir(), h.GetRoot())
	}
	return h.BuildCmd
}

func (h *HotReload) GetBin() string {
	if h.Bin == "" {
		return fmt.Sprintf("%s/main", h.GetTmpDir())
	}
	return h.Bin
}

func (h *HotReload) GetTmpDir() string {
	if h.TmpDir == "" {
		return "tmp"
	}
	return h.TmpDir
}

func (h *HotReload) IsTmpDir(dir string) bool {
	return h.GetTmpDir() == dir
}
