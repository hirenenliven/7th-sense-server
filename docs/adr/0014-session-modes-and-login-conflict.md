# 0014. Session modes and the login conflict, on the existing auth_tokens sessions

- **Date:** 2026-11-17
- **Status:** Accepted

## Context

Logins already create server-side sessions: a row in `auth_tokens` holding the SHA-256 of a random token that travels only in the HttpOnly `ts_session` cookie (a "log in as" tab has its own `ts_login_as` cookie, ADR 0008). `middleware.Authenticate` checks the row on every request and logout soft-deletes it. Nothing limited how many sessions a user had, nobody could see them, and a client could not tell an expired session from one that was ended on purpose. Customers asked for "one device at a time" on some installations and for a way to end a lost device's session.

## Decision

- **Extend `auth_tokens`, no new table and no JWT** (migration 034): `status` (`active`, `logged_out`, `revoked`, `expired`; CHECK), `ip_address`, `revoked_at`, `revoked_by`, `revoke_reason` (`replaced_by_new_login`, `user_revoked`, `admin_revoked`), index `(user_id, status) WHERE token_type = 'session'`. Logout and revoke still set `is_deleted` / `deleted_at`, so every existing query stays right. Only active, unexpired rows authenticate.
- **Mode from the environment, read once at start** (`config.SessionConfig`): `SESSION_MODE=single|multiple` (default multiple) and `MAX_SESSIONS` (multiple only, 0 = unlimited).
- **Conflict instead of silent replacement.** A login with valid credentials that would exceed the limit answers 409 `session_conflict`, creates nothing and lists the sessions in the way. The same request with `forceLogoutOthers: true` re-checks the password, revokes them (single: all others; multiple: the least recently used ones that make room, normally one) and creates the new session in **one transaction that starts with `SELECT … FOR UPDATE` on the user's row**, so two parallel logins of one user run one after the other and cannot both pass the limit. Login throttling (429) stays in front of the check.
- **"Log in as" sessions never count** toward either limit and never conflict, and they are left out of the user's own session list: they are the superadmin's view of the account, not the customer's devices. The superadmin's own sessions are counted on the superadmin's account only.
- **Fixed message codes the client words itself.** 401 `session_revoked` (data `{reason}`) and 401 `session_expired` are told apart from the generic "Please log in to continue" (unknown token, logged out, user deactivated), by reading the row's status only when the active lookup fails. 409 `session_conflict` has data `{mode, maxSessions, sessions}`.
- **Live sockets follow.** Every revoke closes the session's sockets at once with `session:expired {reason, revokeReason}`, emitted on each socket before its disconnect packet (a room broadcast could run after the disconnect and never arrive). The `LIVE_REVALIDATE_SEC` check stays as the safety net.
- **Session management endpoints**: `/sessions` (mine; `/sessions/others` before `/sessions/{id}`), `/admin/users/{id}/sessions` and `/admin/sessions/{id}` (`superAdminOnly`). Ownership is in the SQL (`WHERE id = $1 AND user_id = $2`); another user's session answers 404, never 403. No token or hash is ever returned.
- **Superadmin only for other users' sessions.** Customers and admins do not manage their company's users' sessions for now; adding it later means the `GET /users?membersOnly=true` company scoping in `SessionService.checkSessionUser`.
- **Smaller pieces:** `last_used_at` is written at most once a minute (a data-modifying CTE inside the session lookup); the client IP is the TCP peer unless it is one of `TRUSTED_PROXIES`, then `X-Forwarded-For` read from the right; no location lookup; the device is parsed from the user agent on read (`dto.ParseUserAgent`, a short list, not a library). A goroutine started from `main.go` marks overdue sessions `expired` every `SESSION_CLEANUP_MINUTES` and deletes ended ones after `SESSION_RETENTION_DAYS` (30).

## Alternatives considered

- **JWT with a `jti` denylist.** The app already has opaque hashed tokens checked in the database on every request, which revoke instantly. A JWT would add signing keys and a denylist to get back what exists. Rejected.
- **A new `sessions` table.** It would duplicate `auth_tokens` and need a migration of every live session and of `impersonated_by`. Rejected.
- **Kill the old session silently on a new login (last login wins).** Simple, but a stolen password then logs the owner out without a word, and a user switching tabs loses work. Rejected for the explicit conflict.
- **A unique partial index for single mode.** It cannot express MAX_SESSIONS > 1, and it would make "log in as" sessions conflict. Rejected for the row lock.
- **Customers / admins manage their members' sessions now.** Wanted later perhaps, but it adds a permission surface nobody asked for yet. Deferred (see Decision).
- **GeoIP location.** A database to ship and keep current, and wrong behind VPNs. Location is shown as unavailable.

## Consequences

- A changed `SESSION_MODE` / `MAX_SESSIONS` applies after a restart; sessions already over a new limit stay until they end, the next login then hits the conflict.
- Every 401 now costs one extra indexed query on the failure path (the state lookup); successful requests write at most once a minute.
- `auth_tokens` rows are now deleted after the retention period: auditing older logins needs the logs (`login`, `session revoked by superadmin`).
- The client must keep the three message codes; changing their text breaks it (`internal/response`, `constants/session.js`).
