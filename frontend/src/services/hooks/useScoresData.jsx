import { useState, useEffect } from 'react';
import { fetchScores } from '../api/scoreService';

export const useScoresData = ({ user_id=null, forecast_id=null, shouldFetch=true }) => {
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
        // Pass the userId to get scores for a specific user
        const data = await fetchScores(user_id, forecast_id);
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
  }, [user_id, forecast_id, shouldFetch]);

  return { score, scoreLoading, error };
}; 