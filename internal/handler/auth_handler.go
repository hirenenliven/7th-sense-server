package handler

import (
	"errors"
	"net/http"
	"net/netip"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/middleware"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/response"
	"7th-sense-server/internal/service"
)

// sessionCookiePath limits the session cookie to API requests.
const sessionCookiePath = "/api"

// AuthHandler handles login, logout, "log in as" and the logged-in user. The
// session token travels only in an HttpOnly cookie, never in a response body.
type AuthHandler struct {
	service        *service.AuthService
	secureCookie   bool           // send the cookie over HTTPS only
	live           *LiveHandler   // closes a session's live sockets at logout; nil = live data off
	trustedProxies []netip.Prefix // proxies whose X-Forwarded-For is believed (TRUSTED_PROXIES)
}

// NewAuthHandler creates an AuthHandler that uses authService. secureCookie
// must be true when the API is served over HTTPS.
func NewAuthHandler(authService *service.AuthService, secureCookie bool) *AuthHandler {
	return &AuthHandler{service: authService, secureCookie: secureCookie}
}

// SetLiveHandler makes Logout close the ended session's live sockets, and Login those of the
// sessions it revoked.
func (authHandler *AuthHandler) SetLiveHandler(liveHandler *LiveHandler) {
	authHandler.live = liveHandler
}

// SetTrustedProxies makes the client IP stored with a session come from X-Forwarded-For when the
// request arrives through one of trustedProxies (TRUSTED_PROXIES). Without it the TCP peer's
// address is used.
func (authHandler *AuthHandler) SetTrustedProxies(trustedProxies []netip.Prefix) {
	authHandler.trustedProxies = trustedProxies
}

// Login handles POST /api/v1/auth/login: checks the credentials and the session limits, sets the
// session cookie and returns the user. A login over the limit answers 409 session_conflict with
// the sessions in the way; forceLogoutOthers revokes them and logs in.
//
// @Summary      Login
// @Description  Checks the credentials and starts a session. The browser's user agent and IP address (X-Forwarded-For only behind TRUSTED_PROXIES) are stored with the session.
// @Description  Session limits (SESSION_MODE, MAX_SESSIONS; "log in as" sessions never count): with valid credentials, a login while the user already has an active session (single mode) or MAX_SESSIONS of them (multiple mode) answers 409 with message session_conflict, creates no session and keeps the others; data lists the sessions in the way (id, device, browser, os, ipAddress, createdAt, lastUsedAt).
// @Description  forceLogoutOthers: true re-checks the password, revokes the sessions in the way (single mode: all others; multiple mode: the least recently used) with reason replaced_by_new_login and creates the new session in one transaction; the revoked browsers get 401 session_revoked and their live sockets session:expired.
// @Description  Too many wrong passwords (429) is checked before the session limits.
// @Description  Validation: Invalid username or password; Please enter username; Please enter password.
// @Description
// @Description  Used by pages: Login
// @Tags         Auth
// @ID           login
// @Accept       json
// @Produce      json
// @Param        request body dto.LoginRequest true "Request body"
// @Success      200 {object} response.Envelope{data=model.AuthUser}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Invalid username or password"
// @Failure      403 {object} response.Envelope "Account inactive"
// @Failure      409 {object} response.Envelope{data=dto.SessionConflictResponse} "session_conflict: the user already has as many sessions as allowed"
// @Failure      429 {object} response.Envelope "Too many login attempts"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/auth/login [post]
func (authHandler *AuthHandler) Login(writer http.ResponseWriter, request *http.Request) {
	// A password login from a "log in as" tab would replace the superadmin's own cookie
	if _, ok := middleware.LoginAsKey(request.Context()); ok {
		response.SendResponse(writer, request, response.Forbidden, nil)
		return
	}
	var loginRequest dto.LoginRequest
	if err := decodeAndValidate(writer, request, &loginRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	session, err := authHandler.service.Login(request.Context(), loginRequest, authHandler.sessionClient(request))
	var conflictError *service.SessionConflictError
	if errors.As(err, &conflictError) {
		response.SendResponse(writer, request, err, conflictError.Conflict)
		return
	}
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	if authHandler.live != nil {
		authHandler.live.CloseRevokedSessions(session.RevokedSessions, model.RevokeReasonReplacedByNewLogin)
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

// Logout handles POST /api/v1/auth/logout: ends the session, closes its live
// sockets and clears the cookie. It succeeds without a session, so it is safe to call twice. Under
// /api/as/<key>/ it ends only that "log in as" session and clears its cookie;
// the superadmin's own session is not touched.
//
// @Summary      Logout
// @Description  Ends the session of token: its status becomes logged_out (auth_tokens, also soft-deleted). An unknown or empty token is not an error. Called as /api/as/{key}/v1/auth/logout it ends only that "log in as" session.
// @Description  The session's live data sockets (Socket.IO, /api/v1/socket.io/) get session:expired {reason: logged_out} and are closed; the user's other sessions keep theirs.
// @Description
// @Description  Used by pages: App shell (session, logout, theme)
// @Tags         Auth
// @ID           logout
// @Produce      json
// @Success      200 {object} response.Envelope
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/auth/logout [post]
func (authHandler *AuthHandler) Logout(writer http.ResponseWriter, request *http.Request) {
	token := middleware.SessionToken(request)
	if err := authHandler.service.Logout(request.Context(), token); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	// Every tab of this session stops at once; the user's other sessions keep their sockets.
	if authHandler.live != nil {
		authHandler.live.DisconnectSession(token)
	}

	cookie := authHandler.sessionCookie("")
	if key, ok := middleware.LoginAsKey(request.Context()); ok {
		cookie = authHandler.loginAsCookie(key, "")
	}
	cookie.MaxAge = -1 // delete
	http.SetCookie(writer, cookie)
	response.SendResponse(writer, request, response.Success, "Logged out")
}

// LoginAs handles POST /api/v1/auth/login-as: the superadmin opens a customer's
// account in a new browser tab. It starts a session of the customer whose cookie
// is scoped to /api/as/<sessionKey>, so that tab calls the API under that path
// and the superadmin's own session and cookie stay as they are.
//
// @Summary      Log in as customer
// @Description  Starts a session of the customer for the superadmin, in a new browser tab. The session cookie (HttpOnly) is scoped to /api/as/{sessionKey}: that tab calls /api/as/{sessionKey}/v1/... instead of /api/v1/..., so logging out of either tab never ends the other session. The session records the superadmin (auth_tokens.impersonated_by) and lasts one working day.
// @Description  Roles: superadmin (others get 403).
// @Description  Validation: Please select customer; Customer not found (missing, deleted or not a customer); Account inactive.
// @Description
// @Description  Used by pages: Customers (tab)
// @Tags         Auth
// @ID           loginAs
// @Accept       json
// @Produce      json
// @Security     cookieAuth
// @Param        request body dto.LoginAsRequest true "Request body"
// @Success      200 {object} response.Envelope{data=dto.LoginAsResponse}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed or account inactive"
// @Failure      404 {object} response.Envelope "Customer not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/auth/login-as [post]
func (authHandler *AuthHandler) LoginAs(writer http.ResponseWriter, request *http.Request) {
	actor, err := currentUser(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	var loginAsRequest dto.LoginAsRequest
	if err := decodeAndValidate(writer, request, &loginAsRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	session, err := authHandler.service.LoginAs(request.Context(), actor, loginAsRequest.CustomerID, authHandler.sessionClient(request))
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	// No MaxAge: the cookie ends with the browser session; the server also expires the session
	http.SetCookie(writer, authHandler.loginAsCookie(session.Key, session.Token))
	response.SendResponse(writer, request, response.Success, dto.LoginAsResponse{SessionKey: session.Key, User: session.User})
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

// sessionClient returns the user agent and client IP stored with a new session.
func (authHandler *AuthHandler) sessionClient(request *http.Request) dto.SessionClient {
	return dto.SessionClient{UserAgent: request.UserAgent(), IPAddress: middleware.ClientIP(request, authHandler.trustedProxies)}
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

// loginAsCookie returns the cookie of the "log in as" session key with value:
// the session cookie's attributes, scoped to that session's API path.
func (authHandler *AuthHandler) loginAsCookie(key, value string) *http.Cookie {
	cookie := authHandler.sessionCookie(value)
	cookie.Name = middleware.LoginAsCookieName
	cookie.Path = middleware.LoginAsPathPrefix + key
	return cookie
}
