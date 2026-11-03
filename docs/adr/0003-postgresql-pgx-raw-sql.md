# 0003. PostgreSQL with pgx and hand-written SQL

- **Date:** 2026-11-03
- **Status:** Accepted (Inferred from existing code)

## Context

Time-series machine data (mode periods, cycles, variable values) and relational masters live together.

## Decision

Use PostgreSQL through `jackc/pgx/v5` (pool) with SQL written in the repositories; schema changes are SQL files in `migrations/`.

## Alternatives considered

- GORM / ent
- sqlc

## Consequences

- Full control over queries and performance.
- No compile-time query checking; keep repository / API tests.
