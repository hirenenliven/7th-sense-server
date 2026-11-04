# Graph Report - 7th-sense-server  (2026-11-04)

## Corpus Check
- 276 files · ~206,865 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 4 file(s) not represented in the graph (top: (none) 3, .toml 1)

## Summary
- 2360 nodes · 18936 edges · 101 communities (85 shown, 16 thin omitted)
- Extraction: 95% EXTRACTED · 5% INFERRED · 0% AMBIGUOUS · INFERRED: 896 edges (avg confidence: 0.87)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `3c85ad6a`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- user_access_test.go
- VariableConfig
- UserCellService
- MachineTelemetryService
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
- AuthUser
- context.Context
- testing.T
- time.Time
- mode-chart.js
- net/http.ServeMux
- CNCConnect(ip, port, dataCb, errCb, machineType) int32
- BreakRepository
- §10 Diagnostics & Troubleshooting
- executeAffected
- machine_mode.go
- BranchService
- UserService
- NewError
- UserLocationService
- run
- MachineModelService
- Worker (internal/monitor/worker.go)
- cycleReading
- validText
- validateCatalogFields
- LicenceKeyService
- Client Certificate & Key (client_cert.der / client_key.der)
- PartService
- formatDuration
- CncCommunication.dll API & Integration G
- ControllerService
- Category
- CncCommunication.dll API & Integration G
- AuthUser
- .Login()
- MachineCycleRepository
- CNC Data Collection (in-process)
- Machines API (/machines)
- .Update
- Product
- drawChart
- ClienteleService
- README.md
- AuthHandler
- containsMachineID
- net/http.ResponseWriter
- Role
- apidocs/main.go
- BrandService
- Address
- Test list
- net/http.Request
- MachineService
- .runConnectionLoop
- .validate
- Execute
- UserProfileRepository
- .validate
- net/http.Handler
- .validate
- onDataCallback(ptr uintptr) uintptr
- HolidayService
- .Validate
- UserShopService
- containsID
- UserCountryService
- validateUserLocation
- add_superadmin_raw.go
- CompanyProfile
- Event-driven data flow (no polling in client)
- .validate
- LocationRepository
- Fixed values and why they are needed
- .Create
- newUserServices
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
1. `Execute()` - 181 edges
2. `SendResponse()` - 148 edges
3. `NewError()` - 134 edges
4. `run()` - 96 edges
5. `pathID()` - 80 edges
6. `AuthUser` - 79 edges
7. `Querier` - 60 edges
8. `currentUser()` - 60 edges
9. `decodeAndValidate()` - 58 edges
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

## Communities (101 total, 16 thin omitted)

### Community 0 - "user_access_test.go"
Cohesion: 0.31
Nodes (5): LanguageResponse, LicenceKeyDuration, ThemeResponse, licenceKeyDurations(), Envelope

### Community 1 - "VariableConfig"
Cohesion: 0.12
Nodes (13): 5. Connect and collect, per machine: steps 10 – 17, Step 10 — Get ready to connect, Step 11 — Connect, a few machines at a time, Step 14 — The data collection loop, FieldIndex(), TestFieldIndexAxisPositions(), VariableConfig, formatPreviousValue() (+5 more)

### Community 2 - "UserCellService"
Cohesion: 0.07
Nodes (19): CreateUserCellRequest, UpdateUserCellRequest, validateUserCell(), CreateUserLineRequest, UpdateUserLineRequest, validateUserLine(), NewUserCellHandler(), NewUserLineHandler() (+11 more)

### Community 3 - "MachineTelemetryService"
Cohesion: 0.11
Nodes (13): Step 5 — Prepare the database, MachineModesResponse, NewHealthHandler(), MachineEvent, MachineLiveStatus, MachineModeStatus, ModeChange, MachineTelemetryRepository (+5 more)

### Community 4 - "SendResponse()"
Cohesion: 0.09
Nodes (10): BreakHandler, ClienteleHandler, DowntimeReasonHandler, HolidayHandler, IdleReasonHandler, LicenceKeyHandler, PartHandler, ProductionScheduleHandler (+2 more)

### Community 5 - "pathID()"
Cohesion: 0.10
Nodes (12): MachineCycleHandler, pathID(), targetUserID(), UserCellHandler, UserCityHandler, UserCountryHandler, UserLineHandler, UserLocationHandler (+4 more)

### Community 6 - "net/http.ResponseWriter"
Cohesion: 0.10
Nodes (7): BrandHandler, CategoryHandler, ControllerHandler, HealthHandler, ProductHandler, registerHealthRoutes(), statusRecorder

### Community 7 - "UserCityService"
Cohesion: 0.16
Nodes (8): CreateUserCityRequest, UpdateUserCityRequest, validateUserCity(), UserCity, UserCityRepository, mapUserCityError(), UserCityService, toUserCityResponseError()

### Community 8 - "api_helpers_test.go"
Cohesion: 0.44
Nodes (3): cncMasterSnapshot, authUserKey, machineModelRecord

### Community 9 - "CncMasterSnapshot (34,440 bytes, #pragma"
Cohesion: 0.10
Nodes (43): CncMasterSnapshot (34,440 bytes, #pragma pack(8), x86), AbsPos ([5]float64, Inches / Degrees, Work Coordinate System (WCS) positions), ActiveGCode ([16]byte, ASCII, Currently active G-code group), ActiveProgPath ([256]byte, ASCII Path, Storage path of active NC artifact), AlarmMsg ([128]byte, UTF-8 / ASCII, Alarm description message), AlarmNo (int32, Code (0 = OK), Primary alarm or fault number), AxisLoadXPercentage (int16, Percent (%), Servo X axis drive load (AxisLoadX/Y/ZPercentage)), AxisLoadYPercentage (int16, Percent (%), Servo Y axis drive load (AxisLoadX/Y/ZPercentage)) (+35 more)

### Community 10 - "CNC Data Collection Flow"
Cohesion: 0.08
Nodes (25): Supervisor, 1. The big picture, 3. Startup: steps 1 – 4, 6. When data is edited: step 18, 7. Shutdown: steps 19 – 20, 8. Reference, 9. Data Reading Flow (Function Names Only), docs/cnc-communication.md (+17 more)

### Community 11 - "Connection"
Cohesion: 0.10
Nodes (17): dllDriver, slot, Step 13 — The DLL sends data, Step 17 — Disconnect and try again, Connection, Reading, Stats, Status (+9 more)

### Community 12 - "UserSupervisorService"
Cohesion: 0.15
Nodes (10): CreateUserSupervisorRequest, UpdateUserSupervisorRequest, validateUserSupervisor(), NewUserSupervisorHandler(), UserSupervisor, UserSupervisorRepository, mapUserSupervisorError(), UserSupervisorService (+2 more)

### Community 13 - "Collector"
Cohesion: 0.11
Nodes (22): 2. All steps at a glance, 4. Find the machines: steps 5 – 9, Log lines (in `logs/app.log`, search for the message), Step 12 — Call the DLL's CNCConnect, Step 18 — React to machine and variable changes, Step 20 — Stop the workers and mark machines offline, Step 6 — Record the start and begin following the database, Step 7 — Listen for database changes (+14 more)

### Community 14 - "worker"
Cohesion: 0.09
Nodes (21): Known issue, Step 15 — Store the reading, Step 16 — Every second: dead connection check and live status, FinishedCycle, Value, FinishCyclesResponse, NewFinishedCycle(), CycleStatusFor() (+13 more)

### Community 15 - "AuthUser"
Cohesion: 0.13
Nodes (14): CreateMachineRequest, MachineFields, UpdateMachineRequest, validCoordinate(), validHealthIDs(), validIPv4(), SeesAllCompanyMachines(), Machine (+6 more)

### Community 16 - "context.Context"
Cohesion: 0.16
Nodes (5): ProductionScheduleItem, mapIdleReasonError(), ProductionScheduleRepository, mapProductionScheduleError(), hashToken()

### Community 17 - "testing.T"
Cohesion: 0.09
Nodes (71): TestSnapshotModeFieldOffsets(), TestSnapshotSize(), CreateUserRequest, decodeData(), expectStatus(), newCatalogRouter(), runCatalogCRUD(), sendRequest() (+63 more)

### Community 18 - "time.Time"
Cohesion: 0.05
Nodes (45): 0006. Page-wise dashboard endpoints returning only chart fields, Alternatives considered, Consequences, Context, Decision, FinishCyclesSummary, TimelinePeriod, CycleCountRow (+37 more)

### Community 19 - "mode-chart.js"
Cohesion: 0.21
Nodes (23): chartData(), countByGroup(), drawDashboard(), drawHeader(), drawMachineList(), drawTiles(), drawTotals(), fetchData() (+15 more)

### Community 20 - "net/http.ServeMux"
Cohesion: 0.15
Nodes (33): Decision, RequireRole(), registerBranchRoutes(), registerBrandRoutes(), registerBreakRoutes(), registerCategoryRoutes(), registerClienteleRoutes(), registerControllerRoutes() (+25 more)

### Community 21 - "CNCConnect(ip, port, dataCb, errCb, machineType) int32"
Cohesion: 0.20
Nodes (11): CNCConnect(ip, port, dataCb, errCb, machineType) int32, CNCConnect param dataCb (uintptr stdcall OnData(const CncMasterSnapshot*)), CNCConnect param ip (const char*, ANSI NUL-terminated), CNCConnect param port (int32, e.g. 4840), CNCdisconnect(machineId int32), LazyDLL binding (syscall.NewLazyDLL + eager mod.Load in init), §4 API Reference, §5 DLL Loading & Binding (Go) (+3 more)

### Community 22 - "BreakRepository"
Cohesion: 0.20
Nodes (3): ShiftBreak, BreakRepository, mapBreakError()

### Community 23 - "§10 Diagnostics & Troubleshooting"
Cohesion: 0.12
Nodes (25): Check mode: go run sinumerik_client.go check (OPC UA node existence report), Go Client (sinumerik_client.go, callback -> copy -> render), Prerequisite: 32-bit x86 DLL, compile Go with GOARCH=386, Prerequisite: SIEMENS_USER / SIEMENS_PASS env vars before CNCConnect, Run: go run sinumerik_client.go (live dashboard), §3.1 Environment & Prerequisites, §8 Running the Go Client, §10 Diagnostics & Troubleshooting (+17 more)

### Community 24 - "executeAffected"
Cohesion: 0.16
Nodes (5): PartProcess, PartSubProgram, executeAffected(), PartRepository, mapPartError()

### Community 25 - "machine_mode.go"
Cohesion: 0.15
Nodes (19): engineName(), decideMachineMode(), holdReason(), machineStateFromReading(), numberingEngine(), programStateFromProcStatus(), ncReading(), TestDecideMachineMode() (+11 more)

### Community 26 - "BranchService"
Cohesion: 0.15
Nodes (9): CreateBranchRequest, UpdateBranchRequest, validateBranchFields(), Branch, BranchRepository, mapBranchError(), cityOwner(), BranchService (+1 more)

### Community 27 - "UserService"
Cohesion: 0.16
Nodes (9): UserListFilter, UserResponse, AssignableRoles(), CanAssignRole(), canSee(), UserService, toUserResponseError(), userAccess (+1 more)

### Community 28 - "NewError"
Cohesion: 0.11
Nodes (21): CreateAddressRequest, UpdateAddressRequest, UpdateLanguageRequest, UpdateThemeRequest, isTimeZone(), isWebsite(), normalizeFilters(), normalizeTheme() (+13 more)

### Community 29 - "UserLocationService"
Cohesion: 0.20
Nodes (4): UserLocation, UserLocationRepository, mapUserLocationError(), toUserLocationResponseError()

### Community 30 - "run"
Cohesion: 0.07
Nodes (73): run(), Querier, NewBranchHandler(), NewBrandHandler(), NewCategoryHandler(), NewControllerHandler(), NewDashboardHandler(), NewMachineModelHandler() (+65 more)

### Community 31 - "MachineModelService"
Cohesion: 0.20
Nodes (6): MachineModelRequest, MachineModel, MachineModelRepository, mapMachineModelError(), MachineModelService, toMachineModelResponseError()

### Community 32 - "Worker (internal/monitor/worker.go)"
Cohesion: 0.07
Nodes (22): 7th Sense Server Claude Context, response.SendResponse Format, User CRUD Golden Reference, Collector (internal/monitor/collector.go), CNC Driver (internal/cnc), Known Issue: nil live-value cache panic, machine_events, machine_telementry hypertable (+14 more)

### Community 33 - "cycleReading"
Cohesion: 0.17
Nodes (13): cycleReadingFrom(), decideCycleStep(), isProgramLoop(), resetReason(), readingAt(), TestDecideCycleStepSubProgramCallIsNotALoop(), TestFollowEmergency(), TestIsProgramLoop() (+5 more)

### Community 34 - "validText"
Cohesion: 0.14
Nodes (7): PartProcessFields, PartSubProgramFields, UploadPartDrawingRequest, uniquePositiveIDs(), validText(), validWorkingDays(), PartDrawing

### Community 35 - "validateCatalogFields"
Cohesion: 0.11
Nodes (11): CreateBrandRequest, UpdateBrandRequest, validateCatalogFields(), CreateCategoryRequest, UpdateCategoryRequest, CreateControllerRequest, UpdateControllerRequest, validateControllerFields() (+3 more)

### Community 36 - "LicenceKeyService"
Cohesion: 0.15
Nodes (11): NewLicenceKey, CreateLicenceKeysRequest, GenerateLicenceKeysRequest, GeneratedLicenceKeys, LicenceKey, LicenceKeyRepository, mapLicenceKeyError(), LicenceKeyService (+3 more)

### Community 37 - "Client Certificate & Key (client_cert.der / client_key.der)"
Cohesion: 0.48
Nodes (7): client_cert.der/.pem (X.509 client cert for Siemens OPC UA), client_key.der/.pem (client private key), Siemens certs: client_cert.der + client_key.der in binary directory, Siemens trust: cert must be trusted in NCU Access Administration, Troubleshoot: BadRequestNotAllowed -> Export client_cert, BadRequestNotAllowed (certificate trust rejection), Client Certificate & Key (client_cert.der / client_key.der)

### Community 38 - "PartService"
Cohesion: 0.24
Nodes (10): CreatePartRequest, PartFields, UpdatePartRequest, ValidPartDrawingFile(), Part, NewPartRepository(), PartService, partFromFields() (+2 more)

### Community 39 - "formatDuration"
Cohesion: 0.15
Nodes (30): card(), cell(), cycleData(), cycleEnd(), cycleEndReason(), cycleJob(), cycleList(), cycleNumber() (+22 more)

### Community 40 - "CncCommunication.dll API & Integration G"
Cohesion: 0.19
Nodes (20): CncCommunication.dll (C++ core runtime, x86, OPC UA inside), CncCommunication.dll API & Integration Guide (Go Client) - text version, Fanuc profile (engine 1, port 8193, Fwlib32.dll + fwlibe1.dll), Fwlib32.dll (Fanuc FOCAS2), fwlibe1.dll (Fanuc FOCAS2), libcrypto-3.dll (OpenSSL), libssl-3.dll (OpenSSL), Engine 1: Fanuc (FOCAS2 Ethernet, fwlib32) (+12 more)

### Community 41 - "ControllerService"
Cohesion: 0.21
Nodes (5): Controller, ControllerRepository, mapControllerError(), ControllerService, toControllerResponseError()

### Community 42 - "Category"
Cohesion: 0.25
Nodes (5): Category, CategoryRepository, mapCategoryError(), CategoryService, toCategoryResponseError()

### Community 43 - "CncCommunication.dll API & Integration G"
Cohesion: 0.13
Nodes (18): CNCConnect param machineType (int32 engine id, see §4.2), Engine 0: Simulator (local loopback mock engine), Engine 2: Siemens (Sinumerik OPC UA, active profile), Engine 3: Syntec (Open API / TCP), Engine 4: Mitsubishi (driver not specified), Machine identifier mapping (engine id enum 0-4), §9 Controller-Specific Profiles, Siemens default ports: 4840 (native NCU) / 4841 (custom) (+10 more)

### Community 44 - "AuthUser"
Cohesion: 0.06
Nodes (43): WithTransaction(), clockMinutes(), CreateShiftRequest, ShiftFields, UpdateShiftRequest, MinutesBetween(), NewDowntimeReasonHandler(), NewIdleReasonHandler() (+35 more)

### Community 45 - ".Login()"
Cohesion: 0.18
Nodes (7): LoginRequest, newSessionToken(), truncate(), newLoginThrottle(), loginFailures, loginThrottle, Session

### Community 46 - "MachineCycleRepository"
Cohesion: 0.23
Nodes (3): CycleJob, MachineCycleRepository, mapMachineCycleError()

### Community 47 - "CNC Data Collection (in-process)"
Cohesion: 0.15
Nodes (12): CNC Data Collection (in-process), writeLog CNC Log Types and Switches, .env Configuration (internal/config), Health Mode History API (public), CncCommunication.dll (32-bit vendor DLL), machine_mode_log, Mode Decision (machine_mode.go), ECharts 5.5.1 (+4 more)

### Community 48 - "Machines API (/machines)"
Cohesion: 0.16
Nodes (18): Clienteles API, Company Scope (companyOwner / ownCompanyCustomer), Current User API (/me theme, language), Customer Display Names (customerName / companyName), Holidays API, Licence Keys API, Machine Catalog Brand-Controller-Model, Machines API (/machines) (+10 more)

### Community 49 - ".Update"
Cohesion: 0.35
Nodes (8): CustomerDetails, UpdateUserRequest, ParseDayStart(), TestParseDayStart(), companyProfileFrom(), setCustomerAddress(), setCustomerSettings(), themeOrDefault()

### Community 50 - "Product"
Cohesion: 0.23
Nodes (6): Product, ProductRepository, mapProductError(), ProductService, productImageOrDefault(), toProductResponseError()

### Community 51 - "drawChart"
Cohesion: 0.24
Nodes (11): chartRows(), currentSpan(), currentWindow(), drawChart(), drawZoomLevels(), markZoomLevel(), renderBar(), renderCycleBar() (+3 more)

### Community 52 - "ClienteleService"
Cohesion: 0.13
Nodes (10): ClienteleFields, CreateClienteleRequest, UpdateClienteleRequest, validClienteleEmail(), Clientele, ClienteleRepository, mapClienteleError(), clienteleFromFields() (+2 more)

### Community 53 - "README.md"
Cohesion: 0.05
Nodes (30): 0001. Go standard library net/http, no web framework, Alternatives considered, Consequences, Context, Decision, 0002. Layered architecture: handler, DTO, service, repository, Alternatives considered, Consequences (+22 more)

### Community 54 - "AuthHandler"
Cohesion: 0.38
Nodes (3): AuthHandler, CurrentUser(), UserID()

### Community 56 - "net/http.ResponseWriter"
Cohesion: 0.16
Nodes (5): validator, BranchHandler, currentUserID(), decodeAndValidate(), UserHandler

### Community 57 - "Role"
Cohesion: 0.33
Nodes (3): Role, mapRoleError(), toRoleResponseError()

### Community 58 - "apidocs/main.go"
Cohesion: 0.06
Nodes (45): request, main(), run(), main(), Config, MachineSelection, GenerateCertificates(), getEnv() (+37 more)

### Community 59 - "BrandService"
Cohesion: 0.25
Nodes (5): Brand, BrandRepository, mapBrandError(), BrandService, toBrandResponseError()

### Community 60 - "Address"
Cohesion: 0.22
Nodes (7): Address, AddressRepository, mapAddressError(), NewAddressRepository(), toAddressResponseError(), primaryAddress(), AddressService

### Community 61 - "Test list"
Cohesion: 0.08
Nodes (26): Before the tests (once), Check query, Floor tests on the HUMA panel (Fanuc), HUMA buttons used, Result sheet, TEST-01 · Start a cycle, TEST-02 · Normal end (M30), TEST-03 · Three parts one after another (+18 more)

### Community 62 - "net/http.Request"
Cohesion: 0.12
Nodes (8): dashboardQuery(), DashboardHandler, serveDashboard(), MachineHandler, MachineModelHandler, listUserID(), queryID(), ShiftHandler

### Community 63 - "MachineService"
Cohesion: 0.09
Nodes (32): NewBreakHandler(), NewClienteleHandler(), NewHolidayHandler(), NewLicenceKeyHandler(), NewMachineHandler(), NewPartHandler(), NewProductionScheduleHandler(), NewShiftHandler() (+24 more)

### Community 64 - ".runConnectionLoop"
Cohesion: 0.21
Nodes (6): connectFanucCounter(), main(), run(), sleepUnlessCancelled(), backoff, fanucCounter

### Community 65 - ".validate"
Cohesion: 0.47
Nodes (3): BreakFields, CreateBreakRequest, UpdateBreakRequest

### Community 66 - "Execute"
Cohesion: 0.10
Nodes (5): Execute(), User, mapAuthError(), mapMachineTelemetryError(), mapUserError()

### Community 67 - "UserProfileRepository"
Cohesion: 0.19
Nodes (5): UserProfile, UserProfileRepository, mapProfileError(), toUserProfileResponseError(), UserProfileService

### Community 68 - ".validate"
Cohesion: 0.47
Nodes (3): CreateDowntimeReasonRequest, DowntimeReasonFields, UpdateDowntimeReasonRequest

### Community 69 - "net/http.Handler"
Cohesion: 0.14
Nodes (19): Handler(), NewAuthHandler(), Authenticate(), WithAuthUser(), CrossOrigin(), Logger(), Chain(), Recoverer() (+11 more)

### Community 70 - ".validate"
Cohesion: 0.47
Nodes (3): CreateIdleReasonRequest, IdleReasonFields, UpdateIdleReasonRequest

### Community 71 - "onDataCallback(ptr uintptr) uintptr"
Cohesion: 0.36
Nodes (9): Package-level callback singletons (syscall.NewCallback), CNCConnect param errCb (uintptr stdcall OnErr(int32 code, const char* msg)), Mutex-protected copy: dataMutex sync.RWMutex -> lastRecord/hasRecord, onDataCallback(ptr uintptr) uintptr, onErrorCallback(code, messagePtr uintptr) uintptr, §7 Asynchronous Callback Architecture, onDataCallback, onErrorCallback (+1 more)

### Community 72 - "HolidayService"
Cohesion: 0.16
Nodes (10): CreateHolidayRequest, HolidayFields, UpdateHolidayRequest, Holiday, HolidayRepository, mapHolidayError(), HolidayService, holidayFromFields() (+2 more)

### Community 73 - ".Validate"
Cohesion: 0.40
Nodes (3): UpdateLicenceKeyRequest, validDuration(), ValidLicenceKey()

### Community 74 - "UserShopService"
Cohesion: 0.15
Nodes (9): CreateUserShopRequest, UpdateUserShopRequest, validateUserShop(), NewUserShopHandler(), UserShop, UserShopRepository, mapUserShopError(), UserShopService (+1 more)

### Community 76 - "UserCountryService"
Cohesion: 0.20
Nodes (4): CreateUserCountryRequest, UserCountry, mapUserCountryError(), toUserCountryResponseError()

### Community 77 - "validateUserLocation"
Cohesion: 0.40
Nodes (3): CreateUserLocationRequest, UpdateUserLocationRequest, validateUserLocation()

### Community 78 - "add_superadmin_raw.go"
Cohesion: 0.36
Nodes (10): checkNotTaken(), createSuperadmin(), ensureProfile(), execute(), main(), parseFlags(), run(), superadminRoleID() (+2 more)

### Community 79 - "CompanyProfile"
Cohesion: 0.28
Nodes (5): CompanyProfile, CompanyProfileRepository, mapCompanyProfileError(), toCompanyProfileResponseError(), CompanyProfileService

### Community 80 - "Event-driven data flow (no polling in client)"
Cohesion: 0.60
Nodes (5): Flow 1: load DLL, register two stdcall callbacks (syscall.NewCallback), Flow 2: CNCConnect starts DLL background worker thread, Flow 3: worker polls controller, assembles snapshot, invokes data callback, Flow 4: callback copies payload, signals renderer via buffered channel, Flow 5: main goroutine renders live ANSI dashboard

### Community 81 - ".validate"
Cohesion: 0.18
Nodes (10): validDate(), CreateProductionScheduleRequest, ProductionScheduleFields, ProductionScheduleItemFields, UpdateProductionScheduleRequest, ScheduleDay(), asInvalidSelection(), checkJobsOverlap() (+2 more)

### Community 82 - "LocationRepository"
Cohesion: 0.50
Nodes (3): mapLocationError(), NewLocationRepository(), LocationRepository

### Community 83 - "Fixed values and why they are needed"
Cohesion: 0.09
Nodes (21): 1. The rows are consistent (SQL, no machine needed), 2. The rows agree with the mode log (SQL), 3. The raw signals around one cycle (telemetry and logs), 4. At the machine (floor test), `CNC_CYCLE_MAX_HOURS` (24 hours), `CNC_STALE_SEC` (30 s) and `CNC_QUIET_SEC` (300 s), Completed, `cycleCheckpointInterval` (5 minutes) (+13 more)

### Community 84 - ".Create"
Cohesion: 0.18
Nodes (5): CreateUserStateRequest, validateUserState(), UserState, mapUserStateError(), toUserStateResponseError()

### Community 85 - "newUserServices"
Cohesion: 0.16
Nodes (14): NewCompanyProfileRepository(), UserMachineAccessRepository, NewUserMachineAccessRepository(), NewUserProfileRepository(), NewAddressService(), NewCompanyProfileService(), NewUserProfileService(), newUserServices() (+6 more)

### Community 86 - "auth_tokens SHA-256 Hashed Sessions (mig"
Cohesion: 0.67
Nodes (3): auth_tokens SHA-256 Hashed Sessions (migration 014), Per-Request Session and User Re-check, ts_session Cookie (HttpOnly, SameSite=Strict)

## Ambiguous Edges - Review These
- `FanucCncCommunication.dll` → `CncCommunication.dll (32-bit vendor DLL)`  [AMBIGUOUS]
  docs/setup-guide.txt · relation: semantically_similar_to
- `Snapshot Fields (positions, spindle, feed, program, tool, alarms)` → `Machines API (/machines)`  [AMBIGUOUS]
  docs/CncCommunication.dll — API & Integration Guide (Go Client).pdf · relation: conceptually_related_to

## Knowledge Gaps
- **35 isolated node(s):** `7th-sense-server`, `docs/cnc-communication.md`, `middleware.RequireRole / routes.adminOnly`, `response.SendResponse`, `Test Auth (asTestUser, tests/auth_test.go)` (+30 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 155 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **16 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `CNC Data Collection Flow` connect `CNC Data Collection Flow` to `Worker (internal/monitor/worker.go)`, `VariableConfig`, `Collector`, `CNC Data Collection (in-process)`?**
  _High betweenness centrality (0.202) - this node is a cross-community bridge._
- **What connects `7th-sense-server`, `docs/cnc-communication.md`, `middleware.RequireRole / routes.adminOnly` to the rest of the system?**
  _35 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `VariableConfig` be split into smaller, more focused modules?**
  _Cohesion score 0.12380952380952381 - nodes in this community are weakly interconnected._
- **Why does `2. All steps at a glance` connect `Collector` to `.runConnectionLoop`, `VariableConfig`, `CNC Data Collection Flow`, `Connection`, `worker`?**
  _High betweenness centrality (0.175) - this node is a cross-community bridge._
- **Should `UserCellService` be split into smaller, more focused modules?**
  _Cohesion score 0.07400555041628122 - nodes in this community are weakly interconnected._

### Low-confidence Hints
_AMBIGUOUS edges — the extractor was unsure. Verify before acting on these._

- **What is the exact relationship between `FanucCncCommunication.dll` and `CncCommunication.dll (32-bit vendor DLL)`?**
  _Edge tagged AMBIGUOUS (relation: semantically_similar_to) - confidence is low._
- **What is the exact relationship between `Snapshot Fields (positions, spindle, feed, program, tool, alarms)` and `Machines API (/machines)`?**
  _Edge tagged AMBIGUOUS (relation: conceptually_related_to) - confidence is low._