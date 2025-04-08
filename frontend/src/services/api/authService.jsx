// Authentication service for JWT token management
const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

// Login user and get JWT token
export const login = async (username, password) => {
  try {
    const response = await fetch(`${API_BASE_URL}/users/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ username, password }),
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => null);
      throw new Error(errorData?.message || 'Login failed');
    }

    const data = await response.json();
    // Store token and user info in local storage
    localStorage.setItem('token', data.token);
    localStorage.setItem('username', data.username);
    
    return data;
  } catch (error) {
    console.error('Login error:', error);
    throw error;
  }
};

// Register new user
export const register = async (username, password) => {
  try {
    const response = await fetch(`${API_BASE_URL}/users`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ username, password }),
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => null);
      throw new Error(errorData?.message || 'Registration failed');
    }

    return await response.json();
  } catch (error) {
    console.error('Registration error:', error);
    throw error;
  }
};

// Logout user by removing token
export const logout = () => {
  localStorage.removeItem('token');
  localStorage.removeItem('username');
};

// Get current user info from local storage
export const getCurrentUser = () => {
  return {
    username: localStorage.getItem('username'),
    token: localStorage.getItem('token'),
  };
};

// Check if user is logged in
export const isLoggedIn = () => {
  const token = localStorage.getItem('token');
  return !!token;
};

// Add token to requests
export const authHeader = () => {
  const token = localStorage.getItem('token');
  
  if (token) {
    return { 'Authorization': `Bearer ${token}` };
  } else {
    return {};
  }
};

// Handle 401 Unauthorized errors (token expired)
export const handleUnauthorized = (error) => {
  if (error.response && error.response.status === 401) {
    logout();
    window.location.href = '/login';
  }
  return Promise.reject(error);
};