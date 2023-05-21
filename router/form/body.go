package form

import (
	"encoding/json"

	"github.com/opensaucerer/barf/router/body"
)

// Body prepares the non-file part of the form-data for further formatting
func (m M) Body() body.B {
	// write r.MultipartForm.Value to f.r.Body
	val := make(map[string]string)
	for k, v := range m.r.MultipartForm.Value {
		val[k] = v[0]
	}
	b, _ := json.Marshal(val)
	return b
}
