# Graph Report - 7th-sense-server  (2026-10-16)

## Corpus Check
- 207 files · ~96,423 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 1435 nodes · 4525 edges · 53 communities (46 shown, 7 thin omitted)
- Extraction: 88% EXTRACTED · 12% INFERRED · 0% AMBIGUOUS · INFERRED: 553 edges (avg confidence: 0.8)
- Token cost: 95,381 input · 0 output

## Community Hubs (Navigation)
- HTTP Handlers Layer
- Integration & DTO Tests
- DTO Validation & Transactions
- Profiles & Shared Models
- Route Registration & Guards
- Catalog & Brand DTOs
- Org Location/Shop DTOs
- Production Schedules
- Engineering Principles & API Docs
- Machine DTOs
- Server Wiring (master data)
- Server Wiring (customer entities)
- Part Process DTOs
- User Cell Module
- User Supervisor Module
- User Line Module
- Part & Shift Repositories
- Branch Module
- Product Module
- User City Module
- Company Scope & Holidays
- Config, DB & Seeding
- Controller Module & db.Execute
- User Service
- Clientele Module
- User Country Module
- Machine Model Module
- CRUD Pattern Guidelines
- User Repository & Licence Service
- User & Theme DTOs
- Querier Repositories
- Address Module
- Authentication Flow
- Roles Module
- Licence Key DTOs
- Part Service
- Setup & Tooling Docs
- User State Service
- User State Repository
- Licence Key Repository
- Auth Session Repository
- Part Drawing Storage
- Layered Architecture Principle
- Address DTOs
- Entry Point & Logger
- Role Assignment Rules
- Testing Principles
- COALESCE Rule
- Logging & Observability
- Context Propagation
- Profile Before Optimizing
- Go Module Root

## God Nodes (most connected - your core abstractions)
1. `Execute()` - 157 edges
2. `SendResponse()` - 130 edges
3. `NewError()` - 121 edges
4. `run()` - 78 edges
5. `pathID()` - 71 edges
6. `AuthUser` - 59 edges
7. `decodeAndValidate()` - 52 edges
8. `Querier` - 51 edges
9. `currentUser()` - 47 edges
10. `targetUserID()` - 35 edges

## Surprising Connections (you probably didn't know these)
- `README Layered Layout (Node.js analogy)` --semantically_similar_to--> `Strict Layered Architecture (Handler -> Service -> Repository -> Model)`  [INFERRED] [semantically similar]
  README.md → .claude/CLAUDE.md
- `Golden Reference: User CRUD` --semantically_similar_to--> `Golden Reference Pattern (User CRUD)`  [INFERRED] [semantically similar]
  .claude/CLAUDE.md → 7th_server_skills.md.md
- `Steps for New Entity (Model, DTO, Repository, Service, Responses, Handler, Routes, Wiring)` --semantically_similar_to--> `CRUD Implementation Checklist (11 steps)`  [INFERRED] [semantically similar]
  .claude/CLAUDE.md → 7th_server_skills.md.md
- `Soft Delete (is_deleted, deleted_at)` --semantically_similar_to--> `Soft Delete is Non-Negotiable`  [INFERRED] [semantically similar]
  .claude/CLAUDE.md → 7th_server_skills.md.md
- `COALESCE for Nullable Columns` --semantically_similar_to--> `Nullable Columns Must Use COALESCE`  [INFERRED] [semantically similar]
  .claude/CLAUDE.md → 7th_server_skills.md.md

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Company-scoped customer entities sharing companyOwner access** — _claude_claude_parts_module, _claude_claude_clienteles_module, _claude_claude_shifts_module, _claude_claude_breaks_module, _claude_claude_holidays_module, _claude_claude_licence_keys_module, _claude_claude_production_schedules_module, _claude_claude_company_scoping [EXTRACTED 1.00]
- **Session authentication and authorization flow** — readme_session_cookie, readme_auth_tokens_hashed, readme_per_request_session_check, readme_login_lockout, readme_csrf_cross_origin, _claude_claude_role_guards [INFERRED 0.85]
- **Golden CRUD pattern rules** — 7th_server_skills_md_golden_reference_pattern, 7th_server_skills_md_soft_delete_pattern, 7th_server_skills_md_coalesce_nullables, 7th_server_skills_md_error_mapping_consistency, 7th_server_skills_md_validation_at_boundary, 7th_server_skills_md_concrete_types_over_interfaces, 7th_server_skills_md_db_execute [EXTRACTED 1.00]

## Communities (53 total, 7 thin omitted)

### Community 0 - "HTTP Handlers Layer"
Cohesion: 0.05
Nodes (39): net/http.Request, net/http.ResponseWriter, validator, BranchHandler, BrandHandler, BreakHandler, CategoryHandler, ClienteleHandler (+31 more)

### Community 1 - "Integration & DTO Tests"
Cohesion: 0.07
Nodes (93): encoding/json.RawMessage, net/http.Cookie, net/http/httptest.ResponseRecorder, testing.T, CreateUserRequest, TestCreateUserRequestValidate(), TestCreateUserRequestValidateDropsDuplicateMachines(), TestCreateUserRequestValidateNormalizes() (+85 more)

### Community 2 - "DTO Validation & Transactions"
Cohesion: 0.07
Nodes (33): github.com/jackc/pgx/v5/pgxpool.Pool, pgx.Tx, WithTransaction(), BreakFields, CreateBreakRequest, UpdateBreakRequest, clockMinutes(), CreateShiftRequest (+25 more)

### Community 3 - "Profiles & Shared Models"
Cohesion: 0.06
Nodes (26): sync.Mutex, time.Time, CompanyProfile, UserProfile, CompanyProfileRepository, mapCompanyProfileError(), NewCompanyProfileRepository(), UserProfileRepository (+18 more)

### Community 4 - "Route Registration & Guards"
Cohesion: 0.06
Nodes (33): net/http.HandlerFunc, net/http.ServeMux, CurrentUser(), RequireRole(), UserID(), registerBranchRoutes(), registerBrandRoutes(), registerBreakRoutes() (+25 more)

### Community 5 - "Catalog & Brand DTOs"
Cohesion: 0.06
Nodes (20): CreateBrandRequest, UpdateBrandRequest, validateCatalogFields(), CreateCategoryRequest, UpdateCategoryRequest, CreateControllerRequest, UpdateControllerRequest, validateControllerFields() (+12 more)

### Community 6 - "Org Location/Shop DTOs"
Cohesion: 0.06
Nodes (19): CreateUserLocationRequest, UpdateUserLocationRequest, validateUserLocation(), CreateUserShopRequest, UpdateUserShopRequest, validateUserShop(), NewUserShopHandler(), UserLocation (+11 more)

### Community 7 - "Production Schedules"
Cohesion: 0.10
Nodes (20): validDate(), CreateProductionScheduleRequest, ProductionScheduleFields, ProductionScheduleItemFields, UpdateProductionScheduleRequest, ScheduleDay(), ProductionSchedule, ProductionScheduleItem (+12 more)

### Community 8 - "Engineering Principles & API Docs"
Cohesion: 0.06
Nodes (42): Database Queries are Sacred (no N+1, pooling, short transactions), db.Execute Query Helper, Defense in Depth, Idempotent PUT and DELETE, Never Trust Input (prepared statements, no sensitive logging), Transaction Boundaries in Service Layer, Validation at the Boundary (DTO Validate), model.AssignableRoles (+34 more)

### Community 9 - "Machine DTOs"
Cohesion: 0.10
Nodes (19): CreateMachineRequest, MachineFields, UpdateMachineRequest, validCoordinate(), validHealthIDs(), validIPv4(), Machine, MachineScope (+11 more)

### Community 10 - "Server Wiring (master data)"
Cohesion: 0.11
Nodes (30): run(), NewBranchHandler(), NewBrandHandler(), NewCategoryHandler(), NewControllerHandler(), NewMachineModelHandler(), NewProductHandler(), NewUserCityHandler() (+22 more)

### Community 11 - "Server Wiring (customer entities)"
Cohesion: 0.09
Nodes (22): net/http.Handler, NewBreakHandler(), NewClienteleHandler(), NewHolidayHandler(), NewLicenceKeyHandler(), NewMachineHandler(), NewPartHandler(), NewProductionScheduleHandler() (+14 more)

### Community 12 - "Part Process DTOs"
Cohesion: 0.11
Nodes (16): PartProcessFields, PartSubProgramFields, ClienteleFields, CreateClienteleRequest, UpdateClienteleRequest, validClienteleEmail(), CreateHolidayRequest, HolidayFields (+8 more)

### Community 13 - "User Cell Module"
Cohesion: 0.12
Nodes (11): CreateUserCellRequest, UpdateUserCellRequest, validateUserCell(), NewUserCellHandler(), UserCell, UserCellRepository, mapUserCellError(), NewUserCellRepository() (+3 more)

### Community 14 - "User Supervisor Module"
Cohesion: 0.13
Nodes (11): CreateUserSupervisorRequest, UpdateUserSupervisorRequest, validateUserSupervisor(), NewUserSupervisorHandler(), UserSupervisor, UserSupervisorRepository, mapUserSupervisorError(), NewUserSupervisorRepository() (+3 more)

### Community 15 - "User Line Module"
Cohesion: 0.13
Nodes (10): CreateUserLineRequest, UpdateUserLineRequest, validateUserLine(), UserLine, UserLineRepository, mapUserLineError(), NewUserLineRepository(), UserLineService (+2 more)

### Community 16 - "Part & Shift Repositories"
Cohesion: 0.16
Nodes (4): context.Context, executeAffected(), PartRepository, mapPartError()

### Community 17 - "Branch Module"
Cohesion: 0.14
Nodes (9): CreateBranchRequest, UpdateBranchRequest, validateBranchFields(), Branch, BranchRepository, mapBranchError(), cityOwner(), BranchService (+1 more)

### Community 18 - "Product Module"
Cohesion: 0.14
Nodes (11): CreateProductRequest, UpdateProductRequest, validateProductFields(), Product, ProductRepository, mapProductError(), NewProductRepository(), ProductService (+3 more)

### Community 19 - "User City Module"
Cohesion: 0.14
Nodes (8): CreateUserCityRequest, UpdateUserCityRequest, validateUserCity(), UserCity, UserCityRepository, mapUserCityError(), UserCityService, toUserCityResponseError()

### Community 20 - "Company Scope & Holidays"
Cohesion: 0.18
Nodes (8): AuthUser, Holiday, companyOwner(), ownCompanyCustomer(), HolidayService, holidayFromFields(), NewHolidayService(), toHolidayResponseError()

### Community 21 - "Config, DB & Seeding"
Cohesion: 0.18
Nodes (20): main(), run(), Config, testing.M, getEnv(), DBConfig, Load(), Connect() (+12 more)

### Community 22 - "Controller Module & db.Execute"
Cohesion: 0.16
Nodes (7): Execute(), Controller, ControllerRepository, mapControllerError(), ControllerService, toControllerResponseError(), pgx.Rows

### Community 23 - "User Service"
Cohesion: 0.22
Nodes (13): UserResponse, canSee(), companyProfileFrom(), UserService, NewUserService(), newUserServices(), primaryAddress(), setCustomerAddress() (+5 more)

### Community 24 - "Clientele Module"
Cohesion: 0.17
Nodes (8): Clientele, ClienteleRepository, mapClienteleError(), NewClienteleRepository(), clienteleFromFields(), ClienteleService, NewClienteleService(), toClienteleResponseError()

### Community 25 - "User Country Module"
Cohesion: 0.13
Nodes (8): CreateUserCountryRequest, NewUserCountryHandler(), UserCountry, UserCountryRepository, mapUserCountryError(), UserCountryService, NewUserCountryService(), toUserCountryResponseError()

### Community 26 - "Machine Model Module"
Cohesion: 0.16
Nodes (6): MachineModelRequest, MachineModel, MachineModelRepository, mapMachineModelError(), MachineModelService, toMachineModelResponseError()

### Community 27 - "CRUD Pattern Guidelines"
Cohesion: 0.10
Nodes (22): Code Reviews Are for Patterns, Common Mistakes to Avoid, Concrete Repository Structs Over Interfaces, Consistency Over Individual Preferences, CRUD Implementation Checklist (11 steps), Error Handling is First-Class, Repository -> Response Error Mapping, Golden Reference Pattern (User CRUD) (+14 more)

### Community 28 - "User Repository & Licence Service"
Cohesion: 0.16
Nodes (8): UserListFilter, User, UserRepository, mapUserError(), LicenceKeyService, licenceExpiresOn(), NewLicenceKeyService(), toLicenceKeyResponseError()

### Community 29 - "User & Theme DTOs"
Cohesion: 0.20
Nodes (14): ThemeResponse, UpdateThemeRequest, CustomerDetails, UpdateUserRequest, isTimeZone(), isWebsite(), normalizeFilters(), normalizeTheme() (+6 more)

### Community 30 - "Querier Repositories"
Cohesion: 0.16
Nodes (10): Querier, HolidayRepository, mapHolidayError(), NewHolidayRepository(), mapLocationError(), NewLocationRepository(), execute(), UserMachineAccessRepository (+2 more)

### Community 31 - "Address Module"
Cohesion: 0.18
Nodes (8): Address, AddressRepository, mapAddressError(), NewAddressRepository(), NewAddressService(), toAddressResponseError(), AddressService, TestAddressCRUD()

### Community 32 - "Authentication Flow"
Cohesion: 0.18
Nodes (11): LoginRequest, AuthHandler, NewAuthHandler(), Authenticate(), AuthService, hashToken(), NewAuthService(), newSessionToken() (+3 more)

### Community 33 - "Roles Module"
Cohesion: 0.19
Nodes (9): RoleHandler, NewRoleHandler(), Role, RoleRepository, mapRoleError(), NewRoleRepository(), RoleService, NewRoleService() (+1 more)

### Community 34 - "Licence Key DTOs"
Cohesion: 0.17
Nodes (10): LicenceKeyDuration, NewLicenceKey, CreateLicenceKeysRequest, GenerateLicenceKeysRequest, UpdateLicenceKeyRequest, licenceKeyDurations(), validDuration(), ValidLicenceKey() (+2 more)

### Community 35 - "Part Service"
Cohesion: 0.30
Nodes (10): Part, PartDrawing, PartProcess, PartSubProgram, NewPartRepository(), PartService, NewPartService(), partFromFields() (+2 more)

### Community 36 - "Setup & Tooling Docs"
Cohesion: 0.22
Nodes (11): Documentation at Implementation Time, Environment Parity (externalized config), Tech Stack (Go, PostgreSQL + TimescaleDB, pgx/v5, Air), Air as go.mod Tool (hot reload), VS Code Debugging with Delve (launch.json configs), .env from .env.example (config defaults), GOARCH=386 Pitfall, Install Go 1.27.1+ (+3 more)

### Community 37 - "User State Service"
Cohesion: 0.29
Nodes (4): CreateUserStateRequest, validateUserState(), UserStateService, toUserStateResponseError()

### Community 38 - "User State Repository"
Cohesion: 0.25
Nodes (3): UserState, UserStateRepository, mapUserStateError()

### Community 39 - "Licence Key Repository"
Cohesion: 0.38
Nodes (4): LicenceKey, LicenceKeyRepository, mapLicenceKeyError(), NewLicenceKeyRepository()

### Community 40 - "Auth Session Repository"
Cohesion: 0.31
Nodes (4): time.Duration, AuthRepository, mapAuthError(), NewAuthRepository()

### Community 41 - "Part Drawing Storage"
Cohesion: 0.33
Nodes (3): os.File, PartDrawingStorage, NewPartDrawingStorage()

### Community 42 - "Layered Architecture Principle"
Cohesion: 0.40
Nodes (5): Dependency Direction (inward) and Unidirectional Data Flow, Layered Architecture is Sacred, Predictability Over Cleverness, Strict Layered Architecture (Handler -> Service -> Repository -> Model), README Layered Layout (Node.js analogy)

### Community 45 - "Role Assignment Rules"
Cohesion: 0.67
Nodes (3): AssignableRoles(), CanAssignRole(), SeesAllCompanyMachines()

### Community 46 - "Testing Principles"
Cohesion: 0.67
Nodes (3): Test Behavior, Not Implementation, Integration Tests (tests/ setup helpers), Test Auth (asTestUser, tests/auth_test.go)

## Ambiguous Edges - Review These
- `Role Route Guards (superAdminOnly, managersOnly)` → `middleware.RequireRole / routes.adminOnly`  [AMBIGUOUS]
  .claude/CLAUDE.md · relation: conceptually_related_to
- `Tech Stack (Go, PostgreSQL + TimescaleDB, pgx/v5, Air)` → `Documentation at Implementation Time`  [AMBIGUOUS]
  .claude/CLAUDE.md · relation: references
- `Tech Stack (Go, PostgreSQL + TimescaleDB, pgx/v5, Air)` → `PostgreSQL Install and Migration`  [AMBIGUOUS]
  docs/setup.md · relation: conceptually_related_to

## Knowledge Gaps
- **21 isolated node(s):** `7th-sense-server`, `ThemeResponse`, `authUserKey`, `body`, `machineModelRecord` (+16 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **7 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **What is the exact relationship between `Role Route Guards (superAdminOnly, managersOnly)` and `middleware.RequireRole / routes.adminOnly`?**
  _Edge tagged AMBIGUOUS (relation: conceptually_related_to) - confidence is low._
- **What is the exact relationship between `Tech Stack (Go, PostgreSQL + TimescaleDB, pgx/v5, Air)` and `Documentation at Implementation Time`?**
  _Edge tagged AMBIGUOUS (relation: references) - confidence is low._
- **What is the exact relationship between `Tech Stack (Go, PostgreSQL + TimescaleDB, pgx/v5, Air)` and `PostgreSQL Install and Migration`?**
  _Edge tagged AMBIGUOUS (relation: conceptually_related_to) - confidence is low._
- **Why does `NewError()` connect `User & Theme DTOs` to `HTTP Handlers Layer`, `Integration & DTO Tests`, `DTO Validation & Transactions`, `Profiles & Shared Models`, `Catalog & Brand DTOs`, `Org Location/Shop DTOs`, `Production Schedules`, `Machine DTOs`, `Part Process DTOs`, `User Cell Module`, `User Supervisor Module`, `User Line Module`, `Branch Module`, `Product Module`, `User City Module`, `Company Scope & Holidays`, `Controller Module & db.Execute`, `User Service`, `Clientele Module`, `User Country Module`, `Machine Model Module`, `User Repository & Licence Service`, `Authentication Flow`, `Roles Module`, `Licence Key DTOs`, `Part Service`, `User State Service`, `Address DTOs`?**
  _High betweenness centrality (0.181) - this node is a cross-community bridge._
- **Why does `run()` connect `Server Wiring (master data)` to `HTTP Handlers Layer`, `DTO Validation & Transactions`, `Catalog & Brand DTOs`, `Org Location/Shop DTOs`, `Production Schedules`, `Machine DTOs`, `Server Wiring (customer entities)`, `User Cell Module`, `User Supervisor Module`, `User Line Module`, `Product Module`, `Company Scope & Holidays`, `Config, DB & Seeding`, `User Service`, `Clientele Module`, `User Country Module`, `User Repository & Licence Service`, `Querier Repositories`, `Authentication Flow`, `Roles Module`, `Part Service`, `Licence Key Repository`, `Auth Session Repository`, `Part Drawing Storage`, `Entry Point & Logger`?**
  _High betweenness centrality (0.051) - this node is a cross-community bridge._
- **Why does `CurrentUser()` connect `Route Registration & Guards` to `HTTP Handlers Layer`, `Part & Shift Repositories`, `Company Scope & Holidays`?**
  _High betweenness centrality (0.046) - this node is a cross-community bridge._
- **What connects `7th-sense-server`, `ThemeResponse`, `authUserKey` to the rest of the system?**
  _21 weakly-connected nodes found - possible documentation gaps or missing edges._