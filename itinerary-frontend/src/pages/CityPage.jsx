import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { tripPostsService, citiesService } from '../services/api';
import { ArrowLeft, Loader, MapPin, DollarSign, Calendar, Heart, Eye, Map } from 'lucide-react';

const CityPage = () => {
  const { cityId } = useParams();
  const navigate = useNavigate();
  const [city, setCity] = useState(null);
  const [tripPosts, setTripPosts] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [sortBy, setSortBy] = useState('latest');
  const pageSize = 10;

  useEffect(() => {
    fetchCityAndPosts();
  }, [cityId, page, sortBy]);

  const fetchCityAndPosts = async () => {
    try {
      setIsLoading(true);
      setError('');

      // Fetch city details and trip posts in parallel
      const [cityResponse, postsResponse] = await Promise.all([
        citiesService.getCityById(cityId),
        tripPostsService.getTripPostsByCity(cityId, page, pageSize),
      ]);

      setCity(cityResponse);
      
      // Sort posts based on selected criteria
      let sortedPosts = postsResponse.data || [];
      if (sortBy === 'popular') {
        sortedPosts.sort((a, b) => (b.likes || 0) - (a.likes || 0));
      } else if (sortBy === 'trending') {
        sortedPosts.sort((a, b) => (b.views || 0) - (a.views || 0));
      }
      // 'latest' is default from API
      
      setTripPosts(sortedPosts);
      setTotalPages(Math.ceil((postsResponse.total || 0) / pageSize));
    } catch (err) {
      setError('Failed to load city or trip posts. Please try again.');
      console.error(err);
    } finally {
      setIsLoading(false);
    }
  };

  const handleAddToItinerary = async (postId) => {
    try {
      await tripPostsService.addTripPostToItinerary(postId);
      // Show success message or toast
      alert('Trip added to your itinerary! Navigate to "My Trips" to plan your journey.');
    } catch (err) {
      alert('Failed to add trip. Please try again.');
    }
  };

  // Trip Post Card Component
  const TripPostCard = ({ post }) => (
    <div className="bg-white rounded-xl overflow-hidden shadow-md hover:shadow-lg transition-shadow">
      {/* Post Image */}
      <div className="relative h-48 bg-gradient-to-br from-blue-300 to-purple-400 overflow-hidden group">
        {post.cover_image ? (
          <img
            src={post.cover_image}
            alt={post.title}
            className="w-full h-full object-cover group-hover:scale-105 transition-transform"
          />
        ) : (
          <div className="w-full h-full flex items-center justify-center text-white text-4xl">
            ✈️
          </div>
        )}
        
        {/* Overlay Stats */}
        <div className="absolute inset-0 bg-black bg-opacity-20 group-hover:bg-opacity-40 transition opacity-0 group-hover:opacity-100 flex items-center justify-center">
          <button
            onClick={() => handleAddToItinerary(post.id)}
            className="px-6 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-bold transition"
          >
            Add to Itinerary
          </button>
        </div>

        {/* Like Badge */}
        <div className="absolute top-3 right-3 bg-white rounded-full px-3 py-1 flex items-center gap-1 shadow-md">
          <Heart className="w-4 h-4 text-red-500 fill-current" />
          <span className="text-sm font-bold text-gray-800">{post.likes || 0}</span>
        </div>
      </div>

      {/* Post Info */}
      <div className="p-4">
        {/* Title */}
        <h3
          className="text-lg font-bold text-gray-800 mb-2 cursor-pointer hover:text-blue-600 transition truncate"
          onClick={() => navigate(`/trip-posts/${post.id}`)}
        >
          {post.title}
        </h3>

        {/* Description */}
        <p className="text-gray-600 text-sm mb-3 line-clamp-2">{post.description}</p>

        {/* Trip Stats */}
        <div className="grid grid-cols-3 gap-3 mb-3 py-3 border-t border-b border-gray-200">
          {/* Days */}
          <div className="text-center">
            <Calendar className="w-4 h-4 text-blue-500 mx-auto mb-1" />
            <p className="text-xs font-bold text-gray-800">{post.duration || '?'} Days</p>
          </div>

          {/* Budget */}
          <div className="text-center">
            <DollarSign className="w-4 h-4 text-green-500 mx-auto mb-1" />
            <p className="text-xs font-bold text-gray-800">${post.total_expense || '0'}</p>
          </div>

          {/* Places Count */}
          <div className="text-center">
            <MapPin className="w-4 h-4 text-red-500 mx-auto mb-1" />
            <p className="text-xs font-bold text-gray-800">{post.segments?.length || 0} Places</p>
          </div>
        </div>

        {/* Engagement Stats */}
        <div className="flex items-center justify-between text-xs text-gray-600 mb-3">
          <span className="flex items-center gap-1">
            <Eye className="w-4 h-4" /> {post.views || 0} views
          </span>
          <span className="text-gray-500">
            Published {new Date(post.published_at).toLocaleDateString()}
          </span>
        </div>

        {/* Author Info */}
        <div className="flex items-center gap-2 p-2 bg-gray-50 rounded-lg">
          <div className="w-8 h-8 bg-gradient-to-br from-blue-400 to-purple-500 rounded-full flex items-center justify-center text-white text-xs font-bold">
            {post.user_name?.charAt(0) || '?'}
          </div>
          <div className="flex-1 min-w-0">
            <p className="text-xs font-medium text-gray-800 truncate">{post.user_name || 'Anonymous'}</p>
            <p className="text-xs text-gray-500">Traveler</p>
          </div>
        </div>

        {/* CTA Buttons */}
        <div className="flex gap-2 mt-3">
          <button
            onClick={() => navigate(`/trip-posts/${post.id}`)}
            className="flex-1 px-3 py-2 bg-blue-50 hover:bg-blue-100 text-blue-600 font-medium rounded-lg transition text-sm"
          >
            View Details
          </button>
          <button
            onClick={() => handleAddToItinerary(post.id)}
            className="flex-1 px-3 py-2 bg-gradient-to-r from-blue-500 to-blue-600 hover:from-blue-600 hover:to-blue-700 text-white font-medium rounded-lg transition text-sm"
          >
            Add Trip
          </button>
        </div>
      </div>
    </div>
  );

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100">
      {/* Header */}
      <header className="bg-white shadow-sm sticky top-0 z-40">
        <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-4 flex items-center justify-between">
          <div className="flex items-center gap-4">
            <button
              onClick={() => navigate('/dashboard')}
              className="p-2 hover:bg-gray-100 rounded-lg transition"
            >
              <ArrowLeft className="w-6 h-6 text-gray-700" />
            </button>
            <div>
              <h1 className="text-2xl font-bold text-gray-800">{city?.name || 'City'}</h1>
              <p className="text-sm text-gray-600">{city?.country || ''}</p>
            </div>
          </div>

          {/* Actions */}
          <button
            onClick={() => navigate('/my-trips')}
            className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-lg transition"
          >
            My Trips
          </button>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* City Description */}
        {city?.description && (
          <div className="mb-8 p-6 bg-white rounded-xl shadow-sm border-l-4 border-blue-500">
            <p className="text-gray-700">{city.description}</p>
            {city.best_time_to_visit && (
              <p className="text-sm text-gray-600 mt-2">
                <strong>Best time to visit:</strong> {city.best_time_to_visit}
              </p>
            )}
          </div>
        )}

        {/* Filter & Sort Bar */}
        <div className="mb-6 flex items-center justify-between">
          <h2 className="text-xl font-bold text-gray-800">
            Trip Posts ({tripPosts.length})
          </h2>
          <div className="flex items-center gap-3">
            <label className="text-sm font-medium text-gray-700">Sort by:</label>
            <select
              value={sortBy}
              onChange={(e) => {
                setSortBy(e.target.value);
                setPage(1);
              }}
              className="px-3 py-2 rounded-lg border border-gray-300 focus:ring-2 focus:ring-blue-500 outline-none"
            >
              <option value="latest">Latest</option>
              <option value="popular">Most Popular</option>
              <option value="trending">Trending</option>
            </select>
          </div>
        </div>

        {/* Error Message */}
        {error && (
          <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg">
            <p className="text-red-700 font-medium">{error}</p>
            <button
              onClick={fetchCityAndPosts}
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
            <p className="text-gray-600 font-medium">Loading amazing trip posts...</p>
          </div>
        ) : tripPosts.length === 0 ? (
          <div className="text-center py-20 bg-white rounded-xl">
            <div className="text-6xl mb-4">📭</div>
            <p className="text-gray-600 text-lg font-medium">
              No trip posts available for this city yet.
            </p>
            <p className="text-gray-500 text-sm mt-2">
              Be the first to create and share your trip here!
            </p>
          </div>
        ) : (
          <>
            {/* Trip Posts Grid */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
              {tripPosts.map((post) => (
                <TripPostCard key={post.id} post={post} />
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
        <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
            <div>
              <h3 className="text-lg font-bold text-gray-800 mb-2">{city?.name}</h3>
              <p className="text-sm text-gray-600">{city?.description}</p>
            </div>
            <div>
              <h4 className="font-semibold text-gray-800 mb-3">Quick Links</h4>
              <ul className="text-sm text-gray-600 space-y-2">
                <li><a href="#" className="hover:text-blue-600">Back to Dashboard</a></li>
                <li><a href="#" className="hover:text-blue-600">My Trips</a></li>
                <li><a href="#" className="hover:text-blue-600">Create Trip</a></li>
              </ul>
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

export default CityPage;
