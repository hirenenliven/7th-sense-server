# 0002. Layered architecture: handler, DTO, service, repository

- **Date:** 2026-11-03
- **Status:** Accepted (Inferred from existing code)

## Context

Business rules (company scoping, roles) must not leak into HTTP or SQL code.

## Decision

Handlers decode and validate DTOs (`internal/dto`, `Validate()`), services hold the rules and access checks, repositories run SQL; every response goes through `response.SendResponse` (envelope `{statusCode, success, message, data}`) with predefined `response.*` errors.

## Alternatives considered

- Fat handlers
- Generic repository / ORM layer

## Consequences

- Each layer is testable; errors are consistent.
- More files per resource.
