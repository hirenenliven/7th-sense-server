# Graph Report - 7th-sense-server  (2026-10-28)

## Corpus Check
- 250 files · ~129,987 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 4 file(s) not represented in the graph (top: (none) 3, .toml 1)

## Summary
- 2007 nodes · 16524 edges · 88 communities (72 shown, 16 thin omitted)
- Extraction: 95% EXTRACTED · 5% INFERRED · 0% AMBIGUOUS · INFERRED: 807 edges (avg confidence: 0.86)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `8d6e55a7`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- HTTP Handlers Layer
- testing.T
- .collectReadings
- worker
- Route Registration & Guards
- MachineService
- net/http.Handler
- CNC Snapshot Struct Fields
- .validate
- Querier
- UserCityService
- Integration & DTO Tests (2)
- UserShopService
- CncCommunication.dll API & Integration Guide (Go Client)
- LicenceKeyService
- run
- NewError
- mode-chart.js
- add_superadmin_raw.go
- Connection
- UserService
- validateCatalogFields
- HTTP Handlers Layer (2)
- DowntimeReasonService
- CNC Go Client Run & Check Mode
- user_access_test.go
- ClienteleService
- UserLineService
- time.Time
- Handlers
- Category
- AuthUser
- newUserHandler
- context.Context
- UserLocationService
- sinumerik_client.go
- Execute
- .log
- CNC Connect API & Callbacks
- Server API Domain Docs
- net/http.Request
- SendResponse
- config.go
- Collector
- UserCountryService
- IdleReasonService
- CRUD Implementation Pattern (8 steps)
- auth_tokens SHA-256 Hashed Sessions (migration 014)
- GET /health
- MachineModelService
- ControllerService
- CncCommunication.dll (C++ core runtime, x86, OPC UA inside)
- .Login
- generate.go
- CNC Event-Driven Data Flow
- demo_monitor/main.go
- Check-mode OPC UA node mapping
- .Generate
- Create Customer in One Transaction (user, company profile, user profile, address)
- machineState
- Sinumerik DLL Go Client (use case) (2)
- Adding a New Entity (README steps)
- middleware.CrossOrigin (http.CrossOriginProtection)
- UserStateService
- Go Module Root
- Current User Middleware
- Role Guard Middleware
- SendResponse Helper
- Test Auth Helpers
- CncCommunication.dll API & Integration Guide (Go Client) - text version
- CNCConfig
- VariableConfig
- Login Lockout (5 wrong passwords, 15 min, 429)
- Naming Conventions
- validText
- BranchService
- AuthHandler
- UserCellService
- .validate
- UserRepository
- newPartRouter
- go_pkg_log_slog
- .Validate
- .validate
- AuthService
- .UploadDrawing
- CreateUserStateRequest

## God Nodes (most connected - your core abstractions)
1. `Execute()` - 172 edges
2. `SendResponse()` - 143 edges
3. `NewError()` - 130 edges
4. `run()` - 90 edges
5. `pathID()` - 78 edges
6. `AuthUser` - 69 edges
7. `Querier` - 58 edges
8. `decodeAndValidate()` - 58 edges
9. `currentUser()` - 57 edges
10. `worker` - 42 edges

## Surprising Connections (you probably didn't know these)
- `Machine Identifier Mapping (0 Simulator,1 Fanuc,2 Siemens,3 Syntec,4 Mitsubishi)` --semantically_similar_to--> `Machine Catalog Brand -> Controller -> Model`  [INFERRED] [semantically similar]
  docs/CncCommunication.dll — API & Integration Guide (Go Client).pdf → .claude/CLAUDE.md
- `Snapshot Fields (positions, spindle, feed, program, tool, alarms)` --conceptually_related_to--> `Machines API (/machines)`  [AMBIGUOUS]
  docs/CncCommunication.dll — API & Integration Guide (Go Client).pdf → .claude/CLAUDE.md
- `Step 13 — The DLL sends data` --references--> `Reading`  [INFERRED]
  docs/cnc-data-collection-flow.md → internal/cnc/cnc.go
- `2. All steps at a glance` --references--> `LoadDriver()`  [INFERRED]
  docs/cnc-data-collection-flow.md → internal/cnc/cnc.go
- `Step 3 — Run the collector in the background` --references--> `Collector`  [INFERRED]
  docs/cnc-data-collection-flow.md → internal/monitor/collector.go

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

## Communities (88 total, 16 thin omitted)

### Community 0 - "HTTP Handlers Layer"
Cohesion: 0.09
Nodes (10): BreakHandler, ClienteleHandler, DowntimeReasonHandler, HolidayHandler, IdleReasonHandler, PartHandler, ProductionScheduleHandler, currentUser() (+2 more)

### Community 1 - "testing.T"
Cohesion: 0.05
Nodes (103): UpdateThemeRequest, TestSnapshotSize(), CreateUserRequest, isTimeZone(), isWebsite(), normalizeFilters(), normalizeTheme(), TestCreateUserRequestValidate() (+95 more)

### Community 2 - ".collectReadings"
Cohesion: 0.10
Nodes (16): 8. Reference, Events in `machine_events`, Related documents, Retry and timing summary, Settings (`.env`, defaults in `internal/config/config.go`), Step 14 — The data collection loop, Step 16 — Every second: dead connection check and live status, Tables (+8 more)

### Community 3 - "worker"
Cohesion: 0.11
Nodes (13): Known issue, Step 8 — Load all active machines, Reading, Status, Value, equalEngine(), MonitoredMachine, snapshotStatusOf() (+5 more)

### Community 4 - "Route Registration & Guards"
Cohesion: 0.17
Nodes (30): RequireRole(), registerBranchRoutes(), registerBrandRoutes(), registerBreakRoutes(), registerCategoryRoutes(), registerClienteleRoutes(), registerControllerRoutes(), registerDowntimeReasonRoutes() (+22 more)

### Community 5 - "MachineService"
Cohesion: 0.11
Nodes (18): CreateMachineRequest, MachineFields, UpdateMachineRequest, validCoordinate(), validHealthIDs(), validIPv4(), SeesAllCompanyMachines(), Machine (+10 more)

### Community 6 - "net/http.Handler"
Cohesion: 0.18
Nodes (18): NewHolidayHandler(), NewMachineHandler(), NewProductionScheduleHandler(), NewShiftHandler(), CrossOrigin(), Logger(), Chain(), Middleware (+10 more)

### Community 7 - "CNC Snapshot Struct Fields"
Cohesion: 0.11
Nodes (36): CncMasterSnapshot (34,440 bytes, #pragma pack(8), x86), AbsPos ([5]float64, Inches / Degrees, Work Coordinate System (WCS) positions), ActiveGCode ([16]byte, ASCII, Currently active G-code group), ActiveProgPath ([256]byte, ASCII Path, Storage path of active NC artifact), AlarmMsg ([128]byte, UTF-8 / ASCII, Alarm description message), AlarmNo (int32, Code (0 = OK), Primary alarm or fault number), AxisLoadXPercentage (int16, Percent (%), Servo X axis drive load (AxisLoadX/Y/ZPercentage)), AxisLoadYPercentage (int16, Percent (%), Servo Y axis drive load (AxisLoadX/Y/ZPercentage)) (+28 more)

### Community 8 - ".validate"
Cohesion: 0.47
Nodes (3): CreateHolidayRequest, HolidayFields, UpdateHolidayRequest

### Community 9 - "Querier"
Cohesion: 0.11
Nodes (21): Querier, NewBranchHandler(), NewUserCityHandler(), NewUserLocationHandler(), NewUserStateHandler(), NewDowntimeReasonRepository(), NewIdleReasonRepository(), NewMachineTelemetryRepository() (+13 more)

### Community 10 - "UserCityService"
Cohesion: 0.16
Nodes (8): CreateUserCityRequest, UpdateUserCityRequest, validateUserCity(), UserCity, UserCityRepository, mapUserCityError(), UserCityService, toUserCityResponseError()

### Community 11 - "Integration & DTO Tests (2)"
Cohesion: 0.43
Nodes (3): cncMasterSnapshot, authUserKey, machineModelRecord

### Community 12 - "UserShopService"
Cohesion: 0.15
Nodes (9): CreateUserShopRequest, UpdateUserShopRequest, validateUserShop(), NewUserShopHandler(), UserShop, UserShopRepository, mapUserShopError(), UserShopService (+1 more)

### Community 13 - "CncCommunication.dll API & Integration Guide (Go Client)"
Cohesion: 0.13
Nodes (24): Machine Catalog Brand -> Controller -> Model, client_cert.der/.pem (X.509 client cert for Siemens OPC UA), client_key.der/.pem (client private key), Engine 0: Simulator (local loopback mock engine), Engine 2: Siemens (Sinumerik OPC UA, active profile), Engine 3: Syntec (Open API / TCP), Engine 4: Mitsubishi (driver not specified), Machine identifier mapping (engine id enum 0-4) (+16 more)

### Community 14 - "LicenceKeyService"
Cohesion: 0.20
Nodes (8): NewLicenceKeyHandler(), LicenceKey, LicenceKeyRepository, mapLicenceKeyError(), LicenceKeyService, licenceExpiresOn(), NewLicenceKeyService(), toLicenceKeyResponseError()

### Community 15 - "run"
Cohesion: 0.18
Nodes (24): run(), NewBrandHandler(), NewCategoryHandler(), NewControllerHandler(), NewMachineModelHandler(), NewProductHandler(), NewBranchRepository(), NewBrandRepository() (+16 more)

### Community 16 - "NewError"
Cohesion: 0.19
Nodes (7): CreateAddressRequest, UpdateAddressRequest, UpdateLanguageRequest, NewError(), newResponse(), asInvalidSelection(), Response

### Community 17 - "mode-chart.js"
Cohesion: 0.14
Nodes (39): cell(), chartData(), chartRows(), countByGroup(), currentSpan(), currentWindow(), drawChart(), drawDashboard() (+31 more)

### Community 18 - "add_superadmin_raw.go"
Cohesion: 0.38
Nodes (10): checkNotTaken(), createSuperadmin(), ensureProfile(), execute(), main(), parseFlags(), run(), superadminRoleID() (+2 more)

### Community 19 - "Connection"
Cohesion: 0.10
Nodes (15): dllDriver, slot, 5. Connect and collect, per machine: steps 10 – 17, Step 11 — Connect, a few machines at a time, Step 12 — Call the DLL's CNCConnect, Step 13 — The DLL sends data, Step 17 — Disconnect and try again, Connection (+7 more)

### Community 20 - "UserService"
Cohesion: 0.06
Nodes (34): CustomerDetails, UpdateUserRequest, UserResponse, ParseDayStart(), TestParseDayStart(), Address, AssignableRoles(), CanAssignRole() (+26 more)

### Community 21 - "validateCatalogFields"
Cohesion: 0.11
Nodes (11): CreateBrandRequest, UpdateBrandRequest, validateCatalogFields(), CreateCategoryRequest, UpdateCategoryRequest, CreateControllerRequest, UpdateControllerRequest, validateControllerFields() (+3 more)

### Community 22 - "HTTP Handlers Layer (2)"
Cohesion: 0.10
Nodes (12): validator, ControllerHandler, decodeAndValidate(), pathID(), targetUserID(), UserCellHandler, UserCityHandler, UserLineHandler (+4 more)

### Community 23 - "DowntimeReasonService"
Cohesion: 0.15
Nodes (11): CreateDowntimeReasonRequest, DowntimeReasonFields, UpdateDowntimeReasonRequest, NewDowntimeReasonHandler(), DowntimeReason, DowntimeReasonRepository, mapDowntimeReasonError(), downtimeReasonFromFields() (+3 more)

### Community 24 - "CNC Go Client Run & Check Mode"
Cohesion: 0.13
Nodes (22): CncCommunication.dll API & Integration Guide (Go Client) - text version, IsConnected (int16, Boolean (0/1), Network connection state), Prerequisite: 32-bit x86 DLL, compile Go with GOARCH=386, §1 Overview, §3.1 Environment & Prerequisites, §6 Snapshot Data Structure, §10 Diagnostics & Troubleshooting, Struct size assertion: unsafe.Sizeof(CncMasterSnapshot{}) == 34440 (+14 more)

### Community 25 - "user_access_test.go"
Cohesion: 0.34
Nodes (5): LanguageResponse, LicenceKeyDuration, ThemeResponse, licenceKeyDurations(), body

### Community 26 - "ClienteleService"
Cohesion: 0.13
Nodes (12): ClienteleFields, CreateClienteleRequest, UpdateClienteleRequest, validClienteleEmail(), NewClienteleHandler(), Clientele, ClienteleRepository, mapClienteleError() (+4 more)

### Community 27 - "UserLineService"
Cohesion: 0.15
Nodes (9): CreateUserLineRequest, UpdateUserLineRequest, validateUserLine(), NewUserLineHandler(), UserLine, UserLineRepository, mapUserLineError(), UserLineService (+1 more)

### Community 28 - "time.Time"
Cohesion: 0.09
Nodes (16): MachineModeHistoryResponse, MachineModesResponse, MachineModeTotal, LiveValue, MachineEvent, MachineLiveStatus, MachineModePeriod, MachineModeStatus (+8 more)

### Community 29 - "Handlers"
Cohesion: 0.18
Nodes (6): HealthHandler, NewHealthHandler(), UserCountryHandler, registerHealthRoutes(), HealthService, Handlers

### Community 30 - "Category"
Cohesion: 0.25
Nodes (5): Category, CategoryRepository, mapCategoryError(), CategoryService, toCategoryResponseError()

### Community 31 - "AuthUser"
Cohesion: 0.07
Nodes (25): ValidPartDrawingFile(), CreateProductionScheduleRequest, ProductionScheduleFields, UpdateProductionScheduleRequest, AuthUser, Holiday, ProductionSchedule, ProductionScheduleItem (+17 more)

### Community 32 - "newUserHandler"
Cohesion: 0.36
Nodes (10): RoleHandler, NewRoleHandler(), NewUserHandler(), RoleRepository, NewRoleRepository(), RoleService, NewRoleService(), NewUserService() (+2 more)

### Community 33 - "context.Context"
Cohesion: 0.07
Nodes (25): NewBreakHandler(), Break, Part, PartProcess, PartSubProgram, ShiftBreak, UserProfile, BreakRepository (+17 more)

### Community 34 - "UserLocationService"
Cohesion: 0.16
Nodes (8): CreateUserLocationRequest, UpdateUserLocationRequest, validateUserLocation(), UserLocation, UserLocationRepository, mapUserLocationError(), UserLocationService, toUserLocationResponseError()

### Community 35 - "sinumerik_client.go"
Cohesion: 0.20
Nodes (17): axisColor(), cstr(), enableVT(), findLibDir(), initDLL(), isRotary(), lbl(), logf() (+9 more)

### Community 36 - "Execute"
Cohesion: 0.06
Nodes (20): Execute(), Brand, Product, Role, AuthRepository, mapAuthError(), BrandRepository, mapBrandError() (+12 more)

### Community 37 - ".log"
Cohesion: 0.44
Nodes (6): Log lines (in `logs/app.log`, search for the message), CNCLogSettings, isLogTypeOn(), TestWriteLog(), writeLog(), logType

### Community 38 - "CNC Connect API & Callbacks"
Cohesion: 0.42
Nodes (8): Package-level callback singletons (syscall.NewCallback), Mutex-protected copy: dataMutex sync.RWMutex -> lastRecord/hasRecord, onDataCallback(ptr uintptr) uintptr, onErrorCallback(code, messagePtr uintptr) uintptr, §7 Asynchronous Callback Architecture, onDataCallback, onErrorCallback, stdcall Callbacks via syscall.NewCallback

### Community 39 - "Server API Domain Docs"
Cohesion: 0.23
Nodes (13): Clienteles API (/clienteles), Company Scoping (companyOwner / ownCompanyCustomer), Customer Display Names (customerName / companyName), Licence Keys API (/licence-keys), Machines API (/machines), Organization Hierarchy APIs (/user-country..supervisor), Parts API (/parts with processes, drawings), Production Schedules API (/production-schedules) (+5 more)

### Community 40 - "net/http.Request"
Cohesion: 0.13
Nodes (6): BranchHandler, MachineHandler, currentUserID(), listUserID(), queryID(), UserHandler

### Community 41 - "SendResponse"
Cohesion: 0.10
Nodes (6): BrandHandler, CategoryHandler, LicenceKeyHandler, MachineModelHandler, ProductHandler, SendResponse()

### Community 42 - "config.go"
Cohesion: 0.20
Nodes (12): main(), run(), Config, MachineSelection, getEnv(), getEnvBool(), getEnvInt(), getEnvMachineSelection() (+4 more)

### Community 43 - "Collector"
Cohesion: 0.14
Nodes (18): 1. The big picture, 2. All steps at a glance, 4. Find the machines: steps 5 – 9, 6. When data is edited: step 18, 7. Shutdown: steps 19 – 20, 9. Data Reading Flow (Function Names Only), CNC Data Collection Flow, Main parts (+10 more)

### Community 44 - "UserCountryService"
Cohesion: 0.18
Nodes (9): CreateUserCountryRequest, NewUserCountryHandler(), UserCountry, UserCountryRepository, mapUserCountryError(), NewUserCountryRepository(), UserCountryService, NewUserCountryService() (+1 more)

### Community 45 - "IdleReasonService"
Cohesion: 0.15
Nodes (11): CreateIdleReasonRequest, IdleReasonFields, UpdateIdleReasonRequest, NewIdleReasonHandler(), IdleReason, IdleReasonRepository, mapIdleReasonError(), IdleReasonService (+3 more)

### Community 46 - "CRUD Implementation Pattern (8 steps)"
Cohesion: 0.14
Nodes (8): Air Hot Reload, 7th Sense Server - Claude Context, Engineering Principles Summary, PostgreSQL + TimescaleDB via pgx/v5, SendResponse JSON Response Format, User CRUD Golden Reference, SampleTimestampMs (int64, Milliseconds, Epoch timestamp of snapshot emission), README Layered Layout (Node.js analogy)

### Community 47 - "auth_tokens SHA-256 Hashed Sessions (migration 014)"
Cohesion: 0.67
Nodes (3): auth_tokens SHA-256 Hashed Sessions (migration 014), Per-Request Session and User Re-check, ts_session Cookie (HttpOnly, SameSite=Strict)

### Community 49 - "MachineModelService"
Cohesion: 0.20
Nodes (6): MachineModelRequest, MachineModel, MachineModelRepository, mapMachineModelError(), MachineModelService, toMachineModelResponseError()

### Community 50 - "ControllerService"
Cohesion: 0.21
Nodes (5): Controller, ControllerRepository, mapControllerError(), ControllerService, toControllerResponseError()

### Community 51 - "CncCommunication.dll (C++ core runtime, x86, OPC UA inside)"
Cohesion: 0.22
Nodes (18): CncCommunication.dll (C++ core runtime, x86, OPC UA inside), Fanuc profile (engine 1, port 8193, Fwlib32.dll + fwlibe1.dll), Fwlib32.dll (Fanuc FOCAS2), fwlibe1.dll (Fanuc FOCAS2), libcrypto-3.dll (OpenSSL), libssl-3.dll (OpenSSL), Engine 1: Fanuc (FOCAS2 Ethernet, fwlib32), OPC UA (Siemens transport, port 4840) (+10 more)

### Community 52 - ".Login"
Cohesion: 0.19
Nodes (6): LoginRequest, newSessionToken(), truncate(), newLoginThrottle(), loginFailures, loginThrottle

### Community 55 - "CNC Event-Driven Data Flow"
Cohesion: 0.60
Nodes (5): Flow 1: load DLL, register two stdcall callbacks (syscall.NewCallback), Flow 2: CNCConnect starts DLL background worker thread, Flow 3: worker polls controller, assembles snapshot, invokes data callback, Flow 4: callback copies payload, signals renderer via buffered channel, Flow 5: main goroutine renders live ANSI dashboard

### Community 56 - "demo_monitor/main.go"
Cohesion: 0.22
Nodes (7): CncMasterSnapshot, MachineDef, cstr(), enableVT(), initDLL(), main(), render()

### Community 57 - "Check-mode OPC UA node mapping"
Cohesion: 0.19
Nodes (14): Check mode: go run sinumerik_client.go check (OPC UA node existence report), GcodeGroup ([512]byte, ASCII (reused), DLL reuses this field for the ';'-separated missing-node list), MachineMode (int16, Enum (0=JOG, 1=MDA, 2=AUTO), Controller operating mode), Go Client (sinumerik_client.go, callback -> copy -> render), Operating mode -> /Bag/State/opMode[u1], Check-mode OPC UA node mapping, Run: go run sinumerik_client.go (live dashboard), §8 Running the Go Client (+6 more)

### Community 58 - ".Generate"
Cohesion: 0.40
Nodes (3): GenerateLicenceKeysRequest, GeneratedLicenceKeys, newLicenceKey()

### Community 60 - "machineState"
Cohesion: 0.22
Nodes (11): Step 15 — Store the reading, decideMachineMode(), holdReason(), isInProduction(), machineStateFromReading(), programStateFromProcStatus(), TestMachineStateFromReading(), modeName() (+3 more)

### Community 64 - "UserStateService"
Cohesion: 0.27
Nodes (4): UserState, mapUserStateError(), UserStateService, toUserStateResponseError()

### Community 70 - "CncCommunication.dll API & Integration Guide (Go Client) - text version"
Cohesion: 0.20
Nodes (11): CNCConnect(ip, port, dataCb, errCb, machineType) int32, CNCConnect param dataCb (uintptr stdcall OnData(const CncMasterSnapshot*)), CNCConnect param errCb (uintptr stdcall OnErr(int32 code, const char* msg)), CNCConnect param ip (const char*, ANSI NUL-terminated), CNCConnect param machineType (int32 engine id, see §4.2), CNCConnect param port (int32, e.g. 4840), CNCdisconnect(machineId int32), LazyDLL binding (syscall.NewLazyDLL + eager mod.Load in init) (+3 more)

### Community 71 - "CNCConfig"
Cohesion: 0.22
Nodes (10): Supervisor, 3. Startup: steps 1 – 4, Step 1 — Start collection with the server, Step 2 — Load the DLL, Step 3 — Run the collector in the background, Driver, LoadDriver(), NewSupervisor() (+2 more)

### Community 72 - "VariableConfig"
Cohesion: 0.15
Nodes (10): Step 10 — Get ready to connect, FieldIndex(), TestFieldIndexAxisPositions(), VariableConfig, formatPreviousValue(), formatValue(), newValueChangeLog(), boundVariable (+2 more)

### Community 75 - "validText"
Cohesion: 0.29
Nodes (6): PartProcessFields, PartSubProgramFields, CreatePartRequest, PartFields, UpdatePartRequest, validText()

### Community 76 - "BranchService"
Cohesion: 0.15
Nodes (9): CreateBranchRequest, UpdateBranchRequest, validateBranchFields(), Branch, BranchRepository, mapBranchError(), cityOwner(), BranchService (+1 more)

### Community 77 - "AuthHandler"
Cohesion: 0.38
Nodes (3): AuthHandler, CurrentUser(), UserID()

### Community 78 - "UserCellService"
Cohesion: 0.07
Nodes (20): CreateUserCellRequest, UpdateUserCellRequest, validateUserCell(), CreateUserSupervisorRequest, UpdateUserSupervisorRequest, validateUserSupervisor(), NewUserCellHandler(), NewUserSupervisorHandler() (+12 more)

### Community 79 - ".validate"
Cohesion: 0.24
Nodes (6): validDate(), ProductionScheduleItemFields, ScheduleDay(), checkJobsOverlap(), shiftSpan(), jobSpan

### Community 80 - "UserRepository"
Cohesion: 0.08
Nodes (27): WithTransaction(), uniquePositiveIDs(), clockMinutes(), CreateShiftRequest, ShiftFields, UpdateShiftRequest, MinutesBetween(), validWorkingDays() (+19 more)

### Community 81 - "newPartRouter"
Cohesion: 0.29
Nodes (6): NewPartHandler(), NewClienteleRepository(), PartDrawingStorage, NewPartDrawingStorage(), NewPartService(), newPartRouter()

### Community 82 - "go_pkg_log_slog"
Cohesion: 0.33
Nodes (3): main(), Init(), InitFile()

### Community 83 - ".Validate"
Cohesion: 0.29
Nodes (5): NewLicenceKey, CreateLicenceKeysRequest, UpdateLicenceKeyRequest, validDuration(), ValidLicenceKey()

### Community 84 - ".validate"
Cohesion: 0.47
Nodes (3): BreakFields, CreateBreakRequest, UpdateBreakRequest

### Community 85 - "AuthService"
Cohesion: 0.47
Nodes (6): NewAuthHandler(), Authenticate(), NewAuthRepository(), AuthService, NewAuthService(), newAuthRouter()

## Ambiguous Edges - Review These
- `Snapshot Fields (positions, spindle, feed, program, tool, alarms)` → `Machines API (/machines)`  [AMBIGUOUS]
  docs/CncCommunication.dll — API & Integration Guide (Go Client).pdf · relation: conceptually_related_to

## Knowledge Gaps
- **28 isolated node(s):** `hello`, `7th-sense-server`, `Siemens default ports: 4840 (native NCU) / 4841 (custom)`, `Engine 0: Simulator (local loopback mock engine)`, `Engine 4: Mitsubishi (driver not specified)` (+23 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 87 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **16 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `NewError()` connect `NewError` to `HTTP Handlers Layer`, `testing.T`, `MachineService`, `.validate`, `UserCityService`, `UserShopService`, `LicenceKeyService`, `UserService`, `HTTP Handlers Layer (2)`, `DowntimeReasonService`, `user_access_test.go`, `ClienteleService`, `UserLineService`, `Category`, `AuthUser`, `context.Context`, `UserLocationService`, `Execute`, `net/http.Request`, `UserCountryService`, `IdleReasonService`, `MachineModelService`, `ControllerService`, `.Login`, `.Generate`, `UserStateService`, `validText`, `BranchService`, `UserCellService`, `.validate`, `UserRepository`, `.Validate`, `.validate`, `.UploadDrawing`, `CreateUserStateRequest`?**
  _High betweenness centrality (0.056) - this node is a cross-community bridge._
- **What connects `hello`, `7th-sense-server`, `Siemens default ports: 4840 (native NCU) / 4841 (custom)` to the rest of the system?**
  _28 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `HTTP Handlers Layer` be split into smaller, more focused modules?**
  _Cohesion score 0.08506493506493507 - nodes in this community are weakly interconnected._
- **Why does `SendResponse()` connect `SendResponse` to `HTTP Handlers Layer`, `Route Registration & Guards`, `net/http.Handler`, `net/http.Request`, `AuthHandler`, `AuthService`, `HTTP Handlers Layer (2)`, `user_access_test.go`, `Handlers`?**
  _High betweenness centrality (0.030) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.0518018018018018 - nodes in this community are weakly interconnected._
- **Why does `CurrentUser()` connect `AuthHandler` to `HTTP Handlers Layer`, `context.Context`, `Route Registration & Guards`, `net/http.Request`, `Integration & DTO Tests (2)`, `HTTP Handlers Layer (2)`, `AuthUser`?**
  _High betweenness centrality (0.023) - this node is a cross-community bridge._

### Low-confidence Hints
_AMBIGUOUS edges — the extractor was unsure. Verify before acting on these._

- **What is the exact relationship between `Snapshot Fields (positions, spindle, feed, program, tool, alarms)` and `Machines API (/machines)`?**
  _Edge tagged AMBIGUOUS (relation: conceptually_related_to) - confidence is low._