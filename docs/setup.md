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

Tools installed with `go install` (like Air) go into `%GOPATH%\bin`.

- **Windows:** System Properties → Environment Variables → User variables → `Path` → New →
  `C:\Users\<you>\go\bin`. Open a new terminal afterwards.
- **macOS/Linux:** add to `~/.zshrc` or `~/.bashrc`:
  ```sh
  export PATH="$PATH:$(go env GOPATH)/bin"
  ```

## 4. Install Air (hot reload)

Run this from **outside** the project folder, so it does not touch the project's `go.mod`:

```powershell
cd ~
go install github.com/air-verse/air@latest
air -v
```

If `air` is "not recognized as an internal or external command", step 3 was not applied:
check that `air.exe` exists in `C:\Users\<you>\go\bin` and that this folder is on PATH, then
open a new terminal.

## 5. Install PostgreSQL and create the database

Download from https://www.postgresql.org/download/ and install (remember the password you set
for the `postgres` user). Make sure the `bin` folder (e.g. `C:\Program Files\PostgreSQL\18\bin`)
is on PATH, then:

```powershell
pg_isready                                   # localhost:5432 - accepting connections
psql -U postgres -c "CREATE DATABASE \"7th_sense\";"
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
air
```

Without hot reload:

```powershell
go run ./cmd/server
```

Check it is up:

```powershell
curl http://localhost:8080/health
# {"statusCode":200,"data":{"message":"OK"}}
```

Stop with `Ctrl+C`.

## Troubleshooting

| Problem | Fix |
|---------|-----|
| `air` not recognized | Step 3 (PATH) and step 4; open a new terminal. |
| `air.exe` is in `go\bin\windows_386\` | `GOARCH=386` is set, see step 2. Delete that folder and reinstall Air. |
| `database connection failed: ... password authentication failed` | Wrong `DB_USER` / `DB_PASSWORD` in `.env`. |
| `database connection failed: ... database "7th_sense" does not exist` | Create it (step 5) or fix `DB_DATABASE`. |
| `database connection failed: ... connect: connection refused` | PostgreSQL is not running: `pg_isready`, or start the service in `services.msc`. |
| `bind: Only one usage of each socket address` | Port 8080 is in use. Change `PORT` in `.env`, or find the process: `Get-NetTCPConnection -LocalPort 8080` |
| "Process Exit with Code: 1" in Air output on reload | Normal on Windows: it is the old server being stopped. |
| `go.mod` suddenly lists many unrelated packages | A `go get` / `go install` was run inside the project. Run `go mod tidy`. |
