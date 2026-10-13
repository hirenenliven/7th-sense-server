package middleware

import (
	"context"
	"net/http"
	"slices"

	"7th-sense-server/internal/model"
	"7th-sense-server/internal/response"
	"7th-sense-server/internal/service"
)

// SessionCookieName is the cookie that carries the session token. It is
// HttpOnly, so page scripts can never read it.
const SessionCookieName = "ts_session"

type authUserKey struct{}

// Authenticate resolves the session cookie to the logged-in user and puts it
// in the request context (read it with CurrentUser or UserID). Requests
// without a valid session get 401 and never reach the handler.
func Authenticate(authService *service.AuthService) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			token := ""
			if cookie, err := request.Cookie(SessionCookieName); err == nil {
				token = cookie.Value
			}

			user, err := authService.Authenticate(request.Context(), token)
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
