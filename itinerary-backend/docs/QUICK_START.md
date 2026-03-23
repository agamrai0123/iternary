# Quick Start - Go Templates Project

## Prerequisites

- Go 1.21 or higher
- MySQL 8.0+ (for production)
- A terminal/command prompt

## Installation & Setup

### 1. Navigate to Project
```bash
cd d:\Learn\iternary\itinerary-backend
```

### 2. Download Dependencies
```bash
go mod download
```

### 3. Build the Project
```bash
go build -o itinerary-backend.exe
```

### 4. Run Directly
```bash
go run main.go
```

## Accessing the Application

Once running, open your browser to:

```
http://localhost:8080
```

### Available Pages

- **Home**: `http://localhost:8080/` - Browse destinations
- **Destination**: `http://localhost:8080/destination/dest-001` - View itineraries
- **Itinerary**: `http://localhost:8080/itinerary/itin-001` - Full plan
- **Create**: `http://localhost:8080/create` - Create new itinerary
- **Search**: `http://localhost:8080/search` - Search itineraries
- **Health**: `http://localhost:8080/api/health` - API health check

## Project Structure at a Glance

```
├── main.go                    Entry point - starts server
├── go.mod                     Dependencies
│
├── config/
│   └── config.json           Configuration (port, database, etc.)
│
├── itinerary/
│   ├── config.go            Config loading
│   ├── models.go            Data structures (User, Itinerary, etc.)
│   ├── database.go          Database operations (currently mocked)
│   ├── handlers.go          HTTP handlers (renders templates or returns JSON)
│   ├── service.go           Business logic
│   ├── logger.go            Logging
│   ├── routes.go            Route registration + template setup
│   └── template_helpers.go  Template functions (add, divide, typeIcon, etc.)
│
├── templates/               HTML pages (Go templates)
│   ├── index.html
│   ├── destination-detail.html
│   ├── itinerary-detail.html
│   ├── create-itinerary.html
│   └── search.html
│
├── static/                  Static assets
│   ├── css/style.css       Responsive styling
│   └── js/app.js           Client-side JavaScript
│
└── docs/
    ├── TEMPLATES_GUIDE.md  (this file)
    └── schema.sql         Database schema
```

## Project Flow

### 1. User Visits Home Page
```
http://localhost:8080/
  ↓
routes.go: router.GET("/", handlers.Index)
  ↓
handlers.go: Index() function
  ↓
Calls h.service.GetDestinations(1, 12)
  ↓
Returns destinations data
  ↓
Renders template: c.HTML(200, "index.html", gin.H{"destinations": destinations})
  ↓
Browser receives HTML with {{range .destinations}}...{{end}}
```

### 2. User Clicks "View Itineraries"
```
Click link to /destination/dest-001
  ↓
routes.go: router.GET("/destination/:id", handlers.DestinationDetail)
  ↓
handlers.go: DestinationDetail() function  
  ↓
Query itineraries for destination
  ↓
Render destination-detail.html with itineraries data
  ↓
User sees list of plans sorted by likes
```

### 3. User Clicks "Like"
```
JavaScript button click
  ↓
app.js: likeItinerary(itineraryId)
  ↓
Fetch POST /api/itineraries/:id/like
  ↓
API handler updates like count
  ↓
Returns JSON response
  ↓
JavaScript shows notification "Thanks for liking!"
```

## Development Workflow

### Edit Templates
1. Make changes to `templates/*.html`
2. Save file
3. Refresh browser (no rebuild needed!)

### Edit CSS
1. Edit `static/css/style.css`
2. Save file
3. Refresh browser (change already visible!)

### Edit Go Code
1. Make changes to `.go` files
2. Stop running server (Ctrl+C)
3. Run `go run main.go` again
4. Refresh browser

### Edit JavaScript
1. Edit `static/js/app.js`
2. Save file
3. Refresh browser (Ctrl+Shift+R for hard refresh)

## Configuration

Edit `config/config.json`:

```json
{
  "server": {
    "port": "8080",      ← Change server port here
    "timeout": 30,
    "mode": "debug"      ← Change to "release" for production
  },
  "database": {
    "host": "localhost",
    "port": "3306",
    "user": "root",
    "database": "itinerary_db",
    "password": ""       ← Will be overridden by env var
  }
}
```

## Adding a New Feature

### Example: Add "Trending" Page

#### 1. Add Handler (itinerary/handlers.go)
```go
func (h *Handlers) Trending(c *gin.Context) {
    itineraries, _, err := h.service.GetItinerariesByDestination("", 1, 20)
    if err != nil {
        c.HTML(500, "trending.html", gin.H{"error": err.Error()})
        return
    }
    c.HTML(200, "trending.html", gin.H{
        "itineraries": itineraries,
    })
}
```

#### 2. Add Route (itinerary/routes.go)
```go
router.GET("/trending", handlers.Trending)
```

#### 3. Create Template (templates/trending.html)
```html
<!DOCTYPE html>
<html>
<head>
    <title>Trending</title>
    <link rel="stylesheet" href="/static/css/style.css">
</head>
<body>
    <nav class="navbar">...</nav>
    <div class="container">
        <h1>Trending Itineraries</h1>
        {{range .itineraries}}
            <!-- Show itinerary -->
        {{end}}
    </div>
</body>
</html>
```

#### 4. Add Navigation Link
Add to navbar in templates:
```html
<li><a href="/trending">Trending</a></li>
```

Done! No build step required.

## Using Template Functions

### Example: Display costs
```html
<!-- Before -->
<span>₹{{.Budget}}</span>

<!-- After (with formatting) -->
<span>{{formatPrice .Budget}}</span>
```

### Example: Calculate per-day cost
```html
<span>₹{{divide .Budget .Duration}} per day</span>
```

### Example: Show activity type with icon
```html
<span>{{typeIcon .Type}} {{toUpper .Type}}</span>
```

See more functions in `template_helpers.go` and `TEMPLATES_GUIDE.md`.

## Troubleshooting

### Server won't start
```
error: DB_PASSWORD environment variable not set
```
→ Set the environment variable:
```bash
# Windows (Command Prompt)
set DB_PASSWORD=your_password

# Windows (PowerShell)
$env:DB_PASSWORD="your_password"

# Linux/Mac
export DB_PASSWORD=your_password
```

### Templates not found
```
error: no such file or directory: templates/index.html
```
→ Make sure you're running from project root:
```bash
cd d:\Learn\iternary\itinerary-backend
go run main.go
```

### Static files not loading
```
404: /static/css/style.css not found
```
→ Verify `static/css/style.css` exists and server has read permission

### Port already in use
```
error: listen tcp :8080: bind: address already in use
```
→ Either:
1. Kill other process using port 8080
2. Change port in `config/config.json`
3. Set `SERVER_PORT` environment variable

## Status Dashboard

Check server health:
```
GET http://localhost:8080/api/health

Response: {"status":"healthy"}
```

## Next Steps

1. ✅ **View pages** - Open http://localhost:8080/
2. ⏭️ **Connect database** - Implement database.go functions
3. ⏭️ **Add authentication** - User registration/login
4. ⏭️ **Deploy** - Run on production server

## Useful Commands

```bash
# Test API endpoint
curl http://localhost:8080/api/health

# Check if server is binding to port
netstat -an | grep 8080  (Linux/Mac)
netstat -ano | find "8080"  (Windows)

# Build executable
go build -o app.exe

# Run with output
go run main.go 2>&1

# Format code
go fmt ./...

# Run tests
go test ./...
```

## Tips & Tricks

### Reload page automatically on code change
Use a tool like `nodemon` or `air`:
```bash
# Install air
go install github.com/cosmtrek/air@latest

# Run project with auto-reload
air
```

### Debug template data
Add this to handlers:
```go
fmt.Println("Template data:", data)
```

### Check rendered HTML
Right-click page → View Source to see final HTML

### Use browser DevTools
- F12 or Ctrl+Shift+I to open DevTools
- Check Console for JavaScript errors
- Check Network tab for API calls

## Resources

- [Go text/template Docs](https://golang.org/pkg/text/template/)
- [Gin Framework](https://gin-gonic.com/)
- [HTML/CSS Reference](https://developer.mozilla.org/en-US/docs/Web/HTML)
- [Fetch API](https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API)

## Need Help?

Check these files:
- `TEMPLATES_GUIDE.md` - Deep dive into templates
- `README.md` - Project overview
- `config/config.json` - Configuration reference
- `schema.sql` - Database structure

**Happy coding!** 🚀
