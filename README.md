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

`POST /api/v1/users` creates a customer: the user, its company profile, user profile and
address are saved in one transaction (all or nothing).

```json
{"username": "hiren_c", "email": "hiren@example.com", "password": "secret123", "roleId": 4,
 "branchId": 1, "companyName": "Enliven Pvt Ltd", "finixId": 77,
 "gstNumber": "24AAACE1234F1Z5", "panNumber": "AAACE1234F",
 "website": "www.example.com", "companyLogo": "uploads/logo.png",
 "addressLine1": "Plot 12, GIDC", "addressLine2": "", "country": "India", "state": "Gujarat",
 "city": "Rajkot", "zipCode": "360001",
 "contactPerson": "Hiren", "primaryPhone": "+919876543210", "secondaryPhone": "",
 "timeZone": "Asia/Kolkata", "dayStartFrom": "06:00"}
```

| Form field      | JSON             | Stored in                              |
|-----------------|------------------|----------------------------------------|
| Branch *        | `branchId`       | `users.branch_id` (must exist)         |
| Company Name *  | `companyName`    | `company_profiles.company_name`        |
| Finix ID        | `finixId`        | `company_profiles.finix_id`            |
| GST No. *       | `gstNumber`      | `company_profiles.gst_number`          |
| PAN No. *       | `panNumber`      | `company_profiles.pan_number`          |
| Website         | `website`        | `company_profiles.website`             |
| Company Logo    | `companyLogo`    | `user_profiles.profile_picture` (path/URL) |
| Email *         | `email`          | `users.email` (unique, lowercase)      |
| User Name *     | `username`       | `users.username` (unique)              |
| Password *      | `password`       | `users.password` (bcrypt, 8–72 chars)  |
| Address - 1 *   | `addressLine1`   | `addresses.address_line1`              |
| Address - 2     | `addressLine2`   | `addresses.address_line2`              |
| Country *       | `country`        | `addresses.country_id` (looked up by name) |
| State *         | `state`          | `addresses.state_id` (looked up by name, within the country) |
| City *          | `city`           | `addresses.city`                       |
| Zip Code *      | `zipCode`        | `addresses.pin_code` (3–10 digits)     |
| Contact Person *| `contactPerson`  | `user_profiles.contact_person`, `addresses.contact_person` |
| Primary Phone * | `primaryPhone`   | `users.mobile`, `addresses.contact_phone` |
| Secondary Phone | `secondaryPhone` | `user_profiles.contact_phone`          |
| Time Zone *     | `timeZone`       | `user_profiles.timezone` (IANA name)   |
| Day Start From *| `dayStartFrom`   | `user_profiles.start_day_time` (`HH:MM`) |

`filters` is the filter checklist: any of `country`, `state`, `branch`, `city`, `location`, `shop`, `line`, `cell`, `supervisor` (case-insensitive, duplicates ignored). It is stored in `user_filters`, one row per filter; on update the list replaces the chosen filters and `[]` clears them. Responses include it as `"filters": [...]`.

`roleId` is required (4 = customer); `firstName`, `lastName` and `companyId` are optional.

`GET /api/v1/users/{id}` and the create/update responses return the user with `profile`,
`companyProfile` and `address` objects (`null` when missing). `GET /api/v1/users` lists users only.

`PUT /api/v1/users/{id}` takes the same fields without `password`, plus an optional
`status` (`active` / `inactive`; empty keeps the current one). It replaces the form fields;
user profile fields that are not on the form (designation, birth date, theme, ...) are kept.

`DELETE` is a soft delete of the user, its profiles and addresses: the rows stay with `is_deleted = true` and `deleted_at` set, and is
hidden from every other endpoint. Its email and username can be used again by a new user.

The same five routes exist for `/api/v1/items`. They are registered only once a database-backed
repository is wired in `cmd/server/main.go`; until then only `/health` responds.
