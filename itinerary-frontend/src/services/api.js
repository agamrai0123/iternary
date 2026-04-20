import axios from 'axios';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080/api';
const AUTH_BASE_URL = process.env.REACT_APP_AUTH_URL || 'http://localhost:8080/auth';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add auth token to requests
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Handle unauthorized responses
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// ==================== Authentication ====================

export const authService = {
  login: async (email, password) => {
    const response = await axios.post(`${AUTH_BASE_URL}/login`, {
      email,
      password,
    });
    const { token, user } = response.data;
    localStorage.setItem('token', token);
    localStorage.setItem('user', JSON.stringify(user));
    return { token, user };
  },

  logout: () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
  },

  getCurrentUser: () => {
    const user = localStorage.getItem('user');
    return user ? JSON.parse(user) : null;
  },

  isAuthenticated: () => {
    return !!localStorage.getItem('token');
  },
};

// ==================== Cities & Destinations ====================

export const citiesService = {
  getCities: async (page = 1, pageSize = 12) => {
    const response = await api.get('/cities', {
      params: { page, pageSize },
    });
    return response.data;
  },

  getCityById: async (cityId) => {
    const response = await api.get(`/cities/${cityId}`);
    return response.data;
  },
};

// ==================== Trip Posts (Community Feed) ====================

export const tripPostsService = {
  getTripPostsByCity: async (cityId, page = 1, pageSize = 10) => {
    const response = await api.get(`/cities/${cityId}/trip-posts`, {
      params: { page, pageSize },
    });
    return response.data;
  },

  getTripPostById: async (postId) => {
    const response = await api.get(`/trip-posts/${postId}`);
    return response.data;
  },

  addTripPostToItinerary: async (tripPostId) => {
    const response = await api.post('/user-trips/add-from-post', {
      trip_post_id: tripPostId,
    });
    return response.data;
  },
};

// ==================== User Trips & Itineraries ====================

export const userTripsService = {
  getUserTrips: async () => {
    const response = await api.get('/user-trips');
    return response.data;
  },

  getUserTrip: async (tripId) => {
    const response = await api.get(`/user-trips/${tripId}`);
    return response.data;
  },

  createUserTrip: async (tripData) => {
    const response = await api.post('/user-trips', tripData);
    return response.data;
  },

  updateUserTrip: async (tripId, tripData) => {
    const response = await api.put(`/user-trips/${tripId}`, tripData);
    return response.data;
  },

  publishTrip: async (tripId, postData) => {
    const response = await api.post(`/user-trips/${tripId}/publish`, postData);
    return response.data;
  },
};

// ==================== Trip Segments & Reviews ====================

export const tripSegmentsService = {
  markSegmentVisited: async (segmentId) => {
    const response = await api.post(`/trip-segments/${segmentId}/mark-visited`, {
      completed: true,
    });
    return response.data;
  },

  submitReview: async (segmentId, rating, review) => {
    const response = await api.post('/reviews', {
      segment_id: segmentId,
      rating,
      review,
    });
    return response.data;
  },
};

export default api;
