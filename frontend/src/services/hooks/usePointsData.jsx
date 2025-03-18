import { useState, useEffect } from 'react';
import { pointsService } from '../api/index';

export function usePointsData({ id, userId, useOrderedEndpoint = true, useLatestPoints = true } = {}) {
  const [points, setPoints] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    setLoading(true);
    const user = userId === 'all' ? null : userId;
    let pointsPromise; 
    
    if (useLatestPoints && !user) {
      pointsPromise = pointsService.fetchLatestPoints();
    } else if (useLatestPoints && user) {
      pointsPromise = pointsService.fetchLatestPointsByUser(user);
    } else {
      pointsPromise = useOrderedEndpoint 
        ? pointsService.fetchOrderedPointsByID(id)
        : pointsService.fetchPointsByID(id);
    }
    
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
  }, [id, useOrderedEndpoint, userId, useLatestPoints]);

  return { points, loading, error };
}