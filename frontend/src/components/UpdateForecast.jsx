import React, { useState } from 'react';
import { useParams } from 'react-router-dom';
import {
  Button,
  TextField,
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
  Typography,
  Alert,
  Snackbar,
  useTheme,
  Box
} from '@mui/material';
import { createPoint } from '../services/api/pointsService';

const UpdateForecast = ({ onSubmitSuccess }) => {
  let { id } = useParams();
  const theme = useTheme();
  const [updateData, setUpdateData] = useState({
    point_forecast: '',
    reason: '',
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
      reason: updateData.reason,
      forecast_id: parseInt(id)
    };

    try {
      const response = await createPoint(dataToSubmit.forecast_id, dataToSubmit.point_forecast, dataToSubmit.reason);
      
      if (response) {
        setSubmitStatus('Update added successfully');
        setUpdateData({
          point_forecast: '',
          reason: ''
        });
        handleClose(); // Close the dialog on success
        
        // Call the callback function to refetch data
        if (onSubmitSuccess) {
          onSubmitSuccess();
        }
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
    <>
      <Button 
        variant="contained" 
        color="primary"
        onClick={handleClickOpen}
        sx={{ mt: 2 }}
      >
        Update Forecast
      </Button>
      
      <Dialog open={open} onClose={handleClose} maxWidth="sm" fullWidth>
        <DialogTitle>Update Forecast</DialogTitle>
        <DialogContent>
          <DialogContentText sx={{ mb: 2 }}>
            Add a new forecast point with your latest prediction.
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
            Submit Update
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

export default UpdateForecast;
