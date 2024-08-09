package form

import (
	"net/http"
)

type M struct {
	r *http.Request
}

// Form prepares the barf request with the request form or multipart/form-data for further formatting
func Form(r *http.Request, maxMemory int64) M {
	r.ParseMultipartForm(maxMemory)
	return M{r}
}
