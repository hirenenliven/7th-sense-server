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

// NewError creates a new error response with custom message
func NewError(statusCode int, message string) *Response {
	return newResponse(statusCode, message)
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
	InvalidJSON       = newResponse(http.StatusBadRequest, "Invalid JSON body")
	TryAgain          = newResponse(http.StatusBadRequest, "Something went wrong, Please try again later")
	// MachineNotFound: the machine does not exist, is deleted or is not active.
	MachineNotFound = newResponse(http.StatusNotFound, "Machine not found")
	// SelectMachine: a request that needs a machine names none.
	SelectMachine = newResponse(http.StatusBadRequest, "Please select machine")
	// MachineCycleNotFound: the cycle does not exist or its machine is not one the user may read.
	MachineCycleNotFound = newResponse(http.StatusNotFound, "Machine cycle not found")

	// Authentication and authorization responses
	InvalidCredentials   = newResponse(http.StatusUnauthorized, "Invalid username or password")
	LoginRequired        = newResponse(http.StatusUnauthorized, "Please log in to continue")
	Forbidden            = newResponse(http.StatusForbidden, "You do not have permission to do this")
	AccountInactive      = newResponse(http.StatusForbidden, "Your account is inactive. Please contact your administrator")
	TooManyLoginAttempts = newResponse(http.StatusTooManyRequests, "Too many login attempts. Please try again in 15 minutes")

	// Validation error responses
	ValidName              = newResponse(http.StatusBadRequest, "Please enter valid name")
	ValidCode              = newResponse(http.StatusBadRequest, "Please enter valid code")
	ValidDescription       = newResponse(http.StatusBadRequest, "Please enter valid description")
	ValidCategory          = newResponse(http.StatusBadRequest, "Please enter valid category")
	ValidImage             = newResponse(http.StatusBadRequest, "Please enter valid image")
	ValidBrand             = newResponse(http.StatusBadRequest, "Please enter valid brand")
	SelectBrand            = newResponse(http.StatusBadRequest, "Please select brand")
	SelectController       = newResponse(http.StatusBadRequest, "Please select controller")
	UserProfileNotFound    = newResponse(http.StatusNotFound, "User profile not found")
	CompanyProfileNotFound = newResponse(http.StatusNotFound, "Company profile not found")
)
