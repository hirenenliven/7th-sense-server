# 7th-sense-server

CRUD REST API in Go using only the standard library.
Data is stored in PostgreSQL (pgx); see `migrations/` for the schema.

## Layout

Layered structure (same idea as our Node.js projects: routes → controller → service → db):

```
cmd/server/main.go              entry point: loads config, wires layers, starts server, graceful shutdown
internal/
├── config/config.go            env config (PORT, DB_*)                   ~ dotenv
├── db/postgres.go              PostgreSQL connection pool (pgx)          ~ db connection
├── routes/                     route registration                        ~ routes.js
│   ├── routes.go               main router + global middleware + /health
│   ├── user_routes.go
│   └── item_routes.go
├── middleware/                 logger, panic recoverer, Chain helper     ~ middleware/
├── handler/                    HTTP layer: parse request, call service   ~ controllers/
│   ├── request.go              decodeAndValidate helper
│   ├── user_handler.go
│   └── item_handler.go
├── dto/                        request bodies + Validate()               ~ validationMiddleware
├── service/                    business rules (e.g. unique email)        ~ x.service.js
├── repository/                 SQL queries only, no business rules       ~ DB calls
├── model/                      stored entities                           ~ model/*.model.js
└── response/                   API responses                             ~ common/
    ├── responses.go            predefined responses (Success, Created, NotFound, ...)  ~ Responses
    └── send.go                 response.SendResponse – writes JSON responses   ~ sendResponse
```

Request flow: `routes → middleware → handler → dto.Validate → service → repository`.

Handlers write every response with `response.SendResponse(w, r, statusCode, message)`:

```go
response.SendResponse(w, r, http.StatusOK, items)        // 200 with data
response.SendResponse(w, r, http.StatusCreated, item)     // 201 with data
response.SendResponse(w, r, http.StatusOK, nil)          // 200, message only
response.SendResponse(w, r, http.StatusInternalServerError, err)  // error from a service
```

Services and DTOs return failure responses directly (e.g. `return response.EmailExists`). Any other
error becomes a 500 and is logged, so internals never reach the client.

### Naming conventions

- Packages: short, lowercase, singular (`handler`, `service`, not `handlers`, `userService`).
- Files: `snake_case`, prefixed by entity (`user_handler.go`, `user_service.go`).
- Types: `UserHandler`, `UserService`, `UserRepository`; constructors `NewXxx`.
- Errors: `ErrXxx`.

## Adding a new entity (e.g. product)

1. `internal/model/product.go` – `Product` struct
2. `internal/dto/product_dto.go` – `ProductRequest` + `Validate()`
3. `internal/repository/product_repository.go` – `ProductRepository` struct with the SQL queries
4. A migration in `migrations/` for the table
5. `internal/service/product_service.go` – `ProductService`
6. `internal/handler/product_handler.go` – `ProductHandler`
7. `internal/routes/product_routes.go` – `registerProductRoutes`, then call it in `routes.go`
8. Wire it in `cmd/server/main.go`

## Run

First time on a new machine? Follow [docs/setup.md](docs/setup.md) (Go, Air, `.env`).

Copy `.env.example` to `.env` and configure as needed:

```sh
cp .env.example .env
```

```sh
go tool air                    # hot reload (Air is a go.mod tool)
go run ./cmd/server            # port 8080
PORT=3000 go run ./cmd/server  # custom port
```

## API

Every response has the same shape (`data` is left out when there is none):

```json
{"statusCode": 200, "success": true, "message": "success", "data": {"id": "1", "name": "Hiren"}}
{"statusCode": 404, "success": false, "message": "User not found"}
```

| Method | Path                | Description |
|--------|---------------------|-------------|
| GET    | /health             | Health check |
| GET    | /api/v1/users       | List users  |
| GET    | /api/v1/users/{id}  | Get user    |
| POST   | /api/v1/users       | Create user |
| PUT    | /api/v1/users/{id}  | Update user |
| DELETE | /api/v1/users/{id}  | Delete user |

### Users

`POST /api/v1/users`

```json
{"username": "hiren_c", "email": "hiren@example.com", "password": "secret123",
 "firstName": "Hiren", "lastName": "Chavda", "mobile": "+919876543210"}
```

`PUT /api/v1/users/{id}` takes the same fields without `password`, plus an optional
`status` (`active` / `inactive`; empty keeps the current one). PUT replaces the profile:
omitted `firstName`, `lastName` and `mobile` are cleared.

- `username`: 3–50 letters, digits, `_` or `.`; unique
- `email`: valid, unique (stored lowercase)
- `password`: 8–72 characters, stored as a bcrypt hash and never returned
- `mobile`: optional, 7–15 digits with an optional leading `+`

`DELETE` is a soft delete: the row stays with `is_deleted = true` and `deleted_at` set, and is
hidden from every other endpoint. Its email and username can be used again by a new user.

The same five routes exist for `/api/v1/items`. They are registered only once a database-backed
repository is wired in `cmd/server/main.go`; until then only `/health` responds.
