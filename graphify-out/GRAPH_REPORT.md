# Graph Report - 7th-sense-server  (2026-11-24)

## Corpus Check
- 390 files · ~336,471 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 4 file(s) not represented in the graph (top: (none) 3, .toml 1)

## Summary
- 3495 nodes · 32612 edges · 134 communities (106 shown, 28 thin omitted)
- Extraction: 95% EXTRACTED · 5% INFERRED · 0% AMBIGUOUS · INFERRED: 1475 edges (avg confidence: 0.87)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `9073efdc`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- HomeDashboardResponse
- UserCityRepository
- telemetry_load_test.go
- ControllerRepository
- .Validate
- 0014. Session modes and the login conflict, on the existing auth_tokens sessions
- apidocs/main.go
- MachineModeReasonService
- api_helpers_test.go
- CncMasterSnapshot (34,440 bytes, #pragma
- AuthHandler
- MachineModelRepository
- UserCountryRepository
- MachineModeStatus
- worker
- Hub
- pathID
- mode-chart.js
- UserLocationService
- net/http.Request
- net/http.ServeMux
- PlcCollector
- SendResponse
- §10 Diagnostics & Troubleshooting
- OperatorTaskService
- machine_mode.go
- MachineTelemetryService
- Product
- plcWorker
- github.com/jackc/pgx/v5/pgxpool.Pool
- TelemetryQueue
- SessionConflictError
- Worker (internal/monitor/worker.go)
- 0018. Machine Health X / Y / Z charts: downsampled in SQL, timed, on the machine-health routes
- run
- testing.T
- BrandRepository
- logger/logger.go
- 0016. A Task List screen for the operator tasks' master list
- MachineHealthRepository
- CncCommunication.dll API & Integration G
- AuthService
- TestLiveData
- CncCommunication.dll API & Integration G
- UserShopService
- config.go
- UserLineRepository
- NewError
- Machines API (/machines)
- LiveValue
- Address
- CompanyProfile
- BranchService
- UserRepository
- time.Time
- Execute
- MachineCycle
- UserService
- DashboardRequest
- PartService
- UserProfileRepository
- Test list
- 0001. Go standard library net/http, no web framework
- 0003. PostgreSQL with pgx and hand-written SQL
- dto_test.go
- VirtualScreenResponse
- OpenModePeriod
- VariableConfig
- README.md
- LiveHandler
- WithTransaction
- onDataCallback(ptr uintptr) uintptr
- DashboardService
- UserMachineFilter
- Decision
- 0013. Machine diagnostics: reasons from what the collector recorded, network tested from the server
- 0008. "Log in as customer" uses a separate session with a path-scoped cookie
- 0020. Machine PLC node IDs: the machine's own node ID in machine_variables, default otherwise
- MonitoredMachine
- run
- Event-driven data flow (no polling in client)
- statusRecorder
- Role
- ProductionScheduleRepository
- 0004. Opaque session token in an HttpOnly cookie
- MachineTelemetryRepository
- auth_tokens SHA-256 Hashed Sessions (mig
- AuthUser
- Live data — Socket.IO event contract
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
- connection.go
- 0005. CNC collector runs inside the API process
- 0006. Page-wise dashboard endpoints returning only chart fields
- 0011. Ponytail plugin for Claude Code, under this repo's rules
- ClienteleService
- 0010. Machine access by an organization filter, worked out on every request
- 0022. Test PLC node IDs in a short OPC UA session of their own
- containsID
- apidocs.go
- 0025. A machine is saved at one level of its customer's organization; the path above is read, not stored
- 0007. OpenAPI 3 from swag annotations, Swagger UI behind API_DOCS, generated Postman
- 0002. Layered architecture: handler, DTO, service, repository
- HomeTime
- context.Context
- UserStateService
- 0017. Machine accuracy as one shared list, managed by the superadmin
- LocationRepository
- 0015. Operator tasks: a shared task list, assigned per customer and machine model
- dll_windows_386_test.go
- MachineCycleHandler
- LicenceKeyDuration
- LiveUnsubscribeRequest
- Connection
- SessionRevokedError
- MachineService
- filterTree
- LicenceKeyService
- .GetDiagnostics
- pathID
- TestPLCCollector

## God Nodes (most connected - your core abstractions)
1. `Execute()` - 228 edges
2. `SendResponse()` - 192 edges
3. `NewError()` - 162 edges
4. `run()` - 120 edges
5. `AuthUser` - 114 edges
6. `pathID()` - 99 edges
7. `currentUser()` - 86 edges
8. `Querier` - 80 edges
9. `decodeAndValidate()` - 69 edges
10. `worker` - 61 edges

## Surprising Connections (you probably didn't know these)
- `Machine Identifier Mapping (0 Simulator,1 Fanuc,2 Siemens,3 Syntec,4 Mitsubishi)` --semantically_similar_to--> `Machine Catalog Brand-Controller-Model`  [INFERRED] [semantically similar]
  docs/CncCommunication.dll — API & Integration Guide (Go Client).pdf → .claude/CLAUDE.md
- `Step 7: Wiring, start and stop order` --references--> `Supervisor`  [INFERRED]
  docs/rabbitmq-telemetry-queue-guide.md → internal/collector/supervisor.go
- `Step 1: ADR and configuration` --references--> `RabbitMQConfig`  [INFERRED]
  docs/rabbitmq-telemetry-queue-guide.md → internal/config/config.go
- `4. Socket.IO → browser` --references--> `LiveHandler`  [INFERRED]
  docs/live-data-flow.md → internal/handler/live_handler.go
- `Context` --references--> `Authenticate()`  [INFERRED]
  docs/adr/0014-session-modes-and-login-conflict.md → internal/middleware/auth.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **CNC collection pipeline: Supervisor -> Collector -> Worker -> Driver -> DLL** — docs_cnc_data_collection_flow_supervisor, docs_cnc_data_collection_flow_collector, docs_cnc_data_collection_flow_worker, docs_cnc_data_collection_flow_driver, docs_cnc_data_collection_flow_cnc_communication_dll [EXTRACTED 1.00]
- **Event-driven telemetry pipeline (CNCConnect -> worker -> callback -> mutex copy -> render)** — docs_cnc_communication_cncconnect, docs_cnc_communication_flow_step1_register_callbacks, docs_cnc_communication_flow_step2_worker_thread, docs_cnc_communication_flow_step3_poll_and_invoke, docs_cnc_communication_flow_step4_copy_and_signal, docs_cnc_communication_flow_step5_render_dashboard, docs_cnc_communication_ondatacallback, docs_cnc_communication_mutex_copy, docs_cnc_communication_cncmastersnapshot [EXTRACTED 1.00]
- **Siemens OPC UA connection requirements (GOARCH=386, credentials, trusted certs)** — docs_cnc_communication_prereq_goarch_386, docs_cnc_communication_prereq_siemens_credentials, docs_cnc_communication_siemens_certificates, docs_cnc_communication_siemens_trust, docs_cnc_communication_client_cert_der, docs_cnc_communication_client_key_der, docs_cnc_communication_siemens_sinumerik_profile [EXTRACTED 1.00]
- **CNC telemetry callback data flow** — docs_cnccommunication_dll_api_integration_guide_go_client_cncconnect, docs_cnccommunication_dll_api_integration_guide_go_client_stdcall_callbacks, docs_cnccommunication_dll_api_integration_guide_go_client_ondatacallback, docs_cnccommunication_dll_api_integration_guide_go_client_cncmastersnapshot, docs_cnccommunication_dll_api_integration_guide_go_client_event_driven_data_flow [EXTRACTED 1.00]
- **Company-scoped customer entities sharing companyOwner access** — _claude_claude_company_scope, _claude_claude_parts_api, _claude_claude_clienteles_api, _claude_claude_shifts_breaks_api, _claude_claude_holidays_api, _claude_claude_reasons_api [EXTRACTED 1.00]
- **Struct ABI compatibility (pack(8) x86, size 34440, GOARCH=386)** — docs_cnc_communication_cncmastersnapshot, docs_cnc_communication_struct_size_assertion, docs_cnc_communication_prereq_goarch_386, docs_cnc_communication_ts_memory_panic [INFERRED 0.85]
- **Session authentication and authorization flow** — readme_session_cookie, readme_auth_tokens_hashed, readme_per_request_session_check, readme_login_lockout, readme_csrf_cross_origin, _claude_claude_role_guards [INFERRED 0.85]

## Communities (134 total, 28 thin omitted)

### Community 0 - "HomeDashboardResponse"
Cohesion: 0.10
Nodes (19): HomeAlarm, HomeCycles, HomeFleet, HomeMachine, HomePlan, HomeDashboardResponse, HomeHour, HomeModeCount (+11 more)

### Community 1 - "UserCityRepository"
Cohesion: 0.16
Nodes (8): CreateUserCityRequest, UpdateUserCityRequest, validateUserCity(), UserCity, UserCityRepository, mapUserCityError(), UserCityService, toUserCityResponseError()

### Community 2 - "telemetry_load_test.go"
Cohesion: 0.06
Nodes (55): RabbitMQConfig, NewTelemetryQueue(), newTestWriter(), runWriter(), testDelivery(), TestWriterAcksOnlyAfterSuccess(), TestWriterDeadLettersRefusedAndUndecodableReadings(), TestWriterFlushesBySize() (+47 more)

### Community 3 - "ControllerRepository"
Cohesion: 0.15
Nodes (8): CreateControllerRequest, UpdateControllerRequest, validateControllerFields(), Controller, ControllerRepository, mapControllerError(), ControllerService, toControllerResponseError()

### Community 4 - ".Validate"
Cohesion: 0.19
Nodes (13): UpdateThemeRequest, MachineFilter, normalizeTheme(), TestCreateUserRequestValidate(), TestCreateUserRequestValidateDropsDuplicateMachines(), TestCreateUserRequestValidateNormalizes(), TestCustomerFiltersValidate(), TestUpdateThemeRequestValidate() (+5 more)

### Community 5 - "0014. Session modes and the login conflict, on the existing auth_tokens sessions"
Cohesion: 0.29
Nodes (6): 0014. Session modes and the login conflict, on the existing auth_tokens sessions, Alternatives considered, Consequences, Context, Decision, ParseUserAgent()

### Community 6 - "apidocs/main.go"
Cohesion: 0.09
Nodes (35): request, EnsureCertificates(), existingCertificateExpiry(), generateCertificates(), Paths, writePEM(), Settings, addExamples() (+27 more)

### Community 7 - "MachineModeReasonService"
Cohesion: 0.06
Nodes (30): Decision, 0012. Downtime and idle reasons on a period: one table, tied to machine_mode_log, Alternatives considered, Consequences, Context, Decision, LoginAsResponse, CreateMachineModeReasonRequest (+22 more)

### Community 8 - "api_helpers_test.go"
Cohesion: 0.31
Nodes (13): cncMasterSnapshot, LanguageResponse, LiveAccessChanged, MachineLocationFiltersResponse, RevokeSessionsResponse, ThemeResponse, authUserKey, loginAsKey (+5 more)

### Community 9 - "CncMasterSnapshot (34,440 bytes, #pragma"
Cohesion: 0.10
Nodes (42): CncMasterSnapshot (34,440 bytes, #pragma pack(8), x86), AbsPos ([5]float64, Inches / Degrees, Work Coordinate System (WCS) positions), ActiveGCode ([16]byte, ASCII, Currently active G-code group), ActiveProgPath ([256]byte, ASCII Path, Storage path of active NC artifact), AlarmMsg ([128]byte, UTF-8 / ASCII, Alarm description message), AlarmNo (int32, Code (0 = OK), Primary alarm or fault number), AxisLoadXPercentage (int16, Percent (%), Servo X axis drive load (AxisLoadX/Y/ZPercentage)), AxisLoadYPercentage (int16, Percent (%), Servo Y axis drive load (AxisLoadX/Y/ZPercentage)) (+34 more)

### Community 10 - "AuthHandler"
Cohesion: 0.15
Nodes (12): Decision, AuthHandler, CurrentUser(), isLoginAsKey(), LoginAsKey(), RouteLoginAs(), SessionToken(), UserID() (+4 more)

### Community 11 - "MachineModelRepository"
Cohesion: 0.20
Nodes (6): MachineModelRequest, MachineModel, MachineModelRepository, mapMachineModelError(), MachineModelService, toMachineModelResponseError()

### Community 12 - "UserCountryRepository"
Cohesion: 0.15
Nodes (7): CreateUserCountryRequest, UserCountry, UserCountryRepository, mapUserCountryError(), UserCountryService, toUserCountryResponseError(), levelInMachineFilterError()

### Community 14 - "worker"
Cohesion: 0.07
Nodes (24): Step 10 — Get ready to connect, Step 14 — The data collection loop, Step 16 — Every second: dead connection check and live status, `cycleCheckpointInterval` (5 minutes), Reading, Status, Value, FieldIndex() (+16 more)

### Community 15 - "Hub"
Cohesion: 0.07
Nodes (41): 1. Machine → DLL → server, 2. Worker → hub (before the database), 3. Hub → rooms, 4. Socket.IO → browser, Connecting the socket, Files, How the variable list works, Live data — how a machine value reaches the screen (+33 more)

### Community 16 - "pathID"
Cohesion: 0.07
Nodes (15): BreakHandler, ClienteleHandler, DowntimeReasonHandler, NewDowntimeReasonHandler(), HolidayHandler, MachineDiagnosticsHandler, MachineModeReasonHandler, MachinePlcVariableHandler (+7 more)

### Community 17 - "mode-chart.js"
Cohesion: 0.10
Nodes (64): card(), cell(), chartData(), chartRows(), countByGroup(), currentSpan(), currentWindow(), cycleData() (+56 more)

### Community 18 - "UserLocationService"
Cohesion: 0.07
Nodes (18): CreateUserLocationRequest, UpdateUserLocationRequest, validateUserLocation(), CreateUserSupervisorRequest, UpdateUserSupervisorRequest, validateUserSupervisor(), NewUserSupervisorHandler(), UserLocation (+10 more)

### Community 19 - "net/http.Request"
Cohesion: 0.10
Nodes (12): axisValuesQuery(), dashboardQuery(), DashboardHandler, queryIDList(), readDashboardRequest(), IdleReasonHandler, isSameOrigin(), MachineHandler (+4 more)

### Community 20 - "net/http.ServeMux"
Cohesion: 0.12
Nodes (41): Decision, RequireRole(), registerAuthRoutes(), registerBranchRoutes(), registerBrandRoutes(), registerBreakRoutes(), registerCategoryRoutes(), registerClienteleRoutes() (+33 more)

### Community 21 - "PlcCollector"
Cohesion: 0.30
Nodes (4): samePLCConnection(), newPlcWorker(), sleepUnlessCancelled(), PlcCollector

### Community 22 - "SendResponse"
Cohesion: 0.06
Nodes (21): validator, BranchHandler, BrandHandler, CategoryHandler, ControllerHandler, HealthHandler, LicenceKeyHandler, MachineAccuracyHandler (+13 more)

### Community 23 - "§10 Diagnostics & Troubleshooting"
Cohesion: 0.12
Nodes (27): CNCdisconnect(machineId int32), CncCommunication.dll API & Integration Guide (Go Client) - text version, LazyDLL binding (syscall.NewLazyDLL + eager mod.Load in init), Prerequisite: 32-bit x86 DLL, compile Go with GOARCH=386, §4 API Reference, §5 DLL Loading & Binding (Go), §1 Overview, §6 Snapshot Data Structure (+19 more)

### Community 24 - "OperatorTaskService"
Cohesion: 0.09
Nodes (15): CreateMaintenanceTaskRequest, CreateOperatorTaskRequest, OperatorTaskFields, UpdateMaintenanceTaskRequest, UpdateOperatorTaskRequest, MaintenanceTask, OperatorTask, MaintenanceTaskRepository (+7 more)

### Community 25 - "machine_mode.go"
Cohesion: 0.15
Nodes (20): decideMachineMode(), holdReason(), machineStateFromReading(), numberingEngine(), programStateFromProcStatus(), ncReading(), TestDecideMachineMode(), TestDecideMachineModeConditions() (+12 more)

### Community 26 - "MachineTelemetryService"
Cohesion: 0.21
Nodes (3): Step 5 — Prepare the database, MachineTelemetryService, toMachineTelemetryResponseError()

### Community 27 - "Product"
Cohesion: 0.16
Nodes (9): CreateProductRequest, UpdateProductRequest, validateProductFields(), Product, ProductRepository, mapProductError(), ProductService, productImageOrDefault() (+1 more)

### Community 28 - "plcWorker"
Cohesion: 0.15
Nodes (9): 0023. PLC reader: short session timeout, and read the nodes the controller has no room to monitor, Alternatives considered, Consequences, Context, Decision, previousText(), Endpoint(), backoff (+1 more)

### Community 29 - "github.com/jackc/pgx/v5/pgxpool.Pool"
Cohesion: 0.07
Nodes (31): BreakFields, CreateBreakRequest, UpdateBreakRequest, uniquePositiveIDs(), clockMinutes(), CreateShiftRequest, ShiftFields, UpdateShiftRequest (+23 more)

### Community 30 - "TelemetryQueue"
Cohesion: 0.10
Nodes (13): TelemetryReading, TelemetryQueue, decodeReading(), TestWriterKeepsNewestLiveValue(), roundToTenth(), ReadingValue, batchActions, batchEntry (+5 more)

### Community 32 - "Worker (internal/monitor/worker.go)"
Cohesion: 0.06
Nodes (30): CNC Data Collection (in-process), writeLog CNC Log Types and Switches, .env Configuration (internal/config), Health Mode History API (public), 7th Sense Server Claude Context, response.SendResponse Format, User CRUD Golden Reference, CncCommunication.dll (32-bit vendor DLL) (+22 more)

### Community 33 - "0018. Machine Health X / Y / Z charts: downsampled in SQL, timed, on the machine-health routes"
Cohesion: 0.18
Nodes (9): 0018. Machine Health X / Y / Z charts: downsampled in SQL, timed, on the machine-health routes, Alternatives considered, Consequences, Context, 0019. Axis values: every stored value when it fits, and zoom by loading the window again, Alternatives considered, Consequences, Context (+1 more)

### Community 34 - "run"
Cohesion: 0.04
Nodes (128): run(), SessionConfig, Querier, NewAuthHandler(), NewBranchHandler(), NewBrandHandler(), NewBreakHandler(), NewCategoryHandler() (+120 more)

### Community 35 - "testing.T"
Cohesion: 0.07
Nodes (124): CreateUserRequest, TestParseDayStart(), NewMachineHandler(), WithAuthUser(), Logger(), asTestUser(), decodeData(), expectStatus() (+116 more)

### Community 36 - "BrandRepository"
Cohesion: 0.08
Nodes (15): CreateBrandRequest, UpdateBrandRequest, validateCatalogFields(), CreateCategoryRequest, UpdateCategoryRequest, Brand, Category, BrandRepository (+7 more)

### Community 37 - "logger/logger.go"
Cohesion: 0.40
Nodes (3): main(), Init(), InitFile()

### Community 38 - "0016. A Task List screen for the operator tasks' master list"
Cohesion: 0.40
Nodes (5): 0016. A Task List screen for the operator tasks' master list, Alternatives considered, Consequences, Context, Decision

### Community 39 - "MachineHealthRepository"
Cohesion: 0.12
Nodes (15): Decision, AxisPoint, FindAxisSet(), AxisSeries, AxisSet, AxisStats, AxisValuesRequest, AxisValuesResponse (+7 more)

### Community 40 - "CncCommunication.dll API & Integration G"
Cohesion: 0.27
Nodes (15): CncCommunication.dll (C++ core runtime, x86, OPC UA inside), Fanuc profile (engine 1, port 8193, Fwlib32.dll + fwlibe1.dll), Fwlib32.dll (Fanuc FOCAS2), fwlibe1.dll (Fanuc FOCAS2), libcrypto-3.dll (OpenSSL), libssl-3.dll (OpenSSL), Engine 1: Fanuc (FOCAS2 Ethernet, fwlib32), open62541.dll (OPC UA runtime) (+7 more)

### Community 41 - "AuthService"
Cohesion: 0.13
Nodes (11): LoginRequest, SessionClient, hashToken(), newLoginAsKey(), newSessionToken(), sessionIDs(), truncate(), newLoginThrottle() (+3 more)

### Community 42 - "TestLiveData"
Cohesion: 0.25
Nodes (7): dialLive(), newLiveTestServer(), TestLiveData(), liveTestClient, liveTestEvent, liveTestServer, TestSessionRevokeClosesSockets()

### Community 43 - "CncCommunication.dll API & Integration G"
Cohesion: 0.11
Nodes (26): client_cert.der/.pem (X.509 client cert for Siemens OPC UA), client_key.der/.pem (client private key), CNCConnect param machineType (int32 engine id, see §4.2), Engine 0: Simulator (local loopback mock engine), Engine 2: Siemens (Sinumerik OPC UA, active profile), Engine 3: Syntec (Open API / TCP), Engine 4: Mitsubishi (driver not specified), Machine identifier mapping (engine id enum 0-4) (+18 more)

### Community 44 - "UserShopService"
Cohesion: 0.15
Nodes (8): CreateUserShopRequest, UpdateUserShopRequest, validateUserShop(), UserShop, UserShopRepository, mapUserShopError(), UserShopService, toUserShopResponseError()

### Community 45 - "config.go"
Cohesion: 0.07
Nodes (35): main(), run(), Supervisor, MachineSelection, 0024. Telemetry rows through a RabbitMQ queue, written in batches, Alternatives considered, Consequences, Context (+27 more)

### Community 46 - "UserLineRepository"
Cohesion: 0.07
Nodes (18): CreateUserCellRequest, UpdateUserCellRequest, validateUserCell(), CreateUserLineRequest, UpdateUserLineRequest, validateUserLine(), NewUserCellHandler(), UserCell (+10 more)

### Community 47 - "NewError"
Cohesion: 0.11
Nodes (14): CreateAddressRequest, UpdateAddressRequest, UpdateLanguageRequest, UpdateLicenceKeyRequest, validDuration(), ValidLicenceKey(), isTimeZone(), isWebsite() (+6 more)

### Community 48 - "Machines API (/machines)"
Cohesion: 0.17
Nodes (17): Clienteles API, Company Scope (companyOwner / ownCompanyCustomer), Current User API (/me theme, language), Customer Display Names (customerName / companyName), Holidays API, Licence Keys API, Machine Catalog Brand-Controller-Model, Machines API (/machines) (+9 more)

### Community 49 - "LiveValue"
Cohesion: 0.09
Nodes (27): 0026. Telemetry writes: one asynchronous commit per batch, one index less, PostgreSQL settings, Alternatives considered, Consequences, Context, Measurements (this PC, PostgreSQL alone, 200 batches × 5,000 values, insert + upsert), Step 15 — Store the reading, 1. Why, and what today looks like, 2. Scope (+19 more)

### Community 50 - "Address"
Cohesion: 0.27
Nodes (5): Address, AddressRepository, mapAddressError(), toAddressResponseError(), AddressService

### Community 51 - "CompanyProfile"
Cohesion: 0.28
Nodes (5): CompanyProfile, CompanyProfileRepository, mapCompanyProfileError(), toCompanyProfileResponseError(), CompanyProfileService

### Community 52 - "BranchService"
Cohesion: 0.14
Nodes (9): CreateBranchRequest, UpdateBranchRequest, validateBranchFields(), Branch, BranchRepository, mapBranchError(), cityOwner(), BranchService (+1 more)

### Community 53 - "UserRepository"
Cohesion: 0.31
Nodes (4): Decision, User, UserRepository, mapUserError()

### Community 54 - "time.Time"
Cohesion: 0.11
Nodes (21): PlcNodeTestResult, TimelinePeriod, ModeTimelineResponse, ModeTotalsResponse, NewTimelinePeriod(), MachineModeHistoryResponse, MachineModeTotal, AxisBucketSeconds() (+13 more)

### Community 55 - "Execute"
Cohesion: 0.24
Nodes (5): Execute(), SessionState, SessionUser, AuthRepository, mapAuthError()

### Community 56 - "MachineCycle"
Cohesion: 0.15
Nodes (8): MachineCyclesResponse, CycleJob, MachineCycle, MachineCycleRepository, mapMachineCycleError(), MachineCycleService, sameID(), toMachineCycleResponseError()

### Community 57 - "UserService"
Cohesion: 0.21
Nodes (5): AssignableRoles(), CanAssignRole(), canSee(), UserService, toUserResponseError()

### Community 58 - "DashboardRequest"
Cohesion: 0.33
Nodes (3): DashboardRequest, HomeDashboardRequest, MachineCycleListRequest

### Community 59 - "PartService"
Cohesion: 0.05
Nodes (29): Decision, PartProcessFields, PartSubProgramFields, MachineAccuracyRequest, CreatePartRequest, PartFields, UpdatePartRequest, UploadPartDrawingRequest (+21 more)

### Community 60 - "UserProfileRepository"
Cohesion: 0.19
Nodes (5): UserProfile, UserProfileRepository, mapProfileError(), toUserProfileResponseError(), UserProfileService

### Community 61 - "Test list"
Cohesion: 0.04
Nodes (46): 1. The rows are consistent (SQL, no machine needed), 2. The rows agree with the mode log (SQL), 3. The raw signals around one cycle (telemetry and logs), 4. At the machine (floor test), Before the tests (once), Check query, `CNC_CYCLE_MAX_HOURS` (24 hours), `CNC_STALE_SEC` (30 s) and `CNC_QUIET_SEC` (300 s) (+38 more)

### Community 62 - "0001. Go standard library net/http, no web framework"
Cohesion: 0.33
Nodes (5): 0001. Go standard library net/http, no web framework, Alternatives considered, Consequences, Context, Decision

### Community 63 - "0003. PostgreSQL with pgx and hand-written SQL"
Cohesion: 0.33
Nodes (5): 0003. PostgreSQL with pgx and hand-written SQL, Alternatives considered, Consequences, Context, Decision

### Community 64 - "dto_test.go"
Cohesion: 0.20
Nodes (9): TestCreateBrandRequestValidate(), TestCreateCategoryRequestValidate(), TestCreateControllerRequestValidate(), TestCreateProductRequestValidate(), TestUpdateBrandRequestValidate(), TestUpdateCategoryRequestValidate(), TestUpdateControllerRequestValidate(), TestUpdateProductRequestValidate() (+1 more)

### Community 65 - "VirtualScreenResponse"
Cohesion: 0.20
Nodes (7): VirtualScreenCycle, VirtualScreenLive, VirtualScreenResponse, NewVirtualScreenLive(), JobProgress, LiveVariable, MachineAlarm

### Community 67 - "VariableConfig"
Cohesion: 0.24
Nodes (7): VariableConfig, formatPreviousValue(), formatValue(), newValueChangeLog(), boundVariable, valueChange, valueChangeLog

### Community 69 - "LiveHandler"
Cohesion: 0.10
Nodes (14): Browser → server, liveSocket, liveSubscription, LiveSessionExpired, LiveSubscribeResponse, channelRoom(), closeLiveSocket(), decodeLiveMessage() (+6 more)

### Community 70 - "WithTransaction"
Cohesion: 0.27
Nodes (10): WithTransaction(), CustomerDetails, UpdateUserRequest, UserResponse, ParseDayStart(), companyProfileFrom(), primaryAddress(), setCustomerAddress() (+2 more)

### Community 71 - "onDataCallback(ptr uintptr) uintptr"
Cohesion: 0.18
Nodes (15): Package-level callback singletons (syscall.NewCallback), CNCConnect(ip, port, dataCb, errCb, machineType) int32, CNCConnect param dataCb (uintptr stdcall OnData(const CncMasterSnapshot*)), CNCConnect param errCb (uintptr stdcall OnErr(int32 code, const char* msg)), CNCConnect param ip (const char*, ANSI NUL-terminated), CNCConnect param port (int32, e.g. 4840), Mutex-protected copy: dataMutex sync.RWMutex -> lastRecord/hasRecord, onDataCallback(ptr uintptr) uintptr (+7 more)

### Community 72 - "DashboardService"
Cohesion: 0.20
Nodes (14): CycleCountRow, MachineTimeRow, cyclesIn(), dayWindows(), DashboardService, hours(), machineDayOEE(), modeGroup() (+6 more)

### Community 73 - "UserMachineFilter"
Cohesion: 0.32
Nodes (4): UserMachineFilter, UserMachineFilterRepository, userAccess, userServices

### Community 74 - "Decision"
Cohesion: 0.40
Nodes (6): Decision, FinishCyclesSummary, FinishedCycle, FinishCyclesResponse, NewFinishCyclesSummary(), NewFinishedCycle()

### Community 75 - "0013. Machine diagnostics: reasons from what the collector recorded, network tested from the server"
Cohesion: 0.33
Nodes (5): 0013. Machine diagnostics: reasons from what the collector recorded, network tested from the server, Alternatives considered, Consequences, Context, Decision

### Community 76 - "0008. "Log in as customer" uses a separate session with a path-scoped cookie"
Cohesion: 0.40
Nodes (4): 0008. "Log in as customer" uses a separate session with a path-scoped cookie, Alternatives considered, Consequences, Context

### Community 77 - "0020. Machine PLC node IDs: the machine's own node ID in machine_variables, default otherwise"
Cohesion: 0.33
Nodes (5): 0020. Machine PLC node IDs: the machine's own node ID in machine_variables, default otherwise, Alternatives considered, Consequences, Context, Decision

### Community 78 - "MonitoredMachine"
Cohesion: 0.08
Nodes (35): 1. The big picture, 2. All steps at a glance, 4. Find the machines: steps 5 – 9, 6. When data is edited: step 18, 7. Shutdown: steps 19 – 20, 8. Reference, 9. Data Reading Flow (Function Names Only), docs/cnc-communication.md (+27 more)

### Community 79 - "run"
Cohesion: 0.36
Nodes (10): checkNotTaken(), createSuperadmin(), ensureProfile(), execute(), main(), parseFlags(), run(), superadminRoleID() (+2 more)

### Community 80 - "Event-driven data flow (no polling in client)"
Cohesion: 0.21
Nodes (13): Check mode: go run sinumerik_client.go check (OPC UA node existence report), Flow 1: load DLL, register two stdcall callbacks (syscall.NewCallback), Flow 2: CNCConnect starts DLL background worker thread, Flow 3: worker polls controller, assembles snapshot, invokes data callback, Flow 4: callback copies payload, signals renderer via buffered channel, Flow 5: main goroutine renders live ANSI dashboard, Go Client (sinumerik_client.go, callback -> copy -> render), OPC UA (Siemens transport, port 4840) (+5 more)

### Community 82 - "Role"
Cohesion: 0.33
Nodes (3): Role, mapRoleError(), toRoleResponseError()

### Community 83 - "ProductionScheduleRepository"
Cohesion: 0.08
Nodes (22): validDate(), CreateProductionScheduleRequest, OperatorJobFilter, ProductionScheduleFields, ProductionScheduleItemFields, UpdateProductionScheduleRequest, ScheduleDay(), MachineQuantity (+14 more)

### Community 84 - "0004. Opaque session token in an HttpOnly cookie"
Cohesion: 0.40
Nodes (4): 0004. Opaque session token in an HttpOnly cookie, Alternatives considered, Consequences, Context

### Community 85 - "MachineTelemetryRepository"
Cohesion: 0.09
Nodes (9): Decision, MachineEvent, MachineLiveStatus, ModeChange, MachineTelemetryRepository, mapMachineTelemetryError(), mapTelemetryError(), execute() (+1 more)

### Community 86 - "auth_tokens SHA-256 Hashed Sessions (mig"
Cohesion: 0.67
Nodes (3): auth_tokens SHA-256 Hashed Sessions (migration 014), Per-Request Session and User Re-check, ts_session Cookie (HttpOnly, SameSite=Strict)

### Community 88 - "AuthUser"
Cohesion: 0.05
Nodes (33): CreateDowntimeReasonRequest, DowntimeReasonFields, UpdateDowntimeReasonRequest, CreateHolidayRequest, HolidayFields, UpdateHolidayRequest, CreateIdleReasonRequest, IdleReasonFields (+25 more)

### Community 89 - "Live data — Socket.IO event contract"
Cohesion: 0.14
Nodes (11): 0009. Socket.IO live data: the collector pushes each screen's variables, checked per user, Alternatives considered, Consequences, Context, Adding a variable to a screen, Connecting, Live data — Socket.IO event contract, Ordering (+3 more)

### Community 104 - "connection.go"
Cohesion: 0.07
Nodes (22): polledChanges(), TestPolledChanges(), changeOf(), forwardNotifications(), Change, Connection, Item, ItemResult (+14 more)

### Community 105 - "0005. CNC collector runs inside the API process"
Cohesion: 0.33
Nodes (5): 0005. CNC collector runs inside the API process, Alternatives considered, Consequences, Context, Decision

### Community 106 - "0006. Page-wise dashboard endpoints returning only chart fields"
Cohesion: 0.40
Nodes (4): 0006. Page-wise dashboard endpoints returning only chart fields, Alternatives considered, Consequences, Context

### Community 107 - "0011. Ponytail plugin for Claude Code, under this repo's rules"
Cohesion: 0.33
Nodes (5): 0011. Ponytail plugin for Claude Code, under this repo's rules, Alternatives considered, Consequences, Context, Decision

### Community 108 - "ClienteleService"
Cohesion: 0.13
Nodes (11): ClienteleFields, CreateClienteleRequest, UpdateClienteleRequest, validClienteleEmail(), Clientele, ClienteleRepository, mapClienteleError(), clienteleFromFields() (+3 more)

### Community 109 - "0010. Machine access by an organization filter, worked out on every request"
Cohesion: 0.40
Nodes (4): 0010. Machine access by an organization filter, worked out on every request, Alternatives considered, Consequences, Context

### Community 110 - "0022. Test PLC node IDs in a short OPC UA session of their own"
Cohesion: 0.40
Nodes (4): 0022. Test PLC node IDs in a short OPC UA session of their own, Alternatives considered, Consequences, Context

### Community 113 - "0025. A machine is saved at one level of its customer's organization; the path above is read, not stored"
Cohesion: 0.40
Nodes (4): 0025. A machine is saved at one level of its customer's organization; the path above is read, not stored, Alternatives considered, Consequences, Context

### Community 114 - "0007. OpenAPI 3 from swag annotations, Swagger UI behind API_DOCS, generated Postman"
Cohesion: 0.33
Nodes (5): 0007. OpenAPI 3 from swag annotations, Swagger UI behind API_DOCS, generated Postman, Alternatives considered, Consequences, Context, Decision

### Community 115 - "0002. Layered architecture: handler, DTO, service, repository"
Cohesion: 0.33
Nodes (5): 0002. Layered architecture: handler, DTO, service, repository, Alternatives considered, Consequences, Context, Decision

### Community 116 - "HomeTime"
Cohesion: 0.40
Nodes (5): HomeTime, HomeTimeGroup, HomeTimeMode, homeTimeGroups(), homeTimeModes()

### Community 117 - "context.Context"
Cohesion: 0.12
Nodes (6): ListenForDatabaseNotifications(), TestPlcNodeIDsRequest, Connect(), levelHasMachines(), SessionRepository, levelInMachineFilter()

### Community 118 - "UserStateService"
Cohesion: 0.18
Nodes (7): CreateUserStateRequest, validateUserState(), UserState, UserStateRepository, mapUserStateError(), UserStateService, toUserStateResponseError()

### Community 119 - "0017. Machine accuracy as one shared list, managed by the superadmin"
Cohesion: 0.40
Nodes (4): 0017. Machine accuracy as one shared list, managed by the superadmin, Alternatives considered, Consequences, Context

### Community 120 - "LocationRepository"
Cohesion: 0.50
Nodes (3): mapLocationError(), NewLocationRepository(), LocationRepository

### Community 121 - "0015. Operator tasks: a shared task list, assigned per customer and machine model"
Cohesion: 0.50
Nodes (4): 0015. Operator tasks: a shared task list, assigned per customer and machine model, Alternatives considered, Consequences, Context

### Community 126 - "Connection"
Cohesion: 0.06
Nodes (27): connectFanucCounter(), main(), run(), dllDriver, slot, 3. Startup: steps 1 – 4, 5. Connect and collect, per machine: steps 10 – 17, Step 11 — Connect, a few machines at a time (+19 more)

### Community 128 - "MachineService"
Cohesion: 0.06
Nodes (37): Decision, Decision, PlcReaderResponse, CreateMachineRequest, MachineFields, UpdateMachineRequest, validCoordinate(), validHealthIDs() (+29 more)

### Community 130 - "LicenceKeyService"
Cohesion: 0.13
Nodes (12): NewLicenceKey, CreateLicenceKeysRequest, GenerateLicenceKeysRequest, GeneratedLicenceKeys, LicenceKey, LicenceKeyRepository, mapLicenceKeyError(), LicenceKeyService (+4 more)

### Community 133 - ".GetDiagnostics"
Cohesion: 0.07
Nodes (49): Decision, DiagnosticPLCVariable, DiagnosticCheck, DiagnosticMachine, DiagnosticPLC, DiagnosticVariableRow, MachineDiagnosticsResponse, MachineNetworkResponse (+41 more)

### Community 134 - "pathID"
Cohesion: 0.08
Nodes (11): targetUserID(), RoleHandler, UserCellHandler, UserCityHandler, UserCountryHandler, UserLineHandler, UserLocationHandler, UserShopHandler (+3 more)

### Community 138 - "TestPLCCollector"
Cohesion: 0.10
Nodes (25): 0021. Read PLC variables over OPC UA with gopcua, in-process beside the CNC collector, Alternatives considered, Consequences, Context, Decision, NewDashboardHandler(), NewMachineDiagnosticsHandler(), NewPlcCollector() (+17 more)

## Ambiguous Edges - Review These
- `CncCommunication.dll (32-bit vendor DLL)` → `FanucCncCommunication.dll`  [AMBIGUOUS]
  docs/setup-guide.txt · relation: semantically_similar_to
- `Machines API (/machines)` → `Snapshot Fields (positions, spindle, feed, program, tool, alarms)`  [AMBIGUOUS]
  docs/CncCommunication.dll — API & Integration Guide (Go Client).pdf · relation: conceptually_related_to

## Knowledge Gaps
- **33 isolated node(s):** `7th-sense-server`, `middleware.RequireRole / routes.adminOnly`, `response.SendResponse`, `Test Auth (asTestUser, tests/auth_test.go)`, `user_filters Filter Checklist` (+28 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 272 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **28 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `CNC Data Collection Flow` connect `MonitoredMachine` to `Worker (internal/monitor/worker.go)`, `Live data — Socket.IO event contract`, `Connection`?**
  _High betweenness centrality (0.092) - this node is a cross-community bridge._
- **What connects `7th-sense-server`, `middleware.RequireRole / routes.adminOnly`, `response.SendResponse` to the rest of the system?**
  _33 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `HomeDashboardResponse` be split into smaller, more focused modules?**
  _Cohesion score 0.1038961038961039 - nodes in this community are weakly interconnected._
- **Why does `7th Sense Server Install and Setup Guide` connect `Worker (internal/monitor/worker.go)` to `MonitoredMachine`?**
  _High betweenness centrality (0.082) - this node is a cross-community bridge._
- **Should `telemetry_load_test.go` be split into smaller, more focused modules?**
  _Cohesion score 0.05701754385964912 - nodes in this community are weakly interconnected._

### Low-confidence Hints
_AMBIGUOUS edges — the extractor was unsure. Verify before acting on these._

- **What is the exact relationship between `CncCommunication.dll (32-bit vendor DLL)` and `FanucCncCommunication.dll`?**
  _Edge tagged AMBIGUOUS (relation: semantically_similar_to) - confidence is low._
- **What is the exact relationship between `Machines API (/machines)` and `Snapshot Fields (positions, spindle, feed, program, tool, alarms)`?**
  _Edge tagged AMBIGUOUS (relation: conceptually_related_to) - confidence is low._