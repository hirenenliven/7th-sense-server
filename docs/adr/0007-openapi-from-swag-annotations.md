# 0007. OpenAPI 3 from swag annotations, Swagger UI behind API_DOCS, generated Postman

- **Date:** 2026-11-03
- **Status:** Accepted

## Context

API docs must not drift from the code and must work in Swagger UI and Postman. The server is stdlib Go, so framework-native generators (NestJS, FastAPI) do not apply.

## Decision

Annotate handlers with swaggo/swag comments (summary, description ending "Used by pages: ...", tags, operationId, params, request / response schemas from the Go DTOs). `go generate ./docs/api` runs `go tool swag`, then `tools/apidocs` converts Swagger 2.0 to OpenAPI 3.0.3 (cookie auth, servers, examples), writes `docs/api/openapi.yaml` and builds the Postman collection and environment in `docs/postman/`. Swagger UI (`swaggo/http-swagger`) is served at `/api-docs/` only when `API_DOCS=true`.

## Alternatives considered

- Hand-written openapi.yaml (drifts)
- swag v2 with native OpenAPI 3.1 (still a release candidate)
- swagger2openapi via npx (adds Node to the Go toolchain)

## Consequences

- Schemas come from the real DTO structs; one command regenerates spec and Postman.
- Annotations must change with the handler; so far only endpoints the client uses are documented.
