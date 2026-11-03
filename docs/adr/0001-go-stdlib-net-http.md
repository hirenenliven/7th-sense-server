# 0001. Go standard library net/http, no web framework

- **Date:** 2026-11-03
- **Status:** Accepted (Inferred from existing code)

## Context

The API is a REST service plus an in-process CNC collector; Go 1.22+ ServeMux supports method and path patterns.

## Decision

Use `net/http` with `http.ServeMux` patterns (`GET /api/v1/shifts/{id}`) and small middleware (`middleware.Chain`: logger, recoverer, CORS, authenticate, RequireRole).

## Alternatives considered

- Gin / Echo / Fiber
- chi router

## Consequences

- No framework lock-in; standard tooling.
- What a framework would give (binding, docs) is done explicitly: DTO decode / Validate, swag annotations.
