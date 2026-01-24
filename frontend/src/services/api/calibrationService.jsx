import { API_BASE_URL } from './index';
import { DATE_RANGE_OPTIONS, getStartDateForRange } from './scoreService';

/**
 * Fetch overall calibration data
 * @param {Object} options
 * @param {string|number} options.user_id - Filter by user ID
 * @param {string} options.category - Filter by category
 * @param {string} options.dateRange - Predefined date range (from DATE_RANGE_OPTIONS)
 * @param {Date|string} options.start_date - Custom start date
 * @param {Date|string} options.end_date - Custom end date
 */
export const fetchCalibrationData = async ({
  user_id,
  category,
  dateRange,
  start_date,
  end_date,
} = {}) => {
  const params = new URLSearchParams();
  if (user_id && user_id !== 'all') params.append('user_id', user_id);
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

  const response = await fetch(`${API_BASE_URL}/calibration?${params.toString()}`, {
    method: 'GET',
    headers: { 'Accept': 'application/json' },
  });

  if (!response.ok) {
    throw new Error(`Error fetching calibration data: ${response.status}`);
  }

  return response.json();
};

/**
 * Fetch calibration data grouped by users
 * @param {Object} options
 * @param {string} options.category - Filter by category
 * @param {string} options.dateRange - Predefined date range (from DATE_RANGE_OPTIONS)
 * @param {Date|string} options.start_date - Custom start date
 * @param {Date|string} options.end_date - Custom end date
 */
export const fetchCalibrationDataByUsers = async ({
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

  const response = await fetch(`${API_BASE_URL}/calibration/users?${params.toString()}`, {
    method: 'GET',
    headers: { 'Accept': 'application/json' },
  });

  if (!response.ok) {
    throw new Error(`Error fetching calibration data by users: ${response.status}`);
  }

  return response.json();
};
