import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import {
  Box,
  Container,
  Typography,
  Paper,
  List,
  ListItem,
  ListItemText,
  Divider,
  CircularProgress,
  Alert,
  Chip
} from '@mui/material';
import ForecastGraph from '../components/ForecastGraph';
import UpdateForecast from '../components/UpdateForecast';
import ResolveForecast from '../components/ResolveForecast';

function SpecificForecast() {
  const [forecastData, setForecastData] = useState(null);
  const [forecastPoints, setForecastPoints] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [isAdmin, setIsAdmin] = useState(false);
  let { id } = useParams();

  useEffect(() => {
    Promise.all([
      fetch(`http://localhost:8080/forecasts/${id}`, {
        headers: {
          "Accept": "application/json"
        }
      }),
      fetch(`http://localhost:8080/forecast-points/${id}`, {
        headers: {
          "Accept": "application/json"
        }
      })
    ])
      .then(async ([idData, pointsData]) => {
        if (!idData.ok || !pointsData.ok) {
          throw new Error('Error fetching data');
        }
        const idJson = await idData.json();
        const pointsJson = await pointsData.json();
        return [idJson, pointsJson];
      })
      .then(([idJson, pointsJson]) => {
        setForecastData(idJson);
        if (pointsJson) {
          const sortedPoints = pointsJson.sort((a, b) => new Date(a.created) - new Date(b.created));
          setForecastPoints(sortedPoints);
        } else {
          setForecastPoints([]);
        }
        setLoading(false);
      })
      .catch(error => {
        setError(error);
        setLoading(false);
      });

    const checkAdminStatus = () => {
      const expirationTime = localStorage.getItem('adminLoginExpiration');
      setIsAdmin(expirationTime && new Date().getTime() < parseInt(expirationTime, 10));
    };

    checkAdminStatus();
  }, [id]);

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="200px">
        <CircularProgress />
      </Box>
    );
  }

  if (error) {
    return (
      <Box m={2}>
        <Alert severity="error">Error loading the forecast: {error.message}</Alert>
      </Box>
    );
  }

  const chartData = forecastPoints && forecastPoints.length > 0 ? {
    labels: forecastPoints.map(point => new Date(point.created).toLocaleDateString('en-CA')),
    datasets: [
      {
        label: 'Prediction',
        data: forecastPoints.map(point => point.point_forecast),
        fill: false,
        borderColor: 'rgb(75, 192, 192)',
        tension: 0.1
      }
    ]
  } : null;

  const chartOptions = {
    scales: {
      y: {
        min: 0,
        max: 1,
      }
    }
  };

  const formatDate = (dateString) => dateString.split('T')[0];
  const reversedForecastpoints = [...(forecastPoints || [])].reverse();
  const resolution = forecastData.resolution === "1" ? "Yes" :
    forecastData.resolution === "0" ? "No" :
      "Ambiguous";

  const getResolutionColor = (res) => {
    switch (res) {
      case "Yes": return "success";
      case "No": return "error";
      default: return "warning";
    }
  };

  return (
    <Container maxWidth="lg" sx={{ py: 4, mt: { xs: 8, sm: 10}, mb: 4 }}>
      <Paper elevation={3} sx={{ p: 3, mb: 3 }}>
        <Typography variant="h4" component="h1" gutterBottom>
          {forecastData.question}
        </Typography>
        {forecastData.resolved && (
          <Chip
            label={`Resolved as: ${resolution}`}
            color={getResolutionColor(resolution)}
            sx={{ mb: 2 }}
          />
        )}
      </Paper>

      {chartData && (
        <Paper elevation={3} sx={{ p: 3, mb: 3 }}>
          <ForecastGraph data={chartData} options={chartOptions} />
        </Paper>
      )}

      {isAdmin && forecastData.resolved == null && (
        <Paper elevation={3} sx={{ p: 3, mb: 3 }}>
          <ResolveForecast forecastPoints={forecastPoints} />
        </Paper>
      )}

      {forecastData.resolved != null && (
        <Paper elevation={3} sx={{ p: 3, mb: 3 }}>
          <Typography variant="h6" gutterBottom>Resolution Details</Typography>
          <Typography variant="body1" paragraph>
            Resolved on {formatDate(forecastData.resolved)} with a Brier score of {!forecastData.brier_score ? 0 : forecastData.brier_score}.
          </Typography>
          {forecastData.comment && (
            <Typography variant="body1" paragraph>
              <strong>Comment:</strong> {forecastData.comment}
            </Typography>
          )}
        </Paper>
      )}

      <Paper elevation={3} sx={{ p: 3, mb: 3 }}>
        <Typography variant="h6" gutterBottom>Resolution Criteria</Typography>
        <Typography variant="body1">{forecastData.resolution_criteria}</Typography>
      </Paper>

      <Paper elevation={3} sx={{ p: 3, mb: 3 }}>
        <Typography variant="h6" gutterBottom>Forecast Updates</Typography>
        {forecastPoints && forecastPoints.length > 0 ? (
          <List>
            {reversedForecastpoints.map((forecast, index) => (
              <React.Fragment key={forecast.forecast_id}>
                {index > 0 && <Divider />}
                <ListItem>
                  <ListItemText
                    primary={`Update to ${(forecast.point_forecast * 100).toFixed(1)}% on ${formatDate(forecast.created)}`}
                    secondary={forecast.reason}
                  />
                </ListItem>
              </React.Fragment>
            ))}
          </List>
        ) : (
          <Typography variant="body1">No forecast points available yet.</Typography>
        )}
      </Paper>

      {isAdmin && forecastData.resolved == null && (
        <Paper elevation={3} sx={{ p: 3 }}>
          <UpdateForecast />
        </Paper>
      )}
    </Container>
  );
}

export default SpecificForecast;
