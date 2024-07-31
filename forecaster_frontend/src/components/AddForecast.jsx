import React, { useState } from 'react';
import './ForecastForms.css';

const AddForecast = () => {
    const [ForecastData, setForecastData] = useState({
        question : '', 
        short_question : '',
        category : '',
        resolution_criteria : ''
    });
    const [submitStatus,setSubmitStatus] = useState('');

    const handleChange = (e) => {
        const {name,value} = e.target;
        setForecastData(prevState => ({
            ...prevState, 
            [name]: value
        }));
    };

    const getDate = () => {
        const currentDate = new Date();
        return currentDate.toISOString().split('T')[0] + ' 00:00:00'
    }

    const handleSubmit = async (e) => {
        e.preventDefault();
        setSubmitStatus('Submitting forecast...');

        const dataToSubmit = {
            ...ForecastData,
            creation_date: getDate()
        };

        try {
            const response = await fetch(`https://forecasting-389105.ey.r.appspot.com/forecaster/api/forecasts/`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Token ${process.env.REACT_APP_API_TOKEN}`
                }, 
                body: JSON.stringify(dataToSubmit)
            });
            
            if (response.ok) {
                setSubmitStatus('Forecast added');
                setForecastData({question: '', short_question:'',
                    category: '',resolution_criteria: ''
                });
            } else {
                setSubmitStatus('Forecast not added. Try again.');    
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
              <label htmlFor="question">Question</label>
              <input
                type="text"
                id="question"
                name="question"
                value={ForecastData.question}
                onChange={handleChange}
                required
              />
            </div>
            <div className="form-group">
              <label htmlFor="short_question">Short Question</label>
              <input
                type="text"
                id="short_question"
                name="short_question"
                value={ForecastData.short_question}
                onChange={handleChange}
                required
              />
            </div>
            <div className="form-group">
              <label htmlFor="category">Category</label>
              <input
                id="category"
                name="category"
                value={ForecastData.category}
                onChange={handleChange}
                required
              ></input>
            </div>
            <div className="form-group">
              <label htmlFor="resolution_criteria">Resolution Criteria</label>
              <textarea
                id="resolution_criteria"
                name="resolution_criteria"
                value={ForecastData.resolution_criteria}
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
    
    export default AddForecast;