package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"7th-sense-server/internal/middleware"
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
		return response.NewError(http.StatusBadRequest, "Please enter valid json")
	}
	return requestBody.Validate()
}

// pathID reads the {id} path value as a positive integer.
func pathID(request *http.Request) (int64, error) {
	id, err := strconv.ParseInt(request.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		return 0, response.NewError(http.StatusBadRequest, "Please enter valid id")
	}
	return id, nil
}

// queryID reads the optional query parameter name as a positive integer. It
// returns nil when the parameter is missing or empty.
func queryID(request *http.Request, name string) (*int64, error) {
	value := request.URL.Query().Get(name)
	if value == "" {
		return nil, nil
	}
	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil || id <= 0 {
		return nil, response.NewError(http.StatusBadRequest, "Please enter valid "+name)
	}
	return &id, nil
}

// currentUserID returns the id of the user making the request (set by
// middleware.Authenticate), or Unauthorized when there is none.
func currentUserID(request *http.Request) (int64, error) {
	userID, ok := middleware.UserID(request.Context())
	if !ok {
		return 0, response.Unauthorized
	}
	return userID, nil
}
