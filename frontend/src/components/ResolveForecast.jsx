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
  Alert,
  Snackbar,
  Box,
  RadioGroup,
  FormControl,
  FormLabel,
  FormControlLabel,
  Radio
} from '@mui/material';
import { resolveForecast } from '../services/api/forecastService';

const ResolveForecast = () => {
  let { id } = useParams();
  const [resolveData, setResolveData] = useState({
    resolution: '',
    comment: '',
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
    setResolveData(prevState => ({
      ...prevState,
      [name]: value
    }));
  };  

  const handleSubmit = async (e) => {
    e.preventDefault();
    setSubmitStatus('Submitting update...');
    setSnackbarOpen(true);

    const dataToSubmit = {
      ...resolveData,
      resolution: resolveData.resolution,
      comment: resolveData.comment,
      forecast_id: parseInt(id)
    };

    try {
      const response = await resolveForecast(dataToSubmit.forecast_id, dataToSubmit.resolution, dataToSubmit.comment);

      if (response.ok) {
        setSubmitStatus('Forecast resolved successfully');
        setResolveData({
          resolution: '',
          comment: ''
        });
        handleClose(); // Close the dialog on success
      } else {
        setSubmitStatus('Forecast resolution failed. Please try again.');
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
        Resolve Forecast
      </Button>
      
      <Dialog open={open} onClose={handleClose} maxWidth="sm" fullWidth>
        <DialogTitle>Resolve Forecast</DialogTitle>
        <DialogContent>
          <DialogContentText sx={{ mb: 2 }}>
            Resolve the forecast as Yes, No, or Ambiguous.
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
            <FormControl required>
          <FormLabel id="resolution-group-label">Resolution</FormLabel>
          <RadioGroup
            aria-labelledby="resolution-group-label"
            name="resolution"
            value={resolveData.resolution}
            onChange={handleChange}
          >
            <FormControlLabel 
              value="1" 
              control={<Radio />} 
              label="Resolved as Yes"
            />
            <FormControlLabel 
              value="0" 
              control={<Radio />} 
              label="Resolved as No"
            />
            <FormControlLabel 
              value="-" 
              control={<Radio />} 
              label="Resolved as Ambiguous"
            />
          </RadioGroup>
        </FormControl>

        <TextField
          label="Comment on Resolution"
          id="comment"
          name="comment"
          value={resolveData.comment}
          onChange={handleChange}
          required
          fullWidth
          multiline
          rows={3}
          helperText="Provide details about the resolution"
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
            Submit resolution 
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

export default ResolveForecast;