import { API_BASE_URL } from './index';

// Predefined date range options for aggregate scores
export const DATE_RANGE_OPTIONS = {
  ALL_TIME: 'all_time',
  LAST_12_MONTHS: 'last_12_months',
  LAST_6_MONTHS: 'last_6_months',
  LAST_3_MONTHS: 'last_3_months',
};

// Helper to get start date for predefined ranges
export const getStartDateForRange = (rangeOption) => {
  const now = new Date();
  switch (rangeOption) {
    case DATE_RANGE_OPTIONS.LAST_12_MONTHS:
      return new Date(now.getFullYear() - 1, now.getMonth(), now.getDate());
    case DATE_RANGE_OPTIONS.LAST_6_MONTHS:
      return new Date(now.getFullYear(), now.getMonth() - 6, now.getDate());
    case DATE_RANGE_OPTIONS.LAST_3_MONTHS:
      return new Date(now.getFullYear(), now.getMonth() - 3, now.getDate());
    case DATE_RANGE_OPTIONS.ALL_TIME:
    default:
      return null;
  }
};
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
// Accepts optional date range: either a predefined range option or custom start_date/end_date
export const fetchAggregateScores = async ({ user_id, forecast_id, category, dateRange, start_date, end_date } = {}) => {
  const params = new URLSearchParams();
  if (user_id) params.append('user_id', user_id);
  if (forecast_id) params.append('forecast_id', forecast_id);
  if (category) params.append('category', category);

  // Handle date range - either predefined option or custom dates
  if (dateRange && dateRange !== DATE_RANGE_OPTIONS.ALL_TIME) {
    const startDate = getStartDateForRange(dateRange);
    if (startDate) {
      params.append('start_date', startDate.toISOString());
    }
  } else if (start_date) {
    params.append('start_date', start_date instanceof Date ? start_date.toISOString() : start_date);
    if (end_date) {
      params.append('end_date', end_date instanceof Date ? end_date.toISOString() : end_date);
    }
  }

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

export const fetchAggregateScoresByUsers = async ({ category, dateRange, start_date, end_date } = {}) => {
  const params = new URLSearchParams();
  if (category) params.append('category', category);

  // Handle date range - either predefined option or custom dates
  if (dateRange && dateRange !== DATE_RANGE_OPTIONS.ALL_TIME) {
    const startDate = getStartDateForRange(dateRange);
    if (startDate) {
      params.append('start_date', startDate.toISOString());
    }
  } else if (start_date) {
    params.append('start_date', start_date instanceof Date ? start_date.toISOString() : start_date);
    if (end_date) {
      params.append('end_date', end_date instanceof Date ? end_date.toISOString() : end_date);
    }
  }

  const response = await fetch(`${API_BASE_URL}/scores/aggregate/users?${params.toString()}`, {
    method: 'GET',
    headers: { "Accept": "application/json" },
  });

  if (!response.ok) {
    throw new Error(`Error fetching aggregate scores by users: ${response.status}`);
  }
  return response.json();
};
