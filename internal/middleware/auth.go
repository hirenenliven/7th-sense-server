package middleware

import (
	"context"
	"encoding/hex"
	"net/http"
	"slices"
	"strings"

	"7th-sense-server/internal/model"
	"7th-sense-server/internal/response"
	"7th-sense-server/internal/service"
)

// SessionCookieName is the cookie that carries the session token. It is
// HttpOnly, so page scripts can never read it.
const SessionCookieName = "ts_session"

// LoginAsCookieName is the cookie of a "log in as" session (the superadmin
// viewing a customer's account in another tab). It is HttpOnly too, and scoped
// to that session's path (LoginAsPathPrefix + key), so the browser sends it
// only with that tab's requests and it never replaces SessionCookieName.
const LoginAsCookieName = "ts_login_as"

// LoginAsPathPrefix starts the API path of a "log in as" session:
// /api/as/<key>/v1/... serves the same routes as /api/v1/...
const LoginAsPathPrefix = "/api/as/"

// loginAsKeyLength is the length of a "log in as" key (16 random bytes as hex).
const loginAsKeyLength = 32

type authUserKey struct{}

type loginAsKey struct{}

// RouteLoginAs serves /api/as/<key>/... as /api/... and marks the request as
// belonging to the "log in as" session key (read it with LoginAsKey), so
// Authenticate uses that session's cookie instead of the normal one. A
// malformed key answers 404.
func RouteLoginAs(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		rest, found := strings.CutPrefix(request.URL.Path, LoginAsPathPrefix)
		if !found {
			next.ServeHTTP(writer, request)
			return
		}
		key, apiPath, hasPath := strings.Cut(rest, "/")
		if !hasPath || !isLoginAsKey(key) {
			response.SendResponse(writer, request, response.NotFound, nil)
			return
		}
		routed := request.Clone(context.WithValue(request.Context(), loginAsKey{}, key))
		routed.URL.Path = "/api/" + apiPath
		routed.URL.RawPath = ""
		next.ServeHTTP(writer, routed)
	})
}

// LoginAsKey returns the "log in as" session key of the request (set by
// RouteLoginAs). ok is false for a normal /api request.
func LoginAsKey(ctx context.Context) (key string, ok bool) {
	key, ok = ctx.Value(loginAsKey{}).(string)
	return key, ok
}

// isLoginAsKey reports whether key has the form of a "log in as" key.
func isLoginAsKey(key string) bool {
	if len(key) != loginAsKeyLength {
		return false
	}
	_, err := hex.DecodeString(key)
	return err == nil
}

// SessionToken returns the session token of the request: the "log in as"
// cookie on a /api/as/<key>/ request (never the normal cookie, so an ended
// "log in as" session cannot fall back to the superadmin's), else the normal
// session cookie. Empty when there is none.
func SessionToken(request *http.Request) string {
	cookieName := SessionCookieName
	if _, ok := LoginAsKey(request.Context()); ok {
		cookieName = LoginAsCookieName
	}
	if cookie, err := request.Cookie(cookieName); err == nil {
		return cookie.Value
	}
	return ""
}

// Authenticate resolves the session cookie (see SessionToken) to the logged-in user and puts it
// in the request context (read it with CurrentUser or UserID). Requests
// without a valid session get 401 and never reach the handler.
func Authenticate(authService *service.AuthService) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			user, err := authService.Authenticate(request.Context(), SessionToken(request))
			if err != nil {
				response.SendResponse(writer, request, err, nil)
				return
			}
			next.ServeHTTP(writer, request.WithContext(WithAuthUser(request.Context(), user)))
		})
	}
}

// RequireRole lets only users with one of roles through; others get 403. It
// must run after Authenticate.
func RequireRole(roles ...string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			user, ok := CurrentUser(request.Context())
			if !ok {
				response.SendResponse(writer, request, response.LoginRequired, nil)
				return
			}
			if !slices.Contains(roles, user.Role) {
				response.SendResponse(writer, request, response.Forbidden, nil)
				return
			}
			next.ServeHTTP(writer, request)
		})
	}
}

// WithAuthUser returns ctx carrying user as the logged-in user.
func WithAuthUser(ctx context.Context, user model.AuthUser) context.Context {
	return context.WithValue(ctx, authUserKey{}, user)
}

// CurrentUser returns the logged-in user set by Authenticate. ok is false when
// the request was not authenticated.
func CurrentUser(ctx context.Context) (user model.AuthUser, ok bool) {
	user, ok = ctx.Value(authUserKey{}).(model.AuthUser)
	return user, ok && user.ID > 0
}

// UserID returns the logged-in user's id set by Authenticate. ok is false when
// the request was not authenticated. A user id is never taken from the request body.
func UserID(ctx context.Context) (userID int64, ok bool) {
	user, ok := CurrentUser(ctx)
	return user.ID, ok
}
