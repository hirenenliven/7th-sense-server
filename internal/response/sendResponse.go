package response

import (
	"encoding/json"
	"errors"
	"net/http"

	"7th-sense-server/internal/logger"
)

// body is the shape of every API response.
type body struct {
	StatusCode int    `json:"statusCode"`
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	Data       any    `json:"data,omitempty"`
}

// SendResponse writes res as a JSON response. data is the payload; a string
// replaces the response message instead, and nil sends no data.
//
//	response.SendResponse(writer, request, response.Success, items)
//	response.SendResponse(writer, request, response.Created, item)
//	response.SendResponse(writer, request, response.Success, "User deleted")
//	response.SendResponse(writer, request, err, nil) // error returned by a service
//
// res is a predefined *Response or any error. Errors that are not a *Response
// are logged and sent as ServerError so internals never leak.
func SendResponse(writer http.ResponseWriter, request *http.Request, res error, data any) {
	var resp *Response
	if !errors.As(res, &resp) {
		logger.Log(logger.LevelError, "unexpected error",
			"method", request.Method,
			"path", request.URL.Path,
			"error", res,
		)
		resp = ServerError
	}

	msg := resp.Message
	if s, ok := data.(string); ok {
		msg, data = s, nil
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(resp.StatusCode)
	err := json.NewEncoder(writer).Encode(body{
		StatusCode: resp.StatusCode,
		Success:    resp.StatusCode < http.StatusBadRequest,
		Message:    msg,
		Data:       data,
	})
	if err != nil {
		logger.Log(logger.LevelError, "response encode failed", "error", err)
	}
}
