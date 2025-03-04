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
import { useForecastCombinedData } from '../services/hooks/useForecastCombinedData';
import { useSearchFilter } from '../services/hooks/useSearchFilter';
import ScoreDisplay from '../components/ScoreDisplay';
import { useAggregateScoresData } from '../services/hooks/useAggregateScoresData';

function ForecastPage() {
    const { category } = useParams();
    const location = useLocation();
    const [selectedUserId, setSelectedUserId] = useState('all');
    
    const categoryFilter = category || null;
    
    const listType = location.pathname.endsWith('/resolved') ? 'resolved' : 'open';
    
    const { combinedForecasts, loading, error } = useForecastCombinedData({category: categoryFilter, list_type: listType});
    
    const { handleSearch, sortedForecasts } = useSearchFilter(
      combinedForecasts, 
      { userId: selectedUserId, category: categoryFilter }
    );
    
    const { scores, scoresLoading } = useAggregateScoresData(categoryFilter, selectedUserId, false);
    
    const handleUserChange = (userId) => {
      setSelectedUserId(userId);
    };

    const formatDate = (dateString) => dateString.split('T')[0];

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
        }}
      >
        {error && (
          <Typography color="error">Error loading data: {error.message}</Typography>
        )}
        
        <Grid container spacing={3}>
          {/* Search and Header Section */}
          <Grid item xs={12}>
            <Box sx={{ mb: 4 }}>
              <Typography variant="h4" sx={{ color: 'primary.light', mb: 2 }}>
                {getPageTitle()}
              </Typography>
              
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
              <Grid item xs={12} md={6} lg={4} key={index}>
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
            sortedForecasts.map(forecast => (
              <Grid item xs={12} md={6} lg={4} key={forecast.id}>
                <ForecastCard 
                  forecast={forecast}  
                  isResolved={listType === 'resolved'}
                />
              </Grid>
            ))
          )}
        </Grid>
      </Box>
    </Box>
  );
}

export default ForecastPage;
