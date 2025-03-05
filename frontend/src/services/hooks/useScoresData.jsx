import { useState, useEffect } from 'react';
import { fetchScores, fetchAverageScoresById } from '../api/scoreService';

export const useScoresData = ({ user_id=null, forecast_id=null, shouldFetch=true, useAverageEndpoint=false }) => {
  const [score, setScore] = useState(null);
  const [scoreLoading, setScoreLoading] = useState(shouldFetch);
  const [error, setError] = useState(null);

  useEffect(() => {
    // Skip fetching if shouldFetch is false
    if (!shouldFetch) {
      setScoreLoading(false);
      return;
    }
    
    const fetchScoreData = async () => {
      try {
        setScoreLoading(true);
        
        let data;
        
        // If useAverageEndpoint is true and we have a forecast_id, use the average scores endpoint
        if (useAverageEndpoint && forecast_id) {
          data = await fetchAverageScoresById(forecast_id);
        } else {
          // Otherwise use the regular scores endpoint
          data = await fetchScores(user_id, forecast_id);
        }
        
        setScore(data);
        setError(null);
      } catch (err) {
        console.error('Error fetching scores:', err);
        setError(err.message || 'Failed to load scores');
        setScore(null);
      } finally {
        setScoreLoading(false);
      }
    };

    fetchScoreData();
  }, [user_id, forecast_id, shouldFetch, useAverageEndpoint]);

  return { score, scoreLoading, error };
}; 