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

/**
 * Fetch basic scores filtered by user and/or forecast
 */
export const fetchScores = async (user_id, forecast_id) => {
  const params = new URLSearchParams();
  if (user_id) params.append('user_id', user_id);
  if (forecast_id) params.append('forecast_id', forecast_id);

  const response = await fetch(`${API_BASE_URL}/scores?${params.toString()}`, {
    method: 'GET',
    headers: { 'Accept': 'application/json' },
  });

  if (!response.ok) {
    throw new Error(`Error fetching scores: ${response.status}`);
  }

  return response.json();
};

/**
 * Fetch all average scores
 */
export const fetchAverageScores = async () => {
  const response = await fetch(`${API_BASE_URL}/scores/average`, {
    headers: { 'Accept': 'application/json' },
  });

  if (!response.ok) {
    throw new Error(`Error fetching average scores: ${response.status}`);
  }

  return response.json();
};

/**
 * Fetch aggregate scores with flexible filtering
 * @param {Object} options
 * @param {string|number} options.user_id - Filter by user ID
 * @param {string|number} options.forecast_id - Filter by forecast ID
 * @param {string} options.category - Filter by category
 * @param {string} options.dateRange - Predefined date range (from DATE_RANGE_OPTIONS)
 * @param {Date|string} options.start_date - Custom start date
 * @param {Date|string} options.end_date - Custom end date
 */
export const fetchAggregateScores = async ({
  user_id,
  forecast_id,
  category,
  dateRange,
  start_date,
  end_date,
} = {}) => {
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
    headers: { 'Accept': 'application/json' },
  });

  if (!response.ok) {
    throw new Error(`Error fetching aggregate scores: ${response.status}`);
  }

  return response.json();
};

/**
 * Fetch aggregate scores grouped by users (for leaderboard)
 * @param {Object} options
 * @param {string} options.category - Filter by category
 * @param {string} options.dateRange - Predefined date range (from DATE_RANGE_OPTIONS)
 * @param {Date|string} options.start_date - Custom start date
 * @param {Date|string} options.end_date - Custom end date
 */
export const fetchAggregateScoresByUsers = async ({
  category,
  dateRange,
  start_date,
  end_date,
} = {}) => {
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
    headers: { 'Accept': 'application/json' },
  });

  if (!response.ok) {
    throw new Error(`Error fetching aggregate scores by users: ${response.status}`);
  }

  return response.json();
};
