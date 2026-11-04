package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"7th-sense-server/internal/middleware"
	"7th-sense-server/internal/model"
)

// loginAsCookieOf returns the "log in as" cookie set by recorder, or nil.
func loginAsCookieOf(recorder *httptest.ResponseRecorder) *http.Cookie {
	for _, cookie := range recorder.Result().Cookies() {
		if cookie.Name == middleware.LoginAsCookieName {
			return cookie
		}
	}
	return nil
}

// TestLoginAs covers the superadmin opening a customer's account in another
// tab: the customer session has its own path-scoped cookie, so logging out of
// either session never ends the other.
func TestLoginAs(t *testing.T) {
	router := newAuthRouter()
	suffix := uniqueSuffix()
	superAdminName := "loginas_super_" + suffix
	customerName := "loginas_customer_" + suffix
	createAuthTestUser(t, superAdminName, model.RoleSuperAdmin)
	customerID := createAuthTestUser(t, customerName, model.RoleCustomer)
	operatorID := createAuthTestUser(t, "loginas_operator_"+suffix, model.RoleOperator)

	superSession := login(t, router, superAdminName, false)

	t.Run("only the superadmin, only customers", func(t *testing.T) {
		customerSession := login(t, router, customerName, false)
		result, _ := authRequest(t, router, http.MethodPost, apiV1+"/auth/login-as", map[string]any{"customerId": customerID}, customerSession, nil)
		expectStatus(t, result, http.StatusForbidden, "You do not have permission to do this")

		result, _ = authRequest(t, router, http.MethodPost, apiV1+"/auth/login-as", map[string]any{}, superSession, nil)
		expectStatus(t, result, http.StatusBadRequest, "Please select customer")

		result, _ = authRequest(t, router, http.MethodPost, apiV1+"/auth/login-as", map[string]any{"customerId": operatorID}, superSession, nil)
		expectStatus(t, result, http.StatusNotFound, "Customer not found")
	})

	result, recorder := authRequest(t, router, http.MethodPost, apiV1+"/auth/login-as", map[string]any{"customerId": customerID}, superSession, nil)
	expectStatus(t, result, http.StatusOK, "success")
	var started struct {
		SessionKey string         `json:"sessionKey"`
		User       model.AuthUser `json:"user"`
	}
	if err := json.Unmarshal(result.Data, &started); err != nil {
		t.Fatalf("decode login-as data: %v", err)
	}
	asCookie := loginAsCookieOf(recorder)
	if asCookie == nil || asCookie.Value == "" || !asCookie.HttpOnly {
		t.Fatalf("login-as cookie = %+v; want an HttpOnly cookie", asCookie)
	}
	asPrefix := middleware.LoginAsPathPrefix + started.SessionKey
	if asCookie.Path != asPrefix || started.User.ID != customerID {
		t.Fatalf("cookie path %q, user %d; want %q, %d", asCookie.Path, started.User.ID, asPrefix, customerID)
	}
	if sessionCookieOf(recorder) != nil {
		t.Fatal("login-as must not touch the superadmin's session cookie")
	}
	asV1 := asPrefix + "/v1"

	t.Run("the tab works as the customer", func(t *testing.T) {
		result, _ := authRequest(t, router, http.MethodGet, asV1+"/auth/me", nil, asCookie, nil)
		expectStatus(t, result, http.StatusOK, "success")
		var me model.AuthUser
		if err := json.Unmarshal(result.Data, &me); err != nil || me.ID != customerID || me.Role != model.RoleCustomer {
			t.Fatalf("me = %+v (%v); want customer %d", me, err, customerID)
		}
	})

	t.Run("no fallback to the superadmin's cookie", func(t *testing.T) {
		result, _ := authRequest(t, router, http.MethodGet, asV1+"/auth/me", nil, superSession, nil)
		expectStatus(t, result, http.StatusUnauthorized, "Please log in to continue")
	})

	t.Run("password login is refused in the tab", func(t *testing.T) {
		result, _ := authRequest(t, router, http.MethodPost, asV1+"/auth/login",
			map[string]any{"username": superAdminName, "password": authTestPassword}, nil, nil)
		expectStatus(t, result, http.StatusForbidden, "You do not have permission to do this")
	})

	t.Run("malformed key", func(t *testing.T) {
		result, _ := authRequest(t, router, http.MethodGet, middleware.LoginAsPathPrefix+"not-a-key/v1/auth/me", nil, asCookie, nil)
		if result.StatusCode != http.StatusNotFound {
			t.Fatalf("got %d; want 404", result.StatusCode)
		}
	})

	t.Run("superadmin logout keeps the customer tab", func(t *testing.T) {
		otherSuper := login(t, router, superAdminName, false)
		_, recorder := authRequest(t, router, http.MethodPost, apiV1+"/auth/logout", nil, otherSuper, nil)
		if loginAsCookieOf(recorder) != nil {
			t.Fatal("superadmin logout must not clear the login-as cookie")
		}
		result, _ := authRequest(t, router, http.MethodGet, asV1+"/auth/me", nil, asCookie, nil)
		expectStatus(t, result, http.StatusOK, "success")
	})

	t.Run("customer tab logout keeps the superadmin", func(t *testing.T) {
		_, recorder := authRequest(t, router, http.MethodPost, asV1+"/auth/logout", nil, asCookie, nil)
		cleared := loginAsCookieOf(recorder)
		if cleared == nil || cleared.MaxAge >= 0 || cleared.Path != asPrefix || sessionCookieOf(recorder) != nil {
			t.Fatalf("logout cleared %+v; want only the login-as cookie at %q", cleared, asPrefix)
		}
		result, _ := authRequest(t, router, http.MethodGet, asV1+"/auth/me", nil, asCookie, nil)
		expectStatus(t, result, http.StatusUnauthorized, "Please log in to continue")

		result, _ = authRequest(t, router, http.MethodGet, apiV1+"/auth/me", nil, superSession, nil)
		expectStatus(t, result, http.StatusOK, "success")
		if !strings.Contains(string(result.Data), superAdminName) {
			t.Fatalf("me = %s; want the superadmin", result.Data)
		}
	})
}
