package response

import "net/http"

// APIError is a predefined API failure. Message is for logs; ClientMessage is
// what the client sees. It implements error so services can return it directly.
type APIError struct {
	StatusCode    int
	Message       string
	ClientMessage string
}

func (e *APIError) Error() string { return e.Message }

func newAPIError(statusCode int, clientMessage string) *APIError {
	return &APIError{
		StatusCode:    statusCode,
		Message:       http.StatusText(statusCode),
		ClientMessage: clientMessage,
	}
}

// Predefined errors (Go equivalent of the Node Responses object).
var (
	ErrBadRequest            = newAPIError(http.StatusBadRequest, "Bad request")
	ErrInvalidJSON           = newAPIError(http.StatusBadRequest, "Invalid JSON body")
	ErrInvalidEmail          = newAPIError(http.StatusBadRequest, "Please enter a valid email")
	ErrInvalidName           = newAPIError(http.StatusBadRequest, "Please enter a valid name")
	ErrInvalidPrice          = newAPIError(http.StatusBadRequest, "Please enter a valid price")
	ErrEmailExists           = newAPIError(http.StatusBadRequest, "Email already exists")
	ErrAlreadyExists         = newAPIError(http.StatusBadRequest, "Already exists")
	ErrTryAgain              = newAPIError(http.StatusBadRequest, "Something went wrong, please try again later")
	ErrUnauthorized          = newAPIError(http.StatusUnauthorized, "Unauthorized")
	ErrPaymentRequired       = newAPIError(http.StatusPaymentRequired, "Payment required")
	ErrForbidden             = newAPIError(http.StatusForbidden, "Forbidden")
	ErrNotFound              = newAPIError(http.StatusNotFound, "Not found")
	ErrUserNotFound          = newAPIError(http.StatusNotFound, "User not found")
	ErrItemNotFound          = newAPIError(http.StatusNotFound, "Item not found")
	ErrRequestTimeout        = newAPIError(http.StatusRequestTimeout, "Request timeout")
	ErrRequestEntityTooLarge = newAPIError(http.StatusRequestEntityTooLarge, "Request entity too large")
	ErrRequestURITooLong     = newAPIError(http.StatusRequestURITooLong, "Request URI too long")
	ErrUnsupportedMediaType  = newAPIError(http.StatusUnsupportedMediaType, "Unsupported media type")
	ErrExpectationFailed     = newAPIError(http.StatusExpectationFailed, "Expectation failed")
	ErrInternalServer        = newAPIError(http.StatusInternalServerError, "Internal server error")
	ErrServiceUnavailable    = newAPIError(http.StatusServiceUnavailable, "Service unavailable")
)
