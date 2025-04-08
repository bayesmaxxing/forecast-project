import React, { useState } from 'react';
import { useParams, useLocation } from 'react-router-dom';
import {
  Box,
  Card,
  CardContent,
  Typography,
  Skeleton,
  Paper,
  Grid,
} from '@mui/material';
import Sidebar from '../components/Sidebar';
import SearchBar from '../components/SearchBar';
import ForecastCard from '../components/ForecastCard';
import UserSelector from '../components/UserSelector';
import { useForecastList } from '../services/hooks/useForecastList';
import { usePointsData } from '../services/hooks/usePointsData';
import { useSearchFilter } from '../services/hooks/useSearchFilter';
import ScoreDisplay from '../components/ScoreDisplay';
import { useAggregateScoresData } from '../services/hooks/useAggregateScoresData';
import AddForecast from '../components/AddForecast';

function ForecastPage() {
    const { category } = useParams();
    const location = useLocation();
    const isLoggedIn = localStorage.getItem('token') !== null;
    const [selectedUserId, setSelectedUserId] = useState('all');
    
    // Get the category and list type from the url
    const categoryFilter = category || null;
    const listType = location.pathname.endsWith('/resolved') ? 'resolved' : 'open';
    
    // Fetch the data
    const { forecasts = [], loading, error, refetchForecasts } = useForecastList({
      category: categoryFilter, 
      list_type: listType
    });
    const { points = [], loading: pointsLoading, error: pointsError, refetchPoints } = usePointsData({
      userId: selectedUserId, 
      useLatestPoints: true, 
      useOrderedEndpoint: false
    });

    // Function to refetch all relevant data
    const refetchAllData = () => {
      refetchForecasts();
      refetchPoints();
    };

    // Combine the forecasts and points
    const combined = Array.isArray(forecasts) 
      ? forecasts.map(forecast => {
          const matchingPoint = points.find(point => point.forecast_id === forecast.id);
          return { ...forecast, latestPoint: matchingPoint || null};
        })
      : [];
    
    const { handleSearch, sortedForecasts } = useSearchFilter(
      combined, 
      { userId: selectedUserId, category: categoryFilter }
    );
    
    const { scores, scoresLoading } = useAggregateScoresData(categoryFilter, selectedUserId, false);
    
    const handleUserChange = (userId) => {
      setSelectedUserId(userId);
    };

    const getPageTitle = () => {
      if (!category) return "ALL QUESTIONS";
      return `${category.toUpperCase()} QUESTIONS`;
    };

    return (
    <Box sx={{ display: 'flex' }}>
      <Sidebar />
      <Box
        component="main"
        sx={{
          flexGrow: 1,
          p: 3,
          width: { sm: `calc(100% - 240px)` },
          mt: '64px',
          maxWidth: '1000px',
          mx: 'auto'                   
        }}
      >
        {(error || pointsError) && (
          <Typography color="error">Error loading data: {error?.message || pointsError?.message}</Typography>
        )}
        
        <Grid container spacing={3}>
          {/* Search and Header Section */}
          <Grid item xs={12}>
            <Box sx={{ mb: 4 }}>
              <Grid container spacing={2} alignItems="center" sx={{ mb: 2 }}>
                <Grid item xs={12} sm={isLoggedIn ? 6 : 12}>
                  <Typography variant="h4" sx={{ color: 'primary.light' }}>
                    {getPageTitle()}             
                  </Typography>
                </Grid>
                {isLoggedIn && (
                  <Grid item xs={12} sm={6} sx={{ display: 'flex', justifyContent: { xs: 'flex-start', sm: 'flex-end' } }}>
                    <AddForecast onSubmitSuccess={refetchAllData} />
                  </Grid>
                )}
              </Grid>
              <Grid container spacing={2} alignItems="center">
                <Grid item xs={12} md={8}>
                  <SearchBar 
                    onSearch={handleSearch}
                    placeholder="Search forecasts..."
                  />
                </Grid>
                <Grid item xs={12} md={4}>
                  <UserSelector 
                    onUserChange={handleUserChange}
                    selectedUserId={selectedUserId}
                  />
                </Grid>
              </Grid>
              
              {/* Score Display Section */}
              <Grid container spacing={2} sx={{ mt: 2 }}>
                <Grid item xs={12} md={4}>
                  <ScoreDisplay
                    type="brier"
                    value={scores?.brier_score || null}
                    loading={loading || scoresLoading}
                  />
                </Grid>
                <Grid item xs={12} md={4}>
                  <ScoreDisplay
                    type="base2log"
                    value={scores?.log2_score || null}
                    loading={loading || scoresLoading}
                  />
                </Grid>
                <Grid item xs={12} md={4}>
                  <ScoreDisplay
                    type="baseNlog"
                    value={scores?.logn_score || null}
                    loading={loading || scoresLoading}
                  />
                </Grid>
              </Grid>
            </Box>
          </Grid>

          {/* Forecasts Grid */}
          {loading ? (
            [...Array(6)].map((_, index) => (
              <Grid item xs={12} key={index}>
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
              </Grid>
            ))
          ) : (
            Array.isArray(sortedForecasts) ? sortedForecasts.map(forecast => (
              <Grid item xs={12} key={forecast.id}>
                <ForecastCard 
                  forecast={forecast}  
                  isResolved={listType === 'resolved'}
                />
              </Grid>
            )) : (
              <Grid item xs={12}>
                <Typography>No forecasts available</Typography>
              </Grid>
            )
          )}
        </Grid>
      </Box>
    </Box>
  );
}

export default ForecastPage;
