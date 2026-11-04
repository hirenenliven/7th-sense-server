# 0008. "Log in as customer" uses a separate session with a path-scoped cookie

- **Date:** 2026-11-04
- **Status:** Accepted

## Context

The superadmin must be able to open a customer's account (Account Center › Customer › View) in a new browser tab, see the app as that customer, and log out there without ending the superadmin's own session. Sessions are hashed rows in `auth_tokens` carried by one HttpOnly cookie, `ts_session` (path `/api`). Every tab of a browser shares that cookie, so a second normal login in the new tab would replace the superadmin's session, and logging out there would end it.

## Decision

- `POST /api/v1/auth/login-as {customerId}` (superadmin only, called on the superadmin's session) starts an ordinary session of the customer (`auth_tokens`, one working day, `impersonated_by` = the superadmin, migration 031) and returns a random `sessionKey` (32 hex characters, not a secret).
- Its token goes into a second HttpOnly cookie, `ts_login_as`, whose **Path is `/api/as/<sessionKey>`**. The browser sends it only with that tab's requests, and it never replaces `ts_session`.
- `middleware.RouteLoginAs` serves `/api/as/<key>/...` as `/api/...` and marks the request. `middleware.SessionToken` then reads only `ts_login_as`, never falling back to `ts_session`, so an ended customer session cannot turn into the superadmin's.
- Logout under `/api/as/<key>/` ends only that session and clears only that cookie. A password login under that path is refused (403), because its cookie would replace the superadmin's.
- The client keeps `sessionKey` per tab (sessionStorage) and calls the API under `/api/as/<sessionKey>/v1` (client ADR 0017).

## Alternatives considered

- Bearer token in sessionStorage plus an `Authorization` header: simpler routing, but page scripts could read the session token. The project keeps tokens out of reach of scripts.
- Replace the superadmin's cookie and offer "switch back": only one session per browser, and logging out of the customer would log the superadmin out.
- A separate subdomain per session: needs DNS / TLS set-up for every deployment.

## Consequences

- Every route works unchanged under `/api/as/<key>/` (the rewrite happens before the router). The request log shows the original path.
- A link opened outside the tab (middle-click, a copied URL) has no key, so it uses the superadmin's session.
- `auth_tokens.impersonated_by` records who started each customer session, for auditing.
