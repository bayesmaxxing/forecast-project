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

function ForecastGraph({ data, options = {} }) {
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'));

  const transformData = () => {
    if (!data || !data.labels || !data.datasets || data.datasets.length === 0) {
      console.log('Invalid chart data format:', data);
      return [];
    }
    
    // Log for debugging
    console.log('Transforming chart data', {
      labels: data.labels.length,
      datasets: data.datasets.length,
      sampleDataset: data.datasets[0].data.length
    });
    
    // Handle sequence-based data (using prediction number on x-axis)
    if (data._isSequenced) {
      return data.labels.map((label, index) => {
        // Start with name for x-axis
        const dataPoint = { name: label };
        
        // For each dataset, add its value if it exists at this index
        data.datasets.forEach(dataset => {
          if (index < dataset.data.length) {
            dataPoint[dataset.label] = dataset.data[index];
            // Store the date string for tooltip display
            if (dataset.dates && dataset.dates[index]) {
              dataPoint[`${dataset.label}_date`] = dataset.dates[index];
            }
          }
        });
        
        return dataPoint;
      });
    }
    
    // Original date-based implementation
    return data.labels.map((label, index) => {
      const dataPoint = { name: label };
      data.datasets.forEach(dataset => {
        dataPoint[dataset.label] = dataset.data[index];
      });
      return dataPoint;
    });
  };

  const rechartsData = transformData();

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
                bottom: isMobile ? 10 : 20, // Reduced bottom margin on mobile since we're hiding labels
              }}
            >
              <CartesianGrid 
                strokeDasharray="3 3" 
                stroke={theme.palette.divider}
              />
              <XAxis 
                dataKey="name"
                tick={isMobile ? false : { 
                  fill: theme.palette.text.primary, 
                  fontSize: 12
                }}
                angle={0}
                textAnchor="middle"
                height={isMobile ? 30 : 50}
                dy={10}
              />
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
                labelFormatter={(label) => {
                  return `Prediction: ${label}`;
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