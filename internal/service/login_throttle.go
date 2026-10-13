package service

import (
	"sync"
	"time"
)

// Login throttling: after maxLoginFailures wrong attempts for one username
// within loginFailureWindow, further attempts are refused until the window
// ends. This slows down password guessing.
const (
	maxLoginFailures   = 5
	loginFailureWindow = 15 * time.Minute
)

// loginThrottle counts failed logins per username in memory. Counts reset when
// the server restarts, which is acceptable for slowing down guessing.
type loginThrottle struct {
	mutex    sync.Mutex
	failures map[string]loginFailures
}

// loginFailures is the failure count of one username since firstFailure.
type loginFailures struct {
	count        int
	firstFailure time.Time
}

func newLoginThrottle() *loginThrottle {
	return &loginThrottle{failures: make(map[string]loginFailures)}
}

// isLocked reports whether username has too many recent failures.
func (throttle *loginThrottle) isLocked(username string) bool {
	throttle.mutex.Lock()
	defer throttle.mutex.Unlock()

	entry, ok := throttle.failures[username]
	if !ok {
		return false
	}
	if time.Since(entry.firstFailure) > loginFailureWindow {
		delete(throttle.failures, username)
		return false
	}
	return entry.count >= maxLoginFailures
}

// recordFailure counts a failed login for username.
func (throttle *loginThrottle) recordFailure(username string) {
	throttle.mutex.Lock()
	defer throttle.mutex.Unlock()

	now := time.Now()
	throttle.removeExpired(now)
	entry, ok := throttle.failures[username]
	if !ok {
		entry = loginFailures{firstFailure: now}
	}
	entry.count++
	throttle.failures[username] = entry
}

// reset clears the failures of username after a successful login.
func (throttle *loginThrottle) reset(username string) {
	throttle.mutex.Lock()
	defer throttle.mutex.Unlock()
	delete(throttle.failures, username)
}

// removeExpired drops entries older than the window so the map cannot grow
// without bound. The caller holds the mutex.
func (throttle *loginThrottle) removeExpired(now time.Time) {
	for username, entry := range throttle.failures {
		if now.Sub(entry.firstFailure) > loginFailureWindow {
			delete(throttle.failures, username)
		}
	}
}
