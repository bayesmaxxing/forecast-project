import React from 'react';
import {
  Card,
  CardContent,
  Typography,
  Skeleton,
  Box,
  Grid2
} from '@mui/material';
import ForecastCard from './ForecastCard';

const SKELETON_COUNT = 6;

function ForecastsList({ forecasts, loading, listType }) {
  if (loading) {
    return (
      <>
        {[...Array(SKELETON_COUNT)].map((_, index) => (
          <Grid2 xs={12} key={index}>
            <Card sx={{ 
              backgroundColor: 'background.paper',
              height: '100%',
            }}>
              <CardContent>
                <Skeleton variant="text" height={60} />
                <Skeleton variant="text" width="40%" />
                <Box sx={{ mt: 2 }}>
                  <Skeleton variant="text" width="30%" />
                  <Skeleton variant="text" width="40%" />
                </Box>
              </CardContent>
            </Card>
          </Grid2>
        ))}
      </>
    );
  }

  if (!Array.isArray(forecasts) || forecasts.length === 0) {
    return (
      <Grid2 xs={12}>
        <Typography>No forecasts available</Typography>
      </Grid2>
    );
  }

  return (
    <>
      {forecasts.map(forecast => (
        <Grid2 xs={12} key={forecast.id}>
          <ForecastCard 
            forecast={forecast}  
            isResolved={listType === 'resolved'}
          />
        </Grid2>
      ))}
    </>
  );
}

export default ForecastsList;