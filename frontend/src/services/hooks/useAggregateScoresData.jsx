import { useState, useEffect } from 'react';
import { fetchAggregateScoresByUserID, 
  fetchAggregateScoresAllUsers, 
  fetchAggregateScoresByUserIDAndCategory, 
  fetchAggregateScoresByCategory,
  fetchAggregateScoresByUsers,
} from '../api/scoreService';

export const useAggregateScoresData = (userId = null) => {
  const [scores, setScores] = useState([]);
  const [scoresLoading, setScoresLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchScores = async () => {
      try {
        setScoresLoading(true);
        // Pass the userId to get scores for a specific user
        const data = userId !== null && userId !== 'all'
          ? await fetchAggregateScoresByUserID(userId)
          : await fetchAggregateScoresAllUsers();
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
  }, [userId]);

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