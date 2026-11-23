# Graph Report - 7th-sense-server  (2026-11-23)

## Corpus Check
- 388 files · ~334,124 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 4 file(s) not represented in the graph (top: (none) 3, .toml 1)

## Summary
- 3476 nodes · 32562 edges · 137 communities (117 shown, 20 thin omitted)
- Extraction: 95% EXTRACTED · 5% INFERRED · 0% AMBIGUOUS · INFERRED: 1468 edges (avg confidence: 0.87)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `099bf270`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- telemetry_load_test.go
- UserCityRepository
- .GetMachineSummaries
- ControllerService
- MachinePlcVariableHandler
- newFilterLevelRouter
- apidocs/main.go
- MachineModeReasonService
- api_helpers_test.go
- CncMasterSnapshot (34,440 bytes, #pragma
- AuthHandler
- MachineModelService
- UserCountryRepository
- MachineTelemetryService
- worker
- Hub
- net/http.Request
- mode-chart.js
- validText
- net/http.Request
- net/http.ServeMux
- time.Duration
- SendResponse
- §10 Diagnostics & Troubleshooting
- OperatorTaskService
- machine_mode.go
- UserLineService
- Product
- PlcCollector
- github.com/jackc/pgx/v5/pgxpool.Pool
- TelemetryQueue
- AuthUser
- Worker (internal/monitor/worker.go)
- README.md
- run
- testing.T
- Category
- MachineModePeriod
- 0016. A Task List screen for the operator tasks' master list
- .GetAxisValues
- CncCommunication.dll API & Integration G
- AuthService
- uniqueSuffix
- CncCommunication.dll API & Integration G
- UserShopService
- config.go
- UserCellService
- NewError
- Machines API (/machines)
- TestLiveData
- Address
- MachineAccuracyService
- BranchService
- README.md
- DashboardRequest
- CompanyProfile
- MachineCycle
- UserService
- DashboardService
- executeAffected
- UserProfileRepository
- Test list
- 0001. Go standard library net/http, no web framework
- 0003. PostgreSQL with pgx and hand-written SQL
- BrandRepository
- VirtualScreenResponse
- ControllerRepository
- VariableConfig
- toUserCityResponseError
- liveSocket
- 0014. Session modes and the login conflict, on the existing auth_tokens sessions
- onDataCallback(ptr uintptr) uintptr
- .GetHomeDashboard
- newLiveTestServerWith
- BrandService
- run
- 0008. "Log in as customer" uses a separate session with a path-scoped cookie
- 0022. Test PLC node IDs in a short OPC UA session of their own
- Supervisor
- run
- Event-driven data flow (no polling in client)
- AxisValuesRequest
- Role
- ProductionScheduleRepository
- 0023. PLC reader: short session timeout, and read the nodes the controller has no room to monitor
- Execute
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
- validateUserLine
- 0020. Machine PLC node IDs: the machine's own node ID in machine_variables, default otherwise
- ClienteleService
- 0011. Ponytail plugin for Claude Code, under this repo's rules
- validateUserShop
- expectStatus
- plc_collector_test.go
- PartService
- 0007. OpenAPI 3 from swag annotations, Swagger UI behind API_DOCS, generated Postman
- 0002. Layered architecture: handler, DTO, service, repository
- 0018. Machine Health X / Y / Z charts: downsampled in SQL, timed, on the machine-health routes
- context.Context
- UserStateService
- 0017. Machine accuracy as one shared list, managed by the superadmin
- LocationRepository
- 0006. Page-wise dashboard endpoints returning only chart fields
- UserLocationService
- newUserServices
- NewUserService
- Connection
- SessionRevokedError
- MachineService
- 0010. Machine access by an organization filter, worked out on every request
- LicenceKeyService
- 0013. Machine diagnostics: reasons from what the collector recorded, network tested from the server
- .GetDiagnostics
- pathID
- net/http.Handler
- 0025. A machine is saved at one level of its customer's organization; the path above is read, not stored
- PartDrawingStorage
- SessionConflictError

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

## Communities (137 total, 20 thin omitted)

### Community 0 - "telemetry_load_test.go"
Cohesion: 0.17
Nodes (17): Load test after migration 039 (2026-11-23), Load test (Step 8), NewTelemetryQueue(), loadPlan, loadResult, envInts(), envRatio(), loadSummary() (+9 more)

### Community 1 - "UserCityRepository"
Cohesion: 0.27
Nodes (3): UserCity, UserCityRepository, mapUserCityError()

### Community 2 - ".GetMachineSummaries"
Cohesion: 0.23
Nodes (7): MachineSummary, OEEResponse, OEERow, NewMachineSummary(), round1(), machineOEE, oeeSum

### Community 3 - "ControllerService"
Cohesion: 0.21
Nodes (7): CreateControllerRequest, UpdateControllerRequest, validateControllerFields(), NewControllerHandler(), ControllerService, toControllerResponseError(), NewMachineService()

### Community 5 - "newFilterLevelRouter"
Cohesion: 0.18
Nodes (22): NewBranchHandler(), NewUserCityHandler(), NewUserCountryHandler(), NewUserHandler(), NewUserLineHandler(), NewUserLocationHandler(), NewUserShopHandler(), NewUserStateHandler() (+14 more)

### Community 6 - "apidocs/main.go"
Cohesion: 0.08
Nodes (37): request, main(), EnsureCertificates(), existingCertificateExpiry(), generateCertificates(), Paths, writePEM(), Init() (+29 more)

### Community 7 - "MachineModeReasonService"
Cohesion: 0.10
Nodes (20): 0012. Downtime and idle reasons on a period: one table, tied to machine_mode_log, Alternatives considered, Consequences, Context, Decision, CreateMachineModeReasonRequest, MachineModeReasonFields, MachineModeReasonsResponse (+12 more)

### Community 8 - "api_helpers_test.go"
Cohesion: 0.40
Nodes (8): cncMasterSnapshot, RevokeSessionsResponse, authUserKey, loginAsKey, sessionIDKey, liveValues, machineModelRecord, subscribeAnswer

### Community 9 - "CncMasterSnapshot (34,440 bytes, #pragma"
Cohesion: 0.09
Nodes (47): CncMasterSnapshot (34,440 bytes, #pragma pack(8), x86), AbsPos ([5]float64, Inches / Degrees, Work Coordinate System (WCS) positions), ActiveGCode ([16]byte, ASCII, Currently active G-code group), ActiveProgPath ([256]byte, ASCII Path, Storage path of active NC artifact), AlarmMsg ([128]byte, UTF-8 / ASCII, Alarm description message), AlarmNo (int32, Code (0 = OK), Primary alarm or fault number), AxisLoadXPercentage (int16, Percent (%), Servo X axis drive load (AxisLoadX/Y/ZPercentage)), AxisLoadYPercentage (int16, Percent (%), Servo Y axis drive load (AxisLoadX/Y/ZPercentage)) (+39 more)

### Community 10 - "AuthHandler"
Cohesion: 0.09
Nodes (18): Decision, 0009. Socket.IO live data: the collector pushes each screen's variables, checked per user, Alternatives considered, Consequences, Context, Decision, AuthHandler, CurrentUser() (+10 more)

### Community 11 - "MachineModelService"
Cohesion: 0.16
Nodes (8): Decision, MachineModelRequest, NewMachineModelHandler(), MachineModel, MachineModelRepository, mapMachineModelError(), MachineModelService, toMachineModelResponseError()

### Community 12 - "UserCountryRepository"
Cohesion: 0.21
Nodes (6): CreateUserCountryRequest, UserCountry, UserCountryRepository, mapUserCountryError(), NewUserCountryRepository(), toUserCountryResponseError()

### Community 13 - "MachineTelemetryService"
Cohesion: 0.08
Nodes (14): Step 5 — Prepare the database, MachineModesResponse, NewHealthHandler(), MachineEvent, MachineLiveStatus, MachineModeStatus, ModeChange, mapTelemetryError() (+6 more)

### Community 14 - "worker"
Cohesion: 0.06
Nodes (38): Step 15 — Store the reading, `cycleCheckpointInterval` (5 minutes), 1. Why, and what today looks like, Step 2: Message model, Step 6: Workers publish instead of writing, FinishCyclesSummary, FinishedCycle, TimelinePeriod (+30 more)

### Community 15 - "Hub"
Cohesion: 0.07
Nodes (40): 1. Machine → DLL → server, 2. Worker → hub (before the database), 3. Hub → rooms, 4. Socket.IO → browser, Step by step, LiveSubscribeResponse, NewMachineDiagnosticsHandler(), NewMachinePlcVariableHandler() (+32 more)

### Community 16 - "net/http.Request"
Cohesion: 0.07
Nodes (15): BreakHandler, ClienteleHandler, DowntimeReasonHandler, HolidayHandler, IdleReasonHandler, MachineHandler, MachineModeReasonHandler, PartHandler (+7 more)

### Community 17 - "mode-chart.js"
Cohesion: 0.10
Nodes (64): card(), cell(), chartData(), chartRows(), countByGroup(), currentSpan(), currentWindow(), cycleData() (+56 more)

### Community 18 - "validText"
Cohesion: 0.19
Nodes (8): PartProcessFields, PartSubProgramFields, CreatePartRequest, PartFields, UpdatePartRequest, UploadPartDrawingRequest, validText(), PartDrawing

### Community 19 - "net/http.Request"
Cohesion: 0.10
Nodes (10): axisValuesQuery(), dashboardQuery(), DashboardHandler, readDashboardRequest(), HealthHandler, isSameOrigin(), OperatorTaskHandler, ProductHandler (+2 more)

### Community 20 - "net/http.ServeMux"
Cohesion: 0.12
Nodes (41): Decision, RequireRole(), registerAuthRoutes(), registerBranchRoutes(), registerBrandRoutes(), registerBreakRoutes(), registerCategoryRoutes(), registerClienteleRoutes() (+33 more)

### Community 21 - "time.Duration"
Cohesion: 0.07
Nodes (41): Decision, PLCConfig, RabbitMQConfig, secondsOf(), newTestWriter(), runWriter(), testDelivery(), TestWriterAcksOnlyAfterSuccess() (+33 more)

### Community 22 - "SendResponse"
Cohesion: 0.07
Nodes (12): validator, BranchHandler, BrandHandler, CategoryHandler, ControllerHandler, LicenceKeyHandler, MachineAccuracyHandler, MachineModelHandler (+4 more)

### Community 23 - "§10 Diagnostics & Troubleshooting"
Cohesion: 0.12
Nodes (24): Check mode: go run sinumerik_client.go check (OPC UA node existence report), Go Client (sinumerik_client.go, callback -> copy -> render), Prerequisite: 32-bit x86 DLL, compile Go with GOARCH=386, Prerequisite: SIEMENS_USER / SIEMENS_PASS env vars before CNCConnect, Run: go run sinumerik_client.go (live dashboard), §3.1 Environment & Prerequisites, §8 Running the Go Client, §10 Diagnostics & Troubleshooting (+16 more)

### Community 24 - "OperatorTaskService"
Cohesion: 0.08
Nodes (20): CreateMaintenanceTaskRequest, CreateOperatorTaskRequest, OperatorTaskFields, UpdateMaintenanceTaskRequest, UpdateOperatorTaskRequest, NewMaintenanceTaskHandler(), NewOperatorTaskHandler(), MaintenanceTask (+12 more)

### Community 25 - "machine_mode.go"
Cohesion: 0.15
Nodes (19): decideMachineMode(), holdReason(), machineStateFromReading(), numberingEngine(), programStateFromProcStatus(), ncReading(), TestDecideMachineMode(), TestDecideMachineModeConditions() (+11 more)

### Community 26 - "UserLineService"
Cohesion: 0.21
Nodes (5): UserLine, UserLineRepository, mapUserLineError(), UserLineService, toUserLineResponseError()

### Community 27 - "Product"
Cohesion: 0.14
Nodes (12): CreateProductRequest, UpdateProductRequest, validateProductFields(), NewProductHandler(), Product, ProductRepository, mapProductError(), NewProductRepository() (+4 more)

### Community 28 - "PlcCollector"
Cohesion: 0.15
Nodes (12): 0021. Read PLC variables over OPC UA with gopcua, in-process beside the CNC collector, Alternatives considered, Consequences, Context, Decision, samePLCConnection(), newPlcWorker(), MachineStatus (+4 more)

### Community 29 - "github.com/jackc/pgx/v5/pgxpool.Pool"
Cohesion: 0.05
Nodes (39): WithTransaction(), BreakFields, CreateBreakRequest, UpdateBreakRequest, uniquePositiveIDs(), clockMinutes(), CreateShiftRequest, ShiftFields (+31 more)

### Community 30 - "TelemetryQueue"
Cohesion: 0.10
Nodes (12): TelemetryReading, TelemetryQueue, decodeReading(), TestWriterKeepsNewestLiveValue(), roundToTenth(), ReadingValue, batchActions, batchEntry (+4 more)

### Community 31 - "AuthUser"
Cohesion: 0.13
Nodes (13): Decision, LoginAsResponse, describeDevice(), SessionResponse, NewSessionResponse(), NewSessionResponses(), ParseUserAgent(), AuthUser (+5 more)

### Community 32 - "Worker (internal/monitor/worker.go)"
Cohesion: 0.06
Nodes (29): CNC Data Collection (in-process), writeLog CNC Log Types and Switches, .env Configuration (internal/config), Health Mode History API (public), 7th Sense Server Claude Context, response.SendResponse Format, User CRUD Golden Reference, CncCommunication.dll (32-bit vendor DLL) (+21 more)

### Community 33 - "README.md"
Cohesion: 0.24
Nodes (5): 0015. Operator tasks: a shared task list, assigned per customer and machine model, Alternatives considered, Consequences, Context, Architecture Decision Records

### Community 34 - "run"
Cohesion: 0.16
Nodes (24): run(), Querier, NewClienteleRepository(), NewDowntimeReasonRepository(), NewHolidayRepository(), NewIdleReasonRepository(), NewLicenceKeyRepository(), NewMachineAccuracyRepository() (+16 more)

### Community 35 - "testing.T"
Cohesion: 0.07
Nodes (89): TestSnapshotModeFieldOffsets(), TestSnapshotSize(), NewDashboardHandler(), WithAuthUser(), decodeData(), sendRequest(), apiResponse, mustDecode() (+81 more)

### Community 36 - "Category"
Cohesion: 0.15
Nodes (10): CreateCategoryRequest, UpdateCategoryRequest, NewCategoryHandler(), Category, CategoryRepository, mapCategoryError(), NewCategoryRepository(), CategoryService (+2 more)

### Community 37 - "MachineModePeriod"
Cohesion: 0.16
Nodes (15): Decision, MachineTimeRow, NewFinishedCycle(), NewTimelinePeriod(), MachineModeHistoryResponse, MachineModeTotal, MachineModePeriod, hours() (+7 more)

### Community 38 - "0016. A Task List screen for the operator tasks' master list"
Cohesion: 0.40
Nodes (5): 0016. A Task List screen for the operator tasks' master list, Alternatives considered, Consequences, Context, Decision

### Community 39 - ".GetAxisValues"
Cohesion: 0.16
Nodes (11): AxisPoint, AxisBucketSeconds(), FindAxisSet(), AxisSeries, AxisSet, AxisStats, AxisValuesResponse, AxisBucket (+3 more)

### Community 40 - "CncCommunication.dll API & Integration G"
Cohesion: 0.22
Nodes (18): CncCommunication.dll (C++ core runtime, x86, OPC UA inside), Fanuc profile (engine 1, port 8193, Fwlib32.dll + fwlibe1.dll), Fwlib32.dll (Fanuc FOCAS2), fwlibe1.dll (Fanuc FOCAS2), libcrypto-3.dll (OpenSSL), libssl-3.dll (OpenSSL), Engine 1: Fanuc (FOCAS2 Ethernet, fwlib32), OPC UA (Siemens transport, port 4840) (+10 more)

### Community 41 - "AuthService"
Cohesion: 0.12
Nodes (14): SessionClient, NewLiveHandler(), SessionUser, AuthRepository, NewSessionRepository(), AuthService, hashToken(), newLoginAsKey() (+6 more)

### Community 42 - "uniqueSuffix"
Cohesion: 0.22
Nodes (30): expectStatus(), uniqueSuffix(), authRequest(), createAuthTestUser(), login(), newAuthRouter(), sessionCookieOf(), TestAuthentication() (+22 more)

### Community 43 - "CncCommunication.dll API & Integration G"
Cohesion: 0.13
Nodes (23): client_cert.der/.pem (X.509 client cert for Siemens OPC UA), client_key.der/.pem (client private key), CNCConnect param machineType (int32 engine id, see §4.2), Engine 0: Simulator (local loopback mock engine), Engine 2: Siemens (Sinumerik OPC UA, active profile), Engine 3: Syntec (Open API / TCP), Engine 4: Mitsubishi (driver not specified), Machine identifier mapping (engine id enum 0-4) (+15 more)

### Community 44 - "UserShopService"
Cohesion: 0.21
Nodes (5): UserShop, UserShopRepository, mapUserShopError(), UserShopService, toUserShopResponseError()

### Community 45 - "config.go"
Cohesion: 0.06
Nodes (36): main(), run(), MachineSelection, 0024. Telemetry rows through a RabbitMQ queue, written in batches, Alternatives considered, Consequences, Context, 2. Scope (+28 more)

### Community 46 - "UserCellService"
Cohesion: 0.07
Nodes (20): CreateUserCellRequest, UpdateUserCellRequest, validateUserCell(), CreateUserSupervisorRequest, UpdateUserSupervisorRequest, validateUserSupervisor(), NewUserCellHandler(), NewUserSupervisorHandler() (+12 more)

### Community 47 - "NewError"
Cohesion: 0.09
Nodes (24): CreateAddressRequest, UpdateAddressRequest, UpdateLanguageRequest, UpdateThemeRequest, LoginRequest, MachineFilter, isTimeZone(), isWebsite() (+16 more)

### Community 48 - "Machines API (/machines)"
Cohesion: 0.16
Nodes (18): Clienteles API, Company Scope (companyOwner / ownCompanyCustomer), Current User API (/me theme, language), Customer Display Names (customerName / companyName), Holidays API, Licence Keys API, Machine Catalog Brand-Controller-Model, Machines API (/machines) (+10 more)

### Community 49 - "TestLiveData"
Cohesion: 0.26
Nodes (5): dialLive(), newLiveTestServer(), liveTestClient, liveTestEvent, liveTestServer

### Community 50 - "Address"
Cohesion: 0.22
Nodes (8): Address, AddressRepository, mapAddressError(), NewAddressRepository(), NewAddressService(), toAddressResponseError(), AddressService, TestAddressCRUD()

### Community 51 - "MachineAccuracyService"
Cohesion: 0.19
Nodes (7): MachineAccuracyRequest, NewMachineAccuracyHandler(), MachineAccuracy, MachineAccuracyRepository, mapMachineAccuracyError(), MachineAccuracyService, toMachineAccuracyResponseError()

### Community 52 - "BranchService"
Cohesion: 0.15
Nodes (9): CreateBranchRequest, UpdateBranchRequest, validateBranchFields(), Branch, BranchRepository, mapBranchError(), cityOwner(), BranchService (+1 more)

### Community 53 - "README.md"
Cohesion: 0.40
Nodes (4): 0004. Opaque session token in an HttpOnly cookie, Alternatives considered, Consequences, Context

### Community 54 - "DashboardRequest"
Cohesion: 0.22
Nodes (5): DashboardRequest, HomeDashboardRequest, ModeTotalsResponse, MachineCycleListRequest, ModeDuration

### Community 55 - "CompanyProfile"
Cohesion: 0.28
Nodes (5): CompanyProfile, CompanyProfileRepository, mapCompanyProfileError(), toCompanyProfileResponseError(), CompanyProfileService

### Community 56 - "MachineCycle"
Cohesion: 0.14
Nodes (10): MachineCyclesResponse, NewMachineCycleHandler(), CycleJob, MachineCycle, MachineCycleRepository, mapMachineCycleError(), MachineCycleService, NewMachineCycleService() (+2 more)

### Community 57 - "UserService"
Cohesion: 0.14
Nodes (17): CreateUserRequest, CustomerDetails, UpdateUserRequest, UserResponse, ParseDayStart(), TestParseDayStart(), AssignableRoles(), CanAssignRole() (+9 more)

### Community 58 - "DashboardService"
Cohesion: 0.33
Nodes (8): CycleCountRow, clipWindows(), cyclesIn(), dayWindows(), DashboardService, sortedKeys(), summarizeCycles(), timeWindow

### Community 59 - "executeAffected"
Cohesion: 0.17
Nodes (7): Part, PartProcess, PartSubProgram, executeAffected(), PartRepository, mapPartError(), saveProcesses()

### Community 60 - "UserProfileRepository"
Cohesion: 0.18
Nodes (6): UserProfile, UserProfileRepository, mapProfileError(), NewUserProfileService(), toUserProfileResponseError(), UserProfileService

### Community 61 - "Test list"
Cohesion: 0.05
Nodes (39): 1. The rows are consistent (SQL, no machine needed), 2. The rows agree with the mode log (SQL), 3. The raw signals around one cycle (telemetry and logs), 4. At the machine (floor test), Before the tests (once), Check query, Completed, Floor tests on the HUMA panel (Fanuc) (+31 more)

### Community 62 - "0001. Go standard library net/http, no web framework"
Cohesion: 0.33
Nodes (5): 0001. Go standard library net/http, no web framework, Alternatives considered, Consequences, Context, Decision

### Community 63 - "0003. PostgreSQL with pgx and hand-written SQL"
Cohesion: 0.33
Nodes (5): 0003. PostgreSQL with pgx and hand-written SQL, Alternatives considered, Consequences, Context, Decision

### Community 64 - "BrandRepository"
Cohesion: 0.27
Nodes (3): Brand, BrandRepository, mapBrandError()

### Community 65 - "VirtualScreenResponse"
Cohesion: 0.20
Nodes (7): VirtualScreenCycle, VirtualScreenLive, VirtualScreenResponse, NewVirtualScreenLive(), JobProgress, LiveVariable, MachineAlarm

### Community 66 - "ControllerRepository"
Cohesion: 0.31
Nodes (3): Controller, ControllerRepository, mapControllerError()

### Community 67 - "VariableConfig"
Cohesion: 0.22
Nodes (7): VariableConfig, formatPreviousValue(), formatValue(), newValueChangeLog(), boundVariable, valueChange, valueChangeLog

### Community 68 - "toUserCityResponseError"
Cohesion: 0.27
Nodes (4): CreateUserCityRequest, UpdateUserCityRequest, validateUserCity(), toUserCityResponseError()

### Community 69 - "liveSocket"
Cohesion: 0.07
Nodes (23): Connecting the socket, Files, How the variable list works, Live data — how a machine value reaches the screen, Showing a new variable on a screen, The big picture, What a screen asks for (subscribing), What the server reads (the machine's variable list) (+15 more)

### Community 70 - "0014. Session modes and the login conflict, on the existing auth_tokens sessions"
Cohesion: 0.40
Nodes (4): 0014. Session modes and the login conflict, on the existing auth_tokens sessions, Alternatives considered, Consequences, Context

### Community 71 - "onDataCallback(ptr uintptr) uintptr"
Cohesion: 0.14
Nodes (21): Package-level callback singletons (syscall.NewCallback), CNCConnect(ip, port, dataCb, errCb, machineType) int32, CNCConnect param dataCb (uintptr stdcall OnData(const CncMasterSnapshot*)), CNCConnect param errCb (uintptr stdcall OnErr(int32 code, const char* msg)), CNCConnect param ip (const char*, ANSI NUL-terminated), CNCConnect param port (int32, e.g. 4840), CNCdisconnect(machineId int32), CncCommunication.dll API & Integration Guide (Go Client) - text version (+13 more)

### Community 72 - ".GetHomeDashboard"
Cohesion: 0.09
Nodes (24): HomeAlarm, HomeCycles, HomeFleet, HomeMachine, HomePlan, HomeTime, HomeDashboardResponse, HomeHour (+16 more)

### Community 73 - "newLiveTestServerWith"
Cohesion: 0.44
Nodes (9): SessionConfig, NewAuthHandler(), NewSessionHandler(), Authenticate(), NewAuthRepository(), NewAuthService(), NewSessionService(), newLiveTestServerWith() (+1 more)

### Community 74 - "BrandService"
Cohesion: 0.23
Nodes (6): CreateBrandRequest, UpdateBrandRequest, validateCatalogFields(), NewBrandHandler(), BrandService, toBrandResponseError()

### Community 75 - "run"
Cohesion: 0.36
Nodes (4): connectFanucCounter(), main(), run(), fanucCounter

### Community 76 - "0008. "Log in as customer" uses a separate session with a path-scoped cookie"
Cohesion: 0.40
Nodes (4): 0008. "Log in as customer" uses a separate session with a path-scoped cookie, Alternatives considered, Consequences, Context

### Community 77 - "0022. Test PLC node IDs in a short OPC UA session of their own"
Cohesion: 0.40
Nodes (4): 0022. Test PLC node IDs in a short OPC UA session of their own, Alternatives considered, Consequences, Context

### Community 78 - "Supervisor"
Cohesion: 0.05
Nodes (54): dllDriver, slot, Supervisor, 1. The big picture, 2. All steps at a glance, 3. Startup: steps 1 – 4, 4. Find the machines: steps 5 – 9, 6. When data is edited: step 18 (+46 more)

### Community 79 - "run"
Cohesion: 0.36
Nodes (10): checkNotTaken(), createSuperadmin(), ensureProfile(), execute(), main(), parseFlags(), run(), superadminRoleID() (+2 more)

### Community 80 - "Event-driven data flow (no polling in client)"
Cohesion: 0.48
Nodes (5): Flow 1: load DLL, register two stdcall callbacks (syscall.NewCallback), Flow 2: CNCConnect starts DLL background worker thread, Flow 3: worker polls controller, assembles snapshot, invokes data callback, Flow 4: callback copies payload, signals renderer via buffered channel, Flow 5: main goroutine renders live ANSI dashboard

### Community 82 - "Role"
Cohesion: 0.33
Nodes (3): Role, mapRoleError(), toRoleResponseError()

### Community 83 - "ProductionScheduleRepository"
Cohesion: 0.07
Nodes (23): validDate(), CreateProductionScheduleRequest, OperatorJobFilter, ProductionScheduleFields, ProductionScheduleItemFields, UpdateProductionScheduleRequest, ScheduleDay(), NewProductionScheduleHandler() (+15 more)

### Community 84 - "0023. PLC reader: short session timeout, and read the nodes the controller has no room to monitor"
Cohesion: 0.40
Nodes (5): 0023. PLC reader: short session timeout, and read the nodes the controller has no room to monitor, Alternatives considered, Consequences, Context, Decision

### Community 85 - "Execute"
Cohesion: 0.11
Nodes (5): Execute(), AxisRow, MachineDiagnosticsRepository, MachineTelemetryRepository, mapMachineTelemetryError()

### Community 86 - "auth_tokens SHA-256 Hashed Sessions (mig"
Cohesion: 0.67
Nodes (3): auth_tokens SHA-256 Hashed Sessions (migration 014), Per-Request Session and User Re-check, ts_session Cookie (HttpOnly, SameSite=Strict)

### Community 88 - "AuthUser"
Cohesion: 0.05
Nodes (36): CreateDowntimeReasonRequest, DowntimeReasonFields, UpdateDowntimeReasonRequest, CreateHolidayRequest, HolidayFields, UpdateHolidayRequest, CreateIdleReasonRequest, IdleReasonFields (+28 more)

### Community 89 - "Live data — Socket.IO event contract"
Cohesion: 0.20
Nodes (7): Adding a variable to a screen, Connecting, Live data — Socket.IO event contract, Ordering, Rooms (server side), Server → browser, Settings

### Community 104 - "connection.go"
Cohesion: 0.05
Nodes (33): PlcNodeTestResult, checkPlcNodeID(), PlcNodeTestResponse, TestPlcNodeIDsRequest, UpdatePlcNodeIDRequest, polledChanges(), previousText(), TestPolledChanges() (+25 more)

### Community 105 - "0005. CNC collector runs inside the API process"
Cohesion: 0.33
Nodes (5): 0005. CNC collector runs inside the API process, Alternatives considered, Consequences, Context, Decision

### Community 106 - "validateUserLine"
Cohesion: 0.40
Nodes (3): CreateUserLineRequest, UpdateUserLineRequest, validateUserLine()

### Community 107 - "0020. Machine PLC node IDs: the machine's own node ID in machine_variables, default otherwise"
Cohesion: 0.33
Nodes (5): 0020. Machine PLC node IDs: the machine's own node ID in machine_variables, default otherwise, Alternatives considered, Consequences, Context, Decision

### Community 108 - "ClienteleService"
Cohesion: 0.13
Nodes (12): ClienteleFields, CreateClienteleRequest, UpdateClienteleRequest, validClienteleEmail(), NewClienteleHandler(), Clientele, ClienteleRepository, mapClienteleError() (+4 more)

### Community 109 - "0011. Ponytail plugin for Claude Code, under this repo's rules"
Cohesion: 0.33
Nodes (5): 0011. Ponytail plugin for Claude Code, under this repo's rules, Alternatives considered, Consequences, Context, Decision

### Community 110 - "validateUserShop"
Cohesion: 0.40
Nodes (3): CreateUserShopRequest, UpdateUserShopRequest, validateUserShop()

### Community 111 - "expectStatus"
Cohesion: 0.23
Nodes (15): containsID(), newCatalogRouter(), runCatalogCRUD(), createTestRecord(), TestBrandCRUD(), TestCategoryCRUD(), TestControllerBrand(), TestControllerCRUD() (+7 more)

### Community 112 - "plc_collector_test.go"
Cohesion: 0.33
Nodes (8): LanguageResponse, LicenceKeyDuration, LiveAccessChanged, LiveUnsubscribeRequest, MachineLocationFiltersResponse, ThemeResponse, licenceKeyDurations(), Envelope

### Community 113 - "PartService"
Cohesion: 0.21
Nodes (7): ValidPartDrawingFile(), NewPartHandler(), NewPartRepository(), PartService, NewPartService(), partFromFields(), toPartResponseError()

### Community 114 - "0007. OpenAPI 3 from swag annotations, Swagger UI behind API_DOCS, generated Postman"
Cohesion: 0.33
Nodes (5): 0007. OpenAPI 3 from swag annotations, Swagger UI behind API_DOCS, generated Postman, Alternatives considered, Consequences, Context, Decision

### Community 115 - "0002. Layered architecture: handler, DTO, service, repository"
Cohesion: 0.33
Nodes (5): 0002. Layered architecture: handler, DTO, service, repository, Alternatives considered, Consequences, Context, Decision

### Community 116 - "0018. Machine Health X / Y / Z charts: downsampled in SQL, timed, on the machine-health routes"
Cohesion: 0.17
Nodes (10): 0018. Machine Health X / Y / Z charts: downsampled in SQL, timed, on the machine-health routes, Alternatives considered, Consequences, Context, Decision, 0019. Axis values: every stored value when it fits, and zoom by loading the window again, Alternatives considered, Consequences (+2 more)

### Community 117 - "context.Context"
Cohesion: 0.12
Nodes (6): ListenForDatabaseNotifications(), SessionState, mapAuthError(), levelHasMachines(), SessionRepository, levelInMachineFilter()

### Community 118 - "UserStateService"
Cohesion: 0.15
Nodes (10): CreateUserStateRequest, validateUserState(), UserState, UserStateRepository, mapUserStateError(), NewUserStateRepository(), levelHasMachinesError(), levelInMachineFilterError() (+2 more)

### Community 119 - "0017. Machine accuracy as one shared list, managed by the superadmin"
Cohesion: 0.33
Nodes (5): 0017. Machine accuracy as one shared list, managed by the superadmin, Alternatives considered, Consequences, Context, Decision

### Community 120 - "LocationRepository"
Cohesion: 0.50
Nodes (3): mapLocationError(), NewLocationRepository(), LocationRepository

### Community 121 - "0006. Page-wise dashboard endpoints returning only chart fields"
Cohesion: 0.40
Nodes (4): 0006. Page-wise dashboard endpoints returning only chart fields, Alternatives considered, Consequences, Context

### Community 122 - "UserLocationService"
Cohesion: 0.15
Nodes (8): CreateUserLocationRequest, UpdateUserLocationRequest, validateUserLocation(), UserLocation, UserLocationRepository, mapUserLocationError(), UserLocationService, toUserLocationResponseError()

### Community 123 - "newUserServices"
Cohesion: 0.14
Nodes (14): UserMachineFilter, NewCompanyProfileRepository(), UserMachineFilterRepository, NewUserMachineFilterRepository(), NewUserProfileRepository(), NewCompanyProfileService(), newUserServices(), userServices (+6 more)

### Community 124 - "NewUserService"
Cohesion: 0.47
Nodes (6): RoleHandler, NewRoleHandler(), RoleRepository, RoleService, NewRoleService(), NewUserService()

### Community 126 - "Connection"
Cohesion: 0.09
Nodes (21): 5. Connect and collect, per machine: steps 10 – 17, Step 10 — Get ready to connect, Step 11 — Connect, a few machines at a time, Step 12 — Call the DLL's CNCConnect, Step 13 — The DLL sends data, Step 14 — The data collection loop, Step 16 — Every second: dead connection check and live status, Step 17 — Disconnect and try again (+13 more)

### Community 128 - "MachineService"
Cohesion: 0.06
Nodes (34): Decision, Decision, PlcReaderResponse, CreateMachineRequest, MachineFields, UpdateMachineRequest, validCoordinate(), validHealthIDs() (+26 more)

### Community 129 - "0010. Machine access by an organization filter, worked out on every request"
Cohesion: 0.40
Nodes (4): 0010. Machine access by an organization filter, worked out on every request, Alternatives considered, Consequences, Context

### Community 130 - "LicenceKeyService"
Cohesion: 0.10
Nodes (16): NewLicenceKey, CreateLicenceKeysRequest, GenerateLicenceKeysRequest, UpdateLicenceKeyRequest, validDuration(), ValidLicenceKey(), NewLicenceKeyHandler(), GeneratedLicenceKeys (+8 more)

### Community 132 - "0013. Machine diagnostics: reasons from what the collector recorded, network tested from the server"
Cohesion: 0.33
Nodes (5): 0013. Machine diagnostics: reasons from what the collector recorded, network tested from the server, Alternatives considered, Consequences, Context, Decision

### Community 133 - ".GetDiagnostics"
Cohesion: 0.07
Nodes (47): Decision, DiagnosticPLCVariable, DiagnosticCheck, DiagnosticMachine, DiagnosticPLC, DiagnosticVariableRow, MachineDiagnosticsResponse, MachineNetworkResponse (+39 more)

### Community 134 - "pathID"
Cohesion: 0.07
Nodes (15): MachineCycleHandler, MachineDiagnosticsHandler, pathID(), targetUserID(), SessionHandler, UserCellHandler, UserCityHandler, UserCountryHandler (+7 more)

### Community 138 - "net/http.Handler"
Cohesion: 0.13
Nodes (38): Handler(), NewMachineHandler(), CrossOrigin(), Logger(), Chain(), Middleware, Recoverer(), NewBrandRepository() (+30 more)

### Community 139 - "0025. A machine is saved at one level of its customer's organization; the path above is read, not stored"
Cohesion: 0.40
Nodes (4): 0025. A machine is saved at one level of its customer's organization; the path above is read, not stored, Alternatives considered, Consequences, Context

## Ambiguous Edges - Review These
- `CncCommunication.dll (32-bit vendor DLL)` → `FanucCncCommunication.dll`  [AMBIGUOUS]
  docs/setup-guide.txt · relation: semantically_similar_to
- `Machines API (/machines)` → `Snapshot Fields (positions, spindle, feed, program, tool, alarms)`  [AMBIGUOUS]
  docs/CncCommunication.dll — API & Integration Guide (Go Client).pdf · relation: conceptually_related_to

## Knowledge Gaps
- **33 isolated node(s):** `7th-sense-server`, `middleware.RequireRole / routes.adminOnly`, `response.SendResponse`, `Test Auth (asTestUser, tests/auth_test.go)`, `user_filters Filter Checklist` (+28 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 266 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **20 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `CNC Data Collection Flow` connect `Supervisor` to `Worker (internal/monitor/worker.go)`, `Live data — Socket.IO event contract`, `Connection`?**
  _High betweenness centrality (0.094) - this node is a cross-community bridge._
- **What connects `7th-sense-server`, `middleware.RequireRole / routes.adminOnly`, `response.SendResponse` to the rest of the system?**
  _33 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `apidocs/main.go` be split into smaller, more focused modules?**
  _Cohesion score 0.07740112994350283 - nodes in this community are weakly interconnected._
- **Why does `7th Sense Server Install and Setup Guide` connect `Worker (internal/monitor/worker.go)` to `Machines API (/machines)`, `Supervisor`?**
  _High betweenness centrality (0.083) - this node is a cross-community bridge._
- **Should `MachineModeReasonService` be split into smaller, more focused modules?**
  _Cohesion score 0.09871794871794871 - nodes in this community are weakly interconnected._

### Low-confidence Hints
_AMBIGUOUS edges — the extractor was unsure. Verify before acting on these._

- **What is the exact relationship between `CncCommunication.dll (32-bit vendor DLL)` and `FanucCncCommunication.dll`?**
  _Edge tagged AMBIGUOUS (relation: semantically_similar_to) - confidence is low._
- **What is the exact relationship between `Machines API (/machines)` and `Snapshot Fields (positions, spindle, feed, program, tool, alarms)`?**
  _Edge tagged AMBIGUOUS (relation: conceptually_related_to) - confidence is low._