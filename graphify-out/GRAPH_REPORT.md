# Graph Report - 7th-sense-server  (2026-10-27)

## Corpus Check
- 250 files · ~128,977 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 4 file(s) not represented in the graph (top: (none) 3, .toml 1)

## Summary
- 1989 nodes · 16349 edges · 86 communities (70 shown, 16 thin omitted)
- Extraction: 95% EXTRACTED · 5% INFERRED · 0% AMBIGUOUS · INFERRED: 792 edges (avg confidence: 0.87)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `61214aa4`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- HTTP Handlers Layer
- testing.T
- NewError
- UserShopService
- Route Registration & Guards
- Machine DTOs
- LicenceKeyService
- CNC Snapshot Struct Fields
- DowntimeReasonService
- context.Context
- UserProfileRepository
- Integration & DTO Tests (2)
- context.Context
- CncCommunication.dll API & Integration Guide (Go Client)
- UserRepository
- AuthUser
- UserLineService
- mode-chart.js
- UserCityService
- Connection
- UserService
- 2. All steps at a glance
- HTTP Handlers Layer (2)
- UserSupervisorService
- CNC Go Client Run & Check Mode
- user_access_test.go
- ClienteleRepository
- validateCatalogFields
- ControllerService
- AuthService
- UserLineService
- LocationRepository
- UserLocationService
- UserStateService
- decodeAndValidate
- sinumerik_client.go
- time.Time
- PartService
- CNC Connect API & Callbacks
- Server API Domain Docs
- .Validate
- SendResponse
- IdleReasonService
- 2. All steps at a glance
- run
- .Generate
- 7th Sense Server Engineering Principles (Skills Doc)
- time.Time
- GET /health
- VariableConfig
- .processReading
- CncCommunication.dll (C++ core runtime, x86, OPC UA inside)
- .Login
- runCatalogCRUD
- CNC Event-Driven Data Flow
- Product
- Check-mode OPC UA node mapping
- containsMachineID
- Execute
- Category
- Sinumerik DLL Go Client (use case) (2)
- MachineModelService
- Go Module Root
- Current User Middleware
- Role Guard Middleware
- SendResponse Helper
- Test Auth Helpers
- CncCommunication.dll API & Integration Guide (Go Client) - text version
- MachineTelemetryService
- github.com/jackc/pgx/v5/pgxpool.Pool
- run
- BrandService
- BranchService
- AuthHandler
- Address
- CompanyProfile
- executeAffected
- Querier
- HolidayRepository
- RoleService
- MachineModelHandler
- statusRecorder

## God Nodes (most connected - your core abstractions)
1. `Execute()` - 171 edges
2. `SendResponse()` - 142 edges
3. `NewError()` - 130 edges
4. `run()` - 90 edges
5. `pathID()` - 78 edges
6. `AuthUser` - 69 edges
7. `Querier` - 58 edges
8. `decodeAndValidate()` - 58 edges
9. `currentUser()` - 57 edges
10. `CncMasterSnapshot (34,440 bytes, #pragma pack(8), x86)` - 38 edges

## Surprising Connections (you probably didn't know these)
- `Machine Identifier Mapping (0 Simulator,1 Fanuc,2 Siemens,3 Syntec,4 Mitsubishi)` --semantically_similar_to--> `Machine Catalog Brand -> Controller -> Model`  [INFERRED] [semantically similar]
  docs/CncCommunication.dll — API & Integration Guide (Go Client).pdf → .claude/CLAUDE.md
- `Snapshot Fields (positions, spindle, feed, program, tool, alarms)` --conceptually_related_to--> `Machines API (/machines)`  [AMBIGUOUS]
  docs/CncCommunication.dll — API & Integration Guide (Go Client).pdf → .claude/CLAUDE.md
- `2. All steps at a glance` --references--> `LoadDriver()`  [INFERRED]
  docs/cnc-data-collection-flow.md → internal/cnc/cnc.go
- `Step 7 — Listen for database changes` --references--> `ListenForDatabaseNotifications()`  [INFERRED]
  docs/cnc-data-collection-flow.md → internal/db/postgres.go
- `Step 4 — Run the collector in the background` --references--> `Collector`  [INFERRED]
  docs/cnc-data-collection-flow.md → internal/monitor/collector.go

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

## Communities (86 total, 16 thin omitted)

### Community 0 - "HTTP Handlers Layer"
Cohesion: 0.10
Nodes (11): BreakHandler, ClienteleHandler, DowntimeReasonHandler, HolidayHandler, IdleReasonHandler, PartHandler, ProductionScheduleHandler, currentUser() (+3 more)

### Community 1 - "testing.T"
Cohesion: 0.10
Nodes (64): TestSnapshotSize(), CreateUserRequest, PartDrawing, TestDecideMachineMode(), decodeData(), expectStatus(), newCatalogRouter(), sendRequest() (+56 more)

### Community 2 - "NewError"
Cohesion: 0.05
Nodes (33): CreateAddressRequest, PartProcessFields, PartSubProgramFields, UpdateAddressRequest, UpdateLanguageRequest, BreakFields, CreateBreakRequest, UpdateBreakRequest (+25 more)

### Community 3 - "UserShopService"
Cohesion: 0.15
Nodes (9): CreateUserShopRequest, UpdateUserShopRequest, validateUserShop(), NewUserShopHandler(), UserShop, UserShopRepository, mapUserShopError(), UserShopService (+1 more)

### Community 4 - "Route Registration & Guards"
Cohesion: 0.17
Nodes (30): RequireRole(), registerBranchRoutes(), registerBrandRoutes(), registerBreakRoutes(), registerCategoryRoutes(), registerClienteleRoutes(), registerControllerRoutes(), registerDowntimeReasonRoutes() (+22 more)

### Community 5 - "Machine DTOs"
Cohesion: 0.12
Nodes (16): CreateMachineRequest, MachineFields, UpdateMachineRequest, validCoordinate(), validHealthIDs(), validIPv4(), Machine, MachineScope (+8 more)

### Community 6 - "LicenceKeyService"
Cohesion: 0.19
Nodes (8): NewLicenceKey, CreateLicenceKeysRequest, LicenceKey, LicenceKeyRepository, mapLicenceKeyError(), LicenceKeyService, licenceExpiresOn(), toLicenceKeyResponseError()

### Community 7 - "CNC Snapshot Struct Fields"
Cohesion: 0.11
Nodes (38): CncMasterSnapshot (34,440 bytes, #pragma pack(8), x86), AbsPos ([5]float64, Inches / Degrees, Work Coordinate System (WCS) positions), ActiveGCode ([16]byte, ASCII, Currently active G-code group), ActiveProgPath ([256]byte, ASCII Path, Storage path of active NC artifact), AlarmMsg ([128]byte, UTF-8 / ASCII, Alarm description message), AlarmNo (int32, Code (0 = OK), Primary alarm or fault number), AxisLoadXPercentage (int16, Percent (%), Servo X axis drive load (AxisLoadX/Y/ZPercentage)), AxisLoadYPercentage (int16, Percent (%), Servo Y axis drive load (AxisLoadX/Y/ZPercentage)) (+30 more)

### Community 8 - "DowntimeReasonService"
Cohesion: 0.15
Nodes (12): CreateDowntimeReasonRequest, DowntimeReasonFields, UpdateDowntimeReasonRequest, NewDowntimeReasonHandler(), DowntimeReason, DowntimeReasonRepository, mapDowntimeReasonError(), NewDowntimeReasonRepository() (+4 more)

### Community 9 - "context.Context"
Cohesion: 0.17
Nodes (3): ProductionScheduleItem, ProductionScheduleRepository, mapProductionScheduleError()

### Community 10 - "UserProfileRepository"
Cohesion: 0.17
Nodes (7): UserProfile, UserProfileRepository, mapProfileError(), NewUserProfileRepository(), NewUserProfileService(), toUserProfileResponseError(), UserProfileService

### Community 11 - "Integration & DTO Tests (2)"
Cohesion: 0.43
Nodes (5): cncMasterSnapshot, LanguageResponse, ThemeResponse, authUserKey, machineModelRecord

### Community 12 - "context.Context"
Cohesion: 0.15
Nodes (7): ListenForDatabaseNotifications(), ModeChange, MachineTelemetryRepository, mapMachineTelemetryError(), NewMachineTelemetryRepository(), execute(), hashToken()

### Community 13 - "CncCommunication.dll API & Integration Guide (Go Client)"
Cohesion: 0.13
Nodes (21): Machine Catalog Brand -> Controller -> Model, CNCConnect param machineType (int32 engine id, see §4.2), Engine 0: Simulator (local loopback mock engine), Engine 2: Siemens (Sinumerik OPC UA, active profile), Engine 3: Syntec (Open API / TCP), Engine 4: Mitsubishi (driver not specified), Machine identifier mapping (engine id enum 0-4), §9 Controller-Specific Profiles (+13 more)

### Community 14 - "UserRepository"
Cohesion: 0.12
Nodes (26): NewBreakHandler(), NewClienteleHandler(), NewHolidayHandler(), NewLicenceKeyHandler(), NewPartHandler(), NewProductionScheduleHandler(), NewShiftHandler(), CrossOrigin() (+18 more)

### Community 15 - "AuthUser"
Cohesion: 0.08
Nodes (36): ValidPartDrawingFile(), CreateProductionScheduleRequest, ProductionScheduleFields, UpdateProductionScheduleRequest, AuthUser, Break, Clientele, Holiday (+28 more)

### Community 16 - "UserLineService"
Cohesion: 0.15
Nodes (9): CreateUserLineRequest, UpdateUserLineRequest, validateUserLine(), NewUserLineHandler(), UserLine, UserLineRepository, mapUserLineError(), UserLineService (+1 more)

### Community 17 - "mode-chart.js"
Cohesion: 0.23
Nodes (21): applyRange(), cell(), chartData(), chartRows(), currentSpan(), drawChart(), drawHeader(), drawPeriodList() (+13 more)

### Community 18 - "UserCityService"
Cohesion: 0.16
Nodes (8): CreateUserCityRequest, UpdateUserCityRequest, validateUserCity(), UserCity, UserCityRepository, mapUserCityError(), UserCityService, toUserCityResponseError()

### Community 19 - "Connection"
Cohesion: 0.11
Nodes (15): dllDriver, slot, Step 12 — Call the DLL's CNCConnect, Step 13 — The DLL sends data, Step 17 — Disconnect and try again, Connection, Reading, Stats (+7 more)

### Community 20 - "UserService"
Cohesion: 0.18
Nodes (13): CustomerDetails, UpdateUserRequest, UserResponse, SeesAllCompanyMachines(), canSee(), companyProfileFrom(), UserService, setCustomerAddress() (+5 more)

### Community 21 - "2. All steps at a glance"
Cohesion: 0.12
Nodes (14): 1. The big picture, 7. Shutdown: steps 19 – 20, 8. Reference, 9. Data Reading Flow (Function Names Only), CNC Data Collection Flow, Events in `machine_events`, Log lines (in `logs/app.log`, search for the message), Main parts (+6 more)

### Community 22 - "HTTP Handlers Layer (2)"
Cohesion: 0.09
Nodes (13): ControllerHandler, pathID(), targetUserID(), RoleHandler, UserCellHandler, UserCityHandler, UserCountryHandler, UserLineHandler (+5 more)

### Community 23 - "UserSupervisorService"
Cohesion: 0.15
Nodes (10): CreateUserSupervisorRequest, UpdateUserSupervisorRequest, validateUserSupervisor(), UserSupervisor, UserSupervisorRepository, mapUserSupervisorError(), NewUserSupervisorRepository(), UserSupervisorService (+2 more)

### Community 24 - "CNC Go Client Run & Check Mode"
Cohesion: 0.14
Nodes (21): Prerequisite: 32-bit x86 DLL, compile Go with GOARCH=386, Prerequisite: SIEMENS_USER / SIEMENS_PASS env vars before CNCConnect, §3.1 Environment & Prerequisites, §6 Snapshot Data Structure, §10 Diagnostics & Troubleshooting, Siemens credentials: SIEMENS_USER / SIEMENS_PASS env vars, Siemens trust: cert must be trusted in NCU Access Administration, Struct size assertion: unsafe.Sizeof(CncMasterSnapshot{}) == 34440 (+13 more)

### Community 25 - "user_access_test.go"
Cohesion: 0.34
Nodes (3): LicenceKeyDuration, licenceKeyDurations(), body

### Community 27 - "validateCatalogFields"
Cohesion: 0.11
Nodes (11): CreateBrandRequest, UpdateBrandRequest, validateCatalogFields(), CreateCategoryRequest, UpdateCategoryRequest, CreateControllerRequest, UpdateControllerRequest, validateControllerFields() (+3 more)

### Community 28 - "ControllerService"
Cohesion: 0.21
Nodes (5): Controller, ControllerRepository, mapControllerError(), ControllerService, toControllerResponseError()

### Community 29 - "AuthService"
Cohesion: 0.22
Nodes (12): NewAuthHandler(), Authenticate(), WithAuthUser(), NewAuthRepository(), AuthService, NewAuthService(), asTestUser(), authRequest() (+4 more)

### Community 30 - "UserLineService"
Cohesion: 0.14
Nodes (10): CreateUserCellRequest, UpdateUserCellRequest, validateUserCell(), UserCell, UserCellRepository, mapUserCellError(), NewUserCellRepository(), UserCellService (+2 more)

### Community 31 - "LocationRepository"
Cohesion: 0.50
Nodes (3): mapLocationError(), NewLocationRepository(), LocationRepository

### Community 32 - "UserLocationService"
Cohesion: 0.16
Nodes (8): CreateUserLocationRequest, UpdateUserLocationRequest, validateUserLocation(), UserLocation, UserLocationRepository, mapUserLocationError(), UserLocationService, toUserLocationResponseError()

### Community 33 - "UserStateService"
Cohesion: 0.09
Nodes (13): CreateUserCountryRequest, CreateUserStateRequest, validateUserState(), UserCountry, UserState, UserCountryRepository, mapUserCountryError(), UserStateRepository (+5 more)

### Community 34 - "decodeAndValidate"
Cohesion: 0.13
Nodes (6): validator, BranchHandler, MachineHandler, currentUserID(), decodeAndValidate(), UserHandler

### Community 35 - "sinumerik_client.go"
Cohesion: 0.05
Nodes (36): main(), run(), main(), Config, CncMasterSnapshot, MachineDef, cstr(), enableVT() (+28 more)

### Community 36 - "time.Time"
Cohesion: 0.13
Nodes (14): 5. Connect and collect, per machine: steps 10 – 17, Known issue, Step 10 — Get ready to connect, Step 11 — Connect, a few machines at a time, Step 14 — The data collection loop, Step 16 — Every second: dead connection check and live status, FieldIndex(), TestFieldIndexAxisPositions() (+6 more)

### Community 37 - "PartService"
Cohesion: 0.17
Nodes (5): PartProcess, PartSubProgram, PartRepository, mapPartError(), saveProcesses()

### Community 38 - "CNC Connect API & Callbacks"
Cohesion: 0.31
Nodes (10): Package-level callback singletons (syscall.NewCallback), Mutex-protected copy: dataMutex sync.RWMutex -> lastRecord/hasRecord, onDataCallback(ptr uintptr) uintptr, onErrorCallback(code, messagePtr uintptr) uintptr, §7 Asynchronous Callback Architecture, CncMasterSnapshot (34,440 bytes, #pragma pack(8)), onDataCallback, onErrorCallback (+2 more)

### Community 39 - "Server API Domain Docs"
Cohesion: 0.23
Nodes (13): Clienteles API (/clienteles), Company Scoping (companyOwner / ownCompanyCustomer), Customer Display Names (customerName / companyName), Licence Keys API (/licence-keys), Machines API (/machines), Organization Hierarchy APIs (/user-country..supervisor), Parts API (/parts with processes, drawings), Production Schedules API (/production-schedules) (+5 more)

### Community 40 - ".Validate"
Cohesion: 0.14
Nodes (16): UpdateThemeRequest, isTimeZone(), isWebsite(), normalizeFilters(), normalizeTheme(), ParseDayStart(), TestCreateUserRequestValidate(), TestCreateUserRequestValidateDropsDuplicateMachines() (+8 more)

### Community 41 - "SendResponse"
Cohesion: 0.13
Nodes (5): BrandHandler, CategoryHandler, LicenceKeyHandler, ProductHandler, SendResponse()

### Community 42 - "IdleReasonService"
Cohesion: 0.15
Nodes (12): CreateIdleReasonRequest, IdleReasonFields, UpdateIdleReasonRequest, NewIdleReasonHandler(), IdleReason, IdleReasonRepository, mapIdleReasonError(), NewIdleReasonRepository() (+4 more)

### Community 43 - "2. All steps at a glance"
Cohesion: 0.15
Nodes (16): 2. All steps at a glance, 4. Find the machines: steps 5 – 9, 6. When data is edited: step 18, Step 18 — React to machine and variable changes, Step 20 — Stop the workers and mark machines offline, Step 5 — Prepare the database, Step 6 — Record the start and begin following the database, Step 7 — Listen for database changes (+8 more)

### Community 44 - "run"
Cohesion: 0.09
Nodes (52): run(), NewBranchHandler(), NewBrandHandler(), NewCategoryHandler(), NewControllerHandler(), NewMachineHandler(), NewMachineModelHandler(), NewProductHandler() (+44 more)

### Community 45 - ".Generate"
Cohesion: 0.40
Nodes (3): GenerateLicenceKeysRequest, GeneratedLicenceKeys, newLicenceKey()

### Community 46 - "7th Sense Server Engineering Principles (Skills Doc)"
Cohesion: 0.05
Nodes (22): Common Mistakes to Avoid, CRUD Implementation Checklist (11 steps), 7th Sense Server Engineering Principles (Skills Doc), Golden Reference Pattern: User CRUD, Unique Constraint Error Mapping (pgconn 23505), updated_at Trigger in Migrations, Air Hot Reload, 7th Sense Server - Claude Context (+14 more)

### Community 47 - "time.Time"
Cohesion: 0.10
Nodes (20): Value, MachineModeHistoryResponse, MachineModeTotal, HealthHandler, NewHealthHandler(), LiveValue, MachineEvent, MachineLiveStatus (+12 more)

### Community 49 - "VariableConfig"
Cohesion: 0.22
Nodes (7): VariableConfig, formatPreviousValue(), formatValue(), newValueChangeLog(), boundVariable, valueChange, valueChangeLog

### Community 50 - ".processReading"
Cohesion: 0.29
Nodes (10): Step 15 — Store the reading, decideMachineMode(), holdReason(), isInProduction(), machineStateFromReading(), programStateFromProcStatus(), TestMachineStateFromReading(), channelState (+2 more)

### Community 51 - "CncCommunication.dll (C++ core runtime, x86, OPC UA inside)"
Cohesion: 0.24
Nodes (16): client_cert.der/.pem (X.509 client cert for Siemens OPC UA), client_key.der/.pem (client private key), CncCommunication.dll (C++ core runtime, x86, OPC UA inside), Fanuc profile (engine 1, port 8193, Fwlib32.dll + fwlibe1.dll), Fwlib32.dll (Fanuc FOCAS2), fwlibe1.dll (Fanuc FOCAS2), libcrypto-3.dll (OpenSSL), libssl-3.dll (OpenSSL) (+8 more)

### Community 52 - ".Login"
Cohesion: 0.18
Nodes (7): LoginRequest, newSessionToken(), truncate(), newLoginThrottle(), loginFailures, loginThrottle, Session

### Community 53 - "runCatalogCRUD"
Cohesion: 0.31
Nodes (8): containsID(), runCatalogCRUD(), TestBrandCRUD(), TestCategoryCRUD(), TestControllerCRUD(), TestProductCRUD(), catalogCRUDCase, catalogRecord

### Community 55 - "CNC Event-Driven Data Flow"
Cohesion: 0.36
Nodes (8): Flow 1: load DLL, register two stdcall callbacks (syscall.NewCallback), Flow 2: CNCConnect starts DLL background worker thread, Flow 3: worker polls controller, assembles snapshot, invokes data callback, Flow 4: callback copies payload, signals renderer via buffered channel, Flow 5: main goroutine renders live ANSI dashboard, OPC UA (Siemens transport, port 4840), §2 System Architecture, SINUMERIK CNC (PLC/NCU)

### Community 56 - "Product"
Cohesion: 0.23
Nodes (6): Product, ProductRepository, mapProductError(), ProductService, productImageOrDefault(), toProductResponseError()

### Community 57 - "Check-mode OPC UA node mapping"
Cohesion: 0.19
Nodes (14): Check mode: go run sinumerik_client.go check (OPC UA node existence report), Go Client (sinumerik_client.go, callback -> copy -> render), Program status -> /Channel/State/acProg[u1,1], Operating mode -> /Bag/State/opMode[u1], Tool -> /Channel/State/actTNumber[u1,1], Check-mode OPC UA node mapping, Run: go run sinumerik_client.go (live dashboard), §8 Running the Go Client (+6 more)

### Community 59 - "Execute"
Cohesion: 0.16
Nodes (6): Execute(), UserListFilter, User, AuthRepository, mapAuthError(), mapUserError()

### Community 60 - "Category"
Cohesion: 0.25
Nodes (5): Category, CategoryRepository, mapCategoryError(), CategoryService, toCategoryResponseError()

### Community 64 - "MachineModelService"
Cohesion: 0.20
Nodes (6): MachineModelRequest, MachineModel, MachineModelRepository, mapMachineModelError(), MachineModelService, toMachineModelResponseError()

### Community 70 - "CncCommunication.dll API & Integration Guide (Go Client) - text version"
Cohesion: 0.22
Nodes (11): CNCConnect(ip, port, dataCb, errCb, machineType) int32, CNCConnect param dataCb (uintptr stdcall OnData(const CncMasterSnapshot*)), CNCConnect param errCb (uintptr stdcall OnErr(int32 code, const char* msg)), CNCConnect param ip (const char*, ANSI NUL-terminated), CNCConnect param port (int32, e.g. 4840), CNCdisconnect(machineId int32), CncCommunication.dll API & Integration Guide (Go Client) - text version, LazyDLL binding (syscall.NewLazyDLL + eager mod.Load in init) (+3 more)

### Community 71 - "MachineTelemetryService"
Cohesion: 0.19
Nodes (11): Supervisor, 3. Startup: steps 1 – 4, Step 1 — Start collection with the server, Step 2 — Check that collection is enabled, Step 3 — Load the DLL, Step 4 — Run the collector in the background, Driver, LoadDriver() (+3 more)

### Community 72 - "github.com/jackc/pgx/v5/pgxpool.Pool"
Cohesion: 0.11
Nodes (23): WithTransaction(), uniquePositiveIDs(), clockMinutes(), CreateShiftRequest, ShiftFields, UpdateShiftRequest, MinutesBetween(), validWorkingDays() (+15 more)

### Community 74 - "run"
Cohesion: 0.32
Nodes (11): Connect(), checkNotTaken(), createSuperadmin(), ensureProfile(), execute(), main(), parseFlags(), run() (+3 more)

### Community 75 - "BrandService"
Cohesion: 0.25
Nodes (5): Brand, BrandRepository, mapBrandError(), BrandService, toBrandResponseError()

### Community 76 - "BranchService"
Cohesion: 0.15
Nodes (9): CreateBranchRequest, UpdateBranchRequest, validateBranchFields(), Branch, BranchRepository, mapBranchError(), cityOwner(), BranchService (+1 more)

### Community 77 - "AuthHandler"
Cohesion: 0.38
Nodes (3): AuthHandler, CurrentUser(), UserID()

### Community 78 - "Address"
Cohesion: 0.24
Nodes (6): Address, AddressRepository, mapAddressError(), toAddressResponseError(), primaryAddress(), AddressService

### Community 79 - "CompanyProfile"
Cohesion: 0.28
Nodes (5): CompanyProfile, CompanyProfileRepository, mapCompanyProfileError(), toCompanyProfileResponseError(), CompanyProfileService

### Community 80 - "executeAffected"
Cohesion: 0.22
Nodes (4): ShiftBreak, BreakRepository, mapBreakError(), executeAffected()

### Community 81 - "Querier"
Cohesion: 0.27
Nodes (11): Querier, NewAddressRepository(), NewCompanyProfileRepository(), UserMachineAccessRepository, NewUserMachineAccessRepository(), NewAddressService(), NewCompanyProfileService(), newUserServices() (+3 more)

### Community 83 - "RoleService"
Cohesion: 0.21
Nodes (8): AssignableRoles(), CanAssignRole(), Role, RoleRepository, mapRoleError(), RoleService, toRoleResponseError(), NewUserService()

## Ambiguous Edges - Review These
- `Snapshot Fields (positions, spindle, feed, program, tool, alarms)` → `Machines API (/machines)`  [AMBIGUOUS]
  docs/CncCommunication.dll — API & Integration Guide (Go Client).pdf · relation: conceptually_related_to

## Knowledge Gaps
- **25 isolated node(s):** `hello`, `7th-sense-server`, `Siemens default ports: 4840 (native NCU) / 4841 (custom)`, `Engine 0: Simulator (local loopback mock engine)`, `Engine 4: Mitsubishi (driver not specified)` (+20 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 93 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **16 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `NewError()` connect `NewError` to `HTTP Handlers Layer`, `UserShopService`, `Machine DTOs`, `LicenceKeyService`, `DowntimeReasonService`, `UserProfileRepository`, `AuthUser`, `UserLineService`, `UserCityService`, `UserService`, `HTTP Handlers Layer (2)`, `UserSupervisorService`, `user_access_test.go`, `ControllerService`, `UserLineService`, `UserLocationService`, `UserStateService`, `decodeAndValidate`, `PartService`, `.Validate`, `IdleReasonService`, `.Generate`, `.Login`, `Product`, `Category`, `MachineModelService`, `github.com/jackc/pgx/v5/pgxpool.Pool`, `BrandService`, `BranchService`, `CompanyProfile`, `RoleService`?**
  _High betweenness centrality (0.064) - this node is a cross-community bridge._
- **What connects `hello`, `7th-sense-server`, `Siemens default ports: 4840 (native NCU) / 4841 (custom)` to the rest of the system?**
  _25 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `HTTP Handlers Layer` be split into smaller, more focused modules?**
  _Cohesion score 0.10144927536231885 - nodes in this community are weakly interconnected._
- **Why does `SendResponse()` connect `SendResponse` to `HTTP Handlers Layer`, `decodeAndValidate`, `Route Registration & Guards`, `AuthHandler`, `UserRepository`, `MachineModelHandler`, `HTTP Handlers Layer (2)`, `user_access_test.go`, `AuthService`?**
  _High betweenness centrality (0.027) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.10442890442890443 - nodes in this community are weakly interconnected._
- **Why does `CurrentUser()` connect `AuthHandler` to `HTTP Handlers Layer`, `Route Registration & Guards`, `Integration & DTO Tests (2)`, `context.Context`, `AuthUser`, `HTTP Handlers Layer (2)`?**
  _High betweenness centrality (0.022) - this node is a cross-community bridge._

### Low-confidence Hints
_AMBIGUOUS edges — the extractor was unsure. Verify before acting on these._

- **What is the exact relationship between `Snapshot Fields (positions, spindle, feed, program, tool, alarms)` and `Machines API (/machines)`?**
  _Edge tagged AMBIGUOUS (relation: conceptually_related_to) - confidence is low._