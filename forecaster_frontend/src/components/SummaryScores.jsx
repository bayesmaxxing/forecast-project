import React, {useState, useEffect} from 'react';
import './SummaryScores.css';
import { Line } from 'react-chartjs-2';
import 'chart.js/auto';
import { useNavigate } from 'react-router-dom';
import 'chartjs-plugin-annotation';


function SummaryScores() {
    const [scores, setScores] = useState([]);
    const [resolutions, setResolutions] = useState([]);
    const [selectedMetric, setSelectedMetric] = useState('brier_score');
  
    const navigate = useNavigate();

    useEffect(() => {
       // Fetch the list of resolutions from the API if cache is empty
        Promise.all([
        fetch('https://forecasting-389105.ey.r.appspot.com/scores', {
            headers : {
              'Accept': 'application/json'
            }
          }),
          fetch('https://forecasting-389105.ey.r.appspot.com/forecasts?type=resolved', {
            headers : {
              'Accept': 'application/json'
            }
          })
        ])
          .then(async ([scoreData, forecastData]) => {
            const scoreDataJson = await scoreData.json();
            const forecastDataJson = await forecastData.json();
            return [scoreDataJson, forecastDataJson];
          })
          .then(([scoreDataJson, forecastDataJson]) => {
            setScores(scoreDataJson);
            setResolutions(forecastDataJson);
          })
          .catch(error => console.error('Error fetching data: ', error));
    }, []);

    const getScore = () => {
      if (!scores) return 0;

      switch(selectedMetric) {
        case 'brier_score':
          return scores.AggBrierScore ?? 0;
        case 'log2_score':
          return scores.AggLog2Score ?? 0;
        case 'logn_score':
          return scores.AggLogNScore ?? 0;
        default:
          return scores.AggBrierScore ?? 0;
      }
    };

    const avgScore = getScore();

    
    const metricScores = resolutions.map(resolution => resolution[selectedMetric]);
    const labels = resolutions.map(resolution => `Forecast ${resolution.id}`);

  
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
          const forecastId = resolutions[elementIndex].id; 
          navigate(`/forecast/${forecastId}`);
        }
      },
      scales: {
        x: {
          ticks: {
            display: false
          },
          grid: {
            display: false 
          }
        },
        y: {
          title: {
            display: true,
            text: 'Score'
          },
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
        <p>Average Score: {avgScore.toFixed(3)}</p>
        <Line data={data} options={options} />
      </div>
    );
  }
  
  export default SummaryScores;
