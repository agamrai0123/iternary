# Itinerary Project

A comprehensive travel itinerary planning platform built with Go, allowing users to discover, share, and customize travel plans with transparent pricing.

## 🏗️ Project Structure

```
iternary/
├── .gitignore                 # Git ignore rules
├── README.md                  # Main project README (this file)
├── idea.txt                   # Original project idea & feasibility analysis
│
├── itinerary-backend/         # Go backend application
│   ├── main.go               # Application entry point
│   ├── go.mod                # Go module definition
│   ├── go.sum                # Go dependencies lock file
│   ├── .env.example          # Example environment variables
│   │
│   ├── itinerary/            # Core business logic
│   │   ├── auth_*.go         # Authentication & authorization
│   │   ├── group_*.go        # Group management features
│   │   ├── models.go         # Data models
│   │   ├── service.go        # Business logic services
│   │   ├── handlers.go       # HTTP request handlers
│   │   ├── routes.go         # API routes definition
│   │   ├── database.go       # Database operations
│   │   ├── config.go         # Configuration management
│   │   ├── logger.go         # Logging utilities
│   │   ├── metrics.go        # Performance metrics
│   │   ├── error.go          # Error handling
│   │   └── *_test.go         # Unit tests
│   │
│   ├── config/               # Configuration files
│   │   └── config.json       # Application configuration
│   │
│   ├── docs/                 # Documentation & schemas
│   │   ├── schema.sql        # Database schema
│   │   ├── DATABASE_SETUP.md # Database setup guide
│   │   ├── QUICK_START.md    # Quick start guide
│   │   ├── TEMPLATES_GUIDE.md# Template usage guide
│   │   └── *_SCHEMA.sql      # Feature-specific schemas
│   │
│   ├── static/               # Static assets
│   │   ├── css/              # Stylesheets
│   │   └── js/               # Client-side JavaScript
│   │
│   ├── templates/            # HTML templates
│   │   ├── index.html        # Homepage
│   │   ├── login.html        # Login page
│   │   ├── dashboard.html    # User dashboard
│   │   ├── plan-trip.html    # Trip planning page
│   │   ├── itinerary-detail.html
│   │   └── ...
│   │
│   └── README.md             # Backend-specific documentation
│
├── docs/                     # Project documentation
│   ├── PHASE_A_WEEK_2_DAY_*.md      # Development phase notes
│   ├── PHASE_A_WEEK_2_DOCUMENTATION_INDEX.md
│   └── ...
│
├── archives/                 # Archived documentation (historical)
│   └── PHASE_A_*.md
│
├── API_REFERENCE.md          # API documentation
├── PROJECT_REQUIREMENTS.md   # Project requirements
├── DOCUMENTATION_INDEX.md    # Documentation guide
├── GETTING_STARTED.md        # Getting started guide
└── ...README files
```

## 🚀 Getting Started

1. **Setup Backend**
   ```bash
   cd itinerary-backend
   go mod download
   go build
   ./itinerary-backend
   ```

2. **Database Setup**
   - See [Database Setup](itinerary-backend/docs/DATABASE_SETUP.md)

3. **Configuration**
   - Copy `.env.example` to `.env`
   - Update configuration in `config/config.json`

## 📚 Documentation

For detailed information, refer to:
- [Documentation Index](DOCUMENTATION_INDEX.md) - Navigation guide for all docs
- [Getting Started](GETTING_STARTED.md) - Quick setup guide
- [API Reference](API_REFERENCE.md) - API endpoints documentation
- [Project Requirements](PROJECT_REQUIREMENTS.md) - Feature specifications
- [Backend README](itinerary-backend/README.md) - Backend-specific details

## 🏷️ Key Features

- **User Authentication** - Secure login and registration
- **Itinerary Browsing** - Discover shared travel plans
- **Social Features** - Like, comment, and rate itineraries
- **Copy & Customize** - Duplicate itineraries and modify them
- **Group Management** - Create and manage travel groups
- **Multi-Currency** - Support for multiple currencies
- **Transparent Pricing** - See all costs upfront

## 🔧 Development

**Technology Stack**
- Backend: Go + Gorilla/mux
- Frontend: HTML5/CSS3/JavaScript  
- Database: SQLite3
- Testing: Go testing package

**Code Organization**
- `itinerary/` - Core business logic modules
- `*_test.go` - Unit tests (run with `go test ./...`)
- `handlers.go` - HTTP request handlers
- `service.go` - Business logic layer
- `database.go` - Data persistence layer

## 📋 Project Files

**Active Documentation** (in root)
- `PROJECT_REQUIREMENTS.md` - What to build
- `API_REFERENCE.md` - API specifications
- `GETTING_STARTED.md` - Setup & running the project

**Archived Documentation** (in `archives/`)
- Historical phase documentation
- Development milestones
- Verification reports

## 🔒 .gitignore

The project includes a comprehensive `.gitignore` that excludes:
- Build artifacts (`.exe`, `.o`, etc.)
- Databases (`*.db`, `*.sqlite`)
- Environment files (`.env`)
- Logs (`*.log`)
- IDE settings
- OS files

## 📞 Next Steps

1. Review [Project Requirements](PROJECT_REQUIREMENTS.md)
2. Read [Getting Started](GETTING_STARTED.md)
3. Check [Backend README](itinerary-backend/README.md) for backend details
4. Run tests: `cd itinerary-backend && go test ./...`

---

**Last Updated:** March 24, 2026
