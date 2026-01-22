// Array of colors for different user lines
export const lineColors = [
  'rgb(75, 192, 192)',   // teal
  'rgb(255, 99, 132)',   // pink
  'rgb(54, 162, 235)',   // blue
  'rgb(255, 159, 64)',   // orange
  'rgb(153, 102, 255)',  // purple
  'rgb(255, 205, 86)',   // yellow
  'rgb(201, 203, 207)'   // grey
];

/**
 * Transforms user forecast points into chart-ready data (time-based view)
 * @param {Array} pointsData - Array of forecast points from multiple users
 * @param {Boolean} multiUser - Whether to process as multi-user or single user
 * @param {Number} minTimeWindowHours - Minimum hours between points (to avoid clustering)
 * @param {Array} users - Array of user objects for name lookup
 * @returns {Object} Formatted data for recharts
 */
export const prepareChartData = (pointsData, multiUser = true, minTimeWindowHours = 0, users = []) => {
  if (!pointsData || pointsData.length === 0) return null;

  const formatDate = (dateString) => {
    const date = new Date(dateString);
    return `${date.getMonth() + 1}/${date.getDate()}`;
  };

  // Check if the data is in the simplified format from the ordered endpoint
  const isSimplifiedFormat = pointsData.length > 0 && 'created_at' in pointsData[0];

  // Get the date field based on format
  const getDateField = (point) => isSimplifiedFormat ? point.created_at : point.created;

  // Filter points within a minimum time window
  const filterPointsByTimeWindow = (points) => {
    if (points.length <= 1 || minTimeWindowHours === 0) {
      return [...points].sort((a, b) => new Date(getDateField(a)) - new Date(getDateField(b)));
    }

    const minTimeWindowMs = minTimeWindowHours * 60 * 60 * 1000;
    const sortedPoints = [...points].sort((a, b) => new Date(getDateField(a)) - new Date(getDateField(b)));

    let filteredPoints = [];
    let lastTimestamp = null;

    sortedPoints.forEach(point => {
      const timestamp = new Date(getDateField(point)).getTime();
      if (lastTimestamp === null || (timestamp - lastTimestamp) >= minTimeWindowMs) {
        filteredPoints.push(point);
        lastTimestamp = timestamp;
      }
    });

    // Always include the most recent point
    const lastPoint = sortedPoints[sortedPoints.length - 1];
    const lastPointTime = new Date(getDateField(lastPoint)).getTime();
    if (filteredPoints.length === 0 ||
        lastPointTime !== new Date(getDateField(filteredPoints[filteredPoints.length - 1])).getTime()) {
      filteredPoints.push(lastPoint);
    }

    return filteredPoints;
  };

  if (multiUser) {
    // Group points by user
    const userPointsMap = {};
    pointsData.forEach(point => {
      const userId = point.user_id || 'anonymous';
      if (!userPointsMap[userId]) {
        userPointsMap[userId] = [];
      }
      userPointsMap[userId].push({ ...point });
    });

    // Create a dataset for each user
    const datasets = Object.entries(userPointsMap).map(([userId, userPoints], index) => {
      const sortedUserPoints = filterPointsByTimeWindow(userPoints);
      const user = users.find(u => u.id === parseInt(userId, 10));
      const userName = user ? user.username : (sortedUserPoints[0]?.username || `User ${userId}`);

      return {
        label: userName,
        data: sortedUserPoints.map(point => point.point_forecast),
        timestamps: sortedUserPoints.map(point => new Date(getDateField(point)).getTime()),
        dates: sortedUserPoints.map(point => formatDate(getDateField(point))),
        fill: false,
        borderColor: lineColors[index % lineColors.length],
        tension: 0.1
      };
    });

    // Collect all timestamps from all users
    let allTimestamps = [];
    datasets.forEach(dataset => {
      allTimestamps = [...allTimestamps, ...dataset.timestamps];
    });

    const uniqueTimestamps = [...new Set(allTimestamps)].sort((a, b) => a - b);
    const labels = uniqueTimestamps.map(ts => formatDate(new Date(ts)));

    return {
      labels,
      datasets,
      _timestamps: uniqueTimestamps
    };
  }

  // Single user mode
  const sortedPoints = filterPointsByTimeWindow(pointsData);
  const timestamps = sortedPoints.map(point => new Date(getDateField(point)).getTime());

  return {
    labels: sortedPoints.map(point => formatDate(getDateField(point))),
    datasets: [{
      label: 'Prediction',
      data: sortedPoints.map(point => point.point_forecast),
      timestamps,
      dates: sortedPoints.map(point => formatDate(getDateField(point))),
      fill: false,
      borderColor: lineColors[0],
      tension: 0.1
    }],
    _timestamps: timestamps
  };
};
