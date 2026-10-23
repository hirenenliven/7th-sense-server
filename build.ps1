# Builds the server the official way: a 32-bit build (GOARCH=386), because the CNC DLL in Lib\
# is 32-bit. Output: bin\server.exe; deploy it with the Lib\ folder and the client_cert / client_key files.
$ErrorActionPreference = 'Stop'
Set-Location $PSScriptRoot
$env:GOARCH = '386'
go build -o bin/server.exe ./cmd/server
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
Write-Host 'built bin\server.exe (GOARCH=386)'
