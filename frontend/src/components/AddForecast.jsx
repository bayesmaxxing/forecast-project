import React, { useState } from 'react';
import {
  Paper,
  Box,
  TextField,
  Button,
  Typography,
  Alert,
  Snackbar,
  useTheme,
  Autocomplete
} from '@mui/material';

const AddForecast = () => {
  const theme = useTheme();
  const [forecastData, setForecastData] = useState({
    question: '',
    short_question: '',
    category: '',
    resolution_criteria: ''
  });
  const [submitStatus, setSubmitStatus] = useState('');
  const [snackbarOpen, setSnackbarOpen] = useState(false);

  // Predefined categories - you can adjust these based on your needs
  const categories = [
    'ai',
    'economy',
    'politics',
    'x-risk',
    'personal',
    'sports',
    'technology',
    'society',
    'other'
  ];

  const handleChange = (e) => {
    const { name, value } = e.target;
    setForecastData(prevState => ({
      ...prevState,
      [name]: value
    }));
  };

  const handleCategoryChange = (event, newValue) => {
    setForecastData(prevState => ({
      ...prevState,
      category: newValue || ''
    }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setSubmitStatus('Submitting forecast...');
    setSnackbarOpen(true);

    try {
      const response = await fetch(`https://forecasting-389105.ey.r.appspot.com/forecasts`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(forecastData)
      });

      if (response.ok) {
        setSubmitStatus('Forecast added successfully');
        setForecastData({
          question: '',
          short_question: '',
          category: '',
          resolution_criteria: ''
        });
      } else {
        setSubmitStatus('Failed to add forecast. Please try again.');
      }
    } catch (error) {
      console.error('Error:', error);
      setSubmitStatus('An error occurred. Please try again.');
    }
  };

  const handleSnackbarClose = () => {
    setSnackbarOpen(false);
  };

  return (
    <Paper
      elevation={3}
      sx={{
        p: 3,
        mt: 2,
        backgroundColor: theme.palette.background.paper,
        maxWidth: 800,
        mx: 'auto'
      }}
    >
      <Typography variant="h6" gutterBottom>
        Add New Forecast
      </Typography>

      <Box
        component="form"
        onSubmit={handleSubmit}
        sx={{
          '& .MuiTextField-root': { mb: 2 },
          display: 'flex',
          flexDirection: 'column',
          gap: 2
        }}
      >
        <TextField
          label="Question"
          id="question"
          name="question"
          value={forecastData.question}
          onChange={handleChange}
          required
          fullWidth
          helperText="Enter the full question text"
        />

        <TextField
          label="Short Question"
          id="short_question"
          name="short_question"
          value={forecastData.short_question}
          onChange={handleChange}
          required
          fullWidth
          helperText="Enter a brief version of the question"
        />

        <Autocomplete
          freeSolo
          options={categories}
          value={forecastData.category}
          onChange={handleCategoryChange}
          renderInput={(params) => (
            <TextField
              {...params}
              label="Category"
              name="category"
              required
              fullWidth
              helperText="Select or enter a category"
              onChange={(e) => {
                handleChange(e);
              }}
            />
          )}
        />

        <TextField
          label="Resolution Criteria"
          id="resolution_criteria"
          name="resolution_criteria"
          value={forecastData.resolution_criteria}
          onChange={handleChange}
          required
          fullWidth
          multiline
          rows={4}
          helperText="Specify how this forecast will be resolved"
        />

        <Button 
          type="submit" 
          variant="contained" 
          color="primary"
          sx={{ mt: 2 }}
        >
          Add Forecast
        </Button>
      </Box>

      <Snackbar
        open={snackbarOpen}
        autoHideDuration={6000}
        onClose={handleSnackbarClose}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
      >
        <Alert 
          onClose={handleSnackbarClose} 
          severity={submitStatus.includes('success') ? 'success' : submitStatus.includes('error') ? 'error' : 'info'}
          sx={{ width: '100%' }}
        >
          {submitStatus}
        </Alert>
      </Snackbar>
    </Paper>
  );
};

export default AddForecast;
