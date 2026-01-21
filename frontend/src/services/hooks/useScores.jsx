import { useState, useEffect, useCallback } from 'react';
import {
  fetchScores,
  fetchAverageScores,
  fetchAverageScoresById,
  fetchAggregateScores,
  fetchAggregateScoresByUsers,
} from '../api/scoreService';

/**
 * Consolidated hook for fetching all types of scores.
 *
 * @param {Object} options
 * @param {'basic' | 'average' | 'aggregate' | 'aggregate-by-users'} options.type - Type of scores to fetch
 * @param {string|number|null} options.userId - User ID filter (use 'all' or null for all users)
 * @param {string|number|null} options.forecastId - Forecast ID filter
 * @param {string|null} options.category - Category filter (for aggregate type)
 * @param {string|null} options.dateRange - Date range filter (for aggregate types)
 * @param {boolean} options.shouldFetch - Whether to fetch data (default: true)
 *
 * Usage examples:
 *
 * // Basic scores for a specific user and forecast
 * useScores({ type: 'basic', userId: 1, forecastId: 5 })
 *
 * // Average scores for a specific forecast (used in multi-user mode)
 * useScores({ type: 'average', forecastId: 5 })
 *
 * // Aggregate scores for dashboard
 * useScores({ type: 'aggregate', userId: 1, dateRange: 'last_3_months' })
 *
 * // Aggregate scores by category
 * useScores({ type: 'aggregate', userId: 1, category: 'tech' })
 *
 * // Leaderboard data - all users' aggregate scores
 * useScores({ type: 'aggregate-by-users', dateRange: 'last_12_months' })
 */
export const useScores = ({
  type = 'aggregate',
  userId = null,
  forecastId = null,
  category = null,
  dateRange = null,
  shouldFetch = true,
} = {}) => {
  const [scores, setScores] = useState(type === 'aggregate-by-users' ? [] : null);
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
        case 'basic':
          // Basic scores endpoint - requires at least userId or forecastId
          if (normalizedUserId || forecastId) {
            data = await fetchScores(normalizedUserId, forecastId);
          } else {
            data = [];
          }
          break;

        case 'average':
          // Average scores - optionally filtered by forecast ID
          if (forecastId) {
            data = await fetchAverageScoresById(forecastId);
          } else {
            data = await fetchAverageScores();
          }
          break;

        case 'aggregate':
          // Aggregate scores with flexible filtering
          data = await fetchAggregateScores({
            user_id: normalizedUserId,
            forecast_id: forecastId,
            category,
            dateRange,
          });
          break;

        case 'aggregate-by-users':
          // Aggregate scores grouped by users (for leaderboard)
          data = await fetchAggregateScoresByUsers({
            category,
            dateRange,
          });
          break;

        default:
          console.warn(`Unknown score type: ${type}`);
          data = null;
      }

      setScores(data ?? (type === 'aggregate-by-users' ? [] : null));
    } catch (err) {
      console.error(`Error fetching ${type} scores:`, err);
      setError(err.message || 'Failed to load scores');
      setScores(type === 'aggregate-by-users' ? [] : null);
    } finally {
      setLoading(false);
    }
  }, [type, normalizedUserId, forecastId, category, dateRange, shouldFetch]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  return {
    scores,
    loading,
    error,
    refetch: fetchData,
  };
};
