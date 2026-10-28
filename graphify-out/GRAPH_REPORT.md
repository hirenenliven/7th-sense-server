# Graph Report - 7th-sense-server  (2026-10-28)

## Corpus Check
- 251 files · ~131,200 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 4 file(s) not represented in the graph (top: (none) 3, .toml 1)

## Summary
- 2009 nodes · 16483 edges · 82 communities (64 shown, 18 thin omitted)
- Extraction: 95% EXTRACTED · 5% INFERRED · 0% AMBIGUOUS · INFERRED: 822 edges (avg confidence: 0.86)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `957d745e`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- HTTP Handlers Layer
- testing.T
- worker.go
- UserShopService
- Route Registration & Guards
- MachineService
- LicenceKeyService
- CNC Snapshot Struct Fields
- UserLineService
- UserSupervisorService
- sinumerik_client.go
- Integration & DTO Tests (2)
- CompanyProfile
- CncCommunication.dll API & Integration Guide (Go Client)
- HolidayService
- UserStateService
- UserCellService
- mode-chart.js
- UserCityService
- Connection
- UserService
- validateCatalogFields
- HTTP Handlers Layer (2)
- UserProfileRepository
- CNC Go Client Run & Check Mode
- LicenceKeyDuration
- .validate
- Address
- ControllerService
- NewError
- generate.go
- context.Context
- validText
- AuthUser
- .accessFor
- config.go
- Product
- .validate
- CNC Connect API & Callbacks
- Server API Domain Docs
- UserHandler
- SendResponse
- DowntimeReasonService
- 2. All steps at a glance
- run
- .validate
- 7th Sense Server Engineering Principles (Skills Doc)
- .Validate
- GET /health
- time.Time
- time.Time
- CncCommunication.dll (C++ core runtime, x86, OPC UA inside)
- .Login
- Querier
- CNC Event-Driven Data Flow
- .Validate
- Check-mode OPC UA node mapping
- .Generate
- .UploadDrawing
- machineState
- Sinumerik DLL Go Client (use case) (2)
- MachineModelService
- Go Module Root
- Current User Middleware
- Role Guard Middleware
- SendResponse Helper
- Test Auth Helpers
- CncCommunication.dll API & Integration Guide (Go Client) - text version
- CNCConfig
- CreateAddressRequest
- UpdateAddressRequest
- BranchService
- AuthHandler
- Execute
- WithTransaction
- UserCountryService
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
10. `worker` - 40 edges

## Surprising Connections (you probably didn't know these)
- `Machine Identifier Mapping (0 Simulator,1 Fanuc,2 Siemens,3 Syntec,4 Mitsubishi)` --semantically_similar_to--> `Machine Catalog Brand -> Controller -> Model`  [INFERRED] [semantically similar]
  docs/CncCommunication.dll — API & Integration Guide (Go Client).pdf → .claude/CLAUDE.md
- `Snapshot Fields (positions, spindle, feed, program, tool, alarms)` --conceptually_related_to--> `Machines API (/machines)`  [AMBIGUOUS]
  docs/CncCommunication.dll — API & Integration Guide (Go Client).pdf → .claude/CLAUDE.md
- `2. All steps at a glance` --references--> `LoadDriver()`  [INFERRED]
  docs/cnc-data-collection-flow.md → internal/cnc/cnc.go
- `Step 3 — Run the collector in the background` --references--> `Collector`  [INFERRED]
  docs/cnc-data-collection-flow.md → internal/monitor/collector.go
- `Step 15 — Store the reading` --references--> `machineStateFromReading()`  [INFERRED]
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

## Communities (82 total, 18 thin omitted)

### Community 0 - "HTTP Handlers Layer"
Cohesion: 0.10
Nodes (11): BreakHandler, ClienteleHandler, DowntimeReasonHandler, HolidayHandler, IdleReasonHandler, PartHandler, ProductionScheduleHandler, currentUser() (+3 more)

### Community 1 - "testing.T"
Cohesion: 0.07
Nodes (98): CreateUserRequest, NewMachineHandler(), Authenticate(), WithAuthUser(), CrossOrigin(), Logger(), Chain(), Middleware (+90 more)

### Community 2 - "worker.go"
Cohesion: 0.09
Nodes (12): CncMasterSnapshot, MachineDef, cstr(), enableVT(), initDLL(), main(), render(), TestSnapshotSize() (+4 more)

### Community 3 - "UserShopService"
Cohesion: 0.14
Nodes (10): CreateUserShopRequest, UpdateUserShopRequest, validateUserShop(), NewUserShopHandler(), UserShop, UserShopRepository, mapUserShopError(), NewUserShopRepository() (+2 more)

### Community 4 - "Route Registration & Guards"
Cohesion: 0.17
Nodes (30): RequireRole(), registerBranchRoutes(), registerBrandRoutes(), registerBreakRoutes(), registerCategoryRoutes(), registerClienteleRoutes(), registerControllerRoutes(), registerDowntimeReasonRoutes() (+22 more)

### Community 5 - "MachineService"
Cohesion: 0.12
Nodes (16): CreateMachineRequest, MachineFields, UpdateMachineRequest, validCoordinate(), validHealthIDs(), validIPv4(), Machine, MachineScope (+8 more)

### Community 6 - "LicenceKeyService"
Cohesion: 0.16
Nodes (11): NewLicenceKey, CreateLicenceKeysRequest, NewLicenceKeyHandler(), LicenceKey, LicenceKeyRepository, mapLicenceKeyError(), NewLicenceKeyRepository(), LicenceKeyService (+3 more)

### Community 7 - "CNC Snapshot Struct Fields"
Cohesion: 0.11
Nodes (38): CncMasterSnapshot (34,440 bytes, #pragma pack(8), x86), AbsPos ([5]float64, Inches / Degrees, Work Coordinate System (WCS) positions), ActiveGCode ([16]byte, ASCII, Currently active G-code group), ActiveProgPath ([256]byte, ASCII Path, Storage path of active NC artifact), AlarmMsg ([128]byte, UTF-8 / ASCII, Alarm description message), AlarmNo (int32, Code (0 = OK), Primary alarm or fault number), AxisLoadXPercentage (int16, Percent (%), Servo X axis drive load (AxisLoadX/Y/ZPercentage)), AxisLoadYPercentage (int16, Percent (%), Servo Y axis drive load (AxisLoadX/Y/ZPercentage)) (+30 more)

### Community 8 - "UserLineService"
Cohesion: 0.14
Nodes (11): CreateUserLineRequest, UpdateUserLineRequest, validateUserLine(), NewUserLineHandler(), UserLine, UserLineRepository, mapUserLineError(), NewUserLineRepository() (+3 more)

### Community 9 - "UserSupervisorService"
Cohesion: 0.14
Nodes (11): CreateUserSupervisorRequest, UpdateUserSupervisorRequest, validateUserSupervisor(), NewUserSupervisorHandler(), UserSupervisor, UserSupervisorRepository, mapUserSupervisorError(), NewUserSupervisorRepository() (+3 more)

### Community 10 - "sinumerik_client.go"
Cohesion: 0.20
Nodes (17): axisColor(), cstr(), enableVT(), findLibDir(), initDLL(), isRotary(), lbl(), logf() (+9 more)

### Community 11 - "Integration & DTO Tests (2)"
Cohesion: 0.31
Nodes (5): LanguageResponse, ThemeResponse, authUserKey, body, machineModelRecord

### Community 12 - "CompanyProfile"
Cohesion: 0.23
Nodes (8): CompanyProfile, CompanyProfileRepository, mapCompanyProfileError(), NewCompanyProfileRepository(), NewCompanyProfileService(), toCompanyProfileResponseError(), CompanyProfileService, TestCompanyProfileCRUD()

### Community 13 - "CncCommunication.dll API & Integration Guide (Go Client)"
Cohesion: 0.13
Nodes (21): Machine Catalog Brand -> Controller -> Model, CNCConnect param machineType (int32 engine id, see §4.2), Engine 0: Simulator (local loopback mock engine), Engine 2: Siemens (Sinumerik OPC UA, active profile), Engine 3: Syntec (Open API / TCP), Engine 4: Mitsubishi (driver not specified), Machine identifier mapping (engine id enum 0-4), §9 Controller-Specific Profiles (+13 more)

### Community 14 - "HolidayService"
Cohesion: 0.14
Nodes (13): CreateHolidayRequest, HolidayFields, UpdateHolidayRequest, NewHolidayHandler(), Holiday, HolidayRepository, mapHolidayError(), NewHolidayRepository() (+5 more)

### Community 15 - "UserStateService"
Cohesion: 0.16
Nodes (8): CreateUserStateRequest, validateUserState(), NewUserStateHandler(), UserState, UserStateRepository, mapUserStateError(), UserStateService, toUserStateResponseError()

### Community 16 - "UserCellService"
Cohesion: 0.14
Nodes (11): CreateUserCellRequest, UpdateUserCellRequest, validateUserCell(), NewUserCellHandler(), UserCell, UserCellRepository, mapUserCellError(), NewUserCellRepository() (+3 more)

### Community 17 - "mode-chart.js"
Cohesion: 0.23
Nodes (21): applyRange(), cell(), chartData(), chartRows(), currentSpan(), drawChart(), drawHeader(), drawPeriodList() (+13 more)

### Community 18 - "UserCityService"
Cohesion: 0.15
Nodes (9): CreateUserCityRequest, UpdateUserCityRequest, validateUserCity(), NewUserCityHandler(), UserCity, UserCityRepository, mapUserCityError(), UserCityService (+1 more)

### Community 19 - "Connection"
Cohesion: 0.10
Nodes (15): cncMasterSnapshot, dllDriver, slot, Step 12 — Call the DLL's CNCConnect, Step 13 — The DLL sends data, Connection, Reading, Stats (+7 more)

### Community 20 - "UserService"
Cohesion: 0.16
Nodes (15): CustomerDetails, UpdateUserRequest, UserListFilter, UserResponse, ParseDayStart(), TestParseDayStart(), NewUserHandler(), canSee() (+7 more)

### Community 21 - "validateCatalogFields"
Cohesion: 0.11
Nodes (10): CreateBrandRequest, UpdateBrandRequest, validateCatalogFields(), CreateCategoryRequest, UpdateCategoryRequest, CreateControllerRequest, UpdateControllerRequest, validateControllerFields() (+2 more)

### Community 22 - "HTTP Handlers Layer (2)"
Cohesion: 0.09
Nodes (13): validator, ControllerHandler, decodeAndValidate(), pathID(), targetUserID(), UserCellHandler, UserCityHandler, UserCountryHandler (+5 more)

### Community 23 - "UserProfileRepository"
Cohesion: 0.19
Nodes (5): UserProfile, UserProfileRepository, mapProfileError(), toUserProfileResponseError(), UserProfileService

### Community 24 - "CNC Go Client Run & Check Mode"
Cohesion: 0.14
Nodes (21): Prerequisite: 32-bit x86 DLL, compile Go with GOARCH=386, Prerequisite: SIEMENS_USER / SIEMENS_PASS env vars before CNCConnect, §3.1 Environment & Prerequisites, §6 Snapshot Data Structure, §10 Diagnostics & Troubleshooting, Siemens credentials: SIEMENS_USER / SIEMENS_PASS env vars, Siemens trust: cert must be trusted in NCU Access Administration, Struct size assertion: unsafe.Sizeof(CncMasterSnapshot{}) == 34440 (+13 more)

### Community 26 - ".validate"
Cohesion: 0.38
Nodes (4): ClienteleFields, CreateClienteleRequest, UpdateClienteleRequest, validClienteleEmail()

### Community 27 - "Address"
Cohesion: 0.23
Nodes (7): Address, AddressRepository, mapAddressError(), NewAddressService(), toAddressResponseError(), AddressService, TestAddressCRUD()

### Community 28 - "ControllerService"
Cohesion: 0.21
Nodes (5): Controller, ControllerRepository, mapControllerError(), ControllerService, toControllerResponseError()

### Community 29 - "NewError"
Cohesion: 0.16
Nodes (13): UpdateLanguageRequest, UpdateThemeRequest, isTimeZone(), isWebsite(), normalizeFilters(), normalizeTheme(), TestUpdateThemeRequestValidate(), TestUpdateUserRequestValidateStatus() (+5 more)

### Community 31 - "context.Context"
Cohesion: 0.06
Nodes (22): Category, MachineEvent, MachineLiveStatus, ModeChange, PartProcess, PartSubProgram, ProductionScheduleItem, BreakRepository (+14 more)

### Community 32 - "validText"
Cohesion: 0.20
Nodes (5): PartProcessFields, PartSubProgramFields, uniquePositiveIDs(), validText(), validWorkingDays()

### Community 33 - "AuthUser"
Cohesion: 0.07
Nodes (34): CreatePartRequest, PartFields, UpdatePartRequest, ValidPartDrawingFile(), NewIdleReasonHandler(), NewPartHandler(), AuthUser, Break (+26 more)

### Community 34 - ".accessFor"
Cohesion: 0.22
Nodes (5): AssignableRoles(), CanAssignRole(), SeesAllCompanyMachines(), asInvalidSelection(), userAccess

### Community 35 - "config.go"
Cohesion: 0.11
Nodes (16): main(), run(), main(), Config, MachineSelection, getEnv(), getEnvBool(), getEnvInt() (+8 more)

### Community 36 - "Product"
Cohesion: 0.21
Nodes (7): CreateProductRequest, Product, ProductRepository, mapProductError(), ProductService, productImageOrDefault(), toProductResponseError()

### Community 37 - ".validate"
Cohesion: 0.47
Nodes (3): BreakFields, CreateBreakRequest, UpdateBreakRequest

### Community 38 - "CNC Connect API & Callbacks"
Cohesion: 0.31
Nodes (10): Package-level callback singletons (syscall.NewCallback), Mutex-protected copy: dataMutex sync.RWMutex -> lastRecord/hasRecord, onDataCallback(ptr uintptr) uintptr, onErrorCallback(code, messagePtr uintptr) uintptr, §7 Asynchronous Callback Architecture, CncMasterSnapshot (34,440 bytes, #pragma pack(8)), onDataCallback, onErrorCallback (+2 more)

### Community 39 - "Server API Domain Docs"
Cohesion: 0.23
Nodes (13): Clienteles API (/clienteles), Company Scoping (companyOwner / ownCompanyCustomer), Customer Display Names (customerName / companyName), Licence Keys API (/licence-keys), Machines API (/machines), Organization Hierarchy APIs (/user-country..supervisor), Parts API (/parts with processes, drawings), Production Schedules API (/production-schedules) (+5 more)

### Community 40 - "UserHandler"
Cohesion: 0.18
Nodes (3): BranchHandler, currentUserID(), UserHandler

### Community 41 - "SendResponse"
Cohesion: 0.13
Nodes (5): BrandHandler, CategoryHandler, LicenceKeyHandler, ProductHandler, SendResponse()

### Community 42 - "DowntimeReasonService"
Cohesion: 0.15
Nodes (12): CreateDowntimeReasonRequest, DowntimeReasonFields, UpdateDowntimeReasonRequest, NewDowntimeReasonHandler(), DowntimeReason, DowntimeReasonRepository, mapDowntimeReasonError(), NewDowntimeReasonRepository() (+4 more)

### Community 43 - "2. All steps at a glance"
Cohesion: 0.08
Nodes (34): 1. The big picture, 2. All steps at a glance, 4. Find the machines: steps 5 – 9, 6. When data is edited: step 18, 7. Shutdown: steps 19 – 20, 8. Reference, 9. Data Reading Flow (Function Names Only), CNC Data Collection Flow (+26 more)

### Community 44 - "run"
Cohesion: 0.06
Nodes (55): run(), Connect(), CreateUserLocationRequest, UpdateUserLocationRequest, validateUserLocation(), NewBrandHandler(), NewBreakHandler(), NewCategoryHandler() (+47 more)

### Community 45 - ".validate"
Cohesion: 0.47
Nodes (3): CreateIdleReasonRequest, IdleReasonFields, UpdateIdleReasonRequest

### Community 46 - "7th Sense Server Engineering Principles (Skills Doc)"
Cohesion: 0.05
Nodes (22): Common Mistakes to Avoid, CRUD Implementation Checklist (11 steps), 7th Sense Server Engineering Principles (Skills Doc), Golden Reference Pattern: User CRUD, Unique Constraint Error Mapping (pgconn 23505), updated_at Trigger in Migrations, Air Hot Reload, 7th Sense Server - Claude Context (+14 more)

### Community 47 - ".Validate"
Cohesion: 0.53
Nodes (5): TestCreateUserRequestValidate(), TestCreateUserRequestValidateDropsDuplicateMachines(), TestCreateUserRequestValidateNormalizes(), TestCustomerFiltersValidate(), validCustomerDetails()

### Community 49 - "time.Time"
Cohesion: 0.08
Nodes (23): Step 5 — Prepare the database, Value, MachineModeHistoryResponse, MachineModeTotal, HealthHandler, NewHealthHandler(), LiveValue, MachineModePeriod (+15 more)

### Community 50 - "time.Time"
Cohesion: 0.17
Nodes (12): 5. Connect and collect, per machine: steps 10 – 17, Known issue, Retry and timing summary, Step 10 — Get ready to connect, Step 11 — Connect, a few machines at a time, Step 14 — The data collection loop, Step 15 — Store the reading, Step 16 — Every second: dead connection check and live status (+4 more)

### Community 51 - "CncCommunication.dll (C++ core runtime, x86, OPC UA inside)"
Cohesion: 0.24
Nodes (16): client_cert.der/.pem (X.509 client cert for Siemens OPC UA), client_key.der/.pem (client private key), CncCommunication.dll (C++ core runtime, x86, OPC UA inside), Fanuc profile (engine 1, port 8193, Fwlib32.dll + fwlibe1.dll), Fwlib32.dll (Fanuc FOCAS2), fwlibe1.dll (Fanuc FOCAS2), libcrypto-3.dll (OpenSSL), libssl-3.dll (OpenSSL) (+8 more)

### Community 52 - ".Login"
Cohesion: 0.14
Nodes (11): LoginRequest, NewAuthHandler(), AuthService, hashToken(), NewAuthService(), newSessionToken(), truncate(), newLoginThrottle() (+3 more)

### Community 53 - "Querier"
Cohesion: 0.18
Nodes (14): Querier, NewAddressRepository(), mapLocationError(), NewLocationRepository(), UserMachineAccessRepository, NewUserMachineAccessRepository(), NewUserProfileRepository(), NewUserProfileService() (+6 more)

### Community 55 - "CNC Event-Driven Data Flow"
Cohesion: 0.36
Nodes (8): Flow 1: load DLL, register two stdcall callbacks (syscall.NewCallback), Flow 2: CNCConnect starts DLL background worker thread, Flow 3: worker polls controller, assembles snapshot, invokes data callback, Flow 4: callback copies payload, signals renderer via buffered channel, Flow 5: main goroutine renders live ANSI dashboard, OPC UA (Siemens transport, port 4840), §2 System Architecture, SINUMERIK CNC (PLC/NCU)

### Community 56 - ".Validate"
Cohesion: 0.40
Nodes (3): UpdateLicenceKeyRequest, validDuration(), ValidLicenceKey()

### Community 57 - "Check-mode OPC UA node mapping"
Cohesion: 0.19
Nodes (14): Check mode: go run sinumerik_client.go check (OPC UA node existence report), Go Client (sinumerik_client.go, callback -> copy -> render), Program status -> /Channel/State/acProg[u1,1], Operating mode -> /Bag/State/opMode[u1], Tool -> /Channel/State/actTNumber[u1,1], Check-mode OPC UA node mapping, Run: go run sinumerik_client.go (live dashboard), §8 Running the Go Client (+6 more)

### Community 58 - ".Generate"
Cohesion: 0.40
Nodes (3): GenerateLicenceKeysRequest, GeneratedLicenceKeys, newLicenceKey()

### Community 60 - "machineState"
Cohesion: 0.29
Nodes (10): decideMachineMode(), holdReason(), isInProduction(), machineStateFromReading(), programStateFromProcStatus(), TestMachineStateFromReading(), modeName(), channelState (+2 more)

### Community 64 - "MachineModelService"
Cohesion: 0.20
Nodes (6): MachineModelRequest, MachineModel, MachineModelRepository, mapMachineModelError(), MachineModelService, toMachineModelResponseError()

### Community 70 - "CncCommunication.dll API & Integration Guide (Go Client) - text version"
Cohesion: 0.22
Nodes (11): CNCConnect(ip, port, dataCb, errCb, machineType) int32, CNCConnect param dataCb (uintptr stdcall OnData(const CncMasterSnapshot*)), CNCConnect param errCb (uintptr stdcall OnErr(int32 code, const char* msg)), CNCConnect param ip (const char*, ANSI NUL-terminated), CNCConnect param port (int32, e.g. 4840), CNCdisconnect(machineId int32), CncCommunication.dll API & Integration Guide (Go Client) - text version, LazyDLL binding (syscall.NewLazyDLL + eager mod.Load in init) (+3 more)

### Community 71 - "CNCConfig"
Cohesion: 0.20
Nodes (10): Supervisor, 3. Startup: steps 1 – 4, Step 1 — Start collection with the server, Step 2 — Load the DLL, Step 3 — Run the collector in the background, Driver, LoadDriver(), NewSupervisor() (+2 more)

### Community 76 - "BranchService"
Cohesion: 0.14
Nodes (10): CreateBranchRequest, UpdateBranchRequest, validateBranchFields(), NewBranchHandler(), Branch, BranchRepository, mapBranchError(), cityOwner() (+2 more)

### Community 77 - "AuthHandler"
Cohesion: 0.20
Nodes (5): AuthHandler, MachineHandler, CurrentUser(), UserID(), Handlers

### Community 78 - "Execute"
Cohesion: 0.07
Nodes (26): Execute(), Brand, User, AuthRepository, mapAuthError(), NewAuthRepository(), BrandRepository, mapBrandError() (+18 more)

### Community 80 - "WithTransaction"
Cohesion: 0.07
Nodes (40): WithTransaction(), validDate(), CreateProductionScheduleRequest, ProductionScheduleFields, ProductionScheduleItemFields, UpdateProductionScheduleRequest, ScheduleDay(), clockMinutes() (+32 more)

### Community 83 - "UserCountryService"
Cohesion: 0.18
Nodes (7): CreateUserCountryRequest, NewUserCountryHandler(), UserCountry, UserCountryRepository, mapUserCountryError(), UserCountryService, toUserCountryResponseError()

## Ambiguous Edges - Review These
- `Snapshot Fields (positions, spindle, feed, program, tool, alarms)` → `Machines API (/machines)`  [AMBIGUOUS]
  docs/CncCommunication.dll — API & Integration Guide (Go Client).pdf · relation: conceptually_related_to

## Knowledge Gaps
- **25 isolated node(s):** `hello`, `7th-sense-server`, `Siemens default ports: 4840 (native NCU) / 4841 (custom)`, `Engine 0: Simulator (local loopback mock engine)`, `Engine 4: Mitsubishi (driver not specified)` (+20 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 91 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **18 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `NewError()` connect `NewError` to `HTTP Handlers Layer`, `UserShopService`, `MachineService`, `LicenceKeyService`, `UserLineService`, `UserSupervisorService`, `Integration & DTO Tests (2)`, `CompanyProfile`, `HolidayService`, `UserStateService`, `UserCellService`, `UserCityService`, `UserService`, `HTTP Handlers Layer (2)`, `UserProfileRepository`, `.validate`, `ControllerService`, `context.Context`, `validText`, `AuthUser`, `.accessFor`, `Product`, `.validate`, `DowntimeReasonService`, `run`, `.validate`, `.Validate`, `.Login`, `.Validate`, `.Generate`, `.UploadDrawing`, `MachineModelService`, `CreateAddressRequest`, `UpdateAddressRequest`, `BranchService`, `Execute`, `WithTransaction`, `UserCountryService`?**
  _High betweenness centrality (0.064) - this node is a cross-community bridge._
- **What connects `hello`, `7th-sense-server`, `Siemens default ports: 4840 (native NCU) / 4841 (custom)` to the rest of the system?**
  _25 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `HTTP Handlers Layer` be split into smaller, more focused modules?**
  _Cohesion score 0.09979296066252588 - nodes in this community are weakly interconnected._
- **Why does `SendResponse()` connect `SendResponse` to `HTTP Handlers Layer`, `testing.T`, `Route Registration & Guards`, `UserHandler`, `Integration & DTO Tests (2)`, `AuthHandler`, `MachineModelHandler`, `HTTP Handlers Layer (2)`?**
  _High betweenness centrality (0.034) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.06556776556776557 - nodes in this community are weakly interconnected._
- **Why does `CurrentUser()` connect `AuthHandler` to `HTTP Handlers Layer`, `AuthUser`, `Route Registration & Guards`, `Integration & DTO Tests (2)`, `HTTP Handlers Layer (2)`, `context.Context`?**
  _High betweenness centrality (0.020) - this node is a cross-community bridge._

### Low-confidence Hints
_AMBIGUOUS edges — the extractor was unsure. Verify before acting on these._

- **What is the exact relationship between `Snapshot Fields (positions, spindle, feed, program, tool, alarms)` and `Machines API (/machines)`?**
  _Edge tagged AMBIGUOUS (relation: conceptually_related_to) - confidence is low._