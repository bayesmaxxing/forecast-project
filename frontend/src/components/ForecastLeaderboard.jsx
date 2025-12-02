import React, { useState, useEffect } from 'react';
import {
  Paper,
  Box,
  Typography,
  List,
  ListItem,
  FormControl,
  Select,
  MenuItem,
  Chip,
  useTheme,
  CircularProgress,
  Divider
} from '@mui/material';
import { EmojiEvents } from '@mui/icons-material';
import { fetchScores } from '../services/api/scoreService';
import { useUserData } from '../services/hooks/useUserData';

function ForecastLeaderboard({ forecastId, isResolved }) {
  const [selectedMetric, setSelectedMetric] = useState('brier_score_time_weighted');
  const [scores, setScores] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const theme = useTheme();

  const { users, usersLoading } = useUserData();

  useEffect(() => {
    // Only fetch if forecast is resolved
    if (!isResolved) {
      setLoading(false);
      return;
    }

    const fetchForecastScores = async () => {
      try {
        setLoading(true);
        // Fetch all scores for this specific forecast
        const data = await fetchScores(null, forecastId);
        setScores(data || []);
        setError(null);
      } catch (err) {
        console.error('Error fetching forecast scores:', err);
        setError(err.message || 'Failed to load scores');
        setScores([]);
      } finally {
        setLoading(false);
      }
    };

    fetchForecastScores();
  }, [forecastId, isResolved]);

  // Don't render if forecast is not resolved
  if (!isResolved) {
    return null;
  }

  if (loading || usersLoading) {
    return (
      <Paper elevation={0} sx={{ p: 2.5, height: '100%' }}>
        <Box display="flex" justifyContent="center" alignItems="center">
          <CircularProgress size={24} />
          <Typography sx={{ ml: 2 }}>Loading leaderboard...</Typography>
        </Box>
      </Paper>
    );
  }

  if (error) {
    return (
      <Paper elevation={0} sx={{ p: 2.5, height: '100%' }}>
        <Typography color="error">Error loading leaderboard: {error}</Typography>
      </Paper>
    );
  }

  // Map scores to users
  const userScoreData = scores.map(score => {
    const user = users?.find(u => u.id === score.user_id);
    return {
      id: score.user_id,
      username: user?.username || `User ${score.user_id}`,
      brier_score: score.brier_score,
      brier_score_time_weighted: score.brier_score_time_weighted,
      log2_score: score.log2_score,
      log2_score_time_weighted: score.log2_score_time_weighted,
      logn_score: score.logn_score,
      logn_score_time_weighted: score.logn_score_time_weighted
    };
  });

  // Sort users by selected metric
  const sortedUsers = [...userScoreData].sort((a, b) => {
    // For Brier score, lower is better
    if (selectedMetric.includes('brier')) {
      return a[selectedMetric] - b[selectedMetric];
    } else {
      // For log scores, higher is better
      return b[selectedMetric] - a[selectedMetric];
    }
  }).slice(0, 10); // Only show top 10

  const getMetricLabel = (metric) => {
    switch (metric) {
      case 'brier_score':
        return 'Brier Score';
      case 'brier_score_time_weighted':
        return 'Brier Score (Time-Weighted)';
      case 'log2_score':
        return 'Binary Log Score';
      case 'log2_score_time_weighted':
        return 'Binary Log Score (Time-Weighted)';
      case 'logn_score':
        return 'Natural Log Score';
      case 'logn_score_time_weighted':
        return 'Natural Log Score (Time-Weighted)';
      default:
        return 'Score';
    }
  };

  const getMedalColor = (index) => {
    switch (index) {
      case 0:
        return '#FFD700'; // Gold
      case 1:
        return '#C0C0C0'; // Silver
      case 2:
        return '#CD7F32'; // Bronze
      default:
        return 'transparent';
    }
  };

  return (
    <Paper elevation={0} sx={{ p: 2.5, height: '100%' }}>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
          <EmojiEvents sx={{ color: 'primary.main' }} />
          <Typography variant="subtitle1" sx={{ fontWeight: 600 }}>
            Top 10 Forecasters
          </Typography>
        </Box>

        <FormControl size="small" sx={{ minWidth: 180 }}>
          <Select
            value={selectedMetric}
            onChange={(e) => setSelectedMetric(e.target.value)}
          >
            <MenuItem value="brier_score">Brier</MenuItem>
            <MenuItem value="brier_score_time_weighted">Brier (Time-Weighted)</MenuItem>
            <MenuItem value="log2_score">Binary Log</MenuItem>
            <MenuItem value="log2_score_time_weighted">Binary Log (Time-Weighted)</MenuItem>
            <MenuItem value="logn_score">Natural Log</MenuItem>
            <MenuItem value="logn_score_time_weighted">Natural Log (Time-Weighted)</MenuItem>
          </Select>
        </FormControl>
      </Box>

      <List sx={{ pt: 0 }}>
        {sortedUsers.map((user, index) => (
          <React.Fragment key={user.id}>
            <ListItem
              sx={{
                px: 1.5,
                py: 1,
                borderRadius: 1,
                mb: 0.5,
                bgcolor: index < 3 ? `${getMedalColor(index)}10` : 'transparent',
                border: index < 3 ? `1px solid ${getMedalColor(index)}40` : 'none',
                transition: 'all 0.2s ease-in-out',
                '&:hover': {
                  bgcolor: 'action.hover',
                  transform: 'translateX(4px)',
                }
              }}
            >
              <Box sx={{ display: 'flex', alignItems: 'center', width: '100%', gap: 1.5 }}>
                <Typography
                  variant="body2"
                  sx={{
                    minWidth: 24,
                    fontWeight: 700,
                    color: index < 3 ? getMedalColor(index) : 'text.secondary',
                    fontSize: index < 3 ? '1rem' : '0.875rem'
                  }}
                >
                  {index + 1}
                </Typography>

                <Box sx={{ flex: 1, minWidth: 0 }}>
                  <Typography
                    variant="body2"
                    sx={{
                      fontWeight: index < 3 ? 600 : 400,
                      overflow: 'hidden',
                      textOverflow: 'ellipsis',
                      whiteSpace: 'nowrap'
                    }}
                  >
                    {user.username}
                  </Typography>
                </Box>

                <Chip
                  label={user[selectedMetric]?.toFixed(3) ?? 'N/A'}
                  size="small"
                  sx={{
                    fontWeight: 600,
                    minWidth: 60,
                    background: index < 3
                      ? `linear-gradient(135deg, ${theme.palette.primary.main}, ${theme.palette.secondary.main})`
                      : 'divider',
                    color: index < 3 ? 'white' : 'text.primary'
                  }}
                />
              </Box>
            </ListItem>
            {index < sortedUsers.length - 1 && index === 2 && (
              <Divider sx={{ my: 1 }} />
            )}
          </React.Fragment>
        ))}
        {sortedUsers.length === 0 && (
          <Typography variant="body2" color="text.secondary" align="center" sx={{ py: 2 }}>
            No scores available for this forecast
          </Typography>
        )}
      </List>
    </Paper>
  );
}

export default ForecastLeaderboard;
