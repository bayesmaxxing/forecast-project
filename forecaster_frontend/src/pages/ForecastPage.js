import React, {useState, useEffect} from 'react';
import { Link } from 'react-router-dom';
import './ForecastPage.css';
import Sidebar from '../components/Sidebar';


function ForecastPage() {
    const [searchQuery, setsearchQuery] = useState('');
    const [combinedForecasts, setCombinedForecasts] = useState([]);

    useEffect(() => {
      const CACHE_DURATION = 5 * 60 * 10; // Cache duration in milliseconds, e.g., 5 minutes
      const now = new Date().getTime(); // Current time
      
      const cachedForecasts = localStorage.getItem('open_w_latest');
      
      const forecastsDataValid = cachedForecasts && now - JSON.parse(cachedForecasts).timestamp < CACHE_DURATION;
      
      if (forecastsDataValid) {
        setCombinedForecasts(JSON.parse(cachedForecasts).data);
      } else {
        Promise.all([
          fetch(`https://forecasting-389105.ey.r.appspot.com/forecasts?type=open`, {
            headers : {
              "Accept": "application/json"
            }
          }),
          fetch(`https://forecasting-389105.ey.r.appspot.com/forecast-points/latest`, {
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
          setCombinedForecasts(combined);
          
          localStorage.setItem('open_w_latest', JSON.stringify({data: combinedForecasts, timestamp: now}));
        })
        // If there is some error fetching the data
        .catch(error => console.error('Error fetching data: ', error));
      }
    },[]);
    
    // handler for the query in search field
    const handleSearchChange = (e) => {
      setsearchQuery(e.target.value.toLowerCase());
    };

    // Filtering the list of forecasts based on query in search field
    // query can match to question, short_question, category, or resolution criteria.
    const filteredForecasts = combinedForecasts.filter(forecast => 
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
