import { useState, useEffect } from 'react';
import { forecastService, pointsService } from '../api/index';

export function useForecastCombinedData({ userId = null, category = null, list_type = 'open'} = {}) {
  const [combinedForecasts, setCombinedForecasts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    setLoading(true);
    
    // For debugging 
    console.log('Fetching forecasts with params:', { userId, category, list_type });
    
    // Choose which points API to call based on whether userId is provided
    const pointsPromise = userId 
      ? pointsService.fetchLatestPointsByUser(userId)
      : pointsService.fetchLatestPoints();
    
    Promise.all([
      forecastService.fetchForecasts(list_type, category),
      pointsPromise
    ])
    .then(([forecastDataJson, pointsDataJson]) => {
      const combined = forecastDataJson.map(forecast => {
        const matchingPoint = pointsDataJson.find(point => point.forecast_id === forecast.id);
        return { ...forecast, latestPoint: matchingPoint || null};
      });
      setCombinedForecasts(combined);
      setLoading(false);
    })
    .catch(error => {
      console.error('Error fetching data: ', error);
      setError(error);
      setLoading(false);
    });
  }, [userId, category, list_type]); // Added list_type to dependency array

  return { combinedForecasts, loading, error };
}