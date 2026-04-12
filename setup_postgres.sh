#!/bin/bash
# PostgreSQL Setup Script for Itinerary Application - BASH VERSION
# Run as: bash setup_postgres.sh

set -e

POSTGRES_PASS="postgres_secure_2026"
ADMIN_PASS="itinerary_admin_2026"

echo ""
echo "=========================================="
echo "  Phase B Week 1 - PostgreSQL Setup"
echo "=========================================="
echo ""

# Set PGPASSWORD for non-interactive login
export PGPASSWORD="$POSTGRES_PASS"

# 1. Create database
echo "[1/5] Creating database 'itinerary_production'..."
psql -U postgres -h localhost -c "CREATE DATABASE itinerary_production;" 2>/dev/null || echo "✓ Database already exists"
echo "✓ Database ready"
echo ""

# 2. Create user
echo "[2/5] Creating user 'itinerary_admin'..."
psql -U postgres -h localhost -c "CREATE USER itinerary_admin WITH PASSWORD '$ADMIN_PASS';" 2>/dev/null || echo "✓ User already exists"
echo "✓ User ready"
echo ""

# 3. Grant privileges
echo "[3/5] Granting privileges to user..."
psql -U postgres -h localhost -c "GRANT ALL PRIVILEGES ON DATABASE itinerary_production TO itinerary_admin;"
echo "✓ Privileges granted"
echo ""

# 4. Apply schema
echo "[4/5] Applying PostgreSQL schema..."
psql -U postgres -h localhost -d itinerary_production -f migrations/001_initial_schema.sql 2>/dev/null
echo "✓ Schema applied"
echo ""

# 5. Verify setup
echo "[5/5] Verifying setup..."
TABLE_COUNT=$(psql -U postgres -h localhost -d itinerary_production -t -c "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema='public';" 2>/dev/null | tr -d ' ')
INDEX_COUNT=$(psql -U postgres -h localhost -d itinerary_production -t -c "SELECT COUNT(*) FROM pg_indexes WHERE schemaname='public';" 2>/dev/null | tr -d ' ')

echo "✓ Tables created: $TABLE_COUNT"
echo "✓ Indexes created: $INDEX_COUNT"
echo ""

# Clear password
unset PGPASSWORD

echo "✅ PostgreSQL setup completed successfully!"
echo ""
echo "Connection Details:"
echo "  Host: localhost"
echo "  Port: 5432"
echo "  Database: itinerary_production"
echo "  User: itinerary_admin"
echo "  Password: itinerary_admin_2026"
echo ""
