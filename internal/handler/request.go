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
func decodeAndValidate(writer http.ResponseWriter, request *http.Request, dst validator) error {
	request.Body = http.MaxBytesReader(writer, request.Body, maxBodyBytes)
	if err := json.NewDecoder(request.Body).Decode(dst); err != nil {
		return response.InvalidJSON
	}
	return dst.Validate()
}
