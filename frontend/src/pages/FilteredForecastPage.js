import React, {useState, useEffect} from 'react';
import { useParams } from 'react-router-dom';
import { Link } from 'react-router-dom';
import './ForecastPage.css';
import Sidebar from '../components/Sidebar';


function ForecastPage() {
    const [searchQuery, setsearchQuery] = useState('');
    const [combinedForecasts, setCombinedForecasts] = useState([]);
    const [scores, setScores] = useState([]);
    const [loading, setLoading] = useState(true);

    // set category based on URL
    let { category } = useParams()

    useEffect(() => {
       // Fetch the list of forecasts from the API
        Promise.all([
          fetch(`https://forecasting-389105.ey.r.appspot.com/forecasts?category=${category}&type=open`, {
            headers : {
              "Accept": "application/json"
            }
          }),
          fetch(`https://forecasting-389105.ey.r.appspot.com/forecast-points/latest`, {
            headers : {
              "Accept": "application/json"
            }
          }), 
          fetch(`https://forecasting-389105.ey.r.appspot.com/scores?category=${category}`, {
            headers : {
              "Accept": "application/json"
            }
          })
        ])
        .then(async ([forecastData, pointsData, scoresData]) => {
          const forecastDataJson = await forecastData.json();
          const pointsDataJson = await pointsData.json();
          const scoresDataJson = await scoresData.json();
          return [forecastDataJson, pointsDataJson, scoresDataJson];
        })
        .then(([forecastDataJson, pointsDataJson, scoresDataJson]) => {
          const combined = forecastDataJson.map(forecast => {
            const matchingPoint = pointsDataJson.find(point => point.forecast_id === forecast.id);
            return { ...forecast, latestPoint: matchingPoint || null};
          });
          setCombinedForecasts(combined)
          setScores(scoresDataJson)
          setLoading(false)
        })
        .catch(error => console.error('Error fetching data: ', error));
    }, [category]);
    
  const handleSearchChange = (e) => {
    setsearchQuery(e.target.value.toLowerCase());
  };

  const filteredForecasts = combinedForecasts.filter(forecast => 
    forecast.question.toLowerCase().includes(searchQuery) ||
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

      {loading ? (
        <p>Brier score: loading...</p>
      ) : (
        scores && scores.AggBrierScore > 0.0 ? (
          <p>Brier score: {(scores.AggBrierScore).toFixed(4)}</p>
        ) : (
          <p>No Brier score available.</p>
        )
      )}

      <ul className="forecast-list">
        {loading ? (
          [...Array(5)].map((_, index) => (
            <li key={index} className="forecast-item">
              <div className="question-container">
                <div className="placeholder-text">Loading forecast...</div>
                <div className="recent-forecast-point">
                  <span>--%</span>
                </div>
              </div>
              <div>
                <p>Category: --</p>
                <p>Created: --</p>
              </div>
            </li>
          ))
        ) : (
           sortedForecasts.map(forecast => (
            <li key={forecast.id} className="forecast-item">
              <div className="question-container">
                <Link to={`/forecast/${forecast.id}`} className="question-link">\
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
          ))
          )}
        </ul>
      </div>
    );
};
  
  export default ForecastPage;
