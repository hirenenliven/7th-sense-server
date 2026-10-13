package handler

import (
	"net/http"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/middleware"
	"7th-sense-server/internal/response"
	"7th-sense-server/internal/service"
)

// sessionCookiePath limits the session cookie to API requests.
const sessionCookiePath = "/api"

// AuthHandler handles login, logout and the logged-in user. The session token
// travels only in an HttpOnly cookie, never in a response body.
type AuthHandler struct {
	service      *service.AuthService
	secureCookie bool // send the cookie over HTTPS only
}

// NewAuthHandler creates an AuthHandler that uses authService. secureCookie
// must be true when the API is served over HTTPS.
func NewAuthHandler(authService *service.AuthService, secureCookie bool) *AuthHandler {
	return &AuthHandler{service: authService, secureCookie: secureCookie}
}

// Login handles POST /api/v1/auth/login: checks the credentials, sets the
// session cookie and returns the user.
func (authHandler *AuthHandler) Login(writer http.ResponseWriter, request *http.Request) {
	var loginRequest dto.LoginRequest
	if err := decodeAndValidate(writer, request, &loginRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	session, err := authHandler.service.Login(request.Context(), loginRequest, request.UserAgent())
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	cookie := authHandler.sessionCookie(session.Token)
	// Without Remember the cookie ends with the browser session; the server
	// also expires the session after its lifetime.
	if session.Remember {
		cookie.MaxAge = int(session.Lifetime.Seconds())
	}
	http.SetCookie(writer, cookie)
	response.SendResponse(writer, request, response.Success, session.User)
}

// Logout handles POST /api/v1/auth/logout: ends the session and clears the
// cookie. It succeeds without a session, so it is safe to call twice.
func (authHandler *AuthHandler) Logout(writer http.ResponseWriter, request *http.Request) {
	if cookie, err := request.Cookie(middleware.SessionCookieName); err == nil {
		if err := authHandler.service.Logout(request.Context(), cookie.Value); err != nil {
			response.SendResponse(writer, request, err, nil)
			return
		}
	}

	cookie := authHandler.sessionCookie("")
	cookie.MaxAge = -1 // delete
	http.SetCookie(writer, cookie)
	response.SendResponse(writer, request, response.Success, "Logged out")
}

// Me handles GET /api/v1/auth/me: the logged-in user with its role.
func (authHandler *AuthHandler) Me(writer http.ResponseWriter, request *http.Request) {
	user, ok := middleware.CurrentUser(request.Context())
	if !ok {
		response.SendResponse(writer, request, response.LoginRequired, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, user)
}

// sessionCookie returns the session cookie with value and the shared attributes.
func (authHandler *AuthHandler) sessionCookie(value string) *http.Cookie {
	return &http.Cookie{
		Name:     middleware.SessionCookieName,
		Value:    value,
		Path:     sessionCookiePath,
		HttpOnly: true,
		Secure:   authHandler.secureCookie,
		SameSite: http.SameSiteStrictMode,
	}
}
