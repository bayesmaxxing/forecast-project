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
      elevation={0}
      sx={{
        p: 3,
        backgroundColor: 'background.paper',
        mb: 2,
        borderRadius: 3,
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'space-between',
        border: '1px solid',
        borderColor: 'divider',
        background: 'linear-gradient(135deg, rgba(255, 107, 107, 0.03) 0%, rgba(255, 160, 122, 0.03) 100%)',
        transition: 'all 0.2s ease-in-out',
        '&:hover': {
          borderColor: 'primary.main',
          transform: 'translateY(-2px)',
          boxShadow: '0 8px 16px -4px rgba(255, 107, 107, 0.2)',
        }
      }}
    >
      {loading ? (
        <Skeleton width={200} height={24} />
      ) : (
        <Stack direction="row" spacing={2} alignItems="center" width="100%">
          <Box>
            <Typography
              variant="subtitle2"
              color="text.secondary"
              sx={{
                display: 'flex',
                alignItems: 'center',
                fontWeight: 600,
                textTransform: 'uppercase',
                letterSpacing: '0.05em',
                fontSize: '0.75rem',
                mb: 0.5
              }}
            >
              {scoreLabel}
              <Tooltip title={scoreDescriptions[type] || 'Score information'}>
                <InfoIcon fontSize="small" sx={{ ml: 1, opacity: 0.6, cursor: 'help' }} />
              </Tooltip>
            </Typography>
            <Typography
              variant="h4"
              sx={{
                background: 'linear-gradient(135deg, #FF6B6B 0%, #FFA07A 100%)',
                backgroundClip: 'text',
                WebkitBackgroundClip: 'text',
                WebkitTextFillColor: 'transparent',
                fontWeight: 700,
              }}
            >
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