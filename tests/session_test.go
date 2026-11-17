package tests

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"7th-sense-server/internal/config"
	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/handler"
	"7th-sense-server/internal/middleware"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/routes"
	"7th-sense-server/internal/service"
)

// sessionTestUserAgent is the user agent the session tests log in with.
const sessionTestUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0 Safari/537.36"

// sessionTestRouter is the real router with the auth, session and user handlers and the real
// session middleware, with the session limits of one test.
type sessionTestRouter struct {
	router         http.Handler
	sessionService *service.SessionService
}

// newSessionRouter builds a sessionTestRouter with the limits of sessionConfig.
func newSessionRouter(sessionConfig config.SessionConfig) sessionTestRouter {
	authService := service.NewAuthService(testDBPool, sessionConfig)
	sessionService := service.NewSessionService(testDBPool, sessionConfig)
	router := routes.New(routes.Handlers{
		Auth:    handler.NewAuthHandler(authService, false),
		Session: handler.NewSessionHandler(sessionService),
		User:    newUserHandler(),
	}, middleware.Authenticate(authService))
	return sessionTestRouter{router: router, sessionService: sessionService}
}

// sessionLogin posts a login of username (forceLogoutOthers when force) and returns the answer
// and the session cookie (nil when none was set).
func sessionLogin(t *testing.T, router http.Handler, username string, force bool) (apiResponse, *http.Cookie) {
	t.Helper()
	result, recorder := authRequest(t, router, http.MethodPost, apiV1+"/auth/login",
		map[string]any{"username": username, "password": authTestPassword, "forceLogoutOthers": force}, nil,
		map[string]string{"User-Agent": sessionTestUserAgent})
	return result, sessionCookieOf(recorder)
}

// sessionIDOf returns the auth_tokens id of cookie's session.
func sessionIDOf(t *testing.T, cookie *http.Cookie) int64 {
	t.Helper()
	sum := sha256.Sum256([]byte(cookie.Value))
	var id int64
	if err := testDBPool.QueryRow(context.Background(), `SELECT id FROM auth_tokens WHERE token = $1`, hex.EncodeToString(sum[:])).Scan(&id); err != nil {
		t.Fatalf("find session: %v", err)
	}
	return id
}

// expectRevoked fails unless cookie's session answers 401 session_revoked with reason.
func expectRevoked(t *testing.T, router http.Handler, cookie *http.Cookie, reason string) {
	t.Helper()
	result, _ := authRequest(t, router, http.MethodGet, apiV1+"/auth/me", nil, cookie, nil)
	expectStatus(t, result, http.StatusUnauthorized, "session_revoked")
	var data struct {
		Reason string `json:"reason"`
	}
	decodeData(t, result, &data)
	if data.Reason != reason {
		t.Fatalf("revoke reason = %q; want %q", data.Reason, reason)
	}
}

// expectSignedIn fails unless cookie's session still authenticates.
func expectSignedIn(t *testing.T, router http.Handler, cookie *http.Cookie) {
	t.Helper()
	result, _ := authRequest(t, router, http.MethodGet, apiV1+"/auth/me", nil, cookie, nil)
	expectStatus(t, result, http.StatusOK, "success")
}

// activeSessionCount returns userID's active password-login sessions.
func activeSessionCount(t *testing.T, userID int64) int {
	t.Helper()
	var count int
	err := testDBPool.QueryRow(context.Background(), `SELECT COUNT(*) FROM auth_tokens
		WHERE user_id = $1 AND token_type = 'session' AND status = 'active' AND is_deleted = false
			AND expires_at > NOW() AND impersonated_by IS NULL`, userID).Scan(&count)
	if err != nil {
		t.Fatalf("count sessions: %v", err)
	}
	return count
}

// TestSessionSingleMode covers the single mode: a second login answers the conflict and keeps
// the first session; forceLogoutOthers revokes it, and that browser gets session_revoked.
func TestSessionSingleMode(t *testing.T) {
	sessions := newSessionRouter(config.SessionConfig{Mode: config.SessionModeSingle})
	username := "session_single_" + uniqueSuffix()
	userID := createAuthTestUser(t, username, model.RoleOperator)

	result, first := sessionLogin(t, sessions.router, username, false)
	expectStatus(t, result, http.StatusOK, "success")

	result, cookie := sessionLogin(t, sessions.router, username, false)
	expectStatus(t, result, http.StatusConflict, "session_conflict")
	if cookie != nil {
		t.Fatal("a conflicting login must not set a cookie")
	}
	var conflict dto.SessionConflictResponse
	decodeData(t, result, &conflict)
	if conflict.Mode != config.SessionModeSingle || len(conflict.Sessions) != 1 || conflict.Sessions[0].ID != sessionIDOf(t, first) {
		t.Fatalf("conflict = %+v; want single mode with the first session", conflict)
	}
	if existing := conflict.Sessions[0]; existing.Device != "Chrome on Windows" || existing.IPAddress == "" || existing.CreatedAt.IsZero() || existing.LastUsedAt.IsZero() {
		t.Fatalf("conflict session = %+v; want device, IP and times", existing)
	}
	if strings.Contains(string(result.Data), "token") {
		t.Fatalf("conflict leaks a token: %s", result.Data)
	}
	expectSignedIn(t, sessions.router, first)
	if count := activeSessionCount(t, userID); count != 1 {
		t.Fatalf("active sessions after conflict = %d; want 1", count)
	}

	t.Run("force logout needs the right password", func(t *testing.T) {
		result, _ := authRequest(t, sessions.router, http.MethodPost, apiV1+"/auth/login",
			map[string]any{"username": username, "password": "wrong-password", "forceLogoutOthers": true}, nil, nil)
		expectStatus(t, result, http.StatusUnauthorized, "Invalid username or password")
		expectSignedIn(t, sessions.router, first)
	})

	result, second := sessionLogin(t, sessions.router, username, true)
	expectStatus(t, result, http.StatusOK, "success")
	expectSignedIn(t, sessions.router, second)
	expectRevoked(t, sessions.router, first, model.RevokeReasonReplacedByNewLogin)

	var status, reason string
	var revokedBy int64
	var isDeleted bool
	err := testDBPool.QueryRow(context.Background(), `SELECT status, revoke_reason, revoked_by, is_deleted FROM auth_tokens WHERE id = $1`,
		sessionIDOf(t, first)).Scan(&status, &reason, &revokedBy, &isDeleted)
	if err != nil || status != model.SessionStatusRevoked || reason != model.RevokeReasonReplacedByNewLogin || revokedBy != userID || !isDeleted {
		t.Fatalf("revoked row = %s %s %d deleted=%v (%v)", status, reason, revokedBy, isDeleted, err)
	}

	t.Run("logout marks the session logged out", func(t *testing.T) {
		result, _ := authRequest(t, sessions.router, http.MethodPost, apiV1+"/auth/logout", nil, second, nil)
		expectStatus(t, result, http.StatusOK, "Logged out")
		var status string
		_ = testDBPool.QueryRow(context.Background(), `SELECT status FROM auth_tokens WHERE id = $1`, sessionIDOf(t, second)).Scan(&status)
		if status != model.SessionStatusLoggedOut {
			t.Fatalf("status after logout = %q", status)
		}
		result, _ = authRequest(t, sessions.router, http.MethodGet, apiV1+"/auth/me", nil, second, nil)
		expectStatus(t, result, http.StatusUnauthorized, "Please log in to continue")

		// After logout nothing is in the way any more
		result, _ = sessionLogin(t, sessions.router, username, false)
		expectStatus(t, result, http.StatusOK, "success")
	})
}

// TestSessionParallelLogins checks that concurrent logins of the same user in single mode cannot
// both win: the user's row lock serializes them, so exactly one creates a session.
func TestSessionParallelLogins(t *testing.T) {
	sessions := newSessionRouter(config.SessionConfig{Mode: config.SessionModeSingle})
	username := "session_race_" + uniqueSuffix()
	userID := createAuthTestUser(t, username, model.RoleOperator)

	const parallelLogins = 6
	statuses := make([]int, parallelLogins)
	var wait sync.WaitGroup
	for index := range parallelLogins {
		wait.Go(func() {
			request := httptest.NewRequest(http.MethodPost, apiV1+"/auth/login",
				strings.NewReader(`{"username":"`+username+`","password":"`+authTestPassword+`"}`))
			recorder := httptest.NewRecorder()
			sessions.router.ServeHTTP(recorder, request)
			statuses[index] = recorder.Code
		})
	}
	wait.Wait()

	succeeded, conflicts := 0, 0
	for _, status := range statuses {
		switch status {
		case http.StatusOK:
			succeeded++
		case http.StatusConflict:
			conflicts++
		}
	}
	if succeeded != 1 || conflicts != parallelLogins-1 {
		t.Fatalf("statuses = %v; want one 200 and %d 409", statuses, parallelLogins-1)
	}
	if count := activeSessionCount(t, userID); count != 1 {
		t.Fatalf("active sessions = %d; want 1", count)
	}

	t.Run("parallel forced logins leave one session", func(t *testing.T) {
		for index := range parallelLogins {
			wait.Go(func() {
				request := httptest.NewRequest(http.MethodPost, apiV1+"/auth/login",
					strings.NewReader(`{"username":"`+username+`","password":"`+authTestPassword+`","forceLogoutOthers":true}`))
				recorder := httptest.NewRecorder()
				sessions.router.ServeHTTP(recorder, request)
				statuses[index] = recorder.Code
			})
		}
		wait.Wait()
		for _, status := range statuses {
			if status != http.StatusOK {
				t.Fatalf("forced statuses = %v; want all 200", statuses)
			}
		}
		if count := activeSessionCount(t, userID); count != 1 {
			t.Fatalf("active sessions after forced logins = %d; want 1", count)
		}
	})
}

// TestSessionMultipleMode covers multiple mode: unlimited logins, and with MAX_SESSIONS the
// conflict listing every session and forceLogoutOthers revoking only the least recently used.
func TestSessionMultipleMode(t *testing.T) {
	t.Run("unlimited", func(t *testing.T) {
		sessions := newSessionRouter(config.SessionConfig{Mode: config.SessionModeMultiple})
		username := "session_multi_" + uniqueSuffix()
		userID := createAuthTestUser(t, username, model.RoleOperator)
		for range 4 {
			result, _ := sessionLogin(t, sessions.router, username, false)
			expectStatus(t, result, http.StatusOK, "success")
		}
		if count := activeSessionCount(t, userID); count != 4 {
			t.Fatalf("active sessions = %d; want 4", count)
		}
	})

	t.Run("max sessions", func(t *testing.T) {
		sessions := newSessionRouter(config.SessionConfig{Mode: config.SessionModeMultiple, MaxSessions: 2})
		username := "session_max_" + uniqueSuffix()
		createAuthTestUser(t, username, model.RoleOperator)
		_, oldest := sessionLogin(t, sessions.router, username, false)
		_, newer := sessionLogin(t, sessions.router, username, false)
		// Make the order of use certain: the first session was used an hour ago.
		if _, err := testDBPool.Exec(context.Background(), `UPDATE auth_tokens SET last_used_at = NOW() - INTERVAL '1 hour' WHERE id = $1`, sessionIDOf(t, oldest)); err != nil {
			t.Fatalf("age session: %v", err)
		}

		result, cookie := sessionLogin(t, sessions.router, username, false)
		expectStatus(t, result, http.StatusConflict, "session_conflict")
		var conflict dto.SessionConflictResponse
		decodeData(t, result, &conflict)
		if cookie != nil || conflict.MaxSessions != 2 || len(conflict.Sessions) != 2 || conflict.Sessions[0].ID != sessionIDOf(t, newer) {
			t.Fatalf("conflict = %+v; want both sessions, most recently used first", conflict)
		}

		result, newest := sessionLogin(t, sessions.router, username, true)
		expectStatus(t, result, http.StatusOK, "success")
		expectRevoked(t, sessions.router, oldest, model.RevokeReasonReplacedByNewLogin)
		expectSignedIn(t, sessions.router, newer)
		expectSignedIn(t, sessions.router, newest)
	})
}

// TestMySessions covers a user's own sessions: the list marks the current one and never shows a
// token, another user's session is 404, and a user can end one session or all the others.
func TestMySessions(t *testing.T) {
	sessions := newSessionRouter(config.SessionConfig{Mode: config.SessionModeMultiple})
	suffix := uniqueSuffix()
	username := "session_mine_" + suffix
	otherName := "session_other_" + suffix
	createAuthTestUser(t, username, model.RoleOperator)
	createAuthTestUser(t, otherName, model.RoleOperator)
	_, current := sessionLogin(t, sessions.router, username, false)
	_, laptop := sessionLogin(t, sessions.router, username, false)
	_, phone := sessionLogin(t, sessions.router, username, false)
	_, otherUser := sessionLogin(t, sessions.router, otherName, false)

	result, _ := authRequest(t, sessions.router, http.MethodGet, apiV1+"/sessions", nil, current, nil)
	expectStatus(t, result, http.StatusOK, "success")
	var mine []dto.SessionResponse
	decodeData(t, result, &mine)
	currentCount := 0
	for _, session := range mine {
		if session.IsCurrent {
			currentCount++
			if session.ID != sessionIDOf(t, current) {
				t.Fatalf("current session = %d; want %d", session.ID, sessionIDOf(t, current))
			}
		}
	}
	if len(mine) != 3 || currentCount != 1 {
		t.Fatalf("my sessions = %+v; want 3 with one current", mine)
	}
	if strings.Contains(string(result.Data), `"token`) {
		t.Fatalf("session list leaks a token: %s", result.Data)
	}

	t.Run("another user's session is not found", func(t *testing.T) {
		otherID := strconv.FormatInt(sessionIDOf(t, otherUser), 10)
		result, _ := authRequest(t, sessions.router, http.MethodDelete, apiV1+"/sessions/"+otherID, nil, current, nil)
		expectStatus(t, result, http.StatusNotFound, "Session not found")
		expectSignedIn(t, sessions.router, otherUser)
	})

	t.Run("log out one session", func(t *testing.T) {
		laptopID := strconv.FormatInt(sessionIDOf(t, laptop), 10)
		result, _ := authRequest(t, sessions.router, http.MethodDelete, apiV1+"/sessions/"+laptopID, nil, current, nil)
		expectStatus(t, result, http.StatusOK, "Session logged out")
		expectRevoked(t, sessions.router, laptop, model.RevokeReasonUserRevoked)

		result, _ = authRequest(t, sessions.router, http.MethodDelete, apiV1+"/sessions/"+laptopID, nil, current, nil)
		expectStatus(t, result, http.StatusNotFound, "Session not found")
	})

	t.Run("log out all other sessions", func(t *testing.T) {
		result, _ := authRequest(t, sessions.router, http.MethodDelete, apiV1+"/sessions/others", nil, current, nil)
		expectStatus(t, result, http.StatusOK, "success")
		var revoked dto.RevokeSessionsResponse
		decodeData(t, result, &revoked)
		if revoked.Revoked != 1 {
			t.Fatalf("revoked = %d; want 1 (the phone)", revoked.Revoked)
		}
		expectRevoked(t, sessions.router, phone, model.RevokeReasonUserRevoked)
		expectSignedIn(t, sessions.router, current)
		expectSignedIn(t, sessions.router, otherUser)
	})
}

// TestAdminSessions covers the superadmin's session endpoints and that nobody else may use them.
func TestAdminSessions(t *testing.T) {
	sessions := newSessionRouter(config.SessionConfig{Mode: config.SessionModeMultiple})
	suffix := uniqueSuffix()
	superName := "session_super_" + suffix
	customerName := "session_customer_" + suffix
	superID := createAuthTestUser(t, superName, model.RoleSuperAdmin)
	customerID := createAuthTestUser(t, customerName, model.RoleCustomer)
	superSession := login(t, sessions.router, superName, false)
	_, first := sessionLogin(t, sessions.router, customerName, false)
	_, second := sessionLogin(t, sessions.router, customerName, false)
	_, third := sessionLogin(t, sessions.router, customerName, false)
	userPath := apiV1 + "/admin/users/" + strconv.FormatInt(customerID, 10) + "/sessions"

	t.Run("others get 403", func(t *testing.T) {
		firstPath := apiV1 + "/admin/sessions/" + strconv.FormatInt(sessionIDOf(t, second), 10)
		for _, call := range []struct{ method, path string }{
			{http.MethodGet, userPath}, {http.MethodDelete, userPath}, {http.MethodDelete, firstPath},
		} {
			result, _ := authRequest(t, sessions.router, call.method, call.path, nil, first, nil)
			expectStatus(t, result, http.StatusForbidden, "You do not have permission to do this")
		}
	})

	t.Run("list a user's sessions", func(t *testing.T) {
		result, _ := authRequest(t, sessions.router, http.MethodGet, userPath, nil, superSession, nil)
		expectStatus(t, result, http.StatusOK, "success")
		var list []dto.SessionResponse
		decodeData(t, result, &list)
		if len(list) != 3 {
			t.Fatalf("user sessions = %d; want 3", len(list))
		}
		result, _ = authRequest(t, sessions.router, http.MethodGet, apiV1+"/admin/users/999999999/sessions", nil, superSession, nil)
		expectStatus(t, result, http.StatusNotFound, "User not found")
	})

	t.Run("force log out one session", func(t *testing.T) {
		path := apiV1 + "/admin/sessions/" + strconv.FormatInt(sessionIDOf(t, first), 10)
		result, _ := authRequest(t, sessions.router, http.MethodDelete, path, nil, superSession, nil)
		expectStatus(t, result, http.StatusOK, "Session logged out")
		expectRevoked(t, sessions.router, first, model.RevokeReasonAdminRevoked)
		var revokedBy int64
		_ = testDBPool.QueryRow(context.Background(), `SELECT revoked_by FROM auth_tokens WHERE id = $1`, sessionIDOf(t, first)).Scan(&revokedBy)
		if revokedBy != superID {
			t.Fatalf("revoked_by = %d; want the superadmin %d", revokedBy, superID)
		}
		result, _ = authRequest(t, sessions.router, http.MethodDelete, path, nil, superSession, nil)
		expectStatus(t, result, http.StatusNotFound, "Session not found")
	})

	t.Run("force log out all of a user's sessions", func(t *testing.T) {
		result, _ := authRequest(t, sessions.router, http.MethodDelete, userPath, nil, superSession, nil)
		expectStatus(t, result, http.StatusOK, "success")
		var revoked dto.RevokeSessionsResponse
		decodeData(t, result, &revoked)
		if revoked.Revoked != 2 {
			t.Fatalf("revoked = %d; want 2", revoked.Revoked)
		}
		expectRevoked(t, sessions.router, second, model.RevokeReasonAdminRevoked)
		expectRevoked(t, sessions.router, third, model.RevokeReasonAdminRevoked)
		expectSignedIn(t, sessions.router, superSession)
	})
}

// TestSessionLoginAsDoesNotCount checks that the superadmin's "log in as" session never counts
// toward the customer's limit, is not in the customer's own list and leaves the superadmin signed in.
func TestSessionLoginAsDoesNotCount(t *testing.T) {
	sessions := newSessionRouter(config.SessionConfig{Mode: config.SessionModeSingle})
	suffix := uniqueSuffix()
	superName := "session_lasuper_" + suffix
	customerName := "session_lacustomer_" + suffix
	createAuthTestUser(t, superName, model.RoleSuperAdmin)
	customerID := createAuthTestUser(t, customerName, model.RoleCustomer)
	superSession := login(t, sessions.router, superName, false)

	result, recorder := authRequest(t, sessions.router, http.MethodPost, apiV1+"/auth/login-as", map[string]any{"customerId": customerID}, superSession, nil)
	expectStatus(t, result, http.StatusOK, "success")
	var started dto.LoginAsResponse
	decodeData(t, result, &started)
	asCookie := loginAsCookieOf(recorder)
	asV1 := middleware.LoginAsPathPrefix + started.SessionKey + "/v1"

	result, customerSession := sessionLogin(t, sessions.router, customerName, false)
	expectStatus(t, result, http.StatusOK, "success")

	result, _ = authRequest(t, sessions.router, http.MethodGet, apiV1+"/sessions", nil, customerSession, nil)
	var mine []dto.SessionResponse
	decodeData(t, result, &mine)
	if len(mine) != 1 || mine[0].IsLoginAs {
		t.Fatalf("customer's own sessions = %+v; want only the password login", mine)
	}

	// The customer forcing a new login ends only their own session, not the superadmin's view.
	result, _ = sessionLogin(t, sessions.router, customerName, true)
	expectStatus(t, result, http.StatusOK, "success")
	result, _ = authRequest(t, sessions.router, http.MethodGet, asV1+"/auth/me", nil, asCookie, nil)
	expectStatus(t, result, http.StatusOK, "success")
	expectSignedIn(t, sessions.router, superSession)
}

// TestSessionExpiry covers an expired session (401 session_expired), the throttled last_used_at
// write and the cleanup that marks overdue sessions expired and purges old ended ones.
func TestSessionExpiry(t *testing.T) {
	sessionConfig := config.SessionConfig{Mode: config.SessionModeMultiple, Retention: 30 * 24 * time.Hour}
	sessions := newSessionRouter(sessionConfig)
	username := "session_expiry_" + uniqueSuffix()
	createAuthTestUser(t, username, model.RoleOperator)
	ctx := context.Background()

	t.Run("last used is written at most once a minute", func(t *testing.T) {
		cookie := login(t, sessions.router, username, false)
		id := sessionIDOf(t, cookie)
		lastUsed := func() time.Time {
			var value time.Time
			_ = testDBPool.QueryRow(ctx, `SELECT last_used_at FROM auth_tokens WHERE id = $1`, id).Scan(&value)
			return value
		}
		_, _ = testDBPool.Exec(ctx, `UPDATE auth_tokens SET last_used_at = NOW() - INTERVAL '10 seconds' WHERE id = $1`, id)
		before := lastUsed()
		expectSignedIn(t, sessions.router, cookie)
		if !lastUsed().Equal(before) {
			t.Fatal("last_used_at was written within a minute of the last write")
		}
		_, _ = testDBPool.Exec(ctx, `UPDATE auth_tokens SET last_used_at = NOW() - INTERVAL '2 minutes' WHERE id = $1`, id)
		before = lastUsed()
		expectSignedIn(t, sessions.router, cookie)
		if !lastUsed().After(before) {
			t.Fatal("last_used_at was not written after a minute")
		}
	})

	t.Run("expired session and cleanup", func(t *testing.T) {
		cookie := login(t, sessions.router, username, false)
		id := sessionIDOf(t, cookie)
		if _, err := testDBPool.Exec(ctx, `UPDATE auth_tokens SET expires_at = NOW() - INTERVAL '1 minute' WHERE id = $1`, id); err != nil {
			t.Fatalf("expire session: %v", err)
		}
		result, _ := authRequest(t, sessions.router, http.MethodGet, apiV1+"/auth/me", nil, cookie, nil)
		expectStatus(t, result, http.StatusUnauthorized, "session_expired")

		old := login(t, sessions.router, username, false)
		oldID := sessionIDOf(t, old)
		if _, err := testDBPool.Exec(ctx, `UPDATE auth_tokens SET status = 'logged_out', is_deleted = true,
			deleted_at = NOW() - INTERVAL '40 days' WHERE id = $1`, oldID); err != nil {
			t.Fatalf("age session: %v", err)
		}

		sessions.sessionService.CleanUp(ctx)

		var status string
		_ = testDBPool.QueryRow(ctx, `SELECT status FROM auth_tokens WHERE id = $1`, id).Scan(&status)
		if status != model.SessionStatusExpired {
			t.Fatalf("status after cleanup = %q; want expired", status)
		}
		result, _ = authRequest(t, sessions.router, http.MethodGet, apiV1+"/auth/me", nil, cookie, nil)
		expectStatus(t, result, http.StatusUnauthorized, "session_expired")
		var remaining int
		_ = testDBPool.QueryRow(ctx, `SELECT COUNT(*) FROM auth_tokens WHERE id = $1`, oldID).Scan(&remaining)
		if remaining != 0 {
			t.Fatal("a session ended 40 days ago was not purged")
		}
	})
}

// TestSessionRevokeClosesSockets checks that revoking a session closes its live sockets with
// session:expired and the reason: a forced login in single mode and the superadmin's revoke.
func TestSessionRevokeClosesSockets(t *testing.T) {
	server := newLiveTestServerWith(t, config.SessionConfig{Mode: config.SessionModeSingle})
	suffix := uniqueSuffix()
	username := "session_socket_" + suffix
	superName := "session_socketsuper_" + suffix
	createAuthTestUser(t, username, model.RoleOperator)
	createAuthTestUser(t, superName, model.RoleSuperAdmin)

	// dialRegistered opens a socket and waits for a subscribe acknowledgement: the browser hears
	// "connect" before the server has registered the socket, and only a registered one is closed.
	dialRegistered := func(cookie *http.Cookie) *liveTestClient {
		t.Helper()
		client, status := dialLive(t, server, cookie, "")
		if status != http.StatusSwitchingProtocols {
			t.Fatalf("handshake = %d", status)
		}
		var answer subscribeAnswer
		client.call("subscribe", map[string]any{"machineIds": []int64{}}, &answer)
		return client
	}
	expectEnded := func(client *liveTestClient, revokeReason string) {
		t.Helper()
		var ended dto.LiveSessionExpired
		if err := json.Unmarshal(client.expectEvent("session:expired").payload, &ended); err != nil {
			t.Fatalf("decode session:expired: %v", err)
		}
		if ended.Reason != dto.LiveSessionReasonRevoked || ended.RevokeReason != revokeReason {
			t.Fatalf("session:expired = %+v; want session_revoked / %s", ended, revokeReason)
		}
		client.expectClosed()
	}

	t.Run("forced login", func(t *testing.T) {
		_, first := sessionLogin(t, server.router, username, false)
		client := dialRegistered(first)
		result, _ := sessionLogin(t, server.router, username, true)
		expectStatus(t, result, http.StatusOK, "success")
		expectEnded(client, model.RevokeReasonReplacedByNewLogin)
	})

	t.Run("superadmin revoke", func(t *testing.T) {
		_, current := sessionLogin(t, server.router, username, true)
		client := dialRegistered(current)
		superSession := login(t, server.router, superName, false)
		path := apiV1 + "/admin/sessions/" + strconv.FormatInt(sessionIDOf(t, current), 10)
		result, _ := authRequest(t, server.router, http.MethodDelete, path, nil, superSession, nil)
		expectStatus(t, result, http.StatusOK, "Session logged out")
		expectEnded(client, model.RevokeReasonAdminRevoked)
	})
}

// TestClientIP checks that X-Forwarded-For is believed only from a trusted proxy.
func TestClientIP(t *testing.T) {
	trusted := []netip.Prefix{netip.MustParsePrefix("10.0.0.0/8")}
	cases := []struct {
		name, remote, forwarded string
		trusted                 []netip.Prefix
		want                    string
	}{
		{"no proxy configured ignores the header", "203.0.113.7:5000", "198.51.100.1", nil, "203.0.113.7"},
		{"untrusted peer ignores the header", "203.0.113.7:5000", "198.51.100.1", trusted, "203.0.113.7"},
		{"trusted proxy gives the client", "10.0.0.2:5000", "198.51.100.1", trusted, "198.51.100.1"},
		{"forged left entry is skipped", "10.0.0.2:5000", "1.2.3.4, 198.51.100.1, 10.0.0.3", trusted, "198.51.100.1"},
		{"IPv6 peer", "[2001:db8::1]:443", "", nil, "2001:db8::1"},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			request.RemoteAddr = testCase.remote
			if testCase.forwarded != "" {
				request.Header.Set("X-Forwarded-For", testCase.forwarded)
			}
			if got := middleware.ClientIP(request, testCase.trusted); got != testCase.want {
				t.Fatalf("ClientIP = %q; want %q", got, testCase.want)
			}
		})
	}
}
