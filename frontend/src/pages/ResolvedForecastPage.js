import React, {useState, useEffect} from 'react';
import { Link } from 'react-router-dom';
import {
  Box,
  Card,
  CardContent,
  Typography,
  Grid,
  Chip,
  Paper,
  Skeleton,
  useTheme
} from '@mui/material';
import Sidebar from '../components/Sidebar';


function ResolvedForecastPage() {
  const [forecasts, setForecasts] = useState([]);
  const [searchQuery, setsearchQuery] = useState('');
  const [loading, setLoading] = useState(true);
  const theme = useTheme();

  useEffect(() => {
      fetch(`https://forecasting-389105.ey.r.appspot.com/forecasts?type=resolved`, {
        headers : {
          "Accept" : "application/json"
        }
      })
      .then(response => response.json())
      .then(data => {
        setForecasts(data);
        setLoading(false);
      })
      .catch(error => {console.error('Error fetching data: ', error);
      setLoading(false);
  });
  }, []);

  const handleSearchChange = (e) => {
    setsearchQuery(e.target.value.toLowerCase());
  };

  const filteredForecasts = forecasts.filter(forecast => 
    forecast.question.toLowerCase().includes(searchQuery) ||
    forecast.category.toLowerCase().includes(searchQuery) ||
    forecast.resolution_criteria.toLowerCase().includes(searchQuery)
    );

  const sortedForecasts = [...filteredForecasts].sort((a, b)=>{
    return b.id - a.id;
  })

  const formatDate = (dateString) => dateString.split('T')[0];
  const getResolutionDetails = (resolution) => {
    switch (resolution) {
      case "1":
        return { text: "Yes", color: theme.palette.success.main };
      case "0":
        return { text: "No", color: theme.palette.error.main };
      default:
        return { text: "Ambiguous", color: theme.palette.warning.main };
    }
  };

  return (
    <Box sx={{ display: 'flex' }}>
      <Sidebar />
      <Box
        component="main"
        sx={{
          flexGrow: 1,
          p: 3,
          width: { sm: `calc(100%-240px)` },
          ml: { sm: '240px' },
          mt: { xs: '104px', sm: '64px' },
        }}
      >
        <Typography variant="h4" sx={{ color: 'primary.light', mb: 3 }}>
          RESOLVED QUESTIONS
        </Typography>

        <Grid container spacing={3}>
          {loading ? (
            [...Array(6)].map((_, index) => (
              <Grid item xs={12} md={6} lg={4} key={index}>
                <Card sx={{ backgroundColor: 'background.paper' }}>
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
                <Card sx={{
                  backgroundColor: 'background.paper',
                  height: '100%',
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
                      <Chip
                        label={getResolutionDetails(forecast.resolution).text}
                        sx={{
                          backgroundColor: getResolutionDetails(forecast.resolution).color,
                          color: 'primary.light',
                          minWidth: '90px'
                        }}
                      />
                    </Box>
                    
                    <Box sx={{ mt: 'auto' }}>
                      <Typography sx={{ color: 'primary.light', opacity: 0.8, mb: 1 }}>
                        Resolved on: {formatDate(forecast.resolved)}
                      </Typography>
                      {forecast.brier_score ? (
                        <Typography sx={{ color: 'primary.light', opacity: 0.8 }}>
                          Brier score: {forecast.brier_score.toFixed(3)}
                        </Typography>
                      ) : (
                        <Typography sx={{ color: 'primary.light', opacity: 0.8, fontStyle: 'italic' }}>
                          No scores on Ambiguous forecasts.
                        </Typography>
                      )}
                    </Box>
                  </CardContent>
                </Card>
              </Grid>
            ))
          )}
        </Grid>
      </Box>
    </Box>
  );
} 

export default ResolvedForecastPage;
