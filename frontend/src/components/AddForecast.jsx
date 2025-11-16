import React, { useState } from 'react';
import {
  Box,
  TextField,
  Button,
  Alert,
  Snackbar,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogContentText,
  DialogActions
} from '@mui/material';
import { createForecast } from '../services/api/forecastService';

const AddForecast = ({ onSubmitSuccess }) => {
  const [forecastData, setForecastData] = useState({
    question: '',
    category: '',
    resolution_criteria: '',
    closing_date: ''
  });
  const [submitStatus, setSubmitStatus] = useState('');
  const [snackbarOpen, setSnackbarOpen] = useState(false);
  const [open, setOpen] = useState(false);

  const handleClickOpen = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
  };

  const handleChange = (e) => {
    const { name, value } = e.target;
    setForecastData(prevState => ({
      ...prevState,
      [name]: value
    }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setSubmitStatus('Submitting forecast...');
    setSnackbarOpen(true);

    try {
      const payload = {
        ...forecastData,
        closing_date: forecastData.closing_date
          ? `${forecastData.closing_date}T00:00:00Z`
          : forecastData.closing_date
      };

      const response = await createForecast(payload);

      if (response==='forecast created') {
        setSubmitStatus('Forecast added successfully');
        setForecastData({
          question: '',
          category: '',
          resolution_criteria: '',
          closing_date: ''
        });
        handleClose();
        
        if (onSubmitSuccess) {
          onSubmitSuccess();
        }
      } else {
        setSubmitStatus('Failed to create forecast. Please try again.');
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
    <>
      <Button 
        variant="contained" 
        color="primary"
        onClick={handleClickOpen}
        sx={{ mt: 2 }}
      >
        Create Forecast
      </Button>
      
      <Dialog open={open} onClose={handleClose} maxWidth="sm" fullWidth>
        <DialogTitle>Create forecast question</DialogTitle>
        <DialogContent>
          <DialogContentText sx={{ mb: 2 }}>
            Create a new forecast question.
          </DialogContentText>
          
          <Box
            component="form"
            onSubmit={handleSubmit}
            sx={{
              '& .MuiTextField-root': { mb: 2 },
              display: 'flex',
              flexDirection: 'column',
              gap: 2,
              mt: 1
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
              size="small"
              type="text"
              helperText="Enter the forecast question"
            />
            <TextField
              label="Category"
              id="category"
              name="category"
              value={forecastData.category}
              onChange={handleChange}
              required
              fullWidth
              size="small"
              type="text"
              helperText="Enter the forecast category"
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
              rows={3}
              helperText="Enter the resolution criteria"
            />
            <TextField
              label="Closing Date"
              id="closing_date"
              name="closing_date"
              type="date"
              value={forecastData.closing_date}
              onChange={handleChange}
              required
              fullWidth
              size="small"
              InputLabelProps={{
                shrink: true,
              }}
              helperText="Select the date when the forecast closes"
            />
          </Box>
        </DialogContent>
        <DialogActions sx={{ px: 3, pb: 3 }}>
          <Button onClick={handleClose} color="secondary">
            Cancel
          </Button>
          <Button 
            onClick={handleSubmit} 
            variant="contained" 
            color="primary"
          >
            Submit Forecast
          </Button>
        </DialogActions>
      </Dialog>

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
    </>
  );
};

export default AddForecast;
