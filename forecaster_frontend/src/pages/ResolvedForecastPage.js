import React, {useState, useEffect} from 'react';
import { Link } from 'react-router-dom';
import Sidebar from '../components/Sidebar';


function ResolvedForecastPage() {
  const [forecasts, setForecasts] = useState([]);
  const [resolutions, setResolutions] = useState([]);
  useEffect(() => {
    const forecastsCacheKey = 'forecasts_resolved';
    const resolutionsCacheKey = 'resolutions';
    const forecastsCached = localStorage.getItem(forecastsCacheKey);
    const resolutionsCached = localStorage.getItem(resolutionsCacheKey);
  
    if (forecastsCached && resolutionsCached) {
      setForecasts(JSON.parse(forecastsCached));
      setResolutions(JSON.parse(resolutionsCached));
    } else {
      Promise.all([
        fetch(`https://forecast-project-backend.vercel.app/forecaster/api/forecasts/?resolved=True`),
        fetch(`https://forecast-project-backend.vercel.app/forecaster/api/resolutions/`) 
      ])
      .then(async ([forecastData, resData]) => {
        const forecastDataJson = await forecastData.json();
        const resDataJson = await resData.json();
        return [forecastDataJson, resDataJson];
      })
      .then(([forecastDataJson, resDataJson]) => {
        setForecasts(forecastDataJson);
        setResolutions(resDataJson);
        localStorage.setItem(forecastsCacheKey, JSON.stringify(forecastDataJson));
        localStorage.setItem(resolutionsCacheKey, JSON.stringify(resDataJson));
      })
      .catch(error => console.error('Error fetching data: ', error));
    }
  }, []);

  const getResolution = (forecastId) => {
    const resolutionForForecast = resolutions.find(resolution => resolution.forecast === forecastId);
  
    if (!resolutionForForecast) {
      // Handle the case when there is no resolution for the forecast
      return { Res: 'N/A', brier: 'N/A', log2: 'N/A', logn: 'N/A' };
    }
  
    const Res = resolutionForForecast.resolution === '0' ? 'NO' : 'YES';
    const brier = resolutionForForecast.brier_score;
    const log2 = resolutionForForecast.log2_score;
    const logn = resolutionForForecast.logn_score;
  
    return { Res, brier, log2, logn };
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
                <p>{getResolution(forecast.id).Res}</p>
              </div>
            </div>
            <div> 
                <p>Brier score: {getResolution(forecast.id).brier.toFixed(3)}</p>
                <p>Natural log score: {getResolution(forecast.id).logn.toFixed(3)}</p>
                <p>Log_2 score: {getResolution(forecast.id).log2.toFixed(3)}</p>
            </div>
          </li>
        ))}
      </ul>
    </div>
  );
};
  
  export default ResolvedForecastPage;

  