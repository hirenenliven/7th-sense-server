package middleware

import (
	"net/http"

	"7th-sense-server/internal/response"
)

// CrossOrigin rejects state-changing requests (POST, PUT, DELETE, ...) sent by
// a browser from another site, so a malicious page cannot use the session
// cookie (CSRF). Same-origin browser requests and non-browser clients pass.
func CrossOrigin(nextHandler http.Handler) http.Handler {
	protection := http.NewCrossOriginProtection()
	protection.SetDenyHandler(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		response.SendResponse(writer, request, response.Forbidden, nil)
	}))
	return protection.Handler(nextHandler)
}
