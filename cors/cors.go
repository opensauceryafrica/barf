package cors

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/opensaucerer/barf/helper"
	"github.com/opensaucerer/barf/typing"
)

type cors struct {
	allowedOrigins       []string
	allowedMethods       []string
	allowedHeaders       []string
	exposedHeaders       []string
	allowCredentials     bool
	maxAge               int
	optionsPassthrough   bool
	optionsSuccessStatus int
	allowAllOrigins      bool
	// allowAllMethods      bool
	allowAllHeaders bool

	allowedOriginFunc            func(origin string) bool
	allowedOriginWithRequestFunc func(origin string, r *http.Request) bool

	allowedOriginsWildCard []originWildCard
}

// originAllowed returns true if origin is allowed, otherwise false.
func (c *cors) originAllowed(origin string, r *http.Request) bool {
	if c.allowedOriginFunc != nil {
		return c.allowedOriginFunc(origin)
	}
	if c.allowedOriginWithRequestFunc != nil {
		return c.allowedOriginWithRequestFunc(origin, r)
	}
	if c.allowAllOrigins {
		return c.allowAllOrigins
	}
	origin = strings.ToLower(origin)
	for _, o := range c.allowedOrigins {
		if o == origin {
			return true
		}
	}
	for _, o := range c.allowedOriginsWildCard {
		if o.matchOrigin(origin) {
			return true
		}
	}
	return false
}

// methodAllowed returns true if method is allowed, otherwise false.
func (c *cors) methodAllowed(method string) bool {
	if len(c.allowedMethods) == 0 {
		return false
	}
	method = strings.ToUpper(method)
	if method == http.MethodOptions {
		return true
	}
	for _, m := range c.allowedMethods {
		if m == method {
			return true
		}
	}
	return false
}

// headersAllowed returns true if headers are allowed, otherwise false.
func (c *cors) headersAllowed(headers []string) bool {
	if c.allowAllHeaders && len(headers) == 0 {
		return true
	}
	for _, h := range headers {
		h = http.CanonicalHeaderKey(h)
		if !helper.StringArray(c.allowedHeaders).Contains(h) {
			return false
		}
	}
	return true
}

// preflight handles pre-flight CORS requests.
func (c *cors) preflight(w http.ResponseWriter, r *http.Request) {
	headers := w.Header()
	origin := r.Header.Get("Origin")

	if r.Method != http.MethodOptions {
		return
	}

	headers.Set("Vary", "Origin")
	headers.Set("Vary", "Access-Control-Request-Method")
	headers.Set("Vary", "Access-Control-Request-Headers")

	// maybe allow private networks

	if origin == "" || !c.originAllowed(origin, r) {
		return
	}

	requestMethod := r.Header.Get("Access-Control-Request-Method")
	if requestMethod == "" || !c.methodAllowed(requestMethod) {
		return
	}

	rheaders := strings.Join(r.Header.Values("Access-Control-Request-Headers"), ",")
	parsedHeaders := parseHeaders(rheaders)
	if !c.headersAllowed(parsedHeaders) {
		return
	}
	if c.allowAllOrigins {
		headers.Set("Access-Control-Allow-Origin", "*")
	} else {
		headers.Set("Access-Control-Allow-Origin", origin)
	}

	headers.Set("Access-Control-Allow-Methods", strings.ToUpper(requestMethod))

	if len(parsedHeaders) > 0 {
		headers.Set("Access-Control-Allow-Headers", strings.Join(parsedHeaders, ","))
	}
	if c.allowCredentials {
		headers.Set("Access-Control-Allow-Credentials", "true")
	}
	if c.maxAge > 0 {
		headers.Set("Access-Control-Max-Age", strconv.Itoa(c.maxAge))
	}
}

// request handles cors requests, http requests or http redirects
func (c *cors) request(w http.ResponseWriter, r *http.Request) {
	headers := w.Header()
	origin := r.Header.Get("Origin")

	headers.Set("Vary", "Origin")
	if origin == "" || !c.originAllowed(origin, r) {
		return
	}
	if !c.methodAllowed(r.Method) {
		return
	}
	if c.allowAllOrigins {
		headers.Set("Access-Control-Allow-Origin", "*")
	} else {
		headers.Set("Access-Control-Allow-Origin", origin)
	}
	if c.allowCredentials {
		headers.Set("Access-Control-Allow-Credentials", "true")
	}
	if len(c.exposedHeaders) > 0 {
		headers.Set("Access-Control-Expose-Headers", strings.Join(c.exposedHeaders, ","))
	}
}

type originWildCard struct {
	prefix string
	suffix string
}

// matchOrigin returns true if the origin matches the allowed origins.
func (ow *originWildCard) matchOrigin(origin string) bool {
	return len(origin) >= len(ow.prefix)+len(ow.suffix) && strings.HasPrefix(origin, ow.prefix) && strings.HasSuffix(origin, ow.suffix)
}

// Prepare sets up the CORS middleware based on the given options.
func Prepare(options typing.CORS) *cors {

	c := &cors{
		exposedHeaders:               helper.StringArray(options.ExposedHeaders).Transform(http.CanonicalHeaderKey),
		allowCredentials:             options.AllowCredentials,
		maxAge:                       options.MaxAge,
		optionsPassthrough:           options.OptionsPassthrough,
		allowedOriginFunc:            options.AllowedOriginFunc,
		allowedOriginWithRequestFunc: options.AllowedOriginWithRequestFunc,
	}

	// allowed origins
	if len(options.AllowedOrigins) > 0 {
		c.allowedOrigins = make([]string, 0)
		c.allowedOriginsWildCard = make([]originWildCard, 0)

		for _, origin := range options.AllowedOrigins {
			if origin == "*" {
				c.allowAllOrigins = true
				c.allowedOrigins = nil
				c.allowedOriginsWildCard = nil
				break
			} else if index := strings.IndexByte(origin, '*'); index >= 0 {
				c.allowedOriginsWildCard = append(c.allowedOriginsWildCard, originWildCard{
					prefix: origin[:index],
					suffix: origin[index+1:],
				})
			} else {
				c.allowedOrigins = append(c.allowedOrigins, origin)
			}
		}
	} else {
		if options.AllowedOriginFunc == nil && options.AllowedOriginWithRequestFunc == nil {
			c.allowAllOrigins = true
		}
	}

	// allowed methods
	if len(options.AllowedMethods) > 0 {
		methods := helper.StringArray(options.AllowedMethods).Transform(strings.ToUpper)
		c.allowedMethods = methods
	} else {
		c.allowedMethods = []string{
			http.MethodGet,
			http.MethodHead,
			http.MethodPost,
			http.MethodPut,
			http.MethodHead,
			http.MethodPatch,
			http.MethodDelete,
		}
	}

	// allowed headers
	if len(options.AllowedHeaders) > 0 {
		headers := helper.StringArray(options.AllowedHeaders).Transform(http.CanonicalHeaderKey)
		c.allowedHeaders = headers
		// handle wildcard
		for _, header := range options.AllowedHeaders {
			if header == "*" {
				c.allowAllHeaders = true
				c.allowedHeaders = nil
				break
			}
		}
	} else {
		// default allowed headers
		c.allowedHeaders = []string{"Origin", "Accept", "Content-Type", "X-Requested-With", "Authorization"}
	}

	// options success status
	if options.OptionsSuccessStatus > 0 {
		c.optionsSuccessStatus = options.OptionsSuccessStatus
	} else {
		c.optionsSuccessStatus = http.StatusNoContent
	}

	return c
}

// ParseHeaders normalizes a string containing a list of headers
func parseHeaders(headers string) []string {
	const lower = 'a' - 'A'
	length := len(headers)
	heigth := make([]byte, 0, length)
	uppercase := true
	total := 0
	for i := 0; i < length; i++ {
		if headers[i] == ',' {
			total++
		}
	}

	result := make([]string, 0, total)

	for i := 0; i < length; i++ {
		bit := headers[i]

		switch {
		case bit >= 'a' && bit <= 'z':
			if uppercase {
				heigth = append(heigth, bit-lower)
			} else {
				heigth = append(heigth, bit)
			}
		case bit >= 'A' && bit <= 'Z':
			if !uppercase {
				heigth = append(heigth, bit+lower)
			} else {
				heigth = append(heigth, bit)
			}
		case bit == '-' || bit == '_' || bit == '.' || (bit >= '0' && bit <= '9'):
			heigth = append(heigth, bit)
		}

		if bit == ' ' || bit == ',' || i == length-1 {
			if len(heigth) > 0 {
				result = append(result, string(heigth))
				heigth = heigth[:0]
				uppercase = true
			}
		} else {
			uppercase = bit == '-' || bit == '_'
		}
	}
	return result
}
