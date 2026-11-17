package model

import "time"

// Session statuses (auth_tokens.status, migration 034). Only an active session authenticates;
// logged out and revoked sessions are also soft-deleted (is_deleted).
const (
	SessionStatusActive    = "active"
	SessionStatusLoggedOut = "logged_out"
	SessionStatusRevoked   = "revoked"
	SessionStatusExpired   = "expired"
)

// Why a session was revoked (auth_tokens.revoke_reason).
const (
	// RevokeReasonReplacedByNewLogin: the user logged in elsewhere and chose to log this one out.
	RevokeReasonReplacedByNewLogin = "replaced_by_new_login"
	// RevokeReasonUserRevoked: the user ended it from their own session list.
	RevokeReasonUserRevoked = "user_revoked"
	// RevokeReasonAdminRevoked: the superadmin ended it.
	RevokeReasonAdminRevoked = "admin_revoked"
)

// Session is a login session: an auth_tokens row of token_type 'session'. TokenHash names it to
// the live sockets and is never sent to a client.
type Session struct {
	ID         int64     `db:"id" json:"id"`
	UserID     int64     `db:"user_id" json:"userId"`
	TokenHash  string    `db:"token" json:"-"`
	DeviceInfo string    `db:"device_info" json:"deviceInfo"`
	IPAddress  string    `db:"ip_address" json:"ipAddress"`
	CreatedAt  time.Time `db:"created_at" json:"createdAt"`
	LastUsedAt time.Time `db:"last_used_at" json:"lastUsedAt"`
	ExpiresAt  time.Time `db:"expires_at" json:"expiresAt"`
	// ImpersonatedBy is the superadmin who started it with "log in as"; nil for a password login.
	ImpersonatedBy *int64 `db:"impersonated_by" json:"impersonatedBy"`
}

// IsLoginAs reports whether the superadmin started the session with "log in as". Such sessions
// never count toward the user's session limits.
func (session Session) IsLoginAs() bool {
	return session.ImpersonatedBy != nil
}

// SessionState is what is left of a session that no longer authenticates: its status, whether its
// lifetime is over and why it was revoked.
type SessionState struct {
	Status       string `db:"status"`
	Expired      bool   `db:"expired"`
	RevokeReason string `db:"revoke_reason"`
}

// SessionUser is the user of an active session together with the session's id.
type SessionUser struct {
	AuthUser
	SessionID int64 `db:"session_id"`
}
