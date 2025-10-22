import React from 'react';
import { Grid2 } from '@mui/material';
import ScoreDisplay from './ScoreDisplay';

function ScoresSection({ scores, loading }) {
  return (
    <Grid2 container spacing={2} sx={{ mt: 2 }}>
      <Grid2 xs={12} md={4}>
        <ScoreDisplay
          type="brier"
          value={scores?.brier_score_time_weighted || null}
          loading={loading}
        />
      </Grid2>
      <Grid2 xs={12} md={4}>
        <ScoreDisplay
          type="base2log"
          value={scores?.log2_score_time_weighted || null}
          loading={loading}
        />
      </Grid2>
      <Grid2 xs={12} md={4}>
        <ScoreDisplay
          type="baseNlog"
          value={scores?.logn_score_time_weighted || null}
          loading={loading}
        />
      </Grid2>
    </Grid2>
  );
}

export default ScoresSection;