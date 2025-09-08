import React from 'react';
import { Link } from 'react-router-dom';
import {
  Box,
  Card,
  CardContent,
  Typography,
  Chip,
} from '@mui/material';

function ForecastCard({ forecast, isResolved = false}) {
  const resolution = forecast.resolution === '0' ? 'No' : 
                     forecast.resolution === '-' ? 'Ambiguous' : 'Yes';
  const resolutionColor = forecast.resolution === '0' ? 'secondary.main' : 
                          forecast.resolution === '-' ? 'warning.main' : 'primary.main';

  const formatDate = (dateString) => dateString.split('T')[0];
  return (
    <Card sx={{ 
      backgroundColor: 'background.paper',
      height: '100%',
      width: '100%',
      transition: 'transform 0.2s',
      '&:hover': {
        transform: 'translateY(-4px)',
      }
    }}>
      <CardContent>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 2 }}>
          <Typography
            component={Link}
            to={`/forecast/${forecast.id}`}
            sx={{
              color: 'primary.light',
              textDecoration: 'none',
              '&:hover': {
                color: 'primary.main',
              },
              flex: 1,
              mr: 2
            }}
            variant="h6"
          >
            {forecast.question}
          </Typography>
          
          {isResolved ? (
            <Chip
              label={resolution}
              sx={{
                backgroundColor: resolutionColor,
                color: 'primary.light',
                minWidth: '90px'
              }}
            />
          ) : (
            <Chip
              label={forecast.latestPoint ? 
                `${(forecast.latestPoint.point_forecast * 100).toFixed(1)}%` : 
                'Not forecasted'
              }
              sx={{
                backgroundColor: forecast.latestPoint ? 'primary.main' : 'secondary.main',
                color: 'primary.light',
                minWidth: '90px'
              }}
            />
          )}
        </Box>
        
        <Box sx={{ mt: 'auto' }}>
          <Typography sx={{ color: 'primary.light', opacity: 0.8 }}>
            Category: {forecast.category}
          </Typography>
          
          {isResolved ? (
            <>
              <Typography sx={{ color: 'primary.light', opacity: 0.8 }}>
                Resolved on: {formatDate(forecast.resolved ? forecast.resolved : forecast.created)}
              </Typography>
            </>
          ) : (
            <Typography sx={{ color: 'primary.light', opacity: 0.8 }}>
              Created: {formatDate(forecast.created)}
            </Typography>
          )}
        </Box>
      </CardContent>
    </Card>
  );
}

export default ForecastCard;