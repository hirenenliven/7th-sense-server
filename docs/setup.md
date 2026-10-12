# Setup on a fresh system

Steps to get `7th-sense-server` running with hot reload on a new machine.
Commands are for Windows (PowerShell); macOS/Linux differences are noted where they matter.

## 1. Install Git

Download from https://git-scm.com/downloads and install with default options.

```powershell
git --version
```

## 2. Install Go (1.27.1 or newer)

Download the installer for your OS from https://go.dev/dl/ (Windows: the `windows-amd64.msi`) and run it.
The installer adds Go to PATH. Open a **new** terminal, then check:

```powershell
go version          # go version go1.27.1 windows/amd64
go env GOARCH       # must print amd64 (or arm64 on Apple Silicon)
go env GOPATH       # usually C:\Users\<you>\go
```

If `GOARCH` prints `386`, remove the `GOARCH` environment variable (System Properties →
Environment Variables, or your terminal / VS Code settings). Otherwise Go builds 32-bit binaries
and `go install` puts tools in `go\bin\windows_386\`, where the terminal cannot find them.

## 3. Add Go's bin folder to PATH

Tools installed with `go install` (like the Delve debugger) go into `%GOPATH%\bin`.

- **Windows:** System Properties → Environment Variables → User variables → `Path` → New →
  `C:\Users\<you>\go\bin`. Open a new terminal afterwards.
- **macOS/Linux:** add to `~/.zshrc` or `~/.bashrc`:
  ```sh
  export PATH="$PATH:$(go env GOPATH)/bin"
  ```

## 4. Air (hot reload): nothing to install

Air is declared as a tool in `go.mod` (`tool github.com/air-verse/air`), so its version is pinned
for everyone and it is downloaded automatically the first time you run it:

```powershell
go tool air -v
```

To upgrade it for the whole team: `go get -tool github.com/air-verse/air@latest`, then commit
`go.mod` and `go.sum`.

## 5. Install PostgreSQL and create the database

Download from https://www.postgresql.org/download/ and install (remember the password you set
for the `postgres` user). Make sure the `bin` folder (e.g. `C:\Program Files\PostgreSQL\18\bin`)
is on PATH, then:

```powershell
pg_isready                                   # localhost:5432 - accepting connections
psql -U postgres -c "CREATE DATABASE \"7th_sense\";"
```

Run the database migration to create all tables:

```powershell
psql -U postgres -d 7th_sense -f migrations/001_create_tables.sql
```

The server connects on startup and exits with `database connection failed: ...` if it cannot.

## 6. Clone the project

```powershell
git clone https://github.com/hirenenliven/7th-sense-server.git
cd 7th-sense-server
```

## 7. Create the `.env` file

`.env` is not committed to git, but `.env.example` is provided as a template:

```powershell
cp .env.example .env
```

Then edit `.env` if needed. The default values are:

```env
# Server Configuration
PORT=8080

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_DATABASE=7th_sense
DB_SSLMODE=disable
```

Any missing value falls back to the defaults in `internal/config/config.go`.

## 8. Download dependencies

```powershell
go mod download
```

## 9. Run

With hot reload (rebuilds and restarts on every `.go` file save, using `.air.toml`):

```powershell
go tool air
```

Without hot reload:

```powershell
go run ./cmd/server
```

Check it is up:

```powershell
curl http://localhost:8080/health
# {"statusCode":200,"success":true,"message":"success"}
```

Stop with `Ctrl+C`.

## Debugging (breakpoints) in VS Code

One-time setup: install the Go extension. The first time you press `F5` without Delve (`dlv`),
it offers to install it; click **Install**. Or install it yourself:

```powershell
code --install-extension golang.go
cd ~
go install github.com/go-delve/delve/cmd/dlv@latest
```

Unlike Air, Delve is not a `go.mod` tool: VS Code needs `dlv.exe` on PATH (step 3).

`.vscode/launch.json` has three configurations (Run and Debug panel, `Ctrl+Shift+D`):

| Configuration | Use it to |
|---------------|-----------|
| **Debug Server** | Start the server under the debugger (`F5`) |
| **Debug Current Test** | Debug the tests in the package of the open `_test.go` file |
| **Attach to Running Server** | Attach to a server that is already running |

1. Click left of a line number to set a breakpoint (red dot), e.g. in `internal/handler/user_handler.go`.
2. Pick **Debug Server** and press `F5`.
3. Send a request (Postman / curl). Execution stops at the breakpoint; inspect variables in the
   left panel or hover over them.
4. `F10` step over, `F11` step into, `Shift+F11` step out, `F5` continue, `Shift+F5` stop.

Stop `go tool air` first: both use the same port. The configurations set `cwd` to the project root, so
`.env` is found and logs go to `logs/`.

## Troubleshooting

| Problem | Fix |
|---------|-----|
| `dlv.exe` is in `go\bin\windows_386\` | `GOARCH=386` is set, see step 2. Delete that folder and reinstall Delve. |
| `database connection failed: ... password authentication failed` | Wrong `DB_USER` / `DB_PASSWORD` in `.env`. |
| `database connection failed: ... database "7th_sense" does not exist` | Create it (step 5) or fix `DB_DATABASE`. |
| `database connection failed: ... connect: connection refused` | PostgreSQL is not running: `pg_isready`, or start the service in `services.msc`. |
| `bind: Only one usage of each socket address` | Port 8080 is in use. Change `PORT` in `.env`, or find the process: `Get-NetTCPConnection -LocalPort 8080` |
| "Process Exit with Code: 1" in Air output on reload | Normal on Windows: it is the old server being stopped. |
| `go.mod` suddenly lists many unrelated packages | A `go get` / `go install` was run inside the project. Run `go mod tidy`. |
