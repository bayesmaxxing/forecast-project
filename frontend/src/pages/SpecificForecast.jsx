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
import { useUserData } from '../services/hooks/useUserData';

function SpecificForecast() {
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [selectedUserId, setSelectedUserId] = useState('all');
  let { id } = useParams();
  const numericId = parseInt(id, 10);

  // USER SELECTION
  // Always use multi-user mode when all users are selected
  const isMultiUserMode = selectedUserId === 'all';

  // Handle user selection change
  const handleUserChange = (userId) => {
    setSelectedUserId(userId);
  };

  // AUTH CHECK
  useEffect(() => {
    const token = localStorage.getItem('token');
    setIsLoggedIn(!!token);
  }, []);

  // DATA FETCHING
  const { forecast, forecastLoading, forecastError, refetchForecast } = useForecastData({id: id});
  const { points, pointsLoading, pointsError, refetchPoints } = usePointsData({
    id: id,
    useOrderedEndpoint: false,
    useLatestPoints: false
  });
  
  const { scores, scoreLoading, scoreError, refetchScores } = useScoresData({
    user_id: isMultiUserMode ? null : selectedUserId,
    forecast_id: numericId,
    useAverageEndpoint: isMultiUserMode,
    shouldFetch: forecast?.resolved != null // Only fetch if the forecast has been resolved
  });
  
  const { users, usersLoading, usersError } = useUserData();
  
  // Function to refetch all data
  const refetchAllData = () => {
    refetchForecast();
    refetchPoints();
    if (forecast?.resolved != null) {
      refetchScores();
    }
  };

  const isLoading = forecastLoading || pointsLoading || (forecast?.resolved != null && scoreLoading);
  
  if (isLoading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="200px">
        <CircularProgress />
      </Box>
    );
  }

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
  
  // CHART PREP
  // Filter points based on selected user if not showing all users
  const filteredPoints = selectedUserId === 'all' 
    ? [...(points || [])] 
    : [...(points || [])].filter(point => String(point.user_id) === String(selectedUserId));
    
  const sortedPoints = [...filteredPoints].sort((a, b) => new Date(a.created) - new Date(b.created));
  
  // Use a 4-hour minimum time window between points
  const chartData = filteredPoints && filteredPoints.length > 0 ? 
    prepareChartData(filteredPoints, isMultiUserMode, false, 0) : null;
  
  const chartOptions = {
    title: {
      text: isMultiUserMode ? 'User Predictions Over Time' : 'Prediction Over Time'
    },
    scales: {
      y: {
        min: 0,
        max: 1,
      }
    },
    useSequential: false // Use date-based x-axis
  };

  return (
    <Container maxWidth="lg" sx={{ py: 4, mt: { xs: 8, sm: 10}, mb: 4 }}>
      <Paper elevation={3} sx={{ p: 3, mb: 3 }}>
        <Typography variant="h4" component="h1" gutterBottom>
          {forecast.question}
        </Typography>
        <Box display="flex" justifyContent="space-between" alignItems="center" mt={2} mb={2}>
          <UserSelector onUserChange={handleUserChange} selectedUserId={selectedUserId} />
          {isLoggedIn && forecast.resolved == null && (
            <ResolveForecast onSubmitSuccess={refetchAllData} />
          )}
        </Box>
        
      </Paper>
      {forecast.resolved && (
          <Paper elevation={3} sx={{ p: 3, mb: 3 }}>
            <ResolutionDetails forecast={forecast} score={scores} />
          </Paper>
        )}

      {chartData ? (
        <ForecastGraph data={chartData} options={chartOptions} />
      ) : (
        <Paper elevation={2} sx={{ p: 3, mb: 3, textAlign: 'center' }}>
          <Typography variant="h6" color="text.secondary">
            No forecast data available for the selected user
          </Typography>
        </Paper>
      )}

      <Paper elevation={3} sx={{ p: 3, mb: 3 , mt: 3}}>
        <Typography variant="h6" gutterBottom>Resolution Criteria</Typography>
        <Typography variant="body1">{forecast.resolution_criteria}</Typography>
      </Paper>

      <Paper elevation={3} sx={{ p: 3, mb: 3 }}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
          <Typography variant="h6">Forecast Updates</Typography>
          {isLoggedIn && forecast.resolved == null && (
            <UpdateForecast onSubmitSuccess={refetchAllData} />
          )}
        </Box>
        {sortedPoints.length > 0 ? (
          <ForecastPointList points={sortedPoints} users={users}/>
        ) : (
          <Typography variant="body1" color="text.secondary" sx={{ textAlign: 'center', mt: 2 }}>
            No forecast updates available for the selected user
          </Typography>
        )}
      </Paper>

      
    </Container>
  );
}

export default SpecificForecast;
