# Graph Report - 7th-sense-server  (2026-11-17)

## Corpus Check
- 327 files · ~267,073 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 4 file(s) not represented in the graph (top: (none) 3, .toml 1)

## Summary
- 2890 nodes · 25291 edges · 130 communities (102 shown, 28 thin omitted)
- Extraction: 95% EXTRACTED · 5% INFERRED · 0% AMBIGUOUS · INFERRED: 1183 edges (avg confidence: 0.87)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `82e99e11`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- New
- toUserCityResponseError
- UserLocationService
- MachineService
- net/http.Request
- pathID()
- ControllerService
- MachineModeReasonService
- api_helpers_test.go
- CncMasterSnapshot (34,440 bytes, #pragma
- LiveHandler
- Connection
- ProductionScheduleRepository
- Collector
- worker
- hub.go
- run
- mode-chart.js
- .GetDiagnostics
- net/http.Handler
- net/http.ServeMux
- CNCConnect(ip, port, dataCb, errCb, machineType) int32
- net/http.Request
- §10 Diagnostics & Troubleshooting
- decodeAndValidate
- machine_mode.go
- IdleReasonService
- BrandService
- UserSupervisorService
- cycleReading
- WithTransaction
- testing.T
- Worker (internal/monitor/worker.go)
- validateCatalogFields
- newLiveTestServerWith
- MachineService
- LicenceKeyService
- Client Certificate & Key (client_cert.der / client_key.der)
- MachineTelemetryService
- formatDuration
- CncCommunication.dll API & Integration G
- ClienteleService
- validText
- CncCommunication.dll API & Integration G
- DowntimeReasonService
- SessionService
- NewError
- .Validate
- Machines API (/machines)
- containsMachineID
- Product
- newMachineService
- BranchService
- README.md
- AuthHandler
- .validate
- MachineCycleRepository
- Category
- apidocs/main.go
- 0012. Downtime and idle reasons on a period: one table, tied to machine_mode_log
- UserProfileRepository
- Test list
- 0001. Go standard library net/http, no web framework
- 0003. PostgreSQL with pgx and hand-written SQL
- CNC Data Collection (in-process)
- TestLiveData
- UserLineService
- .validate
- run
- drawChart
- Fixed values and why they are needed
- onDataCallback(ptr uintptr) uintptr
- .validate
- Querier
- LicenceKeyDuration
- runCatalogCRUD
- 0008. "Log in as customer" uses a separate session with a path-scoped cookie
- time.Time
- AuthRepository
- LiveUnsubscribeRequest
- Event-driven data flow (no polling in client)
- run
- Live data — Socket.IO event contract
- Execute
- UserCountry
- SessionHandler
- auth_tokens SHA-256 Hashed Sessions (mig
- executeAffected
- UserCityRepository
- Create Customer in One Transaction (user
- 7th-sense-server
- LoginAsRequest
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
- ClienteleRepository
- 0005. CNC collector runs inside the API process
- SessionRepository
- HolidayService
- Break
- 0011. Ponytail plugin for Claude Code, under this repo's rules
- 0013. Machine diagnostics: reasons from what the collector recorded, network tested from the server
- UserShopService
- apidocs.go
- context.Context
- 0007. OpenAPI 3 from swag annotations, Swagger UI behind API_DOCS, generated Postman
- 0002. Layered architecture: handler, DTO, service, repository
- .validate
- .validate
- TestClienteleCRUD
- Role
- UserState
- 0006. Page-wise dashboard endpoints returning only chart fields
- .Validate
- MachineModeStatus
- UserService
- LocationRepository
- ListenForDatabaseNotifications
- SessionRevokedError
- filterTree
- containsOrgID

## God Nodes (most connected - your core abstractions)
1. `Execute()` - 210 edges
2. `SendResponse()` - 171 edges
3. `NewError()` - 151 edges
4. `run()` - 105 edges
5. `AuthUser` - 103 edges
6. `pathID()` - 88 edges
7. `currentUser()` - 76 edges
8. `Querier` - 69 edges
9. `decodeAndValidate()` - 61 edges
10. `worker` - 58 edges

## Surprising Connections (you probably didn't know these)
- `Machine Identifier Mapping (0 Simulator,1 Fanuc,2 Siemens,3 Syntec,4 Mitsubishi)` --semantically_similar_to--> `Machine Catalog Brand-Controller-Model`  [INFERRED] [semantically similar]
  docs/CncCommunication.dll — API & Integration Guide (Go Client).pdf → .claude/CLAUDE.md
- `4. Socket.IO → browser` --references--> `LiveHandler`  [INFERRED]
  docs/live-data-flow.md → internal/handler/live_handler.go
- `Context` --references--> `Authenticate()`  [INFERRED]
  docs/adr/0014-session-modes-and-login-conflict.md → internal/middleware/auth.go
- `Decision` --references--> `SendResponse()`  [INFERRED]
  docs/adr/0002-layered-handler-service-repository.md → internal/response/sendResponse.go
- `Alternatives considered` --references--> `managersOnly()`  [INFERRED]
  docs/adr/0012-machine-mode-reasons-one-table.md → internal/routes/routes.go

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

## Communities (130 total, 28 thin omitted)

### Community 0 - "New"
Cohesion: 0.14
Nodes (18): NewClienteleHandler(), NewLicenceKeyHandler(), NewMachineHandler(), NewMachineModeReasonHandler(), NewPartHandler(), NewProductionScheduleHandler(), Middleware, NewClienteleRepository() (+10 more)

### Community 1 - "toUserCityResponseError"
Cohesion: 0.27
Nodes (4): CreateUserCityRequest, UpdateUserCityRequest, validateUserCity(), toUserCityResponseError()

### Community 2 - "UserLocationService"
Cohesion: 0.16
Nodes (8): CreateUserLocationRequest, UpdateUserLocationRequest, validateUserLocation(), UserLocation, UserLocationRepository, mapUserLocationError(), UserLocationService, toUserLocationResponseError()

### Community 3 - "MachineService"
Cohesion: 0.06
Nodes (27): 0010. Machine access by an organization filter, worked out on every request, Alternatives considered, Consequences, Context, Decision, MachineCyclesResponse, CreateMachineRequest, MachineFields (+19 more)

### Community 4 - "net/http.Request"
Cohesion: 0.08
Nodes (11): BreakHandler, ClienteleHandler, DowntimeReasonHandler, HolidayHandler, IdleReasonHandler, MachineDiagnosticsHandler, MachineModeReasonHandler, PartHandler (+3 more)

### Community 5 - "pathID()"
Cohesion: 0.08
Nodes (13): MachineCycleHandler, ProductHandler, pathID(), targetUserID(), UserCellHandler, UserCityHandler, UserCountryHandler, UserLineHandler (+5 more)

### Community 6 - "ControllerService"
Cohesion: 0.10
Nodes (12): MachineModelRequest, Controller, MachineModel, ControllerRepository, mapControllerError(), MachineModelRepository, mapMachineModelError(), ControllerService (+4 more)

### Community 7 - "MachineModeReasonService"
Cohesion: 0.13
Nodes (14): CreateMachineModeReasonRequest, MachineModeReasonFields, MachineModeReasonsResponse, UpdateMachineModeReasonRequest, MachineModeReason, ReasonPeriod, MachineModeReasonRepository, mapMachineModeReasonError() (+6 more)

### Community 8 - "api_helpers_test.go"
Cohesion: 0.28
Nodes (12): cncMasterSnapshot, LanguageResponse, LiveAccessChanged, RevokeSessionsResponse, ThemeResponse, authUserKey, loginAsKey, sessionIDKey (+4 more)

### Community 9 - "CncMasterSnapshot (34,440 bytes, #pragma"
Cohesion: 0.10
Nodes (43): CncMasterSnapshot (34,440 bytes, #pragma pack(8), x86), AbsPos ([5]float64, Inches / Degrees, Work Coordinate System (WCS) positions), ActiveGCode ([16]byte, ASCII, Currently active G-code group), ActiveProgPath ([256]byte, ASCII Path, Storage path of active NC artifact), AlarmMsg ([128]byte, UTF-8 / ASCII, Alarm description message), AlarmNo (int32, Code (0 = OK), Primary alarm or fault number), AxisLoadXPercentage (int16, Percent (%), Servo X axis drive load (AxisLoadX/Y/ZPercentage)), AxisLoadYPercentage (int16, Percent (%), Servo Y axis drive load (AxisLoadX/Y/ZPercentage)) (+35 more)

### Community 10 - "LiveHandler"
Cohesion: 0.12
Nodes (11): Browser → server, liveSocket, liveSubscription, LiveSessionExpired, closeLiveSocket(), decodeLiveMessage(), LiveHandler, liveAck() (+3 more)

### Community 11 - "Connection"
Cohesion: 0.10
Nodes (15): dllDriver, slot, Step 12 — Call the DLL's CNCConnect, Step 13 — The DLL sends data, Step 2 — Load the DLL, Connection, Reading, Stats (+7 more)

### Community 12 - "ProductionScheduleRepository"
Cohesion: 0.09
Nodes (19): validDate(), CreateProductionScheduleRequest, ProductionScheduleFields, ProductionScheduleItemFields, UpdateProductionScheduleRequest, ScheduleDay(), MachineQuantity, ProductionSchedule (+11 more)

### Community 13 - "Collector"
Cohesion: 0.06
Nodes (43): Supervisor, 1. The big picture, 2. All steps at a glance, 3. Startup: steps 1 – 4, 4. Find the machines: steps 5 – 9, 6. When data is edited: step 18, 7. Shutdown: steps 19 – 20, 8. Reference (+35 more)

### Community 14 - "worker"
Cohesion: 0.07
Nodes (28): 5. Connect and collect, per machine: steps 10 – 17, Known issue, Step 10 — Get ready to connect, Step 11 — Connect, a few machines at a time, Step 14 — The data collection loop, Step 15 — Store the reading, Step 16 — Every second: dead connection check and live status, Step 17 — Disconnect and try again (+20 more)

### Community 15 - "hub.go"
Cohesion: 0.06
Nodes (45): 1. Machine → DLL → server, 2. Worker → hub (before the database), 3. Hub → rooms, 4. Socket.IO → browser, Connecting the socket, Files, How the variable list works, Live data — how a machine value reaches the screen (+37 more)

### Community 16 - "run"
Cohesion: 0.17
Nodes (22): run(), Connect(), NewBrandHandler(), NewCategoryHandler(), NewControllerHandler(), NewMachineModelHandler(), NewProductHandler(), NewUserHandler() (+14 more)

### Community 17 - "mode-chart.js"
Cohesion: 0.21
Nodes (23): chartData(), countByGroup(), drawDashboard(), drawHeader(), drawMachineList(), drawTiles(), drawTotals(), fetchData() (+15 more)

### Community 18 - ".GetDiagnostics"
Cohesion: 0.09
Nodes (34): DiagnosticCheck, DiagnosticMachine, DiagnosticVariableRow, MachineDiagnosticsResponse, MachineNetworkResponse, DiagnosticEvent, DiagnosticLiveStatus, DiagnosticMode (+26 more)

### Community 19 - "net/http.Handler"
Cohesion: 0.22
Nodes (28): CrossOrigin(), Logger(), Recoverer(), asTestUser(), uniqueSuffix(), authRequest(), createAuthTestUser(), login() (+20 more)

### Community 20 - "net/http.ServeMux"
Cohesion: 0.13
Nodes (38): Decision, RequireRole(), registerAuthRoutes(), registerBranchRoutes(), registerBrandRoutes(), registerBreakRoutes(), registerCategoryRoutes(), registerClienteleRoutes() (+30 more)

### Community 21 - "CNCConnect(ip, port, dataCb, errCb, machineType) int32"
Cohesion: 0.20
Nodes (11): CNCConnect(ip, port, dataCb, errCb, machineType) int32, CNCConnect param dataCb (uintptr stdcall OnData(const CncMasterSnapshot*)), CNCConnect param ip (const char*, ANSI NUL-terminated), CNCConnect param port (int32, e.g. 4840), CNCdisconnect(machineId int32), LazyDLL binding (syscall.NewLazyDLL + eager mod.Load in init), §4 API Reference, §5 DLL Loading & Binding (Go) (+3 more)

### Community 22 - "net/http.Request"
Cohesion: 0.08
Nodes (17): BrandHandler, CategoryHandler, ControllerHandler, dashboardQuery(), DashboardHandler, queryIDList(), readDashboardRequest(), HealthHandler (+9 more)

### Community 23 - "§10 Diagnostics & Troubleshooting"
Cohesion: 0.12
Nodes (25): Check mode: go run sinumerik_client.go check (OPC UA node existence report), Go Client (sinumerik_client.go, callback -> copy -> render), Prerequisite: 32-bit x86 DLL, compile Go with GOARCH=386, Prerequisite: SIEMENS_USER / SIEMENS_PASS env vars before CNCConnect, Run: go run sinumerik_client.go (live dashboard), §3.1 Environment & Prerequisites, §8 Running the Go Client, §10 Diagnostics & Troubleshooting (+17 more)

### Community 24 - "decodeAndValidate"
Cohesion: 0.16
Nodes (5): validator, BranchHandler, currentUserID(), decodeAndValidate(), UserHandler

### Community 25 - "machine_mode.go"
Cohesion: 0.15
Nodes (20): decideMachineMode(), holdReason(), machineStateFromReading(), numberingEngine(), programStateFromProcStatus(), ncReading(), TestDecideMachineMode(), TestDecideMachineModeConditions() (+12 more)

### Community 26 - "IdleReasonService"
Cohesion: 0.22
Nodes (7): NewIdleReasonHandler(), IdleReason, IdleReasonRepository, mapIdleReasonError(), IdleReasonService, idleReasonFromFields(), toIdleReasonResponseError()

### Community 27 - "BrandService"
Cohesion: 0.20
Nodes (5): Brand, BrandRepository, mapBrandError(), BrandService, toBrandResponseError()

### Community 28 - "UserSupervisorService"
Cohesion: 0.15
Nodes (10): CreateUserSupervisorRequest, UpdateUserSupervisorRequest, validateUserSupervisor(), NewUserSupervisorHandler(), UserSupervisor, UserSupervisorRepository, mapUserSupervisorError(), UserSupervisorService (+2 more)

### Community 29 - "cycleReading"
Cohesion: 0.16
Nodes (10): One running cycle per machine, CycleStatusFor(), MachineCycle, cycleReadingFrom(), decideCycleStep(), isProgramLoop(), resetReason(), cycleReading (+2 more)

### Community 30 - "WithTransaction"
Cohesion: 0.16
Nodes (16): clockMinutes(), CreateShiftRequest, ShiftFields, UpdateShiftRequest, MinutesBetween(), Shift, NewBreakRepository(), ShiftRepository (+8 more)

### Community 31 - "testing.T"
Cohesion: 0.09
Nodes (28): TestSnapshotModeFieldOffsets(), TestSnapshotSize(), mustDecode(), TestDashboardCharts(), TestHomeDashboard(), TestCreateBrandRequestValidate(), TestCreateCategoryRequestValidate(), TestCreateControllerRequestValidate() (+20 more)

### Community 32 - "Worker (internal/monitor/worker.go)"
Cohesion: 0.07
Nodes (22): 7th Sense Server Claude Context, response.SendResponse Format, User CRUD Golden Reference, Collector (internal/monitor/collector.go), CNC Driver (internal/cnc), Known Issue: nil live-value cache panic, machine_events, machine_telementry hypertable (+14 more)

### Community 33 - "validateCatalogFields"
Cohesion: 0.11
Nodes (11): CreateBrandRequest, UpdateBrandRequest, validateCatalogFields(), CreateCategoryRequest, UpdateCategoryRequest, CreateControllerRequest, UpdateControllerRequest, validateControllerFields() (+3 more)

### Community 34 - "newLiveTestServerWith"
Cohesion: 0.18
Nodes (18): SessionConfig, NewAuthHandler(), NewLiveHandler(), RoleHandler, NewRoleHandler(), NewSessionHandler(), Authenticate(), NewAuthRepository() (+10 more)

### Community 35 - "MachineService"
Cohesion: 0.18
Nodes (16): NewBreakHandler(), NewHolidayHandler(), NewShiftHandler(), NewHolidayRepository(), NewShiftRepository(), UserRepository, NewBreakService(), NewDowntimeReasonService() (+8 more)

### Community 36 - "LicenceKeyService"
Cohesion: 0.14
Nodes (11): NewLicenceKey, CreateLicenceKeysRequest, GenerateLicenceKeysRequest, GeneratedLicenceKeys, LicenceKey, LicenceKeyRepository, mapLicenceKeyError(), LicenceKeyService (+3 more)

### Community 37 - "Client Certificate & Key (client_cert.der / client_key.der)"
Cohesion: 0.48
Nodes (7): client_cert.der/.pem (X.509 client cert for Siemens OPC UA), client_key.der/.pem (client private key), Siemens certs: client_cert.der + client_key.der in binary directory, Siemens trust: cert must be trusted in NCU Access Administration, Troubleshoot: BadRequestNotAllowed -> Export client_cert, BadRequestNotAllowed (certificate trust rejection), Client Certificate & Key (client_cert.der / client_key.der)

### Community 38 - "MachineTelemetryService"
Cohesion: 0.14
Nodes (17): NewDashboardHandler(), NewHealthHandler(), NewMachineDiagnosticsHandler(), NewMachineCycleRepository(), NewMachineDiagnosticsRepository(), NewMachineTelemetryRepository(), NewProductionScheduleRepository(), NewDashboardService() (+9 more)

### Community 39 - "formatDuration"
Cohesion: 0.15
Nodes (30): card(), cell(), cycleData(), cycleEnd(), cycleEndReason(), cycleJob(), cycleList(), cycleNumber() (+22 more)

### Community 40 - "CncCommunication.dll API & Integration G"
Cohesion: 0.19
Nodes (20): CncCommunication.dll (C++ core runtime, x86, OPC UA inside), CncCommunication.dll API & Integration Guide (Go Client) - text version, Fanuc profile (engine 1, port 8193, Fwlib32.dll + fwlibe1.dll), Fwlib32.dll (Fanuc FOCAS2), fwlibe1.dll (Fanuc FOCAS2), libcrypto-3.dll (OpenSSL), libssl-3.dll (OpenSSL), Engine 1: Fanuc (FOCAS2 Ethernet, fwlib32) (+12 more)

### Community 41 - "ClienteleService"
Cohesion: 0.32
Nodes (7): LoginAsResponse, AuthUser, Clientele, clienteleFromFields(), ClienteleService, toClienteleResponseError(), ownCompanyCustomer()

### Community 42 - "validText"
Cohesion: 0.16
Nodes (5): UploadPartDrawingRequest, uniquePositiveIDs(), validText(), validWorkingDays(), PartDrawing

### Community 43 - "CncCommunication.dll API & Integration G"
Cohesion: 0.13
Nodes (18): CNCConnect param machineType (int32 engine id, see §4.2), Engine 0: Simulator (local loopback mock engine), Engine 2: Siemens (Sinumerik OPC UA, active profile), Engine 3: Syntec (Open API / TCP), Engine 4: Mitsubishi (driver not specified), Machine identifier mapping (engine id enum 0-4), §9 Controller-Specific Profiles, Siemens default ports: 4840 (native NCU) / 4841 (custom) (+10 more)

### Community 44 - "DowntimeReasonService"
Cohesion: 0.22
Nodes (7): NewDowntimeReasonHandler(), DowntimeReason, DowntimeReasonRepository, mapDowntimeReasonError(), downtimeReasonFromFields(), DowntimeReasonService, toDowntimeReasonResponseError()

### Community 45 - "SessionService"
Cohesion: 0.05
Nodes (26): 0014. Session modes and the login conflict, on the existing auth_tokens sessions, Alternatives considered, Consequences, Context, Decision, LoginRequest, describeDevice(), SessionClient (+18 more)

### Community 46 - "NewError"
Cohesion: 0.13
Nodes (11): CreateAddressRequest, UpdateAddressRequest, UpdateLanguageRequest, CreateUserStateRequest, validateUserState(), NewError(), newResponse(), toUserCountryResponseError() (+3 more)

### Community 47 - ".Validate"
Cohesion: 0.20
Nodes (12): UpdateThemeRequest, normalizeTheme(), TestCreateUserRequestValidate(), TestCreateUserRequestValidateDropsDuplicateMachines(), TestCreateUserRequestValidateNormalizes(), TestCustomerFiltersValidate(), TestUpdateThemeRequestValidate(), TestUpdateUserRequestValidateStatus() (+4 more)

### Community 48 - "Machines API (/machines)"
Cohesion: 0.16
Nodes (18): Clienteles API, Company Scope (companyOwner / ownCompanyCustomer), Current User API (/me theme, language), Customer Display Names (customerName / companyName), Holidays API, Licence Keys API, Machine Catalog Brand-Controller-Model, Machines API (/machines) (+10 more)

### Community 50 - "Product"
Cohesion: 0.23
Nodes (6): Product, ProductRepository, mapProductError(), ProductService, productImageOrDefault(), toProductResponseError()

### Community 51 - "newMachineService"
Cohesion: 0.20
Nodes (28): NewBranchHandler(), NewUserCityHandler(), NewUserCountryHandler(), NewUserLocationHandler(), NewUserStateHandler(), NewBranchRepository(), NewUserCityRepository(), UserCountryRepository (+20 more)

### Community 52 - "BranchService"
Cohesion: 0.15
Nodes (9): CreateBranchRequest, UpdateBranchRequest, validateBranchFields(), Branch, BranchRepository, mapBranchError(), cityOwner(), BranchService (+1 more)

### Community 53 - "README.md"
Cohesion: 0.29
Nodes (5): 0004. Opaque session token in an HttpOnly cookie, Alternatives considered, Consequences, Context, Architecture Decision Records

### Community 54 - "AuthHandler"
Cohesion: 0.09
Nodes (18): Decision, 0009. Socket.IO live data: the collector pushes each screen's variables, checked per user, Alternatives considered, Consequences, Context, Decision, AuthHandler, CurrentUser() (+10 more)

### Community 55 - ".validate"
Cohesion: 0.38
Nodes (4): ClienteleFields, CreateClienteleRequest, UpdateClienteleRequest, validClienteleEmail()

### Community 56 - "MachineCycleRepository"
Cohesion: 0.21
Nodes (3): CycleJob, MachineCycleRepository, mapMachineCycleError()

### Community 57 - "Category"
Cohesion: 0.25
Nodes (5): Category, CategoryRepository, mapCategoryError(), CategoryService, toCategoryResponseError()

### Community 58 - "apidocs/main.go"
Cohesion: 0.06
Nodes (48): request, main(), run(), main(), Config, MachineSelection, GenerateCertificates(), getEnv() (+40 more)

### Community 59 - "0012. Downtime and idle reasons on a period: one table, tied to machine_mode_log"
Cohesion: 0.33
Nodes (5): 0012. Downtime and idle reasons on a period: one table, tied to machine_mode_log, Alternatives considered, Consequences, Context, Decision

### Community 60 - "UserProfileRepository"
Cohesion: 0.15
Nodes (9): UserProfile, UserProfileRepository, mapProfileError(), NewUserProfileService(), toUserProfileResponseError(), UserProfileService, sameClock(), sameDate() (+1 more)

### Community 61 - "Test list"
Cohesion: 0.05
Nodes (38): 1. The rows are consistent (SQL, no machine needed), 2. The rows agree with the mode log (SQL), 3. The raw signals around one cycle (telemetry and logs), 4. At the machine (floor test), Before the tests (once), Check query, Completed, Floor tests on the HUMA panel (Fanuc) (+30 more)

### Community 62 - "0001. Go standard library net/http, no web framework"
Cohesion: 0.29
Nodes (6): 0001. Go standard library net/http, no web framework, Alternatives considered, Consequences, Context, Decision, Chain()

### Community 63 - "0003. PostgreSQL with pgx and hand-written SQL"
Cohesion: 0.33
Nodes (5): 0003. PostgreSQL with pgx and hand-written SQL, Alternatives considered, Consequences, Context, Decision

### Community 64 - "CNC Data Collection (in-process)"
Cohesion: 0.15
Nodes (12): CNC Data Collection (in-process), writeLog CNC Log Types and Switches, .env Configuration (internal/config), Health Mode History API (public), CncCommunication.dll (32-bit vendor DLL), machine_mode_log, Mode Decision (machine_mode.go), ECharts 5.5.1 (+4 more)

### Community 65 - "TestLiveData"
Cohesion: 0.25
Nodes (6): dialLive(), newLiveTestServer(), TestLiveData(), liveTestClient, liveTestEvent, liveTestServer

### Community 66 - "UserLineService"
Cohesion: 0.07
Nodes (18): CreateUserCellRequest, UpdateUserCellRequest, validateUserCell(), CreateUserLineRequest, UpdateUserLineRequest, validateUserLine(), NewUserCellHandler(), UserCell (+10 more)

### Community 67 - ".validate"
Cohesion: 0.47
Nodes (3): BreakFields, CreateBreakRequest, UpdateBreakRequest

### Community 68 - "run"
Cohesion: 0.36
Nodes (4): connectFanucCounter(), main(), run(), fanucCounter

### Community 69 - "drawChart"
Cohesion: 0.24
Nodes (11): chartRows(), currentSpan(), currentWindow(), drawChart(), drawZoomLevels(), markZoomLevel(), renderBar(), renderCycleBar() (+3 more)

### Community 70 - "Fixed values and why they are needed"
Cohesion: 0.29
Nodes (7): `CNC_CYCLE_MAX_HOURS` (24 hours), `CNC_STALE_SEC` (30 s) and `CNC_QUIET_SEC` (300 s), `cycleCounterGrace` and `cycleStartGrace` (5 seconds each), `emergencyResetGrace` (2 seconds), Fixed values and why they are needed, `loopMinimumDropFraction` (0.5), `maxPartCountStep` (1)

### Community 71 - "onDataCallback(ptr uintptr) uintptr"
Cohesion: 0.36
Nodes (9): Package-level callback singletons (syscall.NewCallback), CNCConnect param errCb (uintptr stdcall OnErr(int32 code, const char* msg)), Mutex-protected copy: dataMutex sync.RWMutex -> lastRecord/hasRecord, onDataCallback(ptr uintptr) uintptr, onErrorCallback(code, messagePtr uintptr) uintptr, §7 Asynchronous Callback Architecture, onDataCallback, onErrorCallback (+1 more)

### Community 72 - ".validate"
Cohesion: 0.47
Nodes (3): CreateHolidayRequest, HolidayFields, UpdateHolidayRequest

### Community 73 - "Querier"
Cohesion: 0.23
Nodes (14): Querier, NewAddressRepository(), NewCompanyProfileRepository(), NewDowntimeReasonRepository(), NewIdleReasonRepository(), NewMachineModeReasonRepository(), NewUserCellRepository(), NewUserMachineAccessRepository() (+6 more)

### Community 75 - "runCatalogCRUD"
Cohesion: 0.23
Nodes (15): containsID(), newCatalogRouter(), runCatalogCRUD(), createTestRecord(), TestBrandCRUD(), TestCategoryCRUD(), TestControllerBrand(), TestControllerCRUD() (+7 more)

### Community 76 - "0008. "Log in as customer" uses a separate session with a path-scoped cookie"
Cohesion: 0.40
Nodes (4): 0008. "Log in as customer" uses a separate session with a path-scoped cookie, Alternatives considered, Consequences, Context

### Community 77 - "time.Time"
Cohesion: 0.04
Nodes (74): Decision, FinishCyclesSummary, FinishedCycle, HomeAlarm, HomeCycles, HomeFleet, HomeMachine, HomePlan (+66 more)

### Community 78 - "AuthRepository"
Cohesion: 0.18
Nodes (4): SessionState, SessionUser, AuthRepository, mapAuthError()

### Community 80 - "Event-driven data flow (no polling in client)"
Cohesion: 0.60
Nodes (5): Flow 1: load DLL, register two stdcall callbacks (syscall.NewCallback), Flow 2: CNCConnect starts DLL background worker thread, Flow 3: worker polls controller, assembles snapshot, invokes data callback, Flow 4: callback copies payload, signals renderer via buffered channel, Flow 5: main goroutine renders live ANSI dashboard

### Community 81 - "run"
Cohesion: 0.36
Nodes (10): checkNotTaken(), createSuperadmin(), ensureProfile(), execute(), main(), parseFlags(), run(), superadminRoleID() (+2 more)

### Community 82 - "Live data — Socket.IO event contract"
Cohesion: 0.20
Nodes (7): Adding a variable to a screen, Connecting, Live data — Socket.IO event contract, Ordering, Rooms (server side), Server → browser, Settings

### Community 83 - "Execute"
Cohesion: 0.13
Nodes (4): Execute(), User, MachineDiagnosticsRepository, mapUserError()

### Community 84 - "UserCountry"
Cohesion: 0.22
Nodes (3): CreateUserCountryRequest, UserCountry, mapUserCountryError()

### Community 86 - "auth_tokens SHA-256 Hashed Sessions (mig"
Cohesion: 0.67
Nodes (3): auth_tokens SHA-256 Hashed Sessions (migration 014), Per-Request Session and User Re-check, ts_session Cookie (HttpOnly, SameSite=Strict)

### Community 88 - "executeAffected"
Cohesion: 0.09
Nodes (19): PartProcessFields, PartSubProgramFields, WithTransaction(), CreatePartRequest, PartFields, UpdatePartRequest, ValidPartDrawingFile(), Part (+11 more)

### Community 89 - "UserCityRepository"
Cohesion: 0.31
Nodes (3): UserCity, UserCityRepository, mapUserCityError()

### Community 105 - "0005. CNC collector runs inside the API process"
Cohesion: 0.33
Nodes (5): 0005. CNC collector runs inside the API process, Alternatives considered, Consequences, Context, Decision

### Community 107 - "HolidayService"
Cohesion: 0.24
Nodes (6): Holiday, HolidayRepository, mapHolidayError(), HolidayService, holidayFromFields(), toHolidayResponseError()

### Community 108 - "Break"
Cohesion: 0.16
Nodes (8): Break, ShiftBreak, BreakRepository, mapBreakError(), breakFromFields(), BreakService, toBreakResponseError(), companyOwner()

### Community 109 - "0011. Ponytail plugin for Claude Code, under this repo's rules"
Cohesion: 0.33
Nodes (5): 0011. Ponytail plugin for Claude Code, under this repo's rules, Alternatives considered, Consequences, Context, Decision

### Community 110 - "0013. Machine diagnostics: reasons from what the collector recorded, network tested from the server"
Cohesion: 0.33
Nodes (5): 0013. Machine diagnostics: reasons from what the collector recorded, network tested from the server, Alternatives considered, Consequences, Context, Decision

### Community 111 - "UserShopService"
Cohesion: 0.16
Nodes (8): CreateUserShopRequest, UpdateUserShopRequest, validateUserShop(), UserShop, UserShopRepository, mapUserShopError(), UserShopService, toUserShopResponseError()

### Community 113 - "context.Context"
Cohesion: 0.10
Nodes (9): Step 5 — Prepare the database, LiveValue, MachineEvent, MachineLiveStatus, MachineTelemetryRepository, mapMachineTelemetryError(), execute(), UserMachineAccessRepository (+1 more)

### Community 114 - "0007. OpenAPI 3 from swag annotations, Swagger UI behind API_DOCS, generated Postman"
Cohesion: 0.33
Nodes (5): 0007. OpenAPI 3 from swag annotations, Swagger UI behind API_DOCS, generated Postman, Alternatives considered, Consequences, Context, Decision

### Community 115 - "0002. Layered architecture: handler, DTO, service, repository"
Cohesion: 0.33
Nodes (5): 0002. Layered architecture: handler, DTO, service, repository, Alternatives considered, Consequences, Context, Decision

### Community 116 - ".validate"
Cohesion: 0.47
Nodes (3): CreateDowntimeReasonRequest, DowntimeReasonFields, UpdateDowntimeReasonRequest

### Community 117 - ".validate"
Cohesion: 0.47
Nodes (3): CreateIdleReasonRequest, IdleReasonFields, UpdateIdleReasonRequest

### Community 118 - "TestClienteleCRUD"
Cohesion: 0.14
Nodes (53): WithAuthUser(), decodeData(), expectStatus(), sendRequest(), apiResponse, clienteleBody(), TestClienteleCRUD(), generateLicenceKeys() (+45 more)

### Community 119 - "Role"
Cohesion: 0.33
Nodes (3): Role, mapRoleError(), toRoleResponseError()

### Community 121 - "0006. Page-wise dashboard endpoints returning only chart fields"
Cohesion: 0.40
Nodes (4): 0006. Page-wise dashboard endpoints returning only chart fields, Alternatives considered, Consequences, Context

### Community 122 - ".Validate"
Cohesion: 0.40
Nodes (3): UpdateLicenceKeyRequest, validDuration(), ValidLicenceKey()

### Community 124 - "UserService"
Cohesion: 0.05
Nodes (41): CreateUserRequest, CustomerDetails, MachineFilter, UpdateUserRequest, UserListFilter, UserResponse, isTimeZone(), isWebsite() (+33 more)

### Community 125 - "LocationRepository"
Cohesion: 0.50
Nodes (3): mapLocationError(), NewLocationRepository(), LocationRepository

## Ambiguous Edges - Review These
- `FanucCncCommunication.dll` → `CncCommunication.dll (32-bit vendor DLL)`  [AMBIGUOUS]
  docs/setup-guide.txt · relation: semantically_similar_to
- `Snapshot Fields (positions, spindle, feed, program, tool, alarms)` → `Machines API (/machines)`  [AMBIGUOUS]
  docs/CncCommunication.dll — API & Integration Guide (Go Client).pdf · relation: conceptually_related_to

## Knowledge Gaps
- **35 isolated node(s):** `7th-sense-server`, `middleware.RequireRole / routes.adminOnly`, `response.SendResponse`, `Test Auth (asTestUser, tests/auth_test.go)`, `user_filters Filter Checklist` (+30 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 210 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **28 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `CNC Data Collection Flow` connect `Collector` to `CNC Data Collection (in-process)`, `Live data — Socket.IO event contract`, `Worker (internal/monitor/worker.go)`, `worker`?**
  _High betweenness centrality (0.169) - this node is a cross-community bridge._
- **What connects `7th-sense-server`, `middleware.RequireRole / routes.adminOnly`, `response.SendResponse` to the rest of the system?**
  _35 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `New` be split into smaller, more focused modules?**
  _Cohesion score 0.1437908496732026 - nodes in this community are weakly interconnected._
- **Why does `2. All steps at a glance` connect `Collector` to `Connection`, `worker`?**
  _High betweenness centrality (0.141) - this node is a cross-community bridge._
- **Should `MachineService` be split into smaller, more focused modules?**
  _Cohesion score 0.0649895178197065 - nodes in this community are weakly interconnected._

### Low-confidence Hints
_AMBIGUOUS edges — the extractor was unsure. Verify before acting on these._

- **What is the exact relationship between `FanucCncCommunication.dll` and `CncCommunication.dll (32-bit vendor DLL)`?**
  _Edge tagged AMBIGUOUS (relation: semantically_similar_to) - confidence is low._
- **What is the exact relationship between `Snapshot Fields (positions, spindle, feed, program, tool, alarms)` and `Machines API (/machines)`?**
  _Edge tagged AMBIGUOUS (relation: conceptually_related_to) - confidence is low._