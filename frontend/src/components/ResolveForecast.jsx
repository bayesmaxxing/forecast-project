import React, { useState } from 'react';
import { useParams } from 'react-router-dom';

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

        const dataToSubmit = {
            ...resolveData,
            resolution: resolveData.resolution,
            comment: resolveData.comment,          
            id: parseInt(id)
        };

        try {
            const response = await fetch(`https://forecasting-389105.ey.r.appspot.com/resolve/${id}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
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
              <label htmlFor="resolution-ambiguous">Resolved as Ambiguous</label>
              <input
                type="radio"
                id="resolution-ambiguous"
                name="resolution"
                value="-"
                checked={resolveData.resolution==="-"}
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
