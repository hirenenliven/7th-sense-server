package response

import "net/http"

// Response is a predefined API response: an HTTP status code and the message
// the client sees. It implements error, so services and DTOs can return
// failure responses directly.
type Response struct {
	StatusCode int
	Message    string
}

func (apiResponse *Response) Error() string { return apiResponse.Message }

func newResponse(statusCode int, message string) *Response {
	return &Response{StatusCode: statusCode, Message: message}
}

// Common responses (same set and names as the Node.js RESPONSES object).
var (
	Success           = newResponse(http.StatusOK, "success")
	Created           = newResponse(http.StatusCreated, "created")
	BadRequest        = newResponse(http.StatusBadRequest, "bad request")
	Unauthorized      = newResponse(http.StatusUnauthorized, "unauthorized")
	PaymentRequired   = newResponse(http.StatusPaymentRequired, "payment required")
	NotFound          = newResponse(http.StatusNotFound, "not found")
	ServerTimeout     = newResponse(http.StatusRequestTimeout, "server timeout")
	RequestEntryLarge = newResponse(http.StatusRequestEntityTooLarge, "request entry to large")
	RequestURLToLong  = newResponse(http.StatusRequestURITooLong, "request-URL too long")
	UnsupportedMedia  = newResponse(http.StatusUnsupportedMediaType, "unsupported media type")
	ExpectationFailed = newResponse(http.StatusExpectationFailed, "expectation failed")
	ServerError       = newResponse(http.StatusInternalServerError, "internal server error")
	Unavailable       = newResponse(http.StatusServiceUnavailable, "service unavailable")
	ValidEmail        = newResponse(http.StatusBadRequest, "Please enter valid email")
	ValidUserName     = newResponse(http.StatusBadRequest, "Please enter valid UserName")
	ValidPassword     = newResponse(http.StatusBadRequest, "Please enter valid password")
	TryAgain          = newResponse(http.StatusBadRequest, "Something went wrong, Please try again later")
)

// Responses specific to this project's entities.
var (
	InvalidJSON            = newResponse(http.StatusBadRequest, "Invalid JSON body")
	InvalidID              = newResponse(http.StatusBadRequest, "Please enter valid id")
	ValidName              = newResponse(http.StatusBadRequest, "Please enter valid name")
	ValidMobile            = newResponse(http.StatusBadRequest, "Please enter valid mobile")
	ValidStatus            = newResponse(http.StatusBadRequest, "Please enter valid status")
	ValidPrice             = newResponse(http.StatusBadRequest, "Please enter valid price")
	ValidRole              = newResponse(http.StatusBadRequest, "Please enter valid role")
	EmailExists            = newResponse(http.StatusBadRequest, "Email already exists")
	UsernameExists         = newResponse(http.StatusBadRequest, "Username already exists")
	UserNotFound           = newResponse(http.StatusNotFound, "User not found")
	ItemNotFound           = newResponse(http.StatusNotFound, "Item not found")
	ValidCode              = newResponse(http.StatusBadRequest, "Please enter valid code")
	ValidAddress           = newResponse(http.StatusBadRequest, "Please enter valid address")
	ValidCity              = newResponse(http.StatusBadRequest, "Please enter valid city")
	ValidGSTNumber         = newResponse(http.StatusBadRequest, "Please enter valid GST number")
	ValidPANNumber         = newResponse(http.StatusBadRequest, "Please enter valid PAN number")
	ValidColorCode         = newResponse(http.StatusBadRequest, "Please enter valid color code")
	BranchNotFound         = newResponse(http.StatusNotFound, "Branch not found")
	CodeExists             = newResponse(http.StatusBadRequest, "Code already exists")
	UserProfileNotFound    = newResponse(http.StatusNotFound, "User profile not found")
	CompanyProfileNotFound = newResponse(http.StatusNotFound, "Company profile not found")
)
