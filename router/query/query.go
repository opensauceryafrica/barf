package query

import (
	"encoding/json"
	"net/http"
)

type Q []byte

// Query prepares the barf request with the request query for further formatting
func Query(r *http.Request) []byte {
	q := make(map[string]string)
	for k, v := range r.URL.Query() {
		q[k] = v[0]
	}
	json, _ := json.Marshal(q)
	return json
}

// JSON formats the request query as map[string]string.
// It returns an error if the query is not a valid JSON.
func (q Q) JSON() (map[string]string, error) {
	var data map[string]string
	err := json.Unmarshal(q, &data)
	return data, err
}

// Format formats the request query into the given interface. v must be a pointer.
// It returns an error if the query is not a valid JSON.
func (q Q) Format(v interface{}) error {
	return json.Unmarshal(q, v)
}
