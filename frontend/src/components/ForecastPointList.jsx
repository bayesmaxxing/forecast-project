import React from 'react';
import {
  Typography,
  List,
  Box,
  Chip,
  Paper,
} from '@mui/material';
import ReactMarkdown from 'react-markdown';

function ForecastPointsList({ points, users }) {
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
                    label={users?.find(user => user.id === point.user_id)?.username || 'Unknown User'}
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
                  <Box sx={{
                    mt: 1,
                    '& p': {
                      margin: 0,
                      fontSize: '0.875rem',
                      color: 'text.secondary'
                    },
                    '& h1, & h2, & h3, & h4, & h5, & h6': {
                      fontSize: '1rem',
                      fontWeight: 'medium',
                      margin: '0.5rem 0 0.25rem 0'
                    },
                    '& ul, & ol': {
                      paddingLeft: '1.5rem',
                      margin: '0.25rem 0'
                    },
                    '& blockquote': {
                      borderLeft: 3,
                      borderColor: 'primary.main',
                      paddingLeft: '0.75rem',
                      margin: '0.25rem 0',
                      fontStyle: 'italic'
                    },
                    '& code': {
                      backgroundColor: 'grey.100',
                      padding: '0.125rem 0.25rem',
                      borderRadius: '0.25rem',
                      fontSize: '0.8125rem'
                    },
                    '& pre': {
                      backgroundColor: 'grey.100',
                      padding: '0.5rem',
                      borderRadius: '0.25rem',
                      overflow: 'auto',
                      margin: '0.25rem 0'
                    }
                  }}>
                    <ReactMarkdown>{point.reason}</ReactMarkdown>
                  </Box>
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