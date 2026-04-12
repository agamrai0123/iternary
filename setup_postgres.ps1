# PostgreSQL Setup Script for Itinerary Application
# Run as: powershell -ExecutionPolicy Bypass -File setup_postgres.ps1

param(
    [string]$PostgresPassword = "postgres_secure_2026",
    [string]$AdminPassword = "itinerary_admin_2026"
)

Write-Host "=" -MessageParameter Separator * -NoNewline; Write-Host " Phase B Week 1 - PostgreSQL Setup"
Write-Host ""

# Set PGPASSWORD environment variable to avoid password prompts
$env:PGPASSWORD = $PostgresPassword

# 1. Create database
Write-Host "[1/5] Creating database 'itinerary_production'..." -ForegroundColor Cyan
psql -U postgres -h localhost -c "CREATE DATABASE itinerary_production;" 2>$null
if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ Database created successfully" -ForegroundColor Green
} else {
    Write-Host "✓ Database already exists" -ForegroundColor Yellow
}
Write-Host ""

# 2. Create user
Write-Host "[2/5] Creating user 'itinerary_admin'..." -ForegroundColor Cyan
psql -U postgres -h localhost -c "CREATE USER itinerary_admin WITH PASSWORD '$AdminPassword';" 2>$null
if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ User created successfully" -ForegroundColor Green
} else {
    Write-Host "✓ User already exists" -ForegroundColor Yellow
}
Write-Host ""

# 3. Grant privileges
Write-Host "[3/5] Granting privileges to user..." -ForegroundColor Cyan
psql -U postgres -h localhost -c "GRANT ALL PRIVILEGES ON DATABASE itinerary_production TO itinerary_admin;" 2>$null
Write-Host "✓ Privileges granted" -ForegroundColor Green
Write-Host ""

# 4. Apply schema
Write-Host "[4/5] Applying PostgreSQL schema..." -ForegroundColor Cyan
psql -U postgres -h localhost -d itinerary_production -f migrations/001_initial_schema.sql 2>$null
Write-Host "✓ Schema applied successfully" -ForegroundColor Green
Write-Host ""

# 5. Verify setup  
Write-Host "[5/5] Verifying setup..." -ForegroundColor Cyan
$tables = psql -U postgres -h localhost -d itinerary_production -t -c "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema='public';"
Write-Host "✓ Total tables created: $tables" -ForegroundColor Green
Write-Host ""

# Verify data
psql -U postgres -h localhost -d itinerary_production << 'EOF'
\echo ''
\echo '==== PostgreSQL Setup Complete ===='
\echo ''
\echo 'Database: itinerary_production'
\echo 'User: itinerary_admin'
\echo ''
\echo 'Tables:'
SELECT COUNT(*) as table_count FROM information_schema.tables WHERE table_schema='public';
\echo ''
\echo 'Indexes:'
SELECT COUNT(*) as index_count FROM pg_indexes WHERE schemaname='public';
\echo ''
EOF

# Clear password variable
Remove-Item env:PGPASSWORD

Write-Host "✅ PostgreSQL setup completed successfully!" -ForegroundColor Green
Write-Host ""
Write-Host "Connection Details:" -ForegroundColor Cyan
Write-Host "  Host: localhost"
Write-Host "  Port: 5432"
Write-Host "  Database: itinerary_production"
Write-Host "  User: itinerary_admin"
Write-Host "  Password: itinerary_admin_2026"
Write-Host ""
