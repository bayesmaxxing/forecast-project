import React, {useState, useEffect} from 'react';
import { Link } from 'react-router-dom';
import Sidebar from '../components/Sidebar';


function ResolvedForecastPage() {
  const [forecasts, setForecasts] = useState([]);
  const [searchQuery, setsearchQuery] = useState('');

  useEffect(() => {
      fetch(`https://forecasting-389105.ey.r.appspot.com/forecasts?type=resolved`, {
        headers : {
          "Accept" : "application/json"
        }
      })
      .then(response => response.json())
      .then(data => {
        setForecasts(data);
      })
      .catch(error => console.error('Error fetching data: ', error));
  }, []);

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
  })

  const formatDate = (dateString) => dateString.split('T')[0];
  const getResolution = (resolution) => {
    if (resolution === "1") {
      return "Yes"
    } else if (resolution === "0") {
      return "No"
    } else {
      return "Ambiguous"
    }
  };

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
                <p>{getResolution(forecast.resolution)}</p>
              </div>
            </div>
            <div> 
                <p>Resolved on: {formatDate(forecast.resolved)}</p>
                {forecast.brier_score ? (
                <p>Brier score: {forecast.brier_score.toFixed(3)}</p>
              ) : (
                <p>No scores on Ambiguous forecasts.</p>
              )}
            </div>
          </li>
        ))}
      </ul>
    </div>
  );
};
  
  export default ResolvedForecastPage;

  