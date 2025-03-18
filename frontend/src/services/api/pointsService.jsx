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


