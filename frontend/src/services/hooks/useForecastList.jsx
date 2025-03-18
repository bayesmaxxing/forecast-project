import { useState, useEffect } from 'react';
import { forecastService } from '../api/index';

export function useForecastList({list_type, category} = {}) {
  const [forecasts, setForecasts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    setLoading(true);
    
    Promise.all([
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

  return { forecasts, loading, error };
}