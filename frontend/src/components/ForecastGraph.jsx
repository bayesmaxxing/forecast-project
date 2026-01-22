import React from 'react';
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
import { Paper, Typography, useTheme, Box, useMediaQuery } from '@mui/material';
import { format } from 'date-fns';
import UserSelector from './UserSelector';

function ForecastGraph({ data, options = {}, selectedUserId, onUserChange }) {
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'));

  const transformData = () => {
    if (!data || !data.datasets || data.datasets.length === 0) {
      return [];
    }

    // Time-based view using timestamps
    if (data._timestamps && data._timestamps.length > 0) {
      return data._timestamps.map((timestamp) => {
        const dataPoint = {
          timestamp,
          formattedTime: format(new Date(timestamp), 'MM/dd/yyyy HH:mm')
        };

        data.datasets.forEach(dataset => {
          const datasetIndex = dataset.timestamps?.indexOf(timestamp) ?? -1;
          if (datasetIndex !== -1) {
            dataPoint[dataset.label] = dataset.data[datasetIndex];
            if (dataset.dates?.[datasetIndex]) {
              dataPoint[`${dataset.label}_date`] = dataset.dates[datasetIndex];
            }
          }
        });

        return dataPoint;
      });
    }

    // Fallback for simple label-based data
    return data.labels?.map((label, index) => {
      const dataPoint = { name: label };
      data.datasets.forEach(dataset => {
        dataPoint[dataset.label] = dataset.data[index];
      });
      return dataPoint;
    }) || [];
  };

  const rechartsData = transformData();

  // Calculate time axis configuration
  const getXAxisConfig = () => {
    if (!data?._timestamps?.length) {
      return {
        dataKey: "name",
        type: "category",
        tick: { fill: theme.palette.text.primary, fontSize: 12 },
        height: isMobile ? 30 : 50
      };
    }

    const minTime = Math.min(...data._timestamps);
    const maxTime = Math.max(...data._timestamps);
    const timeRange = maxTime - minTime;
    const padding = Math.max(timeRange * 0.1, 24 * 60 * 60 * 1000); // 10% or 1 day min

    return {
      dataKey: "timestamp",
      type: "number",
      scale: "time",
      domain: [minTime - padding, maxTime + padding],
      tickFormatter: (tick) => format(new Date(tick), 'M/d'),
      tick: { fill: theme.palette.text.primary, fontSize: 12 },
      height: isMobile ? 30 : 50,
      dy: 10,
      ticks: calculateTimeTicks(minTime - padding, maxTime + padding)
    };
  };

  const calculateTimeTicks = (start, end) => {
    const range = end - start;
    const dayMs = 24 * 60 * 60 * 1000;

    let tickCount;
    if (range <= dayMs * 3) {
      tickCount = Math.ceil(range / dayMs) + 1;
    } else if (range <= dayMs * 14) {
      tickCount = 5;
    } else {
      tickCount = 7;
    }

    const step = range / (tickCount - 1);
    return Array.from({ length: tickCount }, (_, i) => start + i * step);
  };

  return (
    <Paper
      elevation={0}
      sx={{
        p: { xs: 1.5, sm: 2 },
        width: '100%',
        backgroundColor: theme.palette.background.paper,
        mb: 2
      }}
    >
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 1.5 }}>
        {options.title && (
          <Typography variant="subtitle1" component="h3" sx={{ fontWeight: 600 }}>
            {options.title.text}
          </Typography>
        )}

        {selectedUserId !== undefined && onUserChange && (
          <UserSelector onUserChange={onUserChange} selectedUserId={selectedUserId} />
        )}
      </Box>

      <Box
        sx={{
          width: '100%',
          position: 'relative',
          '&:before': {
            content: '""',
            display: 'block',
            paddingTop: { xs: '70%', sm: '55%', md: '45%' }
          }
        }}
      >
        <Box sx={{ position: 'absolute', top: 0, left: 0, right: 0, bottom: 0 }}>
          <ResponsiveContainer width="100%" height="100%">
            <LineChart
              data={rechartsData}
              margin={{ top: 10, right: 10, left: 0, bottom: isMobile ? 10 : 20 }}
            >
              <CartesianGrid strokeDasharray="3 3" stroke={theme.palette.divider} />
              <XAxis {...getXAxisConfig()} />
              <YAxis
                tick={{ fill: theme.palette.text.primary, fontSize: 12 }}
                domain={options.scales?.y?.min !== undefined
                  ? [options.scales.y.min, options.scales.y.max]
                  : ['auto', 'auto']}
              />
              <Tooltip
                contentStyle={{
                  backgroundColor: theme.palette.background.paper,
                  border: `1px solid ${theme.palette.divider}`,
                  borderRadius: theme.shape.borderRadius,
                  boxShadow: theme.shadows[1]
                }}
                formatter={(value, name, props) => {
                  const dateField = `${name}_date`;
                  const date = props.payload[dateField];
                  return date ? [`${value} (${date})`, name] : [value, name];
                }}
                labelFormatter={(label, payload) => {
                  if (payload?.[0]?.payload?.formattedTime) {
                    return payload[0].payload.formattedTime;
                  }
                  return typeof label === 'number'
                    ? format(new Date(label), 'MM/dd/yyyy HH:mm')
                    : label;
                }}
              />
              <Legend />
              {data.datasets.map((dataset, index) => (
                <Line
                  key={index}
                  type="monotone"
                  dataKey={dataset.label}
                  stroke={dataset.borderColor || theme.palette.primary.main}
                  strokeWidth={dataset.borderWidth || 2}
                  dot={{ r: 3 }}
                  activeDot={{ r: 5 }}
                  connectNulls={true}
                />
              ))}
            </LineChart>
          </ResponsiveContainer>
        </Box>
      </Box>
    </Paper>
  );
}

export default ForecastGraph;
