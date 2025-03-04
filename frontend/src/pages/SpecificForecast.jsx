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
import { useForecastData } from '../services/hooks/useForecastData';
import { usePointsData } from '../services/hooks/usePointsData';
import { useScoresData } from '../services/hooks/useScoresData';
import { prepareChartData } from '../utils/chartDataUtils';

function SpecificForecast() {
  const [isAdmin, setIsAdmin] = useState(false);
  let { id } = useParams();
  const numericId = parseInt(id, 10);

  const { forecast, forecastLoading, forecastError } = useForecastData({id: id});
  const { points, pointsLoading, pointsError } = usePointsData({id: id});
  
  const { score, scoreLoading, scoreError } = useScoresData({forecast_id: numericId});

  // Determine whether we should use multi-user mode
  // This could be based on a prop, context, or derived from the data
  const isMultiUserMode = true; // For testing, set this to true or false

  if (forecastLoading || pointsLoading || (forecast?.resolved != null && scoreLoading)) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="200px">
        <CircularProgress />
      </Box>
    );
  }

  if (forecastError || pointsError || scoreError || !forecast) {
    return (
      <Box m={2}>
        <Alert severity="error">
          Error loading the forecast: {forecastError?.message || pointsError?.message || scoreError?.message || "Forecast not found"}
        </Alert>
      </Box>
    );
  }
  
  const sortedPoints = [...(points || [])].sort((a, b) => new Date(a.created) - new Date(b.created));
  
  // Use the utility function to prepare chart data
  const chartData = prepareChartData(points, isMultiUserMode);
  
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
        {forecast.resolved && <ResolutionDetails forecast={forecast} score={score} />}
      </Paper>

      <ForecastGraph data={chartData} options={chartOptions} />

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
        <ForecastPointList points={sortedPoints} />
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
