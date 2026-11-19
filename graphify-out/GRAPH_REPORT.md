# Graph Report - 7th-sense-server  (2026-11-19)

## Corpus Check
- 370 files · ~305,708 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 4 file(s) not represented in the graph (top: (none) 3, .toml 1)

## Summary
- 3223 nodes · 30690 edges · 139 communities (114 shown, 25 thin omitted)
- Extraction: 96% EXTRACTED · 4% INFERRED · 0% AMBIGUOUS · INFERRED: 1327 edges (avg confidence: 0.87)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `188077aa`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- WithTransaction
- UserCityService
- UserLocationService
- Machine
- net/http.Request
- newLiveTestServerWith
- Execute
- MachineModeReasonService
- api_helpers_test.go
- CncMasterSnapshot (34,440 bytes, #pragma
- LiveHandler
- config.go
- ProductionScheduleRepository
- Collector
- worker
- Hub
- run
- mode-chart.js
- .GetDiagnostics
- net/http.Handler
- net/http.ServeMux
- UserSupervisorService
- SendResponse
- §10 Diagnostics & Troubleshooting
- OperatorTaskService
- machine_mode.go
- LocationRepository
- .loadVariableConfig
- NewError
- companyOwner
- ownCompanyCustomer
- Connection
- Worker (internal/monitor/worker.go)
- README.md
- newUserServices
- MachineTelemetryService
- UserRepository
- .collectReadings
- SessionRepository
- 0019. Axis values: every stored value when it fits, and zoom by loading the window again
- CncCommunication.dll API & Integration G
- AuthService
- ownCompanyCustomer
- CncCommunication.dll API & Integration G
- ProductionScheduleService
- SessionService
- UserShopService
- .Validate
- Machines API (/machines)
- DashboardRequest
- Address
- validText
- .Create
- README.md
- AuthHandler
- CompanyProfile
- MachineCycle
- MachineModelService
- apidocs/main.go
- generate.go
- UserProfileRepository
- Test list
- 0001. Go standard library net/http, no web framework
- 0003. PostgreSQL with pgx and hand-written SQL
- MachinePlcVariableService
- time.Duration
- UserLineService
- run
- MachineAccuracyService
- UserHandler
- ControllerRepository
- onDataCallback(ptr uintptr) uintptr
- time.Time
- HomeDashboardResponse
- Brand
- testing.T
- 0008. "Log in as customer" uses a separate session with a path-scoped cookie
- WithTransaction
- UserCountryRepository
- AuthUser
- Event-driven data flow (no polling in client)
- DashboardService
- Live data — how a machine value reaches the screen
- PartRepository
- CNCConnect(ip, port, dataCb, errCb, machineType) int32
- context.Context
- auth_tokens SHA-256 Hashed Sessions (mig
- ClienteleRepository
- HealthHandler
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
- .Update
- 0011. Ponytail plugin for Claude Code, under this repo's rules
- Client Certificate & Key (client_cert.der / client_key.der)
- MachineCycleSummary
- user_access_test.go
- Role
- 0007. OpenAPI 3 from swag annotations, Swagger UI behind API_DOCS, generated Postman
- 0002. Layered architecture: handler, DTO, service, repository
- 0018. Machine Health X / Y / Z charts: downsampled in SQL, timed, on the machine-health routes
- SessionHandler
- MachineCycleListRequest
- 0017. Machine accuracy as one shared list, managed by the superadmin
- BranchRepository
- 0006. Page-wise dashboard endpoints returning only chart fields
- run
- connection_test.go
- .validate
- UserService
- .Create
- SessionRevokedError
- logger/logger.go
- 0016. A Task List screen for the operator tasks' master list
- LicenceKeyService
- ClienteleService
- 0012. Downtime and idle reasons on a period: one table, tied to machine_mode_log
- PartDrawingStorage
- MachinePlcVariableHandler
- UserState
- 0013. Machine diagnostics: reasons from what the collector recorded, network tested from the server
- .GetOperatorJobs
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
- `Alternatives considered` --references--> `managersOnly()`  [INFERRED]
  docs/adr/0012-machine-mode-reasons-one-table.md → internal/routes/routes.go
- `Decision` --references--> `superAdminOnly()`  [INFERRED]
  docs/adr/0016-operator-task-list-screen.md → internal/routes/routes.go
- `Consequences` --references--> `explainNotConnected()`  [INFERRED]
  docs/adr/0013-machine-diagnostics-from-collector-evidence.md → internal/service/machine_diagnostics_service.go

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

## Communities (139 total, 25 thin omitted)

### Community 0 - "WithTransaction"
Cohesion: 0.08
Nodes (25): BreakFields, CreateBreakRequest, UpdateBreakRequest, clockMinutes(), CreateShiftRequest, ShiftFields, UpdateShiftRequest, MinutesBetween() (+17 more)

### Community 1 - "UserCityService"
Cohesion: 0.17
Nodes (6): CreateUserCityRequest, UpdateUserCityRequest, validateUserCity(), UserCity, mapUserCityError(), toUserCityResponseError()

### Community 2 - "UserLocationService"
Cohesion: 0.15
Nodes (8): CreateUserLocationRequest, UpdateUserLocationRequest, validateUserLocation(), UserLocation, UserLocationRepository, mapUserLocationError(), UserLocationService, toUserLocationResponseError()

### Community 3 - "Machine"
Cohesion: 0.13
Nodes (14): CreateMachineRequest, MachineFields, UpdateMachineRequest, validCoordinate(), validHealthIDs(), validIPv4(), Machine, MachineScope (+6 more)

### Community 4 - "net/http.Request"
Cohesion: 0.07
Nodes (20): BreakHandler, ClienteleHandler, axisValuesQuery(), dashboardQuery(), DashboardHandler, queryIDList(), readDashboardRequest(), DowntimeReasonHandler (+12 more)

### Community 5 - "newLiveTestServerWith"
Cohesion: 0.10
Nodes (42): NewBranchHandler(), NewRoleHandler(), NewSessionHandler(), NewUserCityHandler(), NewUserCountryHandler(), NewUserHandler(), NewUserLineHandler(), NewUserLocationHandler() (+34 more)

### Community 6 - "Execute"
Cohesion: 0.12
Nodes (3): Execute(), MachineDiagnosticsRepository, MachineHealthRepository

### Community 7 - "MachineModeReasonService"
Cohesion: 0.13
Nodes (14): CreateMachineModeReasonRequest, MachineModeReasonFields, MachineModeReasonsResponse, UpdateMachineModeReasonRequest, MachineModeReason, ReasonPeriod, MachineModeReasonRepository, mapMachineModeReasonError() (+6 more)

### Community 8 - "api_helpers_test.go"
Cohesion: 0.42
Nodes (8): cncMasterSnapshot, RevokeSessionsResponse, authUserKey, loginAsKey, sessionIDKey, liveValues, machineModelRecord, subscribeAnswer

### Community 9 - "CncMasterSnapshot (34,440 bytes, #pragma"
Cohesion: 0.10
Nodes (43): CncMasterSnapshot (34,440 bytes, #pragma pack(8), x86), AbsPos ([5]float64, Inches / Degrees, Work Coordinate System (WCS) positions), ActiveGCode ([16]byte, ASCII, Currently active G-code group), ActiveProgPath ([256]byte, ASCII Path, Storage path of active NC artifact), AlarmMsg ([128]byte, UTF-8 / ASCII, Alarm description message), AlarmNo (int32, Code (0 = OK), Primary alarm or fault number), AxisLoadXPercentage (int16, Percent (%), Servo X axis drive load (AxisLoadX/Y/ZPercentage)), AxisLoadYPercentage (int16, Percent (%), Servo Y axis drive load (AxisLoadX/Y/ZPercentage)) (+35 more)

### Community 10 - "LiveHandler"
Cohesion: 0.10
Nodes (14): Browser → server, liveSocket, liveSubscription, LiveSessionExpired, LiveSubscribeResponse, closeLiveSocket(), decodeLiveMessage(), LiveHandler (+6 more)

### Community 11 - "config.go"
Cohesion: 0.16
Nodes (16): main(), run(), Config, DBConfig, LiveConfig, MachineSelection, getEnv(), getEnvBool() (+8 more)

### Community 12 - "ProductionScheduleRepository"
Cohesion: 0.17
Nodes (4): MachineQuantity, ProductionScheduleItem, ProductionScheduleRepository, mapProductionScheduleError()

### Community 13 - "Collector"
Cohesion: 0.08
Nodes (32): 1. The big picture, 2. All steps at a glance, 4. Find the machines: steps 5 – 9, 6. When data is edited: step 18, 7. Shutdown: steps 19 – 20, 9. Data Reading Flow (Function Names Only), docs/cnc-communication.md, CNC Data Collection Flow (+24 more)

### Community 14 - "worker"
Cohesion: 0.09
Nodes (18): Known issue, Step 15 — Store the reading, Status, Value, LiveValue, TelemetryValue, engineName(), snapshotStatusOf() (+10 more)

### Community 15 - "Hub"
Cohesion: 0.08
Nodes (35): Decision, 1. Machine → DLL → server, 2. Worker → hub (before the database), 3. Hub → rooms, 4. Socket.IO → browser, Connecting the socket, Files, How the variable list works (+27 more)

### Community 16 - "run"
Cohesion: 0.05
Nodes (111): run(), Supervisor, Handler(), NewSupervisor(), Connect(), Querier, NewBrandHandler(), NewBreakHandler() (+103 more)

### Community 17 - "mode-chart.js"
Cohesion: 0.10
Nodes (64): card(), cell(), chartData(), chartRows(), countByGroup(), currentSpan(), currentWindow(), cycleData() (+56 more)

### Community 18 - ".GetDiagnostics"
Cohesion: 0.08
Nodes (42): DiagnosticPLCVariable, DiagnosticCheck, DiagnosticMachine, DiagnosticPLC, DiagnosticVariableRow, MachineDiagnosticsResponse, MachineNetworkResponse, DiagnosticEvent (+34 more)

### Community 19 - "net/http.Handler"
Cohesion: 0.17
Nodes (26): authRequest(), createAuthTestUser(), login(), newAuthRouter(), sessionCookieOf(), TestAuthentication(), dialLive(), TestLiveData() (+18 more)

### Community 20 - "net/http.ServeMux"
Cohesion: 0.12
Nodes (41): Decision, RequireRole(), registerAuthRoutes(), registerBranchRoutes(), registerBrandRoutes(), registerBreakRoutes(), registerCategoryRoutes(), registerClienteleRoutes() (+33 more)

### Community 21 - "UserSupervisorService"
Cohesion: 0.16
Nodes (9): CreateUserSupervisorRequest, UpdateUserSupervisorRequest, validateUserSupervisor(), NewUserSupervisorHandler(), UserSupervisor, UserSupervisorRepository, mapUserSupervisorError(), UserSupervisorService (+1 more)

### Community 22 - "SendResponse"
Cohesion: 0.05
Nodes (30): validator, BrandHandler, CategoryHandler, ControllerHandler, LicenceKeyHandler, MachineAccuracyHandler, NewMachineAccuracyHandler(), MachineCycleHandler (+22 more)

### Community 23 - "§10 Diagnostics & Troubleshooting"
Cohesion: 0.12
Nodes (25): Check mode: go run sinumerik_client.go check (OPC UA node existence report), Go Client (sinumerik_client.go, callback -> copy -> render), Prerequisite: 32-bit x86 DLL, compile Go with GOARCH=386, Prerequisite: SIEMENS_USER / SIEMENS_PASS env vars before CNCConnect, Run: go run sinumerik_client.go (live dashboard), §3.1 Environment & Prerequisites, §8 Running the Go Client, §10 Diagnostics & Troubleshooting (+17 more)

### Community 24 - "OperatorTaskService"
Cohesion: 0.09
Nodes (15): CreateMaintenanceTaskRequest, CreateOperatorTaskRequest, OperatorTaskFields, UpdateMaintenanceTaskRequest, UpdateOperatorTaskRequest, MaintenanceTask, OperatorTask, MaintenanceTaskRepository (+7 more)

### Community 25 - "machine_mode.go"
Cohesion: 0.15
Nodes (19): decideMachineMode(), holdReason(), machineStateFromReading(), numberingEngine(), programStateFromProcStatus(), ncReading(), TestDecideMachineMode(), TestDecideMachineModeConditions() (+11 more)

### Community 26 - "LocationRepository"
Cohesion: 0.50
Nodes (3): mapLocationError(), NewLocationRepository(), LocationRepository

### Community 27 - ".loadVariableConfig"
Cohesion: 0.22
Nodes (7): VariableConfig, formatPreviousValue(), formatValue(), newValueChangeLog(), boundVariable, valueChange, valueChangeLog

### Community 28 - "NewError"
Cohesion: 0.13
Nodes (10): CreateAddressRequest, UpdateAddressRequest, UpdateLanguageRequest, UpdateLicenceKeyRequest, validDuration(), ValidLicenceKey(), NewError(), newResponse() (+2 more)

### Community 29 - "companyOwner"
Cohesion: 0.08
Nodes (20): CreateDowntimeReasonRequest, DowntimeReasonFields, UpdateDowntimeReasonRequest, CreateIdleReasonRequest, IdleReasonFields, UpdateIdleReasonRequest, NewIdleReasonHandler(), DowntimeReason (+12 more)

### Community 30 - "ownCompanyCustomer"
Cohesion: 0.17
Nodes (13): Decision, PartProcessFields, PartSubProgramFields, CreatePartRequest, PartFields, UpdatePartRequest, ValidPartDrawingFile(), PartDrawing (+5 more)

### Community 31 - "Connection"
Cohesion: 0.09
Nodes (17): dllDriver, slot, 5. Connect and collect, per machine: steps 10 – 17, Step 11 — Connect, a few machines at a time, Step 12 — Call the DLL's CNCConnect, Step 13 — The DLL sends data, Step 17 — Disconnect and try again, Connection (+9 more)

### Community 32 - "Worker (internal/monitor/worker.go)"
Cohesion: 0.06
Nodes (29): CNC Data Collection (in-process), writeLog CNC Log Types and Switches, .env Configuration (internal/config), Health Mode History API (public), 7th Sense Server Claude Context, response.SendResponse Format, User CRUD Golden Reference, CncCommunication.dll (32-bit vendor DLL) (+21 more)

### Community 33 - "README.md"
Cohesion: 0.24
Nodes (5): 0015. Operator tasks: a shared task list, assigned per customer and machine model, Alternatives considered, Consequences, Context, Architecture Decision Records

### Community 34 - "newUserServices"
Cohesion: 0.15
Nodes (12): UserMachineFilter, NewAddressRepository(), execute(), UserMachineAccessRepository, NewUserMachineAccessRepository(), UserMachineFilterRepository, NewUserMachineFilterRepository(), NewAddressService() (+4 more)

### Community 35 - "MachineTelemetryService"
Cohesion: 0.09
Nodes (23): 3. Startup: steps 1 – 4, Step 1 — Start collection with the server, Step 2 — Load the DLL, Driver, LoadDriver(), CNCConfig, NewMachineDiagnosticsHandler(), NewCollector() (+15 more)

### Community 36 - "UserRepository"
Cohesion: 0.32
Nodes (3): UserListFilter, User, mapUserError()

### Community 37 - ".collectReadings"
Cohesion: 0.09
Nodes (22): 8. Reference, Events in `machine_events`, Related documents, Retry and timing summary, Settings (`.env`, defaults in `internal/config/config.go`), Step 10 — Get ready to connect, Step 14 — The data collection loop, Step 16 — Every second: dead connection check and live status (+14 more)

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
Cohesion: 0.06
Nodes (34): 0014. Session modes and the login conflict, on the existing auth_tokens sessions, Alternatives considered, Consequences, Context, Decision, SessionConfig, LoginRequest, describeDevice() (+26 more)

### Community 42 - "ownCompanyCustomer"
Cohesion: 0.16
Nodes (10): CreateHolidayRequest, HolidayFields, UpdateHolidayRequest, Holiday, HolidayRepository, mapHolidayError(), HolidayService, holidayFromFields() (+2 more)

### Community 43 - "CncCommunication.dll API & Integration G"
Cohesion: 0.13
Nodes (18): CNCConnect param machineType (int32 engine id, see §4.2), Engine 0: Simulator (local loopback mock engine), Engine 2: Siemens (Sinumerik OPC UA, active profile), Engine 3: Syntec (Open API / TCP), Engine 4: Mitsubishi (driver not specified), Machine identifier mapping (engine id enum 0-4), §9 Controller-Specific Profiles, Siemens default ports: 4840 (native NCU) / 4841 (custom) (+10 more)

### Community 44 - "ProductionScheduleService"
Cohesion: 0.18
Nodes (9): validDate(), CreateProductionScheduleRequest, ProductionScheduleFields, ProductionScheduleItemFields, UpdateProductionScheduleRequest, ScheduleDay(), checkJobsOverlap(), shiftSpan() (+1 more)

### Community 46 - "UserShopService"
Cohesion: 0.16
Nodes (8): CreateUserShopRequest, UpdateUserShopRequest, validateUserShop(), UserShop, UserShopRepository, mapUserShopError(), UserShopService, toUserShopResponseError()

### Community 47 - ".Validate"
Cohesion: 0.16
Nodes (16): UpdateThemeRequest, MachineFilter, UpdateUserRequest, normalizeTheme(), TestCreateUserRequestValidate(), TestCreateUserRequestValidateDropsDuplicateMachines(), TestCreateUserRequestValidateNormalizes(), TestCustomerFiltersValidate() (+8 more)

### Community 48 - "Machines API (/machines)"
Cohesion: 0.16
Nodes (18): Clienteles API, Company Scope (companyOwner / ownCompanyCustomer), Current User API (/me theme, language), Customer Display Names (customerName / companyName), Holidays API, Licence Keys API, Machine Catalog Brand-Controller-Model, Machines API (/machines) (+10 more)

### Community 49 - "DashboardRequest"
Cohesion: 0.21
Nodes (11): MachineTimeRow, MachineModeHistoryResponse, MachineModeTotal, MachineModePeriod, hours(), maxTime(), minTime(), modeGroup() (+3 more)

### Community 50 - "Address"
Cohesion: 0.27
Nodes (5): Address, AddressRepository, mapAddressError(), toAddressResponseError(), AddressService

### Community 51 - "validText"
Cohesion: 0.20
Nodes (4): UploadPartDrawingRequest, uniquePositiveIDs(), validText(), validWorkingDays()

### Community 52 - ".Create"
Cohesion: 0.23
Nodes (5): CreateBranchRequest, UpdateBranchRequest, validateBranchFields(), cityOwner(), toBranchResponseError()

### Community 53 - "README.md"
Cohesion: 0.40
Nodes (4): 0004. Opaque session token in an HttpOnly cookie, Alternatives considered, Consequences, Context

### Community 54 - "AuthHandler"
Cohesion: 0.14
Nodes (12): Decision, AuthHandler, CurrentUser(), isLoginAsKey(), LoginAsKey(), RouteLoginAs(), SessionToken(), UserID() (+4 more)

### Community 55 - "CompanyProfile"
Cohesion: 0.23
Nodes (8): CompanyProfile, CompanyProfileRepository, mapCompanyProfileError(), NewCompanyProfileRepository(), NewCompanyProfileService(), toCompanyProfileResponseError(), CompanyProfileService, TestCompanyProfileCRUD()

### Community 56 - "MachineCycle"
Cohesion: 0.15
Nodes (7): One running cycle per machine, CycleJob, MachineCycle, MachineCycleRepository, mapMachineCycleError(), sameID(), toMachineCycleResponseError()

### Community 57 - "MachineModelService"
Cohesion: 0.31
Nodes (3): MachineModel, MachineModelRepository, mapMachineModelError()

### Community 58 - "apidocs/main.go"
Cohesion: 0.16
Nodes (29): request, addExamples(), atoi(), authFolder(), blockStyle(), convertOperation(), dropNulls(), exampleOf() (+21 more)

### Community 59 - "generate.go"
Cohesion: 0.18
Nodes (6): EnsureCertificates(), existingCertificateExpiry(), generateCertificates(), Paths, writePEM(), Settings

### Community 60 - "UserProfileRepository"
Cohesion: 0.14
Nodes (10): UserProfile, UserProfileRepository, mapProfileError(), NewUserProfileRepository(), NewUserProfileService(), toUserProfileResponseError(), UserProfileService, sameClock() (+2 more)

### Community 61 - "Test list"
Cohesion: 0.05
Nodes (38): 1. The rows are consistent (SQL, no machine needed), 2. The rows agree with the mode log (SQL), 3. The raw signals around one cycle (telemetry and logs), 4. At the machine (floor test), Before the tests (once), Check query, Completed, Floor tests on the HUMA panel (Fanuc) (+30 more)

### Community 62 - "0001. Go standard library net/http, no web framework"
Cohesion: 0.33
Nodes (5): 0001. Go standard library net/http, no web framework, Alternatives considered, Consequences, Context, Decision

### Community 63 - "0003. PostgreSQL with pgx and hand-written SQL"
Cohesion: 0.33
Nodes (5): 0003. PostgreSQL with pgx and hand-written SQL, Alternatives considered, Consequences, Context, Decision

### Community 64 - "MachinePlcVariableService"
Cohesion: 0.23
Nodes (9): MachinePlcVariableResponse, MachinePlcVariablesResponse, UpdatePlcNodeIDRequest, NewMachinePlcVariableResponse(), NewMachinePlcVariablesResponse(), MachinePlcVariable, plcDefaultOrder(), MachinePlcVariableService (+1 more)

### Community 65 - "time.Duration"
Cohesion: 0.09
Nodes (16): PLCConfig, 0021. Read PLC variables over OPC UA with gopcua, in-process beside the CNC collector, Alternatives considered, Consequences, Context, Decision, ReadsPLCOverOPCUA(), newPlcWorker() (+8 more)

### Community 66 - "UserLineService"
Cohesion: 0.07
Nodes (19): CreateUserCellRequest, UpdateUserCellRequest, validateUserCell(), CreateUserLineRequest, UpdateUserLineRequest, validateUserLine(), NewUserCellHandler(), UserCell (+11 more)

### Community 67 - "run"
Cohesion: 0.36
Nodes (10): checkNotTaken(), createSuperadmin(), ensureProfile(), execute(), main(), parseFlags(), run(), superadminRoleID() (+2 more)

### Community 68 - "MachineAccuracyService"
Cohesion: 0.19
Nodes (6): MachineAccuracyRequest, MachineAccuracy, MachineAccuracyRepository, mapMachineAccuracyError(), MachineAccuracyService, toMachineAccuracyResponseError()

### Community 69 - "UserHandler"
Cohesion: 0.18
Nodes (3): BranchHandler, currentUserID(), UserHandler

### Community 70 - "ControllerRepository"
Cohesion: 0.10
Nodes (10): Decision, CreateControllerRequest, UpdateControllerRequest, validateControllerFields(), MachineModelRequest, Controller, mapControllerError(), toControllerResponseError() (+2 more)

### Community 71 - "onDataCallback(ptr uintptr) uintptr"
Cohesion: 0.36
Nodes (9): Package-level callback singletons (syscall.NewCallback), CNCConnect param errCb (uintptr stdcall OnErr(int32 code, const char* msg)), Mutex-protected copy: dataMutex sync.RWMutex -> lastRecord/hasRecord, onDataCallback(ptr uintptr) uintptr, onErrorCallback(code, messagePtr uintptr) uintptr, §7 Asynchronous Callback Architecture, onDataCallback, onErrorCallback (+1 more)

### Community 72 - "time.Time"
Cohesion: 0.05
Nodes (42): AxisPoint, FinishCyclesSummary, FinishedCycle, TimelinePeriod, VirtualScreenCycle, VirtualScreenLive, FinishCyclesResponse, ModeTimelineResponse (+34 more)

### Community 73 - "HomeDashboardResponse"
Cohesion: 0.10
Nodes (21): HomeAlarm, HomeFleet, HomeMachine, HomePlan, HomeTime, HomeDashboardResponse, HomeHour, HomeModeCount (+13 more)

### Community 74 - "Brand"
Cohesion: 0.06
Nodes (22): CreateBrandRequest, UpdateBrandRequest, validateCatalogFields(), CreateCategoryRequest, UpdateCategoryRequest, CreateProductRequest, UpdateProductRequest, validateProductFields() (+14 more)

### Community 75 - "testing.T"
Cohesion: 0.07
Nodes (102): TestSnapshotModeFieldOffsets(), TestSnapshotSize(), CreateUserRequest, TestParseDayStart(), containsID(), decodeData(), expectStatus(), newCatalogRouter() (+94 more)

### Community 76 - "0008. "Log in as customer" uses a separate session with a path-scoped cookie"
Cohesion: 0.40
Nodes (4): 0008. "Log in as customer" uses a separate session with a path-scoped cookie, Alternatives considered, Consequences, Context

### Community 77 - "WithTransaction"
Cohesion: 0.32
Nodes (7): WithTransaction(), ProductionSchedule, ProductionScheduleService, saveScheduleItems(), scheduleDetails(), scheduleFromFields(), toProductionScheduleResponseError()

### Community 78 - "UserCountryRepository"
Cohesion: 0.18
Nodes (6): CreateUserCountryRequest, UserCountry, UserCountryRepository, mapUserCountryError(), UserCountryService, toUserCountryResponseError()

### Community 79 - "AuthUser"
Cohesion: 0.23
Nodes (4): LoginAsResponse, MachineCyclesResponse, AuthUser, machineScope()

### Community 80 - "Event-driven data flow (no polling in client)"
Cohesion: 0.60
Nodes (5): Flow 1: load DLL, register two stdcall callbacks (syscall.NewCallback), Flow 2: CNCConnect starts DLL background worker thread, Flow 3: worker polls controller, assembles snapshot, invokes data callback, Flow 4: callback copies payload, signals renderer via buffered channel, Flow 5: main goroutine renders live ANSI dashboard

### Community 81 - "DashboardService"
Cohesion: 0.22
Nodes (13): CycleCountRow, DashboardRequest, clipWindows(), cyclesIn(), dayWindows(), DashboardService, machineDayOEE(), round1() (+5 more)

### Community 82 - "Live data — how a machine value reaches the screen"
Cohesion: 0.14
Nodes (11): 0009. Socket.IO live data: the collector pushes each screen's variables, checked per user, Alternatives considered, Consequences, Context, Adding a variable to a screen, Connecting, Live data — Socket.IO event contract, Ordering (+3 more)

### Community 83 - "PartRepository"
Cohesion: 0.14
Nodes (7): Part, PartProcess, PartSubProgram, executeAffected(), PartRepository, mapPartError(), saveProcesses()

### Community 84 - "CNCConnect(ip, port, dataCb, errCb, machineType) int32"
Cohesion: 0.20
Nodes (11): CNCConnect(ip, port, dataCb, errCb, machineType) int32, CNCConnect param dataCb (uintptr stdcall OnData(const CncMasterSnapshot*)), CNCConnect param ip (const char*, ANSI NUL-terminated), CNCConnect param port (int32, e.g. 4840), CNCdisconnect(machineId int32), LazyDLL binding (syscall.NewLazyDLL + eager mod.Load in init), §4 API Reference, §5 DLL Loading & Binding (Go) (+3 more)

### Community 85 - "context.Context"
Cohesion: 0.09
Nodes (6): MachineEvent, MachineLiveStatus, ModeChange, MachineTelemetryRepository, mapMachineTelemetryError(), levelInMachineFilter()

### Community 86 - "auth_tokens SHA-256 Hashed Sessions (mig"
Cohesion: 0.67
Nodes (3): auth_tokens SHA-256 Hashed Sessions (migration 014), Per-Request Session and User Re-check, ts_session Cookie (HttpOnly, SameSite=Strict)

### Community 89 - "HealthHandler"
Cohesion: 0.16
Nodes (3): HealthHandler, registerHealthRoutes(), statusRecorder

### Community 104 - "connection.go"
Cohesion: 0.12
Nodes (11): changeOf(), Connect(), forwardNotifications(), Change, Item, ItemResult, isBad(), securityOf() (+3 more)

### Community 105 - "0005. CNC collector runs inside the API process"
Cohesion: 0.33
Nodes (5): 0005. CNC collector runs inside the API process, Alternatives considered, Consequences, Context, Decision

### Community 107 - "0020. Machine PLC node IDs: the machine's own node ID in machine_variables, default otherwise"
Cohesion: 0.33
Nodes (5): 0020. Machine PLC node IDs: the machine's own node ID in machine_variables, default otherwise, Alternatives considered, Consequences, Context, Decision

### Community 108 - ".Update"
Cohesion: 0.20
Nodes (12): Decision, CustomerDetails, UserResponse, isTimeZone(), isWebsite(), normalizeFilters(), ParseDayStart(), companyProfileFrom() (+4 more)

### Community 109 - "0011. Ponytail plugin for Claude Code, under this repo's rules"
Cohesion: 0.33
Nodes (5): 0011. Ponytail plugin for Claude Code, under this repo's rules, Alternatives considered, Consequences, Context, Decision

### Community 110 - "Client Certificate & Key (client_cert.der / client_key.der)"
Cohesion: 0.48
Nodes (7): client_cert.der/.pem (X.509 client cert for Siemens OPC UA), client_key.der/.pem (client private key), Siemens certs: client_cert.der + client_key.der in binary directory, Siemens trust: cert must be trusted in NCU Access Administration, Troubleshoot: BadRequestNotAllowed -> Export client_cert, BadRequestNotAllowed (certificate trust rejection), Client Certificate & Key (client_cert.der / client_key.der)

### Community 111 - "MachineCycleSummary"
Cohesion: 0.18
Nodes (10): Decision, HomeCycles, MachineSummary, NewFinishCyclesSummary(), NewFinishedCycle(), NewHomeCycles(), NewMachineSummary(), NewTimelinePeriod() (+2 more)

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
Cohesion: 0.40
Nodes (5): 0018. Machine Health X / Y / Z charts: downsampled in SQL, timed, on the machine-health routes, Alternatives considered, Consequences, Context, Decision

### Community 119 - "0017. Machine accuracy as one shared list, managed by the superadmin"
Cohesion: 0.40
Nodes (4): 0017. Machine accuracy as one shared list, managed by the superadmin, Alternatives considered, Consequences, Context

### Community 120 - "BranchRepository"
Cohesion: 0.31
Nodes (3): Branch, BranchRepository, mapBranchError()

### Community 121 - "0006. Page-wise dashboard endpoints returning only chart fields"
Cohesion: 0.40
Nodes (4): 0006. Page-wise dashboard endpoints returning only chart fields, Alternatives considered, Consequences, Context

### Community 122 - "run"
Cohesion: 0.36
Nodes (4): connectFanucCounter(), main(), run(), fanucCounter

### Community 123 - "connection_test.go"
Cohesion: 0.32
Nodes (3): TestValueOf(), ValueOf(), Value

### Community 124 - ".validate"
Cohesion: 0.38
Nodes (4): ClienteleFields, CreateClienteleRequest, UpdateClienteleRequest, validClienteleEmail()

### Community 125 - "UserService"
Cohesion: 0.17
Nodes (7): Decision, AssignableRoles(), CanAssignRole(), SeesAllCompanyMachines(), canSee(), UserService, toUserResponseError()

### Community 126 - ".Create"
Cohesion: 0.33
Nodes (3): CreateUserStateRequest, validateUserState(), toUserStateResponseError()

### Community 128 - "logger/logger.go"
Cohesion: 0.40
Nodes (3): main(), Init(), InitFile()

### Community 129 - "0016. A Task List screen for the operator tasks' master list"
Cohesion: 0.40
Nodes (5): 0016. A Task List screen for the operator tasks' master list, Alternatives considered, Consequences, Context, Decision

### Community 130 - "LicenceKeyService"
Cohesion: 0.13
Nodes (11): NewLicenceKey, CreateLicenceKeysRequest, GenerateLicenceKeysRequest, GeneratedLicenceKeys, LicenceKey, LicenceKeyRepository, mapLicenceKeyError(), LicenceKeyService (+3 more)

### Community 131 - "ClienteleService"
Cohesion: 0.44
Nodes (4): Clientele, clienteleFromFields(), ClienteleService, toClienteleResponseError()

### Community 132 - "0012. Downtime and idle reasons on a period: one table, tied to machine_mode_log"
Cohesion: 0.33
Nodes (5): 0012. Downtime and idle reasons on a period: one table, tied to machine_mode_log, Alternatives considered, Consequences, Context, Decision

### Community 136 - "0013. Machine diagnostics: reasons from what the collector recorded, network tested from the server"
Cohesion: 0.40
Nodes (4): 0013. Machine diagnostics: reasons from what the collector recorded, network tested from the server, Alternatives considered, Consequences, Context

### Community 139 - "0010. Machine access by an organization filter, worked out on every request"
Cohesion: 0.40
Nodes (4): 0010. Machine access by an organization filter, worked out on every request, Alternatives considered, Consequences, Context

## Ambiguous Edges - Review These
- `CncCommunication.dll (32-bit vendor DLL)` → `FanucCncCommunication.dll`  [AMBIGUOUS]
  docs/setup-guide.txt · relation: semantically_similar_to
- `Snapshot Fields (positions, spindle, feed, program, tool, alarms)` → `Machines API (/machines)`  [AMBIGUOUS]
  docs/CncCommunication.dll — API & Integration Guide (Go Client).pdf · relation: conceptually_related_to

## Knowledge Gaps
- **33 isolated node(s):** `7th-sense-server`, `middleware.RequireRole / routes.adminOnly`, `response.SendResponse`, `Test Auth (asTestUser, tests/auth_test.go)`, `user_filters Filter Checklist` (+28 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 244 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **25 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `CNC Data Collection Flow` connect `Collector` to `Worker (internal/monitor/worker.go)`, `MachineTelemetryService`, `.collectReadings`, `Live data — how a machine value reaches the screen`, `Connection`?**
  _High betweenness centrality (0.123) - this node is a cross-community bridge._
- **What connects `7th-sense-server`, `middleware.RequireRole / routes.adminOnly`, `response.SendResponse` to the rest of the system?**
  _33 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `WithTransaction` be split into smaller, more focused modules?**
  _Cohesion score 0.07676767676767676 - nodes in this community are weakly interconnected._
- **Why does `2. All steps at a glance` connect `Collector` to `MachineTelemetryService`, `.collectReadings`, `worker`, `Connection`?**
  _High betweenness centrality (0.101) - this node is a cross-community bridge._
- **Should `Machine` be split into smaller, more focused modules?**
  _Cohesion score 0.12561576354679804 - nodes in this community are weakly interconnected._

### Low-confidence Hints
_AMBIGUOUS edges — the extractor was unsure. Verify before acting on these._

- **What is the exact relationship between `CncCommunication.dll (32-bit vendor DLL)` and `FanucCncCommunication.dll`?**
  _Edge tagged AMBIGUOUS (relation: semantically_similar_to) - confidence is low._
- **What is the exact relationship between `Snapshot Fields (positions, spindle, feed, program, tool, alarms)` and `Machines API (/machines)`?**
  _Edge tagged AMBIGUOUS (relation: conceptually_related_to) - confidence is low._