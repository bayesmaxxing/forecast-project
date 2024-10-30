import React, { useState } from 'react';
import { useParams } from 'react-router-dom';
import {
  Paper,
  Box,
  TextField,
  Button,
  Typography,
  Alert,
  Snackbar,
  useTheme
} from '@mui/material';

const UpdateForecast = () => {
  let { id } = useParams();
  const theme = useTheme();
  const [updateData, setUpdateData] = useState({
    point_forecast: '',
    upper_ci: '',
    lower_ci: '',
    reason: '',
  });
  const [submitStatus, setSubmitStatus] = useState('');
  const [snackbarOpen, setSnackbarOpen] = useState(false);

  const handleChange = (e) => {
    const { name, value } = e.target;
    if (['point_forecast', 'upper_ci', 'lower_ci'].includes(name)) {
      const regex = /^-?\d*\.?\d*$/;
      if (value === '' || regex.test(value)) {
        setUpdateData(prevState => ({
          ...prevState,
          [name]: value
        }));
      }
    } else {
      setUpdateData(prevState => ({
        ...prevState,
        [name]: value
      }));
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setSubmitStatus('Submitting update...');
    setSnackbarOpen(true);

    const dataToSubmit = {
      ...updateData,
      point_forecast: parseFloat(updateData.point_forecast),
      upper_ci: parseFloat(updateData.upper_ci),
      lower_ci: parseFloat(updateData.lower_ci),
      reason: updateData.reason,
      forecast_id: parseInt(id)
    };

    try {
      const response = await fetch(`https://forecasting-389105.ey.r.appspot.com/forecast-points`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(dataToSubmit)
      });

      if (response.ok) {
        setSubmitStatus('Update added successfully');
        setUpdateData({
          point_forecast: '',
          upper_ci: '',
          lower_ci: '',
          reason: ''
        });
      } else {
        setSubmitStatus('Update failed. Please try again.');
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
        backgroundColor: theme.palette.background.paper
      }}
    >
      <Typography variant="h6" gutterBottom>
        Update Forecast
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
          label="Point Forecast"
          id="point_forecast"
          name="point_forecast"
          value={updateData.point_forecast}
          onChange={handleChange}
          required
          fullWidth
          size="small"
          type="text"
          helperText="Enter a number between 0 and 1"
        />

        <TextField
          label="Upper Confidence Interval"
          id="upper_ci"
          name="upper_ci"
          value={updateData.upper_ci}
          onChange={handleChange}
          required
          fullWidth
          size="small"
          type="text"
          helperText="Must be greater than point forecast"
        />

        <TextField
          label="Lower Confidence Interval"
          id="lower_ci"
          name="lower_ci"
          value={updateData.lower_ci}
          onChange={handleChange}
          required
          fullWidth
          size="small"
          type="text"
          helperText="Must be less than point forecast"
        />

        <TextField
          label="Reason for Update"
          id="reason"
          name="reason"
          value={updateData.reason}
          onChange={handleChange}
          required
          fullWidth
          multiline
          rows={3}
          helperText="Explain your reasoning for this update"
        />

        <Button 
          type="submit" 
          variant="contained" 
          color="primary"
          sx={{ mt: 2 }}
        >
          Submit Update
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

export default UpdateForecast;
