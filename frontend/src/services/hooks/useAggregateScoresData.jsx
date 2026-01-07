import { useState, useEffect } from 'react';
import { fetchAggregateScores,
  fetchAggregateScoresByUsers,
} from '../api/scoreService';

export const useAggregateScoresData = (userId = null, dateRange = null) => {
  const [scores, setScores] = useState([]);
  const [scoresLoading, setScoresLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchScoresData = async () => {
      try {
        setScoresLoading(true);
        // Pass the userId and dateRange to get scores
        const params = { dateRange };
        if (userId !== null && userId !== 'all') {
          params.user_id = userId;
        }
        const data = await fetchAggregateScores(params);
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

    fetchScoresData();
  }, [userId, dateRange]);

  return { scores, scoresLoading, error };
}; 

export const useAggregateScoresDataByCategory = (userId = null, category = null) => {
  const [scores, setScores] = useState([]);
  const [scoresLoading, setScoresLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchScores = async () => {
      try {
        setScoresLoading(true);
        // Pass the userId to get scores for a specific user
        const data = userId !== null && userId !== 'all'
          ? await fetchAggregateScoresByUserIDAndCategory(userId, category)
          : await fetchAggregateScoresByCategory(category);
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
  }, [userId, category]);

  return { scores, scoresLoading, error };
};

export const useAggregateScoresDataByUsers = () => {
  const [scores, setScores] = useState([]);
  const [scoresLoading, setScoresLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchScores = async () => {
      try {
        setScoresLoading(true);
        const data = await fetchAggregateScoresByUsers();
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
  }, []);

  return { scores, scoresLoading, error };
};