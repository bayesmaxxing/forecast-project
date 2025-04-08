import { useState, useEffect } from 'react';
import { fetchUsers } from '../api/userService';

export const useUserData = () => {
  const [users, setUsers] = useState([]);
  const [userLoading, setUserLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchUserData = async () => {
      try {
        const data = await fetchUsers();
        setUsers(data);
        setError(null);
      } catch (err) {
        console.error('Error fetching users:', err);
        setError(err.message || 'Failed to load users');
        setUsers([]);
      } finally {
        setUserLoading(false);
      }
    };

    fetchUserData();
  }, []);

  return { users, userLoading, error };
}; 