import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer
} from 'recharts';
import {
  Paper,
  Box,
  Typography,
  FormControl,
  Select,
  MenuItem,
  useTheme,
  useMediaQuery
} from '@mui/material';

function SummaryScores() {
  const [scores, setScores] = useState([]);
  const [resolutions, setResolutions] = useState([]);
  const [selectedMetric, setSelectedMetric] = useState('brier_score');
  const navigate = useNavigate();
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'));

  useEffect(() => {
    Promise.all([
      fetch('https://forecasting-389105.ey.r.appspot.com/scores', {
        headers: {
          'Accept': 'application/json'
        }
      }),
      fetch('https://forecasting-389105.ey.r.appspot.com/forecasts?type=resolved', {
        headers: {
          'Accept': 'application/json'
        }
      })
    ])
      .then(async ([scoreData, forecastData]) => {
        const scoreDataJson = await scoreData.json();
        const forecastDataJson = await forecastData.json();
        return [scoreDataJson, forecastDataJson];
      })
      .then(([scoreDataJson, forecastDataJson]) => {
        setScores(scoreDataJson);
        setResolutions(forecastDataJson);
      })
      .catch(error => console.error('Error fetching data: ', error));
  }, []);

  const getScore = () => {
    if (!scores) return 0;

    switch (selectedMetric) {
      case 'brier_score':
        return scores.AggBrierScore ?? 0;
      case 'log2_score':
        return scores.AggLog2Score ?? 0;
      case 'logn_score':
        return scores.AggLogNScore ?? 0;
      default:
        return scores.AggBrierScore ?? 0;
    }
  };

  const avgScore = getScore();

  const sortedResolutions = [...resolutions].sort((a, b) => {
    return new Date(b.resolved) - new Date(a.resolved);
  });

  const chartData = sortedResolutions.map(resolution => ({
    id: resolution.id,
    score: resolution[selectedMetric],
    label: `Forecast ${resolution.id}`
  }));

  const getMetricLabel = (metric) => {
    switch (metric) {
      case 'brier_score':
        return 'Brier Score';
      case 'log2_score':
        return 'Binary Log Score';
      case 'logn_score':
        return 'Natural Log Score';
      default:
        return 'Score';
    }
  };

  const handleChartClick = (data) => {
    if (data && data.activePayload && data.activePayload[0]) {
      const clickedPoint = data.activePayload[0].payload;
      navigate(`/forecast/${clickedPoint.id}`);
    }
  };

  return (
    <Paper 
      elevation={3} 
      sx={{ 
        p: { xs: 2, sm: 3 },
        m: { xs: 2, sm: 3 },
        bgcolor: 'background.paper'
      }}
    >
      <Box sx={{ mb: 3 }}>
        <FormControl fullWidth sx={{ mb: 2 }}>
          <Select
            value={selectedMetric}
            onChange={(e) => setSelectedMetric(e.target.value)}
            size="small"
          >
            <MenuItem value="brier_score">Brier Score</MenuItem>
            <MenuItem value="log2_score">Binary Log Score</MenuItem>
            <MenuItem value="logn_score">Natural Log Score</MenuItem>
          </Select>
        </FormControl>
        <Typography variant="h6" align="center" sx={{ mb: 2 }}>
          Average {getMetricLabel(selectedMetric)}: {avgScore.toFixed(3)}
        </Typography>
      </Box>

      <Box
        sx={{
          width: '100%',
          height: { xs: 300, sm: 400 },
          '& .recharts-tooltip-wrapper': {
            outline: 'none'
          }
        }}
      >
        <ResponsiveContainer>
          <LineChart
            data={chartData}
            margin={{
              top: 10,
              right: 30,
              left: 0,
              bottom: 5,
            }}
            onClick={handleChartClick}
          >
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis 
              dataKey="label"
              tick={false}
              height={20}
            />
            <YAxis
              tick={{ 
                fill: theme.palette.text.primary, 
                fontSize: { xs: 10, sm: 12 }
              }}
              label={{ 
                value: getMetricLabel(selectedMetric), 
                angle: -90, 
                position: 'insideLeft',
                style: { textAnchor: 'middle' }
              }}
            />
            <Tooltip
              contentStyle={{
                backgroundColor: theme.palette.background.paper,
                border: `1px solid ${theme.palette.divider}`,
                borderRadius: theme.shape.borderRadius,
              }}
              labelStyle={{
                color: theme.palette.text.primary
              }}
              formatter={(value) => [value.toFixed(3), getMetricLabel(selectedMetric)]}
            />
            <Line
              type="monotone"
              dataKey="score"
              stroke={theme.palette.primary.main}
              dot={{ r: 3 }}
              activeDot={{ r: 8 }}
              strokeWidth={2}
            />
          </LineChart>
        </ResponsiveContainer>
      </Box>
    </Paper>
  );
}

export default SummaryScores;
