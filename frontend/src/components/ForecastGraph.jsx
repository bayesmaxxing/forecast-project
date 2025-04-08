import React from 'react';
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
  ReferenceArea
} from 'recharts';
import { Paper, Typography, useTheme, Box, useMediaQuery, Switch, FormControlLabel } from '@mui/material';
import { format } from 'date-fns';

function ForecastGraph({ data, options = {} }) {
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'));
  const [useSequentialView, setUseSequentialView] = React.useState(options.useSequential !== false);

  // Toggle between date and sequential views
  const handleViewChange = (event) => {
    setUseSequentialView(event.target.checked);
  };

  const transformData = () => {
    if (!data || !data.labels || !data.datasets || data.datasets.length === 0) {
      console.log('Invalid chart data format:', data);
      return [];
    }
    
    // For sequential view, use the existing sequential data
    if (useSequentialView && data._isSequenced) {
      return data.labels.map((label, index) => {
        const dataPoint = { name: label };
        data.datasets.forEach(dataset => {
          if (index < dataset.data.length) {
            dataPoint[dataset.label] = dataset.data[index];
            if (dataset.dates && dataset.dates[index]) {
              dataPoint[`${dataset.label}_date`] = dataset.dates[index];
            }
          }
        });
        return dataPoint;
      });
    }
    
    // For date-based view
    if (data._timestamps && data._timestamps.length > 0) {
      return data._timestamps.map((timestamp, index) => {
        const dataPoint = { 
          timestamp: timestamp,
          // Format the time for display in tooltips
          formattedTime: format(new Date(timestamp), 'MM/dd/yyyy HH:mm')
        };
        
        // Add dataset values matching this timestamp
        data.datasets.forEach(dataset => {
          const datasetIndex = dataset.timestamps ? 
            dataset.timestamps.indexOf(timestamp) : -1;
          
          if (datasetIndex !== -1) {
            dataPoint[dataset.label] = dataset.data[datasetIndex];
            if (dataset.dates && dataset.dates[datasetIndex]) {
              dataPoint[`${dataset.label}_date`] = dataset.dates[datasetIndex];
            }
          }
        });
        
        return dataPoint;
      });
    }
    
    // Fallback implementation
    return data.labels.map((label, index) => {
      const dataPoint = { name: label };
      data.datasets.forEach(dataset => {
        dataPoint[dataset.label] = dataset.data[index];
      });
      return dataPoint;
    });
  };

  const rechartsData = transformData();

  // Customize the x-axis tick display
  const formatXAxis = (tickItem) => {
    if (useSequentialView) return tickItem;
    
    const date = new Date(tickItem);
    return format(date, 'M/d');
  };

  // Create a properly scaled time axis for date view
  const getXAxisConfig = () => {
    if (useSequentialView) {
      return {
        dataKey: "name",
        type: "category",
        tick: { 
          fill: theme.palette.text.primary, 
          fontSize: 12 
        },
        angle: 0,
        textAnchor: "middle",
        height: isMobile ? 30 : 50
      };
    } else {
      // Find min and max timestamps
      let minTime = Infinity;
      let maxTime = -Infinity;
      
      if (data && data._timestamps) {
        data._timestamps.forEach(timestamp => {
          minTime = Math.min(minTime, timestamp);
          maxTime = Math.max(maxTime, timestamp);
        });
      }
      
      // Calculate proper domain with padding
      const timeRange = maxTime - minTime;
      const padding = Math.max(timeRange * 0.1, 24 * 60 * 60 * 1000); // 10% or 1 day min
      
      return {
        dataKey: "timestamp",
        type: "number",
        scale: "time",
        domain: [minTime - padding, maxTime + padding],
        tickFormatter: formatXAxis,
        tick: { 
          fill: theme.palette.text.primary, 
          fontSize: 12 
        },
        angle: 0,
        textAnchor: "middle",
        height: isMobile ? 30 : 50,
        dy: 10,
        // Customize the number of ticks to avoid crowding
        ticks: calculateTimeTicks(minTime - padding, maxTime + padding)
      };
    }
  };
  
  // Calculate appropriate time ticks based on the range
  const calculateTimeTicks = (start, end) => {
    const range = end - start;
    const dayMs = 24 * 60 * 60 * 1000;
    
    // Determine number of ticks based on range
    let tickCount;
    if (range <= dayMs * 3) {
      // If 3 days or less, show ticks for each day
      tickCount = Math.ceil(range / dayMs) + 1;
    } else if (range <= dayMs * 14) {
      // If 2 weeks or less, show ~5 ticks
      tickCount = 5;
    } else {
      // For longer ranges, show ~7 ticks
      tickCount = 7;
    }
    
    // Calculate evenly spaced ticks
    const step = range / (tickCount - 1);
    return Array.from({length: tickCount}, (_, i) => start + i * step);
  };

  return (
    <Paper 
      elevation={2} 
      sx={{ 
        p: { xs: 1, sm: 2, md: 3 },
        width: '100%',
        backgroundColor: theme.palette.background.paper 
      }}
    >
      {options.title && (
        <Typography 
          variant="h6" 
          component="h3" 
          align="center" 
          sx={{ mb: { xs: 1, sm: 2 } }}
        >
          {options.title.text}
        </Typography>
      )}
      
      <Box
        sx={{
          display: 'flex',
          justifyContent: 'flex-end',
          mb: 1
        }}
      >
        <FormControlLabel
          control={
            <Switch
              checked={useSequentialView}
              onChange={handleViewChange}
              color="primary"
              size="small"
            />
          }
          label={useSequentialView ? "Sequential View" : "Time View"}
          labelPlacement="start"
          sx={{ ml: 0, fontSize: '0.8rem' }}
        />
      </Box>
      
      <Box
        sx={{
          width: '100%',
          position: 'relative',
          '&:before': {
            content: '""',
            display: 'block',
            paddingTop: {
              xs: '85%',
              sm: '75%',
              md: '56.25%'
            }
          }
        }}
      >
        <Box
          sx={{
            position: 'absolute',
            top: 0,
            left: 0,
            right: 0,
            bottom: 0
          }}
        >
          <ResponsiveContainer width="100%" height="100%">
            <LineChart
              data={rechartsData}
              margin={{
                top: 10,
                right: 10,
                left: 0,
                bottom: isMobile ? 10 : 20, 
              }}
            >
              <CartesianGrid 
                strokeDasharray="3 3" 
                stroke={theme.palette.divider}
              />
              <XAxis {...getXAxisConfig()} />
              <YAxis
                tick={{ 
                  fill: theme.palette.text.primary, 
                  fontSize: { xs: 10, sm: 12 }
                }}
                domain={options.scales?.y?.min !== undefined ? 
                  [options.scales.y.min, options.scales.y.max] : 
                  ['auto', 'auto']
                }
              />
              <Tooltip 
                contentStyle={{
                  backgroundColor: theme.palette.background.paper,
                  border: `1px solid ${theme.palette.divider}`,
                  borderRadius: theme.shape.borderRadius,
                  boxShadow: theme.shadows[1],
                  fontSize: { xs: 10, sm: 12 }
                }}
                formatter={(value, name, props) => {
                  // Find the date field for this dataset
                  const dateField = `${name}_date`;
                  const date = props.payload[dateField];
                  
                  // Return the formatted value with date if available
                  if (date) {
                    return [`${value} (${date})`, name];
                  }
                  return [value, name];
                }}
                labelFormatter={(label, payload) => {
                  if (!useSequentialView && payload && payload.length > 0) {
                    return payload[0].payload.formattedTime || format(new Date(label), 'MM/dd/yyyy HH:mm');
                  }
                  return useSequentialView ? `Prediction: ${label}` : label;
                }}
              />
              <Legend 
                wrapperStyle={{
                  fontSize: { xs: 10, sm: 12 }
                }}
              />
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