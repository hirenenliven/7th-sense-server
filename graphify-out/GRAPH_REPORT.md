# Graph Report - 7th-sense-server  (2026-10-20)

## Corpus Check
- 1637 files · ~0 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 1637 nodes · 13211 edges · 69 communities (59 shown, 10 thin omitted)
- Extraction: 95% EXTRACTED · 5% INFERRED · 0% AMBIGUOUS · INFERRED: 620 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- HTTP Handlers Layer
- Integration & DTO Tests
- Production Schedules (2)
- Integration & DTO Tests (2)
- Authentication Flow
- Controller Module & db.Execute
- Server Wiring (customer entities)
- Holiday DTOs & Model
- Branch Module
- Catalog & Brand DTOs
- Product Module
- User Supervisor Module
- User Service
- Server Wiring (master data)
- HTTP Handlers Layer (2)
- User City Module
- Route Registration & Guards (2)
- User & Theme DTOs
- HTTP Handlers Layer (3)
- Profiles & Shared Models
- Machine Model Module
- Org Location/Shop DTOs
- User Supervisor Module (2)
- Part Process DTOs (3)
- Org Location/Shop DTOs (2)
- User Line Module
- HTTP Handlers Layer (4)
- User Country Module
- Integration & DTO Tests (3)
- Roles Module
- Sinumerik DLL Go Client (use case)
- HTTP Handlers Layer (5)
- Config, DB & Seeding (2)
- User State Module
- Integration & DTO Tests (4)
- Querier Repositories
- Integration & DTO Tests (5)
- User & Theme DTOs (2)
- Machine DTOs
- Machine DTOs (2)
- DTO Validation & Transactions (4)
- Part Process DTOs (3) (2)
- Licence Key DTOs (2)
- HTTP Handlers Layer (6)
- Server Wiring (master data) (2)
- HTTP Handlers Layer (7)
- DTO Validation & Transactions
- Licence Key DTOs
- Clientele Module
- Route Registration & Guards
- Integration & DTO Tests (6)
- Sinumerik DLL Go Client (use case) (2)
- Go Module Root
- CNC DLL Artifacts & Fanuc
- CNC Go Client Run & Check Mode
- CNC Connect API & Callbacks
- Server API Domain Docs
- CRUD Pattern Guidelines
- Setup & Tooling Docs
- Siemens Certificates & Trust
- CNC Event-Driven Data Flow
- Transactions & DB Query Rules
- Current User Middleware
- Role Guard Middleware
- SendResponse Helper
- Test Auth Helpers
- CNC Snapshot Struct Fields

## God Nodes (most connected - your core abstractions)
1. `Execute()` - 158 edges
2. `SendResponse()` - 130 edges
3. `NewError()` - 121 edges
4. `run()` - 78 edges
5. `pathID()` - 71 edges
6. `AuthUser` - 59 edges
7. `decodeAndValidate()` - 53 edges
8. `Querier` - 52 edges
9. `currentUser()` - 47 edges
10. `CncMasterSnapshot (34,440 bytes, #pragma pack(8), x86)` - 38 edges

## Surprising Connections (you probably didn't know these)
- `Machine Identifier Mapping (0 Simulator,1 Fanuc,2 Siemens,3 Syntec,4 Mitsubishi)` --semantically_similar_to--> `Machine Catalog Brand -> Controller -> Model`  [INFERRED] [semantically similar]
  docs/CncCommunication.dll — API & Integration Guide (Go Client).pdf → .claude/CLAUDE.md
- `Snapshot Fields (positions, spindle, feed, program, tool, alarms)` --conceptually_related_to--> `Machines API (/machines)`  [AMBIGUOUS]
  docs/CncCommunication.dll — API & Integration Guide (Go Client).pdf → .claude/CLAUDE.md
- `OPC UA (Siemens transport, port 4840)` --semantically_similar_to--> `OPC UA Protocol`  [INFERRED] [semantically similar]
  docs/cnc-communication.md → docs/CncCommunication.dll — API & Integration Guide (Go Client).pdf
- `Prerequisite: SIEMENS_USER / SIEMENS_PASS env vars before CNCConnect` --semantically_similar_to--> `SIEMENS_USER / SIEMENS_PASS Environment Variables`  [INFERRED] [semantically similar]
  docs/cnc-communication.md → docs/CncCommunication.dll — API & Integration Guide (Go Client).pdf
- `Siemens Sinumerik profile (engine 2, ports 4840/4841)` --semantically_similar_to--> `Siemens Sinumerik OPC UA Profile (Engine 2, port 4840/4841)`  [INFERRED] [semantically similar]
  docs/cnc-communication.md → docs/CncCommunication.dll — API & Integration Guide (Go Client).pdf

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Event-driven telemetry pipeline (CNCConnect -> worker -> callback -> mutex copy -> render)** — docs_cnc_communication_cncconnect, docs_cnc_communication_flow_step1_register_callbacks, docs_cnc_communication_flow_step2_worker_thread, docs_cnc_communication_flow_step3_poll_and_invoke, docs_cnc_communication_flow_step4_copy_and_signal, docs_cnc_communication_flow_step5_render_dashboard, docs_cnc_communication_ondatacallback, docs_cnc_communication_mutex_copy, docs_cnc_communication_cncmastersnapshot [EXTRACTED 1.00]
- **Siemens OPC UA connection requirements (GOARCH=386, credentials, trusted certs)** — docs_cnc_communication_prereq_goarch_386, docs_cnc_communication_prereq_siemens_credentials, docs_cnc_communication_siemens_certificates, docs_cnc_communication_siemens_trust, docs_cnc_communication_client_cert_der, docs_cnc_communication_client_key_der, docs_cnc_communication_siemens_sinumerik_profile [EXTRACTED 1.00]
- **CNC telemetry callback data flow** — docs_cnccommunication_dll_api_integration_guide_go_client_cncconnect, docs_cnccommunication_dll_api_integration_guide_go_client_stdcall_callbacks, docs_cnccommunication_dll_api_integration_guide_go_client_ondatacallback, docs_cnccommunication_dll_api_integration_guide_go_client_cncmastersnapshot, docs_cnccommunication_dll_api_integration_guide_go_client_event_driven_data_flow [EXTRACTED 1.00]
- **Company-scoped customer resources** — _claude_claude_parts_api, _claude_claude_clienteles_api, _claude_claude_shifts_breaks_holidays, _claude_claude_licence_keys_api, _claude_claude_production_schedules_api, _claude_claude_company_scope [EXTRACTED 1.00]
- **CRUD layered pattern rules** — _claude_claude_layered_architecture, _claude_claude_dto_validation, _claude_claude_service_error_mapping, _claude_claude_coalesce_nullables, _claude_claude_soft_delete, _claude_claude_user_crud_golden_reference [EXTRACTED 1.00]
- **Struct ABI compatibility (pack(8) x86, size 34440, GOARCH=386)** — docs_cnc_communication_cncmastersnapshot, docs_cnc_communication_struct_size_assertion, docs_cnc_communication_prereq_goarch_386, docs_cnc_communication_ts_memory_panic [INFERRED 0.85]
- **Session authentication and authorization flow** — readme_session_cookie, readme_auth_tokens_hashed, readme_per_request_session_check, readme_login_lockout, readme_csrf_cross_origin, _claude_claude_role_guards [INFERRED 0.85]

## Communities (69 total, 10 thin omitted)

### Community 0 - "HTTP Handlers Layer"
Cohesion: 0.11
Nodes (8): BreakHandler, ClienteleHandler, HolidayHandler, ProductionScheduleHandler, ShiftHandler, currentUser(), listUserID(), queryID()

### Community 1 - "Integration & DTO Tests"
Cohesion: 0.24
Nodes (21): catalogCRUDCase, catalogRecord, containsID(), expectStatus(), newCatalogRouter(), runCatalogCRUD(), sendRequest(), uniqueSuffix() (+13 more)

### Community 10 - "Production Schedules (2)"
Cohesion: 0.08
Nodes (21): CreateProductionScheduleRequest, ProductionScheduleFields, UpdateProductionScheduleRequest, Controller, ProductionSchedule, ProductionScheduleItem, ControllerRepository, PartRepository (+13 more)

### Community 11 - "Integration & DTO Tests (2)"
Cohesion: 0.48
Nodes (3): ThemeResponse, authUserKey, machineModelRecord

### Community 12 - "Authentication Flow"
Cohesion: 0.13
Nodes (11): LoginRequest, AuthService, loginFailures, loginThrottle, Session, NewAuthHandler(), hashToken(), NewAuthService() (+3 more)

### Community 14 - "Controller Module & db.Execute"
Cohesion: 0.26
Nodes (6): CreateControllerRequest, UpdateControllerRequest, ControllerService, validateControllerFields(), toControllerResponseError(), asInvalidSelection()

### Community 15 - "Server Wiring (customer entities)"
Cohesion: 0.15
Nodes (25): Middleware, NewClienteleHandler(), NewHolidayHandler(), NewPartHandler(), NewProductionScheduleHandler(), Authenticate(), WithAuthUser(), CrossOrigin() (+17 more)

### Community 16 - "Holiday DTOs & Model"
Cohesion: 0.15
Nodes (11): CreateHolidayRequest, HolidayFields, UpdateHolidayRequest, Holiday, HolidayRepository, HolidayService, mapHolidayError(), holidayFromFields() (+3 more)

### Community 17 - "Branch Module"
Cohesion: 0.23
Nodes (5): CreateBranchRequest, UpdateBranchRequest, validateBranchFields(), cityOwner(), toBranchResponseError()

### Community 18 - "Catalog & Brand DTOs"
Cohesion: 0.09
Nodes (15): CreateBrandRequest, UpdateBrandRequest, CreateCategoryRequest, UpdateCategoryRequest, Brand, Category, BrandRepository, CategoryRepository (+7 more)

### Community 19 - "Product Module"
Cohesion: 0.16
Nodes (9): CreateProductRequest, UpdateProductRequest, Product, ProductRepository, ProductService, validateProductFields(), mapProductError(), productImageOrDefault() (+1 more)

### Community 2 - "User Supervisor Module"
Cohesion: 0.14
Nodes (10): CreateUserCellRequest, UpdateUserCellRequest, UserCell, UserCellRepository, UserCellService, validateUserCell(), NewUserCellHandler(), mapUserCellError() (+2 more)

### Community 20 - "User Service"
Cohesion: 0.06
Nodes (39): CustomerDetails, UpdateUserRequest, UserListFilter, UserResponse, Address, CompanyProfile, AddressRepository, CompanyProfileRepository (+31 more)

### Community 21 - "Server Wiring (master data)"
Cohesion: 0.20
Nodes (23): MachineRepository, run(), NewBrandHandler(), NewCategoryHandler(), NewControllerHandler(), NewMachineModelHandler(), NewProductHandler(), NewBranchRepository() (+15 more)

### Community 22 - "HTTP Handlers Layer (2)"
Cohesion: 0.10
Nodes (10): validator, UserCellHandler, UserCityHandler, UserLineHandler, UserLocationHandler, UserShopHandler, UserStateHandler, UserSupervisorHandler (+2 more)

### Community 23 - "User City Module"
Cohesion: 0.16
Nodes (8): CreateUserCityRequest, UpdateUserCityRequest, UserCity, UserCityRepository, UserCityService, validateUserCity(), mapUserCityError(), toUserCityResponseError()

### Community 25 - "Route Registration & Guards (2)"
Cohesion: 0.35
Nodes (3): LicenceKeyDuration, body, licenceKeyDurations()

### Community 26 - "User & Theme DTOs"
Cohesion: 0.16
Nodes (14): UpdateThemeRequest, isTimeZone(), isWebsite(), normalizeFilters(), normalizeTheme(), TestCreateUserRequestValidate(), TestCreateUserRequestValidateDropsDuplicateMachines(), TestCreateUserRequestValidateNormalizes() (+6 more)

### Community 27 - "HTTP Handlers Layer (3)"
Cohesion: 0.12
Nodes (7): BranchHandler, ControllerHandler, MachineHandler, MachineModelHandler, PartHandler, Handlers, pathID()

### Community 28 - "Profiles & Shared Models"
Cohesion: 0.07
Nodes (22): LicenceKey, PartProcess, PartSubProgram, UserProfile, LicenceKeyRepository, UserProfileRepository, LicenceKeyService, UserMachine (+14 more)

### Community 29 - "Machine Model Module"
Cohesion: 0.20
Nodes (6): MachineModelRequest, MachineModel, MachineModelRepository, MachineModelService, mapMachineModelError(), toMachineModelResponseError()

### Community 3 - "Org Location/Shop DTOs"
Cohesion: 0.14
Nodes (10): CreateUserShopRequest, UpdateUserShopRequest, UserShop, UserShopRepository, UserShopService, validateUserShop(), NewUserShopHandler(), mapUserShopError() (+2 more)

### Community 30 - "User Supervisor Module (2)"
Cohesion: 0.15
Nodes (10): CreateUserSupervisorRequest, UpdateUserSupervisorRequest, UserSupervisor, UserSupervisorRepository, UserSupervisorService, validateUserSupervisor(), NewUserSupervisorHandler(), mapUserSupervisorError() (+2 more)

### Community 31 - "Part Process DTOs (3)"
Cohesion: 0.22
Nodes (4): PartProcessFields, PartSubProgramFields, UploadPartDrawingRequest, validText()

### Community 32 - "Org Location/Shop DTOs (2)"
Cohesion: 0.16
Nodes (8): CreateUserLocationRequest, UpdateUserLocationRequest, UserLocation, UserLocationRepository, UserLocationService, validateUserLocation(), mapUserLocationError(), toUserLocationResponseError()

### Community 33 - "User Line Module"
Cohesion: 0.15
Nodes (9): CreateUserLineRequest, UpdateUserLineRequest, UserLine, UserLineRepository, UserLineService, validateUserLine(), NewUserLineHandler(), mapUserLineError() (+1 more)

### Community 34 - "HTTP Handlers Layer (4)"
Cohesion: 0.16
Nodes (4): CategoryHandler, UserHandler, statusRecorder, currentUserID()

### Community 35 - "User Country Module"
Cohesion: 0.18
Nodes (9): CreateUserCountryRequest, UserCountry, UserCountryRepository, UserCountryService, NewUserCountryHandler(), mapUserCountryError(), NewUserCountryRepository(), NewUserCountryService() (+1 more)

### Community 36 - "Integration & DTO Tests (3)"
Cohesion: 0.12
Nodes (19): orgRecord, createAuthTestUser(), TestCreateBrandRequestValidate(), TestCreateCategoryRequestValidate(), TestCreateControllerRequestValidate(), TestCreateProductRequestValidate(), TestUpdateBrandRequestValidate(), TestUpdateCategoryRequestValidate() (+11 more)

### Community 37 - "Roles Module"
Cohesion: 0.16
Nodes (14): RoleHandler, Role, RoleRepository, RoleService, NewRoleHandler(), NewUserHandler(), mapRoleError(), NewRoleRepository() (+6 more)

### Community 40 - "Sinumerik DLL Go Client (use case)"
Cohesion: 0.06
Nodes (27): Config, checkItem, CncMasterSnapshot, DBConfig, main(), run(), main(), axisColor() (+19 more)

### Community 41 - "HTTP Handlers Layer (5)"
Cohesion: 0.17
Nodes (4): BrandHandler, LicenceKeyHandler, ProductHandler, SendResponse()

### Community 42 - "Config, DB & Seeding (2)"
Cohesion: 0.06
Nodes (25): Branch, User, AuthRepository, BranchRepository, ClienteleRepository, UserRepository, LocationRepository, options (+17 more)

### Community 43 - "User State Module"
Cohesion: 0.17
Nodes (10): CreateUserStateRequest, UserState, UserStateRepository, UserStateService, validateUserState(), NewUserStateHandler(), mapUserStateError(), NewUserStateRepository() (+2 more)

### Community 44 - "Integration & DTO Tests (4)"
Cohesion: 0.22
Nodes (18): PartDrawing, machineFixture, machineRecord, testCompany, clienteleBody(), TestClienteleCRUD(), actAs(), containsMachineID() (+10 more)

### Community 45 - "Querier Repositories"
Cohesion: 0.20
Nodes (10): Querier, NewAuthRepository(), NewClienteleRepository(), NewHolidayRepository(), NewLicenceKeyRepository(), NewUserCellRepository(), NewUserLineRepository(), NewUserLocationRepository() (+2 more)

### Community 47 - "Integration & DTO Tests (5)"
Cohesion: 0.27
Nodes (11): CreateUserRequest, apiResponse, userRecord, decodeData(), createUser(), listIDs(), TestCompanyUserManagement(), TestNobodyManagesTheSuperadmin() (+3 more)

### Community 49 - "User & Theme DTOs (2)"
Cohesion: 0.22
Nodes (7): CreateAddressRequest, UpdateAddressRequest, UpdateLicenceKeyRequest, Response, validDuration(), NewError(), newResponse()

### Community 5 - "Machine DTOs"
Cohesion: 0.13
Nodes (15): CreateMachineRequest, MachineFields, UpdateMachineRequest, Machine, MachineScope, MachineService, validCoordinate(), validHealthIDs() (+7 more)

### Community 50 - "Machine DTOs (2)"
Cohesion: 0.24
Nodes (6): ProductionScheduleItemFields, jobSpan, validDate(), ScheduleDay(), checkJobsOverlap(), shiftSpan()

### Community 52 - "DTO Validation & Transactions (4)"
Cohesion: 0.47
Nodes (3): BreakFields, CreateBreakRequest, UpdateBreakRequest

### Community 53 - "Part Process DTOs (3) (2)"
Cohesion: 0.32
Nodes (5): CreateShiftRequest, ShiftFields, UpdateShiftRequest, uniquePositiveIDs(), validWorkingDays()

### Community 54 - "Licence Key DTOs (2)"
Cohesion: 0.29
Nodes (6): NewLicenceKey, CreateLicenceKeysRequest, ValidLicenceKey(), generateLicenceKeys(), licenceKeyItems(), TestLicenceKeyCRUD()

### Community 56 - "HTTP Handlers Layer (6)"
Cohesion: 0.38
Nodes (3): AuthHandler, CurrentUser(), UserID()

### Community 57 - "Server Wiring (master data) (2)"
Cohesion: 0.38
Nodes (7): BranchService, NewBranchHandler(), NewUserCityHandler(), NewUserLocationHandler(), NewBranchService(), NewUserLocationService(), newOrganizationRouter()

### Community 6 - "DTO Validation & Transactions"
Cohesion: 0.05
Nodes (46): CreatePartRequest, PartFields, UpdatePartRequest, AuthUser, Break, Clientele, Part, Shift (+38 more)

### Community 8 - "Licence Key DTOs"
Cohesion: 0.40
Nodes (3): GenerateLicenceKeysRequest, GeneratedLicenceKeys, newLicenceKey()

### Community 9 - "Clientele Module"
Cohesion: 0.38
Nodes (4): ClienteleFields, CreateClienteleRequest, UpdateClienteleRequest, validClienteleEmail()

### Community 4 - "Route Registration & Guards"
Cohesion: 0.18
Nodes (28): RequireRole(), registerBranchRoutes(), registerBrandRoutes(), registerBreakRoutes(), registerCategoryRoutes(), registerClienteleRoutes(), registerControllerRoutes(), registerHolidayRoutes() (+20 more)

### Community 58 - "Integration & DTO Tests (6)"
Cohesion: 0.67
Nodes (4): authRequest(), login(), sessionCookieOf(), TestAuthentication()

### Community 13 - "CNC DLL Artifacts & Fanuc"
Cohesion: 0.09
Nodes (42): Machine identifier mapping (engine id enum 0-4), §1 Overview, CNCConnect param ip (const char*, ANSI NUL-terminated), CNCConnect param port (int32, e.g. 4840), Fanuc profile (engine 1, port 8193, Fwlib32.dll + fwlibe1.dll), Syntec profile (engine 3, port 5001 or controller-assigned), CNCConnect param machineType (int32 engine id, see §4.2), Engine 0: Simulator (local loopback mock engine) (+34 more)

### Community 24 - "CNC Go Client Run & Check Mode"
Cohesion: 0.16
Nodes (18): Prerequisite: SIEMENS_USER / SIEMENS_PASS env vars before CNCConnect, Siemens credentials: SIEMENS_USER / SIEMENS_PASS env vars, §10 Diagnostics & Troubleshooting, Troubleshoot: %1 is not a valid Win32 application -> Compile the Go client with `set GOARCH=386`; validate GOARCH against the DLL architecture, Troubleshoot: 0xC0000135 / DLL load failed: module not found -> All 6 DLLs must be in Lib\ next to the client, Troubleshoot: BadNodeIdUnknown -> Node unavailable for this hardware model; field defaults to 0, Troubleshoot: CNCConnect returns -1 -> Run `ping <ip>` and `Test-NetConnection -ComputerName <ip> -Port <port>` to check the network path, Troubleshoot: Connects but CONNECTED NO -> Verify SIEMENS_USER / SIEMENS_PASS and endpoint address; confirm the OPC UA server is enabled on the NCU (+10 more)

### Community 38 - "CNC Connect API & Callbacks"
Cohesion: 0.36
Nodes (9): §7 Asynchronous Callback Architecture, onErrorCallback(code, messagePtr uintptr) uintptr, onDataCallback(ptr uintptr) uintptr, Package-level callback singletons (syscall.NewCallback), Mutex-protected copy: dataMutex sync.RWMutex -> lastRecord/hasRecord, CNCConnect param errCb (uintptr stdcall OnErr(int32 code, const char* msg)), onDataCallback, onErrorCallback (+1 more)

### Community 39 - "Server API Domain Docs"
Cohesion: 0.23
Nodes (13): Clienteles API (/clienteles), Company Scoping (companyOwner / ownCompanyCustomer), Customer Display Names (customerName / companyName), Licence Keys API (/licence-keys), Machines API (/machines), Organization Hierarchy APIs (/user-country..supervisor), Parts API (/parts with processes, drawings), Production Schedules API (/production-schedules) (+5 more)

### Community 46 - "CRUD Pattern Guidelines"
Cohesion: 0.06
Nodes (23): Common Mistakes to Avoid, CRUD Implementation Checklist (11 steps), Golden Reference Pattern: User CRUD, Unique Constraint Error Mapping (pgconn 23505), updated_at Trigger in Migrations, Air Hot Reload, Engineering Principles Summary, PostgreSQL + TimescaleDB via pgx/v5 (+15 more)

### Community 48 - "Setup & Tooling Docs"
Cohesion: 0.22
Nodes (8): Air as go.mod Tool (hot reload), .env from .env.example (config defaults), GOARCH=386 Pitfall, GET /health, VS Code Debugging with Delve (launch.json configs), Install Go 1.27.1+, PostgreSQL Install and Migration, Setup Troubleshooting Table

### Community 51 - "Siemens Certificates & Trust"
Cohesion: 0.31
Nodes (10): Siemens default ports: 4840 (native NCU) / 4841 (custom), Siemens Sinumerik profile (engine 2, ports 4840/4841), Engine 2: Siemens (Sinumerik OPC UA, active profile), Siemens certs: client_cert.der + client_key.der in binary directory, client_cert.der/.pem (X.509 client cert for Siemens OPC UA), client_key.der/.pem (client private key), Troubleshoot: BadRequestNotAllowed -> Export client_cert, Siemens trust: cert must be trusted in NCU Access Administration (+2 more)

### Community 55 - "CNC Event-Driven Data Flow"
Cohesion: 0.19
Nodes (13): Flow 2: CNCConnect starts DLL background worker thread, Check mode: go run sinumerik_client.go check (OPC UA node existence report), Go Client (sinumerik_client.go, callback -> copy -> render), OPC UA (Siemens transport, port 4840), §2 System Architecture, SINUMERIK CNC (PLC/NCU), Troubleshoot: Garbled console output -> Close and reopen the terminal window, Flow 1: load DLL, register two stdcall callbacks (syscall.NewCallback) (+5 more)

### Community 7 - "CNC Snapshot Struct Fields"
Cohesion: 0.10
Nodes (44): §6 Snapshot Data Structure, AbsPos ([5]float64, Inches / Degrees, Work Coordinate System (WCS) positions), ActiveGCode ([16]byte, ASCII, Currently active G-code group), ActiveProgPath ([256]byte, ASCII Path, Storage path of active NC artifact), AlarmMsg ([128]byte, UTF-8 / ASCII, Alarm description message), AlarmNo (int32, Code (0 = OK), Primary alarm or fault number), AxisLoadXPercentage (int16, Percent (%), Servo X axis drive load (AxisLoadX/Y/ZPercentage)), AxisLoadYPercentage (int16, Percent (%), Servo Y axis drive load (AxisLoadX/Y/ZPercentage)) (+36 more)

## Ambiguous Edges - Review These
- `Machines API (/machines)` → `Snapshot Fields (positions, spindle, feed, program, tool, alarms)`  [AMBIGUOUS]
  docs/CncCommunication.dll — API & Integration Guide (Go Client).pdf · relation: conceptually_related_to

## Knowledge Gaps
- **27 isolated node(s):** `hello`, `7th-sense-server`, `CNCConnect param ip (const char*, ANSI NUL-terminated)`, `CNCConnect param port (int32, e.g. 4840)`, `Engine 0: Simulator (local loopback mock engine)` (+22 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 72 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **10 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `NewError()` connect `User & Theme DTOs (2)` to `HTTP Handlers Layer`, `User Supervisor Module`, `Org Location/Shop DTOs`, `Machine DTOs`, `DTO Validation & Transactions`, `Licence Key DTOs`, `Clientele Module`, `Production Schedules (2)`, `Authentication Flow`, `Controller Module & db.Execute`, `Holiday DTOs & Model`, `Branch Module`, `Catalog & Brand DTOs`, `Product Module`, `User Service`, `HTTP Handlers Layer (2)`, `User City Module`, `Route Registration & Guards (2)`, `User & Theme DTOs`, `HTTP Handlers Layer (3)`, `Profiles & Shared Models`, `Machine Model Module`, `User Supervisor Module (2)`, `Part Process DTOs (3)`, `Org Location/Shop DTOs (2)`, `User Line Module`, `User Country Module`, `Roles Module`, `User State Module`, `Machine DTOs (2)`, `DTO Validation & Transactions (4)`, `Part Process DTOs (3) (2)`, `Licence Key DTOs (2)`?**
  _High betweenness centrality (0.065) - this node is a cross-community bridge._
- **What connects `hello`, `7th-sense-server`, `CNCConnect param ip (const char*, ANSI NUL-terminated)` to the rest of the system?**
  _27 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `HTTP Handlers Layer` be split into smaller, more focused modules?**
  _Cohesion score 0.11081560283687943 - nodes in this community are weakly interconnected._
- **Why does `SendResponse()` connect `HTTP Handlers Layer (5)` to `HTTP Handlers Layer`, `HTTP Handlers Layer (4)`, `Route Registration & Guards`, `HTTP Handlers Layer (7)`, `Server Wiring (customer entities)`, `HTTP Handlers Layer (2)`, `HTTP Handlers Layer (6)`, `Route Registration & Guards (2)`, `HTTP Handlers Layer (3)`?**
  _High betweenness centrality (0.023) - this node is a cross-community bridge._
- **Should `Production Schedules (2)` be split into smaller, more focused modules?**
  _Cohesion score 0.08415300546448087 - nodes in this community are weakly interconnected._
- **Why does `Execute()` connect `Config, DB & Seeding (2)` to `User Supervisor Module`, `Org Location/Shop DTOs`, `Machine DTOs`, `DTO Validation & Transactions`, `Production Schedules (2)`, `Integration & DTO Tests (2)`, `Holiday DTOs & Model`, `Catalog & Brand DTOs`, `Product Module`, `User Service`, `User City Module`, `Profiles & Shared Models`, `Machine Model Module`, `User Supervisor Module (2)`, `Org Location/Shop DTOs (2)`, `User Line Module`, `User Country Module`, `Roles Module`, `Sinumerik DLL Go Client (use case)`, `User State Module`, `Querier Repositories`?**
  _High betweenness centrality (0.020) - this node is a cross-community bridge._

### Low-confidence Hints
_AMBIGUOUS edges — the extractor was unsure. Verify before acting on these._

- **What is the exact relationship between `Machines API (/machines)` and `Snapshot Fields (positions, spindle, feed, program, tool, alarms)`?**
  _Edge tagged AMBIGUOUS (relation: conceptually_related_to) - confidence is low._