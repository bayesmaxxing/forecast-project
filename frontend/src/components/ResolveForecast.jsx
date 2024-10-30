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
  RadioGroup,
  FormControlLabel,
  Radio,
  FormControl,
  FormLabel,
  useTheme
} from '@mui/material';

const ResolveForecast = ({ forecastPoints }) => {
  let { id } = useParams();
  const theme = useTheme();
  const [resolveData, setResolveData] = useState({
    resolution: '',
    comment: '',
  });
  const [submitStatus, setSubmitStatus] = useState('');
  const [snackbarOpen, setSnackbarOpen] = useState(false);

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
      id: parseInt(id)
    };

    try {
      const response = await fetch(`https://forecasting-389105.ey.r.appspot.com/resolve/${id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(dataToSubmit)
      });

      if (response.ok) {
        setSubmitStatus('Forecast resolved successfully');
        setResolveData({
          resolution: '',
          comment: '',
        });
      } else {
        setSubmitStatus('Resolution failed. Please try again.');
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
        Resolve Forecast
      </Typography>

      <Box
        component="form"
        onSubmit={handleSubmit}
        sx={{
          display: 'flex',
          flexDirection: 'column',
          gap: 3
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

        <Button 
          type="submit" 
          variant="contained" 
          color="primary"
          disabled={!resolveData.resolution} // Disable if no resolution selected
          sx={{ mt: 2 }}
        >
          Resolve Forecast
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

export default ResolveForecast;
