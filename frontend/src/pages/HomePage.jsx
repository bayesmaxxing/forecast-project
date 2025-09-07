import React from 'react';
import { Link } from 'react-router-dom';
import SummaryScores from '../components/SummaryScores';
import {
  Container,
  Typography,
  Box,
  Grid2,
  useTheme
} from '@mui/material';

function HomePage() {
  const theme = useTheme();

  return (
    <Container 
      maxWidth="lg"
      sx={{
        mt: { xs: 6, sm: 8 }, // Add top margin to prevent header overlap
        mb: 4,
        '& a': {
          textDecoration: 'none',
          color: theme.palette.primary.main, // Make all links use primary color
          '&:hover': {
            textDecoration: 'underline'
          }
        }
      }}
    >
      <Box component="section" sx={{ mb: 0 , mt: { xs: 6, sm: 8 } }}>
        <Typography variant="h3" component="h1" gutterBottom>
          Forecasting to understand Reality
        </Typography>
        <Typography variant="body1">
          I forecast to improve my models of the world. On this website, I'll display my current and previous forecasts along with
          my track record.
        </Typography>
      </Box>

      <Box component="section" sx={{ mb: 6 }}>
        <Typography variant="h3" component="h1" gutterBottom>
          Scores
        </Typography>
        <Typography variant="body1" sx={{ mb: 2 }}>
          Forecasts are scored on their accuracy. The closer each score is to 0, the better. For more information, see
          <Link to="/faq"> FAQ</Link>. Click on a datapoint to see information about the forecast.
        </Typography>
        <SummaryScores />
      </Box>

      <Typography variant="h4" component="h2" gutterBottom>
        Focus Areas
      </Typography>

      <Grid2 container spacing={4} sx={{ mb: 4 }}>
        <Grid2 xs={12} sm={6} md={4}>
          <Typography variant="h5" component="h3" gutterBottom>
            <Link to="/questions/category/ai">
              Artificial Intelligence
            </Link>
          </Typography>
          <Typography variant="body1" paragraph>
            Current and future AI models have the potential to radically change society and humanity for both better and worse. These forecasts
            explore future AI technology and impacts on humanity.
          </Typography>
          <Link to="/questions/category/ai">See AI Forecasts →</Link>
        </Grid2>

        <Grid2 xs={12} sm={6} md={4}>
          <Typography variant="h5" component="h3" gutterBottom>
            <Link to="/questions/category/economy">
              Economy
            </Link>
          </Typography>
          <Typography variant="body1" paragraph>
            The economy affects peoples lives daily. These forecasts model economic development to help aid economic decision-making.
          </Typography>
          <Link to="/questions/category/economy">See Economy Forecasts →</Link>
        </Grid2>

        <Grid2 xs={12} sm={6} md={4}>
          <Typography variant="h5" component="h3" gutterBottom>
            <Link to="/questions/category/politics">
              Politics
            </Link>
          </Typography>
          <Typography variant="body1" paragraph>
            Political developments and changes affect both country and world developments. These forecasts anticipate political developments
            to understand future changes.
          </Typography>
          <Link to="/questions/category/politics">See Politics Forecasts →</Link>
        </Grid2>

        <Grid2 xs={12} sm={6} md={4}>
          <Typography variant="h5" component="h3" gutterBottom>
            <Link to="/questions/category/x-risk">
              Existential Risks
            </Link>
          </Typography>
          <Typography variant="body1" paragraph>
            As humanity becomes more and more technologically mature, our technologies are becoming strong enough to kill us all.
            <Link to="/questions/category/nuclear"> Nuclear weapons</Link>,{' '}
            <Link to="/questions/category/nuclear"> Artificial Intelligence</Link> and Biological weapons are examples of such technologies.
          </Typography>
          <Link to="/questions/category/x-risk">See X-Risk Forecasts →</Link>
        </Grid2>

        <Grid2 xs={12} sm={6} md={4}>
          <Typography variant="h5" component="h3" gutterBottom>
            <Link to="/questions">
              Other Categories
            </Link>
          </Typography>
          <Typography variant="body1" paragraph>
            Apart from the focus areas, I also forecast personal questions, Sweden-specific topics, Sports, among many other types of questions.
          </Typography>
          <Link to="/questions">See All Forecasts →</Link>
        </Grid2>
      </Grid2>
    </Container>
  );
}

export default HomePage;
