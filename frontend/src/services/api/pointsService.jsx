import { API_BASE_URL } from './index';
export const fetchPointsByID = async (id) => {
  const response = await fetch(`${API_BASE_URL}/forecast-points/${id}`, {
    headers: { "Accept": "application/json" }
  });
  
  return response.json();
};

export const fetchOrderedPointsByID = async (id) => {
  const response = await fetch(`${API_BASE_URL}/forecast-points/ordered/${id}`, {
    headers: { "Accept": "application/json" }
  });
  
  if (!response.ok) {
    throw new Error(`Error fetching ordered points: ${response.status}`);
  }
  
  return response.json();
};

export const fetchAllPoints = async () => {
  const response = await fetch(`${API_BASE_URL}/forecast-points`, {
    headers: { "Accept": "application/json" }
  });
  
  return response.json();
};

export const fetchLatestPoints = async () => {
  const response = await fetch(`${API_BASE_URL}/forecast-points/latest`, {
    headers: { "Accept": "application/json" }
  });
  
  if (!response.ok) {
    throw new Error(`Error fetching points: ${response.status}`);
  }
  
  return response.json();
};

export const fetchLatestPointsByUser = async (user_id) => {
  const response = await fetch(`${API_BASE_URL}/forecast-points/latest_by_user?user_id=${user_id}`, {
    headers: { "Accept": "application/json" }
  });
  
  if (!response.ok) {
    throw new Error(`Error fetching points: ${response.status}`);
  }
  
  return response.json();
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

