package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/opensaucerer/barf/constant"
	"github.com/opensaucerer/barf/typing"
)

// JSON writes a JSON response to the response writer
func JSON(w http.ResponseWriter, status bool, statusCode int, message string, data map[string]interface{}) {
	log.Println("JSON", status, statusCode, message, data)
	Response(w).Status(statusCode).JSON(typing.Response{
		Status:  status,
		Message: message,
		Data:    data,
	})
}

type response struct {
	code int
	rw   http.ResponseWriter
}

// JSON writes a JSON response to the response writer
func (r *response) JSON(data interface{}) {
	r.rw.Header().Set(constant.ContentType, constant.ApplicationJSON)
	r.rw.WriteHeader(r.code)
	json.NewEncoder(r.rw).Encode(data)
	Write(r.rw) // write to underlying response writer
}

// Status loads a barf response with the given status code
func (r *response) Status(code int) *response {
	r.code = code
	return r
}

// Response prepares a barf response with the given writer
func Response(w http.ResponseWriter) *response {
	return &response{
		rw: w,
	}
}

type ResponseWriter struct {
	rw      http.ResponseWriter
	Written bool // set to true if the underlying response writer has been written to
	Code    int
	Body    []byte
	Headers http.Header
}

func (w *ResponseWriter) Header() http.Header {
	if w.Headers == nil {
		w.Headers = make(http.Header)
	}
	return w.Headers
}

func (w *ResponseWriter) WriteHeader(code int) {
	w.Code = code
	// w.ResponseWriter.WriteHeader(code) // don't write to underlying response writer until Recover() has reported error or nil, if none.
}

func (w *ResponseWriter) Write(b []byte) (int, error) {
	w.Body = b
	return len(b), nil
}

// Loaded returns true if the response has been overloaded by barf's server.ResponseWriter
func Loaded(w http.ResponseWriter) bool {
	_, ok := w.(*ResponseWriter)
	return ok
}

// Written returns true if the response has been written to
func Written(w http.ResponseWriter) bool {
	if Loaded(w) {
		return w.(*ResponseWriter).Written
	}
	return false
}

// Load overloads the response writer with barf's server.ResponseWriter
func Load(w http.ResponseWriter) http.ResponseWriter {
	if Loaded(w) {
		return w
	}
	return &ResponseWriter{
		rw:      w,
		Written: false,
	}
}

// Write performs the actual write to the underlying response writer using barf's server.ResponseWriter if it has not been written to yet
func Write(w http.ResponseWriter) {
	if Loaded(w) && !Written(w) {

		// ensure that it becomes impossible to make a superfluous call to the underlying response writer
		w.(*ResponseWriter).Written = true

		// use the underlying response writer. if at
		// any point the writer has not been as a barf response writer, this
		// will panic (ideally, this should never happen)
		for k, v := range w.(*ResponseWriter).Headers {
			for _, vv := range v {
				w.(*ResponseWriter).rw.Header().Add(k, vv)
			}
		}
		w.(*ResponseWriter).rw.WriteHeader(w.(*ResponseWriter).Code)
		w.(*ResponseWriter).rw.Write(w.(*ResponseWriter).Body)
	}
}

// Status returns the status code written into the response writer using barf's server.ResponseWriter
func Status(w http.ResponseWriter) int {
	if Loaded(w) {
		return w.(*ResponseWriter).Code
	}
	return 200
}
