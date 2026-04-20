# Frontend Setup & Development Guide

## Project Created: Itinerary Frontend ✅

Modern React-based frontend for the Itinerary trip planning platform.

## 📁 Directory Structure

```
itinerary-frontend/
├── public/
│   └── index.html                 # Main HTML template
├── src/
│   ├── components/                # Reusable UI components (future)
│   ├── pages/
│   │   ├── LoginPage.jsx          # ✅ Login page (email/password)
│   │   └── DashboardPage.jsx      # ✅ Cities list dashboard
│   ├── services/
│   │   └── api.js                 # ✅ Axios API client & services
│   ├── hooks/
│   │   └── useAuth.js             # ✅ Authentication context & hook
│   ├── styles/
│   │   └── (custom styles)
│   ├── App.jsx                    # ✅ Main app component with routing
│   ├── index.js                   # ✅ React entry point
│   └── index.css                  # ✅ Global styles
├── .env.example                   # Environment variables template
├── .gitignore                     # Git ignore file
├── .dockerignore                  # Docker ignore file
├── Dockerfile                     # Multi-stage Docker build
├── package.json                   # Project dependencies
├── vite.config.js                 # Vite build configuration
└── README.md                      # Frontend documentation
```

## ✅ Completed Components

### 1. **Login Page** (`LoginPage.jsx`)
- Email & password input fields
- Show/hide password toggle
- Form validation
- Demo credentials display
- Beautiful gradient background with animated blobs
- Error message display
- Loading state
- "Forgot password?" link (placeholder)

**Features**:
- React form handling with state
- API integration with error handling
- Local storage token management
- Responsive design (mobile-first)
- Accessible form elements

### 2. **Dashboard Page** (`DashboardPage.jsx`)
- List all cities/destinations from backend
- Search functionality (client-side filtering)
- City cards with images, description, stats
- Pagination support (previous/next buttons)
- User profile in header
- Logout functionality
- Responsive grid (1 col mobile, 2 col tablet, 3 col desktop)
- Loading skeleton states
- Error handling with retry

**Features**:
- Fetch cities from `GET /api/cities`
- Real-time search filtering
- Pagination (page, pageSize)
- City card click → navigate to city detail page
- User greeting in header
- Beautiful footer with quick links

### 3. **API Service** (`api.js`)
Comprehensive Axios-based API client with services for:

**Authentication**:
- `authService.login(email, password)` - Login user
- `authService.logout()` - Logout user
- `authService.getCurrentUser()` - Get user from localStorage
- `authService.isAuthenticated()` - Check auth status

**Cities & Destinations**:
- `citiesService.getCities(page, pageSize)` - Get all cities
- `citiesService.getCityById(cityId)` - Get single city

**Trip Posts (Community)**:
- `tripPostsService.getTripPostsByCity(cityId, page, pageSize)` - Get posts for city
- `tripPostsService.getTripPostById(postId)` - Get post details
- `tripPostsService.addTripPostToItinerary(tripPostId)` - Add post to itinerary

**User Trips**:
- `userTripsService.getUserTrips()` - List user's trips
- `userTripsService.getUserTrip(tripId)` - Get trip details
- `userTripsService.createUserTrip(tripData)` - Create new trip
- `userTripsService.publishTrip(tripId, postData)` - Publish trip

**Segments & Reviews**:
- `tripSegmentsService.markSegmentVisited(segmentId)` - Mark place visited
- `tripSegmentsService.submitReview(segmentId, rating, review)` - Submit review

**Features**:
- Automatic JWT token injection in headers
- Auto-redirect to login on 401 errors
- Base URL configuration from environment
- Error interceptors
- Consistent response handling

### 4. **Authentication Hook** (`useAuth.js`)
React Context-based authentication management

**Provided**:
- `user` - Current user object
- `isLoading` - Loading state
- `error` - Error messages
- `login(email, password)` - Login function
- `logout()` - Logout function
- `isAuthenticated` - Boolean auth status

**Features**:
- Auto-checks localStorage on mount
- Persistent login across page reloads
- Error handling with user feedback
- Loading states for async operations

### 5. **Main App Component** (`App.jsx`)
- React Router setup with routes
- Protected routes for authenticated pages
- Login page at `/login`
- Dashboard at `/dashboard`
- Root `/` redirects to dashboard
- Loading spinner while checking auth
- Automatic redirect to login if not authenticated

### 6. **Routing Structure**
```
/              → Redirects to /dashboard
/login         → Login page (public)
/dashboard     → Cities list (protected)
/* (404)       → Redirects to /dashboard
```

### 7. **Styling**
- **Framework**: Tailwind CSS (via CDN and npm)
- **Color Scheme**: Blue & purple gradients
- **Responsive**: Mobile-first responsive design
- **Components**: Prebuilt Tailwind classes
- **Animations**: Fade-in, slide-in, pulse effects
- **Icons**: Lucide React icons

## 🚀 Quick Start

### 1. Install Dependencies
```bash
cd itinerary-frontend
npm install
```

### 2. Setup Environment
```bash
cp .env.example .env
# Update .env with your API endpoints
```

### 3. Run Development Server
```bash
npm run dev
```
Opens at `http://localhost:3000`

### 4. Build for Production
```bash
npm run build
```

### 5. Preview Production Build
```bash
npm run preview
```

## 🐳 Docker Setup

### Build Docker Image
```bash
docker build -t itinerary-frontend:latest .
```

### Run Docker Container
```bash
docker run -p 3000:3000 \
  -e REACT_APP_API_URL=http://localhost:8080/api \
  -e REACT_APP_AUTH_URL=http://localhost:8080/auth \
  itinerary-frontend:latest
```

### Docker Compose (Backend + Frontend)
See root `docker-compose.yml` for full stack setup.

## 📝 Environment Variables

Create `.env` file with:
```
REACT_APP_API_URL=http://localhost:8080/api
REACT_APP_AUTH_URL=http://localhost:8080/auth
REACT_APP_GOOGLE_MAPS_API_KEY=your_key_here
```

## 🔐 Demo Credentials

Login with:
- Email: `demo@example.com`
- Password: `demo123456`

## 📋 Pages Created

| Page | Route | Status | Features |
|------|-------|--------|----------|
| Login | `/login` | ✅ Done | Email/password, demo creds, gradient UI |
| Dashboard | `/dashboard` | ✅ Done | City list, search, pagination, logout |
| City Posts | `/city/:cityId` | ⏳ Next | Trip posts feed |
| Trip Details | `/trip-posts/:postId` | ⏳ Next | All places, map, reviews |
| My Itinerary | `/my-itinerary/:tripId` | ⏳ Next | Plan trip, mark visited |
| My Trips | `/my-trips` | ⏳ Next | List user trips |

## 🔗 API Integration

All pages integrate with backend API endpoints:

```
Authentication:
  POST /auth/login

Cities:
  GET /api/cities

Trip Posts:
  GET /api/cities/:cityId/trip-posts
  POST /api/user-trips/add-from-post

User Trips:
  GET /api/user-trips
  POST /api/user-trips
```

## 🎨 UI/UX Highlights

### Login Page
- Beautiful gradient background (blue to purple)
- Animated floating blobs
- Card-based form layout
- Password visibility toggle
- Demo credentials info box
- Feature icons at bottom
- Mobile responsive

### Dashboard Page
- Header with user profile
- Search bar with icon
- City grid cards (3-2-1 responsive)
- Hover effects on cards
- Loading spinner
- Error messages with retry
- Pagination controls
- Footer with links
- Logout button in header

## 🛠️ Available Commands

```bash
npm run dev      # Start development server
npm run build    # Build for production
npm run preview  # Preview production build
npm install      # Install dependencies
npm update       # Update dependencies
```

## 📦 Dependencies

**Production**:
- `react@18.2.0` - UI library
- `react-dom@18.2.0` - DOM rendering
- `react-router-dom@6.16.0` - Routing
- `axios@1.6.0` - HTTP client
- `lucide-react@0.263.1` - Icons
- `tailwindcss@3.3.0` - Styling

**Development**:
- `vite@4.4.0` - Build tool
- `@vitejs/plugin-react@4.0.0` - React plugin

## ✨ Key Features

✅ **Completed**:
- JWT authentication flow
- Protected routes
- API service layer
- Reusable hooks
- Responsive design
- Error handling
- Loading states
- Search functionality
- Pagination

⏳ **Next to Build**:
- City detail page with trip posts
- Trip post detail page
- Itinerary planning interface
- Review modal component
- My trips page
- Community feed page

## 🔄 Next Steps

1. **City Page** - Display trip posts for a city
2. **Trip Post Detail** - Show all places in a trip
3. **Itinerary Interface** - Time-based planning (morning/afternoon/evening/night)
4. **Review Modal** - Submit reviews after marking visited
5. **Integration Testing** - Test with running backend

## 📚 Tech Stack Overview

```
Frontend: React 18 + Vite
Styling: Tailwind CSS
Routing: React Router v6
HTTP: Axios
Icons: Lucide React
Build: Vite (fast development)
Deployment: Docker
```

## 💡 Design Philosophy

- **Modern**: Latest React patterns (hooks, context)
- **Responsive**: Mobile-first approach
- **Accessible**: Semantic HTML, ARIA attributes
- **Performance**: Code splitting, lazy loading ready
- **Maintainable**: Component-based architecture
- **User-Friendly**: Beautiful UI with smooth interactions

## 🎯 User Flow

```
1. User visits /login
   ↓
2. Enters email/password
   ↓
3. Authenticated → Redirected to /dashboard
   ↓
4. Views cities grid
   ↓
5. Searches or clicks city
   ↓
6. Navigates to /city/:cityId
   ↓
7. Views trip posts from other travelers
   ↓
8. Clicks trip → /trip-posts/:postId
   ↓
9. Sees all places, map, reviews
   ↓
10. "Add to Itinerary" → Creates trip
    ↓
11. Views /my-itinerary/:tripId
    ↓
12. Plans by time of day
    ↓
13. Marks places visited → Review modal
    ↓
14. Submits reviews
    ↓
15. Reviews auto-appear in trip
```

## 🚀 Deployment

Frontend can be deployed to:
- Vercel
- Netlify
- AWS S3 + CloudFront
- Docker (provided)
- Any static hosting

## 📞 Support

For issues or questions:
1. Check `.env` configuration
2. Verify backend is running at `REACT_APP_API_URL`
3. Check browser console for errors
4. Verify JWT token in localStorage

---

**Status**: ✅ Frontend framework complete with login and dashboard pages
**Ready for**: City page, trip posts feed, and itinerary planning implementation
