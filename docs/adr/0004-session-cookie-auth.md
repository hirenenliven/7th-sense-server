# 0004. Opaque session token in an HttpOnly cookie

- **Date:** 2026-11-03
- **Status:** Accepted (Inferred from existing code)

## Context

Browser clients only; logout and deactivation must take effect at once.

## Decision

Login checks the bcrypt password hash and stores a random session token in the database; the token travels in the HttpOnly `ts_session` cookie (path `/api`). `middleware.Authenticate` resolves it on every request; `RequireRole` (`managersOnly`, `superAdminOnly`) checks roles and services scope data to the user's company.

## Alternatives considered

- JWT access + refresh tokens
- Third-party identity provider

## Consequences

- Sessions can be revoked server side; no token in JavaScript.
- One database lookup per request; API tools must keep the cookie (Postman cookie jar).
