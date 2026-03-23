# Go Templates Frontend Guide

## Overview

This project uses **Go's built-in `text/template` package** for server-side HTML rendering instead of React or any frontend framework. This means:

1. **All code is in Go** - Backend and frontend logic share the same language
2. **No build step** - No webpack, no npm, no yarn required  
3. **Server-side rendering** - HTML is generated on the server and sent to browser
4. **Minimal JavaScript** - Only used for interactions (likes, comments, etc.)

## Project Structure

```
templates/          ← HTML pages with Go template syntax
├── index.html      ← Home page
├── destination-detail.html
├── itinerary-detail.html
├── create-itinerary.html
└── search.html

static/            ← Static assets
├── css/
│   └── style.css   ← All styling (responsive design)
└── js/
    └── app.js      ← Client-side JavaScript (minimal)

itinerary/
├── routes.go       ← Register routes + template functions
├── handlers.go     ← Render templates with data
└── template_helpers.go ← Template functions (add, divide, etc.)
```

## How It Works

### 1. Request Flow
```
Browser Request
    ↓
Gin Router (routes.go)
    ↓
Handler Function (handlers.go)
    ↓
Database Query via Service
    ↓
Pass data to Template
    ↓
Template renders HTML with {{ .data }}
    ↓
Browser receives complete HTML page
```

### 2. Template Examples

#### Simple Variable
```html
<h1>{{.title}}</h1>
```

#### Loop Through Items
```html
{{range .destinations}}
  <h3>{{.Name}}</h3>
  <p>{{.Country}}</p>
{{end}}
```

#### Template Functions
```html
<!-- Built-in Go functions -->
{{add .page 1}}        ← Page 1 + 1 = 2
{{divide .total 7}}    ← Calculate average cost per day
{{toUpper .name}}      ← Convert to uppercase
{{typeIcon .type}}     ← Get emoji for activity type

<!-- Custom formatting -->
{{formatPrice .cost}}  ← Format as currency
{{truncate .text 50}}  ← Truncate long text
```

#### Conditions
```html
{{if .error}}
  <div class="alert alert-error">{{.error}}</div>
{{end}}

{{if gt .currentPage 1}}
  <a href="/?page={{sub .currentPage 1}}">← Previous</a>
{{end}}
```

## Available Template Functions

All defined in `template_helpers.go`:

| Function | Example | Result |
|----------|---------|--------|
| `add` | `{{add 5 3}}` | 8 |
| `sub` | `{{sub 10 2}}` | 8 |
| `divide` | `{{divide 100 5}}` | 20 |
| `multiply` | `{{multiply 10 2}}` | 20 |
| `toUpper` | `{{toUpper "hello"}}` | HELLO |
| `toLower` | `{{toLower "HELLO"}}` | hello |
| `typeIcon` | `{{typeIcon "stay"}}` | 🏨 |
| `formatPrice` | `{{formatPrice 5000}}` | ₹5000 |
| `truncate` | `{{truncate "hello world" 5}}` | hello... |

## Handler Pattern

```go
// handlers.go

// Index handles GET /
func (h *Handlers) Index(c *gin.Context) {
    // 1. Get data from service/database
    destinations, total, err := h.service.GetDestinations(page, pageSize)
    
    // 2. Prepare template data
    templateData := gin.H{
        "destinations": destinations,
        "total":        total,
        "page":         page,
    }
    
    // 3. Render template with data
    c.HTML(http.StatusOK, "index.html", templateData)
}
```

## Styling Approach

### Structure
- **One main CSS file** (`static/css/style.css`)
- **Responsive design** - Works on mobile, tablet, desktop
- **CSS Grid & Flexbox** - Modern layout without Bootstrap

### Key Classes
- `.container` - Max-width container
- `.grid` - Responsive grid layout
- `.btn-primary`, `.btn-secondary` - Button styles
- `.card` - Card components
- `.alert` - Alert messages

### Customizing Styles
1. Edit `static/css/style.css`
2. Changes apply instantly on page reload
3. No build step needed

## JavaScript Interactions

Most interactions use vanilla JavaScript with API calls:

```javascript
// Like an itinerary
async function likeItinerary(itineraryId) {
    const result = await apiCall('POST', `/api/itineraries/${itineraryId}/like`);
    showNotification('Added like!', 'success');
}

// Post comment
async function postComment(event, itineraryId) {
    const content = event.target.querySelector('textarea').value;
    await apiCall('POST', `/api/itineraries/${itineraryId}/comments`, {comment_id: content});
}
```

### Minimal JavaScript Libraries
Currently using **zero** external JS libraries:
- No jQuery
- No Bootstrap JS
- No React/Vue/Angular
- Pure browser APIs (fetch, DOM) + vanilla JS

## Adding New Pages

### Step 1: Create Handler
```go
// In handlers.go
func (h *Handlers) NewPage(c *gin.Context) {
    data := gin.H{
        "title": "My New Page",
    }
    c.HTML(http.StatusOK, "new-page.html", data)
}
```

### Step 2: Add Route
```go
// In routes.go
router.GET("/new-page", handlers.NewPage)
```

### Step 3: Create Template
```html
<!-- templates/new-page.html -->
<!DOCTYPE html>
<html>
<head>
    <title>{{.title}}</title>
    <link rel="stylesheet" href="/static/css/style.css">
</head>
<body>
    <nav class="navbar"><!-- navbar --></nav>
    <div class="container">
        <h1>{{.title}}</h1>
    </div>
    <footer><!-- footer --></footer>
</body>
</html>
```

## API Endpoints for JavaScript

All endpoints return JSON for client-side interactions:

```javascript
// Get destinations
GET /api/destinations

// Get itineraries
GET /api/destinations/:id/itineraries

// Like itinerary
POST /api/itineraries/:id/like

// Post comment
POST /api/itineraries/:id/comments
```

## Enhancing Interactivity

### Option 1: Stay with Vanilla JS
Current setup is perfect for basic interactions. Add more functions to `app.js` as needed.

### Option 2: Add HTMX
Drop-in library to enable dynamic interactions:
```html
<button hx-post="/api/itineraries/:id/like" hx-swap="innerHTML">Like</button>
```

### Option 3: Add Alpine.js
Lightweight reactive framework:
```html
<div x-data="{ open: false }">
  <button @click="open = !open">Toggle</button>
</div>
```

## Performance Tips

1. **Template Caching** - Go caches compiled templates automatically
2. **CSS Minification** - Consider minifying CSS in production
3. **Lazy Loading** - Images load on-demand
4. **CSS Grid** - Efficient layout rendering
5. **No JavaScript frameworks** - Faster page loads

## Debugging

### See Template Data
```go
fmt.Println("Template data:", templateData)
```

### Check Template Rendering
Look at browser's "View Source" to see final HTML

### Browser Console
`app.js` exports functions to `window` for testing:
```javascript
window.likeItinerary(id)
window.formatPrice(100)
```

## Common Issues

### Templates not found
- Check path: Must be `templates/*.html` relative to working directory
- Solution: Run from project root: `cd itinerary-backend && go run main.go`

### Styles not loading
- Check path: Must use `/static/css/style.css`
- Solution: Verify static files are in `static/` directory

### JavaScript functions undefined
- Check `app.js` is loaded: `<script src="/static/js/app.js"></script>`
- Functions must be exported to `window` object

## Next Steps

1. ✅ Basic templates created
2. ✅ Responsive CSS included  
3. ⏭️ Connect to database (next)
4. ⏭️ Implement pagination (next)
5. ⏭️ Add user authentication (later)

## Resources

- [Go text/template docs](https://golang.org/pkg/text/template/)
- [Gin HTML Rendering](https://gin-gonic.com/#html-rendering)
- [Modern CSS Guide](https://developer.mozilla.org/en-US/docs/Web/CSS)
- [Vanilla JS Tips](https://vanillajstoolkit.com/)
