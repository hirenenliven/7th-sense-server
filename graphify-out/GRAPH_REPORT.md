# Graph Report - 7th-sense-server  (2026-11-03)

## Corpus Check
- 266 files · ~151,733 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 4 file(s) not represented in the graph (top: (none) 3, .toml 1)

## Summary
- 2275 nodes · 18671 edges · 90 communities (69 shown, 21 thin omitted)
- Extraction: 95% EXTRACTED · 5% INFERRED · 0% AMBIGUOUS · INFERRED: 881 edges (avg confidence: 0.86)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `04c96001`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- LicenceKeyDuration
- MachineTelemetryService
- UserSupervisorService
- WithTransaction
- SendResponse()
- pathID()
- net/http.ResponseWriter
- UserCityService
- api_helpers_test.go
- CncMasterSnapshot (34,440 bytes, #pragma
- context.Context
- Connection
- UserLineService
- Collector
- worker
- AuthUser
- newMachineService
- testing.T
- time.Time
- How a cycle ends
- net/http.ServeMux
- CNCConnect(ip, port, dataCb, errCb, machineType) int32
- CompanyProfile
- §10 Diagnostics & Troubleshooting
- dto_test.go
- machine_mode.go
- BranchService
- UserService
- CNC Data Collection Flow
- toUserLocationResponseError
- .Create
- AuthUser
- Worker (internal/monitor/worker.go)
- Querier
- MachineModelService
- NewError
- LicenceKeyService
- Client Certificate & Key (client_cert.der / client_key.der)
- .validate
- mode-chart.js
- CncCommunication.dll API & Integration G
- sinumerik_client.go
- run
- CncCommunication.dll API & Integration G
- .validate
- .Login()
- .validate
- MachineSelection
- Machines API (/machines)
- MachineCycle
- UserShopService
- Execute
- containsID
- AuthHandler
- containsMachineID
- net/http.ResponseWriter
- Category
- .validate
- config.go
- Test list
- validText()
- net/http.Handler
- cycleTracker
- ProductionScheduleRepository
- generate.go
- .Validate
- demo_monitor/main.go
- AuthService
- onDataCallback(ptr uintptr) uintptr
- HolidayService
- logger/logger.go
- .validate
- Event-driven data flow (no polling in client)
- auth_tokens SHA-256 Hashed Sessions (mig
- Create Customer in One Transaction (user
- 7th-sense-server
- hello
- Adding a New Entity (README steps)
- middleware.CrossOrigin (http.CrossOrigin
- middleware.CurrentUser / currentUserID
- GET /health
- Login Lockout (5 wrong passwords, 15 min
- Naming Conventions
- middleware.RequireRole / routes.adminOnl
- response.SendResponse
- Test Auth (asTestUser, tests/auth_test.g
- user_filters Filter Checklist

## God Nodes (most connected - your core abstractions)
1. `Execute()` - 180 edges
2. `SendResponse()` - 147 edges
3. `NewError()` - 134 edges
4. `run()` - 96 edges
5. `pathID()` - 80 edges
6. `AuthUser` - 78 edges
7. `Querier` - 60 edges
8. `currentUser()` - 60 edges
9. `decodeAndValidate()` - 58 edges
10. `worker` - 56 edges

## Surprising Connections (you probably didn't know these)
- `Machine Identifier Mapping (0 Simulator,1 Fanuc,2 Siemens,3 Syntec,4 Mitsubishi)` --semantically_similar_to--> `Machine Catalog Brand-Controller-Model`  [INFERRED] [semantically similar]
  docs/CncCommunication.dll — API & Integration Guide (Go Client).pdf → .claude/CLAUDE.md
- `Snapshot Fields (positions, spindle, feed, program, tool, alarms)` --conceptually_related_to--> `Machines API (/machines)`  [AMBIGUOUS]
  docs/CncCommunication.dll — API & Integration Guide (Go Client).pdf → .claude/CLAUDE.md
- `How a cycle ends` --references--> `decideCycleStep()`  [INFERRED]
  docs/machine-cycles.md → internal/monitor/machine_cycle.go
- `One running cycle per machine` --references--> `decideCycleStep()`  [INFERRED]
  docs/machine-cycles.md → internal/monitor/machine_cycle.go
- `Step 15 — Store the reading` --references--> `machineStateFromReading()`  [INFERRED]
  docs/cnc-data-collection-flow.md → internal/monitor/machine_mode.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **CNC collection pipeline: Supervisor -> Collector -> Worker -> Driver -> DLL** — docs_cnc_data_collection_flow_supervisor, docs_cnc_data_collection_flow_collector, docs_cnc_data_collection_flow_worker, docs_cnc_data_collection_flow_driver, docs_cnc_data_collection_flow_cnc_communication_dll [EXTRACTED 1.00]
- **Event-driven telemetry pipeline (CNCConnect -> worker -> callback -> mutex copy -> render)** — docs_cnc_communication_cncconnect, docs_cnc_communication_flow_step1_register_callbacks, docs_cnc_communication_flow_step2_worker_thread, docs_cnc_communication_flow_step3_poll_and_invoke, docs_cnc_communication_flow_step4_copy_and_signal, docs_cnc_communication_flow_step5_render_dashboard, docs_cnc_communication_ondatacallback, docs_cnc_communication_mutex_copy, docs_cnc_communication_cncmastersnapshot [EXTRACTED 1.00]
- **Siemens OPC UA connection requirements (GOARCH=386, credentials, trusted certs)** — docs_cnc_communication_prereq_goarch_386, docs_cnc_communication_prereq_siemens_credentials, docs_cnc_communication_siemens_certificates, docs_cnc_communication_siemens_trust, docs_cnc_communication_client_cert_der, docs_cnc_communication_client_key_der, docs_cnc_communication_siemens_sinumerik_profile [EXTRACTED 1.00]
- **CNC telemetry callback data flow** — docs_cnccommunication_dll_api_integration_guide_go_client_cncconnect, docs_cnccommunication_dll_api_integration_guide_go_client_stdcall_callbacks, docs_cnccommunication_dll_api_integration_guide_go_client_ondatacallback, docs_cnccommunication_dll_api_integration_guide_go_client_cncmastersnapshot, docs_cnccommunication_dll_api_integration_guide_go_client_event_driven_data_flow [EXTRACTED 1.00]
- **Company-scoped customer entities sharing companyOwner access** — _claude_claude_company_scope, _claude_claude_parts_api, _claude_claude_clienteles_api, _claude_claude_shifts_breaks_api, _claude_claude_holidays_api, _claude_claude_reasons_api [EXTRACTED 1.00]
- **Struct ABI compatibility (pack(8) x86, size 34440, GOARCH=386)** — docs_cnc_communication_cncmastersnapshot, docs_cnc_communication_struct_size_assertion, docs_cnc_communication_prereq_goarch_386, docs_cnc_communication_ts_memory_panic [INFERRED 0.85]
- **Mode decided, logged, served and charted** — docs_cnc_data_collection_flow_mode_decision, docs_cnc_data_collection_flow_machine_mode_log, _claude_claude_health_mode_history, web_mode_chart_index_page [INFERRED 0.85]
- **Session authentication and authorization flow** — readme_session_cookie, readme_auth_tokens_hashed, readme_per_request_session_check, readme_login_lockout, readme_csrf_cross_origin, _claude_claude_role_guards [INFERRED 0.85]

## Communities (90 total, 21 thin omitted)

### Community 1 - "MachineTelemetryService"
Cohesion: 0.07
Nodes (19): 4. Find the machines: steps 5 – 9, Step 5 — Prepare the database, Step 7 — Listen for database changes, Step 8 — Load all active machines, Value, ListenForDatabaseNotifications(), LiveValue, MachineEvent (+11 more)

### Community 2 - "UserSupervisorService"
Cohesion: 0.15
Nodes (10): CreateUserSupervisorRequest, UpdateUserSupervisorRequest, validateUserSupervisor(), NewUserSupervisorHandler(), UserSupervisor, UserSupervisorRepository, mapUserSupervisorError(), UserSupervisorService (+2 more)

### Community 3 - "WithTransaction"
Cohesion: 0.07
Nodes (25): WithTransaction(), uniquePositiveIDs(), clockMinutes(), CreateShiftRequest, ShiftFields, UpdateShiftRequest, MinutesBetween(), validWorkingDays() (+17 more)

### Community 4 - "SendResponse()"
Cohesion: 0.08
Nodes (14): BreakHandler, ClienteleHandler, dashboardQuery(), DashboardHandler, queryIDList(), serveDashboard(), DowntimeReasonHandler, HolidayHandler (+6 more)

### Community 5 - "pathID()"
Cohesion: 0.09
Nodes (12): MachineCycleHandler, PartHandler, pathID(), targetUserID(), UserCellHandler, UserCityHandler, UserLineHandler, UserLocationHandler (+4 more)

### Community 6 - "net/http.ResponseWriter"
Cohesion: 0.08
Nodes (12): BrandHandler, CategoryHandler, ControllerHandler, HealthHandler, LicenceKeyHandler, MachineHandler, MachineModelHandler, ProductHandler (+4 more)

### Community 7 - "UserCityService"
Cohesion: 0.16
Nodes (8): CreateUserCityRequest, UpdateUserCityRequest, validateUserCity(), UserCity, UserCityRepository, mapUserCityError(), UserCityService, toUserCityResponseError()

### Community 8 - "api_helpers_test.go"
Cohesion: 0.29
Nodes (6): cncMasterSnapshot, LanguageResponse, ThemeResponse, authUserKey, body, machineModelRecord

### Community 9 - "CncMasterSnapshot (34,440 bytes, #pragma"
Cohesion: 0.10
Nodes (43): CncMasterSnapshot (34,440 bytes, #pragma pack(8), x86), AbsPos ([5]float64, Inches / Degrees, Work Coordinate System (WCS) positions), ActiveGCode ([16]byte, ASCII, Currently active G-code group), ActiveProgPath ([256]byte, ASCII Path, Storage path of active NC artifact), AlarmMsg ([128]byte, UTF-8 / ASCII, Alarm description message), AlarmNo (int32, Code (0 = OK), Primary alarm or fault number), AxisLoadXPercentage (int16, Percent (%), Servo X axis drive load (AxisLoadX/Y/ZPercentage)), AxisLoadYPercentage (int16, Percent (%), Servo Y axis drive load (AxisLoadX/Y/ZPercentage)) (+35 more)

### Community 10 - "context.Context"
Cohesion: 0.06
Nodes (20): Brand, Controller, ModeChange, UserProfile, BrandRepository, mapBrandError(), ControllerRepository, mapControllerError() (+12 more)

### Community 11 - "Connection"
Cohesion: 0.11
Nodes (15): dllDriver, slot, Step 12 — Call the DLL's CNCConnect, Step 13 — The DLL sends data, Step 17 — Disconnect and try again, Connection, Reading, Stats (+7 more)

### Community 12 - "UserLineService"
Cohesion: 0.07
Nodes (19): CreateUserCellRequest, UpdateUserCellRequest, validateUserCell(), CreateUserLineRequest, UpdateUserLineRequest, validateUserLine(), NewUserCellHandler(), NewUserLineHandler() (+11 more)

### Community 13 - "Collector"
Cohesion: 0.07
Nodes (36): Supervisor, 2. All steps at a glance, 3. Startup: steps 1 – 4, 7. Shutdown: steps 19 – 20, Log lines (in `logs/app.log`, search for the message), Retry and timing summary, Step 11 — Connect, a few machines at a time, Step 18 — React to machine and variable changes (+28 more)

### Community 14 - "worker"
Cohesion: 0.13
Nodes (12): 5. Connect and collect, per machine: steps 10 – 17, Known issue, Step 10 — Get ready to connect, Step 14 — The data collection loop, Step 15 — Store the reading, Step 16 — Every second: dead connection check and live status, FieldIndex(), TestFieldIndexAxisPositions() (+4 more)

### Community 15 - "AuthUser"
Cohesion: 0.11
Nodes (16): CreateMachineRequest, MachineFields, UpdateMachineRequest, validCoordinate(), validHealthIDs(), validIPv4(), SeesAllCompanyMachines(), Machine (+8 more)

### Community 16 - "newMachineService"
Cohesion: 0.16
Nodes (25): NewBranchHandler(), NewUserCityHandler(), NewUserCountryHandler(), NewUserHandler(), NewUserLocationHandler(), NewUserStateHandler(), NewBranchRepository(), NewUserCityRepository() (+17 more)

### Community 17 - "testing.T"
Cohesion: 0.12
Nodes (62): CreateUserRequest, TestParseDayStart(), PartDrawing, decodeData(), expectStatus(), newCatalogRouter(), runCatalogCRUD(), sendRequest() (+54 more)

### Community 18 - "time.Time"
Cohesion: 0.08
Nodes (34): CycleCountRow, DashboardRequest, FinishCyclesResponse, MachineTimeRow, ModeTimelineResponse, OEEResponse, OEERow, MachineModeHistoryResponse (+26 more)

### Community 19 - "How a cycle ends"
Cohesion: 0.14
Nodes (13): 1. The rows are consistent (SQL, no machine needed), 2. The rows agree with the mode log (SQL), 3. The raw signals around one cycle (telemetry and logs), 4. At the machine (floor test), Completed, How a cycle ends, How to verify the cycles, Machine cycles (+5 more)

### Community 20 - "net/http.ServeMux"
Cohesion: 0.16
Nodes (32): RequireRole(), registerBranchRoutes(), registerBrandRoutes(), registerBreakRoutes(), registerCategoryRoutes(), registerClienteleRoutes(), registerControllerRoutes(), registerDashboardRoutes() (+24 more)

### Community 21 - "CNCConnect(ip, port, dataCb, errCb, machineType) int32"
Cohesion: 0.20
Nodes (11): CNCConnect(ip, port, dataCb, errCb, machineType) int32, CNCConnect param dataCb (uintptr stdcall OnData(const CncMasterSnapshot*)), CNCConnect param ip (const char*, ANSI NUL-terminated), CNCConnect param port (int32, e.g. 4840), CNCdisconnect(machineId int32), LazyDLL binding (syscall.NewLazyDLL + eager mod.Load in init), §4 API Reference, §5 DLL Loading & Binding (Go) (+3 more)

### Community 22 - "CompanyProfile"
Cohesion: 0.28
Nodes (5): CompanyProfile, CompanyProfileRepository, mapCompanyProfileError(), toCompanyProfileResponseError(), CompanyProfileService

### Community 23 - "§10 Diagnostics & Troubleshooting"
Cohesion: 0.12
Nodes (25): Check mode: go run sinumerik_client.go check (OPC UA node existence report), Go Client (sinumerik_client.go, callback -> copy -> render), Prerequisite: 32-bit x86 DLL, compile Go with GOARCH=386, Prerequisite: SIEMENS_USER / SIEMENS_PASS env vars before CNCConnect, Run: go run sinumerik_client.go (live dashboard), §3.1 Environment & Prerequisites, §8 Running the Go Client, §10 Diagnostics & Troubleshooting (+17 more)

### Community 24 - "dto_test.go"
Cohesion: 0.20
Nodes (9): TestCreateBrandRequestValidate(), TestCreateCategoryRequestValidate(), TestCreateControllerRequestValidate(), TestCreateProductRequestValidate(), TestUpdateBrandRequestValidate(), TestUpdateCategoryRequestValidate(), TestUpdateControllerRequestValidate(), TestUpdateProductRequestValidate() (+1 more)

### Community 25 - "machine_mode.go"
Cohesion: 0.15
Nodes (20): engineName(), cycleReadingFrom(), decideMachineMode(), holdReason(), machineStateFromReading(), numberingEngine(), programStateFromProcStatus(), ncReading() (+12 more)

### Community 26 - "BranchService"
Cohesion: 0.15
Nodes (9): CreateBranchRequest, UpdateBranchRequest, validateBranchFields(), Branch, BranchRepository, mapBranchError(), cityOwner(), BranchService (+1 more)

### Community 27 - "UserService"
Cohesion: 0.10
Nodes (21): CustomerDetails, UpdateUserRequest, UserListFilter, UserResponse, ParseDayStart(), Address, AssignableRoles(), CanAssignRole() (+13 more)

### Community 28 - "CNC Data Collection Flow"
Cohesion: 0.15
Nodes (12): 1. The big picture, 6. When data is edited: step 18, 8. Reference, 9. Data Reading Flow (Function Names Only), docs/cnc-communication.md, CNC Data Collection Flow, Events in `machine_events`, Main parts (+4 more)

### Community 29 - "toUserLocationResponseError"
Cohesion: 0.16
Nodes (8): CreateUserLocationRequest, UpdateUserLocationRequest, validateUserLocation(), UserLocation, UserLocationRepository, mapUserLocationError(), UserLocationService, toUserLocationResponseError()

### Community 30 - ".Create"
Cohesion: 0.29
Nodes (3): CreateUserStateRequest, validateUserState(), toUserStateResponseError()

### Community 31 - "AuthUser"
Cohesion: 0.06
Nodes (34): ValidPartDrawingFile(), NewIdleReasonHandler(), AuthUser, Clientele, DowntimeReason, IdleReason, Part, PartProcess (+26 more)

### Community 32 - "Worker (internal/monitor/worker.go)"
Cohesion: 0.05
Nodes (34): CNC Data Collection (in-process), writeLog CNC Log Types and Switches, .env Configuration (internal/config), Health Mode History API (public), 7th Sense Server Claude Context, response.SendResponse Format, User CRUD Golden Reference, CncCommunication.dll (32-bit vendor DLL) (+26 more)

### Community 33 - "Querier"
Cohesion: 0.10
Nodes (32): Querier, RoleHandler, NewRoleHandler(), NewAddressRepository(), NewCompanyProfileRepository(), NewDowntimeReasonRepository(), NewIdleReasonRepository(), RoleRepository (+24 more)

### Community 34 - "MachineModelService"
Cohesion: 0.20
Nodes (6): MachineModelRequest, MachineModel, MachineModelRepository, mapMachineModelError(), MachineModelService, toMachineModelResponseError()

### Community 35 - "NewError"
Cohesion: 0.14
Nodes (10): CreateAddressRequest, UpdateAddressRequest, UpdateLanguageRequest, UpdateLicenceKeyRequest, validDuration(), ValidLicenceKey(), NewError(), newResponse() (+2 more)

### Community 36 - "LicenceKeyService"
Cohesion: 0.13
Nodes (11): NewLicenceKey, CreateLicenceKeysRequest, GenerateLicenceKeysRequest, GeneratedLicenceKeys, LicenceKey, LicenceKeyRepository, mapLicenceKeyError(), LicenceKeyService (+3 more)

### Community 37 - "Client Certificate & Key (client_cert.der / client_key.der)"
Cohesion: 0.48
Nodes (7): client_cert.der/.pem (X.509 client cert for Siemens OPC UA), client_key.der/.pem (client private key), Siemens certs: client_cert.der + client_key.der in binary directory, Siemens trust: cert must be trusted in NCU Access Administration, Troubleshoot: BadRequestNotAllowed -> Export client_cert, BadRequestNotAllowed (certificate trust rejection), Client Certificate & Key (client_cert.der / client_key.der)

### Community 38 - ".validate"
Cohesion: 0.47
Nodes (3): CreateDowntimeReasonRequest, DowntimeReasonFields, UpdateDowntimeReasonRequest

### Community 39 - "mode-chart.js"
Cohesion: 0.10
Nodes (64): card(), cell(), chartData(), chartRows(), countByGroup(), currentSpan(), currentWindow(), cycleData() (+56 more)

### Community 40 - "CncCommunication.dll API & Integration G"
Cohesion: 0.19
Nodes (20): CncCommunication.dll (C++ core runtime, x86, OPC UA inside), CncCommunication.dll API & Integration Guide (Go Client) - text version, Fanuc profile (engine 1, port 8193, Fwlib32.dll + fwlibe1.dll), Fwlib32.dll (Fanuc FOCAS2), fwlibe1.dll (Fanuc FOCAS2), libcrypto-3.dll (OpenSSL), libssl-3.dll (OpenSSL), Engine 1: Fanuc (FOCAS2 Ethernet, fwlib32) (+12 more)

### Community 41 - "sinumerik_client.go"
Cohesion: 0.20
Nodes (17): axisColor(), cstr(), enableVT(), findLibDir(), initDLL(), isRotary(), lbl(), logf() (+9 more)

### Community 42 - "run"
Cohesion: 0.08
Nodes (50): run(), Connect(), NewBrandHandler(), NewBreakHandler(), NewCategoryHandler(), NewControllerHandler(), NewDashboardHandler(), NewDowntimeReasonHandler() (+42 more)

### Community 43 - "CncCommunication.dll API & Integration G"
Cohesion: 0.13
Nodes (18): CNCConnect param machineType (int32 engine id, see §4.2), Engine 0: Simulator (local loopback mock engine), Engine 2: Siemens (Sinumerik OPC UA, active profile), Engine 3: Syntec (Open API / TCP), Engine 4: Mitsubishi (driver not specified), Machine identifier mapping (engine id enum 0-4), §9 Controller-Specific Profiles, Siemens default ports: 4840 (native NCU) / 4841 (custom) (+10 more)

### Community 44 - ".validate"
Cohesion: 0.47
Nodes (3): CreateIdleReasonRequest, IdleReasonFields, UpdateIdleReasonRequest

### Community 45 - ".Login()"
Cohesion: 0.18
Nodes (7): LoginRequest, newSessionToken(), truncate(), newLoginThrottle(), loginFailures, loginThrottle, Session

### Community 46 - ".validate"
Cohesion: 0.50
Nodes (3): isTimeZone(), isWebsite(), normalizeFilters()

### Community 48 - "Machines API (/machines)"
Cohesion: 0.16
Nodes (18): Clienteles API, Company Scope (companyOwner / ownCompanyCustomer), Current User API (/me theme, language), Customer Display Names (customerName / companyName), Holidays API, Licence Keys API, Machine Catalog Brand-Controller-Model, Machines API (/machines) (+10 more)

### Community 50 - "MachineCycle"
Cohesion: 0.15
Nodes (9): MachineCyclesResponse, NewMachineCycleHandler(), CycleJob, MachineCycle, MachineCycleRepository, mapMachineCycleError(), MachineCycleService, sameID() (+1 more)

### Community 51 - "UserShopService"
Cohesion: 0.15
Nodes (9): CreateUserShopRequest, UpdateUserShopRequest, validateUserShop(), NewUserShopHandler(), UserShop, UserShopRepository, mapUserShopError(), UserShopService (+1 more)

### Community 52 - "Execute"
Cohesion: 0.05
Nodes (28): Execute(), CreateUserCountryRequest, Role, UserCountry, User, UserState, mapAuthError(), mapLocationError() (+20 more)

### Community 54 - "AuthHandler"
Cohesion: 0.38
Nodes (3): AuthHandler, CurrentUser(), UserID()

### Community 56 - "net/http.ResponseWriter"
Cohesion: 0.15
Nodes (5): validator, BranchHandler, currentUserID(), decodeAndValidate(), UserHandler

### Community 57 - "Category"
Cohesion: 0.06
Nodes (22): CreateBrandRequest, UpdateBrandRequest, validateCatalogFields(), CreateCategoryRequest, UpdateCategoryRequest, CreateControllerRequest, UpdateControllerRequest, validateControllerFields() (+14 more)

### Community 58 - ".validate"
Cohesion: 0.47
Nodes (3): BreakFields, CreateBreakRequest, UpdateBreakRequest

### Community 60 - "config.go"
Cohesion: 0.19
Nodes (11): main(), run(), Config, getEnv(), getEnvBool(), getEnvInt(), getEnvMachineSelection(), DBConfig (+3 more)

### Community 61 - "Test list"
Cohesion: 0.08
Nodes (26): Before the tests (once), Check query, Floor tests on the HUMA panel (Fanuc), HUMA buttons used, Result sheet, TEST-01 · Start a cycle, TEST-02 · Normal end (M30), TEST-03 · Three parts one after another (+18 more)

### Community 62 - "validText()"
Cohesion: 0.21
Nodes (7): PartProcessFields, PartSubProgramFields, CreatePartRequest, PartFields, UpdatePartRequest, UploadPartDrawingRequest, validText()

### Community 63 - "net/http.Handler"
Cohesion: 0.15
Nodes (20): NewClienteleHandler(), NewLicenceKeyHandler(), NewPartHandler(), NewProductionScheduleHandler(), CrossOrigin(), Logger(), Chain(), Middleware (+12 more)

### Community 64 - "cycleTracker"
Cohesion: 0.15
Nodes (13): CycleStatusFor(), decideCycleStep(), isProgramLoop(), resetReason(), readingAt(), TestAwaitNextCountsFromItsStart(), TestDecideCycleStepSubProgramCallIsNotALoop(), TestIsProgramLoop() (+5 more)

### Community 65 - "ProductionScheduleRepository"
Cohesion: 0.09
Nodes (19): validDate(), CreateProductionScheduleRequest, ProductionScheduleFields, ProductionScheduleItemFields, UpdateProductionScheduleRequest, ScheduleDay(), ProductionSchedule, ProductionScheduleItem (+11 more)

### Community 67 - ".Validate"
Cohesion: 0.21
Nodes (11): UpdateThemeRequest, normalizeTheme(), TestCreateUserRequestValidate(), TestCreateUserRequestValidateDropsDuplicateMachines(), TestCreateUserRequestValidateNormalizes(), TestCustomerFiltersValidate(), TestUpdateThemeRequestValidate(), TestUpdateUserRequestValidateStatus() (+3 more)

### Community 68 - "demo_monitor/main.go"
Cohesion: 0.15
Nodes (9): CncMasterSnapshot, MachineDef, cstr(), enableVT(), initDLL(), main(), render(), TestSnapshotModeFieldOffsets() (+1 more)

### Community 69 - "AuthService"
Cohesion: 0.22
Nodes (13): NewAuthHandler(), Authenticate(), WithAuthUser(), AuthRepository, NewAuthRepository(), AuthService, NewAuthService(), asTestUser() (+5 more)

### Community 71 - "onDataCallback(ptr uintptr) uintptr"
Cohesion: 0.36
Nodes (9): Package-level callback singletons (syscall.NewCallback), CNCConnect param errCb (uintptr stdcall OnErr(int32 code, const char* msg)), Mutex-protected copy: dataMutex sync.RWMutex -> lastRecord/hasRecord, onDataCallback(ptr uintptr) uintptr, onErrorCallback(code, messagePtr uintptr) uintptr, §7 Asynchronous Callback Architecture, onDataCallback, onErrorCallback (+1 more)

### Community 72 - "HolidayService"
Cohesion: 0.17
Nodes (9): CreateHolidayRequest, HolidayFields, UpdateHolidayRequest, Holiday, HolidayRepository, mapHolidayError(), HolidayService, holidayFromFields() (+1 more)

### Community 75 - "logger/logger.go"
Cohesion: 0.40
Nodes (3): main(), Init(), InitFile()

### Community 78 - ".validate"
Cohesion: 0.38
Nodes (4): ClienteleFields, CreateClienteleRequest, UpdateClienteleRequest, validClienteleEmail()

### Community 80 - "Event-driven data flow (no polling in client)"
Cohesion: 0.60
Nodes (5): Flow 1: load DLL, register two stdcall callbacks (syscall.NewCallback), Flow 2: CNCConnect starts DLL background worker thread, Flow 3: worker polls controller, assembles snapshot, invokes data callback, Flow 4: callback copies payload, signals renderer via buffered channel, Flow 5: main goroutine renders live ANSI dashboard

### Community 86 - "auth_tokens SHA-256 Hashed Sessions (mig"
Cohesion: 0.67
Nodes (3): auth_tokens SHA-256 Hashed Sessions (migration 014), Per-Request Session and User Re-check, ts_session Cookie (HttpOnly, SameSite=Strict)

## Ambiguous Edges - Review These
- `CncCommunication.dll (32-bit vendor DLL)` → `FanucCncCommunication.dll`  [AMBIGUOUS]
  docs/setup-guide.txt · relation: semantically_similar_to
- `Machines API (/machines)` → `Snapshot Fields (positions, spindle, feed, program, tool, alarms)`  [AMBIGUOUS]
  docs/CncCommunication.dll — API & Integration Guide (Go Client).pdf · relation: conceptually_related_to

## Knowledge Gaps
- **36 isolated node(s):** `hello`, `7th-sense-server`, `middleware.RequireRole / routes.adminOnly`, `response.SendResponse`, `Test Auth (asTestUser, tests/auth_test.go)` (+31 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 130 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **21 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `CNC Data Collection Flow` connect `CNC Data Collection Flow` to `Worker (internal/monitor/worker.go)`, `MachineTelemetryService`, `Collector`, `worker`?**
  _High betweenness centrality (0.215) - this node is a cross-community bridge._
- **What connects `hello`, `7th-sense-server`, `middleware.RequireRole / routes.adminOnly` to the rest of the system?**
  _36 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `MachineTelemetryService` be split into smaller, more focused modules?**
  _Cohesion score 0.06970128022759602 - nodes in this community are weakly interconnected._
- **Why does `2. All steps at a glance` connect `Collector` to `Connection`, `CNC Data Collection Flow`, `worker`?**
  _High betweenness centrality (0.189) - this node is a cross-community bridge._
- **Should `WithTransaction` be split into smaller, more focused modules?**
  _Cohesion score 0.07467532467532467 - nodes in this community are weakly interconnected._

### Low-confidence Hints
_AMBIGUOUS edges — the extractor was unsure. Verify before acting on these._

- **What is the exact relationship between `CncCommunication.dll (32-bit vendor DLL)` and `FanucCncCommunication.dll`?**
  _Edge tagged AMBIGUOUS (relation: semantically_similar_to) - confidence is low._
- **What is the exact relationship between `Machines API (/machines)` and `Snapshot Fields (positions, spindle, feed, program, tool, alarms)`?**
  _Edge tagged AMBIGUOUS (relation: conceptually_related_to) - confidence is low._