
param(
    [int]$Values = 40000,  # telemetry values per second
    [int]$Seconds = 120,   # how long to send
    [int]$Rate = 1,       # readings per machine per second; machines = Values / (40 x Rate)
    [int]$Change = 100    # percent of values that change per reading; 100 = worst case, 20 = realistic
)
Set-Location $PSScriptRoot

if (Test-Path .env) {
    Write-Host "Loading .env:"
    foreach ($line in Get-Content .env) {
        $line = $line.Trim()
        if ($line -eq '' -or $line.StartsWith('#')) { continue }
        $eq = $line.IndexOf('=')
        if ($eq -lt 1) { continue }
        $name = $line.Substring(0, $eq).Trim()
        $value = $line.Substring($eq + 1).Trim()
        # Strip an unquoted inline comment: KEY=value # note (a quoted value keeps its #).
        if (-not ($value.StartsWith('"') -or $value.StartsWith("'"))) {
            $hash = $value.IndexOf(' #')
            if ($hash -ge 0) { $value = $value.Substring(0, $hash) }
            $value = $value.Trim()
        }
        $value = $value.Trim('"').Trim("'")
        # .env wins over anything already in this session: Process-scope variables set by an
        # earlier run would otherwise silently shadow later .env edits.
        $existing = [Environment]::GetEnvironmentVariable($name)
        [Environment]::SetEnvironmentVariable($name, $value, 'Process')
        $overridden = (-not [string]::IsNullOrEmpty($existing) -and $existing -ne $value)
        # Mask credentials: anything with PASSWORD/SECRET/TOKEN/KEY, and the password inside URLs.
        $shown = $value
        if ($name -match 'PASSWORD|SECRET|TOKEN|KEY') {
            $shown = '*****'
        } elseif ($value -match '://[^:/]+:[^@]+@') {
            $shown = $value -replace '(://[^:/]+:)[^@]+(@)', '$1*****$2'
        }
        if ($overridden) { Write-Host "  $name = $shown  (overrode session value)" } else { Write-Host "  $name = $shown" }
    }
} else {
    Write-Host "No .env file found, using session environment variables."
}

$suffix = if ($Change -lt 100) { "-c$Change" } else { '' }
$log = "logs\load-test-$Values$suffix-$(Get-Date -Format yyyyMMdd-HHmmss).log"
New-Item -ItemType Directory -Force logs | Out-Null

$env:TELEMETRY_LOAD = '1'
$env:TELEMETRY_LOAD_VALUES = "$Values"
$env:TELEMETRY_LOAD_SECONDS = "$Seconds"
$env:TELEMETRY_LOAD_RATE = "$Rate"
$env:TELEMETRY_LOAD_PHASES = 'target'
$env:TELEMETRY_LOAD_TARGET_CHANGE = ([double]$Change / 100).ToString('0.####', [Globalization.CultureInfo]::InvariantCulture)
try {
    # cmd merges the test's log lines (stderr) with its results, without PowerShell error records.
    cmd /c "go test ./tests -run TestTelemetryLoad$ -count=1 -v -timeout 120m 2>&1" | Tee-Object -FilePath $log
    $code = $LASTEXITCODE
}
finally {
    Remove-Item Env:TELEMETRY_LOAD, Env:TELEMETRY_LOAD_VALUES, Env:TELEMETRY_LOAD_SECONDS, Env:TELEMETRY_LOAD_RATE, Env:TELEMETRY_LOAD_PHASES, Env:TELEMETRY_LOAD_TARGET_CHANGE -ErrorAction SilentlyContinue
}
Write-Host "Saved: $log"
exit $code
