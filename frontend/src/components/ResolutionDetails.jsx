import React from 'react';
import { Typography, Chip, Box } from '@mui/material';

function ResolutionDetails({ forecast, score }) {
  if (!forecast || forecast.resolved == null) {
    return null;
  }

  const formatDate = (dateString) => dateString.split('T')[0];
  const resolution = forecast.resolution === "1" ? "Yes" :
    forecast.resolution === "0" ? "No" : "Ambiguous";

  const getResolutionColor = (res) => {
    switch (res) {
      case "Yes": return "success";
      case "No": return "error";
      default: return "warning";
    }
  };
  
  return (
    <>
      <Box sx={{ mb: 2 }}>
        <Chip
          label={`Resolved as: ${resolution}`}
          color={getResolutionColor(resolution)}
        />
      </Box>
      
      <Typography variant="h6" gutterBottom>Resolution Details</Typography>
      <Typography variant="body1" paragraph>
        Resolved on {formatDate(forecast.resolved)}
        {score ? 
          ` with a Brier score of ${score.brier_score ? (score.brier_score).toFixed(4) : 0}.` :
          `.`
        }
      </Typography>
      {forecast.comment && (
        <Typography variant="body1" paragraph>
          <strong>Comment:</strong> {forecast.comment}
        </Typography>
      )}
    </>
  );
}

export default ResolutionDetails;