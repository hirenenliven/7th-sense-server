# Graph Report - 7th-sense-server  (2026-10-30)

## Corpus Check
- 259 files · ~141,215 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 4 file(s) not represented in the graph (top: (none) 3, .toml 1)

## Summary
- 2169 nodes · 17718 edges · 95 communities (74 shown, 21 thin omitted)
- Extraction: 95% EXTRACTED · 5% INFERRED · 0% AMBIGUOUS · INFERRED: 858 edges (avg confidence: 0.86)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `068a3317`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- LicenceKeyDuration
- Shift
- executeAffected
- MachineCycle
- SendResponse()
- pathID()
- SendResponse
- UserCityService
- api_helpers_test.go
- CncMasterSnapshot (34,440 bytes, #pragma
- ControllerService
- dll_windows_386.go
- UserLineService
- Collector
- worker
- MachineService
- demo_monitor/main.go
- testing.T
- MachineTelemetryService
- time.Time
- net/http.ServeMux
- HolidayService
- UserSupervisorService
- §10 Diagnostics & Troubleshooting
- UserShopService
- machine_mode.go
- BranchService
- config.go
- context.Context
- UserLocationService
- Execute
- Handlers
- Worker (internal/monitor/worker.go)
- Querier
- MachineModelService
- NewError()
- LicenceKeyService
- github.com/jackc/pgx/v5/pgxpool.Pool
- AuthUser
- mode-chart.js
- CncCommunication.dll API & Integration G
- sinumerik_client.go
- UserProfileRepository
- CncCommunication.dll API & Integration G
- validateCatalogFields()
- .Login()
- Product
- AuthService
- Machines API (/machines)
- 7th-sense-server/internal/model.AuthUser
- newUserHandler
- DowntimeReasonService
- UserCountryService
- New
- AuthHandler
- CNC Data Collection Flow
- net/http.ResponseWriter
- ProductionScheduleRepository
- .validate
- .Validate
- .UploadDrawing
- .Validate
- validText()
- run
- VariableConfig
- generate.go
- .validate()
- logger/logger.go
- DowntimeReasonFields
- CNCConnect(ip, port, dataCb, errCb, mach
- Response
- onDataCallback(ptr uintptr) uintptr
- Client Certificate & Key (client_cert.de
- .validate()
- Event-driven data flow (no polling in cl
- .validate()
- newPartRouter
- .Generate()
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
1. `Execute()` - 179 edges
2. `SendResponse()` - 146 edges
3. `NewError()` - 131 edges
4. `run()` - 93 edges
5. `pathID()` - 80 edges
6. `AuthUser` - 71 edges
7. `Querier` - 60 edges
8. `currentUser()` - 59 edges
9. `decodeAndValidate()` - 58 edges
10. `worker` - 54 edges

## Surprising Connections (you probably didn't know these)
- `Machine Identifier Mapping (0 Simulator,1 Fanuc,2 Siemens,3 Syntec,4 Mitsubishi)` --semantically_similar_to--> `Machine Catalog Brand-Controller-Model`  [INFERRED] [semantically similar]
  docs/CncCommunication.dll — API & Integration Guide (Go Client).pdf → .claude/CLAUDE.md
- `Snapshot Fields (positions, spindle, feed, program, tool, alarms)` --conceptually_related_to--> `Machines API (/machines)`  [AMBIGUOUS]
  docs/CncCommunication.dll — API & Integration Guide (Go Client).pdf → .claude/CLAUDE.md
- `2. All steps at a glance` --references--> `LoadDriver()`  [INFERRED]
  docs/cnc-data-collection-flow.md → internal/cnc/cnc.go
- `Step 3 — Run the collector in the background` --references--> `Collector`  [INFERRED]
  docs/cnc-data-collection-flow.md → internal/monitor/collector.go
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

## Communities (95 total, 21 thin omitted)

### Community 1 - "Shift"
Cohesion: 0.15
Nodes (19): clockMinutes(), CreateShiftRequest, ShiftFields, UpdateShiftRequest, MinutesBetween(), Shift, NewBreakRepository(), ShiftRepository (+11 more)

### Community 2 - "executeAffected"
Cohesion: 0.20
Nodes (5): PartProcess, PartSubProgram, executeAffected(), PartRepository, mapPartError()

### Community 3 - "MachineCycle"
Cohesion: 0.12
Nodes (12): MachineCycleListRequest, MachineCyclesResponse, MachineCycleSummary, NewMachineCycleHandler(), CycleJob, MachineCycle, MachineCycleRepository, mapMachineCycleError() (+4 more)

### Community 4 - "SendResponse()"
Cohesion: 0.09
Nodes (11): BreakHandler, ClienteleHandler, DowntimeReasonHandler, HolidayHandler, IdleReasonHandler, PartHandler, ProductionScheduleHandler, currentUser() (+3 more)

### Community 5 - "pathID()"
Cohesion: 0.10
Nodes (12): validator, ControllerHandler, decodeAndValidate(), pathID(), targetUserID(), UserCellHandler, UserCityHandler, UserLineHandler (+4 more)

### Community 6 - "SendResponse"
Cohesion: 0.11
Nodes (6): BrandHandler, CategoryHandler, LicenceKeyHandler, MachineModelHandler, UserCountryHandler, SendResponse()

### Community 7 - "UserCityService"
Cohesion: 0.08
Nodes (18): CreateUserCityRequest, UpdateUserCityRequest, validateUserCity(), CreateUserStateRequest, validateUserState(), UserCity, UserState, UserCityRepository (+10 more)

### Community 8 - "api_helpers_test.go"
Cohesion: 0.29
Nodes (6): cncMasterSnapshot, LanguageResponse, ThemeResponse, authUserKey, body, machineModelRecord

### Community 9 - "CncMasterSnapshot (34,440 bytes, #pragma"
Cohesion: 0.10
Nodes (43): CncMasterSnapshot (34,440 bytes, #pragma pack(8), x86), AbsPos ([5]float64, Inches / Degrees, Work Coordinate System (WCS) positions), ActiveGCode ([16]byte, ASCII, Currently active G-code group), ActiveProgPath ([256]byte, ASCII Path, Storage path of active NC artifact), AlarmMsg ([128]byte, UTF-8 / ASCII, Alarm description message), AlarmNo (int32, Code (0 = OK), Primary alarm or fault number), AxisLoadXPercentage (int16, Percent (%), Servo X axis drive load (AxisLoadX/Y/ZPercentage)), AxisLoadYPercentage (int16, Percent (%), Servo Y axis drive load (AxisLoadX/Y/ZPercentage)) (+35 more)

### Community 10 - "ControllerService"
Cohesion: 0.10
Nodes (13): NewBrandHandler(), Brand, Controller, BrandRepository, mapBrandError(), ControllerRepository, mapControllerError(), NewControllerRepository() (+5 more)

### Community 11 - "dll_windows_386.go"
Cohesion: 0.10
Nodes (16): dllDriver, slot, Step 12 — Call the DLL's CNCConnect, Step 13 — The DLL sends data, Connection, Reading, Status, newConnection() (+8 more)

### Community 12 - "UserLineService"
Cohesion: 0.14
Nodes (11): CreateUserLineRequest, UpdateUserLineRequest, validateUserLine(), NewUserLineHandler(), UserLine, UserLineRepository, mapUserLineError(), NewUserLineRepository() (+3 more)

### Community 13 - "Collector"
Cohesion: 0.11
Nodes (21): 4. Find the machines: steps 5 – 9, Log lines (in `logs/app.log`, search for the message), Step 18 — React to machine and variable changes, Step 20 — Stop the workers and mark machines offline, Step 5 — Prepare the database, Step 6 — Record the start and begin following the database, Step 7 — Listen for database changes, Step 8 — Load all active machines (+13 more)

### Community 14 - "worker"
Cohesion: 0.11
Nodes (15): 2. All steps at a glance, 5. Connect and collect, per machine: steps 10 – 17, Step 10 — Get ready to connect, Step 11 — Connect, a few machines at a time, Step 14 — The data collection loop, Step 16 — Every second: dead connection check and live status, Step 17 — Disconnect and try again, Stats (+7 more)

### Community 15 - "MachineService"
Cohesion: 0.11
Nodes (19): CreateMachineRequest, MachineFields, UpdateMachineRequest, validCoordinate(), validHealthIDs(), validIPv4(), NewHealthHandler(), Machine (+11 more)

### Community 16 - "demo_monitor/main.go"
Cohesion: 0.22
Nodes (7): CncMasterSnapshot, MachineDef, cstr(), enableVT(), initDLL(), main(), render()

### Community 17 - "testing.T"
Cohesion: 0.06
Nodes (91): TestSnapshotModeFieldOffsets(), TestSnapshotSize(), WithAuthUser(), CrossOrigin(), Logger(), Chain(), Recoverer(), NewAddressService() (+83 more)

### Community 18 - "MachineTelemetryService"
Cohesion: 0.08
Nodes (15): MachineModeHistoryResponse, MachineModesResponse, MachineModeTotal, LiveValue, MachineEvent, MachineLiveStatus, MachineModePeriod, MachineModeStatus (+7 more)

### Community 19 - "time.Time"
Cohesion: 0.12
Nodes (15): Known issue, Step 15 — Store the reading, Value, CycleStatusFor(), decideCycleStep(), isProgramLoop(), resetReason(), secondsOf() (+7 more)

### Community 20 - "net/http.ServeMux"
Cohesion: 0.16
Nodes (31): RequireRole(), registerBranchRoutes(), registerBrandRoutes(), registerBreakRoutes(), registerCategoryRoutes(), registerClienteleRoutes(), registerControllerRoutes(), registerDowntimeReasonRoutes() (+23 more)

### Community 21 - "HolidayService"
Cohesion: 0.21
Nodes (8): Holiday, HolidayRepository, mapHolidayError(), NewHolidayRepository(), HolidayService, holidayFromFields(), toHolidayResponseError(), optionalText()

### Community 22 - "UserSupervisorService"
Cohesion: 0.14
Nodes (11): CreateUserSupervisorRequest, UpdateUserSupervisorRequest, validateUserSupervisor(), NewUserSupervisorHandler(), UserSupervisor, UserSupervisorRepository, mapUserSupervisorError(), NewUserSupervisorRepository() (+3 more)

### Community 23 - "§10 Diagnostics & Troubleshooting"
Cohesion: 0.12
Nodes (25): Check mode: go run sinumerik_client.go check (OPC UA node existence report), Go Client (sinumerik_client.go, callback -> copy -> render), Prerequisite: 32-bit x86 DLL, compile Go with GOARCH=386, Prerequisite: SIEMENS_USER / SIEMENS_PASS env vars before CNCConnect, Run: go run sinumerik_client.go (live dashboard), §3.1 Environment & Prerequisites, §8 Running the Go Client, §10 Diagnostics & Troubleshooting (+17 more)

### Community 24 - "UserShopService"
Cohesion: 0.14
Nodes (11): CreateUserShopRequest, UpdateUserShopRequest, validateUserShop(), NewUserShopHandler(), UserShop, UserShopRepository, mapUserShopError(), NewUserShopRepository() (+3 more)

### Community 25 - "machine_mode.go"
Cohesion: 0.14
Nodes (22): engineName(), cycleReadingFrom(), decideMachineMode(), holdReason(), machineStateFromReading(), numberingEngine(), programStateFromProcStatus(), ncReading() (+14 more)

### Community 26 - "BranchService"
Cohesion: 0.15
Nodes (9): CreateBranchRequest, UpdateBranchRequest, validateBranchFields(), Branch, BranchRepository, mapBranchError(), cityOwner(), BranchService (+1 more)

### Community 27 - "config.go"
Cohesion: 0.16
Nodes (12): main(), run(), Config, MachineSelection, getEnv(), getEnvBool(), getEnvInt(), getEnvMachineSelection() (+4 more)

### Community 28 - "context.Context"
Cohesion: 0.17
Nodes (5): ShiftBreak, BreakRepository, mapBreakError(), ClienteleRepository, mapClienteleError()

### Community 29 - "UserLocationService"
Cohesion: 0.14
Nodes (11): CreateUserLocationRequest, UpdateUserLocationRequest, validateUserLocation(), NewUserLocationHandler(), UserLocation, UserLocationRepository, mapUserLocationError(), NewUserLocationRepository() (+3 more)

### Community 30 - "Execute"
Cohesion: 0.08
Nodes (26): Connect(), Execute(), CreateUserCellRequest, UpdateUserCellRequest, validateUserCell(), NewUserCellHandler(), UserCell, User (+18 more)

### Community 31 - "Handlers"
Cohesion: 0.13
Nodes (5): HealthHandler, MachineCycleHandler, MachineHandler, registerHealthRoutes(), Handlers

### Community 32 - "Worker (internal/monitor/worker.go)"
Cohesion: 0.06
Nodes (27): CNC Data Collection (in-process), writeLog CNC Log Types and Switches, .env Configuration (internal/config), 7th Sense Server Claude Context, response.SendResponse Format, User CRUD Golden Reference, CncCommunication.dll (32-bit vendor DLL), Collector (internal/monitor/collector.go) (+19 more)

### Community 33 - "Querier"
Cohesion: 0.05
Nodes (46): Querier, WithTransaction(), CreateUserRequest, CustomerDetails, UpdateUserRequest, UserListFilter, UserResponse, ParseDayStart() (+38 more)

### Community 34 - "MachineModelService"
Cohesion: 0.18
Nodes (8): MachineModelRequest, MachineModel, MachineModelRepository, mapMachineModelError(), NewMachineModelRepository(), MachineModelService, NewMachineModelService(), toMachineModelResponseError()

### Community 35 - "NewError()"
Cohesion: 0.16
Nodes (16): CreateAddressRequest, UpdateAddressRequest, UpdateLanguageRequest, isTimeZone(), isWebsite(), normalizeFilters(), normalizeTheme(), TestCreateUserRequestValidate() (+8 more)

### Community 36 - "LicenceKeyService"
Cohesion: 0.16
Nodes (11): NewLicenceKey, CreateLicenceKeysRequest, NewLicenceKeyHandler(), LicenceKey, LicenceKeyRepository, mapLicenceKeyError(), NewLicenceKeyRepository(), LicenceKeyService (+3 more)

### Community 37 - "github.com/jackc/pgx/v5/pgxpool.Pool"
Cohesion: 0.21
Nodes (7): NewCategoryHandler(), Category, CategoryRepository, mapCategoryError(), CategoryService, NewCategoryService(), toCategoryResponseError()

### Community 38 - "AuthUser"
Cohesion: 0.09
Nodes (32): CreatePartRequest, PartFields, UpdatePartRequest, ValidPartDrawingFile(), CreateProductionScheduleRequest, ProductionScheduleFields, UpdateProductionScheduleRequest, AuthUser (+24 more)

### Community 39 - "mode-chart.js"
Cohesion: 0.08
Nodes (71): Health Mode History API (public), machine_mode_log, ECharts 5.5.1, Mode List Table and Time per Mode, Mode Count Tiles (filter), Machine Mode Dashboard Page, Zoom Level Bar, card() (+63 more)

### Community 40 - "CncCommunication.dll API & Integration G"
Cohesion: 0.19
Nodes (20): CncCommunication.dll (C++ core runtime, x86, OPC UA inside), CncCommunication.dll API & Integration Guide (Go Client) - text version, Fanuc profile (engine 1, port 8193, Fwlib32.dll + fwlibe1.dll), Fwlib32.dll (Fanuc FOCAS2), fwlibe1.dll (Fanuc FOCAS2), libcrypto-3.dll (OpenSSL), libssl-3.dll (OpenSSL), Engine 1: Fanuc (FOCAS2 Ethernet, fwlib32) (+12 more)

### Community 41 - "sinumerik_client.go"
Cohesion: 0.20
Nodes (17): axisColor(), cstr(), enableVT(), findLibDir(), initDLL(), isRotary(), lbl(), logf() (+9 more)

### Community 42 - "UserProfileRepository"
Cohesion: 0.19
Nodes (5): UserProfile, UserProfileRepository, mapProfileError(), toUserProfileResponseError(), UserProfileService

### Community 43 - "CncCommunication.dll API & Integration G"
Cohesion: 0.13
Nodes (18): CNCConnect param machineType (int32 engine id, see §4.2), Engine 0: Simulator (local loopback mock engine), Engine 2: Siemens (Sinumerik OPC UA, active profile), Engine 3: Syntec (Open API / TCP), Engine 4: Mitsubishi (driver not specified), Machine identifier mapping (engine id enum 0-4), §9 Controller-Specific Profiles, Siemens default ports: 4840 (native NCU) / 4841 (custom) (+10 more)

### Community 44 - "validateCatalogFields()"
Cohesion: 0.11
Nodes (11): CreateBrandRequest, UpdateBrandRequest, validateCatalogFields(), CreateCategoryRequest, UpdateCategoryRequest, CreateControllerRequest, UpdateControllerRequest, validateControllerFields() (+3 more)

### Community 45 - ".Login()"
Cohesion: 0.19
Nodes (6): LoginRequest, newSessionToken(), truncate(), newLoginThrottle(), loginFailures, loginThrottle

### Community 46 - "Product"
Cohesion: 0.20
Nodes (8): NewProductHandler(), Product, ProductRepository, mapProductError(), ProductService, NewProductService(), productImageOrDefault(), toProductResponseError()

### Community 47 - "AuthService"
Cohesion: 0.21
Nodes (9): NewAuthHandler(), Authenticate(), AuthRepository, mapAuthError(), NewAuthRepository(), AuthService, hashToken(), NewAuthService() (+1 more)

### Community 48 - "Machines API (/machines)"
Cohesion: 0.16
Nodes (18): Clienteles API, Company Scope (companyOwner / ownCompanyCustomer), Current User API (/me theme, language), Customer Display Names (customerName / companyName), Holidays API, Licence Keys API, Machine Catalog Brand-Controller-Model, Machines API (/machines) (+10 more)

### Community 49 - "7th-sense-server/internal/model.AuthUser"
Cohesion: 0.15
Nodes (12): CreateIdleReasonRequest, IdleReasonFields, UpdateIdleReasonRequest, NewIdleReasonHandler(), IdleReason, IdleReasonRepository, mapIdleReasonError(), NewIdleReasonRepository() (+4 more)

### Community 50 - "newUserHandler"
Cohesion: 0.18
Nodes (14): NewMachineHandler(), RoleHandler, NewRoleHandler(), NewUserHandler(), Role, RoleRepository, mapRoleError(), NewRoleRepository() (+6 more)

### Community 51 - "DowntimeReasonService"
Cohesion: 0.20
Nodes (9): NewDowntimeReasonHandler(), DowntimeReason, DowntimeReasonRepository, mapDowntimeReasonError(), NewDowntimeReasonRepository(), downtimeReasonFromFields(), DowntimeReasonService, NewDowntimeReasonService() (+1 more)

### Community 52 - "UserCountryService"
Cohesion: 0.15
Nodes (10): CreateUserCountryRequest, NewUserCountryHandler(), UserCountry, UserCountryRepository, mapUserCountryError(), NewUserCountryRepository(), UserCountryService, NewUserCountryService() (+2 more)

### Community 53 - "New"
Cohesion: 0.24
Nodes (13): NewBreakHandler(), NewHolidayHandler(), NewProductionScheduleHandler(), NewShiftHandler(), Middleware, NewUserRepository(), registerAuthRoutes(), New() (+5 more)

### Community 54 - "AuthHandler"
Cohesion: 0.38
Nodes (3): AuthHandler, CurrentUser(), UserID()

### Community 55 - "CNC Data Collection Flow"
Cohesion: 0.08
Nodes (25): Supervisor, 1. The big picture, 3. Startup: steps 1 – 4, 6. When data is edited: step 18, 7. Shutdown: steps 19 – 20, 8. Reference, 9. Data Reading Flow (Function Names Only), docs/cnc-communication.md (+17 more)

### Community 56 - "net/http.ResponseWriter"
Cohesion: 0.14
Nodes (5): BranchHandler, ProductHandler, currentUserID(), UserHandler, statusRecorder

### Community 57 - "ProductionScheduleRepository"
Cohesion: 0.16
Nodes (4): ProductionScheduleItem, ProductionScheduleRepository, mapProductionScheduleError(), saveScheduleItems()

### Community 58 - ".validate"
Cohesion: 0.47
Nodes (3): BreakFields, CreateBreakRequest, UpdateBreakRequest

### Community 59 - ".Validate"
Cohesion: 0.40
Nodes (3): UpdateLicenceKeyRequest, validDuration(), ValidLicenceKey()

### Community 62 - "validText()"
Cohesion: 0.20
Nodes (5): PartProcessFields, PartSubProgramFields, uniquePositiveIDs(), validText(), validWorkingDays()

### Community 63 - "run"
Cohesion: 0.21
Nodes (18): run(), NewBranchHandler(), NewControllerHandler(), NewMachineModelHandler(), NewUserCityHandler(), NewUserStateHandler(), NewBranchRepository(), NewBrandRepository() (+10 more)

### Community 64 - "VariableConfig"
Cohesion: 0.20
Nodes (7): VariableConfig, formatPreviousValue(), formatValue(), newValueChangeLog(), boundVariable, valueChange, valueChangeLog

### Community 66 - ".validate()"
Cohesion: 0.24
Nodes (6): validDate(), ProductionScheduleItemFields, ScheduleDay(), checkJobsOverlap(), shiftSpan(), jobSpan

### Community 67 - "logger/logger.go"
Cohesion: 0.40
Nodes (3): main(), Init(), InitFile()

### Community 68 - "DowntimeReasonFields"
Cohesion: 0.47
Nodes (3): CreateDowntimeReasonRequest, DowntimeReasonFields, UpdateDowntimeReasonRequest

### Community 69 - "CNCConnect(ip, port, dataCb, errCb, mach"
Cohesion: 0.20
Nodes (11): CNCConnect(ip, port, dataCb, errCb, machineType) int32, CNCConnect param dataCb (uintptr stdcall OnData(const CncMasterSnapshot*)), CNCConnect param ip (const char*, ANSI NUL-terminated), CNCConnect param port (int32, e.g. 4840), CNCdisconnect(machineId int32), LazyDLL binding (syscall.NewLazyDLL + eager mod.Load in init), §4 API Reference, §5 DLL Loading & Binding (Go) (+3 more)

### Community 71 - "onDataCallback(ptr uintptr) uintptr"
Cohesion: 0.36
Nodes (9): Package-level callback singletons (syscall.NewCallback), CNCConnect param errCb (uintptr stdcall OnErr(int32 code, const char* msg)), Mutex-protected copy: dataMutex sync.RWMutex -> lastRecord/hasRecord, onDataCallback(ptr uintptr) uintptr, onErrorCallback(code, messagePtr uintptr) uintptr, §7 Asynchronous Callback Architecture, onDataCallback, onErrorCallback (+1 more)

### Community 77 - "Client Certificate & Key (client_cert.de"
Cohesion: 0.48
Nodes (7): client_cert.der/.pem (X.509 client cert for Siemens OPC UA), client_key.der/.pem (client private key), Siemens certs: client_cert.der + client_key.der in binary directory, Siemens trust: cert must be trusted in NCU Access Administration, Troubleshoot: BadRequestNotAllowed -> Export client_cert, BadRequestNotAllowed (certificate trust rejection), Client Certificate & Key (client_cert.der / client_key.der)

### Community 78 - ".validate()"
Cohesion: 0.38
Nodes (4): ClienteleFields, CreateClienteleRequest, UpdateClienteleRequest, validClienteleEmail()

### Community 80 - "Event-driven data flow (no polling in cl"
Cohesion: 0.60
Nodes (5): Flow 1: load DLL, register two stdcall callbacks (syscall.NewCallback), Flow 2: CNCConnect starts DLL background worker thread, Flow 3: worker polls controller, assembles snapshot, invokes data callback, Flow 4: callback copies payload, signals renderer via buffered channel, Flow 5: main goroutine renders live ANSI dashboard

### Community 82 - ".validate()"
Cohesion: 0.47
Nodes (3): CreateHolidayRequest, HolidayFields, UpdateHolidayRequest

### Community 83 - "newPartRouter"
Cohesion: 0.18
Nodes (8): NewClienteleHandler(), NewPartHandler(), NewClienteleRepository(), PartDrawingStorage, NewPartDrawingStorage(), NewClienteleService(), NewPartService(), newPartRouter()

### Community 84 - ".Generate()"
Cohesion: 0.40
Nodes (3): GenerateLicenceKeysRequest, GeneratedLicenceKeys, newLicenceKey()

### Community 86 - "auth_tokens SHA-256 Hashed Sessions (mig"
Cohesion: 0.67
Nodes (3): auth_tokens SHA-256 Hashed Sessions (migration 014), Per-Request Session and User Re-check, ts_session Cookie (HttpOnly, SameSite=Strict)

## Ambiguous Edges - Review These
- `CncCommunication.dll (32-bit vendor DLL)` → `FanucCncCommunication.dll`  [AMBIGUOUS]
  docs/setup-guide.txt · relation: semantically_similar_to
- `Snapshot Fields (positions, spindle, feed, program, tool, alarms)` → `Machines API (/machines)`  [AMBIGUOUS]
  docs/CncCommunication.dll — API & Integration Guide (Go Client).pdf · relation: conceptually_related_to

## Knowledge Gaps
- **36 isolated node(s):** `hello`, `7th-sense-server`, `middleware.RequireRole / routes.adminOnly`, `response.SendResponse`, `Test Auth (asTestUser, tests/auth_test.go)` (+31 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 97 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **21 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `CNC Data Collection Flow` connect `CNC Data Collection Flow` to `Worker (internal/monitor/worker.go)`, `Collector`, `worker`?**
  _High betweenness centrality (0.158) - this node is a cross-community bridge._
- **What connects `hello`, `7th-sense-server`, `middleware.RequireRole / routes.adminOnly` to the rest of the system?**
  _36 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Shift` be split into smaller, more focused modules?**
  _Cohesion score 0.14942528735632185 - nodes in this community are weakly interconnected._
- **Why does `2. All steps at a glance` connect `worker` to `time.Time`, `dll_windows_386.go`, `Collector`, `CNC Data Collection Flow`?**
  _High betweenness centrality (0.137) - this node is a cross-community bridge._
- **Should `MachineCycle` be split into smaller, more focused modules?**
  _Cohesion score 0.11932773109243698 - nodes in this community are weakly interconnected._

### Low-confidence Hints
_AMBIGUOUS edges — the extractor was unsure. Verify before acting on these._

- **What is the exact relationship between `CncCommunication.dll (32-bit vendor DLL)` and `FanucCncCommunication.dll`?**
  _Edge tagged AMBIGUOUS (relation: semantically_similar_to) - confidence is low._
- **What is the exact relationship between `Snapshot Fields (positions, spindle, feed, program, tool, alarms)` and `Machines API (/machines)`?**
  _Edge tagged AMBIGUOUS (relation: conceptually_related_to) - confidence is low._