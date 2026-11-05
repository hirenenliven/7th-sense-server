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

**Users**: GET/POST/PUT/DELETE `/users` — company-scoped (`service.UserService`). A customer is a company; its members (admin, supervisor, operator, qc, shop_head) have `users.company_id` = the customer. The superadmin manages every user (`?filterRole=customer`, `?companyId=` filter the list; `filterRole` is the listed users' role, never the caller's; `?membersOnly=true` lists only company members — no customers, no superadmin — and is what the client's Users list sends, while the Account Center Customer list sends `?filterRole=customer`) and names a member's `companyId`; a customer or its admins see and change only their company's members (`managersOnly` route guard, others get 403; other companies' users answer 404). PUT takes an optional `password` (8–72 characters, bcrypt-hashed): empty keeps the current password, so a manager resets a user's password from the edit form. Roles a user may give: `model.AssignableRoles` — superadmin: customer + member roles; customer: member roles; admin: member roles but admin; nobody gives superadmin. Body `machineIds`: machines (of the company) a supervisor/operator/qc/shop_head may access, stored in `user_machine_access` (migration 015; no soft delete — removing an assignment deletes the row); admins and customers see every company machine and store none. Responses include `role` (name) and `machineIds`.
**Roles**: GET `/roles` — the roles the current user may give (role dropdown of the user form); empty for users who manage no users.
**Items**: GET/POST/PUT/DELETE `/items`
**Branches**: GET/POST/PUT/DELETE `/branches`
**Brands**: GET/POST/PUT/DELETE `/brands` — `{name, code, description}` (name required, code optional; both unique). Every logged-in user reads; only the superadmin writes (`superAdminOnly`; client Account Center › Machine Catalog › Brands). Cannot be deleted while it has controllers or machines use it.
**Categories**: GET/POST/PUT/DELETE `/categories`
**Products**: GET/POST/PUT/DELETE `/products`
**Controllers**: GET/POST/PUT/DELETE `/controllers` — machine catalog Brand → Controller → Model. Body `{brandId, name, description}`: no code or variant, the variant is part of the name (e.g. "Sinumerik 828D", migration 017). Brand required, name unique per brand (case-insensitive), `?brandId=` on list. Cannot be deleted while it has models or machines use it; its brand cannot change while machines use it.
**Models**: GET/POST/PUT/DELETE `/models` (`models` table, `MachineModel` in Go — `model` is the package name) — `{controllerId, name}`, name unique per controller (case-insensitive), `?controllerId=` on list. Cannot be deleted, or moved to another controller, while machines use it. Brands, controllers and models: every logged-in user reads (machine form dropdowns); only the superadmin writes (`superAdminOnly` in `internal/routes`).
**Addresses**: GET/POST/PUT/DELETE `/addresses`
**Catalog**: GET/POST `/catalog`
**Auth**: POST `/auth/login`, POST `/auth/logout` (public), GET `/auth/me` — HttpOnly `ts_session` cookie, hashed sessions in `auth_tokens` (see README › Authentication). POST `/auth/login-as {customerId}` (`superAdminOnly`, ADR 0008): starts a customer session for a new browser tab (`auth_tokens.impersonated_by`, migration 031) and returns `{sessionKey, user}`; its HttpOnly cookie `ts_login_as` has Path `/api/as/<sessionKey>`. `middleware.RouteLoginAs` serves `/api/as/<key>/...` as `/api/...`, and `middleware.SessionToken` then reads only that cookie (no fallback to `ts_session`). Logout there ends only that session; password login there is 403.
**Machines**: GET/POST/PUT/DELETE `/machines` — each machine belongs to a customer (`user_id`) and has a brand, a controller of that brand and a model (`model_id`) of that controller; `product_id` is no longer used. The superadmin reads all (`?userId=` narrows to one customer); add / edit are for managers (`managersOnly`: superadmin, customer, admin — a customer or admin always on its own company, whatever `userId` says), delete is superadmin-only; a customer and its admins read every machine of the company, other members only the machines assigned to them (`user_machine_access`, see `MachineService.machineScope`); anything else answers 404. The password is write-only (never returned; blank on PUT keeps it). Form-only writes: runtime columns (online status, mode, …) are not changed through the API. Accuracy / health options are fixed lists in the client (`constants/machines.js`).
**Parts**: GET/POST/PUT/DELETE `/parts` (migration 016: `parts`, `part_processes`, `part_process_machines`, `part_sub_programs`) — a part belongs to a customer (`user_id`) and is saved with its processes in one transaction (`service.PartService`). Body `{userId, clienteleId, partNo, partName, partCount, processes: [{id?, process, programName, machineIds, drawingNo, drawingRevisionNo, drawingFile, drawingName, accuracyId, setupTime, cycleTime, loadingTime, hasSubProgram, subPrograms: [{name, cycleTime}]}]}`, times in seconds. On PUT a process with `id` is updated, one without is added, saved processes left out are soft-deleted; machines and sub programs of a process are replaced. Production times are computed, never sent: process = setup + cycle + loading + sub program cycle times, part = sum of processes. Applicable machines must be the customer's; Part No. + name unique per customer (case-insensitive). `clienteleId` is optional and must be one of the customer's clienteles. Everyone reads their company's parts (superadmin all, `?userId=` narrows); managers (`managersOnly`) write their own company's. Drawings: POST `/parts/drawings` (multipart `file`, PDF ≤ 10 MB, managers) stores `UPLOAD_DIR/part-drawings/<32 hex>.pdf` and returns `{file, name}` for a process's `drawingFile`/`drawingName`; GET `/parts/{id}/processes/{processId}/drawing` streams it to anyone who may read the part.
**Clienteles**: GET/POST/PUT/DELETE `/clienteles` (migration 017) — the clients a customer makes parts for (`service.ClienteleService`). Body `{userId, clientCode, name, email, contactNumber, address, countryName, stateName, city}` (+ `status` on PUT, blank keeps): Client ID (`clientCode`) unique per customer (case-insensitive), email optional, `contactNumber` as `+<code><number>`, country and state as names ("India", "Gujarat"; migration 018). Same access as parts: everyone reads their company's (superadmin all, `?userId=` narrows), managers write their own company's. A clientele used by parts cannot be deleted or moved to another customer. Parts and clienteles share `companyOwner` / `ownCompanyCustomer` (`service/company_scope.go`).
**Shifts**: GET/POST/PUT/DELETE `/shifts` (migrations 019–021) — a customer's shifts (`service.ShiftService`). Body `{userId, name, shiftType: day|night, startTime, endTime, workingDays: [1..7], breakIds}` (+ `status` on PUT, blank keeps), times `"HH:MM"` 24 h; an end not after the start runs past midnight. `breakIds` are the customer's breaks the shift takes: saving sets their `breaks.shift_id` and releases (NULL) the shift's breaks left out; deleting the shift releases all of them. A break taken by another shift is refused; the breaks must lie inside the shift and not overlap (`checkBreaksFit`, minutes from the shift start). Responses add `breaks` (ordered from the shift start), `durationMinutes`, `breakMinutes` (computed, not stored). Name unique per customer (case-insensitive).
**Breaks**: GET/POST/PUT/DELETE `/breaks` (migrations 020–021, `service.BreakService`) — a customer's breaks `{userId, shiftId?, name, startTime, endTime}` (+ `status`), name unique per customer. `shiftId` (optional, one of the customer's shifts) is the shift that takes it; responses add `shiftName`. A break on a shift must fit it beside the shift's other breaks. Deleting a break takes it off its shift.
**Holidays**: GET/POST/PUT/DELETE `/holidays` (migration 020, `service.HolidayService`) — a customer's holidays `{holidayType: "holiday", name, startDate, endDate}` and week offs `{holidayType: "week_off", weekDay: 1..7, name?}` (+ `status`); the fields the kind does not use are cleared. A weekday is a customer's week off once; the same holiday name and start date once. Shifts, breaks and holidays have the same access as parts and clienteles (`companyOwner`).
**Downtime / idle reasons**: GET/POST/PUT/DELETE `/downtime-reasons`, `/idle-reasons` (migration 027: `downtime_reasons`, `idle_reasons`; replace the old `mst_downtime_reason` / `mst_idle_reason`; one set of files per table, `service.DowntimeReasonService` / `IdleReasonService`) — a customer's reasons `{userId, reasonCode, reason}` (+ `status` on PUT, blank keeps). `reasonCode` is optional (≤ 10, NULL when blank) and unique per customer when given; `reason` (≤ 50) unique per customer (both case-insensitive). Same access as clienteles (`companyOwner`).
**Licence keys**: GET/POST/PUT/DELETE `/licence-keys` (migration 024, `service.LicenceKeyService`) — keys the superadmin generates for a customer, format `K7M2Q-XR9TB-4HWZP-N8CDF` (4 groups of 5 from `dto.LicenceKeyAlphabet`, no I/O/0/1, `crypto/rand`). POST `/licence-keys/generate {count: 1..100}` returns `{licenceKeys}` without saving (the "Generate Key" button); POST `/licence-keys {userId, licenceKeys: [{licenceKey, duration}]}` saves them in one statement, all or none (a key already saved refuses the batch). `duration` is a subscription plan code (`dto.LicenceKeyDurations`, migration 025: M1…M12, Y1…Y3, LIFETIME, TRIAL = `LicenceKeyTrialDays` 30 days, TRIAL_1Y); the server stores `expires_on` (last valid day, counted from the day the key is saved; NULL = lifetime) and responses add `expiresOn`, `isExpired` (computed). PUT `{userId, duration, status}` changes the customer / plan / status (the end date is recounted from `created_at`; the key never changes). Generate, save, edit and delete are `superAdminOnly`; managers read their company's (`companyOwner`, `?userId=` narrows for the superadmin). Key unique among keys not deleted.
**Production schedules**: GET/POST/PUT/DELETE `/production-schedules` (migration 026: `production_schedules`, `production_schedule_items`; replaces the old `production_schedule` / `schedule_items`, `service.ProductionScheduleService`) — a customer's machine planned for one day (`scheduleType: single`, `endDate` = `startDate`) or a range (`multiple`), saved with its jobs in one transaction. Body `{userId, machineId, scheduleType, startDate, endDate, items: [{id?, scheduleDate, jobId, supervisorId, operatorId, shiftId, partId, processId, quantity, startTime, endTime}]}` (+ `status` on PUT), dates `YYYY-MM-DD`, times `"HH:MM"`; single day sets every `scheduleDate` to `startDate`. Items are kept / added / soft-deleted like part processes; `totalMinutes` is computed. Rules: machine, shift, part are the customer's; supervisor / operator are company members with that role; the process is the part's and lists the machine among its applicable machines; each job lies inside its shift (`scheduleDate` = the day the shift starts, so 02:00 in a 22:00–06:00 shift is next morning); jobs do not overlap; a machine has one schedule per day. Same access as parts. Responses add `customerName`, `machineName`, `itemCount`, `totalQuantity`; items add `supervisorName`, `operatorName`, `shiftName`, `partName`, `processName`.
**Health (mode history)**: GET `/health/machines/{id}/modes` (`service.HealthService`) — the machine's `machine_mode_log` periods of the last `dto.MachineModeHistoryDays` (10) days, oldest first: every period still running at or after the window start (one that began earlier is returned whole). Response `{machineId, machineName, currentModeId, onlineStatus, from, to, periods: [{id, modeId, modeName, programName, startedAt, endedAt, durationSec}], totals: [{modeId, modeName, periodCount, durationSec}]}`; the open period has `endedAt` null and `durationSec` up to now. **Public, no session** (registered on the outer router, `registerHealthRoutes`): any machine not deleted, others 404. GET `/health/machines` (public too, `HealthService.GetMachineModes`) — every machine not deleted with its current mode: `{at, machines: [{machineId, machineName, customerName, companyName, modeId, modeName, onlineStatus, programName, modeSince, modeDurationSec}]}` from the open `machine_mode_log` period (`machines.mode_id` when none; `modeId` null = never collected), ordered by name. Chart page: GET `/mode-chart/` (`?machineId=` selects one) serves `web/mode-chart/index.html` + `mode-chart.js`, embedded with `go:embed` (`web.ModeChart`) — a dashboard on the left (a count tile per mode that filters, search, machines grouped by current mode, refreshed every 30 s) and, for the selected machine, a timeline per mode with zoom level buttons (`ZOOM_LEVELS` in `mode-chart.js`: 30s, 1m, 2m, 5m, 10m, 30m, 1h, 3h, 8h, 24h, 3d, 10d; a level ends where the current view ends, the matching level is highlighted, 30 s is the closest zoom; a selected machine opens on the latest 5 minutes, `DEFAULT_ZOOM_MS` = the 5m level) (ECharts 5.5.1 from jsDelivr) and the list of periods; mode colours are in `mode-chart.js`. The timeline has a **Cycles** row (each cycle of the last 24 h coloured by status: completed / unfinished / running) and below it a Cycles section in plain words: summary cards, "where the cycle time went" (split by mode), a cycle time bar chart per cycle (stacked by mode, planned cycle time as a dashed line, opens on the newest `CYCLE_BARS_SHOWN` 40) and the cycle list (newest `CYCLE_TABLE_LIMIT` 200; end reasons in words, `CYCLE_END_REASONS`); a click on a cycle zooms the timeline to it. GET `/health/machines/{id}/cycles` (public too, `MachineCycleService.GetMachineCycleHistory`) — the machine's cycles that started in the last `dto.MachineCycleHistoryHours` (24) hours, same shape as GET `/machine-cycles`. Not the `GET /health` liveness check.
**Machine cycles**: GET `/machine-cycles?machineId=&from=YYYY-MM-DD&to=YYYY-MM-DD`, GET `/machine-cycles/{id}` (`service.MachineCycleService`, read only — the collector writes them) — the cycles of a machine the user may read (`MachineService`, others 404), started from the from day to the to day (to defaults to today, from to to; at most `dto.MachineCycleMaxDays` 31). Response `{machineId, machineName, from, to, cycles: [MachineCycle with machineName, jobId, partName, processName, shiftName, durationSec], summary: {cycleCount, completedCount, unfinishedCount, runningCount, partsMade, productionSec, holdSec, downtimeSec, offlineSec, averageCycleSec, averageProductionSec}}` (averages over completed cycles: whole length / In Production time).
**Dashboards (page-wise charts)**: `service.DashboardService`, `handler.DashboardHandler`, `routes/dashboard_routes.go` — GET `/dashboard/home?date&userId` (the client Home page; `dto.HomeDashboardRequest`: date not after today, `userId` = customer for the superadmin, others 403 for another company): the **active** machines of `MachineService.GetAll` together — `fleet` now (`total`, `online`, `offline`, `awaitingData` = no mode, `inactive` counted only, `modes: [{modeId, count}]` from the open `machine_mode_log` period, `machines.mode_id` without one: `FindOpenModePeriods`), and for the day up to now (whole day, no shifts, the `machine-summaries` rules): `oee` (fleet `oeeSum`, null before any time), `cycles {cycleCount, completedCount, unfinishedCount, runningCount, partsMade}`, `plan {quantity, partsMade}` (null without jobs: `ProductionScheduleRepository.SumMachineQuantities`, active schedules; only machines with a job count), `time {spanSec, groups: [{group, durationSec}]}` (mode groups, `modeGroup`), `hourly` (24 rows `{hour, partsMade, previousPartsMade}`, completed cycles by end hour, the day and the day before), `machines` (`dto.HomeMachine`: mode now + `modeSince`, OEE parts, cycles, parts, `plannedQuantity`), `alarms` (top `dto.HomeTopAlarms` 10 `alarm_on` codes: `{code, message (latest), count, machineCount, lastAt}`, `SumAlarms`). The /productivity and /machine-health charts: read only, every logged-in user, for machines they may read (`MachineService.GetByID`, others 404); calendar days in the server's time zone; query `machineId` (or `machineIds=1,2`), `date` (one day) or `from`/`to` (≤ 31 days, to defaults to today), `shiftIds=1,2` (the machine customer's shifts; a shift past midnight belongs to the day it starts), `program`, `schedule=with`. GET `/productivity/mode-timeline?machineId&date` or `?machineId&minutes=30` (last 1 … `dto.ModeTimelineMaxMinutes` 1440 minutes up to now) — the `machine_mode_log` periods clipped to that window (`{programName (latest period's), opMode (machine_variable_values of cnc_variables `opMode`, named by `model.OpModeNamesByEngine`), from, to, serverTime, periods: [dto.TimelinePeriod {id, modeId, modeName, startedAt, endedAt}]}`; the open period has `endedAt` null; the client's Gantt grows it live and reloads every 30 s). GET `/productivity/mode-totals?machineId` + `date`, `from`/`to` (≤ 31 days) or `all=true` (every recorded period, dates ignored; `DashboardRequest.AllTime` skips the date checks) — seconds per mode summed in SQL (`SumModeDurations`, every period clipped to the range, the open one up to now): `dto.ModeTotalsResponse {from (all time: start of the first period, null when none), to (never after now), serverTime, totals: [model.ModeDuration {modeId, modeName, durationSec}]}` (the Machine Modes card). GET `/productivity/oee` — per day, over the machines together (`oeeSum` adds their planned / run / standard / production seconds and completed / closed cycles, then divides, so a machine weighs by its time and cycles; for one machine the same as before): availability = run time (modes 3, 12) / planned time (the day or the shifts picked; `with` = the customer's active shifts on their working days; planned stops 7, 13 left out), performance = standard cycle time × parts (`parts_made`, from ProgramPartCount only; a cycle never counts as a part) / production seconds of completed cycles with a standard time, 0-part ones included (100 without), quality = completed / closed cycles; performance / quality / oee null without closed cycles; plus `programs`. GET `/productivity/cycle-counts` — per day `{produced, finished, unfinished}` (cycles started in the day / shifts). GET `/productivity/finish-cycles?machineId&date` — `{cycles: [dto.FinishedCycle {endedAt, programName, durationSec, productionSec, holdSec, downtimeSec, offlineSec}] (completed), summary: dto.FinishCyclesSummary {cycleCount, completedCount, unfinishedCount, averageCycleSec, averageProductionSec}, lastFinishedAt, programs}`. GET `/productivity/machine-summaries?date` (default today) — every machine the user may read (`MachineService.GetAll`, no `machineIds`), each `dto.MachineSummary {machineId, availability, performance, quality, oee, cycleCount, completedCount, unfinishedCount}` of the whole day (OEE rules as above for one machine; availability null before any planned time; the Monitoring machine list). GET `/productivity/virtual-screen?machineId` (`virtual_screen_handler.go` / `_service.go` / `_repository.go`, `dto.VirtualScreenResponse`; the client Virtual Screen tab polls it every 5 s) — what the machine does now: mode / program from the open `machine_mode_log` period (`machines.mode_id` without one), `opMode`, `live` (latest `machine_variable_values` by `cnc_variables` code, `dto.VirtualScreenVariableCodes`: spindle load / rpm / override, feed rate / override + unit from `FeedRev`, axis load X/Y/Z, `toolNumber`, `isEmergency`; null = never read; `updatedAt` = newest), `alarm` (last `alarm_on`, `active` while no `alarm_off` after it), `cycle` (the running `machine_cycles` row: elapsed, standard seconds) and `job` (`model.JobProgress`: the running cycle's schedule item, else `FindScheduledJob` now; quantity, times, part / process / clientele, standard cycle, completed cycles). Test: `tests/virtual_screen_test.go`. GET `/machine-health/machine-time-by-day?machineId&from&to` — hours per mode group per day (`dto.ModeGroup*`: production 3/12, hold 4, idle 1/2/5, downtime 6/7/13, offline 9 + unrecorded time). Responses carry only what the charts draw: each row type in `dto/dashboard_dto.go` is filled by one `New…` constructor, so adding / dropping a field is that struct + its constructor (+ the client chart). Test: `tests/dashboard_test.go`.
**CNC data collection** (`docs/cnc-telemetry-collector-plan.md`, migration `028_cnc_machine_telementry.sql`; read API: Health above):
- **What starts it:** the API server itself, in-process (`internal/collector.Supervisor` runs `monitor.Collector` in a goroutine; no separate process or service). `CNC_ENABLED` defaults to `true`; `false` = API only. The DLL is 32-bit, so the server must be built `GOARCH=386`: `.\run.ps1`, `.\build.ps1` (→ `bin\server.exe`) and `go tool air` do. A 64-bit build refuses to start while `CNC_ENABLED=true` (`collector.ErrWrongBuild`); with `CNC_ENABLED=false` it serves the API only.
- **Machines:** the server connects every machine (shard 0 of 1, at most `CNC_SLOTS` ≤ 990) with `status = 'active'`, not deleted, through `CncCommunication.dll` in `Lib\` (`CNC_DLL_NAME`; the older `FanucCncCommunication.dll` reads no Siemens alarm node, so `AlarmNo` stays 0 there) (`internal/cnc`: `CNCConnect`, one callback slot per machine, Siemens login via `SIEMENS_USER`/`SIEMENS_PASS` under a lock). The engine comes from `controllers.cnc_engine` (NULL = `config_error` event, not connected).
- **Variables:** `cnc_variables` (old `mm_variables` ids), `cnc_variable_addresses` (source `snapshot` = a DLL snapshot field, the only one readable today; `plc`/`nc` wait for a DLL read function), `machine_variables` (per machine, `id` = variable id, key `(machine_id, id)`, hard delete). `TimerMinute` / `TimerMsec` (the snapshot's NC timer; which timer the DLL fills is undocumented) are snapshot fields a variable can read. Axis positions of X, Y, Z (MachPos / AbsPos / RelPos / DistToGo, ids 208–219, addresses `MachPos[0]` = X … `[2]` = Z) come from migration 029; the code reads all 5 snapshot axes, a variable decides which are stored.
- **Flow:** `internal/monitor`: one worker per machine loads live values from `machine_variable_values` (numeric only), upserts every reading to detect changes, and inserts **only changed numeric values** into time-series (`MachineTelemetryService.UpsertLiveValues` → `InsertTelemetry`; no batching, but change detection is used: `store_mode` / deadband / `min_interval_ms` are not used). Text variables are not stored at all. Decides the `machine_modes` id on **every connected reading** from the reading's NC state, not from the stored variables (`internal/monitor/machine_mode.go`: `machineStateFromReading` → `decideMachineMode`, first match wins; the same rules as the vendor sample `G0_V0.0.2_DLL/main.go`, tested live on Fanuc and Siemens): no connection → Offline 9; alarm (`AlarmNo` ≠ 0 or `IsEmergency`) and program not running → Downtime 6; AUTO + program running + feed override (`FeedRatePercetnage`) > 0 + no alarm → In Production 3; program held, or running with feed override ≤ 0, an alarm or not in AUTO → Hold 4; program reset and no alarm, in any operation mode → Idle 5 (reason `program reset, no alarm (any mode)`; every state matches a rule, the `no rule matched` fallback is never reached). `MachineMode` / `ProcStatus` are read per engine (`autoOpModes`, `programStatesByEngine`): AUTO = Fanuc 1 (MEM), Siemens 2; Fanuc ProcStatus 3 START → running, 1 Stop / 2 HOLD / 4 MSTR → held, 0 Reset; Siemens 2 Active → running (as Fanuc START), 1 Stop / 3 Waiting (as Fanuc HOLD) / 4 Interrupted → held, 0 Reset (so 2 is Hold on Fanuc but running on Siemens); a value not in the table counts as reset; engines without a table (simulator, Syntec, Mitsubishi, none) are read with the Siemens numbering. The logs name both per engine (`internal/monitor/nc_status.go`: `opModeName` / `procStatusName`; mode logs show `opMode`, `procStatus` by name and `programState` = running / held / reset). chanStatus is not in the snapshot and is not used. Setup / Rapid Mode are no longer set. Lost connection (`worker.checkConnectionLost`): the DLL stops calling back while nothing changes (a held / idle machine is silent), so silence is not Offline — Offline comes from an `IsConnected=0` reading, a failed connect, or no connected reading for `CNC_STALE_SEC` (30) while the DLL's last word was "not connected"; after a connected reading, only no callback for `CNC_QUIET_SEC` (300) reopens the connection, keeping the mode (`stale` event "mode kept"). At start every machine is set Offline (`MarkShardOffline`) until its first connected reading decides its mode. Mode changes, events and statuses are written directly too. Logs: `cnc collector: worker started`, `cnc: variables loaded` (variable → field map) and `cnc: active variables not readable` (FindVariableConfig returns those with an empty address), `cnc: machine subscribed` on connect, `cnc: controller not sending data, no values read` once per connection while the DLL reports `IsConnected=0` (no values are read or logged then), `cnc: controller data started` / `stopped`, `cnc: reading rate` every 10 s per machine (callbacks / processed / skipped / stored). Every collector log line goes through one function, `writeLog(settings, logType, level, machine, message, attributes...)` in `internal/monitor/cnc_log.go` (called as `collector.log` / `collector.logMachine` / `machineWorker.log`; no direct `slog` calls in `internal/monitor`): log types `logTypeEvent` / `logTypeMode` / `logTypeModeCheck` / `logTypeVariable` / `logTypeSnapshot` / `logTypeDowntime` / `logTypeCycle` (`CNC_LOG_CYCLES`) / `logTypeOther`, levels `slog.LevelInfo` / `LevelWarn` / `LevelError`; a line about a machine starts with `[machine <id> <name>]` and has `machineId` + `machine`, every line has `logType`; errors are always written, each type by its switch (`CNC_LOG_OTHER` for `other`; `cnc collector: settings` is always written): `CNC_LOG_EVENTS` → `[machine …] cnc machine event <event>` (failures as Warn); `CNC_LOG_MODES` → `[machine …] cnc: mode <final mode> (from <previous>)` (`mode`, `from`, `reason` = the rule that matched, then only the inputs that rule checks: `modeDecision.conditions` from `decideMachineMode` — e.g. Offline `connected`; Downtime `alarmNo`, `emergency`, `procStatus`, `programState`; In Production `opMode`, `procStatus`, `programState`, `feedOverride`, `alarmNo`, `emergency`; Hold per reason; Idle `procStatus`, `programState`, `alarmNo`, `emergency`); `CNC_LOG_MODE_CHECKS` → `[machine …] cnc: mode check <mode>` whenever the decided rule or its inputs change (`mode`, `current`, `modeChanges`, `reason` + the same inputs); `CNC_LOG_SNAPSHOT` (default false) → `[machine …] cnc: snapshot status` whenever a raw DLL status field changes (DLL names: `IsConnected`, `ResultStatus`, `MachineMode` + `machineModeName`, `ProcStatus` + `procStatusName`, `IsEmergency`, `TimerMinute` / `TimerMsec` (logged, not compared), `FeedRatePercetnage`, `SpindeOverridePercentage`, `CmdFeedOvr`, `CmdSpOvr`, `AlarmNo`, `AlarmMsg`, `OpMsgNo`, `OpMsg`, `ActiveGCode`, `ProgName`, `missingNodes`; `cnc.Reading.Status`); `CNC_LOG_VARIABLES` → `[machine …] cnc: variable values changed` once per reading that changed values (`count` + `changes: ["208:X_MachPos 12.5 -> 12.75", …]`, previous `none` for a first value); `CNC_LOG_DOWNTIME` (default true) → `[machine …] cnc: downtime check downtime=<true|false>` whenever a Downtime input changes (`machineType` Fanuc/Siemens/…, `downtime`, `alarmActive`, `alarmNo`, `alarmMsg`, `emergency`, `programRunning`, `procStatus`, `mode`, `reason`; `logDowntimeCheck`). `cnc collector: settings` logs the machines and switches at start. `CNC_MACHINE_IDS=225,247` makes the collector connect only those machines (`collector.go`: the machine load and `machines_changed` skip the others, whose workers stop); `ALL` / empty = every active machine.
- **Tables written:** `machine_variable_values` (one row per machine-variable: latest numeric value only, upserted every reading to detect changes; text variables not stored), `machine_telementry` (hypertable; only changed numeric values, `value`; the latest value is the newest row), `machine_live_status`, `machine_mode_log`, `machine_events`, plus `machines.mode_id`/`online_status`.
- **Cycles** (migration `030_add_machine_cycles.sql`, replaces the old `mm_cycletime` / `mm_cycletime_history`; written by `internal/monitor/machine_cycle.go`: `decideCycleStep` decides, `cycleTracker` counts, the worker calls `followCycle` on every connected reading before the mode, `checkCycle` on its tick; the job / shift / loading time are worked out by `MachineCycleService.Create`): `machine_cycles`, one row per run of the part program — opens on program state reset → running, held time stays in the cycle; completed on a ProgramPartCount step (`part_count`), CurProgLine jumping back to the program start while running with the cycle's own `ProgName` in both readings (`program_loop`, M99 programs whose counter never moves; a sub program call M98 shows the sub program's name and is not a loop) or running → reset (`program_end`); the cycle after a part count / loop waits `cycleStartGrace` (5 s) and is not written if the program resets meanwhile (the counter counts a reading before M30); unfinished on `alarm_reset` (an `IsEmergency` without `AlarmNo` only once on for `emergencyResetGrace` 2 s: the Fanuc RESET key shows it for < 1 s, so such a reset first ends as `program_end` / `operator_reset` and `confirmEmergencyStop` changes it if it stays on), `operator_reset` (held → reset), `program_changed`, `offline_reset`, `collector_stop`, `timeout` (CHECK ties each reason to its status). The all-zero first snapshot after a connect is skipped; a counter drop or a step > 1 (`maxPartCountStep`: one part per M30) only moves the baseline. After a long Offline `ended_at` = when it went offline. The program name is compared only after a gap (a sub program call may show its own name); at most one part per cycle and one part may take several cycles (0-part cycles stay `completed`): a counter step while no cycle is open (the DLL can send it only with the next start) or within 5 s of the next cycle's start is the part of the cycle closed last when it has none (`lateParts`; the open cycle's `part_count_start` moves); open longer than `CNC_CYCLE_MAX_HOURS` (24) → `timeout`. Inserted at open, seconds kept in the worker and written at each mode change (`recordModeChange` → `saveCycleProgress`), every 5 min while a mode lasts, when a waiting cycle opens, and at close (a cycle left open is rebuilt from `machine_mode_log` when the worker starts); the job is the schedule item running at the start (the one with the cycle's program first), else the only part process of the machine with that program; the shift date falls back to the calendar day; `loading_sec` is NULL across a shift / day change or an Offline gap; length = `ended_at - started_at` (not stored). Seconds per mode (`production_sec`, `hold_sec`, `downtime_sec`, `offline_sec`), `loading_sec` (gap since the previous cycle), `part_count_start` / `_end` / `parts_made` (ProgramPartCount), and the job matched at start (`schedule_item_id`, `part_id`, `process_id`, `shift_id`, `shift_date`, `standard_cycle_sec`). One running cycle per machine (`machine_cycles_open_key`; a refused open — `response.MachineCycleRunning` — makes the worker reload and close the stale running row); every end reason and how it is identified: `docs/machine-cycles.md` (its *Fixed values and why they are needed* section explains `cycleCheckpointInterval`, `loopMinimumDropFraction`, `maxPartCountStep`, `cycleCounterGrace`, `cycleStartGrace`, `CNC_CYCLE_MAX_HOURS`, `CNC_STALE_SEC`, `CNC_QUIET_SEC`); no soft delete (collector data, like `machine_mode_log`).
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
API_DOCS=false        # true serves Swagger UI + OpenAPI at /api-docs/ (dev only)
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
CNC_LOG_DOWNTIME=true     # Downtime inputs (machine type, alarm, E-stop, program status) + downtime yes/no, whenever one changes
CNC_LOG_CYCLES=true       # each cycle opened / closed / restored (machine_cycle.go)
CNC_CYCLE_MAX_HOURS=24    # a cycle open this long is closed as unfinished (timeout)
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
- MachineCycle (`machine_cycles`, collector data), CycleJob (what a cycle is matched to)
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

## Graphify (MANDATORY FIRST STEP FOR EVERY TASK)

This project has a knowledge graph in `graphify-out/` containing architecture, feature communities, god nodes, and cross-file relationships.

Rules:
- Before investigating or changing code for **any task**, read this `CLAUDE.md` and consult the latest Graphify knowledge graph.
- Start with `graphify query "<question>"` to locate relevant server modules and dependencies. Use `graphify path "<A>" "<B>"` to trace relationships and `graphify explain "<concept>"` for focused context.
- If `graphify-out/wiki/index.md` exists, use it for broad navigation. Review `graphify-out/GRAPH_REPORT.md` when broader architecture context is needed or targeted queries are insufficient.
- Identify all potentially affected files and dependent flows before editing. For server changes, trace the relevant route, middleware, handler, DTO/validation, service, repository, model, database/migration, background processing, and shared utility dependencies where applicable.
- Follow the strict layered architecture documented above. Inspect the actual source code of every relevant file before editing; Graphify is a navigation aid, not the source of truth.
- If the graph is incomplete, stale, or does not reveal the full flow, use code search and inspect additional files. Do not assume the graph lists every required change.
- Do not begin implementation until the existing flow and the complete set of required changes are understood. Make the smallest complete change set and avoid unrelated modifications.
- Preserve existing APIs, database schema, business rules, and architecture unless the request explicitly requires a change. Do not bypass or skip application layers.
- After adding, renaming, or removing code files, run `graphify update .` to keep the graph current. Never hand-edit files inside `graphify-out/`.


## 📚 Documentation Rules (MANDATORY)

API docs: `docs/api/openapi.yaml` (OpenAPI 3.0.3, generated — never edit by hand), Swagger UI at `/api-docs/` when `API_DOCS=true`, Postman `docs/postman/7th-sense.postman_collection.json` + `7th-sense.postman_environment.json` (generated). Server ADRs: `docs/adr/`. Page flows that show server steps live in `../7th-sense-client/docs/pages/`.

- **Meaningful change (document it):** a new page or a change to a page's flow or actions; a new or changed endpoint or API call (route, payload, response, status, auth); business logic changes; schema / model changes; auth, permissions, middleware, guards; state management; a new library, integration or pattern.
- **Skip:** styling, copy, renames, formatting, logs, refactors that change no behaviour, small fixes that don't change the flow or API contract.
- **On a meaningful change:** update the swag annotations on the handler (`// @Summary`, `@Description` ending with a line `Used by pages: <Page titles from ../7th-sense-client/docs/pages/README.md>`, `@Tags` = resource, `@ID` = the client service function name, `@Security cookieAuth`, `@Param` for path / query / body DTO, `@Success` with `response.Envelope{data=<type>}`, `@Failure` 400/401/403/404/500 where relevant, `@Router /v1/...`). Then run `go generate ./docs/api` — it rebuilds `openapi.yaml` (swag → `tools/apidocs`: OpenAPI 3, cookie auth, examples) and the Postman collection (requests keep their names; never duplicate). Validate: `go test ./tests/ -run TestAPIDocs`.
- **Also update every page in `../7th-sense-client/docs/pages/`** that uses the changed endpoint (find them through "Used by pages" in the annotation), because the page flows show the server steps too.
- **ADRs** (Michael Nygard template: Title, Date, Status, Context, Decision, Alternatives considered, Consequences): add `docs/adr/NNNN-<slug>.md` + a row in `docs/adr/README.md` when a real choice was made (library, pattern, architecture). Never edit an accepted ADR's content: write a new one and mark the old one "Superseded by NNNN".
- Update docs in the SAME task, before saying it's done. Edit existing sections; don't duplicate them. Make sure every relative file link still points to a real file.
- **End-of-task line in every reply:** `Docs: pages [..] | openapi [..] | postman [..] | adr [..]` — or `Docs: none: <reason>`.
