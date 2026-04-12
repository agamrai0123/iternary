# Triply Project - Quick Reference Card for Developers

## 🚀 Common Commands

### Setup & Running
```bash
cd itinerary-backend
go mod download              # Download dependencies
go run main.go              # Run development server
go build                    # Compile application
./itinerary-backend.exe     # Run compiled build (Windows)
```

### Testing
```bash
go test ./itinerary              # Run all tests
go test -v ./itinerary           # Run tests with verbose output
go test -cover ./itinerary       # Run with coverage report
go test -run TestName ./itinerary # Run specific test function
```

### Database
```bash
sqlite3 itinerary.db        # Open database in SQL shell
.tables                     # List all tables (in sqlite3)
SELECT * FROM users;        # View data
.schema users               # View table structure
.exit                       # Exit sqlite3
```

### Debugging
```bash
# View logs
tail -f log/itinerary-2026-04-12.log

# Make API request
curl http://localhost:8080/api/destinations
curl -v http://localhost:8080/api/destinations    # With headers
curl -H "Authorization: Bearer TOKEN" http://...  # With auth

# Search for error in logs
grep ERROR log/itinerary-*.log
```

---

## 📁 File Structure at a Glance

```
itinerary-backend/
├── main.go                 # Entry point - initializes everything
├── itinerary/
│   ├── routes.go          # API endpoint definitions
│   ├── handlers.go        # HTTP request/response handling
│   ├── service.go         # Business logic
│   ├── database.go        # Database operations
│   ├── models.go          # Data structures
│   ├── auth_*.go          # Authentication code
│   ├── group_*.go         # Group trip features
│   ├── config.go          # Configuration loading
│   ├── logger.go          # Logging setup
│   ├── error.go           # Error handling
│   ├── metrics.go         # Monitoring/metrics
│   └── *_test.go          # Unit tests
├── config/config.json     # Configuration file
├── templates/             # HTML pages
├── static/                # CSS, JS, images
├── docs/                  # Documentation
├── go.mod & go.sum        # Dependency management
└── itinerary.db          # SQLite database (auto-created)
```

---

## 🏗️ Architecture Layers (Data Flow)

```
REQUEST COMES IN
        ↓
┌───────────────────────────────────────────────┐
│ HANDLER LAYER (handlers.go)                  │
│ • Parse request parameters                   │
│ • Call service methods                       │
│ • Return JSON/HTML responses                 │
└────────────┬────────────────────────────────┘
             ↓
┌───────────────────────────────────────────────┐
│ SERVICE LAYER (service.go)                   │
│ • Validate inputs                            │
│ • Business logic                             │
│ • Call database methods                      │
│ • Log operations                             │
└────────────┬────────────────────────────────┘
             ↓
┌───────────────────────────────────────────────┐
│ DATABASE LAYER (database.go)                 │
│ • Execute SQL queries                        │
│ • Return data or errors                      │
└────────────┬────────────────────────────────┘
             ↓
        [Database]
```

---

## 🔑 Key Code Patterns

### Creating an HTTP Handler

```go
func (h *Handlers) GetDestinations(c *gin.Context) {
    // 1. Get parameters
    page := c.Query("page")
    
    // 2. Call service
    destinations, err := h.service.GetDestinations(page)
    
    // 3. Handle error
    if err != nil {
        c.JSON(500, gin.H{"error": "internal server error"})
        return
    }
    
    // 4. Return success
    c.JSON(200, gin.H{"data": destinations})
}
```

### Adding a Route

```go
// In routes.go
// Public route
router.GET("/api/destinations", handlers.GetDestinations)

// Protected route (requires auth)
router.POST("/api/user-trips", 
    authMiddleware.RequireAuth(),  // Add this middleware
    handlers.CreateUserTrip)
```

### Creating a Service Method

```go
func (s *Service) GetDestinations(page int) ([]Destination, error) {
    // Validate
    if page < 1 {
        return nil, fmt.Errorf("invalid page")
    }
    
    // Call database
    destinations, err := s.db.GetDestinations(page)
    if err != nil {
        s.logger.Error("database_error", "err", err.Error())
        return nil, err
    }
    
    return destinations, nil
}
```

### Database Query

```go
func (db *Database) GetDestinations() ([]Destination, error) {
    rows, err := db.conn.Query("SELECT id, name, country FROM destinations")
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var destinations []Destination
    for rows.Next() {
        var d Destination
        err := rows.Scan(&d.ID, &d.Name, &d.Country)
        if err != nil {
            return nil, err
        }
        destinations = append(destinations, d)
    }
    return destinations, nil
}
```

### Logging

```go
logger.Info("user_logged_in", "user_id", "123")
logger.Error("database_error", "error", err.Error())
logger.Debug("fetching_data", "page", 1, "size", 10)
logger.Warn("deprecated_endpoint", "endpoint", "/old-api")
```

### Error Handling

```go
// Create error
apiErr := NewNotFoundError("user not found")

// Return in handler
c.JSON(apiErr.StatusCode, apiErr.ToJSON())

// Common errors
NewNotFoundError(msg)        // 404
NewUnauthorizedError(msg)    // 401
NewBadRequestError(msg)      // 400
NewConflictError(msg)        // 409
```

---

## 🗄️ Common Database Queries

### User Operations
```sql
-- Get user by ID
SELECT * FROM users WHERE id = 'user_123';

-- List all users
SELECT id, username, email FROM users;

-- Create user (handled by Go code)
INSERT INTO users (id, username, email) VALUES ('...', '...', '...');
```

### Destination Operations
```sql
-- Get all destinations
SELECT * FROM destinations;

-- Get specific destination
SELECT * FROM destinations WHERE id = 'dest_456';
```

### Itinerary Operations
```sql
-- Get itineraries for destination
SELECT * FROM itineraries WHERE destination_id = 'dest_456';

-- Get itinerary items by day
SELECT * FROM itinerary_items 
WHERE itinerary_id = 'iter_789' 
ORDER BY day;
```

### Group Trip Operations
```sql
-- Get group trip details
SELECT * FROM group_trips WHERE id = 'group_123';

-- Get group members
SELECT u.*, gm.role FROM group_members gm
JOIN users u ON gm.user_id = u.id
WHERE gm.group_trip_id = 'group_123';
```

---

## 🔒 Authentication Pattern

### Login Flow
```
User enters credentials
    ↓
POST /api/auth/login
    ↓
Handler validates credentials
    ↓
AuthService generates token
    ↓
Store session in database
    ↓
Return token to client
    ↓
Client stores token (cookie/localStorage)
```

### Using Protected Endpoint
```
Client includes token in header:
Authorization: Bearer <token>
    ↓
Request hits authMiddleware.RequireAuth()
    ↓
Middleware validates token
    ↓
If valid: set user_id in context, call handler
If invalid: return 401 Unauthorized
```

### Getting User ID in Handler
```go
func (h *Handlers) GetUserTrips(c *gin.Context) {
    userID := c.GetString("user_id")  // From middleware
    // Now filter trips for this user
}
```

---

## 📊 API Response Format

### Success Response
```json
{
  "data": [
    {"id": "123", "name": "Paris", ...}
  ],
  "total": 50,
  "page": 1,
  "page_size": 10
}
```

### Error Response
```json
{
  "code": "NOT_FOUND",
  "message": "Destination not found",
  "status_code": 404,
  "details": "No destination with ID dest_999"
}
```

---

## 🔄 Request/Response Cycle

### GET Request (Get Data)
```bash
curl http://localhost:8080/api/destinations?page=1

Handler: GetDestinations()
  ↓ Parse page param
  ↓ Call service.GetDestinations()
  ↓ Log operation
  ↓ Return JSON response
```

### POST Request (Create Data)
```bash
curl -X POST http://localhost:8080/api/itineraries \
  -H "Content-Type: application/json" \
  -d '{"title": "...", "budget": 5000}'

Handler: CreateItinerary()
  ↓ Parse JSON body into struct
  ↓ Call service.CreateItinerary()
  ↓ Service validates
  ↓ Call database.CreateItinerary()
  ↓ Return created object in JSON
```

---

## 🧪 Testing Pattern

```go
func TestGetDestinations(t *testing.T) {
    // 1. ARRANGE - Setup
    mockDB := &mockDatabase{}
    service := NewService(mockDB, logger)
    
    // 2. ACT - Call function
    destinations, err := service.GetDestinations(1)
    
    // 3. ASSERT - Verify results
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }
    if len(destinations) == 0 {
        t.Errorf("Expected destinations, got none")
    }
}
```

Run test:
```bash
go test -run TestGetDestinations -v ./itinerary
```

---

## 📝 Adding New Feature Checklist

- [ ] Create model struct (in models.go)
- [ ] Create database method (in database.go)
- [ ] Create service method (in service.go)
- [ ] Create handler (in handlers.go)
- [ ] Add route (in routes.go)
- [ ] Add unit tests
- [ ] Test with curl
- [ ] Check logs for errors
- [ ] Update documentation

---

## 🚨 Common Errors & Fixes

| Error | Cause | Fix |
|-------|-------|-----|
| `database is locked` | Concurrent access | Ensure single DB connection |
| `connection refused` | Server not running | Run `go run main.go` |
| `json: cannot unmarshal` | Wrong JSON format | Check request body structure |
| `404 Not Found` | Route doesn't exist | Verify route in routes.go |
| `401 Unauthorized` | Missing auth token | Add `Authorization: Bearer TOKEN` |
| `undefined: SomeType` | Missing import | Add `import "package/name"` |
| `cannot find module` | Dependencies missing | Run `go mod download` |

---

## 🎯 File Modification Workflow

### When you need to...

**Add a new API endpoint:**
1. handlers.go - add handler function
2. routes.go - add route mapping
3. Test with curl

**Fix a bug in business logic:**
1. Locate service method
2. Add debug logging
3. Fix logic
4. Run tests
5. Restart server

**Add a database field:**
1. Update models.go
2. Modify database.go queries
3. Update handlers/services that use it
4. Add database migration

**Change authentication:**
1. Modify auth.go / auth_service.go
2. Update auth_middleware.go
3. Update routes.go for affected routes
4. Test authentication flow

---

## 🔍 Debugging Workflow

```
Issue Found
    ↓
1. Check logs: tail -f log/itinerary-*.log
    ↓
2. Add debug logging: logger.Debug("msg", fields...)
    ↓
3. Test endpoint: curl -v http://...
    ↓
4. Check database: sqlite3 itinerary.db
    ↓
5. Verify code logic: trace through handler → service → database
    ↓
Issue Fixed ✓
```

---

## 🌍 Environment Setup

### config/config.json
```json
{
  "server": {
    "port": "8080"          // Change port here
  },
  "logging": {
    "level": "info"         // Change to "debug" for verbose logs
  }
}
```

### .env File
```
DB_PATH=itinerary.db
LOG_LEVEL=info
SERVER_PORT=8080
```

---

## ⚡ Performance Tips

1. **Add indexes** to frequently queried columns
2. **Limit result sets** - use pagination (page, page_size)
3. **Cache** destinations (rarely change)
4. **Use connection pooling** (already configured)
5. **Monitor logs** for slow queries

---

## 📚 Key Type Definitions

```go
// User in the system
type User struct {
    ID       string
    Username string
    Email    string
}

// Travel destination
type Destination struct {
    ID          string
    Name        string
    Country     string
}

// Complete travel plan
type Itinerary struct {
    ID       string
    Title    string
    Budget   float64
    Items    []ItineraryItem  // Activities, stays, etc
}

// Single activity/stay in itinerary
type ItineraryItem struct {
    Day    int      // Day number
    Type   string   // stay, food, activity, etc
    Name   string   // Activity name
    Price  float64  // Cost
}

// Collaborative group trip
type GroupTrip struct {
    ID      string         // Unique ID
    Title   string         // Trip name
    Owner   User           // Who created it
    Members []GroupMember  // Who's participating
}

// User's membership in a group
type GroupMember struct {
    UserID string  // User ID
    Role   string  // owner, editor, member, viewer
    Status string  // pending, active, declined
}
```

---

## 🔗 Useful Links

- Go Documentation: https://golang.org/doc/
- Gin Framework: https://github.com/gin-gonic/gin
- SQLite Documentation: https://www.sqlite.org/docs.html
- REST API Best Practices: https://restfulapi.net/

---

## 📞 Quick Help

- **Server won't start?** → Check port 8080 isn't in use: `netstat -ano | findstr :8080`
- **Tests failing?** → Run `go test -v ./itinerary` to see detailed errors
- **Database errors?** → Verify itinerary.db file exists and is writable
- **Can't find something?** → Use `grep -r "SearchTerm" itinerary/` to search code
- **Want to see all routes?** → Search for `router.GET`, `router.POST` in routes.go

---

**Created:** April 2026  
**For:** Quick reference during development  
**Keep this handy!** 📌
