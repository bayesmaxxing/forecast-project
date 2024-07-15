import React, { useState } from 'react';

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
        const year = currentDate.getFullYear();
        const month = String(currentDate.getMonth()+1).padStart(2,'0');
        const day = String(currentDate.getDate()).padStart(2,'0');
        return `${year}-${month}-${day} 00:00:00`;
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setSubmitStatus('Submitting forecast...');

        const dataToSubmit = {
            ...ForecastData,
            creation_date: getDate()
        };

        try {
            const response = await fetch(`https://forecast-project-backend.vercel.app/forecaster/api/forecasts/`, {
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
        <div>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label>Question</label>
              <input
                type="text"
                id="question"
                name="question"
                value={ForecastData.question}
                onChange={handleChange}
                required
              />
            </div>
            <div>
              <label>Short Question</label>
              <input
                type="text"
                id="short_question"
                name="short_question"
                value={ForecastData.short_question}
                onChange={handleChange}
                required
              />
            </div>
            <div>
              <label>Category</label>
              <input
                id="category"
                name="category"
                value={ForecastData.category}
                onChange={handleChange}
                required
              ></input>
            </div>
            <div>
              <label >Resolution Criteria</label>
              <textarea
                id="resolution_criteria"
                name="resolution_criteria"
                value={ForecastData.resolution_criteria}
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
    
    export default AddForecast;