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
 * Optimized to work with the ordered forecast points endpoint
 * @param {Array} pointsData - Array of forecast points from multiple users
 * @param {Boolean} multiUser - Whether to process as multi-user or single user
 * @returns {Object} Formatted data for recharts
 */
export const prepareChartData = (pointsData, multiUser = true) => {
  if (!pointsData || pointsData.length === 0) return null;
  
  // Format date consistently for display
  const formatDate = (dateString) => {
    const date = new Date(dateString);
    return `${date.getMonth() + 1}/${date.getDate()}`;
  };
  
  console.log('Chart data preparation', { 
    pointsCount: pointsData.length, 
    multiUser, 
    userIds: [...new Set(pointsData.map(p => p.user_id))]
  });
  
  // Check if the data is in the simplified format from the ordered endpoint
  const isSimplifiedFormat = pointsData.length > 0 && 'created_at' in pointsData[0];
  
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
      // Sort user points chronologically 
      const sortedUserPoints = [...userPoints].sort((a, b) => {
        const dateA = new Date(isSimplifiedFormat ? a.created_at : a.created);
        const dateB = new Date(isSimplifiedFormat ? b.created_at : b.created);
        return dateA - dateB;
      });
      
      // Get user name or use a generic name
      const userName = sortedUserPoints[0].username || `User ${userId}`;
      
      return {
        label: userName,
        data: sortedUserPoints.map(point => point.point_forecast),
        // Store dates for reference - handle both formats
        dates: sortedUserPoints.map(point => formatDate(isSimplifiedFormat ? point.created_at : point.created)),
        fill: false,
        borderColor: lineColors[index % lineColors.length],
        tension: 0.1
      };
    });
    
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
  }
  
  // Single user mode
  // Sort points chronologically first
  const sortedPoints = [...pointsData].sort((a, b) => {
    const dateA = new Date(isSimplifiedFormat ? a.created_at : a.created);
    const dateB = new Date(isSimplifiedFormat ? b.created_at : b.created);
    return dateA - dateB;
  });
  
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
};