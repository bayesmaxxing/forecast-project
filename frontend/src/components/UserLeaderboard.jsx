import React, { useState } from 'react';
import LeaderboardBase from './LeaderboardBase';
import { useUserData } from '../services/hooks/useUserData';
import { useScores } from '../services/hooks/useScores';

const METRICS = [
  { value: 'brier_score', label: 'Brier' },
  { value: 'log2_score', label: 'Binary Log' },
  { value: 'logn_score', label: 'Natural Log' }
];

function UserLeaderboard({ dateRange = null }) {
  const [selectedMetric, setSelectedMetric] = useState('brier_score');
  const { users, usersLoading } = useUserData();
  const { scores, loading: scoresLoading, error: scoresError } = useScores({
    type: 'aggregate-by-users',
    dateRange
  });

  // Transform scores to leaderboard items
  const items = (scores || [])
    .filter(scoreData => scoreData.total_forecasts > 0)
    .map(scoreData => {
      const user = users?.find(u => u.id === scoreData.user_id);
      return {
        id: scoreData.user_id,
        name: user?.username || `User ${scoreData.user_id}`,
        subtitle: `${scoreData.total_forecasts} forecast${scoreData.total_forecasts !== 1 ? 's' : ''}`,
        brier_score: scoreData.brier_score_time_weighted,
        log2_score: scoreData.log2_score_time_weighted,
        logn_score: scoreData.logn_score_time_weighted
      };
    });

  return (
    <LeaderboardBase
      items={items}
      title="Leaderboard"
      metrics={METRICS}
      selectedMetric={selectedMetric}
      onMetricChange={setSelectedMetric}
      loading={usersLoading || scoresLoading}
      error={scoresError}
      emptyMessage="No data available"
    />
  );
}

export default UserLeaderboard;
