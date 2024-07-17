import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import ForecastGraph from '../components/ForecastGraph';
import './SpecificForecast.css';
import UpdateForecast from '../components/UpdateForecast';


function SpecificForecast() {
    const [forecastData, setForecastData] = useState(null);
    const [forecastPoints, setForecastPoints] = useState(null);
    const [resolutionData, setResolution] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [isAdmin, setIsAdmin] = useState(false);
    let { id } = useParams();
  
    useEffect(() => {
      
      Promise.all([
        fetch(`https://forecast-project-backend.vercel.app/forecaster/api/forecasts/${id}/`, {
          headers : {
            'Authorization': `Token ${process.env.REACT_APP_API_TOKEN}`
          }
        }),
        fetch(`https://forecast-project-backend.vercel.app/forecaster/api/forecast_points/?forecast=${id}`, {
          headers : {
            'Authorization': `Token ${process.env.REACT_APP_API_TOKEN}`
          }
        }),
        fetch(`https://forecast-project-backend.vercel.app/forecaster/api/resolutions/?forecast=${id}`, {
          headers : {
            'Authorization': `Token ${process.env.REACT_APP_API_TOKEN}`
          }
        })
      ])
      .then(async ([idData, pointsData, resData]) => {
        if (!idData.ok || !pointsData.ok) {
          throw new Error('Error fetching data');
        }
        const idJson = await idData.json();
        const pointsJson = await pointsData.json();
        const resJson = await resData.json();
        return [idJson, pointsJson, resJson];
      })
      .then(([idJson, pointsJson, resJson]) => {
        setForecastData(idJson);
        setForecastPoints(pointsJson);
        setResolution(resJson[0]);
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
      labels: forecastPoints.map(point => new Date(point.date_added).toLocaleDateString('en-CA')),
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

    return (
        <div>
          <div>
          <div className='question-header'>{forecastData.question}</div>
          <div>
            {resolutionData != null ? (
            <div className='question-header'>Resolved as: {resolutionData.resolution === "1" ? "Yes":"No"}</div>): (null)}
          </div>
          </div>
          <div className='chart-box'>
          <ForecastGraph data = {chartData} options={chartOptions} />
          </div>
          <div>
            {resolutionData != null ? (
            <div className='info-box'>
            <div className='info-header'>Question resolved as: {resolutionData.resolution === "1" ? "Yes" : "No"}</div>
            <div className='info-item'>It resolved on {formatDate(resolutionData.resolution_date)} and it resulted in a 
            Brier score of {resolutionData.brier_score}.</div>
            {resolutionData.comment != null ? (
              <div className='info-item'>Comment: {resolutionData.comment} </div>
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
            {[...forecastPoints].reverse().map(forecast => (
                <li className='update-item'>
                  <div className="update-container">
                    <div className='info-header'>Update to {(forecast.point_forecast * 100).toFixed(1)}% on {formatDate(forecast.date_added)}:</div>
                    <div className='info-item'>{forecast.reason}</div>
                </div>
                
                </li>
              ))}
            </ul>
          </div>
          {isAdmin && <UpdateForecast />}
        </div>
    );
}

export default SpecificForecast;