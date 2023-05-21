package form

import (
	"mime/multipart"
)

type F map[string][]*multipart.FileHeader

// File prepares the file part of the form-data for further formatting
func (m M) File() F {
	return F(m.r.MultipartForm.File)
}

// Get returns the first file for the given key.
// It returns nil if the key does not exist.
func (f F) Get(key string) *multipart.FileHeader {
	if f[key] == nil {
		return nil
	}
	return f[key][0]
}

// All returns all files for the given key.
// It returns nil if the key does not exist.
func (f F) All(key string) []*multipart.FileHeader {
	return f[key]
}
