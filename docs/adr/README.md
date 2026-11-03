# Architecture Decision Records

Server (backend) decisions, one file each (Michael Nygard template). Never edit an accepted ADR: write a new one and mark the old one "Superseded by NNNN".

| # | Title | Status | Date |
|---|---|---|---|
| 0001 | [Go standard library net/http, no web framework](0001-go-stdlib-net-http.md) | Accepted (Inferred from existing code) | 2026-11-03 |
| 0002 | [Layered architecture: handler, DTO, service, repository](0002-layered-handler-service-repository.md) | Accepted (Inferred from existing code) | 2026-11-03 |
| 0003 | [PostgreSQL with pgx and hand-written SQL](0003-postgresql-pgx-raw-sql.md) | Accepted (Inferred from existing code) | 2026-11-03 |
| 0004 | [Opaque session token in an HttpOnly cookie](0004-session-cookie-auth.md) | Accepted (Inferred from existing code) | 2026-11-03 |
| 0005 | [CNC collector runs inside the API process](0005-in-process-cnc-collector.md) | Accepted (Inferred from existing code) | 2026-11-03 |
| 0006 | [Page-wise dashboard endpoints returning only chart fields](0006-page-wise-dashboard-endpoints.md) | Accepted | 2026-11-03 |
| 0007 | [OpenAPI 3 from swag annotations, Swagger UI behind API_DOCS, generated Postman](0007-openapi-from-swag-annotations.md) | Accepted | 2026-11-03 |

New ADR: copy the newest file, take the next number, status "Proposed" until agreed, then "Accepted"; add a row here.
