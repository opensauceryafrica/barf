package server

import (
	"encoding/json"
	"net/http"

	"github.com/opensaucerer/barf/typing"
)

// JSON writes a JSON response to the response writer
func JSON(w http.ResponseWriter, status bool, statusCode int, message string, data map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(typing.Response{
		Status:  status,
		Message: message,
		Data:    data,
	})
}
