import { useState, useEffect, useMemo } from 'react';

export const useSearchFilter = (forecasts, options = {}) => {
  const [searchTerm, setSearchTerm] = useState('');
  const { userId = 'all' } = options;

  const handleSearch = (term) => {
    setSearchTerm(term);
  };

  const sortedForecasts = useMemo(() => {
    if (!forecasts) return [];
    
    // First, filter by user if a specific user is selected
    let filtered = forecasts;
    if (userId !== 'all') {
      filtered = forecasts.filter(forecast => 
        forecast.user_id === userId || forecast.created_by === userId
      );
    }
    
    // Then apply search term filtering
    if (searchTerm.trim() !== '') {
      const term = searchTerm.toLowerCase();
      filtered = filtered.filter(forecast => 
        forecast.question?.toLowerCase().includes(term) || 
        forecast.resolution_criteria?.toLowerCase().includes(term) ||
        forecast.category?.toLowerCase().includes(term)
      );
    }
    
    // Sort by creation date (newest first)
    return [...filtered].sort((a, b) => {
      return new Date(b.created) - new Date(a.created);
    });
  }, [forecasts, searchTerm, userId]);

  return {
    handleSearch,
    searchTerm,
    sortedForecasts,
  };
}; 