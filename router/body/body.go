package body

import (
	"encoding/json"
	"net/http"
)

type B []byte

// Body prepares the barf request with the request body for further formatting
func Body(r *http.Request) []byte {
	body := make([]byte, r.ContentLength)
	r.Body.Read(body)
	return body
}

// JSON formats the body as map[string]interface{}.
// It returns an error if the body is not a valid JSON.
func (b B) JSON() (map[string]interface{}, error) {
	var data map[string]interface{}
	err := json.Unmarshal(b, &data)
	return data, err
}

// Format formats the request params into the given interface v which must be a pointer.
// It returns an error if the body is not a valid JSON.
func (b B) Format(v interface{}) error {
	return json.Unmarshal(b, v)
}
