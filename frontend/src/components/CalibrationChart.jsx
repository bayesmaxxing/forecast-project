import React from 'react';
import {
  ComposedChart,
  Scatter,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from 'recharts';
import {
  Paper,
  Box,
  Typography,
  useTheme,
  CircularProgress,
  Alert
} from '@mui/material';
import { useCalibration } from '../services/hooks/useCalibration';

function CalibrationChart({ userId = 'all', dateRange = null, category = null }) {
  const theme = useTheme();

  const { calibration, loading, error } = useCalibration({
    type: 'overall',
    userId,
    category,
    dateRange,
  });

  // Transform data for the scatter chart - combine with perfect calibration line
  const chartData = React.useMemo(() => {
    const buckets = calibration?.buckets || [];

    // Create data points for both the scatter and the reference line
    // We'll use the same data array but with different keys
    const points = buckets.map(bucket => ({
      avgPrediction: bucket.avg_prediction * 100,
      actualRate: bucket.actual_rate * 100,
      predictionCount: bucket.prediction_count,
      bucketStart: bucket.bucket_start * 100,
      bucketEnd: bucket.bucket_end * 100,
    }));

    // Add perfect calibration line points
    const linePoints = [
      { avgPrediction: 0, perfectCalibration: 0 },
      { avgPrediction: 100, perfectCalibration: 100 }
    ];

    // Merge - scatter points get actualRate, line points get perfectCalibration
    return [...linePoints, ...points];
  }, [calibration]);

  const scatterData = calibration?.buckets?.map(bucket => ({
    avgPrediction: bucket.avg_prediction * 100,
    actualRate: bucket.actual_rate * 100,
    predictionCount: bucket.prediction_count,
    bucketStart: bucket.bucket_start * 100,
    bucketEnd: bucket.bucket_end * 100,
  })) || [];

  const CustomTooltip = ({ active, payload }) => {
    if (active && payload && payload.length) {
      const data = payload[0].payload;
      // Skip tooltip for the reference line points
      if (data.perfectCalibration !== undefined && data.actualRate === undefined) {
        return null;
      }
      return (
        <Box
          sx={{
            backgroundColor: theme.palette.background.paper,
            border: `1px solid ${theme.palette.divider}`,
            borderRadius: 1,
            p: 1.5,
            boxShadow: theme.shadows[2],
          }}
        >
          <Typography variant="body2" sx={{ fontWeight: 600, mb: 0.5 }}>
            Bucket: {data.bucketStart?.toFixed(0)}% - {data.bucketEnd?.toFixed(0)}%
          </Typography>
          <Typography variant="body2">
            Avg Prediction: {data.avgPrediction?.toFixed(1)}%
          </Typography>
          <Typography variant="body2">
            Actual Rate: {data.actualRate?.toFixed(1)}%
          </Typography>
          <Typography variant="body2" color="text.secondary">
            Predictions: {data.predictionCount}
          </Typography>
        </Box>
      );
    }
    return null;
  };

  if (loading) {
    return (
      <Paper
        elevation={0}
        sx={{
          p: 2.5,
          bgcolor: 'background.paper',
          height: '100%',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          minHeight: 300,
        }}
      >
        <CircularProgress />
      </Paper>
    );
  }

  if (error) {
    return (
      <Paper
        elevation={0}
        sx={{
          p: 2.5,
          bgcolor: 'background.paper',
          height: '100%',
        }}
      >
        <Alert severity="error">{error}</Alert>
      </Paper>
    );
  }

  if (!scatterData.length) {
    return (
      <Paper
        elevation={0}
        sx={{
          p: 2.5,
          bgcolor: 'background.paper',
          height: '100%',
        }}
      >
        <Typography variant="h6" gutterBottom>
          Calibration
        </Typography>
        <Typography color="text.secondary">
          No calibration data available for the selected filters.
        </Typography>
      </Paper>
    );
  }

  return (
    <Paper
      elevation={0}
      sx={{
        p: 2.5,
        bgcolor: 'background.paper',
        height: '100%',
        display: 'flex',
        flexDirection: 'column',
      }}
    >
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
        <Typography variant="h6">
          Calibration
        </Typography>
        <Typography variant="body2" color="text.secondary">
          {calibration?.total_predictions} predictions across {calibration?.total_forecasts} forecasts
        </Typography>
      </Box>

      <Box sx={{ height: 300 }}>
        <ResponsiveContainer width="100%" height="100%">
          <ComposedChart
            data={chartData}
            margin={{ top: 10, right: 20, left: 0, bottom: 20 }}
          >
            <CartesianGrid strokeDasharray="3 3" stroke={theme.palette.divider} />
            <XAxis
              type="number"
              dataKey="avgPrediction"
              domain={[0, 100]}
              tickFormatter={(value) => `${value}%`}
              tick={{ fill: theme.palette.text.primary, fontSize: 12 }}
              label={{
                value: 'Predicted Probability',
                position: 'bottom',
                offset: 0,
                style: { fill: theme.palette.text.secondary, fontSize: 12 }
              }}
            />
            <YAxis
              type="number"
              domain={[0, 100]}
              tickFormatter={(value) => `${value}%`}
              tick={{ fill: theme.palette.text.primary, fontSize: 12 }}
              label={{
                value: 'Actual Rate',
                angle: -90,
                position: 'insideLeft',
                style: { fill: theme.palette.text.secondary, fontSize: 12, textAnchor: 'middle' }
              }}
            />
            <Tooltip content={<CustomTooltip />} />
            {/* Perfect calibration reference line */}
            <Line
              type="linear"
              dataKey="perfectCalibration"
              stroke={theme.palette.text.disabled}
              strokeDasharray="5 5"
              strokeWidth={2}
              dot={false}
              activeDot={false}
              legendType="none"
            />
            <Scatter
              name="Calibration"
              dataKey="actualRate"
              data={scatterData}
              fill={theme.palette.primary.main}
            />
          </ComposedChart>
        </ResponsiveContainer>
      </Box>

      <Typography variant="caption" color="text.secondary" sx={{ mt: 1, textAlign: 'center' }}>
        Points on the dashed line indicate perfect calibration. Point size reflects prediction count.
      </Typography>
    </Paper>
  );
}

export default CalibrationChart;
