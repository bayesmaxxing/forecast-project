import React from 'react';
import { 
  Paper, 
  Typography, 
  Skeleton, 
  Box, 
  Tooltip, 
  Stack,
  Chip
} from '@mui/material';
import InfoIcon from '@mui/icons-material/Info';

// Score descriptions for tooltips
const scoreDescriptions = {
  brier: "The Brier score measures the accuracy of probabilistic predictions. Lower scores indicate better accuracy (0 is perfect).",
  base2log: "Base 2 logarithmic score measures prediction accuracy. Higher scores are better.",
  baseNlog: "Base N logarithmic score is similar to Base 2 but with a different base. Higher scores are better."
};

// Color coding based on score ranges (customize these thresholds as needed)
const getScoreColor = (type, value) => {
  if (value === null) return 'default';
  
  switch(type) {
    case 'brier':
      // For Brier score, lower is better
      if (value < 0.1) return 'success';
      if (value < 0.2) return 'primary';
      if (value < 0.3) return 'warning';
      return 'error';
    case 'base2log':
    case 'baseNlog':
      // For log scores, higher is better
      if (value > 0.8) return 'success';
      if (value > 0.5) return 'primary';
      if (value > 0.3) return 'warning';
      return 'error';
    default:
      return 'primary';
  }
};

function ScoreDisplay({ 
  type = 'brier', 
  value = null, 
  loading = false, 
  label = null,
  decimals = 4
}) {
  // Default labels if not provided
  const scoreLabel = label || {
    brier: 'Brier Score',
    base2log: 'Base 2 Log Score',
    baseNlog: 'Base N Log Score'
  }[type] || 'Score';
  
  // Format the score value appropriately
  const formattedValue = value !== null ? value.toFixed(decimals) : 'N/A';
  
  // Determine color based on score type and value
  const color = getScoreColor(type, value);
  
  return (
    <Paper 
      elevation={2} 
      sx={{ 
        p: 2, 
        backgroundColor: 'background.paper', 
        mb: 2,
        borderRadius: 2,
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'space-between'
      }}
    >
      {loading ? (
        <Skeleton width={200} height={24} />
      ) : (
        <Stack direction="row" spacing={2} alignItems="center" width="100%">
          <Box>
            <Typography variant="subtitle1" color="text.secondary" sx={{ display: 'flex', alignItems: 'center' }}>
              {scoreLabel}
              <Tooltip title={scoreDescriptions[type] || 'Score information'}>
                <InfoIcon fontSize="small" sx={{ ml: 1, opacity: 0.7 }} />
              </Tooltip>
            </Typography>
            <Typography variant="h5" color="primary.light" fontWeight="medium">
              {formattedValue}
            </Typography>
          </Box>
          
          <Box flexGrow={1} />
          
          {value !== null && (
            <Chip 
              label={type === 'brier' ? 'Lower is better' : 'Higher is better'} 
              size="small"
              color={color}
              variant="outlined"
              sx={{ fontWeight: 'medium' }}
            />
          )}
        </Stack>
      )}
    </Paper>
  );
}

export default ScoreDisplay;