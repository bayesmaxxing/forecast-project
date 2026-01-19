import { useState, useEffect, useCallback } from 'react';
import { pointsService } from '../api/index';

export function useLatestForecastPoints({ userId = 'all', limit = 10 } = {}) {
  const [points, setPoints] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const fetchData = useCallback(() => {
    setLoading(true);

    // Build filter options
    const options = {
      createdDirection: 'DESC'
    };

    // Handle user ID (convert 'all' to null)
    const user = userId === 'all' ? null : userId;
    if (user !== null && user !== undefined) {
      options.userId = user;
    }

    return pointsService.fetchForecastPoints(options)
      .then((pointsData) => {
        // Sort by created date descending to ensure latest points are first
        const sortedPoints = [...pointsData].sort((a, b) =>
          new Date(b.created) - new Date(a.created)
        );
        // Slice to limit results since backend doesn't support limit param
        setPoints(sortedPoints.slice(0, limit));
        setLoading(false);
      })
      .catch(error => {
        console.error('Error fetching latest forecast points: ', error);
        setError(error);
        setLoading(false);
      });
  }, [userId, limit]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  return {
    points,
    loading,
    error,
    refetchPoints: fetchData
  };
}
