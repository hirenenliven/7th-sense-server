package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"7th-sense-server/internal/config"
	"7th-sense-server/internal/db"
	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/response"
)

// Session lifetimes. A remembered session survives closing the browser; any
// other session lasts one working day.
const (
	RememberSessionLifetime = 30 * 24 * time.Hour
	DefaultSessionLifetime  = 12 * time.Hour
)

// sessionTokenBytes is the size of a session token before encoding (256 bits).
const sessionTokenBytes = 32

// loginAsKeyBytes is the size of a "log in as" session key before hex encoding.
const loginAsKeyBytes = 16

// sessionTouchInterval is how often a session's last_used_at is written at most, so requests of a
// busy screen do not each write to auth_tokens.
const sessionTouchInterval = time.Minute

// deviceInfoMaxLength is the size of auth_tokens.device_info.
const deviceInfoMaxLength = 255

// ipAddressMaxLength is the size of auth_tokens.ip_address (an IPv6 address in text).
const ipAddressMaxLength = 45

// dummyPasswordHash is compared when the username does not exist, so a wrong
// username takes as long as a wrong password and usernames cannot be probed.
var dummyPasswordHash, _ = bcrypt.GenerateFromPassword([]byte("7th-sense-dummy-password"), bcrypt.DefaultCost)

// Session is a logged-in session: the token for the client's cookie and the user.
// Key is set only for a "log in as" session: the path segment its cookie is
// scoped to (/api/as/<Key>). RevokedSessions are the token hashes of the sessions
// the login revoked (forceLogoutOthers), whose live sockets must be closed.
type Session struct {
	Token           string
	Key             string
	Lifetime        time.Duration
	Remember        bool
	User            model.AuthUser
	RevokedSessions []string
}

// SessionConflictError is the 409 session_conflict of a login that would exceed the user's
// sessions; Conflict is its data. It unwraps to response.SessionConflict.
type SessionConflictError struct {
	Conflict dto.SessionConflictResponse
}

func (conflictError *SessionConflictError) Error() string { return response.SessionConflict.Error() }

// Unwrap returns response.SessionConflict, so SendResponse sends the 409.
func (conflictError *SessionConflictError) Unwrap() error { return response.SessionConflict }

// SessionRevokedError is the 401 session_revoked of a request whose session was revoked;
// RevokeReason (model.RevokeReason*) tells the client why. It unwraps to response.SessionRevoked.
type SessionRevokedError struct {
	RevokeReason string
}

func (revokedError *SessionRevokedError) Error() string { return response.SessionRevoked.Error() }

// Unwrap returns response.SessionRevoked, so SendResponse sends the 401.
func (revokedError *SessionRevokedError) Unwrap() error { return response.SessionRevoked }

// AuthService holds the login rules: checking credentials, the session limits
// (SESSION_MODE / MAX_SESSIONS), creating and resolving sessions, and logging
// out. Only a hash of each session token is stored, so the database never holds
// a usable token.
type AuthService struct {
	dbPool         *pgxpool.Pool
	authRepository *repository.AuthRepository
	sessionConfig  config.SessionConfig
	throttle       *loginThrottle
}

// NewAuthService creates an AuthService that stores sessions in dbPool with the limits of sessionConfig.
func NewAuthService(dbPool *pgxpool.Pool, sessionConfig config.SessionConfig) *AuthService {
	return &AuthService{
		dbPool:         dbPool,
		authRepository: repository.NewAuthRepository(dbPool),
		sessionConfig:  sessionConfig,
		throttle:       newLoginThrottle(),
	}
}

// Login checks the credentials and starts a session for client (its user agent and IP are stored
// with the session). When the user already has as many sessions as allowed it answers a
// SessionConflictError, unless loginRequest.ForceLogoutOthers asks to revoke the sessions in the way.
func (authService *AuthService) Login(ctx context.Context, loginRequest dto.LoginRequest, client dto.SessionClient) (Session, error) {
	throttleKey := strings.ToLower(loginRequest.Username)
	if authService.throttle.isLocked(throttleKey) {
		return Session{}, response.TooManyLoginAttempts
	}

	user, err := authService.authRepository.FindUserByUsername(ctx, loginRequest.Username)
	switch {
	case errors.Is(err, repository.ErrNotFound):
		_ = bcrypt.CompareHashAndPassword(dummyPasswordHash, []byte(loginRequest.Password))
		authService.throttle.recordFailure(throttleKey)
		return Session{}, response.InvalidCredentials
	case err != nil:
		return Session{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		authService.throttle.recordFailure(throttleKey)
		slog.Warn("login failed", "userId", user.ID, "reason", "wrong password")
		return Session{}, response.InvalidCredentials
	}
	// Checked after the password, so the status of an account is only told to its owner.
	if user.Status != "" && user.Status != dto.UserStatusActive {
		slog.Warn("login failed", "userId", user.ID, "reason", "inactive")
		return Session{}, response.AccountInactive
	}
	authService.throttle.reset(throttleKey)

	token, err := newSessionToken()
	if err != nil {
		return Session{}, err
	}
	lifetime := DefaultSessionLifetime
	if loginRequest.Remember {
		lifetime = RememberSessionLifetime
	}
	revoked, err := authService.startSession(ctx, user.ID, hashToken(token), client, lifetime, loginRequest.ForceLogoutOthers)
	if err != nil {
		return Session{}, err
	}

	slog.Info("login", "userId", user.ID, "role", user.Role, "remember", loginRequest.Remember, "revokedSessions", len(revoked))
	user.Password = ""
	return Session{Token: token, Lifetime: lifetime, Remember: loginRequest.Remember, User: user, RevokedSessions: revoked}, nil
}

// startSession creates the session of userID within the session limits, in one transaction that
// first locks the user's row: two logins of the same user then run one after the other, so both
// cannot pass the limit. Sessions in the way answer a SessionConflictError, or are revoked when
// forceLogoutOthers is set; their token hashes are returned.
func (authService *AuthService) startSession(ctx context.Context, userID int64, tokenHash string, client dto.SessionClient, lifetime time.Duration, forceLogoutOthers bool) ([]string, error) {
	var revoked []string
	err := db.WithTransaction(ctx, authService.dbPool, func(tx pgx.Tx) error {
		sessionRepository := repository.NewSessionRepository(tx)
		if err := sessionRepository.LockUser(ctx, userID); err != nil {
			return err
		}
		active, err := sessionRepository.FindAllActive(ctx, userID)
		if err != nil {
			return err
		}
		// "Log in as" sessions are the superadmin's view of the account, not the user's devices.
		own := slices.DeleteFunc(active, model.Session.IsLoginAs)
		inTheWay := authService.sessionsInTheWay(own)
		if len(inTheWay) > 0 {
			if !forceLogoutOthers {
				return &SessionConflictError{Conflict: dto.SessionConflictResponse{
					Mode: authService.sessionConfig.Mode, MaxSessions: authService.sessionConfig.MaxSessions,
					Sessions: dto.NewSessionResponses(own, 0),
				}}
			}
			revoked, err = sessionRepository.RevokeByIDs(ctx, sessionIDs(inTheWay), userID, model.RevokeReasonReplacedByNewLogin)
			if err != nil {
				return err
			}
		}
		_, err = repository.NewAuthRepository(tx).CreateSession(ctx, userID, tokenHash,
			truncate(client.UserAgent, deviceInfoMaxLength), truncate(client.IPAddress, ipAddressMaxLength), lifetime, nil)
		return err
	})
	return revoked, err
}

// sessionsInTheWay returns the sessions of own (most recently used first) a new login must end:
// single mode every one; multiple mode with MAX_SESSIONS the least recently used ones that leave
// room for the new session (normally just one); none otherwise.
func (authService *AuthService) sessionsInTheWay(own []model.Session) []model.Session {
	if authService.sessionConfig.Mode == config.SessionModeSingle {
		return own
	}
	maxSessions := authService.sessionConfig.MaxSessions
	if maxSessions <= 0 || len(own) < maxSessions {
		return nil
	}
	return own[maxSessions-1:]
}

// sessionIDs returns the ids of sessions.
func sessionIDs(sessions []model.Session) []int64 {
	ids := make([]int64, 0, len(sessions))
	for _, session := range sessions {
		ids = append(ids, session.ID)
	}
	return ids
}

// LoginAs starts a session of customer customerID for actor (the superadmin),
// who opens that account in a new browser tab from client. It is an ordinary customer
// session that records actor as impersonated_by and lasts one working day. It
// gets its own key, so its cookie never replaces actor's own session.
func (authService *AuthService) LoginAs(ctx context.Context, actor model.AuthUser, customerID int64, client dto.SessionClient) (Session, error) {
	if actor.Role != model.RoleSuperAdmin {
		return Session{}, response.Forbidden
	}
	customer, err := authService.authRepository.FindUserByID(ctx, customerID)
	if errors.Is(err, repository.ErrNotFound) || (err == nil && customer.Role != model.RoleCustomer) {
		return Session{}, response.CustomerNotFound
	}
	if err != nil {
		return Session{}, err
	}
	if customer.Status != "" && customer.Status != dto.UserStatusActive {
		return Session{}, response.AccountInactive
	}

	token, err := newSessionToken()
	if err != nil {
		return Session{}, err
	}
	key, err := newLoginAsKey()
	if err != nil {
		return Session{}, err
	}
	// No session limit: a "log in as" session never conflicts with the customer's own.
	if _, err := authService.authRepository.CreateSession(ctx, customer.ID, hashToken(token),
		truncate(client.UserAgent, deviceInfoMaxLength), truncate(client.IPAddress, ipAddressMaxLength), DefaultSessionLifetime, &actor.ID); err != nil {
		return Session{}, err
	}
	slog.Info("login as", "userId", customer.ID, "impersonatedBy", actor.ID)
	return Session{Token: token, Key: key, Lifetime: DefaultSessionLifetime, User: customer}, nil
}

// Authenticate returns the user of the session token. See AuthenticateSession for the errors.
func (authService *AuthService) Authenticate(ctx context.Context, token string) (model.AuthUser, error) {
	sessionUser, err := authService.AuthenticateSession(ctx, token)
	return sessionUser.AuthUser, err
}

// AuthenticateSession returns the user of the session token with the session's id. A session
// that no longer authenticates answers SessionRevokedError (revoked), SessionExpired (its
// lifetime is over) or LoginRequired (unknown token, logged out, user inactive or deleted).
func (authService *AuthService) AuthenticateSession(ctx context.Context, token string) (model.SessionUser, error) {
	if token == "" {
		return model.SessionUser{}, response.LoginRequired
	}
	sessionUser, err := authService.authRepository.FindUserBySession(ctx, hashToken(token), sessionTouchInterval)
	if errors.Is(err, repository.ErrNotFound) {
		return model.SessionUser{}, authService.sessionEndedError(ctx, hashToken(token))
	}
	return sessionUser, err
}

// sessionEndedError tells why the session of tokenHash no longer authenticates, so the client
// can say so: revoked, expired, or anything else (LoginRequired).
func (authService *AuthService) sessionEndedError(ctx context.Context, tokenHash string) error {
	state, err := authService.authRepository.FindSessionState(ctx, tokenHash)
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return response.LoginRequired
	case err != nil:
		return err
	case state.Status == model.SessionStatusRevoked:
		return &SessionRevokedError{RevokeReason: state.RevokeReason}
	case state.Status == model.SessionStatusExpired, state.Status == model.SessionStatusActive && state.Expired:
		return response.SessionExpired
	}
	// Logged out, or active but its user was deactivated or deleted.
	return response.LoginRequired
}

// Logout ends the session of token. An unknown or empty token is not an error.
func (authService *AuthService) Logout(ctx context.Context, token string) error {
	if token == "" {
		return nil
	}
	return authService.authRepository.DeleteSession(ctx, hashToken(token))
}

// SessionHash returns the id under which token's session is stored (auth_tokens.token), so
// other parts of the server can name a session without keeping its token.
func (authService *AuthService) SessionHash(token string) string {
	return hashToken(token)
}

// newSessionToken returns a random URL-safe session token.
func newSessionToken() (string, error) {
	tokenBytes := make([]byte, sessionTokenBytes)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", fmt.Errorf("generate session token: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(tokenBytes), nil
}

// newLoginAsKey returns a random hex key for the API path of a "log in as" session.
func newLoginAsKey() (string, error) {
	keyBytes := make([]byte, loginAsKeyBytes)
	if _, err := rand.Read(keyBytes); err != nil {
		return "", fmt.Errorf("generate login as key: %w", err)
	}
	return hex.EncodeToString(keyBytes), nil
}

// hashToken returns the hex SHA-256 of token, the value stored in auth_tokens.token.
func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// truncate shortens value to at most maxLength bytes.
func truncate(value string, maxLength int) string {
	if len(value) > maxLength {
		return value[:maxLength]
	}
	return value
}
