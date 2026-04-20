import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { tripPostsService } from '../services/api';
import {
  ArrowLeft,
  Loader,
  MapPin,
  DollarSign,
  Calendar,
  Heart,
  Eye,
  Clock,
  Star,
  Share2,
  Map,
  Camera,
  AlertCircle,
} from 'lucide-react';

const TripPostDetailPage = () => {
  const { postId } = useParams();
  const navigate = useNavigate();
  const [tripPost, setTripPost] = useState(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');
  const [activeDay, setActiveDay] = useState(1);
  const [selectedPhoto, setSelectedPhoto] = useState(null);
  const [addingToItinerary, setAddingToItinerary] = useState(false);

  useEffect(() => {
    fetchTripPostDetails();
  }, [postId]);

  const fetchTripPostDetails = async () => {
    try {
      setIsLoading(true);
      setError('');
      const response = await tripPostsService.getTripPostById(postId);
      setTripPost(response.data || response);
    } catch (err) {
      setError('Failed to load trip details. Please try again.');
      console.error(err);
    } finally {
      setIsLoading(false);
    }
  };

  const handleAddToItinerary = async () => {
    try {
      setAddingToItinerary(true);
      await tripPostsService.addTripPostToItinerary(postId);
      alert('Trip added to your itinerary! Navigate to "My Trips" to customize it.');
      navigate('/my-trips');
    } catch (err) {
      alert('Failed to add trip. Please try again.');
    } finally {
      setAddingToItinerary(false);
    }
  };

  // Group places by day
  const getPlacesByDay = () => {
    if (!tripPost?.places) return {};
    const grouped = {};
    tripPost.places.forEach((place) => {
      if (!grouped[place.day]) {
        grouped[place.day] = [];
      }
      grouped[place.day].push(place);
    });
    return grouped;
  };

  // Sort places by time of day
  const timeOfDayOrder = {
    morning: 1,
    afternoon: 2,
    evening: 3,
    night: 4,
  };

  const getSortedPlaces = (places) => {
    return [...places].sort(
      (a, b) =>
        (timeOfDayOrder[a.time_of_day] || 0) - (timeOfDayOrder[b.time_of_day] || 0)
    );
  };

  // Get current day's places
  const placesGrouped = getPlacesByDay();
  const currentDayPlaces = placesGrouped[activeDay] ? getSortedPlaces(placesGrouped[activeDay]) : [];
  const allDays = Object.keys(placesGrouped).sort((a, b) => a - b);

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="text-center">
          <Loader className="w-12 h-12 text-blue-500 animate-spin mx-auto mb-4" />
          <p className="text-gray-600 font-medium">Loading trip details...</p>
        </div>
      </div>
    );
  }

  if (error || !tripPost) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="bg-white rounded-xl p-8 max-w-md w-full shadow-lg">
          <div className="flex items-center gap-3 mb-4">
            <AlertCircle className="w-6 h-6 text-red-500" />
            <h2 className="text-lg font-bold text-gray-800">Error Loading Trip</h2>
          </div>
          <p className="text-gray-600 mb-4">{error}</p>
          <div className="flex gap-2">
            <button
              onClick={fetchTripPostDetails}
              className="flex-1 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-lg transition"
            >
              Try Again
            </button>
            <button
              onClick={() => navigate('/dashboard')}
              className="flex-1 px-4 py-2 bg-gray-200 hover:bg-gray-300 text-gray-800 font-medium rounded-lg transition"
            >
              Go Back
            </button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white shadow-sm sticky top-0 z-40">
        <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-4 flex items-center justify-between">
          <button
            onClick={() => navigate(-1)}
            className="p-2 hover:bg-gray-100 rounded-lg transition"
          >
            <ArrowLeft className="w-6 h-6 text-gray-700" />
          </button>
          <h1 className="text-2xl font-bold text-gray-800 flex-1 ml-4 truncate">{tripPost.title}</h1>
          <button className="p-2 hover:bg-gray-100 rounded-lg transition">
            <Share2 className="w-6 h-6 text-gray-700" />
          </button>
        </div>
      </header>

      {/* Hero Image */}
      <div className="h-96 bg-gradient-to-br from-blue-400 to-purple-500 relative overflow-hidden">
        {tripPost.cover_image && (
          <img
            src={tripPost.cover_image}
            alt={tripPost.title}
            className="w-full h-full object-cover"
          />
        )}
        <div className="absolute inset-0 bg-black bg-opacity-20"></div>
      </div>

      {/* Trip Info Card */}
      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 -mt-20 relative z-10 mb-8">
        <div className="bg-white rounded-xl shadow-lg p-6">
          <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
            {/* Duration */}
            <div className="text-center p-4 bg-blue-50 rounded-lg">
              <div className="flex items-center justify-center mb-2">
                <Calendar className="w-6 h-6 text-blue-500" />
              </div>
              <p className="text-2xl font-bold text-gray-800">{tripPost.duration}</p>
              <p className="text-sm text-gray-600">Days</p>
            </div>

            {/* Total Expense */}
            <div className="text-center p-4 bg-green-50 rounded-lg">
              <div className="flex items-center justify-center mb-2">
                <DollarSign className="w-6 h-6 text-green-500" />
              </div>
              <p className="text-2xl font-bold text-gray-800">${tripPost.total_expense || 0}</p>
              <p className="text-sm text-gray-600">Total Cost</p>
            </div>

            {/* Places Count */}
            <div className="text-center p-4 bg-red-50 rounded-lg">
              <div className="flex items-center justify-center mb-2">
                <MapPin className="w-6 h-6 text-red-500" />
              </div>
              <p className="text-2xl font-bold text-gray-800">{tripPost.places?.length || 0}</p>
              <p className="text-sm text-gray-600">Places</p>
            </div>

            {/* Engagement */}
            <div className="text-center p-4 bg-purple-50 rounded-lg">
              <div className="flex items-center justify-center gap-2 mb-2">
                <Heart className="w-5 h-5 text-red-500 fill-current" />
                <Eye className="w-5 h-5 text-blue-500" />
              </div>
              <p className="text-sm text-gray-800">
                <span className="font-bold">{tripPost.likes || 0}</span> likes
              </p>
              <p className="text-sm text-gray-600">
                <span className="font-bold">{tripPost.views || 0}</span> views
              </p>
            </div>
          </div>

          {/* Description */}
          {tripPost.description && (
            <div className="mt-6 pt-6 border-t border-gray-200">
              <p className="text-gray-700 leading-relaxed">{tripPost.description}</p>
            </div>
          )}

          {/* CTA Button */}
          <div className="mt-6 flex gap-3">
            <button
              onClick={handleAddToItinerary}
              disabled={addingToItinerary}
              className="flex-1 px-6 py-3 bg-gradient-to-r from-blue-600 to-blue-700 hover:from-blue-700 hover:to-blue-800 disabled:opacity-50 text-white font-bold rounded-lg transition flex items-center justify-center gap-2"
            >
              {addingToItinerary ? (
                <>
                  <Loader className="w-5 h-5 animate-spin" />
                  Adding...
                </>
              ) : (
                <>
                  <Heart className="w-5 h-5" />
                  Add to My Itinerary
                </>
              )}
            </button>
            <button className="px-6 py-3 bg-gray-200 hover:bg-gray-300 text-gray-800 font-medium rounded-lg transition">
              Save Trip
            </button>
          </div>
        </div>
      </div>

      {/* Main Content */}
      <main className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 pb-12">
        {/* Day Selector */}
        {allDays.length > 1 && (
          <div className="mb-8">
            <h2 className="text-xl font-bold text-gray-800 mb-4">Select Day</h2>
            <div className="flex gap-2 overflow-x-auto pb-2">
              {allDays.map((day) => (
                <button
                  key={day}
                  onClick={() => setActiveDay(parseInt(day))}
                  className={`px-6 py-3 rounded-lg font-medium whitespace-nowrap transition ${
                    activeDay === parseInt(day)
                      ? 'bg-blue-600 text-white'
                      : 'bg-white text-gray-700 border border-gray-300 hover:border-blue-400'
                  }`}
                >
                  Day {day}
                </button>
              ))}
            </div>
          </div>
        )}

        {/* Places for Selected Day */}
        <div className="space-y-6">
          <h2 className="text-2xl font-bold text-gray-800">
            Day {activeDay} ({currentDayPlaces.length} places)
          </h2>

          {currentDayPlaces.length === 0 ? (
            <div className="bg-white rounded-xl p-8 text-center">
              <p className="text-gray-600">No places scheduled for this day</p>
            </div>
          ) : (
            currentDayPlaces.map((place, index) => (
              <div key={place.id} className="bg-white rounded-xl shadow-md overflow-hidden hover:shadow-lg transition">
                {/* Time Badge */}
                <div className="px-6 py-3 bg-gradient-to-r from-blue-50 to-purple-50 border-b border-gray-200 flex items-center gap-2">
                  <Clock className="w-5 h-5 text-blue-600" />
                  <span className="font-bold text-gray-800 capitalize">{place.time_of_day || 'Anytime'}</span>
                  <span className="text-gray-600">•</span>
                  <span className="text-gray-700">{index + 1} of {currentDayPlaces.length}</span>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-3 gap-6 p-6">
                  {/* Photos */}
                  <div className="md:col-span-1">
                    <div className="space-y-3">
                      {place.photos && place.photos.length > 0 ? (
                        <>
                          <div
                            onClick={() => setSelectedPhoto(place.photos[0])}
                            className="h-48 rounded-lg overflow-hidden bg-gray-300 cursor-pointer group"
                          >
                            <img
                              src={place.photos[0].url}
                              alt={place.name}
                              className="w-full h-full object-cover group-hover:scale-110 transition"
                            />
                          </div>
                          {place.photos.length > 1 && (
                            <div className="flex gap-2 overflow-x-auto">
                              {place.photos.map((photo, idx) => (
                                <button
                                  key={idx}
                                  onClick={() => setSelectedPhoto(photo)}
                                  className="w-12 h-12 rounded-lg overflow-hidden flex-shrink-0 border-2 border-gray-300 hover:border-blue-500"
                                >
                                  <img
                                    src={photo.url}
                                    alt="photo"
                                    className="w-full h-full object-cover"
                                  />
                                </button>
                              ))}
                            </div>
                          )}
                        </>
                      ) : (
                        <div className="h-48 rounded-lg bg-gray-300 flex items-center justify-center text-gray-600">
                          <Camera className="w-8 h-8" />
                        </div>
                      )}
                    </div>
                  </div>

                  {/* Details */}
                  <div className="md:col-span-2">
                    {/* Place Name & Type */}
                    <div className="mb-4">
                      <h3 className="text-2xl font-bold text-gray-800">{place.name}</h3>
                      <span className="inline-block mt-2 px-3 py-1 bg-blue-100 text-blue-700 rounded-full text-sm font-medium capitalize">
                        {place.type || 'Activity'}
                      </span>
                    </div>

                    {/* Location */}
                    <div className="mb-4 flex items-start gap-2">
                      <MapPin className="w-5 h-5 text-red-500 flex-shrink-0 mt-1" />
                      <div>
                        <p className="text-gray-700">{place.location}</p>
                        {place.latitude && place.longitude && (
                          <a
                            href={`https://www.google.com/maps/@${place.latitude},${place.longitude},15z`}
                            target="_blank"
                            rel="noopener noreferrer"
                            className="text-sm text-blue-600 hover:text-blue-700 font-medium"
                          >
                            View on Google Maps
                          </a>
                        )}
                      </div>
                    </div>

                    {/* Cost */}
                    <div className="mb-4 flex items-center gap-2">
                      <DollarSign className="w-5 h-5 text-green-500" />
                      <span className="text-lg font-bold text-gray-800">${place.expense || 0}</span>
                      {place.time_of_day && (
                        <span className="text-gray-600">• {place.time_of_day}</span>
                      )}
                    </div>

                    {/* Best Time to Visit */}
                    {place.best_time_to_visit && (
                      <div className="mb-4 p-3 bg-amber-50 border border-amber-200 rounded-lg">
                        <p className="text-sm font-medium text-amber-900">
                          💡 Best time to visit: {place.best_time_to_visit}
                        </p>
                      </div>
                    )}

                    {/* Review */}
                    {place.review && (
                      <div className="mt-4 p-4 bg-purple-50 rounded-lg border border-purple-200">
                        <div className="flex items-center gap-2 mb-2">
                          <div className="flex items-center gap-1">
                            {[1, 2, 3, 4, 5].map((star) => (
                              <Star
                                key={star}
                                className={`w-4 h-4 ${
                                  star <= place.review.rating
                                    ? 'text-yellow-400 fill-current'
                                    : 'text-gray-300'
                                }`}
                              />
                            ))}
                          </div>
                          <span className="font-bold text-gray-800">
                            {place.review.rating}/5
                          </span>
                        </div>
                        <p className="text-gray-700 text-sm">{place.review.review}</p>
                      </div>
                    )}

                    {/* Mark as Visited (if user owns trip) */}
                    <div className="mt-6 flex gap-2">
                      <button className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-lg transition">
                        Mark as Visited
                      </button>
                      <button className="px-4 py-2 bg-gray-200 hover:bg-gray-300 text-gray-800 font-medium rounded-lg transition">
                        Add Notes
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            ))
          )}
        </div>

        {/* Summary */}
        <div className="mt-12 p-8 bg-white rounded-xl shadow-md">
          <h3 className="text-xl font-bold text-gray-800 mb-4">Trip Summary</h3>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            <div>
              <p className="text-gray-600 text-sm">Total Duration</p>
              <p className="text-2xl font-bold text-gray-800">{tripPost.duration} days</p>
            </div>
            <div>
              <p className="text-gray-600 text-sm">Total Budget</p>
              <p className="text-2xl font-bold text-green-600">${tripPost.total_expense || 0}</p>
            </div>
            <div>
              <p className="text-gray-600 text-sm">Total Places</p>
              <p className="text-2xl font-bold text-blue-600">{tripPost.places?.length || 0}</p>
            </div>
            <div>
              <p className="text-gray-600 text-sm">Average/Day</p>
              <p className="text-2xl font-bold text-purple-600">
                ${tripPost.total_expense && tripPost.duration ? (tripPost.total_expense / tripPost.duration).toFixed(2) : 0}
              </p>
            </div>
          </div>
        </div>
      </main>

      {/* Photo Modal */}
      {selectedPhoto && (
        <div
          className="fixed inset-0 bg-black bg-opacity-75 z-50 flex items-center justify-center p-4"
          onClick={() => setSelectedPhoto(null)}
        >
          <div className="max-w-4xl w-full" onClick={(e) => e.stopPropagation()}>
            <button
              onClick={() => setSelectedPhoto(null)}
              className="mb-4 text-white hover:text-gray-200 transition"
            >
              ✕ Close
            </button>
            <img
              src={selectedPhoto.url}
              alt="Full view"
              className="w-full rounded-lg"
            />
            {selectedPhoto.caption && (
              <p className="mt-4 text-white text-center">{selectedPhoto.caption}</p>
            )}
          </div>
        </div>
      )}
    </div>
  );
};

export default TripPostDetailPage;
