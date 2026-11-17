package handler

import (
	"net/http"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/middleware"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/response"
	"7th-sense-server/internal/service"
)

// SessionHandler handles session management (server ADR 0014): the logged-in user's own sessions
// under /sessions and any user's under /admin (superadmin). Every revoke also closes the live
// sockets of the sessions it ended.
type SessionHandler struct {
	service *service.SessionService
	live    *LiveHandler // closes revoked sessions' live sockets; nil = live data off
}

// NewSessionHandler creates a SessionHandler that uses sessionService.
func NewSessionHandler(sessionService *service.SessionService) *SessionHandler {
	return &SessionHandler{service: sessionService}
}

// SetLiveHandler makes every revoke close the ended sessions' live sockets.
func (sessionHandler *SessionHandler) SetLiveHandler(liveHandler *LiveHandler) {
	sessionHandler.live = liveHandler
}

// GetMine handles GET /api/v1/sessions: the logged-in user's active sessions.
//
// @Summary      List my sessions
// @Description  Returns the logged-in user's active sessions, most recently used first; the session of the request has isCurrent true. "Log in as" sessions the superadmin opened on the account are left out (except when the request comes from one). No token or hash is returned; there is no location (no IP lookup).
// @Description
// @Description  Used by pages: My sessions
// @Tags         Sessions
// @ID           getMySessions
// @Produce      json
// @Security     cookieAuth
// @Success      200 {object} response.Envelope{data=[]dto.SessionResponse}
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/sessions [get]
func (sessionHandler *SessionHandler) GetMine(writer http.ResponseWriter, request *http.Request) {
	actor, err := currentUser(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	sessionID, _ := middleware.SessionID(request.Context())
	sessions, err := sessionHandler.service.GetMine(request.Context(), actor, sessionID)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, sessions)
}

// DeleteMine handles DELETE /api/v1/sessions/{id}: ends one of the logged-in user's sessions.
//
// @Summary      Log out one of my sessions
// @Description  Revokes the logged-in user's session id (reason user_revoked). That browser gets 401 session_revoked on its next request and its live sockets session:expired. A session of another user, a "log in as" session or one no longer active answers 404 (never 403), so ids cannot be probed.
// @Description
// @Description  Used by pages: My sessions
// @Tags         Sessions
// @ID           revokeMySession
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Session id"
// @Success      200 {object} response.Envelope
// @Failure      400 {object} response.Envelope "Please enter valid id"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      404 {object} response.Envelope "Session not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/sessions/{id} [delete]
func (sessionHandler *SessionHandler) DeleteMine(writer http.ResponseWriter, request *http.Request) {
	actor, err := currentUser(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	id, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	hash, err := sessionHandler.service.DeleteMine(request.Context(), actor, id)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	sessionHandler.closeLive([]string{hash}, model.RevokeReasonUserRevoked)
	response.SendResponse(writer, request, response.Success, "Session logged out")
}

// DeleteMyOthers handles DELETE /api/v1/sessions/others: ends all the logged-in user's sessions
// but the current one.
//
// @Summary      Log out all my other sessions
// @Description  Revokes every active session of the logged-in user except the session of the request (reason user_revoked); "log in as" sessions are not touched. Returns how many were revoked.
// @Description
// @Description  Used by pages: My sessions
// @Tags         Sessions
// @ID           revokeMyOtherSessions
// @Produce      json
// @Security     cookieAuth
// @Success      200 {object} response.Envelope{data=dto.RevokeSessionsResponse}
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/sessions/others [delete]
func (sessionHandler *SessionHandler) DeleteMyOthers(writer http.ResponseWriter, request *http.Request) {
	actor, err := currentUser(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	sessionID, _ := middleware.SessionID(request.Context())
	hashes, err := sessionHandler.service.DeleteMyOthers(request.Context(), actor, sessionID)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	sessionHandler.closeLive(hashes, model.RevokeReasonUserRevoked)
	response.SendResponse(writer, request, response.Success, dto.RevokeSessionsResponse{Revoked: len(hashes)})
}

// GetByUser handles GET /api/v1/admin/users/{id}/sessions: a user's active sessions (superadmin).
//
// @Summary      List a user's sessions
// @Description  Returns every active session of user id, most recently used first, "log in as" sessions included (isLoginAs).
// @Description  Roles: superadmin (others get 403).
// @Description
// @Description  Used by pages: User sessions
// @Tags         Sessions
// @ID           getUserSessions
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "User id"
// @Success      200 {object} response.Envelope{data=[]dto.SessionResponse}
// @Failure      400 {object} response.Envelope "Please enter valid id"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      404 {object} response.Envelope "User not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/admin/users/{id}/sessions [get]
func (sessionHandler *SessionHandler) GetByUser(writer http.ResponseWriter, request *http.Request) {
	actor, err := currentUser(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	userID, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	sessions, err := sessionHandler.service.GetByUser(request.Context(), actor, userID)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, sessions)
}

// DeleteByID handles DELETE /api/v1/admin/sessions/{id}: force-logs out one session (superadmin).
//
// @Summary      Force log out a session
// @Description  Revokes active session id of any user (revoked_by = the superadmin, reason admin_revoked). That browser gets 401 session_revoked and its live sockets session:expired.
// @Description  Roles: superadmin (others get 403).
// @Description
// @Description  Used by pages: User sessions
// @Tags         Sessions
// @ID           revokeUserSession
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Session id"
// @Success      200 {object} response.Envelope
// @Failure      400 {object} response.Envelope "Please enter valid id"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      404 {object} response.Envelope "Session not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/admin/sessions/{id} [delete]
func (sessionHandler *SessionHandler) DeleteByID(writer http.ResponseWriter, request *http.Request) {
	actor, err := currentUser(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	id, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	hash, err := sessionHandler.service.DeleteByID(request.Context(), actor, id)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	sessionHandler.closeLive([]string{hash}, model.RevokeReasonAdminRevoked)
	response.SendResponse(writer, request, response.Success, "Session logged out")
}

// DeleteByUser handles DELETE /api/v1/admin/users/{id}/sessions: force-logs out every session of
// a user (superadmin).
//
// @Summary      Force log out all of a user's sessions
// @Description  Revokes every active session of user id, "log in as" ones included (revoked_by = the superadmin, reason admin_revoked). Returns how many were revoked.
// @Description  Roles: superadmin (others get 403).
// @Description
// @Description  Used by pages: User sessions
// @Tags         Sessions
// @ID           revokeAllUserSessions
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "User id"
// @Success      200 {object} response.Envelope{data=dto.RevokeSessionsResponse}
// @Failure      400 {object} response.Envelope "Please enter valid id"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      404 {object} response.Envelope "User not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/admin/users/{id}/sessions [delete]
func (sessionHandler *SessionHandler) DeleteByUser(writer http.ResponseWriter, request *http.Request) {
	actor, err := currentUser(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	userID, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	hashes, err := sessionHandler.service.DeleteByUser(request.Context(), actor, userID)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	sessionHandler.closeLive(hashes, model.RevokeReasonAdminRevoked)
	response.SendResponse(writer, request, response.Success, dto.RevokeSessionsResponse{Revoked: len(hashes)})
}

// closeLive closes the live sockets of the revoked sessions, when live data is on.
func (sessionHandler *SessionHandler) closeLive(sessionHashes []string, revokeReason string) {
	if sessionHandler.live != nil {
		sessionHandler.live.CloseRevokedSessions(sessionHashes, revokeReason)
	}
}
