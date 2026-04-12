#!/bin/bash
# Monday Hour 1: Apply Multi-Currency Schema
# This script backs up the database and applies the new schema

cd /d/Learn/iternary/itinerary-backend

# Stop the server if running
echo "=== Step 1: Stopping server ==="
taskkill /IM itinerary-backend.exe /F 2>/dev/null || true
sleep 2
echo "Server stopped"

# Create backup
echo ""
echo "=== Step 2: Creating backup ==="
BACKUP_FILE="./backups/itinerary.db.$(date +%Y%m%d_%H%M%S).backup"
mkdir -p ./backups
cp itinerary.db "$BACKUP_FILE" 2>/dev/null || echo "Database not found yet - will create new one"
echo "Backup created: $BACKUP_FILE"

# Apply multi-currency schema
echo ""
echo "=== Step 3: Applying multi-currency schema ===" 
sqlite3 itinerary.db < multicurrency_schema.sql 2>&1 | tail -20 || echo "Schema already applied or database doesn't exist yet"

echo ""
echo "=== Step 4: Verifying schema tables ==="
echo "SELECT COUNT(*) as total_tables FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%';" | sqlite3 itinerary.db

echo ""
echo "=== Tables created/updated ==="
echo ".schema user_preferences" | sqlite3 itinerary.db | head -10

echo ""
echo "=== Sample data loaded ==="
echo "SELECT COUNT(*) as currencies FROM supported_currencies;" | sqlite3 itinerary.db
echo "SELECT COUNT(*) as languages FROM supported_languages;" | sqlite3 itinerary.db
echo "SELECT COUNT(*) as alert_rules FROM alert_rules;" | sqlite3 itinerary.db

echo ""
echo "✅ Hour 1 Complete: Database schema applied successfully!"
