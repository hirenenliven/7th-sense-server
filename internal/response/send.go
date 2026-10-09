package response

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

// envelope is the shape of every API response: {"statusCode": 200, "data": ...}.
type envelope struct {
	StatusCode int `json:"statusCode"`
	Data       any `json:"data"`
}

// Send writes data wrapped in the standard envelope. A string is sent as
// {"message": "..."}.
func Send(w http.ResponseWriter, statusCode int, data any) {
	if message, ok := data.(string); ok {
		data = map[string]string{"message": message}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(envelope{StatusCode: statusCode, Data: data}); err != nil {
		log.Printf("response: encode failed: %v", err)
	}
}

// SendError writes err as a response. An *APIError keeps its status code;
// any other error is logged and sent as a 500 so internals never leak.
func SendError(w http.ResponseWriter, err error) {
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		log.Printf("response: unexpected error: %v", err)
		apiErr = ErrInternalServer
	}
	Send(w, apiErr.StatusCode, map[string]any{
		"success": false,
		"error":   apiErr.ClientMessage,
	})
}
