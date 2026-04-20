import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { userTripsService } from '../services/api';
import {
  ArrowLeft,
  Loader,
  Plus,
  Trash2,
  Edit,
  MapPin,
  DollarSign,
  Calendar,
  CheckCircle,
  Clock,
  AlertCircle,
  Eye,
} from 'lucide-react';

const MyTripsPage = () => {
  const navigate = useNavigate();
  const [trips, setTrips] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');
  const [filter, setFilter] = useState('all'); // all, draft, planning, ongoing, completed
  const [deleteConfirm, setDeleteConfirm] = useState(null);

  useEffect(() => {
    fetchUserTrips();
  }, []);

  const fetchUserTrips = async () => {
    try {
      setIsLoading(true);
      setError('');
      const response = await userTripsService.getUserTrips();
      setTrips(response.data || response || []);
    } catch (err) {
      setError('Failed to load your trips. Please try again.');
      console.error(err);
    } finally {
      setIsLoading(false);
    }
  };

  const getFilteredTrips = () => {
    if (filter === 'all') return trips;
    return trips.filter((trip) => trip.status === filter);
  };

  const getStatusColor = (status) => {
    switch (status) {
      case 'draft':
        return 'bg-gray-100 text-gray-800';
      case 'planning':
        return 'bg-blue-100 text-blue-800';
      case 'ongoing':
        return 'bg-purple-100 text-purple-800';
      case 'completed':
        return 'bg-green-100 text-green-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  const getStatusIcon = (status) => {
    switch (status) {
      case 'draft':
        return '📝';
      case 'planning':
        return '📋';
      case 'ongoing':
        return '✈️';
      case 'completed':
        return '✅';
      default:
        return '📌';
    }
  };

  const filteredTrips = getFilteredTrips();

  // Trip Card Component
  const TripCard = ({ trip }) => (
    <div className="bg-white rounded-xl shadow-md hover:shadow-lg transition-shadow overflow-hidden">
      {/* Card Header with Status */}
      <div className="px-6 py-4 bg-gradient-to-r from-blue-50 to-purple-50 border-b border-gray-200 flex items-center justify-between">
        <div className="flex items-center gap-2">
          <span className="text-2xl">{getStatusIcon(trip.status)}</span>
          <span className={`px-3 py-1 rounded-full text-xs font-bold capitalize ${getStatusColor(trip.status)}`}>
            {trip.status}
          </span>
        </div>
        <div className="text-sm text-gray-600">
          Created {new Date(trip.created_at).toLocaleDateString()}
        </div>
      </div>

      {/* Card Content */}
      <div className="p-6">
        {/* Trip Title */}
        <h3 className="text-xl font-bold text-gray-800 mb-2 truncate">{trip.title}</h3>

        {/* Trip Stats */}
        <div className="grid grid-cols-3 gap-3 mb-4 py-3 border-y border-gray-200">
          <div className="text-center">
            <MapPin className="w-4 h-4 text-red-500 mx-auto mb-1" />
            <p className="text-xs text-gray-600">Destination</p>
            <p className="text-sm font-bold text-gray-800">{trip.destination_id?.slice(0, 10) || '?'}</p>
          </div>
          <div className="text-center">
            <Calendar className="w-4 h-4 text-blue-500 mx-auto mb-1" />
            <p className="text-xs text-gray-600">Duration</p>
            <p className="text-sm font-bold text-gray-800">{trip.duration} days</p>
          </div>
          <div className="text-center">
            <DollarSign className="w-4 h-4 text-green-500 mx-auto mb-1" />
            <p className="text-xs text-gray-600">Budget</p>
            <p className="text-sm font-bold text-gray-800">${trip.budget}</p>
          </div>
        </div>

        {/* Places Count */}
        {trip.segments && (
          <div className="mb-4 flex items-center gap-2 text-gray-700">
            <Eye className="w-4 h-4 text-gray-500" />
            <span className="text-sm">{trip.segments.length} places planned</span>
          </div>
        )}

        {/* Start Date */}
        {trip.start_date && (
          <div className="mb-4 p-3 bg-blue-50 rounded-lg border border-blue-200">
            <p className="text-sm text-blue-900">
              <strong>Start Date:</strong> {new Date(trip.start_date).toLocaleDateString()}
            </p>
          </div>
        )}

        {/* Action Buttons */}
        <div className="flex gap-2">
          <button
            onClick={() => navigate(`/my-itinerary/${trip.id}`)}
            className="flex-1 px-4 py-2 bg-gradient-to-r from-blue-500 to-blue-600 hover:from-blue-600 hover:to-blue-700 text-white font-medium rounded-lg transition text-sm"
          >
            Plan Trip
          </button>
          <button
            onClick={() => {
              // Edit trip - navigate to edit page (not created yet)
              alert('Edit functionality coming soon');
            }}
            className="px-4 py-2 bg-gray-200 hover:bg-gray-300 text-gray-800 font-medium rounded-lg transition"
            title="Edit trip"
          >
            <Edit className="w-4 h-4" />
          </button>
          <button
            onClick={() => setDeleteConfirm(trip.id)}
            className="px-4 py-2 bg-red-100 hover:bg-red-200 text-red-600 font-medium rounded-lg transition"
            title="Delete trip"
          >
            <Trash2 className="w-4 h-4" />
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
            <h1 className="text-2xl font-bold text-gray-800">My Trips</h1>
          </div>

          {/* Create Trip Button */}
          <button
            onClick={() => {
              // Create new trip - navigate to create page (not created yet)
              alert('Create trip functionality coming soon');
            }}
            className="flex items-center gap-2 px-4 py-2 bg-gradient-to-r from-blue-600 to-blue-700 hover:from-blue-700 hover:to-blue-800 text-white font-medium rounded-lg transition"
          >
            <Plus className="w-5 h-5" />
            New Trip
          </button>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Status Summary */}
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-8">
          <div
            onClick={() => setFilter('all')}
            className={`p-4 rounded-lg cursor-pointer transition ${
              filter === 'all'
                ? 'bg-blue-600 text-white shadow-lg'
                : 'bg-white text-gray-800 hover:shadow-md'
            }`}
          >
            <p className="text-2xl font-bold">{trips.length}</p>
            <p className="text-sm">All Trips</p>
          </div>

          <div
            onClick={() => setFilter('draft')}
            className={`p-4 rounded-lg cursor-pointer transition ${
              filter === 'draft'
                ? 'bg-gray-600 text-white shadow-lg'
                : 'bg-white text-gray-800 hover:shadow-md'
            }`}
          >
            <p className="text-2xl font-bold">{trips.filter((t) => t.status === 'draft').length}</p>
            <p className="text-sm">Draft</p>
          </div>

          <div
            onClick={() => setFilter('planning')}
            className={`p-4 rounded-lg cursor-pointer transition ${
              filter === 'planning'
                ? 'bg-blue-600 text-white shadow-lg'
                : 'bg-white text-gray-800 hover:shadow-md'
            }`}
          >
            <p className="text-2xl font-bold">{trips.filter((t) => t.status === 'planning').length}</p>
            <p className="text-sm">Planning</p>
          </div>

          <div
            onClick={() => setFilter('completed')}
            className={`p-4 rounded-lg cursor-pointer transition ${
              filter === 'completed'
                ? 'bg-green-600 text-white shadow-lg'
                : 'bg-white text-gray-800 hover:shadow-md'
            }`}
          >
            <p className="text-2xl font-bold">{trips.filter((t) => t.status === 'completed').length}</p>
            <p className="text-sm">Completed</p>
          </div>
        </div>

        {/* Filter Tabs */}
        <div className="mb-8 flex gap-2 overflow-x-auto pb-2">
          {['all', 'draft', 'planning', 'ongoing', 'completed'].map((status) => (
            <button
              key={status}
              onClick={() => setFilter(status)}
              className={`px-4 py-2 rounded-lg font-medium whitespace-nowrap transition capitalize ${
                filter === status
                  ? 'bg-blue-600 text-white'
                  : 'bg-white text-gray-700 border border-gray-300 hover:border-blue-400'
              }`}
            >
              {status === 'all' ? 'All Trips' : status}
            </button>
          ))}
        </div>

        {/* Error Message */}
        {error && (
          <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg flex items-start gap-3">
            <AlertCircle className="w-5 h-5 text-red-600 flex-shrink-0 mt-0.5" />
            <div>
              <p className="text-red-700 font-medium">{error}</p>
              <button
                onClick={fetchUserTrips}
                className="mt-2 text-red-600 hover:text-red-700 font-medium text-sm"
              >
                Try Again
              </button>
            </div>
          </div>
        )}

        {/* Loading State */}
        {isLoading ? (
          <div className="flex flex-col items-center justify-center py-20">
            <Loader className="w-12 h-12 text-blue-500 animate-spin mb-4" />
            <p className="text-gray-600 font-medium">Loading your trips...</p>
          </div>
        ) : filteredTrips.length === 0 ? (
          <div className="bg-white rounded-xl p-12 text-center">
            <div className="text-6xl mb-4">🗺️</div>
            <h3 className="text-xl font-bold text-gray-800 mb-2">No trips yet</h3>
            <p className="text-gray-600 mb-6">
              {filter === 'all'
                ? "You haven't created any trips yet. Go explore some destinations!"
                : `No ${filter} trips found.`}
            </p>
            <button
              onClick={() => navigate('/dashboard')}
              className="px-6 py-3 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-lg transition"
            >
              Explore Destinations
            </button>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {filteredTrips.map((trip) => (
              <TripCard key={trip.id} trip={trip} />
            ))}
          </div>
        )}
      </main>

      {/* Delete Confirmation Modal */}
      {deleteConfirm && (
        <div className="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-center justify-center p-4">
          <div className="bg-white rounded-xl p-6 max-w-sm w-full shadow-2xl">
            <div className="flex items-center gap-3 mb-4">
              <AlertCircle className="w-6 h-6 text-red-600" />
              <h3 className="text-lg font-bold text-gray-800">Delete Trip?</h3>
            </div>
            <p className="text-gray-600 mb-6">
              Are you sure you want to delete this trip? This action cannot be undone.
            </p>
            <div className="flex gap-3">
              <button
                onClick={() => setDeleteConfirm(null)}
                className="flex-1 px-4 py-2 bg-gray-200 hover:bg-gray-300 text-gray-800 font-medium rounded-lg transition"
              >
                Cancel
              </button>
              <button
                onClick={() => {
                  // Delete trip
                  setTrips(trips.filter((t) => t.id !== deleteConfirm));
                  setDeleteConfirm(null);
                  alert('Trip deleted successfully');
                }}
                className="flex-1 px-4 py-2 bg-red-600 hover:bg-red-700 text-white font-medium rounded-lg transition"
              >
                Delete
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default MyTripsPage;
