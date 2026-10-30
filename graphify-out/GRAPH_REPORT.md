# Graph Report - 7th-sense-server  (2026-10-30)

## Corpus Check
- 260 files · ~142,153 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 4 file(s) not represented in the graph (top: (none) 3, .toml 1)

## Summary
- 2177 nodes · 17727 edges · 106 communities (84 shown, 22 thin omitted)
- Extraction: 95% EXTRACTED · 5% INFERRED · 0% AMBIGUOUS · INFERRED: 860 edges (avg confidence: 0.86)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `8294d5ff`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- LicenceKeyDuration
- UserCellService
- testing.T
- MachineCycle
- SendResponse()
- pathID()
- SendResponse
- toUserCityResponseError
- api_helpers_test.go
- CncMasterSnapshot (34,440 bytes, #pragma
- ControllerService
- dll_windows_386.go
- UserLineService
- Collector
- worker
- MachineService
- newOrganizationRouter
- expectStatus
- MachineTelemetryService
- .runConnectionLoop
- net/http.ServeMux
- HolidayRepository
- Querier
- §10 Diagnostics & Troubleshooting
- mode-chart.js
- machine_mode.go
- BranchService
- actAs
- CNC Data Collection Flow
- UserLocationService
- Execute
- config.go
- Worker (internal/monitor/worker.go)
- UserService
- MachineModelService
- .Validate
- LicenceKeyService
- BrandService
- AuthUser
- formatDuration
- CncCommunication.dll API & Integration G
- sinumerik_client.go
- MachineTelemetryRepository
- CncCommunication.dll API & Integration G
- validateCatalogFields()
- .Login()
- Product
- demo_monitor/main.go
- Machines API (/machines)
- 7th-sense-server/internal/model.AuthUser
- ProductionScheduleRepository
- DowntimeReasonService
- UserCountryService
- net/http.Handler
- AuthHandler
- Category
- net/http.Request
- context.Context
- Break
- NewError
- CNC Data Collection (in-process)
- ClienteleService
- validText()
- run
- cycleReading
- runCatalogCRUD
- generate.go
- TestAuthentication
- TestProductionScheduleCRUD
- CNCConnect(ip, port, dataCb, errCb, mach
- newUserHandler
- onDataCallback(ptr uintptr) uintptr
- drawChart
- VariableConfig
- AuthService
- ClienteleRepository
- .Create
- Client Certificate & Key (client_cert.de
- .validate()
- NewPartService
- Event-driven data flow (no polling in cl
- logger/logger.go
- .validate
- BreakHandler
- ShiftHandler
- MachineModePeriod
- auth_tokens SHA-256 Hashed Sessions (mig
- .Generate
- MachineModeStatus
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
- MachineSelection

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
- `Step 15 — Store the reading` --references--> `machineStateFromReading()`  [INFERRED]
  docs/cnc-data-collection-flow.md → internal/monitor/machine_mode.go
- `Step 15 — Store the reading` --references--> `programStateFromProcStatus()`  [INFERRED]
  docs/cnc-data-collection-flow.md → internal/monitor/machine_mode.go
- `Step 15 — Store the reading` --references--> `decideMachineMode()`  [INFERRED]
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

## Communities (106 total, 22 thin omitted)

### Community 1 - "UserCellService"
Cohesion: 0.07
Nodes (20): CreateUserCellRequest, UpdateUserCellRequest, validateUserCell(), CreateUserSupervisorRequest, UpdateUserSupervisorRequest, validateUserSupervisor(), NewUserCellHandler(), NewUserSupervisorHandler() (+12 more)

### Community 2 - "testing.T"
Cohesion: 0.13
Nodes (20): NewUserProfileRepository(), NewAddressService(), NewUserProfileService(), TestAddressCRUD(), TestCreateBrandRequestValidate(), TestCreateCategoryRequestValidate(), TestCreateControllerRequestValidate(), TestCreateProductRequestValidate() (+12 more)

### Community 3 - "MachineCycle"
Cohesion: 0.12
Nodes (13): MachineCycleListRequest, MachineCyclesResponse, MachineCycleSummary, NewMachineCycleHandler(), CycleJob, MachineCycle, MachineCycleRepository, mapMachineCycleError() (+5 more)

### Community 4 - "SendResponse()"
Cohesion: 0.09
Nodes (9): DowntimeReasonHandler, HolidayHandler, IdleReasonHandler, MachineCycleHandler, PartHandler, ProductionScheduleHandler, currentUser(), RoleHandler (+1 more)

### Community 5 - "pathID()"
Cohesion: 0.11
Nodes (12): validator, decodeAndValidate(), pathID(), targetUserID(), UserCellHandler, UserCityHandler, UserCountryHandler, UserLineHandler (+4 more)

### Community 6 - "SendResponse"
Cohesion: 0.09
Nodes (11): BrandHandler, CategoryHandler, ControllerHandler, HealthHandler, LicenceKeyHandler, MachineHandler, MachineModelHandler, ProductHandler (+3 more)

### Community 7 - "toUserCityResponseError"
Cohesion: 0.24
Nodes (4): CreateUserCityRequest, UpdateUserCityRequest, validateUserCity(), toUserCityResponseError()

### Community 8 - "api_helpers_test.go"
Cohesion: 0.29
Nodes (6): cncMasterSnapshot, LanguageResponse, ThemeResponse, authUserKey, body, machineModelRecord

### Community 9 - "CncMasterSnapshot (34,440 bytes, #pragma"
Cohesion: 0.10
Nodes (43): CncMasterSnapshot (34,440 bytes, #pragma pack(8), x86), AbsPos ([5]float64, Inches / Degrees, Work Coordinate System (WCS) positions), ActiveGCode ([16]byte, ASCII, Currently active G-code group), ActiveProgPath ([256]byte, ASCII Path, Storage path of active NC artifact), AlarmMsg ([128]byte, UTF-8 / ASCII, Alarm description message), AlarmNo (int32, Code (0 = OK), Primary alarm or fault number), AxisLoadXPercentage (int16, Percent (%), Servo X axis drive load (AxisLoadX/Y/ZPercentage)), AxisLoadYPercentage (int16, Percent (%), Servo Y axis drive load (AxisLoadX/Y/ZPercentage)) (+35 more)

### Community 10 - "ControllerService"
Cohesion: 0.21
Nodes (5): Controller, ControllerRepository, mapControllerError(), ControllerService, toControllerResponseError()

### Community 11 - "dll_windows_386.go"
Cohesion: 0.09
Nodes (19): dllDriver, slot, Step 12 — Call the DLL's CNCConnect, Step 13 — The DLL sends data, Step 17 — Disconnect and try again, Connection, Reading, Stats (+11 more)

### Community 12 - "UserLineService"
Cohesion: 0.09
Nodes (17): CreateUserLineRequest, UpdateUserLineRequest, validateUserLine(), CreateUserShopRequest, UpdateUserShopRequest, validateUserShop(), NewUserLineHandler(), NewUserShopHandler() (+9 more)

### Community 13 - "Collector"
Cohesion: 0.09
Nodes (27): Supervisor, 2. All steps at a glance, 3. Startup: steps 1 – 4, Log lines (in `logs/app.log`, search for the message), Step 18 — React to machine and variable changes, Step 1 — Start collection with the server, Step 20 — Stop the workers and mark machines offline, Step 2 — Load the DLL (+19 more)

### Community 14 - "worker"
Cohesion: 0.12
Nodes (14): Known issue, Step 15 — Store the reading, Step 16 — Every second: dead connection check and live status, Value, CycleStatusFor(), LiveValue, TelemetryValue, roundToTenth() (+6 more)

### Community 15 - "MachineService"
Cohesion: 0.12
Nodes (16): CreateMachineRequest, MachineFields, UpdateMachineRequest, validCoordinate(), validHealthIDs(), validIPv4(), Machine, MachineScope (+8 more)

### Community 16 - "newOrganizationRouter"
Cohesion: 0.18
Nodes (14): NewBranchHandler(), NewUserCityHandler(), NewUserLocationHandler(), NewUserStateHandler(), UserState, UserStateRepository, mapUserStateError(), NewUserStateRepository() (+6 more)

### Community 17 - "expectStatus"
Cohesion: 0.28
Nodes (23): decodeData(), expectStatus(), newCatalogRouter(), sendRequest(), uniqueSuffix(), createAuthTestUser(), createTestRecord(), TestControllerBrand() (+15 more)

### Community 18 - "MachineTelemetryService"
Cohesion: 0.13
Nodes (9): 4. Find the machines: steps 5 – 9, Step 5 — Prepare the database, Step 7 — Listen for database changes, Step 8 — Load all active machines, ListenForDatabaseNotifications(), NewHealthHandler(), HealthService, NewHealthService() (+1 more)

### Community 19 - ".runConnectionLoop"
Cohesion: 0.16
Nodes (10): 5. Connect and collect, per machine: steps 10 – 17, Retry and timing summary, Step 10 — Get ready to connect, Step 11 — Connect, a few machines at a time, Step 14 — The data collection loop, FieldIndex(), TestFieldIndexAxisPositions(), secondsOf() (+2 more)

### Community 20 - "net/http.ServeMux"
Cohesion: 0.16
Nodes (31): RequireRole(), registerBranchRoutes(), registerBrandRoutes(), registerBreakRoutes(), registerCategoryRoutes(), registerClienteleRoutes(), registerControllerRoutes(), registerDowntimeReasonRoutes() (+23 more)

### Community 21 - "HolidayRepository"
Cohesion: 0.38
Nodes (3): HolidayRepository, mapHolidayError(), NewHolidayService()

### Community 22 - "Querier"
Cohesion: 0.18
Nodes (11): Querier, NewClienteleRepository(), NewDowntimeReasonRepository(), NewHolidayRepository(), NewIdleReasonRepository(), NewLicenceKeyRepository(), NewMachineCycleRepository(), NewMachineTelemetryRepository() (+3 more)

### Community 23 - "§10 Diagnostics & Troubleshooting"
Cohesion: 0.12
Nodes (25): Check mode: go run sinumerik_client.go check (OPC UA node existence report), Go Client (sinumerik_client.go, callback -> copy -> render), Prerequisite: 32-bit x86 DLL, compile Go with GOARCH=386, Prerequisite: SIEMENS_USER / SIEMENS_PASS env vars before CNCConnect, Run: go run sinumerik_client.go (live dashboard), §3.1 Environment & Prerequisites, §8 Running the Go Client, §10 Diagnostics & Troubleshooting (+17 more)

### Community 24 - "mode-chart.js"
Cohesion: 0.21
Nodes (23): chartData(), countByGroup(), drawDashboard(), drawHeader(), drawMachineList(), drawTiles(), drawTotals(), fetchData() (+15 more)

### Community 25 - "machine_mode.go"
Cohesion: 0.14
Nodes (21): engineName(), cycleReadingFrom(), decideMachineMode(), holdReason(), machineStateFromReading(), numberingEngine(), programStateFromProcStatus(), ncReading() (+13 more)

### Community 26 - "BranchService"
Cohesion: 0.15
Nodes (9): CreateBranchRequest, UpdateBranchRequest, validateBranchFields(), Branch, BranchRepository, mapBranchError(), cityOwner(), BranchService (+1 more)

### Community 27 - "actAs"
Cohesion: 0.18
Nodes (16): CreateUserRequest, WithAuthUser(), asTestUser(), apiResponse, actAs(), actAsMember(), createUser(), listIDs() (+8 more)

### Community 28 - "CNC Data Collection Flow"
Cohesion: 0.12
Nodes (14): 1. The big picture, 6. When data is edited: step 18, 7. Shutdown: steps 19 – 20, 8. Reference, 9. Data Reading Flow (Function Names Only), docs/cnc-communication.md, CNC Data Collection Flow, Events in `machine_events` (+6 more)

### Community 29 - "UserLocationService"
Cohesion: 0.16
Nodes (8): CreateUserLocationRequest, UpdateUserLocationRequest, validateUserLocation(), UserLocation, UserLocationRepository, mapUserLocationError(), UserLocationService, toUserLocationResponseError()

### Community 30 - "Execute"
Cohesion: 0.06
Nodes (26): Execute(), UserCity, UserProfile, AuthRepository, mapAuthError(), mapLocationError(), NewLocationRepository(), UserCityRepository (+18 more)

### Community 31 - "config.go"
Cohesion: 0.19
Nodes (12): main(), run(), Config, getEnv(), getEnvBool(), getEnvInt(), getEnvMachineSelection(), DBConfig (+4 more)

### Community 32 - "Worker (internal/monitor/worker.go)"
Cohesion: 0.07
Nodes (22): 7th Sense Server Claude Context, response.SendResponse Format, User CRUD Golden Reference, Collector (internal/monitor/collector.go), CNC Driver (internal/cnc), Known Issue: nil live-value cache panic, machine_events, machine_telementry hypertable (+14 more)

### Community 33 - "UserService"
Cohesion: 0.06
Nodes (34): CustomerDetails, UpdateUserRequest, UserResponse, ParseDayStart(), TestParseDayStart(), Address, AssignableRoles(), CanAssignRole() (+26 more)

### Community 34 - "MachineModelService"
Cohesion: 0.20
Nodes (6): MachineModelRequest, MachineModel, MachineModelRepository, mapMachineModelError(), MachineModelService, toMachineModelResponseError()

### Community 35 - ".Validate"
Cohesion: 0.16
Nodes (14): UpdateThemeRequest, isTimeZone(), isWebsite(), normalizeFilters(), normalizeTheme(), TestCreateUserRequestValidate(), TestCreateUserRequestValidateDropsDuplicateMachines(), TestCreateUserRequestValidateNormalizes() (+6 more)

### Community 36 - "LicenceKeyService"
Cohesion: 0.11
Nodes (15): NewLicenceKey, CreateLicenceKeysRequest, UserListFilter, NewLicenceKeyHandler(), LicenceKey, User, LicenceKeyRepository, mapLicenceKeyError() (+7 more)

### Community 37 - "BrandService"
Cohesion: 0.25
Nodes (5): Brand, BrandRepository, mapBrandError(), BrandService, toBrandResponseError()

### Community 38 - "AuthUser"
Cohesion: 0.06
Nodes (53): WithTransaction(), validDate(), CreatePartRequest, PartFields, UpdatePartRequest, ValidPartDrawingFile(), CreateProductionScheduleRequest, ProductionScheduleFields (+45 more)

### Community 39 - "formatDuration"
Cohesion: 0.15
Nodes (30): card(), cell(), cycleData(), cycleEnd(), cycleEndReason(), cycleJob(), cycleList(), cycleNumber() (+22 more)

### Community 40 - "CncCommunication.dll API & Integration G"
Cohesion: 0.19
Nodes (20): CncCommunication.dll (C++ core runtime, x86, OPC UA inside), CncCommunication.dll API & Integration Guide (Go Client) - text version, Fanuc profile (engine 1, port 8193, Fwlib32.dll + fwlibe1.dll), Fwlib32.dll (Fanuc FOCAS2), fwlibe1.dll (Fanuc FOCAS2), libcrypto-3.dll (OpenSSL), libssl-3.dll (OpenSSL), Engine 1: Fanuc (FOCAS2 Ethernet, fwlib32) (+12 more)

### Community 41 - "sinumerik_client.go"
Cohesion: 0.20
Nodes (17): axisColor(), cstr(), enableVT(), findLibDir(), initDLL(), isRotary(), lbl(), logf() (+9 more)

### Community 42 - "MachineTelemetryRepository"
Cohesion: 0.13
Nodes (6): MachineEvent, MachineLiveStatus, ModeChange, MachineTelemetryRepository, mapMachineTelemetryError(), execute()

### Community 43 - "CncCommunication.dll API & Integration G"
Cohesion: 0.13
Nodes (18): CNCConnect param machineType (int32 engine id, see §4.2), Engine 0: Simulator (local loopback mock engine), Engine 2: Siemens (Sinumerik OPC UA, active profile), Engine 3: Syntec (Open API / TCP), Engine 4: Mitsubishi (driver not specified), Machine identifier mapping (engine id enum 0-4), §9 Controller-Specific Profiles, Siemens default ports: 4840 (native NCU) / 4841 (custom) (+10 more)

### Community 44 - "validateCatalogFields()"
Cohesion: 0.11
Nodes (11): CreateBrandRequest, UpdateBrandRequest, validateCatalogFields(), CreateCategoryRequest, UpdateCategoryRequest, CreateControllerRequest, UpdateControllerRequest, validateControllerFields() (+3 more)

### Community 45 - ".Login()"
Cohesion: 0.18
Nodes (7): LoginRequest, newSessionToken(), truncate(), newLoginThrottle(), loginFailures, loginThrottle, Session

### Community 46 - "Product"
Cohesion: 0.23
Nodes (6): Product, ProductRepository, mapProductError(), ProductService, productImageOrDefault(), toProductResponseError()

### Community 47 - "demo_monitor/main.go"
Cohesion: 0.15
Nodes (9): CncMasterSnapshot, MachineDef, cstr(), enableVT(), initDLL(), main(), render(), TestSnapshotModeFieldOffsets() (+1 more)

### Community 48 - "Machines API (/machines)"
Cohesion: 0.16
Nodes (18): Clienteles API, Company Scope (companyOwner / ownCompanyCustomer), Current User API (/me theme, language), Customer Display Names (customerName / companyName), Holidays API, Licence Keys API, Machine Catalog Brand-Controller-Model, Machines API (/machines) (+10 more)

### Community 49 - "7th-sense-server/internal/model.AuthUser"
Cohesion: 0.17
Nodes (11): CreateIdleReasonRequest, IdleReasonFields, UpdateIdleReasonRequest, NewIdleReasonHandler(), IdleReason, IdleReasonRepository, mapIdleReasonError(), IdleReasonService (+3 more)

### Community 50 - "ProductionScheduleRepository"
Cohesion: 0.19
Nodes (4): ProductionScheduleItem, ProductionScheduleRepository, mapProductionScheduleError(), saveScheduleItems()

### Community 51 - "DowntimeReasonService"
Cohesion: 0.17
Nodes (11): CreateDowntimeReasonRequest, DowntimeReasonFields, UpdateDowntimeReasonRequest, NewDowntimeReasonHandler(), DowntimeReason, DowntimeReasonRepository, mapDowntimeReasonError(), downtimeReasonFromFields() (+3 more)

### Community 52 - "UserCountryService"
Cohesion: 0.18
Nodes (9): CreateUserCountryRequest, NewUserCountryHandler(), UserCountry, UserCountryRepository, mapUserCountryError(), NewUserCountryRepository(), UserCountryService, NewUserCountryService() (+1 more)

### Community 53 - "net/http.Handler"
Cohesion: 0.28
Nodes (15): NewMachineHandler(), CrossOrigin(), Logger(), Chain(), Middleware, Recoverer(), NewUserRepository(), registerAuthRoutes() (+7 more)

### Community 54 - "AuthHandler"
Cohesion: 0.38
Nodes (3): AuthHandler, CurrentUser(), UserID()

### Community 55 - "Category"
Cohesion: 0.25
Nodes (5): Category, CategoryRepository, mapCategoryError(), CategoryService, toCategoryResponseError()

### Community 56 - "net/http.Request"
Cohesion: 0.14
Nodes (6): BranchHandler, ClienteleHandler, currentUserID(), listUserID(), queryID(), UserHandler

### Community 57 - "context.Context"
Cohesion: 0.14
Nodes (4): PartSubProgram, executeAffected(), PartRepository, mapPartError()

### Community 58 - "Break"
Cohesion: 0.15
Nodes (11): BreakFields, CreateBreakRequest, UpdateBreakRequest, NewBreakHandler(), Break, ShiftBreak, BreakRepository, mapBreakError() (+3 more)

### Community 59 - "NewError"
Cohesion: 0.15
Nodes (9): CreateAddressRequest, UpdateAddressRequest, UpdateLanguageRequest, UpdateLicenceKeyRequest, validDuration(), ValidLicenceKey(), NewError(), newResponse() (+1 more)

### Community 60 - "CNC Data Collection (in-process)"
Cohesion: 0.15
Nodes (12): CNC Data Collection (in-process), writeLog CNC Log Types and Switches, .env Configuration (internal/config), Health Mode History API (public), CncCommunication.dll (32-bit vendor DLL), machine_mode_log, Mode Decision (machine_mode.go), ECharts 5.5.1 (+4 more)

### Community 61 - "ClienteleService"
Cohesion: 0.38
Nodes (5): NewClienteleHandler(), Clientele, clienteleFromFields(), ClienteleService, toClienteleResponseError()

### Community 62 - "validText()"
Cohesion: 0.15
Nodes (6): PartProcessFields, PartSubProgramFields, UploadPartDrawingRequest, uniquePositiveIDs(), validText(), validWorkingDays()

### Community 63 - "run"
Cohesion: 0.13
Nodes (32): run(), NewBrandHandler(), NewCategoryHandler(), NewControllerHandler(), NewHolidayHandler(), NewMachineModelHandler(), NewPartHandler(), NewProductHandler() (+24 more)

### Community 64 - "cycleReading"
Cohesion: 0.17
Nodes (12): Completed, How a cycle ends, Machine cycles, One running cycle per machine, Readings that are skipped, Unfinished, Where to see it, decideCycleStep() (+4 more)

### Community 65 - "runCatalogCRUD"
Cohesion: 0.31
Nodes (8): containsID(), runCatalogCRUD(), TestBrandCRUD(), TestCategoryCRUD(), TestControllerCRUD(), TestProductCRUD(), catalogCRUDCase, catalogRecord

### Community 67 - "TestAuthentication"
Cohesion: 0.67
Nodes (4): authRequest(), login(), sessionCookieOf(), TestAuthentication()

### Community 68 - "TestProductionScheduleCRUD"
Cohesion: 0.20
Nodes (18): PartDrawing, clienteleBody(), TestClienteleCRUD(), containsMachineID(), machineBody(), newMachineFixture(), TestMachineCRUD(), with() (+10 more)

### Community 69 - "CNCConnect(ip, port, dataCb, errCb, mach"
Cohesion: 0.20
Nodes (11): CNCConnect(ip, port, dataCb, errCb, machineType) int32, CNCConnect param dataCb (uintptr stdcall OnData(const CncMasterSnapshot*)), CNCConnect param ip (const char*, ANSI NUL-terminated), CNCConnect param port (int32, e.g. 4840), CNCdisconnect(machineId int32), LazyDLL binding (syscall.NewLazyDLL + eager mod.Load in init), §4 API Reference, §5 DLL Loading & Binding (Go) (+3 more)

### Community 70 - "newUserHandler"
Cohesion: 0.20
Nodes (11): NewRoleHandler(), NewUserHandler(), Role, RoleRepository, mapRoleError(), NewRoleRepository(), RoleService, NewRoleService() (+3 more)

### Community 71 - "onDataCallback(ptr uintptr) uintptr"
Cohesion: 0.36
Nodes (9): Package-level callback singletons (syscall.NewCallback), CNCConnect param errCb (uintptr stdcall OnErr(int32 code, const char* msg)), Mutex-protected copy: dataMutex sync.RWMutex -> lastRecord/hasRecord, onDataCallback(ptr uintptr) uintptr, onErrorCallback(code, messagePtr uintptr) uintptr, §7 Asynchronous Callback Architecture, onDataCallback, onErrorCallback (+1 more)

### Community 72 - "drawChart"
Cohesion: 0.24
Nodes (11): chartRows(), currentSpan(), currentWindow(), drawChart(), drawZoomLevels(), markZoomLevel(), renderBar(), renderCycleBar() (+3 more)

### Community 73 - "VariableConfig"
Cohesion: 0.24
Nodes (7): VariableConfig, formatPreviousValue(), formatValue(), newValueChangeLog(), boundVariable, valueChange, valueChangeLog

### Community 74 - "AuthService"
Cohesion: 0.31
Nodes (7): NewAuthHandler(), Authenticate(), NewAuthRepository(), AuthService, hashToken(), NewAuthService(), newAuthRouter()

### Community 75 - "ClienteleRepository"
Cohesion: 0.28
Nodes (3): ClienteleRepository, mapClienteleError(), NewClienteleService()

### Community 76 - ".Create"
Cohesion: 0.29
Nodes (3): CreateUserStateRequest, validateUserState(), toUserStateResponseError()

### Community 77 - "Client Certificate & Key (client_cert.de"
Cohesion: 0.48
Nodes (7): client_cert.der/.pem (X.509 client cert for Siemens OPC UA), client_key.der/.pem (client private key), Siemens certs: client_cert.der + client_key.der in binary directory, Siemens trust: cert must be trusted in NCU Access Administration, Troubleshoot: BadRequestNotAllowed -> Export client_cert, BadRequestNotAllowed (certificate trust rejection), Client Certificate & Key (client_cert.der / client_key.der)

### Community 78 - ".validate()"
Cohesion: 0.38
Nodes (4): ClienteleFields, CreateClienteleRequest, UpdateClienteleRequest, validClienteleEmail()

### Community 79 - "NewPartService"
Cohesion: 0.25
Nodes (4): PartDrawingStorage, NewPartDrawingStorage(), NewPartRepository(), NewPartService()

### Community 80 - "Event-driven data flow (no polling in cl"
Cohesion: 0.60
Nodes (5): Flow 1: load DLL, register two stdcall callbacks (syscall.NewCallback), Flow 2: CNCConnect starts DLL background worker thread, Flow 3: worker polls controller, assembles snapshot, invokes data callback, Flow 4: callback copies payload, signals renderer via buffered channel, Flow 5: main goroutine renders live ANSI dashboard

### Community 81 - "logger/logger.go"
Cohesion: 0.40
Nodes (3): main(), Init(), InitFile()

### Community 82 - ".validate"
Cohesion: 0.47
Nodes (3): CreateHolidayRequest, HolidayFields, UpdateHolidayRequest

### Community 85 - "MachineModePeriod"
Cohesion: 0.53
Nodes (4): MachineModeHistoryResponse, MachineModeTotal, MachineModePeriod, modeTotals()

### Community 86 - "auth_tokens SHA-256 Hashed Sessions (mig"
Cohesion: 0.67
Nodes (3): auth_tokens SHA-256 Hashed Sessions (migration 014), Per-Request Session and User Re-check, ts_session Cookie (HttpOnly, SameSite=Strict)

### Community 88 - ".Generate"
Cohesion: 0.40
Nodes (3): GenerateLicenceKeysRequest, GeneratedLicenceKeys, newLicenceKey()

## Ambiguous Edges - Review These
- `CncCommunication.dll (32-bit vendor DLL)` → `FanucCncCommunication.dll`  [AMBIGUOUS]
  docs/setup-guide.txt · relation: semantically_similar_to
- `Snapshot Fields (positions, spindle, feed, program, tool, alarms)` → `Machines API (/machines)`  [AMBIGUOUS]
  docs/CncCommunication.dll — API & Integration Guide (Go Client).pdf · relation: conceptually_related_to

## Knowledge Gaps
- **36 isolated node(s):** `hello`, `7th-sense-server`, `middleware.RequireRole / routes.adminOnly`, `response.SendResponse`, `Test Auth (asTestUser, tests/auth_test.go)` (+31 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 102 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **22 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `CNC Data Collection Flow` connect `CNC Data Collection Flow` to `Worker (internal/monitor/worker.go)`, `Collector`, `MachineTelemetryService`, `.runConnectionLoop`, `CNC Data Collection (in-process)`?**
  _High betweenness centrality (0.149) - this node is a cross-community bridge._
- **What connects `hello`, `7th-sense-server`, `middleware.RequireRole / routes.adminOnly` to the rest of the system?**
  _36 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `UserCellService` be split into smaller, more focused modules?**
  _Cohesion score 0.07493061979648474 - nodes in this community are weakly interconnected._
- **Why does `2. All steps at a glance` connect `Collector` to `.runConnectionLoop`, `dll_windows_386.go`, `CNC Data Collection Flow`, `worker`?**
  _High betweenness centrality (0.130) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.1341991341991342 - nodes in this community are weakly interconnected._

### Low-confidence Hints
_AMBIGUOUS edges — the extractor was unsure. Verify before acting on these._

- **What is the exact relationship between `CncCommunication.dll (32-bit vendor DLL)` and `FanucCncCommunication.dll`?**
  _Edge tagged AMBIGUOUS (relation: semantically_similar_to) - confidence is low._
- **What is the exact relationship between `Snapshot Fields (positions, spindle, feed, program, tool, alarms)` and `Machines API (/machines)`?**
  _Edge tagged AMBIGUOUS (relation: conceptually_related_to) - confidence is low._