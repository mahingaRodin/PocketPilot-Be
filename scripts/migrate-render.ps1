
Write-Host "Running migrations on RENDER database..."

# Set the database URL
$RENDER_DB_URL = "postgresql://pocket_pilot_db_user:iCPdG1GDELK3NNfKSzUfuXpBr7gJAefP@dpg-d4ikn5ali9vc73ekjf3g-a.oregon-postgres.render.com/pocket_pilot_db"

# Resolve absolute path to migrations folder and convert backslashes to forward slashes
$migrationsPath = (Resolve-Path "./migrations").Path -replace '\\','/'

# Build migrate file URL
$migrationsUrl = "file:///$migrationsPath"

# Run migrations
migrate -path $migrationsUrl -database $RENDER_DB_URL up

Write-Host "Render migrations completed!"