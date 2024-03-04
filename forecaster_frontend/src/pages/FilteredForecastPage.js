import React, {useState, useEffect} from 'react';
import { useParams } from 'react-router-dom';
import { Link } from 'react-router-dom';
import './ForecastPage.css';
import Sidebar from '/Users/samuelsvensson/Documents/forecasting_project/forecaster_frontend/src/components/Sidebar';


function ForecastPage() {
    const [forecasts, setForecasts] = useState([]);
    const [forecastPoints, setForecastPoints]=useState([]);
    let { category } = useParams()

    useEffect(() => {
      // Fetch the list of forecasts from the API
      Promise.all([
        fetch(`http://127.0.0.1:8000/forecaster/api/forecasts/?category=${category}&resolved=False`),
        fetch(`http://127.0.0.1:8000/forecaster/api/forecast_points/`)
      ])
      .then( async ([forecastData, pointsData]) => {
        const forecastDataJson = await forecastData.json();
        const pointsDataJson = await pointsData.json();
        return [forecastDataJson, pointsDataJson];
      })
      .then(([forecastDataJson, pointsDataJson]) => {
        setForecasts(forecastDataJson);
        setForecastPoints(pointsDataJson);
      }) 
      .catch(error => console.error('Error fetching data: ', error));
  }, [category]);
    
  const getRecentForecastPoint = (forecastId) => {
    const pointsForForecast = forecastPoints.filter(point => point.forecast === forecastId);
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