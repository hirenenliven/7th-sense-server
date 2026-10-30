# Graph Report - 7th-sense-server  (2026-10-30)

## Corpus Check
- 252 files · ~131,888 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 4 file(s) not represented in the graph (top: (none) 3, .toml 1)

## Summary
- 2027 nodes · 16661 edges · 99 communities (81 shown, 18 thin omitted)
- Extraction: 95% EXTRACTED · 5% INFERRED · 0% AMBIGUOUS · INFERRED: 826 edges (avg confidence: 0.86)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `c977f856`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- HTTP Handlers Layer
- testing.T
- .collectReadings
- Querier
- Route Registration & Guards
- MachineService
- UserLineService
- CNC Snapshot Struct Fields
- .checkReferences
- RoleService
- Execute
- Integration & DTO Tests (2)
- UserShopService
- CncCommunication.dll API & Integration Guide (Go Client)
- UserLineService
- Address
- UserLocationService
- mode-chart.js
- add_superadmin_raw.go
- Connection
- UserService
- Product
- HTTP Handlers Layer (2)
- IdleReasonRepository
- CNC Go Client Run & Check Mode
- user_access_test.go
- validateCatalogFields
- UserCountryService
- ControllerService
- UserProfileRepository
- run
- PartService
- run
- HolidayService
- .validate
- sinumerik_client.go
- Execute
- UserRepository
- CNC Connect API & Callbacks
- Server API Domain Docs
- decodeAndValidate
- net/http.ResponseWriter
- config.go
- Collector
- BrandService
- AuthUser
- CRUD Implementation Pattern (8 steps)
- auth_tokens SHA-256 Hashed Sessions (migration 014)
- GET /health
- MachineModelService
- validText
- CncCommunication.dll (C++ core runtime, x86, OPC UA inside)
- AuthService
- generate.go
- CNC Event-Driven Data Flow
- demo_monitor/main.go
- Check-mode OPC UA node mapping
- ProductionScheduleRepository
- Create Customer in One Transaction (user, company profile, user profile, address)
- decideMachineMode
- Sinumerik DLL Go Client (use case) (2)
- Adding a New Entity (README steps)
- middleware.CrossOrigin (http.CrossOriginProtection)
- .Create
- Go Module Root
- Current User Middleware
- Role Guard Middleware
- SendResponse Helper
- Test Auth Helpers
- CncCommunication.dll API & Integration Guide (Go Client) - text version
- MachineTelemetryRepository
- .flush
- Login Lockout (5 wrong passwords, 15 min, 429)
- Naming Conventions
- .validate
- .validate
- AuthHandler
- UserCellService
- toUserCityResponseError
- AuthUser
- CNCConfig
- logger/logger.go
- Category
- .accessFor
- newOrganizationRouter
- UserRepository
- ClienteleRepository
- .Generate
- .validate
- time.Time
- .Validate
- AuthService
- NewError
- Role
- AuthRepository
- LocationRepository
- .UploadDrawing
- .GetAll

## God Nodes (most connected - your core abstractions)
1. `Execute()` - 172 edges
2. `SendResponse()` - 143 edges
3. `NewError()` - 130 edges
4. `run()` - 90 edges
5. `pathID()` - 78 edges
6. `AuthUser` - 69 edges
7. `Querier` - 58 edges
8. `decodeAndValidate()` - 58 edges
9. `currentUser()` - 57 edges
10. `worker` - 43 edges

## Surprising Connections (you probably didn't know these)
- `Machine Identifier Mapping (0 Simulator,1 Fanuc,2 Siemens,3 Syntec,4 Mitsubishi)` --semantically_similar_to--> `Machine Catalog Brand -> Controller -> Model`  [INFERRED] [semantically similar]
  docs/CncCommunication.dll — API & Integration Guide (Go Client).pdf → .claude/CLAUDE.md
- `Snapshot Fields (positions, spindle, feed, program, tool, alarms)` --conceptually_related_to--> `Machines API (/machines)`  [AMBIGUOUS]
  docs/CncCommunication.dll — API & Integration Guide (Go Client).pdf → .claude/CLAUDE.md
- `2. All steps at a glance` --references--> `LoadDriver()`  [INFERRED]
  docs/cnc-data-collection-flow.md → internal/cnc/cnc.go
- `Step 15 — Store the reading` --references--> `machineStateFromReading()`  [INFERRED]
  docs/cnc-data-collection-flow.md → internal/monitor/machine_mode.go
- `Step 15 — Store the reading` --references--> `programStateFromProcStatus()`  [INFERRED]
  docs/cnc-data-collection-flow.md → internal/monitor/machine_mode.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Event-driven telemetry pipeline (CNCConnect -> worker -> callback -> mutex copy -> render)** — docs_cnc_communication_cncconnect, docs_cnc_communication_flow_step1_register_callbacks, docs_cnc_communication_flow_step2_worker_thread, docs_cnc_communication_flow_step3_poll_and_invoke, docs_cnc_communication_flow_step4_copy_and_signal, docs_cnc_communication_flow_step5_render_dashboard, docs_cnc_communication_ondatacallback, docs_cnc_communication_mutex_copy, docs_cnc_communication_cncmastersnapshot [EXTRACTED 1.00]
- **Siemens OPC UA connection requirements (GOARCH=386, credentials, trusted certs)** — docs_cnc_communication_prereq_goarch_386, docs_cnc_communication_prereq_siemens_credentials, docs_cnc_communication_siemens_certificates, docs_cnc_communication_siemens_trust, docs_cnc_communication_client_cert_der, docs_cnc_communication_client_key_der, docs_cnc_communication_siemens_sinumerik_profile [EXTRACTED 1.00]
- **CNC telemetry callback data flow** — docs_cnccommunication_dll_api_integration_guide_go_client_cncconnect, docs_cnccommunication_dll_api_integration_guide_go_client_stdcall_callbacks, docs_cnccommunication_dll_api_integration_guide_go_client_ondatacallback, docs_cnccommunication_dll_api_integration_guide_go_client_cncmastersnapshot, docs_cnccommunication_dll_api_integration_guide_go_client_event_driven_data_flow [EXTRACTED 1.00]
- **Company-scoped customer resources** — _claude_claude_parts_api, _claude_claude_clienteles_api, _claude_claude_shifts_breaks_holidays, _claude_claude_licence_keys_api, _claude_claude_production_schedules_api, _claude_claude_company_scope [EXTRACTED 1.00]
- **CRUD layered pattern rules** — _claude_claude_layered_architecture, _claude_claude_dto_validation, _claude_claude_service_error_mapping, _claude_claude_coalesce_nullables, _claude_claude_soft_delete, _claude_claude_user_crud_golden_reference [EXTRACTED 1.00]
- **Struct ABI compatibility (pack(8) x86, size 34440, GOARCH=386)** — docs_cnc_communication_cncmastersnapshot, docs_cnc_communication_struct_size_assertion, docs_cnc_communication_prereq_goarch_386, docs_cnc_communication_ts_memory_panic [INFERRED 0.85]
- **Session authentication and authorization flow** — readme_session_cookie, readme_auth_tokens_hashed, readme_per_request_session_check, readme_login_lockout, readme_csrf_cross_origin, _claude_claude_role_guards [INFERRED 0.85]

## Communities (99 total, 18 thin omitted)

### Community 0 - "HTTP Handlers Layer"
Cohesion: 0.09
Nodes (14): BreakHandler, ClienteleHandler, DowntimeReasonHandler, HolidayHandler, IdleReasonHandler, LicenceKeyHandler, MachineHandler, PartHandler (+6 more)

### Community 1 - "testing.T"
Cohesion: 0.06
Nodes (87): UpdateThemeRequest, TestSnapshotModeFieldOffsets(), TestSnapshotSize(), TestCreateUserRequestValidate(), TestCreateUserRequestValidateDropsDuplicateMachines(), TestCreateUserRequestValidateNormalizes(), TestCustomerFiltersValidate(), TestUpdateThemeRequestValidate() (+79 more)

### Community 2 - ".collectReadings"
Cohesion: 0.16
Nodes (14): 2. All steps at a glance, 5. Connect and collect, per machine: steps 10 – 17, Known issue, Step 10 — Get ready to connect, Step 11 — Connect, a few machines at a time, Step 14 — The data collection loop, Step 15 — Store the reading, Step 16 — Every second: dead connection check and live status (+6 more)

### Community 3 - "Querier"
Cohesion: 0.20
Nodes (9): Querier, NewDowntimeReasonRepository(), NewIdleReasonRepository(), NewMachineTelemetryRepository(), NewUserCellRepository(), NewUserLineRepository(), NewUserShopRepository(), NewUserSupervisorRepository() (+1 more)

### Community 4 - "Route Registration & Guards"
Cohesion: 0.17
Nodes (30): RequireRole(), registerBranchRoutes(), registerBrandRoutes(), registerBreakRoutes(), registerCategoryRoutes(), registerClienteleRoutes(), registerControllerRoutes(), registerDowntimeReasonRoutes() (+22 more)

### Community 5 - "MachineService"
Cohesion: 0.10
Nodes (20): MachineModesResponse, CreateMachineRequest, MachineFields, UpdateMachineRequest, validCoordinate(), validHealthIDs(), validIPv4(), NewHealthHandler() (+12 more)

### Community 6 - "UserLineService"
Cohesion: 0.15
Nodes (9): CreateUserLineRequest, UpdateUserLineRequest, validateUserLine(), NewUserLineHandler(), UserLine, UserLineRepository, mapUserLineError(), UserLineService (+1 more)

### Community 7 - "CNC Snapshot Struct Fields"
Cohesion: 0.11
Nodes (36): CncMasterSnapshot (34,440 bytes, #pragma pack(8), x86), AbsPos ([5]float64, Inches / Degrees, Work Coordinate System (WCS) positions), ActiveGCode ([16]byte, ASCII, Currently active G-code group), ActiveProgPath ([256]byte, ASCII Path, Storage path of active NC artifact), AlarmMsg ([128]byte, UTF-8 / ASCII, Alarm description message), AlarmNo (int32, Code (0 = OK), Primary alarm or fault number), AxisLoadXPercentage (int16, Percent (%), Servo X axis drive load (AxisLoadX/Y/ZPercentage)), AxisLoadYPercentage (int16, Percent (%), Servo Y axis drive load (AxisLoadX/Y/ZPercentage)) (+28 more)

### Community 8 - ".checkReferences"
Cohesion: 0.22
Nodes (9): validDate(), CreateProductionScheduleRequest, ProductionScheduleFields, ProductionScheduleItemFields, UpdateProductionScheduleRequest, ScheduleDay(), checkJobsOverlap(), shiftSpan() (+1 more)

### Community 9 - "RoleService"
Cohesion: 0.38
Nodes (7): RoleHandler, NewRoleHandler(), RoleRepository, NewRoleRepository(), RoleService, NewRoleService(), NewUserService()

### Community 10 - "Execute"
Cohesion: 0.16
Nodes (8): CreateUserCityRequest, UpdateUserCityRequest, validateUserCity(), UserCity, UserCityRepository, mapUserCityError(), UserCityService, toUserCityResponseError()

### Community 11 - "Integration & DTO Tests (2)"
Cohesion: 0.42
Nodes (5): cncMasterSnapshot, LanguageResponse, ThemeResponse, authUserKey, machineModelRecord

### Community 12 - "UserShopService"
Cohesion: 0.15
Nodes (9): CreateUserShopRequest, UpdateUserShopRequest, validateUserShop(), NewUserShopHandler(), UserShop, UserShopRepository, mapUserShopError(), UserShopService (+1 more)

### Community 13 - "CncCommunication.dll API & Integration Guide (Go Client)"
Cohesion: 0.13
Nodes (24): Machine Catalog Brand -> Controller -> Model, client_cert.der/.pem (X.509 client cert for Siemens OPC UA), client_key.der/.pem (client private key), Engine 0: Simulator (local loopback mock engine), Engine 2: Siemens (Sinumerik OPC UA, active profile), Engine 3: Syntec (Open API / TCP), Engine 4: Mitsubishi (driver not specified), Machine identifier mapping (engine id enum 0-4) (+16 more)

### Community 14 - "UserLineService"
Cohesion: 0.23
Nodes (6): LicenceKey, LicenceKeyRepository, mapLicenceKeyError(), LicenceKeyService, licenceExpiresOn(), toLicenceKeyResponseError()

### Community 15 - "Address"
Cohesion: 0.24
Nodes (6): Address, AddressRepository, mapAddressError(), toAddressResponseError(), primaryAddress(), AddressService

### Community 16 - "UserLocationService"
Cohesion: 0.16
Nodes (8): CreateUserLocationRequest, UpdateUserLocationRequest, validateUserLocation(), UserLocation, UserLocationRepository, mapUserLocationError(), UserLocationService, toUserLocationResponseError()

### Community 17 - "mode-chart.js"
Cohesion: 0.14
Nodes (39): cell(), chartData(), chartRows(), countByGroup(), currentSpan(), currentWindow(), drawChart(), drawDashboard() (+31 more)

### Community 18 - "add_superadmin_raw.go"
Cohesion: 0.38
Nodes (10): checkNotTaken(), createSuperadmin(), ensureProfile(), execute(), main(), parseFlags(), run(), superadminRoleID() (+2 more)

### Community 19 - "Connection"
Cohesion: 0.08
Nodes (22): dllDriver, slot, Step 12 — Call the DLL's CNCConnect, Step 13 — The DLL sends data, Step 17 — Disconnect and try again, Connection, Reading, Stats (+14 more)

### Community 20 - "UserService"
Cohesion: 0.18
Nodes (14): CreateUserRequest, CustomerDetails, UpdateUserRequest, UserResponse, ParseDayStart(), TestParseDayStart(), canSee(), companyProfileFrom() (+6 more)

### Community 21 - "Product"
Cohesion: 0.23
Nodes (6): Product, ProductRepository, mapProductError(), ProductService, productImageOrDefault(), toProductResponseError()

### Community 22 - "HTTP Handlers Layer (2)"
Cohesion: 0.08
Nodes (14): CategoryHandler, HealthHandler, pathID(), targetUserID(), UserCellHandler, UserCityHandler, UserCountryHandler, UserLineHandler (+6 more)

### Community 23 - "IdleReasonRepository"
Cohesion: 0.16
Nodes (16): NewAddressRepository(), NewCompanyProfileRepository(), UserMachineAccessRepository, NewUserMachineAccessRepository(), NewUserProfileRepository(), NewAddressService(), NewCompanyProfileService(), NewUserProfileService() (+8 more)

### Community 24 - "CNC Go Client Run & Check Mode"
Cohesion: 0.13
Nodes (22): CncCommunication.dll API & Integration Guide (Go Client) - text version, IsConnected (int16, Boolean (0/1), Network connection state), Prerequisite: 32-bit x86 DLL, compile Go with GOARCH=386, §1 Overview, §3.1 Environment & Prerequisites, §6 Snapshot Data Structure, §10 Diagnostics & Troubleshooting, Struct size assertion: unsafe.Sizeof(CncMasterSnapshot{}) == 34440 (+14 more)

### Community 25 - "user_access_test.go"
Cohesion: 0.36
Nodes (3): LicenceKeyDuration, licenceKeyDurations(), body

### Community 26 - "validateCatalogFields"
Cohesion: 0.11
Nodes (11): CreateBrandRequest, UpdateBrandRequest, validateCatalogFields(), CreateCategoryRequest, UpdateCategoryRequest, CreateControllerRequest, UpdateControllerRequest, validateControllerFields() (+3 more)

### Community 27 - "UserCountryService"
Cohesion: 0.18
Nodes (9): CreateUserCountryRequest, NewUserCountryHandler(), UserCountry, UserCountryRepository, mapUserCountryError(), NewUserCountryRepository(), UserCountryService, NewUserCountryService() (+1 more)

### Community 28 - "ControllerService"
Cohesion: 0.21
Nodes (5): Controller, ControllerRepository, mapControllerError(), ControllerService, toControllerResponseError()

### Community 29 - "UserProfileRepository"
Cohesion: 0.36
Nodes (3): UserProfile, toUserProfileResponseError(), UserProfileService

### Community 30 - "run"
Cohesion: 0.14
Nodes (30): run(), NewBrandHandler(), NewCategoryHandler(), NewClienteleHandler(), NewControllerHandler(), NewMachineModelHandler(), NewPartHandler(), NewProductHandler() (+22 more)

### Community 31 - "PartService"
Cohesion: 0.10
Nodes (18): PartProcessFields, PartSubProgramFields, CreatePartRequest, PartFields, UpdatePartRequest, ValidPartDrawingFile(), Part, PartProcess (+10 more)

### Community 32 - "run"
Cohesion: 0.15
Nodes (10): CreateUserSupervisorRequest, UpdateUserSupervisorRequest, validateUserSupervisor(), NewUserSupervisorHandler(), UserSupervisor, UserSupervisorRepository, mapUserSupervisorError(), UserSupervisorService (+2 more)

### Community 33 - "HolidayService"
Cohesion: 0.15
Nodes (11): CreateHolidayRequest, HolidayFields, UpdateHolidayRequest, Holiday, HolidayRepository, mapHolidayError(), HolidayService, holidayFromFields() (+3 more)

### Community 34 - ".validate"
Cohesion: 0.47
Nodes (3): CreateIdleReasonRequest, IdleReasonFields, UpdateIdleReasonRequest

### Community 35 - "sinumerik_client.go"
Cohesion: 0.20
Nodes (17): axisColor(), cstr(), enableVT(), findLibDir(), initDLL(), isRotary(), lbl(), logf() (+9 more)

### Community 36 - "Execute"
Cohesion: 0.15
Nodes (5): Execute(), UserState, UserProfileRepository, mapProfileError(), mapUserStateError()

### Community 37 - "UserRepository"
Cohesion: 0.15
Nodes (22): NewBreakHandler(), NewHolidayHandler(), NewLicenceKeyHandler(), NewMachineHandler(), NewProductionScheduleHandler(), NewShiftHandler(), CrossOrigin(), Logger() (+14 more)

### Community 38 - "CNC Connect API & Callbacks"
Cohesion: 0.42
Nodes (8): Package-level callback singletons (syscall.NewCallback), Mutex-protected copy: dataMutex sync.RWMutex -> lastRecord/hasRecord, onDataCallback(ptr uintptr) uintptr, onErrorCallback(code, messagePtr uintptr) uintptr, §7 Asynchronous Callback Architecture, onDataCallback, onErrorCallback, stdcall Callbacks via syscall.NewCallback

### Community 39 - "Server API Domain Docs"
Cohesion: 0.23
Nodes (13): Clienteles API (/clienteles), Company Scoping (companyOwner / ownCompanyCustomer), Customer Display Names (customerName / companyName), Licence Keys API (/licence-keys), Machines API (/machines), Organization Hierarchy APIs (/user-country..supervisor), Parts API (/parts with processes, drawings), Production Schedules API (/production-schedules) (+5 more)

### Community 40 - "decodeAndValidate"
Cohesion: 0.15
Nodes (5): validator, BranchHandler, currentUserID(), decodeAndValidate(), UserHandler

### Community 41 - "net/http.ResponseWriter"
Cohesion: 0.12
Nodes (5): BrandHandler, ControllerHandler, MachineModelHandler, ProductHandler, statusRecorder

### Community 42 - "config.go"
Cohesion: 0.16
Nodes (13): main(), run(), Config, MachineSelection, getEnv(), getEnvBool(), getEnvInt(), getEnvMachineSelection() (+5 more)

### Community 43 - "Collector"
Cohesion: 0.12
Nodes (21): 4. Find the machines: steps 5 – 9, Log lines (in `logs/app.log`, search for the message), Step 18 — React to machine and variable changes, Step 20 — Stop the workers and mark machines offline, Step 3 — Run the collector in the background, Step 6 — Record the start and begin following the database, Step 7 — Listen for database changes, Step 8 — Load all active machines (+13 more)

### Community 44 - "BrandService"
Cohesion: 0.25
Nodes (5): Brand, BrandRepository, mapBrandError(), BrandService, toBrandResponseError()

### Community 45 - "AuthUser"
Cohesion: 0.07
Nodes (31): NewDowntimeReasonHandler(), NewIdleReasonHandler(), AuthUser, Clientele, DowntimeReason, IdleReason, ProductionSchedule, DowntimeReasonRepository (+23 more)

### Community 46 - "CRUD Implementation Pattern (8 steps)"
Cohesion: 0.14
Nodes (8): Air Hot Reload, 7th Sense Server - Claude Context, Engineering Principles Summary, PostgreSQL + TimescaleDB via pgx/v5, SendResponse JSON Response Format, User CRUD Golden Reference, SampleTimestampMs (int64, Milliseconds, Epoch timestamp of snapshot emission), README Layered Layout (Node.js analogy)

### Community 47 - "auth_tokens SHA-256 Hashed Sessions (migration 014)"
Cohesion: 0.67
Nodes (3): auth_tokens SHA-256 Hashed Sessions (migration 014), Per-Request Session and User Re-check, ts_session Cookie (HttpOnly, SameSite=Strict)

### Community 49 - "MachineModelService"
Cohesion: 0.20
Nodes (6): MachineModelRequest, MachineModel, MachineModelRepository, mapMachineModelError(), MachineModelService, toMachineModelResponseError()

### Community 50 - "validText"
Cohesion: 0.24
Nodes (3): uniquePositiveIDs(), validText(), validWorkingDays()

### Community 51 - "CncCommunication.dll (C++ core runtime, x86, OPC UA inside)"
Cohesion: 0.22
Nodes (18): CncCommunication.dll (C++ core runtime, x86, OPC UA inside), Fanuc profile (engine 1, port 8193, Fwlib32.dll + fwlibe1.dll), Fwlib32.dll (Fanuc FOCAS2), fwlibe1.dll (Fanuc FOCAS2), libcrypto-3.dll (OpenSSL), libssl-3.dll (OpenSSL), Engine 1: Fanuc (FOCAS2 Ethernet, fwlib32), OPC UA (Siemens transport, port 4840) (+10 more)

### Community 52 - "AuthService"
Cohesion: 0.13
Nodes (9): LoginRequest, sleepUnlessCancelled(), newSessionToken(), truncate(), newLoginThrottle(), backoff, loginFailures, loginThrottle (+1 more)

### Community 55 - "CNC Event-Driven Data Flow"
Cohesion: 0.60
Nodes (5): Flow 1: load DLL, register two stdcall callbacks (syscall.NewCallback), Flow 2: CNCConnect starts DLL background worker thread, Flow 3: worker polls controller, assembles snapshot, invokes data callback, Flow 4: callback copies payload, signals renderer via buffered channel, Flow 5: main goroutine renders live ANSI dashboard

### Community 56 - "demo_monitor/main.go"
Cohesion: 0.24
Nodes (7): CncMasterSnapshot, MachineDef, cstr(), enableVT(), initDLL(), main(), render()

### Community 57 - "Check-mode OPC UA node mapping"
Cohesion: 0.19
Nodes (14): Check mode: go run sinumerik_client.go check (OPC UA node existence report), GcodeGroup ([512]byte, ASCII (reused), DLL reuses this field for the ';'-separated missing-node list), MachineMode (int16, Enum (0=JOG, 1=MDA, 2=AUTO), Controller operating mode), Go Client (sinumerik_client.go, callback -> copy -> render), Operating mode -> /Bag/State/opMode[u1], Check-mode OPC UA node mapping, Run: go run sinumerik_client.go (live dashboard), §8 Running the Go Client (+6 more)

### Community 58 - "ProductionScheduleRepository"
Cohesion: 0.20
Nodes (4): ProductionScheduleItem, ProductionScheduleRepository, mapProductionScheduleError(), hashToken()

### Community 60 - "decideMachineMode"
Cohesion: 0.18
Nodes (16): engineName(), decideMachineMode(), holdReason(), machineStateFromReading(), numberingEngine(), programStateFromProcStatus(), ncReading(), TestDecideMachineMode() (+8 more)

### Community 64 - ".Create"
Cohesion: 0.29
Nodes (3): CreateUserStateRequest, validateUserState(), toUserStateResponseError()

### Community 70 - "CncCommunication.dll API & Integration Guide (Go Client) - text version"
Cohesion: 0.20
Nodes (11): CNCConnect(ip, port, dataCb, errCb, machineType) int32, CNCConnect param dataCb (uintptr stdcall OnData(const CncMasterSnapshot*)), CNCConnect param errCb (uintptr stdcall OnErr(int32 code, const char* msg)), CNCConnect param ip (const char*, ANSI NUL-terminated), CNCConnect param machineType (int32 engine id, see §4.2), CNCConnect param port (int32, e.g. 4840), CNCdisconnect(machineId int32), LazyDLL binding (syscall.NewLazyDLL + eager mod.Load in init) (+3 more)

### Community 71 - "MachineTelemetryRepository"
Cohesion: 0.13
Nodes (6): MachineEvent, MachineLiveStatus, ModeChange, MachineTelemetryRepository, mapMachineTelemetryError(), execute()

### Community 72 - ".flush"
Cohesion: 0.20
Nodes (7): VariableConfig, formatPreviousValue(), formatValue(), newValueChangeLog(), boundVariable, valueChange, valueChangeLog

### Community 75 - ".validate"
Cohesion: 0.47
Nodes (3): CreateDowntimeReasonRequest, DowntimeReasonFields, UpdateDowntimeReasonRequest

### Community 76 - ".validate"
Cohesion: 0.15
Nodes (9): CreateBranchRequest, UpdateBranchRequest, validateBranchFields(), Branch, BranchRepository, mapBranchError(), cityOwner(), BranchService (+1 more)

### Community 77 - "AuthHandler"
Cohesion: 0.38
Nodes (3): AuthHandler, CurrentUser(), UserID()

### Community 78 - "UserCellService"
Cohesion: 0.14
Nodes (10): CreateUserCellRequest, UpdateUserCellRequest, validateUserCell(), NewUserCellHandler(), UserCell, UserCellRepository, mapUserCellError(), UserCellService (+2 more)

### Community 79 - "toUserCityResponseError"
Cohesion: 0.28
Nodes (5): CompanyProfile, CompanyProfileRepository, mapCompanyProfileError(), toCompanyProfileResponseError(), CompanyProfileService

### Community 80 - "AuthUser"
Cohesion: 0.07
Nodes (30): WithTransaction(), BreakFields, CreateBreakRequest, UpdateBreakRequest, clockMinutes(), CreateShiftRequest, ShiftFields, UpdateShiftRequest (+22 more)

### Community 81 - "CNCConfig"
Cohesion: 0.07
Nodes (24): Supervisor, 1. The big picture, 3. Startup: steps 1 – 4, 6. When data is edited: step 18, 7. Shutdown: steps 19 – 20, 8. Reference, 9. Data Reading Flow (Function Names Only), CNC Data Collection Flow (+16 more)

### Community 82 - "logger/logger.go"
Cohesion: 0.40
Nodes (3): main(), Init(), InitFile()

### Community 83 - "Category"
Cohesion: 0.25
Nodes (5): Category, CategoryRepository, mapCategoryError(), CategoryService, toCategoryResponseError()

### Community 84 - ".accessFor"
Cohesion: 0.22
Nodes (5): AssignableRoles(), CanAssignRole(), SeesAllCompanyMachines(), asInvalidSelection(), userAccess

### Community 85 - "newOrganizationRouter"
Cohesion: 0.19
Nodes (16): NewBranchHandler(), NewUserCityHandler(), NewUserHandler(), NewUserLocationHandler(), NewUserStateHandler(), NewUserCityRepository(), NewUserLocationRepository(), UserStateRepository (+8 more)

### Community 86 - "UserRepository"
Cohesion: 0.36
Nodes (3): User, UserRepository, mapUserError()

### Community 88 - ".Generate"
Cohesion: 0.40
Nodes (3): GenerateLicenceKeysRequest, GeneratedLicenceKeys, newLicenceKey()

### Community 89 - ".validate"
Cohesion: 0.38
Nodes (4): ClienteleFields, CreateClienteleRequest, UpdateClienteleRequest, validClienteleEmail()

### Community 90 - "time.Time"
Cohesion: 0.14
Nodes (11): Step 5 — Prepare the database, Value, MachineModeHistoryResponse, MachineModeTotal, LiveValue, MachineModePeriod, TelemetryValue, modeTotals() (+3 more)

### Community 91 - ".Validate"
Cohesion: 0.29
Nodes (5): NewLicenceKey, CreateLicenceKeysRequest, UpdateLicenceKeyRequest, validDuration(), ValidLicenceKey()

### Community 92 - "AuthService"
Cohesion: 0.47
Nodes (6): NewAuthHandler(), Authenticate(), NewAuthRepository(), AuthService, NewAuthService(), newAuthRouter()

### Community 93 - "NewError"
Cohesion: 0.15
Nodes (13): CreateAddressRequest, UpdateAddressRequest, UpdateLanguageRequest, isTimeZone(), isWebsite(), normalizeFilters(), normalizeTheme(), TestUpdateUserRequestValidateStatus() (+5 more)

### Community 94 - "Role"
Cohesion: 0.33
Nodes (3): Role, mapRoleError(), toRoleResponseError()

### Community 96 - "LocationRepository"
Cohesion: 0.50
Nodes (3): mapLocationError(), NewLocationRepository(), LocationRepository

## Ambiguous Edges - Review These
- `Snapshot Fields (positions, spindle, feed, program, tool, alarms)` → `Machines API (/machines)`  [AMBIGUOUS]
  docs/CncCommunication.dll — API & Integration Guide (Go Client).pdf · relation: conceptually_related_to

## Knowledge Gaps
- **28 isolated node(s):** `hello`, `7th-sense-server`, `Siemens default ports: 4840 (native NCU) / 4841 (custom)`, `Engine 0: Simulator (local loopback mock engine)`, `Engine 4: Mitsubishi (driver not specified)` (+23 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 87 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **18 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `NewError()` connect `NewError` to `HTTP Handlers Layer`, `testing.T`, `MachineService`, `UserLineService`, `.checkReferences`, `Execute`, `UserShopService`, `UserLineService`, `UserLocationService`, `UserService`, `Product`, `HTTP Handlers Layer (2)`, `user_access_test.go`, `UserCountryService`, `ControllerService`, `UserProfileRepository`, `PartService`, `run`, `HolidayService`, `.validate`, `decodeAndValidate`, `BrandService`, `AuthUser`, `MachineModelService`, `validText`, `AuthService`, `.Create`, `.validate`, `.validate`, `UserCellService`, `toUserCityResponseError`, `AuthUser`, `Category`, `.accessFor`, `.Generate`, `.validate`, `.Validate`, `Role`, `.UploadDrawing`?**
  _High betweenness centrality (0.079) - this node is a cross-community bridge._
- **What connects `hello`, `7th-sense-server`, `Siemens default ports: 4840 (native NCU) / 4841 (custom)` to the rest of the system?**
  _28 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `HTTP Handlers Layer` be split into smaller, more focused modules?**
  _Cohesion score 0.09157894736842105 - nodes in this community are weakly interconnected._
- **Why does `SendResponse()` connect `HTTP Handlers Layer` to `Route Registration & Guards`, `UserRepository`, `decodeAndValidate`, `net/http.ResponseWriter`, `AuthHandler`, `HTTP Handlers Layer (2)`, `user_access_test.go`, `AuthService`?**
  _High betweenness centrality (0.028) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.06447368421052632 - nodes in this community are weakly interconnected._
- **Why does `MonitoredMachine` connect `Collector` to `.collectReadings`, `MachineTelemetryRepository`, `.flush`, `Integration & DTO Tests (2)`, `Connection`?**
  _High betweenness centrality (0.024) - this node is a cross-community bridge._

### Low-confidence Hints
_AMBIGUOUS edges — the extractor was unsure. Verify before acting on these._

- **What is the exact relationship between `Snapshot Fields (positions, spindle, feed, program, tool, alarms)` and `Machines API (/machines)`?**
  _Edge tagged AMBIGUOUS (relation: conceptually_related_to) - confidence is low._