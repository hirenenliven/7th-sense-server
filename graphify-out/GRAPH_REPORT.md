# Graph Report - 7th-sense-server  (2026-10-20)

## Corpus Check
- 209 files · ~98,410 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 1564 nodes · 4815 edges · 70 communities (57 shown, 13 thin omitted)
- Extraction: 87% EXTRACTED · 13% INFERRED · 0% AMBIGUOUS · INFERRED: 604 edges (avg confidence: 0.81)
- Token cost: 68,371 input · 0 output

## Community Hubs (Navigation)
- HTTP Handlers Layer
- Integration & DTO Tests
- User Supervisor Module
- Org Location/Shop DTOs
- Route Registration & Guards
- Machine DTOs
- DTO Validation & Transactions
- CNC Snapshot Struct Fields
- Licence Key DTOs
- Clientele Module
- Production Schedules (2)
- Part Process DTOs
- Authentication Flow
- CNC DLL Artifacts & Fanuc
- Controller Module & db.Execute
- Server Wiring (customer entities)
- Holiday DTOs & Model
- Branch Module
- Catalog & Brand DTOs
- Product Module
- User Service
- Server Wiring (master data)
- Catalog & Brand DTOs (2)
- User City Module
- CNC Go Client Run & Check Mode
- Server Wiring (master data) (2)
- User & Theme DTOs
- Part Process DTOs (2)
- Profiles & Shared Models
- Machine Model Module
- DTO Validation & Transactions (2)
- Part Process DTOs (3)
- Address Module
- User Line Module
- CNC Machine Types & Catalog
- User Country Module
- Profiles & Shared Models (2)
- Roles Module
- CNC Connect API & Callbacks
- Server API Domain Docs
- Config, DB & Seeding
- User Country & Core Models
- Config, DB & Seeding (2)
- User State Module
- User State Module (2)
- Querier Repositories
- CRUD Pattern Guidelines
- Engineering Principles & API Docs (2)
- Setup & Tooling Docs
- Engineering Principles
- Layered Architecture & COALESCE
- Siemens Certificates & Trust
- DTO Validation & Transactions (4)
- Holiday DTOs & Model (2)
- CRUD Pattern Guidelines (3)
- CNC Event-Driven Data Flow
- User Line Module (2)
- Location Repository
- Project Context & Tech Stack
- Authentication Flow (2)
- Transactions & DB Query Rules
- Role Assignment Rules
- Go Module Root
- Current User Middleware
- Role Guard Middleware
- SendResponse Helper
- Test Auth Helpers

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
10. `CncMasterSnapshot (34,440 bytes, #pragma pack(8), x86)` - 38 edges

## Surprising Connections (you probably didn't know these)
- `Machine Identifier Mapping (0 Simulator,1 Fanuc,2 Siemens,3 Syntec,4 Mitsubishi)` --semantically_similar_to--> `Machine Catalog Brand -> Controller -> Model`  [INFERRED] [semantically similar]
  docs/CncCommunication.dll — API & Integration Guide (Go Client).pdf → .claude/CLAUDE.md
- `Snapshot Fields (positions, spindle, feed, program, tool, alarms)` --conceptually_related_to--> `Machines API (/machines)`  [AMBIGUOUS]
  docs/CncCommunication.dll — API & Integration Guide (Go Client).pdf → .claude/CLAUDE.md
- `README Layered Layout (Node.js analogy)` --semantically_similar_to--> `Layered Architecture (Handler -> Service -> Repository -> Model)`  [INFERRED] [semantically similar]
  README.md → .claude/CLAUDE.md
- `CNCConnect(ip, port, dataCb, errCb, machineType) int32` --semantically_similar_to--> `CNCConnect(ip, port, dataCb, errCb, machineType)`  [INFERRED] [semantically similar]
  docs/cnc-communication.md → docs/CncCommunication.dll — API & Integration Guide (Go Client).pdf
- `CncMasterSnapshot (34,440 bytes, #pragma pack(8), x86)` --semantically_similar_to--> `CncMasterSnapshot (34,440 bytes, #pragma pack(8))`  [INFERRED] [semantically similar]
  docs/cnc-communication.md → docs/CncCommunication.dll — API & Integration Guide (Go Client).pdf

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **CNC telemetry callback data flow** — docs_cnccommunication_dll_api_integration_guide_go_client_cncconnect, docs_cnccommunication_dll_api_integration_guide_go_client_stdcall_callbacks, docs_cnccommunication_dll_api_integration_guide_go_client_ondatacallback, docs_cnccommunication_dll_api_integration_guide_go_client_cncmastersnapshot, docs_cnccommunication_dll_api_integration_guide_go_client_event_driven_data_flow [EXTRACTED 1.00]
- **Company-scoped customer resources** — _claude_claude_parts_api, _claude_claude_clienteles_api, _claude_claude_shifts_breaks_holidays, _claude_claude_licence_keys_api, _claude_claude_production_schedules_api, _claude_claude_company_scope [EXTRACTED 1.00]
- **CRUD layered pattern rules** — _claude_claude_layered_architecture, _claude_claude_dto_validation, _claude_claude_service_error_mapping, _claude_claude_coalesce_nullables, _claude_claude_soft_delete, _claude_claude_user_crud_golden_reference [EXTRACTED 1.00]
- **Session authentication and authorization flow** — readme_session_cookie, readme_auth_tokens_hashed, readme_per_request_session_check, readme_login_lockout, readme_csrf_cross_origin, _claude_claude_role_guards [INFERRED 0.85]
- **Event-driven telemetry pipeline (CNCConnect -> worker -> callback -> mutex copy -> render)** — docs_cnc_communication_cncconnect, docs_cnc_communication_flow_step1_register_callbacks, docs_cnc_communication_flow_step2_worker_thread, docs_cnc_communication_flow_step3_poll_and_invoke, docs_cnc_communication_flow_step4_copy_and_signal, docs_cnc_communication_flow_step5_render_dashboard, docs_cnc_communication_ondatacallback, docs_cnc_communication_mutex_copy, docs_cnc_communication_cncmastersnapshot [EXTRACTED 1.00]
- **Siemens OPC UA connection requirements (GOARCH=386, credentials, trusted certs)** — docs_cnc_communication_prereq_goarch_386, docs_cnc_communication_prereq_siemens_credentials, docs_cnc_communication_siemens_certificates, docs_cnc_communication_siemens_trust, docs_cnc_communication_client_cert_der, docs_cnc_communication_client_key_der, docs_cnc_communication_siemens_sinumerik_profile [EXTRACTED 1.00]
- **Struct ABI compatibility (pack(8) x86, size 34440, GOARCH=386)** — docs_cnc_communication_cncmastersnapshot, docs_cnc_communication_struct_size_assertion, docs_cnc_communication_prereq_goarch_386, docs_cnc_communication_ts_memory_panic [INFERRED 0.85]

## Communities (70 total, 13 thin omitted)

### Community 0 - "HTTP Handlers Layer"
Cohesion: 0.05
Nodes (42): validator, AuthHandler, BranchHandler, BrandHandler, BreakHandler, CategoryHandler, ClienteleHandler, ControllerHandler (+34 more)

### Community 1 - "Integration & DTO Tests"
Cohesion: 0.07
Nodes (89): CreateUserRequest, TestCreateUserRequestValidate(), TestCreateUserRequestValidateDropsDuplicateMachines(), TestCreateUserRequestValidateNormalizes(), TestCustomerFiltersValidate(), TestParseDayStart(), TestUpdateThemeRequestValidate(), TestUpdateUserRequestValidateStatus() (+81 more)

### Community 2 - "User Supervisor Module"
Cohesion: 0.06
Nodes (22): CreateUserCellRequest, UpdateUserCellRequest, validateUserCell(), CreateUserSupervisorRequest, UpdateUserSupervisorRequest, validateUserSupervisor(), NewUserCellHandler(), NewUserSupervisorHandler() (+14 more)

### Community 3 - "Org Location/Shop DTOs"
Cohesion: 0.07
Nodes (18): CreateUserLocationRequest, UpdateUserLocationRequest, validateUserLocation(), CreateUserShopRequest, UpdateUserShopRequest, validateUserShop(), NewUserShopHandler(), UserLocation (+10 more)

### Community 4 - "Route Registration & Guards"
Cohesion: 0.07
Nodes (27): registerBranchRoutes(), registerBrandRoutes(), registerBreakRoutes(), registerCategoryRoutes(), registerClienteleRoutes(), registerControllerRoutes(), registerHolidayRoutes(), registerLicenceKeyRoutes() (+19 more)

### Community 5 - "Machine DTOs"
Cohesion: 0.08
Nodes (25): CreateMachineRequest, MachineFields, UpdateMachineRequest, validCoordinate(), validDate(), validHealthIDs(), validIPv4(), CreateProductionScheduleRequest (+17 more)

### Community 6 - "DTO Validation & Transactions"
Cohesion: 0.13
Nodes (21): WithTransaction(), clockMinutes(), CreateShiftRequest, ShiftFields, UpdateShiftRequest, MinutesBetween(), Shift, ShiftBreak (+13 more)

### Community 7 - "CNC Snapshot Struct Fields"
Cohesion: 0.11
Nodes (38): CncMasterSnapshot (34,440 bytes, #pragma pack(8), x86), AbsPos ([5]float64, Inches / Degrees, Work Coordinate System (WCS) positions), ActiveGCode ([16]byte, ASCII, Currently active G-code group), ActiveProgPath ([256]byte, ASCII Path, Storage path of active NC artifact), AlarmMsg ([128]byte, UTF-8 / ASCII, Alarm description message), AlarmNo (int32, Code (0 = OK), Primary alarm or fault number), AxisLoadXPercentage (int16, Percent (%), Servo X axis drive load (AxisLoadX/Y/ZPercentage)), AxisLoadYPercentage (int16, Percent (%), Servo Y axis drive load (AxisLoadX/Y/ZPercentage)) (+30 more)

### Community 8 - "Licence Key DTOs"
Cohesion: 0.11
Nodes (17): LicenceKeyDuration, NewLicenceKey, CreateLicenceKeysRequest, GenerateLicenceKeysRequest, UpdateLicenceKeyRequest, licenceKeyDurations(), validDuration(), ValidLicenceKey() (+9 more)

### Community 9 - "Clientele Module"
Cohesion: 0.12
Nodes (11): ClienteleFields, CreateClienteleRequest, UpdateClienteleRequest, validClienteleEmail(), Clientele, ClienteleRepository, mapClienteleError(), clienteleFromFields() (+3 more)

### Community 10 - "Production Schedules (2)"
Cohesion: 0.14
Nodes (14): ProductionSchedule, ProductionScheduleItem, ProductionScheduleRepository, mapProductionScheduleError(), NewProductionScheduleRepository(), checkJobsOverlap(), ProductionScheduleService, NewProductionScheduleService() (+6 more)

### Community 11 - "Part Process DTOs"
Cohesion: 0.13
Nodes (13): Part, PartProcess, PartSubProgram, PartDrawingStorage, NewPartDrawingStorage(), NewPartRepository(), PartService, NewPartService() (+5 more)

### Community 12 - "Authentication Flow"
Cohesion: 0.10
Nodes (12): LoginRequest, AuthRepository, mapAuthError(), AuthService, hashToken(), NewAuthService(), newSessionToken(), truncate() (+4 more)

### Community 13 - "CNC DLL Artifacts & Fanuc"
Cohesion: 0.13
Nodes (29): CncCommunication.dll (C++ core runtime, x86, OPC UA inside), CNCdisconnect(machineId int32), CncCommunication.dll API & Integration Guide (Go Client) - text version, Fanuc profile (engine 1, port 8193, Fwlib32.dll + fwlibe1.dll), Fwlib32.dll (Fanuc FOCAS2), fwlibe1.dll (Fanuc FOCAS2), LazyDLL binding (syscall.NewLazyDLL + eager mod.Load in init), libcrypto-3.dll (OpenSSL) (+21 more)

### Community 14 - "Controller Module & db.Execute"
Cohesion: 0.13
Nodes (8): CreateControllerRequest, UpdateControllerRequest, validateControllerFields(), Controller, ControllerRepository, mapControllerError(), ControllerService, toControllerResponseError()

### Community 15 - "Server Wiring (customer entities)"
Cohesion: 0.12
Nodes (18): NewBreakHandler(), NewClienteleHandler(), NewLicenceKeyHandler(), NewMachineHandler(), NewProductionScheduleHandler(), CrossOrigin(), Chain(), Middleware (+10 more)

### Community 16 - "Holiday DTOs & Model"
Cohesion: 0.17
Nodes (10): AuthUser, Holiday, HolidayRepository, mapHolidayError(), companyOwner(), ownCompanyCustomer(), HolidayService, holidayFromFields() (+2 more)

### Community 17 - "Branch Module"
Cohesion: 0.14
Nodes (9): CreateBranchRequest, UpdateBranchRequest, validateBranchFields(), Branch, BranchRepository, mapBranchError(), cityOwner(), BranchService (+1 more)

### Community 18 - "Catalog & Brand DTOs"
Cohesion: 0.13
Nodes (10): validateCatalogFields(), CreateCategoryRequest, UpdateCategoryRequest, Category, CategoryRepository, mapCategoryError(), NewCategoryRepository(), CategoryService (+2 more)

### Community 19 - "Product Module"
Cohesion: 0.14
Nodes (11): CreateProductRequest, UpdateProductRequest, validateProductFields(), Product, ProductRepository, mapProductError(), NewProductRepository(), ProductService (+3 more)

### Community 20 - "User Service"
Cohesion: 0.21
Nodes (15): CustomerDetails, UpdateUserRequest, UserResponse, canSee(), companyProfileFrom(), UserService, NewUserService(), newUserServices() (+7 more)

### Community 21 - "Server Wiring (master data)"
Cohesion: 0.10
Nodes (16): main(), run(), NewBrandHandler(), NewCategoryHandler(), NewControllerHandler(), NewHolidayHandler(), NewMachineModelHandler(), NewPartHandler() (+8 more)

### Community 22 - "Catalog & Brand DTOs (2)"
Cohesion: 0.13
Nodes (10): CreateBrandRequest, UpdateBrandRequest, Brand, BrandRepository, mapBrandError(), NewBrandRepository(), BrandService, NewBrandService() (+2 more)

### Community 23 - "User City Module"
Cohesion: 0.14
Nodes (8): CreateUserCityRequest, UpdateUserCityRequest, validateUserCity(), UserCity, UserCityRepository, mapUserCityError(), UserCityService, toUserCityResponseError()

### Community 24 - "CNC Go Client Run & Check Mode"
Cohesion: 0.12
Nodes (24): Check mode: go run sinumerik_client.go check (OPC UA node existence report), GcodeGroup ([512]byte, ASCII (reused), DLL reuses this field for the ';'-separated missing-node list), Go Client (sinumerik_client.go, callback -> copy -> render), Prerequisite: 32-bit x86 DLL, compile Go with GOARCH=386, Run: go run sinumerik_client.go (live dashboard), §8 Running the Go Client, §6 Snapshot Data Structure, §10 Diagnostics & Troubleshooting (+16 more)

### Community 25 - "Server Wiring (master data) (2)"
Cohesion: 0.14
Nodes (19): NewBranchHandler(), NewUserCityHandler(), NewUserCountryHandler(), NewUserHandler(), NewUserLocationHandler(), NewUserStateHandler(), NewBranchRepository(), NewUserCityRepository() (+11 more)

### Community 26 - "User & Theme DTOs"
Cohesion: 0.16
Nodes (14): CreateAddressRequest, ThemeResponse, UpdateAddressRequest, UpdateThemeRequest, isTimeZone(), isWebsite(), normalizeFilters(), normalizeTheme() (+6 more)

### Community 27 - "Part Process DTOs (2)"
Cohesion: 0.18
Nodes (3): executeAffected(), PartRepository, mapPartError()

### Community 28 - "Profiles & Shared Models"
Cohesion: 0.15
Nodes (10): UserProfile, UserProfileRepository, mapProfileError(), NewUserProfileRepository(), NewUserProfileService(), toUserProfileResponseError(), UserProfileService, sameClock() (+2 more)

### Community 29 - "Machine Model Module"
Cohesion: 0.16
Nodes (6): MachineModelRequest, MachineModel, MachineModelRepository, mapMachineModelError(), MachineModelService, toMachineModelResponseError()

### Community 30 - "DTO Validation & Transactions (2)"
Cohesion: 0.20
Nodes (6): Break, BreakRepository, mapBreakError(), breakFromFields(), BreakService, toBreakResponseError()

### Community 31 - "Part Process DTOs (3)"
Cohesion: 0.15
Nodes (11): PartProcessFields, PartSubProgramFields, CreatePartRequest, PartFields, UpdatePartRequest, UploadPartDrawingRequest, uniquePositiveIDs(), ValidPartDrawingFile() (+3 more)

### Community 32 - "Address Module"
Cohesion: 0.17
Nodes (9): Address, AddressRepository, mapAddressError(), NewAddressRepository(), NewAddressService(), toAddressResponseError(), AddressService, TestAddressCRUD() (+1 more)

### Community 33 - "User Line Module"
Cohesion: 0.19
Nodes (6): UserLine, UserLineRepository, mapUserLineError(), UserLineService, NewUserLineService(), toUserLineResponseError()

### Community 34 - "CNC Machine Types & Catalog"
Cohesion: 0.13
Nodes (20): Machine Catalog Brand -> Controller -> Model, IsConnected (int16, Boolean (0/1), Network connection state), Engine 0: Simulator (local loopback mock engine), Engine 2: Siemens (Sinumerik OPC UA, active profile), Engine 3: Syntec (Open API / TCP), Engine 4: Mitsubishi (driver not specified), Machine identifier mapping (engine id enum 0-4), Prerequisite: SIEMENS_USER / SIEMENS_PASS env vars before CNCConnect (+12 more)

### Community 35 - "User Country Module"
Cohesion: 0.15
Nodes (6): CreateUserCountryRequest, UserCountry, UserCountryRepository, mapUserCountryError(), UserCountryService, toUserCountryResponseError()

### Community 36 - "Profiles & Shared Models (2)"
Cohesion: 0.18
Nodes (8): CompanyProfile, CompanyProfileRepository, mapCompanyProfileError(), NewCompanyProfileRepository(), NewCompanyProfileService(), toCompanyProfileResponseError(), CompanyProfileService, TestCompanyProfileCRUD()

### Community 37 - "Roles Module"
Cohesion: 0.20
Nodes (9): RoleHandler, NewRoleHandler(), Role, RoleRepository, mapRoleError(), NewRoleRepository(), RoleService, NewRoleService() (+1 more)

### Community 38 - "CNC Connect API & Callbacks"
Cohesion: 0.21
Nodes (14): Package-level callback singletons (syscall.NewCallback), CNCConnect(ip, port, dataCb, errCb, machineType) int32, CNCConnect param dataCb (uintptr stdcall OnData(const CncMasterSnapshot*)), CNCConnect param errCb (uintptr stdcall OnErr(int32 code, const char* msg)), CNCConnect param ip (const char*, ANSI NUL-terminated), CNCConnect param machineType (int32 engine id, see §4.2), CNCConnect param port (int32, e.g. 4840), Mutex-protected copy: dataMutex sync.RWMutex -> lastRecord/hasRecord (+6 more)

### Community 39 - "Server API Domain Docs"
Cohesion: 0.21
Nodes (14): Clienteles API (/clienteles), Company Scoping (companyOwner / ownCompanyCustomer), Customer Display Names (customerName / companyName), Licence Keys API (/licence-keys), Machines API (/machines), Organization Hierarchy APIs (/user-country..supervisor), Parts API (/parts with processes, drawings), Production Schedules API (/production-schedules) (+6 more)

### Community 40 - "Config, DB & Seeding"
Cohesion: 0.23
Nodes (8): main(), run(), Config, getEnv(), DBConfig, Load(), Connect(), TestMain()

### Community 41 - "User Country & Core Models"
Cohesion: 0.23
Nodes (4): UserListFilter, User, UserRepository, mapUserError()

### Community 42 - "Config, DB & Seeding (2)"
Cohesion: 0.44
Nodes (10): checkNotTaken(), createSuperadmin(), ensureProfile(), execute(), main(), parseFlags(), run(), superadminRoleID() (+2 more)

### Community 43 - "User State Module"
Cohesion: 0.26
Nodes (4): CreateUserStateRequest, validateUserState(), UserStateService, toUserStateResponseError()

### Community 44 - "User State Module (2)"
Cohesion: 0.29
Nodes (4): Execute(), UserState, UserStateRepository, mapUserStateError()

### Community 45 - "Querier Repositories"
Cohesion: 0.23
Nodes (9): Querier, NewClienteleRepository(), NewHolidayRepository(), NewLicenceKeyRepository(), NewUserLineRepository(), execute(), UserMachineAccessRepository, NewUserMachineAccessRepository() (+1 more)

### Community 46 - "CRUD Pattern Guidelines"
Cohesion: 0.27
Nodes (6): CRUD Implementation Checklist (11 steps), Golden Reference Pattern: User CRUD, Unique Constraint Error Mapping (pgconn 23505), SendResponse JSON Response Format, User CRUD Golden Reference, Adding a New Entity (README steps)

### Community 47 - "Engineering Principles & API Docs (2)"
Cohesion: 0.22
Nodes (5): auth_tokens SHA-256 Hashed Sessions (migration 014), middleware.CrossOrigin (http.CrossOriginProtection), Login Lockout (5 wrong passwords, 15 min, 429), Per-Request Session and User Re-check, ts_session Cookie (HttpOnly, SameSite=Strict)

### Community 48 - "Setup & Tooling Docs"
Cohesion: 0.22
Nodes (8): Air as go.mod Tool (hot reload), VS Code Debugging with Delve (launch.json configs), .env from .env.example (config defaults), GOARCH=386 Pitfall, Install Go 1.27.1+, PostgreSQL Install and Migration, Setup Troubleshooting Table, GET /health

### Community 51 - "Siemens Certificates & Trust"
Cohesion: 0.48
Nodes (7): client_cert.der/.pem (X.509 client cert for Siemens OPC UA), client_key.der/.pem (client private key), Siemens certs: client_cert.der + client_key.der in binary directory, Siemens trust: cert must be trusted in NCU Access Administration, Troubleshoot: BadRequestNotAllowed -> Export client_cert, BadRequestNotAllowed (certificate trust rejection), Client Certificate & Key (client_cert.der / client_key.der)

### Community 52 - "DTO Validation & Transactions (4)"
Cohesion: 0.48
Nodes (3): BreakFields, CreateBreakRequest, UpdateBreakRequest

### Community 53 - "Holiday DTOs & Model (2)"
Cohesion: 0.48
Nodes (3): CreateHolidayRequest, HolidayFields, UpdateHolidayRequest

### Community 55 - "CNC Event-Driven Data Flow"
Cohesion: 0.60
Nodes (5): Flow 1: load DLL, register two stdcall callbacks (syscall.NewCallback), Flow 2: CNCConnect starts DLL background worker thread, Flow 3: worker polls controller, assembles snapshot, invokes data callback, Flow 4: callback copies payload, signals renderer via buffered channel, Flow 5: main goroutine renders live ANSI dashboard

### Community 56 - "User Line Module (2)"
Cohesion: 0.47
Nodes (3): CreateUserLineRequest, UpdateUserLineRequest, validateUserLine()

### Community 57 - "Location Repository"
Cohesion: 0.53
Nodes (3): mapLocationError(), NewLocationRepository(), LocationRepository

### Community 58 - "Project Context & Tech Stack"
Cohesion: 0.40
Nodes (5): Air Hot Reload, 7th Sense Server - Claude Context, Engineering Principles Summary, PostgreSQL + TimescaleDB via pgx/v5, SampleTimestampMs (int64, Milliseconds, Epoch timestamp of snapshot emission)

### Community 59 - "Authentication Flow (2)"
Cohesion: 0.40
Nodes (4): NewAuthHandler(), Authenticate(), NewAuthRepository(), newAuthRouter()

## Ambiguous Edges - Review These
- `Machines API (/machines)` → `Snapshot Fields (positions, spindle, feed, program, tool, alarms)`  [AMBIGUOUS]
  docs/CncCommunication.dll — API & Integration Guide (Go Client).pdf · relation: conceptually_related_to

## Knowledge Gaps
- **26 isolated node(s):** `7th-sense-server`, `SendResponse JSON Response Format`, `Adding a New Entity (README steps)`, `middleware.CrossOrigin (http.CrossOriginProtection)`, `Login Lockout (5 wrong passwords, 15 min, 429)` (+21 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 114 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **13 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `NewError()` connect `User & Theme DTOs` to `HTTP Handlers Layer`, `Integration & DTO Tests`, `User Supervisor Module`, `Org Location/Shop DTOs`, `Machine DTOs`, `DTO Validation & Transactions`, `Licence Key DTOs`, `Clientele Module`, `Production Schedules (2)`, `Part Process DTOs`, `Authentication Flow`, `Controller Module & db.Execute`, `Holiday DTOs & Model`, `Branch Module`, `Catalog & Brand DTOs`, `Product Module`, `User Service`, `Catalog & Brand DTOs (2)`, `User City Module`, `Profiles & Shared Models`, `Machine Model Module`, `DTO Validation & Transactions (2)`, `Part Process DTOs (3)`, `User Line Module`, `User Country Module`, `Profiles & Shared Models (2)`, `Roles Module`, `User State Module`, `DTO Validation & Transactions (4)`, `Holiday DTOs & Model (2)`, `User Line Module (2)`?**
  _High betweenness centrality (0.140) - this node is a cross-community bridge._
- **What connects `7th-sense-server`, `SendResponse JSON Response Format`, `Adding a New Entity (README steps)` to the rest of the system?**
  _26 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `HTTP Handlers Layer` be split into smaller, more focused modules?**
  _Cohesion score 0.05215600292339379 - nodes in this community are weakly interconnected._
- **Why does `run()` connect `Server Wiring (master data)` to `HTTP Handlers Layer`, `User Supervisor Module`, `Org Location/Shop DTOs`, `Machine DTOs`, `DTO Validation & Transactions`, `Licence Key DTOs`, `Clientele Module`, `Production Schedules (2)`, `Part Process DTOs`, `Authentication Flow`, `Server Wiring (customer entities)`, `Holiday DTOs & Model`, `Catalog & Brand DTOs`, `Product Module`, `User Service`, `Catalog & Brand DTOs (2)`, `Server Wiring (master data) (2)`, `User Line Module`, `Roles Module`, `Config, DB & Seeding`, `Querier Repositories`, `Authentication Flow (2)`?**
  _High betweenness centrality (0.061) - this node is a cross-community bridge._
- **Should `Integration & DTO Tests` be split into smaller, more focused modules?**
  _Cohesion score 0.06542056074766354 - nodes in this community are weakly interconnected._
- **Why does `WithAuthUser()` connect `Integration & DTO Tests` to `HTTP Handlers Layer`, `Part Process DTOs (2)`, `Authentication Flow (2)`, `Holiday DTOs & Model`?**
  _High betweenness centrality (0.031) - this node is a cross-community bridge._

### Low-confidence Hints
_AMBIGUOUS edges — the extractor was unsure. Verify before acting on these._

- **What is the exact relationship between `Machines API (/machines)` and `Snapshot Fields (positions, spindle, feed, program, tool, alarms)`?**
  _Edge tagged AMBIGUOUS (relation: conceptually_related_to) - confidence is low._