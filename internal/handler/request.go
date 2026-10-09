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

// decodeAndValidate reads the JSON body into requestBody and runs its validation.
func decodeAndValidate(writer http.ResponseWriter, request *http.Request, requestBody validator) error {
	request.Body = http.MaxBytesReader(writer, request.Body, maxBodyBytes)
	if err := json.NewDecoder(request.Body).Decode(requestBody); err != nil {
		return response.InvalidJSON
	}
	return requestBody.Validate()
}
