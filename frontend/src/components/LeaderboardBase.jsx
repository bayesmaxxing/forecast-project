import React from 'react';
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
  Divider,
  CircularProgress,
  useTheme
} from '@mui/material';
import { EmojiEvents } from '@mui/icons-material';

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

/**
 * Reusable leaderboard component for displaying ranked lists of scores.
 *
 * @param {Object} props
 * @param {Array} props.items - Array of items to display: { id, name, subtitle?, [metricKey]: number }
 * @param {string} props.title - Leaderboard title
 * @param {Array} props.metrics - Array of { value, label } for metric options
 * @param {string} props.selectedMetric - Currently selected metric key
 * @param {Function} props.onMetricChange - Callback when metric changes
 * @param {boolean} props.loading - Loading state
 * @param {string} props.error - Error message
 * @param {string} props.emptyMessage - Message when no items
 * @param {number} props.maxItems - Maximum items to display (optional)
 */
function LeaderboardBase({
  items = [],
  title = 'Leaderboard',
  metrics = [],
  selectedMetric,
  onMetricChange,
  loading = false,
  error = null,
  emptyMessage = 'No data available',
  maxItems = null
}) {
  const theme = useTheme();

  // Sort items by selected metric
  const sortedItems = [...items].sort((a, b) => {
    // For Brier score, lower is better
    if (selectedMetric.includes('brier')) {
      return a[selectedMetric] - b[selectedMetric];
    }
    // For log scores, higher is better
    return b[selectedMetric] - a[selectedMetric];
  });

  // Apply max items limit if specified
  const displayItems = maxItems ? sortedItems.slice(0, maxItems) : sortedItems;

  if (loading) {
    return (
      <Paper elevation={0} sx={{ p: 2.5, height: '100%', width: '100%' }}>
        <Box display="flex" justifyContent="center" alignItems="center">
          <CircularProgress size={24} />
          <Typography sx={{ ml: 2 }}>Loading leaderboard...</Typography>
        </Box>
      </Paper>
    );
  }

  if (error) {
    return (
      <Paper elevation={0} sx={{ p: 2.5, height: '100%', width: '100%' }}>
        <Typography color="error">Error loading leaderboard: {error}</Typography>
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
            {title}
          </Typography>
        </Box>

        <FormControl size="small" sx={{ minWidth: 120 }}>
          <Select
            value={selectedMetric}
            onChange={(e) => onMetricChange(e.target.value)}
          >
            {metrics.map(({ value, label }) => (
              <MenuItem key={value} value={value}>{label}</MenuItem>
            ))}
          </Select>
        </FormControl>
      </Box>

      <List sx={{ pt: 0, flex: 1, overflow: 'auto' }}>
        {displayItems.map((item, index) => (
          <React.Fragment key={item.id}>
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
                    {item.name}
                  </Typography>
                  {item.subtitle && (
                    <Typography variant="caption" color="text.secondary">
                      {item.subtitle}
                    </Typography>
                  )}
                </Box>

                <Chip
                  label={item[selectedMetric]?.toFixed(3) ?? 'N/A'}
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
            {index < displayItems.length - 1 && index === 2 && (
              <Divider sx={{ my: 1 }} />
            )}
          </React.Fragment>
        ))}
        {displayItems.length === 0 && (
          <Typography variant="body2" color="text.secondary" align="center" sx={{ py: 2 }}>
            {emptyMessage}
          </Typography>
        )}
      </List>
    </Paper>
  );
}

export default LeaderboardBase;
