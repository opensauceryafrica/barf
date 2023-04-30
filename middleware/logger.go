package middleware

import (
	"net/http"
	"reflect"
	"strconv"
	"time"
	"unsafe"

	logger "github.com/opensaucerer/barf/log"
	"github.com/opensaucerer/barf/server"
)

// Logger logs the request to the console
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if server.Augment.Logging { // create an addressable copy of the response writer
			wt := reflect.ValueOf(w).Elem()
			aw := reflect.New(wt.Type()).Elem()
			aw.Set(wt)
			// field 16 of the struct embeded in http.ResponseWriter is the status code
			uf := aw.Field(16)
			uf = reflect.NewAt(uf.Type(), unsafe.Pointer(uf.UnsafeAddr())).Elem()
			// get status code
			code := uf.Interface().(int)
			// format: iso timestamp: user-agent - http/version: method - path - status code - status text
			msg := time.Now().Format(time.RFC3339) + ": " + r.UserAgent() + " - " + r.Proto + ": " + r.Method + " - " + r.URL.Path + " - " + strconv.Itoa(code) + " - " + http.StatusText(code)
			// log request
			logger.Code(msg, code)

			// call next middleware
			// next.ServeHTTP(w, r)
		}
	})
}
