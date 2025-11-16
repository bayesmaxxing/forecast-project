import { useState, useEffect, useCallback } from 'react';
import { pointsService } from '../api/index';

export function usePointsData({ id, userId, useOrderedEndpoint = true, useLatestPoints = true } = {}) {
  const [points, setPoints] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  // Create a fetchData function using useCallback
  const fetchData = useCallback(() => {
    setLoading(true);
    
    // Build filter options based on hook parameters
    const options = {};
    
    // Handle user ID (convert 'all' to null)
    const user = userId === 'all' ? null : userId;
    if (user !== null && user !== undefined) {
      options.userId = user;
    }
    
    // Handle forecast ID
    if (id !== null && id !== undefined) {
      options.forecastId = id;
    }
    
    // Handle latest points (distinct on forecast_id with ordering)
    if (useLatestPoints) {
      options.distinct = true;
      options.orderByForecastId = true;
    }
    
    // Use the unified fetchForecastPoints function
    return pointsService.fetchForecastPoints(options)
      .then((pointsDataJson) => {
        setPoints(pointsDataJson);
        setLoading(false);
      })
      .catch(error => {
        console.error('Error fetching data: ', error);
        setError(error);
        setLoading(false);
      });
  }, [id, useOrderedEndpoint, userId, useLatestPoints]);

  // Initial data fetch
  useEffect(() => {
    fetchData();
  }, [fetchData]);

  // Return the refetch function along with the data
  return { 
    points, 
    loading: loading, 
    error, 
    refetchPoints: fetchData 
  };
}