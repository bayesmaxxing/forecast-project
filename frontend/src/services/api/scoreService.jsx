import { API_BASE_URL } from './index';
export const fetchScores = async (user_id, forecast_id) => {
    const requestBody = {};
    if (user_id) requestBody.user_id = user_id;
    if (forecast_id) requestBody.forecast_id = forecast_id;
    
    const response = await fetch(`${API_BASE_URL}/scores`, {
      method: 'POST',
      headers: {
        "Accept": "application/json",
        "Content-Type": "application/json"
      },
      body: JSON.stringify(requestBody)
    });
    
    if (!response.ok) {
        throw new Error(`Error fetching scores: ${response.status}`);
    }
  
    return response.json();
};

export const fetchAllScores = async () => {
  const response = await fetch(`${API_BASE_URL}/scores/all`, {
    headers: { "Accept": "application/json" }
  });
  
  if (!response.ok) {
    throw new Error(`Error fetching all scores: ${response.status}`);
  }

  return response.json();
};

export const fetchAggregateScores = async (category,user_id,by_user) => {
  const requestBody = {};
  if (category) requestBody.category = category;
  if (user_id) requestBody.user_id = user_id;
  if (by_user) requestBody.by_user = by_user;

  const response = await fetch(`${API_BASE_URL}/scores/aggregate`, {
    method: 'POST',
    headers: { "Accept": "application/json", "Content-Type": "application/json" },
    body: JSON.stringify(requestBody)
  });
  
  if (!response.ok) {
    throw new Error(`Error fetching aggregate scores: ${response.status}`);
  }

  return response.json();
};
