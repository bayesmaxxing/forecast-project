import React, { useState, useEffect } from 'react';
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer
} from 'recharts';
import {
  Paper,
  Box,
  Typography,
  useTheme,
  useMediaQuery
} from '@mui/material';

function CalibrationChart() {
  const [resolutions, setResolutions] = useState([]);
  const [calibrationData, setCalibrationData] = useState([]);
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'));

  useEffect(() => {
    Promise.all([
      fetch('https://forecasting-389105.ey.r.appspot.com/forecasts?type=resolved', {
        headers: {
          "Accept": "application/json"
        }
      }),
      fetch('https://forecasting-389105.ey.r.appspot.com/forecast-points/latest', {
        headers: {
          "Accept": "application/json"
        }
      })
    ])
      .then(async ([resolved, latest]) => {
        const resolvedData = await resolved.json();
        const latestData = await latest.json();
        return [resolvedData, latestData];
      })
      .then(([resolvedData, latestData]) => {
        const combined = resolvedData.map(forecast => {
          const matchingPoint = latestData.find(point => point.forecast_id === forecast.id);
          return { ...forecast, latestPoint: matchingPoint || null };
        });
        setResolutions(combined);

        const calibrationData = createCalibrationData(combined);
        setCalibrationData(calibrationData);
      })
      .catch(error => console.error('Error fetching data: ', error));
  }, []);

  const createCalibrationData = (data) => {
    const bins = Array(10).fill().map((_, i) => ({
      min: i / 10,
      max: (i + 1) / 10,
      predictions: 0,
      occurrences: 0
    }));

    data.forEach(item => {
      if (item.latestPoint && item.latestPoint.point_forecast !== null) {
        const prediction = item.latestPoint.point_forecast;
        const binIndex = Math.min(Math.floor(prediction * 10), 9);
        if (item.resolution != "-") {
          bins[binIndex].predictions++;
        }
        if (item.resolution === "1") {
          bins[binIndex].occurrences++;
        }
      }
    });

    return bins.map(bin => ({
      binRange: `${(bin.min * 100).toFixed(0)}-${(bin.max * 100).toFixed(0)}%`,
      actualProbability: bin.predictions > 0 ? bin.occurrences / bin.predictions : 0,
      perfectCalibration: (bin.min + bin.max) / 2, 
      predictions: bin.predictions 
    }));
  };

  const CustomTooltip = ({ active, payload, label }) => {
    if (active && payload && payload.length > 0) {
      return (
        <Paper
          elevation={3}
          sx={{
            p: 1.5,
            backgroundColor: 'background.paper',
            border: `1px solid ${theme.palette.divider}`
          }}
        >
          <Typography variant="subtitle2">Bin: {label}</Typography>
          <Typography variant="body2" color="primary">
            Actual: {(payload[0].value * 100).toFixed(1)}%
          </Typography>
          <Typography variant="body2" color="secondary">
            Expected: {(payload[1].value * 100).toFixed(1)}%
          </Typography>
          <Typography variant="body2">
            Predictions: {payload[0].payload.predictions}
          </Typography>
        </Paper>
      );
    }
    return null;
  };

  return (
    <Paper
      elevation={3}
      sx={{
        p: { xs: 2, sm: 3 },
        m: { xs: 1, sm: 2 },
        bgcolor: 'background.paper',
        maxWidth: { xs: '100%', sm: '95%', md: '95%' },
        mx: 'auto'
      }}
    >
      <Typography
        variant="h6"
        align="center"
        sx={{ mb: 2 }}
      >
        Calibration Chart
      </Typography>

      <Box
        sx={{
          width: '100%',
          height: { xs: 300, sm: 450 },
          '& .recharts-tooltip-wrapper': {
            outline: 'none'
          }
        }}
      >
        <ResponsiveContainer>
          <LineChart
            data={calibrationData}
            margin={{
              top: 20,
              right: 30,
              left: isMobile ? 0 : 20,
              bottom: 20,
            }}
          >
            <CartesianGrid strokeDasharray="3 3" stroke={theme.palette.divider} />
            <XAxis
              dataKey="binRange"
              tick={{
                fill: theme.palette.text.primary,
                fontSize: isMobile ? 10 : 12
              }}
              label={{
                value: "Predicted Probability",
                position: "bottom",
                offset: 0,
                style: { textAnchor: 'middle' }
              }}
              angle={isMobile ? -45 : 0}
              textAnchor={isMobile ? "end" : "middle"}
              height={isMobile ? 60 : 40}
            />
            <YAxis
              tick={{
                fill: theme.palette.text.primary,
                fontSize: isMobile ? 10 : 12
              }}
              label={{
                value: "Actual Probability",
                angle: -90,
                position: 'insideLeft',
                style: { textAnchor: 'middle' }
              }}
              domain={[0, 1]}
              ticks={[0, 0.2, 0.4, 0.6, 0.8, 1.0]}
            />
            <Tooltip content={<CustomTooltip />} />
            <Line
              type="monotone"
              dataKey="actualProbability"
              stroke={theme.palette.primary.main}
              strokeWidth={2.5}
              dot={{ r: 4 }}
              activeDot={{ r: 6 }}
              name="My Calibration"
            />
            <Line
              type="monotone"
              dataKey="perfectCalibration"
              stroke={theme.palette.secondary.main}
              strokeWidth={1.5}
              strokeDasharray="5 5"
              dot={false}
              name="Perfect Calibration"
            />
          </LineChart>
        </ResponsiveContainer>
      </Box>
    </Paper>
  );
}

export default CalibrationChart;
