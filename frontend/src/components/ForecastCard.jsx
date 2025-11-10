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
    <Card
      elevation={0}
      sx={{
        backgroundColor: 'background.paper',
        height: '100%',
        width: '100%',
        border: '1px solid',
        borderColor: 'divider',
        borderRadius: 3,
        transition: 'all 0.3s cubic-bezier(0.4, 0, 0.2, 1)',
        position: 'relative',
        overflow: 'hidden',
        '&::before': {
          content: '""',
          position: 'absolute',
          top: 0,
          left: 0,
          right: 0,
          height: '3px',
          background: 'linear-gradient(90deg, #FF6B6B, #FFA07A)',
          transform: 'scaleX(0)',
          transformOrigin: 'left',
          transition: 'transform 0.3s ease-in-out',
        },
        '&:hover': {
          transform: 'translateY(-4px)',
          borderColor: 'primary.main',
          boxShadow: '0 12px 24px -8px rgba(255, 107, 107, 0.25)',
          '&::before': {
            transform: 'scaleX(1)',
          },
        },
      }}
    >
      <CardContent sx={{ p: 3 }}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 2, gap: 2 }}>
          <Typography
            component={Link}
            to={`/forecast/${forecast.id}`}
            sx={{
              color: 'text.primary',
              textDecoration: 'none',
              fontWeight: 600,
              fontSize: '1.1rem',
              transition: 'color 0.2s ease-in-out',
              '&:hover': {
                color: 'primary.main',
              },
              flex: 1,
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
                color: 'white',
                minWidth: '90px',
                fontWeight: 600,
                boxShadow: '0 2px 8px rgba(0, 0, 0, 0.15)',
              }}
            />
          ) : (
            <Chip
              label={forecast.latestPoint ?
                `${(forecast.latestPoint.point_forecast * 100).toFixed(1)}%` :
                'Not forecasted'
              }
              sx={{
                background: forecast.latestPoint
                  ? 'linear-gradient(135deg, #FF6B6B, #FFA07A)'
                  : 'secondary.main',
                color: 'white',
                minWidth: '90px',
                fontWeight: 600,
                boxShadow: '0 2px 8px rgba(0, 0, 0, 0.15)',
              }}
            />
          )}
        </Box>

        <Box sx={{ mt: 'auto', display: 'flex', flexDirection: 'column', gap: 0.5 }}>
          <Typography
            sx={{
              color: 'text.secondary',
              fontSize: '0.875rem',
              display: 'flex',
              alignItems: 'center',
              gap: 0.5,
            }}
          >
            <Box
              component="span"
              sx={{
                width: 6,
                height: 6,
                borderRadius: '50%',
                backgroundColor: 'primary.main',
                display: 'inline-block',
              }}
            />
            {forecast.category}
          </Typography>

          {isResolved ? (
            <Typography sx={{ color: 'text.secondary', fontSize: '0.875rem' }}>
              Resolved: {formatDate(forecast.resolved ? forecast.resolved : forecast.created)}
            </Typography>
          ) : (
            <Typography sx={{ color: 'text.secondary', fontSize: '0.875rem' }}>
              Created: {formatDate(forecast.created)}
            </Typography>
          )}
          {forecast.closing_date && (
            <Typography sx={{ color: 'text.secondary', fontSize: '0.875rem' }}>
              Closing on: {formatDate(forecast.closing_date)}
            </Typography>
          )}
        </Box>
      </CardContent>
    </Card>
  );
}

export default ForecastCard;