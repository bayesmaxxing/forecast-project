import React, { useState } from 'react';
import LeaderboardBase from './LeaderboardBase';
import { useUserData } from '../services/hooks/useUserData';
import { useScores } from '../services/hooks/useScores';

const METRICS = [
  { value: 'brier_score', label: 'Brier' },
  { value: 'brier_score_time_weighted', label: 'Brier (TW)' },
  { value: 'log2_score', label: 'Binary Log' },
  { value: 'log2_score_time_weighted', label: 'Binary Log (TW)' },
  { value: 'logn_score', label: 'Natural Log' },
  { value: 'logn_score_time_weighted', label: 'Natural Log (TW)' }
];

function ForecastLeaderboard({ forecastId, isResolved }) {
  const [selectedMetric, setSelectedMetric] = useState('brier_score_time_weighted');
  const { users, usersLoading } = useUserData();
  const { scores, loading: scoresLoading, error: scoresError } = useScores({
    type: 'basic',
    forecastId,
    shouldFetch: isResolved
  });

  // Don't render if forecast is not resolved
  if (!isResolved) {
    return null;
  }

  // Transform scores to leaderboard items
  const items = (scores || []).map(score => {
    const user = users?.find(u => u.id === score.user_id);
    return {
      id: score.user_id,
      name: user?.username || `User ${score.user_id}`,
      brier_score: score.brier_score,
      brier_score_time_weighted: score.brier_score_time_weighted,
      log2_score: score.log2_score,
      log2_score_time_weighted: score.log2_score_time_weighted,
      logn_score: score.logn_score,
      logn_score_time_weighted: score.logn_score_time_weighted
    };
  });

  return (
    <LeaderboardBase
      items={items}
      title="Top Forecasters"
      metrics={METRICS}
      selectedMetric={selectedMetric}
      onMetricChange={setSelectedMetric}
      loading={usersLoading || scoresLoading}
      error={scoresError}
      emptyMessage="No scores available for this forecast"
      maxItems={10}
    />
  );
}

export default ForecastLeaderboard;
