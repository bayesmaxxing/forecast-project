import { useState, useEffect, useCallback } from 'react';
import {
  fetchCalibrationData,
  fetchCalibrationDataByUsers,
} from '../api/calibrationService';

/**
 * Hook for fetching calibration data.
 *
 * @param {Object} options
 * @param {'overall' | 'by-users'} options.type - Type of calibration data to fetch
 * @param {string|number|null} options.userId - User ID filter (use 'all' or null for all users)
 * @param {string|null} options.category - Category filter
 * @param {string|null} options.dateRange - Date range filter
 * @param {boolean} options.shouldFetch - Whether to fetch data (default: true)
 *
 * Usage examples:
 *
 * // Overall calibration data
 * useCalibration({ type: 'overall' })
 *
 * // Calibration for a specific user
 * useCalibration({ type: 'overall', userId: 1 })
 *
 * // Calibration by users (for comparison)
 * useCalibration({ type: 'by-users' })
 */
export const useCalibration = ({
  type = 'overall',
  userId = null,
  category = null,
  dateRange = null,
  shouldFetch = true,
} = {}) => {
  const [calibration, setCalibration] = useState(null);
  const [loading, setLoading] = useState(shouldFetch);
  const [error, setError] = useState(null);

  // Normalize userId - treat 'all' as null
  const normalizedUserId = userId === 'all' ? null : userId;

  const fetchData = useCallback(async () => {
    if (!shouldFetch) {
      setLoading(false);
      return;
    }

    try {
      setLoading(true);
      setError(null);
      let data;

      switch (type) {
        case 'overall':
          data = await fetchCalibrationData({
            user_id: normalizedUserId,
            category,
            dateRange,
          });
          break;

        case 'by-users':
          data = await fetchCalibrationDataByUsers({
            category,
            dateRange,
          });
          break;

        default:
          console.warn(`Unknown calibration type: ${type}`);
          data = null;
      }

      setCalibration(data);
    } catch (err) {
      console.error(`Error fetching ${type} calibration:`, err);
      setError(err.message || 'Failed to load calibration data');
      setCalibration(null);
    } finally {
      setLoading(false);
    }
  }, [type, normalizedUserId, category, dateRange, shouldFetch]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  return {
    calibration,
    loading,
    error,
    refetch: fetchData,
  };
};
