import React, {useState, useEffect} from 'react';
import { useParams } from 'react-router-dom';
import { Link } from 'react-router-dom';
import {
  Box,
  Card,
  CardContent,
  Typography,
  TextField,
  Skeleton,
  Chip,
  Grid,
  Paper,
  InputAdornment,
  useTheme
} from '@mui/material';
import { Search as SearchIcon } from '@mui/icons-material';
import Sidebar from '../components/Sidebar';


function ForecastPage() {
    const [searchQuery, setsearchQuery] = useState('');
    const [combinedForecasts, setCombinedForecasts] = useState([]);
    const [scores, setScores] = useState([]);
    const [loading, setLoading] = useState(true);

    // set category based on URL
    let { category } = useParams()

    useEffect(() => {
       // Fetch the list of forecasts from the API
        Promise.all([
          fetch(`https://forecasting-389105.ey.r.appspot.com/forecasts?category=${category}&type=open`, {
            headers : {
              "Accept": "application/json"
            }
          }),
          fetch(`https://forecasting-389105.ey.r.appspot.com/forecast-points/latest`, {
            headers : {
              "Accept": "application/json"
            }
          }), 
          fetch(`https://forecasting-389105.ey.r.appspot.com/scores?category=${category}`, {
            headers : {
              "Accept": "application/json"
            }
          })
        ])
        .then(async ([forecastData, pointsData, scoresData]) => {
          const forecastDataJson = await forecastData.json();
          const pointsDataJson = await pointsData.json();
          const scoresDataJson = await scoresData.json();
          return [forecastDataJson, pointsDataJson, scoresDataJson];
        })
        .then(([forecastDataJson, pointsDataJson, scoresDataJson]) => {
          const combined = forecastDataJson.map(forecast => {
            const matchingPoint = pointsDataJson.find(point => point.forecast_id === forecast.id);
            return { ...forecast, latestPoint: matchingPoint || null};
          });
          setCombinedForecasts(combined)
          setScores(scoresDataJson)
          setLoading(false)
        })
        .catch(error => console.error('Error fetching data: ', error));
    }, [category]);
    
  const handleSearchChange = (e) => {
    setsearchQuery(e.target.value.toLowerCase());
  };

  const filteredForecasts = combinedForecasts.filter(forecast => 
    forecast.question.toLowerCase().includes(searchQuery) ||
    forecast.category.toLowerCase().includes(searchQuery) ||
    forecast.resolution_criteria.toLowerCase().includes(searchQuery)
    );

  const sortedForecasts = [...filteredForecasts].sort((a, b)=>{
    return b.id - a.id;
  });

  const formatDate = (dateString) => dateString.split('T')[0];

  return (
    <Box sx={{ p: 3, pt: 10 }}>
      <Grid container spacing={3}>
        {/* Search and Header Section */}
        <Grid item xs={12}>
          <Box sx={{ mb: 4 }}>
            <Typography variant="h4" sx={{ color: 'primary.light', mb: 2 }}>
              {category?.toUpperCase() || "CATEGORY NOT FOUND"}
            </Typography>
            
            <TextField
              fullWidth
              variant="outlined"
              placeholder="Search forecasts..."
              onChange={handleSearchChange}
              sx={{
                mb: 2,
                '& .MuiOutlinedInput-root': {
                  backgroundColor: 'background.paper',
                  '& fieldset': {
                    borderColor: 'primary.main',
                  },
                  '&:hover fieldset': {
                    borderColor: 'primary.light',
                  },
                },
                '& input': {
                  color: 'primary.light',
                }
              }}
              InputProps={{
                startAdornment: (
                  <InputAdornment position="start">
                    <SearchIcon sx={{ color: 'primary.main' }} />
                  </InputAdornment>
                ),
              }}
            />

            <Paper sx={{ p: 2, backgroundColor: 'background.paper', mb: 3 }}>
              {loading ? (
                <Skeleton width={200} height={24} />
              ) : (
                scores && scores.AggBrierScore > 0.0 ? (
                  <Typography variant="h6" sx={{ color: 'primary.light' }}>
                    Brier score: {(scores.AggBrierScore).toFixed(4)}
                  </Typography>
                ) : (
                  <Typography variant="h6" sx={{ color: 'primary.light' }}>
                    No Brier score available.
                  </Typography>
                )
              )}
            </Paper>
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
                        }
                      }}
                      variant="h6"
                    >
                      {forecast.question}
                    </Typography>
                    <Chip
                      label={forecast.latestPoint ? 
                        `${(forecast.latestPoint.point_forecast * 100).toFixed(1)}%` : 
                        'Not forecasted'
                      }
                      sx={{
                        backgroundColor: forecast.latestPoint ? 'primary.main' : 'secondary.main',
                        color: 'primary.light',
                        ml: 2,
                        minWidth: '90px'
                      }}
                    />
                  </Box>
                  <Box sx={{ mt: 'auto' }}>
                    <Typography sx={{ color: 'primary.light', opacity: 0.8 }}>
                      Category: {forecast.category}
                    </Typography>
                    <Typography sx={{ color: 'primary.light', opacity: 0.8 }}>
                      Created: {formatDate(forecast.created)}
                    </Typography>
                  </Box>
                </CardContent>
              </Card>
            </Grid>
          ))
        )}
      </Grid>
    </Box>
  );
};
  
  export default ForecastPage;
