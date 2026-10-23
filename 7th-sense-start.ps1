# Runs the server the official way: a 32-bit build (GOARCH=386), because the CNC DLL in Lib\
# is 32-bit. Arguments are passed to the server. For hot reload use: go tool air
$ErrorActionPreference = 'Stop'
Set-Location $PSScriptRoot
$env:GOARCH = '386'
go run ./cmd/server @args
exit $LASTEXITCODE
