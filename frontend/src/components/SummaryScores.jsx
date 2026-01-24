import React, { useState } from 'react';
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
import { useScores } from '../services/hooks/useScores';
import { useForecastList } from '../services/hooks/useForecastList';
import { getStartDateForRange } from '../services/api/scoreService';

function SummaryScores({user_id='all', dateRange=null}) {

  const [selectedMetric, setSelectedMetric] = useState('brier_score');
  const navigate = useNavigate();
  const theme = useTheme();

  const { scores: aggregateScores, loading: aggregateScoresLoading, error: aggregateScoresError } = useScores({
    type: 'aggregate',
    userId: user_id,
    dateRange
  });
  const { forecasts = [], loading: forecastsLoading, error: forecastsError } = useForecastList({list_type: 'resolved'});
  const { scores = [], loading: scoresLoading, error: scoresError } = useScores({
    type: user_id === 'all' ? 'average' : 'basic',
    userId: user_id
  });
  
  const getScore = () => {
    if (!aggregateScores) return 0;

    switch (selectedMetric) {
      case 'brier_score':
        return aggregateScores.brier_score_time_weighted ?? 0;
      case 'log2_score':
        return aggregateScores.log2_score_time_weighted ?? 0;
      case 'logn_score':
        return aggregateScores.logn_score_time_weighted ?? 0;
      default:
        return aggregateScores.brier_score_time_weighted ?? 0;
    }
  };

  const avgScore = getScore();

  const sortedForecasts = [...forecasts].sort((a, b) => {
    return new Date(b.resolved) - new Date(a.resolved);
  });

  // Filter forecasts by date range
  const startDate = getStartDateForRange(dateRange);
  const filteredForecasts = startDate
    ? sortedForecasts.filter(f => new Date(f.resolved) >= startDate)
    : sortedForecasts;

  // Limit to most recent 100 forecasts to keep chart readable
  const limitedForecasts = filteredForecasts.slice(0, 100);

  const combined = limitedForecasts.map(forecast => ({
    ...forecast,
    score: (scores || []).find(score => score.forecast_id === forecast.id)?.[`${selectedMetric}_time_weighted`] ?? 0
  }));

  const chartData = combined.map(forecast => ({
    id: forecast.id,
    score: forecast.score,
    label: `${forecast.question}`
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
      elevation={0}
      sx={{
        p: 2.5,
        bgcolor: 'background.paper',
        height: '100%',
        width: '100%',
        display: 'flex',
        flexDirection: 'column'
      }}
    >
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2, flexWrap: 'wrap', gap: 2 }}>
        <FormControl sx={{ minWidth: 150 }}>
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
        <Typography variant="subtitle1" sx={{ fontWeight: 600 }}>
          Average: {avgScore.toFixed(3)}
        </Typography>
      </Box>

      <Box
        sx={{
          width: '100%',
          flex: 1,
          minHeight: 0,
          cursor: 'pointer',
          '& .recharts-tooltip-wrapper': {
            outline: 'none'
          },
          '& .recharts-line': {
            transition: 'stroke-width 0.2s ease-in-out',
          },
          '&:hover .recharts-line': {
            strokeWidth: 3,
          }
        }}
      >
        <ResponsiveContainer>
          <LineChart
            data={chartData}
            margin={{
              top: 5,
              right: 20,
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
              dot={{ r: 4, cursor: 'pointer' }}
              activeDot={{ r: 7, cursor: 'pointer' }}
              strokeWidth={2}
            />
          </LineChart>
        </ResponsiveContainer>
      </Box>
    </Paper>
  );
}

export default SummaryScores;
