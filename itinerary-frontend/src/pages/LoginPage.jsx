import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';
import { LogIn, Plane, Check } from 'lucide-react';

const LoginPage = () => {
  const navigate = useNavigate();
  const { login, isLoading } = useAuth();
  const [formData, setFormData] = useState({
    email: '',
    password: '',
  });
  const [error, setError] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [rememberMe, setRememberMe] = useState(false);
  const [emailSuggestions, setEmailSuggestions] = useState([]);
  const [showSuggestions, setShowSuggestions] = useState(false);
  const [successMessage, setSuccessMessage] = useState('');

  // Load saved emails and last login data on mount
  useEffect(() => {
    const savedEmails = JSON.parse(localStorage.getItem('savedEmails') || '[]');
    setEmailSuggestions(savedEmails);

    // Auto-load last login if remember me was checked
    const lastLogin = localStorage.getItem('lastLogin');
    if (lastLogin) {
      try {
        const { email, password, rememberMe: wasRemembered } = JSON.parse(lastLogin);
        if (wasRemembered) {
          setFormData({ email, password });
          setRememberMe(true);
        }
      } catch (e) {
        console.log('Could not load saved login');
      }
    }
  }, []);

  // Filter email suggestions based on input
  const handleEmailChange = (e) => {
    const value = e.target.value;
    setFormData((prev) => ({
      ...prev,
      email: value,
    }));
    setError('');
    setSuccessMessage('');

    // Show suggestions if input is not empty
    if (value.length > 0) {
      const filtered = emailSuggestions.filter((email) =>
        email.toLowerCase().includes(value.toLowerCase())
      );
      setShowSuggestions(filtered.length > 0);
    } else {
      setShowSuggestions(false);
    }
  };

  const handlePasswordChange = (e) => {
    const { name, value } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }));
    setError('');
    setSuccessMessage('');
  };

  // Handle email suggestion click
  const handleSuggestionClick = (email) => {
    setFormData((prev) => ({
      ...prev,
      email,
    }));
    setShowSuggestions(false);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setSuccessMessage('');

    // Validation
    if (!formData.email) {
      setError('Email is required');
      return;
    }
    if (!formData.password) {
      setError('Password is required');
      return;
    }

    try {
      const result = await login(formData.email, formData.password);
      if (result.success) {
        // Save email to history
        const savedEmails = JSON.parse(localStorage.getItem('savedEmails') || '[]');
        if (!savedEmails.includes(formData.email)) {
          savedEmails.unshift(formData.email);
          // Keep only last 10 emails
          if (savedEmails.length > 10) {
            savedEmails.pop();
          }
          localStorage.setItem('savedEmails', JSON.stringify(savedEmails));
        }

        // Save login credentials if remember me is checked
        if (rememberMe) {
          localStorage.setItem(
            'lastLogin',
            JSON.stringify({
              email: formData.email,
              password: formData.password,
              rememberMe: true,
            })
          );
        } else {
          localStorage.removeItem('lastLogin');
        }

        setSuccessMessage('Login successful! Redirecting...');
        setTimeout(() => {
          navigate('/dashboard');
        }, 1000);
      } else {
        setError(result.error || 'Login failed');
      }
    } catch (err) {
      setError('An error occurred during login');
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-600 via-blue-500 to-purple-600 flex items-center justify-center p-4">
      {/* Background Pattern */}
      <div className="absolute inset-0 overflow-hidden">
        <div className="absolute top-20 left-10 w-72 h-72 bg-white rounded-full mix-blend-multiply filter blur-3xl opacity-10 animate-pulse"></div>
        <div className="absolute top-40 right-10 w-72 h-72 bg-purple-300 rounded-full mix-blend-multiply filter blur-3xl opacity-10 animate-pulse" style={{ animationDelay: '2s' }}></div>
        <div className="absolute bottom-20 left-40 w-72 h-72 bg-blue-300 rounded-full mix-blend-multiply filter blur-3xl opacity-10 animate-pulse" style={{ animationDelay: '4s' }}></div>
      </div>

      <div className="relative z-10 w-full max-w-md">
        {/* Logo & Title */}
        <div className="text-center mb-8">
          <div className="flex items-center justify-center gap-3 mb-4">
            <Plane className="w-10 h-10 text-white" />
            <h1 className="text-4xl font-bold text-white">Itinerary</h1>
          </div>
          <p className="text-blue-100 text-lg">Plan Your Perfect Trip</p>
        </div>

        {/* Login Card */}
        <div className="bg-white rounded-2xl shadow-2xl p-8 backdrop-blur-lg bg-opacity-95">
          <h2 className="text-2xl font-bold text-gray-800 mb-2">Welcome Back</h2>
          <p className="text-gray-600 mb-6">Sign in to your account</p>

          {/* Error Message */}
          {error && (
            <div className="mb-4 p-4 bg-red-50 border border-red-200 rounded-lg">
              <p className="text-red-700 text-sm font-medium">{error}</p>
            </div>
          )}

          {/* Success Message */}
          {successMessage && (
            <div className="mb-4 p-4 bg-green-50 border border-green-200 rounded-lg flex items-center gap-2">
              <Check className="w-5 h-5 text-green-600" />
              <p className="text-green-700 text-sm font-medium">{successMessage}</p>
            </div>
          )}

          {/* Login Form */}
          <form onSubmit={handleSubmit} className="space-y-4">
            {/* Email Input with Suggestions */}
            <div>
              <label htmlFor="email" className="block text-sm font-medium text-gray-700 mb-2">
                Email Address
              </label>
              <div className="relative">
                <input
                  type="email"
                  id="email"
                  name="email"
                  value={formData.email}
                  onChange={handleEmailChange}
                  onFocus={() => formData.email.length > 0 && setShowSuggestions(true)}
                  placeholder="you@example.com"
                  autoComplete="email"
                  className="w-full px-4 py-3 rounded-lg border border-gray-300 focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none transition"
                  disabled={isLoading}
                />

                {/* Email Suggestions Dropdown */}
                {showSuggestions && emailSuggestions.length > 0 && (
                  <div className="absolute top-full left-0 right-0 mt-1 bg-white border border-gray-300 rounded-lg shadow-lg z-10">
                    {emailSuggestions
                      .filter((email) =>
                        email.toLowerCase().includes(formData.email.toLowerCase())
                      )
                      .map((email, index) => (
                        <button
                          key={index}
                          type="button"
                          onClick={() => handleSuggestionClick(email)}
                          className="w-full text-left px-4 py-2 hover:bg-blue-50 border-b border-gray-100 last:border-b-0 text-gray-700 flex items-center gap-2 transition"
                        >
                          <span className="text-gray-400">📧</span>
                          <span className="text-sm">{email}</span>
                        </button>
                      ))}
                  </div>
                )}
              </div>
            </div>

            {/* Password Input */}
            <div>
              <label htmlFor="password" className="block text-sm font-medium text-gray-700 mb-2">
                Password
              </label>
              <div className="relative">
                <input
                  type={showPassword ? 'text' : 'password'}
                  id="password"
                  name="password"
                  value={formData.password}
                  onChange={handlePasswordChange}
                  placeholder="••••••••"
                  autoComplete={rememberMe ? 'on' : 'off'}
                  className="w-full px-4 py-3 rounded-lg border border-gray-300 focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none transition"
                  disabled={isLoading}
                />
                <button
                  type="button"
                  onClick={() => setShowPassword(!showPassword)}
                  className="absolute right-3 top-3 text-gray-500 hover:text-gray-700"
                >
                  {showPassword ? '🙈' : '👁️'}
                </button>
              </div>
            </div>

            {/* Remember Me & Forgot Password */}
            <div className="flex items-center justify-between">
              <label className="flex items-center gap-2 cursor-pointer">
                <input
                  type="checkbox"
                  checked={rememberMe}
                  onChange={(e) => setRememberMe(e.target.checked)}
                  className="w-4 h-4 rounded border-gray-300 text-blue-600 focus:ring-2 focus:ring-blue-500"
                  disabled={isLoading}
                />
                <span className="text-sm text-gray-700 font-medium">Remember me</span>
              </label>
              <a href="#" className="text-sm text-blue-600 hover:text-blue-700 font-medium">
                Forgot password?
              </a>
            </div>

            {/* Submit Button */}
            <button
              type="submit"
              disabled={isLoading}
              className="w-full bg-gradient-to-r from-blue-600 to-blue-700 hover:from-blue-700 hover:to-blue-800 disabled:opacity-50 disabled:cursor-not-allowed text-white font-bold py-3 rounded-lg flex items-center justify-center gap-2 transition"
            >
              <LogIn className="w-5 h-5" />
              {isLoading ? 'Signing in...' : 'Sign In'}
            </button>
          </form>

          {/* Demo Credentials Info */}
          <div className="mt-6 p-4 bg-blue-50 rounded-lg border border-blue-200">
            <p className="text-sm text-blue-800 font-medium mb-2">Demo Credentials:</p>
            <p className="text-sm text-blue-700">Email: demo@example.com</p>
            <p className="text-sm text-blue-700">Password: demo123456</p>
          </div>

          {/* Footer */}
          <p className="text-center text-gray-600 text-sm mt-6">
            Don't have an account?{' '}
            <span className="text-gray-500 text-xs">
              Contact administrator for access
            </span>
          </p>
        </div>

        {/* Features */}
        <div className="mt-8 grid grid-cols-3 gap-4">
          <div className="text-center text-white">
            <div className="text-3xl mb-2">🌍</div>
            <p className="text-sm font-medium">Explore</p>
          </div>
          <div className="text-center text-white">
            <div className="text-3xl mb-2">📋</div>
            <p className="text-sm font-medium">Plan</p>
          </div>
          <div className="text-center text-white">
            <div className="text-3xl mb-2">⭐</div>
            <p className="text-sm font-medium">Share</p>
          </div>
        </div>
      </div>
    </div>
  );
};

export default LoginPage;
