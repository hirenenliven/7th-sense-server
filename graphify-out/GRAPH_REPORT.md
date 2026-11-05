# Graph Report - 7th-sense-server  (2026-11-04)

## Corpus Check
- 277 files · ~212,760 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 4 file(s) not represented in the graph (top: (none) 3, .toml 1)

## Summary
- 2385 nodes · 19067 edges · 92 communities (75 shown, 17 thin omitted)
- Extraction: 95% EXTRACTED · 5% INFERRED · 0% AMBIGUOUS · INFERRED: 905 edges (avg confidence: 0.87)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `1ab02b8e`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- net/http.Handler
- UserCityService
- UserShopService
- Machine
- SendResponse()
- pathID()
- SendResponse
- MachineCycleService
- api_helpers_test.go
- CncMasterSnapshot (34,440 bytes, #pragma
- run
- Connection
- context.Context
- Collector
- worker
- UserLineService
- context.Context
- Category
- Execute
- executeAffected
- net/http.ServeMux
- CNCConnect(ip, port, dataCb, errCb, machineType) int32
- queryID
- §10 Diagnostics & Troubleshooting
- SendResponse
- machine_mode.go
- RoleRepository
- UserService
- mode-chart.js
- UserLocationService
- NewError
- validateCatalogFields
- Worker (internal/monitor/worker.go)
- cycleReading
- ControllerService
- DashboardService
- LicenceKeyService
- Client Certificate & Key (client_cert.der / client_key.der)
- testing.T
- formatDuration
- CncCommunication.dll API & Integration G
- expectStatus
- UserProfileRepository
- CncCommunication.dll API & Integration G
- github.com/jackc/pgx/v5/pgxpool.Pool
- .Login()
- .runConnectionLoop
- VariableConfig
- Machines API (/machines)
- MachineCycle
- Product
- CNC Data Collection Flow
- BranchService
- README.md
- AuthHandler
- TestProductionScheduleCRUD
- net/http.Request
- CompanyProfile
- apidocs/main.go
- .Generate
- CNC Data Collection (in-process)
- Test list
- .GetMachineSummaries
- Address
- run
- createAuthTestUser
- UserCellService
- drawChart
- companyOwner
- Querier
- Decision
- onDataCallback(ptr uintptr) uintptr
- .Validate
- LicenceKeyDuration
- time.Time
- LoginAsRequest
- Event-driven data flow (no polling in client)
- auth_tokens SHA-256 Hashed Sessions (mig
- Create Customer in One Transaction (user
- 7th-sense-server
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
1. `Execute()` - 183 edges
2. `SendResponse()` - 156 edges
3. `NewError()` - 136 edges
4. `run()` - 96 edges
5. `AuthUser` - 83 edges
6. `pathID()` - 80 edges
7. `currentUser()` - 62 edges
8. `Querier` - 60 edges
9. `decodeAndValidate()` - 59 edges
10. `worker` - 57 edges

## Surprising Connections (you probably didn't know these)
- `Machine Identifier Mapping (0 Simulator,1 Fanuc,2 Siemens,3 Syntec,4 Mitsubishi)` --semantically_similar_to--> `Machine Catalog Brand-Controller-Model`  [INFERRED] [semantically similar]
  docs/CncCommunication.dll — API & Integration Guide (Go Client).pdf → .claude/CLAUDE.md
- `Decision` --references--> `SendResponse()`  [INFERRED]
  docs/adr/0002-layered-handler-service-repository.md → internal/response/sendResponse.go
- `Snapshot Fields (positions, spindle, feed, program, tool, alarms)` --conceptually_related_to--> `Machines API (/machines)`  [AMBIGUOUS]
  docs/CncCommunication.dll — API & Integration Guide (Go Client).pdf → .claude/CLAUDE.md
- `Step 13 — The DLL sends data` --references--> `Reading`  [INFERRED]
  docs/cnc-data-collection-flow.md → internal/cnc/cnc.go
- `Decision` --references--> `Authenticate()`  [INFERRED]
  docs/adr/0004-session-cookie-auth.md → internal/middleware/auth.go

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

## Communities (92 total, 17 thin omitted)

### Community 0 - "net/http.Handler"
Cohesion: 0.08
Nodes (30): Handler(), NewBreakHandler(), NewClienteleHandler(), NewHolidayHandler(), NewLicenceKeyHandler(), NewPartHandler(), NewProductionScheduleHandler(), NewShiftHandler() (+22 more)

### Community 1 - "UserCityService"
Cohesion: 0.08
Nodes (15): CreateUserCityRequest, UpdateUserCityRequest, validateUserCity(), CreateUserStateRequest, validateUserState(), UserCity, UserState, UserCityRepository (+7 more)

### Community 2 - "UserShopService"
Cohesion: 0.15
Nodes (9): CreateUserShopRequest, UpdateUserShopRequest, validateUserShop(), NewUserShopHandler(), UserShop, UserShopRepository, mapUserShopError(), UserShopService (+1 more)

### Community 3 - "Machine"
Cohesion: 0.12
Nodes (17): CreateMachineRequest, MachineFields, UpdateMachineRequest, validCoordinate(), validHealthIDs(), validIPv4(), SeesAllCompanyMachines(), Machine (+9 more)

### Community 4 - "SendResponse()"
Cohesion: 0.08
Nodes (10): BreakHandler, DowntimeReasonHandler, IdleReasonHandler, MachineCycleHandler, PartHandler, ProductionScheduleHandler, currentUser(), RoleHandler (+2 more)

### Community 5 - "pathID()"
Cohesion: 0.11
Nodes (12): validator, decodeAndValidate(), pathID(), targetUserID(), UserCellHandler, UserCityHandler, UserCountryHandler, UserLineHandler (+4 more)

### Community 6 - "SendResponse"
Cohesion: 0.11
Nodes (5): BrandHandler, ControllerHandler, LicenceKeyHandler, MachineModelHandler, statusRecorder

### Community 7 - "MachineCycleService"
Cohesion: 0.22
Nodes (6): One running cycle per machine, MachineCyclesResponse, NewMachineCycleHandler(), MachineCycleService, sameID(), toMachineCycleResponseError()

### Community 8 - "api_helpers_test.go"
Cohesion: 0.28
Nodes (7): cncMasterSnapshot, LanguageResponse, ThemeResponse, authUserKey, loginAsKey, Envelope, machineModelRecord

### Community 9 - "CncMasterSnapshot (34,440 bytes, #pragma"
Cohesion: 0.10
Nodes (43): CncMasterSnapshot (34,440 bytes, #pragma pack(8), x86), AbsPos ([5]float64, Inches / Degrees, Work Coordinate System (WCS) positions), ActiveGCode ([16]byte, ASCII, Currently active G-code group), ActiveProgPath ([256]byte, ASCII Path, Storage path of active NC artifact), AlarmMsg ([128]byte, UTF-8 / ASCII, Alarm description message), AlarmNo (int32, Code (0 = OK), Primary alarm or fault number), AxisLoadXPercentage (int16, Percent (%), Servo X axis drive load (AxisLoadX/Y/ZPercentage)), AxisLoadYPercentage (int16, Percent (%), Servo Y axis drive load (AxisLoadX/Y/ZPercentage)) (+35 more)

### Community 10 - "run"
Cohesion: 0.08
Nodes (55): run(), NewBranchHandler(), NewBrandHandler(), NewCategoryHandler(), NewControllerHandler(), NewDashboardHandler(), NewMachineHandler(), NewMachineModelHandler() (+47 more)

### Community 11 - "Connection"
Cohesion: 0.11
Nodes (14): dllDriver, slot, Step 12 — Call the DLL's CNCConnect, Step 13 — The DLL sends data, Step 17 — Disconnect and try again, Connection, Stats, newConnection() (+6 more)

### Community 12 - "context.Context"
Cohesion: 0.07
Nodes (24): CreateProductionScheduleRequest, ProductionScheduleFields, UpdateProductionScheduleRequest, Brand, ProductionSchedule, ProductionScheduleItem, ShiftBreak, BrandRepository (+16 more)

### Community 13 - "Collector"
Cohesion: 0.10
Nodes (26): Supervisor, 2. All steps at a glance, 3. Startup: steps 1 – 4, 7. Shutdown: steps 19 – 20, Log lines (in `logs/app.log`, search for the message), Step 18 — React to machine and variable changes, Step 19 — Stop collection with the server, Step 1 — Start collection with the server (+18 more)

### Community 14 - "worker"
Cohesion: 0.12
Nodes (15): Known issue, Step 15 — Store the reading, Step 16 — Every second: dead connection check and live status, `CNC_CYCLE_MAX_HOURS` (24 hours), `CNC_STALE_SEC` (30 s) and `CNC_QUIET_SEC` (300 s), `cycleCheckpointInterval` (5 minutes), `cycleCounterGrace` and `cycleStartGrace` (5 seconds each), `emergencyResetGrace` (2 seconds) (+7 more)

### Community 15 - "UserLineService"
Cohesion: 0.15
Nodes (9): CreateUserLineRequest, UpdateUserLineRequest, validateUserLine(), NewUserLineHandler(), UserLine, UserLineRepository, mapUserLineError(), UserLineService (+1 more)

### Community 16 - "context.Context"
Cohesion: 0.06
Nodes (24): 4. Find the machines: steps 5 – 9, Step 5 — Prepare the database, Step 7 — Listen for database changes, Step 8 — Load all active machines, ListenForDatabaseNotifications(), MachineModesResponse, NewHealthHandler(), equalEngine() (+16 more)

### Community 17 - "Category"
Cohesion: 0.25
Nodes (5): Category, CategoryRepository, mapCategoryError(), CategoryService, toCategoryResponseError()

### Community 18 - "Execute"
Cohesion: 0.06
Nodes (31): Connect(), Execute(), MachineModelRequest, CreateUserCountryRequest, MachineModel, UserCountry, AuthRepository, mapAuthError() (+23 more)

### Community 19 - "executeAffected"
Cohesion: 0.08
Nodes (19): PartProcessFields, PartSubProgramFields, CreatePartRequest, PartFields, UpdatePartRequest, UploadPartDrawingRequest, ValidPartDrawingFile(), validText() (+11 more)

### Community 20 - "net/http.ServeMux"
Cohesion: 0.15
Nodes (34): Decision, RequireRole(), registerAuthRoutes(), registerBranchRoutes(), registerBrandRoutes(), registerBreakRoutes(), registerCategoryRoutes(), registerClienteleRoutes() (+26 more)

### Community 21 - "CNCConnect(ip, port, dataCb, errCb, machineType) int32"
Cohesion: 0.20
Nodes (11): CNCConnect(ip, port, dataCb, errCb, machineType) int32, CNCConnect param dataCb (uintptr stdcall OnData(const CncMasterSnapshot*)), CNCConnect param ip (const char*, ANSI NUL-terminated), CNCConnect param port (int32, e.g. 4840), CNCdisconnect(machineId int32), LazyDLL binding (syscall.NewLazyDLL + eager mod.Load in init), §4 API Reference, §5 DLL Loading & Binding (Go) (+3 more)

### Community 22 - "queryID"
Cohesion: 0.17
Nodes (5): ClienteleHandler, dashboardQuery(), queryIDList(), listUserID(), queryID()

### Community 23 - "§10 Diagnostics & Troubleshooting"
Cohesion: 0.12
Nodes (25): Check mode: go run sinumerik_client.go check (OPC UA node existence report), Go Client (sinumerik_client.go, callback -> copy -> render), Prerequisite: 32-bit x86 DLL, compile Go with GOARCH=386, Prerequisite: SIEMENS_USER / SIEMENS_PASS env vars before CNCConnect, Run: go run sinumerik_client.go (live dashboard), §3.1 Environment & Prerequisites, §8 Running the Go Client, §10 Diagnostics & Troubleshooting (+17 more)

### Community 24 - "SendResponse"
Cohesion: 0.13
Nodes (7): DashboardHandler, readDashboardRequest(), HealthHandler, HolidayHandler, MachineHandler, SendResponse(), registerHealthRoutes()

### Community 25 - "machine_mode.go"
Cohesion: 0.12
Nodes (25): Reading, Status, engineName(), snapshotStatusOf(), decideMachineMode(), holdReason(), machineStateFromReading(), numberingEngine() (+17 more)

### Community 26 - "RoleRepository"
Cohesion: 0.33
Nodes (4): Role, RoleRepository, mapRoleError(), toRoleResponseError()

### Community 27 - "UserService"
Cohesion: 0.14
Nodes (17): CustomerDetails, UpdateUserRequest, UserResponse, ParseDayStart(), AssignableRoles(), CanAssignRole(), RoleService, canSee() (+9 more)

### Community 28 - "mode-chart.js"
Cohesion: 0.21
Nodes (23): chartData(), countByGroup(), drawDashboard(), drawHeader(), drawMachineList(), drawTiles(), drawTotals(), fetchData() (+15 more)

### Community 29 - "UserLocationService"
Cohesion: 0.16
Nodes (8): CreateUserLocationRequest, UpdateUserLocationRequest, validateUserLocation(), UserLocation, UserLocationRepository, mapUserLocationError(), UserLocationService, toUserLocationResponseError()

### Community 30 - "NewError"
Cohesion: 0.06
Nodes (24): CreateAddressRequest, UpdateAddressRequest, UpdateLanguageRequest, ClienteleFields, CreateClienteleRequest, UpdateClienteleRequest, validClienteleEmail(), CreateHolidayRequest (+16 more)

### Community 31 - "validateCatalogFields"
Cohesion: 0.11
Nodes (11): CreateBrandRequest, UpdateBrandRequest, validateCatalogFields(), CreateCategoryRequest, UpdateCategoryRequest, CreateControllerRequest, UpdateControllerRequest, validateControllerFields() (+3 more)

### Community 32 - "Worker (internal/monitor/worker.go)"
Cohesion: 0.07
Nodes (22): 7th Sense Server Claude Context, response.SendResponse Format, User CRUD Golden Reference, Collector (internal/monitor/collector.go), CNC Driver (internal/cnc), Known Issue: nil live-value cache panic, machine_events, machine_telementry hypertable (+14 more)

### Community 33 - "cycleReading"
Cohesion: 0.18
Nodes (8): CycleStatusFor(), cycleReadingFrom(), decideCycleStep(), isProgramLoop(), resetReason(), cycleReading, cycleStep, cycleTracker

### Community 34 - "ControllerService"
Cohesion: 0.19
Nodes (6): Controller, ControllerRepository, mapControllerError(), ControllerService, toControllerResponseError(), asInvalidSelection()

### Community 35 - "DashboardService"
Cohesion: 0.25
Nodes (9): CycleCountRow, DashboardRequest, MachineTimeRow, cyclesIn(), DashboardService, hours(), minTime(), modeGroup() (+1 more)

### Community 36 - "LicenceKeyService"
Cohesion: 0.21
Nodes (8): NewLicenceKey, CreateLicenceKeysRequest, LicenceKey, LicenceKeyRepository, mapLicenceKeyError(), LicenceKeyService, licenceExpiresOn(), toLicenceKeyResponseError()

### Community 37 - "Client Certificate & Key (client_cert.der / client_key.der)"
Cohesion: 0.48
Nodes (7): client_cert.der/.pem (X.509 client cert for Siemens OPC UA), client_key.der/.pem (client private key), Siemens certs: client_cert.der + client_key.der in binary directory, Siemens trust: cert must be trusted in NCU Access Administration, Troubleshoot: BadRequestNotAllowed -> Export client_cert, BadRequestNotAllowed (certificate trust rejection), Client Certificate & Key (client_cert.der / client_key.der)

### Community 38 - "testing.T"
Cohesion: 0.08
Nodes (31): UpdateThemeRequest, TestSnapshotModeFieldOffsets(), TestSnapshotSize(), isTimeZone(), isWebsite(), normalizeFilters(), normalizeTheme(), TestCreateUserRequestValidate() (+23 more)

### Community 39 - "formatDuration"
Cohesion: 0.15
Nodes (30): card(), cell(), cycleData(), cycleEnd(), cycleEndReason(), cycleJob(), cycleList(), cycleNumber() (+22 more)

### Community 40 - "CncCommunication.dll API & Integration G"
Cohesion: 0.19
Nodes (20): CncCommunication.dll (C++ core runtime, x86, OPC UA inside), CncCommunication.dll API & Integration Guide (Go Client) - text version, Fanuc profile (engine 1, port 8193, Fwlib32.dll + fwlibe1.dll), Fwlib32.dll (Fanuc FOCAS2), fwlibe1.dll (Fanuc FOCAS2), libcrypto-3.dll (OpenSSL), libssl-3.dll (OpenSSL), Engine 1: Fanuc (FOCAS2 Ethernet, fwlib32) (+12 more)

### Community 41 - "expectStatus"
Cohesion: 0.15
Nodes (41): CreateUserRequest, containsID(), decodeData(), expectStatus(), newCatalogRouter(), runCatalogCRUD(), sendRequest(), uniqueSuffix() (+33 more)

### Community 42 - "UserProfileRepository"
Cohesion: 0.19
Nodes (5): UserProfile, UserProfileRepository, mapProfileError(), toUserProfileResponseError(), UserProfileService

### Community 43 - "CncCommunication.dll API & Integration G"
Cohesion: 0.13
Nodes (18): CNCConnect param machineType (int32 engine id, see §4.2), Engine 0: Simulator (local loopback mock engine), Engine 2: Siemens (Sinumerik OPC UA, active profile), Engine 3: Syntec (Open API / TCP), Engine 4: Mitsubishi (driver not specified), Machine identifier mapping (engine id enum 0-4), §9 Controller-Specific Profiles, Siemens default ports: 4840 (native NCU) / 4841 (custom) (+10 more)

### Community 44 - "github.com/jackc/pgx/v5/pgxpool.Pool"
Cohesion: 0.07
Nodes (34): WithTransaction(), uniquePositiveIDs(), ProductionScheduleItemFields, ScheduleDay(), clockMinutes(), CreateShiftRequest, ShiftFields, UpdateShiftRequest (+26 more)

### Community 45 - ".Login()"
Cohesion: 0.15
Nodes (12): LoginRequest, NewAuthHandler(), AuthService, hashToken(), NewAuthService(), newLoginAsKey(), newSessionToken(), truncate() (+4 more)

### Community 46 - ".runConnectionLoop"
Cohesion: 0.19
Nodes (9): 5. Connect and collect, per machine: steps 10 – 17, Step 10 — Get ready to connect, Step 11 — Connect, a few machines at a time, Step 14 — The data collection loop, FieldIndex(), TestFieldIndexAxisPositions(), secondsOf(), sleepUnlessCancelled() (+1 more)

### Community 47 - "VariableConfig"
Cohesion: 0.19
Nodes (8): Value, VariableConfig, formatPreviousValue(), formatValue(), newValueChangeLog(), boundVariable, valueChange, valueChangeLog

### Community 48 - "Machines API (/machines)"
Cohesion: 0.16
Nodes (18): Clienteles API, Company Scope (companyOwner / ownCompanyCustomer), Current User API (/me theme, language), Customer Display Names (customerName / companyName), Holidays API, Licence Keys API, Machine Catalog Brand-Controller-Model, Machines API (/machines) (+10 more)

### Community 49 - "MachineCycle"
Cohesion: 0.27
Nodes (4): CycleJob, MachineCycle, MachineCycleRepository, mapMachineCycleError()

### Community 50 - "Product"
Cohesion: 0.23
Nodes (6): Product, ProductRepository, mapProductError(), ProductService, productImageOrDefault(), toProductResponseError()

### Community 51 - "CNC Data Collection Flow"
Cohesion: 0.14
Nodes (13): 1. The big picture, 6. When data is edited: step 18, 8. Reference, 9. Data Reading Flow (Function Names Only), docs/cnc-communication.md, CNC Data Collection Flow, Events in `machine_events`, Main parts (+5 more)

### Community 52 - "BranchService"
Cohesion: 0.24
Nodes (7): CreateBranchRequest, UpdateBranchRequest, validateBranchFields(), Branch, cityOwner(), BranchService, toBranchResponseError()

### Community 53 - "README.md"
Cohesion: 0.04
Nodes (38): 0001. Go standard library net/http, no web framework, Alternatives considered, Consequences, Context, Decision, 0002. Layered architecture: handler, DTO, service, repository, Alternatives considered, Consequences (+30 more)

### Community 54 - "AuthHandler"
Cohesion: 0.22
Nodes (8): Decision, AuthHandler, CurrentUser(), isLoginAsKey(), LoginAsKey(), RouteLoginAs(), SessionToken(), UserID()

### Community 55 - "TestProductionScheduleCRUD"
Cohesion: 0.20
Nodes (18): PartDrawing, clienteleBody(), TestClienteleCRUD(), containsMachineID(), machineBody(), newMachineFixture(), TestMachineCRUD(), with() (+10 more)

### Community 56 - "net/http.Request"
Cohesion: 0.12
Nodes (6): BranchHandler, CategoryHandler, ProductHandler, currentUserID(), queryBool(), UserHandler

### Community 57 - "CompanyProfile"
Cohesion: 0.28
Nodes (5): CompanyProfile, CompanyProfileRepository, mapCompanyProfileError(), toCompanyProfileResponseError(), CompanyProfileService

### Community 58 - "apidocs/main.go"
Cohesion: 0.06
Nodes (45): request, main(), run(), main(), Config, MachineSelection, GenerateCertificates(), getEnv() (+37 more)

### Community 59 - ".Generate"
Cohesion: 0.40
Nodes (3): GenerateLicenceKeysRequest, GeneratedLicenceKeys, newLicenceKey()

### Community 60 - "CNC Data Collection (in-process)"
Cohesion: 0.15
Nodes (12): CNC Data Collection (in-process), writeLog CNC Log Types and Switches, .env Configuration (internal/config), Health Mode History API (public), CncCommunication.dll (32-bit vendor DLL), machine_mode_log, Mode Decision (machine_mode.go), ECharts 5.5.1 (+4 more)

### Community 61 - "Test list"
Cohesion: 0.05
Nodes (38): 1. The rows are consistent (SQL, no machine needed), 2. The rows agree with the mode log (SQL), 3. The raw signals around one cycle (telemetry and logs), 4. At the machine (floor test), Before the tests (once), Check query, Completed, Floor tests on the HUMA panel (Fanuc) (+30 more)

### Community 62 - ".GetMachineSummaries"
Cohesion: 0.23
Nodes (8): MachineSummary, OEEResponse, OEERow, NewMachineSummary(), MachineCycleSummary, round1(), summarizeCycles(), oeeSum

### Community 63 - "Address"
Cohesion: 0.27
Nodes (5): Address, AddressRepository, mapAddressError(), toAddressResponseError(), AddressService

### Community 64 - "run"
Cohesion: 0.36
Nodes (4): connectFanucCounter(), main(), run(), fanucCounter

### Community 65 - "createAuthTestUser"
Cohesion: 0.33
Nodes (10): NewAuthRepository(), authRequest(), createAuthTestUser(), login(), newAuthRouter(), sessionCookieOf(), TestAuthentication(), loginAsCookieOf() (+2 more)

### Community 66 - "UserCellService"
Cohesion: 0.07
Nodes (20): CreateUserCellRequest, UpdateUserCellRequest, validateUserCell(), CreateUserSupervisorRequest, UpdateUserSupervisorRequest, validateUserSupervisor(), UserCell, UserSupervisor (+12 more)

### Community 67 - "drawChart"
Cohesion: 0.24
Nodes (11): chartRows(), currentSpan(), currentWindow(), drawChart(), drawZoomLevels(), markZoomLevel(), renderBar(), renderCycleBar() (+3 more)

### Community 68 - "companyOwner"
Cohesion: 0.05
Nodes (40): LoginAsResponse, BreakFields, CreateBreakRequest, UpdateBreakRequest, CreateDowntimeReasonRequest, DowntimeReasonFields, UpdateDowntimeReasonRequest, CreateIdleReasonRequest (+32 more)

### Community 69 - "Querier"
Cohesion: 0.18
Nodes (17): Querier, NewAddressRepository(), NewCompanyProfileRepository(), UserMachineAccessRepository, NewUserMachineAccessRepository(), NewUserProfileRepository(), NewAddressService(), NewCompanyProfileService() (+9 more)

### Community 70 - "Decision"
Cohesion: 0.25
Nodes (9): Decision, FinishCyclesSummary, FinishedCycle, TimelinePeriod, FinishCyclesResponse, ModeTimelineResponse, NewFinishCyclesSummary(), NewFinishedCycle() (+1 more)

### Community 71 - "onDataCallback(ptr uintptr) uintptr"
Cohesion: 0.36
Nodes (9): Package-level callback singletons (syscall.NewCallback), CNCConnect param errCb (uintptr stdcall OnErr(int32 code, const char* msg)), Mutex-protected copy: dataMutex sync.RWMutex -> lastRecord/hasRecord, onDataCallback(ptr uintptr) uintptr, onErrorCallback(code, messagePtr uintptr) uintptr, §7 Asynchronous Callback Architecture, onDataCallback, onErrorCallback (+1 more)

### Community 77 - "time.Time"
Cohesion: 0.18
Nodes (16): ModeTotalsResponse, MachineModeHistoryResponse, MachineModeTotal, MachineModePeriod, ModeDuration, clipWindows(), dayWindows(), machineDayOEE() (+8 more)

### Community 80 - "Event-driven data flow (no polling in client)"
Cohesion: 0.60
Nodes (5): Flow 1: load DLL, register two stdcall callbacks (syscall.NewCallback), Flow 2: CNCConnect starts DLL background worker thread, Flow 3: worker polls controller, assembles snapshot, invokes data callback, Flow 4: callback copies payload, signals renderer via buffered channel, Flow 5: main goroutine renders live ANSI dashboard

### Community 86 - "auth_tokens SHA-256 Hashed Sessions (mig"
Cohesion: 0.67
Nodes (3): auth_tokens SHA-256 Hashed Sessions (migration 014), Per-Request Session and User Re-check, ts_session Cookie (HttpOnly, SameSite=Strict)

## Ambiguous Edges - Review These
- `FanucCncCommunication.dll` → `CncCommunication.dll (32-bit vendor DLL)`  [AMBIGUOUS]
  docs/setup-guide.txt · relation: semantically_similar_to
- `Snapshot Fields (positions, spindle, feed, program, tool, alarms)` → `Machines API (/machines)`  [AMBIGUOUS]
  docs/CncCommunication.dll — API & Integration Guide (Go Client).pdf · relation: conceptually_related_to

## Knowledge Gaps
- **35 isolated node(s):** `7th-sense-server`, `middleware.RequireRole / routes.adminOnly`, `response.SendResponse`, `Test Auth (asTestUser, tests/auth_test.go)`, `user_filters Filter Checklist` (+30 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 160 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **17 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `CNC Data Collection Flow` connect `CNC Data Collection Flow` to `Worker (internal/monitor/worker.go)`, `Collector`, `.runConnectionLoop`, `context.Context`, `CNC Data Collection (in-process)`?**
  _High betweenness centrality (0.201) - this node is a cross-community bridge._
- **What connects `7th-sense-server`, `middleware.RequireRole / routes.adminOnly`, `response.SendResponse` to the rest of the system?**
  _35 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `net/http.Handler` be split into smaller, more focused modules?**
  _Cohesion score 0.08403361344537816 - nodes in this community are weakly interconnected._
- **Why does `2. All steps at a glance` connect `Collector` to `.runConnectionLoop`, `CNC Data Collection Flow`, `Connection`, `worker`?**
  _High betweenness centrality (0.176) - this node is a cross-community bridge._
- **Should `UserCityService` be split into smaller, more focused modules?**
  _Cohesion score 0.08333333333333333 - nodes in this community are weakly interconnected._

### Low-confidence Hints
_AMBIGUOUS edges — the extractor was unsure. Verify before acting on these._

- **What is the exact relationship between `FanucCncCommunication.dll` and `CncCommunication.dll (32-bit vendor DLL)`?**
  _Edge tagged AMBIGUOUS (relation: semantically_similar_to) - confidence is low._
- **What is the exact relationship between `Snapshot Fields (positions, spindle, feed, program, tool, alarms)` and `Machines API (/machines)`?**
  _Edge tagged AMBIGUOUS (relation: conceptually_related_to) - confidence is low._