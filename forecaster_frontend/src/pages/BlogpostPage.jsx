import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { Marked } from 'marked';
import './BlogpostPage.css';


function BlogpostPage() {
    const [blogpost, setBlogpost] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    let { slug } = useParams();
  
    useEffect(() => {
        fetch(`https://forecast-project-backend.vercel.app/forecaster/api/blogposts/?slug=${slug}`, {
          headers : {
            'Authorization': `Token ${process.env.REACT_APP_API_TOKEN}`
          }
        })
      .then(async ([blogpostData]) => {
        if (!blogpostData.ok) {
          throw new Error('Error fetching data');
        }
        const blogpostJson = await blogpostData.json();
        return [blogpostJson];
      })
      .then(([blogpostJson]) => {
        setBlogpost(blogpostJson);
        setLoading(false);
      })
      .catch(error => {
        setError(error);
        setLoading(false);
      });
    }, [slug]);
  
    if (loading) return <div>Loading...</div>;
    if (error) return <div>Error loading the forecast: {error.message}</div>;

    return (
        <div>
          <div>
          <div className='question-header'>{forecastData.question}</div>
          </div>
          <div className='chart-box'>
          <ForecastGraph data = {chartData} options={chartOptions} />
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
        </div>
    );
}

export default SpecificForecast;