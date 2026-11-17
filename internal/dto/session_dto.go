package dto

import (
	"strings"
	"time"

	"7th-sense-server/internal/model"
)

// Names used when a user agent does not tell the browser or the operating system.
const (
	UnknownBrowser = "Unknown browser"
	UnknownOS      = "Unknown OS"
	UnknownDevice  = "Unknown device"
)

// SessionClient is what the server knows about the browser that logs in: its user agent and IP
// address (handler.clientIP, X-Forwarded-For only behind TRUSTED_PROXIES).
type SessionClient struct {
	UserAgent string
	IPAddress string
}

// SessionResponse is one login session as the session lists and the login conflict show it. It
// never carries the token or its hash. There is no location: the server does no IP lookup.
type SessionResponse struct {
	ID     int64  `json:"id"`
	UserID int64  `json:"userId"`
	Device string `json:"device"` // "Chrome on Windows"
	// Browser and OS are parsed from the user agent.
	Browser    string    `json:"browser"`
	OS         string    `json:"os"`
	IPAddress  string    `json:"ipAddress"`
	CreatedAt  time.Time `json:"createdAt"`
	LastUsedAt time.Time `json:"lastUsedAt"`
	ExpiresAt  time.Time `json:"expiresAt"`
	// IsCurrent marks the session the request came with.
	IsCurrent bool `json:"isCurrent"`
	// IsLoginAs marks a session the superadmin started with "log in as".
	IsLoginAs bool `json:"isLoginAs"`
}

// NewSessionResponse returns session as the API shows it; currentSessionID is the caller's session.
func NewSessionResponse(session model.Session, currentSessionID int64) SessionResponse {
	browser, operatingSystem := ParseUserAgent(session.DeviceInfo)
	return SessionResponse{
		ID:         session.ID,
		UserID:     session.UserID,
		Device:     describeDevice(browser, operatingSystem),
		Browser:    browser,
		OS:         operatingSystem,
		IPAddress:  session.IPAddress,
		CreatedAt:  session.CreatedAt,
		LastUsedAt: session.LastUsedAt,
		ExpiresAt:  session.ExpiresAt,
		IsCurrent:  session.ID == currentSessionID,
		IsLoginAs:  session.IsLoginAs(),
	}
}

// NewSessionResponses returns sessions as the API shows them (never nil, so JSON has []).
func NewSessionResponses(sessions []model.Session, currentSessionID int64) []SessionResponse {
	responses := make([]SessionResponse, 0, len(sessions))
	for _, session := range sessions {
		responses = append(responses, NewSessionResponse(session, currentSessionID))
	}
	return responses
}

// SessionConflictResponse is the data of a 409 session_conflict login: the session mode, the
// limit and the active sessions in the way, most recently used first. Logging in again with
// forceLogoutOthers revokes all of them (single mode) or the least recently used (multiple mode).
type SessionConflictResponse struct {
	Mode        string            `json:"mode"`
	MaxSessions int               `json:"maxSessions"`
	Sessions    []SessionResponse `json:"sessions"`
}

// RevokeSessionsResponse answers a request that revoked several sessions.
type RevokeSessionsResponse struct {
	Revoked int `json:"revoked"`
}

// userAgentBrowsers maps user agent markers to browser names, checked in order: Edge and Opera
// also say Chrome, and Chrome also says Safari.
var userAgentBrowsers = []struct{ marker, name string }{
	{"Edg/", "Edge"}, {"EdgA/", "Edge"}, {"EdgiOS/", "Edge"},
	{"OPR/", "Opera"}, {"SamsungBrowser/", "Samsung Internet"},
	{"Firefox/", "Firefox"}, {"FxiOS/", "Firefox"},
	{"Chrome/", "Chrome"}, {"CriOS/", "Chrome"},
	{"Safari/", "Safari"},
}

// userAgentSystems maps user agent markers to operating system names, checked in order: Android
// also says Linux, iOS also says "like Mac OS X".
var userAgentSystems = []struct{ marker, name string }{
	{"Windows", "Windows"}, {"Android", "Android"},
	{"iPhone", "iOS"}, {"iPad", "iPadOS"},
	{"CrOS", "ChromeOS"}, {"Mac OS X", "macOS"}, {"Linux", "Linux"},
}

// ParseUserAgent returns the browser and operating system a user agent names, or UnknownBrowser /
// UnknownOS. It is a short list of the browsers in use, not a full parser.
func ParseUserAgent(userAgent string) (browser, operatingSystem string) {
	browser, operatingSystem = UnknownBrowser, UnknownOS
	for _, candidate := range userAgentBrowsers {
		if strings.Contains(userAgent, candidate.marker) {
			browser = candidate.name
			break
		}
	}
	for _, candidate := range userAgentSystems {
		if strings.Contains(userAgent, candidate.marker) {
			operatingSystem = candidate.name
			break
		}
	}
	return browser, operatingSystem
}

// describeDevice names a device for people: "Chrome on Windows", or UnknownDevice.
func describeDevice(browser, operatingSystem string) string {
	switch {
	case browser == UnknownBrowser && operatingSystem == UnknownOS:
		return UnknownDevice
	case operatingSystem == UnknownOS:
		return browser
	case browser == UnknownBrowser:
		return operatingSystem
	}
	return browser + " on " + operatingSystem
}
