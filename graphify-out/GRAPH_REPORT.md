# Graph Report - 7th-sense-server  (2026-11-03)

## Corpus Check
- 261 files · ~147,385 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 4 file(s) not represented in the graph (top: (none) 3, .toml 1)

## Summary
- 2217 nodes · 17829 edges · 86 communities (71 shown, 15 thin omitted)
- Extraction: 95% EXTRACTED · 5% INFERRED · 0% AMBIGUOUS · INFERRED: 868 edges (avg confidence: 0.86)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `4f14d6bc`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- user_access_test.go
- UserCellService
- UserSupervisorService
- Shift
- SendResponse()
- pathID()
- SendResponse
- toUserCityResponseError
- api_helpers_test.go
- CncMasterSnapshot (34,440 bytes, #pragma
- ControllerService
- dll_windows_386.go
- UserLineService
- Collector
- worker
- MachineService
- Execute
- testing.T
- VariableConfig
- Handlers
- net/http.ServeMux
- HolidayService
- CompanyProfile
- §10 Diagnostics & Troubleshooting
- .validate
- machine_mode.go
- BranchService
- UserService
- CNC Data Collection Flow
- toUserLocationResponseError
- UserProfileRepository
- newOrganizationRouter
- Worker (internal/monitor/worker.go)
- newUserServices
- MachineModelService
- .Validate
- LicenceKeyService
- BrandService
- AuthUser
- mode-chart.js
- CncCommunication.dll API & Integration G
- config.go
- github.com/jackc/pgx/v5/pgxpool.Pool
- CncCommunication.dll API & Integration G
- Address
- .Login()
- Product
- .Create
- Machines API (/machines)
- MachineModePeriod
- context.Context
- .validate
- UserCountryService
- .Generate
- AuthHandler
- UserService
- net/http.ResponseWriter
- .validate
- .Validate
- RoleService
- Test list
- validText()
- run
- cycleTracker
- CNCConnect(ip, port, dataCb, errCb, mach
- onDataCallback(ptr uintptr) uintptr
- IdleReasonService
- Client Certificate & Key (client_cert.de
- .validate()
- Event-driven data flow (no polling in cl
- auth_tokens SHA-256 Hashed Sessions (mig
- Create Customer in One Transaction (user
- 7th-sense-server
- hello
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
1. `Execute()` - 179 edges
2. `SendResponse()` - 146 edges
3. `NewError()` - 131 edges
4. `run()` - 93 edges
5. `pathID()` - 80 edges
6. `AuthUser` - 71 edges
7. `Querier` - 60 edges
8. `currentUser()` - 59 edges
9. `decodeAndValidate()` - 58 edges
10. `worker` - 56 edges

## Surprising Connections (you probably didn't know these)
- `Machine Identifier Mapping (0 Simulator,1 Fanuc,2 Siemens,3 Syntec,4 Mitsubishi)` --semantically_similar_to--> `Machine Catalog Brand-Controller-Model`  [INFERRED] [semantically similar]
  docs/CncCommunication.dll — API & Integration Guide (Go Client).pdf → .claude/CLAUDE.md
- `Snapshot Fields (positions, spindle, feed, program, tool, alarms)` --conceptually_related_to--> `Machines API (/machines)`  [AMBIGUOUS]
  docs/CncCommunication.dll — API & Integration Guide (Go Client).pdf → .claude/CLAUDE.md
- `Step 12 — Call the DLL's CNCConnect` --references--> `Connection`  [INFERRED]
  docs/cnc-data-collection-flow.md → internal/cnc/cnc.go
- `2. All steps at a glance` --references--> `LoadDriver()`  [INFERRED]
  docs/cnc-data-collection-flow.md → internal/cnc/cnc.go
- `How a cycle ends` --references--> `decideCycleStep()`  [INFERRED]
  docs/machine-cycles.md → internal/monitor/machine_cycle.go

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

## Communities (86 total, 15 thin omitted)

### Community 0 - "user_access_test.go"
Cohesion: 0.35
Nodes (3): LicenceKeyDuration, licenceKeyDurations(), body

### Community 1 - "UserCellService"
Cohesion: 0.14
Nodes (10): CreateUserCellRequest, UpdateUserCellRequest, validateUserCell(), NewUserCellHandler(), UserCell, UserCellRepository, mapUserCellError(), UserCellService (+2 more)

### Community 2 - "UserSupervisorService"
Cohesion: 0.15
Nodes (10): CreateUserSupervisorRequest, UpdateUserSupervisorRequest, validateUserSupervisor(), NewUserSupervisorHandler(), UserSupervisor, UserSupervisorRepository, mapUserSupervisorError(), UserSupervisorService (+2 more)

### Community 3 - "Shift"
Cohesion: 0.12
Nodes (21): clockMinutes(), CreateShiftRequest, ShiftFields, UpdateShiftRequest, MinutesBetween(), validWorkingDays(), NewShiftHandler(), Shift (+13 more)

### Community 4 - "SendResponse()"
Cohesion: 0.09
Nodes (11): BreakHandler, ClienteleHandler, DowntimeReasonHandler, HolidayHandler, IdleReasonHandler, PartHandler, ProductionScheduleHandler, currentUser() (+3 more)

### Community 5 - "pathID()"
Cohesion: 0.10
Nodes (12): validator, ControllerHandler, decodeAndValidate(), pathID(), targetUserID(), UserCellHandler, UserCityHandler, UserLineHandler (+4 more)

### Community 6 - "SendResponse"
Cohesion: 0.11
Nodes (6): BrandHandler, CategoryHandler, LicenceKeyHandler, MachineModelHandler, UserCountryHandler, SendResponse()

### Community 7 - "toUserCityResponseError"
Cohesion: 0.24
Nodes (4): CreateUserCityRequest, UpdateUserCityRequest, validateUserCity(), toUserCityResponseError()

### Community 8 - "api_helpers_test.go"
Cohesion: 0.41
Nodes (5): cncMasterSnapshot, LanguageResponse, ThemeResponse, authUserKey, machineModelRecord

### Community 9 - "CncMasterSnapshot (34,440 bytes, #pragma"
Cohesion: 0.10
Nodes (43): CncMasterSnapshot (34,440 bytes, #pragma pack(8), x86), AbsPos ([5]float64, Inches / Degrees, Work Coordinate System (WCS) positions), ActiveGCode ([16]byte, ASCII, Currently active G-code group), ActiveProgPath ([256]byte, ASCII Path, Storage path of active NC artifact), AlarmMsg ([128]byte, UTF-8 / ASCII, Alarm description message), AlarmNo (int32, Code (0 = OK), Primary alarm or fault number), AxisLoadXPercentage (int16, Percent (%), Servo X axis drive load (AxisLoadX/Y/ZPercentage)), AxisLoadYPercentage (int16, Percent (%), Servo Y axis drive load (AxisLoadX/Y/ZPercentage)) (+35 more)

### Community 10 - "ControllerService"
Cohesion: 0.14
Nodes (9): CreateControllerRequest, UpdateControllerRequest, validateControllerFields(), NewControllerHandler(), Controller, ControllerRepository, mapControllerError(), ControllerService (+1 more)

### Community 11 - "dll_windows_386.go"
Cohesion: 0.10
Nodes (17): dllDriver, slot, Step 13 — The DLL sends data, Step 17 — Disconnect and try again, Connection, Reading, Stats, Status (+9 more)

### Community 12 - "UserLineService"
Cohesion: 0.07
Nodes (19): CreateUserLineRequest, UpdateUserLineRequest, validateUserLine(), CreateUserShopRequest, UpdateUserShopRequest, validateUserShop(), NewUserLineHandler(), NewUserShopHandler() (+11 more)

### Community 13 - "Collector"
Cohesion: 0.11
Nodes (24): 2. All steps at a glance, 4. Find the machines: steps 5 – 9, 6. When data is edited: step 18, Log lines (in `logs/app.log`, search for the message), Step 12 — Call the DLL's CNCConnect, Step 18 — React to machine and variable changes, Step 20 — Stop the workers and mark machines offline, Step 3 — Run the collector in the background (+16 more)

### Community 14 - "worker"
Cohesion: 0.11
Nodes (19): 5. Connect and collect, per machine: steps 10 – 17, Known issue, Step 10 — Get ready to connect, Step 11 — Connect, a few machines at a time, Step 14 — The data collection loop, Step 15 — Store the reading, Step 16 — Every second: dead connection check and live status, Value (+11 more)

### Community 15 - "MachineService"
Cohesion: 0.06
Nodes (32): One running cycle per machine, NewSupervisor(), MachineCyclesResponse, MachineCycleSummary, CreateMachineRequest, MachineFields, UpdateMachineRequest, validCoordinate() (+24 more)

### Community 16 - "Execute"
Cohesion: 0.06
Nodes (28): Execute(), UserCity, User, UserLocation, AuthRepository, mapAuthError(), ClienteleRepository, mapClienteleError() (+20 more)

### Community 17 - "testing.T"
Cohesion: 0.05
Nodes (104): UpdateThemeRequest, TestSnapshotModeFieldOffsets(), TestSnapshotSize(), TestCreateUserRequestValidateDropsDuplicateMachines(), TestCreateUserRequestValidateNormalizes(), TestCustomerFiltersValidate(), TestUpdateThemeRequestValidate(), validCustomerDetails() (+96 more)

### Community 18 - "VariableConfig"
Cohesion: 0.22
Nodes (7): VariableConfig, formatPreviousValue(), formatValue(), newValueChangeLog(), boundVariable, valueChange, valueChangeLog

### Community 19 - "Handlers"
Cohesion: 0.13
Nodes (5): HealthHandler, MachineCycleHandler, MachineHandler, registerHealthRoutes(), Handlers

### Community 20 - "net/http.ServeMux"
Cohesion: 0.16
Nodes (31): RequireRole(), registerBranchRoutes(), registerBrandRoutes(), registerBreakRoutes(), registerCategoryRoutes(), registerClienteleRoutes(), registerControllerRoutes(), registerDowntimeReasonRoutes() (+23 more)

### Community 21 - "HolidayService"
Cohesion: 0.15
Nodes (12): CreateHolidayRequest, HolidayFields, UpdateHolidayRequest, NewHolidayHandler(), Holiday, HolidayRepository, mapHolidayError(), HolidayService (+4 more)

### Community 22 - "CompanyProfile"
Cohesion: 0.28
Nodes (5): CompanyProfile, CompanyProfileRepository, mapCompanyProfileError(), toCompanyProfileResponseError(), CompanyProfileService

### Community 23 - "§10 Diagnostics & Troubleshooting"
Cohesion: 0.12
Nodes (25): Check mode: go run sinumerik_client.go check (OPC UA node existence report), Go Client (sinumerik_client.go, callback -> copy -> render), Prerequisite: 32-bit x86 DLL, compile Go with GOARCH=386, Prerequisite: SIEMENS_USER / SIEMENS_PASS env vars before CNCConnect, Run: go run sinumerik_client.go (live dashboard), §3.1 Environment & Prerequisites, §8 Running the Go Client, §10 Diagnostics & Troubleshooting (+17 more)

### Community 24 - ".validate"
Cohesion: 0.16
Nodes (11): MachineCycleListRequest, validDate(), CreateProductionScheduleRequest, ProductionScheduleFields, ProductionScheduleItemFields, UpdateProductionScheduleRequest, ScheduleDay(), asInvalidSelection() (+3 more)

### Community 25 - "machine_mode.go"
Cohesion: 0.13
Nodes (22): engineName(), cycleReadingFrom(), decideMachineMode(), holdReason(), machineStateFromReading(), numberingEngine(), programStateFromProcStatus(), ncReading() (+14 more)

### Community 26 - "BranchService"
Cohesion: 0.15
Nodes (9): CreateBranchRequest, UpdateBranchRequest, validateBranchFields(), Branch, BranchRepository, mapBranchError(), cityOwner(), BranchService (+1 more)

### Community 27 - "UserService"
Cohesion: 0.23
Nodes (12): CreateUserRequest, CustomerDetails, UpdateUserRequest, UserResponse, ParseDayStart(), TestParseDayStart(), companyProfileFrom(), primaryAddress() (+4 more)

### Community 28 - "CNC Data Collection Flow"
Cohesion: 0.08
Nodes (21): Supervisor, 1. The big picture, 3. Startup: steps 1 – 4, 7. Shutdown: steps 19 – 20, 8. Reference, 9. Data Reading Flow (Function Names Only), docs/cnc-communication.md, CNC Data Collection Flow (+13 more)

### Community 29 - "toUserLocationResponseError"
Cohesion: 0.20
Nodes (8): CreateUserLocationRequest, UpdateUserLocationRequest, validateUserLocation(), NewUserLocationHandler(), UserLocationService, NewUserLocationService(), toUserLocationResponseError(), NewUserShopService()

### Community 30 - "UserProfileRepository"
Cohesion: 0.18
Nodes (6): UserProfile, UserProfileRepository, mapProfileError(), NewUserProfileService(), toUserProfileResponseError(), UserProfileService

### Community 31 - "newOrganizationRouter"
Cohesion: 0.18
Nodes (16): NewBranchHandler(), NewUserCityHandler(), NewUserHandler(), NewUserStateHandler(), UserState, NewUserCountryRepository(), UserStateRepository, mapUserStateError() (+8 more)

### Community 32 - "Worker (internal/monitor/worker.go)"
Cohesion: 0.06
Nodes (27): CNC Data Collection (in-process), writeLog CNC Log Types and Switches, .env Configuration (internal/config), 7th Sense Server Claude Context, response.SendResponse Format, User CRUD Golden Reference, CncCommunication.dll (32-bit vendor DLL), Collector (internal/monitor/collector.go) (+19 more)

### Community 33 - "newUserServices"
Cohesion: 0.17
Nodes (15): NewAddressRepository(), NewCompanyProfileRepository(), UserMachineAccessRepository, NewUserMachineAccessRepository(), NewUserProfileRepository(), NewAddressService(), NewCompanyProfileService(), newUserServices() (+7 more)

### Community 34 - "MachineModelService"
Cohesion: 0.19
Nodes (7): MachineModelRequest, NewMachineModelHandler(), MachineModel, MachineModelRepository, mapMachineModelError(), MachineModelService, toMachineModelResponseError()

### Community 35 - ".Validate"
Cohesion: 0.15
Nodes (14): CreateAddressRequest, UpdateAddressRequest, UpdateLanguageRequest, isTimeZone(), isWebsite(), normalizeFilters(), normalizeTheme(), TestCreateUserRequestValidate() (+6 more)

### Community 36 - "LicenceKeyService"
Cohesion: 0.17
Nodes (10): NewLicenceKey, CreateLicenceKeysRequest, NewLicenceKeyHandler(), LicenceKey, LicenceKeyRepository, mapLicenceKeyError(), LicenceKeyService, licenceExpiresOn() (+2 more)

### Community 37 - "BrandService"
Cohesion: 0.08
Nodes (16): CreateBrandRequest, UpdateBrandRequest, validateCatalogFields(), CreateCategoryRequest, UpdateCategoryRequest, NewCategoryHandler(), Brand, Category (+8 more)

### Community 38 - "AuthUser"
Cohesion: 0.07
Nodes (38): WithTransaction(), CreatePartRequest, PartFields, UpdatePartRequest, ValidPartDrawingFile(), NewClienteleHandler(), AuthUser, Break (+30 more)

### Community 39 - "mode-chart.js"
Cohesion: 0.08
Nodes (71): Health Mode History API (public), machine_mode_log, ECharts 5.5.1, Mode List Table and Time per Mode, Mode Count Tiles (filter), Machine Mode Dashboard Page, Zoom Level Bar, card() (+63 more)

### Community 40 - "CncCommunication.dll API & Integration G"
Cohesion: 0.19
Nodes (20): CncCommunication.dll (C++ core runtime, x86, OPC UA inside), CncCommunication.dll API & Integration Guide (Go Client) - text version, Fanuc profile (engine 1, port 8193, Fwlib32.dll + fwlibe1.dll), Fwlib32.dll (Fanuc FOCAS2), fwlibe1.dll (Fanuc FOCAS2), libcrypto-3.dll (OpenSSL), libssl-3.dll (OpenSSL), Engine 1: Fanuc (FOCAS2 Ethernet, fwlib32) (+12 more)

### Community 41 - "config.go"
Cohesion: 0.05
Nodes (41): main(), run(), main(), Config, MachineSelection, CncMasterSnapshot, MachineDef, cstr() (+33 more)

### Community 42 - "github.com/jackc/pgx/v5/pgxpool.Pool"
Cohesion: 0.24
Nodes (15): NewBrandHandler(), NewBranchRepository(), NewBrandRepository(), NewCategoryRepository(), NewControllerRepository(), NewMachineModelRepository(), NewMachineRepository(), NewProductRepository() (+7 more)

### Community 43 - "CncCommunication.dll API & Integration G"
Cohesion: 0.13
Nodes (18): CNCConnect param machineType (int32 engine id, see §4.2), Engine 0: Simulator (local loopback mock engine), Engine 2: Siemens (Sinumerik OPC UA, active profile), Engine 3: Syntec (Open API / TCP), Engine 4: Mitsubishi (driver not specified), Machine identifier mapping (engine id enum 0-4), §9 Controller-Specific Profiles, Siemens default ports: 4840 (native NCU) / 4841 (custom) (+10 more)

### Community 44 - "Address"
Cohesion: 0.27
Nodes (5): Address, AddressRepository, mapAddressError(), toAddressResponseError(), AddressService

### Community 45 - ".Login()"
Cohesion: 0.15
Nodes (11): LoginRequest, NewAuthHandler(), AuthService, hashToken(), NewAuthService(), newSessionToken(), truncate(), newLoginThrottle() (+3 more)

### Community 46 - "Product"
Cohesion: 0.15
Nodes (11): CreateProductRequest, UpdateProductRequest, validateProductFields(), NewProductHandler(), Product, ProductRepository, mapProductError(), ProductService (+3 more)

### Community 47 - ".Create"
Cohesion: 0.29
Nodes (3): CreateUserStateRequest, validateUserState(), toUserStateResponseError()

### Community 48 - "Machines API (/machines)"
Cohesion: 0.16
Nodes (18): Clienteles API, Company Scope (companyOwner / ownCompanyCustomer), Current User API (/me theme, language), Customer Display Names (customerName / companyName), Holidays API, Licence Keys API, Machine Catalog Brand-Controller-Model, Machines API (/machines) (+10 more)

### Community 49 - "MachineModePeriod"
Cohesion: 0.53
Nodes (4): MachineModeHistoryResponse, MachineModeTotal, MachineModePeriod, modeTotals()

### Community 50 - "context.Context"
Cohesion: 0.05
Nodes (23): Step 5 — Prepare the database, MachineModesResponse, MachineEvent, MachineLiveStatus, MachineModeStatus, ModeChange, PartProcess, PartSubProgram (+15 more)

### Community 51 - ".validate"
Cohesion: 0.47
Nodes (3): CreateDowntimeReasonRequest, DowntimeReasonFields, UpdateDowntimeReasonRequest

### Community 52 - "UserCountryService"
Cohesion: 0.20
Nodes (8): CreateUserCountryRequest, NewUserCountryHandler(), UserCountry, UserCountryRepository, mapUserCountryError(), UserCountryService, NewUserCountryService(), toUserCountryResponseError()

### Community 53 - ".Generate"
Cohesion: 0.40
Nodes (3): GenerateLicenceKeysRequest, GeneratedLicenceKeys, newLicenceKey()

### Community 54 - "AuthHandler"
Cohesion: 0.38
Nodes (3): AuthHandler, CurrentUser(), UserID()

### Community 55 - "UserService"
Cohesion: 0.15
Nodes (8): UserListFilter, AssignableRoles(), CanAssignRole(), SeesAllCompanyMachines(), canSee(), UserService, toUserResponseError(), userAccess

### Community 56 - "net/http.ResponseWriter"
Cohesion: 0.14
Nodes (5): BranchHandler, ProductHandler, currentUserID(), UserHandler, statusRecorder

### Community 58 - ".validate"
Cohesion: 0.47
Nodes (3): BreakFields, CreateBreakRequest, UpdateBreakRequest

### Community 59 - ".Validate"
Cohesion: 0.40
Nodes (3): UpdateLicenceKeyRequest, validDuration(), ValidLicenceKey()

### Community 60 - "RoleService"
Cohesion: 0.23
Nodes (9): RoleHandler, NewRoleHandler(), Role, RoleRepository, mapRoleError(), RoleService, NewRoleService(), toRoleResponseError() (+1 more)

### Community 61 - "Test list"
Cohesion: 0.05
Nodes (38): 1. The rows are consistent (SQL, no machine needed), 2. The rows agree with the mode log (SQL), 3. The raw signals around one cycle (telemetry and logs), 4. At the machine (floor test), Before the tests (once), Check query, Completed, Floor tests on the HUMA panel (Fanuc) (+30 more)

### Community 62 - "validText()"
Cohesion: 0.18
Nodes (6): PartProcessFields, PartSubProgramFields, UploadPartDrawingRequest, uniquePositiveIDs(), validText(), PartDrawing

### Community 63 - "run"
Cohesion: 0.16
Nodes (22): run(), Querier, NewBreakHandler(), NewDowntimeReasonHandler(), NewPartHandler(), NewProductionScheduleHandler(), NewAuthRepository(), NewClienteleRepository() (+14 more)

### Community 64 - "cycleTracker"
Cohesion: 0.13
Nodes (14): CycleStatusFor(), decideCycleStep(), isProgramLoop(), resetReason(), secondsOf(), readingAt(), TestAwaitNextCountsFromItsStart(), TestDecideCycleStepSubProgramCallIsNotALoop() (+6 more)

### Community 69 - "CNCConnect(ip, port, dataCb, errCb, mach"
Cohesion: 0.20
Nodes (11): CNCConnect(ip, port, dataCb, errCb, machineType) int32, CNCConnect param dataCb (uintptr stdcall OnData(const CncMasterSnapshot*)), CNCConnect param ip (const char*, ANSI NUL-terminated), CNCConnect param port (int32, e.g. 4840), CNCdisconnect(machineId int32), LazyDLL binding (syscall.NewLazyDLL + eager mod.Load in init), §4 API Reference, §5 DLL Loading & Binding (Go) (+3 more)

### Community 71 - "onDataCallback(ptr uintptr) uintptr"
Cohesion: 0.36
Nodes (9): Package-level callback singletons (syscall.NewCallback), CNCConnect param errCb (uintptr stdcall OnErr(int32 code, const char* msg)), Mutex-protected copy: dataMutex sync.RWMutex -> lastRecord/hasRecord, onDataCallback(ptr uintptr) uintptr, onErrorCallback(code, messagePtr uintptr) uintptr, §7 Asynchronous Callback Architecture, onDataCallback, onErrorCallback (+1 more)

### Community 72 - "IdleReasonService"
Cohesion: 0.15
Nodes (11): CreateIdleReasonRequest, IdleReasonFields, UpdateIdleReasonRequest, NewIdleReasonHandler(), IdleReason, IdleReasonRepository, mapIdleReasonError(), IdleReasonService (+3 more)

### Community 77 - "Client Certificate & Key (client_cert.de"
Cohesion: 0.48
Nodes (7): client_cert.der/.pem (X.509 client cert for Siemens OPC UA), client_key.der/.pem (client private key), Siemens certs: client_cert.der + client_key.der in binary directory, Siemens trust: cert must be trusted in NCU Access Administration, Troubleshoot: BadRequestNotAllowed -> Export client_cert, BadRequestNotAllowed (certificate trust rejection), Client Certificate & Key (client_cert.der / client_key.der)

### Community 78 - ".validate()"
Cohesion: 0.38
Nodes (4): ClienteleFields, CreateClienteleRequest, UpdateClienteleRequest, validClienteleEmail()

### Community 80 - "Event-driven data flow (no polling in cl"
Cohesion: 0.60
Nodes (5): Flow 1: load DLL, register two stdcall callbacks (syscall.NewCallback), Flow 2: CNCConnect starts DLL background worker thread, Flow 3: worker polls controller, assembles snapshot, invokes data callback, Flow 4: callback copies payload, signals renderer via buffered channel, Flow 5: main goroutine renders live ANSI dashboard

### Community 86 - "auth_tokens SHA-256 Hashed Sessions (mig"
Cohesion: 0.67
Nodes (3): auth_tokens SHA-256 Hashed Sessions (migration 014), Per-Request Session and User Re-check, ts_session Cookie (HttpOnly, SameSite=Strict)

## Ambiguous Edges - Review These
- `CncCommunication.dll (32-bit vendor DLL)` → `FanucCncCommunication.dll`  [AMBIGUOUS]
  docs/setup-guide.txt · relation: semantically_similar_to
- `Snapshot Fields (positions, spindle, feed, program, tool, alarms)` → `Machines API (/machines)`  [AMBIGUOUS]
  docs/CncCommunication.dll — API & Integration Guide (Go Client).pdf · relation: conceptually_related_to

## Knowledge Gaps
- **36 isolated node(s):** `hello`, `7th-sense-server`, `middleware.RequireRole / routes.adminOnly`, `response.SendResponse`, `Test Auth (asTestUser, tests/auth_test.go)` (+31 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 131 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **15 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `CNC Data Collection Flow` connect `CNC Data Collection Flow` to `Worker (internal/monitor/worker.go)`, `Collector`, `worker`?**
  _High betweenness centrality (0.163) - this node is a cross-community bridge._
- **What connects `hello`, `7th-sense-server`, `middleware.RequireRole / routes.adminOnly` to the rest of the system?**
  _36 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `UserCellService` be split into smaller, more focused modules?**
  _Cohesion score 0.14492753623188406 - nodes in this community are weakly interconnected._
- **Why does `2. All steps at a glance` connect `Collector` to `dll_windows_386.go`, `CNC Data Collection Flow`, `worker`?**
  _High betweenness centrality (0.143) - this node is a cross-community bridge._
- **Should `Shift` be split into smaller, more focused modules?**
  _Cohesion score 0.12299465240641712 - nodes in this community are weakly interconnected._

### Low-confidence Hints
_AMBIGUOUS edges — the extractor was unsure. Verify before acting on these._

- **What is the exact relationship between `CncCommunication.dll (32-bit vendor DLL)` and `FanucCncCommunication.dll`?**
  _Edge tagged AMBIGUOUS (relation: semantically_similar_to) - confidence is low._
- **What is the exact relationship between `Snapshot Fields (positions, spindle, feed, program, tool, alarms)` and `Machines API (/machines)`?**
  _Edge tagged AMBIGUOUS (relation: conceptually_related_to) - confidence is low._