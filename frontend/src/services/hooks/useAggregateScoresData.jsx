import { useState, useEffect } from 'react';
import { fetchAggregateScores } from '../api/scoreService';

export const useAggregateScoresData = (category = null,userId = null, byUser = null) => {
  const [scores, setScores] = useState([]);
  const [scoresLoading, setScoresLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchScores = async () => {
      try {
        setScoresLoading(true);
        // Pass the userId to get scores for a specific user
        const data = userId !== 'all' 
          ? await fetchAggregateScores(category,userId,byUser)
          : await fetchAggregateScores(category,null,byUser);
        setScores(data);
        setError(null);
      } catch (err) {
        console.error('Error fetching scores:', err);
        setError(err.message || 'Failed to load scores');
        setScores(null);
      } finally {
        setScoresLoading(false);
      }
    };

    fetchScores();
  }, [category,userId,byUser]);

  return { scores, scoresLoading, error };
}; 