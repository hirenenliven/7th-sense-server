# Graph Report - 7th-sense-server  (2026-10-30)

## Corpus Check
- 65 files · ~131,906 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 2202 nodes · 15263 edges · 104 communities (85 shown, 19 thin omitted)
- Extraction: 94% EXTRACTED · 6% INFERRED · 0% AMBIGUOUS · INFERRED: 884 edges (avg confidence: 0.87)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- user_access_test.go
- AuthUser
- context.Context
- time.Time
- SendResponse()
- pathID()
- net/http.Request
- Execute()
- supervisor.go
- CncMasterSnapshot (34,440 bytes, #pragma
- ControllerService
- dll_windows_386.go
- UserLineService
- Collector
- worker
- MachineService
- server/main.go
- expectStatus()
- Handlers
- old-database-structure.sql (old MySQL ma
- net/http.ServeMux
- testing.T
- net/http.Handler
- §10 Diagnostics & Troubleshooting
- UserShopService
- machine_mode.go
- BranchService
- config.go
- UserCityService
- UserLocationService
- UserSupervisorService
- old table mm_machines
- Worker (internal/monitor/worker.go)
- .Update()
- MachineModelService
- NewError()
- LicenceKeyService
- github.com/jackc/pgx/v5/pgxpool.Pool
- old table tbl_users
- mode-chart.js
- CncCommunication.dll API & Integration G
- sinumerik_client.go
- UserProfileRepository
- CncCommunication.dll API & Integration G
- validateCatalogFields()
- .Login()
- Product
- drawChart()
- Machines API (/machines)
- 7th-sense-server/internal/model.AuthUser
- UserService
- DowntimeReasonService
- UserCountryService
- AddressRepository
- TestProductionScheduleCRUD()
- CNC Data Collection Flow
- BranchHandler
- WithTransaction()
- UserStateService
- CompanyProfileRepository
- PartService
- Querier
- validText()
- RoleRepository
- newMachineService()
- run()
- .validate()
- actAs()
- old table schedule_items
- CNCConnect(ip, port, dataCb, errCb, mach
- TestOrganizationForAnotherUser()
- onDataCallback(ptr uintptr) uintptr
- 7th Sense Server Install and Setup Guide
- Layered Architecture (Handler-Service-Re
- .validate()
- old table machines_modes
- Machine Mode Dashboard Page
- Client Certificate & Key (client_cert.de
- .validate()
- nc_status.go
- Event-driven data flow (no polling in cl
- .validate()
- .validate()
- PartDrawingStorage
- .Generate()
- CreateUserStateRequest
- auth_tokens SHA-256 Hashed Sessions (mig
- LicenceKeyDuration
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
1. `Execute()` - 172 edges
2. `NewError()` - 130 edges
3. `SendResponse()` - 122 edges
4. `old-database-structure.sql (old MySQL machine_monitor schema)` - 94 edges
5. `pathID()` - 68 edges
6. `Querier` - 58 edges
7. `AuthUser` - 51 edges
8. `decodeAndValidate()` - 50 edges
9. `old table mm_machines` - 45 edges
10. `worker` - 43 edges

## Surprising Connections (you probably didn't know these)
- `Machine Identifier Mapping (0 Simulator,1 Fanuc,2 Siemens,3 Syntec,4 Mitsubishi)` --semantically_similar_to--> `Machine Catalog Brand-Controller-Model`  [INFERRED] [semantically similar]
  docs/CncCommunication.dll — API & Integration Guide (Go Client).pdf → .claude/CLAUDE.md
- `Snapshot Fields (positions, spindle, feed, program, tool, alarms)` --conceptually_related_to--> `Machines API (/machines)`  [AMBIGUOUS]
  docs/CncCommunication.dll — API & Integration Guide (Go Client).pdf → .claude/CLAUDE.md
- `Step 15 — Store the reading` --references--> `with()`  [INFERRED]
  docs/cnc-data-collection-flow.md → tests/machine_test.go
- `2. All steps at a glance` --references--> `LoadDriver()`  [INFERRED]
  docs/cnc-data-collection-flow.md → internal/cnc/cnc.go
- `Step 10 — Get ready to connect` --references--> `FieldIndex()`  [INFERRED]
  docs/cnc-data-collection-flow.md → internal/cnc/fields.go

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
- **Old system cycle tracking (cycles, tool life, batches, mode seconds)** — migrations_old_database_structure_mm_cycletime, migrations_old_database_structure_mm_cycletime_history, migrations_old_database_structure_mm_tool_life, migrations_old_database_structure_batch_tbl, migrations_old_database_structure_machines_data, migrations_old_database_structure_machine_data_count [INFERRED 0.85]
- **Session authentication and authorization flow** — readme_session_cookie, readme_auth_tokens_hashed, readme_per_request_session_check, readme_login_lockout, readme_csrf_cross_origin, _claude_claude_role_guards [INFERRED 0.85]

## Communities (104 total, 19 thin omitted)

### Community 0 - "user_access_test.go"
Cohesion: 0.29
Nodes (5): LanguageResponse, ThemeResponse, authUserKey, body, machineModelRecord

### Community 1 - "AuthUser"
Cohesion: 0.07
Nodes (35): clockMinutes(), MinutesBetween(), AuthUser, Break, Clientele, Holiday, Shift, NewBreakRepository() (+27 more)

### Community 2 - "context.Context"
Cohesion: 0.06
Nodes (20): PartProcess, PartSubProgram, ShiftBreak, BreakRepository, mapBreakError(), ClienteleRepository, mapClienteleError(), NewClienteleRepository() (+12 more)

### Community 3 - "time.Time"
Cohesion: 0.06
Nodes (26): 4. Find the machines: steps 5 – 9, Step 3 — Run the collector in the background, Step 5 — Prepare the database, Step 7 — Listen for database changes, Step 8 — Load all active machines, ListenForDatabaseNotifications(), MachineModeHistoryResponse, MachineModesResponse (+18 more)

### Community 4 - "SendResponse()"
Cohesion: 0.09
Nodes (12): BreakHandler, ClienteleHandler, HolidayHandler, MachineHandler, PartHandler, ProductionScheduleHandler, currentUser(), listUserID() (+4 more)

### Community 5 - "pathID()"
Cohesion: 0.07
Nodes (19): validator, ControllerHandler, LicenceKeyHandler, decodeAndValidate(), pathID(), targetUserID(), UserCellHandler, NewUserCellHandler() (+11 more)

### Community 6 - "net/http.Request"
Cohesion: 0.09
Nodes (9): BrandHandler, CategoryHandler, DowntimeReasonHandler, IdleReasonHandler, MachineModelHandler, ProductHandler, UserCountryHandler, UserHandler (+1 more)

### Community 7 - "Execute()"
Cohesion: 0.08
Nodes (23): Execute(), UserCell, User, AuthRepository, mapAuthError(), NewAuthRepository(), UserCellRepository, mapUserCellError() (+15 more)

### Community 8 - "supervisor.go"
Cohesion: 0.08
Nodes (19): Status, Value, TestSnapshotModeFieldOffsets(), TestSnapshotSize(), FieldIndex(), TestFieldIndexAxisPositions(), VariableConfig, engineName() (+11 more)

### Community 9 - "CncMasterSnapshot (34,440 bytes, #pragma"
Cohesion: 0.10
Nodes (43): CncMasterSnapshot (34,440 bytes, #pragma pack(8), x86), AbsPos ([5]float64, Inches / Degrees, Work Coordinate System (WCS) positions), ActiveGCode ([16]byte, ASCII, Currently active G-code group), ActiveProgPath ([256]byte, ASCII Path, Storage path of active NC artifact), AlarmMsg ([128]byte, UTF-8 / ASCII, Alarm description message), AlarmNo (int32, Code (0 = OK), Primary alarm or fault number), AxisLoadXPercentage (int16, Percent (%), Servo X axis drive load (AxisLoadX/Y/ZPercentage)), AxisLoadYPercentage (int16, Percent (%), Servo Y axis drive load (AxisLoadX/Y/ZPercentage)) (+35 more)

### Community 10 - "ControllerService"
Cohesion: 0.10
Nodes (15): NewBrandHandler(), NewControllerHandler(), Brand, Controller, BrandRepository, mapBrandError(), NewBrandRepository(), ControllerRepository (+7 more)

### Community 11 - "dll_windows_386.go"
Cohesion: 0.09
Nodes (16): cncMasterSnapshot, dllDriver, slot, Step 12 — Call the DLL's CNCConnect, Step 13 — The DLL sends data, Connection, Driver, Reading (+8 more)

### Community 12 - "UserLineService"
Cohesion: 0.09
Nodes (16): CreateUserCellRequest, UpdateUserCellRequest, validateUserCell(), CreateUserLineRequest, UpdateUserLineRequest, validateUserLine(), NewUserLineHandler(), UserLine (+8 more)

### Community 13 - "Collector"
Cohesion: 0.13
Nodes (20): 2. All steps at a glance, 6. When data is edited: step 18, 7. Shutdown: steps 19 – 20, Log lines (in `logs/app.log`, search for the message), Step 18 — React to machine and variable changes, Step 19 — Stop collection with the server, Step 20 — Stop the workers and mark machines offline, Step 6 — Record the start and begin following the database (+12 more)

### Community 14 - "worker"
Cohesion: 0.17
Nodes (12): 5. Connect and collect, per machine: steps 10 – 17, Known issue, Retry and timing summary, Step 10 — Get ready to connect, Step 11 — Connect, a few machines at a time, Step 14 — The data collection loop, Step 15 — Store the reading, Step 16 — Every second: dead connection check and live status (+4 more)

### Community 15 - "MachineService"
Cohesion: 0.11
Nodes (19): CreateMachineRequest, MachineFields, UpdateMachineRequest, validCoordinate(), validHealthIDs(), validIPv4(), SeesAllCompanyMachines(), Machine (+11 more)

### Community 16 - "server/main.go"
Cohesion: 0.10
Nodes (10): main(), CncMasterSnapshot, MachineDef, cstr(), enableVT(), initDLL(), main(), render() (+2 more)

### Community 17 - "expectStatus()"
Cohesion: 0.17
Nodes (32): containsID(), decodeData(), expectStatus(), newCatalogRouter(), runCatalogCRUD(), sendRequest(), uniqueSuffix(), apiResponse (+24 more)

### Community 18 - "Handlers"
Cohesion: 0.06
Nodes (6): HealthHandler, NewHealthHandler(), registerHealthRoutes(), HealthService, NewHealthService(), Handlers

### Community 19 - "old-database-structure.sql (old MySQL ma"
Cohesion: 0.08
Nodes (31): old-database-structure.sql (old MySQL machine_monitor schema), old table branches, old table category, old table dropdown_key, old table dropdown_value, old table helpenter_pages, old table helpenter_pages_bkup, old table maintenance_tasks (+23 more)

### Community 20 - "net/http.ServeMux"
Cohesion: 0.13
Nodes (25): RoleHandler, NewRoleHandler(), registerBranchRoutes(), registerBrandRoutes(), registerBreakRoutes(), registerCategoryRoutes(), registerClienteleRoutes(), registerControllerRoutes() (+17 more)

### Community 21 - "testing.T"
Cohesion: 0.13
Nodes (22): UpdateThemeRequest, normalizeTheme(), TestCreateUserRequestValidate(), TestCreateUserRequestValidateDropsDuplicateMachines(), TestCreateUserRequestValidateNormalizes(), TestCustomerFiltersValidate(), TestUpdateThemeRequestValidate(), TestUpdateUserRequestValidateStatus() (+14 more)

### Community 22 - "net/http.Handler"
Cohesion: 0.14
Nodes (24): NewAuthHandler(), NewBreakHandler(), NewClienteleHandler(), NewHolidayHandler(), NewMachineHandler(), NewPartHandler(), NewProductionScheduleHandler(), NewShiftHandler() (+16 more)

### Community 23 - "§10 Diagnostics & Troubleshooting"
Cohesion: 0.12
Nodes (25): Check mode: go run sinumerik_client.go check (OPC UA node existence report), Go Client (sinumerik_client.go, callback -> copy -> render), Prerequisite: 32-bit x86 DLL, compile Go with GOARCH=386, Prerequisite: SIEMENS_USER / SIEMENS_PASS env vars before CNCConnect, Run: go run sinumerik_client.go (live dashboard), §3.1 Environment & Prerequisites, §8 Running the Go Client, §10 Diagnostics & Troubleshooting (+17 more)

### Community 24 - "UserShopService"
Cohesion: 0.14
Nodes (11): CreateUserShopRequest, UpdateUserShopRequest, validateUserShop(), NewUserShopHandler(), UserShop, UserShopRepository, mapUserShopError(), NewUserLineService() (+3 more)

### Community 25 - "machine_mode.go"
Cohesion: 0.22
Nodes (14): decideMachineMode(), holdReason(), machineStateFromReading(), numberingEngine(), programStateFromProcStatus(), ncReading(), TestDecideMachineMode(), TestDecideMachineModeConditions() (+6 more)

### Community 26 - "BranchService"
Cohesion: 0.14
Nodes (10): CreateBranchRequest, UpdateBranchRequest, validateBranchFields(), NewBranchHandler(), Branch, BranchRepository, mapBranchError(), cityOwner() (+2 more)

### Community 27 - "config.go"
Cohesion: 0.16
Nodes (14): main(), run(), Supervisor, Config, MachineSelection, NewSupervisor(), getEnv(), getEnvBool() (+6 more)

### Community 28 - "UserCityService"
Cohesion: 0.15
Nodes (9): CreateUserCityRequest, UpdateUserCityRequest, validateUserCity(), NewUserCityHandler(), UserCity, UserCityRepository, mapUserCityError(), UserCityService (+1 more)

### Community 29 - "UserLocationService"
Cohesion: 0.15
Nodes (9): CreateUserLocationRequest, UpdateUserLocationRequest, validateUserLocation(), NewUserLocationHandler(), UserLocation, UserLocationRepository, mapUserLocationError(), UserLocationService (+1 more)

### Community 30 - "UserSupervisorService"
Cohesion: 0.15
Nodes (10): CreateUserSupervisorRequest, UpdateUserSupervisorRequest, validateUserSupervisor(), NewUserSupervisorHandler(), UserSupervisor, UserSupervisorRepository, mapUserSupervisorError(), NewUserSupervisorRepository() (+2 more)

### Community 31 - "old table mm_machines"
Cohesion: 0.14
Nodes (23): old table brand, old table controller, old table downtime_reason, old table idle_reason, old table machine_accuracy, old table machine_alarms, old table machine_alarms_history, old table machine_lineno (+15 more)

### Community 32 - "Worker (internal/monitor/worker.go)"
Cohesion: 0.11
Nodes (16): CNC Data Collection (in-process), writeLog CNC Log Types and Switches, .env Configuration (internal/config), CncCommunication.dll (32-bit vendor DLL), Collector (internal/monitor/collector.go), CNC Driver (internal/cnc), Known Issue: nil live-value cache panic, machine_events (+8 more)

### Community 33 - ".Update()"
Cohesion: 0.18
Nodes (13): CustomerDetails, UpdateUserRequest, UserResponse, isTimeZone(), isWebsite(), normalizeFilters(), ParseDayStart(), TestParseDayStart() (+5 more)

### Community 34 - "MachineModelService"
Cohesion: 0.17
Nodes (9): MachineModelRequest, NewMachineModelHandler(), MachineModel, MachineModelRepository, mapMachineModelError(), NewMachineModelRepository(), MachineModelService, NewMachineModelService() (+1 more)

### Community 35 - "NewError()"
Cohesion: 0.14
Nodes (11): CreateAddressRequest, UpdateAddressRequest, UpdateLanguageRequest, UpdateLicenceKeyRequest, validDuration(), ValidLicenceKey(), NewError(), newResponse() (+3 more)

### Community 36 - "LicenceKeyService"
Cohesion: 0.16
Nodes (11): NewLicenceKey, CreateLicenceKeysRequest, NewLicenceKeyHandler(), LicenceKey, LicenceKeyRepository, mapLicenceKeyError(), NewLicenceKeyRepository(), LicenceKeyService (+3 more)

### Community 37 - "github.com/jackc/pgx/v5/pgxpool.Pool"
Cohesion: 0.16
Nodes (10): Connect(), NewCategoryHandler(), Category, CategoryRepository, mapCategoryError(), NewCategoryRepository(), NewControllerRepository(), CategoryService (+2 more)

### Community 38 - "old table tbl_users"
Cohesion: 0.11
Nodes (21): old table breaks, old table ci_sessions, old table clientele, old table downtime_items, old table mm_operator_tasks, old table module_permission, old table mst_downtime_reason, old table mst_downtime_reason_copy (+13 more)

### Community 39 - "mode-chart.js"
Cohesion: 0.22
Nodes (20): countByGroup(), currentSpan(), currentWindow(), drawDashboard(), drawMachineList(), drawTiles(), drawZoomLevels(), groupColor() (+12 more)

### Community 40 - "CncCommunication.dll API & Integration G"
Cohesion: 0.19
Nodes (20): CncCommunication.dll (C++ core runtime, x86, OPC UA inside), CncCommunication.dll API & Integration Guide (Go Client) - text version, Fanuc profile (engine 1, port 8193, Fwlib32.dll + fwlibe1.dll), Fwlib32.dll (Fanuc FOCAS2), fwlibe1.dll (Fanuc FOCAS2), libcrypto-3.dll (OpenSSL), libssl-3.dll (OpenSSL), Engine 1: Fanuc (FOCAS2 Ethernet, fwlib32) (+12 more)

### Community 41 - "sinumerik_client.go"
Cohesion: 0.20
Nodes (17): axisColor(), cstr(), enableVT(), findLibDir(), initDLL(), isRotary(), lbl(), logf() (+9 more)

### Community 42 - "UserProfileRepository"
Cohesion: 0.17
Nodes (10): UserProfileRepository, mapProfileError(), NewUserProfileRepository(), NewUserProfileService(), toUserProfileResponseError(), newUserServices(), UserProfileService, sameClock() (+2 more)

### Community 43 - "CncCommunication.dll API & Integration G"
Cohesion: 0.13
Nodes (18): CNCConnect param machineType (int32 engine id, see §4.2), Engine 0: Simulator (local loopback mock engine), Engine 2: Siemens (Sinumerik OPC UA, active profile), Engine 3: Syntec (Open API / TCP), Engine 4: Mitsubishi (driver not specified), Machine identifier mapping (engine id enum 0-4), §9 Controller-Specific Profiles, Siemens default ports: 4840 (native NCU) / 4841 (custom) (+10 more)

### Community 44 - "validateCatalogFields()"
Cohesion: 0.11
Nodes (11): CreateBrandRequest, UpdateBrandRequest, validateCatalogFields(), CreateCategoryRequest, UpdateCategoryRequest, CreateControllerRequest, UpdateControllerRequest, validateControllerFields() (+3 more)

### Community 45 - ".Login()"
Cohesion: 0.13
Nodes (9): LoginRequest, sleepUnlessCancelled(), newSessionToken(), truncate(), newLoginThrottle(), backoff, loginFailures, loginThrottle (+1 more)

### Community 46 - "Product"
Cohesion: 0.19
Nodes (9): NewProductHandler(), Product, ProductRepository, mapProductError(), NewProductRepository(), ProductService, NewProductService(), productImageOrDefault() (+1 more)

### Community 47 - "drawChart()"
Cohesion: 0.20
Nodes (19): cell(), chartData(), chartRows(), drawChart(), drawHeader(), drawPeriodList(), drawTotals(), escapeHtml() (+11 more)

### Community 48 - "Machines API (/machines)"
Cohesion: 0.16
Nodes (18): Clienteles API, Company Scope (companyOwner / ownCompanyCustomer), Current User API (/me theme, language), Customer Display Names (customerName / companyName), Holidays API, Licence Keys API, Machine Catalog Brand-Controller-Model, Machines API (/machines) (+10 more)

### Community 49 - "7th-sense-server/internal/model.AuthUser"
Cohesion: 0.25
Nodes (7): CreateIdleReasonRequest, IdleReasonFields, UpdateIdleReasonRequest, IdleReason, IdleReasonService, idleReasonFromFields(), toIdleReasonResponseError()

### Community 50 - "UserService"
Cohesion: 0.17
Nodes (6): UserListFilter, canSee(), UserService, NewUserService(), toUserResponseError(), userAccess

### Community 51 - "DowntimeReasonService"
Cohesion: 0.25
Nodes (7): CreateDowntimeReasonRequest, DowntimeReasonFields, UpdateDowntimeReasonRequest, DowntimeReason, downtimeReasonFromFields(), DowntimeReasonService, toDowntimeReasonResponseError()

### Community 52 - "UserCountryService"
Cohesion: 0.18
Nodes (7): CreateUserCountryRequest, NewUserCountryHandler(), UserCountry, UserCountryRepository, mapUserCountryError(), UserCountryService, toUserCountryResponseError()

### Community 53 - "AddressRepository"
Cohesion: 0.20
Nodes (9): Address, AddressRepository, mapAddressError(), NewAddressRepository(), NewAddressService(), toAddressResponseError(), AddressService, TestAddressCRUD() (+1 more)

### Community 54 - "TestProductionScheduleCRUD()"
Cohesion: 0.22
Nodes (17): PartDrawing, clienteleBody(), TestClienteleCRUD(), containsMachineID(), machineBody(), newMachineFixture(), TestMachineCRUD(), with() (+9 more)

### Community 55 - "CNC Data Collection Flow"
Cohesion: 0.12
Nodes (15): 1. The big picture, 3. Startup: steps 1 – 4, 8. Reference, 9. Data Reading Flow (Function Names Only), docs/cnc-communication.md, CNC Data Collection Flow, Events in `machine_events`, Main parts (+7 more)

### Community 56 - "BranchHandler"
Cohesion: 0.15
Nodes (7): AuthHandler, BranchHandler, currentUserID(), CurrentUser(), RequireRole(), UserID(), registerAuthRoutes()

### Community 57 - "WithTransaction()"
Cohesion: 0.27
Nodes (10): WithTransaction(), ProductionSchedule, ProductionScheduleItem, NewProductionScheduleRepository(), ProductionScheduleService, NewProductionScheduleService(), saveScheduleItems(), scheduleDetails() (+2 more)

### Community 58 - "UserStateService"
Cohesion: 0.21
Nodes (6): NewUserStateHandler(), UserState, UserStateRepository, mapUserStateError(), UserStateService, toUserStateResponseError()

### Community 59 - "CompanyProfileRepository"
Cohesion: 0.23
Nodes (8): CompanyProfile, CompanyProfileRepository, mapCompanyProfileError(), NewCompanyProfileRepository(), NewCompanyProfileService(), toCompanyProfileResponseError(), CompanyProfileService, TestCompanyProfileCRUD()

### Community 60 - "PartService"
Cohesion: 0.29
Nodes (8): ValidPartDrawingFile(), Part, NewPartRepository(), PartService, NewPartService(), partFromFields(), saveProcesses(), toPartResponseError()

### Community 61 - "Querier"
Cohesion: 0.19
Nodes (10): Querier, mapLocationError(), NewLocationRepository(), NewUserCellRepository(), NewUserLineRepository(), execute(), UserMachineAccessRepository, NewUserMachineAccessRepository() (+2 more)

### Community 62 - "validText()"
Cohesion: 0.21
Nodes (7): PartProcessFields, PartSubProgramFields, CreatePartRequest, PartFields, UpdatePartRequest, UploadPartDrawingRequest, validText()

### Community 63 - "RoleRepository"
Cohesion: 0.21
Nodes (9): AssignableRoles(), CanAssignRole(), Role, RoleRepository, mapRoleError(), NewRoleRepository(), RoleService, NewRoleService() (+1 more)

### Community 64 - "newMachineService()"
Cohesion: 0.37
Nodes (13): NewBranchRepository(), NewUserCityRepository(), NewUserCountryRepository(), NewUserLocationRepository(), NewUserStateRepository(), NewBranchService(), NewUserCityService(), NewUserCountryService() (+5 more)

### Community 65 - "run()"
Cohesion: 0.18
Nodes (8): run(), GenerateCertificates(), NewDowntimeReasonHandler(), NewIdleReasonHandler(), NewUserHandler(), NewDowntimeReasonService(), NewIdleReasonService(), userServices

### Community 66 - ".validate()"
Cohesion: 0.24
Nodes (8): validDate(), CreateProductionScheduleRequest, ProductionScheduleFields, ProductionScheduleItemFields, UpdateProductionScheduleRequest, ScheduleDay(), shiftSpan(), jobSpan

### Community 67 - "actAs()"
Cohesion: 0.24
Nodes (12): CreateUserRequest, WithAuthUser(), asTestUser(), actAs(), actAsMember(), TestAssignableRoles(), TestNobodyManagesTheSuperadmin(), testRoleID() (+4 more)

### Community 68 - "old table schedule_items"
Cohesion: 0.27
Nodes (12): old table batch_tbl, old table mm_cycletime, old table mm_cycletime_history, old table mm_daily_production, old table mm_daily_production_history, old table mm_tool_life, old table mm_tool_life_history, old table parts (+4 more)

### Community 69 - "CNCConnect(ip, port, dataCb, errCb, mach"
Cohesion: 0.20
Nodes (11): CNCConnect(ip, port, dataCb, errCb, machineType) int32, CNCConnect param dataCb (uintptr stdcall OnData(const CncMasterSnapshot*)), CNCConnect param ip (const char*, ANSI NUL-terminated), CNCConnect param port (int32, e.g. 4840), CNCdisconnect(machineId int32), LazyDLL binding (syscall.NewLazyDLL + eager mod.Load in init), §4 API Reference, §5 DLL Loading & Binding (Go) (+3 more)

### Community 70 - "TestOrganizationForAnotherUser()"
Cohesion: 0.27
Nodes (9): authRequest(), createAuthTestUser(), login(), sessionCookieOf(), TestAuthentication(), containsOrgID(), TestOrganizationForAnotherUser(), orgRecord (+1 more)

### Community 71 - "onDataCallback(ptr uintptr) uintptr"
Cohesion: 0.36
Nodes (9): Package-level callback singletons (syscall.NewCallback), CNCConnect param errCb (uintptr stdcall OnErr(int32 code, const char* msg)), Mutex-protected copy: dataMutex sync.RWMutex -> lastRecord/hasRecord, onDataCallback(ptr uintptr) uintptr, onErrorCallback(code, messagePtr uintptr) uintptr, §7 Asynchronous Callback Architecture, onDataCallback, onErrorCallback (+1 more)

### Community 72 - "7th Sense Server Install and Setup Guide"
Cohesion: 0.28
Nodes (7): 7th Sense Server Install and Setup Guide, FanucCncCommunication.dll, Migrations schema.sql + 024-029, PostgreSQL 18, Production Deploy (server.exe + Lib + certs), Siemens OPC UA Client Certs, TimescaleDB Setup

### Community 73 - "Layered Architecture (Handler-Service-Re"
Cohesion: 0.25
Nodes (4): 7th Sense Server Claude Context, response.SendResponse Format, User CRUD Golden Reference, README Layered Layout (Node.js analogy)

### Community 74 - ".validate()"
Cohesion: 0.32
Nodes (5): uniquePositiveIDs(), CreateShiftRequest, ShiftFields, UpdateShiftRequest, validWorkingDays()

### Community 75 - "old table machines_modes"
Cohesion: 0.43
Nodes (8): old table machine_data_count, old table machine_data_count_history, old table machine_data_count_hourly, old table machines_data, old table machines_data_bkp, old table machines_data_history, old table machines_data_temp, old table machines_modes

### Community 76 - "Machine Mode Dashboard Page"
Cohesion: 0.33
Nodes (7): Health Mode History API (public), machine_mode_log, ECharts 5.5.1, Mode List Table and Time per Mode, Mode Count Tiles (filter), Machine Mode Dashboard Page, Zoom Level Bar

### Community 77 - "Client Certificate & Key (client_cert.de"
Cohesion: 0.48
Nodes (7): client_cert.der/.pem (X.509 client cert for Siemens OPC UA), client_key.der/.pem (client private key), Siemens certs: client_cert.der + client_key.der in binary directory, Siemens trust: cert must be trusted in NCU Access Administration, Troubleshoot: BadRequestNotAllowed -> Export client_cert, BadRequestNotAllowed (certificate trust rejection), Client Certificate & Key (client_cert.der / client_key.der)

### Community 78 - ".validate()"
Cohesion: 0.38
Nodes (4): ClienteleFields, CreateClienteleRequest, UpdateClienteleRequest, validClienteleEmail()

### Community 79 - "nc_status.go"
Cohesion: 0.52
Nodes (5): engineOf(), ncStatusName(), opModeName(), procStatusName(), TestNCStatusNamesPerEngine()

### Community 80 - "Event-driven data flow (no polling in cl"
Cohesion: 0.60
Nodes (5): Flow 1: load DLL, register two stdcall callbacks (syscall.NewCallback), Flow 2: CNCConnect starts DLL background worker thread, Flow 3: worker polls controller, assembles snapshot, invokes data callback, Flow 4: callback copies payload, signals renderer via buffered channel, Flow 5: main goroutine renders live ANSI dashboard

### Community 81 - ".validate()"
Cohesion: 0.47
Nodes (3): BreakFields, CreateBreakRequest, UpdateBreakRequest

### Community 82 - ".validate()"
Cohesion: 0.47
Nodes (3): CreateHolidayRequest, HolidayFields, UpdateHolidayRequest

### Community 84 - ".Generate()"
Cohesion: 0.40
Nodes (3): GenerateLicenceKeysRequest, GeneratedLicenceKeys, newLicenceKey()

### Community 86 - "auth_tokens SHA-256 Hashed Sessions (mig"
Cohesion: 0.67
Nodes (3): auth_tokens SHA-256 Hashed Sessions (migration 014), Per-Request Session and User Re-check, ts_session Cookie (HttpOnly, SameSite=Strict)

## Ambiguous Edges - Review These
- `Snapshot Fields (positions, spindle, feed, program, tool, alarms)` → `Machines API (/machines)`  [AMBIGUOUS]
  docs/CncCommunication.dll — API & Integration Guide (Go Client).pdf · relation: conceptually_related_to
- `CncCommunication.dll (32-bit vendor DLL)` → `FanucCncCommunication.dll`  [AMBIGUOUS]
  docs/setup-guide.txt · relation: semantically_similar_to
- `old table machine_data_count` → `old table schedule_items`  [AMBIGUOUS]
  migrations/old-database-structure.sql · relation: shares_data_with
- `old table machine_data_count_history` → `old table schedule_items`  [AMBIGUOUS]
  migrations/old-database-structure.sql · relation: shares_data_with
- `old table machine_data_count_hourly` → `old table schedule_items`  [AMBIGUOUS]
  migrations/old-database-structure.sql · relation: shares_data_with
- `old table mm_cycletime` → `old table schedule_items`  [AMBIGUOUS]
  migrations/old-database-structure.sql · relation: shares_data_with
- `old table mm_cycletime_history` → `old table schedule_items`  [AMBIGUOUS]
  migrations/old-database-structure.sql · relation: shares_data_with
- `old table mm_daily_production` → `old table schedule_items`  [AMBIGUOUS]
  migrations/old-database-structure.sql · relation: shares_data_with
- `old table mm_daily_production_history` → `old table schedule_items`  [AMBIGUOUS]
  migrations/old-database-structure.sql · relation: shares_data_with
- `old table mm_tool_life` → `old table schedule_items`  [AMBIGUOUS]
  migrations/old-database-structure.sql · relation: shares_data_with
- `old table mm_tool_life_history` → `old table schedule_items`  [AMBIGUOUS]
  migrations/old-database-structure.sql · relation: shares_data_with

## Knowledge Gaps
- **54 isolated node(s):** `7th-sense-server`, `Siemens default ports: 4840 (native NCU) / 4841 (custom)`, `Engine 0: Simulator (local loopback mock engine)`, `Engine 4: Mitsubishi (driver not specified)`, `user_filters Filter Checklist` (+49 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 144 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **19 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `CNC Data Collection Flow` connect `CNC Data Collection Flow` to `Worker (internal/monitor/worker.go)`, `time.Time`, `7th Sense Server Install and Setup Guide`, `Collector`, `worker`?**
  _High betweenness centrality (0.175) - this node is a cross-community bridge._
- **What connects `7th-sense-server`, `Siemens default ports: 4840 (native NCU) / 4841 (custom)`, `Engine 0: Simulator (local loopback mock engine)` to the rest of the system?**
  _54 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `AuthUser` be split into smaller, more focused modules?**
  _Cohesion score 0.07315315315315316 - nodes in this community are weakly interconnected._
- **Why does `2. All steps at a glance` connect `Collector` to `time.Time`, `dll_windows_386.go`, `worker`, `CNC Data Collection Flow`?**
  _High betweenness centrality (0.153) - this node is a cross-community bridge._
- **Should `context.Context` be split into smaller, more focused modules?**
  _Cohesion score 0.05859969558599695 - nodes in this community are weakly interconnected._

### Low-confidence Hints
_AMBIGUOUS edges — the extractor was unsure. Verify before acting on these._

- **What is the exact relationship between `Snapshot Fields (positions, spindle, feed, program, tool, alarms)` and `Machines API (/machines)`?**
  _Edge tagged AMBIGUOUS (relation: conceptually_related_to) - confidence is low._
- **What is the exact relationship between `CncCommunication.dll (32-bit vendor DLL)` and `FanucCncCommunication.dll`?**
  _Edge tagged AMBIGUOUS (relation: semantically_similar_to) - confidence is low._