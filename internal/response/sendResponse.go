package response

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
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
//	response.SendResponse(w, r, response.Success, items)
//	response.SendResponse(w, r, response.Created, item)
//	response.SendResponse(w, r, response.Success, "User deleted")
//	response.SendResponse(w, r, err, nil) // error returned by a service
//
// res is a predefined *Response or any error. Errors that are not a *Response
// are logged and sent as ServerError so internals never leak.
func SendResponse(w http.ResponseWriter, r *http.Request, res error, data any) {
	var resp *Response
	if !errors.As(res, &resp) {
		log.Printf("response: unexpected error on %s %s: %v", r.Method, r.URL.Path, res)
		resp = ServerError
	}

	msg := resp.Message
	if s, ok := data.(string); ok {
		msg, data = s, nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	err := json.NewEncoder(w).Encode(body{
		StatusCode: resp.StatusCode,
		Success:    resp.StatusCode < http.StatusBadRequest,
		Message:    msg,
		Data:       data,
	})
	if err != nil {
		log.Printf("response: encode failed: %v", err)
	}
}
