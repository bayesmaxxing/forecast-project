import React, { useState } from 'react';
import { useParams, useLocation } from 'react-router-dom';
import {
  Box,
  Typography,
  Grid2,
} from '@mui/material';
import Sidebar from '../components/Sidebar';
import SearchAndFilters from '../components/SearchAndFilters';
import ScoresSection from '../components/ScoresSection';
import ForecastsList from '../components/ForecastsList';
import AddForecast from '../components/AddForecast';
import { useForecastPageData } from '../services/hooks/useForecastPageData';

function ForecastPage() {
    const { category } = useParams();
    const location = useLocation();
    const isLoggedIn = localStorage.getItem('token') !== null;
    const [selectedUserId, setSelectedUserId] = useState('all');
    
    // Get the category and list type from the url
    const categoryFilter = category || null;
    const listType = location.pathname.endsWith('/resolved') ? 'resolved' : 'open';
    
    // Fetch and combine all data using the custom hook
    const {
      sortedForecasts,
      scores,
      loading,
      scoresLoading,
      error,
      handleSearch,
      refetchAllData
    } = useForecastPageData({ categoryFilter, listType, selectedUserId });
    
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
        {error && (
          <Typography color="error">Error loading data: {error?.message}</Typography>
        )}
        
        <Grid2 container spacing={3}>
          {/* Search and Header Section */}
          <Grid2 xs={12}>
            <Box sx={{ mb: 4 }}>
              <Grid2 container spacing={2} alignItems="center" sx={{ mb: 2 }}>
                <Grid2 xs={12} sm={isLoggedIn ? 6 : 12}>
                  <Typography variant="h4" sx={{ color: 'primary.light' }}>
                    {getPageTitle()}             
                  </Typography>
                </Grid2>
                {isLoggedIn && (
                  <Grid2 xs={12} sm={6} sx={{ display: 'flex', justifyContent: { xs: 'flex-start', sm: 'flex-end' } }}>
                    <AddForecast onSubmitSuccess={refetchAllData} />
                  </Grid2>
                )}
              </Grid2>
              <SearchAndFilters
                onSearch={handleSearch}
                selectedUserId={selectedUserId}
                onUserChange={handleUserChange}
              />
              
              {/* Score Display Section */}
              <ScoresSection 
                scores={scores}
                loading={loading || scoresLoading}
              />
            </Box>
          </Grid2>

          {/* Forecasts Grid */}
          <Grid2 xs={12}>
            <ForecastsList 
              forecasts={sortedForecasts}
              loading={loading}
              listType={listType}
            />
          </Grid2>
        </Grid2>
      </Box>
    </Box>
  );
}

export default ForecastPage;
