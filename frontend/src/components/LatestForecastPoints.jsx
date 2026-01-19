import React from 'react';
import { Link } from 'react-router-dom';
import {
  Paper,
  Box,
  Typography,
  List,
  ListItem,
  Chip,
  Skeleton,
  useTheme
} from '@mui/material';
import { TrendingUp } from '@mui/icons-material';
import { useLatestForecastPoints } from '../services/hooks/useLatestForecastPoints';
import { useUserData } from '../services/hooks/useUserData';

function LatestForecastPoints({ userId = 'all' }) {
  const theme = useTheme();
  const { points, loading, error } = useLatestForecastPoints({ userId, limit: 10 });
  const { users } = useUserData();

  const getUserName = (userId) => {
    const user = users?.find(u => u.id === userId);
    return user?.username || `User ${userId}`;
  };

  const formatTimestamp = (timestamp) => {
    const date = new Date(timestamp);
    return date.toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  const getProbabilityColor = (probability) => {
    // Gradient from red (0%) through yellow (50%) to green (100%)
    if (probability < 0.5) {
      const ratio = probability / 0.5;
      return `rgb(${255}, ${Math.round(255 * ratio)}, 0)`;
    } else {
      const ratio = (probability - 0.5) / 0.5;
      return `rgb(${Math.round(255 * (1 - ratio))}, ${200}, 0)`;
    }
  };

  const truncateReason = (reason, maxLength = 500) => {
    if (!reason) return null;
    if (reason.length <= maxLength) return reason;
    return reason.substring(0, maxLength).trim() + '...';
  };

  if (loading) {
    return (
      <Paper elevation={0} sx={{ p: 2.5, height: '100%', width: '100%' }}>
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 2 }}>
          <TrendingUp sx={{ color: 'primary.main' }} />
          <Typography variant="subtitle1" sx={{ fontWeight: 600 }}>
            Latest Forecasts
          </Typography>
        </Box>
        {[...Array(5)].map((_, index) => (
          <Skeleton key={index} variant="rectangular" height={80} sx={{ mb: 1, borderRadius: 1 }} />
        ))}
      </Paper>
    );
  }

  if (error) {
    return (
      <Paper elevation={0} sx={{ p: 2.5, height: '100%', width: '100%' }}>
        <Typography color="error">Error loading forecast points: {error.message}</Typography>
      </Paper>
    );
  }

  return (
    <Paper
      elevation={0}
      sx={{
        p: 2.5,
        bgcolor: 'background.paper',
        height: '100%',
        width: '100%',
        display: 'flex',
        flexDirection: 'column'
      }}
    >
      <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 2 }}>
        <TrendingUp sx={{ color: 'primary.main' }} />
        <Typography variant="subtitle1" sx={{ fontWeight: 600 }}>
          Latest Forecasts
        </Typography>
      </Box>

      <List sx={{ pt: 0, flex: 1, overflow: 'auto' }}>
        {points.map((point) => (
          <ListItem
            key={point.id}
            component={Link}
            to={`/forecast/${point.forecast_id}`}
            sx={{
              px: 1.5,
              py: 1.5,
              borderRadius: 1,
              mb: 1,
              textDecoration: 'none',
              color: 'inherit',
              flexDirection: 'column',
              alignItems: 'flex-start',
              transition: 'all 0.2s ease-in-out',
              '&:hover': {
                bgcolor: 'action.hover',
                transform: 'translateX(4px)',
              }
            }}
          >
            <Box sx={{ display: 'flex', alignItems: 'center', width: '100%', gap: 1.5, flexWrap: 'wrap', mb: point.reason ? 1 : 0 }}>
              <Chip
                label={getUserName(point.user_id)}
                size="small"
                sx={{
                  fontWeight: 500,
                  bgcolor: 'primary.main',
                  color: 'primary.contrastText',
                  minWidth: 80
                }}
              />

              <Chip
                label={`${(point.point_forecast * 100).toFixed(0)}%`}
                size="small"
                sx={{
                  fontWeight: 600,
                  minWidth: 50,
                  background: `linear-gradient(135deg, ${getProbabilityColor(point.point_forecast)}, ${theme.palette.grey[400]})`,
                  color: 'white'
                }}
              />

              <Box sx={{ flex: 1, minWidth: 0, textAlign: 'right' }}>
                <Typography variant="caption" color="text.secondary">
                  {formatTimestamp(point.created)}
                </Typography>
                <Typography variant="caption" color="text.secondary" sx={{ display: 'block' }}>
                  Forecast #{point.forecast_id}
                </Typography>
              </Box>
            </Box>

            {point.reason && (
              <Typography
                variant="body2"
                color="text.secondary"
                sx={{
                  width: '100%',
                  mt: 0.5
                }}
              >
                {truncateReason(point.reason)}
              </Typography>
            )}
          </ListItem>
        ))}
        {points.length === 0 && (
          <Typography variant="body2" color="text.secondary" align="center" sx={{ py: 2 }}>
            No forecast points available
          </Typography>
        )}
      </List>
    </Paper>
  );
}

export default LatestForecastPoints;
