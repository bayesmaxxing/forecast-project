import React from 'react';
import {
  Typography,
  List,
  ListItem,
  ListItemText,
  Divider,
  Avatar,
  Box,
  Chip,
  Paper,
} from '@mui/material';

function ForecastPointsList({ points }) {
  const formatDate = (dateString) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', { 
      year: 'numeric', 
      month: 'short', 
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };
  
  // Assume points now includes user information
  const sortedPoints = [...points].sort((a, b) => new Date(b.created) - new Date(a.created));

  if (!points || points.length === 0) {
    return <Typography variant="body1">No forecast points available yet.</Typography>;
  }

  return (
    <Box sx={{ position: 'relative', ml: 4 }}>
      {/* Timeline vertical line */}
      <Box sx={{ 
        position: 'absolute', 
        left: 15, 
        top: 0, 
        bottom: 0, 
        width: 2, 
        bgcolor: 'grey.300',
        zIndex: 0 
      }} />
      
      <List sx={{ width: '100%' }}>
        {sortedPoints.map((point, index) => {
          return (
            <Box key={point.id || index} sx={{ mb: 3, position: 'relative', zIndex: 1 }}>
              {/* Timeline node */}
              <Box sx={{ 
                position: 'absolute', 
                left: -15, 
                width: 14, 
                height: 14, 
                borderRadius: '50%', 
                bgcolor: 'grey.300',
                border: '2px solid white',
                zIndex: 2
              }} />
              
              <Paper elevation={1} sx={{ ml: 3, p: 2, borderRadius: 2 }}>
                {/* User chip/avatar */}
                <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
                  <Chip 
                    avatar={<Avatar sx={{ bgcolor: 'grey.300' }}>{String(point.user_id || 'User').charAt(0)}</Avatar>}
                    label={point.user_id || 'Unknown User'}
                    size="small"
                    sx={{ mr: 1 }}
                  />
                  <Typography variant="caption" color="text.secondary">
                    {formatDate(point.created)}
                  </Typography>
                </Box>
                
                <Typography variant="h6" sx={{ fontWeight: 'medium' }}>
                  {(point.point_forecast * 100).toFixed(1)}%
                </Typography>
                
                {point.reason && (
                  <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
                    {point.reason}
                  </Typography>
                )}
              </Paper>
            </Box>
          );
        })}
      </List>
    </Box>
  );
}

export default ForecastPointsList;