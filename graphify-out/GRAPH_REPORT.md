# Graph Report - 7th-sense-server  (2026-11-05)

## Corpus Check
- 296 files · ~234,138 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 4 file(s) not represented in the graph (top: (none) 3, .toml 1)

## Summary
- 2619 nodes · 21132 edges · 109 communities (91 shown, 18 thin omitted)
- Extraction: 95% EXTRACTED · 5% INFERRED · 0% AMBIGUOUS · INFERRED: 1009 edges (avg confidence: 0.87)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `6e4eb5b4`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- net/http.Handler
- UserCityService
- UserShopService
- Machine
- net/http.Request
- pathID()
- net/http.Request
- UserStateService
- api_helpers_test.go
- CncMasterSnapshot (34,440 bytes, #pragma
- LiveHandler
- Connection
- ProductionScheduleRepository
- Collector
- worker
- Hub
- execute
- mode-chart.js
- MachineModelService
- validText
- net/http.ServeMux
- CNCConnect(ip, port, dataCb, errCb, machineType) int32
- UserSupervisorService
- §10 Diagnostics & Troubleshooting
- decodeAndValidate
- machine_mode.go
- .validate
- newMachineService
- run
- cycleReading
- Execute
- testing.T
- Worker (internal/monitor/worker.go)
- VariableConfig
- Address
- .Range
- LicenceKeyService
- Client Certificate & Key (client_cert.der / client_key.der)
- Live data — how a machine value reaches the screen
- formatDuration
- CncCommunication.dll API & Integration G
- time.Duration
- UserService
- CncCommunication.dll API & Integration G
- AuthUser
- .Login()
- NewError
- .Validate
- Machines API (/machines)
- context.Context
- Product
- 2. All steps at a glance
- BranchService
- README.md
- AuthHandler
- expectStatus
- SendResponse
- User
- apidocs/main.go
- .validate
- UserLocationService
- Test list
- newCatalogRouterAs
- CompanyProfile
- CNC Data Collection (in-process)
- liveTestClient
- UserLineService
- .accessFor
- .validate
- drawChart
- run
- onDataCallback(ptr uintptr) uintptr
- .validate
- .validate
- user_access_test.go
- uniqueSuffix
- Fixed values and why they are needed
- HomeFleet
- time.Time
- MachineTelemetryService
- Event-driven data flow (no polling in client)
- UserMachineFilter
- Live data — Socket.IO event contract
- AuthRepository
- 8. Reference
- Querier
- auth_tokens SHA-256 Hashed Sessions (mig
- PartDrawingStorage
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
- 0005. CNC collector runs inside the API process
- 0007. OpenAPI 3 from swag annotations, Swagger UI behind API_DOCS, generated Postman
- HolidayService
- 0006. Page-wise dashboard endpoints returning only chart fields
- TestClienteleCRUD

## God Nodes (most connected - your core abstractions)
1. `Execute()` - 191 edges
2. `SendResponse()` - 159 edges
3. `NewError()` - 141 edges
4. `run()` - 99 edges
5. `AuthUser` - 89 edges
6. `pathID()` - 80 edges
7. `currentUser()` - 64 edges
8. `Querier` - 62 edges
9. `decodeAndValidate()` - 59 edges
10. `worker` - 58 edges

## Surprising Connections (you probably didn't know these)
- `Machine Identifier Mapping (0 Simulator,1 Fanuc,2 Siemens,3 Syntec,4 Mitsubishi)` --semantically_similar_to--> `Machine Catalog Brand-Controller-Model`  [INFERRED] [semantically similar]
  docs/CncCommunication.dll — API & Integration Guide (Go Client).pdf → .claude/CLAUDE.md
- `4. Socket.IO → browser` --references--> `LiveHandler`  [INFERRED]
  docs/live-data-flow.md → internal/handler/live_handler.go
- `Decision` --references--> `SendResponse()`  [INFERRED]
  docs/adr/0002-layered-handler-service-repository.md → internal/response/sendResponse.go
- `Snapshot Fields (positions, spindle, feed, program, tool, alarms)` --conceptually_related_to--> `Machines API (/machines)`  [AMBIGUOUS]
  docs/CncCommunication.dll — API & Integration Guide (Go Client).pdf → .claude/CLAUDE.md
- `What a screen asks for (subscribing)` --references--> `TestHubSnapshotHasOnlyRequestedCodes()`  [INFERRED]
  docs/live-data-flow.md → internal/live/hub_test.go

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

## Communities (109 total, 18 thin omitted)

### Community 0 - "net/http.Handler"
Cohesion: 0.19
Nodes (19): Handler(), CrossOrigin(), Logger(), Chain(), Middleware, Recoverer(), NewUserRepository(), New() (+11 more)

### Community 1 - "UserCityService"
Cohesion: 0.15
Nodes (9): CreateUserCityRequest, UpdateUserCityRequest, validateUserCity(), NewUserCityHandler(), UserCity, UserCityRepository, mapUserCityError(), UserCityService (+1 more)

### Community 2 - "UserShopService"
Cohesion: 0.13
Nodes (12): CreateUserShopRequest, UpdateUserShopRequest, validateUserShop(), NewUserShopHandler(), UserShop, UserShopRepository, mapUserShopError(), NewUserShopRepository() (+4 more)

### Community 3 - "Machine"
Cohesion: 0.12
Nodes (17): Decision, CreateMachineRequest, MachineFields, UpdateMachineRequest, validCoordinate(), validHealthIDs(), validIPv4(), Machine (+9 more)

### Community 4 - "net/http.Request"
Cohesion: 0.09
Nodes (11): BreakHandler, ClienteleHandler, DowntimeReasonHandler, HolidayHandler, IdleReasonHandler, MachineCycleHandler, PartHandler, ProductionScheduleHandler (+3 more)

### Community 5 - "pathID()"
Cohesion: 0.10
Nodes (12): pathID(), targetUserID(), RoleHandler, UserCellHandler, UserCityHandler, UserCountryHandler, UserLineHandler, UserLocationHandler (+4 more)

### Community 6 - "net/http.Request"
Cohesion: 0.12
Nodes (9): ControllerHandler, dashboardQuery(), DashboardHandler, readDashboardRequest(), isSameOrigin(), MachineHandler, listUserID(), queryBool() (+1 more)

### Community 7 - "UserStateService"
Cohesion: 0.09
Nodes (14): CreateUserCountryRequest, CreateUserStateRequest, validateUserState(), NewUserStateHandler(), UserCountry, UserState, UserCountryRepository, mapUserCountryError() (+6 more)

### Community 8 - "api_helpers_test.go"
Cohesion: 0.40
Nodes (6): cncMasterSnapshot, LanguageResponse, ThemeResponse, authUserKey, loginAsKey, machineModelRecord

### Community 9 - "CncMasterSnapshot (34,440 bytes, #pragma"
Cohesion: 0.10
Nodes (43): CncMasterSnapshot (34,440 bytes, #pragma pack(8), x86), AbsPos ([5]float64, Inches / Degrees, Work Coordinate System (WCS) positions), ActiveGCode ([16]byte, ASCII, Currently active G-code group), ActiveProgPath ([256]byte, ASCII Path, Storage path of active NC artifact), AlarmMsg ([128]byte, UTF-8 / ASCII, Alarm description message), AlarmNo (int32, Code (0 = OK), Primary alarm or fault number), AxisLoadXPercentage (int16, Percent (%), Servo X axis drive load (AxisLoadX/Y/ZPercentage)), AxisLoadYPercentage (int16, Percent (%), Servo Y axis drive load (AxisLoadX/Y/ZPercentage)) (+35 more)

### Community 10 - "LiveHandler"
Cohesion: 0.11
Nodes (13): Browser → server, liveSocket, liveSubscription, LiveSubscribeResponse, channelRoom(), closeLiveSocket(), decodeLiveMessage(), LiveHandler (+5 more)

### Community 11 - "Connection"
Cohesion: 0.10
Nodes (14): dllDriver, slot, Step 12 — Call the DLL's CNCConnect, Step 13 — The DLL sends data, 1. Machine → DLL → server, Connection, Reading, Stats (+6 more)

### Community 12 - "ProductionScheduleRepository"
Cohesion: 0.14
Nodes (12): MachineQuantity, ProductionSchedule, ProductionScheduleItem, ProductionScheduleRepository, mapProductionScheduleError(), NewProductionScheduleRepository(), ProductionScheduleService, NewProductionScheduleService() (+4 more)

### Community 13 - "Collector"
Cohesion: 0.13
Nodes (17): Log lines (in `logs/app.log`, search for the message), Step 18 — React to machine and variable changes, Step 20 — Stop the workers and mark machines offline, Step 6 — Record the start and begin following the database, Step 9 — Start one worker per machine, CNCLogSettings, equalEngine(), MonitoredMachine (+9 more)

### Community 14 - "worker"
Cohesion: 0.10
Nodes (21): 5. Connect and collect, per machine: steps 10 – 17, Known issue, Retry and timing summary, Step 10 — Get ready to connect, Step 11 — Connect, a few machines at a time, Step 14 — The data collection loop, Step 15 — Store the reading, Step 16 — Every second: dead connection check and live status (+13 more)

### Community 15 - "Hub"
Cohesion: 0.11
Nodes (25): 2. Worker → hub (before the database), EventRoom(), Hub, Snapshot, NewHub(), StatusRoom(), newRecordingHub(), TestHubJoinValuesSharesRoomsAndReportsUnknownCodes() (+17 more)

### Community 16 - "execute"
Cohesion: 0.15
Nodes (4): MachineEvent, MachineLiveStatus, ModeChange, execute()

### Community 17 - "mode-chart.js"
Cohesion: 0.21
Nodes (23): chartData(), countByGroup(), drawDashboard(), drawHeader(), drawMachineList(), drawTiles(), drawTotals(), fetchData() (+15 more)

### Community 18 - "MachineModelService"
Cohesion: 0.18
Nodes (8): MachineModelRequest, NewMachineModelHandler(), MachineModel, MachineModelRepository, mapMachineModelError(), NewMachineModelRepository(), MachineModelService, toMachineModelResponseError()

### Community 19 - "validText"
Cohesion: 0.14
Nodes (7): PartProcessFields, PartSubProgramFields, UploadPartDrawingRequest, uniquePositiveIDs(), validText(), validWorkingDays(), PartDrawing

### Community 20 - "net/http.ServeMux"
Cohesion: 0.14
Nodes (35): Decision, RequireRole(), registerAuthRoutes(), registerBranchRoutes(), registerBrandRoutes(), registerBreakRoutes(), registerCategoryRoutes(), registerClienteleRoutes() (+27 more)

### Community 21 - "CNCConnect(ip, port, dataCb, errCb, machineType) int32"
Cohesion: 0.20
Nodes (11): CNCConnect(ip, port, dataCb, errCb, machineType) int32, CNCConnect param dataCb (uintptr stdcall OnData(const CncMasterSnapshot*)), CNCConnect param ip (const char*, ANSI NUL-terminated), CNCConnect param port (int32, e.g. 4840), CNCdisconnect(machineId int32), LazyDLL binding (syscall.NewLazyDLL + eager mod.Load in init), §4 API Reference, §5 DLL Loading & Binding (Go) (+3 more)

### Community 22 - "UserSupervisorService"
Cohesion: 0.14
Nodes (11): CreateUserSupervisorRequest, UpdateUserSupervisorRequest, validateUserSupervisor(), NewUserSupervisorHandler(), UserSupervisor, UserSupervisorRepository, mapUserSupervisorError(), NewUserSupervisorRepository() (+3 more)

### Community 23 - "§10 Diagnostics & Troubleshooting"
Cohesion: 0.12
Nodes (25): Check mode: go run sinumerik_client.go check (OPC UA node existence report), Go Client (sinumerik_client.go, callback -> copy -> render), Prerequisite: 32-bit x86 DLL, compile Go with GOARCH=386, Prerequisite: SIEMENS_USER / SIEMENS_PASS env vars before CNCConnect, Run: go run sinumerik_client.go (live dashboard), §3.1 Environment & Prerequisites, §8 Running the Go Client, §10 Diagnostics & Troubleshooting (+17 more)

### Community 24 - "decodeAndValidate"
Cohesion: 0.17
Nodes (5): validator, BranchHandler, currentUserID(), decodeAndValidate(), UserHandler

### Community 25 - "machine_mode.go"
Cohesion: 0.14
Nodes (20): decideMachineMode(), holdReason(), machineStateFromReading(), numberingEngine(), programStateFromProcStatus(), ncReading(), TestDecideMachineMode(), TestDecideMachineModeConditions() (+12 more)

### Community 26 - ".validate"
Cohesion: 0.47
Nodes (3): CreateIdleReasonRequest, IdleReasonFields, UpdateIdleReasonRequest

### Community 27 - "newMachineService"
Cohesion: 0.19
Nodes (25): NewAuthHandler(), NewUserCountryHandler(), Authenticate(), WithAuthUser(), NewBranchRepository(), NewRoleRepository(), NewUserCityRepository(), NewUserCountryRepository() (+17 more)

### Community 28 - "run"
Cohesion: 0.06
Nodes (41): run(), Connect(), CreateCategoryRequest, UpdateCategoryRequest, clockMinutes(), CreateShiftRequest, ShiftFields, UpdateShiftRequest (+33 more)

### Community 29 - "cycleReading"
Cohesion: 0.19
Nodes (7): CycleStatusFor(), decideCycleStep(), isProgramLoop(), resetReason(), cycleReading, cycleStep, cycleTracker

### Community 30 - "Execute"
Cohesion: 0.13
Nodes (5): Execute(), MachineTelemetryRepository, mapMachineTelemetryError(), UserProfileRepository, mapProfileError()

### Community 31 - "testing.T"
Cohesion: 0.11
Nodes (21): UpdateThemeRequest, TestSnapshotModeFieldOffsets(), TestSnapshotSize(), TestUpdateThemeRequestValidate(), mustDecode(), newDashboardRouter(), TestDashboardCharts(), TestHomeDashboard() (+13 more)

### Community 32 - "Worker (internal/monitor/worker.go)"
Cohesion: 0.07
Nodes (22): 7th Sense Server Claude Context, response.SendResponse Format, User CRUD Golden Reference, Collector (internal/monitor/collector.go), CNC Driver (internal/cnc), Known Issue: nil live-value cache panic, machine_events, machine_telementry hypertable (+14 more)

### Community 33 - "VariableConfig"
Cohesion: 0.22
Nodes (7): VariableConfig, formatPreviousValue(), formatValue(), newValueChangeLog(), boundVariable, valueChange, valueChangeLog

### Community 34 - "Address"
Cohesion: 0.22
Nodes (8): Address, AddressRepository, mapAddressError(), NewAddressRepository(), NewAddressService(), toAddressResponseError(), AddressService, TestAddressCRUD()

### Community 36 - "LicenceKeyService"
Cohesion: 0.12
Nodes (14): NewLicenceKey, CreateLicenceKeysRequest, GenerateLicenceKeysRequest, NewLicenceKeyHandler(), GeneratedLicenceKeys, LicenceKey, LicenceKeyRepository, mapLicenceKeyError() (+6 more)

### Community 37 - "Client Certificate & Key (client_cert.der / client_key.der)"
Cohesion: 0.48
Nodes (7): client_cert.der/.pem (X.509 client cert for Siemens OPC UA), client_key.der/.pem (client private key), Siemens certs: client_cert.der + client_key.der in binary directory, Siemens trust: cert must be trusted in NCU Access Administration, Troubleshoot: BadRequestNotAllowed -> Export client_cert, BadRequestNotAllowed (certificate trust rejection), Client Certificate & Key (client_cert.der / client_key.der)

### Community 38 - "Live data — how a machine value reaches the screen"
Cohesion: 0.13
Nodes (13): 3. Hub → rooms, 4. Socket.IO → browser, Connecting the socket, Files, How the variable list works, Live data — how a machine value reaches the screen, Showing a new variable on a screen, Step by step (+5 more)

### Community 39 - "formatDuration"
Cohesion: 0.15
Nodes (30): card(), cell(), cycleData(), cycleEnd(), cycleEndReason(), cycleJob(), cycleList(), cycleNumber() (+22 more)

### Community 40 - "CncCommunication.dll API & Integration G"
Cohesion: 0.19
Nodes (20): CncCommunication.dll (C++ core runtime, x86, OPC UA inside), CncCommunication.dll API & Integration Guide (Go Client) - text version, Fanuc profile (engine 1, port 8193, Fwlib32.dll + fwlibe1.dll), Fwlib32.dll (Fanuc FOCAS2), fwlibe1.dll (Fanuc FOCAS2), libcrypto-3.dll (OpenSSL), libssl-3.dll (OpenSSL), Engine 1: Fanuc (FOCAS2 Ethernet, fwlib32) (+12 more)

### Community 41 - "time.Duration"
Cohesion: 0.36
Nodes (4): connectFanucCounter(), main(), run(), fanucCounter

### Community 42 - "UserService"
Cohesion: 0.06
Nodes (33): Decision, WithTransaction(), CreateUserRequest, CustomerDetails, UpdateUserRequest, UserListFilter, UserResponse, isTimeZone() (+25 more)

### Community 43 - "CncCommunication.dll API & Integration G"
Cohesion: 0.13
Nodes (18): CNCConnect param machineType (int32 engine id, see §4.2), Engine 0: Simulator (local loopback mock engine), Engine 2: Siemens (Sinumerik OPC UA, active profile), Engine 3: Syntec (Open API / TCP), Engine 4: Mitsubishi (driver not specified), Machine identifier mapping (engine id enum 0-4), §9 Controller-Specific Profiles, Siemens default ports: 4840 (native NCU) / 4841 (custom) (+10 more)

### Community 44 - "AuthUser"
Cohesion: 0.06
Nodes (38): LoginAsResponse, CreatePartRequest, PartFields, UpdatePartRequest, ValidPartDrawingFile(), AuthUser, Break, Clientele (+30 more)

### Community 45 - ".Login()"
Cohesion: 0.15
Nodes (11): LoginRequest, NewLiveHandler(), AuthService, hashToken(), newLoginAsKey(), newSessionToken(), truncate(), newLoginThrottle() (+3 more)

### Community 46 - "NewError"
Cohesion: 0.15
Nodes (10): CreateAddressRequest, UpdateAddressRequest, UpdateLanguageRequest, UpdateLicenceKeyRequest, validDuration(), ValidLicenceKey(), queryIDList(), NewError() (+2 more)

### Community 47 - ".Validate"
Cohesion: 0.25
Nodes (11): MachineFilter, normalizeTheme(), TestCreateUserRequestValidate(), TestCreateUserRequestValidateDropsDuplicateMachines(), TestCreateUserRequestValidateNormalizes(), TestCustomerFiltersValidate(), TestUpdateUserRequestValidateStatus(), validCustomerDetails() (+3 more)

### Community 48 - "Machines API (/machines)"
Cohesion: 0.16
Nodes (18): Clienteles API, Company Scope (companyOwner / ownCompanyCustomer), Current User API (/me theme, language), Customer Display Names (customerName / companyName), Holidays API, Licence Keys API, Machine Catalog Brand-Controller-Model, Machines API (/machines) (+10 more)

### Community 49 - "context.Context"
Cohesion: 0.06
Nodes (23): One running cycle per machine, MachineCyclesResponse, CycleJob, MachineCycle, PartProcess, PartSubProgram, ShiftBreak, BreakRepository (+15 more)

### Community 50 - "Product"
Cohesion: 0.14
Nodes (12): CreateProductRequest, UpdateProductRequest, validateProductFields(), NewProductHandler(), Product, ProductRepository, mapProductError(), NewProductRepository() (+4 more)

### Community 51 - "2. All steps at a glance"
Cohesion: 0.09
Nodes (24): Supervisor, 1. The big picture, 2. All steps at a glance, 3. Startup: steps 1 – 4, 6. When data is edited: step 18, 7. Shutdown: steps 19 – 20, 9. Data Reading Flow (Function Names Only), docs/cnc-communication.md (+16 more)

### Community 52 - "BranchService"
Cohesion: 0.14
Nodes (10): CreateBranchRequest, UpdateBranchRequest, validateBranchFields(), NewBranchHandler(), Branch, BranchRepository, mapBranchError(), cityOwner() (+2 more)

### Community 53 - "README.md"
Cohesion: 0.06
Nodes (28): 0001. Go standard library net/http, no web framework, Alternatives considered, Consequences, Context, Decision, 0002. Layered architecture: handler, DTO, service, repository, Alternatives considered, Consequences (+20 more)

### Community 54 - "AuthHandler"
Cohesion: 0.22
Nodes (8): Decision, AuthHandler, CurrentUser(), isLoginAsKey(), LoginAsKey(), RouteLoginAs(), SessionToken(), UserID()

### Community 55 - "expectStatus"
Cohesion: 0.23
Nodes (28): decodeData(), expectStatus(), sendRequest(), apiResponse, generateLicenceKeys(), licenceKeyItems(), TestLicenceKeyCRUD(), actAs() (+20 more)

### Community 56 - "SendResponse"
Cohesion: 0.08
Nodes (6): BrandHandler, CategoryHandler, LicenceKeyHandler, MachineModelHandler, ProductHandler, statusRecorder

### Community 57 - "User"
Cohesion: 0.31
Nodes (4): User, UserRepository, mapUserError(), userServices

### Community 58 - "apidocs/main.go"
Cohesion: 0.06
Nodes (46): request, main(), run(), main(), Config, MachineSelection, GenerateCertificates(), getEnv() (+38 more)

### Community 59 - ".validate"
Cohesion: 0.38
Nodes (4): ClienteleFields, CreateClienteleRequest, UpdateClienteleRequest, validClienteleEmail()

### Community 60 - "UserLocationService"
Cohesion: 0.14
Nodes (10): CreateUserLocationRequest, UpdateUserLocationRequest, validateUserLocation(), NewUserLocationHandler(), UserLocation, UserLocationRepository, mapUserLocationError(), NewUserLocationRepository() (+2 more)

### Community 61 - "Test list"
Cohesion: 0.05
Nodes (38): 1. The rows are consistent (SQL, no machine needed), 2. The rows agree with the mode log (SQL), 3. The raw signals around one cycle (telemetry and logs), 4. At the machine (floor test), Before the tests (once), Check query, Completed, Floor tests on the HUMA panel (Fanuc) (+30 more)

### Community 62 - "newCatalogRouterAs"
Cohesion: 0.06
Nodes (21): CreateBrandRequest, UpdateBrandRequest, validateCatalogFields(), CreateControllerRequest, UpdateControllerRequest, validateControllerFields(), NewBrandHandler(), NewControllerHandler() (+13 more)

### Community 63 - "CompanyProfile"
Cohesion: 0.23
Nodes (8): CompanyProfile, CompanyProfileRepository, mapCompanyProfileError(), NewCompanyProfileRepository(), NewCompanyProfileService(), toCompanyProfileResponseError(), CompanyProfileService, TestCompanyProfileCRUD()

### Community 64 - "CNC Data Collection (in-process)"
Cohesion: 0.15
Nodes (12): CNC Data Collection (in-process), writeLog CNC Log Types and Switches, .env Configuration (internal/config), Health Mode History API (public), CncCommunication.dll (32-bit vendor DLL), machine_mode_log, Mode Decision (machine_mode.go), ECharts 5.5.1 (+4 more)

### Community 65 - "liveTestClient"
Cohesion: 0.17
Nodes (13): authRequest(), createAuthTestUser(), login(), sessionCookieOf(), TestAuthentication(), dialLive(), TestLiveData(), liveTestClient (+5 more)

### Community 66 - "UserLineService"
Cohesion: 0.07
Nodes (21): CreateUserCellRequest, UpdateUserCellRequest, validateUserCell(), CreateUserLineRequest, UpdateUserLineRequest, validateUserLine(), NewUserCellHandler(), NewUserLineHandler() (+13 more)

### Community 67 - ".accessFor"
Cohesion: 0.25
Nodes (4): SeesAllCompanyMachines(), asInvalidSelection(), checkJobsOverlap(), jobSpan

### Community 68 - ".validate"
Cohesion: 0.47
Nodes (3): CreateDowntimeReasonRequest, DowntimeReasonFields, UpdateDowntimeReasonRequest

### Community 69 - "drawChart"
Cohesion: 0.24
Nodes (11): chartRows(), currentSpan(), currentWindow(), drawChart(), drawZoomLevels(), markZoomLevel(), renderBar(), renderCycleBar() (+3 more)

### Community 70 - "run"
Cohesion: 0.36
Nodes (10): checkNotTaken(), createSuperadmin(), ensureProfile(), execute(), main(), parseFlags(), run(), superadminRoleID() (+2 more)

### Community 71 - "onDataCallback(ptr uintptr) uintptr"
Cohesion: 0.36
Nodes (9): Package-level callback singletons (syscall.NewCallback), CNCConnect param errCb (uintptr stdcall OnErr(int32 code, const char* msg)), Mutex-protected copy: dataMutex sync.RWMutex -> lastRecord/hasRecord, onDataCallback(ptr uintptr) uintptr, onErrorCallback(code, messagePtr uintptr) uintptr, §7 Asynchronous Callback Architecture, onDataCallback, onErrorCallback (+1 more)

### Community 72 - ".validate"
Cohesion: 0.47
Nodes (3): BreakFields, CreateBreakRequest, UpdateBreakRequest

### Community 73 - ".validate"
Cohesion: 0.27
Nodes (7): validDate(), CreateProductionScheduleRequest, ProductionScheduleFields, ProductionScheduleItemFields, UpdateProductionScheduleRequest, ScheduleDay(), shiftSpan()

### Community 74 - "user_access_test.go"
Cohesion: 0.30
Nodes (7): LicenceKeyDuration, LiveAccessChanged, LiveUnsubscribeRequest, licenceKeyDurations(), Envelope, liveValues, subscribeAnswer

### Community 75 - "uniqueSuffix"
Cohesion: 0.25
Nodes (16): containsID(), newCatalogRouter(), runCatalogCRUD(), uniqueSuffix(), createTestRecord(), TestBrandCRUD(), TestCategoryCRUD(), TestControllerBrand() (+8 more)

### Community 76 - "Fixed values and why they are needed"
Cohesion: 0.29
Nodes (7): `CNC_CYCLE_MAX_HOURS` (24 hours), `CNC_STALE_SEC` (30 s) and `CNC_QUIET_SEC` (300 s), `cycleCounterGrace` and `cycleStartGrace` (5 seconds each), `emergencyResetGrace` (2 seconds), Fixed values and why they are needed, `loopMinimumDropFraction` (0.5), `maxPartCountStep` (1)

### Community 77 - "HomeFleet"
Cohesion: 0.50
Nodes (3): HomeFleet, HomeModeCount, homeModeCounts()

### Community 78 - "time.Time"
Cohesion: 0.05
Nodes (75): Decision, FinishCyclesSummary, FinishedCycle, HomeAlarm, HomeCycles, HomeMachine, HomePlan, HomeTime (+67 more)

### Community 79 - "MachineTelemetryService"
Cohesion: 0.10
Nodes (13): 4. Find the machines: steps 5 – 9, Step 5 — Prepare the database, Step 7 — Listen for database changes, Step 8 — Load all active machines, ListenForDatabaseNotifications(), MachineModesResponse, NewHealthHandler(), MachineModeStatus (+5 more)

### Community 80 - "Event-driven data flow (no polling in client)"
Cohesion: 0.60
Nodes (5): Flow 1: load DLL, register two stdcall callbacks (syscall.NewCallback), Flow 2: CNCConnect starts DLL background worker thread, Flow 3: worker polls controller, assembles snapshot, invokes data callback, Flow 4: callback copies payload, signals renderer via buffered channel, Flow 5: main goroutine renders live ANSI dashboard

### Community 81 - "UserMachineFilter"
Cohesion: 0.38
Nodes (3): UserMachineFilter, UserMachineFilterRepository, userAccess

### Community 82 - "Live data — Socket.IO event contract"
Cohesion: 0.14
Nodes (11): 0009. Socket.IO live data: the collector pushes each screen's variables, checked per user, Alternatives considered, Consequences, Context, Adding a variable to a screen, Connecting, Live data — Socket.IO event contract, Ordering (+3 more)

### Community 83 - "AuthRepository"
Cohesion: 0.32
Nodes (3): AuthRepository, mapAuthError(), NewAuthRepository()

### Community 84 - "8. Reference"
Cohesion: 0.33
Nodes (6): 8. Reference, Events in `machine_events`, Related documents, Settings (`.env`, defaults in `internal/config/config.go`), Tables, Troubleshooting

### Community 85 - "Querier"
Cohesion: 0.22
Nodes (10): Querier, mapLocationError(), NewLocationRepository(), NewMachineCycleRepository(), UserMachineAccessRepository, NewUserMachineAccessRepository(), NewUserMachineFilterRepository(), NewUserProfileService() (+2 more)

### Community 86 - "auth_tokens SHA-256 Hashed Sessions (mig"
Cohesion: 0.67
Nodes (3): auth_tokens SHA-256 Hashed Sessions (migration 014), Per-Request Session and User Re-check, ts_session Cookie (HttpOnly, SameSite=Strict)

### Community 105 - "0005. CNC collector runs inside the API process"
Cohesion: 0.33
Nodes (5): 0005. CNC collector runs inside the API process, Alternatives considered, Consequences, Context, Decision

### Community 106 - "0007. OpenAPI 3 from swag annotations, Swagger UI behind API_DOCS, generated Postman"
Cohesion: 0.33
Nodes (5): 0007. OpenAPI 3 from swag annotations, Swagger UI behind API_DOCS, generated Postman, Alternatives considered, Consequences, Context, Decision

### Community 108 - "HolidayService"
Cohesion: 0.14
Nodes (13): CreateHolidayRequest, HolidayFields, UpdateHolidayRequest, NewHolidayHandler(), Holiday, HolidayRepository, mapHolidayError(), NewHolidayRepository() (+5 more)

### Community 110 - "0006. Page-wise dashboard endpoints returning only chart fields"
Cohesion: 0.40
Nodes (4): 0006. Page-wise dashboard endpoints returning only chart fields, Alternatives considered, Consequences, Context

### Community 118 - "TestClienteleCRUD"
Cohesion: 0.20
Nodes (16): NewMachineHandler(), clienteleBody(), TestClienteleCRUD(), containsMachineID(), insertTestRow(), machineBody(), newMachineFixture(), newMachineRouter() (+8 more)

## Ambiguous Edges - Review These
- `FanucCncCommunication.dll` → `CncCommunication.dll (32-bit vendor DLL)`  [AMBIGUOUS]
  docs/setup-guide.txt · relation: semantically_similar_to
- `Snapshot Fields (positions, spindle, feed, program, tool, alarms)` → `Machines API (/machines)`  [AMBIGUOUS]
  docs/CncCommunication.dll — API & Integration Guide (Go Client).pdf · relation: conceptually_related_to

## Knowledge Gaps
- **35 isolated node(s):** `7th-sense-server`, `middleware.RequireRole / routes.adminOnly`, `response.SendResponse`, `Test Auth (asTestUser, tests/auth_test.go)`, `user_filters Filter Checklist` (+30 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 194 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **18 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `CNC Data Collection Flow` connect `2. All steps at a glance` to `CNC Data Collection (in-process)`, `Worker (internal/monitor/worker.go)`, `worker`, `MachineTelemetryService`, `Live data — Socket.IO event contract`, `8. Reference`?**
  _High betweenness centrality (0.180) - this node is a cross-community bridge._
- **What connects `7th-sense-server`, `middleware.RequireRole / routes.adminOnly`, `response.SendResponse` to the rest of the system?**
  _35 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `UserShopService` be split into smaller, more focused modules?**
  _Cohesion score 0.12923076923076923 - nodes in this community are weakly interconnected._
- **Why does `2. All steps at a glance` connect `2. All steps at a glance` to `Connection`, `Collector`, `worker`?**
  _High betweenness centrality (0.147) - this node is a cross-community bridge._
- **Should `Machine` be split into smaller, more focused modules?**
  _Cohesion score 0.12121212121212122 - nodes in this community are weakly interconnected._

### Low-confidence Hints
_AMBIGUOUS edges — the extractor was unsure. Verify before acting on these._

- **What is the exact relationship between `FanucCncCommunication.dll` and `CncCommunication.dll (32-bit vendor DLL)`?**
  _Edge tagged AMBIGUOUS (relation: semantically_similar_to) - confidence is low._
- **What is the exact relationship between `Snapshot Fields (positions, spindle, feed, program, tool, alarms)` and `Machines API (/machines)`?**
  _Edge tagged AMBIGUOUS (relation: conceptually_related_to) - confidence is low._