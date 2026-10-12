package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

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
		var maxBytesError *http.MaxBytesError
		if errors.As(err, &maxBytesError) {
			return response.RequestEntryLarge
		}
		return response.InvalidJSON
	}
	return requestBody.Validate()
}

// pathID reads the {id} path value as a positive integer.
func pathID(request *http.Request) (int64, error) {
	id, err := strconv.ParseInt(request.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		return 0, response.InvalidID
	}
	return id, nil
}
