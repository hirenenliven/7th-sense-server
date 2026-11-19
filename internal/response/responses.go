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
	// MachineCycleRunning: the machine already has a running cycle (one at a time per machine).
	MachineCycleRunning = newResponse(http.StatusConflict, "The machine already has a running cycle")

	// Authentication and authorization responses
	InvalidCredentials   = newResponse(http.StatusUnauthorized, "Invalid username or password")
	LoginRequired        = newResponse(http.StatusUnauthorized, "Please log in to continue")
	Forbidden            = newResponse(http.StatusForbidden, "You do not have permission to do this")
	AccountInactive      = newResponse(http.StatusForbidden, "Your account is inactive. Please contact your administrator")
	TooManyLoginAttempts = newResponse(http.StatusTooManyRequests, "Too many login attempts. Please try again in 15 minutes")
	// CustomerNotFound: log in as a user that does not exist, is deleted or is not a customer.
	CustomerNotFound = newResponse(http.StatusNotFound, "Customer not found")
	// SelectCustomer: log in as without a customer.
	SelectCustomer = newResponse(http.StatusBadRequest, "Please select customer")

	// Session responses (server ADR 0014). The three message codes are fixed strings the
	// client tells apart and words itself; do not change them.
	// SessionRevoked: the session was ended by another login, the user or the superadmin.
	SessionRevoked = newResponse(http.StatusUnauthorized, "session_revoked")
	// SessionExpired: the session outlived its lifetime.
	SessionExpired = newResponse(http.StatusUnauthorized, "session_expired")
	// SessionConflict: the login would exceed the user's sessions (SESSION_MODE / MAX_SESSIONS);
	// data lists the active sessions.
	SessionConflict = newResponse(http.StatusConflict, "session_conflict")
	// SessionNotFound: the session does not exist, is not active or is not the caller's.
	SessionNotFound = newResponse(http.StatusNotFound, "Session not found")
	// SessionUserNotFound: the superadmin names a user that does not exist or is deleted.
	SessionUserNotFound = newResponse(http.StatusNotFound, "User not found")

	// Operator task responses (server ADR 0015)
	OperatorTaskNotFound   = newResponse(http.StatusNotFound, "Operator task not found")
	OperatorTaskExists     = newResponse(http.StatusBadRequest, "This task is already added for this customer and machine type")
	SelectTask             = newResponse(http.StatusBadRequest, "Please select task")
	SelectValidTask        = newResponse(http.StatusBadRequest, "Please select valid task")
	SelectValidCustomer    = newResponse(http.StatusBadRequest, "Please select valid customer")
	SelectValidMachineType = newResponse(http.StatusBadRequest, "Please select valid machine type")
	EnterTask              = newResponse(http.StatusBadRequest, "Please enter task")
	TaskTooLong            = newResponse(http.StatusBadRequest, "Task must be at most 100 characters")
	MaintenanceTaskExists  = newResponse(http.StatusBadRequest, "This task already exists")
	MaintenanceTaskMissing = newResponse(http.StatusNotFound, "Task not found")

	// Machine accuracy responses (Account Center › Machine Accuracy)
	EnterMachineAccuracy       = newResponse(http.StatusBadRequest, "Please enter machine accuracy")
	MachineAccuracyTooLong     = newResponse(http.StatusBadRequest, "Machine accuracy must be at most 50 characters")
	MachineAccuracyExists      = newResponse(http.StatusBadRequest, "This machine accuracy already exists")
	MachineAccuracyMissing     = newResponse(http.StatusNotFound, "Machine accuracy not found")
	SelectValidMachineAccuracy = newResponse(http.StatusBadRequest, "Please select valid machine accuracy")
	SelectValidAccuracyLevel   = newResponse(http.StatusBadRequest, "Please select valid machine accuracy level")

	// Machine PLC node id responses (All Machines › PLC node IDs)
	EnterPlcNodeID      = newResponse(http.StatusBadRequest, "Please enter node ID")
	PlcNodeIDTooLong    = newResponse(http.StatusBadRequest, "Node ID must be at most 255 characters")
	InvalidPlcNodeID    = newResponse(http.StatusBadRequest, "Please enter a valid OPC UA node ID, e.g. ns=2;s=/Random/Plc/DB3300.DBX3.0")
	PlcVariableNotFound = newResponse(http.StatusNotFound, "PLC variable not found for this machine")

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
