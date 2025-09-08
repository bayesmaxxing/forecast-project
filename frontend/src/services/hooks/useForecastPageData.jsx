import { useMemo, useCallback } from 'react';
import { useForecastList } from './useForecastList';
import { usePointsData } from './usePointsData';
import { useSearchFilter } from './useSearchFilter';
import { useAggregateScoresData, useAggregateScoresDataByCategory } from './useAggregateScoresData';

export function useForecastPageData({ categoryFilter, listType, selectedUserId }) {
  // Fetch the data using existing hooks
  const { forecasts = [], loading: forecastsLoading, error: forecastsError, refetchForecasts } = useForecastList({
    category: categoryFilter, 
    list_type: listType
  });
  
  const { points = [], loading: pointsLoading, error: pointsError, refetchPoints } = usePointsData({
    userId: selectedUserId, 
    useLatestPoints: true, 
    useOrderedEndpoint: false
  });

  // Call both hooks (React requires unconditional hook calls)
  const { scores: allScores, scoresLoading: allScoresLoading, error: allScoresError } = useAggregateScoresData(selectedUserId);
  const { scores: categoryScores, scoresLoading: categoryScoresLoading, error: categoryScoresError } = useAggregateScoresDataByCategory(selectedUserId, categoryFilter);
  
  // Select the appropriate scores based on category filter
  const scores = categoryFilter ? categoryScores : allScores;
  const scoresLoading = categoryFilter ? categoryScoresLoading : allScoresLoading;
  const scoresError = categoryFilter ? categoryScoresError : allScoresError;

  // Combine the forecasts and points
  const combined = useMemo(() => {
    return Array.isArray(forecasts) 
      ? forecasts.map(forecast => {
          const matchingPoint = points.find(point => point.forecast_id === forecast.id);
          return { ...forecast, latestPoint: matchingPoint || null};
        })
      : [];
  }, [forecasts, points]);

  // Use search filter on combined data (no user filtering, only search filtering)
  const { handleSearch, sortedForecasts } = useSearchFilter(combined);

  // Function to refetch all relevant data
  const refetchAllData = useCallback(() => {
    refetchForecasts();
    refetchPoints();
  }, [refetchForecasts, refetchPoints]);

  // Consolidated loading and error states
  const loading = forecastsLoading;
  const error = forecastsError || pointsError || scoresError;

  return {
    sortedForecasts,
    scores,
    loading,
    pointsLoading,
    scoresLoading,
    error,
    handleSearch,
    refetchAllData
  };
}