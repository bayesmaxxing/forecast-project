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
 * Transforms user forecast points into chart-ready data
 * Supporting both sequential and time-based views
 * @param {Array} pointsData - Array of forecast points from multiple users
 * @param {Boolean} multiUser - Whether to process as multi-user or single user
 * @param {Boolean} useSequential - Whether to use sequential numbering (#1, #2) or dates
 * @param {Number} minTimeWindowHours - Minimum hours between points (to avoid clustering)
 * @returns {Object} Formatted data for recharts
 */
export const prepareChartData = (pointsData, multiUser = true, useSequential = true, minTimeWindowHours = 0, users = []) => {
  if (!pointsData || pointsData.length === 0) return null;
  
  // Format date consistently for display
  const formatDate = (dateString) => {
    const date = new Date(dateString);
    return `${date.getMonth() + 1}/${date.getDate()}`;
  };
  
  console.log('Chart data preparation', { 
    pointsCount: pointsData.length, 
    multiUser,
    useSequential,
    userIds: [...new Set(pointsData.map(p => p.user_id))]
  });
  
  // Check if the data is in the simplified format from the ordered endpoint
  const isSimplifiedFormat = pointsData.length > 0 && 'created_at' in pointsData[0];
  
  // Filter points within a minimum time window (rather than by day)
  // This keeps more detail when forecasts are close together
  const filterPointsByTimeWindow = (points) => {
    if (points.length <= 1) return points;
    
    let filteredPoints = [];
    let lastTimestamp = null;
    const minTimeWindowMs = minTimeWindowHours * 60 * 60 * 1000;
    
    // Sort chronologically first
    const sortedPoints = [...points].sort((a, b) => {
      const dateA = new Date(isSimplifiedFormat ? a.created_at : a.created);
      const dateB = new Date(isSimplifiedFormat ? b.created_at : b.created);
      return dateA - dateB;
    });
    
    // Keep points that are separated by at least minTimeWindowMs
    sortedPoints.forEach(point => {
      const timestamp = new Date(isSimplifiedFormat ? point.created_at : point.created).getTime();
      
      if (lastTimestamp === null || (timestamp - lastTimestamp) >= minTimeWindowMs) {
        filteredPoints.push(point);
        lastTimestamp = timestamp;
      }
    });
    
    // Always include the most recent point
    const lastPoint = sortedPoints[sortedPoints.length - 1];
    const lastPointTime = new Date(isSimplifiedFormat ? lastPoint.created_at : lastPoint.created).getTime();
    
    if (filteredPoints.length === 0 || 
        lastPointTime !== new Date(isSimplifiedFormat ? 
          filteredPoints[filteredPoints.length - 1].created_at : 
          filteredPoints[filteredPoints.length - 1].created).getTime()) {
      filteredPoints.push(lastPoint);
    }
    
    return filteredPoints;
  };
  
  // Multi-user mode
  if (multiUser) {
    // Group points by user
    const userPointsMap = {};
    
    pointsData.forEach(point => {
      const userId = point.user_id || 'anonymous';
      if (!userPointsMap[userId]) {
        userPointsMap[userId] = [];
      }
      userPointsMap[userId].push({...point});
    });
    
    // Create a dataset for each user
    const datasets = Object.entries(userPointsMap).map(([userId, userPoints], index) => {
      // Apply time window filtering if not using sequential numbering
      let sortedUserPoints = useSequential ? 
        [...userPoints].sort((a, b) => {
          const dateA = new Date(isSimplifiedFormat ? a.created_at : a.created);
          const dateB = new Date(isSimplifiedFormat ? b.created_at : b.created);
          return dateA - dateB;
        }) : 
        filterPointsByTimeWindow(userPoints);
      
      // Get user name or use a generic name
      const user = users.find(u => u.id === parseInt(userId, 10));
      const userName = user ? user.username : (sortedUserPoints[0]?.username || `User ${userId}`);
      
      return {
        label: userName,
        data: sortedUserPoints.map(point => point.point_forecast),
        // Store full date objects for the x-axis scaling
        timestamps: sortedUserPoints.map(point => new Date(isSimplifiedFormat ? point.created_at : point.created).getTime()),
        // Store formatted dates for display
        dates: sortedUserPoints.map(point => formatDate(isSimplifiedFormat ? point.created_at : point.created)),
        fill: false,
        borderColor: lineColors[index % lineColors.length],
        tension: 0.1
      };
    });
    
    if (useSequential) {
      // Find the dataset with the most points to use for labels
      const maxPointsDataset = datasets.reduce((max, ds) => 
        ds.data.length > max.data.length ? ds : max, 
        { data: [] }
      );
      
      // Generate numeric labels (1, 2, 3, ...) matching the length of the longest dataset
      const labels = Array.from({ length: maxPointsDataset.data.length }, (_, i) => `#${i+1}`);
      
      return {
        labels,
        datasets,
        _isSequenced: true
      };
    } else {
      // For date-based x-axis, we need to collect all timestamps from all users
      let allTimestamps = [];
      datasets.forEach(dataset => {
        allTimestamps = [...allTimestamps, ...dataset.timestamps];
      });
      
      // Get unique timestamps and sort
      const uniqueTimestamps = [...new Set(allTimestamps)].sort((a, b) => a - b);
      
      // Create labels from timestamps
      const labels = uniqueTimestamps.map(ts => formatDate(new Date(ts)));
      
      // Create actual data structure with proper timestamp range
      // No need for additional padding - this will be handled by the chart component
      return {
        labels,
        datasets,
        _isSequenced: false,
        _timestamps: uniqueTimestamps
      };
    }
  }
  
  // Single user mode
  // Apply time window filtering if not using sequential mode
  let sortedPoints = useSequential ? 
    [...pointsData].sort((a, b) => {
      const dateA = new Date(isSimplifiedFormat ? a.created_at : a.created);
      const dateB = new Date(isSimplifiedFormat ? b.created_at : b.created);
      return dateA - dateB;
    }) : 
    filterPointsByTimeWindow(pointsData);
  
  const timestamps = sortedPoints.map(point => 
    new Date(isSimplifiedFormat ? point.created_at : point.created).getTime()
  );
  
  if (useSequential) {
    return {
      // For x-axis labels, use simple sequential numbers
      labels: sortedPoints.map((_, i) => `#${i+1}`),
      datasets: [{
        label: 'Prediction',
        data: sortedPoints.map(point => point.point_forecast),
        // Store dates for tooltips - handle both formats
        dates: sortedPoints.map(point => formatDate(isSimplifiedFormat ? point.created_at : point.created)),
        fill: false,
        borderColor: lineColors[0],
        tension: 0.1
      }],
      _isSequenced: true
    };
  } else {
    // Build data structure for time-based view
    // No need for padding - that's handled by the chart component
    return {
      labels: sortedPoints.map(point => formatDate(isSimplifiedFormat ? point.created_at : point.created)),
      datasets: [{
        label: 'Prediction',
        data: sortedPoints.map(point => point.point_forecast),
        timestamps: timestamps,
        dates: sortedPoints.map(point => formatDate(isSimplifiedFormat ? point.created_at : point.created)),
        fill: false,
        borderColor: lineColors[0],
        tension: 0.1
      }],
      _isSequenced: false,
      _timestamps: timestamps
    };
  }
};