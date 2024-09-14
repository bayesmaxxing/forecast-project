import React, {useState, useEffect} from 'react';
import { useParams } from 'react-router-dom';
import { Link } from 'react-router-dom';
import './ForecastPage.css';
import Sidebar from '../components/Sidebar';


function ForecastPage() {
    const [searchQuery, setsearchQuery] = useState('');
    const [combinedForecasts, setCombinedForecasts] = useState([]);

    // set category based on URL
    let { category } = useParams()

    useEffect(() => {
      const CACHE_DURATION = 5 * 60 * 1000; // Cache duration in milliseconds, e.g., 5 minutes
      const now = new Date().getTime(); // Current time

      const forecastCatCacheKey = 'forecasts_${category}_w_latest';
    
      // Try to load data from cache
      const forecastCached = localStorage.getItem(forecastCatCacheKey);
    
      // Check if the cache is older than 5 minutes
      const forecastsDataValid = forecastCached && now - JSON.parse(forecastCached).timestamp < CACHE_DURATION;

      if (forecastsDataValid) {
        setCombinedForecasts(JSON.parse(forecastCached).data);
      } else {
        // Fetch the list of forecasts from the API
        Promise.all([
          fetch(`http://localhost:8080/forecasts?category=${category}&type=open`, {
            headers : {
              "Accept": "application/json"
            }
          }),
          fetch(`http://localhost:8080/forecast-points/latest`, {
            headers : {
              "Accept": "application/json"
            }
          })
        ])
        .then(async ([forecastData, pointsData]) => {
          const forecastDataJson = await forecastData.json();
          const pointsDataJson = await pointsData.json();
          return [forecastDataJson, pointsDataJson];
        })
        .then(([forecastDataJson, pointsDataJson]) => {
          const combined = forecastDataJson.map(forecast => {
            const matchingPoint = pointsDataJson.find(point => point.forecast_id === forecast.id);
            return { ...forecast, latestPoint: matchingPoint || null};
          });
          setCombinedForecasts(combined)

          localStorage.setItem(`forecasts_${category}_unresolved`, JSON.stringify({data: combined, timestamp: now}));
        })
        .catch(error => console.error('Error fetching data: ', error));
      }
    }, [category]);
    
  const handleSearchChange = (e) => {
    setsearchQuery(e.target.value.toLowerCase());
  };

  const filteredForecasts = combinedForecasts.filter(forecast => 
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
                {forecast.latestPoint ? (
                <p>{(forecast.latestPoint.point_forecast * 100).toFixed(1)}%</p>
                ) : ( 
                  <p>Not forecasted</p>
                )}
              </div>
            </div>
            <div>
                <p>Category: {forecast.category}</p>
                <p>Created: {formatDate(forecast.created)}</p>
            </div>
          </li>
        ))}
      </ul>
    </div>
  );
};
  
  export default ForecastPage;