import React, {useState, useEffect} from 'react';
import { useParams } from 'react-router-dom';
import { Link } from 'react-router-dom';
import './ForecastPage.css';
import Sidebar from '../components/Sidebar';


function ForecastPage() {
    const [forecasts, setForecasts] = useState([]);
    const [forecastPoints, setForecastPoints]=useState([]);
    const [searchQuery, setsearchQuery] = useState('');

    // set category based on URL
    let { category } = useParams()

    useEffect(() => {
      const CACHE_DURATION = 5 * 60 * 1000; // Cache duration in milliseconds, e.g., 5 minutes
      const now = new Date().getTime(); // Current time

      const forecastsCacheKey = `forecasts_${category}_unresolved`;
      const forecastPointsCacheKey = 'forecast_points';
    
      // Try to load data from cache
      const forecastsCached = localStorage.getItem(forecastsCacheKey);
      const forecastPointsCached = localStorage.getItem(forecastPointsCacheKey);
    
      // Check if the cache is older than 5 minutes
      const forecastsDataValid = forecastsCached && now - JSON.parse(forecastsCached).timestamp < CACHE_DURATION;
      const forecastPointsDataValid = forecastPointsCached && now - JSON.parse(forecastPointsCached).timestamp < CACHE_DURATION;

      if (forecastsDataValid && forecastPointsDataValid) {
        setForecasts(JSON.parse(forecastsCached).data);
        setForecastPoints(JSON.parse(forecastPointsCached).data);
      } else {
        // Fetch the list of forecasts from the API
        Promise.all([
          fetch(`https://forecasting-389105.ey.r.appspot.com/forecaster/api/forecasts/?category=${category}&resolved=False`, {
            headers : {
              'Authorization': `Token ${process.env.REACT_APP_API_TOKEN}`
            }
          }),
          fetch(`https://forecasting-389105.ey.r.appspot.com/forecaster/api/forecast_points/`, {
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
          setForecasts(forecastDataJson);
          setForecastPoints(pointsDataJson);
          // Cache the new data
          localStorage.setItem(`forecasts_${category}_unresolved`, JSON.stringify({data: forecastDataJson, timestamp: now}));
          localStorage.setItem('forecast_points', JSON.stringify({data: pointsDataJson, timestamp: now}));
        })
        .catch(error => console.error('Error fetching data: ', error));
      }
    }, [category]);
    
  const getRecentForecastPoint = (forecastId) => {
    const pointsForForecast = forecastPoints.filter(point => point.forecast === forecastId);
    const mostRecentPoint = pointsForForecast.reduce((maxPoint, currentPoint) => {
      return currentPoint.date_added > maxPoint.date_added ? currentPoint : maxPoint;
    }, pointsForForecast[0]);

    return mostRecentPoint;
  };

  const handleSearchChange = (e) => {
    setsearchQuery(e.target.value.toLowerCase());
  };

  const filteredForecasts = forecasts.filter(forecast => 
    forecast.question.toLowerCase().includes(searchQuery) ||
    forecast.short_question.toLowerCase().includes(searchQuery) ||
    forecast.category.toLowerCase().includes(searchQuery) ||
    forecast.resolution_criteria.toLowerCase().includes(searchQuery)
    );

  const sortedForecasts = [...filteredForecasts].sort((a, b)=>{
    return b.id - a.id;
  });

  const formatDate = (dateString) => dateString.split('T')[0];

  return (
    <div>
      <Sidebar onSearchChange={handleSearchChange}/>
      <h1 style={{ textTransform: 'uppercase' }}>{category}</h1>
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