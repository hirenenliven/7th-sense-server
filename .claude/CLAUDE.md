# 7th Sense Server - Claude Context

## Project Overview
Go-based REST API backend for the 7th Sense system. Uses PostgreSQL with pgx driver and TimescaleDB for time-series data, following a strict layered architecture pattern designed for scalability and maintainability.

## Tech Stack
- **Language**: Go 1.27.1
- **Database**: PostgreSQL + TimescaleDB
- **Driver**: pgx/v5
- **Architecture**: Layered (Handler → Service → Repository → Model)
- **Hot Reload**: Air (go.mod tool)
- **Key Dependencies**: jackc/pgx/v5, joho/godotenv, golang.org/x/crypto (bcrypt), lumberjack (logging)

## Directory Structure
```
7th-sense-server/
├── cmd/server/              # Application entry point (main.go)
├── internal/
│   ├── config/              # Environment configuration
│   ├── db/                  # Database connection utilities (pgx pool)
│   ├── dto/                 # Data Transfer Objects (request validation)
│   ├── handler/             # HTTP request handlers (controllers)
│   ├── logger/              # Logging utilities
│   ├── middleware/          # HTTP middleware (logging, recovery)
│   ├── model/               # Database models (structs)
│   ├── repository/          # SQL queries and data access
│   ├── response/            # API response definitions
│   ├── routes/              # Route registration
│   └── service/             # Business logic layer
├── migrations/              # Database migration files
├── tests/                   # Integration tests
├── logs/                    # Application logs
└── .air.toml               # Air hot reload config
```

## Architecture Pattern (STRICT)

### Layered Architecture
Request flow: `routes → middleware → handler → dto.Validate → service → repository → db`

**Layers (no skipping allowed):**
1. **Handler** (`internal/handler/`) - HTTP concerns, request parsing, response sending
2. **Service** (`internal/service/`) - Business logic, error mapping, orchestration
3. **Repository** (`internal/repository/`) - SQL queries only, no business logic
4. **Model** (`internal/model/`) - Database structs with db/json tags
5. **DTO** (`internal/dto/`) - Request validation via `Validate()` method

### Critical Rules
- **Never skip layers** - Handlers cannot call repositories directly
- **Concrete types over interfaces** - Use structs, not interfaces (YAGNI)
- **Error mapping** - Repository errors → Response errors in service layer
- **COALESCE for nullables** - Always use `COALESCE(column, '')` in SELECT
- **Soft delete** - Every table has `is_deleted` and `deleted_at`
- **Validation in DTOs only** - Never validate in handlers/services

## CRUD Implementation Pattern

**Golden Reference**: User CRUD (`user_handler.go`, `user_service.go`, `user_repository.go`) is the canonical implementation. Copy this pattern exactly for new entities.

### Steps for New Entity:
1. **Model** (`internal/model/entity.go`) - Struct with db/json tags
2. **DTO** (`internal/dto/entity_dto.go`) - Request structs with `Validate()` method
3. **Repository** (`internal/repository/entity_repository.go`) - SQL queries via `db.Execute()`
4. **Service** (`internal/service/entity_service.go`) - Business logic + error mapping
5. **Response Errors** (`internal/response/responses.go`) - Add entity-specific responses
6. **Handler** (`internal/handler/entity_handler.go`) - HTTP handlers using `response.SendResponse()`
7. **Routes** (`internal/routes/entity_routes.go`) - Register routes in ServeMux
8. **Wiring** (`cmd/server/main.go`) - Wire dependencies in main

### Repository Pattern
```go
func (r *EntityRepository) FindByID(ctx context.Context, id int64) (*model.Entity, error) {
    query := `SELECT id, field1, COALESCE(field2, '') AS field2, ... 
              FROM entities WHERE id = $1 AND is_deleted = false`
    return db.Execute(ctx, r.dbPool, query, id)
}
```

### Service Error Mapping
```go
func toEntityResponseError(err error) error {
    switch {
    case errors.Is(err, repository.ErrNotFound):
        return response.EntityNotFound
    case errors.Is(err, repository.ErrDuplicateField):
        return response.FieldExists
    default:
        return err  // Becomes 500
    }
}
```

### Handler Pattern
```go
func (h *EntityHandler) GetByID(w http.ResponseWriter, r *http.Request) {
    id, err := pathID(r)  // Extract ID from path
    if err != nil {
        response.SendResponse(w, r, http.StatusBadRequest, err)
        return
    }
    
    entity, err := h.service.GetByID(r.Context(), id)
    response.SendResponse(w, r, err, entity)
}
```

## Naming Conventions
- **Packages**: lowercase, singular (`handler`, `service`, not `handlers`)
- **Files**: snake_case, entity-prefixed (`user_handler.go`, `user_service.go`)
- **Types**: PascalCase (`UserHandler`, `UserService`, `UserRepository`)
- **Constructors**: `NewXxx` (`NewUserHandler`)
- **Errors**: `ErrXxx` (`ErrNotFound`, `ErrDuplicateField`)

## Response Format
All responses use `response.SendResponse(w, r, statusCode, data)`:
```json
{"statusCode": 200, "success": true, "message": "success", "data": {...}}
{"statusCode": 404, "success": false, "message": "Entity not found"}
```

## Common Commands

### Development
```bash
go tool air                    # Hot reload (air config in .air.toml)
go run ./cmd/server            # Run without hot reload
PORT=3000 go run ./cmd/server  # Custom port
```

### Testing
```bash
go test ./tests/...            # Run all tests
go test -v ./tests/...         # Verbose output
```

### Database
```bash
psql -U postgres -d 7th_sense -f migrations/schema.sql
```

### Dependencies
```bash
go mod download                # Download dependencies
go mod tidy                    # Clean up go.mod
```

## API Endpoints

Base path: `/api/v1`

**Users**: GET/POST/PUT/DELETE `/users` — company-scoped (`service.UserService`). A customer is a company; its members (admin, supervisor, operator, qc, shop_head) have `users.company_id` = the customer. The superadmin manages every user (`?filterRole=customer`, `?companyId=` filter the list; `filterRole` is the listed users' role, never the caller's) and names a member's `companyId`; a customer or its admins see and change only their company's members (`userManagersOnly` route guard, others get 403; other companies' users answer 404). Roles a user may give: `model.AssignableRoles` — superadmin: customer + member roles; customer: member roles; admin: member roles but admin; nobody gives superadmin. Body `machineIds`: machines (of the company) a supervisor/operator/qc/shop_head may access, stored in `user_machine_access` (migration 015; no soft delete — removing an assignment deletes the row); admins and customers see every company machine and store none. Responses include `role` (name) and `machineIds`.
**Roles**: GET `/roles` — the roles the current user may give (role dropdown of the user form); empty for users who manage no users.
**Items**: GET/POST/PUT/DELETE `/items`
**Branches**: GET/POST/PUT/DELETE `/branches`
**Brands**: GET/POST/PUT/DELETE `/brands` — every logged-in user reads; only the superadmin writes (`superAdminOnly`).
**Categories**: GET/POST/PUT/DELETE `/categories`
**Products**: GET/POST/PUT/DELETE `/products`
**Controllers**: GET/POST/PUT/DELETE `/controllers` — machine catalog Brand → Controller → Model. Body `{brandId, name, description}`: no code or variant, the variant is part of the name (e.g. "Sinumerik 828D", migration 017). Brand required, name unique per brand (case-insensitive), `?brandId=` on list. Cannot be deleted while it has models or machines use it; its brand cannot change while machines use it.
**Models**: GET/POST/PUT/DELETE `/models` (`models` table, `MachineModel` in Go — `model` is the package name) — `{controllerId, name}`, name unique per controller (case-insensitive), `?controllerId=` on list. Cannot be deleted, or moved to another controller, while machines use it. Brands, controllers and models: every logged-in user reads (machine form dropdowns); only the superadmin writes (`superAdminOnly` in `internal/routes`).
**Addresses**: GET/POST/PUT/DELETE `/addresses`
**Catalog**: GET/POST `/catalog`
**Auth**: POST `/auth/login`, POST `/auth/logout` (public), GET `/auth/me` — HttpOnly `ts_session` cookie, hashed sessions in `auth_tokens` (see README › Authentication)
**Machines**: GET/POST/PUT/DELETE `/machines` — each machine belongs to a customer (`user_id`) and has a brand, a controller of that brand and a model (`model_id`) of that controller; `product_id` is no longer used. The superadmin reads all (`?userId=` narrows to one customer) and is the only one who writes (`superAdminOnly`); a customer and its admins read every machine of the company, other members only the machines assigned to them (`user_machine_access`, see `MachineService.machineScope`); anything else answers 404. The password is write-only (never returned; blank on PUT keeps it). Form-only writes: runtime columns (online status, mode, …) are not changed through the API. Accuracy / health options are fixed lists in the client (`constants/machines.js`).
**Current user**: GET `/me`
**Organization** (per user, `?user<Parent>Id=` on list): GET/POST/PUT/DELETE `/user-city|branch|location|shop|line|cell|supervisor`, GET/POST/DELETE `/user-country|state` — one set of files per table (`user_<level>_*.go`). The current user comes from the session (`middleware.Authenticate`); the superadmin may act for another user with `?userId=<id>` on any of these routes (`targetUserID` in `handler/request.go`), everyone else gets 403 for any id but their own. Every `/api` route except login/logout needs a session; role-gated routes are wrapped with `superAdminOnly` / `userManagersOnly` in `internal/routes` (categories, products, machines, brands, controllers and models are written by the superadmin only).

## Configuration

**Environment Variables** (`.env` file):
```
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_DATABASE=7th_sense
DB_SSLMODE=disable
COOKIE_SECURE=false   # true when served over HTTPS
```

Default values in `internal/config/config.go` if missing.

## Key Files Reference

- **Entry Point**: `cmd/server/main.go` - Wiring, server startup, graceful shutdown
- **Routes**: `internal/routes/routes.go` - Main router + middleware
- **Response Helpers**: `internal/response/responses.go` - Predefined responses
- **DB Connection**: `internal/db/postgres.go` - pgx pool setup
- **Migration Schema**: `migrations/schema.sql` - Complete database schema
- **Skills Document**: `7th_server_skills.md.md` - Detailed engineering principles

## Current Entities (Models)

- User, UserProfile, UserMachineAccess (`user_machine_access`; the older `user_machines` table is unused), Role
- Item
- Branch
- Brand, Category, Product
- Controller, MachineModel (models table)
- Address
- CompanyProfile
- Machine

## Important Notes

1. **Soft Delete**: All queries must include `WHERE is_deleted = false`
2. **Unique Indexes**: Must include `WHERE is_deleted = false` for soft-delete compatibility
3. **Context**: Always pass `ctx` through all layers for cancellation support
4. **Logging**: Use structured logging in service layer, never log sensitive data
5. **Transactions**: Service layer only, keep short, rollback on error
6. **Idempotency**: PUT and DELETE must be idempotent

## Testing

Integration tests in `tests/` directory using setup helpers:
- `setup_test.go` - Test database setup
- `api_helpers_test.go` - HTTP test helpers
- Entity-specific test files (e.g., `user_profile_test.go`)

## Debugging

VS Code configurations in `.vscode/launch.json`:
- **Debug Server** - Start under debugger (F5)
- **Debug Current Test** - Debug tests in open file
- **Attach to Running Server** - Attach to running process

Requires Delve debugger: `go install github.com/go-delve/delve/cmd/dlv@latest`

## Engineering Principles Summary

From `7th_server_skills.md.md`:
- Predictability over cleverness
- Consistency over individual preferences
- Simplicity over complexity
- Layered architecture is sacred
- Never trust input (defense in depth)
- Test behavior, not implementation
- Profile before optimizing
- Environment parity (dev/staging/prod)
- Documentation at implementation time

**Discipline in the small things creates freedom in the big things.**
