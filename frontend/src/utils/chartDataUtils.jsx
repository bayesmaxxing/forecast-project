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
   * @param {Array} pointsData - Array of forecast points from multiple users
   * @param {Boolean} multiUser - Whether to process as multi-user or single user
   * @returns {Object} Formatted data for recharts
   */
  export const prepareChartData = (pointsData, multiUser = true) => {
    if (!pointsData || pointsData.length === 0) return null;
    
    // If not multi-user mode, use the original single-user format
    if (!multiUser) {
      const sortedPoints = [...pointsData].sort((a, b) => 
        new Date(a.created) - new Date(b.created)
      );
      
      return {
        labels: sortedPoints.map(point => 
          new Date(point.created).toLocaleDateString('en-CA')
        ),
        datasets: [{
          label: 'Prediction',
          data: sortedPoints.map(point => point.point_forecast),
          fill: false,
          borderColor: lineColors[0],
          tension: 0.1
        }]
      };
    }
    
    // Multi-user processing
    // Group points by user
    const userPointsMap = {};
    
    pointsData.forEach(point => {
      const userId = point.user_id || 'anonymous';
      if (!userPointsMap[userId]) {
        userPointsMap[userId] = [];
      }
      userPointsMap[userId].push(point);
    });
    
    // Get all unique dates across all users' points
    const allDates = new Set();
    pointsData.forEach(point => {
      allDates.add(new Date(point.created).toLocaleDateString('en-CA'));
    });
    
    // Sort dates chronologically
    const sortedDates = Array.from(allDates).sort((a, b) => 
      new Date(a) - new Date(b)
    );
    
    // Create a dataset for each user
    const datasets = Object.entries(userPointsMap).map(([userId, userPoints], index) => {
      // Sort user points chronologically
      const sortedUserPoints = [...userPoints].sort(
        (a, b) => new Date(a.created) - new Date(b.created)
      );
      
      // Create a map of date to forecast point for quick lookup
      const dateToPointMap = {};
      sortedUserPoints.forEach(point => {
        dateToPointMap[new Date(point.created).toLocaleDateString('en-CA')] = point.point_forecast;
      });
      
      // Create data array that aligns with the sortedDates array
      const data = sortedDates.map(date => 
        dateToPointMap[date] !== undefined ? dateToPointMap[date] : null
      );
      
      
      const userName = `User ${userId}`;
      
      // Assign a color from the array, cycling if needed
      const colorIndex = index % lineColors.length;
      
      return {
        label: userName,
        data: data,
        fill: false,
        borderColor: lineColors[colorIndex],
        tension: 0.1,
        // Connect data points even with null values in between
        spanGaps: true
      };
    });
    
    return {
      labels: sortedDates,
      datasets: datasets
    };
  }; 