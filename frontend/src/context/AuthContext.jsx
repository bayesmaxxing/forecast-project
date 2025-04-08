import React, { createContext, useState, useEffect, useContext } from 'react';
import * as authService from '../services/api/authService';

// Create Authentication Context
const AuthContext = createContext();

// Custom hook to use the auth context
export const useAuth = () => useContext(AuthContext);

// Provider component
export const AuthProvider = ({ children }) => {
  const [currentUser, setCurrentUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  // Initialize auth state on app load
  useEffect(() => {
    // Check if user is already logged in
    const user = authService.getCurrentUser();
    if (user && user.token) {
      setCurrentUser(user);
    }
    setLoading(false);
  }, []);

  // Login function
  const login = async (username, password) => {
    setError(null);
    try {
      const userData = await authService.login(username, password);
      setCurrentUser({
        username: userData.username,
        token: userData.token
      });
      return userData;
    } catch (err) {
      setError(err.message || 'Login failed');
      throw err;
    }
  };

  // Register function
  const register = async (username, password) => {
    setError(null);
    try {
      await authService.register(username, password);
      // After registration, automatically log in
      return login(username, password);
    } catch (err) {
      setError(err.message || 'Registration failed');
      throw err;
    }
  };

  // Logout function
  const logout = () => {
    authService.logout();
    setCurrentUser(null);
  };

  // Context value
  const value = {
    currentUser,
    isAuthenticated: !!currentUser,
    login,
    register,
    logout,
    loading,
    error
  };

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
};

export default AuthContext;