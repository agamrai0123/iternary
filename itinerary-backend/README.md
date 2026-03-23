# Travel Itinerary Platform - Backend

A **Go-based backend** with **HTML templates** for a crowd-sourced travel itinerary platform. No React or frontend framework required—everything is built with Go!

## 🎯 Project Highlights

- **100% Go-based** - Backend, routing, and frontend rendering all in Go
- **HTML Templates** - Using Go's `text/template` for server-side rendering
- **Clean Architecture** - Following the auth-service pattern with separation of concerns
- **Modern UI** - Responsive design with vanilla CSS and minimal JavaScript
- **RESTful API** - Both web routes and JSON API endpoints

## 📋 Project Overview

Users can:
- Browse destinations
- View community-created itineraries ranked by popularity
- See detailed itinerary plans with costs for accommodations, food, activities
- Like and comment on itineraries
- Create and customize their own itineraries
- Calculate total trip expenses upfront
- Search itineraries by criteria

## 🏗️ Architecture

```
itinerary-backend/
├── main.go                 # Entry point
├── config/
│   └── config.json        # Configuration file
├── itinerary/
│   ├── config.go          # Config management
│   ├── models.go          # Data models
│   ├── database.go        # Database operations
│   ├── handlers.go        # HTTP handlers (web + API)
│   ├── service.go         # Business logic
│   ├── logger.go          # Logging
│   ├── routes.go          # Route definitions
│   └── template_helpers.go # Template functions
├── templates/             # HTML templates (Go templates)
│   ├── index.html         # Home page - browse destinations
│   ├── destination-detail.html  # List itineraries
│   ├── itinerary-detail.html    # Full itinerary view
│   ├── create-itinerary.html    # Create itinerary form
│   └── search.html        # Search page
├── static/                # Static files
│   ├── css/
│   │   └── style.css      # Responsive styling
│   └── js/
│       └── app.js         # Frontend JS (HTMX-ready, minimal)
├── docs/                  # Documentation
│   └── schema.sql        # Database schema
├── .env.example          # Environment template
├── go.mod                # Go module
└── README.md             # This file
```

## 🚀 Quick Start

### Prerequisites
- Go 1.21 or higher
- Oracle Database 12c+ (XE works great)
- Oracle running on `localhost:1521` with service `XE`

### Installation

1. Clone and setup:
```bash
cd itinerary-backend
go mod download
```

2. Set environment variables:
```bash
export DB_PASSWORD=your_oracle_password
export DB_HOST=localhost
```

3. Initialize database:
```bash
go run init_db.go init
```

This will:
- Create all tables (users, destinations, itineraries, items, comments)
- Insert test data (3 sample users, 3 destinations, 4 itineraries)
- Verify setup complete

4. Run the server:
```bash
go run main.go
```

Server will start on `http://localhost:8080`

**For detailed database setup**, see [docs/DATABASE_SETUP.md](docs/DATABASE_SETUP.md)

## 📡 API Endpoints

### Health Check
- `GET /api/health` - Service health status

### Destinations
- `GET /api/destinations?page=1&page_size=10` - List all destinations

### Itineraries
- `GET /api/destinations/:destinationId/itineraries` - Get itineraries for a destination
- `GET /api/itineraries/:itineraryId` - Get detailed itinerary
- `POST /api/itineraries` - Create new itinerary

### User Interactions
- `POST /api/itineraries/:itineraryId/like` - Like an itinerary
- `POST /api/itineraries/:itineraryId/comments` - Add a comment

## 📊 Database Schema

### Core Tables
- **destinations** - Travel destinations (cities, countries)
- **itineraries** - User-created travel plans
- **itinerary_items** - Individual items in itineraries (stays, food, activities)
- **comments** - User comments on itineraries
- **users** - User profiles
- **user_plans** - Saved/copied itineraries by users

## 🔧 Configuration

Edit `config/config.json` or override with environment variables:

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

**Environment Variable Overrides:**
- `DB_PASSWORD` - Oracle system password
- `DB_HOST` - Database hostname

See `itinerary/config.go` for how configuration is loaded and merged.

## ✨ Key Features (Phase 1)

- [x] API structure setup
- [x] HTML template rendering
- [x] Responsive CSS styling
- [x] Destination browsing
- [x] Itinerary listing and viewing
- [x] Route definitions (web + API)
- [ ] Database migrations/schema
- [ ] User authentication
- [ ] Database persistence
- [ ] Search functionality
- [ ] Price aggregation
- [ ] Ranking algorithm

## 🎨 Frontend Architecture

This project uses **Go's `text/template`** for server-side rendering instead of React:

### Advantages:
- ✅ No build step or Node.js required
- ✅ Single language—no context switching between Go and JavaScript
- ✅ Faster initial page loads (server-side rendering)
- ✅ SEO-friendly (content in HTML)
- ✅ Smaller deployment size
- ✅ simpler development workflow

### Pages:
- **Home** (`/`) - Browse all destinations
- **Destination Detail** (`/destination/:id`) - View itineraries for selected destination
- **Itinerary Detail** (`/itinerary/:id`) - Complete plan with daily breakdown and costs
- **Create Itinerary** (`/create`) - Form to create new itinerary
- **Search** (`/search`) - Search and filter itineraries

### Template Functions:
- `{{add .x 1}}` - Add numbers
- `{{sub .x 1}}` - Subtract
- `{{divide .total .days}}` - Divide floats
- `{{toUpper .text}}` - Uppercase text
- `{{typeIcon .type}}` - Get emoji icon for type
- `{{formatPrice .amount}}` - Format as currency

## 🤖 Planned AI Features

1. **AI Itinerary Generator** - Generate plans based on budget and preferences
2. **Price Freshness Agent** - Auto-update stale prices
3. **Smart Remix** - Compress/expand existing itineraries
4. **Personalized Recommendations** - Based on user filters
5. **Comment Summarization** - Extract pros/cons
6. **Auto-categorization** - Tag itineraries by type/season

## 📝 Environment Variables

```bash
DB_PASSWORD=your_password
DB_HOST=localhost
```

## 🧪 Testing

```bash
go test ./...
```

## � Next Steps

1. **Database Integration** - Connect to MySQL and implement data persistence
2. **User Authentication** - JWT-based auth for users
3. **Form Handling** - Process create/edit forms with validation
4. **Search Implementation** - Full-text search in itineraries
5. **Enhanced Frontend** - Add HTMX for dynamic interactions without page reloads
6. **API Endpoints** - Complete REST API for mobile apps

## 💡 Frontend Enhancement Tips

The current frontend is **HTML + CSS + Vanilla JS**. To add dynamic interactions without React, consider:

### Option 1: Add HTMX
```html
<button hx-post="/api/itineraries/:id/like" 
        hx-swap="innerHTML"
        hx-target="#like-count">
  Like
</button>
```

### Option 2: Keep It Simple
The current setup with vanilla JS is perfect for simple AJAX calls and form submissions.

### Option 3: Add Alpine.js
Lightweight JavaScript framework for interactive components:
```html
<div x-data="{ liked: false }">
  <button @click="liked = !liked">Like</button>
</div>
```

## �📚 Learning Resources

- [Gin Web Framework](https://gin-gonic.com/)
- [Go Database/SQL](https://golang.org/pkg/database/sql/)
- [Project Idea Analysis](../idea.txt)

## 🚀 Next Steps

1. **Database Schema** - Create migrations for all tables
2. **Authentication** - Implement JWT-based auth
3. **Validation** - Add comprehensive input validation
4. **Testing** - Write unit and integration tests
5. **Frontend** - Build React/Next.js UI

## 📄 License

MIT License

## 👤 Author

Your Name

## 🤝 Contributing

1. Create a feature branch
2. Make your changes
3. Submit a pull request

---

**Status:** 🎉 Go Templates Frontend Ready - Database integration next
