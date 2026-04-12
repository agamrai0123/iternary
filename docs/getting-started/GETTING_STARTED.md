# 🚀 Getting Started - 3 Simple Steps

Your Travel Itinerary Platform is ready to run! Here's how to get started:

## Step 1️⃣ Set Your Database Password

Open your terminal and set the Oracle password:

**macOS/Linux:**
```bash
export DB_PASSWORD=your_oracle_system_password
export DB_HOST=localhost
```

**Windows (PowerShell):**
```powershell
$env:DB_PASSWORD="your_oracle_system_password"
$env:DB_HOST="localhost"
```

**Windows (Command Prompt):**
```cmd
set DB_PASSWORD=your_oracle_system_password
set DB_HOST=localhost
```

## Step 2️⃣ Initialize the Database

Navigate to the project and initialize the database:

```bash
cd itinerary-backend
go run init_db.go init
```

You'll see output like:
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

## Step 3️⃣ Start the Application

```bash
go run main.go
```

Wait for:
```
[GIN-debug] Listening and serving HTTP on :8080
```

Then open in your browser: **http://localhost:8080**

---

## ✨ What You'll See

### Home Page (`/`)
- Grid of destinations (Goa, Manali, Bali)
- Featured itineraries
- Browse and search functionality

### Browse Destinations
- Click any destination to see itineraries
- View all itineraries created by travelers
- See popularity (likes) and budget

### View Itineraries
- Complete day-by-day breakdown
- All activities, meals, hotels with costs
- Total trip budget calculated
- Like and comment functionality

### Create New Plan
- Form to build your own itinerary
- Add activities, meals, transport
- Calculate total cost automatically

### Search & Filter
- Find itineraries by destination
- Filter by budget range
- Sort by popularity or recency

---

## 📊 Sample Data Included

The database is pre-populated with realistic test data:

| | Count | Details |
|---|-------|---------|
| Users | 3 | With comments and likes |
| Destinations | 3 | Goa, Manali, Bali |
| Itineraries | 4 | From ₹12K to ₹45K budgets |
| Activities | 10+ | Hotels, meals, tours, transport |
| Comments | 3 | With 4-5 star ratings |

---

## 🛠️ Database Management

### Verify Data
```bash
go run init_db.go verify
```

### Reset Everything
```bash
go run init_db.go clean
go run init_db.go init
```

### Direct Database Access
```bash
sqlplus system/password@localhost:1521/XE
SQL> SELECT * FROM destinations;
```

---

## 📖 Documentation

For more details, see:

| File | Purpose |
|------|---------|
| [SETUP_SUMMARY.md](SETUP_SUMMARY.md) | Complete project overview |
| [docs/DATABASE_SETUP.md](docs/DATABASE_SETUP.md) | Detailed database configuration |
| [docs/TEMPLATES_GUIDE.md](docs/TEMPLATES_GUIDE.md) | HTML template reference |
| [README.md](README.md) | Project information |
| [docs/QUICK_START.md](docs/QUICK_START.md) | Quick reference guide |

---

## 🎯 API Endpoints

### Pages (Browser)
- `GET /` - Home page
- `GET /destination/:id` - Destination detail
- `GET /itinerary/:id` - Itinerary detail
- `GET /create` - Create itinerary form
- `GET /search` - Search interface

### JSON API (Programmatic)
- `GET /api/destinations` - All destinations
- `GET /api/destinations/:id` - One destination
- `GET /api/itineraries` - All itineraries (paginated)
- `GET /api/itineraries/:id` - One itinerary with items
- `POST /api/itineraries/:id/comments` - Add comment
- `POST /api/itineraries/:id/like` - Like an itinerary
- `GET /api/health` - Health check

**Example API Call:**
```bash
curl http://localhost:8080/api/destinations
```

---

## ⚙️ Configuration

The application reads configuration from:

1. **config/config.json** - Default settings:
   ```json
   {
     "server": { "port": ":8080" },
     "database": {
       "user": "system",
       "host": "localhost",
       "port": "1521",
       "service": "XE"
     }
   }
   ```

2. **Environment Variables** (override config.json):
   - `DB_PASSWORD` - Oracle password
   - `DB_HOST` - Database hostname

---

## 🔍 Troubleshooting

| Issue | Solution |
|-------|----------|
| "Connection refused" | Verify Oracle is running on localhost:1521 |
| "ORA-01017" | Check DB_PASSWORD environment variable |
| "Table already exists" | Normal on first init; safe to ignore |
| No data visible | Run `go run init_db.go verify` to check |
| Port 8080 in use | Change port in config/config.json or kill process |

---

## 🎓 Next Steps

1. ✅ **Database initialized** - You're here!
2. ✅ **Application running** - Visit http://localhost:8080
3. → Browse destinations and itineraries
4. → Try adding comments and likes
5. → Create a new itinerary via `/create`
6. → Explore the API endpoints with curl

---

## 📝 Project Structure

```
itinerary-backend/
├── main.go               ← Run this to start server
├── init_db.go            ← Run this to setup database
├── config/config.json    ← Edit for settings
├── itinerary/            ← Go backend code
├── templates/            ← HTML pages (5 files)
├── static/css, static/js ← Styling and scripts
└── docs/                 ← Documentation
```

---

## 💬 Example Usage

### Browse Destinations
```bash
curl http://localhost:8080/api/destinations
```

### Get Itineraries
```bash
curl http://localhost:8080/api/itineraries?page=1&limit=10
```

### Get One Itinerary with Items
```bash
curl http://localhost:8080/api/itineraries/itin-001
```

### Add a Comment
```bash
curl -X POST http://localhost:8080/api/itineraries/itin-001/comments \
  -H "Content-Type: application/json" \
  -d '{"user_id":"user-001","content":"Amazing trip!","rating":5}'
```

### Like an Itinerary
```bash
curl -X POST http://localhost:8080/api/itineraries/itin-001/like \
  -H "Content-Type: application/json" \
  -d '{"user_id":"user-001"}'
```

---

## 🎉 You're All Set!

Your application is now:
- ✅ Connected to Oracle database
- ✅ Loaded with test data
- ✅ Ready to serve web pages
- ✅ Ready to serve as an API

Visit **http://localhost:8080** or start using the APIs!

---

**Need Help?**
- Check [docs/DATABASE_SETUP.md](docs/DATABASE_SETUP.md) for detailed setup
- See [docs/TEMPLATES_GUIDE.md](docs/TEMPLATES_GUIDE.md) for template info
- Review [SETUP_SUMMARY.md](SETUP_SUMMARY.md) for complete project overview
