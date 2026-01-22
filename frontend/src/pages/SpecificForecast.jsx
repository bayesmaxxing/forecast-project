import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import {
  Box,
  Container,
  Typography,
  Paper,
  CircularProgress,
  Alert,
  Collapse,
  IconButton,
  Grid
} from '@mui/material';
import { ExpandMore } from '@mui/icons-material';
import ForecastGraph from '../components/ForecastGraph';
import ForecastPointList from '../components/ForecastPointList';
import ResolutionDetails from '../components/ResolutionDetails';
import UpdateForecast from '../components/UpdateForecast';
import ResolveForecast from '../components/ResolveForecast';
import UserSelector from '../components/UserSelector';
import ForecastLeaderboard from '../components/ForecastLeaderboard';
import { useForecastData } from '../services/hooks/useForecastData';
import { usePointsData } from '../services/hooks/usePointsData';
import { useScoresData } from '../services/hooks/useScoresData';
import { prepareChartData } from '../utils/chartDataUtils';
import { useUserData } from '../services/hooks/useUserData';

function SpecificForecast() {
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [selectedUserId, setSelectedUserId] = useState('all');
  const [criteriaExpanded, setCriteriaExpanded] = useState(false);
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
    prepareChartData(filteredPoints, isMultiUserMode, false, 0, users) : null;
  
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
    <Container maxWidth="lg" sx={{ py: 2, mt: { xs: 5, sm: 7}, mb: 2 }}>
      <Paper elevation={0} sx={{ p: 2.5, mb: 2 }}>
        <Box display="flex" justifyContent="space-between" alignItems="center">
          <Typography variant="h5" component="h1" sx={{ fontWeight: 600 }}>
            {forecast.question}
          </Typography>
          {isLoggedIn && forecast.resolved == null && (
            <ResolveForecast onSubmitSuccess={refetchAllData} />
          )}
        </Box>
      </Paper>
      {forecast.resolved && (
        <Grid container spacing={2} sx={{ mb: 2 }}>
          <Grid item xs={12} md={6}>
            <Paper elevation={0} sx={{ p: 2.5, height: '100%' }}>
              <ResolutionDetails forecast={forecast} score={scores} />
            </Paper>
          </Grid>
          <Grid item xs={12} md={6}>
            <ForecastLeaderboard
              forecastId={numericId}
              isResolved={forecast.resolved != null}
            />
          </Grid>
        </Grid>
      )}

      {chartData ? (
        <ForecastGraph
          data={chartData}
          options={chartOptions}
          selectedUserId={selectedUserId}
          onUserChange={handleUserChange}
        />
      ) : (
        <Paper elevation={0} sx={{ p: 2.5, mb: 2, textAlign: 'center' }}>
          <Typography variant="body1" color="text.secondary">
            No forecast data available for the selected user
          </Typography>
        </Paper>
      )}

      <Paper elevation={0} sx={{ p: 2.5, mb: 2, mt: 2}}>
        <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', cursor: 'pointer' }} onClick={() => setCriteriaExpanded(!criteriaExpanded)}>
          <Typography variant="subtitle1" sx={{ fontWeight: 600 }}>Resolution Criteria</Typography>
          <IconButton
            size="small"
            sx={{
              transform: criteriaExpanded ? 'rotate(180deg)' : 'rotate(0deg)',
              transition: 'transform 0.3s',
            }}
          >
            <ExpandMore />
          </IconButton>
        </Box>
        <Collapse in={criteriaExpanded}>
          <Typography variant="body2" sx={{ lineHeight: 1.6, mt: 1 }}>{forecast.resolution_criteria}</Typography>
        </Collapse>
      </Paper>

      <Paper elevation={0} sx={{ p: 2.5, mb: 2 }}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 1.5 }}>
          <Typography variant="subtitle1" sx={{ fontWeight: 600 }}>Forecast Updates</Typography>
          {isLoggedIn && forecast.resolved == null && (
            <UpdateForecast onSubmitSuccess={refetchAllData} />
          )}
        </Box>
        {sortedPoints.length > 0 ? (
          <ForecastPointList points={sortedPoints} users={users}/>
        ) : (
          <Typography variant="body2" color="text.secondary" sx={{ textAlign: 'center', mt: 1 }}>
            No forecast updates available for the selected user
          </Typography>
        )}
      </Paper>

      
    </Container>
  );
}

export default SpecificForecast;
