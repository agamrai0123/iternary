# Itinerary Platform - Oracle Database Setup Guide

## Overview

This guide covers setting up the Oracle database for the Itinerary Platform. The project connects to Oracle XE using the `godror` driver.

## Prerequisites

- Oracle Database 12c+ (XE works perfectly)
- Database running on `localhost:1521` with service `XE`
- System user with password configured
- Go 1.21+ installed

## Quick Setup (Recommended)

### 1. Set Environment Variables

```bash
# macOS/Linux
export DB_PASSWORD=your_oracle_password
export DB_HOST=localhost

# Windows (PowerShell)
$env:DB_PASSWORD="your_oracle_password"
$env:DB_HOST="localhost"

# Windows (Command Prompt)
set DB_PASSWORD=your_oracle_password
set DB_HOST=localhost
```

### 2. Run Database Initialization

```bash
cd itinerary-backend

# First time: download dependencies
go mod download

# Initialize database with schema and test data
go run init_db.go init
```

Expected output:
```
✓ Connected to Oracle database
✓ Executed statement 1
✓ Executed statement 2
...
✓ Database schema initialized successfully!

📊 Database Summary:
  Users: 3
  Destinations: 3
  Itineraries: 4
  Itinerary Items: 10
  Comments: 3
```

### 3. Verify Setup

```bash
go run init_db.go verify
```

Should show:
```
✓ Connected to Oracle database

📊 Table Counts:
  Users: 3 rows
  Destinations: 3 rows
  Itineraries: 4 rows
  Itinerary Items: 10 rows
  Comments: 3 rows

📋 Sample Data:
  Users:
    - traveler1 (user-001) <traveler1@example.com>
    - explorer2 (user-002) <explorer@example.com>
    - wanderer3 (user-003) <wanderer@example.com>

  Destinations:
    - Goa, India (dest-001)
    - Manali, India (dest-002)
    - Bali, Indonesia (dest-003)

  Itineraries:
    - 5-Day Budget Goa (5 days, ₹15000) [itin-001]
    - Luxury 7-Day Goa Experience (7 days, ₹45000) [itin-002]
    - 4-Day Manali Adventure (4 days, ₹12000) [itin-003]
    - 6-Day Bali Paradise (6 days, ₹18000) [itin-004]
```

### 4. Start Application

```bash
go run main.go
```

Visit `http://localhost:8080`

---

## Manual SQL Setup

If you prefer to set up the database manually using SQL*Plus or SQL Developer:

### Using SQL*Plus

```bash
sqlplus system/password@localhost:1521/XE
SQL> @docs/schema.sql
SQL> COMMIT;
SQL> EXIT;
```

### Using SQL Developer

1. Create new connection:
   - Connection Name: `LocalXE`
   - Username: `system`
   - Password: `your_oracle_password`
   - Hostname: `localhost`
   - Port: `1521`
   - Service Name: `XE`

2. Open File → Open `docs/schema.sql`

3. Run Script (Ctrl+Enter or Cmd+Enter)

4. Verify with:
   ```sql
   SELECT COUNT(*) FROM users;
   SELECT * FROM destinations;
   ```

---

## Database Management Commands

### Initialize Database
```bash
go run init_db.go init
```
- Creates all tables (with CASCADE CONSTRAINTS)
- Inserts test data (3 users, 3 destinations, 4 itineraries, etc.)
- Idempotent - safe to run multiple times

### Verify Current Data
```bash
go run init_db.go verify
```
- Shows row counts for all tables
- Displays sample data (users, destinations, itineraries)
- Useful for checking database health

### Clean Database (Drop All Tables)
```bash
go run init_db.go clean
```
- Drops all tables with CASCADE CONSTRAINTS
- Requires confirmation
- Useful for full reset before re-initialization

---

## Configuration

### config/config.json
```json
{
  "server": {
    "port": ":8080"
  },
  "database": {
    "user": "system",
    "host": "localhost",
    "port": "1521",
    "service": "XE"
  },
  "logging": {
    "file": "logs/app.log",
    "level": "info"
  }
}
```

### Environment Variable Overrides
- `DB_PASSWORD` - Overrides config file (system password)
- `DB_HOST` - Overrides config file (hostname)

These are loaded in `itinerary/config.go` with this precedence:
1. config.json (base)
2. Environment variables (override)

---

## Test Data Included

After initialization, your database includes:

### Users (3)
- `traveler1` - traveler1@example.com
- `explorer2` - explorer@example.com
- `wanderer3` - wanderer@example.com

### Destinations (3)
- **Goa** (India) - Beaches, heritage, nightlife
- **Manali** (India) - Mountains, adventure, trekking
- **Bali** (Indonesia) - Beaches, temples, rice paddies

### Itineraries (4)
1. **5-Day Budget Goa** (₹15,000) - 5 days, 45 likes
   - Beachside accommodation
   - Jet ski & water sports
   - Old Goa heritage tour
   - Scooter rental exploration

2. **Luxury 7-Day Goa** (₹45,000) - 7 days, 32 likes
   - 5-star resort, spa, sunset cruises
   - Upscale dining and activities

3. **4-Day Manali Adventure** (₹12,000) - 4 days, 28 likes
   - Trekking and paragliding
   - Rohtang Pass & local villages
   - Adventure-focused

4. **6-Day Bali Paradise** (₹18,000) - 6 days, 67 likes
   - Comprehensive experience
   - Temples, beaches, rice terraces
   - Great for first-time visitors

### Itinerary Items (10+)
Each itinerary has daily breakdown with:
- Accommodations (stays)
- Restaurants (food)
- Activities & attractions
- Transportation

All with pricing, duration, location, and ratings

### Comments (3)
Sample reviews with 4-5 star ratings on popular itineraries

---

## Working with the Database

### Adding More Test Data (Manually)

Using SQL*Plus or SQL Developer:

```sql
-- Add a new destination
INSERT INTO destinations (id, name, country, description, created_at, updated_at)
VALUES ('dest-004', 'Kerala', 'India', 'Backwaters, beaches and spices', SYSDATE, SYSDATE);
COMMIT;

-- Add a new itinerary
INSERT INTO itineraries (id, user_id, destination_id, title, description, duration, budget, likes, created_at, updated_at)
VALUES ('itin-005', 'user-001', 'dest-004', '3-Day Kerala Backwaters', 
        'Relax on houseboats and explore backwaters', 3, 10000, 0, SYSDATE, SYSDATE);
COMMIT;

-- Add itinerary items
INSERT INTO itinerary_items (id, itinerary_id, day, type, name, description, price, duration, location, created_at, updated_at)
VALUES ('item-201', 'itin-005', 1, 'stay', 'Houseboat Stay', 'Traditional Kerala houseboat', 2000, 24, 'Kumarakom', SYSDATE, SYSDATE);
COMMIT;
```

### Exploring Data

```sql
-- View all destinations with itinerary count
SELECT d.id, d.name, d.country, COUNT(i.id) as itinerary_count
FROM destinations d
LEFT JOIN itineraries i ON d.id = i.destination_id
GROUP BY d.id, d.name, d.country
ORDER BY itinerary_count DESC;

-- View most liked itineraries
SELECT id, title, budget, likes, duration
FROM itineraries
ORDER BY likes DESC;

-- View items for a specific itinerary
SELECT day, type, name, price, location
FROM itinerary_items
WHERE itinerary_id = 'itin-001'
ORDER BY day, type;

-- View reviews for an itinerary
SELECT u.username, c.content, c.rating, c.created_at
FROM comments c
JOIN users u ON c.user_id = u.id
WHERE c.itinerary_id = 'itin-001'
ORDER BY c.created_at DESC;
```

### Modifying Test Data

```sql
-- Update an itinerary title
UPDATE itineraries
SET title = 'Budget Goa: 5 Days Paradise'
WHERE id = 'itin-001';
COMMIT;

-- Update likes count
UPDATE itineraries
SET likes = likes + 1
WHERE id = 'itin-004';
COMMIT;

-- Update destination description
UPDATE destinations
SET description = 'Updated description for Goa'
WHERE id = 'dest-001';
COMMIT;
```

---

## Troubleshooting

### Error: ORA-12514 - TNS:listener could not resolve SERVICE_NAME

**Cause**: Oracle is not running or service name is wrong

**Solutions**:
1. Start Oracle:
   ```bash
   # macOS (using Homebrew)
   brew services start oracle
   
   # Linux/Windows: Use Oracle services or Command Line Tools
   ```

2. Verify service name:
   ```sql
   -- In SQL*Plus, check available services
   SELECT name FROM v$services;
   ```

3. Check `config/config.json` has correct service: `"service": "XE"`

### Error: ORA-01017 - invalid username/password

**Cause**: Wrong system password or DB_PASSWORD environment variable not set

**Solutions**:
```bash
# Set environment variable
export DB_PASSWORD=correct_password

# Verify it's set
echo $DB_PASSWORD

# Or reset Oracle password
sqlplus / as sysdba
SQL> ALTER USER system IDENTIFIED BY newpassword;
SQL> COMMIT;
```

### Error: Tables already exist

**First run** when tables already exist - this is OK:
```
ℹ Statement ...: Table/constraint already exists (OK)
```

**To reset completely**:
```bash
go run init_db.go clean
go run init_db.go init
```

### Port 8080 Already in Use

```bash
# macOS/Linux
lsof -i :8080
kill -9 <PID>

# Windows (PowerShell)
netstat -ano | findstr :8080
taskkill /PID <PID> /F
```

Or change port in `config/config.json`:
```json
{
  "server": {
    "port": ":8081"
  }
}
```

### Connection Timeout

Check these:
1. Oracle is running: `sqlplus system/password@localhost:1521/XE`
2. Network accessibility: `ping localhost`
3. Port is open: `telnet localhost 1521`
4. Firewall not blocking port 1521

### No Data Showing in Web UI

- Run: `go run init_db.go verify`
- Check logs: `tail -f logs/app.log`
- Verify data exists: `SELECT COUNT(*) FROM destinations;` in SQL*Plus

---

## Database Schema Overview

### Core Tables

- **users** - User accounts and authentication
- **destinations** - Travel destinations
- **itineraries** - Travel plans

- **itinerary_items** - Daily activities/attractions in each itinerary
- **comments** - User reviews on itineraries
- **user_plans** - Saved/copied itineraries
- **likes** - Track which users like which itineraries

### Data Types

Oracle-specific types used:
- `VARCHAR2(n)` - Variable-length strings (up to 4000)
- `CLOB` - Large text (descriptions, content)
- `NUMBER(p,s)` - Decimal numbers (budget, price, rating)
- `TIMESTAMP` - Date + time
- `SYSDATE` - Current date/time (Oracle function)

### Indexes

Created for performance:
- Primary keys on all tables
- Foreign key indexes (user_id, destination_id, itinerary_id)
- Query optimization (country, email, username, type, created_at)

---

## Next Steps

1. ✅ Database initialized and populated
2. ✅ Application can fetch and display data
3. → Explore the web interface at http://localhost:8080
4. → Add new destinations and itineraries
5. → Modify test data as needed
6. → Check [TEMPLATES_GUIDE.md](TEMPLATES_GUIDE.md) for frontend customization
7. → See [QUICK_START.md](QUICK_START.md) for general project setup

---

## Additional Resources

- [Oracle Database Documentation](https://docs.oracle.com/database/)
- [godror Driver Documentation](https://github.com/godror/godror)
- [SQL Developer Documentation](https://www.oracle.com/database/sqldeveloper/)
