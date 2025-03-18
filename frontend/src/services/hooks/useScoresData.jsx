import { useState, useEffect } from 'react';
import { fetchScores, fetchAverageScoresById, fetchAverageScores } from '../api/scoreService';

export const useScoresData = ({ user_id=null, forecast_id=null, shouldFetch=true, useAverageEndpoint=false }) => {
  const [scores, setScores] = useState([]);
  const [scoresLoading, setScoresLoading] = useState(shouldFetch);
  const [scoresError, setScoresError] = useState(null);

  useEffect(() => {
    // Skip fetching if shouldFetch is false
    if (!shouldFetch) {
      setScoresLoading(false);
      return;
    }
    
    const fetchScoreData = async () => {
      try {
        setScoresLoading(true);
        
        let data;
        
        // If useAverageEndpoint is true and we have a forecast_id, use the average scores endpoint
        if (useAverageEndpoint && forecast_id) {
          data = await fetchAverageScoresById(forecast_id);
        } else if (user_id && forecast_id) {
          data = await fetchScores(user_id, forecast_id);
        } else if (useAverageEndpoint && !forecast_id) {
          data = await fetchAverageScores();
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
    };

    fetchScoreData();
  }, [user_id, forecast_id, shouldFetch, useAverageEndpoint]);

  return { scores, scoresLoading, scoresError };
}; 