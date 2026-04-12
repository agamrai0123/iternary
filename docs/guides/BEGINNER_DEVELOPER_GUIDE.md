# Triply Project - Beginner Developer Guide

**Project Name:** Triply (Itinerary)  
**Status:** 95% Backend Complete, Ready for Frontend/Testing  
**Language:** Go  
**Last Updated:** April 2026

---

## 📌 Table of Contents

1. [Project Overview](#project-overview)
2. [Technology Stack](#technology-stack)
3. [Architecture Overview](#architecture-overview)
4. [Project Structure](#project-structure)
5. [Key Concepts](#key-concepts)
6. [Development Environment Setup](#development-environment-setup)
7. [Understanding Core Components](#understanding-core-components)
8. [Common Workflows](#common-workflows)
9. [Database Schema](#database-schema)
10. [API Endpoints](#api-endpoints)
11. [Testing Guide](#testing-guide)
12. [Debugging Tips](#debugging-tips)
13. [Troubleshooting](#troubleshooting)

---

## 🎯 Project Overview

**Triply** is a **collaborative trip planning platform** where users can:
- Discover travel destinations
- Create and share travel itineraries
- Plan trips with friends (group trips)
- Vote on destinations and activities
- Split expenses within groups
- Browse community itineraries
- Like and comment on plans from other travelers

### Key Features
- **User Authentication** - Secure user registration and login
- **Destinations** - Browse and search travel destinations
- **Personal Itineraries** - Create custom travel plans
- **Group Trips** - Plan trips collaboratively with multiple users
- **Community Sharing** - Discover and share itineraries
- **Comments & Ratings** - User feedback system
- **Expense Tracking** - Split costs among group members

---

## 🛠️ Technology Stack

### Current Implementation
| Component | Technology | Version |
|-----------|-----------|---------|
| **Language** | Go | 1.21+ |
| **Web Framework** | Gin | v1.10.0 |
| **Database** | SQLite (Local) / PostgreSQL (Prod) | - |
| **Authentication** | Token-based (simplified JWT for MVP) | - |
| **Template Engine** | Go html/template | Built-in |
| **Logging** | zerolog | v1.x |
| **Package Manager** | Go Modules | Built-in |

### Why These Technologies?
- **Go**: Fast, compiled, excellent for building APIs, great concurrency support
- **Gin**: Lightweight, fast web framework with great middleware support
- **SQLite**: Perfect for development and small deployments
- **zerolog**: Structured logging for production-ready error tracking

---

## 🏗️ Architecture Overview

The project follows a **Layered Architecture** pattern:

```
┌──────────────────────────────────────────────┐
│        HTTP Layer (Routing)                  │
│  - routes.go                                 │
│  - Web routes & API endpoints                │
└──────────────────┬──────────────────────────┘
                   │
┌──────────────────▼──────────────────────────┐
│      Middle Layer (Handlers & Auth)          │
│  - handlers.go                               │
│  - auth_handlers.go                          │
│  - auth_middleware.go                        │
└──────────────────┬──────────────────────────┘
                   │
┌──────────────────▼──────────────────────────┐
│      Business Logic (Service Layer)          │
│  - service.go                                │
│  - auth_service.go                           │
│  - group_service.go                          │
│  - Validation, business rules                │
└──────────────────┬──────────────────────────┘
                   │
┌──────────────────▼──────────────────────────┐
│      Data Layer (Database Access)            │
│  - database.go                               │
│  - group_database.go                         │
│  - SQL queries, data models                  │
└──────────────────┬──────────────────────────┘
                   │
┌──────────────────▼──────────────────────────┐
│      Database (SQLite/PostgreSQL)            │
└──────────────────────────────────────────────┘

┌──────────────────────────────────────────────┐
│    Utility & Infrastructure                  │
│  - config.go        (Configuration)          │
│  - logger.go        (Logging)                │
│  - error.go         (Error Handling)         │
│  - metrics.go       (Monitoring)             │
└──────────────────────────────────────────────┘
```

### Key Principle: **Separation of Concerns**
Each layer has a specific responsibility:
- **HTTP Layer**: Handle HTTP requests/responses
- **Handlers**: Parse input, call services, format output
- **Service**: Business logic, validation, orchestration
- **Database**: Data persistence, SQL execution

---

## 📁 Project Structure

```
itinerary-backend/
│
├── main.go                    # 🎯 Application entry point
│   └── Initializes all components and starts server
│
├── itinerary/                # Core application package
│   │
│   ├── 📊 HTTP & Routing
│   │   ├── routes.go         # API and web routes definition
│   │   ├── handlers.go       # HTTP request handlers
│   │   ├── auth_handlers.go  # Authentication endpoints
│   │   │
│   │   └── 🔐 Middleware
│   │       ├── auth_middleware.go      # Authorization checks
│   │       └── metrics_middleware.go   # Request metrics & monitoring
│   │
│   ├── 🔒 Authentication
│   │   ├── auth.go           # User auth models & methods
│   │   ├── auth_service.go   # Session & token management
│   │   └── auth_*_test.go    # Auth tests
│   │
│   ├── 🚀 Business Logic
│   │   ├── service.go        # Main service logic
│   │   ├── group_service.go  # Group trip logic
│   │   └── group_*_test.go   # Group tests
│   │
│   ├── 📂 Data Access
│   │   ├── database.go       # Database initialization & queries
│   │   ├── group_database.go # Group-specific database operations
│   │   ├── models.go         # Data models (User, Destination, etc.)
│   │   ├── group_models.go   # Group models (GroupTrip, Member, etc.)
│   │   └── *_test.go         # Database tests
│   │
│   ├── 🛠️ Infrastructure
│   │   ├── config.go         # Configuration management
│   │   ├── logger.go         # Structured logging
│   │   ├── error.go          # Error handling & codes
│   │   ├── metrics.go        # Performance metrics
│   │   └── template_helpers.go # HTML template utilities
│   │
│   └── 📝 Unit Tests
│       └── *_test.go         # Test files (can be run with: go test ./itinerary)
│
├── config/
│   └── config.json          # Configuration file (loaded at startup)
│
├── docs/                    # Project documentation
│   ├── DATABASE_SETUP.md    # Database initialization guide
│   ├── QUICK_START.md       # Quick start for developers
│   └── schema.sql           # Database schema
│
├── templates/              # HTML templates for web UI
│   ├── login.html          # Login page
│   ├── dashboard.html      # User dashboard
│   ├── plan-trip.html      # Trip creation page
│   └── ...
│
├── static/                 # Static assets
│   ├── css/                # Stylesheets
│   └── js/                 # JavaScript files
│
├── go.mod                  # Go module definition (dependencies)
├── go.sum                  # Dependency checksums
│
├── .env.example           # Example environment variables
├── itinerary.db           # SQLite database file (created at runtime)
│
└── README.md              # Backend-specific README

📂 Log Directory (created at runtime):
log/
└── itinerary-YYYY-MM-DD.log  # Daily log files
```

### 🔤 File Naming Convention

| Pattern | Meaning | Example |
|---------|---------|---------|
| `*_test.go` | Unit tests | `service_test.go` |
| `*_test.go.skip` | Disabled test | `group_integration_test.go.skip` |
| `*_handlers.go` | HTTP handlers | `auth_handlers.go` |
| `*_service.go` | Business logic | `auth_service.go` |
| `*_database.go` | Data layer | `group_database.go` |
| `*_middleware.go` | Request handlers | `auth_middleware.go` |
| `config_` | Configuration | `config.go`, `config.json` |

---

## 💡 Key Concepts

### 1. **Models (Data Structures)**

Located in `models.go` and `group_models.go`:

```go
// User - Represents a user in the system
type User struct {
    ID        string
    Username  string
    Email     string
    CreatedAt time.Time
}

// Destination - A travel destination
type Destination struct {
    ID          string
    Name        string
    Country     string
    Description string
    Image       string
}

// Itinerary - A complete travel plan
type Itinerary struct {
    ID            string
    UserID        string
    DestinationID string
    Title         string
    Budget        float64
    Items         []ItineraryItem  // Activities, stays, etc.
}

// GroupTrip - A trip for multiple users
type GroupTrip struct {
    ID      string
    OwnerID string
    Title   string
    Members []*GroupMember
}
```

**Key insight:** Models define `binding` tags for validation:
- `binding:"required"` - Field is mandatory
- `binding:"email"` - Must be valid email
- `binding:"gt=0"` - Must be greater than 0
- `binding:"oneof=x y z"` - Must be one of these values

### 2. **Service Layer (Business Logic)**

Located in `service.go` and `*_service.go`:

```go
type Service struct {
    db     *Database  // Access to data layer
    logger *Logger    // For logging
}

// Example: Service method validates and calls database
func (s *Service) GetItinerariesByDestination(destID string) error {
    // 1. Validate input
    if destID == "" {
        return fmt.Errorf("destination ID required")
    }
    
    // 2. Call database
    itineraries, err := s.db.GetItinerariesByDestination(destID)
    
    // 3. Enrich data if needed
    for i := range itineraries {
        items, _ := s.db.GetItineraryItems(itineraries[i].ID)
        itineraries[i].Items = items
    }
    
    // 4. Return to handler
    return itineraries
}
```

**Why?** Services contain business logic separate from HTTP handling.

### 3. **Handlers (HTTP Layer)**

Located in `handlers.go` and `*_handlers.go`:

```go
// Handler receives HTTP request, calls service, returns response
func (h *Handlers) GetDestinations(c *gin.Context) {
    // 1. Parse query parameters
    page := c.Query("page")
    
    // 2. Call service
    destinations, total, err := h.service.GetDestinations(page, pageSize)
    
    // 3. Return JSON response
    c.JSON(200, gin.H{
        "data": destinations,
        "total": total,
    })
}
```

**Key insight:** Handlers should be thin - they parse input and call services.

### 4. **Middleware (Request Interceptors)**

Located in `*_middleware.go`:

Middleware intercepts requests BEFORE they reach handlers:

```
Request → Middleware 1 → Middleware 2 → Handler → Response

// Example: Auth middleware
func (am *AuthMiddleware) RequireAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(401, gin.H{"error": "unauthorized"})
            c.Abort()  // Stop processing
            return
        }
        // Middleware passes auth token to handler via c.Set()
        c.Set("user_id", userID)
        c.Next()  // Continue to next middleware/handler
    }
}
```

**Common middleware in this project:**
- `MetricsHandler()` - Tracks request metrics
- `RequireAuth()` - Checks authentication token
- `RequestLogger()` - Logs all requests

### 5. **Dependency Injection**

Notice the pattern:

```go
type Service struct {
    db     *Database  // Injected
    logger *Logger    // Injected
}

type Handlers struct {
    service *Service  // Injected
    logger  *Logger   // Injected
    metrics *Metrics  // Injected
}

// In main.go
service := NewService(db, logger)        // Pass dependencies
handlers := NewHandlers(service, logger, metrics)  // Pass dependencies
```

**Why?** This makes code testable and loosely coupled.

### 6. **Error Handling**

Located in `error.go`:

```go
// Structured error codes
const (
    ErrInvalidInput    ErrorCode = "INVALID_INPUT"
    ErrNotFound        ErrorCode = "NOT_FOUND"
    ErrUnauthorized    ErrorCode = "UNAUTHORIZED"
)

// API errors return structured responses
apiErr := NewNotFoundError("itinerary not found")
c.JSON(apiErr.StatusCode, apiErr.ToJSON())
// Response: {"code": "NOT_FOUND", "message": "...", "status_code": 404}
```

---

## 🚀 Development Environment Setup

### Prerequisites
- Go 1.21 or higher
- Git
- A code editor (VS Code recommended)
- Postman or curl for testing APIs

### Step 1: Clone & Navigate

```bash
git clone https://github.com/yourusername/iternary.git
cd iternary/itinerary-backend
```

### Step 2: Install Dependencies

```bash
go mod download
```

This reads `go.mod` and downloads all required packages.

### Step 3: Configure Application

```bash
# Copy example config
cp .env.example .env

# Review config/config.json
# Server port: 8080
# DB path: itinerary.db
# Log level: info
```

### Step 4: Run Application

```bash
# Build and run
go run main.go

# Output should show:
# Starting Itinerary Service
# Database connection successful
# Starting server on port 8080
```

### Step 5: Verify It's Running

```bash
# In another terminal
curl http://localhost:8080/api/health

# Expected response:
# {"status":"ok","timestamp":"2026-04-12T10:00:00Z"}
```

---

## 🔍 Understanding Core Components

### Component 1: Router & Routes (`routes.go`)

**Purpose:** Define all API endpoints and web routes

```go
// Example of route definition
router.GET("/api/destinations", handlers.GetDestinations)
//     ↑                                ↑
//   Method                          Handler function
```

**Route Organization:**
- **Web Routes** (return HTML): `/`, `/login`, `/dashboard`
- **API Routes** (return JSON): `/api/destinations`, `/api/user-trips`
- **Auth Routes** (public): `/api/auth/login`, `/api/auth/register`

**Protected vs Public:**
```go
// Public route - anyone can access
router.GET("/api/destinations", handlers.GetDestinations)

// Protected route - requires auth token
router.POST("/api/user-trips", authMiddleware.RequireAuth(), handlers.CreateUserTrip)
```

### Component 2: Database Layer (`database.go`)

**Purpose:** Handle all database operations

**Key Methods:**
```go
// Get destinations with pagination
func (db *Database) GetDestinations(page, pageSize int) ([]Destination, int, error)

// Get itineraries for a destination
func (db *Database) GetItinerariesByDestination(destID string) ([]Itinerary, error)

// Create new itinerary
func (db *Database) CreateItinerary(itinerary *Itinerary) error
```

**Pattern:** Database methods execute SQL and return data or errors.

### Component 3: Service Layer (`service.go`)

**Purpose:** Combine database operations with business logic

```go
type Service struct {
    db     *Database  // For data access
    logger *Logger    // For logging
}

// Example: Service method
func (s *Service) GetItinerariesByDestination(destID string, page, size int) {
    // 1. Validate inputs
    // 2. Call database methods
    // 3. Transform/enrich data
    // 4. Log operations
    // 5. Handle errors
}
```

### Component 4: Handlers (`handlers.go`)

**Purpose:** Handle HTTP requests

**Typical Handler Flow:**
```go
func (h *Handlers) GetDestinations(c *gin.Context) {
    // 1. Parse query parameters
    page := c.Query("page")
    pageSize := c.Query("page_size")
    
    // 2. Validate
    if page < 1 { page = 1 }
    
    // 3. Call service
    destinations, total, err := h.service.GetDestinations(page, pageSize)
    
    // 4. Handle error
    if err != nil {
        h.logger.Error("failed to fetch", "error", err.Error())
        c.JSON(500, gin.H{"error": "internal server error"})
        return
    }
    
    // 5. Return success response
    c.JSON(200, gin.H{
        "data": destinations,
        "total": total,
    })
}
```

### Component 5: Authentication (`auth.go`, `auth_service.go`)

**Purpose:** Manage user sessions and tokens

**How it works:**
```
User Login Form
     ↓
/api/auth/login endpoint
     ↓
Validate username/password
     ↓
Generate secure token
     ↓
Store session in database
     ↓
Return token to client
     ↓
Client includes token in subsequent requests
     ↓
Auth middleware validates token
```

**Token Storage:**
```go
type Session struct {
    ID        string
    UserID    string
    Token     string      // Secure random string
    ExpiresAt time.Time   // When session expires
    CreatedAt time.Time
}
```

### Component 6: Configuration (`config.go`)

**Purpose:** Load and manage application configuration

**Configuration values:**
```json
{
  "server": {
    "port": "8080",
    "timeout": 30,
    "mode": "development"
  },
  "database": {
    "host": "localhost",
    "database": "itinerary.db"
  },
  "logging": {
    "level": "info",
    "format": "json"
  }
}
```

**Usage in code:**
```go
config, err := LoadConfig("config/config.json")
port := config.Server.Port  // "8080"
```

### Component 7: Logging (`logger.go`)

**Purpose:** Structured logging for debugging

**Usage:**
```go
logger.Info("user_logged_in", "user_id", "123", "timestamp", time.Now())
// Output: {"level":"info","msg":"user_logged_in","user_id":"123",...}

logger.Error("database_error", "error", err.Error(), "query", "SELECT...")
// Output: {"level":"error","msg":"database_error","error":"..."}
```

**Log Levels:**
- `debug` - Detailed information for debugging
- `info` - General information about operation
- `warn` - Warning conditions
- `error` - Error conditions

**Log Location:** `log/itinerary-YYYY-MM-DD.log`

---

## 🔄 Common Workflows

### Workflow 1: Adding a New API Endpoint

**Step 1: Define Model** (if needed)
```go
// In models.go
type NewFeature struct {
    ID    string
    Name  string
    // ...
}
```

**Step 2: Add Database Method** (database.go)
```go
func (db *Database) GetNewFeatures() ([]NewFeature, error) {
    // SQL query to fetch data
}
```

**Step 3: Add Service Method** (service.go)
```go
func (s *Service) GetNewFeatures() ([]NewFeature, error) {
    // Validation, business logic
    return s.db.GetNewFeatures()
}
```

**Step 4: Add Handler** (handlers.go)
```go
func (h *Handlers) GetNewFeatures(c *gin.Context) {
    features, err := h.service.GetNewFeatures()
    if err != nil {
        c.JSON(500, gin.H{"error": "internal server error"})
        return
    }
    c.JSON(200, gin.H{"data": features})
}
```

**Step 5: Add Route** (routes.go)
```go
router.GET("/api/features", handlers.GetNewFeatures)
```

**Step 6: Test**
```bash
curl http://localhost:8080/api/features
```

### Workflow 2: Adding Authentication to an Endpoint

```go
// In routes.go
// Before: anyone can access
router.GET("/api/user-trips", handlers.GetUserTrips)

// After: only authenticated users
router.GET("/api/user-trips", 
    authMiddleware.RequireAuth(),  // Add this
    handlers.GetUserTrips)

// In handler, access user ID from context
func (h *Handlers) GetUserTrips(c *gin.Context) {
    userID := c.GetString("user_id")  // From auth middleware
    // Now you can filter trips for this user
}
```

### Workflow 3: Debugging a Failed Request

```bash
# Step 1: Check logs
tail -f log/itinerary-2026-04-12.log

# Step 2: Make request with verbose output
curl -v http://localhost:8080/api/destinations

# Step 3: Check database
sqlite3 itinerary.db
sqlite> SELECT * FROM destinations;

# Step 4: Add debug logs to code
logger.Debug("fetching destinations", "page", page, "size", pageSize)
```

### Workflow 4: Making Database Changes

```bash
# Step 1: Write SQL to add migration
# migrations/001_add_new_column.sql
ALTER TABLE users ADD COLUMN phone VARCHAR(20);

# Step 2: Update model to include new field
type User struct {
    // existing fields...
    Phone string `db:"phone"`
}

# Step 3: Update database methods that use User

# Step 4: Run migration (if using migration tool)
# OR manually execute: sqlite3 itinerary.db < migrations/001_*.sql
```

---

## 🗄️ Database Schema

### Core Tables

**Users Table**
```sql
CREATE TABLE users (
    id TEXT PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Destinations Table**
```sql
CREATE TABLE destinations (
    id TEXT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    country VARCHAR(100) NOT NULL,
    description TEXT,
    image_url TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Itineraries Table**
```sql
CREATE TABLE itineraries (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    destination_id TEXT NOT NULL,
    title VARCHAR(255) NOT NULL,
    budget DECIMAL(10,2),
    duration INT,
    likes INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (destination_id) REFERENCES destinations(id)
);
```

**Itinerary Items Table** (activities, stays, etc.)
```sql
CREATE TABLE itinerary_items (
    id TEXT PRIMARY KEY,
    itinerary_id TEXT NOT NULL,
    day INT NOT NULL,
    type VARCHAR(50),  -- stay, food, activity, transport
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10,2),
    created_at TIMESTAMP,
    FOREIGN KEY (itinerary_id) REFERENCES itineraries(id)
);
```

**Group Trips Table**
```sql
CREATE TABLE group_trips (
    id TEXT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    destination_id TEXT NOT NULL,
    owner_id TEXT NOT NULL,
    budget DECIMAL(10,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (owner_id) REFERENCES users(id),
    FOREIGN KEY (destination_id) REFERENCES destinations(id)
);
```

### Query Examples

```sql
-- Get all destinations
SELECT * FROM destinations;

-- Get itineraries for a destination ordered by likes
SELECT * FROM itineraries 
WHERE destination_id = 'dest_123'
ORDER BY likes DESC;

-- Get items for an itinerary sorted by day
SELECT * FROM itinerary_items 
WHERE itinerary_id = 'iter_456'
ORDER BY day, type;

-- Get group trip members
SELECT u.*, gm.role 
FROM group_members gm
JOIN users u ON gm.user_id = u.id
WHERE gm.group_trip_id = 'group_789';
```

---

## 🔌 API Endpoints

### Authentication Endpoints

| Method | Endpoint | Auth | Purpose |
|--------|----------|------|---------|
| POST | `/api/auth/register` | No | Register new user |
| POST | `/api/auth/login` | No | Login user, get token |
| POST | `/api/auth/logout` | Yes | Logout user |

### Destination Endpoints

| Method | Endpoint | Auth | Purpose |
|--------|----------|------|---------|
| GET | `/api/destinations` | No | List all destinations |
| GET | `/api/destinations/:id` | No | Get destination details |

### Itinerary Endpoints

| Method | Endpoint | Auth | Purpose |
|--------|----------|------|---------|
| GET | `/api/destinations/:id/itineraries` | No | List itineraries for destination |
| GET | `/api/itineraries/:id` | No | Get itinerary details |
| POST | `/api/itineraries` | No | Create new itinerary |
| POST | `/api/itineraries/:id/like` | No | Like an itinerary |
| POST | `/api/itineraries/:id/comments` | No | Add comment |

### User Trip Endpoints (Personal Plans)

| Method | Endpoint | Auth | Purpose |
|--------|----------|------|---------|
| POST | `/api/user-trips` | Yes | Create personal trip |
| GET | `/api/user-trips` | Yes | List my trips |
| GET | `/api/user-trips/:id` | Yes | Get trip details |
| PUT | `/api/user-trips/:id` | Yes | Update trip |
| DELETE | `/api/user-trips/:id` | Yes | Delete trip |

### Group Trip Endpoints

| Method | Endpoint | Auth | Purpose |
|--------|----------|------|---------|
| POST | `/api/group-trips` | Yes | Create group trip |
| GET | `/api/group-trips/:id` | Yes | Get group trip |
| PUT | `/api/group-trips/:id` | Yes | Update group trip |
| POST | `/api/group-trips/:id/members` | Yes | Add member |
| POST | `/api/group-trips/:id/expenses` | Yes | Add expense |

### Example Requests

**Get Destinations**
```bash
curl "http://localhost:8080/api/destinations?page=1&page_size=10"
```

**Create Itinerary**
```bash
curl -X POST http://localhost:8080/api/itineraries \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user_123",
    "destination_id": "dest_456",
    "title": "Summer in Paris",
    "budget": 5000,
    "duration": 7
  }'
```

**Protected Endpoint (Requires Token)**
```bash
curl -X POST http://localhost:8080/api/user-trips \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "My Paris Trip",
    "destination_id": "dest_456"
  }'
```

---

## ✅ Testing Guide

### Running Unit Tests

```bash
# Run all tests
go test ./itinerary

# Run tests with coverage
go test -cover ./itinerary

# Output shows:
# coverage: 85.2% of statements

# Run specific test file
go test ./itinerary -run TestAuthService

# Run with verbose output
go test -v ./itinerary
```

### Test File Locations

Tests are in the same directory as code they test:

```
auth_service.go          ─→  auth_service_test.go
group_service.go         ─→  group_service_test.go
database.go              ─→  database_test.go (if exists)
```

### Understanding Test Structure

```go
// tests/auth_service_test.go
func TestGenerateToken(t *testing.T) {
    // 1. Arrange - Setup test data
    authService := NewAuthService(nil, nil)
    
    // 2. Act - Execute the function
    token, err := authService.GenerateToken()
    
    // 3. Assert - Check results
    if err != nil {
        t.Errorf("expected no error, got %v", err)
    }
    if token == "" {
        t.Errorf("expected non-empty token")
    }
}
```

### Integration Testing

```bash
# Test endpoints locally (with server running)
curl http://localhost:8080/api/health
curl http://localhost:8080/api/destinations
```

### Test Coverage Goals

- **Current:** 85%+
- **Target:** 90%+
- **Critical paths:** 100% (auth, payments, group operations)

---

## 🐛 Debugging Tips

### Tip 1: Enable Debug Logging

```bash
# In config/config.json
{
  "logging": {
    "level": "debug"  # Changed from "info"
  }
}
```

### Tip 2: Use Print Debugging

```go
// Add temporary debug prints
fmt.Println("DEBUG: Destination ID:", destID)
fmt.Println("DEBUG: User ID:", userID)
```

### Tip 3: Check Logs

```bash
# View latest logs
tail -50 log/itinerary-2026-04-12.log

# Search logs for errors
grep ERROR log/itinerary-2026-04-12.log

# Follow log as it's written (like tail -f)
tail -f log/itinerary-2026-04-12.log
```

### Tip 4: Test with curl

```bash
# Make requests from command line
curl -v http://localhost:8080/api/destinations
# -v flag shows request and response headers

# Test POST request
curl -X POST http://localhost:8080/api/itineraries \
  -H "Content-Type: application/json" \
  -d '{"user_id": "123", ...}'
```

### Tip 5: Check Database Directly

```bash
# Open SQLite database
sqlite3 itinerary.db

# Inside sqlite3 shell:
sqlite> .tables              # List all tables
sqlite> SELECT * FROM users; # View data
sqlite> .schema users        # View table structure
sqlite> .exit                # Exit
```

### Tip 6: Add Breakpoints (if using VS Code with Go)

1. Click left of line number to add breakpoint (red dot appears)
2. Run: `go run main.go`
3. Make request that hits breakpoint
4. Debug console will pause and show variables

### Tip 7: Common Error Messages & Solutions

| Error | Cause | Solution |
|-------|-------|----------|
| `database is locked` | SQLite conflict | Ensure only one process accessing DB |
| `connection refused` | Server not running | Run `go run main.go` |
| `404 - Not Found` | Wrong endpoint | Check route in `routes.go` |
| `unauthorized` | Missing auth token | Add `Authorization: Bearer TOKEN` header |
| `json: cannot unmarshal` | Wrong request format | Check JSON structure matches model |

---

## 🔧 Troubleshooting

### Issue: Server Won't Start

```bash
# Check if port is in use
# On Windows:
netstat -ano | findstr :8080

# On Mac/Linux:
lsof -i :8080

# Solution: Kill process or use different port in config
```

### Issue: Database Not Found

```bash
# Error: database file not found
# Solution: Ensure you're in correct directory
cd itinerary-backend
# Database will be created automatically on first run
```

### Issue: Imports Not Found

```bash
# Error: cannot find module...
# Solution: Download dependencies
go mod download
```

### Issue: Route Returns 404

```bash
# Check if route is defined
# 1. Open routes.go
# 2. Search for endpoint name
# 3. Verify method (GET/POST) matches request
# 4. Test with: curl http://localhost:8080/api/destinations
```

### Issue: Auth Token Failing

```bash
# Check token format
# Token should be in header: Authorization: Bearer <token>
# Or query param: ?token=<token>

# Test token validation
curl -H "Authorization: Bearer YOUR_TOKEN" http://localhost:8080/api/user-trips
```

### Issue: Can't Connect to Database

```bash
# Check database permissions
ls -la itinerary.db

# Ensure database file is writable
# On Windows: Right-click > Properties > Security
# On Mac/Linux: chmod 644 itinerary.db
```

---

## 📚 Next Steps for New Developers

### 1. **First Day: Understand Structure**
   - [ ] Read this guide completely
   - [ ] Explore the directory structure
   - [ ] Run the application locally
   - [ ] Make requests using curl
   - [ ] View database using sqlite3

### 2. **Second Day: Add Simple Feature**
   - [ ] Add a new simple API endpoint
   - [ ] Follow Workflow 1: "Adding a New API Endpoint"
   - [ ] Test the endpoint with curl
   - [ ] View data in database

### 3. **Third Day: Understand Data Flow**
   - [ ] Trace a request from handler → service → database
   - [ ] Add debug logs along the way
   - [ ] Understand how data is transformed at each layer

### 4. **Week 1: Contribute to Group Features**
   - [ ] Study group functionality (group_service.go, group_database.go)
   - [ ] Fix a small bug or add a small feature
   - [ ] Write tests for your changes
   - [ ] Submit for code review

### 5. **Ongoing: Best Practices**
   - [ ] Always write tests for new code
   - [ ] Keep functions small and focused
   - [ ] Add logging for important operations
   - [ ] Follow existing code patterns
   - [ ] Write clear commit messages

---

## 🤝 Contributing Guidelines

### Before Starting Work

1. Create a new branch: `git checkout -b feature/your-feature-name`
2. Reference the file structure and architecture
3. Follow existing code style and patterns

### After Making Changes

1. Run tests: `go test ./itinerary`
2. Check for compile errors: `go build`
3. Review your code for clarity
4. Commit with clear message: `git commit -m "Add X feature"`

### Code Quality Checklist

- [ ] Tests written for new code
- [ ] No unused imports
- [ ] Functions have clear purposes
- [ ] Error handling implemented
- [ ] Logging added for important operations
- [ ] Documentation updated if needed

---

## 📖 Additional Resources

### Files to Read First

1. [README.md](README.md) - Project overview
2. [API_REFERENCE.md](../API_REFERENCE.md) - API documentation
3. [PROJECT_REQUIREMENTS.md](../PROJECT_REQUIREMENTS.md) - Feature specs

### Go Language Resources

- [Go Official Documentation](https://golang.org/doc/)
- [Effective Go](https://golang.org/doc/effective_go)
- [Go by Example](https://gobyexample.com/)

### Framework Documentation

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [zerolog Logger](https://github.com/rs/zerolog)

### Local Documentation

- [Database Setup](./docs/DATABASE_SETUP.md)
- [Quick Start](./docs/QUICK_START.md)
- [Templates Guide](./docs/TEMPLATES_GUIDE.md)

---

## 🎓 Understanding Go Syntax (Basics)

If you're new to Go, here are key syntax elements used in this project:

### Packages
```go
package main  // Every file belongs to a package
```

### Imports
```go
import (
    "database/sql"
    "github.com/gin-gonic/gin"
)
```

### Structs (Like Classes)
```go
type User struct {
    ID       string
    Email    string
    CreatedAt time.Time
}
```

### Methods (Functions on Structs)
```go
// Method: function that belongs to a type
func (h *Handlers) GetDestinations(c *gin.Context) {
    //  ↑      ↑
    // receiver method body
}
```

### Error Handling
```go
result, err := someFunction()
if err != nil {
    log.Fatal(err)
}
```

### Interfaces
```go
// Contract: any type with these methods satisfies this interface
type Reader interface {
    Read(p []byte) (n int, err error)
}
```

### Goroutines
```go
// Run function concurrently
go fetchData()
go processRequest()
```

---

## ❓ FAQ

**Q: Where do I add new features?**  
A: Follow Workflow 1 - add model, database method, service method, handler, route.

**Q: How do I debug issues?**  
A: Check logs first (`log/` directory), then add debug logging, then check database directly.

**Q: How do I run tests?**  
A: `go test ./itinerary` or `go test -v -cover ./itinerary`

**Q: What if I need to change the database schema?**  
A: Add SQL migration, update models, update database methods, test thoroughly.

**Q: How do I add authentication to an endpoint?**  
A: Add `authMiddleware.RequireAuth()` to route in `routes.go`.

**Q: Where are logs stored?**  
A: In `log/itinerary-YYYY-MM-DD.log` (created at runtime).

**Q: Can I use a different database?**  
A: Currently built for SQLite. PostgreSQL support is planned but requires migration.

---

## 📞 Getting Help

- Check the troubleshooting section above
- Search logs for error messages
- Refer to existing similar code
- Ask team members or create an issue on GitHub

---

**Last Updated:** April 12, 2026  
**Created for:** Beginner Go Developers  
**Version:** 1.0
