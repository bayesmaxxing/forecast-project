import { useState, useEffect, useMemo } from 'react';

export const useSearchFilter = (forecasts) => {
  const [searchTerm, setSearchTerm] = useState('');

  const handleSearch = (term) => {
    setSearchTerm(term);
  };

  const sortedForecasts = useMemo(() => {
    if (!forecasts) return [];
    
    // Only apply search term filtering (no user filtering for forecasts)
    let filtered = forecasts;
    if (searchTerm.trim() !== '') {
      const term = searchTerm.toLowerCase();
      filtered = forecasts.filter(forecast => 
        forecast.question?.toLowerCase().includes(term) || 
        forecast.resolution_criteria?.toLowerCase().includes(term) ||
        forecast.category?.toLowerCase().includes(term)
      );
    }
    
    // Sort by creation date (newest first)
    return [...filtered].sort((a, b) => {
      return new Date(b.created) - new Date(a.created);
    });
  }, [forecasts, searchTerm]);

  return {
    handleSearch,
    searchTerm,
    sortedForecasts,
  };
}; 