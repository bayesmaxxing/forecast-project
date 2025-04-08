import { useState, useEffect, useCallback } from 'react';
import { forecastService } from '../api/index';

export function useForecastList({list_type, category} = {}) {
  const [forecasts, setForecasts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  // Create a fetchData function using useCallback to memoize it
  const fetchData = useCallback(() => {
    setLoading(true);
    
    return Promise.all([
      forecastService.fetchForecasts(list_type, category),
    ])
    .then(([forecastDataJson]) => {
      setForecasts(forecastDataJson || []);
      setLoading(false);
    })
    .catch(error => {
      console.error('Error fetching data: ', error);
      setError(error);
      setLoading(false);
    });
  }, [list_type, category]);

  // Initial data fetch
  useEffect(() => {
    fetchData();
  }, [fetchData]);

  // Return the refetch function along with the data
  return { 
    forecasts, 
    loading, 
    error, 
    refetchForecasts: fetchData 
  };
}