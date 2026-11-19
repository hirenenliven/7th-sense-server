# Graph Report - 7th-sense-server  (2026-11-19)

## Corpus Check
- 370 files · ~305,708 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 4 file(s) not represented in the graph (top: (none) 3, .toml 1)

## Summary
- 3222 nodes · 30742 edges · 131 communities (109 shown, 22 thin omitted)
- Extraction: 96% EXTRACTED · 4% INFERRED · 0% AMBIGUOUS · INFERRED: 1333 edges (avg confidence: 0.87)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `165793b5`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- WithTransaction
- UserCityService
- UserLocationService
- MachineService
- net/http.Request
- pathID()
- Execute
- MachineModeReasonService
- api_helpers_test.go
- CncMasterSnapshot (34,440 bytes, #pragma
- LiveHandler
- DowntimeReasonService
- ProductionScheduleRepository
- Collector
- worker
- Hub
- run
- mode-chart.js
- .GetDiagnostics
- net/http.Handler
- net/http.ServeMux
- newDashboardRouter
- SendResponse
- §10 Diagnostics & Troubleshooting
- OperatorTaskService
- machine_mode.go
- LocationRepository
- .loadVariableConfig
- net/http.Request
- companyOwner
- PartService
- Connection
- Worker (internal/monitor/worker.go)
- README.md
- newUserServices
- MachineTelemetryService
- UserRepository
- Fixed values and why they are needed
- SessionRepository
- 0019. Axis values: every stored value when it fits, and zoom by loading the window again
- CncCommunication.dll API & Integration G
- AuthService
- ownCompanyCustomer
- CncCommunication.dll API & Integration G
- ProductionScheduleService
- SessionService
- UserShopService
- NewError
- Machines API (/machines)
- DashboardRequest
- Address
- validText
- BranchService
- README.md
- AuthHandler
- CompanyProfile
- MachineCycle
- MachineModelService
- apidocs/main.go
- MonitoredMachine
- UserProfileRepository
- Test list
- 0001. Go standard library net/http, no web framework
- 0003. PostgreSQL with pgx and hand-written SQL
- UserLineService
- plcWorker
- UserLineService
- run
- BrandService
- .GetAxisValues
- ControllerRepository
- onDataCallback(ptr uintptr) uintptr
- time.Time
- HomeDashboardResponse
- Brand
- testing.T
- 0008. "Log in as customer" uses a separate session with a path-scoped cookie
- OpenModePeriod
- UserStateService
- AuthUser
- Event-driven data flow (no polling in client)
- DashboardService
- Live data — how a machine value reaches the screen
- PartRepository
- CNCConnect(ip, port, dataCb, errCb, machineType) int32
- MachineTelemetryRepository
- auth_tokens SHA-256 Hashed Sessions (mig
- ClienteleRepository
- cycleReading
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
- SessionRepository
- 0020. Machine PLC node IDs: the machine's own node ID in machine_variables, default otherwise
- UserService
- 0011. Ponytail plugin for Claude Code, under this repo's rules
- Client Certificate & Key (client_cert.der / client_key.der)
- Decision
- user_access_test.go
- Role
- 0007. OpenAPI 3 from swag annotations, Swagger UI behind API_DOCS, generated Postman
- 0002. Layered architecture: handler, DTO, service, repository
- 0018. Machine Health X / Y / Z charts: downsampled in SQL, timed, on the machine-health routes
- HolidayRepository
- MachineCycleListRequest
- 0017. Machine accuracy as one shared list, managed by the superadmin
- MachineHealthRepository
- 0006. Page-wise dashboard endpoints returning only chart fields
- run
- 0021. Read PLC variables over OPC UA with gopcua, in-process beside the CNC collector
- MachineModeStatus
- CanAssignRole
- SessionRevokedError
- 0016. A Task List screen for the operator tasks' master list
- LicenceKeyService
- ClienteleService
- 0010. Machine access by an organization filter, worked out on every request

## God Nodes (most connected - your core abstractions)
1. `Execute()` - 225 edges
2. `SendResponse()` - 190 edges
3. `NewError()` - 159 edges
4. `run()` - 119 edges
5. `AuthUser` - 112 edges
6. `pathID()` - 98 edges
7. `currentUser()` - 84 edges
8. `Querier` - 79 edges
9. `decodeAndValidate()` - 68 edges
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
- **Session authentication and authorization flow** — readme_session_cookie, readme_auth_tokens_hashed, readme_per_request_session_check, readme_login_lockout, readme_csrf_cross_origin, _claude_claude_role_guards [INFERRED 0.85]

## Communities (131 total, 22 thin omitted)

### Community 0 - "WithTransaction"
Cohesion: 0.07
Nodes (30): BreakFields, CreateBreakRequest, UpdateBreakRequest, clockMinutes(), CreateShiftRequest, ShiftFields, UpdateShiftRequest, MinutesBetween() (+22 more)

### Community 1 - "UserCityService"
Cohesion: 0.16
Nodes (8): CreateUserCityRequest, UpdateUserCityRequest, validateUserCity(), UserCity, UserCityRepository, mapUserCityError(), UserCityService, toUserCityResponseError()

### Community 2 - "UserLocationService"
Cohesion: 0.15
Nodes (8): CreateUserLocationRequest, UpdateUserLocationRequest, validateUserLocation(), UserLocation, UserLocationRepository, mapUserLocationError(), UserLocationService, toUserLocationResponseError()

### Community 3 - "MachineService"
Cohesion: 0.05
Nodes (39): Decision, Decision, MachineAccuracyRequest, MachineNetworkResponse, CreateMachineRequest, MachineFields, UpdateMachineRequest, validCoordinate() (+31 more)

### Community 4 - "net/http.Request"
Cohesion: 0.07
Nodes (17): BreakHandler, ClienteleHandler, DowntimeReasonHandler, HolidayHandler, IdleReasonHandler, MachineDiagnosticsHandler, MachineHandler, MachineModeReasonHandler (+9 more)

### Community 5 - "pathID()"
Cohesion: 0.08
Nodes (11): targetUserID(), RoleHandler, UserCellHandler, UserCityHandler, UserCountryHandler, UserLineHandler, UserLocationHandler, UserShopHandler (+3 more)

### Community 6 - "Execute"
Cohesion: 0.21
Nodes (4): Execute(), MachineDiagnosticsRepository, MachineModeReasonRepository, mapMachineModeReasonError()

### Community 7 - "MachineModeReasonService"
Cohesion: 0.14
Nodes (17): 0012. Downtime and idle reasons on a period: one table, tied to machine_mode_log, Alternatives considered, Consequences, Context, Decision, CreateMachineModeReasonRequest, MachineModeReasonFields, MachineModeReasonsResponse (+9 more)

### Community 8 - "api_helpers_test.go"
Cohesion: 0.42
Nodes (8): cncMasterSnapshot, RevokeSessionsResponse, authUserKey, loginAsKey, sessionIDKey, liveValues, machineModelRecord, subscribeAnswer

### Community 9 - "CncMasterSnapshot (34,440 bytes, #pragma"
Cohesion: 0.10
Nodes (43): CncMasterSnapshot (34,440 bytes, #pragma pack(8), x86), AbsPos ([5]float64, Inches / Degrees, Work Coordinate System (WCS) positions), ActiveGCode ([16]byte, ASCII, Currently active G-code group), ActiveProgPath ([256]byte, ASCII Path, Storage path of active NC artifact), AlarmMsg ([128]byte, UTF-8 / ASCII, Alarm description message), AlarmNo (int32, Code (0 = OK), Primary alarm or fault number), AxisLoadXPercentage (int16, Percent (%), Servo X axis drive load (AxisLoadX/Y/ZPercentage)), AxisLoadYPercentage (int16, Percent (%), Servo Y axis drive load (AxisLoadX/Y/ZPercentage)) (+35 more)

### Community 10 - "LiveHandler"
Cohesion: 0.10
Nodes (13): Browser → server, liveSocket, liveSubscription, LiveSessionExpired, channelRoom(), closeLiveSocket(), decodeLiveMessage(), LiveHandler (+5 more)

### Community 11 - "DowntimeReasonService"
Cohesion: 0.15
Nodes (11): CreateDowntimeReasonRequest, DowntimeReasonFields, UpdateDowntimeReasonRequest, NewDowntimeReasonHandler(), DowntimeReason, DowntimeReasonRepository, mapDowntimeReasonError(), downtimeReasonFromFields() (+3 more)

### Community 12 - "ProductionScheduleRepository"
Cohesion: 0.10
Nodes (7): MachineQuantity, OperatorJob, ProductionScheduleItem, executeAffected(), ProductionScheduleRepository, mapProductionScheduleError(), levelInMachineFilter()

### Community 13 - "Collector"
Cohesion: 0.16
Nodes (16): 2. All steps at a glance, 7. Shutdown: steps 19 – 20, Log lines (in `logs/app.log`, search for the message), Step 18 — React to machine and variable changes, Step 19 — Stop collection with the server, Step 20 — Stop the workers and mark machines offline, Step 3 — Run the collector in the background, Step 6 — Record the start and begin following the database (+8 more)

### Community 14 - "worker"
Cohesion: 0.10
Nodes (17): 5. Connect and collect, per machine: steps 10 – 17, Known issue, Retry and timing summary, Step 10 — Get ready to connect, Step 11 — Connect, a few machines at a time, Step 14 — The data collection loop, Step 15 — Store the reading, Step 16 — Every second: dead connection check and live status (+9 more)

### Community 15 - "Hub"
Cohesion: 0.04
Nodes (70): Supervisor, PLCConfig, 0013. Machine diagnostics: reasons from what the collector recorded, network tested from the server, Alternatives considered, Consequences, Context, Decision, Decision (+62 more)

### Community 16 - "run"
Cohesion: 0.05
Nodes (114): run(), Handler(), Querier, NewBranchHandler(), NewBrandHandler(), NewBreakHandler(), NewCategoryHandler(), NewClienteleHandler() (+106 more)

### Community 17 - "mode-chart.js"
Cohesion: 0.10
Nodes (64): card(), cell(), chartData(), chartRows(), countByGroup(), currentSpan(), currentWindow(), cycleData() (+56 more)

### Community 18 - ".GetDiagnostics"
Cohesion: 0.08
Nodes (42): DiagnosticPLCVariable, DiagnosticCheck, DiagnosticMachine, DiagnosticPLC, DiagnosticVariableRow, MachineDiagnosticsResponse, DiagnosticEvent, DiagnosticLiveStatus (+34 more)

### Community 19 - "net/http.Handler"
Cohesion: 0.15
Nodes (29): authRequest(), createAuthTestUser(), login(), newAuthRouter(), sessionCookieOf(), TestAuthentication(), dialLive(), newLiveTestServer() (+21 more)

### Community 20 - "net/http.ServeMux"
Cohesion: 0.12
Nodes (41): Decision, RequireRole(), registerAuthRoutes(), registerBranchRoutes(), registerBrandRoutes(), registerBreakRoutes(), registerCategoryRoutes(), registerClienteleRoutes() (+33 more)

### Community 21 - "newDashboardRouter"
Cohesion: 0.27
Nodes (10): NewDashboardHandler(), NewMachineCycleRepository(), NewMachineHealthRepository(), NewMachineCycleService(), mustDecode(), newDashboardRouter(), TestAxisValues(), TestDashboardCharts() (+2 more)

### Community 22 - "SendResponse"
Cohesion: 0.06
Nodes (18): validator, BranchHandler, BrandHandler, CategoryHandler, ControllerHandler, HealthHandler, MachineAccuracyHandler, NewMachineAccuracyHandler() (+10 more)

### Community 23 - "§10 Diagnostics & Troubleshooting"
Cohesion: 0.12
Nodes (25): Check mode: go run sinumerik_client.go check (OPC UA node existence report), Go Client (sinumerik_client.go, callback -> copy -> render), Prerequisite: 32-bit x86 DLL, compile Go with GOARCH=386, Prerequisite: SIEMENS_USER / SIEMENS_PASS env vars before CNCConnect, Run: go run sinumerik_client.go (live dashboard), §3.1 Environment & Prerequisites, §8 Running the Go Client, §10 Diagnostics & Troubleshooting (+17 more)

### Community 24 - "OperatorTaskService"
Cohesion: 0.09
Nodes (14): CreateMaintenanceTaskRequest, CreateOperatorTaskRequest, OperatorTaskFields, UpdateMaintenanceTaskRequest, UpdateOperatorTaskRequest, MaintenanceTask, OperatorTask, mapMaintenanceTaskError() (+6 more)

### Community 25 - "machine_mode.go"
Cohesion: 0.14
Nodes (20): decideMachineMode(), holdReason(), machineStateFromReading(), numberingEngine(), programStateFromProcStatus(), ncReading(), TestDecideMachineMode(), TestDecideMachineModeConditions() (+12 more)

### Community 26 - "LocationRepository"
Cohesion: 0.50
Nodes (3): mapLocationError(), NewLocationRepository(), LocationRepository

### Community 27 - ".loadVariableConfig"
Cohesion: 0.24
Nodes (7): VariableConfig, formatPreviousValue(), formatValue(), newValueChangeLog(), boundVariable, valueChange, valueChangeLog

### Community 28 - "net/http.Request"
Cohesion: 0.09
Nodes (14): axisValuesQuery(), dashboardQuery(), DashboardHandler, queryIDList(), readDashboardRequest(), LicenceKeyHandler, isSameOrigin(), MachineCycleHandler (+6 more)

### Community 29 - "companyOwner"
Cohesion: 0.15
Nodes (11): CreateIdleReasonRequest, IdleReasonFields, UpdateIdleReasonRequest, NewIdleReasonHandler(), IdleReason, IdleReasonRepository, mapIdleReasonError(), IdleReasonService (+3 more)

### Community 30 - "PartService"
Cohesion: 0.17
Nodes (10): ValidPartDrawingFile(), Part, PartDrawingStorage, NewPartRepository(), asInvalidSelection(), PartService, NewPartService(), partFromFields() (+2 more)

### Community 31 - "Connection"
Cohesion: 0.08
Nodes (20): dllDriver, slot, Step 12 — Call the DLL's CNCConnect, Step 13 — The DLL sends data, Step 17 — Disconnect and try again, 1. Machine → DLL → server, Connection, Reading (+12 more)

### Community 32 - "Worker (internal/monitor/worker.go)"
Cohesion: 0.05
Nodes (41): CNC Data Collection (in-process), writeLog CNC Log Types and Switches, .env Configuration (internal/config), Health Mode History API (public), 7th Sense Server Claude Context, response.SendResponse Format, User CRUD Golden Reference, 1. The big picture (+33 more)

### Community 33 - "README.md"
Cohesion: 0.24
Nodes (5): 0015. Operator tasks: a shared task list, assigned per customer and machine model, Alternatives considered, Consequences, Context, Architecture Decision Records

### Community 34 - "newUserServices"
Cohesion: 0.16
Nodes (13): UserMachineFilter, NewAddressRepository(), NewCompanyProfileRepository(), UserMachineFilterRepository, NewUserMachineFilterRepository(), NewAddressService(), NewCompanyProfileService(), newUserServices() (+5 more)

### Community 35 - "MachineTelemetryService"
Cohesion: 0.16
Nodes (6): 4. Find the machines: steps 5 – 9, Step 5 — Prepare the database, Step 7 — Listen for database changes, Step 8 — Load all active machines, ListenForDatabaseNotifications(), MachineTelemetryService

### Community 36 - "UserRepository"
Cohesion: 0.22
Nodes (6): Decision, UserListFilter, User, UserRepository, mapUserError(), NewClienteleService()

### Community 37 - "Fixed values and why they are needed"
Cohesion: 0.10
Nodes (20): 1. The rows are consistent (SQL, no machine needed), 2. The rows agree with the mode log (SQL), 3. The raw signals around one cycle (telemetry and logs), 4. At the machine (floor test), `CNC_CYCLE_MAX_HOURS` (24 hours), `CNC_STALE_SEC` (30 s) and `CNC_QUIET_SEC` (300 s), Completed, `cycleCounterGrace` and `cycleStartGrace` (5 seconds each) (+12 more)

### Community 38 - "SessionRepository"
Cohesion: 0.24
Nodes (4): SessionState, SessionUser, AuthRepository, mapAuthError()

### Community 39 - "0019. Axis values: every stored value when it fits, and zoom by loading the window again"
Cohesion: 0.40
Nodes (5): 0019. Axis values: every stored value when it fits, and zoom by loading the window again, Alternatives considered, Consequences, Context, Decision

### Community 40 - "CncCommunication.dll API & Integration G"
Cohesion: 0.19
Nodes (20): CncCommunication.dll (C++ core runtime, x86, OPC UA inside), CncCommunication.dll API & Integration Guide (Go Client) - text version, Fanuc profile (engine 1, port 8193, Fwlib32.dll + fwlibe1.dll), Fwlib32.dll (Fanuc FOCAS2), fwlibe1.dll (Fanuc FOCAS2), libcrypto-3.dll (OpenSSL), libssl-3.dll (OpenSSL), Engine 1: Fanuc (FOCAS2 Ethernet, fwlib32) (+12 more)

### Community 41 - "AuthService"
Cohesion: 0.08
Nodes (24): 0014. Session modes and the login conflict, on the existing auth_tokens sessions, Alternatives considered, Consequences, Context, Decision, SessionConfig, LoginRequest, SessionClient (+16 more)

### Community 42 - "ownCompanyCustomer"
Cohesion: 0.22
Nodes (9): CreateHolidayRequest, HolidayFields, UpdateHolidayRequest, Holiday, ownCompanyCustomer(), HolidayService, holidayFromFields(), toHolidayResponseError() (+1 more)

### Community 43 - "CncCommunication.dll API & Integration G"
Cohesion: 0.13
Nodes (18): CNCConnect param machineType (int32 engine id, see §4.2), Engine 0: Simulator (local loopback mock engine), Engine 2: Siemens (Sinumerik OPC UA, active profile), Engine 3: Syntec (Open API / TCP), Engine 4: Mitsubishi (driver not specified), Machine identifier mapping (engine id enum 0-4), §9 Controller-Specific Profiles, Siemens default ports: 4840 (native NCU) / 4841 (custom) (+10 more)

### Community 44 - "ProductionScheduleService"
Cohesion: 0.16
Nodes (17): validDate(), CreateProductionScheduleRequest, ProductionScheduleFields, ProductionScheduleItemFields, UpdateProductionScheduleRequest, ScheduleDay(), ProductionSchedule, NewProductionScheduleRepository() (+9 more)

### Community 45 - "SessionService"
Cohesion: 0.13
Nodes (9): describeDevice(), SessionConflictResponse, SessionResponse, NewSessionResponse(), NewSessionResponses(), Session, SessionService, toSessionResponseError() (+1 more)

### Community 46 - "UserShopService"
Cohesion: 0.16
Nodes (8): CreateUserShopRequest, UpdateUserShopRequest, validateUserShop(), UserShop, UserShopRepository, mapUserShopError(), UserShopService, toUserShopResponseError()

### Community 47 - "NewError"
Cohesion: 0.09
Nodes (25): CreateAddressRequest, UpdateAddressRequest, UpdateLanguageRequest, UpdateThemeRequest, UpdateLicenceKeyRequest, validDuration(), ValidLicenceKey(), MachineFilter (+17 more)

### Community 48 - "Machines API (/machines)"
Cohesion: 0.16
Nodes (18): Clienteles API, Company Scope (companyOwner / ownCompanyCustomer), Current User API (/me theme, language), Customer Display Names (customerName / companyName), Holidays API, Licence Keys API, Machine Catalog Brand-Controller-Model, Machines API (/machines) (+10 more)

### Community 49 - "DashboardRequest"
Cohesion: 0.18
Nodes (14): TimelinePeriod, ModeTimelineResponse, NewTimelinePeriod(), MachineModeHistoryResponse, MachineModeTotal, MachineModePeriod, clipWindows(), machineDayOEE() (+6 more)

### Community 50 - "Address"
Cohesion: 0.27
Nodes (5): Address, AddressRepository, mapAddressError(), toAddressResponseError(), AddressService

### Community 51 - "validText"
Cohesion: 0.17
Nodes (9): PartProcessFields, PartSubProgramFields, CreatePartRequest, PartFields, UpdatePartRequest, UploadPartDrawingRequest, uniquePositiveIDs(), validText() (+1 more)

### Community 52 - "BranchService"
Cohesion: 0.14
Nodes (9): CreateBranchRequest, UpdateBranchRequest, validateBranchFields(), Branch, BranchRepository, mapBranchError(), cityOwner(), BranchService (+1 more)

### Community 53 - "README.md"
Cohesion: 0.40
Nodes (4): 0004. Opaque session token in an HttpOnly cookie, Alternatives considered, Consequences, Context

### Community 54 - "AuthHandler"
Cohesion: 0.09
Nodes (18): Decision, 0009. Socket.IO live data: the collector pushes each screen's variables, checked per user, Alternatives considered, Consequences, Context, Decision, AuthHandler, CurrentUser() (+10 more)

### Community 55 - "CompanyProfile"
Cohesion: 0.28
Nodes (5): CompanyProfile, CompanyProfileRepository, mapCompanyProfileError(), toCompanyProfileResponseError(), CompanyProfileService

### Community 56 - "MachineCycle"
Cohesion: 0.15
Nodes (8): MachineCyclesResponse, CycleJob, MachineCycle, MachineCycleRepository, mapMachineCycleError(), MachineCycleService, sameID(), toMachineCycleResponseError()

### Community 57 - "MachineModelService"
Cohesion: 0.20
Nodes (6): MachineModelRequest, MachineModel, MachineModelRepository, mapMachineModelError(), MachineModelService, toMachineModelResponseError()

### Community 58 - "apidocs/main.go"
Cohesion: 0.05
Nodes (53): request, main(), run(), main(), Config, MachineSelection, EnsureCertificates(), existingCertificateExpiry() (+45 more)

### Community 59 - "MonitoredMachine"
Cohesion: 0.21
Nodes (7): equalEngine(), MonitoredMachine, samePLCConnection(), newPlcWorker(), newMachineWorker(), toMachineTelemetryResponseError(), PlcCollector

### Community 60 - "UserProfileRepository"
Cohesion: 0.14
Nodes (10): UserProfile, UserProfileRepository, mapProfileError(), NewUserProfileRepository(), NewUserProfileService(), toUserProfileResponseError(), UserProfileService, sameClock() (+2 more)

### Community 61 - "Test list"
Cohesion: 0.08
Nodes (26): Before the tests (once), Check query, Floor tests on the HUMA panel (Fanuc), HUMA buttons used, Result sheet, TEST-01 · Start a cycle, TEST-02 · Normal end (M30), TEST-03 · Three parts one after another (+18 more)

### Community 62 - "0001. Go standard library net/http, no web framework"
Cohesion: 0.33
Nodes (5): 0001. Go standard library net/http, no web framework, Alternatives considered, Consequences, Context, Decision

### Community 63 - "0003. PostgreSQL with pgx and hand-written SQL"
Cohesion: 0.33
Nodes (5): 0003. PostgreSQL with pgx and hand-written SQL, Alternatives considered, Consequences, Context, Decision

### Community 64 - "UserLineService"
Cohesion: 0.15
Nodes (8): CreateUserLineRequest, UpdateUserLineRequest, validateUserLine(), UserLine, UserLineRepository, mapUserLineError(), UserLineService, toUserLineResponseError()

### Community 65 - "plcWorker"
Cohesion: 0.22
Nodes (4): previousText(), Endpoint(), backoff, plcWorker

### Community 66 - "UserLineService"
Cohesion: 0.07
Nodes (20): CreateUserCellRequest, UpdateUserCellRequest, validateUserCell(), CreateUserSupervisorRequest, UpdateUserSupervisorRequest, validateUserSupervisor(), NewUserCellHandler(), NewUserSupervisorHandler() (+12 more)

### Community 67 - "run"
Cohesion: 0.36
Nodes (10): checkNotTaken(), createSuperadmin(), ensureProfile(), execute(), main(), parseFlags(), run(), superadminRoleID() (+2 more)

### Community 68 - "BrandService"
Cohesion: 0.20
Nodes (5): Brand, BrandRepository, mapBrandError(), BrandService, toBrandResponseError()

### Community 69 - ".GetAxisValues"
Cohesion: 0.22
Nodes (9): AxisPoint, FindAxisSet(), AxisSeries, AxisSet, AxisStats, AxisValuesRequest, AxisValuesResponse, minTimeOf() (+1 more)

### Community 70 - "ControllerRepository"
Cohesion: 0.15
Nodes (8): CreateControllerRequest, UpdateControllerRequest, validateControllerFields(), Controller, ControllerRepository, mapControllerError(), ControllerService, toControllerResponseError()

### Community 71 - "onDataCallback(ptr uintptr) uintptr"
Cohesion: 0.36
Nodes (9): Package-level callback singletons (syscall.NewCallback), CNCConnect param errCb (uintptr stdcall OnErr(int32 code, const char* msg)), Mutex-protected copy: dataMutex sync.RWMutex -> lastRecord/hasRecord, onDataCallback(ptr uintptr) uintptr, onErrorCallback(code, messagePtr uintptr) uintptr, §7 Asynchronous Callback Architecture, onDataCallback, onErrorCallback (+1 more)

### Community 72 - "time.Time"
Cohesion: 0.09
Nodes (18): VirtualScreenCycle, VirtualScreenLive, ModeTotalsResponse, AxisBucketSeconds(), VirtualScreenResponse, NewVirtualScreenLive(), AxisBucket, AxisRow (+10 more)

### Community 73 - "HomeDashboardResponse"
Cohesion: 0.08
Nodes (24): HomeAlarm, HomeCycles, HomeFleet, HomeMachine, HomePlan, HomeTime, HomeDashboardResponse, HomeHour (+16 more)

### Community 74 - "Brand"
Cohesion: 0.07
Nodes (19): CreateBrandRequest, UpdateBrandRequest, validateCatalogFields(), CreateCategoryRequest, UpdateCategoryRequest, CreateProductRequest, UpdateProductRequest, validateProductFields() (+11 more)

### Community 75 - "testing.T"
Cohesion: 0.07
Nodes (102): TestSnapshotModeFieldOffsets(), TestSnapshotSize(), containsID(), decodeData(), expectStatus(), newCatalogRouter(), runCatalogCRUD(), sendRequest() (+94 more)

### Community 76 - "0008. "Log in as customer" uses a separate session with a path-scoped cookie"
Cohesion: 0.40
Nodes (4): 0008. "Log in as customer" uses a separate session with a path-scoped cookie, Alternatives considered, Consequences, Context

### Community 78 - "UserStateService"
Cohesion: 0.08
Nodes (14): CreateUserCountryRequest, CreateUserStateRequest, validateUserState(), UserCountry, UserState, UserCountryRepository, mapUserCountryError(), UserStateRepository (+6 more)

### Community 79 - "AuthUser"
Cohesion: 0.23
Nodes (4): LoginAsResponse, OperatorJobFilter, AuthUser, companyOwner()

### Community 80 - "Event-driven data flow (no polling in client)"
Cohesion: 0.60
Nodes (5): Flow 1: load DLL, register two stdcall callbacks (syscall.NewCallback), Flow 2: CNCConnect starts DLL background worker thread, Flow 3: worker polls controller, assembles snapshot, invokes data callback, Flow 4: callback copies payload, signals renderer via buffered channel, Flow 5: main goroutine renders live ANSI dashboard

### Community 81 - "DashboardService"
Cohesion: 0.21
Nodes (12): CycleCountRow, MachineTimeRow, cyclesIn(), dayWindows(), DashboardService, hours(), modeGroup(), round1() (+4 more)

### Community 82 - "Live data — how a machine value reaches the screen"
Cohesion: 0.20
Nodes (7): Adding a variable to a screen, Connecting, Live data — Socket.IO event contract, Ordering, Rooms (server side), Server → browser, Settings

### Community 83 - "PartRepository"
Cohesion: 0.23
Nodes (4): PartProcess, PartSubProgram, PartRepository, mapPartError()

### Community 84 - "CNCConnect(ip, port, dataCb, errCb, machineType) int32"
Cohesion: 0.20
Nodes (11): CNCConnect(ip, port, dataCb, errCb, machineType) int32, CNCConnect param dataCb (uintptr stdcall OnData(const CncMasterSnapshot*)), CNCConnect param ip (const char*, ANSI NUL-terminated), CNCConnect param port (int32, e.g. 4840), CNCdisconnect(machineId int32), LazyDLL binding (syscall.NewLazyDLL + eager mod.Load in init), §4 API Reference, §5 DLL Loading & Binding (Go) (+3 more)

### Community 85 - "MachineTelemetryRepository"
Cohesion: 0.08
Nodes (10): LiveValue, MachineEvent, MachineLiveStatus, ModeChange, TelemetryValue, MachineTelemetryRepository, mapMachineTelemetryError(), execute() (+2 more)

### Community 86 - "auth_tokens SHA-256 Hashed Sessions (mig"
Cohesion: 0.67
Nodes (3): auth_tokens SHA-256 Hashed Sessions (migration 014), Per-Request Session and User Re-check, ts_session Cookie (HttpOnly, SameSite=Strict)

### Community 89 - "cycleReading"
Cohesion: 0.21
Nodes (7): CycleStatusFor(), decideCycleStep(), isProgramLoop(), resetReason(), cycleReading, cycleStep, cycleTracker

### Community 104 - "connection.go"
Cohesion: 0.08
Nodes (17): changeOf(), Connect(), forwardNotifications(), Change, Item, ItemResult, isBad(), securityOf() (+9 more)

### Community 105 - "0005. CNC collector runs inside the API process"
Cohesion: 0.33
Nodes (5): 0005. CNC collector runs inside the API process, Alternatives considered, Consequences, Context, Decision

### Community 107 - "0020. Machine PLC node IDs: the machine's own node ID in machine_variables, default otherwise"
Cohesion: 0.33
Nodes (5): 0020. Machine PLC node IDs: the machine's own node ID in machine_variables, default otherwise, Alternatives considered, Consequences, Context, Decision

### Community 108 - "UserService"
Cohesion: 0.16
Nodes (15): WithTransaction(), CreateUserRequest, CustomerDetails, UpdateUserRequest, UserResponse, ParseDayStart(), TestParseDayStart(), canSee() (+7 more)

### Community 109 - "0011. Ponytail plugin for Claude Code, under this repo's rules"
Cohesion: 0.33
Nodes (5): 0011. Ponytail plugin for Claude Code, under this repo's rules, Alternatives considered, Consequences, Context, Decision

### Community 110 - "Client Certificate & Key (client_cert.der / client_key.der)"
Cohesion: 0.48
Nodes (7): client_cert.der/.pem (X.509 client cert for Siemens OPC UA), client_key.der/.pem (client private key), Siemens certs: client_cert.der + client_key.der in binary directory, Siemens trust: cert must be trusted in NCU Access Administration, Troubleshoot: BadRequestNotAllowed -> Export client_cert, BadRequestNotAllowed (certificate trust rejection), Client Certificate & Key (client_cert.der / client_key.der)

### Community 111 - "Decision"
Cohesion: 0.40
Nodes (6): Decision, FinishCyclesSummary, FinishedCycle, FinishCyclesResponse, NewFinishCyclesSummary(), NewFinishedCycle()

### Community 112 - "user_access_test.go"
Cohesion: 0.33
Nodes (7): LanguageResponse, LicenceKeyDuration, LiveAccessChanged, LiveUnsubscribeRequest, ThemeResponse, licenceKeyDurations(), Envelope

### Community 113 - "Role"
Cohesion: 0.33
Nodes (3): Role, mapRoleError(), toRoleResponseError()

### Community 114 - "0007. OpenAPI 3 from swag annotations, Swagger UI behind API_DOCS, generated Postman"
Cohesion: 0.33
Nodes (5): 0007. OpenAPI 3 from swag annotations, Swagger UI behind API_DOCS, generated Postman, Alternatives considered, Consequences, Context, Decision

### Community 115 - "0002. Layered architecture: handler, DTO, service, repository"
Cohesion: 0.33
Nodes (5): 0002. Layered architecture: handler, DTO, service, repository, Alternatives considered, Consequences, Context, Decision

### Community 116 - "0018. Machine Health X / Y / Z charts: downsampled in SQL, timed, on the machine-health routes"
Cohesion: 0.50
Nodes (4): 0018. Machine Health X / Y / Z charts: downsampled in SQL, timed, on the machine-health routes, Alternatives considered, Consequences, Context

### Community 118 - "MachineCycleListRequest"
Cohesion: 0.33
Nodes (3): DashboardRequest, HomeDashboardRequest, MachineCycleListRequest

### Community 119 - "0017. Machine accuracy as one shared list, managed by the superadmin"
Cohesion: 0.40
Nodes (4): 0017. Machine accuracy as one shared list, managed by the superadmin, Alternatives considered, Consequences, Context

### Community 120 - "MachineHealthRepository"
Cohesion: 0.40
Nodes (4): Decision, AxisVariable, MachineHealthRepository, NewDashboardService()

### Community 121 - "0006. Page-wise dashboard endpoints returning only chart fields"
Cohesion: 0.40
Nodes (4): 0006. Page-wise dashboard endpoints returning only chart fields, Alternatives considered, Consequences, Context

### Community 122 - "run"
Cohesion: 0.36
Nodes (4): connectFanucCounter(), main(), run(), fanucCounter

### Community 123 - "0021. Read PLC variables over OPC UA with gopcua, in-process beside the CNC collector"
Cohesion: 0.40
Nodes (4): 0021. Read PLC variables over OPC UA with gopcua, in-process beside the CNC collector, Alternatives considered, Consequences, Context

### Community 129 - "0016. A Task List screen for the operator tasks' master list"
Cohesion: 0.40
Nodes (5): 0016. A Task List screen for the operator tasks' master list, Alternatives considered, Consequences, Context, Decision

### Community 130 - "LicenceKeyService"
Cohesion: 0.14
Nodes (11): NewLicenceKey, CreateLicenceKeysRequest, GenerateLicenceKeysRequest, GeneratedLicenceKeys, LicenceKey, LicenceKeyRepository, mapLicenceKeyError(), LicenceKeyService (+3 more)

### Community 131 - "ClienteleService"
Cohesion: 0.23
Nodes (8): ClienteleFields, CreateClienteleRequest, UpdateClienteleRequest, validClienteleEmail(), Clientele, clienteleFromFields(), ClienteleService, toClienteleResponseError()

### Community 139 - "0010. Machine access by an organization filter, worked out on every request"
Cohesion: 0.40
Nodes (4): 0010. Machine access by an organization filter, worked out on every request, Alternatives considered, Consequences, Context

## Ambiguous Edges - Review These
- `CncCommunication.dll (32-bit vendor DLL)` → `FanucCncCommunication.dll`  [AMBIGUOUS]
  docs/setup-guide.txt · relation: semantically_similar_to
- `Machines API (/machines)` → `Snapshot Fields (positions, spindle, feed, program, tool, alarms)`  [AMBIGUOUS]
  docs/CncCommunication.dll — API & Integration Guide (Go Client).pdf · relation: conceptually_related_to

## Knowledge Gaps
- **33 isolated node(s):** `7th-sense-server`, `middleware.RequireRole / routes.adminOnly`, `response.SendResponse`, `Test Auth (asTestUser, tests/auth_test.go)`, `user_filters Filter Checklist` (+28 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 243 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **22 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `CNC Data Collection Flow` connect `Worker (internal/monitor/worker.go)` to `MachineTelemetryService`, `Collector`, `worker`, `Hub`, `Live data — how a machine value reaches the screen`?**
  _High betweenness centrality (0.121) - this node is a cross-community bridge._
- **What connects `7th-sense-server`, `middleware.RequireRole / routes.adminOnly`, `response.SendResponse` to the rest of the system?**
  _33 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `WithTransaction` be split into smaller, more focused modules?**
  _Cohesion score 0.07103825136612021 - nodes in this community are weakly interconnected._
- **Why does `7th Sense Server Install and Setup Guide` connect `Worker (internal/monitor/worker.go)` to `Machines API (/machines)`?**
  _High betweenness centrality (0.103) - this node is a cross-community bridge._
- **Should `MachineService` be split into smaller, more focused modules?**
  _Cohesion score 0.050239234449760764 - nodes in this community are weakly interconnected._

### Low-confidence Hints
_AMBIGUOUS edges — the extractor was unsure. Verify before acting on these._

- **What is the exact relationship between `CncCommunication.dll (32-bit vendor DLL)` and `FanucCncCommunication.dll`?**
  _Edge tagged AMBIGUOUS (relation: semantically_similar_to) - confidence is low._
- **What is the exact relationship between `Machines API (/machines)` and `Snapshot Fields (positions, spindle, feed, program, tool, alarms)`?**
  _Edge tagged AMBIGUOUS (relation: conceptually_related_to) - confidence is low._