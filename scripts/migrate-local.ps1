# scripts/migrate-local.ps1

Write-Host "Running migrations on LOCAL database..."

# Load environment variables from .env.local
if (Test-Path ".env.local") {
    Get-Content .env.local | ForEach-Object {
        if ($_ -and $_ -notmatch "^\s*#") {
            $parts = $_ -split "=", 2
            if ($parts.Count -eq 2) {
                [System.Environment]::SetEnvironmentVariable($parts[0], $parts[1], "Process")
            }
        }
    }
} else {
    Write-Error ".env.local not found!"
    exit 1
}

# Resolve absolute path and convert backslashes to forward slashes
$migrationsPath = (Resolve-Path "./migrations").Path -replace '\\','/'

# Build migrate file URL
$migrationsUrl = "file:///$migrationsPath"

# Run migrations
migrate -path $migrationsUrl -database $env:DATABASE_URL up

Write-Host "Local migrations completed!"
