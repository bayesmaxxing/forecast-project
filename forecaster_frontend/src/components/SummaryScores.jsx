import React, {useState, useEffect} from 'react';
import './SummaryScores.css';
import { Line } from 'react-chartjs-2';
import 'chart.js/auto';
import { useNavigate } from 'react-router-dom';
import 'chartjs-plugin-annotation';


function SummaryScores() {
    const [resolutions, setResolutions] = useState([]);
    const [selectedMetric, setSelectedMetric] = useState('brier_score');
    const [averageScore, setAverageScore] = useState(0);

    useEffect(() => {
      const CACHE_DURATION = 5 * 60 * 1000; // Cache duration in milliseconds, e.g., 5 minutes
      const now = new Date().getTime(); // Current time

      const resolutionsCache = localStorage.getItem('resolutions_scores');

      const resolutionsDataValid = resolutionsCache && now - JSON.parse(resolutionsCache).timestamp < CACHE_DURATION;
      // Try to load data from cache
      if (resolutionsDataValid) {
        setResolutions(JSON.parse(resolutionsCache).data);
      } else {
        // Fetch the list of resolutions from the API if cache is empty
        fetch('https://forecast-project-backend.vercel.app/forecaster/api/resolutions/', {
            headers : {
              'Authorization': `Token ${process.env.REACT_APP_API_TOKEN}`
            }
          })
          .then(response => response.json())
          .then(data => {
            // Update state with fetched data
            setResolutions(data);
            // Update cache with new data
            localStorage.setItem('resolutions_scores', JSON.stringify({data: data, timestamp: now}));
          })
          .catch(error => console.error('Error fetching data: ', error));
      }
    }, []);

    useEffect(() => {
      // Depending on the selected metric, calculate the appropriate average score
      let currentScores;
      switch(selectedMetric) {
          case 'brier_score':
              currentScores = resolutions.map(resolution => resolution.brier_score);
              break;
          case 'log2_score':
              currentScores = resolutions.map(resolution => resolution.log2_score);
              break;
          case 'logn_score':
              currentScores = resolutions.map(resolution => resolution.logn_score);
              break;
          default:
              currentScores = [];
      }
      const average = calculateAverage(currentScores);
      setAverageScore(average);
  }, [selectedMetric, resolutions]);
    
    const navigate = useNavigate();

    const calculateAverage = (scores) => {
      if (scores.length === 0) {
        return 0;
      }
      const sum = scores.reduce((a, b) => a + b, 0);
      return sum / scores.length;
    };
 
    const metricScores = resolutions.map(resolution => resolution[selectedMetric]);
    const labels = resolutions.map(resolution => `Forecast ${resolution.forecast}`);

  
    const data = {
      labels,
      datasets: [
        {
          label: 'Score',
          data: metricScores,
          fill: false,
          backgroundColor: 'rgb(75, 192, 192)',
          borderColor: 'rgba(75, 192, 192, 0.2)',
        },
      ],
    };
  
    const options = {
      onClick: (event, elements, chart) => {
        if (elements.length > 0) {
          const elementIndex = elements[0].index;
          const forecastId = resolutions[elementIndex].forecast; // Ensure this matches your data structure
          navigate(`/forecast/${forecastId}`); // Adjust the path as needed
        }
      },
      scales: {
        x: {
          ticks: {
            display: false // This will hide the X-axis labels
          },
          grid: {
            display: false // Optionally, this hides the X-axis grid lines if you want a cleaner look
          }
        },
        y: {
          beginAtZero: true
        }
      }
    };
  
    return (
      <div className="summary-box">
        <select value={selectedMetric} onChange={(e) => setSelectedMetric(e.target.value)}>
          <option value="brier_score">Brier Score</option>
          <option value="log2_score">Binary Log Score</option>
          <option value="logn_score">Natural Log Score</option>
        </select>
        <p>Average Score: {averageScore.toFixed(3)}</p>
        <Line data={data} options={options} />
      </div>
    );
  }
  
  export default SummaryScores;
