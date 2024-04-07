import React, {useState, useEffect} from 'react';
import { Link } from 'react-router-dom';
import Sidebar from '../components/Sidebar';


function ResolvedForecastPage() {
  const [forecasts, setForecasts] = useState([]);
  const [resolutions, setResolutions] = useState([]);
  const [searchQuery, setsearchQuery] = useState('');

  useEffect(() => {
    const CACHE_DURATION = 5 * 60 * 1000; // Cache duration in milliseconds, e.g., 5 minutes
    const now = new Date().getTime(); // Current time

    const forecastsCacheKey = 'forecasts_resolved';
    const resolutionsCacheKey = 'resolutions';

    const forecastsCached = localStorage.getItem(forecastsCacheKey);
    const resolutionsCached = localStorage.getItem(resolutionsCacheKey);

    const forecastsDataValid = forecastsCached && now - JSON.parse(forecastsCached).timestamp < CACHE_DURATION;
    const resolutionsDataValid = resolutionsCached && now - JSON.parse(resolutionsCached).timestamp < CACHE_DURATION;
  
    if (forecastsDataValid && resolutionsDataValid) {
      setForecasts(JSON.parse(forecastsCached).data);
      setResolutions(JSON.parse(resolutionsCached).data);
    } else {
      Promise.all([
        fetch(`https://forecast-project-backend.vercel.app/forecaster/api/forecasts/?resolved=True`, {
          headers : {
            'Authorization': `Token ${process.env.REACT_APP_API_TOKEN}`
          }
        }),
        fetch(`https://forecast-project-backend.vercel.app/forecaster/api/resolutions/`, {
          headers : {
            'Authorization': `Token ${process.env.REACT_APP_API_TOKEN}`
          }
        }) 
      ])
      .then(async ([forecastData, resData]) => {
        const forecastDataJson = await forecastData.json();
        const resDataJson = await resData.json();
        return [forecastDataJson, resDataJson];
      })
      .then(([forecastDataJson, resDataJson]) => {
        setForecasts(forecastDataJson);
        setResolutions(resDataJson);
        localStorage.setItem('forecasts_resolved', JSON.stringify({data: forecastDataJson, timestamp: now}));
        localStorage.setItem('resolutions', JSON.stringify({data: resDataJson, timestamp: now}));
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

  