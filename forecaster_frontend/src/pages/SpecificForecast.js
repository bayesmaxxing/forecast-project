import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import ForecastGraph from '../components/ForecastGraph';
import './SpecificForecast.css';
import UpdateForecast from '../components/UpdateForecast';
import ResolveForecast from '../components/ResolveForecast';


function SpecificForecast() {
    const [forecastData, setForecastData] = useState(null);
    const [forecastPoints, setForecastPoints] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [isAdmin, setIsAdmin] = useState(false);
    let { id } = useParams();
  
    useEffect(() => {
      
      Promise.all([
        fetch(`https://forecasting-389105.ey.r.appspot.com/forecasts/${id}`, {
          headers : {
            "Accept": "application/json"
          }
        }),
        fetch(`https://forecasting-389105.ey.r.appspot.com/forecast-points/${id}`, {
          headers : {
            "Accept": "application/json"
          }
        })
      ])
      .then(async ([idData, pointsData]) => {
        if (!idData.ok || !pointsData.ok) {
          throw new Error('Error fetching data');
        }
        const idJson = await idData.json();
        const pointsJson = await pointsData.json();
        return [idJson, pointsJson];
      })
      .then(([idJson, pointsJson]) => {
        setForecastData(idJson);
        const sortedPoints = pointsJson.sort((a, b) => new Date(a.created) - new Date(b.created));
        setForecastPoints(sortedPoints);
        setLoading(false);
      })
      .catch(error => {
        setError(error);
        setLoading(false);
      });
      const checkAdminStatus = () => {
        const expirationTime = localStorage.getItem('adminLoginExpiration');
        setIsAdmin(expirationTime && new Date().getTime() < parseInt(expirationTime, 10));
      };

      checkAdminStatus();
    }, [id]);
    

    if (loading) return <div>Loading...</div>;
    if (error) return <div>Error loading the forecast: {error.message}</div>;
    
    const chartData = {
      labels: forecastPoints.map(point => new Date(point.created).toLocaleDateString('en-CA')),
      datasets: [
          {
              label: 'Prediction',
              data: forecastPoints.map(point => point.point_forecast),
              fill: false,
              borderColor: 'rgb(75, 192, 192)',
              tension: 0.1
          }
      ]
    };
    const chartOptions = {
      scales: {
        y: {
          min: 0, 
          max: 1, 
        }
      },

    };
    
    const formatDate = (dateString) => dateString.split('T')[0];
    const reversedForecastpoints = [...forecastPoints].reverse();
    const resolution = forecastData.resolution === "1" ? "Yes":
                       forecastData.resolution === "0" ? "No":
                       "Ambiguous"
                
    return (
        <div>
          <div>
          <div className='question-header'>{forecastData.question}</div>
          <div>
            {forecastData.resolved != null ? (
            <div className='question-header'>Resolved as: {resolution}</div>
            ) : null}
          </div>
          </div>
          <div className='chart-box'>
          <ForecastGraph data = {chartData} options={chartOptions} />
          </div>
          {isAdmin && forecastData.resolved == null && <ResolveForecast forecastPoints={forecastPoints} />}
          <div>
            {forecastData.resolved != null ? (
            <div className='info-box'>
            <div className='info-header'>Question resolved as: {resolution}</div>
            <div className='info-item'>It resolved on {formatDate(forecastData.resolved)} and it resulted in a 
            Brier score of {!forecastData.brier_score ? 0: forecastData.brier_score }.</div>
            {forecastData.comment != null ? (
              <div className='info-item'>Comment: {forecastData.comment} </div>
            ) : (null)
            }
            </div>): (null)}
          </div> 
          <div className='info-box'>
            <div className='info-header'>Resolution Criteria</div>
            <div className='info-item'>{forecastData.resolution_criteria}</div>
          </div>
          <div className='updates-box'>
            <ul className='update-list'>
            {reversedForecastpoints.map(forecast => (
                <li key={forecast.forecast_id} className='update-item'>
                  <div className="update-container">
                    <div className='info-header'>Update to {(forecast.point_forecast * 100).toFixed(1)}% on {formatDate(forecast.created)}:</div>
                    <div className='info-item'>{forecast.reason}</div>
                </div>
                
                </li>
              ))}
            </ul>
          </div>
          {isAdmin && forecastData.resolved == null && <UpdateForecast />}
        </div>
    );
}

export default SpecificForecast;
