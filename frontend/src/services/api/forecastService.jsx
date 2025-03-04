import { API_BASE_URL } from './index';

export const fetchForecasts = async (list_type = null, category = null) => {
  // Build request body with only non-null parameters
  const requestBody = {};
  if (list_type) requestBody.list_type = list_type;
  if (category) requestBody.category = category;
  
  const response = await fetch(`${API_BASE_URL}/forecasts`, {
    method: 'POST',
    headers: {
      "Accept": "application/json",
      "Content-Type": "application/json"
    },
    body: JSON.stringify(requestBody)
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

