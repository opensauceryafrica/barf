package param

import (
	"encoding/json"
	"net/http"

	"github.com/opensaucerer/barf/typing"
)

type P []byte

// Params prepares the barf request with the request params for further formatting
func Params(r *http.Request) []byte {
	p := r.Context().Value(typing.ParamsCtxKey{})
	if p == nil {
		json, _ := json.Marshal(map[string]string{})
		return json
	}
	json, _ := json.Marshal(p)
	return json
}

// JSON formats the request params as map[string]string.
// It returns an error if the params is not a valid JSON.
func (p P) JSON() (map[string]string, error) {
	var data map[string]string
	err := json.Unmarshal(p, &data)
	return data, err
}

// Format formats the request params into the given interface v which must be a pointer.
// It returns an error if the params is not a valid JSON.
func (p P) Format(v interface{}) error {
	return json.Unmarshal(p, v)
}
