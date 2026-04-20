import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { userTripsService } from '../services/api';
import {
  ArrowLeft,
  Loader,
  Plus,
  Trash2,
  Edit2,
  MapPin,
  DollarSign,
  Calendar,
  Clock,
  AlertCircle,
  Save,
  X,
  GripVertical,
  ChevronDown,
  ChevronUp,
  Copy,
} from 'lucide-react';

const MyItineraryPage = () => {
  const { tripId } = useParams();
  const navigate = useNavigate();
  const [trip, setTrip] = useState(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const [error, setError] = useState('');
  const [expandedDay, setExpandedDay] = useState(1);
  const [editingPlace, setEditingPlace] = useState(null);
  const [addingPlace, setAddingPlace] = useState(false);
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(null);

  useEffect(() => {
    fetchTripDetails();
  }, [tripId]);

  const fetchTripDetails = async () => {
    try {
      setIsLoading(true);
      setError('');
      const response = await userTripsService.getUserTripById(tripId);
      setTrip(response.data || response);
    } catch (err) {
      setError('Failed to load trip details. Please try again.');
      console.error(err);
    } finally {
      setIsLoading(false);
    }
  };

  const handleSaveTrip = async () => {
    try {
      setIsSaving(true);
      await userTripsService.updateUserTrip(tripId, trip);
      alert('Trip saved successfully!');
    } catch (err) {
      alert('Failed to save trip. Please try again.');
    } finally {
      setIsSaving(false);
    }
  };

  const handleDeletePlace = (placeId) => {
    setTrip({
      ...trip,
      segments: trip.segments.filter((p) => p.id !== placeId),
    });
    setShowDeleteConfirm(null);
  };

  const getPlacesByDay = () => {
    if (!trip?.segments) return {};
    const grouped = {};
    trip.segments.forEach((place) => {
      if (!grouped[place.day]) {
        grouped[place.day] = [];
      }
      grouped[place.day].push(place);
    });
    return grouped;
  };

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

  const calculateDayExpense = (day) => {
    const places = getPlacesByDay()[day] || [];
    return places.reduce((sum, place) => sum + (place.expense || 0), 0);
  };

  const placesGrouped = getPlacesByDay();
  const allDays = Object.keys(placesGrouped).sort((a, b) => a - b);
  const totalExpense = trip?.segments?.reduce((sum, p) => sum + (p.expense || 0), 0) || 0;

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

  if (error || !trip) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="bg-white rounded-xl p-8 max-w-md w-full shadow-lg">
          <div className="flex items-center gap-3 mb-4">
            <AlertCircle className="w-6 h-6 text-red-500" />
            <h2 className="text-lg font-bold text-gray-800">Error</h2>
          </div>
          <p className="text-gray-600 mb-4">{error}</p>
          <button
            onClick={() => navigate('/my-trips')}
            className="w-full px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-lg transition"
          >
            Go Back
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white shadow-sm sticky top-0 z-40">
        <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
          <div className="flex items-center justify-between mb-2">
            <div className="flex items-center gap-4">
              <button
                onClick={() => navigate('/my-trips')}
                className="p-2 hover:bg-gray-100 rounded-lg transition"
              >
                <ArrowLeft className="w-6 h-6 text-gray-700" />
              </button>
              <div>
                <h1 className="text-2xl font-bold text-gray-800">{trip.title}</h1>
                <p className="text-sm text-gray-600">
                  {trip.segments?.length || 0} places • {trip.duration} days • $
                  {totalExpense}
                </p>
              </div>
            </div>

            {/* Save Button */}
            <button
              onClick={handleSaveTrip}
              disabled={isSaving}
              className="flex items-center gap-2 px-4 py-2 bg-gradient-to-r from-green-500 to-green-600 hover:from-green-600 hover:to-green-700 disabled:opacity-50 text-white font-medium rounded-lg transition"
            >
              {isSaving ? (
                <>
                  <Loader className="w-5 h-5 animate-spin" />
                  Saving...
                </>
              ) : (
                <>
                  <Save className="w-5 h-5" />
                  Save Trip
                </>
              )}
            </button>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Trip Stats */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-8">
          <div className="bg-white rounded-lg p-4 shadow-sm">
            <p className="text-gray-600 text-sm">Total Days</p>
            <p className="text-3xl font-bold text-gray-800">{trip.duration}</p>
          </div>
          <div className="bg-white rounded-lg p-4 shadow-sm">
            <p className="text-gray-600 text-sm">Total Places</p>
            <p className="text-3xl font-bold text-blue-600">{trip.segments?.length || 0}</p>
          </div>
          <div className="bg-white rounded-lg p-4 shadow-sm">
            <p className="text-gray-600 text-sm">Total Budget</p>
            <p className="text-3xl font-bold text-green-600">${totalExpense}</p>
          </div>
          <div className="bg-white rounded-lg p-4 shadow-sm">
            <p className="text-gray-600 text-sm">Per Day</p>
            <p className="text-3xl font-bold text-purple-600">
              ${trip.duration ? (totalExpense / trip.duration).toFixed(2) : 0}
            </p>
          </div>
        </div>

        {/* Days Timeline */}
        <div className="space-y-4">
          {allDays.length === 0 ? (
            <div className="bg-white rounded-xl p-12 text-center">
              <p className="text-6xl mb-4">📍</p>
              <h3 className="text-xl font-bold text-gray-800 mb-2">No places added yet</h3>
              <p className="text-gray-600 mb-6">Start building your itinerary by adding places</p>
              <button
                onClick={() => setAddingPlace(true)}
                className="px-6 py-2 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-lg transition"
              >
                Add First Place
              </button>
            </div>
          ) : (
            allDays.map((day) => {
              const dayPlaces = getSortedPlaces(placesGrouped[day] || []);
              const dayExpense = calculateDayExpense(day);
              const isExpanded = expandedDay === parseInt(day);

              return (
                <div key={day} className="bg-white rounded-xl shadow-md overflow-hidden">
                  {/* Day Header */}
                  <button
                    onClick={() =>
                      setExpandedDay(isExpanded ? null : parseInt(day))
                    }
                    className="w-full px-6 py-4 bg-gradient-to-r from-blue-50 to-purple-50 hover:from-blue-100 hover:to-purple-100 border-b border-gray-200 flex items-center justify-between transition"
                  >
                    <div className="flex items-center gap-4 text-left">
                      <div className="w-10 h-10 rounded-full bg-blue-600 text-white flex items-center justify-center font-bold">
                        {day}
                      </div>
                      <div>
                        <h3 className="font-bold text-gray-800">Day {day}</h3>
                        <p className="text-sm text-gray-600">
                          {dayPlaces.length} places • ${dayExpense}
                        </p>
                      </div>
                    </div>
                    <div className="flex items-center gap-2">
                      {isExpanded ? (
                        <ChevronUp className="w-5 h-5 text-gray-600" />
                      ) : (
                        <ChevronDown className="w-5 h-5 text-gray-600" />
                      )}
                    </div>
                  </button>

                  {/* Day Content */}
                  {isExpanded && (
                    <div className="p-6">
                      {dayPlaces.length === 0 ? (
                        <div className="text-center py-8">
                          <p className="text-gray-600">No places added for this day</p>
                        </div>
                      ) : (
                        <div className="space-y-4">
                          {dayPlaces.map((place, index) => (
                            <div
                              key={place.id}
                              className="p-4 bg-gray-50 rounded-lg border border-gray-200 hover:border-blue-300 transition"
                            >
                              <div className="flex items-start gap-4">
                                {/* Drag Handle */}
                                <div className="pt-2 text-gray-400 cursor-grab active:cursor-grabbing">
                                  <GripVertical className="w-5 h-5" />
                                </div>

                                {/* Place Info */}
                                <div className="flex-1 min-w-0">
                                  {/* Time & Index */}
                                  <div className="flex items-center gap-2 mb-2">
                                    <span className="px-2 py-1 bg-blue-100 text-blue-700 text-xs font-bold rounded capitalize">
                                      {place.time_of_day || 'Anytime'}
                                    </span>
                                    <span className="text-xs text-gray-600">
                                      {index + 1} of {dayPlaces.length}
                                    </span>
                                  </div>

                                  {/* Place Name & Type */}
                                  <h4 className="font-bold text-gray-800 mb-1">{place.name}</h4>
                                  {place.type && (
                                    <p className="text-sm text-gray-600 capitalize">
                                      Type: {place.type}
                                    </p>
                                  )}

                                  {/* Location */}
                                  <div className="flex items-start gap-2 mt-2 text-sm text-gray-700">
                                    <MapPin className="w-4 h-4 text-red-500 flex-shrink-0 mt-0.5" />
                                    <span>{place.location}</span>
                                  </div>

                                  {/* Cost */}
                                  <div className="flex items-center gap-2 mt-2 text-sm">
                                    <DollarSign className="w-4 h-4 text-green-500" />
                                    <span className="font-bold text-gray-800">
                                      ${place.expense || 0}
                                    </span>
                                  </div>

                                  {/* Notes */}
                                  {place.notes && (
                                    <div className="mt-3 p-2 bg-white rounded border-l-4 border-blue-400">
                                      <p className="text-sm text-gray-700">{place.notes}</p>
                                    </div>
                                  )}

                                  {/* Photos */}
                                  {place.photos && place.photos.length > 0 && (
                                    <div className="flex gap-2 mt-3">
                                      {place.photos.slice(0, 3).map((photo, idx) => (
                                        <img
                                          key={idx}
                                          src={photo.url}
                                          alt="place"
                                          className="w-16 h-16 rounded-lg object-cover"
                                        />
                                      ))}
                                      {place.photos.length > 3 && (
                                        <div className="w-16 h-16 rounded-lg bg-gray-300 flex items-center justify-center text-gray-600 text-xs font-bold">
                                          +{place.photos.length - 3}
                                        </div>
                                      )}
                                    </div>
                                  )}
                                </div>

                                {/* Actions */}
                                <div className="flex items-center gap-2 flex-shrink-0">
                                  <button
                                    onClick={() => setEditingPlace(place)}
                                    className="p-2 hover:bg-blue-100 text-blue-600 rounded-lg transition"
                                    title="Edit place"
                                  >
                                    <Edit2 className="w-4 h-4" />
                                  </button>
                                  <button
                                    onClick={() => setShowDeleteConfirm(place.id)}
                                    className="p-2 hover:bg-red-100 text-red-600 rounded-lg transition"
                                    title="Delete place"
                                  >
                                    <Trash2 className="w-4 h-4" />
                                  </button>
                                </div>
                              </div>
                            </div>
                          ))}
                        </div>
                      )}

                      {/* Add Place Button */}
                      <button
                        onClick={() => setAddingPlace(true)}
                        className="mt-4 w-full px-4 py-2 border-2 border-dashed border-gray-300 hover:border-blue-400 text-gray-700 hover:text-blue-600 font-medium rounded-lg transition flex items-center justify-center gap-2"
                      >
                        <Plus className="w-4 h-4" />
                        Add Place to Day {day}
                      </button>
                    </div>
                  )}
                </div>
              );
            })
          )}
        </div>

        {/* Add Place Button */}
        {allDays.length > 0 && (
          <button
            onClick={() => setAddingPlace(true)}
            className="mt-8 w-full px-6 py-3 bg-gradient-to-r from-blue-600 to-blue-700 hover:from-blue-700 hover:to-blue-800 text-white font-medium rounded-lg transition flex items-center justify-center gap-2"
          >
            <Plus className="w-5 h-5" />
            Add New Place
          </button>
        )}
      </main>

      {/* Add/Edit Place Modal */}
      {(addingPlace || editingPlace) && (
        <div className="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-center justify-center p-4">
          <div className="bg-white rounded-xl p-6 max-w-md w-full shadow-2xl max-h-96 overflow-y-auto">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-lg font-bold text-gray-800">
                {editingPlace ? 'Edit Place' : 'Add Place'}
              </h3>
              <button
                onClick={() => {
                  setEditingPlace(null);
                  setAddingPlace(false);
                }}
                className="p-1 hover:bg-gray-100 rounded-lg transition"
              >
                <X className="w-6 h-6 text-gray-600" />
              </button>
            </div>

            {/* Form (Simplified for now) */}
            <div className="space-y-3">
              <input
                type="text"
                placeholder="Place name"
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:border-blue-500"
              />
              <input
                type="text"
                placeholder="Location"
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:border-blue-500"
              />
              <select className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:border-blue-500">
                <option>Select day</option>
                {Array.from({ length: trip.duration }, (_, i) => (
                  <option key={i + 1} value={i + 1}>
                    Day {i + 1}
                  </option>
                ))}
              </select>
              <select className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:border-blue-500">
                <option value="">Select time</option>
                <option value="morning">Morning</option>
                <option value="afternoon">Afternoon</option>
                <option value="evening">Evening</option>
                <option value="night">Night</option>
              </select>
              <input
                type="number"
                placeholder="Cost"
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:border-blue-500"
              />
              <textarea
                placeholder="Notes"
                rows="3"
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:border-blue-500"
              />

              <div className="flex gap-2 pt-4">
                <button
                  onClick={() => {
                    setEditingPlace(null);
                    setAddingPlace(false);
                  }}
                  className="flex-1 px-4 py-2 bg-gray-200 hover:bg-gray-300 text-gray-800 font-medium rounded-lg transition"
                >
                  Cancel
                </button>
                <button
                  onClick={() => {
                    alert('Place added/updated successfully');
                    setEditingPlace(null);
                    setAddingPlace(false);
                  }}
                  className="flex-1 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-lg transition"
                >
                  Save
                </button>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Delete Confirmation Modal */}
      {showDeleteConfirm && (
        <div className="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-center justify-center p-4">
          <div className="bg-white rounded-xl p-6 max-w-sm w-full shadow-2xl">
            <div className="flex items-center gap-3 mb-4">
              <AlertCircle className="w-6 h-6 text-red-600" />
              <h3 className="text-lg font-bold text-gray-800">Delete Place?</h3>
            </div>
            <p className="text-gray-600 mb-6">
              Are you sure you want to remove this place from your itinerary? This action cannot be undone.
            </p>
            <div className="flex gap-3">
              <button
                onClick={() => setShowDeleteConfirm(null)}
                className="flex-1 px-4 py-2 bg-gray-200 hover:bg-gray-300 text-gray-800 font-medium rounded-lg transition"
              >
                Cancel
              </button>
              <button
                onClick={() => handleDeletePlace(showDeleteConfirm)}
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

export default MyItineraryPage;
