package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"7th-sense-server/internal/handler"
	"7th-sense-server/internal/middleware"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/routes"
	"7th-sense-server/internal/service"
)

const authTestPassword = "Secret#123"

// newAuthRouter builds the real router with the real session middleware, the
// auth handler and the (admin-only) user handler.
func newAuthRouter() http.Handler {
	userCountryService := service.NewUserCountryService(repository.NewUserCountryRepository(testDBPool))
	userStateService := service.NewUserStateService(repository.NewUserStateRepository(testDBPool), userCountryService)
	userCityService := service.NewUserCityService(repository.NewUserCityRepository(testDBPool), userStateService)
	branchService := service.NewBranchService(repository.NewBranchRepository(testDBPool), userCityService)
	authService := service.NewAuthService(repository.NewAuthRepository(testDBPool))
	return routes.New(routes.Handlers{
		Auth: handler.NewAuthHandler(authService, false),
		User: handler.NewUserHandler(service.NewUserService(testDBPool, branchService)),
	}, middleware.Authenticate(authService))
}

// createAuthTestUser inserts a user with authTestPassword and role, and
// soft-deletes it when the test ends.
func createAuthTestUser(t *testing.T, username, role string) int64 {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(authTestPassword), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	var id int64
	err = testDBPool.QueryRow(context.Background(),
		`INSERT INTO users (username, email, password, role_id, status)
		VALUES ($1, $2, $3, (SELECT id FROM roles WHERE name = $4), 'active')
		RETURNING id`,
		username, username+"@auth-test.example.com", string(hash), role,
	).Scan(&id)
	if err != nil {
		t.Fatalf("create user %s: %v", username, err)
	}
	t.Cleanup(func() {
		_, _ = testDBPool.Exec(context.Background(), `UPDATE users SET is_deleted = true, deleted_at = NOW() WHERE id = $1`, id)
	})
	return id
}

// authRequest sends method path with body as JSON and the session cookie (if
// any) and returns the decoded envelope and the raw recorder.
func authRequest(t *testing.T, router http.Handler, method, path string, body any, session *http.Cookie, headers map[string]string) (apiResponse, *httptest.ResponseRecorder) {
	t.Helper()
	var requestBody bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&requestBody).Encode(body); err != nil {
			t.Fatalf("encode request body: %v", err)
		}
	}
	request := httptest.NewRequest(method, path, &requestBody)
	if session != nil {
		request.AddCookie(session)
	}
	for name, value := range headers {
		request.Header.Set(name, value)
	}
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	var result apiResponse
	if err := json.NewDecoder(bytes.NewReader(recorder.Body.Bytes())).Decode(&result); err != nil {
		t.Fatalf("%s %s: decode response %q: %v", method, path, recorder.Body.String(), err)
	}
	return result, recorder
}

// sessionCookieOf returns the session cookie set by recorder, or nil.
func sessionCookieOf(recorder *httptest.ResponseRecorder) *http.Cookie {
	for _, cookie := range recorder.Result().Cookies() {
		if cookie.Name == middleware.SessionCookieName {
			return cookie
		}
	}
	return nil
}

// login logs username in and returns its session cookie.
func login(t *testing.T, router http.Handler, username string, remember bool) *http.Cookie {
	t.Helper()
	result, recorder := authRequest(t, router, http.MethodPost, apiV1+"/auth/login",
		map[string]any{"username": username, "password": authTestPassword, "remember": remember}, nil, nil)
	expectStatus(t, result, http.StatusOK, "success")
	cookie := sessionCookieOf(recorder)
	if cookie == nil || cookie.Value == "" {
		t.Fatalf("login %s: no session cookie", username)
	}
	return cookie
}

// TestAuthentication covers login, session use, role checks, CSRF, logout,
// deactivation and login throttling end to end.
func TestAuthentication(t *testing.T) {
	router := newAuthRouter()
	suffix := uniqueSuffix()
	adminName := "auth_admin_" + suffix
	customerName := "auth_customer_" + suffix
	adminID := createAuthTestUser(t, adminName, model.RoleAdmin)
	createAuthTestUser(t, customerName, model.RoleCustomer)

	t.Run("login validation and wrong credentials", func(t *testing.T) {
		result, _ := authRequest(t, router, http.MethodPost, apiV1+"/auth/login", map[string]any{"password": "x"}, nil, nil)
		expectStatus(t, result, http.StatusBadRequest, "Please enter username")

		result, _ = authRequest(t, router, http.MethodPost, apiV1+"/auth/login", map[string]any{"username": adminName, "password": "wrong-password"}, nil, nil)
		expectStatus(t, result, http.StatusUnauthorized, "Invalid username or password")

		result, _ = authRequest(t, router, http.MethodPost, apiV1+"/auth/login", map[string]any{"username": "nobody_" + suffix, "password": "wrong-password"}, nil, nil)
		expectStatus(t, result, http.StatusUnauthorized, "Invalid username or password")
	})

	t.Run("protected routes need a session", func(t *testing.T) {
		for _, path := range []string{"/auth/me", "/me", "/users"} {
			result, _ := authRequest(t, router, http.MethodGet, apiV1+path, nil, nil, nil)
			expectStatus(t, result, http.StatusUnauthorized, "Please log in to continue")
		}
		forged := &http.Cookie{Name: middleware.SessionCookieName, Value: "forged-token"}
		result, _ := authRequest(t, router, http.MethodGet, apiV1+"/auth/me", nil, forged, nil)
		expectStatus(t, result, http.StatusUnauthorized, "Please log in to continue")
	})

	customerSession := login(t, router, customerName, false)
	adminSession := login(t, router, adminName, true)

	t.Run("cookie attributes", func(t *testing.T) {
		if !customerSession.HttpOnly || customerSession.SameSite != http.SameSiteStrictMode || customerSession.Path != "/api" {
			t.Fatalf("session cookie = %+v; want HttpOnly, SameSite=Strict, Path=/api", customerSession)
		}
		if customerSession.MaxAge != 0 {
			t.Fatalf("not remembered: MaxAge = %d; want browser-session cookie", customerSession.MaxAge)
		}
		if adminSession.MaxAge <= 0 {
			t.Fatalf("remembered: MaxAge = %d; want persistent cookie", adminSession.MaxAge)
		}
	})

	t.Run("me returns the logged-in user with role", func(t *testing.T) {
		result, _ := authRequest(t, router, http.MethodGet, apiV1+"/auth/me", nil, customerSession, nil)
		expectStatus(t, result, http.StatusOK, "success")
		var user model.AuthUser
		decodeData(t, result, &user)
		if user.Username != customerName || user.Role != model.RoleCustomer {
			t.Fatalf("me = %+v; want %s with role customer", user, customerName)
		}
		if bytes.Contains(result.Data, []byte("password")) {
			t.Fatalf("me leaks the password field: %s", result.Data)
		}
	})

	t.Run("admin-only routes", func(t *testing.T) {
		result, _ := authRequest(t, router, http.MethodGet, apiV1+"/users", nil, customerSession, nil)
		expectStatus(t, result, http.StatusForbidden, "You do not have permission to do this")

		result, _ = authRequest(t, router, http.MethodGet, apiV1+"/users", nil, adminSession, nil)
		expectStatus(t, result, http.StatusOK, "success")
	})

	t.Run("cross-site writes are rejected", func(t *testing.T) {
		result, _ := authRequest(t, router, http.MethodDelete, apiV1+"/users/999999999", nil, adminSession,
			map[string]string{"Sec-Fetch-Site": "cross-site"})
		expectStatus(t, result, http.StatusForbidden, "You do not have permission to do this")
	})

	t.Run("logout ends the session", func(t *testing.T) {
		result, recorder := authRequest(t, router, http.MethodPost, apiV1+"/auth/logout", nil, customerSession, nil)
		expectStatus(t, result, http.StatusOK, "Logged out")
		if cleared := sessionCookieOf(recorder); cleared == nil || cleared.MaxAge >= 0 {
			t.Fatalf("logout cookie = %+v; want deleted cookie", cleared)
		}
		result, _ = authRequest(t, router, http.MethodGet, apiV1+"/auth/me", nil, customerSession, nil)
		expectStatus(t, result, http.StatusUnauthorized, "Please log in to continue")

		// Idempotent
		result, _ = authRequest(t, router, http.MethodPost, apiV1+"/auth/logout", nil, customerSession, nil)
		expectStatus(t, result, http.StatusOK, "Logged out")
	})

	t.Run("deactivated user loses session and cannot log in", func(t *testing.T) {
		if _, err := testDBPool.Exec(context.Background(), `UPDATE users SET status = 'inactive' WHERE id = $1`, adminID); err != nil {
			t.Fatalf("deactivate: %v", err)
		}
		result, _ := authRequest(t, router, http.MethodGet, apiV1+"/auth/me", nil, adminSession, nil)
		expectStatus(t, result, http.StatusUnauthorized, "Please log in to continue")

		result, _ = authRequest(t, router, http.MethodPost, apiV1+"/auth/login",
			map[string]any{"username": adminName, "password": authTestPassword}, nil, nil)
		expectStatus(t, result, http.StatusForbidden, "Your account is inactive. Please contact your administrator")
	})

	t.Run("too many wrong passwords lock the username", func(t *testing.T) {
		for range 5 {
			result, _ := authRequest(t, router, http.MethodPost, apiV1+"/auth/login",
				map[string]any{"username": customerName, "password": "wrong-password"}, nil, nil)
			expectStatus(t, result, http.StatusUnauthorized, "Invalid username or password")
		}
		result, _ := authRequest(t, router, http.MethodPost, apiV1+"/auth/login",
			map[string]any{"username": customerName, "password": authTestPassword}, nil, nil)
		expectStatus(t, result, http.StatusTooManyRequests, "Too many login attempts. Please try again in 15 minutes")
	})
}
