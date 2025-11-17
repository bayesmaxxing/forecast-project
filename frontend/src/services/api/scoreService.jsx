import { API_BASE_URL } from './index';
export const fetchScores = async (user_id, forecast_id) => {
    const params = new URLSearchParams();
    if (user_id) params.append('user_id', user_id);
    if (forecast_id) params.append('forecast_id', forecast_id);
    
    const response = await fetch(`${API_BASE_URL}/scores?${params.toString()}`, {
      method: 'GET',
      headers: {
        "Accept": "application/json"
      }
    });
    
    if (!response.ok) {
        throw new Error(`Error fetching scores: ${response.status}`);
    }
  
    return response.json();
};

export const fetchAverageScoresById = async (id) => {
  const response = await fetch(`${API_BASE_URL}/scores/average/${id}`, {
    headers: { "Accept": "application/json" }
  });
  
  if (!response.ok) {
    throw new Error(`Error fetching average scores: ${response.status}`);
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

export const fetchAverageScores = async () => {
  const response = await fetch(`${API_BASE_URL}/scores/average`, {
    headers: { "Accept": "application/json" }
  });
  
  if (!response.ok) {
    throw new Error(`Error fetching average scores: ${response.status}`);
  }

  return response.json();
};



// New unified aggregate scores endpoint with flexible query parameters
export const fetchAggregateScores = async ({ user_id, forecast_id, category } = {}) => {
  const params = new URLSearchParams();
  if (user_id) params.append('user_id', user_id);
  if (forecast_id) params.append('forecast_id', forecast_id);
  if (category) params.append('category', category);
  
  const response = await fetch(`${API_BASE_URL}/scores/aggregate?${params.toString()}`, {
    method: 'GET',
    headers: { "Accept": "application/json" },
  });
  
  if (!response.ok) {
    throw new Error(`Error fetching aggregate scores: ${response.status}`);
  }

  return response.json();
};

// Legacy function names for backward compatibility
export const fetchAggregateScoresAllUsers = async () => {
  return fetchAggregateScores();
};

export const fetchAggregateScoresByUserID = async (user_id) => {
  return fetchAggregateScores({ user_id });
};

export const fetchAggregateScoresByUserIDAndCategory = async (user_id, category) => {
  return fetchAggregateScores({ user_id, category });
};

export const fetchAggregateScoresByCategory = async (category) => {
  return fetchAggregateScores({ category });
};

export const fetchAggregateScoresByUsers = async (category = null) => {
  const params = new URLSearchParams();
  if (category) params.append('category', category);
  
  const response = await fetch(`${API_BASE_URL}/scores/aggregate/users?${params.toString()}`, {
    method: 'GET',
    headers: { "Accept": "application/json" },
  });
  
  if (!response.ok) {
    throw new Error(`Error fetching aggregate scores by users: ${response.status}`);
  }
  return response.json();
};
