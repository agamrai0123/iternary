import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';
import { citiesService } from '../services/api';
import { LogOut, Search, Loader } from 'lucide-react';

const DashboardPage = () => {
  const navigate = useNavigate();
  const { user, logout } = useAuth();
  const [cities, setCities] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');
  const [searchQuery, setSearchQuery] = useState('');
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);

  useEffect(() => {
    fetchCities();
  }, [page]);

  const fetchCities = async () => {
    try {
      setIsLoading(true);
      setError('');
      const response = await citiesService.getCities(page, 12);
      setCities(response.data || []);
      setTotalPages(Math.ceil((response.total || 0) / 12));
    } catch (err) {
      setError('Failed to load cities. Please try again.');
      console.error(err);
    } finally {
      setIsLoading(false);
    }
  };

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const handleCityClick = (cityId) => {
    navigate(`/city/${cityId}`);
  };

  const filteredCities = cities.filter((city) =>
    city.name?.toLowerCase().includes(searchQuery.toLowerCase())
  );

  // City card component
  const CityCard = ({ city }) => (
    <div
      onClick={() => handleCityClick(city.id)}
      className="bg-white rounded-xl overflow-hidden shadow-lg hover:shadow-xl transition-all cursor-pointer transform hover:-translate-y-2"
    >
      {/* City Image */}
      <div className="relative h-48 bg-gradient-to-br from-blue-400 to-purple-500 overflow-hidden">
        {city.image ? (
          <img
            src={city.image}
            alt={city.name}
            className="w-full h-full object-cover hover:scale-110 transition-transform"
          />
        ) : (
          <div className="w-full h-full flex items-center justify-center text-white text-4xl">
            📍
          </div>
        )}
        {/* Overlay */}
        <div className="absolute inset-0 bg-black bg-opacity-20 hover:bg-opacity-40 transition"></div>
      </div>

      {/* City Info */}
      <div className="p-4">
        <h3 className="text-xl font-bold text-gray-800 mb-1">{city.name}</h3>
        <p className="text-gray-600 text-sm mb-3 line-clamp-2">{city.description}</p>

        {/* Stats */}
        <div className="flex justify-between text-xs text-gray-500 mb-3 pt-3 border-t border-gray-200">
          <span>📍 {city.country || 'Unknown'}</span>
          <span>✈️ {Math.floor(Math.random() * 50) + 10} posts</span>
        </div>

        {/* CTA Button */}
        <button className="w-full bg-gradient-to-r from-blue-500 to-blue-600 hover:from-blue-600 hover:to-blue-700 text-white font-medium py-2 rounded-lg transition">
          Explore Trips
        </button>
      </div>
    </div>
  );

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100">
      {/* Header */}
      <header className="bg-white shadow-sm sticky top-0 z-40">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 flex items-center justify-between">
          {/* Logo */}
          <div className="flex items-center gap-3">
            <div className="text-3xl">✈️</div>
            <div>
              <h1 className="text-2xl font-bold text-gray-800">Itinerary</h1>
              <p className="text-xs text-gray-500">Discover & Plan Trips</p>
            </div>
          </div>

          {/* User Info & Logout */}
          <div className="flex items-center gap-4">
            <div className="text-right hidden sm:block">
              <p className="text-sm font-medium text-gray-800">{user?.name || user?.email}</p>
              <p className="text-xs text-gray-500">Explorer</p>
            </div>
            <div className="w-10 h-10 bg-gradient-to-br from-blue-500 to-purple-500 rounded-full flex items-center justify-center text-white font-bold">
              {user?.name?.charAt(0) || user?.email?.charAt(0) || 'U'}
            </div>
            <button
              onClick={handleLogout}
              className="flex items-center gap-2 px-4 py-2 text-gray-700 hover:text-red-600 hover:bg-red-50 rounded-lg transition"
              title="Logout"
            >
              <LogOut className="w-5 h-5" />
              <span className="hidden sm:inline text-sm font-medium">Logout</span>
            </button>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        {/* Welcome Section */}
        <div className="mb-12">
          <h2 className="text-4xl font-bold text-gray-900 mb-3">
            Welcome to Your Travel Adventure! 🌍
          </h2>
          <p className="text-gray-600 text-lg">
            Explore amazing trip itineraries shared by travelers around the world. Discover
            destinations, see detailed plans, and start planning your next adventure.
          </p>
        </div>

        {/* Search Bar */}
        <div className="mb-10">
          <div className="relative max-w-md">
            <Search className="absolute left-3 top-3 text-gray-400 w-5 h-5" />
            <input
              type="text"
              placeholder="Search cities..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="w-full pl-10 pr-4 py-3 rounded-lg border border-gray-300 focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none"
            />
          </div>
        </div>

        {/* Error Message */}
        {error && (
          <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg">
            <p className="text-red-700 font-medium">{error}</p>
            <button
              onClick={fetchCities}
              className="mt-2 text-red-600 hover:text-red-700 font-medium text-sm"
            >
              Try Again
            </button>
          </div>
        )}

        {/* Loading State */}
        {isLoading ? (
          <div className="flex flex-col items-center justify-center py-20">
            <Loader className="w-12 h-12 text-blue-500 animate-spin mb-4" />
            <p className="text-gray-600 font-medium">Loading amazing destinations...</p>
          </div>
        ) : filteredCities.length === 0 ? (
          <div className="text-center py-20">
            <div className="text-6xl mb-4">🔍</div>
            <p className="text-gray-600 text-lg font-medium">
              {searchQuery ? 'No cities found matching your search' : 'No cities available'}
            </p>
          </div>
        ) : (
          <>
            {/* Cities Grid */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-10">
              {filteredCities.map((city) => (
                <CityCard key={city.id} city={city} />
              ))}
            </div>

            {/* Pagination */}
            {totalPages > 1 && (
              <div className="flex justify-center items-center gap-2">
                <button
                  onClick={() => setPage(Math.max(1, page - 1))}
                  disabled={page === 1}
                  className="px-4 py-2 rounded-lg bg-white text-gray-700 border border-gray-300 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition"
                >
                  ← Previous
                </button>
                {Array.from({ length: totalPages }).map((_, i) => (
                  <button
                    key={i + 1}
                    onClick={() => setPage(i + 1)}
                    className={`px-3 py-2 rounded-lg transition ${
                      page === i + 1
                        ? 'bg-blue-600 text-white'
                        : 'bg-white text-gray-700 border border-gray-300 hover:bg-gray-50'
                    }`}
                  >
                    {i + 1}
                  </button>
                ))}
                <button
                  onClick={() => setPage(Math.min(totalPages, page + 1))}
                  disabled={page === totalPages}
                  className="px-4 py-2 rounded-lg bg-white text-gray-700 border border-gray-300 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition"
                >
                  Next →
                </button>
              </div>
            )}
          </>
        )}
      </main>

      {/* Footer */}
      <footer className="bg-white border-t border-gray-200 mt-16">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            <div>
              <h3 className="text-lg font-bold text-gray-800 mb-2">Itinerary</h3>
              <p className="text-sm text-gray-600">
                Your gateway to discovering and sharing amazing travel experiences.
              </p>
            </div>
            <div>
              <h4 className="font-semibold text-gray-800 mb-3">Quick Links</h4>
              <ul className="text-sm text-gray-600 space-y-2">
                <li><a href="#" className="hover:text-blue-600">My Trips</a></li>
                <li><a href="#" className="hover:text-blue-600">Community</a></li>
                <li><a href="#" className="hover:text-blue-600">Help</a></li>
              </ul>
            </div>
            <div>
              <h4 className="font-semibold text-gray-800 mb-3">Follow Us</h4>
              <div className="flex gap-3 text-sm">
                <a href="#" className="text-gray-600 hover:text-blue-600">Twitter</a>
                <a href="#" className="text-gray-600 hover:text-blue-600">Facebook</a>
                <a href="#" className="text-gray-600 hover:text-blue-600">Instagram</a>
              </div>
            </div>
          </div>
          <div className="border-t border-gray-200 mt-8 pt-8 text-center text-sm text-gray-600">
            <p>&copy; 2026 Itinerary. All rights reserved.</p>
          </div>
        </div>
      </footer>
    </div>
  );
};

export default DashboardPage;
