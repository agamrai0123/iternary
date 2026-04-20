# Itinerary Frontend

Modern React-based frontend for the Itinerary trip planning application.

## Features

- 🔐 Secure login authentication
- 🌍 Browse cities and destinations
- 📋 Discover trip posts from other travelers
- 🗺️ View trip details with maps
- ✈️ Add trips to your personal itinerary
- ⭐ Rate and review places
- 📱 Fully responsive design

## Tech Stack

- **Framework**: React 18
- **Build Tool**: Vite
- **Routing**: React Router v6
- **HTTP Client**: Axios
- **Styling**: Tailwind CSS
- **Icons**: Lucide React

## Getting Started

### Prerequisites

- Node.js 16+
- npm or yarn

### Installation

1. Install dependencies:
```bash
npm install
```

2. Create `.env` file from `.env.example`:
```bash
cp .env.example .env
```

3. Update `.env` with your API endpoints:
```
REACT_APP_API_URL=http://localhost:8080/api
REACT_APP_AUTH_URL=http://localhost:8080/auth
```

### Development

Run the development server:
```bash
npm run dev
```

The application will open at `http://localhost:3000`

### Building

Create production build:
```bash
npm run build
```

Preview production build:
```bash
npm run preview
```

## Project Structure

```
itinerary-frontend/
├── public/
│   └── index.html
├── src/
│   ├── components/      # Reusable components
│   ├── pages/           # Page components
│   │   ├── LoginPage.jsx
│   │   └── DashboardPage.jsx
│   ├── services/        # API services
│   │   └── api.js
│   ├── hooks/           # Custom React hooks
│   │   └── useAuth.js
│   ├── styles/          # Global styles
│   ├── App.jsx
│   ├── index.js
│   └── index.css
├── package.json
├── vite.config.js
└── README.md
```

## API Endpoints

### Authentication
- `POST /auth/login` - User login
- `POST /auth/logout` - User logout

### Cities & Destinations
- `GET /api/cities` - List all cities
- `GET /api/cities/:cityId` - Get city details

### Trip Posts (Community Feed)
- `GET /api/cities/:cityId/trip-posts` - Get trip posts for a city
- `GET /api/trip-posts/:postId` - Get trip post details
- `POST /api/user-trips/add-from-post` - Add trip to itinerary

### User Trips
- `GET /api/user-trips` - List user's trips
- `GET /api/user-trips/:tripId` - Get trip details
- `POST /api/user-trips` - Create new trip
- `PUT /api/user-trips/:tripId` - Update trip

### Reviews & Segments
- `POST /api/trip-segments/:segmentId/mark-visited` - Mark place as visited
- `POST /api/reviews` - Submit review

## Pages

### Login Page (`/login`)
- Email and password input
- Demo credentials display
- Beautiful gradient background
- Responsive design

### Dashboard Page (`/dashboard`)
- List of all cities/destinations
- Search functionality
- City cards with images
- Pagination support
- User profile and logout

## Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run preview` - Preview production build

## Styling

The project uses Tailwind CSS for styling. All CSS classes follow Tailwind conventions.

### Custom CSS

Custom styles are defined in `src/index.css`:
- Fade-in animations
- Slide-in animations
- Custom scrollbar styling
- Form focus styles

## Components

### Planned Components

- `CityCard` - Display city information
- `TripPostCard` - Display trip post
- `ItineraryTimeline` - Show trip schedule by time
- `ReviewModal` - Submit reviews
- `MapView` - Display location on map

## Authentication

The app uses JWT tokens stored in localStorage. The `useAuth` hook provides authentication context and methods:

```javascript
const { user, login, logout, isAuthenticated } = useAuth();
```

## Error Handling

All API calls include error handling with user-friendly error messages displayed in the UI.

## Future Enhancements

- [ ] City detail page with trip posts feed
- [ ] Trip post detail page
- [ ] My itinerary planning page
- [ ] Review submission modal
- [ ] Community feed page
- [ ] User profile page
- [ ] Settings page
- [ ] Favorites/bookmarks
- [ ] Social sharing
- [ ] Dark mode

## License

MIT License - See LICENSE file for details

## Support

For issues or questions, please create an issue in the GitHub repository.
