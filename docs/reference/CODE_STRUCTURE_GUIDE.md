# Backend Code Structure Guide

**Date:** April 12, 2026  
**Status:** Reorganized using modular Go project architecture

---

## 📁 Code Organization Overview

The backend code is now organized using a **modular, domain-driven** structure following Go best practices:

```
itinerary-backend/
├── itinerary/                    # Main application package
│   ├── auth/                     # 🔐 Authentication module
│   ├── groups/                   # 👥 Groups management module
│   ├── models/                   # 📊 Core data models
│   ├── handlers/                 # 🔗 HTTP request handlers
│   ├── middleware/               # ⚙️ Request middleware (metrics, etc)
│   ├── utils/                    # 🛠️ Utility functions
│   ├── config/                   # ⚙️ Configuration management
│   ├── routes/                   # 🛣️ API routing
│   ├── service/                  # 💼 Business logic layer
│   ├── database/                 # 🗄️ Database optimization
│   ├── cache/                    # 💾 Caching layer
│   ├── integration_tests/        # 🧪 Integration tests
│   ├── static/                   # 🎨 Frontend assets
│   └── templates/                # 📄 HTML templates
│
├── main.go                       # Application entry point
├── go.mod & go.sum
└── config/
    └── config.json
```

---

## 🏗️ Module Breakdown

### 🔐 **auth/** - Authentication
**Responsibility:** User authentication and authorization  
**Files:**
- `auth.go` - Core authentication logic
- `handlers.go` - Auth endpoints (login, register, etc)
- `middleware.go` - Auth middleware for protected routes
- `service.go` - Auth business logic
- `service_test.go` - Auth tests

**Key Exports:** `AuthService`, `AuthHandler`, `AuthMiddleware`

---

### 👥 **groups/** - Group Management
**Responsibility:** Group operations and management  
**Files:**
- `models.go` - Group data structures
- `handlers.go` - Group API endpoints
- `service.go` - Group business logic
- `database.go` - Group-specific database operations
- `routes.go` - Group routing
- `integration_test.go` - Group integration tests

**Key Exports:** `GroupService`, `GroupHandler`

---

### 📊 **models/** - Core Data Models
**Responsibility:** Primary data structures  
**Files:**
- `models.go` - Itinerary, User, Post, Item models
- `models_test.go` - Model tests

**Key Exports:** `Itinerary`, `User`, `Post`, `Item`, etc.

---

### 🔗 **handlers/** - HTTP Handlers
**Responsibility:** Main HTTP request handling  
**Files:**
- `handlers.go` - Main HTTP handlers for core features

**Key Exports:** `Handler`, HTTP handler functions

---

### ⚙️ **middleware/** - Request Middleware
**Responsibility:** Cross-cutting concerns (metrics, logging, etc)  
**Files:**
- `metrics.go` - Metrics collection
- `metrics_middleware.go` - Metrics middleware
- Additional middleware as needed

**Key Exports:** `MetricsMiddleware`, metrics functions

---

### 🛠️ **utils/** - Utility Functions
**Responsibility:** Common utility functions  
**Files:**
- `logger.go` - Logging utilities
- `error.go` - Error handling
- `template_helpers.go` - HTML template helpers
- Corresponding test files

**Key Exports:** `Logger`, `Error`, `TemplateHelper`

---

### ⚙️ **config/** - Configuration
**Responsibility:** Application configuration management  
**Files:**
- `config.go` - Configuration loading and management
- `config_test.go` - Config tests

**Key Exports:** `Config`, configuration functions

---

### 🛣️ **routes/** - API Routing
**Responsibility:** API route definitions  
**Files:**
- `routes.go` - Main router setup, route definitions

**Key Exports:** `SetupRoutes()`, route handler mappings

---

### 💼 **service/** - Business Logic Layer
**Responsibility:** Core business logic  
**Files:**
- `service.go` - Main service implementation
- `database.go` - General database operations

**Key Exports:** `Service`, database utilities

---

### 🗄️ **database/** - Database Optimization
**Responsibility:** Query optimization and database utilities  
**Files:**
- `pool.go` - Connection pooling
- `indexes.go` - Index management
- `optimization_module.go` - Query optimization
- `query_optimizer.go` - Optimization logic
- `query_profiler.go` - Performance profiling
- Additional optimization utilities

**Key Exports:** `QueryOptimizer`, `QueryProfiler`, database pool

---

### 💾 **cache/** - Caching Layer
**Responsibility:** Caching strategy and implementation  
**Files:**
- `factory.go` - Cache factory
- `memory_cache.go` - In-memory cache implementation
- `redis/` - Redis cache implementation
- Examples and documentation

**Key Exports:** `Cache`, `CacheFactory`, cache implementations

---

### 🧪 **integration_tests/** - Integration Tests
**Responsibility:** End-to-end integration testing  
**Files:**
- `integration_test.go` - Basic integration tests
- `day6_integration_test.go` - Additional integration scenarios
- `performance_test.go` - Performance benchmarks
- `security_test.go` - Security validation tests

**Purpose:** Test interactions between modules

---

---

## 📋 Layer Architecture

**Dependency Flow:**
```
HTTP Requests
    ↓
routes/ (routing)
    ↓
handlers/ (request handling)
    ↓
service/ (business logic)
    ↓
database/ + cache/ (persistence)
```

**Cross-Cutting Concerns:**
```
middleware/  (metrics, logging, etc)
utils/       (helpers, errors, logging)
config/      (configuration)
auth/        (authentication)
```

---

## 🔄 Import Organization

**Import Hierarchy (avoid circular dependencies):**

```
models/ 
    ↓
config/, utils/, logger/
    ↓
cache/, database/
    ↓
service/, auth/, groups/
    ↓
handlers/
    ↓
middleware/
    ↓
routes/
    ↓
main.go
```

---

## 🔑 Key Principles

### 1. **Single Responsibility**
Each module has one primary responsibility:
- `auth/` handles authentication
- `groups/` handles groups
- `handlers/` handles HTTP requests
- `service/` handles business logic

### 2. **Clear Boundaries**
Modules communicate through well-defined interfaces:
- `auth.Service` provides auth operations
- `groups.Service` provides group operations
- Each module exports public functions/types

### 3. **Testability**
- Unit tests co-located with modules (`*_test.go`)
- Integration tests in `integration_tests/`
- Mock-friendly interfaces

### 4. **Scalability**
- Easy to add new modules (e.g., `payments/`, `notifications/`)
- Shared utilities in `utils/`
- Centralized database layer

---

## 📦 Package Names Convention

Each module follows the pattern:
```go
package auth      // auth/auth.go
package groups    // groups/models.go
package models    // models/models.go
package service   // service/service.go
// etc.
```

---

## 🚀 Adding New Features

**When adding a new feature:**

1. **If it's a new domain** (e.g., payments):
   ```
   Create: itinerary/payments/
   With:   models.go, service.go, handlers.go, routes.go, database.go
   ```

2. **If it's a utility**:
   ```
   Add to: itinerary/utils/
   ```

3. **If it's crosscutting** (e.g., logging):
   ```
   Add to: itinerary/middleware/ or utils/
   ```

4. **Update imports** in parent files accordingly

---

## 📊 File Organization Summary

| Module | Purpose | Typical Files |
|--------|---------|---------------|
| **auth** | User authentication | auth.go, handlers.go, service.go, middleware.go |
| **groups** | Group management | models.go, handlers.go, service.go, database.go, routes.go |
| **models** | Core data models | models.go, models_test.go |
| **handlers** | HTTP handling | handlers.go |
| **middleware** | Cross-cutting | metrics.go, logging.go, etc. |
| **utils** | Utilities | error.go, logger.go, helpers.go |
| **config** | Configuration | config.go, config_test.go |
| **routes** | Routing | routes.go |
| **service** | Business logic | service.go, database.go |
| **database** | DB optimization | pool.go, optimization.go, profiler.go |
| **cache** | Caching | factory.go, memory_cache.go, redis/ |
| **integration_tests** | E2E tests | *.go test files |

---

## 🔍 Finding Code

**Looking for...**
- Authentication logic → `auth/`
- Group features → `groups/`
- Data models → `models/models.go`
- HTTP endpoints → `handlers/` or module-specific `handlers.go`
- Business logic → `service/service.go` or module `service.go`
- Database queries → `database/` or module `database.go`
- Utilities → `utils/`
- Configuration → `config/`
- Tests → Look for `*_test.go` in each module or `integration_tests/`

---

## 🧪 Testing Strategy

**Unit Tests:** Co-located with modules
```
models/models_test.go
auth/service_test.go
groups/models_test.go
```

**Integration Tests:** Centralized in one place
```
integration_tests/integration_test.go
integration_tests/performance_test.go
integration_tests/security_test.go
```

**Running Tests:**
```bash
# Unit tests only
go test ./itinerary/auth/
go test ./itinerary/groups/

# All tests
go test ./...

# Integration tests
go test ./itinerary/integration_tests/
```

---

## 📈 Structure Benefits

✅ **Clear Organization** - Know where code lives  
✅ **Easy Maintenance** - Find related code together  
✅ **Scalability** - Easy to add new modules  
✅ **Modularity** - Replace implementations without affecting others  
✅ **Testability** - Isolated modules are easier to test  
✅ **Team Collaboration** - Different teams can own different modules  

---

## 🔄 Migration Notes

**What Changed:**
- ✅ Grouped related files into modules
- ✅ Renamed files for clarity (e.g., `auth_handlers.go` → `auth/handlers.go`)
- ✅ Separated integration tests
- ✅ Organized by domain/concern

**What Stayed:**
- ✅ Same business logic
- ✅ Same functionality
- ✅ Same tests (just reorganized)

**Import Updates Needed:**
Files now in modules, so imports need updating:
```go
// Old
import (
    "itinerary-backend/itinerary"
)

// New
import (
    "itinerary-backend/itinerary/auth"
    "itinerary-backend/itinerary/groups"
)
```

---

## 🛠️ Next Steps

1. ✅ Code organized into modules
2. ⏳ Update imports in `main.go` and `routes/`
3. ⏳ Update imports in tests
4. ⏳ Run tests to verify structure
5. ⏳ Update documentation links

---

**Architecture Follows:** Go best practices, domain-driven design  
**Status:** Ready for development  
**Last Updated:** April 12, 2026
