import { API_BASE_URL } from './index';

export const fetchForecasts = async (list_type = null, category = null) => {
  // Build query parameters (list_type maps to status in backend)
  const params = new URLSearchParams();
  
  if (list_type !== null && list_type !== undefined) {
    // Map frontend 'list_type' to backend 'status'
    params.append('status', list_type);
  }
  
  if (category !== null && category !== undefined) {
    params.append('category', category);
  }
  
  const queryString = params.toString();
  const url = queryString 
    ? `${API_BASE_URL}/forecasts?${queryString}`
    : `${API_BASE_URL}/forecasts`;
  
  const response = await fetch(url, {
    method: 'GET',
    headers: {
      "Accept": "application/json"
    }
  });
  
  if (!response.ok) {
    throw new Error(`Error fetching forecasts: ${response.status}`);
  }
  
  return response.json();
};

export const fetchForecastById = async (id) => {
    const response = await fetch(`${API_BASE_URL}/forecasts/${id}`, {
      headers: { "Accept": "application/json" }
    });
    
    if (!response.ok) {
      throw new Error(`Error fetching forecast: ${response.status}`);
    }
    
    return response.json();
  };  

export const resolveForecast = async (forecast_id, resolution, comment) => {
  const token = localStorage.getItem('token');
  if (!token) {
    throw new Error('User needs to login to resolve a forecast');
  }
  const response = await fetch(`${API_BASE_URL}/api/resolve`, {
    method: 'PUT',
    headers: { 
      "Accept": "application/json",
      "Authorization": `Bearer ${token}`
    },
    body: JSON.stringify({
      id: forecast_id,
      resolution: resolution,
      comment: comment
    })
  });

  if (!response.ok) {
    throw new Error(`Error resolving forecast: ${response.status}`);
  }

  return true;
};

export const createForecast = async (forecast) => {
  const token = localStorage.getItem('token');
  if (!token) {
    throw new Error('User needs to login to create a forecast');
  }
  const response = await fetch(`${API_BASE_URL}/api/forecasts/create`, {
    method: 'POST',
    headers: { "Accept": "application/json", "Authorization": `Bearer ${token}` },
    body: JSON.stringify(forecast)
  });

  if (!response.ok) {
    throw new Error(`Error creating forecast: ${response.status}`);
  }

  return response.json(); 
};
