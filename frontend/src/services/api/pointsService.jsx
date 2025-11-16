import { API_BASE_URL } from './index';

// Generic function to fetch forecast points with optional filters
export const fetchForecastPoints = async (options = {}) => {
  const { userId, forecastId, date, distinct, orderByForecastId, createdDirection } = options;
  
  // Build query parameters
  const params = new URLSearchParams();
  
  if (userId !== undefined && userId !== null) {
    params.append('user_id', userId);
  }
  if (forecastId !== undefined && forecastId !== null) {
    params.append('forecast_id', forecastId);
  }
  if (date !== undefined && date !== null) {
    params.append('date', date);
  }
  if (distinct !== undefined && distinct !== null) {
    params.append('distinct', distinct);
  }
  if (orderByForecastId !== undefined && orderByForecastId !== null) {
    params.append('order_by_forecast_id', orderByForecastId);
  }
  if (createdDirection !== undefined && createdDirection !== null) {
    params.append('created_direction', createdDirection);
  }
  
  const queryString = params.toString();
  const url = queryString 
    ? `${API_BASE_URL}/forecast-points?${queryString}`
    : `${API_BASE_URL}/forecast-points`;
  
  const response = await fetch(url, {
    headers: { "Accept": "application/json" }
  });
  
  if (!response.ok) {
    throw new Error(`Error fetching forecast points: ${response.status}`);
  }
  
  return response.json();
};

// Backward compatibility wrappers
export const fetchPointsByID = async (id) => {
  return fetchForecastPoints({ forecastId: id });
};

export const fetchOrderedPointsByID = async (id) => {
  return fetchForecastPoints({ forecastId: id });
};

export const fetchAllPoints = async () => {
  return fetchForecastPoints();
};

export const fetchLatestPoints = async () => {
  return fetchForecastPoints({ distinct: true, orderByForecastId: true });
};

export const fetchLatestPointsByUser = async (user_id) => {
  return fetchForecastPoints({ userId: user_id, distinct: true, orderByForecastId: true });
};

export const createPoint = async (forecast_id, point_forecast, reason) => {
  // Get the token from localStorage
  const token = localStorage.getItem('token');
  if (!token) {
    throw new Error('User needs to login to create a forecast point');
  }
  
  const response = await fetch(`${API_BASE_URL}/api/forecast-points`, {
    method: 'POST',
    headers: { 
      "Accept": "application/json",
      "Content-Type": "application/json",
      "Authorization": `Bearer ${token}`
    },
    body: JSON.stringify({
      forecast_id: forecast_id,
      point_forecast: point_forecast,
      reason: reason,
      user_id: 0,
    })
  });

  if (!response.ok) {
    throw new Error(`Error creating point: ${response.status}`);
  }

  return response.json(); 
};

