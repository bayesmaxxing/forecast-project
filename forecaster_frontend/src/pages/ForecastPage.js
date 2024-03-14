import React, {useState, useEffect} from 'react';
import { Link } from 'react-router-dom';
import './ForecastPage.css';
import Sidebar from '../components/Sidebar';


function ForecastPage() {
    const [forecasts, setForecasts] = useState([]);
    const [forecastpoints, setForecastPoints] = useState([]);
    useEffect(() => {
      const forecastsCache = localStorage.getItem('forecasts');
      const forecastPointsCache = localStorage.getItem('forecastPoints');
  
      // Try to load data from cache
      if (forecastsCache && forecastPointsCache) {
        setForecasts(JSON.parse(forecastsCache));
        setForecastPoints(JSON.parse(forecastPointsCache));
      } else {
        // Fetch the list of forecasts from the API if cache is empty
        Promise.all([
          fetch(`https://forecast-project-backend.vercel.app/forecaster/api/forecasts/?resolved=False`),
          fetch(`https://forecast-project-backend.vercel.app/forecaster/api/forecast_points/`)
        ])
        .then(async ([forecastData, pointsData]) => {
          const forecastDataJson = await forecastData.json();
          const pointsDataJson = await pointsData.json();
          return [forecastDataJson, pointsDataJson];
        })
        .then(([forecastDataJson, pointsDataJson]) => {
          // Update state with fetched data
          setForecasts(forecastDataJson);
          setForecastPoints(pointsDataJson);
          // Update cache with new data
          localStorage.setItem('forecasts', JSON.stringify(forecastDataJson));
          localStorage.setItem('forecastPoints', JSON.stringify(pointsDataJson));
        })
        .catch(error => console.error('Error fetching data: ', error));
      }
    }, []);

    const getRecentForecastPoint = (forecastId) => {
      const pointsForForecast = forecastpoints.filter(point => point.forecast === forecastId);
      const mostRecentPoint = pointsForForecast.reduce((maxPoint, currentPoint) => {
        return currentPoint.date_added > maxPoint.date_added ? currentPoint : maxPoint;
      }, pointsForForecast[0]);
  
      return mostRecentPoint;
    };
    
    const sortedForecasts = [...forecasts].sort((a, b)=>{
      return b.id - a.id;
    });

    const formatDate = (dateString) => dateString.split('T')[0];

    return (
      <div>
        <Sidebar></Sidebar>
        <h1>ALL QUESTIONS</h1>
        <ul className="forecast-list">
          {sortedForecasts.map(forecast => (
            <li key={forecast.id} className="forecast-item">
              <div className="question-container">
                <Link to={`/forecast/${forecast.id}`} className="question-link">
                  {forecast.question}
                </Link>
                <div className="recent-forecast-point">
                  <p>{(getRecentForecastPoint(forecast.id)?.point_forecast * 100).toFixed(1)}%</p>
                </div>
              </div>
              <div>
                <p>Category: {forecast.category}</p>
                <p>Created: {formatDate(forecast.creation_date)}</p>
              </div>
            </li>
          ))}
        </ul>
      </div>
    );
  };
  
  export default ForecastPage;
