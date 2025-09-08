import { API_BASE_URL } from './index';

export const fetchForecastNews = async (forecastQuestion) => {
  const response = await fetch(`${API_BASE_URL}/news`, {
    method: 'POST',
    headers: {
      "Accept": "application/json",
      "Content-Type": "application/json"
    },
    body: JSON.stringify({ query: forecastQuestion })
  });
  
  if (!response.ok) {
    throw new Error(`Error fetching forecast news: ${response.status}`);
  }
  
  return response.json();
};