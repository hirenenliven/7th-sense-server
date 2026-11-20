# Graph Report - 7th-sense-server  (2026-11-19)

## Corpus Check
- 372 files · ~309,202 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 4 file(s) not represented in the graph (top: (none) 3, .toml 1)

## Summary
- 3247 nodes · 30992 edges · 126 communities (107 shown, 19 thin omitted)
- Extraction: 96% EXTRACTED · 4% INFERRED · 0% AMBIGUOUS · INFERRED: 1342 edges (avg confidence: 0.87)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `188077aa`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- WithTransaction
- UserCityService
- UserLocationService
- MachineService
- net/http.Request
- newLiveTestServerWith
- Querier
- MachineModeReasonService
- api_helpers_test.go
- CncMasterSnapshot (34,440 bytes, #pragma
- LiveHandler
- cycleReading
- ProductionScheduleRepository
- Collector
- worker
- Hub
- run
- mode-chart.js
- .GetDiagnostics
- net/http.Handler
- net/http.ServeMux
- TestMachineFilterMatrix
- net/http.Request
- §10 Diagnostics & Troubleshooting
- OperatorTaskService
- machine_mode.go
- UserLineService
- Product
- actAs
- AuthUser
- ownCompanyCustomer
- dllDriver
- Worker (internal/monitor/worker.go)
- README.md
- UserMachineFilter
- net/http.Handler
- decodeData
- Category
- liveTestClient
- .GetAxisValues
- CncCommunication.dll API & Integration G
- AuthService
- .validate
- CncCommunication.dll API & Integration G
- AuthService
- BreakRepository
- UserShopService
- .Validate
- Machines API (/machines)
- MaintenanceTaskRepository
- Address
- validText
- .Create
- README.md
- newUserServices
- CompanyProfile
- MachineCycle
- TestGujaratMachineAccessScenario
- config.go
- .validate
- UserProfileRepository
- Test list
- 0001. Go standard library net/http, no web framework
- 0003. PostgreSQL with pgx and hand-written SQL
- .validate
- time.Duration
- UserLineService
- Execute
- MachineAccuracyService
- .validate
- 0014. Session modes and the login conflict, on the existing auth_tokens sessions
- onDataCallback(ptr uintptr) uintptr
- time.Time
- MachineTelemetryService
- ControllerService
- testing.T
- 0008. "Log in as customer" uses a separate session with a path-scoped cookie
- 0022. Test PLC node IDs in a short OPC UA session of their own
- UserCountryRepository
- MachineModeStatus
- Event-driven data flow (no polling in client)
- TestOrganizationForAnotherUser
- Live data — how a machine value reaches the screen
- PartRepository
- CNCConnect(ip, port, dataCb, errCb, machineType) int32
- MachineTelemetryRepository
- auth_tokens SHA-256 Hashed Sessions (mig
- context.Context
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
- UserService
- 0011. Ponytail plugin for Claude Code, under this repo's rules
- Client Certificate & Key (client_cert.der / client_key.der)
- 0015. Operator tasks: a shared task list, assigned per customer and machine model
- user_access_test.go
- .validate
- 0007. OpenAPI 3 from swag annotations, Swagger UI behind API_DOCS, generated Postman
- 0002. Layered architecture: handler, DTO, service, repository
- 0018. Machine Health X / Y / Z charts: downsampled in SQL, timed, on the machine-health routes
- SessionHandler
- 0017. Machine accuracy as one shared list, managed by the superadmin
- 0006. Page-wise dashboard endpoints returning only chart fields
- .validate
- UserStateService
- SessionRevokedError
- LicenceKeyService
- pathID
- 0013. Machine diagnostics: reasons from what the collector recorded, network tested from the server

## God Nodes (most connected - your core abstractions)
1. `Execute()` - 225 edges
2. `SendResponse()` - 191 edges
3. `NewError()` - 159 edges
4. `run()` - 119 edges
5. `AuthUser` - 113 edges
6. `pathID()` - 99 edges
7. `currentUser()` - 85 edges
8. `Querier` - 79 edges
9. `decodeAndValidate()` - 69 edges
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

## Communities (126 total, 19 thin omitted)

### Community 0 - "WithTransaction"
Cohesion: 0.14
Nodes (21): WithTransaction(), clockMinutes(), CreateShiftRequest, ShiftFields, UpdateShiftRequest, MinutesBetween(), NewShiftHandler(), Shift (+13 more)

### Community 1 - "UserCityService"
Cohesion: 0.16
Nodes (8): CreateUserCityRequest, UpdateUserCityRequest, validateUserCity(), UserCity, UserCityRepository, mapUserCityError(), UserCityService, toUserCityResponseError()

### Community 2 - "UserLocationService"
Cohesion: 0.16
Nodes (8): CreateUserLocationRequest, UpdateUserLocationRequest, validateUserLocation(), UserLocation, UserLocationRepository, mapUserLocationError(), UserLocationService, toUserLocationResponseError()

### Community 3 - "MachineService"
Cohesion: 0.06
Nodes (36): 0010. Machine access by an organization filter, worked out on every request, Alternatives considered, Consequences, Context, Decision, PlcReaderResponse, CreateMachineRequest, MachineFields (+28 more)

### Community 4 - "net/http.Request"
Cohesion: 0.07
Nodes (15): BreakHandler, ClienteleHandler, DowntimeReasonHandler, HolidayHandler, IdleReasonHandler, MachineHandler, MachineModeReasonHandler, PartHandler (+7 more)

### Community 5 - "newLiveTestServerWith"
Cohesion: 0.14
Nodes (37): NewBranchHandler(), RoleHandler, NewRoleHandler(), NewSessionHandler(), NewUserCityHandler(), NewUserCountryHandler(), NewUserHandler(), NewUserLineHandler() (+29 more)

### Community 6 - "Querier"
Cohesion: 0.16
Nodes (16): Querier, NewDowntimeReasonRepository(), NewHolidayRepository(), NewIdleReasonRepository(), NewLicenceKeyRepository(), NewMachineCycleRepository(), MachineHealthRepository, NewMachineHealthRepository() (+8 more)

### Community 7 - "MachineModeReasonService"
Cohesion: 0.10
Nodes (20): 0012. Downtime and idle reasons on a period: one table, tied to machine_mode_log, Alternatives considered, Consequences, Context, Decision, CreateMachineModeReasonRequest, MachineModeReasonFields, MachineModeReasonsResponse (+12 more)

### Community 8 - "api_helpers_test.go"
Cohesion: 0.42
Nodes (10): cncMasterSnapshot, LanguageResponse, RevokeSessionsResponse, ThemeResponse, authUserKey, loginAsKey, sessionIDKey, liveValues (+2 more)

### Community 9 - "CncMasterSnapshot (34,440 bytes, #pragma"
Cohesion: 0.10
Nodes (43): CncMasterSnapshot (34,440 bytes, #pragma pack(8), x86), AbsPos ([5]float64, Inches / Degrees, Work Coordinate System (WCS) positions), ActiveGCode ([16]byte, ASCII, Currently active G-code group), ActiveProgPath ([256]byte, ASCII Path, Storage path of active NC artifact), AlarmMsg ([128]byte, UTF-8 / ASCII, Alarm description message), AlarmNo (int32, Code (0 = OK), Primary alarm or fault number), AxisLoadXPercentage (int16, Percent (%), Servo X axis drive load (AxisLoadX/Y/ZPercentage)), AxisLoadYPercentage (int16, Percent (%), Servo Y axis drive load (AxisLoadX/Y/ZPercentage)) (+35 more)

### Community 10 - "LiveHandler"
Cohesion: 0.07
Nodes (22): Browser → server, liveSocket, liveSubscription, LiveSessionExpired, AuthHandler, channelRoom(), closeLiveSocket(), decodeLiveMessage() (+14 more)

### Community 11 - "cycleReading"
Cohesion: 0.15
Nodes (8): CycleStatusFor(), cycleReadingFrom(), decideCycleStep(), isProgramLoop(), resetReason(), cycleReading, cycleStep, cycleTracker

### Community 12 - "ProductionScheduleRepository"
Cohesion: 0.08
Nodes (24): validDate(), CreateProductionScheduleRequest, OperatorJobFilter, ProductionScheduleFields, ProductionScheduleItemFields, UpdateProductionScheduleRequest, ScheduleDay(), NewProductionScheduleHandler() (+16 more)

### Community 13 - "Collector"
Cohesion: 0.07
Nodes (33): 1. The big picture, 4. Find the machines: steps 5 – 9, 6. When data is edited: step 18, 7. Shutdown: steps 19 – 20, 8. Reference, 9. Data Reading Flow (Function Names Only), docs/cnc-communication.md, CNC Data Collection Flow (+25 more)

### Community 14 - "worker"
Cohesion: 0.06
Nodes (38): 2. All steps at a glance, 5. Connect and collect, per machine: steps 10 – 17, Known issue, Retry and timing summary, Step 10 — Get ready to connect, Step 11 — Connect, a few machines at a time, Step 12 — Call the DLL's CNCConnect, Step 13 — The DLL sends data (+30 more)

### Community 15 - "Hub"
Cohesion: 0.06
Nodes (46): Supervisor, 2. Worker → hub (before the database), 3. Hub → rooms, 4. Socket.IO → browser, Connecting the socket, Files, How the variable list works, Live data — how a machine value reaches the screen (+38 more)

### Community 16 - "run"
Cohesion: 0.12
Nodes (35): run(), Connect(), NewBrandHandler(), NewCategoryHandler(), NewClienteleHandler(), NewControllerHandler(), NewMachineModelHandler(), NewMaintenanceTaskHandler() (+27 more)

### Community 17 - "mode-chart.js"
Cohesion: 0.10
Nodes (64): card(), cell(), chartData(), chartRows(), countByGroup(), currentSpan(), currentWindow(), cycleData() (+56 more)

### Community 18 - ".GetDiagnostics"
Cohesion: 0.08
Nodes (41): DiagnosticPLCVariable, DiagnosticCheck, DiagnosticMachine, DiagnosticPLC, DiagnosticVariableRow, MachineDiagnosticsResponse, MachineNetworkResponse, DiagnosticEvent (+33 more)

### Community 19 - "net/http.Handler"
Cohesion: 0.28
Nodes (26): expectStatus(), uniqueSuffix(), authRequest(), createAuthTestUser(), login(), newAuthRouter(), sessionCookieOf(), TestAuthentication() (+18 more)

### Community 20 - "net/http.ServeMux"
Cohesion: 0.12
Nodes (41): Decision, RequireRole(), registerAuthRoutes(), registerBranchRoutes(), registerBrandRoutes(), registerBreakRoutes(), registerCategoryRoutes(), registerClienteleRoutes() (+33 more)

### Community 21 - "TestMachineFilterMatrix"
Cohesion: 0.17
Nodes (27): CreateUserRequest, filterTree, levelCase, insertModePeriod(), insertReason(), TestMachineModeReasons(), insertTestRow(), machineBody() (+19 more)

### Community 22 - "net/http.Request"
Cohesion: 0.06
Nodes (23): validator, BranchHandler, BrandHandler, CategoryHandler, ControllerHandler, axisValuesQuery(), dashboardQuery(), DashboardHandler (+15 more)

### Community 23 - "§10 Diagnostics & Troubleshooting"
Cohesion: 0.12
Nodes (25): Check mode: go run sinumerik_client.go check (OPC UA node existence report), Go Client (sinumerik_client.go, callback -> copy -> render), Prerequisite: 32-bit x86 DLL, compile Go with GOARCH=386, Prerequisite: SIEMENS_USER / SIEMENS_PASS env vars before CNCConnect, Run: go run sinumerik_client.go (live dashboard), §3.1 Environment & Prerequisites, §8 Running the Go Client, §10 Diagnostics & Troubleshooting (+17 more)

### Community 24 - "OperatorTaskService"
Cohesion: 0.08
Nodes (20): MachineModelRequest, CreateMaintenanceTaskRequest, CreateOperatorTaskRequest, OperatorTaskFields, UpdateMaintenanceTaskRequest, UpdateOperatorTaskRequest, MachineModel, OperatorTask (+12 more)

### Community 25 - "machine_mode.go"
Cohesion: 0.15
Nodes (20): decideMachineMode(), holdReason(), machineStateFromReading(), numberingEngine(), programStateFromProcStatus(), ncReading(), TestDecideMachineMode(), TestDecideMachineModeConditions() (+12 more)

### Community 26 - "UserLineService"
Cohesion: 0.15
Nodes (8): CreateUserLineRequest, UpdateUserLineRequest, validateUserLine(), UserLine, UserLineRepository, mapUserLineError(), UserLineService, toUserLineResponseError()

### Community 27 - "Product"
Cohesion: 0.16
Nodes (9): CreateProductRequest, UpdateProductRequest, validateProductFields(), Product, ProductRepository, mapProductError(), ProductService, productImageOrDefault() (+1 more)

### Community 28 - "actAs"
Cohesion: 0.23
Nodes (19): clienteleBody(), TestClienteleCRUD(), actAs(), containsMachineID(), newMachineFixture(), TestMachineCRUD(), with(), machineFixture (+11 more)

### Community 29 - "AuthUser"
Cohesion: 0.06
Nodes (38): LoginAsResponse, NewBreakHandler(), NewDowntimeReasonHandler(), NewHolidayHandler(), NewIdleReasonHandler(), AuthUser, Break, Clientele (+30 more)

### Community 30 - "ownCompanyCustomer"
Cohesion: 0.19
Nodes (8): ValidPartDrawingFile(), Part, PartDrawingStorage, NewPartRepository(), PartService, partFromFields(), saveProcesses(), toPartResponseError()

### Community 31 - "dllDriver"
Cohesion: 0.08
Nodes (19): connectFanucCounter(), main(), run(), dllDriver, slot, 3. Startup: steps 1 – 4, Step 1 — Start collection with the server, Step 2 — Load the DLL (+11 more)

### Community 32 - "Worker (internal/monitor/worker.go)"
Cohesion: 0.06
Nodes (29): CNC Data Collection (in-process), writeLog CNC Log Types and Switches, .env Configuration (internal/config), Health Mode History API (public), 7th Sense Server Claude Context, response.SendResponse Format, User CRUD Golden Reference, CncCommunication.dll (32-bit vendor DLL) (+21 more)

### Community 33 - "README.md"
Cohesion: 0.22
Nodes (6): 0016. A Task List screen for the operator tasks' master list, Alternatives considered, Consequences, Context, Decision, Architecture Decision Records

### Community 34 - "UserMachineFilter"
Cohesion: 0.24
Nodes (6): UserMachineFilter, UserMachineFilterRepository, NewUserMachineFilterRepository(), userAccess, userServices, userAccess

### Community 35 - "net/http.Handler"
Cohesion: 0.13
Nodes (26): Handler(), NewMachineDiagnosticsHandler(), NewMachinePlcVariableHandler(), CrossOrigin(), Logger(), Chain(), Middleware, Recoverer() (+18 more)

### Community 36 - "decodeData"
Cohesion: 0.25
Nodes (16): containsID(), decodeData(), newCatalogRouter(), runCatalogCRUD(), createTestRecord(), TestBrandCRUD(), TestCategoryCRUD(), TestControllerBrand() (+8 more)

### Community 37 - "Category"
Cohesion: 0.25
Nodes (5): Category, CategoryRepository, mapCategoryError(), CategoryService, toCategoryResponseError()

### Community 38 - "liveTestClient"
Cohesion: 0.26
Nodes (5): dialLive(), newLiveTestServer(), liveTestClient, liveTestEvent, liveTestServer

### Community 39 - ".GetAxisValues"
Cohesion: 0.13
Nodes (15): 0019. Axis values: every stored value when it fits, and zoom by loading the window again, Alternatives considered, Consequences, Context, Decision, AxisPoint, MachineCycleListRequest, FindAxisSet() (+7 more)

### Community 40 - "CncCommunication.dll API & Integration G"
Cohesion: 0.19
Nodes (20): CncCommunication.dll (C++ core runtime, x86, OPC UA inside), CncCommunication.dll API & Integration Guide (Go Client) - text version, Fanuc profile (engine 1, port 8193, Fwlib32.dll + fwlibe1.dll), Fwlib32.dll (Fanuc FOCAS2), fwlibe1.dll (Fanuc FOCAS2), libcrypto-3.dll (OpenSSL), libssl-3.dll (OpenSSL), Engine 1: Fanuc (FOCAS2 Ethernet, fwlib32) (+12 more)

### Community 41 - "AuthService"
Cohesion: 0.06
Nodes (22): Decision, LoginRequest, describeDevice(), SessionClient, SessionConflictResponse, SessionResponse, NewSessionResponse(), NewSessionResponses() (+14 more)

### Community 42 - ".validate"
Cohesion: 0.47
Nodes (3): CreateHolidayRequest, HolidayFields, UpdateHolidayRequest

### Community 43 - "CncCommunication.dll API & Integration G"
Cohesion: 0.13
Nodes (18): CNCConnect param machineType (int32 engine id, see §4.2), Engine 0: Simulator (local loopback mock engine), Engine 2: Siemens (Sinumerik OPC UA, active profile), Engine 3: Syntec (Open API / TCP), Engine 4: Mitsubishi (driver not specified), Machine identifier mapping (engine id enum 0-4), §9 Controller-Specific Profiles, Siemens default ports: 4840 (native NCU) / 4841 (custom) (+10 more)

### Community 44 - "AuthService"
Cohesion: 0.23
Nodes (9): SessionConfig, NewAuthHandler(), NewLiveHandler(), Authenticate(), NewAuthRepository(), NewSessionRepository(), AuthService, NewAuthService() (+1 more)

### Community 45 - "BreakRepository"
Cohesion: 0.20
Nodes (3): ShiftBreak, BreakRepository, mapBreakError()

### Community 46 - "UserShopService"
Cohesion: 0.16
Nodes (8): CreateUserShopRequest, UpdateUserShopRequest, validateUserShop(), UserShop, UserShopRepository, mapUserShopError(), UserShopService, toUserShopResponseError()

### Community 47 - ".Validate"
Cohesion: 0.11
Nodes (20): CreateAddressRequest, UpdateAddressRequest, UpdateLanguageRequest, UpdateThemeRequest, normalizeTheme(), TestCreateUserRequestValidate(), TestCreateUserRequestValidateDropsDuplicateMachines(), TestCreateUserRequestValidateNormalizes() (+12 more)

### Community 48 - "Machines API (/machines)"
Cohesion: 0.16
Nodes (18): Clienteles API, Company Scope (companyOwner / ownCompanyCustomer), Current User API (/me theme, language), Customer Display Names (customerName / companyName), Holidays API, Licence Keys API, Machine Catalog Brand-Controller-Model, Machines API (/machines) (+10 more)

### Community 49 - "MaintenanceTaskRepository"
Cohesion: 0.33
Nodes (3): MaintenanceTask, MaintenanceTaskRepository, mapMaintenanceTaskError()

### Community 50 - "Address"
Cohesion: 0.24
Nodes (6): Address, AddressRepository, mapAddressError(), toAddressResponseError(), primaryAddress(), AddressService

### Community 51 - "validText"
Cohesion: 0.13
Nodes (10): PartProcessFields, PartSubProgramFields, CreatePartRequest, PartFields, UpdatePartRequest, UploadPartDrawingRequest, uniquePositiveIDs(), validText() (+2 more)

### Community 52 - ".Create"
Cohesion: 0.25
Nodes (5): CreateBranchRequest, UpdateBranchRequest, validateBranchFields(), cityOwner(), toBranchResponseError()

### Community 53 - "README.md"
Cohesion: 0.40
Nodes (4): 0004. Opaque session token in an HttpOnly cookie, Alternatives considered, Consequences, Context

### Community 54 - "newUserServices"
Cohesion: 0.28
Nodes (9): NewAddressRepository(), NewUserProfileRepository(), NewAddressService(), NewUserProfileService(), newUserServices(), TestAddressCRUD(), sameClock(), sameDate() (+1 more)

### Community 55 - "CompanyProfile"
Cohesion: 0.23
Nodes (8): CompanyProfile, CompanyProfileRepository, mapCompanyProfileError(), NewCompanyProfileRepository(), NewCompanyProfileService(), toCompanyProfileResponseError(), CompanyProfileService, TestCompanyProfileCRUD()

### Community 56 - "MachineCycle"
Cohesion: 0.12
Nodes (11): One running cycle per machine, MachineCyclesResponse, NewMachineCycleHandler(), CycleJob, MachineCycle, JobProgress, MachineCycleRepository, mapMachineCycleError() (+3 more)

### Community 57 - "TestGujaratMachineAccessScenario"
Cohesion: 0.33
Nodes (7): apiResponse, gujaratSeed, loadGujaratSeed(), logResponse(), machineNames(), softDeleteUser(), TestGujaratMachineAccessScenario()

### Community 58 - "config.go"
Cohesion: 0.06
Nodes (51): request, main(), Config, MachineSelection, EnsureCertificates(), existingCertificateExpiry(), generateCertificates(), Paths (+43 more)

### Community 59 - ".validate"
Cohesion: 0.47
Nodes (3): BreakFields, CreateBreakRequest, UpdateBreakRequest

### Community 60 - "UserProfileRepository"
Cohesion: 0.19
Nodes (5): UserProfile, UserProfileRepository, mapProfileError(), toUserProfileResponseError(), UserProfileService

### Community 61 - "Test list"
Cohesion: 0.04
Nodes (45): 1. The rows are consistent (SQL, no machine needed), 2. The rows agree with the mode log (SQL), 3. The raw signals around one cycle (telemetry and logs), 4. At the machine (floor test), Before the tests (once), Check query, `CNC_CYCLE_MAX_HOURS` (24 hours), `CNC_STALE_SEC` (30 s) and `CNC_QUIET_SEC` (300 s) (+37 more)

### Community 62 - "0001. Go standard library net/http, no web framework"
Cohesion: 0.33
Nodes (5): 0001. Go standard library net/http, no web framework, Alternatives considered, Consequences, Context, Decision

### Community 63 - "0003. PostgreSQL with pgx and hand-written SQL"
Cohesion: 0.33
Nodes (5): 0003. PostgreSQL with pgx and hand-written SQL, Alternatives considered, Consequences, Context, Decision

### Community 64 - ".validate"
Cohesion: 0.47
Nodes (3): CreateDowntimeReasonRequest, DowntimeReasonFields, UpdateDowntimeReasonRequest

### Community 65 - "time.Duration"
Cohesion: 0.06
Nodes (26): 0021. Read PLC variables over OPC UA with gopcua, in-process beside the CNC collector, Alternatives considered, Consequences, Context, Decision, Decision, PlcNodeTestResult, checkPlcNodeID() (+18 more)

### Community 66 - "UserLineService"
Cohesion: 0.07
Nodes (20): CreateUserCellRequest, UpdateUserCellRequest, validateUserCell(), CreateUserSupervisorRequest, UpdateUserSupervisorRequest, validateUserSupervisor(), NewUserCellHandler(), NewUserSupervisorHandler() (+12 more)

### Community 67 - "Execute"
Cohesion: 0.05
Nodes (33): main(), run(), Decision, Execute(), Branch, AxisStartValue, AxisVariable, Role (+25 more)

### Community 68 - "MachineAccuracyService"
Cohesion: 0.18
Nodes (7): MachineAccuracyRequest, NewMachineAccuracyHandler(), MachineAccuracy, MachineAccuracyRepository, mapMachineAccuracyError(), MachineAccuracyService, toMachineAccuracyResponseError()

### Community 69 - ".validate"
Cohesion: 0.47
Nodes (3): CreateIdleReasonRequest, IdleReasonFields, UpdateIdleReasonRequest

### Community 70 - "0014. Session modes and the login conflict, on the existing auth_tokens sessions"
Cohesion: 0.40
Nodes (4): 0014. Session modes and the login conflict, on the existing auth_tokens sessions, Alternatives considered, Consequences, Context

### Community 71 - "onDataCallback(ptr uintptr) uintptr"
Cohesion: 0.36
Nodes (9): Package-level callback singletons (syscall.NewCallback), CNCConnect param errCb (uintptr stdcall OnErr(int32 code, const char* msg)), Mutex-protected copy: dataMutex sync.RWMutex -> lastRecord/hasRecord, onDataCallback(ptr uintptr) uintptr, onErrorCallback(code, messagePtr uintptr) uintptr, §7 Asynchronous Callback Architecture, onDataCallback, onErrorCallback (+1 more)

### Community 72 - "time.Time"
Cohesion: 0.04
Nodes (74): Decision, FinishCyclesSummary, FinishedCycle, HomeAlarm, HomeCycles, HomeFleet, HomeMachine, HomePlan (+66 more)

### Community 73 - "MachineTelemetryService"
Cohesion: 0.16
Nodes (5): Step 7 — Listen for database changes, ListenForDatabaseNotifications(), OpenModePeriod, currentMode(), MachineTelemetryService

### Community 74 - "ControllerService"
Cohesion: 0.06
Nodes (20): CreateBrandRequest, UpdateBrandRequest, validateCatalogFields(), CreateCategoryRequest, UpdateCategoryRequest, CreateControllerRequest, UpdateControllerRequest, validateControllerFields() (+12 more)

### Community 75 - "testing.T"
Cohesion: 0.09
Nodes (39): TestSnapshotModeFieldOffsets(), TestSnapshotSize(), NewDashboardHandler(), sendRequest(), mustDecode(), newDashboardRouter(), TestAxisValues(), TestDashboardCharts() (+31 more)

### Community 76 - "0008. "Log in as customer" uses a separate session with a path-scoped cookie"
Cohesion: 0.40
Nodes (4): 0008. "Log in as customer" uses a separate session with a path-scoped cookie, Alternatives considered, Consequences, Context

### Community 77 - "0022. Test PLC node IDs in a short OPC UA session of their own"
Cohesion: 0.40
Nodes (4): 0022. Test PLC node IDs in a short OPC UA session of their own, Alternatives considered, Consequences, Context

### Community 78 - "UserCountryRepository"
Cohesion: 0.23
Nodes (5): CreateUserCountryRequest, UserCountry, UserCountryRepository, mapUserCountryError(), toUserCountryResponseError()

### Community 80 - "Event-driven data flow (no polling in client)"
Cohesion: 0.60
Nodes (5): Flow 1: load DLL, register two stdcall callbacks (syscall.NewCallback), Flow 2: CNCConnect starts DLL background worker thread, Flow 3: worker polls controller, assembles snapshot, invokes data callback, Flow 4: callback copies payload, signals renderer via buffered channel, Flow 5: main goroutine renders live ANSI dashboard

### Community 81 - "TestOrganizationForAnotherUser"
Cohesion: 0.40
Nodes (5): WithAuthUser(), asTestUser(), containsOrgID(), TestOrganizationForAnotherUser(), orgRecord

### Community 82 - "Live data — how a machine value reaches the screen"
Cohesion: 0.14
Nodes (11): 0009. Socket.IO live data: the collector pushes each screen's variables, checked per user, Alternatives considered, Consequences, Context, Adding a variable to a screen, Connecting, Live data — Socket.IO event contract, Ordering (+3 more)

### Community 83 - "PartRepository"
Cohesion: 0.15
Nodes (5): PartProcess, PartSubProgram, executeAffected(), PartRepository, mapPartError()

### Community 84 - "CNCConnect(ip, port, dataCb, errCb, machineType) int32"
Cohesion: 0.20
Nodes (11): CNCConnect(ip, port, dataCb, errCb, machineType) int32, CNCConnect param dataCb (uintptr stdcall OnData(const CncMasterSnapshot*)), CNCConnect param ip (const char*, ANSI NUL-terminated), CNCConnect param port (int32, e.g. 4840), CNCdisconnect(machineId int32), LazyDLL binding (syscall.NewLazyDLL + eager mod.Load in init), §4 API Reference, §5 DLL Loading & Binding (Go) (+3 more)

### Community 85 - "MachineTelemetryRepository"
Cohesion: 0.11
Nodes (8): MachineEvent, MachineLiveStatus, ModeChange, MachineTelemetryRepository, mapMachineTelemetryError(), execute(), UserMachineAccessRepository, NewUserMachineAccessRepository()

### Community 86 - "auth_tokens SHA-256 Hashed Sessions (mig"
Cohesion: 0.67
Nodes (3): auth_tokens SHA-256 Hashed Sessions (migration 014), Per-Request Session and User Re-check, ts_session Cookie (HttpOnly, SameSite=Strict)

### Community 88 - "context.Context"
Cohesion: 0.21
Nodes (3): ClienteleRepository, mapClienteleError(), levelInMachineFilter()

### Community 89 - "HealthHandler"
Cohesion: 0.22
Nodes (5): Decision, Decision, isLoginAsKey(), RouteLoginAs(), statusRecorder

### Community 104 - "connection.go"
Cohesion: 0.08
Nodes (19): changeOf(), Connect(), forwardNotifications(), Change, Item, ItemResult, isBad(), parseNodeID() (+11 more)

### Community 105 - "0005. CNC collector runs inside the API process"
Cohesion: 0.33
Nodes (5): 0005. CNC collector runs inside the API process, Alternatives considered, Consequences, Context, Decision

### Community 107 - "0020. Machine PLC node IDs: the machine's own node ID in machine_variables, default otherwise"
Cohesion: 0.33
Nodes (5): 0020. Machine PLC node IDs: the machine's own node ID in machine_variables, default otherwise, Alternatives considered, Consequences, Context, Decision

### Community 108 - "UserService"
Cohesion: 0.14
Nodes (16): CustomerDetails, MachineFilter, UpdateUserRequest, UserListFilter, UserResponse, ParseDayStart(), TestParseDayStart(), AssignableRoles() (+8 more)

### Community 109 - "0011. Ponytail plugin for Claude Code, under this repo's rules"
Cohesion: 0.33
Nodes (5): 0011. Ponytail plugin for Claude Code, under this repo's rules, Alternatives considered, Consequences, Context, Decision

### Community 110 - "Client Certificate & Key (client_cert.der / client_key.der)"
Cohesion: 0.48
Nodes (7): client_cert.der/.pem (X.509 client cert for Siemens OPC UA), client_key.der/.pem (client private key), Siemens certs: client_cert.der + client_key.der in binary directory, Siemens trust: cert must be trusted in NCU Access Administration, Troubleshoot: BadRequestNotAllowed -> Export client_cert, BadRequestNotAllowed (certificate trust rejection), Client Certificate & Key (client_cert.der / client_key.der)

### Community 111 - "0015. Operator tasks: a shared task list, assigned per customer and machine model"
Cohesion: 0.50
Nodes (4): 0015. Operator tasks: a shared task list, assigned per customer and machine model, Alternatives considered, Consequences, Context

### Community 112 - "user_access_test.go"
Cohesion: 0.32
Nodes (5): LicenceKeyDuration, LiveAccessChanged, LiveUnsubscribeRequest, licenceKeyDurations(), Envelope

### Community 113 - ".validate"
Cohesion: 0.50
Nodes (3): isTimeZone(), isWebsite(), normalizeFilters()

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
Cohesion: 0.33
Nodes (5): 0017. Machine accuracy as one shared list, managed by the superadmin, Alternatives considered, Consequences, Context, Decision

### Community 121 - "0006. Page-wise dashboard endpoints returning only chart fields"
Cohesion: 0.40
Nodes (4): 0006. Page-wise dashboard endpoints returning only chart fields, Alternatives considered, Consequences, Context

### Community 124 - ".validate"
Cohesion: 0.38
Nodes (4): ClienteleFields, CreateClienteleRequest, UpdateClienteleRequest, validClienteleEmail()

### Community 126 - "UserStateService"
Cohesion: 0.18
Nodes (7): CreateUserStateRequest, validateUserState(), UserState, UserStateRepository, mapUserStateError(), UserStateService, toUserStateResponseError()

### Community 130 - "LicenceKeyService"
Cohesion: 0.10
Nodes (16): NewLicenceKey, CreateLicenceKeysRequest, GenerateLicenceKeysRequest, UpdateLicenceKeyRequest, validDuration(), ValidLicenceKey(), NewLicenceKeyHandler(), GeneratedLicenceKeys (+8 more)

### Community 134 - "pathID"
Cohesion: 0.08
Nodes (14): MachineCycleHandler, MachineDiagnosticsHandler, MachinePlcVariableHandler, pathID(), targetUserID(), UserCellHandler, UserCityHandler, UserCountryHandler (+6 more)

### Community 136 - "0013. Machine diagnostics: reasons from what the collector recorded, network tested from the server"
Cohesion: 0.33
Nodes (5): 0013. Machine diagnostics: reasons from what the collector recorded, network tested from the server, Alternatives considered, Consequences, Context, Decision

## Ambiguous Edges - Review These
- `CncCommunication.dll (32-bit vendor DLL)` → `FanucCncCommunication.dll`  [AMBIGUOUS]
  docs/setup-guide.txt · relation: semantically_similar_to
- `Snapshot Fields (positions, spindle, feed, program, tool, alarms)` → `Machines API (/machines)`  [AMBIGUOUS]
  docs/CncCommunication.dll — API & Integration Guide (Go Client).pdf · relation: conceptually_related_to

## Knowledge Gaps
- **33 isolated node(s):** `7th-sense-server`, `middleware.RequireRole / routes.adminOnly`, `response.SendResponse`, `Test Auth (asTestUser, tests/auth_test.go)`, `user_filters Filter Checklist` (+28 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 246 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **19 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `CNC Data Collection Flow` connect `Collector` to `Worker (internal/monitor/worker.go)`, `Live data — how a machine value reaches the screen`, `worker`, `dllDriver`?**
  _High betweenness centrality (0.130) - this node is a cross-community bridge._
- **What connects `7th-sense-server`, `middleware.RequireRole / routes.adminOnly`, `response.SendResponse` to the rest of the system?**
  _33 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `WithTransaction` be split into smaller, more focused modules?**
  _Cohesion score 0.1350806451612903 - nodes in this community are weakly interconnected._
- **Why does `2. All steps at a glance` connect `worker` to `Collector`, `dllDriver`?**
  _High betweenness centrality (0.109) - this node is a cross-community bridge._
- **Should `MachineService` be split into smaller, more focused modules?**
  _Cohesion score 0.06349206349206349 - nodes in this community are weakly interconnected._

### Low-confidence Hints
_AMBIGUOUS edges — the extractor was unsure. Verify before acting on these._

- **What is the exact relationship between `CncCommunication.dll (32-bit vendor DLL)` and `FanucCncCommunication.dll`?**
  _Edge tagged AMBIGUOUS (relation: semantically_similar_to) - confidence is low._
- **What is the exact relationship between `Snapshot Fields (positions, spindle, feed, program, tool, alarms)` and `Machines API (/machines)`?**
  _Edge tagged AMBIGUOUS (relation: conceptually_related_to) - confidence is low._