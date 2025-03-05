import { useState, useEffect } from 'react';
import { pointsService } from '../api/index';

export function usePointsData({ id, useOrderedEndpoint = true } = {}) {
  const [points, setPoints] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    setLoading(true);
    
    // Choose which points API to call based on the useOrderedEndpoint flag
    const pointsPromise = useOrderedEndpoint 
      ? pointsService.fetchOrderedPointsByID(id)
      : pointsService.fetchPointsByID(id);
    
    Promise.all([
      pointsPromise
    ])
    .then(([pointsDataJson]) => {
      setPoints(pointsDataJson);
      setLoading(false);
    })
    .catch(error => {
      console.error('Error fetching data: ', error);
      setError(error);
      setLoading(false);
    });
  }, [id, useOrderedEndpoint]);

  return { points, loading, error };
}