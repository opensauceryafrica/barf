package constant

import (
	"os"
)

const (

	// EnvTagName is the tag name for environment variables struct
	EnvTag = "barfenv"

	// ShutdownTimeout is the time to wait for the server to shutdown gracefully
	ShutdownTimeout = 5 // seconds

	// WriteTimeout is the maximum duration before timing out
	// writes of the response.
	WriteTimeout = 10 // seconds

	// ReadTimeout is the maximum duration for reading the entire
	// request, including the body.
	ReadTimeout = 10 // seconds

	// MaxHeaderBytes is the maximum number of bytes the server will
	// read parsing the request header's keys and values, including the
	// request line. It does not limit the size of the request body.
	MaxHeaderBytes = 1 << 20 // 1 MB

	// Host is the host for the server to listen on
	Host = ""

	// Port is the port for the server to listen on
	Port = "21186"

	// EnvPath is the path to the environment variables file
	EnvPath = ".env"

	// UseHTTPS enables barf to use https by default.
	UseHTTPS = false
)

var (
	// Logging is for defining whether or not to enable request logging
	Logging = true

	// Recovery is for defining whether or not to enable panic recovery
	Recovery = true

	// ShutdownChan is the channel to listen for shutdown signals
	ShutdownChan = make(chan os.Signal, 1)

	// envTagKeys is the list of keys for environment variables struct
	EnvTagKeys = map[string]interface{}{
		"required": []string{"true", "false"},
		"key":      nil,
	}
)
