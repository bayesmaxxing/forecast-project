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



export const fetchAggregateScoresAllUsers = async () => {
  
  const response = await fetch(`${API_BASE_URL}/scores/aggregate/all`, {
    method: 'GET',
    headers: { "Accept": "application/json", "Content-Type": "application/json" },
  });
  
  if (!response.ok) {
    throw new Error(`Error fetching aggregate scores: ${response.status}`);
  }

  return response.json();
};

export const fetchAggregateScoresByUserID = async (user_id) => {
  const response = await fetch(`${API_BASE_URL}/scores/aggregate/${user_id}`, {
    method: 'GET',
    headers: { "Accept": "application/json", "Content-Type": "application/json" },
  });
  
  if (!response.ok) {
    throw new Error(`Error fetching aggregate scores: ${response.status}`);
  }
  return response.json();
};

export const fetchAggregateScoresByUserIDAndCategory = async (user_id, category) => {
  const response = await fetch(`${API_BASE_URL}/scores/aggregate/${user_id}/${category}`, {
    method: 'GET',
    headers: { "Accept": "application/json", "Content-Type": "application/json" },
  });
  
  if (!response.ok) {
    throw new Error(`Error fetching aggregate scores: ${response.status}`);
  }
  return response.json();
};

export const fetchAggregateScoresByCategory = async (category) => {
  const response = await fetch(`${API_BASE_URL}/scores/aggregate/category/${category}`, {
    method: 'GET',
    headers: { "Accept": "application/json", "Content-Type": "application/json" },
  });
  if (!response.ok) {
    throw new Error(`Error fetching aggregate scores: ${response.status}`);
  }
  return response.json();
};

export const fetchAggregateScoresByUsers = async () => {
  const response = await fetch(`${API_BASE_URL}/scores/aggregate/users`, {
    method: 'GET',
    headers: { "Accept": "application/json", "Content-Type": "application/json" },
  });
  
  if (!response.ok) {
    throw new Error(`Error fetching aggregate scores: ${response.status}`);
  }
  return response.json();
};
