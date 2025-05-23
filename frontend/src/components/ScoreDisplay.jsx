import React from 'react';
import { 
  Paper, 
  Typography, 
  Skeleton, 
  Box, 
  Tooltip, 
  Stack,
} from '@mui/material';
import InfoIcon from '@mui/icons-material/Info';

// Score descriptions for tooltips
const scoreDescriptions = {
  brier: "The Brier score measures the accuracy of probabilistic predictions. Lower scores indicate better accuracy (0 is perfect).",
  base2log: "Base 2 logarithmic score measures prediction accuracy. Higher scores are better.",
  baseNlog: "Base N logarithmic score is similar to Base 2 but with a different base. Higher scores are better."
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
        
        </Stack>
      )}
    </Paper>
  );
}

export default ScoreDisplay;