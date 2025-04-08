import { useState, useEffect, useCallback } from 'react';
import { forecastService } from '../api/index';

export function useForecastData({id} = {}) {
  const [forecast, setForecast] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  // Create a fetchData function using useCallback to memoize it
  const fetchData = useCallback(() => {
    setLoading(true);
    
    return Promise.all([
      forecastService.fetchForecastById(id),
    ])
    .then(([forecastDataJson]) => {
      setForecast(forecastDataJson);
      setLoading(false);
    })
    .catch(error => {
      console.error('Error fetching data: ', error);
      setError(error);
      setLoading(false);
    });
  }, [id]);

  // Initial data fetch
  useEffect(() => {
    fetchData();
  }, [fetchData]);

  // Return the refetch function along with the data
  return { 
    forecast, 
    loading, 
    error, 
    refetchForecast: fetchData 
  };
}