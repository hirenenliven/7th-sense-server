# Graph Report - 7th-sense-server  (2026-11-04)

## Corpus Check
- 278 files · ~210,508 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 4 file(s) not represented in the graph (top: (none) 3, .toml 1)

## Summary
- 2385 nodes · 19097 edges · 89 communities (72 shown, 17 thin omitted)
- Extraction: 95% EXTRACTED · 5% INFERRED · 0% AMBIGUOUS · INFERRED: 911 edges (avg confidence: 0.87)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `515c8ec6`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- user_access_test.go
- VariableConfig
- run
- Machine
- SendResponse()
- pathID()
- net/http.ResponseWriter
- UserCityService
- api_helpers_test.go
- CncMasterSnapshot (34,440 bytes, #pragma
- CNC Data Collection Flow
- Connection
- UserSupervisorService
- Collector
- worker
- MachineCycle
- context.Context
- mode-chart.js
- UserProfileRepository
- PartRepository
- net/http.ServeMux
- CNCConnect(ip, port, dataCb, errCb, machineType) int32
- CNC Data Collection (in-process)
- §10 Diagnostics & Troubleshooting
- Handlers
- machine_mode.go
- BranchService
- UserService
- NewError
- UserLocationService
- .validate
- MachineModelService
- Worker (internal/monitor/worker.go)
- cycleReading
- ControllerService
- validateCatalogFields
- LicenceKeyService
- Client Certificate & Key (client_cert.der / client_key.der)
- PartService
- formatDuration
- CncCommunication.dll API & Integration G
- testing.T
- Execute
- CncCommunication.dll API & Integration G
- Shift
- .Login()
- validText
- .validate
- Machines API (/machines)
- drawChart
- Product
- PartDrawingStorage
- .validate
- README.md
- AuthHandler
- containsMachineID
- net/http.Request
- containsID
- apidocs/main.go
- BrandService
- Test list
- net/http.Handler
- .runConnectionLoop
- UserLineService
- AuthUser
- AuthService
- onDataCallback(ptr uintptr) uintptr
- IdleReasonService
- UserShopService
- DashboardService
- Event-driven data flow (no polling in client)
- NewError
- Fixed values and why they are needed
- UserStateService
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
1. `Execute()` - 182 edges
2. `SendResponse()` - 150 edges
3. `NewError()` - 135 edges
4. `run()` - 96 edges
5. `AuthUser` - 82 edges
6. `pathID()` - 80 edges
7. `currentUser()` - 61 edges
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
- `Step 12 — Call the DLL's CNCConnect` --references--> `Connection`  [INFERRED]
  docs/cnc-data-collection-flow.md → internal/cnc/cnc.go
- `2. All steps at a glance` --references--> `LoadDriver()`  [INFERRED]
  docs/cnc-data-collection-flow.md → internal/cnc/cnc.go

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

## Communities (89 total, 17 thin omitted)

### Community 0 - "user_access_test.go"
Cohesion: 0.32
Nodes (6): LanguageResponse, LicenceKeyDuration, LoginAsRequest, ThemeResponse, licenceKeyDurations(), Envelope

### Community 1 - "VariableConfig"
Cohesion: 0.15
Nodes (10): Step 10 — Get ready to connect, FieldIndex(), TestFieldIndexAxisPositions(), VariableConfig, formatPreviousValue(), formatValue(), newValueChangeLog(), boundVariable (+2 more)

### Community 2 - "run"
Cohesion: 0.07
Nodes (81): run(), Connect(), Querier, NewBranchHandler(), NewBrandHandler(), NewCategoryHandler(), NewControllerHandler(), NewDashboardHandler() (+73 more)

### Community 3 - "Machine"
Cohesion: 0.12
Nodes (15): CreateMachineRequest, MachineFields, UpdateMachineRequest, validCoordinate(), validHealthIDs(), validIPv4(), SeesAllCompanyMachines(), Machine (+7 more)

### Community 4 - "SendResponse()"
Cohesion: 0.08
Nodes (13): BreakHandler, ClienteleHandler, DowntimeReasonHandler, HolidayHandler, IdleReasonHandler, LicenceKeyHandler, MachineHandler, ProductionScheduleHandler (+5 more)

### Community 5 - "pathID()"
Cohesion: 0.10
Nodes (12): validator, PartHandler, decodeAndValidate(), pathID(), targetUserID(), UserCellHandler, UserCityHandler, UserLineHandler (+4 more)

### Community 6 - "net/http.ResponseWriter"
Cohesion: 0.12
Nodes (6): CategoryHandler, DashboardHandler, serveDashboard(), MachineModelHandler, ProductHandler, statusRecorder

### Community 7 - "UserCityService"
Cohesion: 0.16
Nodes (8): CreateUserCityRequest, UpdateUserCityRequest, validateUserCity(), UserCity, UserCityRepository, mapUserCityError(), UserCityService, toUserCityResponseError()

### Community 8 - "api_helpers_test.go"
Cohesion: 0.40
Nodes (4): cncMasterSnapshot, authUserKey, loginAsKey, machineModelRecord

### Community 9 - "CncMasterSnapshot (34,440 bytes, #pragma"
Cohesion: 0.10
Nodes (43): CncMasterSnapshot (34,440 bytes, #pragma pack(8), x86), AbsPos ([5]float64, Inches / Degrees, Work Coordinate System (WCS) positions), ActiveGCode ([16]byte, ASCII, Currently active G-code group), ActiveProgPath ([256]byte, ASCII Path, Storage path of active NC artifact), AlarmMsg ([128]byte, UTF-8 / ASCII, Alarm description message), AlarmNo (int32, Code (0 = OK), Primary alarm or fault number), AxisLoadXPercentage (int16, Percent (%), Servo X axis drive load (AxisLoadX/Y/ZPercentage)), AxisLoadYPercentage (int16, Percent (%), Servo Y axis drive load (AxisLoadX/Y/ZPercentage)) (+35 more)

### Community 10 - "CNC Data Collection Flow"
Cohesion: 0.08
Nodes (25): Supervisor, 1. The big picture, 3. Startup: steps 1 – 4, 6. When data is edited: step 18, 7. Shutdown: steps 19 – 20, 8. Reference, 9. Data Reading Flow (Function Names Only), docs/cnc-communication.md (+17 more)

### Community 11 - "Connection"
Cohesion: 0.09
Nodes (20): dllDriver, slot, 5. Connect and collect, per machine: steps 10 – 17, Step 11 — Connect, a few machines at a time, Step 13 — The DLL sends data, Step 17 — Disconnect and try again, Connection, Reading (+12 more)

### Community 12 - "UserSupervisorService"
Cohesion: 0.15
Nodes (10): CreateUserSupervisorRequest, UpdateUserSupervisorRequest, validateUserSupervisor(), NewUserSupervisorHandler(), UserSupervisor, UserSupervisorRepository, mapUserSupervisorError(), UserSupervisorService (+2 more)

### Community 13 - "Collector"
Cohesion: 0.11
Nodes (22): 2. All steps at a glance, 4. Find the machines: steps 5 – 9, Log lines (in `logs/app.log`, search for the message), Step 12 — Call the DLL's CNCConnect, Step 18 — React to machine and variable changes, Step 20 — Stop the workers and mark machines offline, Step 6 — Record the start and begin following the database, Step 7 — Listen for database changes (+14 more)

### Community 14 - "worker"
Cohesion: 0.12
Nodes (19): Known issue, Step 14 — The data collection loop, Step 15 — Store the reading, Step 16 — Every second: dead connection check and live status, `cycleCheckpointInterval` (5 minutes), FinishCyclesSummary, FinishedCycle, TimelinePeriod (+11 more)

### Community 15 - "MachineCycle"
Cohesion: 0.14
Nodes (10): One running cycle per machine, MachineCyclesResponse, NewMachineCycleHandler(), CycleJob, MachineCycle, MachineCycleRepository, mapMachineCycleError(), MachineCycleService (+2 more)

### Community 16 - "context.Context"
Cohesion: 0.05
Nodes (27): Step 5 — Prepare the database, MachineModesResponse, LiveValue, MachineEvent, MachineLiveStatus, MachineModeStatus, ModeChange, ProductionSchedule (+19 more)

### Community 17 - "mode-chart.js"
Cohesion: 0.21
Nodes (23): chartData(), countByGroup(), drawDashboard(), drawHeader(), drawMachineList(), drawTiles(), drawTotals(), fetchData() (+15 more)

### Community 18 - "UserProfileRepository"
Cohesion: 0.14
Nodes (10): UserProfile, UserProfileRepository, mapProfileError(), NewUserProfileRepository(), NewUserProfileService(), toUserProfileResponseError(), UserProfileService, sameClock() (+2 more)

### Community 19 - "PartRepository"
Cohesion: 0.18
Nodes (5): Part, PartProcess, PartSubProgram, PartRepository, mapPartError()

### Community 20 - "net/http.ServeMux"
Cohesion: 0.15
Nodes (34): Decision, RequireRole(), registerAuthRoutes(), registerBranchRoutes(), registerBrandRoutes(), registerBreakRoutes(), registerCategoryRoutes(), registerClienteleRoutes() (+26 more)

### Community 21 - "CNCConnect(ip, port, dataCb, errCb, machineType) int32"
Cohesion: 0.20
Nodes (11): CNCConnect(ip, port, dataCb, errCb, machineType) int32, CNCConnect param dataCb (uintptr stdcall OnData(const CncMasterSnapshot*)), CNCConnect param ip (const char*, ANSI NUL-terminated), CNCConnect param port (int32, e.g. 4840), CNCdisconnect(machineId int32), LazyDLL binding (syscall.NewLazyDLL + eager mod.Load in init), §4 API Reference, §5 DLL Loading & Binding (Go) (+3 more)

### Community 22 - "CNC Data Collection (in-process)"
Cohesion: 0.15
Nodes (12): CNC Data Collection (in-process), writeLog CNC Log Types and Switches, .env Configuration (internal/config), Health Mode History API (public), CncCommunication.dll (32-bit vendor DLL), machine_mode_log, Mode Decision (machine_mode.go), ECharts 5.5.1 (+4 more)

### Community 23 - "§10 Diagnostics & Troubleshooting"
Cohesion: 0.12
Nodes (25): Check mode: go run sinumerik_client.go check (OPC UA node existence report), Go Client (sinumerik_client.go, callback -> copy -> render), Prerequisite: 32-bit x86 DLL, compile Go with GOARCH=386, Prerequisite: SIEMENS_USER / SIEMENS_PASS env vars before CNCConnect, Run: go run sinumerik_client.go (live dashboard), §3.1 Environment & Prerequisites, §8 Running the Go Client, §10 Diagnostics & Troubleshooting (+17 more)

### Community 24 - "Handlers"
Cohesion: 0.14
Nodes (5): HealthHandler, MachineCycleHandler, UserCountryHandler, registerHealthRoutes(), Handlers

### Community 25 - "machine_mode.go"
Cohesion: 0.15
Nodes (20): engineName(), decideMachineMode(), holdReason(), machineStateFromReading(), numberingEngine(), programStateFromProcStatus(), ncReading(), TestDecideMachineMode() (+12 more)

### Community 26 - "BranchService"
Cohesion: 0.15
Nodes (9): CreateBranchRequest, UpdateBranchRequest, validateBranchFields(), Branch, BranchRepository, mapBranchError(), cityOwner(), BranchService (+1 more)

### Community 27 - "UserService"
Cohesion: 0.06
Nodes (38): WithTransaction(), CustomerDetails, UpdateUserRequest, UserResponse, ParseDayStart(), TestParseDayStart(), Address, AssignableRoles() (+30 more)

### Community 28 - "NewError"
Cohesion: 0.21
Nodes (11): UpdateThemeRequest, normalizeTheme(), TestCreateUserRequestValidate(), TestCreateUserRequestValidateDropsDuplicateMachines(), TestCreateUserRequestValidateNormalizes(), TestCustomerFiltersValidate(), TestUpdateThemeRequestValidate(), TestUpdateUserRequestValidateStatus() (+3 more)

### Community 29 - "UserLocationService"
Cohesion: 0.16
Nodes (8): CreateUserLocationRequest, UpdateUserLocationRequest, validateUserLocation(), UserLocation, UserLocationRepository, mapUserLocationError(), UserLocationService, toUserLocationResponseError()

### Community 30 - ".validate"
Cohesion: 0.24
Nodes (8): validDate(), CreateProductionScheduleRequest, ProductionScheduleFields, ProductionScheduleItemFields, UpdateProductionScheduleRequest, ScheduleDay(), shiftSpan(), jobSpan

### Community 31 - "MachineModelService"
Cohesion: 0.20
Nodes (6): MachineModelRequest, MachineModel, MachineModelRepository, mapMachineModelError(), MachineModelService, toMachineModelResponseError()

### Community 32 - "Worker (internal/monitor/worker.go)"
Cohesion: 0.07
Nodes (22): 7th Sense Server Claude Context, response.SendResponse Format, User CRUD Golden Reference, Collector (internal/monitor/collector.go), CNC Driver (internal/cnc), Known Issue: nil live-value cache panic, machine_events, machine_telementry hypertable (+14 more)

### Community 33 - "cycleReading"
Cohesion: 0.13
Nodes (16): CycleStatusFor(), cycleReadingFrom(), decideCycleStep(), isProgramLoop(), resetReason(), readingAt(), TestAwaitNextCountsFromItsStart(), TestDecideCycleStepCountsOnePartPerStep() (+8 more)

### Community 34 - "ControllerService"
Cohesion: 0.21
Nodes (5): Controller, ControllerRepository, mapControllerError(), ControllerService, toControllerResponseError()

### Community 35 - "validateCatalogFields"
Cohesion: 0.11
Nodes (11): CreateBrandRequest, UpdateBrandRequest, validateCatalogFields(), CreateCategoryRequest, UpdateCategoryRequest, CreateControllerRequest, UpdateControllerRequest, validateControllerFields() (+3 more)

### Community 36 - "LicenceKeyService"
Cohesion: 0.14
Nodes (11): NewLicenceKey, CreateLicenceKeysRequest, GenerateLicenceKeysRequest, GeneratedLicenceKeys, LicenceKey, LicenceKeyRepository, mapLicenceKeyError(), LicenceKeyService (+3 more)

### Community 37 - "Client Certificate & Key (client_cert.der / client_key.der)"
Cohesion: 0.48
Nodes (7): client_cert.der/.pem (X.509 client cert for Siemens OPC UA), client_key.der/.pem (client private key), Siemens certs: client_cert.der + client_key.der in binary directory, Siemens trust: cert must be trusted in NCU Access Administration, Troubleshoot: BadRequestNotAllowed -> Export client_cert, BadRequestNotAllowed (certificate trust rejection), Client Certificate & Key (client_cert.der / client_key.der)

### Community 38 - "PartService"
Cohesion: 0.32
Nodes (6): ValidPartDrawingFile(), NewPartRepository(), PartService, partFromFields(), saveProcesses(), toPartResponseError()

### Community 39 - "formatDuration"
Cohesion: 0.15
Nodes (30): card(), cell(), cycleData(), cycleEnd(), cycleEndReason(), cycleJob(), cycleList(), cycleNumber() (+22 more)

### Community 40 - "CncCommunication.dll API & Integration G"
Cohesion: 0.19
Nodes (20): CncCommunication.dll (C++ core runtime, x86, OPC UA inside), CncCommunication.dll API & Integration Guide (Go Client) - text version, Fanuc profile (engine 1, port 8193, Fwlib32.dll + fwlibe1.dll), Fwlib32.dll (Fanuc FOCAS2), fwlibe1.dll (Fanuc FOCAS2), libcrypto-3.dll (OpenSSL), libssl-3.dll (OpenSSL), Engine 1: Fanuc (FOCAS2 Ethernet, fwlib32) (+12 more)

### Community 41 - "testing.T"
Cohesion: 0.09
Nodes (71): TestSnapshotModeFieldOffsets(), TestSnapshotSize(), CreateUserRequest, decodeData(), expectStatus(), newCatalogRouter(), runCatalogCRUD(), sendRequest() (+63 more)

### Community 42 - "Execute"
Cohesion: 0.05
Nodes (30): Execute(), UserListFilter, Category, Role, User, UserLine, AuthRepository, mapAuthError() (+22 more)

### Community 43 - "CncCommunication.dll API & Integration G"
Cohesion: 0.13
Nodes (18): CNCConnect param machineType (int32 engine id, see §4.2), Engine 0: Simulator (local loopback mock engine), Engine 2: Siemens (Sinumerik OPC UA, active profile), Engine 3: Syntec (Open API / TCP), Engine 4: Mitsubishi (driver not specified), Machine identifier mapping (engine id enum 0-4), §9 Controller-Specific Profiles, Siemens default ports: 4840 (native NCU) / 4841 (custom) (+10 more)

### Community 44 - "Shift"
Cohesion: 0.13
Nodes (18): uniquePositiveIDs(), clockMinutes(), CreateShiftRequest, ShiftFields, UpdateShiftRequest, MinutesBetween(), validWorkingDays(), Shift (+10 more)

### Community 45 - ".Login()"
Cohesion: 0.15
Nodes (9): LoginRequest, hashToken(), newLoginAsKey(), newSessionToken(), truncate(), newLoginThrottle(), loginFailures, loginThrottle (+1 more)

### Community 46 - "validText"
Cohesion: 0.19
Nodes (8): PartProcessFields, PartSubProgramFields, CreatePartRequest, PartFields, UpdatePartRequest, UploadPartDrawingRequest, validText(), PartDrawing

### Community 47 - ".validate"
Cohesion: 0.47
Nodes (3): CreateHolidayRequest, HolidayFields, UpdateHolidayRequest

### Community 48 - "Machines API (/machines)"
Cohesion: 0.16
Nodes (18): Clienteles API, Company Scope (companyOwner / ownCompanyCustomer), Current User API (/me theme, language), Customer Display Names (customerName / companyName), Holidays API, Licence Keys API, Machine Catalog Brand-Controller-Model, Machines API (/machines) (+10 more)

### Community 49 - "drawChart"
Cohesion: 0.24
Nodes (11): chartRows(), currentSpan(), currentWindow(), drawChart(), drawZoomLevels(), markZoomLevel(), renderBar(), renderCycleBar() (+3 more)

### Community 50 - "Product"
Cohesion: 0.23
Nodes (6): Product, ProductRepository, mapProductError(), ProductService, productImageOrDefault(), toProductResponseError()

### Community 52 - ".validate"
Cohesion: 0.38
Nodes (4): ClienteleFields, CreateClienteleRequest, UpdateClienteleRequest, validClienteleEmail()

### Community 53 - "README.md"
Cohesion: 0.04
Nodes (38): 0001. Go standard library net/http, no web framework, Alternatives considered, Consequences, Context, Decision, 0002. Layered architecture: handler, DTO, service, repository, Alternatives considered, Consequences (+30 more)

### Community 54 - "AuthHandler"
Cohesion: 0.22
Nodes (8): Decision, AuthHandler, CurrentUser(), isLoginAsKey(), LoginAsKey(), RouteLoginAs(), SessionToken(), UserID()

### Community 56 - "net/http.Request"
Cohesion: 0.12
Nodes (6): BranchHandler, BrandHandler, ControllerHandler, currentUserID(), queryBool(), UserHandler

### Community 58 - "apidocs/main.go"
Cohesion: 0.06
Nodes (44): request, main(), run(), main(), Config, MachineSelection, GenerateCertificates(), getEnv() (+36 more)

### Community 59 - "BrandService"
Cohesion: 0.25
Nodes (5): Brand, BrandRepository, mapBrandError(), BrandService, toBrandResponseError()

### Community 61 - "Test list"
Cohesion: 0.05
Nodes (38): 1. The rows are consistent (SQL, no machine needed), 2. The rows agree with the mode log (SQL), 3. The raw signals around one cycle (telemetry and logs), 4. At the machine (floor test), Before the tests (once), Check query, Completed, Floor tests on the HUMA panel (Fanuc) (+30 more)

### Community 63 - "net/http.Handler"
Cohesion: 0.12
Nodes (26): Handler(), NewBreakHandler(), NewClienteleHandler(), NewHolidayHandler(), NewLicenceKeyHandler(), NewPartHandler(), NewProductionScheduleHandler(), NewShiftHandler() (+18 more)

### Community 64 - ".runConnectionLoop"
Cohesion: 0.19
Nodes (7): connectFanucCounter(), main(), run(), secondsOf(), sleepUnlessCancelled(), backoff, fanucCounter

### Community 66 - "UserLineService"
Cohesion: 0.09
Nodes (16): CreateUserCellRequest, UpdateUserCellRequest, validateUserCell(), CreateUserLineRequest, UpdateUserLineRequest, validateUserLine(), NewUserCellHandler(), NewUserLineHandler() (+8 more)

### Community 68 - "AuthUser"
Cohesion: 0.07
Nodes (29): LoginAsResponse, BreakFields, CreateBreakRequest, UpdateBreakRequest, CreateDowntimeReasonRequest, DowntimeReasonFields, UpdateDowntimeReasonRequest, NewDowntimeReasonHandler() (+21 more)

### Community 69 - "AuthService"
Cohesion: 0.25
Nodes (13): NewAuthHandler(), Authenticate(), WithAuthUser(), NewAuthRepository(), AuthService, NewAuthService(), authRequest(), login() (+5 more)

### Community 71 - "onDataCallback(ptr uintptr) uintptr"
Cohesion: 0.36
Nodes (9): Package-level callback singletons (syscall.NewCallback), CNCConnect param errCb (uintptr stdcall OnErr(int32 code, const char* msg)), Mutex-protected copy: dataMutex sync.RWMutex -> lastRecord/hasRecord, onDataCallback(ptr uintptr) uintptr, onErrorCallback(code, messagePtr uintptr) uintptr, §7 Asynchronous Callback Architecture, onDataCallback, onErrorCallback (+1 more)

### Community 72 - "IdleReasonService"
Cohesion: 0.16
Nodes (10): CreateIdleReasonRequest, IdleReasonFields, UpdateIdleReasonRequest, NewIdleReasonHandler(), IdleReason, IdleReasonRepository, mapIdleReasonError(), IdleReasonService (+2 more)

### Community 74 - "UserShopService"
Cohesion: 0.13
Nodes (12): CreateUserShopRequest, UpdateUserShopRequest, validateUserShop(), NewUserShopHandler(), UserShop, UserShopRepository, mapUserShopError(), NewUserShopRepository() (+4 more)

### Community 77 - "DashboardService"
Cohesion: 0.08
Nodes (34): Decision, CycleCountRow, DashboardRequest, MachineTimeRow, ModeTotalsResponse, OEEResponse, OEERow, NewFinishCyclesSummary() (+26 more)

### Community 80 - "Event-driven data flow (no polling in client)"
Cohesion: 0.60
Nodes (5): Flow 1: load DLL, register two stdcall callbacks (syscall.NewCallback), Flow 2: CNCConnect starts DLL background worker thread, Flow 3: worker polls controller, assembles snapshot, invokes data callback, Flow 4: callback copies payload, signals renderer via buffered channel, Flow 5: main goroutine renders live ANSI dashboard

### Community 81 - "NewError"
Cohesion: 0.11
Nodes (16): CreateAddressRequest, UpdateAddressRequest, UpdateLanguageRequest, UpdateLicenceKeyRequest, validDuration(), ValidLicenceKey(), isTimeZone(), isWebsite() (+8 more)

### Community 83 - "Fixed values and why they are needed"
Cohesion: 0.29
Nodes (7): `CNC_CYCLE_MAX_HOURS` (24 hours), `CNC_STALE_SEC` (30 s) and `CNC_QUIET_SEC` (300 s), `cycleCounterGrace` and `cycleStartGrace` (5 seconds each), `emergencyResetGrace` (2 seconds), Fixed values and why they are needed, `loopMinimumDropFraction` (0.5), `maxPartCountStep` (1)

### Community 84 - "UserStateService"
Cohesion: 0.09
Nodes (13): CreateUserCountryRequest, CreateUserStateRequest, validateUserState(), UserCountry, UserState, UserCountryRepository, mapUserCountryError(), UserStateRepository (+5 more)

### Community 86 - "auth_tokens SHA-256 Hashed Sessions (mig"
Cohesion: 0.67
Nodes (3): auth_tokens SHA-256 Hashed Sessions (migration 014), Per-Request Session and User Re-check, ts_session Cookie (HttpOnly, SameSite=Strict)

## Ambiguous Edges - Review These
- `CncCommunication.dll (32-bit vendor DLL)` → `FanucCncCommunication.dll`  [AMBIGUOUS]
  docs/setup-guide.txt · relation: semantically_similar_to
- `Snapshot Fields (positions, spindle, feed, program, tool, alarms)` → `Machines API (/machines)`  [AMBIGUOUS]
  docs/CncCommunication.dll — API & Integration Guide (Go Client).pdf · relation: conceptually_related_to

## Knowledge Gaps
- **35 isolated node(s):** `7th-sense-server`, `docs/cnc-communication.md`, `middleware.RequireRole / routes.adminOnly`, `response.SendResponse`, `Test Auth (asTestUser, tests/auth_test.go)` (+30 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 160 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **17 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `CNC Data Collection Flow` connect `CNC Data Collection Flow` to `Worker (internal/monitor/worker.go)`, `Connection`, `Collector`, `CNC Data Collection (in-process)`?**
  _High betweenness centrality (0.202) - this node is a cross-community bridge._
- **What connects `7th-sense-server`, `docs/cnc-communication.md`, `middleware.RequireRole / routes.adminOnly` to the rest of the system?**
  _35 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `run` be split into smaller, more focused modules?**
  _Cohesion score 0.06504065040650407 - nodes in this community are weakly interconnected._
- **Why does `2. All steps at a glance` connect `Collector` to `.runConnectionLoop`, `VariableConfig`, `CNC Data Collection Flow`, `Connection`, `worker`?**
  _High betweenness centrality (0.174) - this node is a cross-community bridge._
- **Should `Machine` be split into smaller, more focused modules?**
  _Cohesion score 0.11954022988505747 - nodes in this community are weakly interconnected._

### Low-confidence Hints
_AMBIGUOUS edges — the extractor was unsure. Verify before acting on these._

- **What is the exact relationship between `CncCommunication.dll (32-bit vendor DLL)` and `FanucCncCommunication.dll`?**
  _Edge tagged AMBIGUOUS (relation: semantically_similar_to) - confidence is low._
- **What is the exact relationship between `Snapshot Fields (positions, spindle, feed, program, tool, alarms)` and `Machines API (/machines)`?**
  _Edge tagged AMBIGUOUS (relation: conceptually_related_to) - confidence is low._