import React, { useState } from 'react';
import { useParams } from 'react-router-dom';

const getScores = (forecastPoints, resolution) => {
        const epsilon = 1e-15;
        const resolutionInt = resolution === "1" ? 1 : 0;
        
        const brierScore = forecastPoints.reduce((sum, currentValue) => sum + (currentValue - resolutionInt) ** 2, 0) / forecastPoints.length;
        
        const logNScore = forecastPoints.reduce((sum, currentValue) => 
          sum + (resolutionInt * Math.log(Math.max(currentValue, epsilon)) 
              + (1-resolutionInt) * Math.log(Math.max(1-currentValue, epsilon))),0) / forecastPoints.length;
        
        const log2Score = forecastPoints.reduce((sum, currentValue) => 
          sum + (resolutionInt * Math.log2(Math.max(currentValue, epsilon)) 
              + (1-resolutionInt) * Math.log2(Math.max(1-currentValue, epsilon))),0) / forecastPoints.length;
        
        return { brierScore, logNScore, log2Score };
};

const getDate = () => {
        const currentDate = new Date();
        return currentDate.toISOString().split('T')[0] + ' 00:00:00'
};

const ResolveForecast = ({ forecastPoints }) => {
  let { id } = useParams();
  const [resolveData, setResolveData] = useState({
        resolution: '',
        comment: '',
    });
  
  const [submitStatus,setSubmitStatus] = useState('');

  const handleChange = (e) => {
        const { name, value } = e.target;
            setResolveData(prevState => ({
                ...prevState,
                [name]: value
            }));
        };

  const handleSubmit = async (e) => {
        e.preventDefault();
        setSubmitStatus('Submitting update...');

        const scores = getScores(forecastPoints, resolveData.resolution);

        const dataToSubmit = {
            ...resolveData,
            resolution: resolveData.resolution,
            brier_score: scores.brierScore,
            logn_score: scores.logNScore,
            log2_score: scores.log2Score,
            comment: resolveData.comment,          
            resolution_date: getDate(),
            forecast: parseInt(id)
        };

        try {
            const response = await fetch(`https://forecasting-389105.ey.r.appspot.com/forecaster/api/resolutions/?forecast=${id}`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Token ${process.env.REACT_APP_API_TOKEN}`
                }, 
                body: JSON.stringify(dataToSubmit)
            });
            
            if (response.ok) {
                setSubmitStatus('Resolved');
                setResolveData({resolution: '', comment:'',
                });
            } else {
                setSubmitStatus('Resolution not added. Try again.');    
            } 
        } catch (error) {
            console.error('Error:',error)
            setSubmitStatus('An error occurred. Try again.')
        }
      };
     
    return (
        <div className="form-container">
          <form onSubmit={handleSubmit} className="forecast-form">
            <div className="form-group">
              <label htmlFor="resolution-yes">Resolved as Yes</label>
              <input
                type="radio"
                id="resolution-yes"
                name="resolution"
                value="1"
                checked={resolveData.resolution==="1"}
                onChange={handleChange}
              />
            </div>
            <div className="form-group">
              <label htmlFor="resolution-no">Resolved as No</label>
              <input
                type="radio"
                id="resolution-no"
                name="resolution"
                value="0"
                checked={resolveData.resolution==="0"}
                onChange={handleChange}
              />
            </div>
            <div className="form-group">
              <label htmlFor="comment">Comment on resolution</label>
              <textarea
                id="comment"
                name="comment"
                value={resolveData.comment}
                onChange={handleChange}
                required
                rows="3"
              ></textarea>
            </div>
            <button type="submit" className="submit-button">
              Submit
            </button>
          </form>
          {submitStatus && <p >{submitStatus}</p>}
        </div>
      );
    };
    export default ResolveForecast;
