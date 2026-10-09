# 7th-sense-server

CRUD REST API in Go using only the standard library.
Uses a repository interface for data access, making it easy to plug in a database implementation later.

## Layout

Layered structure (same idea as our Node.js projects: routes → controller → service → db):

```
cmd/server/main.go              entry point: loads config, wires layers, starts server, graceful shutdown
internal/
├── config/config.go            env config (PORT)                         ~ dotenv
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
├── repository/                 data access (interface only)              ~ DB calls
├── model/                      stored entities                           ~ model/*.model.js
└── response/                   response helpers                          ~ common/
    ├── api_error.go            predefined errors (ErrUserNotFound, ...)  ~ Responses
    └── send.go                 Send / SendError                          ~ sendResponse
```

Request flow: `routes → middleware → handler → dto.Validate → service → repository`.

Services return `*response.APIError` values (e.g. `response.ErrEmailExists`); `response.SendError`
turns them into the right status code. Any other error becomes a 500 and is logged.

### Naming conventions

- Packages: short, lowercase, singular (`handler`, `service`, not `handlers`, `userService`).
- Files: `snake_case`, prefixed by entity (`user_handler.go`, `user_service.go`).
- Types: `UserHandler`, `UserService`, `UserRepository`; constructors `NewXxx`.
- Errors: `ErrXxx`.

## Adding a new entity (e.g. product)

1. `internal/model/product.go` – `Product` struct
2. `internal/dto/product_dto.go` – `ProductRequest` + `Validate()`
3. `internal/repository/product_repository.go` – `ProductRepository` interface
4. `internal/repository/postgres_product_repository.go` – Database implementation (when DB is added)
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
go run ./cmd/server            # port 8080
PORT=3000 go run ./cmd/server  # custom port
```

## API

All responses use the envelope `{"statusCode": <code>, "data": <payload>}`.
Errors: `{"statusCode": 404, "data": {"success": false, "error": "User not found"}}`.

| Method | Path                | Description |
|--------|---------------------|-------------|
| GET    | /health             | Health check |
| GET    | /api/v1/users       | List users  |
| GET    | /api/v1/users/{id}  | Get user    |
| POST   | /api/v1/users       | Create user |
| PUT    | /api/v1/users/{id}  | Update user |
| DELETE | /api/v1/users/{id}  | Delete user |

The same five routes exist for `/api/v1/items`. They are registered only once a database-backed
repository is wired in `cmd/server/main.go`; until then only `/health` responds.
