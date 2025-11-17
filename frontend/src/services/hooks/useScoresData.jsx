import { useState, useEffect, useCallback } from 'react';
import { fetchScores, fetchAverageScoresById, fetchAverageScores } from '../api/scoreService';

export const useScoresData = ({ user_id=null, forecast_id=null, shouldFetch=true, useAverageEndpoint=false }) => {
  const [scores, setScores] = useState([]);
  const [scoresLoading, setScoresLoading] = useState(shouldFetch);
  const [scoresError, setScoresError] = useState(null);

  // Create a fetchData function using useCallback
  const fetchData = useCallback(async () => {
    // Skip fetching if shouldFetch is false
    if (!shouldFetch) {
      setScoresLoading(false);
      return;
    }
    
    try {
      setScoresLoading(true);
      const user = user_id === 'all' ? null : user_id;
      let data;
      
      // If useAverageEndpoint is true and we have a forecast_id, use the average scores endpoint
      if (useAverageEndpoint && forecast_id) {
        data = await fetchAverageScoresById(forecast_id);
      } else if (user_id && forecast_id) {
        data = await fetchScores(user_id, forecast_id);
      } else if (useAverageEndpoint && !forecast_id) {
        data = await fetchAverageScores();
      } else if (user) {
        data = await fetchScores(user, null);
      } else {
        data = [];
      }
      setScores(data || []);
      setScoresError(null);
    } catch (err) {
      console.error('Error fetching scores:', err);
      setScoresError(err.message || 'Failed to load scores');
      setScores([]);
    } finally {
      setScoresLoading(false);
    }
  }, [user_id, forecast_id, shouldFetch, useAverageEndpoint]);

  // Initial data fetch
  useEffect(() => {
    fetchData();
  }, [fetchData]);

  // Return the refetch function along with the data
  return { 
    scores, 
    scoresLoading, 
    scoresError, 
    refetchScores: fetchData 
  };
}; 