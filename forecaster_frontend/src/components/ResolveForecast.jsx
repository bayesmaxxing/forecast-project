import React, { useState } from 'react';

const getScores = (forecastPoints, resolution) => {
        const epsilon = 1e-15
        
        const brierScore = forecastPoints.reduce((sum, currentValue) => sum + (currentValue - resolution) ** 2, 0) / forecastPoints.length;
        
        const logNScore = forecastPoints.reduce((sum, currentValue) => 
          sum + (resolution * Math.log(Math.max(currentValue, epsilon)) 
              + (1-resolution) * Math.log(Math.max(1-currentValue, epsilon))),0) / forecastPoints.length;
        
        const log2Score = forecastPoints.reduce((sum, currentValue) => 
          sum + (resolution * Math.log2(Math.max(currentValue, epsilon)) 
              + (1-resolution) * Math.log2(Math.max(1-currentValue, epsilon))),0) / forecastPoints.length;
        
        return { brierScore, lognScore, log2Score };
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
        }
  };

  const scores = getScores(forecastPoints, resolution)

  const handleSubmit = async (e) => {
        e.preventDefault();
        setSubmitStatus('Submitting update...');

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
            const response = await fetch(`https://forecast-project-backend.vercel.app/forecaster/api/resolutions/?forecast=${id}`, {
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
        <div>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label>Resolution</label>
              <input
                type="text"
                id="resolution"
                name="resolution"
                value={resolveData.resolution}
                onChange={handleChange}
                required
              />
            </div>
            <div>
              <label >Comment on resolution</label>
              <textarea
                id="comment"
                name="comment"
                value={resolveData.comment}
                onChange={handleChange}
                required
                rows="3"
              ></textarea>
            </div>
            <button
              type="submit"
            >
              Submit
            </button>
          </form>
          {submitStatus && <p >{submitStatus}</p>}
        </div>
      );
    };
    export default ResolveForecast;
