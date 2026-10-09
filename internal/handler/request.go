package handler

import (
	"encoding/json"
	"net/http"

	"7th-sense-server/internal/response"
)

const maxBodyBytes = 1 << 20 // 1 MB

// validator is implemented by every request DTO.
type validator interface {
	Validate() error
}

// decodeAndValidate reads the JSON body into dst and runs its validation.
func decodeAndValidate(w http.ResponseWriter, r *http.Request, dst validator) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return response.InvalidJSON
	}
	return dst.Validate()
}
