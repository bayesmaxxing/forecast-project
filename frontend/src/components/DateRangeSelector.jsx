import React from 'react';
import {
  FormControl,
  Select,
  MenuItem,
  Box,
  Typography
} from '@mui/material';
import { CalendarMonth } from '@mui/icons-material';
import { DATE_RANGE_OPTIONS } from '../services/api/scoreService';

const DATE_RANGE_LABELS = {
  [DATE_RANGE_OPTIONS.ALL_TIME]: 'All Time',
  [DATE_RANGE_OPTIONS.LAST_12_MONTHS]: 'Last 12 Months',
  [DATE_RANGE_OPTIONS.LAST_6_MONTHS]: 'Last 6 Months',
  [DATE_RANGE_OPTIONS.LAST_3_MONTHS]: 'Last 3 Months',
};

function DateRangeSelector({ value, onChange, size = 'small', showIcon = true }) {
  return (
    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
      {showIcon && <CalendarMonth sx={{ color: 'text.secondary', fontSize: 20 }} />}
      <FormControl size={size} sx={{ minWidth: 140 }}>
        <Select
          value={value}
          onChange={(e) => onChange(e.target.value)}
        >
          {Object.entries(DATE_RANGE_LABELS).map(([key, label]) => (
            <MenuItem key={key} value={key}>
              {label}
            </MenuItem>
          ))}
        </Select>
      </FormControl>
    </Box>
  );
}

export default DateRangeSelector;
