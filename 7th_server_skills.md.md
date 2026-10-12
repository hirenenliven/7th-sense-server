# 7th Sense Server - Backend API Project

## Engineering Principles

This project follows strict software engineering principles forged from 10+ years of production experience across enterprise Go applications serving millions of requests. These patterns are not theoretical—they are battle-tested solutions to real problems encountered in production at scale.

### Core Philosophy

**Predictability Over Cleverness**
- Code should be boring and predictable, not clever and surprising
- A junior developer should be able to understand any file within 5 minutes
- If you need a comment to explain "why" your code is clever, it's probably wrong
- The best code is the code you don't have to think about twice when reading

**Consistency Over Individual Preferences**
- Team patterns trump personal coding style
- If the pattern is X, don't invent Y unless there's a documented reason
- Copy-paste is a virtue when it means consistency
- Unnecessary innovation is technical debt in disguise

**Simplicity Over Complexity**
- Solve the problem at hand, not future problems that may never materialize
- YAGNI (You Aren't Gonna Need It) is a principle, not a suggestion
- If you can't explain it simply, you don't understand it well enough
- Complex solutions to simple problems are a red flag

### Architectural Principles

**1. Layered Architecture is Sacred**
- Each layer has a single responsibility
- Never violate layer boundaries (no cross-layer imports)
- Handlers know nothing about databases
- Repositories know nothing about HTTP
- Services know nothing about JSON parsing
- This separation enables testing, mocking, and independent evolution

**2. Dependency Direction**
- Dependencies flow inward (toward the core)
- Handler → Service → Repository → Model
- Never the reverse
- Circular dependencies are design failures, not something to work around

**3. Data Flow is Unidirectional**
- Request → Handler → Service → Repository → Database
- Response flows back the same way
- No shortcuts, no side channels
- If you need to skip a layer, your design is wrong

### Code Quality Principles

**4. Readability First**
- Code is read 10x more than it's written
- Variable names should tell a story
- Functions should do one thing and do it well
- Deep nesting is a code smell—refactor early
- If it's hard to understand, it's wrong

**5. Error Handling is First-Class**
- Errors are not an afterthought
- Every error must be handled explicitly
- Never ignore errors (empty if statements are evil)
- Error messages must be actionable
- Wrap errors with context, don't lose the stack

**6. Validation at the Boundary**
- Validate as early as possible (DTO layer)
- Never trust user input
- Fail fast—reject invalid requests before business logic
- Validation logic is reusable, scattered validation is not

### Testing Principles

**7. Test Behavior, Not Implementation**
- Test what the code does, not how it does it
- Mock external dependencies, not internal logic
- Tests should be resilient to refactoring
- If a test breaks when you rename a variable, it's testing implementation

**8. Tests Are Documentation**
- Good tests serve as executable documentation
- Test names should describe the scenario
- Setup should be clear and minimal
- If a test needs a comment to explain what it's testing, rewrite the test

### Performance Principles

**9. Premature Optimization is Evil**
- Profile before optimizing
- The bottleneck is rarely where you think it is
- Readable code is usually fast enough
- Cache only when you have measured the need

**10. Database Queries are Sacred**
- N+1 queries are unacceptable
- Always query exactly what you need
- Use indexes, but don't over-index
- Connection pooling is non-negotiable
- Transactions should be as short as possible

### Security Principles

**11. Never Trust Input**
- Validate everything from the outside
- Sanitize before storing
- Never concatenate SQL strings (use prepared statements)
- Never log sensitive data (passwords, tokens, PII)
- Security is a mindset, not a checklist

**12. Defense in Depth**
- Validate at multiple layers (DTO, service, database constraints)
- Don't rely on a single security control
- Assume every layer can fail
- Auditing is not optional—log who did what when

### Team Collaboration Principles

**13. Code Reviews Are for Patterns**
- Focus on adherence to established patterns
- Style debates are a waste of time
- If it matches the pattern, approve it
- If it doesn't, ask why and document the reason

**14. Documentation at Implementation Time**
- Write documentation when you write the code
- "I'll document it later" means "I'll never document it"
- If it's not documented, it doesn't exist
- README files are for humans, not parsers

### DevOps Principles

**15. Environment Parity**
- Development, staging, and production should be identical
- Use Docker or environment-specific configs
- "It works on my machine" is not an acceptable excuse
- Configuration should be externalized, never hardcoded

**16. Monitoring and Observability**
- If you can't measure it, you can't improve it
- Log structured data, not free-form text
- Metrics should be actionable
- Alerts should be for things that require human intervention

### The "Why" Behind Strict Enforcement

**Why so strict?**
Because I've seen what happens when teams aren't:

- "It's just a small shortcut" → accumulates into technical debt
- "This is a special case" → becomes the new pattern accidentally
- "I'll refactor it later" → later never comes
- "The deadline is tight" → rushed code lives forever in production

**Violations have consequences:**
- Technical debt that requires 10x the time to fix
- Bugs that surface in production because patterns weren't followed
- Onboarding time increases for new developers
- Code reviews become debates instead of approvals
- The codebase becomes a patchwork of incompatible styles

**The payoff for discipline:**
- New features ship faster because patterns are established
- Bugs are caught early because patterns are proven
- Onboarding takes days instead of weeks
- Code reviews are fast and focused
- The codebase ages gracefully instead of rotting

### Your Responsibility

As a developer on this project:
1. **Read the patterns** before writing code
2. **Ask questions** if you don't understand why something is done a certain way
3. **Follow the patterns** even when you think you have a "better way"
4. **Propose changes** through discussion, not by implementing first and asking later
5. **Document deviations** with explicit reasons if approved

This is not about stifling creativity—it's about channeling creativity into solving business problems, not reinventing software engineering fundamentals.

**Discipline in the small things creates freedom in the big things.**

## Golden Reference Pattern
**User CRUD** (`internal/handler/user_handler.go`, `internal/service/user_service.go`, `internal/repository/user_repository.go`) is the **canonical reference implementation** for all CRUD operations. Every new entity must mirror this pattern exactly. This is not a suggestion—it's a requirement for code consistency and maintainability.

## Why This Architecture Matters
After years of maintaining monolithic Go services, I've learned that:
- **Strict layer separation** prevents circular dependencies and makes testing easy
- **Concrete repository structs** (not interfaces) eliminate unnecessary abstraction
- **Consistent error handling** makes debugging in production predictable
- **Soft delete patterns** prevent data loss and enable audit trails
- **COALESCE for nullables** eliminates nil pointer panics in JSON responses

## Project Overview
Go-based REST API backend for the 7th Sense system. Uses PostgreSQL with pgx driver and TimescaleDB for time-series data, following a strict layered architecture pattern designed for scalability and maintainability.

## Tech Stack
- **Language**: Go
- **Database**: PostgreSQL + TimescaleDB
- **Driver**: pgx/v5
- **Architecture**: Layered (Handler → Service → Repository → Model)

## Project Structure
```
7th-sense-server/
├── cmd/server/          # Application entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── db/              # Database connection utilities
│   ├── dto/             # Data Transfer Objects (request bodies with validation)
│   ├── handler/         # HTTP request handlers
│   ├── middleware/      # HTTP middleware (logging, recovery)
│   ├── model/           # Database models (structs matching tables)
│   ├── repository/      # SQL queries and data access
│   ├── response/        # API response definitions
│   ├── routes/          # Route registration
│   └── service/         # Business logic layer
├── migrations/          # Database migration files
└── logs/                # Application logs
```

## Engineering Best Practices (Follow These Strictly)

### 1. Never Skip Layers
- Handlers → Services → Repositories is the only allowed call chain
- Never call repository directly from handler
- Never call handler from service
- Never import lower layers into higher layers (e.g., service importing handler)
- **Why**: Prevents tight coupling, makes unit testing possible, enables easy mocking

### 2. Concrete Types Over Interfaces
- Use concrete structs for repositories (not interfaces)
- Only use interfaces when you need multiple implementations (rare)
- **Why**: YAGNI principle. Interfaces add complexity without benefit in 99% of cases. Testing with concrete types is simpler and more realistic.

### 3. Error Handling Consistency
- Repository errors → Response errors mapping is mandatory
- Never return raw database errors to HTTP layer
- Always wrap errors with context (operation name)
- **Why**: Prevents leaking implementation details, provides actionable error messages to clients, makes logging useful

### 4. Nullable Columns Must Use COALESCE
- Always use `COALESCE(column, '') AS column` in SELECT
- Never return nil for string columns in JSON
- **Why**: Prevents null values in JSON, eliminates frontend null-checking chaos, matches typical API expectations

### 5. Soft Delete is Non-Negotiable
- Every table must have `is_deleted` and `deleted_at`
- Unique indexes must include `WHERE is_deleted = false`
- Never hard delete in production
- **Why**: Enables audit trails, prevents accidental data loss, allows recovery, satisfies compliance requirements

### 5.1 Response and Error Handling Pattern
**Reference**: `internal/response/responses.go`

**Engineering Notes**:
- Use predefined Response objects, never create ad-hoc responses
- Response implements error interface, so can be returned directly
- Always use entity-specific responses (e.g., UserNotFound, not generic NotFound)
- Error messages must be human-readable for clients, not technical
- Status codes must match HTTP semantics (400 for client errors, 500 for server errors)

**Response Pattern**:
```go
// In service layer - return response errors
func toEntityResponseError(err error) error {
    switch {
    case errors.Is(err, repository.ErrNotFound):
        return response.EntityNotFound  // 404, "Entity not found"
    case errors.Is(err, repository.ErrDuplicateField):
        return response.FieldExists  // 400, "Field already exists"
    default:
        return err  // Let other errors pass through (becomes 500)
    }
}

// In handler layer - use response.SendResponse
response.SendResponse(writer, request, err, data)
// If err is a Response, it returns that status code and message
// If err is not a Response, it returns 500 Internal Server Error
```

**Creating New Entity-Specific Responses**:
```go
// In internal/response/responses.go
var (
    EntityNotFound = newResponse(http.StatusNotFound, "Entity not found")
    FieldExists = newResponse(http.StatusBadRequest, "Field already exists")
    ValidField = newResponse(http.StatusBadRequest, "Please enter valid field")
)
```

**Message Guidelines**:
- Use "Please enter valid X" for validation errors
- Use "X not found" for missing resources
- Use "X already exists" for duplicate violations
- Never expose database error messages to clients
- Keep messages short and actionable

**Why**: Consistent error handling makes frontend development predictable, prevents information leakage, and provides good UX.

### 6. Transaction Boundaries
- Transactions go in service layer, not repository
- Keep transactions as short as possible
- Rollback on any error
- **Why**: Prevents deadlocks, ensures data consistency, makes retry logic simple

### 7. Validation in DTOs Only
- Never validate in handlers or services
- DTO `Validate()` method is the single source of truth
- Return response errors for validation failures
- **Why**: Reusable validation logic, consistent error messages, testable in isolation

### 8. Idempotent Operations
- PUT should be idempotent (same result on multiple calls)
- DELETE should be idempotent (already deleted = success)
- Handle concurrent updates gracefully
- **Why**: Makes retry logic work, prevents duplicate data, essential for distributed systems

### 9. Context Propagation
- Always pass context through all layers
- Use context for cancellation and timeouts
- Never ignore context in database calls
- **Why**: Enables request cancellation, prevents resource leaks, makes observability possible

### 10. Logging Consistency
- Log at service layer for business events
- Log at repository layer for SQL queries (in development)
- Never log passwords or sensitive data
- **Why**: Debugging in production, performance analysis, security compliance

## Architecture Pattern

### Layered Architecture (Strict Separation)
1. **Handler Layer** (`internal/handler/`)
   - HTTP request/response handling
   - Uses `decodeAndValidate` for request body parsing
   - Calls service layer
   - Returns responses via `response.SendResponse`
   - **Reference**: `user_handler.go`

2. **Service Layer** (`internal/service/`)
   - Business logic (password hashing, error mapping)
   - Orchestrates repository calls
   - Maps repository errors to response errors
   - No HTTP concerns
   - **Reference**: `user_service.go`

3. **Repository Layer** (`internal/repository/`)
   - SQL queries only using `db.Execute()`
   - No business logic
   - Uses COALESCE for nullable columns
   - Returns repository errors (ErrNotFound, ErrDuplicate*)
   - **Reference**: `user_repository.go`

4. **Model Layer** (`internal/model/`)
   - Structs matching database tables
   - DB tags for column mapping
   - JSON tags for API responses

5. **DTO Layer** (`internal/dto/`)
   - Request body structs
   - Implements `Validate()` method
   - Data validation (regex, length checks)
   - **Reference**: `user_dto.go`

## CRUD Implementation Checklist (For Every New Entity)

**Critical**: Before implementing any CRUD, read and understand the User CRUD files completely. Do not start coding until you understand the pattern. Copy-paste the structure and adapt it for your entity. This is not the time for creativity—consistency beats cleverness every time.

### Step 1: Model (`internal/model/entity.go`)
**Reference**: `internal/model/user.go`
```go
type Entity struct {
    ID         int64      `db:"id" json:"id"`
    Field1     string     `db:"field1" json:"field1"`
    Field2     *string    `db:"field2" json:"field2"`  // nullable
    IsDeleted  bool       `db:"is_deleted" json:"isDeleted"`
    DeletedAt  *time.Time `db:"deleted_at" json:"deletedAt"`
    CreatedAt  time.Time  `db:"created_at" json:"createdAt"`
    UpdatedAt  time.Time  `db:"updated_at" json:"updatedAt"`
}
```

### Step 2: DTO (`internal/dto/entity_dto.go`)
**Reference**: `internal/dto/user_dto.go`
```go
type CreateEntityRequest struct {
    Field1 string `json:"field1"`
    Field2 string `json:"field2"`
}

func (request *CreateEntityRequest) Validate() error {
    // Trim and validate fields
    // Return response errors for validation failures
}

type UpdateEntityRequest struct {
    Field1 string `json:"field1"`
    Field2 string `json:"field2"`
    Status string `json:"status"`  // optional
}

func (request *UpdateEntityRequest) Validate() error {
    // Trim and validate fields
}
```

### Step 3: Repository (`internal/repository/entity_repository.go`)
**Reference**: `internal/repository/user_repository.go`

**Engineering Notes**:
- Repository methods are the ONLY place SQL belongs
- Never add business logic here (no if/else beyond error mapping)
- Use prepared statements via pgx (handled by db.Execute)
- Order of parameters in db.Execute must match $1, $2, $3 in query
- Always use `ctx` as first parameter for context propagation
- Error messages must include operation name for debugging: "find user: %w"

```go
type EntityRepository struct {
    dbPool *pgxpool.Pool
}

func NewEntityRepository(dbPool *pgxpool.Pool) *EntityRepository {
    return &EntityRepository{dbPool: dbPool}
}

// Methods: FindAll, FindByID, Create, Update, Delete
// All methods must:
// - Use db.Execute(ctx, repo.dbPool, query, params...)
// - Use COALESCE for nullable columns in SELECT
// - Use WHERE is_deleted = false
// - Return model.Entity or error
// - Use mapEntityError for error mapping
```

### Step 4: Repository Errors (`internal/repository/errors.go`)
**Reference**: `internal/repository/errors.go`
```go
var (
    ErrNotFound = errors.New("record not found")
    ErrDuplicateField = errors.New("duplicate field")  // for unique constraints
)
```

### Step 5: Service (`internal/service/entity_service.go`)
**Reference**: `internal/service/user_service.go`

**Engineering Notes**:
- Service layer is where business logic lives (hashing, validation, orchestration)
- Never return repository errors directly to HTTP layer
- Error mapping function must be private (toEntityResponseError)
- Use pointer receiver for all methods
- Keep methods simple—if it's complex, break into smaller private functions
- Context should be passed through to repository calls

```go
type EntityService struct {
    entityRepository *repository.EntityRepository
}

func NewEntityService(entityRepository *repository.EntityRepository) *EntityService {
    return &EntityService{entityRepository: entityRepository}
}

// Methods: GetAll, GetByID, Create, Update, Delete
// All methods must:
// - Call repository methods
// - Apply business logic (hashing, etc.)
// - Map repository errors to response errors via toEntityResponseError
```

### Step 6: Response Errors (`internal/response/responses.go`)
**Reference**: `internal/response/responses.go`
```go
var (
    EntityNotFound = newResponse(http.StatusNotFound, "Entity not found")
    FieldExists = newResponse(http.StatusBadRequest, "Field already exists")
    ValidField = newResponse(http.StatusBadRequest, "Please enter valid field")
)
```

### Step 7: Handler (`internal/handler/entity_handler.go`)
**Reference**: `internal/handler/user_handler.go`

**Engineering Notes**:
- Handlers should be thin—just HTTP concerns
- Never do business logic here (no hashing, no validation beyond decode)
- Always check errors before proceeding
- Return early on errors (guard clauses)
- Use response.SendResponse for all responses (never manual JSON encoding)
- request.Context() provides cancellation support

```go
type EntityHandler struct {
    service *service.EntityService
}

func NewEntityHandler(entityService *service.EntityService) *EntityHandler {
    return &EntityHandler{service: entityService}
}

// Methods: GetAll, GetByID, Create, Update, Delete
// All methods must:
// - Use pathID() for ID extraction
// - Use decodeAndValidate() for request body
// - Call service methods
// - Use response.SendResponse() for responses
```

### Step 8: Routes (`internal/routes/entity_routes.go`)
**Reference**: `internal/routes/user_routes.go`
```go
func registerEntityRoutes(router *http.ServeMux, prefix string, entityHandler *handler.EntityHandler) {
    router.HandleFunc("GET "+prefix+"/entities", entityHandler.GetAll)
    router.HandleFunc("GET "+prefix+"/entities/{id}", entityHandler.GetByID)
    router.HandleFunc("POST "+prefix+"/entities", entityHandler.Create)
    router.HandleFunc("PUT "+prefix+"/entities/{id}", entityHandler.Update)
    router.HandleFunc("DELETE "+prefix+"/entities/{id}", entityHandler.Delete)
}
```

### Step 9: Register Routes (`internal/routes/routes.go`)
**Reference**: `internal/routes/routes.go`
```go
type Handlers struct {
    User   *handler.UserHandler
    Entity *handler.EntityHandler  // Add here
}

func registerV1Routes(router *http.ServeMux, handlers Handlers) {
    // ... existing routes
    if handlers.Entity != nil {
        registerEntityRoutes(router, v1Prefix, handlers.Entity)
    }
}
```

### Step 10: Wire Up (`cmd/server/main.go`)
**Reference**: `cmd/server/main.go`
```go
entityRepository := repository.NewEntityRepository(dbPool)
router := routes.New(routes.Handlers{
    User:   handler.NewUserHandler(service.NewUserService(userRepository)),
    Entity: handler.NewEntityHandler(service.NewEntityService(entityRepository)),  // Add here
})
```

### Step 11: Migration (`migrations/xxx_add_entity.sql`)
**Reference**: `migrations/schema.sql` (users table definition)
```sql
CREATE TABLE IF NOT EXISTS entities (
    id SERIAL PRIMARY KEY,
    field1 VARCHAR(100) NOT NULL,
    field2 VARCHAR(100),
    is_deleted BOOLEAN NOT NULL DEFAULT false,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Unique constraint if needed
CREATE UNIQUE INDEX IF NOT EXISTS entities_field1_key ON entities (field1) WHERE is_deleted = false;

-- Indexes
CREATE INDEX IF NOT EXISTS idx_entities_field1 ON entities(field1);
CREATE INDEX IF NOT EXISTS idx_entities_is_deleted ON entities(is_deleted);

-- Trigger for updated_at
CREATE TRIGGER update_entities_updated_at BEFORE UPDATE ON entities
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
```

## Code Conventions (MUST FOLLOW)

### Repository Pattern
```go
// Always use db.Execute for queries
rows, err := db.Execute(ctx, repo.dbPool, query, params...)

// Use COALESCE for nullable columns
SELECT COALESCE(field_name, '') AS field_name

// Soft delete pattern
WHERE is_deleted = false

// Unique constraint error mapping
var pgError *pgconn.PgError
if errors.As(err, &pgError) && pgError.Code == "23505" {
    switch pgError.ConstraintName {
    case "entities_field_key":
        return ErrDuplicateField
    }
}
```

### Service Pattern
```go
// Error mapping function (required)
func toEntityResponseError(err error) error {
    switch {
    case errors.Is(err, repository.ErrNotFound):
        return response.EntityNotFound
    case errors.Is(err, repository.ErrDuplicateField):
        return response.FieldExists
    default:
        return err
    }
}
```

### Handler Pattern
```go
// Decode and validate request
var request dto.CreateEntityRequest
if err := decodeAndValidate(writer, request, &request); err != nil {
    response.SendResponse(writer, request, err, nil)
    return
}

// Call service
result, err := service.Create(ctx, request)
response.SendResponse(writer, request, err, result)
```

## Database Connection
- Uses environment variables (.env file)
- Required: DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_DATABASE
- Connection pooling via pgxpool

## Soft Delete Pattern (MUST USE)
All tables must use soft delete:
- `is_deleted` boolean column (default false)
- `deleted_at` timestamp
- Repository methods exclude deleted records with `WHERE is_deleted = false`
- Delete operation sets `is_deleted=true` and `deleted_at=NOW()`
- Unique indexes must include `WHERE is_deleted = false`

## Important Rules (NEVER BREAK)
- **Never use repository interfaces** - use concrete structs with pgxpool
- **Always use db.Execute()** for database queries
- **Use COALESCE** for nullable columns in SELECT
- **Follow the exact pattern** of User CRUD (golden reference)
- **Add unique constraints** in migrations for error mapping
- **Soft delete everywhere** - use is_deleted flag
- **Every entity needs all 5 layers**: Model → DTO → Repository → Service → Handler
- **Always add error mapping function** in service layer

## Common Mistakes to Avoid (I've Seen These in Production)

### 1. "I'll just add this one field to the handler"
**Don't do it.** Add it to the DTO, validate it properly, then pass through service. Shortcuts accumulate into technical debt.

### 2. "This query is simple, I'll skip the repository"
**Don't do it.** All SQL belongs in repositories. Even simple queries get complex over time. Keep the layer separation.

### 3. "I don't need soft delete for this table"
**You do.** Always. Production requirements change, compliance audits happen, mistakes occur. Soft delete is insurance.

### 4. "I'll return the raw error, it's informative enough"
**Don't do it.** Never leak implementation details. Map repository errors to response errors. Your clients shouldn't know you use PostgreSQL.

### 5. "I'll add the trigger later"
**Do it now.** The updated_at trigger should be part of the migration. Forgetting it leads to inconsistent timestamps.

### 6. "This entity doesn't need validation"
**Everything needs validation.** Even if it's just checking for empty strings. Validate early, validate often.

### 7. "I'll make the repository an interface for testability"
**You don't need it.** Concrete types are more testable than you think. Mocking concrete repos is easy. Interfaces without multiple implementations are YAGNI.

### 8. "I'll skip the error mapping function"
**Don't skip it.** Copy-paste the toEntityResponseError pattern. It's boilerplate for a reason—consistency.

## Final Thoughts

After 10+ years of building Go services, here's what I've learned:

1. **Consistency beats cleverness** every time. A boring, consistent codebase is maintainable. A clever, inconsistent one is a nightmare.

2. **Layers exist for a reason**. When you're tempted to skip a layer, stop and think about why the layer exists. Usually, it's for a good reason.

3. **Error handling is not optional**. In production, you'll thank yourself for consistent error handling. Debugging is 10x harder without it.

4. **Soft delete is free insurance**. The day you need to recover data, you'll be glad you have it.

5. **Code reviews should focus on patterns**. Not style, not performance, but adherence to established patterns. That's what keeps a codebase healthy.

6. **New developers should read the reference files first**. Before touching any code, read user_handler.go, user_service.go, and user_repository.go. Understand the "why" before implementing the "what".

7. **If you're unsure, ask**. Don't guess. The patterns here are intentional. If you don't understand why something is done a certain way, ask the team. There's probably a story behind it.

This codebase is designed to scale with your team and your product. Following these patterns ensures that as the codebase grows, it remains maintainable, testable, and reliable. Deviate from them only with explicit team approval and a documented reason.

**Build software that your future self (and future teammates) will thank you for.**
