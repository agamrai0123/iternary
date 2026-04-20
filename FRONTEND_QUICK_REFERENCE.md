# Frontend Quick Reference Guide

## 🚀 Getting Started

### Run Frontend Locally
```bash
cd itinerary-frontend
npm install
npm start
# Runs on http://localhost:3000
```

### Project Structure
```
itinerary-frontend/
├── src/
│   ├── pages/                    # Page components (6 pages)
│   │   ├── LoginPage.jsx         # Authentication
│   │   ├── DashboardPage.jsx     # Cities list
│   │   ├── CityPage.jsx          # Trip posts feed
│   │   ├── TripPostDetailPage.jsx # Trip details ⭐
│   │   ├── MyTripsPage.jsx       # User trips ⭐
│   │   └── MyItineraryPage.jsx   # Trip planning ⭐
│   ├── services/
│   │   └── api.js               # API client & services
│   ├── hooks/
│   │   └── useAuth.js           # Auth context
│   ├── App.jsx                  # Routing
│   ├── index.css                # Tailwind styles
│   └── index.js                 # Entry point
├── public/
│   └── index.html
├── package.json
└── tailwind.config.js           # Tailwind config
```

---

## 📄 Page Reference

### 1. Login Page `/login`
```javascript
import LoginPage from './pages/LoginPage';

// Features:
- Email/password input
- Remember me checkbox
- "Forgot password" link
- Google OAuth button
- Sign up link
- Form validation
```

### 2. Dashboard `/dashboard`
```javascript
import DashboardPage from './pages/DashboardPage';

// Features:
- Cities grid/list
- Search cities
- Quick stats
- "My Trips" button
```

### 3. City Page `/city/:cityId`
```javascript
import CityPage from './pages/CityPage';

// Features:
- Trip posts grid
- Filter & sort
- "View Details" button
- "Add Trip" button
- Author info cards
```

### 4. Trip Post Detail `/trip-posts/:postId` ⭐
```javascript
import TripPostDetailPage from './pages/TripPostDetailPage';

// Features:
- Hero image
- Trip stats (duration, cost, places, likes, views)
- Trip description
- Day selector tabs
- Places grouped by day & time
- Place photos gallery
- Location with Google Maps link
- Reviews with ratings
- Add to itinerary button
- Save trip button
```

### 5. My Trips `/my-trips` ⭐
```javascript
import MyTripsPage from './pages/MyTripsPage';

// Features:
- Trips list with status
- Status filter tabs
- Trip cards with stats
- Plan trip button
- Edit/delete buttons
- Summary dashboard
```

### 6. My Itinerary `/my-itinerary/:tripId` ⭐
```javascript
import MyItineraryPage from './pages/MyItineraryPage';

// Features:
- Collapsible days
- Places by time of day
- Budget tracking (real-time)
- Edit place modal
- Add place modal
- Delete confirmation
- Save trip button
```

---

## 🔌 API Services

### Authentication Service
```javascript
import { authService } from './services/api';

await authService.login(email, password);
await authService.register(userData);
await authService.loginWithGoogle(token);
await authService.logout();
```

### Cities Service
```javascript
import { citiesService } from './services/api';

const cities = await citiesService.getCities();
const city = await citiesService.getCityById(cityId);
const posts = await citiesService.getTripPostsByCity(cityId);
```

### Trip Posts Service
```javascript
import { tripPostsService } from './services/api';

const post = await tripPostsService.getTripPostById(postId);
const posts = await tripPostsService.searchTripPosts(query);
await tripPostsService.addTripPostToItinerary(postId);
```

### User Trips Service
```javascript
import { userTripsService } from './services/api';

const trips = await userTripsService.getUserTrips();
const trip = await userTripsService.getUserTripById(tripId);
await userTripsService.updateUserTrip(tripId, data);
await userTripsService.deleteUserTrip(tripId);
```

---

## 🎯 Common Tasks

### Add a New Page
1. Create component in `src/pages/NewPage.jsx`
2. Import in `App.jsx`
3. Add route to Routes section
4. Wrap with `<ProtectedRoute>` if needed

### Use Authentication
```javascript
import { useAuth } from './hooks/useAuth';

function MyComponent() {
  const { user, isAuthenticated, logout } = useAuth();
  
  if (!isAuthenticated) {
    return <Navigate to="/login" />;
  }
  
  return <div>Welcome {user.name}</div>;
}
```

### Call API
```javascript
import { tripPostsService } from './services/api';

useEffect(() => {
  const fetchData = async () => {
    try {
      const data = await tripPostsService.getTripPostById(id);
      setData(data);
    } catch (error) {
      setError(error.message);
    }
  };
  fetchData();
}, [id]);
```

### Handle Loading State
```javascript
const [isLoading, setIsLoading] = useState(true);

if (isLoading) {
  return (
    <div className="flex items-center justify-center">
      <Loader className="w-8 h-8 animate-spin" />
    </div>
  );
}
```

### Show Error Message
```javascript
const [error, setError] = useState('');

if (error) {
  return (
    <div className="p-4 bg-red-50 border border-red-200 rounded-lg">
      <p className="text-red-700">{error}</p>
      <button onClick={retry} className="text-red-600">
        Try Again
      </button>
    </div>
  );
}
```

### Navigate
```javascript
import { useNavigate } from 'react-router-dom';

const navigate = useNavigate();

// Go to page
navigate('/dashboard');

// Go back
navigate(-1);

// Go with params
navigate(`/city/${cityId}`);
```

### Use Modal
```javascript
const [showModal, setShowModal] = useState(false);

{showModal && (
  <div className="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-center justify-center">
    <div className="bg-white rounded-lg p-6">
      <h2>Modal Title</h2>
      {/* Content */}
      <button onClick={() => setShowModal(false)}>Close</button>
    </div>
  </div>
)}
```

---

## 🎨 Tailwind Classes

### Colors
```javascript
// Primary
bg-blue-600, text-blue-700, border-blue-500

// Secondary
bg-purple-600, text-purple-700

// Success
bg-green-600, text-green-700

// Danger
bg-red-600, text-red-700

// Neutral
bg-gray-100, text-gray-800, border-gray-300
```

### Common Patterns
```javascript
// Full-width button
className="w-full px-4 py-2 bg-blue-600"

// Card
className="bg-white rounded-lg shadow-md p-6"

// Header
className="sticky top-0 z-40 bg-white shadow-sm"

// Modal
className="fixed inset-0 bg-black bg-opacity-50 z-50"

// Grid
className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4"

// Flex
className="flex items-center justify-between gap-4"
```

---

## 🔧 Debugging

### Check Auth
```javascript
const { user, isAuthenticated } = useAuth();
console.log('User:', user);
console.log('Authenticated:', isAuthenticated);
```

### Check Route Params
```javascript
const { postId } = useParams();
console.log('Post ID:', postId);
```

### Check API Response
```javascript
try {
  const response = await api.get('/endpoint');
  console.log('Response:', response);
} catch (error) {
  console.error('Error:', error.message);
}
```

### Network Tab
Open DevTools → Network → Check API calls and responses

---

## 📱 Responsive Breakpoints

```javascript
// Mobile first (< 640px)
className="..."

// Tablet (≥ 640px)
className="sm:..."

// Small desktop (≥ 768px)
className="md:..."

// Desktop (≥ 1024px)
className="lg:..."

// Large desktop (≥ 1280px)
className="xl:..."
```

---

## 🚨 Common Issues

### Issue: 404 on page reload
**Solution:** Backend needs catch-all route or use HashRouter

### Issue: Blank page
**Solution:** Check browser console for errors, verify API endpoint

### Issue: Images not loading
**Solution:** Check image URL, check CORS headers, use absolute paths

### Issue: API calls failing
**Solution:** Check backend is running, verify API base URL, check network tab

### Issue: Auth not working
**Solution:** Clear localStorage, check tokens, verify auth service

---

## 📊 Performance Tips

1. **Lazy load pages**
   ```javascript
   const TripPostDetail = lazy(() => import('./pages/TripPostDetailPage'));
   ```

2. **Memoize components**
   ```javascript
   export default memo(TripCard);
   ```

3. **Optimize images**
   - Use WEBP format
   - Set width/height
   - Use lazy loading

4. **Debounce search**
   ```javascript
   const debouncedSearch = debounce(handleSearch, 300);
   ```

5. **Virtual lists for long lists**
   - Install: `react-window`
   - Use for 1000+ items

---

## 🐛 Testing

### Unit test example
```javascript
import { render, screen } from '@testing-library/react';
import LoginPage from './pages/LoginPage';

test('renders login form', () => {
  render(<LoginPage />);
  expect(screen.getByText('Login')).toBeInTheDocument();
});
```

### Component test
```javascript
test('submits form', async () => {
  render(<LoginPage />);
  const button = screen.getByText('Login');
  fireEvent.click(button);
  // Assert
});
```

---

## 🚀 Deployment

### Build
```bash
npm run build
# Creates optimized build in build/ folder
```

### Environment Variables
Create `.env`:
```
REACT_APP_API_BASE_URL=https://api.example.com
REACT_APP_GOOGLE_CLIENT_ID=your-client-id
```

### Deploy to Netlify
```bash
npm run build
# Drag build folder to netlify.com
```

### Deploy to Vercel
```bash
vercel --prod
```

---

## 📚 Resources

- [React Documentation](https://react.dev)
- [React Router](https://reactrouter.com)
- [Tailwind CSS](https://tailwindcss.com)
- [Lucide Icons](https://lucide.dev)
- [Axios](https://axios-http.com)

---

## 👥 Team

- **Frontend Lead:** [Your Name]
- **Last Updated:** Phase 2 Sprint 1 Complete
- **Branch:** main
- **Status:** Ready for Backend Integration 🚀

---

## 📞 Quick Commands

```bash
# Start dev server
npm start

# Build for production
npm run build

# Run tests
npm test

# Format code
npm run format

# Lint code
npm run lint
```

**Happy Coding! 🎉**
