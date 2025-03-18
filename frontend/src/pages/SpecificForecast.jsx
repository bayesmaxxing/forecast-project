import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import {
  Box,
  Container,
  Typography,
  Paper,
  CircularProgress,
  Alert
} from '@mui/material';
import ForecastGraph from '../components/ForecastGraph';
import ForecastPointList from '../components/ForecastPointList';
import ResolutionDetails from '../components/ResolutionDetails';
import UpdateForecast from '../components/UpdateForecast';
import ResolveForecast from '../components/ResolveForecast';
import UserSelector from '../components/UserSelector';
import { useForecastData } from '../services/hooks/useForecastData';
import { usePointsData } from '../services/hooks/usePointsData';
import { useScoresData } from '../services/hooks/useScoresData';
import { prepareChartData } from '../utils/chartDataUtils';

function SpecificForecast() {
  const [isAdmin, setIsAdmin] = useState(false);
  const [selectedUserId, setSelectedUserId] = useState('all');
  let { id } = useParams();
  const numericId = parseInt(id, 10);

  // Always use multi-user mode when all users are selected
  // When a specific user is selected, still use single-user mode
  const isMultiUserMode = selectedUserId === 'all';

  const { forecast, forecastLoading, forecastError } = useForecastData({id: id});
  
  // Use the ordered points endpoint for better multi-user graph visualization
  const { points, pointsLoading, pointsError } = usePointsData({
    id: id,
    useOrderedEndpoint: true,
    useLatestPoints: false
  });
  
  // Only fetch scores if the forecast is resolved, using the shouldFetch parameter
  // This prevents sending requests that would result in 500 errors for unresolved forecasts
  const { score, scoreLoading, scoreError } = useScoresData({
    user_id: isMultiUserMode ? null : selectedUserId,
    forecast_id: numericId,
    useAverageEndpoint: isMultiUserMode,
    shouldFetch: forecast?.resolved != null // Only fetch if the forecast has been resolved
  });
  
  // Handle user selection change
  const handleUserChange = (userId) => {
    setSelectedUserId(userId);
  };

  // Check if we're loading any data - only check score loading if the forecast is resolved
  const isLoading = forecastLoading || pointsLoading || (forecast?.resolved != null && scoreLoading);
  
  if (isLoading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="200px">
        <CircularProgress />
      </Box>
    );
  }

  // Only include scoreError in the condition if the forecast is resolved
  // This prevents showing errors for scores when the forecast is not resolved
  const hasError = forecastError || pointsError || (forecast?.resolved != null && scoreError) || !forecast;
  
  if (hasError) {
    const errorMessage = forecastError?.message || 
                         pointsError?.message || 
                         (forecast?.resolved != null ? scoreError?.message : null) || 
                         "Forecast not found";
    
    return (
      <Box m={2}>
        <Alert severity="error">
          Error loading the forecast: {errorMessage}
        </Alert>
      </Box>
    );
  }
  
  // Filter points based on selected user if not showing all users
  const filteredPoints = selectedUserId === 'all' 
    ? [...(points || [])] 
    : [...(points || [])].filter(point => String(point.user_id) === String(selectedUserId));
    
  const sortedPoints = [...filteredPoints].sort((a, b) => new Date(a.created) - new Date(b.created));
  
  // Debug the filtering results
  console.log('Filtering points:', {
    allPoints: points?.length || 0,
    selectedUserId,
    filteredCount: filteredPoints.length,
    multiUserMode: isMultiUserMode
  });
  
  // Use the utility function to prepare chart data - ensure we have points before trying to prepare chart data
  const chartData = filteredPoints && filteredPoints.length > 0 ? prepareChartData(filteredPoints, isMultiUserMode) : null;
  
  const chartOptions = {
    title: {
      text: isMultiUserMode ? 'User Predictions Over Time' : 'Prediction Over Time'
    },
    scales: {
      y: {
        min: 0,
        max: 1,
      }
    }
  };

  return (
    <Container maxWidth="lg" sx={{ py: 4, mt: { xs: 8, sm: 10}, mb: 4 }}>
      <Paper elevation={3} sx={{ p: 3, mb: 3 }}>
        <Typography variant="h4" component="h1" gutterBottom>
          {forecast.question}
        </Typography>
        <Box display="flex" justifyContent="space-between" alignItems="center" mt={2} mb={2}>
          <UserSelector onUserChange={handleUserChange} selectedUserId={selectedUserId} />
        </Box>
        {forecast.resolved && <ResolutionDetails forecast={forecast} score={score} />}
      </Paper>

      {chartData ? (
        <ForecastGraph data={chartData} options={chartOptions} />
      ) : (
        <Paper elevation={2} sx={{ p: 3, mb: 3, textAlign: 'center' }}>
          <Typography variant="h6" color="text.secondary">
            No forecast data available for the selected user
          </Typography>
        </Paper>
      )}

      {isAdmin && forecast.resolved == null && (
        <Paper elevation={3} sx={{ p: 3, mb: 3 }}>
          <ResolveForecast forecastPoints={sortedPoints} />
        </Paper>
      )}

      {forecast.resolved != null && (
        <Paper elevation={3} sx={{ p: 3, mb: 3 }}>
          <ResolutionDetails forecast={forecast} score={score} />
        </Paper>
      )}

      <Paper elevation={3} sx={{ p: 3, mb: 3 }}>
        <Typography variant="h6" gutterBottom>Resolution Criteria</Typography>
        <Typography variant="body1">{forecast.resolution_criteria}</Typography>
      </Paper>

      <Paper elevation={3} sx={{ p: 3, mb: 3 }}>
        <Typography variant="h6" gutterBottom>Forecast Updates</Typography>
        {sortedPoints.length > 0 ? (
          <ForecastPointList points={sortedPoints} />
        ) : (
          <Typography variant="body1" color="text.secondary" sx={{ textAlign: 'center', mt: 2 }}>
            No forecast updates available for the selected user
          </Typography>
        )}
      </Paper>

      {isAdmin && forecast.resolved == null && (
        <Paper elevation={3} sx={{ p: 3 }}>
          <UpdateForecast />
        </Paper>
      )}
    </Container>
  );
}

export default SpecificForecast;
