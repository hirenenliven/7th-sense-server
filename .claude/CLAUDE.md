# 7th Sense Server - Claude Context

## Project Overview
Go-based REST API backend for the 7th Sense system. Uses PostgreSQL with pgx driver and TimescaleDB for time-series data, following a strict layered architecture pattern designed for scalability and maintainability.

## Tech Stack
- **Language**: Go 1.27.1
- **Database**: PostgreSQL + TimescaleDB
- **Driver**: pgx/v5
- **Architecture**: Layered (Handler → Service → Repository → Model)
- **Hot Reload**: Air (go.mod tool)
- **Key Dependencies**: jackc/pgx/v5, joho/godotenv, golang.org/x/crypto (bcrypt), lumberjack (logging)

## Directory Structure
```
7th-sense-server/
├── cmd/server/              # Application entry point (main.go)
├── internal/
│   ├── config/              # Environment configuration
│   ├── db/                  # Database connection utilities (pgx pool)
│   ├── dto/                 # Data Transfer Objects (request validation)
│   ├── handler/             # HTTP request handlers (controllers)
│   ├── logger/              # Logging utilities
│   ├── middleware/          # HTTP middleware (logging, recovery)
│   ├── model/               # Database models (structs)
│   ├── repository/          # SQL queries and data access
│   ├── response/            # API response definitions
│   ├── routes/              # Route registration
│   └── service/             # Business logic layer
├── migrations/              # Database migration files
├── tests/                   # Integration tests
├── web/                     # Static pages embedded in the server (mode-chart/)
├── logs/                    # Application logs
└── .air.toml               # Air hot reload config
```

## Architecture Pattern (STRICT)

### Layered Architecture
Request flow: `routes → middleware → handler → dto.Validate → service → repository → db`

**Layers (no skipping allowed):**
1. **Handler** (`internal/handler/`) - HTTP concerns, request parsing, response sending
2. **Service** (`internal/service/`) - Business logic, error mapping, orchestration
3. **Repository** (`internal/repository/`) - SQL queries only, no business logic
4. **Model** (`internal/model/`) - Database structs with db/json tags
5. **DTO** (`internal/dto/`) - Request validation via `Validate()` method

### Critical Rules
- **Never skip layers** - Handlers cannot call repositories directly
- **Concrete types over interfaces** - Use structs, not interfaces (YAGNI)
- **Error mapping** - Repository errors → Response errors in service layer
- **COALESCE for nullables** - Always use `COALESCE(column, '')` in SELECT
- **Soft delete** - Every table has `is_deleted` and `deleted_at`
- **Validation in DTOs only** - Never validate in handlers/services

## CRUD Implementation Pattern

**Golden Reference**: User CRUD (`user_handler.go`, `user_service.go`, `user_repository.go`) is the canonical implementation. Copy this pattern exactly for new entities.

### Steps for New Entity:
1. **Model** (`internal/model/entity.go`) - Struct with db/json tags
2. **DTO** (`internal/dto/entity_dto.go`) - Request structs with `Validate()` method
3. **Repository** (`internal/repository/entity_repository.go`) - SQL queries via `db.Execute()`
4. **Service** (`internal/service/entity_service.go`) - Business logic + error mapping
5. **Response Errors** (`internal/response/responses.go`) - Add entity-specific responses
6. **Handler** (`internal/handler/entity_handler.go`) - HTTP handlers using `response.SendResponse()`
7. **Routes** (`internal/routes/entity_routes.go`) - Register routes in ServeMux
8. **Wiring** (`cmd/server/main.go`) - Wire dependencies in main

### Repository Pattern
```go
func (r *EntityRepository) FindByID(ctx context.Context, id int64) (*model.Entity, error) {
    query := `SELECT id, field1, COALESCE(field2, '') AS field2, ... 
              FROM entities WHERE id = $1 AND is_deleted = false`
    return db.Execute(ctx, r.dbPool, query, id)
}
```

### Service Error Mapping
```go
func toEntityResponseError(err error) error {
    switch {
    case errors.Is(err, repository.ErrNotFound):
        return response.EntityNotFound
    case errors.Is(err, repository.ErrDuplicateField):
        return response.FieldExists
    default:
        return err  // Becomes 500
    }
}
```

### Handler Pattern
```go
// GetByID handles GET /api/v1/entities/{id}.
func (entityHandler *EntityHandler) GetByID(writer http.ResponseWriter, request *http.Request) {
    id, err := pathID(request)
    if err != nil {
        response.SendResponse(writer, request, err, nil)
        return
    }

    entity, err := entityHandler.service.GetByID(request.Context(), id)
    if err != nil {
        response.SendResponse(writer, request, err, nil)
        return
    }
    response.SendResponse(writer, request, response.Success, entity)
}
```

## Naming Conventions
- **Packages**: lowercase, singular (`handler`, `service`, not `handlers`)
- **Files**: snake_case, entity-prefixed (`user_handler.go`, `user_service.go`)
- **Types**: PascalCase (`UserHandler`, `UserService`, `UserRepository`)
- **Constructors**: `NewXxx` (`NewUserHandler`)
- **Errors**: `ErrXxx` (`ErrNotFound`, `ErrDuplicateField`)

## Coding Standards

These apply to all new and changed code. Match the file you are in; do not refactor existing code just to meet them.

### Naming
- Names say what a thing is or does: `customerID`, `machineScope`, `checkBreaksFit`. No abbreviations or single letters beyond the usual `ctx`, `err`, `id`, `i`.
- Receivers and HTTP parameters use full words, as the existing code does: `(userHandler *UserHandler)`, `(userService *UserService)`, `writer http.ResponseWriter, request *http.Request`.
- A function's name matches what it does. If it also saves, deletes or changes state, the name says so (`SaveTheme`, not `GetTheme`). When a function's behaviour changes, rename it.
- Use the same verbs per layer: repository `FindAll` / `FindByID` / `Create` / `Update` / `Delete`; service and handler `GetAll` / `GetByID` / `Create` / `Update` / `Delete`; checks are `checkXxx` / `isXxx` / `hasXxx`; error mapping is `toEntityResponseError` (service) and `mapEntityError` (repository).
- Use the domain's words consistently: customer, company, member, clientele, shift, break. Don't name one thing two ways (e.g. `client` vs `clientele`).
- DTOs are `CreateXxxRequest` / `UpdateXxxRequest` / `XxxResponse`; JSON fields are camelCase, DB columns snake_case.

### Functions
- One function, one job. If you need "and" to describe it, or it is too long to follow at a glance, extract a named private helper. Don't pull out single-use helpers that only hide two or three plain lines.
- Return early on errors (guard clauses); avoid deep nesting.
- Handle every error. Wrap with the operation name in repositories (`fmt.Errorf("find user: %w", err)`); map to `response.*` errors in services.
- Pass `ctx context.Context` first, then the acting user (`actor`) when the method needs one, then the inputs.
- Prefer plain code over clever code: no generics, interfaces, reflection or goroutines unless the problem needs them (concurrency lives in `internal/collector` / `internal/monitor`).

### Comments
- Every exported type and function gets a short doc comment that starts with its name and says what it does, e.g. `// GetByID handles GET /api/v1/brands/{id}.` (handlers name their route and query parameters).
- Inside functions, comment only the *why*: business rules, non-obvious decisions, multi-step workflows (transactions, keep / add / soft-delete of child rows), edge cases (shift past midnight, DLL limits).
- Don't comment what the code already says (`// Extract ID from path`, `// Call service`), don't leave commented-out code, and don't leave TODOs without an owner or reason.
- Update or delete a comment when the code it describes changes.

### Formatting and checks
- Run `gofmt` (or `goimports`) on every changed file; group imports as standard library, third party, then `7th-sense-server/internal/...`.
- `go build ./...` and `go vet ./...` must pass; run `go test ./tests/...` when you touch behaviour that tests cover.
- No magic numbers or strings: give a value used in more than one place, or a value whose meaning isn't obvious, a named constant (e.g. `dto.LicenceKeyTrialDays`).

### Scope of changes
- Change only what the task needs. Don't rename, reformat or restructure unrelated code in the same change.
- Never change architecture, API contracts, database schema or business rules as part of a "cleanup". These need their own task.
- When you change an endpoint, table or rule, update the matching section of this file in the same change.

## Response Format
All responses use `response.SendResponse(writer, request, result, data)`. `result` is a `response.*` value (`response.Success`, `response.EntityNotFound`, ...) or an error; any other error becomes a 500:
```json
{"statusCode": 200, "success": true, "message": "success", "data": {...}}
{"statusCode": 404, "success": false, "message": "Entity not found"}
```

## Common Commands

### Development
```bash
go tool air                    # Hot reload (air config in .air.toml, GOARCH=386)
.\run.ps1                      # Run without hot reload (GOARCH=386, for the 32-bit CNC DLL)
.\build.ps1                    # Build bin\server.exe (GOARCH=386)
```

### Testing
```bash
go test ./tests/...            # Run all tests
go test -v ./tests/...         # Verbose output
```

### Database
```bash
psql -U postgres -d 7th_sense -f migrations/schema.sql
```

### Dependencies
```bash
go mod download                # Download dependencies
go mod tidy                    # Clean up go.mod
```

## API Endpoints

Base path: `/api/v1`

**Users**: GET/POST/PUT/DELETE `/users` — company-scoped (`service.UserService`). A customer is a company; its members (admin, supervisor, operator, qc, shop_head) have `users.company_id` = the customer. The superadmin manages every user (`?filterRole=customer`, `?companyId=` filter the list; `filterRole` is the listed users' role, never the caller's) and names a member's `companyId`; a customer or its admins see and change only their company's members (`managersOnly` route guard, others get 403; other companies' users answer 404). Roles a user may give: `model.AssignableRoles` — superadmin: customer + member roles; customer: member roles; admin: member roles but admin; nobody gives superadmin. Body `machineIds`: machines (of the company) a supervisor/operator/qc/shop_head may access, stored in `user_machine_access` (migration 015; no soft delete — removing an assignment deletes the row); admins and customers see every company machine and store none. Responses include `role` (name) and `machineIds`.
**Roles**: GET `/roles` — the roles the current user may give (role dropdown of the user form); empty for users who manage no users.
**Items**: GET/POST/PUT/DELETE `/items`
**Branches**: GET/POST/PUT/DELETE `/branches`
**Brands**: GET/POST/PUT/DELETE `/brands` — every logged-in user reads; only the superadmin writes (`superAdminOnly`).
**Categories**: GET/POST/PUT/DELETE `/categories`
**Products**: GET/POST/PUT/DELETE `/products`
**Controllers**: GET/POST/PUT/DELETE `/controllers` — machine catalog Brand → Controller → Model. Body `{brandId, name, description}`: no code or variant, the variant is part of the name (e.g. "Sinumerik 828D", migration 017). Brand required, name unique per brand (case-insensitive), `?brandId=` on list. Cannot be deleted while it has models or machines use it; its brand cannot change while machines use it.
**Models**: GET/POST/PUT/DELETE `/models` (`models` table, `MachineModel` in Go — `model` is the package name) — `{controllerId, name}`, name unique per controller (case-insensitive), `?controllerId=` on list. Cannot be deleted, or moved to another controller, while machines use it. Brands, controllers and models: every logged-in user reads (machine form dropdowns); only the superadmin writes (`superAdminOnly` in `internal/routes`).
**Addresses**: GET/POST/PUT/DELETE `/addresses`
**Catalog**: GET/POST `/catalog`
**Auth**: POST `/auth/login`, POST `/auth/logout` (public), GET `/auth/me` — HttpOnly `ts_session` cookie, hashed sessions in `auth_tokens` (see README › Authentication)
**Machines**: GET/POST/PUT/DELETE `/machines` — each machine belongs to a customer (`user_id`) and has a brand, a controller of that brand and a model (`model_id`) of that controller; `product_id` is no longer used. The superadmin reads all (`?userId=` narrows to one customer); add / edit are for managers (`managersOnly`: superadmin, customer, admin — a customer or admin always on its own company, whatever `userId` says), delete is superadmin-only; a customer and its admins read every machine of the company, other members only the machines assigned to them (`user_machine_access`, see `MachineService.machineScope`); anything else answers 404. The password is write-only (never returned; blank on PUT keeps it). Form-only writes: runtime columns (online status, mode, …) are not changed through the API. Accuracy / health options are fixed lists in the client (`constants/machines.js`).
**Parts**: GET/POST/PUT/DELETE `/parts` (migration 016: `parts`, `part_processes`, `part_process_machines`, `part_sub_programs`) — a part belongs to a customer (`user_id`) and is saved with its processes in one transaction (`service.PartService`). Body `{userId, clienteleId, partNo, partName, partCount, processes: [{id?, process, programName, machineIds, drawingNo, drawingRevisionNo, drawingFile, drawingName, accuracyId, setupTime, cycleTime, loadingTime, hasSubProgram, subPrograms: [{name, cycleTime}]}]}`, times in seconds. On PUT a process with `id` is updated, one without is added, saved processes left out are soft-deleted; machines and sub programs of a process are replaced. Production times are computed, never sent: process = setup + cycle + loading + sub program cycle times, part = sum of processes. Applicable machines must be the customer's; Part No. + name unique per customer (case-insensitive). `clienteleId` is optional and must be one of the customer's clienteles. Everyone reads their company's parts (superadmin all, `?userId=` narrows); managers (`managersOnly`) write their own company's. Drawings: POST `/parts/drawings` (multipart `file`, PDF ≤ 10 MB, managers) stores `UPLOAD_DIR/part-drawings/<32 hex>.pdf` and returns `{file, name}` for a process's `drawingFile`/`drawingName`; GET `/parts/{id}/processes/{processId}/drawing` streams it to anyone who may read the part.
**Clienteles**: GET/POST/PUT/DELETE `/clienteles` (migration 017) — the clients a customer makes parts for (`service.ClienteleService`). Body `{userId, clientCode, name, email, contactNumber, address, countryName, stateName, city}` (+ `status` on PUT, blank keeps): Client ID (`clientCode`) unique per customer (case-insensitive), email optional, `contactNumber` as `+<code><number>`, country and state as names ("India", "Gujarat"; migration 018). Same access as parts: everyone reads their company's (superadmin all, `?userId=` narrows), managers write their own company's. A clientele used by parts cannot be deleted or moved to another customer. Parts and clienteles share `companyOwner` / `ownCompanyCustomer` (`service/company_scope.go`).
**Shifts**: GET/POST/PUT/DELETE `/shifts` (migrations 019–021) — a customer's shifts (`service.ShiftService`). Body `{userId, name, shiftType: day|night, startTime, endTime, workingDays: [1..7], breakIds}` (+ `status` on PUT, blank keeps), times `"HH:MM"` 24 h; an end not after the start runs past midnight. `breakIds` are the customer's breaks the shift takes: saving sets their `breaks.shift_id` and releases (NULL) the shift's breaks left out; deleting the shift releases all of them. A break taken by another shift is refused; the breaks must lie inside the shift and not overlap (`checkBreaksFit`, minutes from the shift start). Responses add `breaks` (ordered from the shift start), `durationMinutes`, `breakMinutes` (computed, not stored). Name unique per customer (case-insensitive).
**Breaks**: GET/POST/PUT/DELETE `/breaks` (migrations 020–021, `service.BreakService`) — a customer's breaks `{userId, shiftId?, name, startTime, endTime}` (+ `status`), name unique per customer. `shiftId` (optional, one of the customer's shifts) is the shift that takes it; responses add `shiftName`. A break on a shift must fit it beside the shift's other breaks. Deleting a break takes it off its shift.
**Holidays**: GET/POST/PUT/DELETE `/holidays` (migration 020, `service.HolidayService`) — a customer's holidays `{holidayType: "holiday", name, startDate, endDate}` and week offs `{holidayType: "week_off", weekDay: 1..7, name?}` (+ `status`); the fields the kind does not use are cleared. A weekday is a customer's week off once; the same holiday name and start date once. Shifts, breaks and holidays have the same access as parts and clienteles (`companyOwner`).
**Downtime / idle reasons**: GET/POST/PUT/DELETE `/downtime-reasons`, `/idle-reasons` (migration 027: `downtime_reasons`, `idle_reasons`; replace the old `mst_downtime_reason` / `mst_idle_reason`; one set of files per table, `service.DowntimeReasonService` / `IdleReasonService`) — a customer's reasons `{userId, reasonCode, reason}` (+ `status` on PUT, blank keeps). `reasonCode` is optional (≤ 10, NULL when blank) and unique per customer when given; `reason` (≤ 50) unique per customer (both case-insensitive). Same access as clienteles (`companyOwner`).
**Licence keys**: GET/POST/PUT/DELETE `/licence-keys` (migration 024, `service.LicenceKeyService`) — keys the superadmin generates for a customer, format `K7M2Q-XR9TB-4HWZP-N8CDF` (4 groups of 5 from `dto.LicenceKeyAlphabet`, no I/O/0/1, `crypto/rand`). POST `/licence-keys/generate {count: 1..100}` returns `{licenceKeys}` without saving (the "Generate Key" button); POST `/licence-keys {userId, licenceKeys: [{licenceKey, duration}]}` saves them in one statement, all or none (a key already saved refuses the batch). `duration` is a subscription plan code (`dto.LicenceKeyDurations`, migration 025: M1…M12, Y1…Y3, LIFETIME, TRIAL = `LicenceKeyTrialDays` 30 days, TRIAL_1Y); the server stores `expires_on` (last valid day, counted from the day the key is saved; NULL = lifetime) and responses add `expiresOn`, `isExpired` (computed). PUT `{userId, duration, status}` changes the customer / plan / status (the end date is recounted from `created_at`; the key never changes). Generate, save, edit and delete are `superAdminOnly`; managers read their company's (`companyOwner`, `?userId=` narrows for the superadmin). Key unique among keys not deleted.
**Production schedules**: GET/POST/PUT/DELETE `/production-schedules` (migration 026: `production_schedules`, `production_schedule_items`; replaces the old `production_schedule` / `schedule_items`, `service.ProductionScheduleService`) — a customer's machine planned for one day (`scheduleType: single`, `endDate` = `startDate`) or a range (`multiple`), saved with its jobs in one transaction. Body `{userId, machineId, scheduleType, startDate, endDate, items: [{id?, scheduleDate, jobId, supervisorId, operatorId, shiftId, partId, processId, quantity, startTime, endTime}]}` (+ `status` on PUT), dates `YYYY-MM-DD`, times `"HH:MM"`; single day sets every `scheduleDate` to `startDate`. Items are kept / added / soft-deleted like part processes; `totalMinutes` is computed. Rules: machine, shift, part are the customer's; supervisor / operator are company members with that role; the process is the part's and lists the machine among its applicable machines; each job lies inside its shift (`scheduleDate` = the day the shift starts, so 02:00 in a 22:00–06:00 shift is next morning); jobs do not overlap; a machine has one schedule per day. Same access as parts. Responses add `customerName`, `machineName`, `itemCount`, `totalQuantity`; items add `supervisorName`, `operatorName`, `shiftName`, `partName`, `processName`.
**Health (mode history)**: GET `/health/machines/{id}/modes` (`service.HealthService`) — the machine's `machine_mode_log` periods of the last `dto.MachineModeHistoryDays` (10) days, oldest first: every period still running at or after the window start (one that began earlier is returned whole). Response `{machineId, machineName, currentModeId, onlineStatus, from, to, periods: [{id, modeId, modeName, programName, startedAt, endedAt, durationSec}], totals: [{modeId, modeName, periodCount, durationSec}]}`; the open period has `endedAt` null and `durationSec` up to now. **Public, no session** (registered on the outer router, `registerHealthRoutes`): any machine not deleted, others 404. GET `/health/machines` (public too, `HealthService.GetMachineModes`) — every machine not deleted with its current mode: `{at, machines: [{machineId, machineName, customerName, companyName, modeId, modeName, onlineStatus, programName, modeSince, modeDurationSec}]}` from the open `machine_mode_log` period (`machines.mode_id` when none; `modeId` null = never collected), ordered by name. Chart page: GET `/mode-chart/` (`?machineId=` selects one) serves `web/mode-chart/index.html` + `mode-chart.js`, embedded with `go:embed` (`web.ModeChart`) — a dashboard on the left (a count tile per mode that filters, search, machines grouped by current mode, refreshed every 30 s) and, for the selected machine, a timeline per mode with zoom level buttons (`ZOOM_LEVELS` in `mode-chart.js`: 30s, 1m, 2m, 5m, 10m, 30m, 1h, 3h, 8h, 24h, 3d, 10d; a level ends where the current view ends, the matching level is highlighted, 30 s is the closest zoom) (ECharts 5.5.1 from jsDelivr) and the list of periods; mode colours are in `mode-chart.js`. Not the `GET /health` liveness check.
**CNC data collection** (`docs/cnc-telemetry-collector-plan.md`, migration `028_cnc_machine_telementry.sql`; read API: Health above):
- **What starts it:** the API server itself, in-process (`internal/collector.Supervisor` runs `monitor.Collector` in a goroutine; no separate process or service). `CNC_ENABLED` defaults to `true`; `false` = API only. The DLL is 32-bit, so the server must be built `GOARCH=386`: `.\run.ps1`, `.\build.ps1` (→ `bin\server.exe`) and `go tool air` do. A 64-bit build refuses to start while `CNC_ENABLED=true` (`collector.ErrWrongBuild`); with `CNC_ENABLED=false` it serves the API only.
- **Machines:** the server connects every machine (shard 0 of 1, at most `CNC_SLOTS` ≤ 990) with `status = 'active'`, not deleted, through `CncCommunication.dll` in `Lib\` (`CNC_DLL_NAME`; the older `FanucCncCommunication.dll` reads no Siemens alarm node, so `AlarmNo` stays 0 there) (`internal/cnc`: `CNCConnect`, one callback slot per machine, Siemens login via `SIEMENS_USER`/`SIEMENS_PASS` under a lock). The engine comes from `controllers.cnc_engine` (NULL = `config_error` event, not connected).
- **Variables:** `cnc_variables` (old `mm_variables` ids), `cnc_variable_addresses` (source `snapshot` = a DLL snapshot field, the only one readable today; `plc`/`nc` wait for a DLL read function), `machine_variables` (per machine, `id` = variable id, key `(machine_id, id)`, hard delete). Axis positions of X, Y, Z (MachPos / AbsPos / RelPos / DistToGo, ids 208–219, addresses `MachPos[0]` = X … `[2]` = Z) come from migration 029; the code reads all 5 snapshot axes, a variable decides which are stored.
- **Flow:** `internal/monitor`: one worker per machine loads live values from `machine_variable_values` (numeric only), upserts every reading to detect changes, and inserts **only changed numeric values** into time-series (`MachineTelemetryService.UpsertLiveValues` → `InsertTelemetry`; no batching, but change detection is used: `store_mode` / deadband / `min_interval_ms` are not used). Text variables are not stored at all. Decides the `machine_modes` id on **every connected reading** from the reading's NC state, not from the stored variables (`internal/monitor/machine_mode.go`: `machineStateFromReading` → `decideMachineMode`, first match wins): no connection → Offline 9; alarm (`AlarmNo` ≠ 0 or `IsEmergency`) and no program running → Downtime 6; AUTO (`MachineMode` 2) + running/waiting + chanStatus active + feed override (`FeedRatePercetnage`) > 0 + no alarm → In Production 3; program stopped/interrupted, chanStatus interrupted, or running with feed override ≤ 0, an alarm or not AUTO → Hold 4; else Idle 5 (JOG, reset). `ProcStatus` is acProg (0 reset, 1 stopped, 2 active, 3 waiting, 4 interrupted — `programStateFromProcStatus`); chanStatus is not in the snapshot (`channelUnknown`, decides nothing). Setup / Rapid Mode are no longer set. Lost connection (`worker.checkConnectionLost`): the DLL stops calling back while nothing changes (a held / idle machine is silent), so silence is not Offline — Offline comes from an `IsConnected=0` reading, a failed connect, or no connected reading for `CNC_STALE_SEC` (30) while the DLL's last word was "not connected"; after a connected reading, only no callback for `CNC_QUIET_SEC` (300) reopens the connection, keeping the mode (`stale` event "mode kept"). At start every machine is set Offline (`MarkShardOffline`) until its first connected reading decides its mode. Mode changes, events and statuses are written directly too. Logs: `cnc collector: worker started`, `cnc: variables loaded` (variable → field map) and `cnc: active variables not readable` (FindVariableConfig returns those with an empty address), `cnc: machine subscribed` on connect, `cnc: controller not sending data, no values read` once per connection while the DLL reports `IsConnected=0` (no values are read or logged then), `cnc: controller data started` / `stopped`, `cnc: reading rate` every 10 s per machine (callbacks / processed / skipped / stored). Every collector log line goes through one function, `writeLog(settings, logType, level, machine, message, attributes...)` in `internal/monitor/cnc_log.go` (called as `collector.log` / `collector.logMachine` / `machineWorker.log`; no direct `slog` calls in `internal/monitor`): log types `logTypeEvent` / `logTypeMode` / `logTypeModeCheck` / `logTypeVariable` / `logTypeSnapshot` / `logTypeOther`, levels `slog.LevelInfo` / `LevelWarn` / `LevelError`; a line about a machine starts with `[machine <id> <name>]` and has `machineId` + `machine`, every line has `logType`; errors are always written, each type by its switch (`CNC_LOG_OTHER` for `other`; `cnc collector: settings` is always written): `CNC_LOG_EVENTS` → `[machine …] cnc machine event <event>` (failures as Warn); `CNC_LOG_MODES` → `[machine …] cnc: mode changed <from> -> <to>` (`reason` = the rule that matched + NC inputs `opMode`, `procStatus`, `chanStatus`, `feedOverride`, `alarmNo`, `emergency`, `ncProgram`) ; `CNC_LOG_MODE_CHECKS` → `[machine …] cnc: mode check <mode>` whenever the NC inputs change (`decided`, `current`, `modeChanges`); `CNC_LOG_SNAPSHOT` (default false) → `[machine …] cnc: snapshot status` whenever a raw DLL status field changes (DLL names: `IsConnected`, `ResultStatus`, `MachineMode`, `ProcStatus`, `IsEmergency`, `FeedRatePercetnage`, `SpindeOverridePercentage`, `CmdFeedOvr`, `CmdSpOvr`, `AlarmNo`, `AlarmMsg`, `OpMsgNo`, `OpMsg`, `ActiveGCode`, `ProgName`, `missingNodes`; `cnc.Reading.Status`); `CNC_LOG_VARIABLES` → `[machine …] cnc: variable values changed` once per reading that changed values (`count` + `changes: ["208:X_MachPos 12.5 -> 12.75", …]`, previous `none` for a first value). `cnc collector: settings` logs the machines and switches at start. `CNC_MACHINE_IDS=225,247` makes the collector connect only those machines (`collector.go`: the machine load and `machines_changed` skip the others, whose workers stop); `ALL` / empty = every active machine.
- **Tables written:** `machine_variable_values` (one row per machine-variable: latest numeric value only, upserted every reading to detect changes; text variables not stored), `machine_telementry` (hypertable; only changed numeric values, `value`; the latest value is the newest row), `machine_live_status`, `machine_mode_log`, `machine_events`, plus `machines.mode_id`/`online_status`.
- **Change notifications:** `machines_changed` / `cnc_config_changed` NOTIFY make collectors follow edits without polling.
- **Builds and logs:** a non-386 build cannot load the DLL. Collection logs go to `logs/app.log` with the API's. A native crash inside the DLL ends the whole server.
**Current user**: GET `/me` (with `profile.theme` / `profile.language`); PUT `/me/theme {theme}` and PUT `/me/language {language}` save the logged-in user's own choice in `user_profiles` (any role; a missing profile is created). Languages: `dto.Languages` (en, fr, de, it — the client's `LANGUAGES`).
**Organization** (per user, `?user<Parent>Id=` on list): GET/POST/PUT/DELETE `/user-city|branch|location|shop|line|cell|supervisor`, GET/POST/DELETE `/user-country|state` — one set of files per table (`user_<level>_*.go`). Countries and states are assigned customer-wise by the superadmin only (POST/DELETE `/user-country|state` are `superAdminOnly`); everyone reads their company's. Names come from the client's country library (`country-region-data`, picked several at once and posted one row each); the server keeps no country list, and POST `/user-state` checks the country is the user's. The current user comes from the session (`middleware.Authenticate`); the organization belongs to the company: a customer works on its own rows and its members on the customer's rows (`targetUserID` in `handler/request.go`); only managers (superadmin, customer, admin) change them (`managersOnly`), other members read. The superadmin may act for another user with `?userId=<id>`; everyone else gets 403 for any other id. Lists (`listUserID`): the superadmin without `?userId=` gets every user's rows (each row has `userId`; the client shows a Customer column and picks the customer on the form). Every `/api` route except login/logout and `/health/machines/{id}/modes` needs a session; role-gated routes are wrapped with `superAdminOnly` / `managersOnly` in `internal/routes`. The role and company always come from the session (`middleware.Authenticate` reads them from the database), never from the request. Shared master data — branches, categories, products, brands, controllers, models — is read by everyone and written by the superadmin only.

## Configuration

**Environment Variables** (`.env` file):
```
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_DATABASE=7th_sense
DB_SSLMODE=disable
COOKIE_SECURE=false   # true when served over HTTPS
UPLOAD_DIR=uploads    # uploaded files (part drawings); keep outside the deploy dir
CNC_ENABLED=true      # the server connects the machines itself (GOARCH=386 build + Lib\); false = API only
CNC_SLOTS=300         # most machines the server holds (max 990)
CNC_LIB_DIR=Lib       # CNC_DLL_NAME=CncCommunication.dll
CNC_MACHINE_IDS=ALL     # ALL / empty = every active machine; 225,247 = connect (collect, log) only these
CNC_LOG_EVENTS=true       # machine events                       } internal/monitor/cnc_log.go,
CNC_LOG_MODES=true        # mode changes                         } all in logs/app.log;
CNC_LOG_MODE_CHECKS=true  # decided mode whenever NC inputs change } errors are always
CNC_LOG_VARIABLES=false   # changed variable values (many lines)  } written
CNC_LOG_OTHER=true        # start/stop, connections, reading rate, …
CNC_LOG_SNAPSHOT=false    # raw DLL status fields (mode, alarm, OpMsg, …) whenever one changes; analysis only
# CNC_STALE_SEC=30 CNC_QUIET_SEC=300 CNC_MAX_CONNECTING=20 CNC_STOP_TIMEOUT_SEC=15
```

Default values in `internal/config/config.go` if missing.

## Key Files Reference

- **Entry Point**: `cmd/server/main.go` - Wiring, server startup, graceful shutdown
- **Routes**: `internal/routes/routes.go` - Main router + middleware
- **Response Helpers**: `internal/response/responses.go` - Predefined responses
- **DB Connection**: `internal/db/postgres.go` - pgx pool setup
- **Migration Schema**: `migrations/schema.sql` - Complete database schema
- **Skills Document**: `7th_server_skills.md.md` - Detailed engineering principles

## Current Entities (Models)

- User, UserProfile, UserMachineAccess (`user_machine_access`; the older `user_machines` table is unused), Role
- Item
- Branch
- Brand, Category, Product
- Controller, MachineModel (models table)
- Address
- CompanyProfile
- Machine
- Clientele
- Shift, ShiftBreak (a break the shift takes), Break (`shift_id`), Holiday (holidays and week offs)
- Part, PartProcess, PartSubProgram (applicable machines in `part_process_machines`, no soft delete)
- LicenceKey
- DowntimeReason, IdleReason
- ProductionSchedule, ProductionScheduleItem
- CNC telemetry (`model/machine_telemetry.go`): MonitoredMachine, VariableConfig, TelemetryValue, ModeChange, MachineModePeriod (a `machine_mode_log` row with its mode name), MachineEvent, MachineLiveStatus

## Important Notes

1. **Soft Delete**: All queries must include `WHERE is_deleted = false`
2. **Unique Indexes**: Must include `WHERE is_deleted = false` for soft-delete compatibility
3. **Context**: Always pass `ctx` through all layers for cancellation support
4. **Logging**: Use structured logging in service layer, never log sensitive data
5. **Transactions**: Service layer only, keep short, rollback on error
6. **Idempotency**: PUT and DELETE must be idempotent
7. **Customer display names**: rows owned by a customer (machines, parts, clienteles, shifts, breaks, holidays, licence keys, production schedules) return `customerName` = the customer's first + last name (falls back to username; the client always shows this) and `companyName` = `company_profiles.company_name` (`''` when none, readability only). The company name comes from `company_profiles.user_id`, not `users.company_id` (that is only set on members).

## Testing

Integration tests in `tests/` directory using setup helpers:
- `setup_test.go` - Test database setup
- `api_helpers_test.go` - HTTP test helpers
- Entity-specific test files (e.g., `user_profile_test.go`)

## Debugging

VS Code configurations in `.vscode/launch.json`:
- **Debug Server** - Start under debugger (F5)
- **Debug Current Test** - Debug tests in open file
- **Attach to Running Server** - Attach to running process

Requires Delve debugger: `go install github.com/go-delve/delve/cmd/dlv@latest`

## Engineering Principles Summary

From `7th_server_skills.md.md`:
- Predictability over cleverness
- Consistency over individual preferences
- Simplicity over complexity
- Layered architecture is sacred
- Never trust input (defense in depth)
- Test behavior, not implementation
- Profile before optimizing
- Environment parity (dev/staging/prod)
- Documentation at implementation time

**Discipline in the small things creates freedom in the big things.**

## graphify

This project has a knowledge graph at graphify-out/ with god nodes, community structure, and cross-file relationships.

Rules:
- For codebase questions, first run `graphify query "<question>"` when graphify-out/graph.json exists. Use `graphify path "<A>" "<B>"` for relationships and `graphify explain "<concept>"` for focused concepts. These return a scoped subgraph, usually much smaller than GRAPH_REPORT.md or raw grep output.
- If graphify-out/wiki/index.md exists, use it for broad navigation instead of raw source browsing.
- Read graphify-out/GRAPH_REPORT.md only for broad architecture review or when query/path/explain do not surface enough context.
- After modifying code, run `graphify update .` to keep the graph current (AST-only, no API cost).
