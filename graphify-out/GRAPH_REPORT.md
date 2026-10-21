# Graph Report - 7th-sense-server  (2026-10-21)

## Corpus Check
- 224 files · ~105,739 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 4 file(s) not represented in the graph (top: (none) 3, .toml 1)

## Summary
- 1733 nodes · 14504 edges · 71 communities (54 shown, 17 thin omitted)
- Extraction: 95% EXTRACTED · 5% INFERRED · 0% AMBIGUOUS · INFERRED: 673 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `b1d97521`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- HTTP Handlers Layer
- testing.T
- UserStateService
- Org Location/Shop DTOs
- Route Registration & Guards
- Machine DTOs
- AuthUser
- CNC Snapshot Struct Fields
- Licence Key DTOs
- Clientele Module
- context.Context
- Integration & DTO Tests (2)
- Steps
- CNC DLL Artifacts & Fanuc
- validateCatalogFields
- run
- HolidayService
- BranchService
- BrandService
- Product Module
- UserService
- github.com/jackc/pgx/v5/pgxpool.Pool
- HTTP Handlers Layer (2)
- User City Module
- CNC Go Client Run & Check Mode
- Route Registration & Guards (2)
- NewError
- decodeAndValidate
- time.Time
- Execute
- User Supervisor Module (2)
- Part Process DTOs (3)
- Org Location/Shop DTOs (2)
- User Line Module
- Handlers
- sinumerik_client.go
- DowntimeReasonService
- CNC Connect API & Callbacks
- Server API Domain Docs
- config.go
- HTTP Handlers Layer (5)
- logger/logger.go
- IdleReasonService
- .UploadDrawing
- 2. `internal/cnc/` — the DLL wrapper only
- CRUD Pattern Guidelines
- Address
- Setup & Tooling Docs
- Response
- Machine DTOs (2)
- Siemens Certificates & Trust
- DTO Validation & Transactions (4)
- .Validate
- LicenceKeyService
- CNC Event-Driven Data Flow
- Category
- CompanyProfile
- CreateAddressRequest
- Transactions & DB Query Rules
- Sinumerik DLL Go Client (use case) (2)
- newUserServices
- Go Module Root
- Current User Middleware
- Role Guard Middleware
- SendResponse Helper
- Test Auth Helpers
- statusRecorder
- .GetAll

## God Nodes (most connected - your core abstractions)
1. `Execute()` - 164 edges
2. `SendResponse()` - 140 edges
3. `NewError()` - 129 edges
4. `run()` - 84 edges
5. `pathID()` - 77 edges
6. `AuthUser` - 69 edges
7. `currentUser()` - 57 edges
8. `decodeAndValidate()` - 57 edges
9. `Querier` - 56 edges
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

## Communities (71 total, 17 thin omitted)

### Community 0 - "HTTP Handlers Layer"
Cohesion: 0.10
Nodes (11): BreakHandler, ClienteleHandler, DowntimeReasonHandler, HolidayHandler, IdleReasonHandler, PartHandler, ProductionScheduleHandler, currentUser() (+3 more)

### Community 1 - "testing.T"
Cohesion: 0.05
Nodes (103): LoginRequest, NewAuthHandler(), NewMachineHandler(), Authenticate(), WithAuthUser(), CrossOrigin(), Logger(), Chain() (+95 more)

### Community 2 - "UserStateService"
Cohesion: 0.18
Nodes (7): CreateUserStateRequest, validateUserState(), UserState, UserStateRepository, mapUserStateError(), UserStateService, toUserStateResponseError()

### Community 3 - "Org Location/Shop DTOs"
Cohesion: 0.15
Nodes (9): CreateUserShopRequest, UpdateUserShopRequest, validateUserShop(), NewUserShopHandler(), UserShop, UserShopRepository, mapUserShopError(), UserShopService (+1 more)

### Community 4 - "Route Registration & Guards"
Cohesion: 0.17
Nodes (30): RequireRole(), registerBranchRoutes(), registerBrandRoutes(), registerBreakRoutes(), registerCategoryRoutes(), registerClienteleRoutes(), registerControllerRoutes(), registerDowntimeReasonRoutes() (+22 more)

### Community 5 - "Machine DTOs"
Cohesion: 0.12
Nodes (16): CreateMachineRequest, MachineFields, UpdateMachineRequest, validCoordinate(), validHealthIDs(), validIPv4(), SeesAllCompanyMachines(), Machine (+8 more)

### Community 6 - "AuthUser"
Cohesion: 0.07
Nodes (34): CreatePartRequest, PartFields, UpdatePartRequest, ValidPartDrawingFile(), NewClienteleHandler(), NewPartHandler(), AuthUser, Clientele (+26 more)

### Community 7 - "CNC Snapshot Struct Fields"
Cohesion: 0.10
Nodes (44): CncMasterSnapshot (34,440 bytes, #pragma pack(8), x86), AbsPos ([5]float64, Inches / Degrees, Work Coordinate System (WCS) positions), ActiveGCode ([16]byte, ASCII, Currently active G-code group), ActiveProgPath ([256]byte, ASCII Path, Storage path of active NC artifact), AlarmMsg ([128]byte, UTF-8 / ASCII, Alarm description message), AlarmNo (int32, Code (0 = OK), Primary alarm or fault number), AxisLoadXPercentage (int16, Percent (%), Servo X axis drive load (AxisLoadX/Y/ZPercentage)), AxisLoadYPercentage (int16, Percent (%), Servo Y axis drive load (AxisLoadX/Y/ZPercentage)) (+36 more)

### Community 8 - "Licence Key DTOs"
Cohesion: 0.40
Nodes (3): GenerateLicenceKeysRequest, GeneratedLicenceKeys, newLicenceKey()

### Community 9 - "Clientele Module"
Cohesion: 0.38
Nodes (4): ClienteleFields, CreateClienteleRequest, UpdateClienteleRequest, validClienteleEmail()

### Community 10 - "context.Context"
Cohesion: 0.07
Nodes (20): NewBreakHandler(), Break, Controller, ShiftBreak, BreakRepository, mapBreakError(), ClienteleRepository, mapClienteleError() (+12 more)

### Community 11 - "Integration & DTO Tests (2)"
Cohesion: 0.49
Nodes (3): ThemeResponse, authUserKey, machineModelRecord

### Community 12 - "Steps"
Cohesion: 0.14
Nodes (13): 0. Check the DLL behaviour on the test machine (no code), 1. Migration `027_add_machine_telemetry.sql`, 3. `internal/collector/` — managing the machines, 4. Data layer (user CRUD pattern), 5. `cmd/collector/main.go`, 6. API (after data is flowing), CNC DLL Integration Plan, DLL limits that shape the design (+5 more)

### Community 13 - "CNC DLL Artifacts & Fanuc"
Cohesion: 0.09
Nodes (42): Machine Catalog Brand -> Controller -> Model, CncCommunication.dll (C++ core runtime, x86, OPC UA inside), CNCConnect(ip, port, dataCb, errCb, machineType) int32, CNCConnect param dataCb (uintptr stdcall OnData(const CncMasterSnapshot*)), CNCConnect param ip (const char*, ANSI NUL-terminated), CNCConnect param machineType (int32 engine id, see §4.2), CNCConnect param port (int32, e.g. 4840), CNCdisconnect(machineId int32) (+34 more)

### Community 14 - "validateCatalogFields"
Cohesion: 0.11
Nodes (11): CreateBrandRequest, UpdateBrandRequest, validateCatalogFields(), CreateCategoryRequest, UpdateCategoryRequest, CreateControllerRequest, UpdateControllerRequest, validateControllerFields() (+3 more)

### Community 15 - "run"
Cohesion: 0.06
Nodes (71): run(), Querier, CreateUserCountryRequest, NewBranchHandler(), NewBrandHandler(), NewCategoryHandler(), NewControllerHandler(), NewMachineModelHandler() (+63 more)

### Community 16 - "HolidayService"
Cohesion: 0.15
Nodes (12): CreateHolidayRequest, HolidayFields, UpdateHolidayRequest, NewHolidayHandler(), Holiday, HolidayRepository, mapHolidayError(), HolidayService (+4 more)

### Community 17 - "BranchService"
Cohesion: 0.15
Nodes (9): CreateBranchRequest, UpdateBranchRequest, validateBranchFields(), Branch, BranchRepository, mapBranchError(), cityOwner(), BranchService (+1 more)

### Community 18 - "BrandService"
Cohesion: 0.25
Nodes (5): Brand, BrandRepository, mapBrandError(), BrandService, toBrandResponseError()

### Community 19 - "Product Module"
Cohesion: 0.23
Nodes (6): Product, ProductRepository, mapProductError(), ProductService, productImageOrDefault(), toProductResponseError()

### Community 20 - "UserService"
Cohesion: 0.14
Nodes (19): WithTransaction(), CreateUserRequest, CustomerDetails, UpdateUserRequest, UserResponse, ParseDayStart(), TestParseDayStart(), AssignableRoles() (+11 more)

### Community 21 - "github.com/jackc/pgx/v5/pgxpool.Pool"
Cohesion: 0.14
Nodes (20): clockMinutes(), CreateShiftRequest, ShiftFields, UpdateShiftRequest, MinutesBetween(), NewShiftHandler(), Shift, NewBreakRepository() (+12 more)

### Community 22 - "HTTP Handlers Layer (2)"
Cohesion: 0.09
Nodes (11): BranchHandler, pathID(), targetUserID(), UserCellHandler, UserCityHandler, UserCountryHandler, UserLineHandler, UserLocationHandler (+3 more)

### Community 23 - "User City Module"
Cohesion: 0.16
Nodes (8): CreateUserCityRequest, UpdateUserCityRequest, validateUserCity(), UserCity, UserCityRepository, mapUserCityError(), UserCityService, toUserCityResponseError()

### Community 24 - "CNC Go Client Run & Check Mode"
Cohesion: 0.16
Nodes (18): Prerequisite: 32-bit x86 DLL, compile Go with GOARCH=386, Prerequisite: SIEMENS_USER / SIEMENS_PASS env vars before CNCConnect, §3.1 Environment & Prerequisites, §10 Diagnostics & Troubleshooting, Siemens credentials: SIEMENS_USER / SIEMENS_PASS env vars, Troubleshoot: BadNodeIdUnknown -> Node unavailable for this hardware model; field defaults to 0, Troubleshoot: CNCConnect returns -1 -> Run `ping <ip>` and `Test-NetConnection -ComputerName <ip> -Port <port>` to check the network path, Troubleshoot: Connects but CONNECTED NO -> Verify SIEMENS_USER / SIEMENS_PASS and endpoint address; confirm the OPC UA server is enabled on the NCU (+10 more)

### Community 25 - "Route Registration & Guards (2)"
Cohesion: 0.35
Nodes (3): LicenceKeyDuration, licenceKeyDurations(), body

### Community 26 - "NewError"
Cohesion: 0.17
Nodes (16): UpdateAddressRequest, UpdateThemeRequest, isTimeZone(), isWebsite(), normalizeFilters(), normalizeTheme(), TestCreateUserRequestValidate(), TestCreateUserRequestValidateDropsDuplicateMachines() (+8 more)

### Community 27 - "decodeAndValidate"
Cohesion: 0.14
Nodes (5): validator, ControllerHandler, LicenceKeyHandler, MachineHandler, decodeAndValidate()

### Community 28 - "time.Time"
Cohesion: 0.14
Nodes (12): UserProfile, UserProfileRepository, mapProfileError(), NewUserProfileRepository(), NewUserProfileService(), toUserProfileResponseError(), UserMachine, UserMachineAccess (+4 more)

### Community 29 - "Execute"
Cohesion: 0.06
Nodes (28): Execute(), MachineModelRequest, MachineModel, User, AuthRepository, mapAuthError(), mapLocationError(), NewLocationRepository() (+20 more)

### Community 30 - "User Supervisor Module (2)"
Cohesion: 0.09
Nodes (17): CreateUserCellRequest, UpdateUserCellRequest, validateUserCell(), CreateUserSupervisorRequest, UpdateUserSupervisorRequest, validateUserSupervisor(), NewUserCellHandler(), NewUserSupervisorHandler() (+9 more)

### Community 31 - "Part Process DTOs (3)"
Cohesion: 0.20
Nodes (5): PartProcessFields, PartSubProgramFields, uniquePositiveIDs(), validText(), validWorkingDays()

### Community 32 - "Org Location/Shop DTOs (2)"
Cohesion: 0.16
Nodes (8): CreateUserLocationRequest, UpdateUserLocationRequest, validateUserLocation(), UserLocation, UserLocationRepository, mapUserLocationError(), UserLocationService, toUserLocationResponseError()

### Community 33 - "User Line Module"
Cohesion: 0.14
Nodes (10): CreateUserLineRequest, UpdateUserLineRequest, validateUserLine(), NewUserLineHandler(), UserLine, UserLineRepository, mapUserLineError(), NewUserCellService() (+2 more)

### Community 34 - "Handlers"
Cohesion: 0.15
Nodes (6): AuthHandler, currentUserID(), UserHandler, CurrentUser(), UserID(), Handlers

### Community 35 - "sinumerik_client.go"
Cohesion: 0.17
Nodes (14): axisColor(), cstr(), enableVT(), isRotary(), lbl(), logf(), main(), nameOr() (+6 more)

### Community 36 - "DowntimeReasonService"
Cohesion: 0.15
Nodes (11): CreateDowntimeReasonRequest, DowntimeReasonFields, UpdateDowntimeReasonRequest, NewDowntimeReasonHandler(), DowntimeReason, DowntimeReasonRepository, mapDowntimeReasonError(), downtimeReasonFromFields() (+3 more)

### Community 38 - "CNC Connect API & Callbacks"
Cohesion: 0.36
Nodes (9): Package-level callback singletons (syscall.NewCallback), CNCConnect param errCb (uintptr stdcall OnErr(int32 code, const char* msg)), Mutex-protected copy: dataMutex sync.RWMutex -> lastRecord/hasRecord, onDataCallback(ptr uintptr) uintptr, onErrorCallback(code, messagePtr uintptr) uintptr, §7 Asynchronous Callback Architecture, onDataCallback, onErrorCallback (+1 more)

### Community 39 - "Server API Domain Docs"
Cohesion: 0.23
Nodes (13): Clienteles API (/clienteles), Company Scoping (companyOwner / ownCompanyCustomer), Customer Display Names (customerName / companyName), Licence Keys API (/licence-keys), Machines API (/machines), Organization Hierarchy APIs (/user-country..supervisor), Parts API (/parts with processes, drawings), Production Schedules API (/production-schedules) (+5 more)

### Community 40 - "config.go"
Cohesion: 0.26
Nodes (8): main(), run(), Config, getEnv(), DBConfig, Load(), Connect(), TestMain()

### Community 41 - "HTTP Handlers Layer (5)"
Cohesion: 0.14
Nodes (5): BrandHandler, CategoryHandler, MachineModelHandler, ProductHandler, SendResponse()

### Community 43 - "IdleReasonService"
Cohesion: 0.15
Nodes (11): CreateIdleReasonRequest, IdleReasonFields, UpdateIdleReasonRequest, NewIdleReasonHandler(), IdleReason, IdleReasonRepository, mapIdleReasonError(), IdleReasonService (+3 more)

### Community 45 - "2. `internal/cnc/` — the DLL wrapper only"
Cohesion: 0.60
Nodes (5): 2. `internal/cnc/` — the DLL wrapper only, findLibDir(), initDLL(), onData(), setDllDirectory()

### Community 46 - "CRUD Pattern Guidelines"
Cohesion: 0.06
Nodes (23): Common Mistakes to Avoid, CRUD Implementation Checklist (11 steps), 7th Sense Server Engineering Principles (Skills Doc), Golden Reference Pattern: User CRUD, Unique Constraint Error Mapping (pgconn 23505), updated_at Trigger in Migrations, Air Hot Reload, 7th Sense Server - Claude Context (+15 more)

### Community 47 - "Address"
Cohesion: 0.27
Nodes (5): Address, AddressRepository, mapAddressError(), toAddressResponseError(), AddressService

### Community 48 - "Setup & Tooling Docs"
Cohesion: 0.22
Nodes (8): Air as go.mod Tool (hot reload), VS Code Debugging with Delve (launch.json configs), .env from .env.example (config defaults), GOARCH=386 Pitfall, Install Go 1.27.1+, PostgreSQL Install and Migration, Setup Troubleshooting Table, GET /health

### Community 50 - "Machine DTOs (2)"
Cohesion: 0.18
Nodes (10): validDate(), CreateProductionScheduleRequest, ProductionScheduleFields, ProductionScheduleItemFields, UpdateProductionScheduleRequest, ScheduleDay(), asInvalidSelection(), checkJobsOverlap() (+2 more)

### Community 51 - "Siemens Certificates & Trust"
Cohesion: 0.31
Nodes (10): client_cert.der/.pem (X.509 client cert for Siemens OPC UA), client_key.der/.pem (client private key), Engine 2: Siemens (Sinumerik OPC UA, active profile), Siemens certs: client_cert.der + client_key.der in binary directory, Siemens default ports: 4840 (native NCU) / 4841 (custom), Siemens Sinumerik profile (engine 2, ports 4840/4841), Siemens trust: cert must be trusted in NCU Access Administration, Troubleshoot: BadRequestNotAllowed -> Export client_cert (+2 more)

### Community 52 - "DTO Validation & Transactions (4)"
Cohesion: 0.47
Nodes (3): BreakFields, CreateBreakRequest, UpdateBreakRequest

### Community 54 - "LicenceKeyService"
Cohesion: 0.16
Nodes (11): NewLicenceKey, CreateLicenceKeysRequest, UpdateLicenceKeyRequest, NewLicenceKeyHandler(), LicenceKey, LicenceKeyRepository, mapLicenceKeyError(), LicenceKeyService (+3 more)

### Community 55 - "CNC Event-Driven Data Flow"
Cohesion: 0.19
Nodes (13): Check mode: go run sinumerik_client.go check (OPC UA node existence report), Flow 1: load DLL, register two stdcall callbacks (syscall.NewCallback), Flow 2: CNCConnect starts DLL background worker thread, Flow 3: worker polls controller, assembles snapshot, invokes data callback, Flow 4: callback copies payload, signals renderer via buffered channel, Flow 5: main goroutine renders live ANSI dashboard, Go Client (sinumerik_client.go, callback -> copy -> render), OPC UA (Siemens transport, port 4840) (+5 more)

### Community 56 - "Category"
Cohesion: 0.25
Nodes (5): Category, CategoryRepository, mapCategoryError(), CategoryService, toCategoryResponseError()

### Community 57 - "CompanyProfile"
Cohesion: 0.28
Nodes (5): CompanyProfile, CompanyProfileRepository, mapCompanyProfileError(), toCompanyProfileResponseError(), CompanyProfileService

### Community 64 - "newUserServices"
Cohesion: 0.23
Nodes (11): NewAddressRepository(), NewCompanyProfileRepository(), UserMachineAccessRepository, NewUserMachineAccessRepository(), NewAddressService(), NewCompanyProfileService(), newUserServices(), userServices (+3 more)

## Ambiguous Edges - Review These
- `Machines API (/machines)` → `Snapshot Fields (positions, spindle, feed, program, tool, alarms)`  [AMBIGUOUS]
  docs/CncCommunication.dll — API & Integration Guide (Go Client).pdf · relation: conceptually_related_to

## Knowledge Gaps
- **27 isolated node(s):** `hello`, `7th-sense-server`, `CNCConnect param ip (const char*, ANSI NUL-terminated)`, `CNCConnect param port (int32, e.g. 4840)`, `Engine 0: Simulator (local loopback mock engine)` (+22 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 83 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **17 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `NewError()` connect `NewError` to `HTTP Handlers Layer`, `testing.T`, `UserStateService`, `Org Location/Shop DTOs`, `Machine DTOs`, `AuthUser`, `Licence Key DTOs`, `Clientele Module`, `context.Context`, `run`, `HolidayService`, `BranchService`, `BrandService`, `Product Module`, `UserService`, `github.com/jackc/pgx/v5/pgxpool.Pool`, `HTTP Handlers Layer (2)`, `User City Module`, `Route Registration & Guards (2)`, `decodeAndValidate`, `time.Time`, `Execute`, `User Supervisor Module (2)`, `Part Process DTOs (3)`, `Org Location/Shop DTOs (2)`, `User Line Module`, `DowntimeReasonService`, `IdleReasonService`, `.UploadDrawing`, `Response`, `Machine DTOs (2)`, `DTO Validation & Transactions (4)`, `.Validate`, `LicenceKeyService`, `Category`, `CompanyProfile`, `CreateAddressRequest`?**
  _High betweenness centrality (0.060) - this node is a cross-community bridge._
- **What connects `hello`, `7th-sense-server`, `CNCConnect param ip (const char*, ANSI NUL-terminated)` to the rest of the system?**
  _27 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `HTTP Handlers Layer` be split into smaller, more focused modules?**
  _Cohesion score 0.09859154929577464 - nodes in this community are weakly interconnected._
- **Why does `SendResponse()` connect `HTTP Handlers Layer (5)` to `HTTP Handlers Layer`, `testing.T`, `Handlers`, `Route Registration & Guards`, `HTTP Handlers Layer (2)`, `Route Registration & Guards (2)`, `decodeAndValidate`?**
  _High betweenness centrality (0.035) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.05244096769520498 - nodes in this community are weakly interconnected._
- **Why does `pathID()` connect `HTTP Handlers Layer (2)` to `HTTP Handlers Layer`, `Handlers`, `HTTP Handlers Layer (5)`, `Route Registration & Guards (2)`, `NewError`, `decodeAndValidate`?**
  _High betweenness centrality (0.022) - this node is a cross-community bridge._

### Low-confidence Hints
_AMBIGUOUS edges — the extractor was unsure. Verify before acting on these._

- **What is the exact relationship between `Machines API (/machines)` and `Snapshot Fields (positions, spindle, feed, program, tool, alarms)`?**
  _Edge tagged AMBIGUOUS (relation: conceptually_related_to) - confidence is low._