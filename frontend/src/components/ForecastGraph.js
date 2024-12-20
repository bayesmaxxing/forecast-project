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
    if (!data || !data.labels) return [];
    
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
                labelFormatter={(label) => `Date: ${label}`}
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
