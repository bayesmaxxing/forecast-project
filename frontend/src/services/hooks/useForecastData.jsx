import { useState, useEffect } from 'react';
import { forecastService } from '../api/index';

export function useForecastData({id} = {}) {
  const [forecast, setForecast] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    setLoading(true);
    
    // Choose which points API to call based on whether userId is provided
    
    Promise.all([
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

  return { forecast, loading, error };
}