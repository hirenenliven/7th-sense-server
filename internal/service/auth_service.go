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
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

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

// dummyPasswordHash is compared when the username does not exist, so a wrong
// username takes as long as a wrong password and usernames cannot be probed.
var dummyPasswordHash, _ = bcrypt.GenerateFromPassword([]byte("7th-sense-dummy-password"), bcrypt.DefaultCost)

// Session is a logged-in session: the token for the client's cookie and the user.
// Key is set only for a "log in as" session: the path segment its cookie is
// scoped to (/api/as/<Key>).
type Session struct {
	Token    string
	Key      string
	Lifetime time.Duration
	Remember bool
	User     model.AuthUser
}

// AuthService holds the login rules: checking credentials, creating and
// resolving sessions, and logging out. Only a hash of each session token is
// stored, so the database never holds a usable token.
type AuthService struct {
	authRepository *repository.AuthRepository
	throttle       *loginThrottle
}

// NewAuthService creates an AuthService that stores sessions with authRepository.
func NewAuthService(authRepository *repository.AuthRepository) *AuthService {
	return &AuthService{
		authRepository: authRepository,
		throttle:       newLoginThrottle(),
	}
}

// Login checks the credentials and starts a session. deviceInfo (the browser's
// user agent) is stored with the session for auditing.
func (authService *AuthService) Login(ctx context.Context, loginRequest dto.LoginRequest, deviceInfo string) (Session, error) {
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
	if err := authService.authRepository.CreateSession(ctx, user.ID, hashToken(token), truncate(deviceInfo, 255), lifetime, nil); err != nil {
		return Session{}, err
	}

	slog.Info("login", "userId", user.ID, "role", user.Role, "remember", loginRequest.Remember)
	user.Password = ""
	return Session{Token: token, Lifetime: lifetime, Remember: loginRequest.Remember, User: user}, nil
}

// LoginAs starts a session of customer customerID for actor (the superadmin),
// who opens that account in a new browser tab. It is an ordinary customer
// session that records actor as impersonated_by and lasts one working day. It
// gets its own key, so its cookie never replaces actor's own session.
func (authService *AuthService) LoginAs(ctx context.Context, actor model.AuthUser, customerID int64, deviceInfo string) (Session, error) {
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
	if err := authService.authRepository.CreateSession(ctx, customer.ID, hashToken(token), truncate(deviceInfo, 255), DefaultSessionLifetime, &actor.ID); err != nil {
		return Session{}, err
	}
	slog.Info("login as", "userId", customer.ID, "impersonatedBy", actor.ID)
	return Session{Token: token, Key: key, Lifetime: DefaultSessionLifetime, User: customer}, nil
}

// Authenticate returns the user of the session token, or LoginRequired when
// the token is unknown, expired or logged out, or its user can no longer log in.
func (authService *AuthService) Authenticate(ctx context.Context, token string) (model.AuthUser, error) {
	if token == "" {
		return model.AuthUser{}, response.LoginRequired
	}
	user, err := authService.authRepository.FindUserBySession(ctx, hashToken(token))
	if errors.Is(err, repository.ErrNotFound) {
		return model.AuthUser{}, response.LoginRequired
	}
	return user, err
}

// Logout ends the session of token. An unknown or empty token is not an error.
func (authService *AuthService) Logout(ctx context.Context, token string) error {
	if token == "" {
		return nil
	}
	return authService.authRepository.DeleteSession(ctx, hashToken(token))
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
