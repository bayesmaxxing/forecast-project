import React, { useState } from 'react';
import {
  Paper,
  Box,
  Typography,
  List,
  ListItem,
  ListItemText,
  FormControl,
  Select,
  MenuItem,
  Divider,
  Chip,
  useTheme
} from '@mui/material';
import { EmojiEvents } from '@mui/icons-material';
import { useUserData } from '../services/hooks/useUserData';
import { useAggregateScoresByUsers } from '../services/hooks/useAggregateScoresByUsers';

function UserLeaderboard({ dateRange = null }) {
  const [selectedMetric, setSelectedMetric] = useState('brier_score');
  const theme = useTheme();
  const { users, usersLoading, usersError } = useUserData();

  // Fetch aggregate scores for all users
  const { scores, loading: scoresLoading, error: scoresError } = useAggregateScoresByUsers(dateRange);

  // Map user data with scores
  const userScores = scores?.map(scoreData => {
    const user = users?.find(u => u.id === scoreData.user_id);
    return {
      id: scoreData.user_id,
      username: user?.username || `User ${scoreData.user_id}`,
      brier_score: scoreData.brier_score_time_weighted,
      log2_score: scoreData.log2_score_time_weighted,
      logn_score: scoreData.logn_score_time_weighted,
      forecastCount: scoreData.total_forecasts
    };
  }) || [];

  // Sort users by score (lower is better for brier, higher for log scores)
  const sortedUsers = [...userScores]
    .filter(user => user.forecastCount > 0)
    .sort((a, b) => {
      if (selectedMetric === 'brier_score') {
        return a[selectedMetric] - b[selectedMetric]; // Lower is better
      } else {
        return b[selectedMetric] - a[selectedMetric]; // Higher is better for log scores
      }
    });

  const getMetricLabel = (metric) => {
    switch (metric) {
      case 'brier_score':
        return 'Brier Score';
      case 'log2_score':
        return 'Binary Log Score';
      case 'logn_score':
        return 'Natural Log Score';
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

  if (usersLoading || scoresLoading) {
    return (
      <Paper elevation={0} sx={{ p: 2.5, height: '100%', width: '100%' }}>
        <Typography>Loading leaderboard...</Typography>
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
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
          <EmojiEvents sx={{ color: 'primary.main' }} />
          <Typography variant="subtitle1" sx={{ fontWeight: 600 }}>
            Leaderboard
          </Typography>
        </Box>

        <FormControl size="small" sx={{ minWidth: 120 }}>
          <Select
            value={selectedMetric}
            onChange={(e) => setSelectedMetric(e.target.value)}
          >
            <MenuItem value="brier_score">Brier</MenuItem>
            <MenuItem value="log2_score">Binary Log</MenuItem>
            <MenuItem value="logn_score">Natural Log</MenuItem>
          </Select>
        </FormControl>
      </Box>

      <List sx={{ pt: 0, flex: 1, overflow: 'auto' }}>
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
                  <Typography variant="caption" color="text.secondary">
                    {user.forecastCount} forecast{user.forecastCount !== 1 ? 's' : ''}
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
            No data available
          </Typography>
        )}
      </List>
    </Paper>
  );
}

export default UserLeaderboard;
