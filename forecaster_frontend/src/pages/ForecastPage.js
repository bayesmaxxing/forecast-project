import React, {useState, useEffect} from 'react';
import { Link } from 'react-router-dom';
import './ForecastPage.css';
import Sidebar from '../components/Sidebar';


function ForecastPage() {
    
    const [forecasts, setForecasts] = useState([]);
    const [forecastpoints, setForecastPoints] = useState([]);
    const [searchQuery, setsearchQuery] = useState('');

    useEffect(() => {
      const CACHE_DURATION = 5 * 60 * 1000; // Cache duration in milliseconds, e.g., 5 minutes
      const now = new Date().getTime(); // Current time
      
      // Get data from local cache
      const forecastsCache = localStorage.getItem('forecasts');
      const forecastPointsCache = localStorage.getItem('forecastPoints');
      
      // Check if the cache is older than 5 minutes
      const forecastsDataValid = forecastsCache && now - JSON.parse(forecastsCache).timestamp < CACHE_DURATION;
      const forecastPointsDataValid = forecastPointsCache && now - JSON.parse(forecastPointsCache).timestamp < CACHE_DURATION;
      
      // Try to load data from cache if it's valid
      if (forecastsDataValid && forecastPointsDataValid) {
        setForecasts(JSON.parse(forecastsCache).data);
        setForecastPoints(JSON.parse(forecastPointsCache).data);
      } else {
        // Fetch the list of forecasts from the API if cache is empty or expired
        Promise.all([
          fetch(`https://forecast-project-backend.vercel.app/forecaster/api/forecasts/?resolved=False`, {
            headers : {
              'Authorization': `Token ${process.env.REACT_APP_API_TOKEN}`
            }
          }),
          fetch(`https://forecast-project-backend.vercel.app/forecaster/api/forecast_points/`, {
            headers : {
              'Authorization': `Token ${process.env.REACT_APP_API_TOKEN}`
            }
          })
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
          
          // Update cache with new data and current timestamp
          localStorage.setItem('forecasts', JSON.stringify({data: forecastDataJson, timestamp: now}));
          localStorage.setItem('forecastPoints', JSON.stringify({data: pointsDataJson, timestamp: now}));
        })
        // If there is some error fetching the data
        .catch(error => console.error('Error fetching data: ', error));
      }
    }, []);

    // Get the most recent forecast for each question in forecasts
    const getRecentForecastPoint = (forecastId) => {
      const pointsForForecast = forecastpoints.filter(point => point.forecast === forecastId);
      const mostRecentPoint = pointsForForecast.reduce((maxPoint, currentPoint) => {
        return currentPoint.date_added > maxPoint.date_added ? currentPoint : maxPoint;
      }, pointsForForecast[0]);
  
      return mostRecentPoint;
    };
    
    // handler for the query in search field
    const handleSearchChange = (e) => {
      setsearchQuery(e.target.value.toLowerCase());
    };

    // Filtering the list of forecasts based on query in search field
    // query can match to question, short_question, category, or resolution criteria.
    const filteredForecasts = forecasts.filter(forecast => 
      forecast.question.toLowerCase().includes(searchQuery) ||
      forecast.short_question.toLowerCase().includes(searchQuery) ||
      forecast.category.toLowerCase().includes(searchQuery) ||
      forecast.resolution_criteria.toLowerCase().includes(searchQuery)
      );
    
    //Sorting forecasts based on id (higher id = more recently added)
    const sortedForecasts = [...filteredForecasts].sort((a, b)=>{
      return b.id - a.id;
    });

    const formatDate = (dateString) => dateString.split('T')[0];

    return (
      <div>
        <Sidebar onSearchChange={handleSearchChange}/>
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
