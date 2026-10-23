# Graph Report - 7th-sense-server  (2026-10-23)

## Corpus Check
- 238 files · ~115,431 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 4 file(s) not represented in the graph (top: (none) 3, .toml 1)

## Summary
- 1867 nodes · 15604 edges · 74 communities (59 shown, 15 thin omitted)
- Extraction: 95% EXTRACTED · 5% INFERRED · 0% AMBIGUOUS · INFERRED: 736 edges (avg confidence: 0.86)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `d9d2c113`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- HTTP Handlers Layer
- testing.T
- UserSupervisorService
- Org Location/Shop DTOs
- Route Registration & Guards
- Machine DTOs
- Shift
- CNC Snapshot Struct Fields
- validateCatalogFields
- ControllerService
- AuthUser
- Integration & DTO Tests (2)
- Collector
- CncCommunication.dll API & Integration Guide (Go Client)
- .UploadDrawing
- run
- HolidayService
- UserCity
- context.Context
- MonitoredMachine
- UserService
- Querier
- HTTP Handlers Layer (2)
- CNC Go Client Run & Check Mode
- user_access_test.go
- Connection
- 5. One machine: connect and collect (`internal/monitor/worker.go`)
- UserLineService
- validText
- Org Location/Shop DTOs (2)
- AuthService
- decodeAndValidate
- add_superadmin_raw.go
- DowntimeReasonService
- time.Duration
- CNC Connect API & Callbacks
- Server API Domain Docs
- HTTP Handlers Layer (5)
- .validate
- UserCountryService
- MachineModelService
- 7th Sense Server Engineering Principles (Skills Doc)
- time.Time
- GET /health
- .validate
- ProductionScheduleRepository
- CncCommunication.dll (C++ core runtime, x86, OPC UA inside)
- .Login
- NewError
- net/http.Handler
- CNC Event-Driven Data Flow
- Product
- Check-mode OPC UA node mapping
- .validate
- Sinumerik DLL Go Client (use case) (2)
- ControllerRepository
- Go Module Root
- Current User Middleware
- Role Guard Middleware
- SendResponse Helper
- Test Auth Helpers
- CncCommunication.dll API & Integration Guide (Go Client) - text version
- MachineTelemetryService
- statusRecorder
- Brand
- AuthHandler
- .Generate
- Execute
- newUserHandler
- MachineHandler

## God Nodes (most connected - your core abstractions)
1. `Execute()` - 169 edges
2. `SendResponse()` - 141 edges
3. `NewError()` - 130 edges
4. `run()` - 85 edges
5. `pathID()` - 77 edges
6. `AuthUser` - 69 edges
7. `Querier` - 58 edges
8. `decodeAndValidate()` - 58 edges
9. `currentUser()` - 57 edges
10. `CncMasterSnapshot (34,440 bytes, #pragma pack(8), x86)` - 38 edges

## Surprising Connections (you probably didn't know these)
- `Machine Identifier Mapping (0 Simulator,1 Fanuc,2 Siemens,3 Syntec,4 Mitsubishi)` --semantically_similar_to--> `Machine Catalog Brand -> Controller -> Model`  [INFERRED] [semantically similar]
  docs/CncCommunication.dll — API & Integration Guide (Go Client).pdf → .claude/CLAUDE.md
- `Snapshot Fields (positions, spindle, feed, program, tool, alarms)` --conceptually_related_to--> `Machines API (/machines)`  [AMBIGUOUS]
  docs/CncCommunication.dll — API & Integration Guide (Go Client).pdf → .claude/CLAUDE.md
- `5. One machine: connect and collect (`internal/monitor/worker.go`)` --references--> `Reading`  [INFERRED]
  docs/cnc-startup-flow.md → internal/cnc/cnc.go
- `5. One machine: connect and collect (`internal/monitor/worker.go`)` --references--> `Connection`  [INFERRED]
  docs/cnc-startup-flow.md → internal/cnc/cnc.go
- `3. Collector run` --references--> `WithTransaction()`  [INFERRED]
  docs/cnc-startup-flow.md → internal/db/postgres.go

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

## Communities (74 total, 15 thin omitted)

### Community 0 - "HTTP Handlers Layer"
Cohesion: 0.10
Nodes (11): BreakHandler, ClienteleHandler, DowntimeReasonHandler, HolidayHandler, IdleReasonHandler, PartHandler, ProductionScheduleHandler, currentUser() (+3 more)

### Community 1 - "testing.T"
Cohesion: 0.06
Nodes (97): UpdateThemeRequest, TestSnapshotSize(), CreateUserRequest, isTimeZone(), isWebsite(), normalizeFilters(), normalizeTheme(), TestCreateUserRequestValidate() (+89 more)

### Community 2 - "UserSupervisorService"
Cohesion: 0.14
Nodes (11): CreateUserSupervisorRequest, UpdateUserSupervisorRequest, validateUserSupervisor(), NewUserSupervisorHandler(), UserSupervisor, UserSupervisorRepository, mapUserSupervisorError(), NewUserSupervisorRepository() (+3 more)

### Community 3 - "Org Location/Shop DTOs"
Cohesion: 0.13
Nodes (13): CreateUserShopRequest, UpdateUserShopRequest, validateUserShop(), NewUserShopHandler(), UserShop, UserShopRepository, mapUserShopError(), NewUserShopRepository() (+5 more)

### Community 4 - "Route Registration & Guards"
Cohesion: 0.17
Nodes (30): RequireRole(), registerBranchRoutes(), registerBrandRoutes(), registerBreakRoutes(), registerCategoryRoutes(), registerClienteleRoutes(), registerControllerRoutes(), registerDowntimeReasonRoutes() (+22 more)

### Community 5 - "Machine DTOs"
Cohesion: 0.12
Nodes (17): CreateMachineRequest, MachineFields, UpdateMachineRequest, validCoordinate(), validHealthIDs(), validIPv4(), Machine, MachineScope (+9 more)

### Community 6 - "Shift"
Cohesion: 0.15
Nodes (19): clockMinutes(), CreateShiftRequest, ShiftFields, UpdateShiftRequest, MinutesBetween(), Shift, NewBreakRepository(), ShiftRepository (+11 more)

### Community 7 - "CNC Snapshot Struct Fields"
Cohesion: 0.11
Nodes (38): CncMasterSnapshot (34,440 bytes, #pragma pack(8), x86), AbsPos ([5]float64, Inches / Degrees, Work Coordinate System (WCS) positions), ActiveGCode ([16]byte, ASCII, Currently active G-code group), ActiveProgPath ([256]byte, ASCII Path, Storage path of active NC artifact), AlarmMsg ([128]byte, UTF-8 / ASCII, Alarm description message), AlarmNo (int32, Code (0 = OK), Primary alarm or fault number), AxisLoadXPercentage (int16, Percent (%), Servo X axis drive load (AxisLoadX/Y/ZPercentage)), AxisLoadYPercentage (int16, Percent (%), Servo Y axis drive load (AxisLoadX/Y/ZPercentage)) (+30 more)

### Community 8 - "validateCatalogFields"
Cohesion: 0.11
Nodes (11): CreateBrandRequest, UpdateBrandRequest, validateCatalogFields(), CreateCategoryRequest, UpdateCategoryRequest, CreateControllerRequest, UpdateControllerRequest, validateControllerFields() (+3 more)

### Community 10 - "AuthUser"
Cohesion: 0.08
Nodes (30): CreatePartRequest, PartFields, UpdatePartRequest, ValidPartDrawingFile(), AuthUser, Break, Clientele, IdleReason (+22 more)

### Community 11 - "Integration & DTO Tests (2)"
Cohesion: 0.43
Nodes (5): cncMasterSnapshot, LanguageResponse, ThemeResponse, authUserKey, machineModelRecord

### Community 12 - "Collector"
Cohesion: 0.20
Nodes (10): 0. The big picture, 3. Collector run, 4. Fetching the machines, 7. Following changes while running (no polling), 8. Server stop (Ctrl+C / SIGTERM), CNC Data Collection: Startup Flow, Logs you see at startup (`logs/app.log`), When it doesn't connect (+2 more)

### Community 13 - "CncCommunication.dll API & Integration Guide (Go Client)"
Cohesion: 0.13
Nodes (21): Machine Catalog Brand -> Controller -> Model, CNCConnect param machineType (int32 engine id, see §4.2), Engine 0: Simulator (local loopback mock engine), Engine 2: Siemens (Sinumerik OPC UA, active profile), Engine 3: Syntec (Open API / TCP), Engine 4: Mitsubishi (driver not specified), Machine identifier mapping (engine id enum 0-4), §9 Controller-Specific Profiles (+13 more)

### Community 15 - "run"
Cohesion: 0.14
Nodes (25): run(), Connect(), NewBreakHandler(), NewCategoryHandler(), NewClienteleHandler(), NewControllerHandler(), NewHolidayHandler(), NewIdleReasonHandler() (+17 more)

### Community 16 - "HolidayService"
Cohesion: 0.16
Nodes (10): CreateHolidayRequest, HolidayFields, UpdateHolidayRequest, Holiday, HolidayRepository, mapHolidayError(), HolidayService, holidayFromFields() (+2 more)

### Community 17 - "UserCity"
Cohesion: 0.18
Nodes (6): CreateUserCityRequest, UpdateUserCityRequest, validateUserCity(), UserCity, mapUserCityError(), toUserCityResponseError()

### Community 18 - "context.Context"
Cohesion: 0.06
Nodes (25): NewLicenceKey, Listen(), CreateLicenceKeysRequest, LicenceKey, PartProcess, PartSubProgram, ShiftBreak, UserProfile (+17 more)

### Community 19 - "MonitoredMachine"
Cohesion: 0.10
Nodes (14): dllDriver, slot, newConnection(), openDLL(), cstr(), cstrAt(), findLibDir(), openDLL() (+6 more)

### Community 20 - "UserService"
Cohesion: 0.06
Nodes (42): WithTransaction(), CustomerDetails, UpdateUserRequest, UserResponse, ParseDayStart(), TestParseDayStart(), Address, AssignableRoles() (+34 more)

### Community 21 - "Querier"
Cohesion: 0.13
Nodes (16): Querier, NewUserCityHandler(), NewUserStateHandler(), UserCityRepository, NewUserCityRepository(), NewUserCountryRepository(), UserMachineAccessRepository, NewUserMachineAccessRepository() (+8 more)

### Community 22 - "HTTP Handlers Layer (2)"
Cohesion: 0.09
Nodes (13): BranchHandler, ControllerHandler, pathID(), targetUserID(), UserCellHandler, UserCityHandler, UserCountryHandler, UserLineHandler (+5 more)

### Community 24 - "CNC Go Client Run & Check Mode"
Cohesion: 0.14
Nodes (21): Prerequisite: 32-bit x86 DLL, compile Go with GOARCH=386, Prerequisite: SIEMENS_USER / SIEMENS_PASS env vars before CNCConnect, §3.1 Environment & Prerequisites, §6 Snapshot Data Structure, §10 Diagnostics & Troubleshooting, Siemens credentials: SIEMENS_USER / SIEMENS_PASS env vars, Siemens trust: cert must be trusted in NCU Access Administration, Struct size assertion: unsafe.Sizeof(CncMasterSnapshot{}) == 34440 (+13 more)

### Community 25 - "user_access_test.go"
Cohesion: 0.38
Nodes (3): LicenceKeyDuration, licenceKeyDurations(), body

### Community 27 - "Connection"
Cohesion: 0.24
Nodes (4): Connection, Reading, Stats, Value

### Community 28 - "5. One machine: connect and collect (`internal/monitor/worker.go`)"
Cohesion: 0.23
Nodes (6): 5. One machine: connect and collect (`internal/monitor/worker.go`), FieldIndex(), VariableConfig, boundVariable, stopReason, worker

### Community 30 - "UserLineService"
Cohesion: 0.07
Nodes (21): CreateUserCellRequest, UpdateUserCellRequest, validateUserCell(), CreateUserLineRequest, UpdateUserLineRequest, validateUserLine(), NewUserCellHandler(), NewUserLineHandler() (+13 more)

### Community 31 - "validText"
Cohesion: 0.20
Nodes (5): PartProcessFields, PartSubProgramFields, uniquePositiveIDs(), validText(), validWorkingDays()

### Community 32 - "Org Location/Shop DTOs (2)"
Cohesion: 0.08
Nodes (18): CreateBranchRequest, UpdateBranchRequest, validateBranchFields(), CreateUserLocationRequest, UpdateUserLocationRequest, validateUserLocation(), NewBranchHandler(), NewUserLocationHandler() (+10 more)

### Community 33 - "AuthService"
Cohesion: 0.27
Nodes (8): NewAuthHandler(), Authenticate(), AuthRepository, NewAuthRepository(), AuthService, hashToken(), NewAuthService(), newAuthRouter()

### Community 34 - "decodeAndValidate"
Cohesion: 0.14
Nodes (5): validator, LicenceKeyHandler, currentUserID(), decodeAndValidate(), UserHandler

### Community 35 - "add_superadmin_raw.go"
Cohesion: 0.05
Nodes (39): main(), run(), main(), Config, 1. Server start, axisColor(), cstr(), enableVT() (+31 more)

### Community 36 - "DowntimeReasonService"
Cohesion: 0.15
Nodes (12): CreateDowntimeReasonRequest, DowntimeReasonFields, UpdateDowntimeReasonRequest, NewDowntimeReasonHandler(), DowntimeReason, DowntimeReasonRepository, mapDowntimeReasonError(), NewDowntimeReasonRepository() (+4 more)

### Community 37 - "time.Duration"
Cohesion: 0.33
Nodes (3): sleep(), backoff, Session

### Community 38 - "CNC Connect API & Callbacks"
Cohesion: 0.31
Nodes (10): Package-level callback singletons (syscall.NewCallback), Mutex-protected copy: dataMutex sync.RWMutex -> lastRecord/hasRecord, onDataCallback(ptr uintptr) uintptr, onErrorCallback(code, messagePtr uintptr) uintptr, §7 Asynchronous Callback Architecture, CncMasterSnapshot (34,440 bytes, #pragma pack(8)), onDataCallback, onErrorCallback (+2 more)

### Community 39 - "Server API Domain Docs"
Cohesion: 0.23
Nodes (13): Clienteles API (/clienteles), Company Scoping (companyOwner / ownCompanyCustomer), Customer Display Names (customerName / companyName), Licence Keys API (/licence-keys), Machines API (/machines), Organization Hierarchy APIs (/user-country..supervisor), Parts API (/parts with processes, drawings), Production Schedules API (/production-schedules) (+5 more)

### Community 41 - "HTTP Handlers Layer (5)"
Cohesion: 0.14
Nodes (5): BrandHandler, CategoryHandler, MachineModelHandler, ProductHandler, SendResponse()

### Community 43 - ".validate"
Cohesion: 0.47
Nodes (3): CreateIdleReasonRequest, IdleReasonFields, UpdateIdleReasonRequest

### Community 44 - "UserCountryService"
Cohesion: 0.17
Nodes (8): CreateUserCountryRequest, NewUserCountryHandler(), UserCountry, UserCountryRepository, mapUserCountryError(), UserCountryService, NewUserCountryService(), toUserCountryResponseError()

### Community 45 - "MachineModelService"
Cohesion: 0.18
Nodes (8): MachineModelRequest, MachineModel, MachineModelRepository, mapMachineModelError(), NewMachineModelRepository(), MachineModelService, NewMachineModelService(), toMachineModelResponseError()

### Community 46 - "7th Sense Server Engineering Principles (Skills Doc)"
Cohesion: 0.05
Nodes (22): Common Mistakes to Avoid, CRUD Implementation Checklist (11 steps), 7th Sense Server Engineering Principles (Skills Doc), Golden Reference Pattern: User CRUD, Unique Constraint Error Mapping (pgconn 23505), updated_at Trigger in Migrations, Air Hot Reload, 7th Sense Server - Claude Context (+14 more)

### Community 47 - "time.Time"
Cohesion: 0.14
Nodes (9): MachineEvent, MachineLiveStatus, ModeChange, TelemetryValue, MachineTelemetryRepository, execute(), MachineTelemetryService, UserMachine (+1 more)

### Community 49 - ".validate"
Cohesion: 0.38
Nodes (4): ClienteleFields, CreateClienteleRequest, UpdateClienteleRequest, validClienteleEmail()

### Community 50 - "ProductionScheduleRepository"
Cohesion: 0.10
Nodes (20): validDate(), CreateProductionScheduleRequest, ProductionScheduleFields, ProductionScheduleItemFields, UpdateProductionScheduleRequest, ScheduleDay(), ProductionSchedule, ProductionScheduleItem (+12 more)

### Community 51 - "CncCommunication.dll (C++ core runtime, x86, OPC UA inside)"
Cohesion: 0.24
Nodes (16): client_cert.der/.pem (X.509 client cert for Siemens OPC UA), client_key.der/.pem (client private key), CncCommunication.dll (C++ core runtime, x86, OPC UA inside), Fanuc profile (engine 1, port 8193, Fwlib32.dll + fwlibe1.dll), Fwlib32.dll (Fanuc FOCAS2), fwlibe1.dll (Fanuc FOCAS2), libcrypto-3.dll (OpenSSL), libssl-3.dll (OpenSSL) (+8 more)

### Community 52 - ".Login"
Cohesion: 0.19
Nodes (6): LoginRequest, newSessionToken(), truncate(), newLoginThrottle(), loginFailures, loginThrottle

### Community 53 - "NewError"
Cohesion: 0.11
Nodes (13): CreateAddressRequest, UpdateAddressRequest, UpdateLanguageRequest, UpdateLicenceKeyRequest, validDuration(), ValidLicenceKey(), CreateUserStateRequest, validateUserState() (+5 more)

### Community 54 - "net/http.Handler"
Cohesion: 0.19
Nodes (15): NewLicenceKeyHandler(), NewMachineHandler(), NewProductionScheduleHandler(), CrossOrigin(), Logger(), Chain(), Middleware, Recoverer() (+7 more)

### Community 55 - "CNC Event-Driven Data Flow"
Cohesion: 0.36
Nodes (8): Flow 1: load DLL, register two stdcall callbacks (syscall.NewCallback), Flow 2: CNCConnect starts DLL background worker thread, Flow 3: worker polls controller, assembles snapshot, invokes data callback, Flow 4: callback copies payload, signals renderer via buffered channel, Flow 5: main goroutine renders live ANSI dashboard, OPC UA (Siemens transport, port 4840), §2 System Architecture, SINUMERIK CNC (PLC/NCU)

### Community 56 - "Product"
Cohesion: 0.20
Nodes (8): NewProductHandler(), Product, ProductRepository, mapProductError(), ProductService, NewProductService(), productImageOrDefault(), toProductResponseError()

### Community 57 - "Check-mode OPC UA node mapping"
Cohesion: 0.19
Nodes (14): Check mode: go run sinumerik_client.go check (OPC UA node existence report), Go Client (sinumerik_client.go, callback -> copy -> render), Program status -> /Channel/State/acProg[u1,1], Operating mode -> /Bag/State/opMode[u1], Tool -> /Channel/State/actTNumber[u1,1], Check-mode OPC UA node mapping, Run: go run sinumerik_client.go (live dashboard), §8 Running the Go Client (+6 more)

### Community 58 - ".validate"
Cohesion: 0.47
Nodes (3): BreakFields, CreateBreakRequest, UpdateBreakRequest

### Community 64 - "ControllerRepository"
Cohesion: 0.22
Nodes (5): Controller, ControllerRepository, mapControllerError(), NewControllerRepository(), NewControllerService()

### Community 70 - "CncCommunication.dll API & Integration Guide (Go Client) - text version"
Cohesion: 0.22
Nodes (11): CNCConnect(ip, port, dataCb, errCb, machineType) int32, CNCConnect param dataCb (uintptr stdcall OnData(const CncMasterSnapshot*)), CNCConnect param errCb (uintptr stdcall OnErr(int32 code, const char* msg)), CNCConnect param ip (const char*, ANSI NUL-terminated), CNCConnect param port (int32, e.g. 4840), CNCdisconnect(machineId int32), CncCommunication.dll API & Integration Guide (Go Client) - text version, LazyDLL binding (syscall.NewLazyDLL + eager mod.Load in init) (+3 more)

### Community 71 - "MachineTelemetryService"
Cohesion: 0.18
Nodes (9): Supervisor, 2. CNC collection starts (still inside `run()`), Driver, Open(), NewSupervisor(), CNCConfig, NewCollector(), NewMachineTelemetryRepository() (+1 more)

### Community 76 - "Brand"
Cohesion: 0.23
Nodes (6): NewBrandHandler(), Brand, BrandRepository, mapBrandError(), BrandService, toBrandResponseError()

### Community 77 - "AuthHandler"
Cohesion: 0.38
Nodes (3): AuthHandler, CurrentUser(), UserID()

### Community 79 - ".Generate"
Cohesion: 0.40
Nodes (3): GenerateLicenceKeysRequest, GeneratedLicenceKeys, newLicenceKey()

### Community 83 - "Execute"
Cohesion: 0.06
Nodes (23): Execute(), UserListFilter, Branch, Category, Role, User, UserState, mapAuthError() (+15 more)

### Community 87 - "newUserHandler"
Cohesion: 0.36
Nodes (10): RoleHandler, NewRoleHandler(), NewUserHandler(), RoleRepository, NewRoleRepository(), RoleService, NewRoleService(), NewUserService() (+2 more)

## Ambiguous Edges - Review These
- `Snapshot Fields (positions, spindle, feed, program, tool, alarms)` → `Machines API (/machines)`  [AMBIGUOUS]
  docs/CncCommunication.dll — API & Integration Guide (Go Client).pdf · relation: conceptually_related_to

## Knowledge Gaps
- **25 isolated node(s):** `hello`, `7th-sense-server`, `Siemens default ports: 4840 (native NCU) / 4841 (custom)`, `Engine 0: Simulator (local loopback mock engine)`, `Engine 4: Mitsubishi (driver not specified)` (+20 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 86 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **15 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `NewError()` connect `NewError` to `HTTP Handlers Layer`, `testing.T`, `UserSupervisorService`, `Org Location/Shop DTOs`, `Machine DTOs`, `Shift`, `ControllerService`, `AuthUser`, `.UploadDrawing`, `HolidayService`, `UserCity`, `context.Context`, `UserService`, `HTTP Handlers Layer (2)`, `user_access_test.go`, `UserLineService`, `validText`, `Org Location/Shop DTOs (2)`, `decodeAndValidate`, `DowntimeReasonService`, `.validate`, `UserCountryService`, `MachineModelService`, `.validate`, `ProductionScheduleRepository`, `.Login`, `Product`, `.validate`, `Brand`, `.Generate`, `Execute`?**
  _High betweenness centrality (0.066) - this node is a cross-community bridge._
- **What connects `hello`, `7th-sense-server`, `Siemens default ports: 4840 (native NCU) / 4841 (custom)` to the rest of the system?**
  _25 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `HTTP Handlers Layer` be split into smaller, more focused modules?**
  _Cohesion score 0.09979296066252588 - nodes in this community are weakly interconnected._
- **Why does `SendResponse()` connect `HTTP Handlers Layer (5)` to `HTTP Handlers Layer`, `AuthService`, `decodeAndValidate`, `Route Registration & Guards`, `AuthHandler`, `HTTP Handlers Layer (2)`, `net/http.Handler`, `user_access_test.go`, `MachineHandler`?**
  _High betweenness centrality (0.027) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.05624438454627134 - nodes in this community are weakly interconnected._
- **Why does `CurrentUser()` connect `AuthHandler` to `HTTP Handlers Layer`, `Route Registration & Guards`, `AuthUser`, `Integration & DTO Tests (2)`, `context.Context`, `HTTP Handlers Layer (2)`?**
  _High betweenness centrality (0.024) - this node is a cross-community bridge._

### Low-confidence Hints
_AMBIGUOUS edges — the extractor was unsure. Verify before acting on these._

- **What is the exact relationship between `Snapshot Fields (positions, spindle, feed, program, tool, alarms)` and `Machines API (/machines)`?**
  _Edge tagged AMBIGUOUS (relation: conceptually_related_to) - confidence is low._