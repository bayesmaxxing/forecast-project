import { useState, useEffect, useCallback } from 'react';
import { fetchAggregateScoresByUsers } from '../api/scoreService';

export const useAggregateScoresByUsers = (dateRange = null) => {
  const [scores, setScores] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const fetchData = useCallback(async () => {
    try {
      setLoading(true);
      const data = await fetchAggregateScoresByUsers({ dateRange });
      setScores(data || []);
      setError(null);
    } catch (err) {
      console.error('Error fetching aggregate scores by users:', err);
      setError(err.message || 'Failed to load scores');
      setScores([]);
    } finally {
      setLoading(false);
    }
  }, [dateRange]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  return {
    scores,
    loading,
    error,
    refetch: fetchData
  };
};
